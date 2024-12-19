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

package prompt

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Prompt is a helper for asking inputs from users.
type Prompt struct {
	bufReader bufio.Reader
}

// Input represents an query to ask from users.
type Input struct {
	// Message is a message to show when asking for an input. e.g. "Which platform to use?"
	Message string
	// TargetPointer is a pointer to the target variable to set the value of the input.
	TargetPointer any
	// Required indicates whether this input is required or not.
	Required bool
}

func NewPrompt(in io.Reader) Prompt {
	return Prompt{
		bufReader: *bufio.NewReader(in),
	}
}

// RunSlice sequentially asks for inputs and set the values to the target pointers.
func (p *Prompt) RunSlice(inputs []Input) error {
	for _, in := range inputs {
		if err := p.Run(in); err != nil {
			return err
		}
	}
	return nil
}

// Run asks for an input and set the value to the target pointer.
func (p *Prompt) Run(input Input) error {
	fmt.Printf("%s: ", input.Message)
	line, err := p.bufReader.ReadString('\n')
	if err != nil {
		return err
	}
	// split by spaces
	fields := strings.Fields(line)

	if len(fields) == 0 {
		if input.Required {
			return fmt.Errorf("this field is required")
		} else {
			return nil
		}
	}

	switch tp := input.TargetPointer.(type) {
	case *string:
		if len(fields) > 1 {
			return fmt.Errorf("too many arguments")
		}
		*tp = fields[0]
	case *[]string:
		*tp = fields
	case *int:
		if len(fields) > 1 {
			return fmt.Errorf("too many arguments")
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return fmt.Errorf("this field should be an int value")
		}
		*tp = n
	case *bool:
		if len(fields) > 1 {
			return fmt.Errorf("too many arguments")
		}

		switch fields[0] {
		case "y", "Y":
			*tp = true
		case "n", "N":
			*tp = false
		default:
			b, err := strconv.ParseBool(fields[0])
			if err != nil {
				return fmt.Errorf("this field should be a bool value")
			}
			*tp = b
		}
	default:
		return fmt.Errorf("unsupported type: %T", tp)
	}
	return nil
}
