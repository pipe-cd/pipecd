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
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
	"github.com/pipe-cd/piped-plugin-sdk-go/logpersister/logpersistertest"
	"github.com/pipe-cd/piped-plugin-sdk-go/toolregistry/toolregistrytest"

	kubeconfig "github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes_multicluster/config"
)

func TestPlugin_executeK8sMultiPrimaryRolloutStage_SingleCluster(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	appCfg := sdk.LoadApplicationConfigForTest[kubeconfig.KubernetesApplicationSpec](t, filepath.Join("testdata", "primary_rollout", "app.pipecd.yaml"), "kubernetes_multicluster")

	testRegistry := toolregistrytest.NewTestToolRegistry(t)

	input := &sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec]{
		Request: sdk.ExecuteStageRequest[kubeconfig.KubernetesApplicationSpec]{
			StageName:   StageK8sMultiPrimaryRollout,
			StageConfig: []byte(`{}`),
			TargetDeploymentSource: sdk.DeploymentSource[kubeconfig.KubernetesApplicationSpec]{
				ApplicationDirectory:      filepath.Join("testdata", "primary_rollout"),
				CommitHash:                "0123456789",
				ApplicationConfig:         appCfg,
				ApplicationConfigFilename: "app.pipecd.yaml",
			},
			Deployment: sdk.Deployment{
				PipedID:       "piped-id",
				ApplicationID: "app-id",
			},
		},
		Client: sdk.NewClient(nil, "kubernetes_multicluster", "app-id", "stage-id", logpersistertest.NewTestLogPersister(t), testRegistry),
		Logger: zaptest.NewLogger(t),
	}

	dtConfig, dynamicClient := setupTestDeployTargetConfigAndDynamicClient(t)

	plugin := &Plugin{}
	status := plugin.executeK8sMultiPrimaryRolloutStage(ctx, input, []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]{
		{Name: "default", Config: *dtConfig},
	})

	assert.Equal(t, sdk.StageStatusSuccess, status)

	deploymentRes := schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}
	deployment, err := dynamicClient.Resource(deploymentRes).Namespace("default").Get(ctx, "simple", metav1.GetOptions{})
	require.NoError(t, err)

	assert.Equal(t, "simple", deployment.GetName())
	assert.Equal(t, "primary", deployment.GetLabels()["pipecd.dev/variant"])
	assert.Equal(t, "primary", deployment.GetAnnotations()["pipecd.dev/variant"])
	assert.Equal(t, "piped-id", deployment.GetLabels()["pipecd.dev/piped"])
	assert.Equal(t, "app-id", deployment.GetLabels()["pipecd.dev/application"])
}

func TestPlugin_executeK8sMultiPrimaryRolloutStage_MultiCluster(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	appCfg := sdk.LoadApplicationConfigForTest[kubeconfig.KubernetesApplicationSpec](t, filepath.Join("testdata", "primary_rollout", "app.pipecd.yaml"), "kubernetes_multicluster")

	testRegistry := toolregistrytest.NewTestToolRegistry(t)

	input := &sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec]{
		Request: sdk.ExecuteStageRequest[kubeconfig.KubernetesApplicationSpec]{
			StageName:   StageK8sMultiPrimaryRollout,
			StageConfig: []byte(`{}`),
			TargetDeploymentSource: sdk.DeploymentSource[kubeconfig.KubernetesApplicationSpec]{
				ApplicationDirectory:      filepath.Join("testdata", "primary_rollout"),
				CommitHash:                "0123456789",
				ApplicationConfig:         appCfg,
				ApplicationConfigFilename: "app.pipecd.yaml",
			},
			Deployment: sdk.Deployment{
				PipedID:       "piped-id",
				ApplicationID: "app-id",
			},
		},
		Client: sdk.NewClient(nil, "kubernetes_multicluster", "app-id", "stage-id", logpersistertest.NewTestLogPersister(t), testRegistry),
		Logger: zaptest.NewLogger(t),
	}

	cluster1 := setupCluster(t, "cluster1")
	cluster2 := setupCluster(t, "cluster2")

	dts := []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]{
		{Name: "cluster1", Config: *cluster1.dtc},
		{Name: "cluster2", Config: *cluster2.dtc},
	}

	plugin := &Plugin{}
	status := plugin.executeK8sMultiPrimaryRolloutStage(ctx, input, dts)

	require.Equal(t, sdk.StageStatusSuccess, status)

	deploymentRes := schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}

	// Both clusters should have the primary deployment.
	for _, cl := range []*cluster{cluster1, cluster2} {
		deployment, err := cl.cli.Resource(deploymentRes).Namespace("default").Get(ctx, "simple", metav1.GetOptions{})
		require.NoError(t, err)

		assert.Equal(t, "simple", deployment.GetName())
		assert.Equal(t, "primary", deployment.GetLabels()["pipecd.dev/variant"])
		assert.Equal(t, "piped-id", deployment.GetLabels()["pipecd.dev/piped"])
		assert.Equal(t, "app-id", deployment.GetLabels()["pipecd.dev/application"])
	}
}

