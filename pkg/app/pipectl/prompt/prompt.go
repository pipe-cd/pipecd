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

type Prompt struct {
	bufReader bufio.Reader
}

type Input struct {
	Message       string
	TargetPointer any
	Required      bool
}

func NewPrompt(in io.Reader) Prompt {
	return Prompt{
		bufReader: *bufio.NewReader(in),
	}
}

func (prompt *Prompt) Run(inputs []Input) error {
	for _, in := range inputs {
		if e := prompt.run(in); e != nil {
			return e
		}
	}
	return nil
}

func (prompt *Prompt) run(in Input) error {
	fmt.Printf("%s ", in.Message)
	line, e := prompt.bufReader.ReadString('\n')
	if e != nil {
		return e
	}
	fields := strings.Fields(line)

	if len(fields) == 0 {
		if in.Required {
			return fmt.Errorf("this field is required")
		} else {
			return nil
		}
	}

	fmt.Println("fields[0]:", fields[0])

	switch p := in.TargetPointer.(type) {
	case *string:
		if len(fields) > 1 {
			return fmt.Errorf("too many arguments")
		}
		*p = fields[0]

		fmt.Println("p:", p)
		fmt.Println("*p:", *p)
		fmt.Println("&p:", &p)

	case *[]string:
		*p = fields
	case *int:
		if len(fields) > 1 {
			return fmt.Errorf("too many arguments")
		}
		n, e := strconv.Atoi(fields[0])
		if e != nil {
			return fmt.Errorf("this field should be an int value: %s", fields[0])
		}
		*p = n
	case *bool:
		if len(fields) > 1 {
			return fmt.Errorf("too many arguments")
		}
		b, e := strconv.ParseBool(fields[0])
		if e != nil {
			return fmt.Errorf("this field should be a bool value: %s", fields[0])
		}
		*p = b
	default:
		return fmt.Errorf("unsupported type: %T", p)
	}
	return nil
}
