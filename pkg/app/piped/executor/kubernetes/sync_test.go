package kubernetes

import (
	"testing"

	"github.com/stretchr/testify/assert"

	provider "github.com/pipe-cd/pipe/pkg/app/piped/cloudprovider/kubernetes"
)

func TestFindRemoveResources(t *testing.T) {
	tests := []struct {
		name          string
		manifests     []provider.Manifest
		liveResources []provider.Manifest
		want          []provider.ResourceKey
	}{
		{
			name: "no resource removed",
			manifests: []provider.Manifest{
				{
					Key: provider.ResourceKey{
						APIVersion: "v1",
						Kind:       "Service",
						Name:       "foo",
					},
				},
			},
			liveResources: []provider.Manifest{
				{
					Key: provider.ResourceKey{
						APIVersion: "v1",
						Kind:       "Service",
						Name:       "foo",
					},
				},
			},
			want: []provider.ResourceKey{},
		},
		{
			name:      "one resource removed",
			manifests: []provider.Manifest{},
			liveResources: []provider.Manifest{
				{
					Key: provider.ResourceKey{
						APIVersion: "v1",
						Kind:       "Service",
						Name:       "foo",
					},
				},
			},
			want: []provider.ResourceKey{
				{
					APIVersion: "v1",
					Kind:       "Service",
					Name:       "foo",
				},
			},
		},
		{
			name: "don't remove resource running in different namespace from manifests",
			manifests: []provider.Manifest{
				{
					Key: provider.ResourceKey{
						APIVersion: "v1",
						Kind:       "Service",
						Namespace:  "different",
						Name:       "foo",
					},
				},
			},
			liveResources: []provider.Manifest{
				{
					Key: provider.ResourceKey{
						APIVersion: "v1",
						Kind:       "Service",
						Namespace:  "namespace",
						Name:       "foo",
					},
				},
			},
			want: []provider.ResourceKey{},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := findRemoveResources(tc.manifests, tc.liveResources)
			assert.Equal(t, tc.want, got)
		})
	}
}
