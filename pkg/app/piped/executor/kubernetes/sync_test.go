package kubernetes

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	provider "github.com/pipe-cd/pipe/pkg/app/piped/cloudprovider/kubernetes"
	"github.com/pipe-cd/pipe/pkg/app/piped/cloudprovider/kubernetes/providertest"
	"github.com/pipe-cd/pipe/pkg/app/piped/executor"
	"github.com/pipe-cd/pipe/pkg/cache"
	"github.com/pipe-cd/pipe/pkg/cache/cachetest"
	"github.com/pipe-cd/pipe/pkg/config"
	"github.com/pipe-cd/pipe/pkg/model"
)

func TestEnsureSync(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testcases := []struct {
		name     string
		executor *Executor
		want     model.StageStatus
	}{
		{
			name: "failed to load manifest",
			want: model.StageStatus_STAGE_FAILURE,
			executor: &Executor{
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
				provider: func() provider.Provider {
					p := providertest.NewMockProvider(ctrl)
					p.EXPECT().LoadManifests(gomock.Any()).Return(nil, fmt.Errorf("error"))
					return p
				}(),
			},
		},
		{
			name: "missing variant selector",
			want: model.StageStatus_STAGE_FAILURE,
			executor: &Executor{
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
						c.EXPECT().Put(gomock.Any(), gomock.Any()).Return(nil)
						return c
					}(),
					Logger: zap.NewNop(),
				},
				provider: func() provider.Provider {
					p := providertest.NewMockProvider(ctrl)
					p.EXPECT().LoadManifests(gomock.Any()).Return([]provider.Manifest{
						provider.MakeManifest(provider.ResourceKey{
							APIVersion: "apps/v1",
							Kind:       provider.KindDeployment,
						}, &unstructured.Unstructured{}),
					}, nil)
					return p
				}(),
				config: &config.KubernetesDeploymentSpec{
					GenericDeploymentSpec: config.GenericDeploymentSpec{
						Pipeline: &config.DeploymentPipeline{},
					},
				},
			},
		},
		{
			name: "unable to apply manifests",
			want: model.StageStatus_STAGE_FAILURE,
			executor: &Executor{
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
				provider: func() provider.Provider {
					p := providertest.NewMockProvider(ctrl)
					p.EXPECT().LoadManifests(gomock.Any()).Return([]provider.Manifest{
						provider.MakeManifest(provider.ResourceKey{
							APIVersion: "apps/v1",
							Kind:       provider.KindDeployment,
						}, &unstructured.Unstructured{
							Object: map[string]interface{}{"spec": map[string]interface{}{}},
						}),
					}, nil)
					p.EXPECT().ApplyManifest(gomock.Any(), gomock.Any()).Return(fmt.Errorf("error"))
					return p
				}(),
				config: &config.KubernetesDeploymentSpec{
					QuickSync: config.K8sSyncStageOptions{
						AddVariantLabelToSelector: true,
					},
				},
			},
		},
		{
			name: "successfully apply manifests",
			want: model.StageStatus_STAGE_SUCCESS,
			executor: &Executor{
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
				provider: func() provider.Provider {
					p := providertest.NewMockProvider(ctrl)
					p.EXPECT().LoadManifests(gomock.Any()).Return([]provider.Manifest{
						provider.MakeManifest(provider.ResourceKey{
							APIVersion: "apps/v1",
							Kind:       provider.KindDeployment,
						}, &unstructured.Unstructured{
							Object: map[string]interface{}{"spec": map[string]interface{}{}},
						}),
					}, nil)
					p.EXPECT().ApplyManifest(gomock.Any(), gomock.Any()).Return(nil)
					return p
				}(),
				config: &config.KubernetesDeploymentSpec{
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

func TestFindRemoveResources(t *testing.T) {
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
