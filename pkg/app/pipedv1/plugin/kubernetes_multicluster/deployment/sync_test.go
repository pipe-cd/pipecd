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
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	kubeconfig "github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes_multicluster/config"
	"github.com/pipe-cd/pipecd/pkg/plugin/logpersister/logpersistertest"
	"github.com/pipe-cd/pipecd/pkg/plugin/sdk"
	"github.com/pipe-cd/pipecd/pkg/plugin/toolregistry/toolregistrytest"
)

func TestPlugin_executeK8sMultiSyncStage(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	// read the application config from the example file
	appCfg := sdk.LoadApplicationConfigForTest[kubeconfig.KubernetesApplicationSpec](t, filepath.Join(examplesDir(), "kubernetes", "simple", "app.pipecd.yaml"), "kubernetes")

	// initialize tool registry
	testRegistry := toolregistrytest.NewTestToolRegistry(t)

	// prepare the input
	input := &sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec]{
		Request: sdk.ExecuteStageRequest[kubeconfig.KubernetesApplicationSpec]{
			StageName:               "K8S_MULTI_SYNC",
			StageConfig:             []byte(``),
			RunningDeploymentSource: sdk.DeploymentSource[kubeconfig.KubernetesApplicationSpec]{},
			TargetDeploymentSource: sdk.DeploymentSource[kubeconfig.KubernetesApplicationSpec]{
				ApplicationDirectory:      filepath.Join(examplesDir(), "kubernetes", "simple"),
				CommitHash:                "0123456789",
				ApplicationConfig:         appCfg,
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

	// initialize deploy target config and dynamic client for assertions with envtest
	dtConfig, dynamicClient := setupTestDeployTargetConfigAndDynamicClient(t)

	plugin := &Plugin{}

	status := plugin.executeK8sMultiSyncStage(ctx, input, []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]{
		{
			Name:   "default",
			Config: *dtConfig,
		},
	})

	assert.Equal(t, sdk.StageStatusSuccess, status)

	deployment, err := dynamicClient.Resource(schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}).Namespace("default").Get(context.Background(), "simple", metav1.GetOptions{})
	require.NoError(t, err)

	assert.Equal(t, "simple", deployment.GetName())
	assert.Equal(t, "simple", deployment.GetLabels()["app"])

	assert.Equal(t, "piped", deployment.GetLabels()["pipecd.dev/managed-by"])
	assert.Equal(t, "piped-id", deployment.GetLabels()["pipecd.dev/piped"])
	assert.Equal(t, "app-id", deployment.GetLabels()["pipecd.dev/application"])
	assert.Equal(t, "0123456789", deployment.GetLabels()["pipecd.dev/commit-hash"])

	assert.Equal(t, "piped", deployment.GetAnnotations()["pipecd.dev/managed-by"])
	assert.Equal(t, "piped-id", deployment.GetAnnotations()["pipecd.dev/piped"])
	assert.Equal(t, "app-id", deployment.GetAnnotations()["pipecd.dev/application"])
	assert.Equal(t, "apps/v1", deployment.GetAnnotations()["pipecd.dev/original-api-version"])
	assert.Equal(t, "apps:Deployment::simple", deployment.GetAnnotations()["pipecd.dev/resource-key"]) // This assertion differs from the non-plugin-arched piped's Kubernetes platform provider, but we decided to change this behavior.
	assert.Equal(t, "0123456789", deployment.GetAnnotations()["pipecd.dev/commit-hash"])
}

