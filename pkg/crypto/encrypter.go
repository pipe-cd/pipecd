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
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io/ioutil"
)

// Encrypter is a interface to encrypt data.
type Encrypter interface {
	Encrypt(text string) (string, error)
}

type encrypter struct {
	key []byte
}

const (
	aes256size = 32
)

// NewEncrypter returns a new encrypter.
func NewEncrypter(keyFile string) (Encrypter, error) {
	key, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return nil, fmt.Errorf("unable to read key file: %v", err)
	}
	if len(key) < aes256size {
		return nil, fmt.Errorf("key size (%d) must be greater than or equal to %d", len(key), aes256size)
	}
	return &encrypter{
		key: key[:aes256size],
	}, nil
}

func (e *encrypter) Encrypt(text string) (string, error) {
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		return "", err
	}

	encrypted := gcm.Seal(nil, nonce, []byte(text), nil)
	encrypted = append(nonce, encrypted...)

	return base64.StdEncoding.EncodeToString(encrypted), nil
}
