// Copyright 2023 The PipeCD Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package prompt

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunStringSlice(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name          string
		in            PromptInput
		str           string // user's input
		expectedValue []string
		expectedErr   bool
	}{
		{
			name: "valid string slice",
			in: PromptInput{
				Message:       "anyPrompt",
				TargetPointer: new([]string),
				Required:      false,
			},
			str:           "foo bar\n",
			expectedValue: []string{"foo", "bar"},
			expectedErr:   false,
		},
		{
			name: "empty but not required",
			in: PromptInput{
				Message:       "anyPrompt",
				TargetPointer: new([]string),
				Required:      false,
			},
			str:           "\n",
			expectedValue: nil,
			expectedErr:   false,
		},
		{
			name: "missing required",
			in: PromptInput{
				Message:       "anyPrompt",
				TargetPointer: new([]string),
				Required:      true,
			},
			str:           "\n",
			expectedValue: nil,
			expectedErr:   true,
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			strReader := strings.NewReader(tc.str)
			p := NewPrompt(strReader)
			e := p.run(tc.in)
			assert.Equal(t, tc.expectedErr, e != nil)
			assert.Equal(t, tc.expectedValue, *tc.in.TargetPointer.(*[]string))
		})
	}
}

func TestRunString(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name          string
		in            PromptInput
		str           string // user's input
		expectedValue string
		expectedErr   bool
	}{
		{
			name: "valid string",
			in: PromptInput{
				Message:       "anyPrompt",
				TargetPointer: new(string),
				Required:      false,
			},
			str:           "foo\n",
			expectedValue: "foo",
			expectedErr:   false,
		},
		{
			name: "empty but not required",
			in: PromptInput{
				Message:       "anyPrompt",
				TargetPointer: new(string),
				Required:      false,
			},
			str:           "\n",
			expectedValue: "",
			expectedErr:   false,
		},
		{
			name: "missing required",
			in: PromptInput{
				Message:       "anyPrompt",
				TargetPointer: new(string),
				Required:      true,
			},
			str:           "\n",
			expectedValue: "",
			expectedErr:   true,
		},
		{
			name: "two many arguments",
			in: PromptInput{
				Message:       "anyPrompt",
				TargetPointer: new(string),
				Required:      false,
			},
			str:           "foo bar\n",
			expectedValue: "",
			expectedErr:   true,
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			strReader := strings.NewReader(tc.str)
			p := NewPrompt(strReader)
			e := p.run(tc.in)
			assert.Equal(t, tc.expectedErr, e != nil)
			assert.Equal(t, tc.expectedValue, *tc.in.TargetPointer.(*string))
		})
	}
}

func TestRunInt(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name          string
		in            PromptInput
		str           string // user's input
		expectedValue int
		expectedErr   bool
	}{
		{
			name: "valid int",
			in: PromptInput{
				Message:       "anyPrompt",
				TargetPointer: new(int),
				Required:      false,
			},
			str:           "123\n",
			expectedValue: 123,
			expectedErr:   false,
		},
		{
			name: "invalid int",
			in: PromptInput{
				Message:       "anyPrompt",
				TargetPointer: new(int),
				Required:      false,
			},
			str:           "abc\n",
			expectedValue: 0,
			expectedErr:   true,
		},
		{
			name: "empty but not required",
			in: PromptInput{
				Message:       "anyPrompt",
				TargetPointer: new(int),
				Required:      false,
			},
			str:           "\n",
			expectedValue: 0,
			expectedErr:   false,
		},
		{
			name: "missing required",
			in: PromptInput{
				Message:       "anyPrompt",
				TargetPointer: new(int),
				Required:      true,
			},
			str:           "\n",
			expectedValue: 0,
			expectedErr:   true,
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			strReader := strings.NewReader(tc.str)
			p := NewPrompt(strReader)
			e := p.run(tc.in)
			assert.Equal(t, tc.expectedErr, e != nil)
			assert.Equal(t, tc.expectedValue, *tc.in.TargetPointer.(*int))
		})
	}
}

func TestRunBool(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name          string
		in            PromptInput
		str           string // user's input
		expectedValue bool
		expectedErr   bool
	}{
		{
			name: "valid bool",
			in: PromptInput{
				Message:       "anyPrompt",
				TargetPointer: new(bool),
				Required:      false,
			},
			str:           "true\n",
			expectedValue: true,
			expectedErr:   false,
		},
		{
			name: "invalid bool",
			in: PromptInput{
				Message:       "anyPrompt",
				TargetPointer: new(bool),
				Required:      false,
			},
			str:           "abc\n",
			expectedValue: false,
			expectedErr:   true,
		},
		{
			name: "empty but not required",
			in: PromptInput{
				Message:       "anyPrompt",
				TargetPointer: new(bool),
				Required:      false,
			},
			str:           "\n",
			expectedValue: false,
			expectedErr:   false,
		},
		{
			name: "missing required",
			in: PromptInput{
				Message:       "anyPrompt",
				TargetPointer: new(bool),
				Required:      true,
			},
			str:           "\n",
			expectedValue: false,
			expectedErr:   true,
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			strReader := strings.NewReader(tc.str)
			p := NewPrompt(strReader)
			e := p.run(tc.in)
			assert.Equal(t, tc.expectedErr, e != nil)
			assert.Equal(t, tc.expectedValue, *tc.in.TargetPointer.(*bool))
		})
	}
}