func TestPlugin_executeK8sMultiSyncStage_withInputNamespace(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	// read the application config from the example file
	appCfg := sdk.LoadApplicationConfigForTest[kubeconfig.KubernetesApplicationSpec](t, filepath.Join(examplesDir(), "kubernetes", "simple", "app.pipecd.yaml"), "kubernetes")

	// override the autoCreateNamespace and namespace
	appCfg.Spec.Input.AutoCreateNamespace = true
	appCfg.Spec.Input.Namespace = "test-namespace"

	// initialize tool registry
	testRegistry := toolregistrytest.NewTestToolRegistry(t)

	// prepare the input
	input := &sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec]{
		Request: sdk.ExecuteStageRequest[kubeconfig.KubernetesApplicationSpec]{
			StageName:               "K8S_MULTI_SYNC",
			StageConfig:             []byte(``),
			RunningDeploymentSource: sdk.DeploymentSource[kubeconfig.KubernetesApplicationSpec]{},
			TargetDeploymentSource: sdk.DeploymentSource[kubeconfig.KubernetesApplicationSpec]{
				ApplicationDirectory:      filepath.Join(examplesDir(), "kubernetes", "simple"),
				CommitHash:                "0123456789",
				ApplicationConfig:         appCfg,
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

	// initialize deploy target config and dynamic client for assertions with envtest
	dtConfig, dynamicClient := setupTestDeployTargetConfigAndDynamicClient(t)

	plugin := &Plugin{}

	status := plugin.executeK8sMultiSyncStage(ctx, input, []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]{
		{
			Name:   "default",
			Config: *dtConfig,
		},
	})

	assert.Equal(t, sdk.StageStatusSuccess, status)

	deployment, err := dynamicClient.Resource(schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}).Namespace("test-namespace").Get(context.Background(), "simple", metav1.GetOptions{})
	require.NoError(t, err)

	assert.Equal(t, "piped", deployment.GetLabels()["pipecd.dev/managed-by"])
	assert.Equal(t, "piped-id", deployment.GetLabels()["pipecd.dev/piped"])
	assert.Equal(t, "app-id", deployment.GetLabels()["pipecd.dev/application"])
	assert.Equal(t, "0123456789", deployment.GetLabels()["pipecd.dev/commit-hash"])

	assert.Equal(t, "simple", deployment.GetName())
	assert.Equal(t, "simple", deployment.GetLabels()["app"])
	assert.Equal(t, "piped", deployment.GetAnnotations()["pipecd.dev/managed-by"])
	assert.Equal(t, "piped-id", deployment.GetAnnotations()["pipecd.dev/piped"])
	assert.Equal(t, "app-id", deployment.GetAnnotations()["pipecd.dev/application"])
	assert.Equal(t, "apps/v1", deployment.GetAnnotations()["pipecd.dev/original-api-version"])
	assert.Equal(t, "apps:Deployment:test-namespace:simple", deployment.GetAnnotations()["pipecd.dev/resource-key"]) // This assertion differs from the non-plugin-arched piped's Kubernetes platform provider, but we decided to change this behavior.
	assert.Equal(t, "0123456789", deployment.GetAnnotations()["pipecd.dev/commit-hash"])
}

