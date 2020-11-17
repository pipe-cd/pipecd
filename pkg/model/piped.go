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

package model

import (
	"math/rand"
	"sort"
	"time"
	"unsafe"

	"golang.org/x/crypto/bcrypt"
)

const (
	pipedMaxKeyNum  = 3
	pipedKeyLength  = 50
	redactedMessage = "redacted"
)

type SealedSecretManagementType string

const (
	SealedSecretManagementNone       SealedSecretManagementType = "NONE"
	SealedSecretManagementSealingKey SealedSecretManagementType = "SEALING_KEY"
	SealedSecretManagementGCPKMS     SealedSecretManagementType = "GCP_KMS"
	SealedSecretManagementAWSKMS     SealedSecretManagementType = "AWS_KMS"
)

func (t SealedSecretManagementType) String() string {
	return string(t)
}

// GeneratePipedKey generates a new key for piped.
// This returns raw key value for used by piped and
// a hash value of the key for storing in datastore.
func GeneratePipedKey() (key, hash string, err error) {
	var encoded []byte
	key = GenerateRandomString(pipedKeyLength)
	encoded, err = bcrypt.GenerateFromPassword([]byte(key), bcrypt.DefaultCost)
	if err != nil {
		return
	}
	hash = string(encoded)
	return
}

// CheckKey checks if the give key is one of the stored keys.
func (p *Piped) CheckKey(key string) (err error) {
	// The KeyHash field was deprecated.
	// And this block will be removed in the future.
	if p.KeyHash != "" {
		err = bcrypt.CompareHashAndPassword([]byte(p.KeyHash), []byte(key))
		if err == nil {
			return nil
		}
	}

	for _, k := range p.Keys {
		err = bcrypt.CompareHashAndPassword([]byte(k.Hash), []byte(key))
		if err == nil {
			return nil
		}
	}
	return
}

// AddKey adds a new key to the list.
// If the piped is already having "pipedMaxKeyNum" number of keys
// the oldest one will be removed before adding.
func (p *Piped) AddKey(hash, creator string, createdAt time.Time) {
	key := &PipedKey{
		Hash:      hash,
		Creator:   creator,
		CreatedAt: createdAt.Unix(),
	}
	if len(p.Keys) == 0 {
		p.Keys = []*PipedKey{key}
		return
	}

	sort.Slice(p.Keys, func(i, j int) bool {
		return p.Keys[i].CreatedAt > p.Keys[j].CreatedAt
	})
	if len(p.Keys) >= pipedMaxKeyNum {
		p.Keys = p.Keys[:pipedMaxKeyNum-1]
	}
	keys := make([]*PipedKey, 0, len(p.Keys)+1)
	keys = append(keys, key)
	keys = append(keys, p.Keys...)

	p.Keys = keys
}

var randomSrc = rand.NewSource(time.Now().UnixNano())

const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // Number of letter indices fitting in 63 bits
	letterBytes   = "abcdefghijklmnopqrstuvwxyz0123456789"
)

// GenerateRandomString makes a random string with the given length.
// This implementation was referenced from the link below.
// https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
func GenerateRandomString(n int) string {
	b := make([]byte, n)

	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters.
	for i, cache, remain := n-1, randomSrc.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = randomSrc.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}

func (p *Piped) RedactSensitiveData() {
	p.KeyHash = redactedMessage
}
