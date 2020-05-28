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

package kubernetes

import (
	"fmt"
	"strings"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

// diffReporter is a custom reporter that only records differences
// detected during comparison.
type diffReporter struct {
	path  cmp.Path
	diffs []string
}

func (r *diffReporter) PushStep(ps cmp.PathStep) {
	r.path = append(r.path, ps)
}

func (r *diffReporter) Report(rs cmp.Result) {
	if !rs.Equal() {
		vx, vy := r.path.Last().Values()
		r.diffs = append(r.diffs, fmt.Sprintf("%#v:\n\t-: %+v\n\t+: %+v\n", r.path, vx, vy))
	}
}

func (r *diffReporter) PopStep() {
	r.path = r.path[:len(r.path)-1]
}

func (r *diffReporter) String() string {
	return strings.Join(r.diffs, "\n")
}

type DiffResult struct {
	Path   string
	Before string
	After  string
}

type DiffResultList []DiffResult

func DiffWorkload(first, second Manifest) string {
	return ""
}

func Diff(first, second Manifest) DiffResult {
	var r diffReporter

	// trans := cmp.Transformer("Sort", func(in []string) []string {
	// 	out := append([]string(nil), in...) // Copy input to avoid mutating it
	// 	sort.Strings(out)
	// 	return out
	// })

	// cmp.Equal(first.u, second.u, trans, cmp.Reporter(&r))

	cmp.Equal(first.u, second.u, cmp.Reporter(&r), cmpopts.SortSlices(sortLess))
	fmt.Println("DIFFFFFFFFFFFFFFFFFFFFFFF")
	fmt.Println(r.String())
	fmt.Println("===============")

	return DiffResult{}
}

func sortLess(i, j interface{}) bool {
	switch i.(type) {
	case string:
		return i.(string) < j.(string)
	case int:
		return i.(int) < j.(int)
	case int32:
		return i.(int32) < j.(int32)
	case int64:
		return i.(int64) < j.(int64)
	case float32:
		return i.(float32) < j.(float32)
	case float64:
		return i.(float64) < j.(float64)
	}
	return false
}