func TestPlugin_executeK8sMultiSyncStage_withPrune(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	// initialize tool registry
	testRegistry := toolregistrytest.NewTestToolRegistry(t)

	// initialize deploy target config and dynamic client for assertions with envtest
	dtConfig, dynamicClient := setupTestDeployTargetConfigAndDynamicClient(t)

	running := filepath.Join("./", "testdata", "prune", "running")

	// read the running application config from the testdata file
	runningCfg := sdk.LoadApplicationConfigForTest[kubeconfig.KubernetesApplicationSpec](t, filepath.Join(running, "app.pipecd.yaml"), "kubernetes")

	ok := t.Run("prepare", func(t *testing.T) {
		// prepare the input to ensure the running deployment exists
		runningInput := &sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec]{
			Request: sdk.ExecuteStageRequest[kubeconfig.KubernetesApplicationSpec]{
				StageName:               "K8S_MULTI_SYNC",
				StageConfig:             []byte(``),
				RunningDeploymentSource: sdk.DeploymentSource[kubeconfig.KubernetesApplicationSpec]{},
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
			Client: sdk.NewClient(nil, "kubernetes", "app-id", "stage-id", logpersistertest.NewTestLogPersister(t), testRegistry),
			Logger: zaptest.NewLogger(t),
		}

		plugin := &Plugin{}
		status := plugin.executeK8sMultiSyncStage(ctx, runningInput, []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]{
			{
				Name:   "default",
				Config: *dtConfig,
			},
		})
		require.Equal(t, sdk.StageStatusSuccess, status)

		service, err := dynamicClient.Resource(schema.GroupVersionResource{Group: "", Version: "v1", Resource: "services"}).Namespace("default").Get(context.Background(), "simple", metav1.GetOptions{})
		require.NoError(t, err)

		require.Equal(t, "piped", service.GetLabels()["pipecd.dev/managed-by"])
		require.Equal(t, "piped-id", service.GetLabels()["pipecd.dev/piped"])
		require.Equal(t, "app-id", service.GetLabels()["pipecd.dev/application"])
		require.Equal(t, "0123456789", service.GetLabels()["pipecd.dev/commit-hash"])

		require.Equal(t, "simple", service.GetName())
		require.Equal(t, "piped", service.GetAnnotations()["pipecd.dev/managed-by"])
		require.Equal(t, "piped-id", service.GetAnnotations()["pipecd.dev/piped"])
		require.Equal(t, "app-id", service.GetAnnotations()["pipecd.dev/application"])
		require.Equal(t, "v1", service.GetAnnotations()["pipecd.dev/original-api-version"])
		require.Equal(t, ":Service::simple", service.GetAnnotations()["pipecd.dev/resource-key"]) // This assertion differs from the non-plugin-arched piped's Kubernetes platform provider, but we decided to change this behavior.
		require.Equal(t, "0123456789", service.GetAnnotations()["pipecd.dev/commit-hash"])
	})
	require.Truef(t, ok, "expected prepare to succeed")

	t.Run("run with prune", func(t *testing.T) {
		target := filepath.Join("./", "testdata", "prune", "target")

		// read the running application config from the testdata file
		targetCfg := sdk.LoadApplicationConfigForTest[kubeconfig.KubernetesApplicationSpec](t, filepath.Join(target, "app.pipecd.yaml"), "kubernetes")

		// prepare the input to ensure the running deployment exists
		targetInput := &sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec]{
			Request: sdk.ExecuteStageRequest[kubeconfig.KubernetesApplicationSpec]{
				StageName:   "K8S_MULTI_SYNC",
				StageConfig: []byte(``),
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
			Client: sdk.NewClient(nil, "kubernetes", "app-id", "stage-id", logpersistertest.NewTestLogPersister(t), testRegistry),
			Logger: zaptest.NewLogger(t),
		}

		plugin := &Plugin{}
		status := plugin.executeK8sMultiSyncStage(ctx, targetInput, []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]{
			{
				Name:   "default",
				Config: *dtConfig,
			},
		})
		assert.Equal(t, sdk.StageStatusSuccess, status)

		_, err := dynamicClient.Resource(schema.GroupVersionResource{Group: "", Version: "v1", Resource: "services"}).Namespace("default").Get(context.Background(), "simple", metav1.GetOptions{})
		require.Error(t, err)
		require.Truef(t, apierrors.IsNotFound(err), "expected error to be NotFound, but got %v", err)
	})
}

