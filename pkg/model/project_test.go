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

package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

type mockEncrypter struct {
}

func (e mockEncrypter) Encrypt(text string) (string, error) {
	return "encrypted-" + text, nil
}

type mockDecrypter struct {
}

func (d mockDecrypter) Decrypt(text string) (string, error) {
	return "decrypted-" + text, nil
}

func TestRedactSensitiveData(t *testing.T) {
	cases := []struct {
		name    string
		project *Project
		expect  *Project
	}{
		{
			name: "redact",
			project: &Project{
				StaticAdmin: &ProjectStaticUser{
					PasswordHash: "raw",
				},
				Sso: &ProjectSSOConfig{
					Github: &ProjectSSOConfig_GitHub{
						ClientId:     "raw",
						ClientSecret: "raw",
					},
				},
			},
			expect: &Project{
				StaticAdmin: &ProjectStaticUser{
					PasswordHash: "redacted",
				},
				Sso: &ProjectSSOConfig{
					Github: &ProjectSSOConfig_GitHub{
						ClientId:     "redacted",
						ClientSecret: "redacted",
					},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.project.RedactSensitiveData()
			assert.Equal(t, tc.expect, tc.project)
		})
	}
}

func TestUpdateProjectStaticUser(t *testing.T) {
	cases := []struct {
		name           string
		username       string
		password       string
		expectUsername string
		expectPassword string
	}{
		{
			name:           "update",
			username:       "foo",
			password:       "bar",
			expectUsername: "foo",
			expectPassword: "bar",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := &ProjectStaticUser{}
			p.Update(tc.username, tc.password)
			assert.Equal(t, tc.expectUsername, p.Username)
			err := bcrypt.CompareHashAndPassword([]byte(p.PasswordHash), []byte(tc.expectPassword))
			assert.Nil(t, err)
		})
	}
}

func TestUpdateProjectSSOConfig(t *testing.T) {
	cases := []struct {
		name   string
		sso    *ProjectSSOConfig
		expect *ProjectSSOConfig
	}{
		{
			name: "update",
			sso: &ProjectSSOConfig{
				Provider: ProjectSSOConfig_GITHUB,
				Github: &ProjectSSOConfig_GitHub{
					ClientId:     "updated-client-id",
					ClientSecret: "updated-client-secret",
					BaseUrl:      "updated-base-url",
					UploadUrl:    "updated-upload-url",
				},
				Google: nil,
			},
			expect: &ProjectSSOConfig{
				Provider: ProjectSSOConfig_GITHUB,
				Github: &ProjectSSOConfig_GitHub{
					ClientId:     "updated-client-id",
					ClientSecret: "updated-client-secret",
					BaseUrl:      "updated-base-url",
					UploadUrl:    "updated-upload-url",
				},
				Google: nil,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := &ProjectSSOConfig{}
			p.Update(tc.sso)
			assert.Equal(t, tc.expect, p)
		})
	}
}

func TestEncrypt(t *testing.T) {
	cases := []struct {
		name   string
		sso    *ProjectSSOConfig
		expect *ProjectSSOConfig
	}{
		{
			name: "encrypt",
			sso: &ProjectSSOConfig{
				Provider: ProjectSSOConfig_GITHUB,
				Github: &ProjectSSOConfig_GitHub{
					ClientId:     "client-id",
					ClientSecret: "client-secret",
					BaseUrl:      "base-url",
					UploadUrl:    "upload-url",
				},
				Google: nil,
			},
			expect: &ProjectSSOConfig{
				Provider: ProjectSSOConfig_GITHUB,
				Github: &ProjectSSOConfig_GitHub{
					ClientId:     "encrypted-client-id",
					ClientSecret: "encrypted-client-secret",
					BaseUrl:      "base-url",
					UploadUrl:    "upload-url",
				},
				Google: nil,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			e := &mockEncrypter{}
			err := tc.sso.Encrypt(e)
			assert.NoError(t, err)
			assert.Equal(t, tc.expect, tc.sso)
		})
	}
}

