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
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io"
)

// HybridEncrypter uses RSA to encrypt a randomly generated key for a symmetric AES-GCM.
// RSA is able to encrypt only a very limited amount of data. In order
// to encrypt reasonable amounts of data a hybrid scheme is commonly used.
type HybridEncrypter struct {
	key *rsa.PublicKey
}

func NewHybridEncrypter(key []byte) (*HybridEncrypter, error) {
	k, err := ParseRSAPublicKeyFromPem(key)
	if err != nil {
		return nil, err
	}
	return &HybridEncrypter{
		key: k,
	}, nil
}

// Encrypt performs a regular AES-GCM + RSA-OAEP encryption.
// The output string is:
//   RSA ciphertext length || RSA ciphertext || AES ciphertext
//
// The implementation of this function was brought from well known Bitnami's SealedSecret library.
// https://github.com/bitnami-labs/sealed-secrets/blob/master/pkg/crypto/crypto.go#L35
func (e *HybridEncrypter) Encrypt(text string) (string, error) {
	// Generate a random symmetric key.
	// The hybrid scheme should use at least a 16-byte symmetric key.
	symKey := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, symKey); err != nil {
		return "", err
	}

	block, err := aes.NewCipher(symKey)
	if err != nil {
		return "", err
	}

	aed, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Encrypt the symmetric key.
	rsaCipherText, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, e.key, symKey, nil)
	if err != nil {
		return "", err
	}

	// First 2 bytes are RSA ciphertext length, so we can separate all the pieces later.
	ciphertext := make([]byte, 2)
	binary.BigEndian.PutUint16(ciphertext, uint16(len(rsaCipherText)))
	ciphertext = append(ciphertext, rsaCipherText...)

	// Symmetric key is only used once, so zero nonce is ok.
	zeroNonce := make([]byte, aed.NonceSize())

	// Append symmetrically encrypted secret.
	ciphertext = aed.Seal(ciphertext, zeroNonce, []byte(text), nil)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

type HybridDecrypter struct {
	key *rsa.PrivateKey
}

func NewHybridDecrypter(key []byte) (*HybridDecrypter, error) {
	k, err := ParseRSAPrivateKeyFromPem(key)
	if err != nil {
		return nil, err
	}
	return &HybridDecrypter{
		key: k,
	}, nil
}

// Decrypt performs a regular AES-GCM + RSA-OAEP decryption.
//
// The implementation of this function was brought from well known Bitnami's SealedSecret library.
// https://github.com/bitnami-labs/sealed-secrets/blob/master/pkg/crypto/crypto.go#L86
func (d *HybridDecrypter) Decrypt(encryptedText string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < 2 {
		return "", fmt.Errorf("data is too short")
	}
	rsaLen := int(binary.BigEndian.Uint16(ciphertext))
	if len(ciphertext) < rsaLen+2 {
		return "", fmt.Errorf("data is too short")
	}

	rsaCiphertext := ciphertext[2 : rsaLen+2]
	aesCiphertext := ciphertext[rsaLen+2:]

	symKey, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, d.key, rsaCiphertext, nil)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(symKey)
	if err != nil {
		return "", err
	}

	aed, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Key is only used once, so zero nonce is ok.
	nonce := make([]byte, aed.NonceSize())
	text, err := aed.Open(nil, nonce, aesCiphertext, nil)
	if err != nil {
		return "", err
	}

	return string(text), nil
}
