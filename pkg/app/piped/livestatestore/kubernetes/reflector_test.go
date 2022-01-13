// Copyright 2020 The PipeCD Authors.
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
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/pipe-cd/pipecd/pkg/config"
)

func TestResourceMatcher(t *testing.T) {
	testcases := []struct {
		name string
		cfg  config.KubernetesAppStateInformer
		gvks map[schema.GroupVersionKind]bool
	}{
		{
			name: "empty config",
			cfg:  config.KubernetesAppStateInformer{},
			gvks: map[schema.GroupVersionKind]bool{
				schema.GroupVersionKind{"pipecd.dev", "v1beta1", "Foo"}:       false,
				schema.GroupVersionKind{"", "v1", "Foo"}:                      false,
				schema.GroupVersionKind{"", "v1", "Service"}:                  true,
				schema.GroupVersionKind{"networking.k8s.io", "v1", "Ingress"}: true,
			},
		},
		{
			name: "include config",
			cfg: config.KubernetesAppStateInformer{
				IncludeResources: []config.KubernetesResourceMatcher{
					{APIVersion: "pipecd.dev/v1beta1"},
					{APIVersion: "pipecd.dev/v1alpha1", Kind: "Foo"},
				},
			},
			gvks: map[schema.GroupVersionKind]bool{
				schema.GroupVersionKind{"pipecd.dev", "v1beta1", "Foo"}:  true,
				schema.GroupVersionKind{"pipecd.dev", "v1alpha1", "Foo"}: true,
				schema.GroupVersionKind{"pipecd.dev", "v1alpha1", "Bar"}: false,
			},
		},
		{
			name: "exclude config",
			cfg: config.KubernetesAppStateInformer{
				ExcludeResources: []config.KubernetesResourceMatcher{
					{APIVersion: "networking.k8s.io/v1"},
					{APIVersion: "apps/v1", Kind: "Deployment"},
				},
			},
			gvks: map[schema.GroupVersionKind]bool{
				schema.GroupVersionKind{"apps", "v1", "ReplicaSet"}:           true,
				schema.GroupVersionKind{"apps", "v1", "Deployment"}:           false,
				schema.GroupVersionKind{"networking.k8s.io", "v1", "Ingress"}: false,
			},
		},
	}

	for _, tc := range testcases {
		m := newResourceMatcher(tc.cfg)
		for gvk, expected := range tc.gvks {
			desc := fmt.Sprintf("%s: %v", tc.name, gvk)
			t.Run(desc, func(t *testing.T) {
				matched := m.Match(gvk)
				assert.Equal(t, expected, matched)
			})
		}
	}
}
