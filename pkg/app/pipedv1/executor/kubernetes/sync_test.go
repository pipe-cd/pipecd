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
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/executor"
	provider "github.com/pipe-cd/pipecd/pkg/app/pipedv1/platformprovider/kubernetes"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/platformprovider/kubernetes/kubernetestest"
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
