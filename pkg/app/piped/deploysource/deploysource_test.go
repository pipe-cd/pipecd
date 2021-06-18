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

package deploysource

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pipe-cd/pipe/pkg/config"
)

type testSecretDecrypter struct {
	prefix string
}

func (d testSecretDecrypter) Decrypt(text string) (string, error) {
	return d.prefix + text, nil
}

func TestDecryptSecrets(t *testing.T) {
	dir, err := ioutil.TempDir("", "test-decrypting-secrets")
	require.NoError(t, err)

	defer os.RemoveAll(dir)

	err = ioutil.WriteFile(filepath.Join(dir, "replacing.yaml"), []byte(`
apiVersion: "pipecd.dev/v1beta1"
kind: SealedSecret
spec:
  template: |
    apiVersion: v1
    kind: Secret
    metadata:
      name: mysecret
    type: Opaque
    data:
      username: {{ .encryptedItems.username }}
      password: {{ .encryptedItems.password }}
  encryptedItems:
    username: encrypted-username
    password: encrypted-password
`),
		0644,
	)
	require.NoError(t, err)

	err = ioutil.WriteFile(filepath.Join(dir, "copy.yaml"), []byte(`
apiVersion: "pipecd.dev/v1beta1"
kind: SealedSecret
spec:
  encryptedData: encrypted-data
`),
		0644,
	)

	require.NoError(t, err)

	secrets := []config.SealedSecretMapping{
		{
			Path: "replacing.yaml",
		},
		{
			Path:        "copy.yaml",
			OutFilename: "new-copy.yaml",
		},
		{
			Path:   "copy.yaml",
			OutDir: ".credentials",
		},
	}
	dcr := testSecretDecrypter{
		prefix: "decrypted-",
	}

	for _, s := range secrets {
		err = decryptSecret(dir, s, dcr)
		require.NoError(t, err)
	}

	data, err := ioutil.ReadFile(filepath.Join(dir, "replacing.yaml"))
	require.NoError(t, err)
	assert.Equal(t,
		`apiVersion: v1
kind: Secret
metadata:
  name: mysecret
type: Opaque
data:
  username: decrypted-encrypted-username
  password: decrypted-encrypted-password
`,
		string(data),
	)

	data, err = ioutil.ReadFile(filepath.Join(dir, "new-copy.yaml"))
	require.NoError(t, err)
	assert.Equal(t,
		`decrypted-encrypted-data`,
		string(data),
	)

	data, err = ioutil.ReadFile(filepath.Join(dir, ".credentials/copy.yaml"))
	require.NoError(t, err)
	assert.Equal(t,
		`decrypted-encrypted-data`,
		string(data),
	)
}
