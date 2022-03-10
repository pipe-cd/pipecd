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

package insight

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/pipe-cd/pipecd/pkg/model"
)

func Test_ExtractDataPoints(t *testing.T) {
	t.Parallel()

	type args struct {
		datapoints []DataPoint
		from       time.Time
		to         time.Time
	}
	tests := []struct {
		name    string
		args    args
		want    []*model.InsightDataPoint
		wantErr bool
	}{
		{
			name: "success with yearly",
			args: args{
				datapoints: func() []DataPoint {
					df := []*DeployFrequency{
						{
							DeployCount: 1000,
							Timestamp:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
						},
						{
							DeployCount: 3000,
							Timestamp:   time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
						},
					}
					dp, _ := ToDataPoints(df)
					return dp
				}(),
				from: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				to:   time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			want: []*model.InsightDataPoint{
				{
					Timestamp: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
					Value:     1000,
				},
				{
					Timestamp: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
					Value:     3000,
				},
			},
		},
		{
			name: "success with monthly",
			args: args{
				datapoints: func() []DataPoint {
					df := []*DeployFrequency{
						{
							DeployCount: 1000,
							Timestamp:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
						},
					}
					dp, _ := ToDataPoints(df)
					return dp
				}(),
				from: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				to:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			want: []*model.InsightDataPoint{
				{
					Timestamp: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
					Value:     1000,
				},
			},
		},
		{
			name: "success with weekly",
			args: args{
				datapoints: func() []DataPoint {
					df := []*DeployFrequency{
						{
							DeployCount: 1000,
							Timestamp:   time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC).Unix(),
						},
						{
							DeployCount: 3000,
							Timestamp:   time.Date(2021, 1, 10, 0, 0, 0, 0, time.UTC).Unix(),
						},
					}
					dp, _ := ToDataPoints(df)
					return dp
				}(),
				from: time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC),
				to:   time.Date(2021, 1, 10, 0, 0, 0, 0, time.UTC),
			},
			want: []*model.InsightDataPoint{
				{
					Timestamp: time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC).Unix(),
					Value:     1000,
				},
				{
					Timestamp: time.Date(2021, 1, 10, 0, 0, 0, 0, time.UTC).Unix(),
					Value:     3000,
				},
			},
		},
		{
			name: "success with daily",
			args: args{
				datapoints: func() []DataPoint {
					df := []*DeployFrequency{
						{
							DeployCount: 1000,
							Timestamp:   time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC).Unix(),
						},
						{
							DeployCount: 3000,
							Timestamp:   time.Date(2021, 1, 4, 0, 0, 0, 0, time.UTC).Unix(),
						},
					}
					dp, _ := ToDataPoints(df)
					return dp
				}(),
				from: time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC),
				to:   time.Date(2021, 1, 4, 0, 0, 0, 0, time.UTC),
			},
			want: []*model.InsightDataPoint{
				{
					Timestamp: time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC).Unix(),
					Value:     1000,
				},
				{
					Timestamp: time.Date(2021, 1, 4, 0, 0, 0, 0, time.UTC).Unix(),
					Value:     3000,
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := extractDataPoints(tt.args.datapoints, tt.args.from, tt.args.to)
			if (err != nil) != tt.wantErr {
				if !tt.wantErr {
					assert.NoError(t, err)
					return
				}
				assert.Error(t, err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
