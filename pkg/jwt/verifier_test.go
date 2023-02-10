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

package jwt

import (
	"fmt"
	"strings"
	"testing"
	"time"

	jwtgo "github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pipe-cd/pipecd/pkg/model"
)

func TestVerify(t *testing.T) {
	now := time.Now()

	testcases := []struct {
		name      string
		claims    *Claims
		fail      bool
		errPrefix string
	}{
		{
			name: "ok",
			claims: NewClaims("user-1", "avatar-url", time.Hour, model.Role{
				ProjectId: "project-1",
			}),
			fail: false,
		},
		{
			name: "wrong issuer",
			claims: &Claims{
				StandardClaims: jwtgo.StandardClaims{
					Issuer:    "test-issuer",
					IssuedAt:  now.Unix(),
					NotBefore: now.Unix(),
					ExpiresAt: now.Add(time.Hour).Unix(),
				},
			},
			fail:      true,
			errPrefix: "invalid issuer",
		},
		{
			name: "expired",
			claims: &Claims{
				StandardClaims: jwtgo.StandardClaims{
					Issuer:    Issuer,
					IssuedAt:  now.Add(-time.Hour).Unix(),
					NotBefore: now.Add(-time.Hour).Unix(),
					ExpiresAt: now.Add(-time.Minute).Unix(),
				},
			},
			fail:      true,
			errPrefix: "unable to parse token: token is expired",
		},
		{
			name: "missing issueAt",
			claims: &Claims{
				StandardClaims: jwtgo.StandardClaims{
					Issuer:    Issuer,
					NotBefore: now.Unix(),
					ExpiresAt: now.Add(time.Hour).Unix(),
				},
			},
			fail:      true,
			errPrefix: "missing issuedAt",
		},
		{
			name: "missing expiresAt",
			claims: &Claims{
				StandardClaims: jwtgo.StandardClaims{
					Issuer:    Issuer,
					IssuedAt:  now.Unix(),
					NotBefore: now.Unix(),
				},
			},
			fail:      true,
			errPrefix: "missing expiresAt",
		},
		{
			name: "missing notBefore",
			claims: &Claims{
				StandardClaims: jwtgo.StandardClaims{
					Issuer:    Issuer,
					IssuedAt:  now.Unix(),
					ExpiresAt: now.Add(time.Hour).Unix(),
				},
			},
			fail:      true,
			errPrefix: "missing notBefore",
		},
	}

	testFunc := func(s Signer, v Verifier) {
		for _, tc := range testcases {
			t.Run(tc.name, func(t *testing.T) {
				token, err := s.Sign(tc.claims)
				require.NoError(t, err)
				require.True(t, len(token) > 0)

				got, err := v.Verify(token)
				if tc.fail {
					require.Error(t, err)
					assert.Nil(t, got)
					if tc.errPrefix != "" && !strings.HasPrefix(err.Error(), tc.errPrefix) {
						assert.Fail(t, fmt.Sprintf("unexpected error prefix, expected: %s, got: %s", tc.errPrefix, err.Error()))
					}
				} else {
					assert.NoError(t, err)
					assert.Equal(t, tc.claims, got)
				}
			})
		}
	}

	rsS, err := NewSigner(jwtgo.SigningMethodRS256, "testdata/private.key")
	require.NoError(t, err)
	require.NotNil(t, rsS)

	rsV, err := NewVerifier(jwtgo.SigningMethodRS256, "testdata/public.key")
	require.NoError(t, err)
	require.NotNil(t, rsV)

	testFunc(rsS, rsV)

	hsS, err := NewSigner(jwtgo.SigningMethodHS256, "testdata/private.key")
	require.NoError(t, err)
	require.NotNil(t, hsS)

	hsV, err := NewVerifier(jwtgo.SigningMethodHS256, "testdata/private.key")
	require.NoError(t, err)
	require.NotNil(t, hsV)

	testFunc(hsS, hsV)

	c := NewClaims("user", "avatar-url", time.Hour, model.Role{
		ProjectId: "project",
	})

	token, err := rsS.Sign(c)
	require.NoError(t, err)
	require.True(t, len(token) > 0)
	got, err := hsV.Verify(token)
	require.Error(t, err)
	require.Nil(t, got)

	token, err = hsS.Sign(c)
	require.NoError(t, err)
	require.True(t, len(token) > 0)
	got, err = rsV.Verify(token)
	require.Error(t, err)
	require.Nil(t, got)
}
