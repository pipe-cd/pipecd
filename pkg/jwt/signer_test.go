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
	"testing"
	"time"

	jwtgo "github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/require"

	"github.com/pipe-cd/pipecd/pkg/model"
)

func TestSign(t *testing.T) {
	claims := NewClaims("user-1", "avatar-url", time.Hour, model.Role{
		ProjectId: "project-1",
	})

	s, err := NewSigner(jwtgo.SigningMethodRS256, "testdata/private.key")
	require.NoError(t, err)
	require.NotNil(t, s)

	token, err := s.Sign(claims)
	require.NoError(t, err)
	require.True(t, len(token) > 0)

	s, err = NewSigner(jwtgo.SigningMethodHS256, "testdata/private.key")
	require.NoError(t, err)
	require.NotNil(t, s)

	token, err = s.Sign(claims)
	require.NoError(t, err)
	require.True(t, len(token) > 0)
}
