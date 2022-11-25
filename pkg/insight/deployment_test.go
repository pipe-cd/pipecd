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

package insight

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoundDay(t *testing.T) {
	testcases := []struct {
		name      string
		timestamp int64
		expected  int64
	}{
		{
			name:      "zero",
			timestamp: 0,
			expected:  0,
		},
		{
			name:      "normal time",
			timestamp: 1668013222,
			expected:  1667952000,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := roundDay(tc.timestamp)
			assert.Equal(t, tc.expected, got)
		})
	}
}
