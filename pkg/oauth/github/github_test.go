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

package github

import (
	"testing"

	"github.com/google/go-github/v29/github"
	"github.com/stretchr/testify/assert"

	"github.com/pipe-cd/pipecd/pkg/model"
)

func stringPointer(s string) *string { return &s }

func TestDecideRole(t *testing.T) {
	cases := []struct {
		name     string
		username string
		oc       *OAuthClient
		teams    []*github.Team
		role     *model.Role
		wantErr  bool
	}{
		{
			name:     "nothing",
			username: "foo",
			oc: &OAuthClient{
				project: &model.Project{
					Id:                 "id",
					AllowStrayAsViewer: false,
					UserGroups: []*model.ProjectUserGroup{
						{
							SsoGroup: "org/team-admin",
							Role:     "Admin",
						},
						{
							SsoGroup: "org/team-editor",
							Role:     "Editor",
						},
						{
							SsoGroup: "org/team-viewer",
							Role:     "Viewer",
						},
					},
				},
			},
			teams: []*github.Team{
				{
					Organization: &github.Organization{Login: stringPointer("org")},
					Slug:         stringPointer("team1"),
				},
			},
			wantErr: true,
		},
		{
			name:     "viewer as default",
			username: "foo",
			oc: &OAuthClient{
				project: &model.Project{
					Id:                 "id",
					AllowStrayAsViewer: true,
					UserGroups: []*model.ProjectUserGroup{
						{
							SsoGroup: "org/team-admin",
							Role:     "Admin",
						},
						{
							SsoGroup: "org/team-editor",
							Role:     "Editor",
						},
						{
							SsoGroup: "org/team-viewer",
							Role:     "Viewer",
						},
					},
				},
			},
			teams: []*github.Team{
				{
					Organization: &github.Organization{Login: stringPointer("org")},
					Slug:         stringPointer("team1"),
				},
			},
			role: &model.Role{
				ProjectId: "id",
				ProjectRbacRoles: []string{
					model.BuiltinRBACRoleViewer.String(),
				},
			},
			wantErr: false,
		},
		{
			name:     "admin",
			username: "foo",
			oc: &OAuthClient{
				project: &model.Project{
					Id:                 "id",
					AllowStrayAsViewer: false,
					UserGroups: []*model.ProjectUserGroup{
						{
							SsoGroup: "org/team-admin",
							Role:     "Admin",
						},
						{
							SsoGroup: "org/team-editor",
							Role:     "Editor",
						},
						{
							SsoGroup: "org/team-viewer",
							Role:     "Viewer",
						},
					},
				},
			},
			teams: []*github.Team{
				{
					Organization: &github.Organization{Login: stringPointer("org")},
					Slug:         stringPointer("team-admin"),
				},
				{
					Organization: &github.Organization{Login: stringPointer("org")},
					Slug:         stringPointer("team-editor"),
				},
				{
					Organization: &github.Organization{Login: stringPointer("org")},
					Slug:         stringPointer("team-viewer"),
				},
			},
			role: &model.Role{
				ProjectId: "id",
				ProjectRbacRoles: []string{
					model.BuiltinRBACRoleAdmin.String(),
					model.BuiltinRBACRoleEditor.String(),
					model.BuiltinRBACRoleViewer.String(),
				},
			},
			wantErr: false,
		},
		{
			name:     "editor",
			username: "foo",
			oc: &OAuthClient{
				project: &model.Project{
					Id:                 "id",
					AllowStrayAsViewer: false,
					UserGroups: []*model.ProjectUserGroup{
						{
							SsoGroup: "org/team-admin",
							Role:     "Admin",
						},
						{
							SsoGroup: "org/team-editor",
							Role:     "Editor",
						},
						{
							SsoGroup: "org/team-viewer",
							Role:     "Viewer",
						},
					},
				},
			},
			teams: []*github.Team{
				{
					Organization: &github.Organization{Login: stringPointer("org")},
					Slug:         stringPointer("team1"),
				},
				{
					Organization: &github.Organization{Login: stringPointer("org")},
					Slug:         stringPointer("team-editor"),
				},
				{
					Organization: &github.Organization{Login: stringPointer("org")},
					Slug:         stringPointer("team-viewer"),
				},
			},
			role: &model.Role{
				ProjectId: "id",
				ProjectRbacRoles: []string{
					model.BuiltinRBACRoleEditor.String(),
					model.BuiltinRBACRoleViewer.String(),
				},
			},
			wantErr: false,
		},
		{
			name:     "viewer",
			username: "foo",
			oc: &OAuthClient{
				project: &model.Project{
					Id:                 "id",
					AllowStrayAsViewer: false,
					UserGroups: []*model.ProjectUserGroup{
						{
							SsoGroup: "org/team-admin",
							Role:     "Admin",
						},
						{
							SsoGroup: "org/team-editor",
							Role:     "Editor",
						},
						{
							SsoGroup: "org/team-viewer",
							Role:     "Viewer",
						},
					},
				},
			},
			teams: []*github.Team{
				{
					Organization: &github.Organization{Login: stringPointer("org")},
					Slug:         stringPointer("team1"),
				},
				{
					Organization: &github.Organization{Login: stringPointer("org")},
					Slug:         stringPointer("team2"),
				},
				{
					Organization: &github.Organization{Login: stringPointer("org")},
					Slug:         stringPointer("team-viewer"),
				},
			},
			role: &model.Role{
				ProjectId: "id",
				ProjectRbacRoles: []string{
					model.BuiltinRBACRoleViewer.String(),
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			role, err := tc.oc.decideRole(tc.username, tc.teams)
			assert.Equal(t, tc.wantErr, err != nil)
			if err == nil {
				assert.Equal(t, tc.role, role)
			}
		})
	}
}
