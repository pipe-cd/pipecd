// Copyright 2025 The PipeCD Authors.
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

package store

import (
	"testing"

	"k8s.io/apimachinery/pkg/runtime/schema"

	kubeconfig "github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/config"
)

func TestResourceMatcher_matchGVK(t *testing.T) {
	tests := []struct {
		name     string
		config   kubeconfig.KubernetesAppStateInformer
		gvk      schema.GroupVersionKind
		expected bool
	}{
		{
			name: "exclude by APIVersion only",
			config: kubeconfig.KubernetesAppStateInformer{
				ExcludeResources: []kubeconfig.KubernetesResourceMatcher{
					{APIVersion: "apps/v1"},
				},
			},
			gvk:      schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "Deployment"},
			expected: false,
		},
		{
			name: "exclude by APIVersion and Kind",
			config: kubeconfig.KubernetesAppStateInformer{
				ExcludeResources: []kubeconfig.KubernetesResourceMatcher{
					{APIVersion: "apps/v1", Kind: "Deployment"},
				},
			},
			gvk:      schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "Deployment"},
			expected: false,
		},
		{
			name: "exclude doesn't match different Kind",
			config: kubeconfig.KubernetesAppStateInformer{
				ExcludeResources: []kubeconfig.KubernetesResourceMatcher{
					{APIVersion: "apps/v1", Kind: "Deployment"},
				},
			},
			gvk:      schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "StatefulSet"},
			expected: true, // Should fall back to whitelist check
		},
		{
			name: "include by APIVersion only",
			config: kubeconfig.KubernetesAppStateInformer{
				IncludeResources: []kubeconfig.KubernetesResourceMatcher{
					{APIVersion: "custom.io/v1alpha1"},
				},
			},
			gvk:      schema.GroupVersionKind{Group: "custom.io", Version: "v1alpha1", Kind: "CustomResource"},
			expected: true,
		},
		{
			name: "include by APIVersion and Kind",
			config: kubeconfig.KubernetesAppStateInformer{
				IncludeResources: []kubeconfig.KubernetesResourceMatcher{
					{APIVersion: "custom.io/v1alpha1", Kind: "CustomResource"},
				},
			},
			gvk:      schema.GroupVersionKind{Group: "custom.io", Version: "v1alpha1", Kind: "CustomResource"},
			expected: true,
		},
		{
			name: "include doesn't match different Kind",
			config: kubeconfig.KubernetesAppStateInformer{
				IncludeResources: []kubeconfig.KubernetesResourceMatcher{
					{APIVersion: "custom.io/v1alpha1", Kind: "CustomResource"},
				},
			},
			gvk:      schema.GroupVersionKind{Group: "custom.io", Version: "v1alpha1", Kind: "AnotherResource"},
			expected: false, // Not in whitelist either
		},
		{
			name:     "whitelisted resource - core v1 Pod",
			config:   kubeconfig.KubernetesAppStateInformer{},
			gvk:      schema.GroupVersionKind{Group: "", Version: "v1", Kind: "Pod"},
			expected: true,
		},
		{
			name:     "whitelisted resource - apps/v1 Deployment",
			config:   kubeconfig.KubernetesAppStateInformer{},
			gvk:      schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "Deployment"},
			expected: true,
		},
		{
			name:     "non-whitelisted Kind",
			config:   kubeconfig.KubernetesAppStateInformer{},
			gvk:      schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "UnknownKind"},
			expected: false,
		},
		{
			name:     "non-whitelisted Group",
			config:   kubeconfig.KubernetesAppStateInformer{},
			gvk:      schema.GroupVersionKind{Group: "unknown.io", Version: "v1", Kind: "Pod"},
			expected: false,
		},
		{
			name:     "non-whitelisted Version",
			config:   kubeconfig.KubernetesAppStateInformer{},
			gvk:      schema.GroupVersionKind{Group: "apps", Version: "v3", Kind: "Deployment"},
			expected: false,
		},
		{
			name: "exclude takes precedence over include",
			config: kubeconfig.KubernetesAppStateInformer{
				IncludeResources: []kubeconfig.KubernetesResourceMatcher{
					{APIVersion: "apps/v1"},
				},
				ExcludeResources: []kubeconfig.KubernetesResourceMatcher{
					{APIVersion: "apps/v1", Kind: "Deployment"},
				},
			},
			gvk:      schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "Deployment"},
			expected: false,
		},
		{
			name: "include different Kind in same APIVersion",
			config: kubeconfig.KubernetesAppStateInformer{
				IncludeResources: []kubeconfig.KubernetesResourceMatcher{
					{APIVersion: "apps/v1"},
				},
				ExcludeResources: []kubeconfig.KubernetesResourceMatcher{
					{APIVersion: "apps/v1", Kind: "Deployment"},
				},
			},
			gvk:      schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "StatefulSet"},
			expected: true, // Included by APIVersion, not excluded by specific Kind
		},
		{
			name:     "v1beta1 version is whitelisted",
			config:   kubeconfig.KubernetesAppStateInformer{},
			gvk:      schema.GroupVersionKind{Group: "extensions", Version: "v1beta1", Kind: "Ingress"},
			expected: true,
		},
		{
			name:     "v1beta2 version is whitelisted",
			config:   kubeconfig.KubernetesAppStateInformer{},
			gvk:      schema.GroupVersionKind{Group: "autoscaling", Version: "v1beta2", Kind: "HorizontalPodAutoscaler"},
			expected: true,
		},
		{
			name:     "v2 version is whitelisted",
			config:   kubeconfig.KubernetesAppStateInformer{},
			gvk:      schema.GroupVersionKind{Group: "autoscaling", Version: "v2", Kind: "HorizontalPodAutoscaler"},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matcher := newResourceMatcher(tt.config)
			result := matcher.matchGVK(tt.gvk)
			if result != tt.expected {
				t.Errorf("matchGVK() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestNewResourceMatcher(t *testing.T) {
	config := kubeconfig.KubernetesAppStateInformer{
		IncludeResources: []kubeconfig.KubernetesResourceMatcher{
			{APIVersion: "apps/v1"},
			{APIVersion: "custom.io/v1alpha1", Kind: "CustomResource"},
		},
		ExcludeResources: []kubeconfig.KubernetesResourceMatcher{
			{APIVersion: "batch/v1"},
			{APIVersion: "apps/v1", Kind: "Deployment"},
		},
	}

	matcher := newResourceMatcher(config)

	// Test includes
	if _, ok := matcher.includes["apps/v1"]; !ok {
		t.Error("Expected 'apps/v1' to be in includes")
	}
	if _, ok := matcher.includes["custom.io/v1alpha1:CustomResource"]; !ok {
		t.Error("Expected 'custom.io/v1alpha1:CustomResource' to be in includes")
	}

	// Test excludes
	if _, ok := matcher.excludes["batch/v1"]; !ok {
		t.Error("Expected 'batch/v1' to be in excludes")
	}
	if _, ok := matcher.excludes["apps/v1:Deployment"]; !ok {
		t.Error("Expected 'apps/v1:Deployment' to be in excludes")
	}

	// Test counts
	if len(matcher.includes) != 2 {
		t.Errorf("Expected 2 includes, got %d", len(matcher.includes))
	}
	if len(matcher.excludes) != 2 {
		t.Errorf("Expected 2 excludes, got %d", len(matcher.excludes))
	}
}
