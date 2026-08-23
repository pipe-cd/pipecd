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
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
)

func TestBuildPipelineStages(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		stages        []sdk.StageConfig
		rollback      bool
		wantNames     []string
		wantIndexes   []int
		wantRollbacks []bool
	}{
		{
			name: "no rollback",
			stages: []sdk.StageConfig{
				{Name: StageECSPrimaryRollout, Index: 0},
				{Name: StageECSCanaryRollout, Index: 1},
				{Name: StageECSTrafficRouting, Index: 2},
			},
			rollback:      false,
			wantNames:     []string{StageECSPrimaryRollout, StageECSCanaryRollout, StageECSTrafficRouting},
			wantIndexes:   []int{0, 1, 2},
			wantRollbacks: []bool{false, false, false},
		},
		{
			name: "with rollback",
			stages: []sdk.StageConfig{
				{Name: StageECSPrimaryRollout, Index: 0},
				{Name: StageECSCanaryRollout, Index: 1},
			},
			rollback: true,
			// The rollback stage reuses the smallest requested index so it
			// stays valid under piped's stage index validation and runs
			// first among all plugins' rollback stages.
			wantNames:     []string{StageECSPrimaryRollout, StageECSCanaryRollout, StageECSRollback},
			wantIndexes:   []int{0, 1, 0},
			wantRollbacks: []bool{false, false, true},
		},
		{
			name: "with rollback and non-contiguous indexes",
			stages: []sdk.StageConfig{
				{Name: StageECSPrimaryRollout, Index: 2},
				{Name: StageECSTrafficRouting, Index: 5},
			},
			rollback:      true,
			wantNames:     []string{StageECSPrimaryRollout, StageECSTrafficRouting, StageECSRollback},
			wantIndexes:   []int{2, 5, 2},
			wantRollbacks: []bool{false, false, true},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			input := &sdk.BuildPipelineSyncStagesInput{
				Request: sdk.BuildPipelineSyncStagesRequest{
					Stages:   tc.stages,
					Rollback: tc.rollback,
				},
			}
			got := buildPipelineStages(input)

			require.Len(t, got, len(tc.wantNames))
			for i, s := range got {
				assert.Equal(t, tc.wantNames[i], s.Name)
				assert.Equal(t, tc.wantIndexes[i], s.Index)
				assert.Equal(t, tc.wantRollbacks[i], s.Rollback)
			}
		})
	}
}

// replica of piped's controller.validateStageIndexes (unexported there), kept in sync manually.
func validateStageIndexes(req []sdk.StageConfig, res []sdk.PipelineStage) error {
	reqIndexes := make(map[int]struct{})
	for _, s := range req {
		reqIndexes[s.Index] = struct{}{}
	}
	for _, s := range res {
		if _, ok := reqIndexes[s.Index]; !ok {
			return fmt.Errorf("stage index %d from plugin is not defined in the request", s.Index)
		}
	}
	return nil
}

// TestBuildPipelineStagesRollbackIndexContract verifies the SDK contract that
// every returned Index must be one of the requested indexes
// (see sdk.PipelineStage.Index comment).
func TestBuildPipelineStagesRollbackIndexContract(t *testing.T) {
	t.Parallel()

	reqStages := []sdk.StageConfig{
		{Name: StageECSPrimaryRollout, Index: 0},
		{Name: StageECSCanaryRollout, Index: 1},
	}
	got := buildPipelineStages(&sdk.BuildPipelineSyncStagesInput{
		Request: sdk.BuildPipelineSyncStagesRequest{
			Stages:   reqStages,
			Rollback: true,
		},
	})

	err := validateStageIndexes(reqStages, got)
	assert.NoError(t, err, "rollback stage index must be one of the requested indexes")
}

func TestBuildQuickSyncPipeline(t *testing.T) {
	t.Parallel()

	t.Run("without rollback", func(t *testing.T) {
		t.Parallel()
		got := buildQuickSyncPipeline(false)
		require.Len(t, got, 1)
		assert.Equal(t, StageECSSync, got[0].Name)
		assert.False(t, got[0].Rollback)
	})

	t.Run("with rollback", func(t *testing.T) {
		t.Parallel()
		got := buildQuickSyncPipeline(true)
		require.Len(t, got, 2)
		assert.Equal(t, StageECSSync, got[0].Name)
		assert.Equal(t, StageECSRollback, got[1].Name)
		assert.True(t, got[1].Rollback)
	})
}
