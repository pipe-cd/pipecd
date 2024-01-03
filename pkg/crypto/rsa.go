// Copyright 2024 The PipeCD Authors.
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
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"encoding/base64"
)

type RSAEncrypter struct {
	key *rsa.PublicKey
}

func NewRSAEncrypter(key []byte) (*RSAEncrypter, error) {
	k, err := ParseRSAPublicKeyFromPem(key)
	if err != nil {
		return nil, err
	}
	return &RSAEncrypter{
		key: k,
	}, nil
}

func (e *RSAEncrypter) Encrypt(text string) (string, error) {
	encryptedBytes, err := rsa.EncryptOAEP(
		sha512.New(),
		rand.Reader,
		e.key,
		[]byte(text),
		nil,
	)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(encryptedBytes), nil
}

type RSADecrypter struct {
	key *rsa.PrivateKey
}

func NewRSADecrypter(key []byte) (*RSADecrypter, error) {
	k, err := ParseRSAPrivateKeyFromPem(key)
	if err != nil {
		return nil, err
	}
	return &RSADecrypter{
		key: k,
	}, nil
}

func (d *RSADecrypter) Decrypt(encryptedText string) (string, error) {
	encryptedBytes, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		return "", err
	}

	text, err := rsa.DecryptOAEP(
		sha512.New(),
		rand.Reader,
		d.key,
		encryptedBytes,
		nil,
	)
	if err != nil {
		return "", err
	}

	return string(text), nil
}
