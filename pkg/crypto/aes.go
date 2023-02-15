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
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
)

type AESEncryptDecrypter struct {
	key []byte
}

const aesKeySize = 32

// NewAESEncryptDecrypter reads the specified key file and returns an AES EncryptDecrypter.
func NewAESEncryptDecrypter(keyFile string) (*AESEncryptDecrypter, error) {
	key, err := os.ReadFile(keyFile)
	if err != nil {
		return nil, fmt.Errorf("unable to read key file: %v", err)
	}
	if len(key) < aesKeySize {
		return nil, fmt.Errorf("key size (%d) must be greater than or equal to %d", len(key), aesKeySize)
	}
	return &AESEncryptDecrypter{
		key: key[:aesKeySize],
	}, nil
}

func (a *AESEncryptDecrypter) Encrypt(text string) (string, error) {
	block, err := aes.NewCipher(a.key)
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

func (a *AESEncryptDecrypter) Decrypt(encryptedText string) (string, error) {
	encrypted, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(a.key)
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
