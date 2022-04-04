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

func TestProject_FilterRBACRoles(t *testing.T) {
	rbac := []*ProjectRBACRole{
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
	}
	p := &Project{RbacRoles: rbac}

	isBuiltin := true
	got := p.FilterRBACRoles(isBuiltin)
	assert.Len(t, got, 3)

	isBuiltin = false
	got = p.FilterRBACRoles(isBuiltin)
	assert.Len(t, got, 1)
}

func TestValidateRBACRoles(t *testing.T) {
	testcases := []struct {
		name    string
		roles   []*ProjectRBACRole
		wantErr bool
	}{
		{
			name: "cannot use built-in role name",
			roles: []*ProjectRBACRole{
				{
					Name: "Admin",
					Policies: []*ProjectRBACPolicy{
						{
							Resources: []*ProjectRBACResource{
								{
									Type: ProjectRBACResource_ALL,
								},
							},
							Actions: []ProjectRBACPolicy_Action{
								ProjectRBACPolicy_ALL,
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "role name must be unique",
			roles: []*ProjectRBACRole{
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
			wantErr: true,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateRBACRoles(tc.roles)
			assert.Equal(t, tc.wantErr, err != nil)
		})
	}
}

func TestProjectRBACRole_IsBuiltinName(t *testing.T) {
	p := &ProjectRBACRole{}
	// Admin
	p.Name = "Admin"
	assert.True(t, p.IsBuiltinName())
	// Editor
	p.Name = "Editor"
	assert.True(t, p.IsBuiltinName())
	// Viewer
	p.Name = "Viewer"
	assert.True(t, p.IsBuiltinName())
	// Other
	p.Name = "Tester"
	assert.False(t, p.IsBuiltinName())
}

func TestValidateUserGroups(t *testing.T) {
	groups := []*ProjectUserGroup{
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
	}
	err := ValidateUserGroups(groups)
	assert.NoError(t, err)

	groups = []*ProjectUserGroup{
		{
			Role:     "Tester",
			SsoGroup: "team/tester",
		},
		{
			Role:     "Owner",
			SsoGroup: "team/tester",
		},
	}
	err = ValidateUserGroups(groups)
	assert.Error(t, err)
}
