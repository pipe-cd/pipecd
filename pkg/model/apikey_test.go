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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
