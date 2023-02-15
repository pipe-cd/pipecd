// Copyright 2023 The PipeCD Authors.
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

	"github.com/pipe-cd/pipecd/pkg/model"
)

func TestAddProject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testcases := []struct {
		name      string
		project   *model.Project
		dsFactory func(*model.Project) DataStore
		wantErr   bool
	}{
		{
			name:      "Invalid project",
			project:   &model.Project{},
			dsFactory: func(d *model.Project) DataStore { return nil },
			wantErr:   true,
		},
		{
			name: "Valid project",
			project: &model.Project{
				Id: "id",
				StaticAdmin: &model.ProjectStaticUser{
					Username:     "username",
					PasswordHash: "password-hash",
				},
				CreatedAt: 1,
				UpdatedAt: 1,
			},
			dsFactory: func(d *model.Project) DataStore {
				ds := NewMockDataStore(ctrl)
				ds.EXPECT().Create(gomock.Any(), gomock.Any(), d.Id, d)
				return ds
			},
			wantErr: false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			s := NewProjectStore(tc.dsFactory(tc.project), TestCommander)
			err := s.Add(context.Background(), tc.project)
			assert.Equal(t, tc.wantErr, err != nil)
		})
	}
}

func TestGetProject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testcases := []struct {
		name    string
		id      string
		ds      DataStore
		wantErr bool
	}{
		{
			name: "successful fetch from datastore",
			id:   "id",
			ds: func() DataStore {
				ds := NewMockDataStore(ctrl)
				ds.EXPECT().
					Get(gomock.Any(), gomock.Any(), "id", &model.Project{}).
					Return(nil)
				return ds
			}(),
			wantErr: false,
		},
		{
			name: "failed fetch from datastore",
			id:   "id",
			ds: func() DataStore {
				ds := NewMockDataStore(ctrl)
				ds.EXPECT().
					Get(gomock.Any(), gomock.Any(), "id", &model.Project{}).
					Return(fmt.Errorf("err"))
				return ds
			}(),
			wantErr: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			s := NewProjectStore(tc.ds, TestCommander)
			_, err := s.Get(context.Background(), tc.id)
			assert.Equal(t, tc.wantErr, err != nil)
		})
	}
}

func TestListProjects(t *testing.T) {
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
			opts: ListOptions{},
			ds: func() DataStore {
				it := NewMockIterator(ctrl)
				it.EXPECT().
					Next(&model.Project{}).
					Return(ErrIteratorDone)

				ds := NewMockDataStore(ctrl)
				ds.EXPECT().
					Find(gomock.Any(), gomock.Any(), ListOptions{}).
					Return(it, nil)
				return ds
			}(),
			wantErr: false,
		},
		{
			name: "unexpected error occurred",
			opts: ListOptions{},
			ds: func() DataStore {
				it := NewMockIterator(ctrl)
				it.EXPECT().
					Next(&model.Project{}).
					Return(fmt.Errorf("err"))

				ds := NewMockDataStore(ctrl)
				ds.EXPECT().
					Find(gomock.Any(), gomock.Any(), ListOptions{}).
					Return(it, nil)
				return ds
			}(),
			wantErr: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			s := NewProjectStore(tc.ds, TestCommander)
			_, err := s.List(context.Background(), tc.opts)
			assert.Equal(t, tc.wantErr, err != nil)
		})
	}
}

func TestValidateProject(t *testing.T) {
	testcases := []struct {
		name    string
		project *model.Project
		wantErr bool
	}{
		{
			name: "ok",
			project: &model.Project{
				Id: "id",
				StaticAdmin: &model.ProjectStaticUser{
					Username:     "username",
					PasswordHash: "password-hash",
				},
				CreatedAt: 1,
				UpdatedAt: 1,
			},
			wantErr: false,
		},
		{
			name: "policies validation error",
			project: &model.Project{
				Id: "id",
				StaticAdmin: &model.ProjectStaticUser{
					Username:     "username",
					PasswordHash: "password-hash",
				},
				RbacRoles: []*model.ProjectRBACRole{
					{
						Name: "Tester",
					},
				},
				CreatedAt: 1,
				UpdatedAt: 1,
			},
			wantErr: true,
		},
		{
			name: "resources validation error",
			project: &model.Project{
				Id: "id",
				StaticAdmin: &model.ProjectStaticUser{
					Username:     "username",
					PasswordHash: "password-hash",
				},
				RbacRoles: []*model.ProjectRBACRole{
					{
						Name: "Tester",
						Policies: []*model.ProjectRBACPolicy{
							{
								Actions: []model.ProjectRBACPolicy_Action{
									model.ProjectRBACPolicy_ALL,
								},
							},
						},
					},
				},
				CreatedAt: 1,
				UpdatedAt: 1,
			},
			wantErr: true,
		},
		{
			name: "resource type validation error",
			project: &model.Project{
				Id: "id",
				StaticAdmin: &model.ProjectStaticUser{
					Username:     "username",
					PasswordHash: "password-hash",
				},
				RbacRoles: []*model.ProjectRBACRole{
					{
						Name: "Tester",
						Policies: []*model.ProjectRBACPolicy{
							{
								Resources: []*model.ProjectRBACResource{
									{
										Type:   99,
										Labels: map[string]string{"key": "value"},
									},
								},
								Actions: []model.ProjectRBACPolicy_Action{
									model.ProjectRBACPolicy_ALL,
								},
							},
						},
					},
				},
				CreatedAt: 1,
				UpdatedAt: 1,
			},
			wantErr: true,
		},
		{
			name: "actions validation error",
			project: &model.Project{
				Id: "id",
				StaticAdmin: &model.ProjectStaticUser{
					Username:     "username",
					PasswordHash: "password-hash",
				},
				RbacRoles: []*model.ProjectRBACRole{
					{
						Name: "Tester",
						Policies: []*model.ProjectRBACPolicy{
							{
								Resources: []*model.ProjectRBACResource{
									{
										Type: model.ProjectRBACResource_ALL,
									},
								},
							},
						},
					},
				},
				CreatedAt: 1,
				UpdatedAt: 1,
			},
			wantErr: true,
		},
		{
			name: "action type validation error",
			project: &model.Project{
				Id: "id",
				StaticAdmin: &model.ProjectStaticUser{
					Username:     "username",
					PasswordHash: "password-hash",
				},
				RbacRoles: []*model.ProjectRBACRole{
					{
						Name: "Tester",
						Policies: []*model.ProjectRBACPolicy{
							{
								Resources: []*model.ProjectRBACResource{
									{
										Type: model.ProjectRBACResource_ALL,
									},
								},
								Actions: []model.ProjectRBACPolicy_Action{99, 100},
							},
						},
					},
				},
				CreatedAt: 1,
				UpdatedAt: 1,
			},
			wantErr: true,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.project.Validate()
			assert.Equal(t, tc.wantErr, err != nil)
		})
	}
}
