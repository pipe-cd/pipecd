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
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"io/ioutil"
)

// Decrypter is a interface to decrypt data.
type Decrypter interface {
	Decrypt(encryptedText string) (string, error)
}

type decrypter struct {
	key []byte
}

// NewDecrypter returns a new decrypter.
func NewDecrypter(keyFile string) (Decrypter, error) {
	key, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return nil, fmt.Errorf("unable to read key file: %v", err)
	}
	if len(key) < aes256size {
		return nil, fmt.Errorf("key size (%d) must be greater than or equal to %d", len(key), aes256size)
	}
	return &decrypter{
		key: key[:aes256size],
	}, nil
}

func (d *decrypter) Decrypt(encryptedText string) (string, error) {
	encrypted, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(d.key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := encrypted[:gcm.NonceSize()]
	text, err := gcm.Open(nil, nonce, encrypted[gcm.NonceSize():], nil)
	if err != nil {
		return "", err
	}

	return string(text), nil
}
