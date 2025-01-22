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
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/yaml"
)

type differ struct {
	ignoreAddingMapKeys           bool
	equateEmpty                   bool
	compareNumberAndNumericString bool
	ignoredPaths                  map[string]struct{}
	ignoreConfig                  map[string][]string

	result *Result
}

type Option func(*differ)

// WithIgnoreAddingMapKeys configures differ to ignore all map keys
// that were added to the second object and missing in the first one.
func WithIgnoreAddingMapKeys() Option {
	return func(d *differ) {
		d.ignoreAddingMapKeys = true
	}
}

// WithEquateEmpty configures differ to consider all maps/slides with a length of zero and nil to be equal.
func WithEquateEmpty() Option {
	return func(d *differ) {
		d.equateEmpty = true
	}
}

// WithCompareNumberAndNumericString configures differ to compare a number with a numeric string.
// Differ parses the string to number before comparing their values.
// e.g. 1.5 == "1.5"
func WithCompareNumberAndNumericString() Option {
	return func(d *differ) {
		d.compareNumberAndNumericString = true
	}
}

// WithIgnoreConfig configures ignored fields.
func WithIgnoreConfig(config map[string][]string) Option {
	return func(d *differ) {
		d.ignoreConfig = config
	}
}

func (d *differ) initIgnoredPaths(key string) {
	paths := d.ignoreConfig[key]
	d.ignoredPaths = make(map[string]struct{}, len(paths))

	for _, path := range paths {
		d.ignoredPaths[path] = struct{}{}
	}
}

// DiffUnstructureds calculates the diff between two unstructured objects.
// If you compare non-k8s manifests, use DiffStructureds instead.
func DiffUnstructureds(x, y unstructured.Unstructured, key string, opts ...Option) (*Result, error) {
	var (
		path = []PathStep{}
		vx   = reflect.ValueOf(x.Object)
		vy   = reflect.ValueOf(y.Object)
		d    = &differ{result: &Result{}}
	)
	for _, opt := range opts {
		opt(d)
	}

	d.initIgnoredPaths(key)

	if err := d.diff(path, vx, vy); err != nil {
		return nil, err
	}

	d.result.sort()
	return d.result, nil
}

// DiffStructureds calulates the diff between non-k8s manifests.
// If you compare k8s manifests, use DiffUnstructureds instead.
func DiffStructureds(x, y interface{}, opts ...Option) (*Result, error) {
	mapX := map[string]interface{}{}
	mapY := map[string]interface{}{}
	ymlX, err := yaml.Marshal(x)
	if err != nil {
		return nil, err
	}
	ymlY, err := yaml.Marshal(y)
	if err != nil {
		return nil, err
	}
	yaml.Unmarshal(ymlX, &mapX)
	yaml.Unmarshal(ymlY, &mapY)

	var (
		path = []PathStep{}
		vx   = reflect.ValueOf(mapX)
		vy   = reflect.ValueOf(mapY)
		d    = &differ{result: &Result{}}
	)
	for _, opt := range opts {
		opt(d)
	}

	if err := d.diff(path, vx, vy); err != nil {
		return nil, err
	}

	d.result.sort()
	return d.result, nil
}

func (d *differ) diff(path []PathStep, vx, vy reflect.Value) error {
	if !vx.IsValid() {
		if d.equateEmpty && isEmptyInterface(vy) {
			return nil
		}

		d.addNode(path, nil, vy.Type(), vx, vy)
		return nil
	}

	if !vy.IsValid() {
		if d.equateEmpty && isEmptyInterface(vx) {
			return nil
		}

		d.addNode(path, vx.Type(), nil, vx, vy)
		return nil
	}

	if isNumberValue(vx) && isNumberValue(vy) {
		return d.diffNumber(path, vx, vy)
	}

	if d.compareNumberAndNumericString {
		if isNumberValue(vx) {
			if y, ok := convertToNumber(vy); ok {
				return d.diffNumber(path, vx, y)
			}
		}

		if isNumberValue(vy) {
			if x, ok := convertToNumber(vx); ok {
				return d.diffNumber(path, x, vy)
			}
		}
	}

	if vx.Type() != vy.Type() {
		d.addNode(path, vx.Type(), vy.Type(), vx, vy)
		return nil
	}

	switch vx.Kind() {
	case reflect.Map:
		return d.diffMap(path, vx, vy)

	case reflect.Slice, reflect.Array:
		return d.diffSlice(path, vx, vy)

	case reflect.Interface:
		return d.diffInterface(path, vx, vy)

	case reflect.String:
		return d.diffString(path, vx, vy)

	case reflect.Bool:
		return d.diffBool(path, vx, vy)

	default:
		return fmt.Errorf("%v kind is not handled", vx.Kind())
	}
}

