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

package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pipe-cd/pipecd/pkg/model"
)

func TestKubernetesApplicationConfig(t *testing.T) {
	testcases := []struct {
		fileName           string
		expectedKind       Kind
		expectedAPIVersion string
		expectedSpec       interface{}
		expectedError      error
	}{
		{
			fileName:           "testdata/application/k8s-app-bluegreen.yaml",
			expectedKind:       KindKubernetesApp,
			expectedAPIVersion: "pipecd.dev/v1beta1",
			expectedSpec: &KubernetesApplicationSpec{
				GenericApplicationSpec: GenericApplicationSpec{
					Description: "application description first string\napplication description second string\n",
					Planner: DeploymentPlanner{
						AlwaysUsePipeline: true,
						AutoRollback:      newBoolPointer(true),
					},
					Pipeline: &DeploymentPipeline{
						Stages: []PipelineStage{
							{
								Name: model.StageK8sCanaryRollout,
								K8sCanaryRolloutStageOptions: &K8sCanaryRolloutStageOptions{
									Replicas: Replicas{
										Number:       100,
										IsPercentage: true,
									},
								},
							},
							{
								Name: model.StageK8sTrafficRouting,
								K8sTrafficRoutingStageOptions: &K8sTrafficRoutingStageOptions{
									Canary: Percentage{
										Number: 100,
									},
								},
							},
							{
								Name:                          model.StageK8sPrimaryRollout,
								K8sPrimaryRolloutStageOptions: &K8sPrimaryRolloutStageOptions{},
							},
							{
								Name: model.StageK8sTrafficRouting,
								K8sTrafficRoutingStageOptions: &K8sTrafficRoutingStageOptions{
									Primary: Percentage{
										Number: 100,
									},
								},
							},
							{
								Name:                       model.StageK8sCanaryClean,
								K8sCanaryCleanStageOptions: &K8sCanaryCleanStageOptions{},
							},
						},
					},
					Timeout: Duration(6 * time.Hour),
					Trigger: Trigger{
						OnCommit: OnCommit{
							Disabled: false,
						},
						OnCommand: OnCommand{
							Disabled: false,
						},
						OnOutOfSync: OnOutOfSync{
							Disabled:  newBoolPointer(true),
							MinWindow: Duration(5 * time.Minute),
						},
						OnChain: OnChain{
							Disabled: newBoolPointer(true),
						},
					},
				},
				Input: KubernetesDeploymentInput{
					AutoRollback: newBoolPointer(true),
				},
				TrafficRouting: &KubernetesTrafficRouting{
					Method: KubernetesTrafficRoutingMethodPodSelector,
				},
				VariantLabel: KubernetesVariantLabel{
					Key:           "pipecd.dev/variant",
					PrimaryValue:  "primary",
					BaselineValue: "baseline",
					CanaryValue:   "canary",
				},
			},
			expectedError: nil,
		},
		{
			fileName:           "testdata/application/k8s-app-resource-route.yaml",
			expectedKind:       KindKubernetesApp,
			expectedAPIVersion: "pipecd.dev/v1beta1",
			expectedSpec: &KubernetesApplicationSpec{
				GenericApplicationSpec: GenericApplicationSpec{
					Timeout: Duration(6 * time.Hour),
					Trigger: Trigger{
						OnCommit: OnCommit{
							Disabled: false,
						},
						OnCommand: OnCommand{
							Disabled: false,
						},
						OnOutOfSync: OnOutOfSync{
							Disabled:  newBoolPointer(true),
							MinWindow: Duration(5 * time.Minute),
						},
						OnChain: OnChain{
							Disabled: newBoolPointer(true),
						},
					},
					Planner: DeploymentPlanner{
						AutoRollback: newBoolPointer(true),
					},
				},
				Input: KubernetesDeploymentInput{
					AutoRollback: newBoolPointer(true),
				},
				VariantLabel: KubernetesVariantLabel{
					Key:           "pipecd.dev/variant",
					PrimaryValue:  "primary",
					BaselineValue: "baseline",
					CanaryValue:   "canary",
				},
				ResourceRoutes: []KubernetesResourceRoute{
					{
						Provider: KubernetesProviderMatcher{
							Name: "ConfigCluster",
						},
						Match: &KubernetesResourceRouteMatcher{
							Kind: "Ingress",
						},
					},
					{
						Provider: KubernetesProviderMatcher{
							Name: "ConfigCluster",
						},
						Match: &KubernetesResourceRouteMatcher{
							Kind: "Service",
							Name: "Foo",
						},
					},
					{
						Provider: KubernetesProviderMatcher{
							Labels: map[string]string{
								"group": "workload",
							},
						},
					},
				},
			},
			expectedError: nil,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.fileName, func(t *testing.T) {
			cfg, err := LoadFromYAML(tc.fileName)
			require.Equal(t, tc.expectedError, err)
			if err == nil {
				assert.Equal(t, tc.expectedKind, cfg.Kind)
				assert.Equal(t, tc.expectedAPIVersion, cfg.APIVersion)
				assert.Equal(t, tc.expectedSpec, cfg.spec)
			}
		})
	}
}
