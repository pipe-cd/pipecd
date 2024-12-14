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
	"testing"

	"github.com/google/go-cmp/cmp"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestManifest_AddStringMapValues(t *testing.T) {
	tests := []struct {
		name     string
		initial  map[string]interface{}
		values   map[string]string
		fields   []string
		expected map[string]interface{}
	}{
		{
			name: "add new values to empty map",
			initial: map[string]interface{}{
				"metadata": map[string]interface{}{
					"annotations": map[string]interface{}{},
				},
			},
			values: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
			fields: []string{"metadata", "annotations"},
			expected: map[string]interface{}{
				"metadata": map[string]interface{}{
					"annotations": map[string]interface{}{
						"key1": "value1",
						"key2": "value2",
					},
				},
			},
		},
		{
			name: "override existing values",
			initial: map[string]interface{}{
				"metadata": map[string]interface{}{
					"annotations": map[string]interface{}{
						"key1": "oldvalue1",
					},
				},
			},
			values: map[string]string{
				"key1": "newvalue1",
				"key2": "value2",
			},
			fields: []string{"metadata", "annotations"},
			expected: map[string]interface{}{
				"metadata": map[string]interface{}{
					"annotations": map[string]interface{}{
						"key1": "newvalue1",
						"key2": "value2",
					},
				},
			},
		},
		{
			name: "add values to non-existing map",
			initial: map[string]interface{}{
				"metadata": map[string]interface{}{},
			},
			values: map[string]string{
				"key1": "value1",
			},
			fields: []string{"metadata", "annotations"},
			expected: map[string]interface{}{
				"metadata": map[string]interface{}{
					"annotations": map[string]interface{}{
						"key1": "value1",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manifest := Manifest{
				Body: &unstructured.Unstructured{
					Object: tt.initial,
				},
			}
			err := manifest.AddStringMapValues(tt.values, tt.fields...)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if diff := cmp.Diff(tt.expected, manifest.Body.Object); diff != "" {
				t.Errorf("unexpected result (-want +got):\n%s", diff)
			}
		})
	}
}
