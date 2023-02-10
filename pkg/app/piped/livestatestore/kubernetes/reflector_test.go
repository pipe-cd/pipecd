// Copyright 2023 The PipeCD Authors.
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
	t.Parallel()

	testcases := []struct {
		name string
		cfg  config.KubernetesAppStateInformer
		gvks map[schema.GroupVersionKind]bool
	}{
		{
			name: "empty config",
			cfg:  config.KubernetesAppStateInformer{},
			gvks: map[schema.GroupVersionKind]bool{
				{"pipecd.dev", "v1beta1", "Foo"}:       false,
				{"", "v1", "Foo"}:                      false,
				{"", "v1", "Service"}:                  true,
				{"networking.k8s.io", "v1", "Ingress"}: true,
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
				{"pipecd.dev", "v1beta1", "Foo"}:  true,
				{"pipecd.dev", "v1alpha1", "Foo"}: true,
				{"pipecd.dev", "v1alpha1", "Bar"}: false,
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
				{"apps", "v1", "ReplicaSet"}:           true,
				{"apps", "v1", "Deployment"}:           false,
				{"networking.k8s.io", "v1", "Ingress"}: false,
			},
		},
	}

	for _, tc := range testcases {
		tc := tc
		m := newResourceMatcher(tc.cfg)
		for gvk, expected := range tc.gvks {
			desc := fmt.Sprintf("%s: %v", tc.name, gvk)
			gvk, expected := gvk, expected
			t.Run(desc, func(t *testing.T) {
				t.Parallel()

				matched := m.Match(gvk)
				assert.Equal(t, expected, matched)
			})
		}
	}
}
