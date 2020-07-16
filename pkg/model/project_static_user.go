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
	"crypto/subtle"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// Auth confirms username and password.
func (p *ProjectStaticUser) Auth(username, password string) error {
	if username == "" || subtle.ConstantTimeCompare([]byte(p.Username), []byte(username)) != 1 {
		return fmt.Errorf("wrong input data")
	}
	if password == "" || bcrypt.CompareHashAndPassword([]byte(p.PasswordHash), []byte(password)) != nil {
		return fmt.Errorf("wrong input data")
	}
	return nil
}