func TestPlugin_executeK8sMultiSyncStage_withPrune_changesNamespace(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	// initialize tool registry
	testRegistry := toolregistrytest.NewTestToolRegistry(t)

	// initialize deploy target config and dynamic client for assertions with envtest
	dtConfig, dynamicClient := setupTestDeployTargetConfigAndDynamicClient(t)

	running := filepath.Join("./", "testdata", "prune_with_change_namespace", "running")

	// read the running application config from the example file
	runningCfg := sdk.LoadApplicationConfigForTest[kubeconfig.KubernetesApplicationSpec](t, filepath.Join(running, "app.pipecd.yaml"), "kubernetes")

	ok := t.Run("prepare", func(t *testing.T) {
		// prepare the input to ensure the running deployment exists
		runningInput := &sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec]{
			Request: sdk.ExecuteStageRequest[kubeconfig.KubernetesApplicationSpec]{
				StageName:               "K8S_MULTI_SYNC",
				StageConfig:             []byte(``),
				RunningDeploymentSource: sdk.DeploymentSource[kubeconfig.KubernetesApplicationSpec]{},
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
			Client: sdk.NewClient(nil, "kubernetes", "app-id", "stage-id", logpersistertest.NewTestLogPersister(t), testRegistry),
			Logger: zaptest.NewLogger(t),
		}

		plugin := &Plugin{}
		status := plugin.executeK8sMultiSyncStage(ctx, runningInput, []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]{
			{
				Name:   "default",
				Config: *dtConfig,
			},
		})
		require.Equal(t, sdk.StageStatusSuccess, status)

		service, err := dynamicClient.Resource(schema.GroupVersionResource{Group: "", Version: "v1", Resource: "services"}).Namespace("test-1").Get(context.Background(), "simple", metav1.GetOptions{})
		require.NoError(t, err)

		require.Equal(t, "piped", service.GetLabels()["pipecd.dev/managed-by"])
		require.Equal(t, "piped-id", service.GetLabels()["pipecd.dev/piped"])
		require.Equal(t, "app-id", service.GetLabels()["pipecd.dev/application"])
		require.Equal(t, "0123456789", service.GetLabels()["pipecd.dev/commit-hash"])

		require.Equal(t, "simple", service.GetName())
		require.Equal(t, "piped", service.GetAnnotations()["pipecd.dev/managed-by"])
		require.Equal(t, "piped-id", service.GetAnnotations()["pipecd.dev/piped"])
		require.Equal(t, "app-id", service.GetAnnotations()["pipecd.dev/application"])
		require.Equal(t, "v1", service.GetAnnotations()["pipecd.dev/original-api-version"])
		require.Equal(t, "0123456789", service.GetAnnotations()["pipecd.dev/commit-hash"])
		require.Equal(t, ":Service:test-1:simple", service.GetAnnotations()["pipecd.dev/resource-key"])
	})
	require.Truef(t, ok, "expected prepare to succeed")

	t.Run("run with prune", func(t *testing.T) {
		target := filepath.Join("./", "testdata", "prune_with_change_namespace", "target")

		// read the running application config from the example file
		targetCfg := sdk.LoadApplicationConfigForTest[kubeconfig.KubernetesApplicationSpec](t, filepath.Join(target, "app.pipecd.yaml"), "kubernetes")

		// prepare the input to ensure the running deployment exists
		targetInput := &sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec]{
			Request: sdk.ExecuteStageRequest[kubeconfig.KubernetesApplicationSpec]{
				StageName:   "K8S_MULTI_SYNC",
				StageConfig: []byte(``),
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
			Client: sdk.NewClient(nil, "kubernetes", "app-id", "stage-id", logpersistertest.NewTestLogPersister(t), testRegistry),
			Logger: zaptest.NewLogger(t),
		}

		plugin := &Plugin{}
		status := plugin.executeK8sMultiSyncStage(ctx, targetInput, []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]{
			{
				Name:   "default",
				Config: *dtConfig,
			},
		})
		require.Equal(t, sdk.StageStatusSuccess, status)

		// The service should be removed from the previous namespace
		_, err := dynamicClient.Resource(schema.GroupVersionResource{Group: "", Version: "v1", Resource: "services"}).Namespace("test-1").Get(context.Background(), "simple", metav1.GetOptions{})
		require.Error(t, err)
		require.Truef(t, apierrors.IsNotFound(err), "expected error to be NotFound, but got %v", err)

		// The service should be created in the new namespace
		service, err := dynamicClient.Resource(schema.GroupVersionResource{Group: "", Version: "v1", Resource: "services"}).Namespace("test-2").Get(context.Background(), "simple", metav1.GetOptions{})
		require.NoError(t, err)

		require.Equal(t, "piped", service.GetLabels()["pipecd.dev/managed-by"])
		require.Equal(t, "piped-id", service.GetLabels()["pipecd.dev/piped"])
		require.Equal(t, "app-id", service.GetLabels()["pipecd.dev/application"])
		require.Equal(t, "0012345678", service.GetLabels()["pipecd.dev/commit-hash"])

		require.Equal(t, "simple", service.GetName())
		require.Equal(t, "piped", service.GetAnnotations()["pipecd.dev/managed-by"])
		require.Equal(t, "piped-id", service.GetAnnotations()["pipecd.dev/piped"])
		require.Equal(t, "app-id", service.GetAnnotations()["pipecd.dev/application"])
		require.Equal(t, "v1", service.GetAnnotations()["pipecd.dev/original-api-version"])
		require.Equal(t, "0012345678", service.GetAnnotations()["pipecd.dev/commit-hash"])
		require.Equal(t, ":Service:test-2:simple", service.GetAnnotations()["pipecd.dev/resource-key"])
	})
}

