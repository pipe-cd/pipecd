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

	// StageK8sSync represents the state where
	// all resources should be synced with the Git state.
	StageK8sSync Stage = "K8S_SYNC"
	// StageK8sPrimaryRollout represents the state where
	// the PRIMARY variant resources has been updated to the new version/configuration.
	StageK8sPrimaryRollout Stage = "K8S_PRIMARY_ROLLOUT"
	// StageK8sCanaryRollout represents the state where
	// the CANARY variant resources has been rolled out with the new version/configuration.
	StageK8sCanaryRollout Stage = "K8S_CANARY_ROLLOUT"
	// StageK8sCanaryClean represents the state where
	// the CANARY variant resources has been cleaned.
	StageK8sCanaryClean Stage = "K8S_CANARY_CLEAN"
	// StageK8sBaselineRollout represents the state where
	// the BASELINE variant resources has been rolled out.
	StageK8sBaselineRollout Stage = "K8S_BASELINE_ROLLOUT"
	// StageK8sBaselineClean represents the state where
	// the BASELINE variant resources has been cleaned.
	StageK8sBaselineClean Stage = "K8S_BASELINE_CLEAN"
	// StageK8sTrafficRouting represents the state where the traffic to application
	// should be splitted as the specified percentage to PRIMARY, CANARY, BASELINE variants.
	StageK8sTrafficRouting Stage = "K8S_TRAFFIC_ROUTING"

	// StageTerraformSync synced infrastructure with all the tf defined in Git.
	// Firstly, it does plan and if there are any changes detected it applies those changes automatically.
	StageTerraformSync Stage = "TERRAFORM_SYNC"
	// StageTerraformPlan shows terraform plan result.
	StageTerraformPlan Stage = "TERRAFORM_PLAN"
	// StageTerraformApply represents the state where
	// the new configuration has been applied.
	StageTerraformApply Stage = "TERRAFORM_APPLY"

	// StageCloudRunSync does quick sync by rolling out the new version
	// and switching all trafic to it.
	StageCloudRunSync Stage = "CLOUDRUN_SYNC"
	// StageCloudRunCanaryRollout represents the state where
	// the workloads of the new version has been rolled out.
	StageCloudRunCanaryRollout Stage = "CLOUDRUN_CANARY_ROLLOUT"
	// StageCloudRunTrafficRouting represents the state where the traffic to application
	// should be splitted as the specified percentage to previous version and new version.
	StageCloudRunTrafficRouting Stage = "CLOUDRUN_TRAFFIC_ROUTING"

	// StageLambdaSync does quick sync by rolling out the new version
	// and switching all trafic to it.
	StageLambdaSync Stage = "LAMBDA_SYNC"
	// StageLambdaCanaryRollout represents the state where
	// the workloads of the new version has been rolled out.
	StageLambdaCanaryRollout Stage = "LAMBDA_CANARY_ROLLOUT"
	// StageLambdaTrafficRouting represents the state where the traffic to application
	// should be splitted as the specified percentage to previous version and new version.
	StageLambdaTrafficRouting Stage = "LAMBDA_TRAFFIC_ROUTING"

	// StageRollback represents a state where
	// the all temporarily created stages will be reverted to
	// bring back the pre-deploy stage.
	// This stage is AUTOMATICALLY GENERATED and can not be used
	// to specify in configuration file.
	StageRollback Stage = "ROLLBACK"
)

func (s Stage) String() string {
	return string(s)
}
