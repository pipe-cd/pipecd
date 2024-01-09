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
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

const (
	maximum_retry = 5 // maximum number of retry for required fields
)

// Reader is an interface that reads input from user or mock.
type Reader interface {
	// ReadString reads a string with showing message. If the input is empty, it returns an empty string.
	// If the input is more than one word, it returns an error.
	ReadString(message string) (string, error)
	// ReadStrings reads a string slice separated by blank space with showing message. If the input is empty, it returns an empty slice.
	ReadStrings(message string) ([]string, error)
	// ReadInt reads an integer with showing message. If the input is empty, it returns 0.
	// If the input is not an integer, it returns an error.
	ReadInt(message string) (int, error)
	// ReadStringRequired reads a string with showing message. If the input is empty, it asks again.
	// If number of retry exceeds maximum_retry, it returns an error.
	// If the input is more than one word, it returns an error.
	ReadStringRequired(message string) (string, error)
}

// stdReader is a Reader implementation that reads input from stdin.
type stdinReader struct {
	reader bufio.Reader
}

func NewReader(in io.Reader) Reader {
	return &stdinReader{
		reader: *bufio.NewReader(in),
	}
}

func (r *stdinReader) ReadString(message string) (string, error) {
	s, e := r.ReadStrings(message)
	if e != nil {
		return "", e
	}

	if len(s) == 0 {
		return "", nil
	}
	if len(s) > 1 {
		return "", fmt.Errorf("too many arguments")
	}

	return s[0], nil
}

func (r *stdinReader) ReadStrings(message string) ([]string, error) {
	fmt.Printf("%s ", message)
	line, e := r.reader.ReadString('\n')
	if e != nil {
		return nil, e
	}
	return strings.Fields(line), nil
}

func (r *stdinReader) ReadInt(message string) (int, error) {
	s, e := r.ReadString(message)
	if e != nil {
		return 0, e
	}
	if len(s) == 0 {
		return 0, nil
	}
	n, e := strconv.Atoi(s)
	if e != nil {
		return 0, e
	}
	return n, nil
}

func (r *stdinReader) ReadStringRequired(message string) (string, error) {
	for i := 0; i < maximum_retry; i++ {
		s, e := r.ReadString(message)
		if e != nil {
			return "", e
		}
		if len(s) > 0 {
			return s, nil
		}
		fmt.Printf("[WARN] This field is required. \n")
	}
	return "", fmt.Errorf("this field is required. Maximum retry(%d) exceeded", maximum_retry)
}
