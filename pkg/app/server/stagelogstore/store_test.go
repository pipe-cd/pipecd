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

	"github.com/pipe-cd/pipecd/pkg/model"
)

func TestMergeBlocks(t *testing.T) {
	testcases := []struct {
		name string

		before []*model.LogBlock
		after  []*model.LogBlock

		expected []*model.LogBlock
	}{
		{
			name:   "before is empty",
			before: []*model.LogBlock{},
			after: []*model.LogBlock{
				{
					Index:     1590000011,
					Log:       "log-1",
					Severity:  model.LogSeverity_SUCCESS,
					CreatedAt: 1590000010,
				},
				{
					Index:     1590000012,
					Log:       "log-2",
					Severity:  model.LogSeverity_ERROR,
					CreatedAt: 1590000020,
				},
			},
			expected: []*model.LogBlock{
				{
					Index:     1590000011,
					Log:       "log-1",
					Severity:  model.LogSeverity_SUCCESS,
					CreatedAt: 1590000010,
				},
				{
					Index:     1590000012,
					Log:       "log-2",
					Severity:  model.LogSeverity_ERROR,
					CreatedAt: 1590000020,
				},
			},
		},
		{
			name: "after is empty",
			before: []*model.LogBlock{
				{
					Index:     1590000011,
					Log:       "log-1",
					Severity:  model.LogSeverity_SUCCESS,
					CreatedAt: 1590000010,
				},
				{
					Index:     1590000012,
					Log:       "log-2",
					Severity:  model.LogSeverity_ERROR,
					CreatedAt: 1590000020,
				},
			},
			after: []*model.LogBlock{},
			expected: []*model.LogBlock{
				{
					Index:     1590000011,
					Log:       "log-1",
					Severity:  model.LogSeverity_SUCCESS,
					CreatedAt: 1590000010,
				},
				{
					Index:     1590000012,
					Log:       "log-2",
					Severity:  model.LogSeverity_ERROR,
					CreatedAt: 1590000020,
				},
			},
		},
		{
			name: "append without duplicating",
			before: []*model.LogBlock{
				{
					Index:     1590000011,
					Log:       "log-1",
					Severity:  model.LogSeverity_SUCCESS,
					CreatedAt: 1590000010,
				},
				{
					Index:     1590000012,
					Log:       "log-2",
					Severity:  model.LogSeverity_ERROR,
					CreatedAt: 1590000020,
				},
			},
			after: []*model.LogBlock{
				{
					Index:     1590000013,
					Log:       "log-3",
					Severity:  model.LogSeverity_SUCCESS,
					CreatedAt: 1590000030,
				},
				{
					Index:     1590000014,
					Log:       "log-4",
					Severity:  model.LogSeverity_ERROR,
					CreatedAt: 1590000040,
				},
			},
			expected: []*model.LogBlock{
				{
					Index:     1590000011,
					Log:       "log-1",
					Severity:  model.LogSeverity_SUCCESS,
					CreatedAt: 1590000010,
				},
				{
					Index:     1590000012,
					Log:       "log-2",
					Severity:  model.LogSeverity_ERROR,
					CreatedAt: 1590000020,
				},
				{
					Index:     1590000013,
					Log:       "log-3",
					Severity:  model.LogSeverity_SUCCESS,
					CreatedAt: 1590000030,
				},
				{
					Index:     1590000014,
					Log:       "log-4",
					Severity:  model.LogSeverity_ERROR,
					CreatedAt: 1590000040,
				},
			},
		},

		{
			name: "append with deduplicating",
			before: []*model.LogBlock{
				{
					Index:     1590000011,
					Log:       "log-1",
					Severity:  model.LogSeverity_SUCCESS,
					CreatedAt: 1590000010,
				},
				{
					Index:     1590000012,
					Log:       "log-2",
					Severity:  model.LogSeverity_ERROR,
					CreatedAt: 1590000020,
				},
			},
			after: []*model.LogBlock{
				{
					Index:     1590000012,
					Log:       "log-2",
					Severity:  model.LogSeverity_ERROR,
					CreatedAt: 1590000020,
				},
				{
					Index:     1590000013,
					Log:       "log-3",
					Severity:  model.LogSeverity_SUCCESS,
					CreatedAt: 1590000030,
				},
			},
			expected: []*model.LogBlock{
				{
					Index:     1590000011,
					Log:       "log-1",
					Severity:  model.LogSeverity_SUCCESS,
					CreatedAt: 1590000010,
				},
				{
					Index:     1590000012,
					Log:       "log-2",
					Severity:  model.LogSeverity_ERROR,
					CreatedAt: 1590000020,
				},
				{
					Index:     1590000013,
					Log:       "log-3",
					Severity:  model.LogSeverity_SUCCESS,
					CreatedAt: 1590000030,
				},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			actual := mergeBlocks(tc.before, tc.after)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
