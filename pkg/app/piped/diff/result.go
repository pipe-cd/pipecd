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

package diff

import (
	"errors"
	"fmt"
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
	leftPadding        int
	redactPathPrefix   string
	redactReplacementX string
	redactReplacementY string

	nodes []DiffNode
}

type DiffNode struct {
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

func (r *Result) SetLeftPadding(p int) {
	r.leftPadding = p
}

func (r *Result) SetRedactPath(prefix, replacementX, replacementY string) {
	r.redactPathPrefix = prefix
	r.redactReplacementX = replacementX
	r.redactReplacementY = replacementY
}

func (r *Result) Num() int {
	return len(r.nodes)
}

func (r *Result) HasDiff() bool {
	return len(r.nodes) > 0
}

func (r *Result) Find(query string) (*DiffNode, error) {
	reg, err := regexp.Compile(query)
	if err != nil {
		return nil, err
	}

	for i := range r.nodes {
		matched := reg.MatchString(r.nodes[i].PathString)
		if !matched {
			continue
		}
		return &r.nodes[i], nil
	}
	return nil, ErrNotFound
}

func (r *Result) FindAll(query string) ([]DiffNode, error) {
	reg, err := regexp.Compile(query)
	if err != nil {
		return nil, err
	}

	list := make([]DiffNode, 0)
	for i := range r.nodes {
		matched := reg.MatchString(r.nodes[i].PathString)
		if !matched {
			continue
		}
		list = append(list, r.nodes[i])
	}
	return list, nil
}

func (r *Result) FindByPrefix(prefix string) []DiffNode {
	list := make([]DiffNode, 0)
	for i := range r.nodes {
		if strings.HasPrefix(r.nodes[i].PathString, prefix) {
			list = append(list, r.nodes[i])
		}
	}
	return list
}

func (r *Result) DiffString() string {
	if len(r.nodes) == 0 {
		return ""
	}

	var prePath []PathStep
	var b strings.Builder

	printValue := func(mark string, v reflect.Value, lastStep PathStep, depth int) {
		nodeString, nl := PrintNodeValue(v)
		if nodeString == "" {
			return
		}

		if lastStep.Type == SliceIndexPathStep {
			nl = false
		}

		if lastStep.Type == SliceIndexPathStep {
			b.WriteString(fmt.Sprintf("%s%*s- ", mark, depth*2-1, ""))
		} else if nl {
			b.WriteString(fmt.Sprintf("%s%*s%s:\n", mark, depth*2-1, "", lastStep.String()))
		} else {
			b.WriteString(fmt.Sprintf("%s%*s%s: ", mark, depth*2-1, "", lastStep.String()))
		}

		parts := strings.Split(nodeString, "\n")
		for i, p := range parts {
			if lastStep.Type != SliceIndexPathStep {
				if nl {
					b.WriteString(fmt.Sprintf("%s%*s%s\n", mark, depth*2+1, "", p))
				} else {
					b.WriteString(fmt.Sprintf("%s\n", p))
				}
				continue
			}
			if i == 0 {
				b.WriteString(fmt.Sprintf("%s\n", p))
				continue
			}
			b.WriteString(fmt.Sprintf("%s%*s%s\n", mark, depth*2+1, "", p))
		}
	}

	for _, n := range r.nodes {
		duplicateDepth := pathDuplicateDepth(n.Path, prePath)
		prePath = n.Path
		b.WriteString(fmt.Sprintf("%*s#%s\n", (r.leftPadding+duplicateDepth)*2, "", n.PathString))

		var array bool
		for i := duplicateDepth; i < len(n.Path)-1; i++ {
			if n.Path[i].Type == SliceIndexPathStep {
				b.WriteString(fmt.Sprintf("%*s- ", (r.leftPadding+i)*2, ""))
				array = true
				continue
			}
			if array {
				b.WriteString(fmt.Sprintf("%s:\n", n.Path[i].String()))
				array = false
				continue
			}
			b.WriteString(fmt.Sprintf("%*s%s:\n", (r.leftPadding+i)*2, "", n.Path[i].String()))
		}

		lastStep := n.Path[len(n.Path)-1]
		valueX, valueY := n.ValueX, n.ValueY
		if r.redactPathPrefix != "" && strings.HasPrefix(n.PathString, r.redactPathPrefix) {
			valueX = reflect.ValueOf(r.redactReplacementX)
			valueY = reflect.ValueOf(r.redactReplacementY)
		}
		printValue("-", valueX, lastStep, r.leftPadding+len(n.Path)-1)
		printValue("+", valueY, lastStep, r.leftPadding+len(n.Path)-1)
	}

	return b.String()
}

func pathDuplicateDepth(x, y []PathStep) int {
	minLen := len(x)
	if minLen > len(y) {
		minLen = len(y)
	}

	for i := 0; i < minLen; i++ {
		if x[i] == y[i] {
			continue
		}
		return i
	}
	return 0
}

func PrintNodeValue(v reflect.Value) (string, bool) {
	return printNodeValue(v, "")
}

func printNodeValue(v reflect.Value, prefix string) (string, bool) {
	if !v.IsValid() {
		return "", false
	}

	switch v.Kind() {
	case reflect.Map:
		out := make([]string, 0, v.Len())
		keys := v.MapKeys()
		sort.Slice(keys, func(i, j int) bool {
			return keys[i].String() < keys[j].String()
		})
		for _, k := range keys {
			sub := v.MapIndex(k)
			subString, nl := printNodeValue(sub, prefix+"  ")
			if !nl {
				out = append(out, fmt.Sprintf("%s%s: %s", prefix, k.String(), subString))
				continue
			}
			out = append(out, fmt.Sprintf("%s%s:\n%s", prefix, k.String(), subString))
		}
		if len(out) == 0 {
			return "", false
		}
		return strings.Join(out, "\n"), true

	case reflect.Slice, reflect.Array:
		out := make([]string, 0, v.Len())
		for i := 0; i < v.Len(); i++ {
			sub, _ := printNodeValue(v.Index(i), prefix+"  ")
			parts := strings.Split(sub, "\n")
			for i, p := range parts {
				p = strings.TrimPrefix(p, prefix+"  ")
				if i == 0 {
					out = append(out, fmt.Sprintf("%s- %s", prefix, p))
					continue
				}
				out = append(out, fmt.Sprintf("%s  %s", prefix, p))
			}
		}
		return strings.Join(out, "\n"), true

	case reflect.Interface:
		return printNodeValue(v.Elem(), prefix)

	case reflect.String:
		return v.String(), false

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10), false

	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(v.Float(), 'f', -1, 64), false

	default:
		return v.String(), false
	}
}

func (r *Result) addNode(path []PathStep, typeX, typeY reflect.Type, valueX, valueY reflect.Value) {
	r.nodes = append(r.nodes, DiffNode{
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
