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

package provider

import (
	"encoding/json"
	"maps"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/yaml"
)

// Manifest represents a Kubernetes resource manifest.
type Manifest struct {
	Key  ResourceKey
	Body *unstructured.Unstructured
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (m *Manifest) UnmarshalJSON(data []byte) error {
	m.Body = new(unstructured.Unstructured)
	return m.Body.UnmarshalJSON(data)
}

// MarshalJSON implements the json.Marshaler interface.
// It marshals the underlying unstructured.Unstructured object into JSON bytes.
func (m *Manifest) MarshalJSON() ([]byte, error) {
	return m.Body.MarshalJSON()
}

// ConvertToStructuredObject converts the manifest into a structured Kubernetes object.
// The provided interface should be a pointer to a concrete Kubernetes type (e.g. *v1.Pod).
// It first marshals the manifest to JSON and then unmarshals it into the provided object.
func (m Manifest) ConvertToStructuredObject(o interface{}) error {
	data, err := m.MarshalJSON()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, o)
}

func (m *Manifest) YamlBytes() ([]byte, error) {
	return yaml.Marshal(m.Body)
}

func (m Manifest) AddAnnotations(annotations map[string]string) {
	if len(annotations) == 0 {
		return
	}

	annos := m.Body.GetAnnotations()
	if annos == nil {
		m.Body.SetAnnotations(annotations)
		return
	}
	for k, v := range annotations {
		annos[k] = v
	}
	m.Body.SetAnnotations(annos)
}

// AddStringMapValues adds or overrides the given key-values into the string map
// that can be found at the specified fields.
func (m Manifest) AddStringMapValues(values map[string]string, fields ...string) error {
	curMap, _, err := unstructured.NestedStringMap(m.Body.Object, fields...)
	if err != nil {
		return err
	}

	if curMap == nil {
		return unstructured.SetNestedStringMap(m.Body.Object, values, fields...)
	}
	maps.Copy(curMap, values)
	return unstructured.SetNestedStringMap(m.Body.Object, curMap, fields...)
}
