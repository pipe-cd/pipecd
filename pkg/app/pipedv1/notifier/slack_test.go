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

package notifier

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	config "github.com/pipe-cd/pipecd/pkg/configv1"
)

func Test_getAccountsAsString(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		accounts []string
		want     string
	}{
		{
			name:     "empty",
			accounts: []string{},
			want:     "",
		},
		{
			name:     "single",
			accounts: []string{"foo"},
			want:     "<@foo>",
		},
		{
			name:     "multiple",
			accounts: []string{"foo", "bar"},
			want:     "<@foo> <@bar>",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := getAccountsAsString(tt.accounts)
			if got != tt.want {
				t.Errorf("getAccountsAsString(): got %s, want %s", got, tt.want)
			}
		})
	}
}

func Test_getGroupsAsString(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		groups []string
		want   string
	}{
		{
			name:   "empty",
			groups: []string{},
			want:   "",
		},
		{
			name:   "single",
			groups: []string{"foo"},
			want:   "<!subteam^foo>",
		},
		{
			name:   "with correct format <!subteam^foo>",
			groups: []string{"<!subteam^foo>"},
			want:   "<!subteam^foo>",
		},
		{
			name:   "multiple",
			groups: []string{"foo", "bar"},
			want:   "<!subteam^foo> <!subteam^bar>",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := getGroupsAsString(tt.groups)
			if got != tt.want {
				t.Errorf("getGroupsAsString(): got %s, want %s", got, tt.want)
			}
		})
	}
}

func TestLoadOAuthToken(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	tokenFile := filepath.Join(dir, "token")
	require.NoError(t, os.WriteFile(tokenFile, []byte("xoxb-token\n"), 0600))
	crlfTokenFile := filepath.Join(dir, "crlf-token")
	require.NoError(t, os.WriteFile(crlfTokenFile, []byte("xoxb-token\r\n"), 0600))

	tests := []struct {
		name    string
		cfg     config.NotificationReceiverSlack
		want    string
		wantErr bool
	}{
		{
			name: "oauthToken",
			cfg:  config.NotificationReceiverSlack{OAuthToken: "xoxb-token"},
			want: "xoxb-token",
		},
		{
			name: "oauthTokenData",
			cfg:  config.NotificationReceiverSlack{OAuthTokenData: base64.StdEncoding.EncodeToString([]byte("xoxb-token"))},
			want: "xoxb-token",
		},
		{
			name: "oauthTokenData with a trailing newline",
			cfg:  config.NotificationReceiverSlack{OAuthTokenData: base64.StdEncoding.EncodeToString([]byte("xoxb-token\n"))},
			want: "xoxb-token",
		},
		{
			name: "oauthTokenFile with a trailing newline",
			cfg:  config.NotificationReceiverSlack{OAuthTokenFile: tokenFile},
			want: "xoxb-token",
		},
		{
			name: "oauthTokenFile with a trailing CRLF",
			cfg:  config.NotificationReceiverSlack{OAuthTokenFile: crlfTokenFile},
			want: "xoxb-token",
		},
		{
			name:    "invalid oauthTokenData",
			cfg:     config.NotificationReceiverSlack{OAuthTokenData: "not base64"},
			wantErr: true,
		},
		{
			name:    "missing oauthTokenFile",
			cfg:     config.NotificationReceiverSlack{OAuthTokenFile: filepath.Join(dir, "missing")},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := loadOAuthToken(tt.cfg)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}
