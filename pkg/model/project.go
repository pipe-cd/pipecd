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
	"context"
	"crypto/subtle"
	"fmt"
	"net/url"
	"strings"

	"github.com/google/go-github/v29/github"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
)

var (
	githubScopes = []string{"read:org"}
)

// UserInfo is the login user information.
type UserInfo struct {
	AvatarURL string
	Username  string
	Role      Role_ProjectRole
}

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

// GenerateAuthCodeURL generates an auth URL for the specified configuration.
func (p *ProjectSingleSignOn) GenerateAuthCodeURL(project, apiURL, callbackPath, state string) (string, error) {
	switch p.Provider {
	case ProjectSingleSignOnProvider_GITHUB:
		if p.Github == nil {
			return "", fmt.Errorf("missing GitHub oauth in the SSO configuration")
		}
		return p.Github.GenerateAuthCodeURL(project, apiURL, callbackPath, state)
	default:
		return "", fmt.Errorf("not implemented")
	}
}

// GenerateAuthCodeURL generates an auth URL for the specified configuration.
func (p *ProjectSingleSignOn_GitHub) GenerateAuthCodeURL(project, apiURL, callbackPath, state string) (string, error) {
	u, err := url.Parse(p.BaseUrl)
	if err != nil {
		return "", err
	}

	cfg := oauth2.Config{
		ClientID: p.ClientId,
		Endpoint: oauth2.Endpoint{AuthURL: fmt.Sprintf("%s://%s%s", u.Scheme, u.Host, "/login/oauth/authorize")},
	}

	cfg.Scopes = githubScopes
	apiURL = strings.TrimSuffix(apiURL, "/")
	cfg.RedirectURL = fmt.Sprintf("%s%s?project=%s", apiURL, callbackPath, project)
	authURL := cfg.AuthCodeURL(state, oauth2.ApprovalForce, oauth2.AccessTypeOnline)

	return authURL, nil
}

// GenerateUserInfo generates a login user information.
func (p *ProjectSingleSignOn) GenerateUserInfo(ctx context.Context, code string) (*UserInfo, error) {
	switch p.Provider {
	case ProjectSingleSignOnProvider_GITHUB:
		if p.Github == nil {
			return nil, fmt.Errorf("missing GitHub oauth in the SSO configuration")
		}
		return p.Github.GenerateUserInfo(ctx, code)
	default:
		return nil, fmt.Errorf("not implemented")
	}
}

// GenerateUserInfo generates a login user information.
func (p *ProjectSingleSignOn_GitHub) GenerateUserInfo(ctx context.Context, code string) (*UserInfo, error) {
	bu, err := url.Parse(p.BaseUrl)
	if err != nil {
		return nil, err
	}
	uu, err := url.Parse(p.UploadUrl)
	if err != nil {
		return nil, err
	}

	cfg := oauth2.Config{
		ClientID:     p.ClientId,
		ClientSecret: p.ClientSecret,
		Endpoint:     oauth2.Endpoint{TokenURL: fmt.Sprintf("%s://%s%s", bu.Scheme, bu.Host, "/login/oauth/access_token")},
	}
	token, err := cfg.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	cli := github.NewClient(cfg.Client(ctx, token))
	if !strings.HasSuffix(bu.Path, "/") {
		bu.Path += "/"
	}
	cli.BaseURL = bu
	if !strings.HasSuffix(uu.Path, "/") {
		bu.Path += "/"
	}
	cli.UploadURL = uu

	user, _, err := cli.Users.Get(ctx, "")
	if err != nil {
		return nil, err
	}
	teams, _, err := cli.Teams.ListUserTeams(ctx, nil)
	if err != nil {
		return nil, err
	}
	role, err := p.decideRole(user.GetLogin(), teams)
	if err != nil {
		return nil, err
	}

	return &UserInfo{
		Username:  user.GetLogin(),
		AvatarURL: user.GetAvatarURL(),
		Role:      *role,
	}, nil
}

func (p *ProjectSingleSignOn_GitHub) decideRole(user string, teams []*github.Team) (*Role_ProjectRole, error) {
	var viewer, editor bool
	for _, team := range teams {
		slug := team.GetSlug()
		if p.Org != team.Organization.GetLogin() || slug == "" {
			continue
		}
		switch slug {
		case p.AdminTeam:
			r := Role_ADMIN
			return &r, nil
		case p.EditorTeam:
			editor = true
		case p.ViewerTeam:
			viewer = true
		}
	}
	if editor {
		r := Role_EDITOR
		return &r, nil
	}
	if viewer {
		r := Role_VIEWER
		return &r, nil
	}
	return nil, fmt.Errorf("user (%s) not found in any of the %d project teams", user, len(teams))
}
