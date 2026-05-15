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

package git

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	configv1 "github.com/pipe-cd/pipecd/pkg/configv1"
)

func TestGenerateSSHConfig(t *testing.T) {
	testcases := []struct {
		name        string
		cfg         configv1.PipedGit
		expected    string
		expectedErr error
	}{
		{
			name: "default",
			cfg: configv1.PipedGit{
				SSHKeyFile: "/tmp/piped-secret/ssh-key",
			},
			expected: `
Host github.com
    Hostname github.com
    User git
    IdentityFile /etc/piped-secret/ssh-key
    UserKnownHostsFile /dev/null
    StrictHostKeyChecking no
`,
			expectedErr: nil,
		},
		{
			name: "host is configured",
			cfg: configv1.PipedGit{
				Host:       "gitlab.com",
				SSHKeyFile: "/tmp/piped-secret/ssh-key",
			},
			expected: `
Host gitlab.com
    Hostname gitlab.com
    User git
    IdentityFile /etc/piped-secret/ssh-key
    UserKnownHostsFile /dev/null
    StrictHostKeyChecking no
`,
			expectedErr: nil,
		},
		{
			name: "host and hostname are configured",
			cfg: configv1.PipedGit{
				Host:       "gitlab.com",
				HostName:   "gitlab.com",
				SSHKeyFile: "/tmp/piped-secret/ssh-key",
			},
			expected: `
Host gitlab.com
    Hostname gitlab.com
    User git
    IdentityFile /etc/piped-secret/ssh-key
    UserKnownHostsFile /dev/null
    StrictHostKeyChecking no
`,
			expectedErr: nil,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			sshKeyFile := "/etc/piped-secret/ssh-key"
			got, err := generateSSHConfig(sshConfig{
				Host:         tc.cfg.Host,
				HostName:     tc.cfg.HostName,
				IdentityFile: sshKeyFile,
			})
			assert.Equal(t, tc.expected, got)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestAddSSHConfig(t *testing.T) {
	tempHome := t.TempDir()
	cfg := configv1.PipedGit{
		SSHConfigFilePath: tempHome + "/.ssh/config",
		Host:              "gitlab.com",
		SSHKeyData:        "bGVnYWN5",
		SSHKeys: []configv1.PipedSSHKeyEntry{
			{
				Host:       "github.com",
				SSHKeyData: "ZXh0cmE=",
			},
		},
	}

	tempDir, err := AddSSHConfig(cfg)
	assert.NoError(t, err)
	assert.NotEmpty(t, tempDir)

	files, err := os.ReadDir(tempDir)
	assert.NoError(t, err)
	assert.Len(t, files, 2)

	for _, file := range files {
		info, err := file.Info()
		assert.NoError(t, err)
		assert.Equal(t, os.FileMode(0600), info.Mode().Perm(), "SSH identity file must have 0600 permissions")
	}

	cfgContent, err := os.ReadFile(tempHome + "/.ssh/config")
	assert.NoError(t, err)
	cfgStr := string(cfgContent)
	gitlabIdx := strings.Index(cfgStr, "Host gitlab.com")
	githubIdx := strings.Index(cfgStr, "Host github.com")
	assert.True(t, gitlabIdx >= 0, "expected Host gitlab.com block in ssh config")
	assert.True(t, githubIdx >= 0, "expected Host github.com block in ssh config")
	assert.True(t, gitlabIdx < githubIdx, "legacy host (gitlab.com) must appear before extra host (github.com)")
}