func TestPlugin_executeK8sMultiSyncStage_withPrune_clusterScoped(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	// initialize tool registry
	testRegistry := toolregistrytest.NewTestToolRegistry(t)

	// initialize deploy target config and dynamic client for assertions with envtest
	dtConfig, dynamicClient := setupTestDeployTargetConfigAndDynamicClient(t)

	// prepare the custom resource definition
	prepare := filepath.Join("./", "testdata", "prune_cluster_scoped_resource", "prepare")

	prepareCfg := sdk.LoadApplicationConfigForTest[kubeconfig.KubernetesApplicationSpec](t, filepath.Join(prepare, "app.pipecd.yaml"), "kubernetes")

	ok := t.Run("prepare crd", func(t *testing.T) {
		// prepare the input to ensure the running deployment exists
		prepareInput := &sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec]{
			Request: sdk.ExecuteStageRequest[kubeconfig.KubernetesApplicationSpec]{
				StageName:               "K8S_MULTI_SYNC",
				StageConfig:             []byte(``),
				RunningDeploymentSource: sdk.DeploymentSource[kubeconfig.KubernetesApplicationSpec]{},
				TargetDeploymentSource: sdk.DeploymentSource[kubeconfig.KubernetesApplicationSpec]{
					ApplicationDirectory:      prepare,
					CommitHash:                "0123456789",
					ApplicationConfig:         prepareCfg,
					ApplicationConfigFilename: "app.pipecd.yaml",
				},
				Deployment: sdk.Deployment{
					PipedID:       "piped-id",
					ApplicationID: "prepare-app-id",
				},
			},
			Client: sdk.NewClient(nil, "kubernetes", "prepare-app-id", "stage-id", logpersistertest.NewTestLogPersister(t), testRegistry),
			Logger: zaptest.NewLogger(t),
		}

		plugin := &Plugin{}
		status := plugin.executeK8sMultiSyncStage(ctx, prepareInput, []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]{
			{
				Name:   "default",
				Config: *dtConfig,
			},
		})
		require.Equal(t, sdk.StageStatusSuccess, status)
	})
	require.Truef(t, ok, "expected prepare to succeed")

	// prepare the running resources
	running := filepath.Join("./", "testdata", "prune_cluster_scoped_resource", "running")

	// read the running application config from the example file
	runningCfg := sdk.LoadApplicationConfigForTest[kubeconfig.KubernetesApplicationSpec](t, filepath.Join(running, "app.pipecd.yaml"), "kubernetes")

	ok = t.Run("prepare running", func(t *testing.T) {
		// prepare the input to ensure the running deployment exists
		runningInput := &sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec]{
			Request: sdk.ExecuteStageRequest[kubeconfig.KubernetesApplicationSpec]{
				StageName:               "K8S_MULTI_SYNC",
				StageConfig:             []byte(``),
				RunningDeploymentSource: sdk.DeploymentSource[kubeconfig.KubernetesApplicationSpec]{},
				TargetDeploymentSource: sdk.DeploymentSource[kubeconfig.KubernetesApplicationSpec]{
					ApplicationDirectory:      running,
					CommitHash:                "0123456789",
					ApplicationConfig:         runningCfg,
					ApplicationConfigFilename: "app.pipecd.yaml",
				},
			},
			Client: sdk.NewClient(nil, "kubernetes", "app-id", "stage-id", logpersistertest.NewTestLogPersister(t), testRegistry),
			Logger: zaptest.NewLogger(t),
		}

		plugin := &Plugin{}
		status := plugin.executeK8sMultiSyncStage(ctx, runningInput, []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]{
			{
				Name:   "default",
				Config: *dtConfig,
			},
		})
		require.Equal(t, sdk.StageStatusSuccess, status)

		// The my-new-cron-object/my-new-cron-object-2/my-new-cron-object-v1beta1 should be created
		_, err := dynamicClient.Resource(schema.GroupVersionResource{Group: "stable.example.com", Version: "v1", Resource: "crontabs"}).Get(context.Background(), "my-new-cron-object", metav1.GetOptions{})
		require.NoError(t, err)
		_, err = dynamicClient.Resource(schema.GroupVersionResource{Group: "stable.example.com", Version: "v1", Resource: "crontabs"}).Get(context.Background(), "my-new-cron-object-2", metav1.GetOptions{})
		require.NoError(t, err)
		_, err = dynamicClient.Resource(schema.GroupVersionResource{Group: "stable.example.com", Version: "v1", Resource: "crontabs"}).Get(context.Background(), "my-new-cron-object-v1beta1", metav1.GetOptions{})
		require.NoError(t, err)
	})
	require.Truef(t, ok, "expected prepare to succeed")

	t.Run("sync", func(t *testing.T) {
		// sync the target resources and assert the prune behavior
		target := filepath.Join("./", "testdata", "prune_cluster_scoped_resource", "target")

		// read the running application config from the example file
		targetCfg := sdk.LoadApplicationConfigForTest[kubeconfig.KubernetesApplicationSpec](t, filepath.Join(target, "app.pipecd.yaml"), "kubernetes")

		// prepare the input to ensure the running deployment exists
		targetInput := &sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec]{
			Request: sdk.ExecuteStageRequest[kubeconfig.KubernetesApplicationSpec]{
				StageName:   "K8S_MULTI_SYNC",
				StageConfig: []byte(``),
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
			},
			Client: sdk.NewClient(nil, "kubernetes", "app-id", "stage-id", logpersistertest.NewTestLogPersister(t), testRegistry),
			Logger: zaptest.NewLogger(t),
		}

		plugin := &Plugin{}
		status := plugin.executeK8sMultiSyncStage(ctx, targetInput, []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]{
			{
				Name:   "default",
				Config: *dtConfig,
			},
		})
		require.Equal(t, sdk.StageStatusSuccess, status)

		// The my-new-cron-object should not be removed
		_, err := dynamicClient.Resource(schema.GroupVersionResource{Group: "stable.example.com", Version: "v1", Resource: "crontabs"}).Get(context.Background(), "my-new-cron-object", metav1.GetOptions{})
		require.NoError(t, err)
		// The my-new-cron-object-2 should be removed
		_, err = dynamicClient.Resource(schema.GroupVersionResource{Group: "stable.example.com", Version: "v1", Resource: "crontabs"}).Get(context.Background(), "my-new-cron-object-2", metav1.GetOptions{})
		require.Error(t, err)
		require.Truef(t, apierrors.IsNotFound(err), "expected error to be NotFound, but got %v", err)
		// The my-new-cron-object-v1beta1 should be removed
		_, err = dynamicClient.Resource(schema.GroupVersionResource{Group: "stable.example.com", Version: "v1", Resource: "crontabs"}).Get(context.Background(), "my-new-cron-object-v1beta1", metav1.GetOptions{})
		require.Error(t, err)
		require.Truef(t, apierrors.IsNotFound(err), "expected error to be NotFound, but got %v", err)
	})
}

