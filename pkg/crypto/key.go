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
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
)

const DefauleRSAKeySize = 2048

// GenerateRSAPems generates RSA key pair and the PEM encoding of them.
func GenerateRSAPems(size int) (private, public []byte, err error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, size)
	if err != nil {
		return
	}

	public = pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: x509.MarshalPKCS1PublicKey(&privateKey.PublicKey),
		},
	)

	private = pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
		},
	)

	return
}

func LoadRSAPublicKey(path string) (*rsa.PublicKey, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	data = bytes.TrimSpace(data)
	return ParseRSAPublicKeyFromPem(data)
}

func ParseRSAPublicKeyFromPem(data []byte) (*rsa.PublicKey, error) {
	var err error
	block, _ := pem.Decode(data)
	bytes := block.Bytes

	if x509.IsEncryptedPEMBlock(block) {
		bytes, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			return nil, err
		}
	}

	key, err := x509.ParsePKCS1PublicKey(bytes)
	if err != nil {
		return nil, err
	}

	return key, nil
}

func LoadRSAPrivateKey(path string) (*rsa.PrivateKey, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	data = bytes.TrimSpace(data)
	return ParseRSAPrivateKeyFromPem(data)
}

func ParseRSAPrivateKeyFromPem(data []byte) (*rsa.PrivateKey, error) {
	var err error
	block, _ := pem.Decode(data)
	bytes := block.Bytes

	if x509.IsEncryptedPEMBlock(block) {
		bytes, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			return nil, err
		}
	}

	key, err := x509.ParsePKCS1PrivateKey(bytes)
	if err != nil {
		return nil, err
	}

	return key, nil
}
