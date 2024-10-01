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
	"testing"
	"time"

	"github.com/pipe-cd/pipecd/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestBuildQuickSyncPipeline(t *testing.T) {
	t.Parallel()

	now := time.Now()

	tests := []struct {
		name         string
		index        int32
		autoRollback bool
		expected     []*model.PipelineStage
	}{
		{
			name:         "without auto rollback",
			index:        0,
			autoRollback: false,
			expected: []*model.PipelineStage{
				{
					Id:         PredefinedStageK8sSync,
					Name:       StageK8sSync.String(),
					Desc:       "Sync by applying all manifests",
					Index:      0,
					Predefined: true,
					Visible:    true,
					Status:     model.StageStatus_STAGE_NOT_STARTED_YET,
					Metadata:   nil,
					CreatedAt:  now.Unix(),
					UpdatedAt:  now.Unix(),
				},
			},
		},
		{
			name:         "with auto rollback",
			index:        0,
			autoRollback: true,
			expected: []*model.PipelineStage{
				{
					Id:         PredefinedStageK8sSync,
					Name:       StageK8sSync.String(),
					Desc:       "Sync by applying all manifests",
					Index:      0,
					Predefined: true,
					Visible:    true,
					Status:     model.StageStatus_STAGE_NOT_STARTED_YET,
					Metadata:   nil,
					CreatedAt:  now.Unix(),
					UpdatedAt:  now.Unix(),
				},
				{
					Id:         PredefinedStageRollback,
					Name:       StageK8sRollback.String(),
					Desc:       "Rollback the deployment",
					Predefined: true,
					Visible:    false,
					Status:     model.StageStatus_STAGE_NOT_STARTED_YET,
					CreatedAt:  now.Unix(),
					UpdatedAt:  now.Unix(),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := buildQuickSyncPipeline(tt.index, tt.autoRollback, now)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
