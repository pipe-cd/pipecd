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
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testStageOptions struct {
	Number int    `json:"number"`
	Text   string `json:"text"`
}

func (t testStageOptions) Validate() error {
	if t.Number > 0 {
		return nil
	} else {
		return errors.New("number must be greater than 0")
	}
}

func TestDecodeStageOptionsYAML(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		title     string
		data      []byte
		expected  testStageOptions
		expectErr bool
	}{
		{
			title: "valid data",
			data: []byte(`
number: 123
text: "test"
`),
			expected: testStageOptions{
				Number: 123,
				Text:   "test",
			},
			expectErr: false,
		},
		{
			title:     "invalid format",
			data:      []byte(`invalid`),
			expected:  testStageOptions{},
			expectErr: true,
		},
		{
			title:     "validation failed",
			data:      []byte(`number: -1`),
			expected:  testStageOptions{},
			expectErr: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.title, func(t *testing.T) {
			t.Parallel()
			got, err := DecodeStageOptionsYAML[testStageOptions](tc.data)
			assert.Equal(t, tc.expectErr, err != nil)
			if err == nil {
				assert.Equal(t, tc.expected, *got)
			}
		})
	}
}
