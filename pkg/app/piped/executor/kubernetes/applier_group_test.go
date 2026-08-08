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
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	provider "github.com/pipe-cd/pipecd/pkg/app/piped/platformprovider/kubernetes"
	"github.com/pipe-cd/pipecd/pkg/app/piped/platformprovider/kubernetes/kubernetestest"
	"github.com/pipe-cd/pipecd/pkg/config"
)

func TestMakeKeyFromProviderLabels(t *testing.T) {
	testcases := []struct {
		name   string
		labels map[string]string
		want   string
	}{
		{
			name: "empty",
			want: "",
		},
		{
			name: "one label",
			labels: map[string]string{
				"foo": "foo-1",
			},
			want: "foo:foo-1",
		},
		{
			name: "multiple labels",
			labels: map[string]string{
				"foo": "foo-1",
				"bar": "bar-1",
			},
			want: "bar:bar-1,foo:foo-1",
		},
		{
			name: "multiple labels in the reverse order",
			labels: map[string]string{
				"bar": "bar-1",
				"foo": "foo-1",
			},
			want: "bar:bar-1,foo:foo-1",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := makeKeyFromProviderLabels(tc.labels)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestApplierGroupGet(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	defaultApplier := kubernetestest.NewMockApplier(ctrl)
	applierA := kubernetestest.NewMockApplier(ctrl)
	applierB := kubernetestest.NewMockApplier(ctrl)

	testcases := []struct {
		name           string
		resourceRoutes []config.KubernetesResourceRoute
		appliers       map[string]provider.Applier
		labeledCps     map[string][]string
		resourceKey    provider.ResourceKey
		wantApplier    provider.Applier
		wantErr        error
	}{
		{
			name:           "empty routes, fallback to default",
			resourceRoutes: nil,
			appliers:       map[string]provider.Applier{"default": defaultApplier},
			resourceKey:    provider.ResourceKey{Kind: "Deployment", Name: "app"},
			wantApplier:    defaultApplier,
			wantErr:        nil,
		},
		{
			name: "matched route by provider name",
			resourceRoutes: []config.KubernetesResourceRoute{
				{
					Match: &config.KubernetesResourceRouteMatcher{
						Kind: "Deployment",
					},
					Provider: config.KubernetesProviderMatcher{
						Name: "provider-a",
					},
				},
			},
			appliers: map[string]provider.Applier{
				"provider-a": applierA,
			},
			resourceKey: provider.ResourceKey{Kind: "Deployment", Name: "app"},
			wantApplier: applierA,
			wantErr:     nil,
		},
		{
			name: "matched route by name but provider applier not found",
			resourceRoutes: []config.KubernetesResourceRoute{
				{
					Match: &config.KubernetesResourceRouteMatcher{
						Kind: "Deployment",
					},
					Provider: config.KubernetesProviderMatcher{
						Name: "provider-a",
					},
				},
			},
			appliers:    map[string]provider.Applier{},
			resourceKey: provider.ResourceKey{Kind: "Deployment", Name: "app"},
			wantApplier: nil,
			wantErr:     fmt.Errorf("provider provider-a specified in resourceRoutes was not found"),
		},
		{
			name: "matched route by provider labels",
			resourceRoutes: []config.KubernetesResourceRoute{
				{
					Match: &config.KubernetesResourceRouteMatcher{
						Kind: "Service",
					},
					Provider: config.KubernetesProviderMatcher{
						Labels: map[string]string{"env": "prod"},
					},
				},
			},
			appliers: map[string]provider.Applier{
				"provider-a": applierA,
				"provider-b": applierB,
			},
			labeledCps: map[string][]string{
				"env:prod": {"provider-a", "provider-b"},
			},
			resourceKey: provider.ResourceKey{Kind: "Service", Name: "srv"},
			wantApplier: provider.NewMultiApplier(applierA, applierB),
			wantErr:     nil,
		},
		{
			name: "matched route by labels but no matching labeled providers",
			resourceRoutes: []config.KubernetesResourceRoute{
				{
					Match: &config.KubernetesResourceRouteMatcher{
						Kind: "Service",
					},
					Provider: config.KubernetesProviderMatcher{
						Labels: map[string]string{"env": "prod"},
					},
				},
			},
			appliers:    map[string]provider.Applier{},
			labeledCps:  map[string][]string{},
			resourceKey: provider.ResourceKey{Kind: "Service", Name: "srv"},
			wantApplier: nil,
			wantErr:     fmt.Errorf("there are no provider that matches the specified labels (map[env:prod])"),
		},
		{
			name: "matched route by labels but applier for matched provider not found",
			resourceRoutes: []config.KubernetesResourceRoute{
				{
					Match: &config.KubernetesResourceRouteMatcher{
						Kind: "Service",
					},
					Provider: config.KubernetesProviderMatcher{
						Labels: map[string]string{"env": "prod"},
					},
				},
			},
			appliers: map[string]provider.Applier{
				"provider-a": applierA,
			},
			labeledCps: map[string][]string{
				"env:prod": {"provider-a", "provider-b"},
			},
			resourceKey: provider.ResourceKey{Kind: "Service", Name: "srv"},
			wantApplier: nil,
			wantErr:     fmt.Errorf("provider provider-b specified in resourceRoutes was not found"),
		},
		{
			name: "route matcher kind mismatch, fallback to default",
			resourceRoutes: []config.KubernetesResourceRoute{
				{
					Match: &config.KubernetesResourceRouteMatcher{
						Kind: "Deployment",
					},
					Provider: config.KubernetesProviderMatcher{
						Name: "provider-a",
					},
				},
			},
			appliers: map[string]provider.Applier{
				"provider-a": applierA,
			},
			resourceKey: provider.ResourceKey{Kind: "Service", Name: "srv"},
			wantApplier: defaultApplier,
			wantErr:     nil,
		},
		{
			name: "route matcher name mismatch, fallback to default",
			resourceRoutes: []config.KubernetesResourceRoute{
				{
					Match: &config.KubernetesResourceRouteMatcher{
						Name: "app-prod",
					},
					Provider: config.KubernetesProviderMatcher{
						Name: "provider-a",
					},
				},
			},
			appliers: map[string]provider.Applier{
				"provider-a": applierA,
			},
			resourceKey: provider.ResourceKey{Kind: "Deployment", Name: "app-dev"},
			wantApplier: defaultApplier,
			wantErr:     nil,
		},
		{
			name: "route matcher nil matches any resource",
			resourceRoutes: []config.KubernetesResourceRoute{
				{
					Match: nil,
					Provider: config.KubernetesProviderMatcher{
						Name: "provider-a",
					},
				},
			},
			appliers: map[string]provider.Applier{
				"provider-a": applierA,
			},
			resourceKey: provider.ResourceKey{Kind: "ConfigMap", Name: "cfg"},
			wantApplier: applierA,
			wantErr:     nil,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ag := &applierGroup{
				resourceRoutes:   tc.resourceRoutes,
				appliers:         tc.appliers,
				labeledProviders: tc.labeledCps,
				defaultApplier:   defaultApplier,
			}

			got, err := ag.Get(tc.resourceKey)
			if tc.wantErr != nil {
				assert.EqualError(t, err, tc.wantErr.Error())
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.wantApplier, got)
			}
		})
	}
}
