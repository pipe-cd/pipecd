// Copyright 2020 The PipeCD Authors.
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

package model

// Stage represents the middle and temporary state of application
// before reaching its final desired state.
type Stage string

const (
	// StageWait represents the waiting state for a specified period of time.
	StageWait Stage = "WAIT"
	// StageWaitApproval represents the waiting state until getting an approval
	// from one of the specified approvers.
	StageWaitApproval Stage = "WAIT_APPROVAL"
	// StageAnalysis represents the waiting state for analysing
	// the application status based on metrics, log, http request...
	StageAnalysis Stage = "ANALYSIS"

	// StageK8sPrimaryUpdate represents the state where
	// the PRIMARY has been updated to the new version/configuration.
	StageK8sPrimaryUpdate Stage = "K8S_PRIMARY_UPDATE"
	// StageK8sStageRollout represents the state where
	// the STAGE workloads has been rolled out with the new version/configuration.
	StageK8sStageRollout Stage = "K8S_STAGE_ROLLOUT"
	// StageK8sStageClean represents the state where
	// the STAGE workloads has been cleaned.
	StageK8sStageClean Stage = "K8S_STAGE_CLEAN"
	// StageK8sBaselineRollout represents the state where
	// the BASELINE workloads has been rolled out with the new version/configuration.
	StageK8sBaselineRollout Stage = "K8S_BASELINE_ROLLOUT"
	// StageK8sBaselineClean represents the state where
	// the BASELINE workloads has been cleaned.
	StageK8sBaselineClean Stage = "K8S_BASELINE_CLEAN"
	// StageK8sTrafficRoute represents the state where the traffic to application
	// should be routed as the specified percentage to PRIMARY, STAGE, BASELINE.
	StageK8sTrafficRoute Stage = "K8S_TRAFFIC_ROUTE"

	// StageTerraformPlan shows terraform plan result.
	StageTerraformPlan Stage = "TERRAFORM_PLAN"
	// StageTerraformApply represents the state where
	// the new configuration has been applied.
	StageTerraformApply Stage = "TERRAFORM_APPLY"

	// StageCloudRunNewVersionRollout represents the state where
	// the workloads of the new version has been rolled out.
	StageCloudRunNewVersionRollout Stage = "CLOUDRUN_NEW_VERSION_ROLLOUT"
	// StageCloudRunTrafficRoute represents the state where the traffic to application
	// should be routed as the specified percentage to previous version and new version.
	StageCloudRunTrafficRoute Stage = "CLOUDRUN_TRAFFIC_ROUTE"

	// StageLambdaNewVersionRollout represents the state where
	// the workloads of the new version has been rolled out.
	StageLambdaNewVersionRollout Stage = "LAMBDA_NEW_VERSION_ROLLOUT"
	// StageLambdaTrafficRoute represents the state where the traffic to application
	// should be routed as the specified percentage to previous version and new version.
	StageLambdaTrafficRoute Stage = "LAMBDA_TRAFFIC_ROUTE"

	// StageRollBack represents a state where
	// the all temporarily created stages will be reverted to
	// bring back the pre-deploy stage.
	// This stage is AUTOMATICALLY GENERATED and can not be used
	// to specify in configuration file.
	StageRollBack Stage = "ROLLBACK"
)

func (s Stage) String() string {
	return string(s)
}
