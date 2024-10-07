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

import "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

// Manifest represents a Kubernetes resource manifest.
type Manifest struct {
	// TODO: define ResourceKey and add as a field here.
	Body *unstructured.Unstructured
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (m *Manifest) UnmarshalJSON(data []byte) error {
	m.Body = new(unstructured.Unstructured)
	return m.Body.UnmarshalJSON(data)
}
