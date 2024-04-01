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
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pipe-cd/pipecd/pkg/config"
)

func TestSourceProcesser(t *testing.T) {
	t.Parallel()

	workspace, err := os.MkdirTemp("", "test-process-data")
	require.NoError(t, err)
	t.Cleanup(func() {
		os.RemoveAll(workspace)
	})
	dcr := testSecretDecrypter{
		prefix: "decrypted-",
	}

	fileData := map[string]string{
		"config.yaml":   "config-data",
		"resource.yaml": "echo {{ .attachment.config }} && echo {{ .encryptedSecrets.secret }}",
	}
	attachConfig := config.Attachment{
		Sources: map[string]string{
			"config": "config.yaml",
		},
		Targets: []string{
			"resource.yaml",
		},
	}
	secretConfig := config.SecretEncryption{
		EncryptedSecrets: map[string]string{
			"secret": "encrypted-secret",
		},
		DecryptionTargets: []string{
			"resource.yaml",
		},
	}

	appDir, err := os.MkdirTemp(workspace, "app-dir")
	require.NoError(t, err)

	for p, c := range fileData {
		p = filepath.Join(appDir, p)
		err = os.MkdirAll(filepath.Dir(p), 0700)
		require.NoError(t, err)
		err = os.WriteFile(p, []byte(c), 0600)
		require.NoError(t, err)
	}

	err = DecryptSecrets(appDir, secretConfig, dcr)
	require.NoError(t, err)
	err = AttachData(appDir, attachConfig)
	require.NoError(t, err)

	data, err := os.ReadFile(filepath.Join(appDir, "resource.yaml"))
	require.NoError(t, err)
	assert.Equal(t, "echo config-data && echo decrypted-encrypted-secret", string(data))
}
