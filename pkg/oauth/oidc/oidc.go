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
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"

	"github.com/pipe-cd/pipecd/pkg/model"
)

var defaultUsernameClaimKeys = []string{"username", "preferred_username", "name", "cognito:username"}
var defaultAvatarURLClaimKeys = []string{"picture", "avatar_url"}
var defaultRoleClaimKeys = []string{"groups", "roles", "cognito:groups", "custom:roles", "custom:groups"}

// OAuthClient is an oauth client for OIDC.
type OAuthClient struct {
	*oidc.Provider
	*oauth2.Token

	sharedSSOConfig *model.ProjectSSOConfig_Oidc
	project         *model.Project
}

// NewOAuthClient creates a new oauth client for OIDC.
func NewOAuthClient(ctx context.Context,
	sso *model.ProjectSSOConfig_Oidc,
	project *model.Project,
	code string,
) (*OAuthClient, error) {
	c := &OAuthClient{
		project:         project,
		sharedSSOConfig: sso,
	}

	if sso.AuthorizationEndpoint != "" || sso.TokenEndpoint != "" || sso.UserInfoEndpoint != "" {
		provider, err := createCustomOIDCProvider(ctx, sso)
		if err != nil {
			return nil, err
		}
		c.Provider = provider
	} else {
		provider, err := oidc.NewProvider(ctx, sso.Issuer)
		if err != nil {
			return nil, err
		}
		c.Provider = provider
	}

	cfg := oauth2.Config{
		ClientID:     sso.ClientId,
		ClientSecret: sso.ClientSecret,
		RedirectURL:  sso.RedirectUri,
		Endpoint:     c.Endpoint(),
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
func (c *OAuthClient) GetUser(ctx context.Context) (*model.User, error) {

	idTokenRAW, ok := c.Token.Extra("id_token").(string)
	if !ok {
		return nil, fmt.Errorf("no id_token in oauth2 token")
	}

	verifier := c.Verifier(&oidc.Config{ClientID: c.sharedSSOConfig.ClientId})
	idToken, err := verifier.Verify(ctx, idTokenRAW)
	if err != nil {
		return nil, err
	}

	var claims jwt.MapClaims
	if err := idToken.Claims(&claims); err != nil {
		return nil, err
	}

	if c.UserInfoEndpoint() != "" {
		userInfo, err := c.UserInfo(ctx, oauth2.StaticTokenSource(c.Token))
		if err != nil {
			return nil, err
		}

		var userInfoClaims map[string]interface{}
		if err := userInfo.Claims(&userInfoClaims); err != nil {
			return nil, err
		}

		for k, v := range userInfoClaims {
			claims[k] = v
		}
	}

	role, err := c.decideRole(claims, c.sharedSSOConfig.RolesClaimKey)
	if err != nil {
		return nil, err
	}

	username, avatarURL, err := c.decideUserInfos(claims, c.sharedSSOConfig.UsernameClaimKey, c.sharedSSOConfig.AvatarUrlClaimKey)
	if err != nil {
		return nil, err
	}
	return &model.User{
		Username:  username,
		AvatarUrl: avatarURL,
		Role:      role,
	}, nil
}

func (c *OAuthClient) decideRole(claims jwt.MapClaims, roleClaimKey string) (role *model.Role, err error) {
	roleStrings := make([]string, 0)

	role = &model.Role{
		ProjectId:        c.project.Id,
		ProjectRbacRoles: roleStrings,
	}

	roleClaimKeys := []string{}
	if roleClaimKey != "" {
		roleClaimKeys = append(roleClaimKeys, roleClaimKey)
	} else {
		roleClaimKeys = defaultRoleClaimKeys
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

func (c *OAuthClient) decideUserInfos(claims jwt.MapClaims, usernameClaimKey, avatarURLClaimKey string) (username, avatarURL string, err error) {

	username = ""
	usernameClaimKeys := []string{}
	if usernameClaimKey != "" {
		usernameClaimKeys = append(usernameClaimKeys, usernameClaimKey)
	} else {
		usernameClaimKeys = defaultUsernameClaimKeys
	}

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

	avatarURL = ""
	avatarURLClaimKeys := []string{}
	if usernameClaimKey != "" {
		avatarURLClaimKeys = append(avatarURLClaimKeys, avatarURLClaimKey)
	} else {
		avatarURLClaimKeys = defaultAvatarURLClaimKeys
	}
	for _, key := range avatarURLClaimKeys {
		val, ok := claims[key]
		if ok && val != nil {
			if str, ok := val.(string); ok && str != "" {
				avatarURL = str
				break
			}
		}
	}

	return username, avatarURL, nil
}

// As the go-oidc package does not provide any method to override fields like UserInfoEndpoint or AuthorizeEndpoint,
// NewOAuthClient needs to create a custom OIDC provider based on the provider created by the go-oidc package.
// createCustomOIDCProvider will first call the openid-configuration endpoint to retrieve all endpoints from the issuer URL,
// then pass user-provided URLs to override the existing URLs in the providerConfig struct.
// https://pkg.go.dev/github.com/coreos/go-oidc/v3@v3.11.0/oidc#NewProvider
// https://pkg.go.dev/github.com/coreos/go-oidc/v3@v3.11.0/oidc#ProviderConfig
func createCustomOIDCProvider(ctx context.Context, sso *model.ProjectSSOConfig_Oidc) (*oidc.Provider, error) {
	// Copied from go-oidc package
	issuer := sso.Issuer

	wellKnown := strings.TrimSuffix(issuer, "/") + "/.well-known/openid-configuration"
	req, err := http.NewRequest("GET", wellKnown, nil)
	if err != nil {
		return nil, err
	}

	client := http.DefaultClient
	if c := getClient(ctx); c != nil {
		client = c
	}
	resp, err := client.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %s", resp.Status, body)
	}

	var p providerJSON
	err = unmarshalResp(resp, body, &p)
	if err != nil {
		return nil, fmt.Errorf("oidc: failed to decode provider discovery object: %v", err)
	}
	// End of Copied from go-oidc package

	// Override the endpoints with the user-provided URLs
	providerConfig := oidc.ProviderConfig{
		IssuerURL: issuer,
		AuthURL: func() string {
			if sso.AuthorizationEndpoint != "" {
				return sso.AuthorizationEndpoint
			}
			return p.AuthURL
		}(),
		TokenURL: func() string {
			if sso.TokenEndpoint != "" {
				return sso.TokenEndpoint
			}
			return p.TokenURL
		}(),
		DeviceAuthURL: p.DeviceAuthURL,
		UserInfoURL: func() string {
			if sso.UserInfoEndpoint != "" {
				return sso.UserInfoEndpoint
			}
			return p.UserInfoURL
		}(),
		JWKSURL: p.JWKSURL,
	}

	return providerConfig.NewProvider(ctx), nil
}

func getClient(ctx context.Context) *http.Client {
	if c, ok := ctx.Value(oauth2.HTTPClient).(*http.Client); ok {
		return c
	}
	return nil
}

type providerJSON struct {
	Issuer        string   `json:"issuer"`
	AuthURL       string   `json:"authorization_endpoint"`
	TokenURL      string   `json:"token_endpoint"`
	DeviceAuthURL string   `json:"device_authorization_endpoint"`
	JWKSURL       string   `json:"jwks_uri"`
	UserInfoURL   string   `json:"userinfo_endpoint"`
	Algorithms    []string `json:"id_token_signing_alg_values_supported"`
}

func unmarshalResp(r *http.Response, body []byte, v interface{}) error {
	err := json.Unmarshal(body, &v)
	if err == nil {
		return nil
	}
	ct := r.Header.Get("Content-Type")
	mediaType, _, parseErr := mime.ParseMediaType(ct)
	if parseErr == nil && mediaType == "application/json" {
		return fmt.Errorf("got Content-Type = application/json, but could not unmarshal as JSON: %v", err)
	}
	return fmt.Errorf("expected Content-Type = application/json, got %q: %v", ct, err)
}
