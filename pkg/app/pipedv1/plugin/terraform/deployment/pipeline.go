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

import (
	"time"

	"github.com/pipe-cd/pipecd/pkg/model"
)

type stage string

const (
	// stageTerraformSync syncs infrastructure with all the terraform files defined in Git.
	// At first, it executes "plan", and if any changes detected, it executes "apply" for those changes automatically.
	stageTerraformSync stage = "TERRAFORM_SYNC"
	// stageTerraformPlan executes "plan" and shows the plan result.
	stageTerraformPlan stage = "TERRAFORM_PLAN"
	// stageTerraformApply executes "apply" to sync infrastructure.
	stageTerraformApply stage = "TERRAFORM_APPLY"

	// stageTerraformRollback rollbacks the deployment.
	stageTerraformRollback stage = "TERRAFORM_ROLLBACK"
)

var allStages = []string{
	string(stageTerraformSync),
	string(stageTerraformPlan),
	string(stageTerraformApply),
	string(stageTerraformRollback),
}

var (
	predefinedStageTerraformSync = model.PipelineStage{
		Id:       "TerraformSync",
		Name:     string(stageTerraformSync),
		Desc:     "Sync by automatically applying all detected changes",
		Rollback: false,
	}
	predefinedStageTerraformRollback = model.PipelineStage{
		Id:       "TerraformRollback",
		Name:     string(stageTerraformRollback),
		Desc:     "Rollback the deployment",
		Rollback: true,
	}
)

func buildQuickSyncStages(autoRollback bool, now time.Time) []*model.PipelineStage {
	out := make([]*model.PipelineStage, 0, 2)

	out = append(out, &model.PipelineStage{
		Id:        predefinedStageTerraformSync.GetId(),
		Name:      predefinedStageTerraformSync.GetName(),
		Desc:      predefinedStageTerraformSync.GetDesc(),
		Rollback:  predefinedStageTerraformSync.GetRollback(),
		Status:    model.StageStatus_STAGE_NOT_STARTED_YET,
		Metadata:  nil,
		CreatedAt: now.Unix(),
		UpdatedAt: now.Unix(),
	})

	// Append ROLLBACK stage if auto rollback is enabled.
	if autoRollback {
		out = append(out, &model.PipelineStage{
			Id:        predefinedStageTerraformRollback.GetId(),
			Name:      predefinedStageTerraformRollback.GetName(),
			Desc:      predefinedStageTerraformRollback.GetDesc(),
			Rollback:  predefinedStageTerraformRollback.GetRollback(),
			Status:    model.StageStatus_STAGE_NOT_STARTED_YET,
			CreatedAt: now.Unix(),
			UpdatedAt: now.Unix(),
		})
	}

	return out
}