func (d *differ) diffSlice(path []PathStep, vx, vy reflect.Value) error {
	if vx.IsNil() || vy.IsNil() {
		d.addNode(path, vx.Type(), vy.Type(), vx, vy)
		return nil
	}

	minLen := vx.Len()
	if minLen > vy.Len() {
		minLen = vy.Len()
	}

	for i := 0; i < minLen; i++ {
		nextPath := newSlicePath(path, i)
		nextValueX := vx.Index(i)
		nextValueY := vy.Index(i)
		if err := d.diff(nextPath, nextValueX, nextValueY); err != nil {
			return err
		}
	}

	for i := minLen; i < vx.Len(); i++ {
		nextPath := newSlicePath(path, i)
		nextValueX := vx.Index(i)
		d.addNode(nextPath, nextValueX.Type(), nextValueX.Type(), nextValueX, reflect.Value{})
	}

	for i := minLen; i < vy.Len(); i++ {
		nextPath := newSlicePath(path, i)
		nextValueY := vy.Index(i)
		d.addNode(nextPath, nextValueY.Type(), nextValueY.Type(), reflect.Value{}, nextValueY)
	}

	return nil
}

func (d *differ) diffMap(path []PathStep, vx, vy reflect.Value) error {
	if vx.IsNil() || vy.IsNil() {
		d.addNode(path, vx.Type(), vy.Type(), vx, vy)
		return nil
	}

	keys := append(vx.MapKeys(), vy.MapKeys()...)
	checks := make(map[string]struct{})

	for _, k := range keys {
		if k.Kind() != reflect.String {
			return fmt.Errorf("unsupport %v as key type of a map", k.Kind())
		}
		if _, ok := checks[k.String()]; ok {
			continue
		}

		nextValueY := vy.MapIndex(k)
		// Don't need to check the key existing in the first(LiveManifest) one but missing in the seccond(GitManifest) one
		// when IgnoreAddingMapKeys is configured.
		if d.ignoreAddingMapKeys && !nextValueY.IsValid() {
			continue
		}

		nextPath := newMapPath(path, k.String())
		nextValueX := vx.MapIndex(k)
		checks[k.String()] = struct{}{}
		if err := d.diff(nextPath, nextValueX, nextValueY); err != nil {
			return err
		}
	}
	return nil
}

func (d *differ) diffInterface(path []PathStep, vx, vy reflect.Value) error {
	if isEmptyInterface(vx) && isEmptyInterface(vy) {
		return nil
	}

	if vx.IsNil() || vy.IsNil() {
		d.addNode(path, vx.Type(), vy.Type(), vx, vy)
		return nil
	}

	vx, vy = vx.Elem(), vy.Elem()
	return d.diff(path, vx, vy)
}

func (d *differ) diffString(path []PathStep, vx, vy reflect.Value) error {
	if vx.String() == vy.String() {
		return nil
	}
	d.addNode(path, vx.Type(), vy.Type(), vx, vy)
	return nil
}

func (d *differ) diffBool(path []PathStep, vx, vy reflect.Value) error {
	if vx.Bool() == vy.Bool() {
		return nil
	}
	d.addNode(path, vx.Type(), vy.Type(), vx, vy)
	return nil
}

func (d *differ) diffNumber(path []PathStep, vx, vy reflect.Value) error {
	if floatNumber(vx) == floatNumber(vy) {
		return nil
	}

	d.addNode(path, vx.Type(), vy.Type(), vx, vy)
	return nil
}

