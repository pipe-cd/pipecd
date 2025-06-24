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
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
	"github.com/pipe-cd/piped-plugin-sdk-go/logpersister/logpersistertest"
	"github.com/pipe-cd/piped-plugin-sdk-go/toolregistry/toolregistrytest"

	kubeConfigPkg "github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/config"
)

func TestPlugin_executeK8sRollbackStage_NoPreviousDeployment(t *testing.T) {
	t.Parallel()

	// Initialize tool registry
	testRegistry := toolregistrytest.NewTestToolRegistry(t)

	// Initialize deploy target config and dynamic client for assertions with envtest
	dtConfig, _ := setupTestDeployTargetConfigAndDynamicClient(t)

	// Read the application config from the example file
	cfg := sdk.LoadApplicationConfigForTest[kubeConfigPkg.KubernetesApplicationSpec](t, filepath.Join("testdata", "simple", "app.pipecd.yaml"), "kubernetes")

	// Prepare the input
	input := &sdk.ExecuteStageInput[kubeConfigPkg.KubernetesApplicationSpec]{
		Request: sdk.ExecuteStageRequest[kubeConfigPkg.KubernetesApplicationSpec]{
			StageName:   "K8S_ROLLBACK",
			StageConfig: []byte(``),
			RunningDeploymentSource: sdk.DeploymentSource[kubeConfigPkg.KubernetesApplicationSpec]{
				CommitHash: "", // Empty commit hash indicates no previous deployment
			},
			TargetDeploymentSource: sdk.DeploymentSource[kubeConfigPkg.KubernetesApplicationSpec]{
				ApplicationDirectory:      filepath.Join("testdata", "simple"),
				CommitHash:                "0123456789",
				ApplicationConfig:         cfg,
				ApplicationConfigFilename: "app.pipecd.yaml",
			},
			Deployment: sdk.Deployment{
				PipedID:       "piped-id",
				ApplicationID: "app-id",
			},
		},
		Client: sdk.NewClient(nil, "kubernetes", "", "", logpersistertest.NewTestLogPersister(t), testRegistry),
		Logger: zaptest.NewLogger(t),
	}

	plugin := &Plugin{}
	status := plugin.executeK8sRollbackStage(t.Context(), input, []*sdk.DeployTarget[kubeConfigPkg.KubernetesDeployTargetConfig]{
		{
			Name:   "default",
			Config: *dtConfig,
		},
	})

	assert.Equal(t, sdk.StageStatusFailure, status)
}

