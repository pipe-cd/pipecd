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

package crypto

import (
	"bytes"
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

	assert.True(t, strings.HasPrefix(string(public), "-----BEGIN RSA PUBLIC KEY-----"))
	assert.True(t, strings.HasSuffix(string(public), "-----END RSA PUBLIC KEY-----"))

	assert.True(t, strings.HasPrefix(string(private), "-----BEGIN RSA PRIVATE KEY-----"))
	assert.True(t, strings.HasSuffix(string(private), "-----END RSA PRIVATE KEY-----"))

	publicKey, err := ParseRSAPublicKeyFromPem(public)
	require.NoError(t, err)
	assert.NotNil(t, publicKey)
}

func TestLoadRSAKey(t *testing.T) {
	privateKey, err := LoadRSAPrivateKey("testdata/private-rsa-pem")
	require.NoError(t, err)
	assert.NotNil(t, privateKey)

	publicKey, err := LoadRSAPublicKey("testdata/public-rsa-pem")
	require.NoError(t, err)
	assert.NotNil(t, publicKey)
}
