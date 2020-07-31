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

package github

import (
	"testing"

	"github.com/google/go-github/v29/github"
	"github.com/stretchr/testify/assert"

	"github.com/pipe-cd/pipe/pkg/model"
)

func stringPointer(s string) *string { return &s }

func TestDecideRole(t *testing.T) {
	cases := []struct {
		name     string
		username string
		teams    []*github.Team
		role     model.Role_ProjectRole
		wantErr  bool
	}{
		{
			name:     "nothing",
			username: "foo",
			teams: []*github.Team{
				{
					Slug: stringPointer("team1"),
				},
			},
			wantErr: true,
		},
		{
			name:     "admin",
			username: "foo",
			teams: []*github.Team{
				{
					Slug: stringPointer("team-admin"),
				},
				{
					Slug: stringPointer("team-editor"),
				},
				{
					Slug: stringPointer("team-viewer"),
				},
			},
			role: model.Role_ADMIN,
		},
		{
			name:     "editor",
			username: "foo",
			teams: []*github.Team{
				{
					Slug: stringPointer("team1"),
				},
				{
					Slug: stringPointer("team-editor"),
				},
				{
					Slug: stringPointer("team-viewer"),
				},
			},
			role: model.Role_EDITOR,
		},
		{
			name:     "viewer",
			username: "foo",
			teams: []*github.Team{
				{
					Slug: stringPointer("team1"),
				},
				{
					Slug: stringPointer("team2"),
				},
				{
					Slug: stringPointer("team-viewer"),
				},
			},
			role: model.Role_VIEWER,
		},
	}

	oc := &OAuthClient{
		adminTeam:  "team-admin",
		editorTeam: "team-editor",
		viewerTeam: "team-viewer",
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			role, err := oc.decideRole(tc.username, tc.teams)
			assert.Equal(t, tc.wantErr, err != nil)
			if err == nil {
				assert.Equal(t, tc.role, *role)
			}
		})
	}
}
