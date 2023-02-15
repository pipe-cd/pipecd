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
	"errors"
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
}

func TestPipedCheckKey(t *testing.T) {
	key, hash, err := GeneratePipedKey()
	require.NoError(t, err)

	p := &Piped{}
	p.AddKey(hash, "user", time.Now())

	err = p.CheckKey(key)
	assert.NoError(t, err)

	err = p.CheckKey("invalid")
	assert.Error(t, err)

	noKeyPiped := &Piped{}
	err = noKeyPiped.CheckKey(key)
	assert.Equal(t, errors.New("piped does not contain any key"), err)
}

func TestAddKey(t *testing.T) {
	p := &Piped{}
	require.Equal(t, 0, len(p.Keys))

	now := time.Now()

	err := p.AddKey("hash-1", "user-1", now)
	assert.NoError(t, err)
	require.Equal(t, []*PipedKey{
		{
			Hash:      "hash-1",
			Creator:   "user-1",
			CreatedAt: now.Unix(),
		},
	}, p.Keys)

	err = p.AddKey("hash-2", "user-1", now.Add(time.Second))
	assert.NoError(t, err)
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

	err = p.AddKey("hash-3", "user-3", now.Add(2*time.Second))
	assert.Equal(t, errors.New("number of keys for each piped must be less than or equal to 2, you may need to delete the old keys before adding a new one"), err)
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
}

func TestPipedDeleteOldPipedKeys(t *testing.T) {
	testcases := []struct {
		name     string
		piped    Piped
		expected Piped
	}{
		{
			name: "no key",
			piped: Piped{
				Keys: []*PipedKey{},
			},
			expected: Piped{
				Keys: []*PipedKey{},
			},
		},
		{
			name: "has 1 key",
			piped: Piped{
				Keys: []*PipedKey{
					{
						Hash:      "hash-1",
						Creator:   "user-1",
						CreatedAt: 1,
					},
				},
			},
			expected: Piped{
				Keys: []*PipedKey{
					{
						Hash:      "hash-1",
						Creator:   "user-1",
						CreatedAt: 1,
					},
				},
			},
		},
		{
			name: "has multiple keys",
			piped: Piped{
				Keys: []*PipedKey{
					{
						Hash:      "hash-1",
						Creator:   "user-1",
						CreatedAt: 1,
					},
					{
						Hash:      "hash-3",
						Creator:   "user-3",
						CreatedAt: 3,
					},
					{
						Hash:      "hash-2",
						Creator:   "user-2",
						CreatedAt: 2,
					},
				},
			},
			expected: Piped{
				Keys: []*PipedKey{
					{
						Hash:      "hash-3",
						Creator:   "user-3",
						CreatedAt: 3,
					},
				},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			tc.piped.DeleteOldPipedKeys()
			assert.Equal(t, tc.expected, tc.piped)
		})
	}
}

func TestPipedRedactSensitiveData(t *testing.T) {
	testcases := []struct {
		name     string
		piped    Piped
		expected Piped
	}{
		{
			name: "contains multiple Keys",
			piped: Piped{
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
		t.Run(tc.name, func(t *testing.T) {
			tc.piped.RedactSensitiveData()
			assert.Equal(t, tc.expected, tc.piped)
		})
	}
}

func TestMakePipedURL(t *testing.T) {
	testcases := []struct {
		name     string
		baseURL  string
		pipedID  string
		expected string
	}{
		{
			name:     "baseURL has no suffix",
			baseURL:  "https://pipecd.dev",
			pipedID:  "piped-id",
			expected: "https://pipecd.dev/settings/piped",
		},
		{
			name:     "baseURL suffixed by /",
			baseURL:  "https://pipecd.dev/",
			pipedID:  "piped-id",
			expected: "https://pipecd.dev/settings/piped",
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := MakePipedURL(tc.baseURL, tc.pipedID)
			assert.Equal(t, tc.expected, got)
		})
	}
}
