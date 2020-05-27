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
	"github.com/kapetaniosci/pipe/pkg/config"
	"github.com/kapetaniosci/pipe/pkg/model"
)

const (
	PredefinedWaitApproval        = "WaitApproval"
	PredefinedStageK8sScale       = "K8sScale"
	PredefinedStageK8sRollback    = "K8sRollback"
	PredefinedStageK8sUpdate      = "K8sUpdate"
	PredefinedStageTerraformPlan  = "TerraformPlan"
	PredefinedStageTerraformApply = "TerraformApply"
)

var predefinedStages = map[string]config.PipelineStage{
	PredefinedWaitApproval: config.PipelineStage{
		Id:   PredefinedWaitApproval,
		Name: model.StageWaitApproval,
		Desc: "Wait for an approval",
	},
	PredefinedStageK8sScale: config.PipelineStage{
		Id:   PredefinedStageK8sScale,
		Name: model.StageK8sPrimaryUpdate,
		Desc: "Scale primary workloads",
	},
	PredefinedStageK8sRollback: config.PipelineStage{
		Id:   PredefinedStageK8sRollback,
		Name: model.StageK8sPrimaryUpdate,
		Desc: "Rollback primary to previous version",
	},
	PredefinedStageK8sUpdate: config.PipelineStage{
		Id:   PredefinedStageK8sUpdate,
		Name: model.StageK8sPrimaryUpdate,
		Desc: "Update primary to new version/configuration",
	},
	PredefinedStageTerraformPlan: config.PipelineStage{
		Id:   PredefinedStageTerraformPlan,
		Name: model.StageTerraformPlan,
		Desc: "Terraform plan",
	},
	PredefinedStageTerraformApply: config.PipelineStage{
		Id:   PredefinedStageTerraformApply,
		Name: model.StageTerraformApply,
		Desc: "Terraform apply",
	},
}

// GetPredefinedStage finds and returns the predefined stage for the given id.
func GetPredefinedStage(id string) (config.PipelineStage, bool) {
	stage, ok := predefinedStages[id]
	return stage, ok
}
