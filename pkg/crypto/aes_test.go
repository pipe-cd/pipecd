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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAESEncryptDecrypterInvalidKey(t *testing.T) {
	ed, err := NewAESEncryptDecrypter("testdata/short-key")
	assert.Nil(t, ed)
	assert.Equal(t, "key size (9) must be greater than or equal to 32", err.Error())
}

func TestAESEncryptDecrypt(t *testing.T) {
	text := "foo-bar-baz"

	ed, err := NewAESEncryptDecrypter("testdata/key")
	require.NoError(t, err)
	require.NotNil(t, ed)

	encryptedText, err := ed.Encrypt(text)
	require.NoError(t, err)
	assert.True(t, len(encryptedText) > 0)

	decrypted, err := ed.Decrypt(encryptedText)
	require.NoError(t, err)
	assert.Equal(t, text, decrypted)
}
