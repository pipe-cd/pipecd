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
	"context"
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"sigs.k8s.io/controller-runtime/pkg/envtest"

	"github.com/pipe-cd/pipecd/pkg/app/piped/executor"
	provider "github.com/pipe-cd/pipecd/pkg/app/piped/platformprovider/kubernetes"
	"github.com/pipe-cd/pipecd/pkg/app/piped/platformprovider/kubernetes/kubernetestest"
	"github.com/pipe-cd/pipecd/pkg/app/piped/toolregistry"
	"github.com/pipe-cd/pipecd/pkg/cache"
	"github.com/pipe-cd/pipecd/pkg/cache/cachetest"
	"github.com/pipe-cd/pipecd/pkg/config"
	"github.com/pipe-cd/pipecd/pkg/model"
)

func TestEnsureSync(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testcases := []struct {
		name     string
		executor *deployExecutor
		want     model.StageStatus
	}{
		{
			name: "failed to load manifest",
			want: model.StageStatus_STAGE_FAILURE,
			executor: &deployExecutor{
				Input: executor.Input{
					Deployment: &model.Deployment{
						Trigger: &model.DeploymentTrigger{
							Commit: &model.Commit{},
						},
					},
					LogPersister: &fakeLogPersister{},
					AppManifestsCache: func() cache.Cache {
						c := cachetest.NewMockCache(ctrl)
						c.EXPECT().Get(gomock.Any()).Return(nil, fmt.Errorf("not found"))
						return c
					}(),
					Logger: zap.NewNop(),
				},
				loader: func() provider.Loader {
					p := kubernetestest.NewMockLoader(ctrl)
					p.EXPECT().LoadManifests(gomock.Any()).Return(nil, fmt.Errorf("error"))
					return p
				}(),
			},
		},
		{
			name: "unable to apply manifests",
			want: model.StageStatus_STAGE_FAILURE,
			executor: &deployExecutor{
				Input: executor.Input{
					Deployment: &model.Deployment{
						Trigger: &model.DeploymentTrigger{
							Commit: &model.Commit{},
						},
					},
					PipedConfig:  &config.PipedSpec{},
					LogPersister: &fakeLogPersister{},
					AppManifestsCache: func() cache.Cache {
						c := cachetest.NewMockCache(ctrl)
						c.EXPECT().Get(gomock.Any()).Return(nil, fmt.Errorf("not found"))
						c.EXPECT().Put(gomock.Any(), gomock.Any()).Return(nil)
						return c
					}(),
					Logger: zap.NewNop(),
				},
				loader: func() provider.Loader {
					p := kubernetestest.NewMockLoader(ctrl)
					p.EXPECT().LoadManifests(gomock.Any()).Return([]provider.Manifest{
						provider.MakeManifest(provider.ResourceKey{
							APIVersion: "apps/v1",
							Kind:       provider.KindDeployment,
						}, &unstructured.Unstructured{
							Object: map[string]interface{}{"spec": map[string]interface{}{}},
						}),
					}, nil)
					return p
				}(),
				applierGetter: &applierGroup{
					defaultApplier: func() provider.Applier {
						p := kubernetestest.NewMockApplier(ctrl)
						p.EXPECT().ApplyManifest(gomock.Any(), gomock.Any()).Return(fmt.Errorf("error"))
						return p
					}(),
				},
				appCfg: &config.KubernetesApplicationSpec{
					QuickSync: config.K8sSyncStageOptions{
						AddVariantLabelToSelector: true,
					},
				},
			},
		},
		{
			name: "successfully apply manifests",
			want: model.StageStatus_STAGE_SUCCESS,
			executor: &deployExecutor{
				Input: executor.Input{
					Deployment: &model.Deployment{
						Trigger: &model.DeploymentTrigger{
							Commit: &model.Commit{},
						},
					},
					PipedConfig:  &config.PipedSpec{},
					LogPersister: &fakeLogPersister{},
					AppManifestsCache: func() cache.Cache {
						c := cachetest.NewMockCache(ctrl)
						c.EXPECT().Get(gomock.Any()).Return(nil, fmt.Errorf("not found"))
						c.EXPECT().Put(gomock.Any(), gomock.Any()).Return(nil)
						return c
					}(),
					Logger: zap.NewNop(),
				},
				loader: func() provider.Loader {
					p := kubernetestest.NewMockLoader(ctrl)
					p.EXPECT().LoadManifests(gomock.Any()).Return([]provider.Manifest{
						provider.MakeManifest(provider.ResourceKey{
							APIVersion: "apps/v1",
							Kind:       provider.KindDeployment,
						}, &unstructured.Unstructured{
							Object: map[string]interface{}{"spec": map[string]interface{}{}},
						}),
					}, nil)
					return p
				}(),
				applierGetter: &applierGroup{
					defaultApplier: func() provider.Applier {
						p := kubernetestest.NewMockApplier(ctrl)
						p.EXPECT().ApplyManifest(gomock.Any(), gomock.Any()).Return(nil)
						return p
					}(),
				},
				appCfg: &config.KubernetesApplicationSpec{
					QuickSync: config.K8sSyncStageOptions{
						AddVariantLabelToSelector: true,
					},
				},
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			got := tc.executor.ensureSync(ctx)
			assert.Equal(t, tc.want, got)
			cancel()
		})
	}
}

