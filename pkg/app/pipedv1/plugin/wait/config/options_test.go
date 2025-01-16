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

package config

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	config "github.com/pipe-cd/pipecd/pkg/configv1"
)

func TestDecode(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		name     string
		data     json.RawMessage
		expected WaitStageOptions
		wantErr  bool
	}{
		{
			name: "valid config",
			data: json.RawMessage(`{"duration":"1m"}`),
			expected: WaitStageOptions{
				Duration: config.Duration(1 * time.Minute),
			},
			wantErr: false,
		},
		{
			name:     "invalid config",
			data:     json.RawMessage(`invalid`),
			expected: WaitStageOptions{},
			wantErr:  true,
		},
		{
			name:     "empty config",
			data:     json.RawMessage(`{}`),
			expected: WaitStageOptions{},
			wantErr:  true,
		},
		{
			name: "negative duration",
			data: json.RawMessage(`{
				"duration":"-1m"
			}`),
			expected: WaitStageOptions{},
			wantErr:  true,
		},
		{
			name: "zero duration",
			data: json.RawMessage(`{
				"duration":"0s"
			}`),
			expected: WaitStageOptions{},
			wantErr:  true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, err := Decode(tc.data)
			assert.Equal(t, tc.wantErr, err != nil)
			assert.Equal(t, tc.expected, got)
		})
	}
}
