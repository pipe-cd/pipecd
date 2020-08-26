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
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

type Renderer struct {
	leftPadding        int
	redactPathPrefix   string
	redactReplacementX string
	redactReplacementY string
}

type RenderOption func(*Renderer)

func WithLeftPadding(p int) RenderOption {
	return func(r *Renderer) {
		r.leftPadding = p
	}
}

func WithRedactPath(prefix, replacementX, replacementY string) RenderOption {
	return func(r *Renderer) {
		r.redactPathPrefix = prefix
		r.redactReplacementX = replacementX
		r.redactReplacementY = replacementY
	}
}

func NewRenderer(opts ...RenderOption) *Renderer {
	r := &Renderer{}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

func (r *Renderer) Render(ns Nodes) string {
	if len(ns) == 0 {
		return ""
	}

	var prePath []PathStep
	var b strings.Builder

	printValue := func(mark string, v reflect.Value, lastStep PathStep, depth int) {
		nodeString, nl := renderNodeValue(v, "")
		if nodeString == "" {
			return
		}

		if lastStep.Type == SliceIndexPathStep {
			nl = false
		}

		switch {
		case lastStep.Type == SliceIndexPathStep:
			b.WriteString(fmt.Sprintf("%s%*s- ", mark, depth*2-1, ""))
		case nl:
			b.WriteString(fmt.Sprintf("%s%*s%s:\n", mark, depth*2-1, "", lastStep.String()))
		default:
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

	for _, n := range ns {
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

func renderNodeValue(v reflect.Value, prefix string) (string, bool) {
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
			subString, nl := renderNodeValue(sub, prefix+"  ")
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
			sub, _ := renderNodeValue(v.Index(i), prefix+"  ")
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
		return renderNodeValue(v.Elem(), prefix)

	case reflect.String:
		return v.String(), false

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatInt(v.Int(), 10), false

	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(v.Float(), 'f', -1, 64), false

	default:
		return v.String(), false
	}
}

func RenderPrimitiveValue(v reflect.Value) string {
	switch v.Kind() {
	case reflect.String:
		return v.String()

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatInt(v.Int(), 10)

	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(v.Float(), 'f', -1, 64)

	default:
		return v.String()
	}
}