func TestExecutor_ensureSync(t *testing.T) {
	ctrl := gomock.NewController(t)

	// initialize tool registry
	toolregistry.InitDefaultRegistry("/tmp/piped-bin", zap.NewNop())

	// initialize envtest
	tEnv := new(envtest.Environment)
	cfg, err := tEnv.Start()
	require.NoError(t, err)
	defer tEnv.Stop()

	kubeconfig, err := kubeconfigFromRestConfig(cfg)
	require.NoError(t, err)

	workDir := t.TempDir()
	kubeconfigPath := path.Join(workDir, "kubeconfig")
	err = os.WriteFile(kubeconfigPath, []byte(kubeconfig), 0755)
	require.NoError(t, err)

	manifests, err := provider.LoadManifestsFromYAMLFile("../../../../../examples/kubernetes/simple/deployment.yaml")
	require.NoError(t, err)

	appCfg, err := config.LoadFromYAML("../../../../../examples/kubernetes/simple/app.pipecd.yaml")
	require.NoError(t, err)

	executor := &deployExecutor{
		Input: executor.Input{
			Deployment: &model.Deployment{
				ApplicationId:    "app-id",
				PlatformProvider: "default",
			},
			LogPersister: &fakeLogPersister{},
			PipedConfig: &config.PipedSpec{
				PipedID: "piped-id",
				PlatformProviders: []config.PipedPlatformProvider{
					{
						Name: "default",
						Type: model.PlatformProviderKubernetes,
						KubernetesConfig: &config.PlatformProviderKubernetesConfig{
							KubeConfigPath: kubeconfigPath,
						},
					},
				},
			},
			AppManifestsCache: func() cache.Cache {
				c := cachetest.NewMockCache(ctrl)
				c.EXPECT().Get(gomock.Any()).Return(nil, fmt.Errorf("not found"))
				c.EXPECT().Put(gomock.Any(), gomock.Any()).Return(nil)
				return c
			}(),
			Logger: zap.NewNop(),
		},
		loader: func() provider.Loader {
			p := kubernetestest.NewMockLoader(ctrl)
			p.EXPECT().LoadManifests(gomock.Any()).Return(func() []provider.Manifest {
				return manifests
			}(), nil)
			return p
		}(),
		appCfg: appCfg.KubernetesApplicationSpec,
		applierGetter: func() applierGetter {
			ag, err := newApplierGroup("default", *appCfg.KubernetesApplicationSpec, &config.PipedSpec{
				PlatformProviders: []config.PipedPlatformProvider{
					{
						Name: "default",
						Type: model.PlatformProviderKubernetes,
						KubernetesConfig: &config.PlatformProviderKubernetesConfig{
							KubeConfigPath: kubeconfigPath,
						},
					},
				},
			}, zap.NewNop())
			require.NoError(t, err)
			return ag
		}(),
		commit: "0123456789",
	}

	status := executor.ensureSync(context.Background())
	assert.Equal(t, model.StageStatus_STAGE_SUCCESS.String(), status.String())

	// check the deployment is created with client-go
	dynamicClient, err := dynamic.NewForConfig(cfg)
	require.NoError(t, err)

	deployment, err := dynamicClient.Resource(schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}).Namespace("default").Get(context.Background(), "simple", metav1.GetOptions{})
	require.NoError(t, err)

	assert.Equal(t, "simple", deployment.GetName())
	assert.Equal(t, "simple", deployment.GetLabels()["app"])
	assert.Equal(t, "piped", deployment.GetAnnotations()["pipecd.dev/managed-by"])
	assert.Equal(t, "piped-id", deployment.GetAnnotations()["pipecd.dev/piped"])
	assert.Equal(t, "app-id", deployment.GetAnnotations()["pipecd.dev/application"])
	assert.Equal(t, "apps/v1", deployment.GetAnnotations()["pipecd.dev/original-api-version"])
	assert.Equal(t, "apps/v1:Deployment:default:simple", deployment.GetAnnotations()["pipecd.dev/resource-key"])
	assert.Equal(t, "0123456789", deployment.GetAnnotations()["pipecd.dev/commit-hash"])
}

func kubeconfigFromRestConfig(restConfig *rest.Config) (string, error) {
	clusters := make(map[string]*clientcmdapi.Cluster)
	clusters["default-cluster"] = &clientcmdapi.Cluster{
		Server:                   restConfig.Host,
		CertificateAuthorityData: restConfig.CAData,
	}
	contexts := make(map[string]*clientcmdapi.Context)
	contexts["default-context"] = &clientcmdapi.Context{
		Cluster:  "default-cluster",
		AuthInfo: "default-user",
	}
	authinfos := make(map[string]*clientcmdapi.AuthInfo)
	authinfos["default-user"] = &clientcmdapi.AuthInfo{
		ClientCertificateData: restConfig.CertData,
		ClientKeyData:         restConfig.KeyData,
	}
	clientConfig := clientcmdapi.Config{
		Kind:           "Config",
		APIVersion:     "v1",
		Clusters:       clusters,
		Contexts:       contexts,
		CurrentContext: "default-context",
		AuthInfos:      authinfos,
	}
	b, err := clientcmd.Write(clientConfig)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func TestFindRemoveResources(t *testing.T) {
	t.Parallel()

	testcases := []struct {
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
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := findRemoveResources(tc.manifests, tc.liveResources)
			assert.Equal(t, tc.want, got)
		})
	}
}