func TestPlugin_executeK8sMultiPrimaryRolloutStage_WithCreateService(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	configDir := filepath.Join("testdata", "primary_rollout_with_create_service")
	appCfg := sdk.LoadApplicationConfigForTest[kubeconfig.KubernetesApplicationSpec](t, filepath.Join(configDir, "app.pipecd.yaml"), "kubernetes_multicluster")

	testRegistry := toolregistrytest.NewTestToolRegistry(t)

	input := &sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec]{
		Request: sdk.ExecuteStageRequest[kubeconfig.KubernetesApplicationSpec]{
			StageName:   StageK8sMultiPrimaryRollout,
			StageConfig: []byte(`{"createService": true}`),
			TargetDeploymentSource: sdk.DeploymentSource[kubeconfig.KubernetesApplicationSpec]{
				ApplicationDirectory:      configDir,
				CommitHash:                "0123456789",
				ApplicationConfig:         appCfg,
				ApplicationConfigFilename: "app.pipecd.yaml",
			},
			Deployment: sdk.Deployment{
				PipedID:       "piped-id",
				ApplicationID: "app-id",
			},
		},
		Client: sdk.NewClient(nil, "kubernetes_multicluster", "app-id", "stage-id", logpersistertest.NewTestLogPersister(t), testRegistry),
		Logger: zaptest.NewLogger(t),
	}

	dtConfig, dynamicClient := setupTestDeployTargetConfigAndDynamicClient(t)

	plugin := &Plugin{}
	status := plugin.executeK8sMultiPrimaryRolloutStage(ctx, input, []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]{
		{Name: "default", Config: *dtConfig},
	})

	assert.Equal(t, sdk.StageStatusSuccess, status)

	// Primary deployment should exist.
	deploymentRes := schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}
	deployment, err := dynamicClient.Resource(deploymentRes).Namespace("default").Get(ctx, "simple", metav1.GetOptions{})
	require.NoError(t, err)
	assert.Equal(t, "simple", deployment.GetName())
	assert.Equal(t, "primary", deployment.GetLabels()["pipecd.dev/variant"])

	// Primary variant service should be created with variant selector added.
	serviceRes := schema.GroupVersionResource{Group: "", Version: "v1", Resource: "services"}
	service, err := dynamicClient.Resource(serviceRes).Namespace("default").Get(ctx, "simple-primary", metav1.GetOptions{})
	require.NoError(t, err)
	assert.Equal(t, "simple-primary", service.GetName())

	selector, found, err := unstructured.NestedStringMap(service.Object, "spec", "selector")
	require.NoError(t, err)
	require.True(t, found)
	assert.Equal(t, map[string]string{"app": "simple", "pipecd.dev/variant": "primary"}, selector)
}

func TestPlugin_executeK8sMultiPrimaryRolloutStage_WithAddVariantLabelToSelector(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	appCfg := sdk.LoadApplicationConfigForTest[kubeconfig.KubernetesApplicationSpec](t, filepath.Join("testdata", "primary_rollout", "app.pipecd.yaml"), "kubernetes_multicluster")

	testRegistry := toolregistrytest.NewTestToolRegistry(t)

	input := &sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec]{
		Request: sdk.ExecuteStageRequest[kubeconfig.KubernetesApplicationSpec]{
			StageName:   StageK8sMultiPrimaryRollout,
			StageConfig: []byte(`{"addVariantLabelToSelector": true}`),
			TargetDeploymentSource: sdk.DeploymentSource[kubeconfig.KubernetesApplicationSpec]{
				ApplicationDirectory:      filepath.Join("testdata", "primary_rollout"),
				CommitHash:                "0123456789",
				ApplicationConfig:         appCfg,
				ApplicationConfigFilename: "app.pipecd.yaml",
			},
			Deployment: sdk.Deployment{
				PipedID:       "piped-id",
				ApplicationID: "app-id",
			},
		},
		Client: sdk.NewClient(nil, "kubernetes_multicluster", "app-id", "stage-id", logpersistertest.NewTestLogPersister(t), testRegistry),
		Logger: zaptest.NewLogger(t),
	}

	dtConfig, dynamicClient := setupTestDeployTargetConfigAndDynamicClient(t)

	plugin := &Plugin{}
	status := plugin.executeK8sMultiPrimaryRolloutStage(ctx, input, []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]{
		{Name: "default", Config: *dtConfig},
	})

	assert.Equal(t, sdk.StageStatusSuccess, status)

	deploymentRes := schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}
	deployment, err := dynamicClient.Resource(deploymentRes).Namespace("default").Get(ctx, "simple", metav1.GetOptions{})
	require.NoError(t, err)

	// Variant label should be present in spec.selector.matchLabels.
	matchLabels, found, err := unstructured.NestedStringMap(deployment.Object, "spec", "selector", "matchLabels")
	require.NoError(t, err)
	require.True(t, found)
	assert.Equal(t, "primary", matchLabels["pipecd.dev/variant"])
}

