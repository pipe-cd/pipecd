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

package deployment

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	sdk "github.com/pipe-cd/piped-plugin-sdk-go"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/ecs/config"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/ecs/provider"
)

const (
	canaryTargetGroupArnKey        = "canary-target-group-arn"
	trafficRoutePrimaryMetadataKey = "primary-percentage"
	trafficRouteCanaryMetadataKey  = "canary-percentage"
	currentListenersKey            = "current-listeners"
)

func (p *ECSPlugin) executeECSTrafficRouting(
	ctx context.Context,
	input *sdk.ExecuteStageInput[config.ECSApplicationSpec],
	deployTarget *sdk.DeployTarget[config.ECSDeployTargetConfig],
) sdk.StageStatus {
	lp := input.Client.LogPersister()

	cfg, err := input.Request.TargetDeploymentSource.AppConfig()
	if err != nil {
		lp.Errorf("Failed to load app config: %v", err)
		return sdk.StageStatusFailure
	}

	accessType := cfg.Spec.Input.AccessType
	if accessType != "ELB" {
		lp.Errorf("Unsupported access type %s in stage Traffic Routing for ECS application", accessType)
		return sdk.StageStatusFailure
	}

	primary, canary, err := provider.LoadTargetGroups(cfg.Spec.Input.TargetGroups)
	if err != nil {
		lp.Errorf("Failed to load target groups: %v", err)
		return sdk.StageStatusFailure
	}

	if primary == nil || canary == nil {
		lp.Errorf("Required both primary and canary target groups for traffic routing")
		return sdk.StageStatusFailure
	}

	if err = input.Client.PutDeploymentPluginMetadata(
		ctx,
		canaryTargetGroupArnKey,
		*canary.TargetGroupArn,
	); err != nil {
		lp.Errorf("Failed to save canary target group ARN for rollback: %v", err)
		return sdk.StageStatusFailure
	}

	client, err := provider.DefaultRegistry().Client(deployTarget.Name, deployTarget.Config)
	if err != nil {
		lp.Errorf("Failed to get ECS client for deploy target %s: %v", deployTarget.Name, err)
		return sdk.StageStatusFailure
	}

	var options config.ECSTrafficRoutingStageOptions
	if err := json.Unmarshal(input.Request.StageConfig, &options); err != nil {
		lp.Errorf("Failed to parse stage option: %v", err)
		return sdk.StageStatusFailure
	}

	lp.Infof("Start performing routing traffic")
	if err = routing(ctx, lp, input.Client, client, *primary, *canary, options); err != nil {
		lp.Errorf("Failed to route traffic: %v", err)
		return sdk.StageStatusFailure
	}

	return sdk.StageStatusSuccess
}

// metadataStore abstracts the deployment plugin metadata operations for testability
type metadataStore interface {
	PutDeploymentPluginMetadataMulti(ctx context.Context, metadata map[string]string) error
	GetDeploymentPluginMetadata(ctx context.Context, key string) (string, bool, error)
	PutDeploymentPluginMetadata(ctx context.Context, key string, value string) error
}

func routing(
	ctx context.Context,
	lp sdk.StageLogPersister,
	mdStore metadataStore,
	providerClient provider.Client,
	primaryTargetGroup types.LoadBalancer,
	canaryTargetGroup types.LoadBalancer,
	options config.ECSTrafficRoutingStageOptions,
) error {
	// Retrieve traffic split of primary and canary
	primaryWeight, canaryWeight := options.Percentages()
	routeTrafficCfg := provider.RoutingTrafficConfig{
		{
			TargetGroupArn: *primaryTargetGroup.TargetGroupArn,
			Weight:         primaryWeight,
		},
		{
			TargetGroupArn: *canaryTargetGroup.TargetGroupArn,
			Weight:         canaryWeight,
		},
	}

	percentageMetadata := map[string]string{
		trafficRoutePrimaryMetadataKey: strconv.FormatInt(int64(primaryWeight), 10),
		trafficRouteCanaryMetadataKey:  strconv.FormatInt(int64(canaryWeight), 10),
	}
	if err := mdStore.PutDeploymentPluginMetadataMulti(ctx, percentageMetadata); err != nil {
		return fmt.Errorf("failed to store percentage metadata: %v", err)
	}

	var currListenerArns []string
	value, ok, err := mdStore.GetDeploymentPluginMetadata(ctx, currentListenersKey)
	if err != nil {
		return fmt.Errorf("failed to get current listener arns: %v", err)
	}
	if ok {
		currListenerArns = strings.Split(value, ",")
	} else {
		currListenerArns, err = providerClient.GetListenerArns(ctx, primaryTargetGroup)
		if err != nil {
			return fmt.Errorf("failed to get current active listeners: %v", err)
		}
	}

	metadataCurrListener := strings.Join(currListenerArns, ",")
	if err := mdStore.PutDeploymentPluginMetadata(ctx, currentListenersKey, metadataCurrListener); err != nil {
		return fmt.Errorf("failed to store listeners to metadata store: %v", err)
	}

	modifiedRules, err := providerClient.ModifyListeners(ctx, currListenerArns, routeTrafficCfg)
	if err != nil {
		return fmt.Errorf("failed to routing traffic to primary and canary variants: %v", err)
	}
	lp.Infof("Modified %d listener rules: %v", len(modifiedRules), modifiedRules)

	return nil
}
