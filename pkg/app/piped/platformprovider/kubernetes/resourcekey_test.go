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

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestGetResourceKeyFromActualResource(t *testing.T) {
	type args struct {
		resource *unstructured.Unstructured
	}

	tests := []struct {
		name string
		args args
		want ResourceKey
	}{
		{
			name: "get resource key from annotation",
			args: args{
				resource: &unstructured.Unstructured{
					Object: map[string]interface{}{
						"apiVersion": "v1",
						"kind":       "Pod",
						"metadata": map[string]interface{}{
							"name": "my-pod",
							"annotations": map[string]interface{}{
								"pipecd.dev/resource-key": "v1:Pod:my-namespace:my-pod",
							},
						},
					},
				},
			},
			want: ResourceKey{
				APIVersion: "v1",
				Kind:       "Pod",
				Namespace:  "my-namespace",
				Name:       "my-pod",
			},
		},
		{
			name: "prioritize the resource key from annotation when both metadata.annotation and metadata.namespace are available",
			args: args{
				resource: &unstructured.Unstructured{
					Object: map[string]interface{}{
						"apiVersion": "v1",
						"kind":       "Pod",
						"metadata": map[string]interface{}{
							"name":      "my-pod",
							"namespace": "test",
							"annotations": map[string]interface{}{
								"pipecd.dev/resource-key": "v1:Pod:my-namespace:my-pod",
							},
						},
					},
				},
			},
			want: ResourceKey{
				APIVersion: "v1",
				Kind:       "Pod",
				Namespace:  "my-namespace",
				Name:       "my-pod",
			},
		},
		{
			name: "make resource key when annotation is not available",
			args: args{
				resource: &unstructured.Unstructured{
					Object: map[string]interface{}{
						"apiVersion": "v1",
						"kind":       "Pod",
						"metadata": map[string]interface{}{
							"name":      "my-pod",
							"namespace": "my-namespace",
						},
					},
				},
			},
			want: ResourceKey{
				APIVersion: "v1",
				Kind:       "Pod",
				Namespace:  "my-namespace",
				Name:       "my-pod",
			},
		},
		{
			name: "make resource key when the value of annotation is invalid",
			args: args{
				resource: &unstructured.Unstructured{
					Object: map[string]interface{}{
						"apiVersion": "v1",
						"kind":       "Pod",
						"metadata": map[string]interface{}{
							"name":      "my-pod",
							"namespace": "my-namespace",
							"annotations": map[string]interface{}{
								"pipecd.dev/resource-key": "invalid-value",
							},
						},
					},
				},
			},
			want: ResourceKey{
				APIVersion: "v1",
				Kind:       "Pod",
				Namespace:  "my-namespace",
				Name:       "my-pod",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got := GetResourceKeyFromActualResource(tt.args.resource)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMakeResourceKey(t *testing.T) {
	type args struct {
		resource *unstructured.Unstructured
	}

	tests := []struct {
		name string
		args args
		want ResourceKey
	}{
		{
			name: "default case",
			args: args{
				resource: &unstructured.Unstructured{
					Object: map[string]interface{}{
						"apiVersion": "v1",
						"kind":       "Pod",
						"metadata": map[string]interface{}{
							"name":      "my-pod",
							"namespace": "my-namespace",
						},
					},
				},
			},
			want: ResourceKey{
				APIVersion: "v1",
				Kind:       "Pod",
				Namespace:  "my-namespace",
				Name:       "my-pod",
			},
		},
		{
			name: "use 'default' when namespace is empty",
			args: args{
				resource: &unstructured.Unstructured{
					Object: map[string]interface{}{
						"apiVersion": "v1",
						"kind":       "Pod",
						"metadata": map[string]interface{}{
							"name": "my-pod",
						},
					},
				},
			},
			want: ResourceKey{
				APIVersion: "v1",
				Kind:       "Pod",
				Namespace:  DefaultNamespace,
				Name:       "my-pod",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got := MakeResourceKey(tt.args.resource)
			assert.Equal(t, tt.want, got)
		})
	}
}
