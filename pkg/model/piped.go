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
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// GeneratePipedKey generates a new key for piped.
// This returns raw key value for used by piped and
// a hash value of the key for storing in datastore.
func GeneratePipedKey() (key, hash string, err error) {
	var encoded []byte
	key = uuid.New().String()
	encoded, err = bcrypt.GenerateFromPassword([]byte(key), bcrypt.DefaultCost)
	if err != nil {
		return
	}
	hash = string(encoded)
	return
}

// CompareKey compares plaintext key with its hashed value storing in piped.
func (p *Piped) CompareKey(key string) error {
	return bcrypt.CompareHashAndPassword([]byte(p.KeyHash), []byte(key))
}
