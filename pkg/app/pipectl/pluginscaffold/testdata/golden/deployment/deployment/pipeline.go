// Copyright 2026 The PipeCD Authors.
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
	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
)

const (
	StageDemoSync = "DEMO_SYNC"
	StageDemoRollback = "DEMO_ROLLBACK"
)

var allStages = []string{
	StageDemoSync,
	StageDemoRollback,
}

func buildPipelineStages(input *sdk.BuildPipelineSyncStagesInput) []sdk.PipelineStage {
	stages := input.Request.Stages
	out := make([]sdk.PipelineStage, 0, len(stages)+1)
	for _, rs := range stages {
		out = append(out, sdk.PipelineStage{
			Index:              rs.Index,
			Name:               rs.Name,
			Rollback:           false,
			Metadata:           map[string]string{},
			AvailableOperation: sdk.ManualOperationNone,
		})
	}
	if input.Request.Rollback {
		out = append(out, sdk.PipelineStage{
			Name:               StageDemoRollback,
			Index:              len(stages) + 1,
			Rollback:           true,
			Metadata:           map[string]string{},
			AvailableOperation: sdk.ManualOperationNone,
		})
	}
	return out
}

func buildQuickSyncPipeline(autoRollback bool) []sdk.QuickSyncStage {
	out := make([]sdk.QuickSyncStage, 0, len(allStages))
	for _, name := range allStages {
		if name == StageDemoRollback {
			continue
		}
		out = append(out, sdk.QuickSyncStage{
			Name:        name,
			Description: name,
			Rollback:    false,
		})
	}
	if autoRollback {
		out = append(out, sdk.QuickSyncStage{
			Name:        StageDemoRollback,
			Description: "rollback",
			Rollback:    true,
		})
	}
	return out
}
