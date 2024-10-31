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

package executor

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pipe-cd/pipecd/pkg/model"
)

func TestDetermineStageStatus(t *testing.T) {
	testcases := []struct {
		name     string
		sig      StopSignalType
		ori      model.StageStatus
		got      model.StageStatus
		expected model.StageStatus
	}{
		{
			name:     "No stop signal, should get got status",
			sig:      StopSignalNone,
			ori:      model.StageStatus_STAGE_RUNNING,
			got:      model.StageStatus_STAGE_SUCCESS,
			expected: model.StageStatus_STAGE_SUCCESS,
		}, {
			name:     "Terminated signal given, should get original status",
			sig:      StopSignalTerminate,
			ori:      model.StageStatus_STAGE_RUNNING,
			got:      model.StageStatus_STAGE_SKIPPED,
			expected: model.StageStatus_STAGE_RUNNING,
		}, {
			name:     "Timeout signal given, should get failed status",
			sig:      StopSignalTimeout,
			ori:      model.StageStatus_STAGE_RUNNING,
			got:      model.StageStatus_STAGE_RUNNING,
			expected: model.StageStatus_STAGE_FAILURE,
		}, {
			name:     "Cancel signal given, should get cancelled status",
			sig:      StopSignalCancel,
			ori:      model.StageStatus_STAGE_RUNNING,
			got:      model.StageStatus_STAGE_RUNNING,
			expected: model.StageStatus_STAGE_CANCELLED,
		}, {
			name:     "Unknown signal type given, should get failed status",
			sig:      StopSignalType("unknown"),
			ori:      model.StageStatus_STAGE_RUNNING,
			got:      model.StageStatus_STAGE_RUNNING,
			expected: model.StageStatus_STAGE_FAILURE,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := DetermineStageStatus(tc.sig, tc.ori, tc.got)
			assert.Equal(t, tc.expected, got)
		})
	}
}
