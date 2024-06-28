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
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

const DefauleRSAKeySize = 2048

// GenerateRSAPems generates RSA key pair and the PEM encoding of them.
func GenerateRSAPems(size int) (private, public []byte, err error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, size)
	if err != nil {
		return
	}
	err = privateKey.Validate()
	if err != nil {
		return
	}

	publicBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return
	}
	public = pem.EncodeToMemory(
		&pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: publicBytes,
		},
	)

	privateBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return
	}
	private = pem.EncodeToMemory(
		&pem.Block{
			Type:  "PRIVATE KEY",
			Bytes: privateBytes,
		},
	)

	return
}

func ParseRSAPublicKeyFromPem(data []byte) (*rsa.PublicKey, error) {
	var err error
	data = bytes.TrimSpace(data)
	block, _ := pem.Decode(data)
	bytes := block.Bytes

	if x509.IsEncryptedPEMBlock(block) {
		bytes, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			return nil, err
		}
	}

	key, err := x509.ParsePKIXPublicKey(bytes)
	if err != nil {
		return nil, err
	}

	if k, ok := key.(*rsa.PublicKey); ok {
		return k, nil
	}
	return nil, errors.New("invalid key format, it must be a public RSA key")
}

func ParseRSAPrivateKeyFromPem(data []byte) (*rsa.PrivateKey, error) {
	var err error
	data = bytes.TrimSpace(data)
	block, _ := pem.Decode(data)
	bytes := block.Bytes

	if x509.IsEncryptedPEMBlock(block) {
		bytes, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			return nil, err
		}
	}

	key, err := x509.ParsePKCS8PrivateKey(bytes)
	if err != nil {
		return nil, err
	}

	if k, ok := key.(*rsa.PrivateKey); ok {
		return k, nil
	}
	return nil, errors.New("invalid key format, it must be a private RSA key")
}
