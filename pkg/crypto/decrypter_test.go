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
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecrypt(t *testing.T) {
	text := "foo-bar-baz"

	e, err := NewEncrypter("testdata/secret")
	require.NoError(t, err)
	require.NotNil(t, e)

	cipherText, err := e.Encrypt(text)
	require.NoError(t, err)
	require.True(t, len(cipherText) > 0)

	d, err := NewDecrypter("testdata/secret")
	require.NoError(t, err)
	require.NotNil(t, d)

	decrepted, err := d.Decrypt(cipherText)
	require.NoError(t, err)
	require.Equal(t, text, decrepted)
}
