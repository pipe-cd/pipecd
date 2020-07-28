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
	"crypto/subtle"
	"fmt"
	"net/url"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
)

var (
	githubScopes = []string{"read:org"}
)

// Auth confirms username and password.
func (p *ProjectStaticUser) Auth(username, password string) error {
	if username == "" {
		return fmt.Errorf("username is empty")
	}
	if subtle.ConstantTimeCompare([]byte(p.Username), []byte(username)) != 1 {
		return fmt.Errorf("wrong username %q", username)
	}
	if password == "" {
		return fmt.Errorf("password is empty")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(p.PasswordHash), []byte(password)); err != nil {
		return fmt.Errorf("wrong password for username %q: %v", username, err)
	}
	return nil
}

// CreateAuthURL creates a auth url.
func (p *ProjectSingleSignOn) CreateAuthURL(project, apiURL, callbackPath, state string) (string, error) {
	switch p.Provider {
	case ProjectSingleSignOnProvider_GITHUB:
		if p.Github == nil {
			return "", fmt.Errorf("missing GitHub oauth in the SSO configuration")
		}
		return p.Github.CreateAuthURL(project, apiURL, callbackPath, state)
	default:
		return "", fmt.Errorf("not implemented")
	}
}

// CreateAuthURL creates a auth url.
func (p *ProjectSingleSignOn_GitHub) CreateAuthURL(project, apiURL, callbackPath, state string) (string, error) {
	u, err := url.Parse(p.BaseUrl)
	if err != nil {
		return "", err
	}

	cfg := oauth2.Config{
		ClientID: p.ClientId,
		Endpoint: oauth2.Endpoint{AuthURL: fmt.Sprintf("%s://%s%s", u.Scheme, u.Host, "/login/oauth/authorize")},
	}

	cfg.Scopes = githubScopes
	cfg.RedirectURL = fmt.Sprintf("%s%s?project=%s", apiURL, callbackPath, project)
	authURL := cfg.AuthCodeURL(state, oauth2.ApprovalForce, oauth2.AccessTypeOnline)

	return authURL, nil
}
