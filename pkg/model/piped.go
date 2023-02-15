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

package model

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const (
	pipedMaxKeyNum  = 2
	pipedKeyLength  = 50
	redactedMessage = "redacted"
)

type SecretManagementType string

const (
	SecretManagementTypeNone    SecretManagementType = "NONE"
	SecretManagementTypeKeyPair SecretManagementType = "KEY_PAIR"
	SecretManagementTypeGCPKMS  SecretManagementType = "GCP_KMS"
	SecretManagementTypeAWSKMS  SecretManagementType = "AWS_KMS"
)

func (t SecretManagementType) String() string {
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
	if len(p.Keys) == 0 {
		return errors.New("piped does not contain any key")
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
// A piped can hold a maximum of "pipedMaxKeyNum" keys.
func (p *Piped) AddKey(hash, creator string, createdAt time.Time) error {
	if len(p.Keys) == pipedMaxKeyNum {
		return fmt.Errorf("number of keys for each piped must be less than or equal to %d, you may need to delete the old keys before adding a new one", pipedMaxKeyNum)
	}

	k := &PipedKey{
		Hash:      hash,
		Creator:   creator,
		CreatedAt: createdAt.Unix(),
	}

	keys := make([]*PipedKey, 0, len(p.Keys)+1)
	keys = append(keys, k)
	keys = append(keys, p.Keys...)
	p.Keys = keys

	return nil
}

// DeleteOldPipedKeys removes all old keys to keep only the latest one.
func (p *Piped) DeleteOldPipedKeys() {
	if len(p.Keys) < 2 {
		return
	}

	var latest *PipedKey
	for i := range p.Keys {
		if latest == nil || latest.CreatedAt < p.Keys[i].CreatedAt {
			latest = p.Keys[i]
		}
	}
	p.Keys = []*PipedKey{latest}
}

func (p *Piped) RedactSensitiveData() {
	for i := range p.Keys {
		p.Keys[i].Hash = redactedMessage
	}
}

func (p *Piped) SetUpdatedAt(t int64) {
	p.UpdatedAt = t
}

func MakePipedURL(baseURL, pipedID string) string {
	return fmt.Sprintf("%s/settings/piped", strings.TrimSuffix(baseURL, "/"))
}
