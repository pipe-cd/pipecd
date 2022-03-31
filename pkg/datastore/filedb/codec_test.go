// Copyright 2022 The PipeCD Authors.
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

package filedb

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pipe-cd/pipecd/pkg/datastore"
)

type fakeModel struct {
	Data      string `json:"data"`
	UpdatedAt int64  `json:"updated_at"`
}

func (fm *fakeModel) GetUpdatedAt() int64 {
	return fm.UpdatedAt
}

func (fm *fakeModel) SetUpdatedAt(t int64) {
	fm.UpdatedAt = t
}

func TestMerge(t *testing.T) {
	testcases := []struct {
		name     string
		parts    map[datastore.Shard][]byte
		expected *fakeModel
	}{
		{
			name: "should merge correctly with time ordered data",
			parts: map[datastore.Shard][]byte{
				datastore.ClientShard: []byte(`{"data":"1","updated_at":1}`),
				datastore.AgentShard:  []byte(`{"data":"1","updated_at":2}`),
				datastore.OpsShard:    []byte(`{"data":"1","updated_at":3}`),
			},
			expected: &fakeModel{
				Data:      "1",
				UpdatedAt: 3,
			},
		},
		{
			name: "should merge correctly with time non ordered data",
			parts: map[datastore.Shard][]byte{
				datastore.ClientShard: []byte(`{"data":"1","updated_at":2}`),
				datastore.AgentShard:  []byte(`{"data":"1","updated_at":3}`),
				datastore.OpsShard:    []byte(`{"data":"1","updated_at":1}`),
			},
			expected: &fakeModel{
				Data:      "1",
				UpdatedAt: 3,
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			m := &fakeModel{}
			merge(m, tc.parts)
			assert.Equal(t, tc.expected, m)
		})
	}
}
