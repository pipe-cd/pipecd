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

package planner

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/planner"
	"github.com/pipe-cd/pipecd/pkg/config"
	"github.com/pipe-cd/pipecd/pkg/model"
)

func buildQuickSyncPipeline(autoRollback bool, now time.Time) []*model.PipelineStage {
	var (
		preStageID = ""
		stage, _   = planner.GetPredefinedStage(planner.PredefinedStageK8sSync)
		stages     = []config.PipelineStage{stage}
		out        = make([]*model.PipelineStage, 0, len(stages))
	)

	for i, s := range stages {
		id := s.ID
		if id == "" {
			id = fmt.Sprintf("stage-%d", i)
		}
		stage := &model.PipelineStage{
			Id:         id,
			Name:       s.Name.String(),
			Desc:       s.Desc,
			Index:      int32(i),
			Predefined: true,
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

		// Add a stage for rolling back script run stages.
		for i, s := range pp.Stages {
			if s.Name == model.StageScriptRun {
				// Use metadata as a way to pass parameters to the stage.
				envStr, _ := json.Marshal(s.ScriptRunStageOptions.Env)
				metadata := map[string]string{
					"baseStageID": out[i].Id,
					"onRollback":  s.ScriptRunStageOptions.OnRollback,
					"env":         string(envStr),
				}
				ss, _ := planner.GetPredefinedStage(planner.PredefinedStageScriptRunRollback)
				out = append(out, &model.PipelineStage{
					Id:         ss.ID,
					Name:       ss.Name.String(),
					Desc:       ss.Desc,
					Predefined: true,
					Visible:    false,
					Status:     model.StageStatus_STAGE_NOT_STARTED_YET,
					Metadata:   metadata,
					CreatedAt:  now.Unix(),
					UpdatedAt:  now.Unix(),
				})
			}
		}
	}

	return out
}
