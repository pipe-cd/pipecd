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
	"strconv"
	"strings"
)

// mockReader is a mock implementation of Reader for unit tests.
// It reads input from the given slice of strings instead of stdin.
type mockReader struct {
	inputs []string
}

func NewMockReader(inputs []string) *mockReader {
	return &mockReader{inputs: inputs}
}

func (r *mockReader) ReadString(message string) (string, error) {
	if len(r.inputs) == 0 {
		return "", fmt.Errorf("no more input")
	}
	input := r.inputs[0]
	r.inputs = r.inputs[1:]
	f := strings.Fields(input)
	if len(f) > 1 {
		return "", fmt.Errorf("too many arguments")
	}

	if len(f) == 0 {
		return "", nil
	}
	return f[0], nil
}

func (r *mockReader) ReadStrings(message string) ([]string, error) {
	if len(r.inputs) == 0 {
		return nil, fmt.Errorf("no more input")
	}
	input := r.inputs[0]
	r.inputs = r.inputs[1:]
	return strings.Fields(input), nil
}

func (r *mockReader) ReadInt(message string) (int, error) {
	if len(r.inputs) == 0 {
		return 0, fmt.Errorf("no more input")
	}
	input := r.inputs[0]
	r.inputs = r.inputs[1:]
	if len(input) == 0 {
		return 0, nil
	}
	return strconv.Atoi(input)
}

func (r *mockReader) ReadStringRequired(message string) (string, error) {
	s, e := r.ReadString(message)
	if e != nil {
		return "", e
	}
	if len(s) == 0 {
		return "", fmt.Errorf("empty input")
	}
	return s, e
}
