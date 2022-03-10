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

package config

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSealedSecretConfig(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		fileName           string
		expectedKind       Kind
		expectedAPIVersion string
		expectedSpec       interface{}
		expectedError      error
	}{
		{
			fileName:           "testdata/sealedsecret/invalid.yaml",
			expectedKind:       KindSealedSecret,
			expectedAPIVersion: "pipecd.dev/v1beta1",
			expectedSpec:       nil,
			expectedError:      fmt.Errorf("either encryptedData or encryptedItems must be set"),
		},
		{
			fileName:           "testdata/sealedsecret/ok.yaml",
			expectedKind:       KindSealedSecret,
			expectedAPIVersion: "pipecd.dev/v1beta1",
			expectedSpec: &SealedSecretSpec{
				Template: `apiVersion: v1
kind: Secret
metadata:
  name: mysecret
type: Opaque
data:
  username: {{ .encryptedItems.username }}
  password: {{ .encryptedItems.password }}
`,
				EncryptedItems: map[string]string{
					"username": "encrypted-username",
					"password": "encrypted-password",
				},
			},
			expectedError: nil,
		},
	}
	for _, tc := range testcases {
		tc := tc
		t.Run(tc.fileName, func(t *testing.T) {
			t.Parallel()

			cfg, err := LoadFromYAML(tc.fileName)
			require.Equal(t, tc.expectedError, err)
			if err == nil {
				assert.Equal(t, tc.expectedKind, cfg.Kind)
				assert.Equal(t, tc.expectedAPIVersion, cfg.APIVersion)
				assert.Equal(t, tc.expectedSpec, cfg.spec)
			}
		})
	}
}

type testSealedSecretDecrypter struct {
	prefix string
}

func (d testSealedSecretDecrypter) Decrypt(text string) (string, error) {
	return d.prefix + text, nil
}

func TestSealedSecretRenderOrifinalContent(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name          string
		spec          *SealedSecretSpec
		expected      string
		expectedError error
	}{
		{
			name: "without template",
			spec: &SealedSecretSpec{
				EncryptedData: "encrypted-data",
			},
			expected: "decrypted-encrypted-data",
		},
		{
			name: "render with the specified template",
			spec: &SealedSecretSpec{
				Template: "Hello {{ .encryptedItems.username }}",
				EncryptedItems: map[string]string{
					"username": "encrypted-username",
				},
			},
			expected: "Hello decrypted-encrypted-username",
		},
		{
			name: "missing data",
			spec: &SealedSecretSpec{
				Template: "Hello {{ .encryptedItems.username }}, {{ .encryptedItems.other }}",
				EncryptedItems: map[string]string{
					"username": "PipeCD",
				},
			},
			expectedError: fmt.Errorf(`template: sealedsecret:1:56: executing "sealedsecret" at <.encryptedItems.other>: map has no entry for key "other"`),
		},
	}

	dcr := &testSealedSecretDecrypter{
		prefix: "decrypted-",
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := tc.spec.RenderOriginalContent(dcr)
			assert.Equal(t, tc.expected, string(got))
			if tc.expectedError != nil && err != nil {
				assert.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				assert.Equal(t, tc.expectedError, err)
			}
		})
	}
}
