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

package insight

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/pipe-cd/pipe/pkg/model"
)

func TestApplicationCount_Find(t *testing.T) {
	applicationCount := &ApplicationCount{
		Counts: []ApplicationCountByLabelSet{
			{
				LabelSet: ApplicationCountLabelSet{
					Kind:   model.ApplicationKind_CLOUDRUN,
					Status: ApplicationStatusEnable,
				},
				Count: 1,
			},
		},
	}
	tests := []struct {
		name     string
		ac       *ApplicationCount
		labelSet ApplicationCountLabelSet
		want     ApplicationCountByLabelSet
		wantErr  bool
	}{
		{
			name: "success",
			ac:   applicationCount,
			labelSet: ApplicationCountLabelSet{
				Kind:   model.ApplicationKind_CLOUDRUN,
				Status: ApplicationStatusEnable,
			},
			want: ApplicationCountByLabelSet{
				LabelSet: ApplicationCountLabelSet{
					Kind:   model.ApplicationKind_CLOUDRUN,
					Status: ApplicationStatusEnable,
				},
				Count: 1,
			},
		},
		{
			name: "not found",
			ac:   applicationCount,
			labelSet: ApplicationCountLabelSet{
				Kind:   model.ApplicationKind_CLOUDRUN,
				Status: ApplicationStatusDisable,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.ac.Find(tt.labelSet)
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

func TestApplicationCount_MigrateApplicationCount(t *testing.T) {
	applicationCount := &ApplicationCount{
		Counts: []ApplicationCountByLabelSet{
			{
				LabelSet: ApplicationCountLabelSet{
					Kind:   model.ApplicationKind_CLOUDRUN,
					Status: ApplicationStatusEnable,
				},
				Count: 1,
			},
		},
	}
	tests := []struct {
		name string
		ac   *ApplicationCount
		want *ApplicationCount
	}{
		{
			name: "success",
			ac:   applicationCount,
			want: &ApplicationCount{
				Counts: []ApplicationCountByLabelSet{
					{
						LabelSet: ApplicationCountLabelSet{
							Kind:   model.ApplicationKind_CLOUDRUN,
							Status: ApplicationStatusEnable,
						},
						Count: 1,
					},
					{
						LabelSet: ApplicationCountLabelSet{
							Kind:   model.ApplicationKind_KUBERNETES,
							Status: ApplicationStatusEnable,
						},
					},
					{
						LabelSet: ApplicationCountLabelSet{
							Kind:   model.ApplicationKind_KUBERNETES,
							Status: ApplicationStatusDisable,
						},
					},
					{
						LabelSet: ApplicationCountLabelSet{
							Kind:   model.ApplicationKind_KUBERNETES,
							Status: ApplicationStatusDeleted,
						},
					},
					{
						LabelSet: ApplicationCountLabelSet{
							Kind:   model.ApplicationKind_TERRAFORM,
							Status: ApplicationStatusEnable,
						},
					},
					{
						LabelSet: ApplicationCountLabelSet{
							Kind:   model.ApplicationKind_TERRAFORM,
							Status: ApplicationStatusDisable,
						},
					},
					{
						LabelSet: ApplicationCountLabelSet{
							Kind:   model.ApplicationKind_TERRAFORM,
							Status: ApplicationStatusDeleted,
						},
					},
					{
						LabelSet: ApplicationCountLabelSet{
							Kind:   model.ApplicationKind_CROSSPLANE,
							Status: ApplicationStatusEnable,
						},
					},
					{
						LabelSet: ApplicationCountLabelSet{
							Kind:   model.ApplicationKind_CROSSPLANE,
							Status: ApplicationStatusDisable,
						},
					},
					{
						LabelSet: ApplicationCountLabelSet{
							Kind:   model.ApplicationKind_CROSSPLANE,
							Status: ApplicationStatusDeleted,
						},
					},
					{
						LabelSet: ApplicationCountLabelSet{
							Kind:   model.ApplicationKind_LAMBDA,
							Status: ApplicationStatusEnable,
						},
					},
					{
						LabelSet: ApplicationCountLabelSet{
							Kind:   model.ApplicationKind_LAMBDA,
							Status: ApplicationStatusDisable,
						},
					},
					{
						LabelSet: ApplicationCountLabelSet{
							Kind:   model.ApplicationKind_LAMBDA,
							Status: ApplicationStatusDeleted,
						},
					},
					{
						LabelSet: ApplicationCountLabelSet{
							Kind:   model.ApplicationKind_CLOUDRUN,
							Status: ApplicationStatusDisable,
						},
					},
					{
						LabelSet: ApplicationCountLabelSet{
							Kind:   model.ApplicationKind_CLOUDRUN,
							Status: ApplicationStatusDeleted,
						},
					},
					{
						LabelSet: ApplicationCountLabelSet{
							Kind:   model.ApplicationKind_ECS,
							Status: ApplicationStatusEnable,
						},
					},
					{
						LabelSet: ApplicationCountLabelSet{
							Kind:   model.ApplicationKind_ECS,
							Status: ApplicationStatusDisable,
						},
					},
					{
						LabelSet: ApplicationCountLabelSet{
							Kind:   model.ApplicationKind_ECS,
							Status: ApplicationStatusDeleted,
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ac := tt.ac
			ac.MigrateApplicationCount()
			assert.Equal(t, tt.want, ac)
		})
	}
}

func TestApplicationCount_UpdateCount(t *testing.T) {
	applicationCount := func() *ApplicationCount {
		// init application count
		ac := NewApplicationCount()

		for i := 0; i < len(ac.Counts); i++ {
			c := &ac.Counts[i]
			enable := ApplicationCountLabelSet{
				Kind:   model.ApplicationKind_CLOUDRUN,
				Status: ApplicationStatusEnable,
			}
			if c.LabelSet == enable {
				c.Count = 1
			}
		}
		return ac
	}()
	tests := []struct {
		name    string
		ac      *ApplicationCount
		apps    []*model.Application
		want    *ApplicationCount
		wantErr bool
	}{
		{
			name: "success",
			ac:   applicationCount,
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
			want: func() *ApplicationCount {
				ac := NewApplicationCount()
				for i := 0; i < len(ac.Counts); i++ {
					c := &ac.Counts[i]
					enable := ApplicationCountLabelSet{
						Kind:   model.ApplicationKind_CLOUDRUN,
						Status: ApplicationStatusEnable,
					}
					delete := ApplicationCountLabelSet{
						Kind:   model.ApplicationKind_CLOUDRUN,
						Status: ApplicationStatusDeleted,
					}
					disable := ApplicationCountLabelSet{
						Kind:   model.ApplicationKind_CLOUDRUN,
						Status: ApplicationStatusDisable,
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

				return ac
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ac := tt.ac
			ac.UpdateCount(tt.apps)
			assert.Equal(t, tt.want, ac)
		})
	}
}
