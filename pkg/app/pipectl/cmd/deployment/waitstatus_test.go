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

package deployment

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pipe-cd/pipe/pkg/model"
)

func TestMakeStatuses(t *testing.T) {
	testcases := []struct {
		name        string
		statuses    []string
		expected    []model.DeploymentStatus
		expectedErr bool
	}{
		{
			name:     "empty",
			expected: []model.DeploymentStatus{},
		},
		{
			name:        "has an invalid status",
			statuses:    []string{"SUCCESS", "INVALID"},
			expectedErr: true,
		},
		{
			name:     "ok",
			statuses: []string{"SUCCESS", "PLANNED"},
			expected: []model.DeploymentStatus{
				model.DeploymentStatus_DEPLOYMENT_SUCCESS,
				model.DeploymentStatus_DEPLOYMENT_PLANNED,
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			statuses, err := makeStatuses(tc.statuses)
			assert.Equal(t, tc.expected, statuses)
			assert.Equal(t, tc.expectedErr, err != nil)
		})
	}
}

func TestAvailableStatuses(t *testing.T) {
	statuses := availableStatuses()
	assert.True(t, len(statuses) > 0)
}
