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

package crypto

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateRSAPems(t *testing.T) {
	private, public, err := GenerateRSAPems(0)
	require.Error(t, err)
	assert.Nil(t, private)
	assert.Nil(t, public)

	private, public, err = GenerateRSAPems(2048)
	require.NoError(t, err)

	public = bytes.TrimSpace(public)
	private = bytes.TrimSpace(private)

	assert.True(t, strings.HasPrefix(string(public), "-----BEGIN PUBLIC KEY-----"))
	assert.True(t, strings.HasSuffix(string(public), "-----END PUBLIC KEY-----"))

	assert.True(t, strings.HasPrefix(string(private), "-----BEGIN PRIVATE KEY-----"))
	assert.True(t, strings.HasSuffix(string(private), "-----END PRIVATE KEY-----"))

	publicKey, err := ParseRSAPublicKeyFromPem(public)
	require.NoError(t, err)
	assert.NotNil(t, publicKey)
}

func TestParseRSAKey(t *testing.T) {
	data, err := os.ReadFile("testdata/private-rsa-pem")
	require.NoError(t, err)

	privateKey, err := ParseRSAPrivateKeyFromPem(data)
	require.NoError(t, err)
	assert.NotNil(t, privateKey)

	data, err = os.ReadFile("testdata/public-rsa-pem")
	require.NoError(t, err)

	publicKey, err := ParseRSAPublicKeyFromPem(data)
	require.NoError(t, err)
	assert.NotNil(t, publicKey)
}
