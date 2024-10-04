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

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestSortManifests(t *testing.T) {
	maker := func(name string, annotations map[string]string) Manifest {
		m := Manifest{
			Key: ResourceKey{Name: name},
			u: &unstructured.Unstructured{
				Object: map[string]interface{}{},
			},
		}
		m.AddAnnotations(annotations)
		return m
	}

	testcases := []struct {
		name      string
		manifests []Manifest
		want      []Manifest
	}{
		{
			name: "empty",
		},
		{
			name: "one manifest",
			manifests: []Manifest{
				maker("name-1", map[string]string{AnnotationOrder: "0"}),
			},
			want: []Manifest{
				maker("name-1", map[string]string{AnnotationOrder: "0"}),
			},
		},
		{
			name: "multiple manifests",
			manifests: []Manifest{
				maker("name-2", map[string]string{AnnotationOrder: "2"}),
				maker("name--1", map[string]string{AnnotationOrder: "-1"}),
				maker("name-nil", nil),
				maker("name-0", map[string]string{AnnotationOrder: "0"}),
				maker("name-1", map[string]string{AnnotationOrder: "1"}),
			},
			want: []Manifest{
				maker("name--1", map[string]string{AnnotationOrder: "-1"}),
				maker("name-nil", nil),
				maker("name-0", map[string]string{AnnotationOrder: "0"}),
				maker("name-1", map[string]string{AnnotationOrder: "1"}),
				maker("name-2", map[string]string{AnnotationOrder: "2"}),
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			sortManifests(tc.manifests)
			assert.Equal(t, tc.want, tc.manifests)
		})
	}
}
