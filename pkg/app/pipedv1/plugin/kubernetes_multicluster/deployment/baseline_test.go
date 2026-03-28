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

// baselineDeployment builds a pre-created baseline Deployment for test setup.
func baselineDeployment(namespace string) *unstructured.Unstructured {
	return &unstructured.Unstructured{
		Object: map[string]any{
			"apiVersion": "apps/v1",
			"kind":       "Deployment",
			"metadata": map[string]any{
				"name":      "simple-baseline",
				"namespace": namespace,
				"labels": map[string]any{
					"app":                    "simple",
					"pipecd.dev/managed-by":  "piped",
					"pipecd.dev/piped":       "piped-id",
					"pipecd.dev/application": "app-id",
					"pipecd.dev/variant":     "baseline",
				},
				"annotations": map[string]any{
					"pipecd.dev/managed-by":  "piped",
					"pipecd.dev/application": "app-id",
					"pipecd.dev/variant":     "baseline",
				},
			},
			"spec": map[string]any{
				"replicas": int64(1),
				"selector": map[string]any{
					"matchLabels": map[string]any{
						"app":                "simple",
						"pipecd.dev/variant": "baseline",
					},
				},
				"template": map[string]any{
					"metadata": map[string]any{
						"labels": map[string]any{
							"app":                "simple",
							"pipecd.dev/variant": "baseline",
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
}

// baselineService builds a pre-created baseline Service for test setup.
func baselineService(namespace string) *unstructured.Unstructured {
	return &unstructured.Unstructured{
		Object: map[string]any{
			"apiVersion": "v1",
			"kind":       "Service",
			"metadata": map[string]any{
				"name":      "simple-baseline",
				"namespace": namespace,
				"labels": map[string]any{
					"app":                    "simple",
					"pipecd.dev/managed-by":  "piped",
					"pipecd.dev/piped":       "piped-id",
					"pipecd.dev/application": "app-id",
					"pipecd.dev/variant":     "baseline",
				},
				"annotations": map[string]any{
					"pipecd.dev/managed-by":  "piped",
					"pipecd.dev/application": "app-id",
					"pipecd.dev/variant":     "baseline",
				},
			},
			"spec": map[string]any{
				"selector": map[string]any{
					"app":                "simple",
					"pipecd.dev/variant": "baseline",
				},
				"ports": []any{
					map[string]any{
						"protocol":   "TCP",
						"port":       int64(9085),
						"targetPort": int64(9085),
					},
				},
			},
		},
	}
}

func TestPlugin_executeK8sMultiBaselineRolloutStage_SingleCluster(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	configDir := filepath.Join("testdata", "baseline_rollout")
	appCfg := sdk.LoadApplicationConfigForTest[kubeconfig.KubernetesApplicationSpec](t, filepath.Join(configDir, "app.pipecd.yaml"), "kubernetes_multicluster")

	testRegistry := toolregistrytest.NewTestToolRegistry(t)

	input := &sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec]{
		Request: sdk.ExecuteStageRequest[kubeconfig.KubernetesApplicationSpec]{
			StageName:   StageK8sMultiBaselineRollout,
			StageConfig: []byte(`{"replicas": "50%", "suffix": "baseline"}`),
			RunningDeploymentSource: sdk.DeploymentSource[kubeconfig.KubernetesApplicationSpec]{
				ApplicationDirectory:      configDir,
				CommitHash:                "0123456789",
				ApplicationConfig:         appCfg,
				ApplicationConfigFilename: "app.pipecd.yaml",
			},
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
	status := plugin.executeK8sMultiBaselineRolloutStage(ctx, input, []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]{
		{Name: "default", Config: *dtConfig},
	})

	assert.Equal(t, sdk.StageStatusSuccess, status)

	// The baseline deployment should be created with "-baseline" suffix.
	deploymentRes := schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}
	deployment, err := dynamicClient.Resource(deploymentRes).Namespace("default").Get(ctx, "simple-baseline", metav1.GetOptions{})
	require.NoError(t, err)

	assert.Equal(t, "simple-baseline", deployment.GetName())

	// Verify variant label is set to "baseline".
	assert.Equal(t, "baseline", deployment.GetLabels()["pipecd.dev/variant"])
	assert.Equal(t, "baseline", deployment.GetAnnotations()["pipecd.dev/variant"])

	// Verify replica count is 1 (50% of 2 = 1).
	spec, ok := deployment.Object["spec"].(map[string]interface{})
	require.True(t, ok)
	replicas, ok := spec["replicas"].(int64)
	require.True(t, ok)
	assert.Equal(t, int64(1), replicas)
}

func TestPlugin_executeK8sMultiBaselineRolloutStage_MultiCluster(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	configDir := filepath.Join("testdata", "baseline_rollout")
	appCfg := sdk.LoadApplicationConfigForTest[kubeconfig.KubernetesApplicationSpec](t, filepath.Join(configDir, "app.pipecd.yaml"), "kubernetes_multicluster")

	testRegistry := toolregistrytest.NewTestToolRegistry(t)

	input := &sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec]{
		Request: sdk.ExecuteStageRequest[kubeconfig.KubernetesApplicationSpec]{
			StageName:   StageK8sMultiBaselineRollout,
			StageConfig: []byte(`{"replicas": 1, "suffix": "baseline"}`),
			RunningDeploymentSource: sdk.DeploymentSource[kubeconfig.KubernetesApplicationSpec]{
				ApplicationDirectory:      configDir,
				CommitHash:                "0123456789",
				ApplicationConfig:         appCfg,
				ApplicationConfigFilename: "app.pipecd.yaml",
			},
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

	cluster1 := setupCluster(t, "cluster1")
	cluster2 := setupCluster(t, "cluster2")

	dts := []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]{
		{Name: "cluster1", Config: *cluster1.dtc},
		{Name: "cluster2", Config: *cluster2.dtc},
	}

	plugin := &Plugin{}
	status := plugin.executeK8sMultiBaselineRolloutStage(ctx, input, dts)

	require.Equal(t, sdk.StageStatusSuccess, status)

	deploymentRes := schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}

	// Both clusters should have a baseline deployment.
	for _, cl := range []*cluster{cluster1, cluster2} {
		deployment, err := cl.cli.Resource(deploymentRes).Namespace("default").Get(ctx, "simple-baseline", metav1.GetOptions{})
		require.NoError(t, err)

		assert.Equal(t, "simple-baseline", deployment.GetName())
		assert.Equal(t, "baseline", deployment.GetLabels()["pipecd.dev/variant"])
		assert.Equal(t, "piped-id", deployment.GetLabels()["pipecd.dev/piped"])
		assert.Equal(t, "app-id", deployment.GetLabels()["pipecd.dev/application"])
	}
}

func TestPlugin_executeK8sMultiBaselineRolloutStage_WithCreateService(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	configDir := filepath.Join("testdata", "baseline_rollout_with_create_service")
	appCfg := sdk.LoadApplicationConfigForTest[kubeconfig.KubernetesApplicationSpec](t, filepath.Join(configDir, "app.pipecd.yaml"), "kubernetes_multicluster")

	testRegistry := toolregistrytest.NewTestToolRegistry(t)

	input := &sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec]{
		Request: sdk.ExecuteStageRequest[kubeconfig.KubernetesApplicationSpec]{
			StageName:   StageK8sMultiBaselineRollout,
			StageConfig: []byte(`{"replicas": "50%", "createService": true}`),
			RunningDeploymentSource: sdk.DeploymentSource[kubeconfig.KubernetesApplicationSpec]{
				ApplicationDirectory:      configDir,
				CommitHash:                "0123456789",
				ApplicationConfig:         appCfg,
				ApplicationConfigFilename: "app.pipecd.yaml",
			},
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
	status := plugin.executeK8sMultiBaselineRolloutStage(ctx, input, []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]{
		{Name: "default", Config: *dtConfig},
	})

	assert.Equal(t, sdk.StageStatusSuccess, status)

	// Baseline deployment should be created with variant labels.
	deploymentRes := schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}
	deployment, err := dynamicClient.Resource(deploymentRes).Namespace("default").Get(ctx, "simple-baseline", metav1.GetOptions{})
	require.NoError(t, err)
	assert.Equal(t, "simple-baseline", deployment.GetName())
	assert.Equal(t, "simple", deployment.GetLabels()["app"])
	assert.Equal(t, "baseline", deployment.GetLabels()["pipecd.dev/variant"])
	assert.Equal(t, "baseline", deployment.GetAnnotations()["pipecd.dev/variant"])
	assert.Equal(t, "piped-id", deployment.GetLabels()["pipecd.dev/piped"])
	assert.Equal(t, "app-id", deployment.GetLabels()["pipecd.dev/application"])

	// Baseline service should be created with variant selector added.
	serviceRes := schema.GroupVersionResource{Group: "", Version: "v1", Resource: "services"}
	service, err := dynamicClient.Resource(serviceRes).Namespace("default").Get(ctx, "simple-baseline", metav1.GetOptions{})
	require.NoError(t, err)
	assert.Equal(t, "simple-baseline", service.GetName())

	selector, found, err := unstructured.NestedStringMap(service.Object, "spec", "selector")
	require.NoError(t, err)
	require.True(t, found)
	assert.Equal(t, map[string]string{"app": "simple", "pipecd.dev/variant": "baseline"}, selector)

	ports, found, err := unstructured.NestedSlice(service.Object, "spec", "ports")
	require.NoError(t, err)
	require.True(t, found)
	require.Len(t, ports, 1)
	port := ports[0].(map[string]any)
	assert.Equal(t, int64(9085), port["port"])
	assert.Equal(t, int64(9085), port["targetPort"])
}

func TestPlugin_executeK8sMultiBaselineRolloutStage_WithoutCreateService(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	configDir := filepath.Join("testdata", "baseline_rollout_without_create_service")
	appCfg := sdk.LoadApplicationConfigForTest[kubeconfig.KubernetesApplicationSpec](t, filepath.Join(configDir, "app.pipecd.yaml"), "kubernetes_multicluster")

	testRegistry := toolregistrytest.NewTestToolRegistry(t)

	input := &sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec]{
		Request: sdk.ExecuteStageRequest[kubeconfig.KubernetesApplicationSpec]{
			StageName:   StageK8sMultiBaselineRollout,
			StageConfig: []byte(`{"replicas": "50%"}`),
			RunningDeploymentSource: sdk.DeploymentSource[kubeconfig.KubernetesApplicationSpec]{
				ApplicationDirectory:      configDir,
				CommitHash:                "0123456789",
				ApplicationConfig:         appCfg,
				ApplicationConfigFilename: "app.pipecd.yaml",
			},
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
	status := plugin.executeK8sMultiBaselineRolloutStage(ctx, input, []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]{
		{Name: "default", Config: *dtConfig},
	})

	assert.Equal(t, sdk.StageStatusSuccess, status)

	// Baseline deployment should be created.
	deploymentRes := schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}
	deployment, err := dynamicClient.Resource(deploymentRes).Namespace("default").Get(ctx, "simple-baseline", metav1.GetOptions{})
	require.NoError(t, err)
	assert.Equal(t, "simple-baseline", deployment.GetName())
	assert.Equal(t, "simple", deployment.GetLabels()["app"])
	assert.Equal(t, "baseline", deployment.GetLabels()["pipecd.dev/variant"])

	// No baseline service should be created when createService is false.
	serviceRes := schema.GroupVersionResource{Group: "", Version: "v1", Resource: "services"}
	_, err = dynamicClient.Resource(serviceRes).Namespace("default").Get(ctx, "simple-baseline", metav1.GetOptions{})
	require.Error(t, err)
	assert.True(t, k8serrors.IsNotFound(err))
}

func TestPlugin_executeK8sMultiBaselineRolloutStage_Failure(t *testing.T) {
	t.Parallel()

	configDir := filepath.Join("testdata", "baseline_rollout")
	appCfg := sdk.LoadApplicationConfigForTest[kubeconfig.KubernetesApplicationSpec](t, filepath.Join(configDir, "app.pipecd.yaml"), "kubernetes_multicluster")

	testRegistry := toolregistrytest.NewTestToolRegistry(t)

	input := &sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec]{
		Request: sdk.ExecuteStageRequest[kubeconfig.KubernetesApplicationSpec]{
			StageName:   StageK8sMultiBaselineRollout,
			StageConfig: []byte(`{"replicas": 1}`),
			RunningDeploymentSource: sdk.DeploymentSource[kubeconfig.KubernetesApplicationSpec]{
				ApplicationDirectory:      configDir,
				CommitHash:                "0123456789",
				ApplicationConfig:         appCfg,
				ApplicationConfigFilename: "app.pipecd.yaml",
			},
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

	// Provide a bad kubeconfig path.
	dts := []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]{
		{
			Name: "bad-cluster",
			Config: kubeconfig.KubernetesDeployTargetConfig{
				KubeConfigPath: "/nonexistent/kubeconfig",
			},
		},
	}

	plugin := &Plugin{}
	status := plugin.executeK8sMultiBaselineRolloutStage(t.Context(), input, dts)

	assert.Equal(t, sdk.StageStatusFailure, status)
}

func TestPlugin_executeK8sMultiBaselineCleanStage(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	testRegistry := toolregistrytest.NewTestToolRegistry(t)

	appCfg := sdk.LoadApplicationConfigForTest[kubeconfig.KubernetesApplicationSpec](t, filepath.Join("testdata", "simple", "app.pipecd.yaml"), "kubernetes_multicluster")

	input := &sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec]{
		Request: sdk.ExecuteStageRequest[kubeconfig.KubernetesApplicationSpec]{
			StageName:   StageK8sMultiBaselineClean,
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

	// Pre-create the baseline deployment (simulating what K8S_BASELINE_ROLLOUT would have done).
	_, err := dynamicClient.Resource(deploymentRes).Namespace("default").Create(ctx, baselineDeployment("default"), metav1.CreateOptions{})
	require.NoError(t, err)

	_, err = dynamicClient.Resource(deploymentRes).Namespace("default").Get(ctx, "simple-baseline", metav1.GetOptions{})
	require.NoError(t, err)

	plugin := &Plugin{}
	status := plugin.executeK8sMultiBaselineCleanStage(ctx, input, []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]{
		{Name: "default", Config: *dtConfig},
	})

	assert.Equal(t, sdk.StageStatusSuccess, status)

	_, err = dynamicClient.Resource(deploymentRes).Namespace("default").Get(ctx, "simple-baseline", metav1.GetOptions{})
	require.Error(t, err)
	assert.True(t, k8serrors.IsNotFound(err))
}

func TestPlugin_executeK8sMultiBaselineCleanStage_multipleTargets(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	testRegistry := toolregistrytest.NewTestToolRegistry(t)

	appCfg := sdk.LoadApplicationConfigForTest[kubeconfig.KubernetesApplicationSpec](t, filepath.Join("testdata", "simple", "app.pipecd.yaml"), "kubernetes_multicluster")

	clusterUS := setupCluster(t, "cluster-us")
	clusterEU := setupCluster(t, "cluster-eu")

	input := &sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec]{
		Request: sdk.ExecuteStageRequest[kubeconfig.KubernetesApplicationSpec]{
			StageName:   StageK8sMultiBaselineClean,
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

	for _, c := range []*cluster{clusterUS, clusterEU} {
		_, err := c.cli.Resource(deploymentRes).Namespace("default").Create(ctx, baselineDeployment("default"), metav1.CreateOptions{})
		require.NoError(t, err)
	}

	plugin := &Plugin{}
	status := plugin.executeK8sMultiBaselineCleanStage(ctx, input, []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]{
		{Name: clusterUS.name, Config: *clusterUS.dtc},
		{Name: clusterEU.name, Config: *clusterEU.dtc},
	})

	assert.Equal(t, sdk.StageStatusSuccess, status)

	for _, c := range []*cluster{clusterUS, clusterEU} {
		_, err := c.cli.Resource(deploymentRes).Namespace("default").Get(ctx, "simple-baseline", metav1.GetOptions{})
		require.Error(t, err)
		assert.True(t, k8serrors.IsNotFound(err), "baseline deployment should be deleted on cluster %s", c.name)
	}
}

func TestPlugin_executeK8sMultiBaselineCleanStage_withCreateService(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	testRegistry := toolregistrytest.NewTestToolRegistry(t)

	configDir := filepath.Join("testdata", "baseline_clean_with_create_service")
	appCfg := sdk.LoadApplicationConfigForTest[kubeconfig.KubernetesApplicationSpec](t, filepath.Join(configDir, "app.pipecd.yaml"), "kubernetes_multicluster")

	input := &sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec]{
		Request: sdk.ExecuteStageRequest[kubeconfig.KubernetesApplicationSpec]{
			StageName:   StageK8sMultiBaselineClean,
			StageConfig: []byte(`{}`),
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
		Client: sdk.NewClient(nil, "kubernetes_multicluster", "", "", logpersistertest.NewTestLogPersister(t), testRegistry),
		Logger: zaptest.NewLogger(t),
	}

	dtConfig, dynamicClient := setupTestDeployTargetConfigAndDynamicClient(t)

	deploymentRes := schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}
	serviceRes := schema.GroupVersionResource{Group: "", Version: "v1", Resource: "services"}

	// Pre-create both baseline deployment and service.
	_, err := dynamicClient.Resource(deploymentRes).Namespace("default").Create(ctx, baselineDeployment("default"), metav1.CreateOptions{})
	require.NoError(t, err)
	_, err = dynamicClient.Resource(serviceRes).Namespace("default").Create(ctx, baselineService("default"), metav1.CreateOptions{})
	require.NoError(t, err)

	plugin := &Plugin{}
	status := plugin.executeK8sMultiBaselineCleanStage(ctx, input, []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]{
		{Name: "default", Config: *dtConfig},
	})

	assert.Equal(t, sdk.StageStatusSuccess, status)

	_, err = dynamicClient.Resource(deploymentRes).Namespace("default").Get(ctx, "simple-baseline", metav1.GetOptions{})
	require.Error(t, err)
	assert.True(t, k8serrors.IsNotFound(err))

	_, err = dynamicClient.Resource(serviceRes).Namespace("default").Get(ctx, "simple-baseline", metav1.GetOptions{})
	require.Error(t, err)
	assert.True(t, k8serrors.IsNotFound(err))
}

