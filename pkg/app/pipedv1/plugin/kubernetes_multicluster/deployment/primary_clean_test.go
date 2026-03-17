// Copyright 2025 The PipeCD Authors.
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

package deployment

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
	"github.com/pipe-cd/piped-plugin-sdk-go/logpersister/logpersistertest"
	"github.com/pipe-cd/piped-plugin-sdk-go/toolregistry/toolregistrytest"

	kubeconfig "github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes_multicluster/config"
)

func TestPlugin_executeK8sMultiPrimaryCleanStage(t *testing.T) {
	t.Parallel()

	ctx := t.Context()

	testRegistry := toolregistrytest.NewTestToolRegistry(t)

	appCfg := sdk.LoadApplicationConfigForTest[kubeconfig.KubernetesApplicationSpec](t, filepath.Join("testdata", "simple", "app.pipecd.yaml"), "kubernetes_multicluster")

	input := &sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec]{
		Request: sdk.ExecuteStageRequest[kubeconfig.KubernetesApplicationSpec]{
			StageName:   "K8S_PRIMARY_CLEAN",
			StageConfig: []byte(`{}`),
			TargetDeploymentSource: sdk.DeploymentSource[kubeconfig.KubernetesApplicationSpec]{
				ApplicationDirectory:      filepath.Join("testdata", "simple"),
				CommitHash:                "0123456789",
				ApplicationConfig:         appCfg,
				ApplicationConfigFilename: "app.pipecd.yaml",
			},
			Deployment: sdk.Deployment{
				PipedID:       "piped-id",
				ApplicationID: "app-id",
			},
		},
		Client: sdk.NewClient(nil, "kubernetes_multicluster", "", "", logpersistertest.NewTestLogPersister(t), testRegistry),
		Logger: zaptest.NewLogger(t),
	}

	dtConfig, dynamicClient := setupTestDeployTargetConfigAndDynamicClient(t)

	deploymentRes := schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}

	// Pre-create a primary deployment resource in the cluster (simulating what K8S_PRIMARY_ROLLOUT would do).
	primaryDeployment := &unstructured.Unstructured{
		Object: map[string]any{
			"apiVersion": "apps/v1",
			"kind":       "Deployment",
			"metadata": map[string]any{
				"name":      "simple-primary",
				"namespace": "default",
				"labels": map[string]any{
					"app":                    "simple",
					"pipecd.dev/managed-by":  "piped",
					"pipecd.dev/piped":       "piped-id",
					"pipecd.dev/application": "app-id",
					"pipecd.dev/variant":     "primary",
				},
				"annotations": map[string]any{
					"pipecd.dev/managed-by":  "piped",
					"pipecd.dev/application": "app-id",
					"pipecd.dev/variant":     "primary",
				},
			},
			"spec": map[string]any{
				"replicas": int64(2),
				"selector": map[string]any{
					"matchLabels": map[string]any{
						"app":                "simple",
						"pipecd.dev/variant": "primary",
					},
				},
				"template": map[string]any{
					"metadata": map[string]any{
						"labels": map[string]any{
							"app":                "simple",
							"pipecd.dev/variant": "primary",
						},
					},
					"spec": map[string]any{
						"containers": []any{
							map[string]any{
								"name":  "helloworld",
								"image": "ghcr.io/pipe-cd/helloworld:v0.32.0",
							},
						},
					},
				},
			},
		},
	}

	_, err := dynamicClient.Resource(deploymentRes).Namespace("default").Create(ctx, primaryDeployment, metav1.CreateOptions{})
	require.NoError(t, err)

	// Verify the primary deployment exists before running the stage.
	_, err = dynamicClient.Resource(deploymentRes).Namespace("default").Get(ctx, "simple-primary", metav1.GetOptions{})
	require.NoError(t, err)

	plugin := &Plugin{}

	status := plugin.executeK8sMultiPrimaryCleanStage(ctx, input, []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]{
		{
			Name:   "default",
			Config: *dtConfig,
		},
	})

	assert.Equal(t, sdk.StageStatusSuccess, status)

	// Assert that the primary deployment has been deleted.
	_, err = dynamicClient.Resource(deploymentRes).Namespace("default").Get(ctx, "simple-primary", metav1.GetOptions{})
	require.Error(t, err)
	assert.True(t, k8serrors.IsNotFound(err))
}

