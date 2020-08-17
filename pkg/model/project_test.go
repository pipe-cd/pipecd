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
				Sso: &ProjectSingleSignOn{
					Github: &ProjectSingleSignOn_GitHub{
						ClientId:     "raw",
						ClientSecret: "raw",
					},
				},
			},
			expect: &Project{
				StaticAdmin: &ProjectStaticUser{
					PasswordHash: "redacted",
				},
				Sso: &ProjectSingleSignOn{
					Github: &ProjectSingleSignOn_GitHub{
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

func TestUpdateProjectSingleSignOn(t *testing.T) {
	cases := []struct {
		name   string
		sso    *ProjectSingleSignOn
		expect *ProjectSingleSignOn
	}{
		{
			name: "update",
			sso: &ProjectSingleSignOn{
				Provider: ProjectSingleSignOnProvider_GITHUB,
				Github: &ProjectSingleSignOn_GitHub{
					ClientId:     "updated",
					ClientSecret: "updated",
					BaseUrl:      "updated",
					UploadUrl:    "updated",
					Org:          "updated",
					AdminTeam:    "updated",
					EditorTeam:   "updated",
					ViewerTeam:   "updated",
				},
				Google: nil,
			},
			expect: &ProjectSingleSignOn{
				Provider: ProjectSingleSignOnProvider_GITHUB,
				Github: &ProjectSingleSignOn_GitHub{
					ClientId:     "updated",
					ClientSecret: "updated",
					BaseUrl:      "updated",
					UploadUrl:    "updated",
					Org:          "updated",
					AdminTeam:    "updated",
					EditorTeam:   "updated",
					ViewerTeam:   "updated",
				},
				Google: nil,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := &ProjectSingleSignOn{}
			p.Update(tc.sso)
			assert.Equal(t, tc.expect, p)
		})
	}
}
