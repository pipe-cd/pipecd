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
		expected []*model.InsightDataPoint
	}{
		{
			name:     "nil",
			expected: []*model.InsightDataPoint{},
		},
		{
			name:     "empty",
			ds:       []*DeploymentData{},
			expected: []*model.InsightDataPoint{},
		},
		{
			name: "ok",
			ds: []*DeploymentData{
				&DeploymentData{
					CompletedAt:    101,
					CompletedAtDay: 100,
				},
				&DeploymentData{
					CompletedAt:    102,
					CompletedAtDay: 100,
				},
				&DeploymentData{
					CompletedAt:    205,
					CompletedAtDay: 200,
				},
			},
			expected: []*model.InsightDataPoint{
				&model.InsightDataPoint{
					Timestamp: 100,
					Value:     2,
				},
				&model.InsightDataPoint{
					Timestamp: 200,
					Value:     1,
				},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := buildDeploymentFrequencyDataPoints(tc.ds, "", nil)
			assert.Equal(t, tc.expected, got)
		})
	}
}

func TestBuildDeploymentChangeFailureRateDataPoint(t *testing.T) {
	testcases := []struct {
		name     string
		ds       []*DeploymentData
		expected []*model.InsightDataPoint
	}{
		{
			name:     "nil",
			expected: []*model.InsightDataPoint{},
		},
		{
			name:     "empty",
			ds:       []*DeploymentData{},
			expected: []*model.InsightDataPoint{},
		},
		{
			name: "ok",
			ds: []*DeploymentData{
				&DeploymentData{
					CompletedAtDay: 100,
				},
				&DeploymentData{
					CompletedAtDay: 100,
				},
				&DeploymentData{
					CompletedAtDay: 200,
					CompleteStatus: model.DeploymentStatus_DEPLOYMENT_FAILURE.String(),
				},
				&DeploymentData{
					CompletedAtDay: 300,
					CompleteStatus: model.DeploymentStatus_DEPLOYMENT_FAILURE.String(),
				},
				&DeploymentData{
					CompletedAtDay: 300,
				},
			},
			expected: []*model.InsightDataPoint{
				&model.InsightDataPoint{
					Timestamp: 100,
					Value:     0,
				},
				&model.InsightDataPoint{
					Timestamp: 200,
					Value:     1,
				},
				&model.InsightDataPoint{
					Timestamp: 300,
					Value:     0.5,
				},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := buildDeploymentChangeFailureRateDataPoints(tc.ds, "", nil)
			assert.Equal(t, tc.expected, got)
		})
	}
}

func TestFillUpDataPoints(t *testing.T) {
	testcases := []struct {
		name     string
		ds       []*model.InsightDataPoint
		from, to int64
		want     []*model.InsightDataPoint
	}{
		{
			name: "missing head part",
			ds: []*model.InsightDataPoint{
				&model.InsightDataPoint{Timestamp: 172800, Value: 2},
				&model.InsightDataPoint{Timestamp: 259200, Value: 3},
			},
			from: 86400,
			to:   259200,
			want: []*model.InsightDataPoint{
				&model.InsightDataPoint{Timestamp: 86400, Value: 0},
				&model.InsightDataPoint{Timestamp: 172800, Value: 2},
				&model.InsightDataPoint{Timestamp: 259200, Value: 3},
			},
		},
		{
			name: "missing tail part",
			ds: []*model.InsightDataPoint{
				&model.InsightDataPoint{Timestamp: 86400, Value: 1},
				&model.InsightDataPoint{Timestamp: 172800, Value: 2},
			},
			from: 86400,
			to:   259200,
			want: []*model.InsightDataPoint{
				&model.InsightDataPoint{Timestamp: 86400, Value: 1},
				&model.InsightDataPoint{Timestamp: 172800, Value: 2},
				&model.InsightDataPoint{Timestamp: 259200, Value: 0},
			},
		},
		{
			name: "missing both parts",
			ds: []*model.InsightDataPoint{
				&model.InsightDataPoint{Timestamp: 172800, Value: 2},
			},
			from: 86400,
			to:   259200,
			want: []*model.InsightDataPoint{
				&model.InsightDataPoint{Timestamp: 86400, Value: 0},
				&model.InsightDataPoint{Timestamp: 172800, Value: 2},
				&model.InsightDataPoint{Timestamp: 259200, Value: 0},
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := fillUpDataPoints(tc.ds, tc.from, tc.to)
			assert.Equal(t, tc.want, got)
		})
	}
}
