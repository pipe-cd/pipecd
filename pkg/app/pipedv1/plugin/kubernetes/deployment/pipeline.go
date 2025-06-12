// Copyright 2024 The PipeCD Authors.
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

const (
	// StageK8sSync represents the state where
	// all resources should be synced with the Git state.
	StageK8sSync = "K8S_SYNC"
	// StageK8sPrimaryRollout represents the state where
	// the PRIMARY variant resources has been updated to the new version/configuration.
	StageK8sPrimaryRollout = "K8S_PRIMARY_ROLLOUT"
	// StageK8sCanaryRollout represents the state where
	// the CANARY variant resources has been rolled out with the new version/configuration.
	StageK8sCanaryRollout = "K8S_CANARY_ROLLOUT"
	// StageK8sCanaryClean represents the state where
	// the CANARY variant resources has been cleaned.
	StageK8sCanaryClean = "K8S_CANARY_CLEAN"
	// StageK8sBaselineRollout represents the state where
	// the BASELINE variant resources has been rolled out.
	StageK8sBaselineRollout = "K8S_BASELINE_ROLLOUT"
	// StageK8sBaselineClean represents the state where
	// the BASELINE variant resources has been cleaned.
	StageK8sBaselineClean = "K8S_BASELINE_CLEAN"
	// StageK8sTrafficRouting represents the state where the traffic to application
	// should be splitted as the specified percentage to PRIMARY, CANARY, BASELINE variants.
	StageK8sTrafficRouting = "K8S_TRAFFIC_ROUTING"
	// StageK8sRollback represents the state where all deployed resources should be rollbacked.
	StageK8sRollback = "K8S_ROLLBACK"
)

var allStages = []string{
	StageK8sSync,
	StageK8sPrimaryRollout,
	StageK8sCanaryRollout,
	StageK8sCanaryClean,
	StageK8sBaselineRollout,
	StageK8sBaselineClean,
	StageK8sTrafficRouting,
	StageK8sRollback,
}

const (
	// StageDescriptionK8sSync represents the description of the K8sSync stage.
	StageDescriptionK8sSync = "Sync by applying all manifests"
	// StageDescriptionK8sRollback represents the description of the K8sRollback stage.
	StageDescriptionK8sRollback = "Rollback the deployment"
)
