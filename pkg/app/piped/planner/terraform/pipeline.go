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

package terraform

import (
	"fmt"
	"time"

	"github.com/pipe-cd/pipecd/pkg/app/piped/planner"
	"github.com/pipe-cd/pipecd/pkg/config"
	"github.com/pipe-cd/pipecd/pkg/model"
)

func buildQuickSyncPipeline(autoRollback bool, now time.Time) []*model.PipelineStage {
	var (
		s, _ = planner.GetPredefinedStage(planner.PredefinedStageTerraformSync)
		out  = make([]*model.PipelineStage, 0, 2)
	)

	// Append SYNC stage.
	id := s.ID
	if id == "" {
		id = "stage-0"
	}
	stage := &model.PipelineStage{
		Id:         id,
		Name:       s.Name.String(),
		Desc:       s.Desc,
		Index:      0,
		Predefined: true,
		Visible:    true,
		Status:     model.StageStatus_STAGE_NOT_STARTED_YET,
		Metadata:   planner.MakeInitialStageMetadata(s),
		CreatedAt:  now.Unix(),
		UpdatedAt:  now.Unix(),
	}
	out = append(out, stage)

	// Append ROLLBACK stage if auto rollback is enabled.
	if autoRollback {
		s, _ := planner.GetPredefinedStage(planner.PredefinedStageRollback)
		out = append(out, &model.PipelineStage{
			Id:         s.ID,
			Name:       s.Name.String(),
			Desc:       s.Desc,
			Predefined: true,
			Visible:    false,
			Status:     model.StageStatus_STAGE_NOT_STARTED_YET,
			CreatedAt:  now.Unix(),
			UpdatedAt:  now.Unix(),
		})
	}

	return out
}

func buildProgressivePipeline(pp *config.DeploymentPipeline, autoRollback bool, now time.Time) []*model.PipelineStage {
	var (
		preStageID = ""
		out        = make([]*model.PipelineStage, 0, len(pp.Stages))
	)

	for i, s := range pp.Stages {
		id := s.ID
		if id == "" {
			id = fmt.Sprintf("stage-%d", i)
		}
		stage := &model.PipelineStage{
			Id:         id,
			Name:       s.Name.String(),
			Desc:       s.Desc,
			Index:      int32(i),
			Predefined: false,
			Visible:    true,
			Status:     model.StageStatus_STAGE_NOT_STARTED_YET,
			Metadata:   planner.MakeInitialStageMetadata(s),
			CreatedAt:  now.Unix(),
			UpdatedAt:  now.Unix(),
		}
		if preStageID != "" {
			stage.Requires = []string{preStageID}
		}
		preStageID = id
		out = append(out, stage)
	}

	if autoRollback {
		s, _ := planner.GetPredefinedStage(planner.PredefinedStageRollback)
		out = append(out, &model.PipelineStage{
			Id:         s.ID,
			Name:       s.Name.String(),
			Desc:       s.Desc,
			Predefined: true,
			Visible:    false,
			Status:     model.StageStatus_STAGE_NOT_STARTED_YET,
			CreatedAt:  now.Unix(),
			UpdatedAt:  now.Unix(),
		})
	}

	return out
}
