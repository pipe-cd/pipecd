// Copyright 2020 The PipeCD Authors.
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
	"bytes"
	"fmt"
	"io"

	goyaml "github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/parser"
)

// GetValue gives back the value placed at a given path. The type of
// returned value can be string, int64, uint64, float64, bool, and []interface{}.
//
// The path requires to start with "$" which represents the root element.
// Available operators are:
// $     : the root object/element
// .     : child operator
// ..    : recursive descent
// [num] : object/element of array by number
// [*]   : all objects/elements for array.
//
// e.g. "$.foo.bar[0].baz"
func GetValue(yml []byte, path string) (interface{}, error) {
	if len(yml) == 0 {
		return nil, fmt.Errorf("empty yaml given")
	}
	if path == "" {
		return nil, fmt.Errorf("no path given")
	}

	p, err := goyaml.PathString(path)
	if err != nil {
		return nil, err
	}
	// Build an AST node based on the given path.
	node, err := p.ReadNode(bytes.NewReader(yml))
	if err != nil {
		return nil, err
	}
	if node == nil {
		return nil, fmt.Errorf("wrong path (%s) given", path)
	}

	// TODO: Disallow yamlprocessor panic if the value is -0
	// TODO: Use goyaml.Read() upon merging our patch
	// https://github.com/goccy/go-yaml/pull/194
	var value interface{}
	if err := goyaml.Unmarshal([]byte(node.String()), &value); err != nil {
		return nil, fmt.Errorf("failed to unmarshal the node placed at the given path: %w", err)
	}
	return value, nil
}

// ReplaceValue replaces the value placed at a given path with
// a given value, and then gives back the new yaml bytes.
//
// For available operators for the path, see GetValue().
func ReplaceValue(yml []byte, path string, value string) ([]byte, error) {
	if len(yml) == 0 {
		return nil, fmt.Errorf("empty yaml given")
	}
	if path == "" {
		return nil, fmt.Errorf("no path given")
	}

	p, err := goyaml.PathString(path)
	if err != nil {
		return nil, err
	}
	file, err := parser.ParseBytes(yml, 0)
	if err != nil {
		return nil, err
	}

	// Retrieve the current node placed at the specified path.
	oldNode, err := p.FilterFile(file)
	if err != nil {
		return nil, err
	}
	newNode := &ast.StringNode{
		BaseNode: &ast.BaseNode{},
		Token:    oldNode.GetToken(),
		Value:    value,
	}

	err = p.ReplaceWithNode(file, newNode)
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, file)
	return buf.Bytes(), err
}
