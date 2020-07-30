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

package github

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/go-github/v29/github"
	"github.com/stretchr/testify/assert"

	"github.com/pipe-cd/pipe/pkg/model"
)

func stringPointer(s string) *string { return &s }

func TestGetUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cases := []struct {
		name    string
		oc      *OAuthClient
		user    *model.User
		wantErr bool
	}{
		{
			name: "nothing",
			oc: &OAuthClient{
				client: func() client {
					c := NewMockclient(ctrl)
					c.EXPECT().
						getUser(gomock.Any()).
						Return(&github.User{
							Login: stringPointer("foo"),
						}, nil, nil)
					c.EXPECT().
						getTeams(gomock.Any()).
						Return([]*github.Team{
							{
								Slug: stringPointer("team1"),
							},
						}, nil, nil)
					return c
				}(),
			},
			wantErr: true,
		},
		{
			name: "admin",
			oc: &OAuthClient{
				client: func() client {
					c := NewMockclient(ctrl)
					c.EXPECT().
						getUser(gomock.Any()).
						Return(&github.User{
							Login: stringPointer("foo"),
						}, nil, nil)
					c.EXPECT().
						getTeams(gomock.Any()).
						Return([]*github.Team{
							{
								Slug: stringPointer("team-admin"),
							},
							{
								Slug: stringPointer("team-editor"),
							},
							{
								Slug: stringPointer("team-viewer"),
							},
						}, nil, nil)
					return c
				}(),
			},
			user: &model.User{
				Username: "foo",
				Role: &model.Role{
					ProjectRole: model.Role_ADMIN,
				},
			},
		},
		{
			name: "editor",
			oc: &OAuthClient{
				client: func() client {
					c := NewMockclient(ctrl)
					c.EXPECT().
						getUser(gomock.Any()).
						Return(&github.User{
							Login: stringPointer("foo"),
						}, nil, nil)
					c.EXPECT().
						getTeams(gomock.Any()).
						Return([]*github.Team{
							{
								Slug: stringPointer("team1"),
							},
							{
								Slug: stringPointer("team-editor"),
							},
							{
								Slug: stringPointer("team-viewer"),
							},
						}, nil, nil)
					return c
				}(),
			},
			user: &model.User{
				Username: "foo",
				Role: &model.Role{
					ProjectRole: model.Role_EDITOR,
				},
			},
		},
		{
			name: "viewer",
			oc: &OAuthClient{
				client: func() client {
					c := NewMockclient(ctrl)
					c.EXPECT().
						getUser(gomock.Any()).
						Return(&github.User{
							Login: stringPointer("foo"),
						}, nil, nil)
					c.EXPECT().
						getTeams(gomock.Any()).
						Return([]*github.Team{
							{
								Slug: stringPointer("team1"),
							},
							{
								Slug: stringPointer("team2"),
							},
							{
								Slug: stringPointer("team-viewer"),
							},
						}, nil, nil)
					return c
				}(),
			},
			user: &model.User{
				Username: "foo",
				Role: &model.Role{
					ProjectRole: model.Role_VIEWER,
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.oc.adminTeam = "team-admin"
			tc.oc.editorTeam = "team-editor"
			tc.oc.viewerTeam = "team-viewer"
			user, err := tc.oc.GetUser(context.Background())
			assert.Equal(t, tc.wantErr, err != nil)
			assert.Equal(t, tc.user, user)
		})
	}
}
