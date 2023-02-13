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
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/pipe-cd/pipecd/pkg/app/piped/executor"
	provider "github.com/pipe-cd/pipecd/pkg/app/piped/platformprovider/kubernetes"
	"github.com/pipe-cd/pipecd/pkg/app/piped/platformprovider/kubernetes/kubernetestest"
	"github.com/pipe-cd/pipecd/pkg/cache"
	"github.com/pipe-cd/pipecd/pkg/cache/cachetest"
	"github.com/pipe-cd/pipecd/pkg/config"
	"github.com/pipe-cd/pipecd/pkg/model"
)

func TestEnsureCanaryRollout(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	appCfg := &config.KubernetesApplicationSpec{
		VariantLabel: config.KubernetesVariantLabel{
			Key:           "pipecd.dev/variant",
			PrimaryValue:  "primary",
			BaselineValue: "baseline",
			CanaryValue:   "canary",
		},
	}
	testcases := []struct {
		name     string
		executor *deployExecutor
		want     model.StageStatus
	}{
		{
			name: "malformed configuration",
			want: model.StageStatus_STAGE_FAILURE,
			executor: &deployExecutor{
				Input: executor.Input{
					Deployment: &model.Deployment{
						Trigger: &model.DeploymentTrigger{
							Commit: &model.Commit{},
						},
					},
					Stage:        &model.PipelineStage{},
					LogPersister: &fakeLogPersister{},
					Logger:       zap.NewNop(),
				},
				appCfg: appCfg,
			},
		},
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
					Stage:        &model.PipelineStage{},
					StageConfig: config.PipelineStage{
						K8sCanaryRolloutStageOptions: &config.K8sCanaryRolloutStageOptions{},
					},
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
				appCfg: appCfg,
			},
		},
		{
			name: "no manifests to handle",
			want: model.StageStatus_STAGE_FAILURE,
			executor: &deployExecutor{
				Input: executor.Input{
					Deployment: &model.Deployment{
						Trigger: &model.DeploymentTrigger{
							Commit: &model.Commit{},
						},
					},
					LogPersister: &fakeLogPersister{},
					Stage:        &model.PipelineStage{},
					StageConfig: config.PipelineStage{
						K8sCanaryRolloutStageOptions: &config.K8sCanaryRolloutStageOptions{},
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
					p.EXPECT().LoadManifests(gomock.Any()).Return([]provider.Manifest{}, nil)
					return p
				}(),
				appCfg: appCfg,
			},
		},
		{
			name: "failed to apply manifests",
			want: model.StageStatus_STAGE_FAILURE,
			executor: &deployExecutor{
				Input: executor.Input{
					Deployment: &model.Deployment{
						Trigger: &model.DeploymentTrigger{
							Commit: &model.Commit{},
						},
					},
					LogPersister: &fakeLogPersister{},
					Stage:        &model.PipelineStage{},
					StageConfig: config.PipelineStage{
						K8sCanaryRolloutStageOptions: &config.K8sCanaryRolloutStageOptions{},
					},
					AppManifestsCache: func() cache.Cache {
						c := cachetest.NewMockCache(ctrl)
						c.EXPECT().Get(gomock.Any()).Return(nil, fmt.Errorf("not found"))
						c.EXPECT().Put(gomock.Any(), gomock.Any()).Return(nil)
						return c
					}(),
					MetadataStore: &fakeMetadataStore{},
					PipedConfig:   &config.PipedSpec{},
					Logger:        zap.NewNop(),
				},
				loader: func() provider.Loader {
					p := kubernetestest.NewMockLoader(ctrl)
					p.EXPECT().LoadManifests(gomock.Any()).Return([]provider.Manifest{
						provider.MakeManifest(provider.ResourceKey{
							APIVersion: "apps/v1",
							Kind:       provider.KindDeployment,
						}, &unstructured.Unstructured{
							Object: map[string]interface{}{
								"apiVersion": "apps/v1",
								"kind":       "Deployment",
								"spec": map[string]interface{}{
									"selector": map[string]interface{}{
										"matchLabels": map[string]interface{}{"app": "foo"},
									},
								},
							},
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
				appCfg: &config.KubernetesApplicationSpec{},
			},
		},
		{
			name: "successfully applying manifests",
			want: model.StageStatus_STAGE_SUCCESS,
			executor: &deployExecutor{
				Input: executor.Input{
					Deployment: &model.Deployment{
						Trigger: &model.DeploymentTrigger{
							Commit: &model.Commit{},
						},
					},
					LogPersister: &fakeLogPersister{},
					Stage:        &model.PipelineStage{},
					StageConfig: config.PipelineStage{
						K8sCanaryRolloutStageOptions: &config.K8sCanaryRolloutStageOptions{},
					},
					AppManifestsCache: func() cache.Cache {
						c := cachetest.NewMockCache(ctrl)
						c.EXPECT().Get(gomock.Any()).Return(nil, fmt.Errorf("not found"))
						c.EXPECT().Put(gomock.Any(), gomock.Any()).Return(nil)
						return c
					}(),
					MetadataStore: &fakeMetadataStore{},
					PipedConfig:   &config.PipedSpec{},
					Logger:        zap.NewNop(),
				},
				loader: func() provider.Loader {
					p := kubernetestest.NewMockLoader(ctrl)
					p.EXPECT().LoadManifests(gomock.Any()).Return([]provider.Manifest{
						provider.MakeManifest(provider.ResourceKey{
							APIVersion: "apps/v1",
							Kind:       provider.KindDeployment,
						}, &unstructured.Unstructured{
							Object: map[string]interface{}{
								"apiVersion": "apps/v1",
								"kind":       "Deployment",
								"spec": map[string]interface{}{
									"selector": map[string]interface{}{
										"matchLabels": map[string]interface{}{"app": "foo"},
									},
								},
							},
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
				appCfg: &config.KubernetesApplicationSpec{},
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			got := tc.executor.ensureCanaryRollout(ctx)
			assert.Equal(t, tc.want, got)
		})
	}
}
