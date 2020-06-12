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

package terraform

import (
	"fmt"
	"time"

	"github.com/pipe-cd/pipe/pkg/app/piped/planner"
	"github.com/pipe-cd/pipe/pkg/config"
	"github.com/pipe-cd/pipe/pkg/model"
)

func builDefaultPipeline(now time.Time) []*model.PipelineStage {
	var (
		planStage, _  = planner.GetPredefinedStage(planner.PredefinedStageTerraformPlan)
		applyStage, _ = planner.GetPredefinedStage(planner.PredefinedStageTerraformApply)
		stages        = []config.PipelineStage{
			planStage,
			applyStage,
		}
	)

	out := make([]*model.PipelineStage, 0, len(stages))
	for i, s := range stages {
		id := s.Id
		if id == "" {
			id = fmt.Sprintf("stage-%d", i)
		}
		stage := &model.PipelineStage{
			Id:         id,
			Name:       s.Name.String(),
			Desc:       s.Desc,
			Index:      int32(i),
			Predefined: true,
			Status:     model.StageStatus_STAGE_NOT_STARTED_YET,
			CreatedAt:  now.Unix(),
			UpdatedAt:  now.Unix(),
		}
		out = append(out, stage)
	}

	return out
}

func buildProgressivePipeline(stages []config.PipelineStage, now time.Time) []*model.PipelineStage {
	out := make([]*model.PipelineStage, 0, len(stages))
	for i, s := range stages {
		id := s.Id
		if id == "" {
			id = fmt.Sprintf("stage-%d", i)
		}
		stage := &model.PipelineStage{
			Id:         id,
			Name:       s.Name.String(),
			Desc:       s.Desc,
			Index:      int32(i),
			Predefined: false,
			Status:     model.StageStatus_STAGE_NOT_STARTED_YET,
			CreatedAt:  now.Unix(),
			UpdatedAt:  now.Unix(),
		}
		out = append(out, stage)
	}

	return out
}
