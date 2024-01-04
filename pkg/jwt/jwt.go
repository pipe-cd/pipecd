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
	"os"
	"time"

	jwtgo "github.com/golang-jwt/jwt"

	"github.com/pipe-cd/pipecd/pkg/model"
)

const (
	// Issuer is the name of the issuer for PipeCD token.
	Issuer = "PipeCD"
	// SignedTokenKey is the name of singed token key in cookie.
	SignedTokenKey = "token"
)

// Claims extends the StandardClaims with the role to access PipeCD resources.
type Claims struct {
	jwtgo.StandardClaims
	AvatarURL string     `json:"avatarUrl,omitempty"`
	Role      model.Role `json:"role,omitempty"`
}

// NewClaims creates a new claims for a given github user.
func NewClaims(githubUserID, avatarURL string, ttl time.Duration, role model.Role) *Claims {
	now := time.Now().UTC()
	return &Claims{
		StandardClaims: jwtgo.StandardClaims{
			Subject:   githubUserID,
			Issuer:    Issuer,
			IssuedAt:  now.Unix(),
			NotBefore: now.Unix(),
			ExpiresAt: now.Add(ttl).Unix(),
			Audience:  "",
		},
		AvatarURL: avatarURL,
		Role:      role,
	}
}

func readKeyFile(method jwtgo.SigningMethod, keyFile string, isSigningKey bool) (interface{}, error) {
	data, err := os.ReadFile(keyFile)
	if err != nil {
		return nil, fmt.Errorf("unable to read key file: %v", err)
	}
	switch method.Alg() {
	case "HS256", "HS384", "HS512":
		return data, nil
	case "RS256", "RS384", "RS512":
		if isSigningKey {
			return jwtgo.ParseRSAPrivateKeyFromPEM(data)
		}
		return jwtgo.ParseRSAPublicKeyFromPEM(data)
	}
	return nil, fmt.Errorf("unexpected signing method: %v", method.Alg())
}
