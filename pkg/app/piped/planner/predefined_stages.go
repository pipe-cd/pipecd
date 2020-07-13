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

package planner

import (
	"github.com/pipe-cd/pipe/pkg/config"
	"github.com/pipe-cd/pipe/pkg/model"
)

const (
	PredefinedWaitApproval        = "WaitApproval"
	PredefinedStageK8sSync        = "K8sSync"
	PredefinedStageTerraformPlan  = "TerraformPlan"
	PredefinedStageTerraformApply = "TerraformApply"
	PredefinedStageRollback       = "Rollback"
)

var predefinedStages = map[string]config.PipelineStage{
	PredefinedWaitApproval: {
		Id:   PredefinedWaitApproval,
		Name: model.StageWaitApproval,
		Desc: "Wait for an approval",
	},
	PredefinedStageK8sSync: {
		Id:   PredefinedStageK8sSync,
		Name: model.StageK8sPrimaryUpdate,
		Desc: "Sync resources with Git state",
	},
	PredefinedStageTerraformPlan: {
		Id:   PredefinedStageTerraformPlan,
		Name: model.StageTerraformPlan,
		Desc: "Terraform plan",
	},
	PredefinedStageTerraformApply: {
		Id:   PredefinedStageTerraformApply,
		Name: model.StageTerraformApply,
		Desc: "Terraform apply",
	},
	PredefinedStageRollback: {
		Id:   PredefinedStageRollback,
		Name: model.StageRollback,
		Desc: "Rollback the deployment",
	},
}

// GetPredefinedStage finds and returns the predefined stage for the given id.
func GetPredefinedStage(id string) (config.PipelineStage, bool) {
	stage, ok := predefinedStages[id]
	return stage, ok
}
