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

package diff

import (
	"errors"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var (
	ErrNotFound = errors.New("not found")
)

type Result struct {
	nodes []Node
}

func (r *Result) HasDiff() bool {
	return len(r.nodes) > 0
}

func (r *Result) NumNodes() int {
	return len(r.nodes)
}

func (r *Result) Nodes() Nodes {
	return r.nodes
}

type Node struct {
	Path       []PathStep
	PathString string
	TypeX      reflect.Type
	TypeY      reflect.Type
	ValueX     reflect.Value
	ValueY     reflect.Value
}

type PathStepType string

const (
	MapIndexPathStep   PathStepType = "MapIndex"
	SliceIndexPathStep PathStepType = "SliceIndex"
)

type PathStep struct {
	Type       PathStepType
	MapIndex   string
	SliceIndex int
}

func (s PathStep) String() string {
	switch s.Type {
	case SliceIndexPathStep:
		return strconv.FormatInt(int64(s.SliceIndex), 10)
	case MapIndexPathStep:
		return s.MapIndex
	default:
		return ""
	}
}

func (n Node) StringX() string {
	return RenderPrimitiveValue(n.ValueX)
}

func (n Node) StringY() string {
	return RenderPrimitiveValue(n.ValueY)
}

type Nodes []Node

func (ns Nodes) FindOne(query string) (*Node, error) {
	reg, err := regexp.Compile(query)
	if err != nil {
		return nil, err
	}

	for i := range ns {
		matched := reg.MatchString(ns[i].PathString)
		if !matched {
			continue
		}
		return &ns[i], nil
	}
	return nil, ErrNotFound
}

func (ns Nodes) Find(query string) (Nodes, error) {
	reg, err := regexp.Compile(query)
	if err != nil {
		return nil, err
	}

	nodes := make([]Node, 0)
	for i := range ns {
		matched := reg.MatchString(ns[i].PathString)
		if !matched {
			continue
		}
		nodes = append(nodes, ns[i])
	}
	return nodes, nil
}

func (ns Nodes) FindByPrefix(prefix string) Nodes {
	nodes := make([]Node, 0)
	for i := range ns {
		if strings.HasPrefix(ns[i].PathString, prefix) {
			nodes = append(nodes, ns[i])
		}
	}
	return nodes
}

func (r *Result) addNode(path []PathStep, typeX, typeY reflect.Type, valueX, valueY reflect.Value) {
	r.nodes = append(r.nodes, Node{
		Path:       path,
		PathString: makePathString(path),
		TypeX:      typeX,
		TypeY:      typeY,
		ValueX:     valueX,
		ValueY:     valueY,
	})
}

func (r *Result) sort() {
	sort.Slice(r.nodes, func(i, j int) bool {
		return r.nodes[i].PathString < r.nodes[j].PathString
	})
}

func makePathString(path []PathStep) string {
	steps := make([]string, 0, len(path))
	for _, s := range path {
		steps = append(steps, s.String())
	}
	return strings.Join(steps, ".")
}