func TestPlugin_executeK8sRollbackStage_SuccessfulRollback(t *testing.T) {
	t.Parallel()

	// Initialize tool registry
	testRegistry := toolregistrytest.NewTestToolRegistry(t)

	// Initialize deploy target config and dynamic client for assertions with envtest
	dtConfig, dynamicClient := setupTestDeployTargetConfigAndDynamicClient(t)

	// Read the application config from the example file
	cfg := sdk.LoadApplicationConfigForTest[kubeConfigPkg.KubernetesApplicationSpec](t, filepath.Join("testdata", "simple", "app.pipecd.yaml"), "kubernetes")

	// Prepare the input
	input := &sdk.ExecuteStageInput[kubeConfigPkg.KubernetesApplicationSpec]{
		Request: sdk.ExecuteStageRequest[kubeConfigPkg.KubernetesApplicationSpec]{
			StageName:   "K8S_ROLLBACK",
			StageConfig: []byte(``),
			RunningDeploymentSource: sdk.DeploymentSource[kubeConfigPkg.KubernetesApplicationSpec]{
				ApplicationDirectory:      filepath.Join("testdata", "simple"),
				CommitHash:                "previous-hash",
				ApplicationConfig:         cfg,
				ApplicationConfigFilename: "app.pipecd.yaml",
			},
			TargetDeploymentSource: sdk.DeploymentSource[kubeConfigPkg.KubernetesApplicationSpec]{
				ApplicationDirectory:      filepath.Join("testdata", "simple"),
				CommitHash:                "0123456789",
				ApplicationConfig:         cfg,
				ApplicationConfigFilename: "app.pipecd.yaml",
			},
			Deployment: sdk.Deployment{
				PipedID:       "piped-id",
				ApplicationID: "app-id",
			},
		},
		Client: sdk.NewClient(nil, "kubernetes", "", "", logpersistertest.NewTestLogPersister(t), testRegistry),
		Logger: zaptest.NewLogger(t),
	}

	plugin := &Plugin{}
	status := plugin.executeK8sRollbackStage(t.Context(), input, []*sdk.DeployTarget[kubeConfigPkg.KubernetesDeployTargetConfig]{
		{
			Name:   "default",
			Config: *dtConfig,
		},
	})

	assert.Equal(t, sdk.StageStatusSuccess, status)

	// Verify the deployment was rolled back
	deployment, err := dynamicClient.Resource(schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}).Namespace("default").Get(t.Context(), "simple", metav1.GetOptions{})
	require.NoError(t, err)

	// Verify labels and annotations
	assert.Equal(t, "piped", deployment.GetLabels()["pipecd.dev/managed-by"])
	assert.Equal(t, "piped-id", deployment.GetLabels()["pipecd.dev/piped"])
	assert.Equal(t, "app-id", deployment.GetLabels()["pipecd.dev/application"])
	assert.Equal(t, "previous-hash", deployment.GetLabels()["pipecd.dev/commit-hash"])

	assert.Equal(t, "piped", deployment.GetAnnotations()["pipecd.dev/managed-by"])
	assert.Equal(t, "piped-id", deployment.GetAnnotations()["pipecd.dev/piped"])
	assert.Equal(t, "app-id", deployment.GetAnnotations()["pipecd.dev/application"])
	assert.Equal(t, "apps/v1", deployment.GetAnnotations()["pipecd.dev/original-api-version"])
	assert.Equal(t, "apps:Deployment::simple", deployment.GetAnnotations()["pipecd.dev/resource-key"])
	assert.Equal(t, "previous-hash", deployment.GetAnnotations()["pipecd.dev/commit-hash"])
}

func TestPlugin_executeK8sRollbackStage_WithVariantLabels(t *testing.T) {
	t.Parallel()

	// Initialize tool registry
	testRegistry := toolregistrytest.NewTestToolRegistry(t)

	// Initialize deploy target config and dynamic client for assertions with envtest
	dtConfig, dynamicClient := setupTestDeployTargetConfigAndDynamicClient(t)

	// Read the application config and modify it to include variant labels
	cfg := sdk.LoadApplicationConfigForTest[kubeConfigPkg.KubernetesApplicationSpec](t, filepath.Join("testdata", "simple", "app.pipecd.yaml"), "kubernetes")

	// Prepare the input
	input := &sdk.ExecuteStageInput[kubeConfigPkg.KubernetesApplicationSpec]{
		Request: sdk.ExecuteStageRequest[kubeConfigPkg.KubernetesApplicationSpec]{
			StageName:   "K8S_ROLLBACK",
			StageConfig: []byte(``),
			RunningDeploymentSource: sdk.DeploymentSource[kubeConfigPkg.KubernetesApplicationSpec]{
				ApplicationDirectory:      filepath.Join("testdata", "simple"),
				CommitHash:                "previous-hash",
				ApplicationConfig:         cfg,
				ApplicationConfigFilename: "app.pipecd.yaml",
			},
			TargetDeploymentSource: sdk.DeploymentSource[kubeConfigPkg.KubernetesApplicationSpec]{
				ApplicationDirectory:      filepath.Join("testdata", "simple"),
				CommitHash:                "0123456789",
				ApplicationConfig:         cfg,
				ApplicationConfigFilename: "app.pipecd.yaml",
			},
			Deployment: sdk.Deployment{
				PipedID:       "piped-id",
				ApplicationID: "app-id",
			},
		},
		Client: sdk.NewClient(nil, "kubernetes", "", "", logpersistertest.NewTestLogPersister(t), testRegistry),
		Logger: zaptest.NewLogger(t),
	}

	plugin := &Plugin{}
	status := plugin.executeK8sRollbackStage(t.Context(), input, []*sdk.DeployTarget[kubeConfigPkg.KubernetesDeployTargetConfig]{
		{
			Name:   "default",
			Config: *dtConfig,
		},
	})

	assert.Equal(t, sdk.StageStatusSuccess, status)

	// Verify the deployment was rolled back with variant labels
	deployment, err := dynamicClient.Resource(schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}).Namespace("default").Get(t.Context(), "simple", metav1.GetOptions{})
	require.NoError(t, err)

	// Verify labels and annotations
	assert.Equal(t, "piped", deployment.GetLabels()["pipecd.dev/managed-by"])
	assert.Equal(t, "piped-id", deployment.GetLabels()["pipecd.dev/piped"])
	assert.Equal(t, "app-id", deployment.GetLabels()["pipecd.dev/application"])
	assert.Equal(t, "previous-hash", deployment.GetLabels()["pipecd.dev/commit-hash"])
	assert.Equal(t, "primary", deployment.GetLabels()["pipecd.dev/variant"])

	assert.Equal(t, "piped", deployment.GetAnnotations()["pipecd.dev/managed-by"])
	assert.Equal(t, "piped-id", deployment.GetAnnotations()["pipecd.dev/piped"])
	assert.Equal(t, "app-id", deployment.GetAnnotations()["pipecd.dev/application"])
	assert.Equal(t, "apps/v1", deployment.GetAnnotations()["pipecd.dev/original-api-version"])
	assert.Equal(t, "apps:Deployment::simple", deployment.GetAnnotations()["pipecd.dev/resource-key"])
	assert.Equal(t, "previous-hash", deployment.GetAnnotations()["pipecd.dev/commit-hash"])
	assert.Equal(t, "primary", deployment.GetAnnotations()["pipecd.dev/variant"])
}

