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
	"fmt"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var (
	githubScopes = []string{"read:org"}
)

// CreateAuthURL creates a auth url.
func (p *ProjectSingleSignOn_GitHub) CreateAuthURL(project, apiURL, callbackPath, state string) (string, error) {
	cfg := oauth2.Config{
		ClientID: p.ClientId,
		Endpoint: github.Endpoint,
	}

	cfg.Scopes = githubScopes
	cfg.RedirectURL = fmt.Sprintf("%s%s?project=%s", apiURL, callbackPath, project)
	authURL := cfg.AuthCodeURL(state, oauth2.ApprovalForce, oauth2.AccessTypeOnline)

	return authURL, nil
}
