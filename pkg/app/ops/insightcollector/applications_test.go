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

package insightcollector

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/pipe-cd/pipe/pkg/datastore"
	"github.com/pipe-cd/pipe/pkg/datastore/datastoretest"
	"github.com/pipe-cd/pipe/pkg/filestore"
	"github.com/pipe-cd/pipe/pkg/insight"
	"github.com/pipe-cd/pipe/pkg/insight/insightstore/insightstoretest"
	"github.com/pipe-cd/pipe/pkg/model"
)

func TestInsightCollector_getApplications(t *testing.T) {
	tests := []struct {
		name                          string
		prepareApplicationStoreMockFn func(m *datastoretest.MockApplicationStore)
		to                            time.Time
		want                          []*model.Application
		wantErr                       bool
	}{
		{
			name: "get less than 50(pageSize) applications",
			to:   time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
			prepareApplicationStoreMockFn: func(m *datastoretest.MockApplicationStore) {
				m.EXPECT().ListApplications(gomock.Any(), datastore.ListOptions{
					Limit: limit,
					Filters: []datastore.ListFilter{
						{
							Field:    "CreatedAt",
							Operator: "<",
							Value:    time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC).Unix(),
						},
					},
					Orders: []datastore.Order{
						{
							Field:     "CreatedAt",
							Direction: datastore.Desc,
						},
					}}).Return([]*model.Application{
					{
						Id: "1",
					},
					{
						Id: "2",
					},
				}, nil)
			},
			want: []*model.Application{
				{
					Id: "1",
				},
				{
					Id: "2",
				},
			},
		},
		{
			name: "get more than 50(pageSize) applications",
			to:   time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
			prepareApplicationStoreMockFn: func(m *datastoretest.MockApplicationStore) {
				m.EXPECT().ListApplications(gomock.Any(), datastore.ListOptions{
					Limit: limit,
					Filters: []datastore.ListFilter{
						{
							Field:    "CreatedAt",
							Operator: "<",
							Value:    time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC).Unix(),
						},
					},
					Orders: []datastore.Order{
						{
							Field:     "CreatedAt",
							Direction: datastore.Desc,
						},
					}}).Return([]*model.Application{
					{}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {},
					{}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {},
					{}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {},
					{
						CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
					},
				}, nil)
				m.EXPECT().ListApplications(gomock.Any(), datastore.ListOptions{
					Limit: limit,
					Filters: []datastore.ListFilter{
						{
							Field:    "CreatedAt",
							Operator: "<",
							Value:    time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
						},
					},
					Orders: []datastore.Order{
						{
							Field:     "CreatedAt",
							Direction: datastore.Desc,
						},
					}}).Return([]*model.Application{
					{}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {},
				}, nil)
			},
			want: []*model.Application{
				{}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {},
				{}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {},
				{}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {},
				{
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
				},
				{}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {},
			},
		},
		{
			name: "return error when ListApplications fail",
			to:   time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
			prepareApplicationStoreMockFn: func(m *datastoretest.MockApplicationStore) {
				m.EXPECT().ListApplications(gomock.Any(), datastore.ListOptions{
					Limit: limit,
					Filters: []datastore.ListFilter{
						{
							Field:    "CreatedAt",
							Operator: "<",
							Value:    time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC).Unix(),
						},
					},
					Orders: []datastore.Order{
						{
							Field:     "CreatedAt",
							Direction: datastore.Desc,
						},
					}}).Return(nil, errors.New("error"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAppStore := datastoretest.NewMockApplicationStore(ctrl)
			tt.prepareApplicationStoreMockFn(mockAppStore)
			i := &InsightCollector{
				applicationStore: mockAppStore,
				logger:           zap.NewNop(),
			}
			got, err := i.getApplications(context.Background(), tt.to)
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

func TestInsightCollector_updateApplicationCount(t *testing.T) {
	applicationCount := func() *insight.ApplicationCount {
		// init application count
		ac := insight.NewApplicationCount()

		for i := 0; i < len(ac.Counts); i++ {
			c := &ac.Counts[i]
			enable := insight.ApplicationCountLabelSet{
				Kind:   model.ApplicationKind_CLOUDRUN,
				Status: insight.ApplicationStatusEnable,
			}
			if c.LabelSet == enable {
				c.Count = 1
			}
		}
		return ac
	}()
	tests := []struct {
		name                      string
		prepareInsightstoreMockFn func(m *insightstoretest.MockStore)
		apps                      []*model.Application
		pid                       string
		target                    time.Time
		wantErr                   bool
	}{
		{
			name: "success",
			prepareInsightstoreMockFn: func(m *insightstoretest.MockStore) {
				m.EXPECT().LoadApplicationCount(gomock.Any(), "12").Return(applicationCount, nil)
				m.EXPECT().PutApplicationCount(gomock.Any(), func() *insight.ApplicationCount {
					ac := insight.NewApplicationCount()
					for i := 0; i < len(ac.Counts); i++ {
						c := &ac.Counts[i]
						enable := insight.ApplicationCountLabelSet{
							Kind:   model.ApplicationKind_CLOUDRUN,
							Status: insight.ApplicationStatusEnable,
						}
						delete := insight.ApplicationCountLabelSet{
							Kind:   model.ApplicationKind_CLOUDRUN,
							Status: insight.ApplicationStatusDeleted,
						}
						disable := insight.ApplicationCountLabelSet{
							Kind:   model.ApplicationKind_CLOUDRUN,
							Status: insight.ApplicationStatusDisable,
						}
						if c.LabelSet == enable {
							c.Count = 2
						}
						if c.LabelSet == delete {
							c.Count = 1
						}
						if c.LabelSet == disable {
							c.Count = 1
						}
					}

					ac.AccumulatedTo = time.Date(2020, 1, 1, 4, 0, 0, 0, time.UTC).Unix()

					return ac
				}(), "12").Return(nil)
			},
			pid: "12",
			apps: []*model.Application{
				{
					Kind:      model.ApplicationKind_CLOUDRUN,
					Disabled:  true,
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC).Unix(),
				},
				{
					Kind:      model.ApplicationKind_CLOUDRUN,
					Deploying: true,
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC).Unix(),
				},
				{
					Kind:      model.ApplicationKind_CLOUDRUN,
					Deleted:   true,
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC).Unix(),
				},
				{
					Kind:      model.ApplicationKind_CLOUDRUN,
					Deploying: true,
					CreatedAt: time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC).Unix(),
				},
			},
			target: time.Date(2020, 1, 1, 4, 0, 0, 0, time.UTC),
		},
		{
			name: "success even if application count was not found",
			prepareInsightstoreMockFn: func(m *insightstoretest.MockStore) {
				m.EXPECT().LoadApplicationCount(gomock.Any(), "12").Return(nil, filestore.ErrNotFound)
				m.EXPECT().PutApplicationCount(gomock.Any(), func() *insight.ApplicationCount {
					ac := insight.NewApplicationCount()
					for i := 0; i < len(ac.Counts); i++ {
						c := &ac.Counts[i]
						enable := insight.ApplicationCountLabelSet{
							Kind:   model.ApplicationKind_CLOUDRUN,
							Status: insight.ApplicationStatusEnable,
						}
						delete := insight.ApplicationCountLabelSet{
							Kind:   model.ApplicationKind_CLOUDRUN,
							Status: insight.ApplicationStatusDeleted,
						}
						disable := insight.ApplicationCountLabelSet{
							Kind:   model.ApplicationKind_CLOUDRUN,
							Status: insight.ApplicationStatusDisable,
						}
						if c.LabelSet == enable {
							c.Count = 2
						}
						if c.LabelSet == delete {
							c.Count = 1
						}
						if c.LabelSet == disable {
							c.Count = 1
						}
					}

					ac.AccumulatedFrom = time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC).Unix()
					ac.AccumulatedTo = time.Date(2020, 1, 1, 4, 0, 0, 0, time.UTC).Unix()

					return ac
				}(), "12").Return(nil)
			},
			pid: "12",
			apps: []*model.Application{
				{
					Kind:      model.ApplicationKind_CLOUDRUN,
					Disabled:  true,
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC).Unix(),
				},
				{
					Kind:      model.ApplicationKind_CLOUDRUN,
					Deploying: true,
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC).Unix(),
				},
				{
					Kind:      model.ApplicationKind_CLOUDRUN,
					Deleted:   true,
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC).Unix(),
				},
				{
					Kind:      model.ApplicationKind_CLOUDRUN,
					Deploying: true,
					CreatedAt: time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC).Unix(),
				},
			},
			target: time.Date(2020, 1, 1, 4, 0, 0, 0, time.UTC),
		},
		{
			name: "fail when failed to put application count",
			prepareInsightstoreMockFn: func(m *insightstoretest.MockStore) {
				m.EXPECT().LoadApplicationCount(gomock.Any(), "12").Return(applicationCount, nil)
				m.EXPECT().PutApplicationCount(gomock.Any(), func() *insight.ApplicationCount {
					ac := insight.NewApplicationCount()
					for i := 0; i < len(ac.Counts); i++ {
						c := &ac.Counts[i]
						enable := insight.ApplicationCountLabelSet{
							Kind:   model.ApplicationKind_CLOUDRUN,
							Status: insight.ApplicationStatusEnable,
						}
						delete := insight.ApplicationCountLabelSet{
							Kind:   model.ApplicationKind_CLOUDRUN,
							Status: insight.ApplicationStatusDeleted,
						}
						disable := insight.ApplicationCountLabelSet{
							Kind:   model.ApplicationKind_CLOUDRUN,
							Status: insight.ApplicationStatusDisable,
						}
						if c.LabelSet == enable {
							c.Count = 2
						}
						if c.LabelSet == delete {
							c.Count = 1
						}
						if c.LabelSet == disable {
							c.Count = 1
						}
					}

					ac.AccumulatedTo = time.Date(2020, 1, 1, 4, 0, 0, 0, time.UTC).Unix()

					return ac
				}(), "12").Return(errors.New("error"))
			},
			pid: "12",
			apps: []*model.Application{
				{
					Kind:      model.ApplicationKind_CLOUDRUN,
					Disabled:  true,
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC).Unix(),
				},
				{
					Kind:      model.ApplicationKind_CLOUDRUN,
					Deploying: true,
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC).Unix(),
				},
				{
					Kind:      model.ApplicationKind_CLOUDRUN,
					Deleted:   true,
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC).Unix(),
				},
				{
					Kind:      model.ApplicationKind_CLOUDRUN,
					Deploying: true,
					CreatedAt: time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC).Unix(),
				},
			},
			target:  time.Date(2020, 1, 1, 4, 0, 0, 0, time.UTC),
			wantErr: true,
		},
		{
			name: "fail when failed to get application count",
			prepareInsightstoreMockFn: func(m *insightstoretest.MockStore) {
				m.EXPECT().LoadApplicationCount(gomock.Any(), "12").Return(nil, errors.New("error"))
			},
			pid: "12",
			apps: []*model.Application{
				{
					Kind:      model.ApplicationKind_CLOUDRUN,
					Deploying: true,
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC).Unix(),
				},
				{
					Kind:      model.ApplicationKind_CLOUDRUN,
					Deploying: true,
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC).Unix(),
				},
				{
					Kind:      model.ApplicationKind_CLOUDRUN,
					Deleted:   true,
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC).Unix(),
				},
				{
					Kind:      model.ApplicationKind_CLOUDRUN,
					Disabled:  true,
					CreatedAt: time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC).Unix(),
				},
			},
			target:  time.Date(2020, 1, 1, 4, 0, 0, 0, time.UTC),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mock := insightstoretest.NewMockStore(ctrl)
			tt.prepareInsightstoreMockFn(mock)
			i := &InsightCollector{
				insightstore: mock,
				logger:       zap.NewNop(),
			}
			err := i.updateApplicationCount(context.Background(), tt.apps, tt.pid, tt.target)
			if (err != nil) != tt.wantErr {
				if !tt.wantErr {
					assert.NoError(t, err)
					return
				}
				assert.Error(t, err, tt.wantErr)
				return
			}
		})
	}
}
