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

package kubernetes

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/yaml"

	"github.com/pipe-cd/pipecd/pkg/model"
)

type Manifest struct {
	Key ResourceKey
	u   *unstructured.Unstructured
}

func MakeManifest(key ResourceKey, u *unstructured.Unstructured) Manifest {
	return Manifest{
		Key: key,
		u:   u,
	}
}

func (m Manifest) Duplicate(name string) Manifest {
	u := m.u.DeepCopy()
	u.SetName(name)

	key := m.Key
	key.Name = name

	return Manifest{
		Key: key,
		u:   u,
	}
}

func (m Manifest) YamlBytes() ([]byte, error) {
	return yaml.Marshal(m.u)
}

func (m Manifest) MarshalJSON() ([]byte, error) {
	return m.u.MarshalJSON()
}

func (m Manifest) AddAnnotations(annotations map[string]string) {
	if len(annotations) == 0 {
		return
	}

	annos := m.u.GetAnnotations()
	if annos == nil {
		m.u.SetAnnotations(annotations)
		return
	}
	for k, v := range annotations {
		annos[k] = v
	}
	m.u.SetAnnotations(annos)
}

func (m Manifest) GetAnnotations() map[string]string {
	return m.u.GetAnnotations()
}

func (m Manifest) GetNestedStringMap(fields ...string) (map[string]string, error) {
	sm, _, err := unstructured.NestedStringMap(m.u.Object, fields...)
	if err != nil {
		return nil, err
	}

	return sm, nil
}

func (m Manifest) GetNestedMap(fields ...string) (map[string]interface{}, error) {
	sm, _, err := unstructured.NestedMap(m.u.Object, fields...)
	if err != nil {
		return nil, err
	}

	return sm, nil
}

// AddStringMapValues adds or overrides the given key-values into the string map
// that can be found at the specified fields.
func (m Manifest) AddStringMapValues(values map[string]string, fields ...string) error {
	curMap, _, err := unstructured.NestedStringMap(m.u.Object, fields...)
	if err != nil {
		return err
	}

	if curMap == nil {
		return unstructured.SetNestedStringMap(m.u.Object, values, fields...)
	}
	for k, v := range values {
		curMap[k] = v
	}
	return unstructured.SetNestedStringMap(m.u.Object, curMap, fields...)
}

func (m Manifest) GetSpec() (interface{}, error) {
	spec, ok, err := unstructured.NestedFieldNoCopy(m.u.Object, "spec")
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("spec was not found")
	}
	return spec, nil
}

func (m Manifest) SetStructuredSpec(spec interface{}) error {
	data, err := yaml.Marshal(spec)
	if err != nil {
		return err
	}

	unstructuredSpec := make(map[string]interface{})
	if err := yaml.Unmarshal(data, &unstructuredSpec); err != nil {
		return err
	}

	return unstructured.SetNestedField(m.u.Object, unstructuredSpec, "spec")
}

func (m Manifest) ConvertToStructuredObject(o interface{}) error {
	data, err := m.MarshalJSON()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, o)
}

func ParseFromStructuredObject(s interface{}) (Manifest, error) {
	data, err := json.Marshal(s)
	if err != nil {
		return Manifest{}, err
	}

	obj := &unstructured.Unstructured{}
	if err := obj.UnmarshalJSON(data); err != nil {
		return Manifest{}, err
	}

	return Manifest{
		Key: MakeResourceKey(obj),
		u:   obj,
	}, nil
}

func LoadPlainYAMLManifests(dir string, names []string, configFileName string) ([]Manifest, error) {
	// If no name was specified we have to walk the app directory to collect the manifest list.
	if len(names) == 0 {
		err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if path == dir {
				return nil
			}
			if f.IsDir() {
				return filepath.SkipDir
			}
			ext := filepath.Ext(f.Name())
			if ext != ".yaml" && ext != ".yml" && ext != ".json" {
				return nil
			}
			if model.IsApplicationConfigFile(f.Name()) {
				return nil
			}
			if f.Name() == configFileName {
				return nil
			}
			names = append(names, f.Name())
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	manifests := make([]Manifest, 0, len(names))
	for _, name := range names {
		path := filepath.Join(dir, name)
		ms, err := LoadManifestsFromYAMLFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to load manifest at %s (%w)", path, err)
		}
		manifests = append(manifests, ms...)
	}

	return manifests, nil
}

func LoadManifestsFromYAMLFile(path string) ([]Manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return ParseManifests(string(data))
}

func ParseManifests(data string) ([]Manifest, error) {
	const separator = "\n---"
	var (
		parts     = strings.Split(data, separator)
		manifests = make([]Manifest, 0, len(parts))
	)

	for i, part := range parts {
		// Ignore all the cases where no content between separator.
		if len(strings.TrimSpace(part)) == 0 {
			continue
		}
		// Append new line which trim by document separator.
		if i != len(parts)-1 {
			part += "\n"
		}
		var obj unstructured.Unstructured
		if err := yaml.Unmarshal([]byte(part), &obj); err != nil {
			return nil, err
		}
		if len(obj.Object) == 0 {
			continue
		}
		manifests = append(manifests, Manifest{
			Key: MakeResourceKey(&obj),
			u:   &obj,
		})
	}
	return manifests, nil
}
