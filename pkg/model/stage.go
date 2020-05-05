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
	// StageK8sPrimaryOut represents the state where the PRIMARY
	// has been updated to the new version/configuration.
	StageK8sPrimaryOut Stage = "K8S_PRIMARY_OUT"
	// StageK8sStageOut represents the state where the STAGE workloads
	// has been rolled out with the new version/configuration.
	StageK8sStageOut Stage = "K8S_STAGE_OUT"
	// StageK8sStageIn represents the state where the STAGE workloads
	// has been cleaned.
	StageK8sStageIn Stage = "K8S_STAGE_IN"
	// StageK8sBaselineOut represents the state where the BASELINE workloads
	// has been rolled out with the new version/configuration.
	StageK8sBaselineOut Stage = "K8S_BASELINE_OUT"
	// StageK8sBaselineIn represents the state where the BASELINE workloads
	// has been cleaned.
	StageK8sBaselineIn Stage = "K8S_BASELINE_IN"
	// StageK8sTrafficRoute represents the state where the traffic to application
	// should be routed as the specified percentage to PRIMARY, STAGE, BASELINE.
	StageK8sTrafficRoute Stage = "K8S_TRAFFIC_ROUTE"
	// StageTerraformPlan shows terraform plan result.
	StageTerraformPlan Stage = "TERRAFORM_PLAN"
	// StageTerraformApply represents the state where
	// the new configuration has been applied.
	StageTerraformApply Stage = "TERRAFORM_APPLY"
)

func (s Stage) String() string {
	return string(s)
}
