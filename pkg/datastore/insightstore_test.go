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

package datastore

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/pipe-cd/pipe/pkg/model"
)

func TestAddInsight(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testcases := []struct {
		name      string
		insight   *model.InsightDataPoint
		dsFactory func(*model.InsightDataPoint) DataStore
		wantErr   bool
	}{
		{
			name:      "Invalid insight",
			insight:   &model.InsightDataPoint{},
			dsFactory: func(d *model.InsightDataPoint) DataStore { return nil },
			wantErr:   true,
		},
		{
			name: "Valid insight",
			insight: &model.InsightDataPoint{
				Id:            "ID",
				Timestamp:     1,
				Value:         0,
				ApplicationId: "AppID",
				EnvId:         "EnvID",
				PipedId:       "PipedId",
				ProjectId:     "ProjectId",
				MetricsKind:   1,
				Step:          1,
				CreatedAt:     1,
			},
			dsFactory: func(d *model.InsightDataPoint) DataStore {
				ds := NewMockDataStore(ctrl)
				ds.EXPECT().Create(gomock.Any(), "Insight", d.Id, d)
				return ds
			},
			wantErr: false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			s := NewInsightStore(tc.dsFactory(tc.insight))
			err := s.AddInsight(context.Background(), tc.insight)
			assert.Equal(t, tc.wantErr, err != nil)
		})
	}
}

func TestListInsight(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testcases := []struct {
		name    string
		opts    ListOptions
		ds      DataStore
		wantErr bool
	}{
		{
			name: "iterator done",
			opts: ListOptions{Page: 1},
			ds: func() DataStore {
				it := NewMockIterator(ctrl)
				it.EXPECT().
					Next(&model.InsightDataPoint{}).
					Return(ErrIteratorDone)

				ds := NewMockDataStore(ctrl)
				ds.EXPECT().
					Find(gomock.Any(), "Insight", ListOptions{Page: 1}).
					Return(it, nil)
				return ds
			}(),
			wantErr: false,
		},
		{
			name: "unexpected error occurred",
			opts: ListOptions{Page: 1},
			ds: func() DataStore {
				it := NewMockIterator(ctrl)
				it.EXPECT().
					Next(&model.InsightDataPoint{}).
					Return(fmt.Errorf("err"))

				ds := NewMockDataStore(ctrl)
				ds.EXPECT().
					Find(gomock.Any(), "Insight", ListOptions{Page: 1}).
					Return(it, nil)
				return ds
			}(),
			wantErr: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			s := NewInsightStore(tc.ds)
			_, err := s.ListInsights(context.Background(), tc.opts)
			assert.Equal(t, tc.wantErr, err != nil)
		})
	}
}
