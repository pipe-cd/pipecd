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

package stagelogstore

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/pipe-cd/pipecd/pkg/cache"
	"github.com/pipe-cd/pipecd/pkg/cache/cachetest"
	"github.com/pipe-cd/pipecd/pkg/model"
)

func TestCacheGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	c := cachetest.NewMockCache(ctrl)

	testcases := []struct {
		name         string
		deploymentID string
		stageID      string
		retriedCount int32

		returnData string
		returnErr  error

		expected    logFragment
		expectedErr error
	}{
		{
			name:       "cache key not found",
			returnData: "",
			returnErr:  cache.ErrNotFound,

			expected:    logFragment{},
			expectedErr: cache.ErrNotFound,
		},
		{
			name: "successfully getting from cache",
			returnData: `
{
    "Blocks": [
        {
            "index": 1,
            "log": "Hello 1",
            "severity": 1,
            "created_at": 1590499431
        },
        {
            "index": 2,
            "log": "Hello 2",
            "severity": 2,
            "created_at": 1590499432
        }
    ],
    "Completed": false
}`,
			returnErr: nil,

			expected: logFragment{
				Blocks: []*model.LogBlock{
					{
						Index:     1,
						Log:       "Hello 1",
						Severity:  model.LogSeverity_SUCCESS,
						CreatedAt: 1590499431,
					},
					{
						Index:     2,
						Log:       "Hello 2",
						Severity:  model.LogSeverity_ERROR,
						CreatedAt: 1590499432,
					},
				},
				Completed: false,
			},
			expectedErr: nil,
		},
	}

	slc := stageLogCache{
		cache: c,
	}

	for _, tc := range testcases {
		key := cacheKey(tc.deploymentID, tc.stageID, tc.retriedCount)
		c.EXPECT().Get(key).Return([]byte(tc.returnData), tc.returnErr)

		lf, err := slc.Get(tc.deploymentID, tc.stageID, tc.retriedCount)
		if err != nil {
			if tc.expectedErr == nil {
				assert.NoError(t, err)
				return
			}
			assert.Error(t, err, tc.expectedErr)
			return
		}
		assert.Equal(t, tc.expected, lf)
	}
}