func TestPlugin_executeK8sMultiSyncStage_multiCluster(t *testing.T) {
	t.Parallel()

	// read the application config from the example file
	cfg := sdk.LoadApplicationConfigForTest[kubeconfig.KubernetesApplicationSpec](t, filepath.Join(examplesDir(), "kubernetes", "simple", "app.pipecd.yaml"), "kubernetes")

	// initialize tool registry
	testRegistry := toolregistrytest.NewTestToolRegistry(t)

	// prepare the input
	input := &sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec]{
		Request: sdk.ExecuteStageRequest[kubeconfig.KubernetesApplicationSpec]{
			StageName: "K8S_MULTI_SYNC",
			Deployment: sdk.Deployment{
				PipedID:       "piped-id",
				ApplicationID: "app-id",
			},
			StageConfig:             []byte(``),
			RunningDeploymentSource: sdk.DeploymentSource[kubeconfig.KubernetesApplicationSpec]{},
			TargetDeploymentSource: sdk.DeploymentSource[kubeconfig.KubernetesApplicationSpec]{
				ApplicationDirectory:      filepath.Join(examplesDir(), "kubernetes", "simple"),
				CommitHash:                "0123456789",
				ApplicationConfig:         cfg,
				ApplicationConfigFilename: "app.pipecd.yaml",
			},
		},
		Client: sdk.NewClient(nil, "kubernetes", "app-id", "stage-id", logpersistertest.NewTestLogPersister(t), testRegistry),
		Logger: zaptest.NewLogger(t),
	}

	// prepare the first cluster
	cluster1 := setupCluster(t, "cluster1")
	cluster2 := setupCluster(t, "cluster2")

	dts := []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]{
		{
			Name:   "cluster1",
			Config: *cluster1.dtc,
		},
		{
			Name:   "cluster2",
			Config: *cluster2.dtc,
		},
	}

	plugin := &Plugin{}
	status := plugin.executeK8sMultiSyncStage(t.Context(), input, dts)

	require.Equal(t, sdk.StageStatusSuccess, status)

	for _, cluster := range []*cluster{cluster1, cluster2} {
		deployment, err := cluster.cli.Resource(schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}).Namespace("default").Get(context.Background(), "simple", metav1.GetOptions{})
		require.NoError(t, err)
		assert.Equal(t, "simple", deployment.GetName())
		assert.Equal(t, "simple", deployment.GetLabels()["app"])

		assert.Equal(t, "piped", deployment.GetLabels()["pipecd.dev/managed-by"])
		assert.Equal(t, "piped-id", deployment.GetLabels()["pipecd.dev/piped"])
		assert.Equal(t, "app-id", deployment.GetLabels()["pipecd.dev/application"])
		assert.Equal(t, "0123456789", deployment.GetLabels()["pipecd.dev/commit-hash"])

		assert.Equal(t, "piped", deployment.GetAnnotations()["pipecd.dev/managed-by"])
		assert.Equal(t, "piped-id", deployment.GetAnnotations()["pipecd.dev/piped"])
		assert.Equal(t, "app-id", deployment.GetAnnotations()["pipecd.dev/application"])
		assert.Equal(t, "apps/v1", deployment.GetAnnotations()["pipecd.dev/original-api-version"])
		assert.Equal(t, "apps:Deployment::simple", deployment.GetAnnotations()["pipecd.dev/resource-key"]) // This assertion differs from the non-plugin-arched piped's Kubernetes platform provider, but we decided to change this behavior.
		assert.Equal(t, "0123456789", deployment.GetAnnotations()["pipecd.dev/commit-hash"])
	}
}

