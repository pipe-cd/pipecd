// Copyright 2026 The PipeCD Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package livestate

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
	"go.uber.org/zap"

	ecsconfig "github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/ecs/config"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/ecs/provider"
)

const (
	resourceTypeService = "ECS:Service"
	resourceTypeTaskSet = "ECS:TaskSet"
	resourceTypeTask    = "ECS:Task"
)

var (
	_ sdk.LivestatePlugin[ecsconfig.ECSPluginConfig, ecsconfig.ECSDeployTargetConfig, ecsconfig.ECSApplicationSpec] = (*ECSLivestatePlugin)(nil)
	_ sdk.Initializer[ecsconfig.ECSPluginConfig, ecsconfig.ECSDeployTargetConfig]                                   = (*ECSLivestatePlugin)(nil)
)

type ECSLivestatePlugin struct {
	fetcher     *ECSFetcher
	initialized sync.Once
}

func (p *ECSLivestatePlugin) Initialize(
	ctx context.Context,
	input *sdk.InitializeInput[ecsconfig.ECSPluginConfig, ecsconfig.ECSDeployTargetConfig],
) error {
	var err error
	p.initialized.Do(func() {
		if len(input.DeployTargets) != 1 {
			err = fmt.Errorf("only 1 deploy target is allowed but got %d", len(input.DeployTargets))
			return
		}
		var (
			dtName   string
			dtConfig ecsconfig.ECSDeployTargetConfig
			client   provider.Client
		)
		for name, cfg := range input.DeployTargets {
			dtName = name
			dtConfig = cfg.Config
		}

		client, err = provider.DefaultRegistry().Client(dtName, dtConfig)
		if err != nil {
			return
		}

		p.fetcher = &ECSFetcher{
			client: client,
		}
	})
	return err
}

// GetLivestate returns the current live state of the ECS application and whether it is in sync with the desired state declared in Git
func (p *ECSLivestatePlugin) GetLivestate(
	ctx context.Context,
	_ *ecsconfig.ECSPluginConfig,
	deployTargets []*sdk.DeployTarget[ecsconfig.ECSDeployTargetConfig],
	input *sdk.GetLivestateInput[ecsconfig.ECSApplicationSpec],
) (*sdk.GetLivestateResponse, error) {
	appCfg, err := input.Request.DeploymentSource.AppConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get app config: %w", err)
	}
	spec := appCfg.Spec.Input

	input.Logger.Info("GetLivestate called",
		zap.String("applicationID", input.Request.ApplicationID),
		zap.String("serviceDefinitionFile", spec.ServiceDefinitionFile),
		zap.Bool("runStandaloneTask", spec.RunStandaloneTask),
		zap.String("appDir", input.Request.DeploymentSource.ApplicationDirectory),
	)

	// Standalone tasks are one-off runs with no persistent service to observe.
	// There is no meaningful live state to report, return UNKNOWN rather than misleading the user with an empty or stale result
	if spec.RunStandaloneTask {
		return &sdk.GetLivestateResponse{
			SyncState: sdk.ApplicationSyncState{
				Status:      sdk.ApplicationSyncStateUnknown,
				ShortReason: "Standalone task does not have a live state",
			},
		}, nil
	}

	// A missing service definition file is a misconfiguration, without it plugin cannot identify which ECS service to inspect
	if spec.ServiceDefinitionFile == "" {
		return &sdk.GetLivestateResponse{
			SyncState: sdk.ApplicationSyncState{
				Status:      sdk.ApplicationSyncStateInvalidConfig,
				ShortReason: "serviceDefinitionFile is required",
			},
		}, nil
	}

	// Parse the service definition from Git to obtain the cluster and service name
	// Intentionally, do not inject deployment tags here as only need the identifiers to locate the live resource on AWS
	appDir := input.Request.DeploymentSource.ApplicationDirectory
	desiredService, err := provider.ParseServiceDefinition(appDir, spec.ServiceDefinitionFile)
	if err != nil {
		return nil, fmt.Errorf("failed to parse service definition: %w", err)
	}

	// Fetch the live state from AWS, an error here (for example network issue, permission denied, ...)
	// is surfaced as OUT_OF_SYNC rather than a hard error
	// so the UI can still display something meaningful instead of going blank
	result, err := p.fetcher.FetchResources(ctx, desiredService)
	if err != nil {
		return &sdk.GetLivestateResponse{
			SyncState: sdk.ApplicationSyncState{
				Status:      sdk.ApplicationSyncStateOutOfSync,
				ShortReason: fmt.Sprintf("Failed to query ECS resources: %v", err),
			},
		}, nil
	}

	deployTargetName := deployTargets[0].Name

	// Build the response by converting raw AWS objects into SDK resource states and
	// computing the sync verdict against the current Git commit
	return &sdk.GetLivestateResponse{
		LiveState: sdk.ApplicationLiveState{
			Resources: buildResourceStates(result, deployTargetName),
		},
		SyncState: computeSyncState(result, desiredService, input.Request.DeploymentSource.CommitHash),
	}, nil
}

