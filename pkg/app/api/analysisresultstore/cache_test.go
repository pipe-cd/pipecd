// Copyright 2021 The PipeCD Authors.
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

package analysisresultstore

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pipe-cd/pipe/pkg/cache"
	"github.com/pipe-cd/pipe/pkg/cache/cachetest"
	"github.com/pipe-cd/pipe/pkg/model"
)

func TestCacheGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	c := cachetest.NewMockCache(ctrl)

	testcases := []struct {
		name string

		applicationID string

		returnData string
		returnErr  error

		expected    *model.AnalysisMetadata
		expectedErr error
	}{
		{
			name:          "cache key not found",
			applicationID: "application-id",
			returnData:    "",
			returnErr:     cache.ErrNotFound,

			expected:    nil,
			expectedErr: cache.ErrNotFound,
		},
		{
			name:          "successfully getting from cache",
			applicationID: "application-id",
			returnData: `{
				"startTime": 1590000000,
				"duration": 3600,
				"interval": 300,
				"query": "foo"
			}`,

			expected: &model.AnalysisMetadata{
				StartTime: 1590000000,
				Duration:  3600,
				Interval:  300,
				Query:     "foo",
			},
			expectedErr: nil,
		},
	}

	acache := analysisCache{
		backend: c,
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			key := cacheKey(tc.applicationID)
			c.EXPECT().Get(key).Return([]byte(tc.returnData), tc.returnErr)
			metadata, err := acache.Get(tc.applicationID)
			require.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expected, metadata)
		})
	}
}