func TestDecrypt(t *testing.T) {
	cases := []struct {
		name   string
		sso    *ProjectSSOConfig
		expect *ProjectSSOConfig
	}{
		{
			name: "decrypt",
			sso: &ProjectSSOConfig{
				Provider: ProjectSSOConfig_GITHUB,
				Github: &ProjectSSOConfig_GitHub{
					ClientId:     "client-id",
					ClientSecret: "client-secret",
					BaseUrl:      "base-url",
					UploadUrl:    "upload-url",
				},
				Google: nil,
			},
			expect: &ProjectSSOConfig{
				Provider: ProjectSSOConfig_GITHUB,
				Github: &ProjectSSOConfig_GitHub{
					ClientId:     "decrypted-client-id",
					ClientSecret: "decrypted-client-secret",
					BaseUrl:      "base-url",
					UploadUrl:    "upload-url",
				},
				Google: nil,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			d := &mockDecrypter{}
			err := tc.sso.Decrypt(d)
			assert.NoError(t, err)
			assert.Equal(t, tc.expect, tc.sso)
		})
	}
}

func TestProject_HasRBACRole(t *testing.T) {
	roles := []*ProjectRBACRole{
		{
			Name: "test",
			Policies: []*ProjectRBACPolicy{
				{
					Resources: []*ProjectRBACResource{
						{
							Type: ProjectRBACResource_APPLICATION,
						},
					},
					Actions: []ProjectRBACPolicy_Action{
						ProjectRBACPolicy_GET,
					},
				},
			},
		},
	}
	p := &Project{RbacRoles: roles}

	// True
	assert.True(t, p.HasRBACRole("test"))

	// False
	assert.False(t, p.HasRBACRole("foo"))
}

func TestProject_HasUserGroup(t *testing.T) {
	groups := []*ProjectUserGroup{
		{
			SsoGroup: "team/tester",
			Role:     "Tester",
		},
	}
	rbac := &ProjectRBACConfig{
		Admin:  "team/admin",
		Editor: "team/editor",
		Viewer: "team/viewer",
	}
	p := &Project{
		UserGroups: groups,
		Rbac:       rbac,
	}

	// True
	assert.True(t, p.HasUserGroup("team/tester"))
	assert.True(t, p.HasUserGroup("team/admin"))

	// False
	assert.False(t, p.HasUserGroup("team/foo"))
}

func TestProject_SetLegacyUserGroups(t *testing.T) {
	testcases := []struct {
		name    string
		project *Project
		want    *Project
	}{
		{
			name:    "empty",
			project: &Project{},
			want:    &Project{},
		},
		{
			name: "merge rbac config and user group",
			project: &Project{
				Rbac: &ProjectRBACConfig{
					Admin:  "team/admin",
					Editor: "team/editor",
					Viewer: "team/viewer",
				},
				UserGroups: []*ProjectUserGroup{
					{
						Role:     "Tester",
						SsoGroup: "team/tester",
					},
					{
						Role:     "Owner",
						SsoGroup: "team/owner",
					},
				},
			},
			want: &Project{
				Rbac: &ProjectRBACConfig{
					Admin:  "team/admin",
					Editor: "team/editor",
					Viewer: "team/viewer",
				},
				UserGroups: []*ProjectUserGroup{
					{
						Role:     "Admin",
						SsoGroup: "team/admin",
					},
					{
						Role:     "Editor",
						SsoGroup: "team/editor",
					},
					{
						Role:     "Viewer",
						SsoGroup: "team/viewer",
					},
					{
						Role:     "Tester",
						SsoGroup: "team/tester",
					},
					{
						Role:     "Owner",
						SsoGroup: "team/owner",
					},
				},
			},
		},
		{
			name: "exists same name rbac config",
			project: &Project{
				Rbac: &ProjectRBACConfig{
					Admin:  "team/admin",
					Editor: "team/admin",
					Viewer: "team/admin",
				},
				UserGroups: []*ProjectUserGroup{
					{
						Role:     "Tester",
						SsoGroup: "team/tester",
					},
					{
						Role:     "Owner",
						SsoGroup: "team/owner",
					},
				},
			},
			want: &Project{
				Rbac: &ProjectRBACConfig{
					Admin:  "team/admin",
					Editor: "team/admin",
					Viewer: "team/admin",
				},
				UserGroups: []*ProjectUserGroup{
					{
						Role:     "Admin",
						SsoGroup: "team/admin",
					},
					{
						Role:     "Tester",
						SsoGroup: "team/tester",
					},
					{
						Role:     "Owner",
						SsoGroup: "team/owner",
					},
				},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			tc.project.SetLegacyUserGroups()
			assert.Equal(t, tc.want, tc.project)
		})
	}
}

