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
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunString(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name          string
		in            Input
		str           string // user's input
		expectedValue string
		expectedErr   error
	}{
		{
			name: "valid string",
			in: Input{
				Message:       "anyPrompt",
				TargetPointer: new(string),
				Required:      false,
			},
			str:           "foo\n",
			expectedValue: "foo",
			expectedErr:   nil,
		},
		{
			name: "empty but not required",
			in: Input{
				Message:       "anyPrompt",
				TargetPointer: new(string),
				Required:      false,
			},
			str:           "\n",
			expectedValue: "",
			expectedErr:   nil,
		},
		{
			name: "missing required",
			in: Input{
				Message:       "anyPrompt",
				TargetPointer: new(string),
				Required:      true,
			},
			str:           "\n",
			expectedValue: "",
			expectedErr:   fmt.Errorf("this field is required"),
		},
		{
			name: "two many arguments",
			in: Input{
				Message:       "anyPrompt",
				TargetPointer: new(string),
				Required:      false,
			},
			str:           "foo bar\n",
			expectedValue: "",
			expectedErr:   fmt.Errorf("too many arguments"),
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			strReader := strings.NewReader(tc.str)
			p := NewPrompt(strReader)
			err := p.RunOne(tc.in)
			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expectedValue, *tc.in.TargetPointer.(*string))
		})
	}
}

func TestRunStringSlice(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name          string
		in            Input
		str           string // user's input
		expectedValue []string
		expectedErr   error
	}{
		{
			name: "valid string slice",
			in: Input{
				Message:       "anyPrompt",
				TargetPointer: new([]string),
				Required:      false,
			},
			str:           "foo bar\n",
			expectedValue: []string{"foo", "bar"},
			expectedErr:   nil,
		},
		{
			name: "empty but not required",
			in: Input{
				Message:       "anyPrompt",
				TargetPointer: new([]string),
				Required:      false,
			},
			str:           "\n",
			expectedValue: nil,
			expectedErr:   nil,
		},
		{
			name: "missing required",
			in: Input{
				Message:       "anyPrompt",
				TargetPointer: new([]string),
				Required:      true,
			},
			str:           "\n",
			expectedValue: nil,
			expectedErr:   fmt.Errorf("this field is required"),
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			strReader := strings.NewReader(tc.str)
			p := NewPrompt(strReader)
			err := p.RunOne(tc.in)
			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expectedValue, *tc.in.TargetPointer.(*[]string))
		})
	}
}

func TestRunInt(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name          string
		in            Input
		str           string // user's input
		expectedValue int
		expectedErr   error
	}{
		{
			name: "valid int",
			in: Input{
				Message:       "anyPrompt",
				TargetPointer: new(int),
				Required:      false,
			},
			str:           "123\n",
			expectedValue: 123,
			expectedErr:   nil,
		},
		{
			name: "invalid int",
			in: Input{
				Message:       "anyPrompt",
				TargetPointer: new(int),
				Required:      false,
			},
			str:           "abc\n",
			expectedValue: 0,
			expectedErr:   fmt.Errorf("this field should be an int value"),
		},
		{
			name: "empty but not required",
			in: Input{
				Message:       "anyPrompt",
				TargetPointer: new(int),
				Required:      false,
			},
			str:           "\n",
			expectedValue: 0,
			expectedErr:   nil,
		},
		{
			name: "missing required",
			in: Input{
				Message:       "anyPrompt",
				TargetPointer: new(int),
				Required:      true,
			},
			str:           "\n",
			expectedValue: 0,
			expectedErr:   fmt.Errorf("this field is required"),
		},
		{
			name: "too many arguments",
			in: Input{
				Message:       "anyPrompt",
				TargetPointer: new(int),
				Required:      true,
			},
			str:           "12 34\n",
			expectedValue: 0,
			expectedErr:   fmt.Errorf("too many arguments"),
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			strReader := strings.NewReader(tc.str)
			p := NewPrompt(strReader)
			err := p.RunOne(tc.in)
			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expectedValue, *tc.in.TargetPointer.(*int))
		})
	}
}

func TestRunBool(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name          string
		in            Input
		str           string // user's input
		expectedValue bool
		expectedErr   error
	}{
		{
			name: "valid bool",
			in: Input{
				Message:       "anyPrompt",
				TargetPointer: new(bool),
				Required:      false,
			},
			str:           "true\n",
			expectedValue: true,
			expectedErr:   nil,
		},
		{
			name: "y means true",
			in: Input{
				Message:       "anyPrompt",
				TargetPointer: new(bool),
				Required:      false,
			},
			str:           "y\n",
			expectedValue: true,
			expectedErr:   nil,
		},
		{
			name: "n means false",
			in: Input{
				Message:       "anyPrompt",
				TargetPointer: new(bool),
				Required:      false,
			},
			str:           "n\n",
			expectedValue: false,
			expectedErr:   nil,
		},
		{
			name: "invalid bool",
			in: Input{
				Message:       "anyPrompt",
				TargetPointer: new(bool),
				Required:      false,
			},
			str:           "abc\n",
			expectedValue: false,
			expectedErr:   fmt.Errorf("this field should be a bool value"),
		},
		{
			name: "empty but not required",
			in: Input{
				Message:       "anyPrompt",
				TargetPointer: new(bool),
				Required:      false,
			},
			str:           "\n",
			expectedValue: false,
			expectedErr:   nil,
		},
		{
			name: "missing required",
			in: Input{
				Message:       "anyPrompt",
				TargetPointer: new(bool),
				Required:      true,
			},
			str:           "\n",
			expectedValue: false,
			expectedErr:   fmt.Errorf("this field is required"),
		},
		{
			name: "too many arguments",
			in: Input{
				Message:       "anyPrompt",
				TargetPointer: new(bool),
				Required:      true,
			},
			str:           "true false\n",
			expectedValue: false,
			expectedErr:   fmt.Errorf("too many arguments"),
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			strReader := strings.NewReader(tc.str)
			p := NewPrompt(strReader)
			err := p.RunOne(tc.in)
			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expectedValue, *tc.in.TargetPointer.(*bool))
		})
	}
}
