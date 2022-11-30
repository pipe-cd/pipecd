// Copyright 2022 The PipeCD Authors.
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

package yamlprocessor

import (
	"errors"
	"fmt"

	goyaml "github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/parser"
)

type Processor struct {
	file *ast.File
	newLineAtEOF bool
}

func NewProcessor(data []byte) (*Processor, error) {
	f, err := parser.ParseBytes(data, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	nl := false
	if len(data) > 0 && data[len(data)-1] == '\n' {
		nl = true
	}

	return &Processor{
		file: f,
		newLineAtEOF: nl,
	}, nil
}

// GetValue gives back the value placed at a given path. The type of
// returned value can be string, int64, uint64, float64 and bool.
//
// The path requires to start with "$" which represents the root element.
// Available operators are:
// $     : the root object/element
// .     : child operator
// ..    : recursive descent
// [num] : object/element of array by number
//
// e.g. "$.foo.bar[0].baz"
func (p *Processor) GetValue(path string) (interface{}, error) {
	if path == "" {
		return nil, fmt.Errorf("no path given")
	}

	yamlPath, err := goyaml.PathString(path)
	if err != nil {
		return nil, fmt.Errorf("failed to parse path %s: %w", path, err)
	}

	node, err := yamlPath.FilterFile(p.file)
	if err != nil {
		return nil, err
	}

	var value interface{}
	if err := goyaml.Unmarshal([]byte(node.String()), &value); err != nil {
		return nil, err
	}

	return value, nil
}

// ReplaceString replaces the value placed at a given path with
// a given string value.
func (p *Processor) ReplaceString(path, value string) error {
	if path == "" {
		return errors.New("no path given")
	}

	yamlPath, err := goyaml.PathString(path)
	if err != nil {
		return fmt.Errorf("failed to parse path %s: %w", path, err)
	}

	// Retrieve the current node placed at the specified path.
	oldNode, err := yamlPath.FilterFile(p.file)
	if err != nil {
		return err
	}

	newNode := &ast.StringNode{
		BaseNode: &ast.BaseNode{},
		Token:    oldNode.GetToken(),
		Value:    value,
	}

	if c := oldNode.GetComment(); c != nil {
		newNode.SetComment(c)
	}

	return yamlPath.ReplaceWithNode(p.file, newNode)
}

func (p *Processor) Bytes() []byte {
	str := p.file.String()
	if p.newLineAtEOF {
		str += "\n"
	}

	return []byte(str)
}
