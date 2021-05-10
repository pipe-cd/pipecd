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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGeneratePipedKey(t *testing.T) {
	key, hash, err := GeneratePipedKey()
	assert.NoError(t, err)
	assert.True(t, len(key) > 0)
	assert.True(t, len(hash) > 0)

	p := &Piped{}
	p.AddKey(hash, "user", time.Now())

	err = p.CheckKey(key)
	assert.NoError(t, err)

	err = p.CheckKey("invalid")
	assert.Error(t, err)
}

func TestAddKey(t *testing.T) {
	p := &Piped{}
	require.Equal(t, 0, len(p.Keys))

	now := time.Now()

	p.AddKey("hash-1", "user-1", now)
	require.Equal(t, []*PipedKey{
		{
			Hash:      "hash-1",
			Creator:   "user-1",
			CreatedAt: now.Unix(),
		},
	}, p.Keys)

	p.AddKey("hash-2", "user-1", now.Add(time.Second))
	require.Equal(t, []*PipedKey{
		{
			Hash:      "hash-2",
			Creator:   "user-1",
			CreatedAt: now.Unix() + 1,
		},
		{
			Hash:      "hash-1",
			Creator:   "user-1",
			CreatedAt: now.Unix(),
		},
	}, p.Keys)

	p.AddKey("hash-3", "user-3", now.Add(2*time.Second))
	require.Equal(t, []*PipedKey{
		{
			Hash:      "hash-3",
			Creator:   "user-3",
			CreatedAt: now.Unix() + 2,
		},
		{
			Hash:      "hash-2",
			Creator:   "user-1",
			CreatedAt: now.Unix() + 1,
		},
		{
			Hash:      "hash-1",
			Creator:   "user-1",
			CreatedAt: now.Unix(),
		},
	}, p.Keys)

	p.AddKey("hash-4", "user-1", now.Add(3*time.Second))
	require.Equal(t, []*PipedKey{
		{
			Hash:      "hash-4",
			Creator:   "user-1",
			CreatedAt: now.Unix() + 3,
		},
		{
			Hash:      "hash-3",
			Creator:   "user-3",
			CreatedAt: now.Unix() + 2,
		},
		{
			Hash:      "hash-2",
			Creator:   "user-1",
			CreatedAt: now.Unix() + 1,
		},
	}, p.Keys)
}

func TestPipedRedactSensitiveData(t *testing.T) {
	testcases := []struct {
		name     string
		piped    Piped
		expected Piped
	}{
		{
			name: "contains both KeyHash and Keys",
			piped: Piped{
				KeyHash: "hash-key",
				Keys: []*PipedKey{
					{
						Hash:      "hash-1",
						Creator:   "user-1",
						CreatedAt: 1,
					},
					{
						Hash:      "hash-2",
						Creator:   "user-2",
						CreatedAt: 2,
					},
					{
						Hash:      "hash-3",
						Creator:   "user-3",
						CreatedAt: 3,
					},
				},
			},
			expected: Piped{
				KeyHash: "redacted",
				Keys: []*PipedKey{
					{
						Hash:      "redacted",
						Creator:   "user-1",
						CreatedAt: 1,
					},
					{
						Hash:      "redacted",
						Creator:   "user-2",
						CreatedAt: 2,
					},
					{
						Hash:      "redacted",
						Creator:   "user-3",
						CreatedAt: 3,
					},
				},
			},
		},
	}

	for _, tc := range testcases {
		tc.piped.RedactSensitiveData()
		assert.Equal(t, tc.expected, tc.piped)
	}
}
