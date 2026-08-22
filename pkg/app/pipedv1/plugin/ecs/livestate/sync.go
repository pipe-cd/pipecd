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
// Rather than diffing configuration fields directly,
// it relies on the commit hash that PipeCD stamps onto the PRIMARY task set during deployment.
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