func TestPlugin_executeK8sRollbackStage_PrunesCanaryAndBaselineVariants(t *testing.T) {
	t.Parallel()

	ctx := t.Context()

	// Initialize tool registry
	testRegistry := toolregistrytest.NewTestToolRegistry(t)

	// Use dedicated testdata for this test
	runningDir := filepath.Join("testdata", "prune_rollback", "running")
	targetDir := filepath.Join("testdata", "prune_rollback", "target")
	runningCfg := sdk.LoadApplicationConfigForTest[kubeConfigPkg.KubernetesApplicationSpec](t, filepath.Join(runningDir, "app.pipecd.yaml"), "kubernetes")
	targetCfg := sdk.LoadApplicationConfigForTest[kubeConfigPkg.KubernetesApplicationSpec](t, filepath.Join(targetDir, "app.pipecd.yaml"), "kubernetes")

	// Initialize deploy target config and dynamic client for assertions with envtest
	dtConfig, dynamicClient := setupTestDeployTargetConfigAndDynamicClient(t)
	plugin := &Plugin{}

	runningDeploymentSource := sdk.DeploymentSource[kubeConfigPkg.KubernetesApplicationSpec]{
		ApplicationDirectory:      runningDir,
		CommitHash:                "running-hash",
		ApplicationConfig:         runningCfg,
		ApplicationConfigFilename: "app.pipecd.yaml",
	}
	targetDeploymentSource := sdk.DeploymentSource[kubeConfigPkg.KubernetesApplicationSpec]{
		ApplicationDirectory:      targetDir,
		CommitHash:                "target-hash",
		ApplicationConfig:         targetCfg,
		ApplicationConfigFilename: "app.pipecd.yaml",
	}
	deployment := sdk.Deployment{
		PipedID:       "piped-id",
		ApplicationID: "app-id",
	}
	deployTargets := []*sdk.DeployTarget[kubeConfigPkg.KubernetesDeployTargetConfig]{
		{
			Name:   "default",
			Config: *dtConfig,
		},
	}
	successResponse := &sdk.ExecuteStageResponse{
		Status: sdk.StageStatusSuccess,
	}

	deploymentRes := schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}
	serviceRes := schema.GroupVersionResource{Group: "", Version: "v1", Resource: "services"}
	configMapRes := schema.GroupVersionResource{Group: "", Version: "v1", Resource: "configmaps"}

	baselineOk := t.Run("baseline rollout", func(t *testing.T) {
		baselineInput := &sdk.ExecuteStageInput[kubeConfigPkg.KubernetesApplicationSpec]{
			Request: sdk.ExecuteStageRequest[kubeConfigPkg.KubernetesApplicationSpec]{
				StageName:               "K8S_BASELINE_ROLLOUT",
				StageConfig:             []byte(`{"createService": true}`),
				RunningDeploymentSource: runningDeploymentSource,
				TargetDeploymentSource:  targetDeploymentSource,
				Deployment:              deployment,
			},
			Client: sdk.NewClient(nil, "kubernetes", "", "", logpersistertest.NewTestLogPersister(t), testRegistry),
			Logger: zaptest.NewLogger(t),
		}
		status, err := plugin.ExecuteStage(ctx, nil, deployTargets, baselineInput)
		require.NoError(t, err)
		assert.Equal(t, successResponse, status)

		// baseline deployment/service should exist
		_, err = dynamicClient.Resource(deploymentRes).Namespace("default").Get(ctx, "simple-baseline", metav1.GetOptions{})
		require.NoError(t, err, "baseline deployment should exist after baseline rollout")
		_, err = dynamicClient.Resource(serviceRes).Namespace("default").Get(ctx, "simple-baseline", metav1.GetOptions{})
		require.NoError(t, err, "baseline service should exist after baseline rollout")

		// canary/primary should not exist yet
		_, err = dynamicClient.Resource(deploymentRes).Namespace("default").Get(ctx, "simple-canary", metav1.GetOptions{})
		assert.Error(t, err, "canary deployment should not exist after baseline rollout")
		_, err = dynamicClient.Resource(serviceRes).Namespace("default").Get(ctx, "simple-canary", metav1.GetOptions{})
		assert.Error(t, err, "canary service should not exist after baseline rollout")
		_, err = dynamicClient.Resource(configMapRes).Namespace("default").Get(ctx, "canary-patch-weight-config-canary", metav1.GetOptions{})
		assert.Error(t, err, "canary configmap should not exist after baseline rollout")
		_, err = dynamicClient.Resource(deploymentRes).Namespace("default").Get(ctx, "simple", metav1.GetOptions{})
		assert.Error(t, err, "primary deployment should not exist after baseline rollout")
		_, err = dynamicClient.Resource(serviceRes).Namespace("default").Get(ctx, "simple", metav1.GetOptions{})
		assert.Error(t, err, "primary service should not exist after baseline rollout")
		_, err = dynamicClient.Resource(configMapRes).Namespace("default").Get(ctx, "canary-patch-weight-config", metav1.GetOptions{})
		assert.Error(t, err, "primary configmap should not exist after baseline rollout")
	})
	require.True(t, baselineOk, "baseline rollout subtest failed, aborting")

	canaryOk := t.Run("canary rollout", func(t *testing.T) {
		canaryInput := &sdk.ExecuteStageInput[kubeConfigPkg.KubernetesApplicationSpec]{
			Request: sdk.ExecuteStageRequest[kubeConfigPkg.KubernetesApplicationSpec]{
				StageName: "K8S_CANARY_ROLLOUT",
				StageConfig: []byte(`
				{
				   "patches": [
					 {
					   "target": {"kind": "ConfigMap", "name": "canary-patch-weight-config", "documentRoot": "$.data.'weight.yaml'"},
					   "ops": [
						{"op": "yaml-replace", "path": "$.primary.weight", "value": "90"},
						{"op": "yaml-replace", "path": "$.canary.weight", "value": "10"}
					   ]
					 }
				   ]
				}`),
				RunningDeploymentSource: runningDeploymentSource,
				TargetDeploymentSource:  targetDeploymentSource,
				Deployment:              deployment,
			},
			Client: sdk.NewClient(nil, "kubernetes", "", "", logpersistertest.NewTestLogPersister(t), testRegistry),
			Logger: zaptest.NewLogger(t),
		}
		status, err := plugin.ExecuteStage(ctx, nil, deployTargets, canaryInput)
		require.NoError(t, err)
		assert.Equal(t, successResponse, status)

		// baseline deployment/service should exist
		_, err = dynamicClient.Resource(deploymentRes).Namespace("default").Get(ctx, "simple-baseline", metav1.GetOptions{})
		require.NoError(t, err, "baseline deployment should exist after baseline rollout")
		_, err = dynamicClient.Resource(serviceRes).Namespace("default").Get(ctx, "simple-baseline", metav1.GetOptions{})
		require.NoError(t, err, "baseline service should exist after baseline rollout")

		// canary deployment/configmap should exist
		_, err = dynamicClient.Resource(deploymentRes).Namespace("default").Get(ctx, "simple-canary", metav1.GetOptions{})
		require.NoError(t, err, "canary deployment should exist after canary rollout")
		_, err = dynamicClient.Resource(configMapRes).Namespace("default").Get(ctx, "canary-patch-weight-config-canary", metav1.GetOptions{})
		require.NoError(t, err, "canary configmap should exist after canary rollout")

		// primary should not exist yet
		_, err = dynamicClient.Resource(deploymentRes).Namespace("default").Get(ctx, "simple", metav1.GetOptions{})
		assert.Error(t, err, "primary deployment should not exist after canary rollout")
		_, err = dynamicClient.Resource(serviceRes).Namespace("default").Get(ctx, "simple", metav1.GetOptions{})
		assert.Error(t, err, "primary service should not exist after canary rollout")
		_, err = dynamicClient.Resource(configMapRes).Namespace("default").Get(ctx, "canary-patch-weight-config", metav1.GetOptions{})
		assert.Error(t, err, "primary configmap should not exist after canary rollout")
	})
	require.True(t, canaryOk, "canary rollout subtest failed, aborting")

	primaryOk := t.Run("primary rollout", func(t *testing.T) {
		primaryInput := &sdk.ExecuteStageInput[kubeConfigPkg.KubernetesApplicationSpec]{
			Request: sdk.ExecuteStageRequest[kubeConfigPkg.KubernetesApplicationSpec]{
				StageName:               "K8S_PRIMARY_ROLLOUT",
				StageConfig:             []byte(`{}`),
				RunningDeploymentSource: runningDeploymentSource,
				TargetDeploymentSource:  targetDeploymentSource,
				Deployment:              deployment,
			},
			Client: sdk.NewClient(nil, "kubernetes", "", "", logpersistertest.NewTestLogPersister(t), testRegistry),
			Logger: zaptest.NewLogger(t),
		}
		status, err := plugin.ExecuteStage(ctx, nil, deployTargets, primaryInput)
		require.NoError(t, err)
		assert.Equal(t, successResponse, status)

		// baseline deployment/service should exist
		_, err = dynamicClient.Resource(deploymentRes).Namespace("default").Get(ctx, "simple-baseline", metav1.GetOptions{})
		require.NoError(t, err, "baseline deployment should exist after baseline rollout")
		_, err = dynamicClient.Resource(serviceRes).Namespace("default").Get(ctx, "simple-baseline", metav1.GetOptions{})
		require.NoError(t, err, "baseline service should exist after baseline rollout")

		// canary deployment/configmap should exist
		_, err = dynamicClient.Resource(deploymentRes).Namespace("default").Get(ctx, "simple-canary", metav1.GetOptions{})
		require.NoError(t, err, "canary deployment should exist after canary rollout")
		_, err = dynamicClient.Resource(configMapRes).Namespace("default").Get(ctx, "canary-patch-weight-config-canary", metav1.GetOptions{})
		require.NoError(t, err, "canary configmap should exist after canary rollout")

		// Assert that the primary resources having the commit-hash of the target deployment
		primaryDep, err := dynamicClient.Resource(deploymentRes).Namespace("default").Get(ctx, "simple", metav1.GetOptions{})
		require.NoError(t, err, "primary deployment should exist")
		assert.Equal(t, "target-hash", primaryDep.GetAnnotations()["pipecd.dev/commit-hash"])
		primarySvc, err := dynamicClient.Resource(serviceRes).Namespace("default").Get(ctx, "simple", metav1.GetOptions{})
		require.NoError(t, err, "primary service should exist")
		assert.Equal(t, "target-hash", primarySvc.GetAnnotations()["pipecd.dev/commit-hash"])
		primaryCfg, err := dynamicClient.Resource(configMapRes).Namespace("default").Get(ctx, "canary-patch-weight-config", metav1.GetOptions{})
		require.NoError(t, err, "primary configmap should exist")
		assert.Equal(t, "target-hash", primaryCfg.GetAnnotations()["pipecd.dev/commit-hash"])

	})
	require.True(t, primaryOk, "primary rollout subtest failed, aborting")

	rollbackOk := t.Run("rollback", func(t *testing.T) {
		rollbackInput := &sdk.ExecuteStageInput[kubeConfigPkg.KubernetesApplicationSpec]{
			Request: sdk.ExecuteStageRequest[kubeConfigPkg.KubernetesApplicationSpec]{
				StageName:               "K8S_ROLLBACK",
				StageConfig:             []byte(``),
				RunningDeploymentSource: runningDeploymentSource,
				TargetDeploymentSource:  targetDeploymentSource,
				Deployment:              deployment,
			},
			Client: sdk.NewClient(nil, "kubernetes", "", "", logpersistertest.NewTestLogPersister(t), testRegistry),
			Logger: zaptest.NewLogger(t),
		}
		status, err := plugin.ExecuteStage(ctx, nil, deployTargets, rollbackInput)
		require.NoError(t, err)
		assert.Equal(t, successResponse, status)

		// Assert that canary and baseline resources are deleted
		deploymentRes := schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}
		serviceRes := schema.GroupVersionResource{Group: "", Version: "v1", Resource: "services"}
		configMapRes := schema.GroupVersionResource{Group: "", Version: "v1", Resource: "configmaps"}

		// baseline deployment/service should not exist
		_, err = dynamicClient.Resource(deploymentRes).Namespace("default").Get(ctx, "simple-baseline", metav1.GetOptions{})
		assert.Error(t, err, "baseline deployment should be deleted")
		_, err = dynamicClient.Resource(serviceRes).Namespace("default").Get(ctx, "simple-baseline", metav1.GetOptions{})
		assert.Error(t, err, "baseline service should be deleted")

		// canary deployment/configmap should not exist
		_, err = dynamicClient.Resource(deploymentRes).Namespace("default").Get(ctx, "simple-canary", metav1.GetOptions{})
		assert.Error(t, err, "canary deployment should be deleted")
		_, err = dynamicClient.Resource(configMapRes).Namespace("default").Get(ctx, "canary-patch-weight-config-canary", metav1.GetOptions{})
		assert.Error(t, err, "canary configmap should be deleted")

		// Assert that primary resources still exist
		primaryDep, err := dynamicClient.Resource(deploymentRes).Namespace("default").Get(ctx, "simple", metav1.GetOptions{})
		require.NoError(t, err, "primary deployment should exist")
		assert.Equal(t, "primary", primaryDep.GetLabels()["pipecd.dev/variant"])
		primarySvc, err := dynamicClient.Resource(serviceRes).Namespace("default").Get(ctx, "simple", metav1.GetOptions{})
		require.NoError(t, err, "primary service should exist")
		assert.Equal(t, "simple", primarySvc.GetName())
		primaryCfg, err := dynamicClient.Resource(configMapRes).Namespace("default").Get(ctx, "canary-patch-weight-config", metav1.GetOptions{})
		require.NoError(t, err, "primary configmap should exist")
		assert.Equal(t, "canary-patch-weight-config", primaryCfg.GetName())

		// Assert that the primary resources having the commit-hash of the running deployment
		assert.Equal(t, "running-hash", primaryDep.GetAnnotations()["pipecd.dev/commit-hash"])
		assert.Equal(t, "running-hash", primarySvc.GetAnnotations()["pipecd.dev/commit-hash"])
		assert.Equal(t, "running-hash", primaryCfg.GetAnnotations()["pipecd.dev/commit-hash"])
	})

	require.True(t, rollbackOk, "rollback subtest failed, aborting")
}
