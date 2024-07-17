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
	"fmt"
	"strings"
	"time"

	"github.com/pipe-cd/pipecd/pkg/config"
	"github.com/pipe-cd/pipecd/pkg/model"
)

const (
	PredefinedStageK8sSync            = "K8sSync"
	PredefinedStageRollback           = "Rollback"
)
	
var predefinedStages = map[string]config.PipelineStage{
	PredefinedStageK8sSync: {
		ID:   PredefinedStageK8sSync,
		Name: model.StageK8sSync,
		Desc: "Sync by applying all manifests",
	},
}

// GetPredefinedStage finds and returns the predefined stage for the given id.
func GetPredefinedStage(id string) (config.PipelineStage, bool) {
	stage, ok := predefinedStages[id]
	return stage, ok
}

// MakeInitialStageMetadata makes the initial metadata for the given state configuration.
func MakeInitialStageMetadata(cfg config.PipelineStage) map[string]string {
	switch cfg.Name {
	case model.StageWaitApproval:
		return map[string]string{
			"Approvers": strings.Join(cfg.WaitApprovalStageOptions.Approvers, ","),
		}
	default:
		return nil
	}
}


func buildQuickSyncPipeline(autoRollback bool, now time.Time) []*model.PipelineStage {
	var (
		preStageID = ""
		stage, _   = GetPredefinedStage(PredefinedStageK8sSync)
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
			Metadata:   MakeInitialStageMetadata(s),
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
		s, _ := GetPredefinedStage(PredefinedStageRollback)
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
