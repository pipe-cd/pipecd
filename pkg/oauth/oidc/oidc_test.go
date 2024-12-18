// Copyright 2024 The PipeCD Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package oidc

import (
	"fmt"
	"testing"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"

	"github.com/pipe-cd/pipecd/pkg/model"
)

func TestDecideRole(t *testing.T) {
	cases := []struct {
		claims   jwt.MapClaims
		oc       *OAuthClient
		expected *model.Role
		err      error
	}{
		{
			claims: jwt.MapClaims{
				"groups": []interface{}{model.BuiltinRBACRoleAdmin.String(), model.BuiltinRBACRoleEditor.String()},
			},
			oc: &OAuthClient{
				project: &model.Project{
					Id:                 "project-id",
					AllowStrayAsViewer: false,
					UserGroups:         nil,
				},
			},
			expected: &model.Role{
				ProjectId:        "project-id",
				ProjectRbacRoles: []string{model.BuiltinRBACRoleAdmin.String(), model.BuiltinRBACRoleEditor.String()},
			},
			err: nil,
		},
		{
			claims: jwt.MapClaims{
				"roles": []interface{}{model.BuiltinRBACRoleEditor.String()},
			},
			oc: &OAuthClient{
				project: &model.Project{
					Id:                 "project-id",
					AllowStrayAsViewer: false,
					UserGroups:         nil,
				},
			},
			expected: &model.Role{
				ProjectId:        "project-id",
				ProjectRbacRoles: []string{model.BuiltinRBACRoleEditor.String()},
			},
			err: nil,
		},
		{
			claims: jwt.MapClaims{
				"custom:groups": []interface{}{model.BuiltinRBACRoleViewer.String()},
			},
			oc: &OAuthClient{
				project: &model.Project{
					Id:                 "project-id",
					AllowStrayAsViewer: false,
					UserGroups:         nil,
				},
			},
			expected: &model.Role{
				ProjectId:        "project-id",
				ProjectRbacRoles: []string{model.BuiltinRBACRoleViewer.String()},
			},
			err: nil,
		},
		{
			claims: jwt.MapClaims{},
			oc: &OAuthClient{
				project: &model.Project{
					Id:                 "project-id",
					AllowStrayAsViewer: true,
					UserGroups:         nil,
				},
			},
			expected: &model.Role{
				ProjectId:        "project-id",
				ProjectRbacRoles: []string{model.BuiltinRBACRoleViewer.String()},
			},
			err: nil,
		},
		{
			claims: jwt.MapClaims{},
			oc: &OAuthClient{
				project: &model.Project{
					Id:                 "project-id",
					AllowStrayAsViewer: false,
					UserGroups:         nil,
				},
			},
			expected: nil,
			err:      fmt.Errorf("no role found in claims"),
		},
	}

	for _, c := range cases {
		role, err := c.oc.decideRole(c.claims)
		if c.err != nil {
			assert.Error(t, err)
			assert.Equal(t, c.err.Error(), err.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, c.expected, role)
		}
	}
}

func TestDecideUserInfos(t *testing.T) {
	client := &OAuthClient{}

	cases := []struct {
		claims         jwt.MapClaims
		expectedUser   string
		expectedAvatar string
		err            error
	}{
		{
			claims: jwt.MapClaims{
				"username": "john_doe",
			},
			expectedUser:   "john_doe",
			expectedAvatar: "",
			err:            nil,
		},
		{
			claims: jwt.MapClaims{
				"name": "John Doe",
			},
			expectedUser:   "John Doe",
			expectedAvatar: "",
			err:            nil,
		},
		{
			claims: jwt.MapClaims{
				"preferred_username": "johnny",
			},
			expectedUser:   "johnny",
			expectedAvatar: "",
			err:            nil,
		},
		{
			claims: jwt.MapClaims{
				"cognito:username": "john_cognito",
			},
			expectedUser:   "john_cognito",
			expectedAvatar: "",
			err:            nil,
		},
		{
			claims: jwt.MapClaims{
				"avatar_url": "http://example.com/avatar.jpg",
			},
			expectedUser:   "",
			expectedAvatar: "http://example.com/avatar.jpg",
			err:            fmt.Errorf("no username found in claims"),
		},
		{
			claims: jwt.MapClaims{
				"username":   "john_doe",
				"avatar_url": "http://example.com/avatar.jpg",
			},
			expectedUser:   "john_doe",
			expectedAvatar: "http://example.com/avatar.jpg",
			err:            nil,
		},
		{
			claims:         jwt.MapClaims{},
			expectedUser:   "",
			expectedAvatar: "",
			err:            fmt.Errorf("no username found in claims"),
		},
	}

	for _, c := range cases {
		username, avatarUrl, err := client.decideUserInfos(c.claims)
		if c.err != nil {
			assert.Error(t, err)
			assert.Equal(t, c.err.Error(), err.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, c.expectedUser, username)
			assert.Equal(t, c.expectedAvatar, avatarUrl)
		}
	}
}
