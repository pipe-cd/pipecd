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
	"time"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/pipe-cd/pipecd/pkg/model"
)

func TestBuildApplicationLiveState(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name      string
		manifests []Manifest
		want      *model.ApplicationLiveState
	}{
		{
			name: "single pod",
			manifests: []Manifest{
				{
					body: &unstructured.Unstructured{
						Object: map[string]interface{}{
							"apiVersion": "v1",
							"kind":       "Pod",
							"metadata": map[string]interface{}{
								"name":              "test-pod",
								"namespace":         "default",
								"uid":               "test-uid",
								"creationTimestamp": now.Format(time.RFC3339),
							},
						},
					},
				},
			},
			want: &model.ApplicationLiveState{
				Resources: []*model.ResourceState{
					{
						Id:           "test-uid",
						Name:         "test-pod",
						ResourceType: "Pod",
						ResourceMetadata: map[string]string{
							"Namespace":   "default",
							"API Version": "v1",
							"Kind":        "Pod",
						},
						CreatedAt: now.Unix(),
						UpdatedAt: now.Unix(),
					},
				},
				HealthStatus: model.ApplicationLiveState_UNKNOWN,
			},
		},
		{
			name: "single pod with owner references",
			manifests: []Manifest{
				{
					body: &unstructured.Unstructured{
						Object: map[string]interface{}{
							"apiVersion": "v1",
							"kind":       "Pod",
							"metadata": map[string]interface{}{
								"name":              "test-pod",
								"namespace":         "default",
								"uid":               "test-uid",
								"creationTimestamp": now.Format(time.RFC3339),
								"ownerReferences": []interface{}{
									map[string]interface{}{
										"uid": "owner-uid",
									},
								},
							},
						},
					},
				},
			},
			want: &model.ApplicationLiveState{
				Resources: []*model.ResourceState{
					{
						Id:           "test-uid",
						Name:         "test-pod",
						ResourceType: "Pod",
						ResourceMetadata: map[string]string{
							"Namespace":   "default",
							"API Version": "v1",
							"Kind":        "Pod",
						},
						ParentIds: []string{"owner-uid"},
						CreatedAt: now.Unix(),
						UpdatedAt: now.Unix(),
					},
				},
				HealthStatus: model.ApplicationLiveState_UNKNOWN,
			},
		},
		{
			name: "multiple resources with owner references",
			manifests: []Manifest{
				{
					body: &unstructured.Unstructured{
						Object: map[string]interface{}{
							"apiVersion": "v1",
							"kind":       "Pod",
							"metadata": map[string]interface{}{
								"name":              "test-pod-1",
								"namespace":         "default",
								"uid":               "test-uid-1",
								"creationTimestamp": now.Format(time.RFC3339),
								"ownerReferences": []interface{}{
									map[string]interface{}{
										"uid": "owner-uid-1",
									},
								},
							},
						},
					},
				},
				{
					body: &unstructured.Unstructured{
						Object: map[string]interface{}{
							"apiVersion": "v1",
							"kind":       "Service",
							"metadata": map[string]interface{}{
								"name":              "test-service",
								"namespace":         "default",
								"uid":               "test-uid-2",
								"creationTimestamp": now.Format(time.RFC3339),
								"ownerReferences": []interface{}{
									map[string]interface{}{
										"uid": "owner-uid-2",
									},
								},
							},
						},
					},
				},
			},
			want: &model.ApplicationLiveState{
				Resources: []*model.ResourceState{
					{
						Id:           "test-uid-1",
						Name:         "test-pod-1",
						ResourceType: "Pod",
						ResourceMetadata: map[string]string{
							"Namespace":   "default",
							"API Version": "v1",
							"Kind":        "Pod",
						},
						ParentIds: []string{"owner-uid-1"},
						CreatedAt: now.Unix(),
						UpdatedAt: now.Unix(),
					},
					{
						Id:           "test-uid-2",
						Name:         "test-service",
						ResourceType: "Service",
						ResourceMetadata: map[string]string{
							"Namespace":   "default",
							"API Version": "v1",
							"Kind":        "Service",
						},
						ParentIds: []string{"owner-uid-2"},
						CreatedAt: now.Unix(),
						UpdatedAt: now.Unix(),
					},
				},
				HealthStatus: model.ApplicationLiveState_UNKNOWN,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BuildApplicationLiveState("test-deploytarget", tt.manifests, now)
			assert.Equal(t, tt.want, got, "expected live state to be equal to the expected one")
		})
	}
}
