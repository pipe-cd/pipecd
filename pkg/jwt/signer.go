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

package jwt

import (
	"fmt"

	jwtgo "github.com/golang-jwt/jwt"
)

type Signer interface {
	Sign(claims *Claims) (string, error)
}

type signer struct {
	key    interface{}
	method jwtgo.SigningMethod
}

// NewSigner returns a new signer using SigningMethodRS256.
func NewSigner(method jwtgo.SigningMethod, keyFile string) (Signer, error) {
	key, err := readKeyFile(method, keyFile, true)
	if err != nil {
		return nil, fmt.Errorf("unable to read key file: %v", err)
	}
	return &signer{
		key:    key,
		method: method,
	}, nil
}

func (s *signer) Sign(claims *Claims) (string, error) {
	token := jwtgo.NewWithClaims(s.method, claims)
	return token.SignedString(s.key)
}