func TestProject_AddUserGroup(t *testing.T) {
	type args struct {
		sso  string
		role string
	}
	testcases := []struct {
		name    string
		args    args
		project *Project
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				sso:  "team/admin",
				role: "Admin",
			},
			project: &Project{
				RbacRoles: []*ProjectRBACRole{
					builtinAdminRBACRole,
				},
			},
			wantErr: false,
		},
		{
			name: "the role is already assigned in rbac config",
			args: args{
				sso:  "team/admin",
				role: "Admin",
			},
			project: &Project{
				Rbac: &ProjectRBACConfig{
					Admin: "team/admin",
				},
				RbacRoles: []*ProjectRBACRole{
					builtinAdminRBACRole,
				},
			},
			wantErr: true,
		},
		{
			name: "the role is already assigned in user group",
			args: args{
				sso:  "team/tester",
				role: "Tester",
			},
			project: &Project{
				UserGroups: []*ProjectUserGroup{
					{
						SsoGroup: "team/tester",
						Role:     "Tester",
					},
				},
				RbacRoles: []*ProjectRBACRole{
					{
						Name: "Tester",
						Policies: []*ProjectRBACPolicy{
							{
								Resources: []*ProjectRBACResource{
									{
										Type: ProjectRBACResource_APPLICATION,
									},
								},
								Actions: []ProjectRBACPolicy_Action{
									ProjectRBACPolicy_GET,
								},
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "the role doesn't exist",
			args: args{
				sso:  "team/tester",
				role: "Tester",
			},
			project: &Project{},
			wantErr: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.project.AddUserGroup(tc.args.sso, tc.args.role)
			assert.Equal(t, tc.wantErr, err != nil)
		})
	}
}

func TestProject_DeleteUserGroup(t *testing.T) {
	type args struct {
		sso string
	}
	testcases := []struct {
		name    string
		args    args
		project *Project
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				sso: "team/admin",
			},
			project: &Project{
				UserGroups: []*ProjectUserGroup{
					{
						SsoGroup: "team/admin",
						Role:     "Admin",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "the user group doen't exist",
			args: args{
				sso: "team/admin",
			},
			project: &Project{},
			wantErr: true,
		},
		{
			name: "delete rbac config role",
			args: args{
				sso: "team/admin",
			},
			project: &Project{
				Rbac: &ProjectRBACConfig{
					Admin: "team/admin",
				},
			},
			wantErr: false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.project.DeleteUserGroup(tc.args.sso)
			assert.Equal(t, tc.wantErr, err != nil)
		})
	}
}

func TestProject_SetBuiltinRBACRoles(t *testing.T) {
	testcases := []struct {
		name    string
		project *Project
		want    *Project
	}{
		{
			name: "ok",
			project: &Project{
				RbacRoles: []*ProjectRBACRole{
					{
						Name: "Tester",
						Policies: []*ProjectRBACPolicy{
							{
								Resources: []*ProjectRBACResource{
									{
										Type: ProjectRBACResource_APPLICATION,
									},
								},
								Actions: []ProjectRBACPolicy_Action{
									ProjectRBACPolicy_GET,
								},
							},
						},
					},
				},
			},
			want: &Project{
				RbacRoles: []*ProjectRBACRole{
					builtinAdminRBACRole,
					builtinEditorRBACRole,
					builtinViewerRBACRole,
					{
						Name: "Tester",
						Policies: []*ProjectRBACPolicy{
							{
								Resources: []*ProjectRBACResource{
									{
										Type: ProjectRBACResource_APPLICATION,
									},
								},
								Actions: []ProjectRBACPolicy_Action{
									ProjectRBACPolicy_GET,
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			tc.project.SetBuiltinRBACRoles()
			assert.Equal(t, tc.want, tc.project)
		})
	}
}

func TestProject_AddRBACRole(t *testing.T) {
	type args struct {
		name     string
		policies []*ProjectRBACPolicy
	}
	testcases := []struct {
		name    string
		args    args
		project *Project
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				name: "Tester",
				policies: []*ProjectRBACPolicy{
					{
						Resources: []*ProjectRBACResource{
							{
								Type: ProjectRBACResource_APPLICATION,
							},
						},
						Actions: []ProjectRBACPolicy_Action{
							ProjectRBACPolicy_GET,
						},
					},
				},
			},
			project: &Project{},
			wantErr: false,
		},
		{
			name: "the name is already used",
			args: args{
				name: "Tester",
				policies: []*ProjectRBACPolicy{
					{
						Resources: []*ProjectRBACResource{
							{
								Type: ProjectRBACResource_APPLICATION,
							},
						},
						Actions: []ProjectRBACPolicy_Action{
							ProjectRBACPolicy_GET,
						},
					},
				},
			},
			project: &Project{
				RbacRoles: []*ProjectRBACRole{
					{
						Name: "Tester",
						Policies: []*ProjectRBACPolicy{
							{
								Resources: []*ProjectRBACResource{
									{
										Type: ProjectRBACResource_APPLICATION,
									},
								},
								Actions: []ProjectRBACPolicy_Action{
									ProjectRBACPolicy_GET,
								},
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "the name of built-in role cannot be used",
			args: args{
				name: "Admin",
				policies: []*ProjectRBACPolicy{
					{
						Resources: []*ProjectRBACResource{
							{
								Type: ProjectRBACResource_APPLICATION,
							},
						},
						Actions: []ProjectRBACPolicy_Action{
							ProjectRBACPolicy_GET,
						},
					},
				},
			},
			project: &Project{},
			wantErr: true,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.project.AddRBACRole(tc.args.name, tc.args.policies)
			assert.Equal(t, tc.wantErr, got != nil)
		})
	}
}

func TestProject_UpdateRBACRole(t *testing.T) {
	type args struct {
		name     string
		policies []*ProjectRBACPolicy
	}
	testcases := []struct {
		name    string
		args    args
		project *Project
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				name: "Tester",
				policies: []*ProjectRBACPolicy{
					{
						Resources: []*ProjectRBACResource{
							{
								Type: ProjectRBACResource_APPLICATION,
							},
						},
						Actions: []ProjectRBACPolicy_Action{
							ProjectRBACPolicy_GET,
						},
					},
				},
			},
			project: &Project{
				RbacRoles: []*ProjectRBACRole{
					{
						Name: "Tester",
						Policies: []*ProjectRBACPolicy{
							{
								Resources: []*ProjectRBACResource{
									{
										Type: ProjectRBACResource_APPLICATION,
									},
								},
								Actions: []ProjectRBACPolicy_Action{
									ProjectRBACPolicy_ALL,
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "the role doesn't exist",
			args: args{
				name: "Tester",
				policies: []*ProjectRBACPolicy{
					{
						Resources: []*ProjectRBACResource{
							{
								Type: ProjectRBACResource_APPLICATION,
							},
						},
						Actions: []ProjectRBACPolicy_Action{
							ProjectRBACPolicy_GET,
						},
					},
				},
			},
			project: &Project{},
			wantErr: true,
		},
		{
			name: "built-in role cannot be updated",
			args: args{
				name: "Admin",
				policies: []*ProjectRBACPolicy{
					{
						Resources: []*ProjectRBACResource{
							{
								Type: ProjectRBACResource_APPLICATION,
							},
						},
						Actions: []ProjectRBACPolicy_Action{
							ProjectRBACPolicy_GET,
						},
					},
				},
			},
			project: &Project{},
			wantErr: true,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.project.UpdateRBACRole(tc.args.name, tc.args.policies)
			assert.Equal(t, tc.wantErr, got != nil)
		})
	}
}

func TestProject_DeleteRBACRole(t *testing.T) {
	type args struct {
		name string
	}
	testcases := []struct {
		name    string
		args    args
		project *Project
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				name: "Tester",
			},
			project: &Project{
				RbacRoles: []*ProjectRBACRole{
					{
						Name: "Tester",
						Policies: []*ProjectRBACPolicy{
							{
								Resources: []*ProjectRBACResource{
									{
										Type: ProjectRBACResource_APPLICATION,
									},
								},
								Actions: []ProjectRBACPolicy_Action{
									ProjectRBACPolicy_GET,
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "the role doesn't exist",
			args: args{
				name: "Tester",
			},
			project: &Project{},
			wantErr: true,
		},
		{
			name: "built-in role cannot be deleted",
			args: args{
				name: "Tester",
			},
			project: &Project{},
			wantErr: true,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.project.DeleteRBACRole(tc.args.name)
			assert.Equal(t, tc.wantErr, got != nil)
		})
	}
}

func TestProjectRBACRoles_Authorize(t *testing.T) {
	type args struct {
		rsc    ProjectRBACResource_ResourceType
		labels map[string]string
	}
	testcases := []struct {
		name  string
		args  args
		roles ProjectRBACRoles
		want  bool
	}{
		{
			name: "no user roles",
			args: args{
				rsc:    ProjectRBACResource_APPLICATION,
				labels: map[string]string{"env": "test"},
			},
			roles: ProjectRBACRoles{},
			want:  false,
		},
		{
			name: "ok",
			args: args{
				rsc:    ProjectRBACResource_APPLICATION,
				labels: map[string]string{},
			},
			roles: ProjectRBACRoles{
				builtinAdminRBACRole,
			},
			want: true,
		},
		{
			name: "necessary labels do not exist",
			args: args{
				rsc:    ProjectRBACResource_APPLICATION,
				labels: map[string]string{},
			},
			roles: ProjectRBACRoles{
				{
					Name: "Tester",
					Policies: []*ProjectRBACPolicy{
						{
							Resources: []*ProjectRBACResource{
								{
									Type:   ProjectRBACResource_ALL,
									Labels: map[string]string{"env": "test"},
								},
							},
							Actions: []ProjectRBACPolicy_Action{
								ProjectRBACPolicy_ALL,
							},
						},
					},
				},
			},
			want: false,
		},
		{
			name: "necessary labels exist",
			args: args{
				rsc:    ProjectRBACResource_APPLICATION,
				labels: map[string]string{"env": "test", "piped": "test"},
			},
			roles: ProjectRBACRoles{
				{
					Name: "Tester",
					Policies: []*ProjectRBACPolicy{
						{
							Resources: []*ProjectRBACResource{
								{
									Type:   ProjectRBACResource_APPLICATION,
									Labels: map[string]string{"env": "test", "piped": "test"},
								},
							},
							Actions: []ProjectRBACPolicy_Action{
								ProjectRBACPolicy_ALL,
							},
						},
					},
				},
			},
			want: true,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.roles.Authorize(tc.args.rsc, tc.args.labels)
			assert.Equal(t, tc.want, got)
		})
	}
}