func TestPlugin_executeK8sMultiBaselineCleanStage_withoutCreateService(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	testRegistry := toolregistrytest.NewTestToolRegistry(t)

	configDir := filepath.Join("testdata", "baseline_clean_without_create_service")
	appCfg := sdk.LoadApplicationConfigForTest[kubeconfig.KubernetesApplicationSpec](t, filepath.Join(configDir, "app.pipecd.yaml"), "kubernetes_multicluster")

	input := &sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec]{
		Request: sdk.ExecuteStageRequest[kubeconfig.KubernetesApplicationSpec]{
			StageName:   StageK8sMultiBaselineClean,
			StageConfig: []byte(`{}`),
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
		Client: sdk.NewClient(nil, "kubernetes_multicluster", "", "", logpersistertest.NewTestLogPersister(t), testRegistry),
		Logger: zaptest.NewLogger(t),
	}

	dtConfig, dynamicClient := setupTestDeployTargetConfigAndDynamicClient(t)

	deploymentRes := schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}

	// Pre-create only a baseline deployment (no service).
	_, err := dynamicClient.Resource(deploymentRes).Namespace("default").Create(ctx, baselineDeployment("default"), metav1.CreateOptions{})
	require.NoError(t, err)

	plugin := &Plugin{}
	status := plugin.executeK8sMultiBaselineCleanStage(ctx, input, []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]{
		{Name: "default", Config: *dtConfig},
	})

	assert.Equal(t, sdk.StageStatusSuccess, status)

	_, err = dynamicClient.Resource(deploymentRes).Namespace("default").Get(ctx, "simple-baseline", metav1.GetOptions{})
	require.Error(t, err)
	assert.True(t, k8serrors.IsNotFound(err))
}
