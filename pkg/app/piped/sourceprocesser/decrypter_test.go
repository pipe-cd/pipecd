// Copyright 2023 The PipeCD Authors.
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

type testSecretDecrypter struct {
	prefix string
}

func (d testSecretDecrypter) Decrypt(text string) (string, error) {
	return d.prefix + text, nil
}

func TestDecryptSecrets(t *testing.T) {
	t.Parallel()

	workspace, err := os.MkdirTemp("", "test-decrypt-secrets")
	require.NoError(t, err)
	t.Cleanup(func() {
		os.RemoveAll(workspace)
	})
	dcr := testSecretDecrypter{
		prefix: "decrypted-",
	}

	testcases := []struct {
		name                string
		sources             map[string]string
		encryption          config.SecretEncryption
		expected            map[string]string
		expectedErrorPrefix string
	}{
		{
			name: "target not found",
			sources: map[string]string{
				"resource.yaml": "resource-data",
			},
			encryption: config.SecretEncryption{
				EncryptedSecrets: map[string]string{
					"password": "encrypted-password",
				},
				DecryptionTargets: []string{
					"not-found-resource.yaml",
				},
			},
			expectedErrorPrefix: "failed to parse decryption target not-found-resource.yaml",
		},
		{
			name: "the target is not using any encrypted secret",
			sources: map[string]string{
				"resource.yaml": "resource-data",
			},
			encryption: config.SecretEncryption{
				EncryptedSecrets: map[string]string{
					"password": "encrypted-password",
				},
				DecryptionTargets: []string{
					"resource.yaml",
				},
			},
			expected: map[string]string{
				"resource.yaml": "resource-data",
			},
		},
		{
			name: "single target",
			sources: map[string]string{
				"resource.yaml": "resource-data: {{ .encryptedSecrets.password }}",
			},
			encryption: config.SecretEncryption{
				EncryptedSecrets: map[string]string{
					"password": "encrypted-password",
				},
				DecryptionTargets: []string{
					"resource.yaml",
				},
			},
			expected: map[string]string{
				"resource.yaml": "resource-data: decrypted-encrypted-password",
			},
		},
		{
			name: "multi targets",
			sources: map[string]string{
				"resource1.yaml": "resource1-data: {{ .encryptedSecrets.password }}",
				"resource2.yaml": "resource2-data: bar is {{ .encryptedSecrets.bar }}, foo is {{ .encryptedSecrets.foo }}",
			},
			encryption: config.SecretEncryption{
				EncryptedSecrets: map[string]string{
					"password": "encrypted-password",
					"foo":      "encrypted-foo",
					"bar":      "encrypted-bar",
				},
				DecryptionTargets: []string{
					"resource1.yaml",
					"resource2.yaml",
				},
			},
			expected: map[string]string{
				"resource1.yaml": "resource1-data: decrypted-encrypted-password",
				"resource2.yaml": "resource2-data: bar is decrypted-encrypted-bar, foo is decrypted-encrypted-foo",
			},
		},
		{
			name: "target is using a nonexistent encrypted secret",
			sources: map[string]string{
				"resource.yaml": "resource-data: {{ .encryptedSecrets.password }}, {{ .encryptedSecrets.nonexistent }}",
			},
			encryption: config.SecretEncryption{
				EncryptedSecrets: map[string]string{
					"password": "encrypted-password",
				},
				DecryptionTargets: []string{
					"resource.yaml",
				},
			},
			expectedErrorPrefix: `failed to render decryption target resource.yaml (template: resource.yaml:1:69: executing "resource.yaml" at <.encryptedSecrets.nonexistent>: map has no entry for key "nonexistent")`,
		},
		{
			name: "sprig functions",
			sources: map[string]string{
				"resource.yaml": "resource-data: {{ .encryptedSecrets.password | b64enc }}",
			},
			encryption: config.SecretEncryption{
				EncryptedSecrets: map[string]string{
					"password": "encrypted-password",
				},
				DecryptionTargets: []string{
					"resource.yaml",
				},
			},
			expected: map[string]string{
				"resource.yaml": "resource-data: ZGVjcnlwdGVkLWVuY3J5cHRlZC1wYXNzd29yZA==",
			},
		},
		{
			name: "sub directory",
			sources: map[string]string{
				"sub/dir/resource.yaml": "resource-data: {{ .encryptedSecrets.password }}",
			},
			encryption: config.SecretEncryption{
				EncryptedSecrets: map[string]string{
					"password": "encrypted-password",
				},
				DecryptionTargets: []string{
					"sub/dir/resource.yaml",
				},
			},
			expected: map[string]string{
				"sub/dir/resource.yaml": "resource-data: decrypted-encrypted-password",
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			appDir, err := os.MkdirTemp(workspace, "app-dir")
			require.NoError(t, err)

			// Prepare source files.
			for p, c := range tc.sources {
				p = filepath.Join(appDir, p)
				err := os.MkdirAll(filepath.Dir(p), 0700)
				require.NoError(t, err)
				err = os.WriteFile(p, []byte(c), 0600)
				require.NoError(t, err)
			}

			err = DecryptSecrets(appDir, tc.encryption, dcr)
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
