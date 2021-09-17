package kubernetes

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	provider "github.com/pipe-cd/pipe/pkg/app/piped/cloudprovider/kubernetes"
	"github.com/pipe-cd/pipe/pkg/config"
)

func TestDecideStrategy(t *testing.T) {
	tests := []struct {
		name            string
		olds            []provider.Manifest
		news            []provider.Manifest
		workloadRefs    []config.K8sResourceReference
		wantProgressive bool
		wantDesc        string
	}{
		{
			name: "no workload in the old commit",
			news: []provider.Manifest{
				{
					Key: provider.ResourceKey{
						APIVersion: "apps/v1",
						Kind:       provider.KindDeployment,
						Name:       "name",
					},
				},
			},
			wantProgressive: false,
			wantDesc:        "Quick sync by applying all manifests because it was unable to find the currently running workloads",
		},
		{
			name: "no workload in the new commit",
			olds: []provider.Manifest{
				{
					Key: provider.ResourceKey{
						APIVersion: "apps/v1",
						Kind:       provider.KindDeployment,
						Name:       "name",
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
					Name:       "name",
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
			name: "mutilple workloads: pod template was changed",
			olds: func() []provider.Manifest {
				m1 := provider.MakeManifest(provider.ResourceKey{
					APIVersion: "apps/v1",
					Kind:       provider.KindDeployment,
					Name:       "name-1",
				}, &unstructured.Unstructured{
					Object: map[string]interface{}{"spec": map[string]interface{}{"template": "foo-1"}}},
				)
				m2 := provider.MakeManifest(provider.ResourceKey{
					APIVersion: "apps/v1",
					Kind:       provider.KindDeployment,
					Name:       "name-2",
				}, &unstructured.Unstructured{
					Object: map[string]interface{}{"spec": map[string]interface{}{"template": "foo-2"}}},
				)
				return []provider.Manifest{m1, m2}
			}(),
			news: func() []provider.Manifest {
				m1 := provider.MakeManifest(provider.ResourceKey{
					APIVersion: "apps/v1",
					Kind:       provider.KindDeployment,
					Name:       "name-1",
				}, &unstructured.Unstructured{
					Object: map[string]interface{}{"spec": map[string]interface{}{"template": "foo-1"}}},
				)
				m2 := provider.MakeManifest(provider.ResourceKey{
					APIVersion: "apps/v1",
					Kind:       provider.KindDeployment,
					Name:       "name-2",
				}, &unstructured.Unstructured{
					Object: map[string]interface{}{"spec": map[string]interface{}{"template": "bar-2"}}},
				)
				return []provider.Manifest{m1, m2}
			}(),
			wantProgressive: true,
			wantDesc:        "Sync progressively because pod template of workload name-2 was changed",
		},
		{
			name: "changed deployment was not the target",
			olds: func() []provider.Manifest {
				m1 := provider.MakeManifest(provider.ResourceKey{
					APIVersion: "apps/v1",
					Kind:       provider.KindDeployment,
					Name:       "name-1",
				}, &unstructured.Unstructured{
					Object: map[string]interface{}{"spec": map[string]interface{}{"template": "foo-1"}}},
				)
				m2 := provider.MakeManifest(provider.ResourceKey{
					APIVersion: "apps/v1",
					Kind:       provider.KindDeployment,
					Name:       "name-2",
				}, &unstructured.Unstructured{
					Object: map[string]interface{}{"spec": map[string]interface{}{"template": "foo-2"}}},
				)
				return []provider.Manifest{m1, m2}
			}(),
			news: func() []provider.Manifest {
				m1 := provider.MakeManifest(provider.ResourceKey{
					APIVersion: "apps/v1",
					Kind:       provider.KindDeployment,
					Name:       "name-1",
				}, &unstructured.Unstructured{
					Object: map[string]interface{}{"spec": map[string]interface{}{"template": "foo-1"}}},
				)
				m2 := provider.MakeManifest(provider.ResourceKey{
					APIVersion: "apps/v1",
					Kind:       provider.KindDeployment,
					Name:       "name-2",
				}, &unstructured.Unstructured{
					Object: map[string]interface{}{"spec": map[string]interface{}{"template": "bar-2"}}},
				)
				return []provider.Manifest{m1, m2}
			}(),
			workloadRefs: []config.K8sResourceReference{
				{
					Kind: provider.KindDeployment,
					Name: "name-1",
				},
			},
			wantProgressive: false,
			wantDesc:        "Quick sync by applying all manifests",
		},
		{
			name: "scale one deployment",
			olds: func() []provider.Manifest {
				m := provider.MakeManifest(provider.ResourceKey{
					APIVersion: "apps/v1",
					Kind:       provider.KindDeployment,
					Name:       "name",
				}, &unstructured.Unstructured{
					Object: map[string]interface{}{"spec": map[string]interface{}{
						"template": "foo",
						"replicas": 1,
					}}},
				)
				return []provider.Manifest{m}
			}(),
			news: func() []provider.Manifest {
				m := provider.MakeManifest(provider.ResourceKey{
					APIVersion: "apps/v1",
					Kind:       provider.KindDeployment,
					Name:       "name",
				}, &unstructured.Unstructured{
					Object: map[string]interface{}{"spec": map[string]interface{}{
						"template": "foo",
						"replicas": 2,
					}}},
				)
				return []provider.Manifest{m}
			}(),
			wantProgressive: false,
			wantDesc:        "Quick sync to scale Deployment/name from 1 to 2",
		},
		{
			name: "scale multiple deployments",
			olds: func() []provider.Manifest {
				m1 := provider.MakeManifest(provider.ResourceKey{
					APIVersion: "apps/v1",
					Kind:       provider.KindDeployment,
					Name:       "name-1",
				}, &unstructured.Unstructured{
					Object: map[string]interface{}{"spec": map[string]interface{}{
						"template": "foo",
						"replicas": 1,
					}}},
				)
				m2 := provider.MakeManifest(provider.ResourceKey{
					APIVersion: "apps/v1",
					Kind:       provider.KindDeployment,
					Name:       "name-2",
				}, &unstructured.Unstructured{
					Object: map[string]interface{}{"spec": map[string]interface{}{
						"template": "bar",
						"replicas": 20,
					}}},
				)
				return []provider.Manifest{m1, m2}
			}(),
			news: func() []provider.Manifest {
				m1 := provider.MakeManifest(provider.ResourceKey{
					APIVersion: "apps/v1",
					Kind:       provider.KindDeployment,
					Name:       "name-1",
				}, &unstructured.Unstructured{
					Object: map[string]interface{}{"spec": map[string]interface{}{
						"template": "foo",
						"replicas": 5,
					}}},
				)
				m2 := provider.MakeManifest(provider.ResourceKey{
					APIVersion: "apps/v1",
					Kind:       provider.KindDeployment,
					Name:       "name-2",
				}, &unstructured.Unstructured{
					Object: map[string]interface{}{"spec": map[string]interface{}{
						"template": "bar",
						"replicas": 10,
					}}},
				)
				return []provider.Manifest{m1, m2}
			}(),
			wantProgressive: false,
			wantDesc:        "Quick sync to scale Deployment/name-1 from 1 to 5, Deployment/name-2 from 20 to 10",
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
			gotProgressive, gotDesc := decideStrategy(tc.olds, tc.news, tc.workloadRefs)
			assert.Equal(t, tc.wantProgressive, gotProgressive)
			assert.Equal(t, tc.wantDesc, gotDesc)
		})
	}
}

func TestDetermineVersion(t *testing.T) {
	testcases := []struct {
		name          string
		manifests     string
		expected      string
		expectedError error
	}{
		// {
		// 	name:      "no workload",
		// 	manifests: "testdata/version_no_workload.yaml",
		// 	expected:  "unknown",
		// },
		// {
		// 	name:      "single container",
		// 	manifests: "testdata/version_single_container.yaml",
		// 	expected:  "v1.0.0",
		// },
		// {
		// 	name:      "multiple containers",
		// 	manifests: "testdata/version_multi_containers.yaml",
		// 	expected:  "helloworld: v1.0.0, my-service: v0.6.0",
		// },
		{
			name:      "multiple workloads",
			manifests: "testdata/version_multi_workloads.yaml",
			expected:  "helloworld: v1.0.0, my-service: v0.5.0",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			manifests, err := provider.LoadManifestsFromYAMLFile(tc.manifests)
			require.NoError(t, err)

			version, err := determineVersion(manifests)
			assert.Equal(t, tc.expected, version)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
