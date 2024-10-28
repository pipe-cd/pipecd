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

package diff

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindOne(t *testing.T) {
	nodes := []Node{
		{PathString: "spec.template.spec"},
	}

	testcases := []struct {
		name           string
		nodes          Nodes
		query          string
		expected       *Node
		exepectedError error
	}{
		{
			name:           "nil list",
			query:          ".+",
			exepectedError: ErrNotFound,
		},
		{
			name:           "not found",
			nodes:          nodes,
			query:          `spec\.not-found\..+`,
			exepectedError: ErrNotFound,
		},
		{
			name:     "found",
			nodes:    nodes,
			query:    `spec\.template\..+`,
			expected: &nodes[0],
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			n, err := tc.nodes.FindOne(tc.query)
			assert.Equal(t, tc.expected, n)
			assert.Equal(t, tc.exepectedError, err)
		})
	}
}

func TestFind(t *testing.T) {
	nodes := []Node{
		{PathString: "spec.replicas"},
		{PathString: "spec.template.spec.containers.0.image"},
		{PathString: "spec.template.spec.containers.1.image"},
	}

	testcases := []struct {
		name     string
		nodes    Nodes
		query    string
		expected []Node
	}{
		{
			name:     "nil list",
			query:    ".+",
			expected: []Node{},
		},
		{
			name:     "not found",
			nodes:    nodes,
			query:    `spec\.not-found\..+`,
			expected: []Node{},
		},
		{
			name:     "found two nodes",
			nodes:    nodes,
			query:    `spec\.template\.spec\.containers\.\d+.image$`,
			expected: []Node{nodes[1], nodes[2]},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ns, err := tc.nodes.Find(tc.query)
			assert.Equal(t, Nodes(tc.expected), ns)
			assert.NoError(t, err)
		})
	}
}