func TestPlugin_executeK8sMultiSyncStage_multiCluster_templateNone(t *testing.T) {
	t.Parallel()

	target := filepath.Join("./", "testdata", "multicluster_template_none", "target")
	// read the application config from the example file
	cfg := sdk.LoadApplicationConfigForTest[kubeconfig.KubernetesApplicationSpec](t, filepath.Join(target, "app.pipecd.yaml"), "kubernetes")

	// initialize tool registry
	testRegistry := toolregistrytest.NewTestToolRegistry(t)

	// prepare the input
	input := &sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec]{
		Request: sdk.ExecuteStageRequest[kubeconfig.KubernetesApplicationSpec]{
			StageName: "K8S_MULTI_SYNC",
			Deployment: sdk.Deployment{
				PipedID:       "piped-id",
				ApplicationID: "app-id",
			},
			StageConfig:             []byte(``),
			RunningDeploymentSource: sdk.DeploymentSource[kubeconfig.KubernetesApplicationSpec]{},
			TargetDeploymentSource: sdk.DeploymentSource[kubeconfig.KubernetesApplicationSpec]{
				ApplicationDirectory:      target,
				CommitHash:                "0123456789",
				ApplicationConfig:         cfg,
				ApplicationConfigFilename: "app.pipecd.yaml",
			},
		},
		Client: sdk.NewClient(nil, "kubernetes", "app-id", "stage-id", logpersistertest.NewTestLogPersister(t), testRegistry),
		Logger: zaptest.NewLogger(t),
	}

	// prepare the first cluster
	cluster1 := setupCluster(t, "cluster1")
	cluster2 := setupCluster(t, "cluster2")

	dts := []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]{
		{
			Name:   "cluster1",
			Config: *cluster1.dtc,
		},
		{
			Name:   "cluster2",
			Config: *cluster2.dtc,
		},
	}

	plugin := &Plugin{}
	status := plugin.executeK8sMultiSyncStage(t.Context(), input, dts)

	require.Equal(t, sdk.StageStatusSuccess, status)

	for _, cluster := range []*cluster{cluster1, cluster2} {
		deployment, err := cluster.cli.Resource(schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}).Namespace("default").Get(context.Background(), fmt.Sprintf("simple-%s", cluster.name), metav1.GetOptions{})
		require.NoError(t, err)
		assert.Equal(t, fmt.Sprintf("simple-%s", cluster.name), deployment.GetName())
		assert.Equal(t, fmt.Sprintf("simple-%s", cluster.name), deployment.GetLabels()["app"])

		assert.Equal(t, "piped", deployment.GetLabels()["pipecd.dev/managed-by"])
		assert.Equal(t, "piped-id", deployment.GetLabels()["pipecd.dev/piped"])
		assert.Equal(t, "app-id", deployment.GetLabels()["pipecd.dev/application"])
		assert.Equal(t, "0123456789", deployment.GetLabels()["pipecd.dev/commit-hash"])

		assert.Equal(t, "piped", deployment.GetAnnotations()["pipecd.dev/managed-by"])
		assert.Equal(t, "piped-id", deployment.GetAnnotations()["pipecd.dev/piped"])
		assert.Equal(t, "app-id", deployment.GetAnnotations()["pipecd.dev/application"])
		assert.Equal(t, "apps/v1", deployment.GetAnnotations()["pipecd.dev/original-api-version"])
		assert.Equal(t, "apps:Deployment::"+fmt.Sprintf("simple-%s", cluster.name), deployment.GetAnnotations()["pipecd.dev/resource-key"]) // This assertion differs from the non-plugin-arched piped's Kubernetes platform provider, but we decided to change this behavior.
		assert.Equal(t, "0123456789", deployment.GetAnnotations()["pipecd.dev/commit-hash"])
	}
}

