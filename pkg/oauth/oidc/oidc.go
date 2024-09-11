// Copyright 2024 The PipeCD Authors.
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

package oidc

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"golang.org/x/oauth2"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/golang-jwt/jwt"
	"github.com/pipe-cd/pipecd/pkg/model"
)

var usernameClaimKeys = []string{"username", "preferred_username", "name", "cognito:username"}
var avatarUrlClaimKeys = []string{"picture", "avatar_url"}
var roleClaimKeys = []string{"groups", "roles", "cognito:groups", "custom:roles", "custom:groups"}

// OAuthClient is an oauth client for OIDC.
type OAuthClient struct {
	*oidc.Provider
	*oauth2.Token

	project *model.Project
}

// NewOAuthClient creates a new oauth client for OIDC.
func NewOAuthClient(ctx context.Context,
	sso *model.ProjectSSOConfig_Oidc,
	project *model.Project,
	code string,
) (*OAuthClient, error) {
	c := &OAuthClient{
		project: project,
	}

	provider, err := oidc.NewProvider(ctx, sso.Issuer)
	if err != nil {
		return nil, err
	}
	c.Provider = provider

	cfg := oauth2.Config{
		ClientID:     sso.ClientId,
		ClientSecret: sso.ClientSecret,
		RedirectURL:  sso.RedirectUri,
		Endpoint:     provider.Endpoint(),
		Scopes:       append(sso.Scopes, oidc.ScopeOpenID),
	}

	if sso.ProxyUrl != "" {
		proxyURL, err := url.Parse(sso.ProxyUrl)
		if err != nil {
			return nil, err
		}

		t := http.DefaultTransport.(*http.Transport).Clone()
		t.Proxy = http.ProxyURL(proxyURL)
		ctx = context.WithValue(ctx, oauth2.HTTPClient, &http.Client{Transport: t})
	}

	oauth2Token, err := cfg.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}
	c.Token = oauth2Token

	return c, nil
}

// GetUser returns a user model.
func (c *OAuthClient) GetUser(ctx context.Context, clientId string) (*model.User, error) {

	idTokenRAW, ok := c.Token.Extra("id_token").(string)
	if !ok {
		return nil, fmt.Errorf("no id_token in oauth2 token")
	}

	verifier := c.Provider.Verifier(&oidc.Config{ClientID: clientId})
	idToken, err := verifier.Verify(ctx, idTokenRAW)
	if err != nil {
		return nil, err
	}

	var claims jwt.MapClaims
	if err := idToken.Claims(&claims); err != nil {
		return nil, err
	}

	role, err := c.decideRole(claims)
	if err != nil {
		return nil, err
	}

	username, avatarUrl, err := c.decideUserInfos(claims)
	if err != nil {
		return nil, err
	}
	return &model.User{
		Username:  username,
		AvatarUrl: avatarUrl,
		Role:      role,
	}, nil
}

func (c *OAuthClient) decideRole(claims jwt.MapClaims) (role *model.Role, err error) {
	roleStrings := make([]string, 0)

	role = &model.Role{
		ProjectId:        c.project.Id,
		ProjectRbacRoles: roleStrings,
	}

	for _, key := range roleClaimKeys {
		val, ok := claims[key]
		if !ok || val == nil {
			continue
		}
		switch val := val.(type) {
		case []interface{}:
			for _, item := range val {
				if str, ok := item.(string); ok {
					roleStrings = append(roleStrings, str)
				}
			}
		case []string:
			roleStrings = append(roleStrings, val...)
		case string:
			if val != "" {
				roleStrings = append(roleStrings, val)
			}
		}
	}

	// Check if the current user belongs to any registered teams.
	for _, r := range roleStrings {
		if r == model.BuiltinRBACRoleAdmin.String() {
			role.ProjectRbacRoles = append(role.ProjectRbacRoles, model.BuiltinRBACRoleAdmin.String())
		}
		if r == model.BuiltinRBACRoleEditor.String() {
			role.ProjectRbacRoles = append(role.ProjectRbacRoles, model.BuiltinRBACRoleEditor.String())
		}
		if r == model.BuiltinRBACRoleViewer.String() {
			role.ProjectRbacRoles = append(role.ProjectRbacRoles, model.BuiltinRBACRoleViewer.String())
		}
	}

	// In case the current user does not have any role
	// if AllowStrayAsViewer option is set, assign Viewer role
	// as user's role.
	if c.project.AllowStrayAsViewer && len(roleStrings) == 0 {
		role.ProjectRbacRoles = []string{model.BuiltinRBACRoleViewer.String()}
		return
	}

	if len(roleStrings) == 0 {
		err = fmt.Errorf("no role found in claims")
		return
	}

	return
}

func (c *OAuthClient) decideUserInfos(claims jwt.MapClaims) (username, avatarUrl string, err error) {

	username = ""
	for _, key := range usernameClaimKeys {
		val, ok := claims[key]
		if ok && val != nil {
			if str, ok := val.(string); ok && str != "" {
				username = str
				break
			}
		}
	}

	if username == "" {
		err = fmt.Errorf("no username found in claims")
		return
	}

	avatarUrl = ""
	for _, key := range avatarUrlClaimKeys {
		val, ok := claims[key]
		if ok && val != nil {
			if str, ok := val.(string); ok && str != "" {
				avatarUrl = str
				break
			}
		}
	}

	return username, avatarUrl, nil
}