func TestPlugin_executeK8sMultiPrimaryCleanStage_multipleTargets(t *testing.T) {
	t.Parallel()

	ctx := t.Context()

	testRegistry := toolregistrytest.NewTestToolRegistry(t)

	appCfg := sdk.LoadApplicationConfigForTest[kubeconfig.KubernetesApplicationSpec](t, filepath.Join("testdata", "simple", "app.pipecd.yaml"), "kubernetes_multicluster")

	clusterUS := setupCluster(t, "cluster-us")
	clusterEU := setupCluster(t, "cluster-eu")

	input := &sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec]{
		Request: sdk.ExecuteStageRequest[kubeconfig.KubernetesApplicationSpec]{
			StageName:   "K8S_PRIMARY_CLEAN",
			StageConfig: []byte(`{}`),
			TargetDeploymentSource: sdk.DeploymentSource[kubeconfig.KubernetesApplicationSpec]{
				ApplicationDirectory:      filepath.Join("testdata", "simple"),
				CommitHash:                "0123456789",
				ApplicationConfig:         appCfg,
				ApplicationConfigFilename: "app.pipecd.yaml",
			},
			Deployment: sdk.Deployment{
				PipedID:       "piped-id",
				ApplicationID: "app-id",
			},
		},
		Client: sdk.NewClient(nil, "kubernetes_multicluster", "", "", logpersistertest.NewTestLogPersister(t), testRegistry),
		Logger: zaptest.NewLogger(t),
	}

	deploymentRes := schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}

	// Pre-create primary deployment resources on both clusters.
	for _, c := range []*cluster{clusterUS, clusterEU} {
		primaryDeployment := &unstructured.Unstructured{
			Object: map[string]any{
				"apiVersion": "apps/v1",
				"kind":       "Deployment",
				"metadata": map[string]any{
					"name":      "simple-primary",
					"namespace": "default",
					"labels": map[string]any{
						"app":                    "simple",
						"pipecd.dev/managed-by":  "piped",
						"pipecd.dev/piped":       "piped-id",
						"pipecd.dev/application": "app-id",
						"pipecd.dev/variant":     "primary",
					},
					"annotations": map[string]any{
						"pipecd.dev/managed-by":  "piped",
						"pipecd.dev/application": "app-id",
						"pipecd.dev/variant":     "primary",
					},
				},
				"spec": map[string]any{
					"replicas": int64(2),
					"selector": map[string]any{
						"matchLabels": map[string]any{
							"app":                "simple",
							"pipecd.dev/variant": "primary",
						},
					},
					"template": map[string]any{
						"metadata": map[string]any{
							"labels": map[string]any{
								"app":                "simple",
								"pipecd.dev/variant": "primary",
							},
						},
						"spec": map[string]any{
							"containers": []any{
								map[string]any{
									"name":  "helloworld",
									"image": "ghcr.io/pipe-cd/helloworld:v0.32.0",
								},
							},
						},
					},
				},
			},
		}
		_, err := c.cli.Resource(deploymentRes).Namespace("default").Create(ctx, primaryDeployment, metav1.CreateOptions{})
		require.NoError(t, err)
	}

	plugin := &Plugin{}

	status := plugin.executeK8sMultiPrimaryCleanStage(ctx, input, []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]{
		{Name: clusterUS.name, Config: *clusterUS.dtc},
		{Name: clusterEU.name, Config: *clusterEU.dtc},
	})

	assert.Equal(t, sdk.StageStatusSuccess, status)

	// Assert that the primary deployments have been deleted on both clusters.
	for _, c := range []*cluster{clusterUS, clusterEU} {
		_, err := c.cli.Resource(deploymentRes).Namespace("default").Get(ctx, "simple-primary", metav1.GetOptions{})
		require.Error(t, err)
		assert.True(t, k8serrors.IsNotFound(err), "primary deployment should be deleted on cluster %s", c.name)
	}
}
