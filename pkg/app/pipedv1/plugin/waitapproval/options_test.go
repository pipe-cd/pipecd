// Copyright 2025 The PipeCD Authors.
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

package main

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecode(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		name     string
		data     json.RawMessage
		expected waitApprovalStageOptions
		wantErr  bool
	}{
		{
			name: "valid config",
			data: json.RawMessage(`{"approvers":["user1@example.com","user2@example.com"],"minApproverNum":1}`),
			expected: waitApprovalStageOptions{
				Approvers:      []string{"user1@example.com", "user2@example.com"},
				MinApproverNum: 1,
			},
			wantErr: false,
		},
		{
			name:     "invalid config",
			data:     json.RawMessage(`invalid`),
			expected: waitApprovalStageOptions{},
			wantErr:  true,
		},
		{
			name:     "empty config",
			data:     json.RawMessage(`{}`),
			expected: waitApprovalStageOptions{},
			wantErr:  true,
		},
		{
			name:     "missing approvers",
			data:     json.RawMessage(`{"minApproverNum":1}`),
			expected: waitApprovalStageOptions{},
			wantErr:  true,
		},
		{
			name:     "minApproverNum greater than approvers",
			data:     json.RawMessage(`{"approvers":["user1@example.com"],"minApproverNum":2}`),
			expected: waitApprovalStageOptions{},
			wantErr:  true,
		},
		{
			name: "minApproverNum defaults to 1",
			data: json.RawMessage(`{"approvers":["user1@example.com","user2@example.com"]}`),
			expected: waitApprovalStageOptions{
				Approvers:      []string{"user1@example.com", "user2@example.com"},
				MinApproverNum: 1,
			},
			wantErr: false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, err := decode(tc.data)
			assert.Equal(t, tc.wantErr, err != nil)
			assert.Equal(t, tc.expected, got)
		})
	}
}
