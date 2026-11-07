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
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	sdk "github.com/pipe-cd/piped-plugin-sdk-go"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/ecs/provider"
)

// computeSyncState determines whether the live ECS state matches what was declared in Git.
//
// Rather than diffing configuration fields directly, it relies on the commit hash that PipeCD
// stamps during deployment. Where the hash is stored depends on the deployment controller of the service.
//
// This keeps the check resilient to AWS-side mutations (e.g. auto-scaling adjustments) that are intentional
// and should not be treated as drift.
func computeSyncState(result *queryResourcesResult, desiredService types.Service, commitHash string) sdk.ApplicationSyncState {
	// A nil service means AWS returned no matching service for the cluster in the desired configuration.
	// The app has never been deployed, or was deleted out-of-band
	if result.Service == nil {
		return sdk.ApplicationSyncState{
			Status:      sdk.ApplicationSyncStateOutOfSync,
			ShortReason: fmt.Sprintf("Service %s not found in cluster", aws.ToString(desiredService.ServiceName)),
		}
	}

	switch deploymentControllerType(result.Service) {
	case types.DeploymentControllerTypeEcs:
		return computeECSControllerSyncState(result.Service, commitHash)
	default:
		return computeExternalControllerSyncState(result, commitHash)
	}
}

// computeExternalControllerSyncState compares the commit hash tag of the PRIMARY task set,
// which PipeCD manages itself with the EXTERNAL deployment controller.
func computeExternalControllerSyncState(result *queryResourcesResult, commitHash string) sdk.ApplicationSyncState {
	// Under the EXTERNAL deployment controller model PipeCD uses,
	// exactly one task set carries the PRIMARY status at any point in time.
	// Its absence means no successful deployment has ever completed, or the service is in a broken intermediate state.
	var primaryTaskSet *types.TaskSet
	for i := range result.TaskSets {
		if aws.ToString(result.TaskSets[i].Status) == "PRIMARY" {
			primaryTaskSet = &result.TaskSets[i]
			break
		}
	}
	if primaryTaskSet == nil {
		return sdk.ApplicationSyncState{
			Status:      sdk.ApplicationSyncStateOutOfSync,
			ShortReason: "No PRIMARY task set found",
		}
	}

	// PipeCD tags every task set it creates with the Git commit hash at the time of deployment.
	// If the PRIMARY task set's hash differs from the current HEAD,
	// the cluster is running an older revision and needs a new deployment to reconcile.
	//
	// When commitHash is empty (e.g. the deployment source has no VCS info),
	// we skip this check rather than producing a false OUT_OF_SYNC.
	if commitHash != "" {
		for _, tag := range primaryTaskSet.Tags {
			if aws.ToString(tag.Key) == provider.LabelCommitHash {
				liveCommit := aws.ToString(tag.Value)
				if liveCommit != commitHash {
					return sdk.ApplicationSyncState{
						Status:      sdk.ApplicationSyncStateOutOfSync,
						ShortReason: "Deployed commit does not match current commit",
						Reason:      fmt.Sprintf("deployed: %s, expected: %s", liveCommit, commitHash),
					}
				}
				break
			}
		}
	}

	return sdk.ApplicationSyncState{Status: sdk.ApplicationSyncStateSynced}
}

// deploymentControllerType returns the deployment controller of the live service.
// A service without a deployment controller is managed by the ECS controller.
func deploymentControllerType(svc *types.Service) types.DeploymentControllerType {
	if svc.DeploymentController == nil || svc.DeploymentController.Type == "" {
		return types.DeploymentControllerTypeEcs
	}
	return svc.DeploymentController.Type
}

// computeECSControllerSyncState compares the commit hash tag of a service managed by the ECS controller.
// PipeCD writes this tag on the service only after a service deployment succeeds, so it is
// the commit that is actually running, even when a newer rollout is in progress or has failed.
func computeECSControllerSyncState(svc *types.Service, commitHash string) sdk.ApplicationSyncState {
	// Without a commit to compare against, the check is skipped to avoid a false OUT_OF_SYNC.
	if commitHash == "" {
		return sdk.ApplicationSyncState{Status: sdk.ApplicationSyncStateSynced}
	}

	for _, tag := range svc.Tags {
		if aws.ToString(tag.Key) != provider.LabelCommitHash {
			continue
		}
		if liveCommit := aws.ToString(tag.Value); liveCommit != commitHash {
			return sdk.ApplicationSyncState{
				Status:      sdk.ApplicationSyncStateOutOfSync,
				ShortReason: "Deployed commit does not match current commit",
				Reason:      fmt.Sprintf("deployed: %s, expected: %s", liveCommit, commitHash),
			}
		}
		return sdk.ApplicationSyncState{Status: sdk.ApplicationSyncStateSynced}
	}

	// The tag is written only after a successful deployment, so its absence means none has succeeded yet.
	return sdk.ApplicationSyncState{
		Status:      sdk.ApplicationSyncStateOutOfSync,
		ShortReason: "No successful deployment found for the service",
	}
}