func TestPlugin_executeK8sMultiSyncStage_multiCluster_failed_one_of_the_sync(t *testing.T) {
	t.Parallel()

	target := filepath.Join("./", "testdata", "multicluster_failed_one_of_the_sync", "target")

	cfg := sdk.LoadApplicationConfigForTest[kubeconfig.KubernetesApplicationSpec](t, filepath.Join(target, "app.pipecd.yaml"), "kubernetes")

	// initialize tool registry
	testRegistry := toolregistrytest.NewTestToolRegistry(t)

	// prepare the input
	input := &sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec]{
		Request: sdk.ExecuteStageRequest[kubeconfig.KubernetesApplicationSpec]{
			StageName: "K8S_MULTI_SYNC",
			Deployment: sdk.Deployment{
				PipedID:       "piped-id",
				ApplicationID: "app-id",
			},
			StageConfig:             []byte(``),
			RunningDeploymentSource: sdk.DeploymentSource[kubeconfig.KubernetesApplicationSpec]{},
			TargetDeploymentSource: sdk.DeploymentSource[kubeconfig.KubernetesApplicationSpec]{
				ApplicationDirectory:      target,
				CommitHash:                "0123456789",
				ApplicationConfig:         cfg,
				ApplicationConfigFilename: "app.pipecd.yaml",
			},
		},
		Client: sdk.NewClient(nil, "kubernetes", "app-id", "stage-id", logpersistertest.NewTestLogPersister(t), testRegistry),
		Logger: zaptest.NewLogger(t),
	}

	// prepare the first cluster
	cluster1 := setupCluster(t, "cluster1")
	cluster2 := setupCluster(t, "cluster2")

	dts := []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]{
		{
			Name:   "cluster1",
			Config: *cluster1.dtc,
		},
		{
			Name:   "cluster2",
			Config: *cluster2.dtc,
		},
	}

	plugin := &Plugin{}
	status := plugin.executeK8sMultiSyncStage(t.Context(), input, dts)

	require.Equal(t, sdk.StageStatusFailure, status)
}
