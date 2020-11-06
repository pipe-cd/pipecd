package kubernetes

import (
	"testing"

	"github.com/stretchr/testify/assert"

	provider "github.com/pipe-cd/pipe/pkg/app/piped/cloudprovider/kubernetes"
)

func TestFindRemoveManifests(t *testing.T) {
	tests := []struct {
		name      string
		prevs     []provider.Manifest
		curs      []provider.Manifest
		namespace string
		want      []provider.ResourceKey
	}{
		{
			name: "no resource removed",
			prevs: []provider.Manifest{
				{
					Key: provider.ResourceKey{
						APIVersion: "v1",
						Kind:       "Service",
						Name:       "foo",
					},
				},
			},
			curs: []provider.Manifest{
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
			name: "one resource removed",
			prevs: []provider.Manifest{
				{
					Key: provider.ResourceKey{
						APIVersion: "v1",
						Kind:       "Service",
						Name:       "foo",
					},
				},
			},
			curs: []provider.Manifest{},
			want: []provider.ResourceKey{
				{
					APIVersion: "v1",
					Kind:       "Service",
					Name:       "foo",
				},
			},
		},
		{
			name: "one resource removed with specified namespace",
			prevs: []provider.Manifest{
				{
					Key: provider.ResourceKey{
						APIVersion: "v1",
						Kind:       "Service",
						Name:       "foo",
					},
				},
			},
			curs:      []provider.Manifest{},
			namespace: "namespace",
			want: []provider.ResourceKey{
				{
					APIVersion: "v1",
					Kind:       "Service",
					Namespace:  "namespace",
					Name:       "foo",
				},
			},
		},
		{
			name: "give namespace different from running one",
			prevs: []provider.Manifest{
				{
					Key: provider.ResourceKey{
						APIVersion: "v1",
						Kind:       "Service",
						Namespace:  "namespace",
						Name:       "foo",
					},
				},
			},
			curs:      []provider.Manifest{},
			namespace: "different",
			want: []provider.ResourceKey{
				{
					APIVersion: "v1",
					Kind:       "Service",
					Namespace:  "namespace",
					Name:       "foo",
				},
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := findRemoveManifests(tc.prevs, tc.curs, tc.namespace)
			assert.Equal(t, tc.want, got)
		})
	}
}