func TestPlugin_executeK8sMultiPrimaryRolloutStage_WithPrune(t *testing.T) {
	t.Parallel()

	ctx := t.Context()

	testRegistry := toolregistrytest.NewTestToolRegistry(t)
	dtConfig, dynamicClient := setupTestDeployTargetConfigAndDynamicClient(t)

	deploymentRes := schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}
	serviceRes := schema.GroupVersionResource{Group: "", Version: "v1", Resource: "services"}

	runningOk := t.Run("prepare running state", func(t *testing.T) {
		running := filepath.Join("testdata", "primary_rollout_prune", "running")
		runningCfg := sdk.LoadApplicationConfigForTest[kubeconfig.KubernetesApplicationSpec](t, filepath.Join(running, "app.pipecd.yaml"), "kubernetes_multicluster")

		runningInput := &sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec]{
			Request: sdk.ExecuteStageRequest[kubeconfig.KubernetesApplicationSpec]{
				StageName:   StageK8sMultiPrimaryRollout,
				StageConfig: []byte(`{"prune": true}`),
				TargetDeploymentSource: sdk.DeploymentSource[kubeconfig.KubernetesApplicationSpec]{
					ApplicationDirectory:      running,
					CommitHash:                "0123456789",
					ApplicationConfig:         runningCfg,
					ApplicationConfigFilename: "app.pipecd.yaml",
				},
				Deployment: sdk.Deployment{
					PipedID:       "piped-id",
					ApplicationID: "app-id",
				},
			},
			Client: sdk.NewClient(nil, "kubernetes_multicluster", "app-id", "stage-id", logpersistertest.NewTestLogPersister(t), testRegistry),
			Logger: zaptest.NewLogger(t),
		}

		plugin := &Plugin{}
		status := plugin.executeK8sMultiPrimaryRolloutStage(ctx, runningInput, []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]{
			{Name: "default", Config: *dtConfig},
		})
		assert.Equal(t, sdk.StageStatusSuccess, status)

		// Both deployment and service should exist after running state deployment.
		_, err := dynamicClient.Resource(deploymentRes).Namespace("default").Get(ctx, "simple", metav1.GetOptions{})
		assert.NoError(t, err)
		_, err = dynamicClient.Resource(serviceRes).Namespace("default").Get(ctx, "simple", metav1.GetOptions{})
		assert.NoError(t, err)
	})
	require.True(t, runningOk, "prepare running state subtest failed, aborting")

	t.Run("prune with target state", func(t *testing.T) {
		target := filepath.Join("testdata", "primary_rollout_prune", "target")
		targetCfg := sdk.LoadApplicationConfigForTest[kubeconfig.KubernetesApplicationSpec](t, filepath.Join(target, "app.pipecd.yaml"), "kubernetes_multicluster")

		running := filepath.Join("testdata", "primary_rollout_prune", "running")
		runningCfg := sdk.LoadApplicationConfigForTest[kubeconfig.KubernetesApplicationSpec](t, filepath.Join(running, "app.pipecd.yaml"), "kubernetes_multicluster")

		targetInput := &sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec]{
			Request: sdk.ExecuteStageRequest[kubeconfig.KubernetesApplicationSpec]{
				StageName:   StageK8sMultiPrimaryRollout,
				StageConfig: []byte(`{"prune": true}`),
				RunningDeploymentSource: sdk.DeploymentSource[kubeconfig.KubernetesApplicationSpec]{
					ApplicationDirectory:      running,
					CommitHash:                "0123456789",
					ApplicationConfig:         runningCfg,
					ApplicationConfigFilename: "app.pipecd.yaml",
				},
				TargetDeploymentSource: sdk.DeploymentSource[kubeconfig.KubernetesApplicationSpec]{
					ApplicationDirectory:      target,
					CommitHash:                "0012345678",
					ApplicationConfig:         targetCfg,
					ApplicationConfigFilename: "app.pipecd.yaml",
				},
				Deployment: sdk.Deployment{
					PipedID:       "piped-id",
					ApplicationID: "app-id",
				},
			},
			Client: sdk.NewClient(nil, "kubernetes_multicluster", "app-id", "stage-id", logpersistertest.NewTestLogPersister(t), testRegistry),
			Logger: zaptest.NewLogger(t),
		}

		plugin := &Plugin{}
		status := plugin.executeK8sMultiPrimaryRolloutStage(ctx, targetInput, []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]{
			{Name: "default", Config: *dtConfig},
		})
		assert.Equal(t, sdk.StageStatusSuccess, status)

		// Deployment should still exist.
		_, err := dynamicClient.Resource(deploymentRes).Namespace("default").Get(ctx, "simple", metav1.GetOptions{})
		assert.NoError(t, err)

		// Service should have been pruned because it's not in the target manifests.
		_, err = dynamicClient.Resource(serviceRes).Namespace("default").Get(ctx, "simple", metav1.GetOptions{})
		require.Error(t, err)
		assert.True(t, apierrors.IsNotFound(err), "expected service to be pruned, but got %v", err)
	})
}

