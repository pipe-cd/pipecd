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
	"strings"
	"testing"
	"time"

	jwtgo "github.com/golang-jwt/jwt/v5"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
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
		errString string
	}{
		{
			name: "ok",
			claims: NewClaims("user-1", "avatar-url", time.Hour, &model.Role{
				ProjectId: "project-1",
			}),
			fail: false,
		},
		{
			name: "wrong issuer",
			claims: &Claims{
				RegisteredClaims: jwtgo.RegisteredClaims{
					Issuer:    "test-issuer",
					IssuedAt:  jwtgo.NewNumericDate(now),
					NotBefore: jwtgo.NewNumericDate(now),
					ExpiresAt: jwtgo.NewNumericDate(now.Add(time.Hour)),
				},
			},
			fail:      true,
			errString: "invalid issuer",
		},
		{
			name: "expired",
			claims: &Claims{
				RegisteredClaims: jwtgo.RegisteredClaims{
					Issuer:    Issuer,
					IssuedAt:  jwtgo.NewNumericDate(now.Add(-time.Hour)),
					NotBefore: jwtgo.NewNumericDate(now.Add(-time.Hour)),
					ExpiresAt: jwtgo.NewNumericDate(now.Add(-time.Minute)),
				},
			},
			fail:      true,
			errString: "unable to parse token: token has invalid claims: token is expired",
		},
		{
			name: "missing expiresAt",
			claims: &Claims{
				RegisteredClaims: jwtgo.RegisteredClaims{
					Issuer:    Issuer,
					IssuedAt:  jwtgo.NewNumericDate(now),
					NotBefore: jwtgo.NewNumericDate(now),
				},
			},
			fail:      true,
			errString: "unable to parse token: token has invalid claims: token is missing required claim: exp claim is required",
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
					if tc.errString != "" && !strings.Contains(err.Error(), tc.errString) {
						assert.Fail(t, fmt.Sprintf("unexpected error, expected: %s, got: %s", tc.errString, err.Error()))
					}
				} else {
					assert.NoError(t, err)
					if !cmp.Equal(tc.claims, got, cmpopts.IgnoreUnexported(jwtgo.RegisteredClaims{}, model.Role{})) {
						assert.Fail(t, fmt.Sprintf("unexpected claims, expected: %v, got: %v", tc.claims, got))
					}
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

	c := NewClaims("user", "avatar-url", time.Hour, &model.Role{
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
