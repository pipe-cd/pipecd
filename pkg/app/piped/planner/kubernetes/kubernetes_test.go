package kubernetes

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	provider "github.com/pipe-cd/pipe/pkg/app/piped/cloudprovider/kubernetes"
)

func TestDecideStrategy(t *testing.T) {
	tests := []struct {
		name            string
		olds            []provider.Manifest
		news            []provider.Manifest
		wantProgressive bool
		wantDesc        string
	}{
		{
			name: "no running workloads found",
			olds: []provider.Manifest{
				{
					Key: provider.ResourceKey{
						APIVersion: "v1",
						Kind:       provider.KindService,
					},
				},
			},
			wantProgressive: false,
			wantDesc:        "Quick sync by applying all manifests because it was unable to find the currently running workloads",
		},
		{
			name: "no workloads found in the new manifests",
			olds: []provider.Manifest{
				{
					Key: provider.ResourceKey{
						APIVersion: "apps/v1",
						Kind:       provider.KindDeployment,
					},
				},
			},
			news: []provider.Manifest{
				{
					Key: provider.ResourceKey{
						APIVersion: "v1",
						Kind:       provider.KindService,
					},
				},
			},
			wantProgressive: false,
			wantDesc:        "Quick sync by applying all manifests because it was unable to find workloads in the new manifests",
		},
		{
			name: "pod template was changed",
			olds: func() []provider.Manifest {
				m := provider.MakeManifest(provider.ResourceKey{
					APIVersion: "apps/v1",
					Kind:       provider.KindDeployment,
				}, &unstructured.Unstructured{
					Object: map[string]interface{}{"spec": map[string]interface{}{"template": "foo"}}},
				)
				return []provider.Manifest{m}
			}(),
			news: func() []provider.Manifest {
				m := provider.MakeManifest(provider.ResourceKey{
					APIVersion: "apps/v1",
					Kind:       provider.KindDeployment,
					Name:       "name",
				}, &unstructured.Unstructured{
					Object: map[string]interface{}{"spec": map[string]interface{}{"template": "bar"}}},
				)
				return []provider.Manifest{m}
			}(),
			wantProgressive: true,
			wantDesc:        "Sync progressively because pod template of workload name was changed",
		},
		{
			name: "configmap deleted",
			olds: func() []provider.Manifest {
				m1 := provider.MakeManifest(provider.ResourceKey{
					APIVersion: "apps/v1",
					Kind:       provider.KindDeployment,
				}, &unstructured.Unstructured{})
				m2 := provider.MakeManifest(provider.ResourceKey{
					APIVersion: "v1",
					Kind:       provider.KindConfigMap,
				}, &unstructured.Unstructured{})
				return []provider.Manifest{m1, m2}
			}(),
			news: func() []provider.Manifest {
				m := provider.MakeManifest(provider.ResourceKey{
					APIVersion: "apps/v1",
					Kind:       provider.KindDeployment,
				}, &unstructured.Unstructured{})
				return []provider.Manifest{m}
			}(),
			wantProgressive: true,
			wantDesc:        "Sync progressively because 1 configmap/secret deleted",
		},
		{
			name: "new configmap added",
			olds: func() []provider.Manifest {
				m := provider.MakeManifest(provider.ResourceKey{
					APIVersion: "apps/v1",
					Kind:       provider.KindDeployment,
				}, &unstructured.Unstructured{})
				return []provider.Manifest{m}
			}(),
			news: func() []provider.Manifest {
				m1 := provider.MakeManifest(provider.ResourceKey{
					APIVersion: "apps/v1",
					Kind:       provider.KindDeployment,
				}, &unstructured.Unstructured{})
				m2 := provider.MakeManifest(provider.ResourceKey{
					APIVersion: "v1",
					Kind:       provider.KindConfigMap,
				}, &unstructured.Unstructured{})
				return []provider.Manifest{m1, m2}
			}(),
			wantProgressive: true,
			wantDesc:        "Sync progressively because new 1 configmap/secret added",
		},
		{
			name: "one configmap updated",
			olds: func() []provider.Manifest {
				m1 := provider.MakeManifest(provider.ResourceKey{
					APIVersion: "apps/v1",
					Kind:       provider.KindDeployment,
				}, &unstructured.Unstructured{})
				m2 := provider.MakeManifest(provider.ResourceKey{
					APIVersion: "v1",
					Kind:       provider.KindConfigMap,
					Name:       "configmap1",
				}, &unstructured.Unstructured{
					Object: map[string]interface{}{"data": "foo"}},
				)
				m3 := provider.MakeManifest(provider.ResourceKey{
					APIVersion: "v1",
					Kind:       provider.KindConfigMap,
					Name:       "configmap2",
				}, &unstructured.Unstructured{
					Object: map[string]interface{}{"data": "baz"}},
				)
				return []provider.Manifest{m1, m2, m3}
			}(),
			news: func() []provider.Manifest {
				m1 := provider.MakeManifest(provider.ResourceKey{
					APIVersion: "apps/v1",
					Kind:       provider.KindDeployment,
				}, &unstructured.Unstructured{})
				m2 := provider.MakeManifest(provider.ResourceKey{
					APIVersion: "v1",
					Kind:       provider.KindConfigMap,
					Name:       "configmap1",
				}, &unstructured.Unstructured{
					Object: map[string]interface{}{"data": "bar"}},
				)
				m3 := provider.MakeManifest(provider.ResourceKey{
					APIVersion: "v1",
					Kind:       provider.KindConfigMap,
					Name:       "configmap2",
				}, &unstructured.Unstructured{
					Object: map[string]interface{}{"data": "baz"}},
				)
				return []provider.Manifest{m1, m2, m3}
			}(),
			wantProgressive: true,
			wantDesc:        "Sync progressively because ConfigMap configmap1 was updated",
		},
		{
			name: "all configmaps as is",
			olds: func() []provider.Manifest {
				m1 := provider.MakeManifest(provider.ResourceKey{
					APIVersion: "apps/v1",
					Kind:       provider.KindDeployment,
				}, &unstructured.Unstructured{})
				m2 := provider.MakeManifest(provider.ResourceKey{
					APIVersion: "v1",
					Kind:       provider.KindConfigMap,
					Name:       "configmap1",
				}, &unstructured.Unstructured{
					Object: map[string]interface{}{"data": "foo"}},
				)
				m3 := provider.MakeManifest(provider.ResourceKey{
					APIVersion: "v1",
					Kind:       provider.KindConfigMap,
					Name:       "configmap2",
				}, &unstructured.Unstructured{
					Object: map[string]interface{}{"data": "baz"}},
				)
				return []provider.Manifest{m1, m2, m3}
			}(),
			news: func() []provider.Manifest {
				m1 := provider.MakeManifest(provider.ResourceKey{
					APIVersion: "apps/v1",
					Kind:       provider.KindDeployment,
				}, &unstructured.Unstructured{})
				m2 := provider.MakeManifest(provider.ResourceKey{
					APIVersion: "v1",
					Kind:       provider.KindConfigMap,
					Name:       "configmap1",
				}, &unstructured.Unstructured{
					Object: map[string]interface{}{"data": "foo"}},
				)
				m3 := provider.MakeManifest(provider.ResourceKey{
					APIVersion: "v1",
					Kind:       provider.KindConfigMap,
					Name:       "configmap2",
				}, &unstructured.Unstructured{
					Object: map[string]interface{}{"data": "baz"}},
				)
				return []provider.Manifest{m1, m2, m3}
			}(),
			wantProgressive: false,
			wantDesc:        "Quick sync by applying all manifests",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotProgressive, gotDesc := decideStrategy(tc.olds, tc.news)
			assert.Equal(t, tc.wantProgressive, gotProgressive)
			assert.Equal(t, tc.wantDesc, gotDesc)
		})
	}
}
