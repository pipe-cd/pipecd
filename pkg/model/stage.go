// Copyright 2023 The PipeCD Authors.
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
	// StageScriptRun represents a state where
	// the specified script will be executed.
	StageScriptRun Stage = "SCRIPT_RUN"

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
	// and switching all traffic to it.
	StageCloudRunSync Stage = "CLOUDRUN_SYNC"
	// StageCloudRunPromote promotes the new version to receive amount of traffic.
	StageCloudRunPromote Stage = "CLOUDRUN_PROMOTE"

	// StageLambdaSync does quick sync by rolling out the new version
	// and switching all traffic to it.
	StageLambdaSync Stage = "LAMBDA_SYNC"
	// StageLambdaCanaryRollout represents the state where
	// the CANARY variant resources has been rolled out with the new version/configuration.
	StageLambdaCanaryRollout Stage = "LAMBDA_CANARY_ROLLOUT"
	// StageLambdaPromote prmotes the new version to receive amount of traffic.
	StageLambdaPromote Stage = "LAMBDA_PROMOTE"

	// StageECSSync does quick sync by rolling out the new version
	// and switching all traffic to it.
	StageECSSync Stage = "ECS_SYNC"
	// StageECSCanaryRollout represents the stage where
	// the CANARY variant resource have been rolled out with the new version/configuration.
	// The CANARY variant will serve amount of traffic set in this stage option.
	StageECSCanaryRollout Stage = "ECS_CANARY_ROLLOUT"
	// StageECSPrimaryRollout represents the stage where
	// the PRIMARY variant resource have been rolled out with the new version/configuration.
	// The PRIMARY variant will serve 100% traffic after it's rolled out.
	StageECSPrimaryRollout Stage = "ECS_PRIMARY_ROLLOUT"
	// StageECSTrafficRouting represents the state where the traffic to application
	// should be splitted as the specified percentage to PRIMARY/CANARY variants.
	StageECSTrafficRouting Stage = "ECS_TRAFFIC_ROUTING"
	// StageECSCanaryClean represents the stage where
	// the CANARY variant resources has been cleaned.
	StageECSCanaryClean Stage = "ECS_CANARY_CLEAN"
	// StageCustomSync represents the stage where users can use their
	// defined scripts to sync the application's state instead of the KIND_SYNC stage.
	StageCustomSync Stage = "CUSTOM_SYNC"

	// StageRollback represents a state where
	// the all temporarily created stages will be reverted to
	// bring back the pre-deploy stage.
	// This stage is AUTOMATICALLY GENERATED and can not be used
	// to specify in configuration file.
	StageRollback Stage = "ROLLBACK"
	// StageCustomSyncRollback represents a state where
	// all changes made by the CUSTOM_SYNC stage will be reverted to
	// bring back the pre-deploy stage.
	StageCustomSyncRollback Stage = "CUSTOM_SYNC_ROLLBACK"
	// StageScriptRunRollback represents a state where
	// all changes made by the SCRIPT_RUN_ROLLBACK stage will be reverted to
	// bring back the pre-deploy stage.
	StageScriptRunRollback Stage = "SCRIPT_RUN_ROLLBACK"
)

func (s Stage) String() string {
	return string(s)
}
