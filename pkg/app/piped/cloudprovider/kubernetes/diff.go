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
	"regexp"
	"strconv"
	"strings"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/pipe-cd/pipe/pkg/app/piped/diff"
)

type diffOption struct {
	PathPrefix              string
	IgnoreOrder             bool
	IgnoreMissingFields     bool
	RedactPathPrefix        string
	RedactFirstReplacement  string
	RedactSecondReplacement string
}

type DiffOption func(*diffOption)

// WithDiffPathPrefix configures differ to returns only results where
// their paths are prefixed by the given string.
func WithDiffPathPrefix(prefix string) DiffOption {
	return func(o *diffOption) {
		o.PathPrefix = prefix
	}
}

// WithDiffIgnoreOrder configures differ to ignore the order of slice while calculating.
func WithDiffIgnoreOrder() DiffOption {
	return func(o *diffOption) {
		o.IgnoreOrder = true
	}
}

// WithDiffIgnoreMissingFields configures differ to ignore all fields
// that appear in the second manifest but not in the first.
func WithDiffIgnoreMissingFields() DiffOption {
	return func(o *diffOption) {
		o.IgnoreMissingFields = true
	}
}

// WithDiffRedactPathPrefix redacts all the diff values of the given path prefix.
func WithDiffRedactPathPrefix(prefix, firstReplacement, secondReplacement string) DiffOption {
	return func(o *diffOption) {
		o.RedactPathPrefix = prefix
		o.RedactFirstReplacement = firstReplacement
		o.RedactSecondReplacement = secondReplacement
	}
}

type PathStepType string

const (
	MapKeyPathStep     PathStepType = "MapKey"
	SliceIndexPathStep PathStepType = "SliceIndex"
)

type PathStep struct {
	Type  PathStepType
	Key   string
	Index int
}

type DiffPath []PathStep

func (p DiffPath) String() string {
	var ss []string
	for _, s := range p {
		switch s.Type {
		case SliceIndexPathStep:
			ss = append(ss, fmt.Sprintf("[%d]", s.Index))
		default:
			ss = append(ss, s.Key)
		}
	}
	return strings.Join(ss, ".")
}

func parseDiffPath(s string) (DiffPath, error) {
	parts := strings.Split(s, ".")
	path := make(DiffPath, 0, len(parts))

	for _, p := range parts {
		if strings.HasPrefix(p, "[") {
			p = p[1 : len(p)-1]
			index, err := strconv.Atoi(p)
			if err != nil {
				return nil, err
			}
			path = append(path, PathStep{
				Type:  SliceIndexPathStep,
				Index: index,
			})
			continue
		}
		path = append(path, PathStep{
			Type: MapKeyPathStep,
			Key:  p,
		})
	}

	return path, nil
}

type DiffResult struct {
	Path       DiffPath
	PathString string
	Before     string
	After      string
}

func (d DiffResult) String() string {
	return fmt.Sprintf("%v:\n\t-: %s\n\t+: %s\n", d.Path, d.Before, d.After)
}

type DiffResultList []DiffResult

func (dl DiffResultList) Find(query string) (result DiffResult, found bool, err error) {
	reg, err := regexp.Compile(query)
	if err != nil {
		return
	}

	for _, d := range dl {
		matched := reg.MatchString(d.PathString)
		if !matched {
			continue
		}
		return d, true, nil
	}
	return
}

func (dl DiffResultList) FindAll(query string) (list []DiffResult) {
	reg, err := regexp.Compile(query)
	if err != nil {
		return
	}

	for _, d := range dl {
		matched := reg.MatchString(d.PathString)
		if !matched {
			continue
		}
		list = append(list, d)
	}
	return
}

func (dl DiffResultList) FindByPrefix(prefix string) (list []DiffResult) {
	for _, d := range dl {
		if strings.HasPrefix(d.PathString, prefix) {
			list = append(list, d)
		}
	}
	return
}

func (dl DiffResultList) String() string {
	ds := make([]string, 0, len(dl))
	for _, d := range dl {
		ds = append(ds, d.String())
	}
	return strings.Join(ds, "\n")
}

func DiffWorkload(first, second Manifest) string {
	return ""
}

func Diff(first, second Manifest, opts ...DiffOption) DiffResultList {
	var options diffOption
	for _, opt := range opts {
		opt(&options)
	}

	reporter := diffReporter{
		options: options,
	}

	if !options.IgnoreOrder {
		cmp.Equal(first.u, second.u, cmp.Reporter(&reporter))
	} else {
		cmp.Equal(first.u, second.u, cmp.Reporter(&reporter), cmpopts.SortSlices(sortLess))
	}
	return reporter.diffs
}

func DiffManifests(first, second Manifest, opts ...diff.DiffOption) (*diff.Result, error) {
	return diff.DiffUnstructureds(*first.u, *second.u, opts...)
}

// diffReporter is a custom reporter that only records differences
// detected during comparison.
type diffReporter struct {
	path    cmp.Path
	options diffOption
	diffs   DiffResultList
}

func (r *diffReporter) PushStep(ps cmp.PathStep) {
	r.path = append(r.path, ps)
}

func (r *diffReporter) Report(rs cmp.Result) {
	if !rs.Equal() {
		var (
			path       = convertDiffPath(r.path)
			pathString = path.String()
			vx, vy     = r.path.Last().Values()
		)
		if !strings.HasPrefix(pathString, r.options.PathPrefix) {
			return
		}
		if r.options.IgnoreMissingFields && !vx.IsValid() {
			return
		}
		var before, after string
		if r.options.RedactPathPrefix != "" && strings.HasPrefix(pathString, r.options.RedactPathPrefix) {
			before = r.options.RedactFirstReplacement
			after = r.options.RedactSecondReplacement
		} else {
			before = fmt.Sprintf("%+v", vx)
			after = fmt.Sprintf("%+v", vy)
		}

		r.diffs = append(r.diffs, DiffResult{
			Path:       path,
			PathString: pathString,
			Before:     before,
			After:      after,
		})
	}
}

func (r *diffReporter) PopStep() {
	r.path = r.path[:len(r.path)-1]
}

func (r *diffReporter) String() string {
	return r.diffs.String()
}

func convertDiffPath(path cmp.Path) DiffPath {
	p := make([]PathStep, 0, len(path))

	for _, s := range path {
		switch s := s.(type) {
		case cmp.SliceIndex:
			p = append(p, PathStep{
				Type:  SliceIndexPathStep,
				Index: s.Key(),
			})
		case cmp.MapIndex:
			p = append(p, PathStep{
				Type: MapKeyPathStep,
				Key:  fmt.Sprintf("%v", s.Key()),
			})
		}
	}
	return p
}

func sortLess(i, j interface{}) bool {
	switch t := i.(type) {
	case string:
		return t < j.(string)
	case int:
		return t < j.(int)
	case int32:
		return t < j.(int32)
	case int64:
		return t < j.(int64)
	case float32:
		return t < j.(float32)
	case float64:
		return t < j.(float64)
	}
	return false
}
