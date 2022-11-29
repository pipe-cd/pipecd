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

	"github.com/stretchr/testify/assert"

	"github.com/pipe-cd/pipecd/pkg/model"
)

func TestBuildDeploymentFrequencyDataPoint(t *testing.T) {
	testcases := []struct {
		name     string
		ds       []*DeploymentData
		step     model.InsightStep
		expected []*model.InsightDataPoint
	}{
		{
			name:     "nil",
			step:     model.InsightStep_DAILY,
			expected: []*model.InsightDataPoint{},
		},
		{
			name:     "empty",
			ds:       []*DeploymentData{},
			step:     model.InsightStep_DAILY,
			expected: []*model.InsightDataPoint{},
		},
		{
			name: "daily step",
			ds: []*DeploymentData{
				&DeploymentData{
					CompletedAt: 1669574625,
				},
				&DeploymentData{
					CompletedAt: 1669574635,
				},
				&DeploymentData{
					CompletedAt: 1669661030,
				},
			},
			step: model.InsightStep_DAILY,
			expected: []*model.InsightDataPoint{
				&model.InsightDataPoint{
					Timestamp: 1669507200,
					Value:     2,
				},
				&model.InsightDataPoint{
					Timestamp: 1669593600,
					Value:     1,
				},
			},
		},
		{
			name: "monthly step",
			ds: []*DeploymentData{
				&DeploymentData{
					CompletedAt: 1666982630,
				},
				&DeploymentData{
					CompletedAt: 1666982635,
				},
				&DeploymentData{
					CompletedAt: 1669661010,
				},
			},
			step: model.InsightStep_MONTHLY,
			expected: []*model.InsightDataPoint{
				&model.InsightDataPoint{
					Timestamp: 1664582400,
					Value:     2,
				},
				&model.InsightDataPoint{
					Timestamp: 1667260800,
					Value:     1,
				},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := buildDeploymentFrequencyDataPoints(tc.ds, "", nil, tc.step)
			assert.Equal(t, tc.expected, got)
		})
	}
}

func TestBuildDeploymentChangeFailureRateDataPoint(t *testing.T) {
	testcases := []struct {
		name     string
		ds       []*DeploymentData
		step     model.InsightStep
		expected []*model.InsightDataPoint
	}{
		{
			name:     "nil",
			step:     model.InsightStep_DAILY,
			expected: []*model.InsightDataPoint{},
		},
		{
			name:     "empty",
			ds:       []*DeploymentData{},
			step:     model.InsightStep_DAILY,
			expected: []*model.InsightDataPoint{},
		},
		{
			name: "daily step",
			ds: []*DeploymentData{
				&DeploymentData{
					CompletedAt: 1669340910,
				},
				&DeploymentData{
					CompletedAt: 1669340920,
				},
				&DeploymentData{
					CompletedAt:    1669600130,
					CompleteStatus: model.DeploymentStatus_DEPLOYMENT_FAILURE.String(),
				},
				&DeploymentData{
					CompletedAt:    1669686600,
					CompleteStatus: model.DeploymentStatus_DEPLOYMENT_FAILURE.String(),
				},
				&DeploymentData{
					CompletedAt: 1669686610,
				},
			},
			step: model.InsightStep_DAILY,
			expected: []*model.InsightDataPoint{
				&model.InsightDataPoint{
					Timestamp: 1669334400,
					Value:     0,
				},
				&model.InsightDataPoint{
					Timestamp: 1669593600,
					Value:     1,
				},
				&model.InsightDataPoint{
					Timestamp: 1669680000,
					Value:     0.5,
				},
			},
		},
		{
			name: "monthly step",
			ds: []*DeploymentData{
				&DeploymentData{
					CompletedAt: 1664416110,
				},
				&DeploymentData{
					CompletedAt: 1664416120,
				},
				&DeploymentData{
					CompletedAt:    1667008110,
					CompleteStatus: model.DeploymentStatus_DEPLOYMENT_FAILURE.String(),
				},
				&DeploymentData{
					CompletedAt:    1668908910,
					CompleteStatus: model.DeploymentStatus_DEPLOYMENT_FAILURE.String(),
				},
				&DeploymentData{
					CompletedAt: 1668908920,
				},
			},
			step: model.InsightStep_MONTHLY,
			expected: []*model.InsightDataPoint{
				&model.InsightDataPoint{
					Timestamp: 1661990400,
					Value:     0,
				},
				&model.InsightDataPoint{
					Timestamp: 1664582400,
					Value:     1,
				},
				&model.InsightDataPoint{
					Timestamp: 1667260800,
					Value:     0.5,
				},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := buildDeploymentChangeFailureRateDataPoints(tc.ds, "", nil, tc.step)
			assert.Equal(t, tc.expected, got)
		})
	}
}

