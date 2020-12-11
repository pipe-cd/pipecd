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

package yamlutil

import (
	"bytes"
	"fmt"

	gyaml "github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/parser"
)

// GetValue gives back the value placed at a given path.
//
// The path requires to start with "$" which represents the root element.
// Available operators are:
// $     : the root object/element
// .     : child operator
// ..    : recursive descent
// [num] : object/element of array by number
// [*]   : all objects/elements for array.
//
// e.g. "$.spec.template.spec.containers[0].image"
func GetValue(yml []byte, path string) (interface{}, error) {
	if len(yml) == 0 {
		return nil, fmt.Errorf("empty yaml given")
	}
	if path == "" {
		return nil, fmt.Errorf("no path given")
	}

	p, err := gyaml.PathString(path)
	if err != nil {
		return nil, err
	}

	var value interface{}
	if err := p.Read(bytes.NewReader(yml), &value); err != nil {
		return nil, err
	}
	return value, nil
}

// ReplaceValue replaces the value placed at a given path with
// a given value, and then gives back the new yaml content.
// The supported types are: string, bool, float64, int64, uint64 and nil.
// An error is returned if other types is given as a value.
func ReplaceValue(yml []byte, path string, value interface{}) ([]byte, error) {
	if len(yml) == 0 {
		return nil, fmt.Errorf("empty yaml given")
	}
	if path == "" {
		return nil, fmt.Errorf("no path given")
	}

	p, err := gyaml.PathString(path)
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

	var newNode ast.Node
	switch v := value.(type) {
	case string:
		newNode = &ast.StringNode{
			BaseNode: &ast.BaseNode{},
			Token:    oldNode.GetToken(),
			Value:    v,
		}
	case bool:
		newNode = &ast.BoolNode{
			BaseNode: &ast.BaseNode{},
			Token:    oldNode.GetToken(),
			Value:    v,
		}
	case float64:
		newNode = &ast.FloatNode{
			BaseNode: &ast.BaseNode{},
			Token:    oldNode.GetToken(),
			Value:    v,
		}
	case int64, uint64:
		newNode = &ast.IntegerNode{
			BaseNode: &ast.BaseNode{},
			Token:    oldNode.GetToken(),
			Value:    v,
		}
	case nil:
		newNode = &ast.NullNode{
			BaseNode: &ast.BaseNode{},
			Token:    oldNode.GetToken(),
		}
	default:
		return nil, fmt.Errorf("the given value is unsupported type")
	}
	err = p.ReplaceWithNode(file, newNode)
	return []byte(file.String()), err
}
