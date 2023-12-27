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

package initialize

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

const (
	maximum_retry = 5
)

// Read a string value.
// func promptString(message string, in io.Reader) (string, error) {
// 	var s string
// 	fmt.Printf("%s ", message)
// 	n, e := fmt.Fscan(in, &s)
// 	fmt.Printf("  [read] [n:%d] %s: %s\n", n, message, s)
// 	if n > 1 {
// 		return s, fmt.Errorf("too many arguments")
// 	}
// 	if e != nil && n > 0 {
// 		// If the input is empty, ignore the error.
// 		return "", e
// 	}

//		return s, nil
//	}

// Read a string value. If the input is empty, return an empty string without error.
// If the input has multiple words, return an error.
func promptString(message string, in io.Reader) (string, error) {
	fmt.Printf("%s ", message)
	reader := bufio.NewReader(in)
	line, e := reader.ReadString('\n')
	// line, _, e := reader.ReadLine()
	if e != nil && e != io.EOF {
		return "", e
	}

	txts := strings.Fields(line)
	if len(txts) > 1 {
		return "", fmt.Errorf("too many arguments")
	}

	fmt.Println("  [read-txts] ", txts)
	fmt.Println("  [read-line] ", line)

	if len(txts) == 0 {
		return "", nil
	}
	return txts[0], nil
}

// [Scanner]
// func promptString(message string, in io.Reader) (string, error) {
// 	fmt.Printf("%s ", message)
// 	scanner := bufio.NewScanner(in)
// 	read := scanner.Scan()
// 	if !read {
// 		return "", fmt.Errorf("failed to read input: %v", scanner.Err())
// 	}

// 	s := scanner.Text()
// 	trim := strings.Fields(s)
// 	if len(trim) == 0 {
// 		// when input is empty
// 		return "", nil
// 	}
// 	if len(trim) > 1 {
// 		// when input contains spaces
// 		return "", fmt.Errorf("too many arguments")
// 	}

// 	return trim[0], nil
// }

// Read string values separated by space.
func promptStrings(message string, in io.Reader) ([]string, error) {
	fmt.Printf("%s ", message)
	reader := bufio.NewReader(in)
	line, e := reader.ReadString('\n')
	if e != nil {
		return nil, e
	}
	return strings.Fields(line), nil
}

// Read a int value and validate it.
func promptInt(message string, in io.Reader) (int, error) {
	fmt.Printf("%s ", message)
	reader := bufio.NewReader(in)
	txt, e := reader.ReadString('\n')
	if e != nil {
		return 0, e
	}
	fields := strings.Fields(txt)
	if len(fields) == 0 {
		return 0, nil
	}

	n, e := strconv.Atoi(fields[0])
	if e != nil {
		return 0, e
	}
	return n, nil
}

// Read a string value and validate it is not empty.
func promptStringRequired(message string, in io.Reader) (string, error) {
	for i := 0; i < maximum_retry; i++ {
		s, e := promptString(message, in)
		if e != nil {
			return "", e
		}
		if s != "" {
			return s, nil
		}
		fmt.Printf("[WARN] This field is required. \n")
	}
	return "", fmt.Errorf("This field is required. Maximum retry exceeded.")
}
