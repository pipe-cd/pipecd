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

package kubernetes

import (
	"testing"

	"github.com/pipe-cd/pipecd/pkg/config"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
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

func Test_loader_determineNamespace(t *testing.T) {
	testcases := []struct {
		name                  string
		manifest              Manifest
		isNamespacedResources map[schema.GroupVersionKind]bool
		cfgK8sInput           config.KubernetesDeploymentInput
		want                  string
		wantErr               bool
	}{
		{
			name: "failed because unknown resource kind",
			manifest: Manifest{
				u: &unstructured.Unstructured{
					Object: map[string]interface{}{
						"apiVersion": "unknown",
						"kind":       "Unknown",
					},
				},
			},
			isNamespacedResources: map[schema.GroupVersionKind]bool{},
			want:                  "",
			wantErr:               true,
		},
		{
			name: "cluster-scoped resource: use '' as namespace",
			manifest: Manifest{
				u: &unstructured.Unstructured{
					Object: map[string]interface{}{
						"apiVersion": "v1",
						"kind":       "Namespace",
					},
				},
			},
			isNamespacedResources: map[schema.GroupVersionKind]bool{
				{Group: "", Version: "v1", Kind: "Namespace"}: false,
			},
			want:    "",
			wantErr: false,
		},
		{
			name: "cluster-scoped resource: use '' even though the app.pipecd.yaml has 'spec.input.namespace'",
			manifest: Manifest{
				u: &unstructured.Unstructured{
					Object: map[string]interface{}{
						"apiVersion": "v1",
						"kind":       "Namespace",
					},
				},
			},
			isNamespacedResources: map[schema.GroupVersionKind]bool{
				{Group: "", Version: "v1", Kind: "Namespace"}: false,
			},
			cfgK8sInput: config.KubernetesDeploymentInput{
				Namespace: "test",
			},
			want:    "",
			wantErr: false,
		},
		{
			name: "namespace-scoped resource: use the namespace set in the app.pipecd.yaml if it is not empty and the manifest has no namespace",
			manifest: Manifest{
				u: &unstructured.Unstructured{
					Object: map[string]interface{}{
						"apiVersion": "v1",
						"kind":       "Pod",
						"metadata":   map[string]interface{}{},
					},
				},
			},
			isNamespacedResources: map[schema.GroupVersionKind]bool{
				{Group: "", Version: "v1", Kind: "Pod"}: true,
			},
			cfgK8sInput: config.KubernetesDeploymentInput{
				Namespace: "inputNamespace",
			},
			want:    "inputNamespace",
			wantErr: false,
		},
		{
			name: "namespace-scoped resource: use the namespace set in the app.pipecd.yaml even though the manifest has a namespace",
			manifest: Manifest{
				u: &unstructured.Unstructured{
					Object: map[string]interface{}{
						"apiVersion": "v1",
						"kind":       "Pod",
						"metadata": map[string]interface{}{
							"namespace": "test",
						},
					},
				},
			},
			isNamespacedResources: map[schema.GroupVersionKind]bool{
				{Group: "", Version: "v1", Kind: "Pod"}: true,
			},
			cfgK8sInput: config.KubernetesDeploymentInput{
				Namespace: "inputNamespace",
			},
			want:    "inputNamespace",
			wantErr: false,
		},
		{
			name: "namespace-scoped resource: use 'default' namespace when input namespace is empty and the namespace in the app.pipecd.yaml is empty",
			manifest: Manifest{
				u: &unstructured.Unstructured{
					Object: map[string]interface{}{
						"apiVersion": "v1",
						"kind":       "Pod",
						"metadata":   map[string]interface{}{},
					},
				},
			},
			isNamespacedResources: map[schema.GroupVersionKind]bool{
				{Group: "", Version: "v1", Kind: "Pod"}: true,
			},
			cfgK8sInput: config.KubernetesDeploymentInput{},
			want:        "default",
			wantErr:     false,
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			l := &loader{
				isNamespacedResources: tc.isNamespacedResources,
				input:                 tc.cfgK8sInput,
			}
			err := l.determineNamespace(&tc.manifest)

			assert.Equal(t, tc.wantErr, err != nil)
			assert.Equal(t, tc.want, tc.manifest.Key.Namespace)
			assert.Equal(t, tc.want, tc.manifest.u.GetNamespace())
		})
	}
}