func buildResourceStates(result *queryResourcesResult, deployTarget string) []sdk.ResourceState {
	if result.Service == nil {
		return nil
	}

	svc := result.Service
	resources := make([]sdk.ResourceState, 0, 1+len(result.TaskSets)+len(result.Tasks))

	// Service
	svcHealth, svcDesc := serviceHealthStatus(svc)
	resources = append(resources, sdk.ResourceState{
		ID:           aws.ToString(svc.ServiceArn),
		Name:         aws.ToString(svc.ServiceName),
		ResourceType: resourceTypeService,
		ResourceMetadata: map[string]string{
			"status":       aws.ToString(svc.Status),
			"runningCount": strconv.Itoa(int(svc.RunningCount)),
			"desiredCount": strconv.Itoa(int(svc.DesiredCount)),
			"pendingCount": strconv.Itoa(int(svc.PendingCount)),
		},
		HealthStatus:      svcHealth,
		HealthDescription: svcDesc,
		DeployTarget:      deployTarget,
		CreatedAt:         derefTime(svc.CreatedAt),
	})

	// TaskSets, build a map startedBy(=taskSet.Id) -> taskSetArn for task parent lookup
	taskSetArnByID := make(map[string]string, len(result.TaskSets))
	for _, ts := range result.TaskSets {
		tsHealth, tsDesc := taskSetHealthStatus(&ts)
		meta := map[string]string{
			"status":         aws.ToString(ts.Status),
			"taskDefinition": aws.ToString(ts.TaskDefinition),
		}
		if ts.Scale != nil {
			meta["scale"] = fmt.Sprintf("%.0f%%", ts.Scale.Value)
		}
		for _, tag := range ts.Tags {
			if aws.ToString(tag.Key) == provider.LabelCommitHash {
				meta["commit"] = aws.ToString(tag.Value)
				break
			}
		}

		taskSetArnByID[aws.ToString(ts.Id)] = aws.ToString(ts.TaskSetArn)
		resources = append(resources, sdk.ResourceState{
			ID:                aws.ToString(ts.TaskSetArn),
			ParentIDs:         []string{aws.ToString(svc.ServiceArn)},
			Name:              aws.ToString(ts.Id),
			ResourceType:      resourceTypeTaskSet,
			ResourceMetadata:  meta,
			HealthStatus:      tsHealth,
			HealthDescription: tsDesc,
			DeployTarget:      deployTarget,
			CreatedAt:         derefTime(ts.CreatedAt),
		})
	}

	// Tasks
	for _, task := range result.Tasks {
		taskHealth, taskDesc := taskHealthStatus(&task)
		meta := map[string]string{
			"lastStatus": aws.ToString(task.LastStatus),
		}
		if task.StartedAt != nil {
			meta["startedAt"] = task.StartedAt.Format(time.RFC3339)
		}

		// Link task to its TaskSet via StartedBy field
		parentIDs := []string{aws.ToString(svc.ServiceArn)}
		if tsArn, ok := taskSetArnByID[aws.ToString(task.StartedBy)]; ok {
			parentIDs = []string{tsArn}
		}

		resources = append(resources, sdk.ResourceState{
			ID:                aws.ToString(task.TaskArn),
			ParentIDs:         parentIDs,
			Name:              aws.ToString(task.TaskArn),
			ResourceType:      resourceTypeTask,
			ResourceMetadata:  meta,
			HealthStatus:      taskHealth,
			HealthDescription: taskDesc,
			DeployTarget:      deployTarget,
			CreatedAt:         derefTime(task.CreatedAt),
		})
	}

	return resources
}

// derefTime safely dereferences a *time.Time returned by the AWS SDK.
//
// AWS represents optional timestamps as pointers, so a nil value indicates the field was not populated.
func derefTime(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	}
	return *t
}

func serviceHealthStatus(svc *types.Service) (sdk.ResourceHealthStatus, string) {
	switch aws.ToString(svc.Status) {
	case "ACTIVE":
		if svc.PendingCount == 0 && svc.RunningCount >= svc.DesiredCount {
			return sdk.ResourceHealthStateHealthy, ""
		}
		return sdk.ResourceHealthStateUnhealthy, fmt.Sprintf(
			"Service has %d running tasks out of %d desired (%d pending)",
			svc.RunningCount, svc.DesiredCount, svc.PendingCount,
		)
	case "DRAINING", "INACTIVE":
		return sdk.ResourceHealthStateUnhealthy, fmt.Sprintf("Service is in %s state", aws.ToString(svc.Status))
	default:
		return sdk.ResourceHealthStateUnknown, ""
	}
}

func taskSetHealthStatus(ts *types.TaskSet) (sdk.ResourceHealthStatus, string) {
	switch ts.StabilityStatus {
	case types.StabilityStatusSteadyState:
		return sdk.ResourceHealthStateHealthy, ""
	case types.StabilityStatusStabilizing:
		return sdk.ResourceHealthStateUnhealthy, fmt.Sprintf(
			"TaskSet is stabilizing: %d pending, %d running",
			ts.PendingCount, ts.RunningCount,
		)
	default:
		return sdk.ResourceHealthStateUnknown, ""
	}
}

func taskHealthStatus(task *types.Task) (sdk.ResourceHealthStatus, string) {
	if task.LastStatus == nil {
		return sdk.ResourceHealthStateUnknown, ""
	}
	switch aws.ToString(task.LastStatus) {
	case "RUNNING":
		if task.HealthStatus == types.HealthStatusUnhealthy {
			return sdk.ResourceHealthStateUnhealthy, "Task container health checks are failing"
		}
		return sdk.ResourceHealthStateHealthy, ""
	case "PENDING":
		return sdk.ResourceHealthStateUnhealthy, "Task is in PENDING state"
	case "STOPPED":
		reason := "Task stopped"
		if task.StoppedReason != nil {
			reason = fmt.Sprintf("Task stopped: %s", aws.ToString(task.StoppedReason))
		}
		return sdk.ResourceHealthStateUnhealthy, reason
	default:
		return sdk.ResourceHealthStateUnknown, ""
	}
}