// isEmptyInterface reports whether v is nil or zero value or its element is an empty map, an empty slice.
func isEmptyInterface(v reflect.Value) bool {
	if !v.IsValid() || v.IsNil() || v.IsZero() {
		return true
	}
	if v.Kind() != reflect.Interface {
		return false
	}

	e := v.Elem()

	// When the value that the interface v contains is a zero value (false boolean, zero number, empty string...).
	if e.IsZero() {
		return true
	}

	switch e.Kind() {
	case reflect.Array, reflect.Slice, reflect.Map:
		return e.Len() == 0
	default:
		return false
	}
}

func floatNumber(v reflect.Value) float64 {
	switch v.Kind() {
	case reflect.Float32, reflect.Float64:
		return v.Float()
	default:
		return float64(v.Int())
	}
}

func isNumberValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Float32, reflect.Float64:
		return true
	default:
		return false
	}
}

func convertToNumber(v reflect.Value) (reflect.Value, bool) {
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Float32, reflect.Float64:
		return v, true
	case reflect.String:
		if n, err := strconv.ParseFloat(v.String(), 64); err == nil {
			return reflect.ValueOf(n), true
		}
		return v, false
	default:
		return v, false
	}
}

func newSlicePath(path []PathStep, index int) []PathStep {
	next := make([]PathStep, len(path))
	copy(next, path)
	next = append(next, PathStep{
		Type:       SliceIndexPathStep,
		SliceIndex: index,
	})
	return next
}

func newMapPath(path []PathStep, index string) []PathStep {
	next := make([]PathStep, len(path))
	copy(next, path)
	next = append(next, PathStep{
		Type:     MapIndexPathStep,
		MapIndex: index,
	})
	return next
}

func (d *differ) addNode(path []PathStep, tx, ty reflect.Type, vx, vy reflect.Value) {
	if len(d.ignoredPaths) > 0 {
		pathString := makePathString(path)
		if d.isIgnoredPath(pathString) {
			return
		}
		nvx := d.ignoredValue(vx, pathString)
		nvy := d.ignoredValue(vy, pathString)

		d.result.addNode(path, tx, ty, nvx, nvy)
		return
	}

	d.result.addNode(path, tx, ty, vx, vy)
}

func (d *differ) ignoredValue(v reflect.Value, prefix string) reflect.Value {
	switch v.Kind() {
	case reflect.Map:
		nv := reflect.MakeMap(v.Type())
		keys := v.MapKeys()
		for _, k := range keys {
			nprefix := prefix + "." + k.String()
			if d.isIgnoredPath(nprefix) {
				continue
			}

			sub := v.MapIndex(k)
			filtered := d.ignoredValue(sub, nprefix)
			if !filtered.IsValid() {
				continue
			}
			nv.SetMapIndex(k, filtered)
		}
		return nv

	case reflect.Slice, reflect.Array:
		nv := reflect.MakeSlice(v.Type(), 0, 0)
		for i := 0; i < v.Len(); i++ {
			nprefix := prefix + "." + strconv.Itoa(i)
			if d.isIgnoredPath(nprefix) {
				continue
			}

			filtered := d.ignoredValue(v.Index(i), nprefix)
			if !filtered.IsValid() {
				continue
			}
			nv = reflect.Append(nv, filtered)
		}
		return nv

	case reflect.Interface:
		nprefix := prefix + "." + v.String()
		if d.isIgnoredPath(nprefix) {
			return reflect.New(v.Type())
		}
		return d.ignoredValue(v.Elem(), prefix)

	case reflect.String:
		return v

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v

	case reflect.Float32, reflect.Float64:
		return v

	default:
		nprefix := prefix + "." + v.String()
		if d.isIgnoredPath(nprefix) {
			return reflect.New(v.Type())
		}
		return v
	}
}

func (d *differ) isIgnoredPath(pathString string) bool {
	var pathSubStr string
	pathElms := strings.Split(pathString, ".")

	for i, path := range pathElms {
		if i != 0 {
			pathSubStr += "."
		}
		pathSubStr += path
		if _, found := d.ignoredPaths[pathSubStr]; found {
			return true
		}
	}
	return false
}
