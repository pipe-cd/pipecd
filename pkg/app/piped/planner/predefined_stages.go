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

package planner

import (
	"github.com/pipe-cd/pipecd/pkg/config"
	"github.com/pipe-cd/pipecd/pkg/model"
)

const (
	PredefinedStageK8sSync       = "K8sSync"
	PredefinedStageTerraformSync = "TerraformSync"
	PredefinedStageCloudRunSync  = "CloudRunSync"
	PredefinedStageLambdaSync    = "LambdaSync"
	PredefinedStageECSSync       = "ECSSync"
	PredefinedStageRollback      = "Rollback"
)

var predefinedStages = map[string]config.PipelineStage{
	PredefinedStageK8sSync: {
		Id:   PredefinedStageK8sSync,
		Name: model.StageK8sSync,
		Desc: "Sync by applying all manifests",
	},
	PredefinedStageTerraformSync: {
		Id:   PredefinedStageTerraformSync,
		Name: model.StageTerraformSync,
		Desc: "Sync by automatically applying any detected changes",
	},
	PredefinedStageCloudRunSync: {
		Id:   PredefinedStageCloudRunSync,
		Name: model.StageCloudRunSync,
		Desc: "Deploy the new version and configure all traffic to it",
	},
	PredefinedStageLambdaSync: {
		Id:   PredefinedStageLambdaSync,
		Name: model.StageLambdaSync,
		Desc: "Deploy the new version and configure all traffic to it",
	},
	PredefinedStageECSSync: {
		Id:   PredefinedStageECSSync,
		Name: model.StageECSSync,
		Desc: "Deploy the new version and configure all traffic to it",
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
