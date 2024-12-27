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

var builtinAPIGroups = map[string]struct{}{
	"admissionregistration.k8s.io": {},
	"apiextensions.k8s.io":         {},
	"apiregistration.k8s.io":       {},
	"apps":                         {},
	"authentication.k8s.io":        {},
	"authorization.k8s.io":         {},
	"autoscaling":                  {},
	"batch":                        {},
	"certificates.k8s.io":          {},
	"coordination.k8s.io":          {},
	"extensions":                   {},
	"internal.autoscaling.k8s.io":  {},
	"metrics.k8s.io":               {},
	"networking.k8s.io":            {},
	"node.k8s.io":                  {},
	"policy":                       {},
	"rbac.authorization.k8s.io":    {},
	"scheduling.k8s.io":            {},
	"storage.k8s.io":               {},
	"":                             {},
}

func isBuiltinAPIGroup(apiGroup string) bool {
	_, ok := builtinAPIGroups[apiGroup]
	return ok
}

// Manifest represents a Kubernetes resource manifest.
type Manifest struct {
	body *unstructured.Unstructured
}

func (m Manifest) Key() ResourceKey {
	return makeResourceKey(m.body)
}

func (m Manifest) Kind() string {
	return m.body.GetKind()
}

func (m Manifest) Name() string {
	return m.body.GetName()
}

// IsDeployment returns true if the manifest is a Deployment.
// It checks the API group and the kind of the manifest.
func (m Manifest) IsDeployment() bool {
	// TODO: check the API group more strictly.
	return isBuiltinAPIGroup(m.body.GroupVersionKind().Group) && m.body.GetKind() == KindDeployment
}

// IsSecret returns true if the manifest is a Secret.
// It checks the API group and the kind of the manifest.
func (m Manifest) IsSecret() bool {
	// TODO: check the API group more strictly.
	return isBuiltinAPIGroup(m.body.GroupVersionKind().Group) && m.body.GetKind() == KindSecret
}

// IsConfigMap returns true if the manifest is a ConfigMap.
// It checks the API group and the kind of the manifest.
func (m Manifest) IsConfigMap() bool {
	// TODO: check the API group more strictly.
	return isBuiltinAPIGroup(m.body.GroupVersionKind().Group) && m.body.GetKind() == KindConfigMap
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (m *Manifest) UnmarshalJSON(data []byte) error {
	m.body = new(unstructured.Unstructured)
	return m.body.UnmarshalJSON(data)
}

// MarshalJSON implements the json.Marshaler interface.
// It marshals the underlying unstructured.Unstructured object into JSON bytes.
func (m *Manifest) MarshalJSON() ([]byte, error) {
	return m.body.MarshalJSON()
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
	return yaml.Marshal(m.body)
}

func (m Manifest) GetAnnotations() map[string]string {
	return m.body.GetAnnotations()
}

func (m Manifest) NestedMap(fields ...string) (map[string]any, bool, error) {
	return unstructured.NestedMap(m.body.Object, fields...)
}

func (m Manifest) AddLabels(labels map[string]string) {
	if len(labels) == 0 {
		return
	}

	lbs := m.body.GetLabels()
	if lbs == nil {
		m.body.SetLabels(labels)
		return
	}
	for k, v := range labels {
		lbs[k] = v
	}
	m.body.SetLabels(lbs)
}

func (m Manifest) AddAnnotations(annotations map[string]string) {
	if len(annotations) == 0 {
		return
	}

	annos := m.body.GetAnnotations()
	if annos == nil {
		m.body.SetAnnotations(annotations)
		return
	}
	for k, v := range annotations {
		annos[k] = v
	}
	m.body.SetAnnotations(annos)
}

// IsManagedByPiped returns true if the manifest is managed by Piped.
func (m Manifest) IsManagedByPiped() bool {
	return len(m.body.GetOwnerReferences()) == 0 && m.body.GetAnnotations()[LabelManagedBy] == ManagedByPiped
}

// AddStringMapValues adds or overrides the given key-values into the string map
// that can be found at the specified fields.
func (m Manifest) AddStringMapValues(values map[string]string, fields ...string) error {
	curMap, _, err := unstructured.NestedStringMap(m.body.Object, fields...)
	if err != nil {
		return err
	}

	if curMap == nil {
		return unstructured.SetNestedStringMap(m.body.Object, values, fields...)
	}
	maps.Copy(curMap, values)
	return unstructured.SetNestedStringMap(m.body.Object, curMap, fields...)
}

// FindConfigsAndSecrets returns the manifests that are ConfigMap or Secret.
func FindConfigsAndSecrets(manifests []Manifest) map[ResourceKey]Manifest {
	configs := make(map[ResourceKey]Manifest)
	for _, m := range manifests {
		if m.IsConfigMap() {
			configs[m.Key()] = m
		}
		if m.IsSecret() {
			configs[m.Key()] = m
		}
	}
	return configs
}

// WorkloadPair represents a pair of old and new manifests.
type WorkloadPair struct {
	Old Manifest
	New Manifest
}

// FindSameManifests returns the pairs of old and new manifests that have the same key.
func FindSameManifests(olds, news []Manifest) []WorkloadPair {
	pairs := make([]WorkloadPair, 0)
	oldMap := make(map[ResourceKey]Manifest, len(olds))
	for _, m := range olds {
		key := m.Key().normalize()
		oldMap[key] = m
	}
	for _, n := range news {
		key := n.Key().normalize()
		if o, ok := oldMap[key]; ok {
			pairs = append(pairs, WorkloadPair{
				Old: o,
				New: n,
			})
		}
	}
	return pairs
}
