// Copyright 2022 The PipeCD Authors.
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

package insight

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/pipe-cd/pipecd/pkg/model"
)

var jst = time.FixedZone("Asia/Tokyo", 9*60*60)

func TestNoramalizeTime(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name     string
		t        int64
		expected int64
	}{
		{
			name:     "Mid time",
			t:        1648339200, // 2022/03/27:09:00:00 JST
			expected: 1648306800, // 2022/03/27:00:00:00 JST
		},
		{
			name:     "No change",
			t:        1648306800, // 2022/03/27:00:00:00 JST
			expected: 1648306800, // 2022/03/27:00:00:00 JST
		},
		{
			name:     "23:59:59",
			t:        1648393199, // 2022/03/27:23:59:59 JST
			expected: 1648306800, // 2022/03/27:00:00:00 JST
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := NormalizeUnixTime(tc.t, jst)
			assert.Equal(t, tc.expected, got)
		})
	}
}

func TestGroupDeploymentByDaily(t *testing.T) {
	testcases := []struct {
		name        string
		deployments []*model.InsightDeployment
		expected    [][]*model.InsightDeployment
	}{
		{
			name: "Success",
			deployments: []*model.InsightDeployment{
				{
					Id:          "0",
					CompletedAt: 1648306799, // 2022/03/26:23:59:59 JST
				},
				{
					Id:          "1",
					CompletedAt: 1648306800, // 2022/03/27:00:00:00 JST
				},
				{
					Id:          "2",
					CompletedAt: 1648339200, // 2022/03/27:09:00:00 JST
				},
				{
					Id:          "3",
					CompletedAt: 1648393199, // 2022/03/27:23:59:59 JST
				},
				{
					Id:          "4",
					CompletedAt: 1648393200, // 2022/03/28:00:00:00 JST
				},
			},
			expected: [][]*model.InsightDeployment{
				{{
					Id:          "0",
					CompletedAt: 1648306799, // 2022/03/26:23:59:59 JST
				}},
				{
					{
						Id:          "1",
						CompletedAt: 1648306800, // 2022/03/27:00:00:00 JST
					},
					{
						Id:          "2",
						CompletedAt: 1648339200, // 2022/03/27:09:00:00 JST
					},
					{
						Id:          "3",
						CompletedAt: 1648393199, // 2022/03/27:23:59:59 JST
					},
				},
				{{
					Id:          "4",
					CompletedAt: 1648393200, // 2022/03/28:00:00:00 JST
				}},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := GroupDeploymentsByDaily(tc.deployments, jst)
			assert.Equal(t, tc.expected, got)
		})
	}
}
