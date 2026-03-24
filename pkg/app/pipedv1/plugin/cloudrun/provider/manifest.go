// Copyright 2026 The PipeCD Authors.
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
	"fmt"
	"os"

	"google.golang.org/api/run/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/yaml"
)

type ServiceManifest struct {
	u *unstructured.Unstructured
}

func (m ServiceManifest) Name() string {
	return m.u.GetName()
}

func (m ServiceManifest) SetRevision(name string) error {
	return unstructured.SetNestedField(m.u.Object, name, "spec", "template", "metadata", "name")
}

type RevisionTraffic struct {
	RevisionName string `json:"revisionName"`
	Percent      int    `json:"percent"`
}

func (m ServiceManifest) UpdateTraffic(revisions []RevisionTraffic) error {
	items := []interface{}{}
	for i := range revisions {
		out, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&revisions[i])
		if err != nil {
			return fmt.Errorf("unable to set traffic for object: %w", err)
		}
		items = append(items, out)
	}

	return unstructured.SetNestedSlice(m.u.Object, items, "spec", "traffic")
}

func (m ServiceManifest) UpdateAllTraffic(revision string) error {
	return m.UpdateTraffic([]RevisionTraffic{
		{
			RevisionName: revision,
			Percent:      100,
		},
	})
}

func (m ServiceManifest) YamlBytes() ([]byte, error) {
	return yaml.Marshal(m.u)
}

func (m ServiceManifest) AddLabels(labels map[string]string) {
	if len(labels) == 0 {
		return
	}

	lbls := m.u.GetLabels()
	if lbls == nil {
		m.u.SetLabels(labels)
		return
	}
	for k, v := range labels {
		lbls[k] = v
	}
	m.u.SetLabels(lbls)
}

func (m ServiceManifest) AddRevisionLabels(labels map[string]string) error {
	if len(labels) == 0 {
		return nil
	}

	fields := []string{"spec", "template", "metadata", "labels"}
	lbls, ok, err := unstructured.NestedStringMap(m.u.Object, fields...)
	if err != nil {
		return err
	}
	if !ok {
		return unstructured.SetNestedStringMap(m.u.Object, labels, fields...)
	}

	for k, v := range labels {
		lbls[k] = v
	}
	return unstructured.SetNestedStringMap(m.u.Object, lbls, fields...)
}

func (m ServiceManifest) RunService() (*run.Service, error) {
	data, err := m.YamlBytes()
	if err != nil {
		return nil, err
	}

	var s run.Service
	if err := yaml.Unmarshal(data, &s); err != nil {
		return nil, err
	}
	return &s, nil
}

func LoadServiceManifest(path string) (ServiceManifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return ServiceManifest{}, err
	}
	return ParseServiceManifest(data)
}

func ParseServiceManifest(data []byte) (ServiceManifest, error) {
	var obj unstructured.Unstructured
	if err := yaml.Unmarshal(data, &obj); err != nil {
		return ServiceManifest{}, err
	}

	return ServiceManifest{
		u: &obj,
	}, nil
}

func DiffManifests(a, b ServiceManifest) (bool, string) {
	// Stub implementation for prototype
	return true, "diff not implemented"
}

func (m ServiceManifest) ExtractImages() ([]string, error) {
	containers, ok, err := unstructured.NestedSlice(m.u.Object, "spec", "template", "spec", "containers")
	if err != nil {
		return nil, err
	}
	if !ok || len(containers) == 0 {
		return nil, fmt.Errorf("spec.template.spec.containers was missing")
	}

	var images []string
	for _, c := range containers {
		container, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&c)
		if err != nil {
			return nil, fmt.Errorf("invalid container format")
		}
		image, ok, err := unstructured.NestedString(container, "image")
		if err != nil {
			return nil, err
		}
		if !ok || image == "" {
			return nil, fmt.Errorf("image was missing")
		}
		images = append(images, image)
	}
	return images, nil
}
