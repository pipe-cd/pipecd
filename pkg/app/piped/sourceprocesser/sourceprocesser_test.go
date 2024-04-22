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

package sourceprocesser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pipe-cd/pipecd/pkg/config"
)

func TestSourceProcessor(t *testing.T) {
	t.Parallel()

	workspace, err := os.MkdirTemp("", "test-process-data")
	require.NoError(t, err)
	t.Cleanup(func() {
		os.RemoveAll(workspace)
	})
	dcr := testSecretDecrypter{
		prefix: "decrypted-",
	}

	testcases := []struct {
		name                string
		fileData            map[string]string
		attachConfig        config.Attachment
		secretConfig        config.SecretEncryption
		expected            map[string]string
		expectedErrorPrefix string
	}{
		{
			name: "target not found",
			fileData: map[string]string{
				"resource.yaml": "resource-data",
			},
			secretConfig: config.SecretEncryption{
				EncryptedSecrets: map[string]string{
					"password": "encrypted-password",
				},
				DecryptionTargets: []string{
					"not-found-resource.yaml",
				},
			},
			expectedErrorPrefix: "failed to parse target file not-found-resource.yaml",
		},
		{
			name: "attachment work with secret decryption",
			fileData: map[string]string{
				"config.yaml":   "config-data",
				"resource.yaml": "echo {{ .attachment.config }} && echo {{ .encryptedSecrets.secret }}",
			},
			attachConfig: config.Attachment{
				Sources: map[string]string{
					"config": "config.yaml",
				},
				Targets: []string{
					"resource.yaml",
				},
			},
			secretConfig: config.SecretEncryption{
				EncryptedSecrets: map[string]string{
					"secret": "encrypted-secret",
				},
				DecryptionTargets: []string{
					"resource.yaml",
				},
			},
			expected: map[string]string{
				"resource.yaml": "echo config-data && echo decrypted-encrypted-secret",
			},
		},
		{
			name: "secret decryption combined attachement",
			fileData: map[string]string{
				"config.yaml":   "{{ .encryptedSecrets.secret }}",
				"resource.yaml": "echo {{ .attachment.config }}",
			},
			attachConfig: config.Attachment{
				Sources: map[string]string{
					"config": "config.yaml",
				},
				Targets: []string{
					"resource.yaml",
				},
			},
			secretConfig: config.SecretEncryption{
				EncryptedSecrets: map[string]string{
					"secret": "encrypted-secret",
				},
				DecryptionTargets: []string{
					"config.yaml",
				},
			},
			expected: map[string]string{
				"resource.yaml": "echo decrypted-encrypted-secret",
			},
		},
		{
			name: "sprig function",
			fileData: map[string]string{
				"config.yaml":   "config-data",
				"resource.yaml": "echo {{ .attachment.config | b64enc }}",
			},
			attachConfig: config.Attachment{
				Sources: map[string]string{
					"config": "config.yaml",
				},
				Targets: []string{
					"resource.yaml",
				},
			},
			expected: map[string]string{
				"resource.yaml": "echo Y29uZmlnLWRhdGE=", // base64 encode of "config-data"
			},
		},
		{
			name: "sub directory",
			fileData: map[string]string{
				"sub/dir/config.yaml": "config-data",
				"resource.yaml":       "echo {{ .attachment.config }}",
			},
			attachConfig: config.Attachment{
				Sources: map[string]string{
					"config": "sub/dir/config.yaml",
				},
				Targets: []string{
					"resource.yaml",
				},
			},
			expected: map[string]string{
				"resource.yaml": "echo config-data",
			},
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			appDir, err := os.MkdirTemp(workspace, "app-dir")
			require.NoError(t, err)

			for p, c := range tc.fileData {
				p = filepath.Join(appDir, p)
				err = os.MkdirAll(filepath.Dir(p), 0700)
				require.NoError(t, err)
				err = os.WriteFile(p, []byte(c), 0600)
				require.NoError(t, err)
			}

			atp := NewAttachmentProcessor(&tc.attachConfig)
			sdp := NewSecretDecrypterProcessor(&tc.secretConfig, dcr)
			sp := NewSourceProcessor(appDir, []SourceTemplateProcessor{sdp, atp}...)

			err = sp.Process()
			if tc.expectedErrorPrefix != "" {
				require.Error(t, err)
				assert.True(t, strings.HasPrefix(err.Error(), tc.expectedErrorPrefix), fmt.Sprintf("Error: %v", err))
			} else {
				require.NoError(t, err)
			}

			for p, c := range tc.expected {
				p = filepath.Join(appDir, p)
				data, err := os.ReadFile(p)
				require.NoError(t, err)
				assert.Equal(t, c, string(data))
			}
		})
	}

}