func TestFillUpDataPoints(t *testing.T) {
	testcases := []struct {
		name     string
		ds       []*model.InsightDataPoint
		step     model.InsightStep
		from, to int64
		want     []*model.InsightDataPoint
	}{
		{
			name: "daily step: missing head part",
			ds: []*model.InsightDataPoint{
				&model.InsightDataPoint{Timestamp: 172800, Value: 2},
				&model.InsightDataPoint{Timestamp: 259200, Value: 3},
			},
			from: 86400,
			to:   259200,
			step: model.InsightStep_DAILY,
			want: []*model.InsightDataPoint{
				&model.InsightDataPoint{Timestamp: 86400, Value: 0},
				&model.InsightDataPoint{Timestamp: 172800, Value: 2},
				&model.InsightDataPoint{Timestamp: 259200, Value: 3},
			},
		},
		{
			name: "daily step: missing tail part",
			ds: []*model.InsightDataPoint{
				&model.InsightDataPoint{Timestamp: 86400, Value: 1},
				&model.InsightDataPoint{Timestamp: 172800, Value: 2},
			},
			from: 86400,
			to:   259200,
			step: model.InsightStep_DAILY,
			want: []*model.InsightDataPoint{
				&model.InsightDataPoint{Timestamp: 86400, Value: 1},
				&model.InsightDataPoint{Timestamp: 172800, Value: 2},
				&model.InsightDataPoint{Timestamp: 259200, Value: 0},
			},
		},
		{
			name: "daily step: missing both parts",
			ds: []*model.InsightDataPoint{
				&model.InsightDataPoint{Timestamp: 172800, Value: 2},
			},
			from: 86400,
			to:   259200,
			step: model.InsightStep_DAILY,
			want: []*model.InsightDataPoint{
				&model.InsightDataPoint{Timestamp: 86400, Value: 0},
				&model.InsightDataPoint{Timestamp: 172800, Value: 2},
				&model.InsightDataPoint{Timestamp: 259200, Value: 0},
			},
		},
		{
			name: "monthly step: missing head part",
			ds: []*model.InsightDataPoint{
				&model.InsightDataPoint{Timestamp: 1664582400, Value: 2}, // 2022/10
				&model.InsightDataPoint{Timestamp: 1667260800, Value: 3}, // 2022/11
			},
			from: 1661990401, // 2022/9
			to:   1667260801, // 2022/11
			step: model.InsightStep_MONTHLY,
			want: []*model.InsightDataPoint{
				&model.InsightDataPoint{Timestamp: 1661990400, Value: 0},
				&model.InsightDataPoint{Timestamp: 1664582400, Value: 2},
				&model.InsightDataPoint{Timestamp: 1667260800, Value: 3},
			},
		},
		{
			name: "monthly step: missing tail part",
			ds: []*model.InsightDataPoint{
				&model.InsightDataPoint{Timestamp: 1661990400, Value: 2}, // 2022/9
				&model.InsightDataPoint{Timestamp: 1664582400, Value: 3}, // 2022/10
			},
			from: 1661990401, // 2022/9
			to:   1667260801, // 2022/11
			step: model.InsightStep_MONTHLY,
			want: []*model.InsightDataPoint{
				&model.InsightDataPoint{Timestamp: 1661990400, Value: 2},
				&model.InsightDataPoint{Timestamp: 1664582400, Value: 3},
				&model.InsightDataPoint{Timestamp: 1667260800, Value: 0},
			},
		},
		{
			name: "monthly step: missing both parts",
			ds: []*model.InsightDataPoint{
				&model.InsightDataPoint{Timestamp: 1664582400, Value: 2}, // 2022/10
				&model.InsightDataPoint{Timestamp: 1667260800, Value: 3}, // 2022/11
			},
			from: 1661990401, // 2022/9
			to:   1673344801, // 2023/1
			step: model.InsightStep_MONTHLY,
			want: []*model.InsightDataPoint{
				&model.InsightDataPoint{Timestamp: 1661990400, Value: 0},
				&model.InsightDataPoint{Timestamp: 1664582400, Value: 2},
				&model.InsightDataPoint{Timestamp: 1667260800, Value: 3},
				&model.InsightDataPoint{Timestamp: 1669852800, Value: 0},
				&model.InsightDataPoint{Timestamp: 1672531200, Value: 0},
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := fillUpDataPoints(tc.ds, tc.from, tc.to, tc.step)
			assert.Equal(t, tc.want, got)
		})
	}
}
