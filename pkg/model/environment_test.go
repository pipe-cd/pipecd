// Copyright 2021 The PipeCD Authors.
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
)

func TestMakeEnvironmentURL(t *testing.T) {
	testcases := []struct {
		name          string
		baseURL       string
		environmentID string
		expected      string
	}{
		{
			name:          "baseURL has no suffix",
			baseURL:       "https://pipecd.dev",
			environmentID: "env-1",
			expected:      "https://pipecd.dev/settings/environment",
		},
		{
			name:          "baseURL suffixed by /",
			baseURL:       "https://pipecd.dev/",
			environmentID: "env-1",
			expected:      "https://pipecd.dev/settings/environment",
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := MakeEnvironmentURL(tc.baseURL, tc.environmentID)
			assert.Equal(t, tc.expected, got)
		})
	}
}
