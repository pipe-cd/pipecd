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

package model

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestGenerateAPIKey(t *testing.T) {
	id := "test-id"
	key, hash, err := GenerateAPIKey(id)
	require.NoError(t, err)
	require.True(t, len(key) > 0)
	require.True(t, len(hash) > 0)

	parsedID, err := ExtractAPIKeyID(key)
	require.NoError(t, err)
	assert.Equal(t, id, parsedID)

	apiKey := &APIKey{
		Id:      id,
		KeyHash: hash,
	}

	err = apiKey.CompareKey(key)
	assert.NoError(t, err)
}

func TestAPIKeyRedactSensitiveData(t *testing.T) {
	apiKey := &APIKey{
		Id:      "id",
		KeyHash: "hash",
	}
	apiKey.RedactSensitiveData()
	assert.Equal(t, apiKey.KeyHash, "redacted")
}

func TestExtractAPIKeyID(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    string
		expectedErr error
	}{
		{
			name:        "Valid API Key",
			input:       "abc.def",
			expected:    "abc",
			expectedErr: nil,
		},
		{
			name:        "Empty Key",
			input:       "",
			expected:    "",
			expectedErr: errors.New("malformed api key"),
		},
		{
			name:        "Only one part",
			input:       "abc",
			expected:    "",
			expectedErr: errors.New("malformed api key"),
		},
		{
			name:        "Multiple periods",
			input:       "abc.def.ghi",
			expected:    "",
			expectedErr: errors.New("malformed api key"),
		},
		{
			name:        "Empty first part",
			input:       ".def",
			expected:    "",
			expectedErr: errors.New("malformed api key"),
		},
		{
			name:        "Empty second part",
			input:       "abc.",
			expected:    "",
			expectedErr: errors.New("malformed api key"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := ExtractAPIKeyID(tt.input)
			assert.Equal(t, tt.expected, actual)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestCompareKey(t *testing.T) {
	password := "my-secret-key"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	hashedStr := string(hashedPassword)
	apiKey := &APIKey{KeyHash: hashedStr}

	tests := []struct {
		name          string
		input         string
		expectedError error
	}{
		{
			name:          "Valid Key",
			input:         password,
			expectedError: nil,
		},
		{
			name:          "Invalid Key",
			input:         "wrong-key",
			expectedError: fmt.Errorf("wrong api key wrong-key: %w", bcrypt.ErrMismatchedHashAndPassword),
		},
		{
			name:          "Empty Key",
			input:         "",
			expectedError: errors.New("key was empty"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := apiKey.CompareKey(tt.input)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}
