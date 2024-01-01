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
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pipe-cd/pipecd/pkg/config"
)

func TestGenerateSSHConfig(t *testing.T) {
	testcases := []struct {
		name        string
		cfg         config.PipedGit
		expected    string
		expectedErr error
	}{
		{
			name: "default",
			cfg: config.PipedGit{
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
			cfg: config.PipedGit{
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
			cfg: config.PipedGit{
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
			got, err := generateSSHConfig(tc.cfg, sshKeyFile)
			assert.Equal(t, tc.expected, got)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}
