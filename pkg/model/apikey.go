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

	"golang.org/x/crypto/bcrypt"
)

const (
	apiKeyLength = 32
)

func GenerateAPIKey(id string) (key, hash string, err error) {
	k := GenerateRandomString(apiKeyLength)
	key = fmt.Sprintf("%s.%s", id, k)

	var encoded []byte
	encoded, err = bcrypt.GenerateFromPassword([]byte(key), bcrypt.DefaultCost)
	if err != nil {
		return
	}

	hash = string(encoded)
	return
}

func ExtractAPIKeyID(key string) (string, error) {
	parts := strings.Split(key, ".")
	if len(parts) != 2 {
		return "", errors.New("malformed api key")
	}

	if parts[0] == "" || parts[1] == "" {
		return "", errors.New("malformed api key")
	}

	return parts[0], nil
}

func (k *APIKey) CompareKey(key string) error {
	if key == "" {
		return errors.New("key was empty")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(k.KeyHash), []byte(key)); err != nil {
		return fmt.Errorf("wrong api key %s: %w", key, err)
	}
	return nil
}

// RedactSensitiveData redacts sensitive data.
func (k *APIKey) RedactSensitiveData() {
	k.KeyHash = redactedMessage
}

func (k *APIKey) SetUpdatedAt(t int64) {
	k.UpdatedAt = t
}