func TestPlugin_executeK8sMultiPrimaryRolloutStage_WithPrune_ManualPreCreate(t *testing.T) {
	t.Parallel()

	ctx := t.Context()

	appCfg := sdk.LoadApplicationConfigForTest[kubeconfig.KubernetesApplicationSpec](t, filepath.Join("testdata", "primary_rollout", "app.pipecd.yaml"), "kubernetes_multicluster")
	testRegistry := toolregistrytest.NewTestToolRegistry(t)

	dtConfig, dynamicClient := setupTestDeployTargetConfigAndDynamicClient(t)

	deploymentRes := schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}

	// Pre-create a stale primary deployment that is NOT in the target manifests.
	staleDeployment := &unstructured.Unstructured{
		Object: map[string]any{
			"apiVersion": "apps/v1",
			"kind":       "Deployment",
			"metadata": map[string]any{
				"name":      "simple-stale",
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
				"replicas": int64(1),
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
								"image": "ghcr.io/pipe-cd/helloworld:v0.31.0",
							},
						},
					},
				},
			},
		},
	}

	_, err := dynamicClient.Resource(deploymentRes).Namespace("default").Create(ctx, staleDeployment, metav1.CreateOptions{})
	require.NoError(t, err)

	input := &sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec]{
		Request: sdk.ExecuteStageRequest[kubeconfig.KubernetesApplicationSpec]{
			StageName:   StageK8sMultiPrimaryRollout,
			StageConfig: []byte(`{"prune": true}`),
			TargetDeploymentSource: sdk.DeploymentSource[kubeconfig.KubernetesApplicationSpec]{
				ApplicationDirectory:      filepath.Join("testdata", "primary_rollout"),
				CommitHash:                "0123456789",
				ApplicationConfig:         appCfg,
				ApplicationConfigFilename: "app.pipecd.yaml",
			},
			Deployment: sdk.Deployment{
				PipedID:       "piped-id",
				ApplicationID: "app-id",
			},
		},
		Client: sdk.NewClient(nil, "kubernetes_multicluster", "app-id", "stage-id", logpersistertest.NewTestLogPersister(t), testRegistry),
		Logger: zaptest.NewLogger(t),
	}

	plugin := &Plugin{}
	status := plugin.executeK8sMultiPrimaryRolloutStage(ctx, input, []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]{
		{Name: "default", Config: *dtConfig},
	})

	assert.Equal(t, sdk.StageStatusSuccess, status)

	// The target deployment should exist.
	_, err = dynamicClient.Resource(deploymentRes).Namespace("default").Get(ctx, "simple", metav1.GetOptions{})
	assert.NoError(t, err)

	// The stale deployment should have been pruned.
	_, err = dynamicClient.Resource(deploymentRes).Namespace("default").Get(ctx, "simple-stale", metav1.GetOptions{})
	require.Error(t, err)
	assert.True(t, apierrors.IsNotFound(err), "expected stale deployment to be pruned, but got: %v", err)
}
