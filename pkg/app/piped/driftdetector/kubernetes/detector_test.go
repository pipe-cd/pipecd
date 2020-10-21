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

package kubernetes

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	provider "github.com/pipe-cd/pipe/pkg/app/piped/cloudprovider/kubernetes"
	"github.com/pipe-cd/pipe/pkg/config"
)

func TestGroupManifests(t *testing.T) {
	testcases := []struct {
		name               string
		heads              []provider.Manifest
		lives              []provider.Manifest
		expectedAdds       []provider.Manifest
		expectedDeletes    []provider.Manifest
		expectedHeadInters []provider.Manifest
		expectedLiveInters []provider.Manifest
	}{
		{
			name: "empty list",
		},
		{
			name: "only adds",
			lives: []provider.Manifest{
				{Key: provider.ResourceKey{Name: "b"}},
				{Key: provider.ResourceKey{Name: "a"}},
			},
			expectedAdds: []provider.Manifest{
				{Key: provider.ResourceKey{Name: "a"}},
				{Key: provider.ResourceKey{Name: "b"}},
			},
		},
		{
			name: "only deletes",
			heads: []provider.Manifest{
				{Key: provider.ResourceKey{Name: "b"}},
				{Key: provider.ResourceKey{Name: "a"}},
			},
			expectedDeletes: []provider.Manifest{
				{Key: provider.ResourceKey{Name: "a"}},
				{Key: provider.ResourceKey{Name: "b"}},
			},
		},
		{
			name: "only inters",
			heads: []provider.Manifest{
				{Key: provider.ResourceKey{Name: "b"}},
				{Key: provider.ResourceKey{Name: "a"}},
			},
			lives: []provider.Manifest{
				{Key: provider.ResourceKey{Name: "a"}},
				{Key: provider.ResourceKey{Name: "b"}},
			},
			expectedHeadInters: []provider.Manifest{
				{Key: provider.ResourceKey{Name: "a"}},
				{Key: provider.ResourceKey{Name: "b"}},
			},
			expectedLiveInters: []provider.Manifest{
				{Key: provider.ResourceKey{Name: "a"}},
				{Key: provider.ResourceKey{Name: "b"}},
			},
		},
		{
			name: "all kinds",
			heads: []provider.Manifest{
				{Key: provider.ResourceKey{Name: "b"}},
				{Key: provider.ResourceKey{Name: "a"}},
				{Key: provider.ResourceKey{Name: "c"}},
			},
			lives: []provider.Manifest{
				{Key: provider.ResourceKey{Name: "a"}},
				{Key: provider.ResourceKey{Name: "d"}},
				{Key: provider.ResourceKey{Name: "b"}},
			},
			expectedAdds: []provider.Manifest{
				{Key: provider.ResourceKey{Name: "d"}},
			},
			expectedDeletes: []provider.Manifest{
				{Key: provider.ResourceKey{Name: "c"}},
			},
			expectedHeadInters: []provider.Manifest{
				{Key: provider.ResourceKey{Name: "a"}},
				{Key: provider.ResourceKey{Name: "b"}},
			},
			expectedLiveInters: []provider.Manifest{
				{Key: provider.ResourceKey{Name: "a"}},
				{Key: provider.ResourceKey{Name: "b"}},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			adds, deletes, headInters, liveInters := groupManifests(tc.heads, tc.lives)
			assert.Equal(t, tc.expectedAdds, adds)
			assert.Equal(t, tc.expectedDeletes, deletes)
			assert.Equal(t, tc.expectedHeadInters, headInters)
			assert.Equal(t, tc.expectedLiveInters, liveInters)
		})
	}
}

type testSealedSecretDecrypter struct {
	prefix string
}

func (d testSealedSecretDecrypter) Decrypt(text string) (string, error) {
	return d.prefix + text, nil
}

func TestDecryptSealedSecrets(t *testing.T) {
	dir, err := ioutil.TempDir("", "test-decrypting-sealed-secrets")
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
	dcr := testSealedSecretDecrypter{
		prefix: "decrypted-",
	}

	err = decryptSealedSecrets(dir, secrets, dcr)
	require.NoError(t, err)

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
