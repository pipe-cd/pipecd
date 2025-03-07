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
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/yaml"

	kubeConfigPkg "github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/config"
	config "github.com/pipe-cd/pipecd/pkg/configv1"
	"github.com/pipe-cd/pipecd/pkg/plugin/logpersister/logpersistertest"
	"github.com/pipe-cd/pipecd/pkg/plugin/sdk"
	"github.com/pipe-cd/pipecd/pkg/plugin/toolregistry/toolregistrytest"
)

func setupTestDeployTargetConfig(t *testing.T, kubeCfg *rest.Config) *kubeConfigPkg.KubernetesDeployTargetConfig {
	t.Helper()

	kubeconfig, err := kubeconfigFromRestConfig(kubeCfg)
	require.NoError(t, err)

	workDir := t.TempDir()
	kubeconfigPath := path.Join(workDir, "kubeconfig")
	err = os.WriteFile(kubeconfigPath, []byte(kubeconfig), 0755)
	require.NoError(t, err)

	return &kubeConfigPkg.KubernetesDeployTargetConfig{
		KubeConfigPath: kubeconfigPath,
	}
}

func setupTestDeployTargetConfigAndDynamicClient(t *testing.T) (*kubeConfigPkg.KubernetesDeployTargetConfig, dynamic.Interface) {
	t.Helper()

	kubeCfg := setupEnvTest(t)
	deployTargetCfg := setupTestDeployTargetConfig(t, kubeCfg)

	dynamicClient, err := dynamic.NewForConfig(kubeCfg)
	require.NoError(t, err)

	return deployTargetCfg, dynamicClient
}

func TestPlugin_executeK8sSyncStage(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	// read the application config from the example file
	cfg, err := os.ReadFile(filepath.Join(examplesDir(), "kubernetes", "simple", "app.pipecd.yaml"))
	require.NoError(t, err)

	// initialize tool registry
	testRegistry := toolregistrytest.NewTestToolRegistry(t)

	// prepare the input
	input := &sdk.ExecuteStageInput{
		Request: sdk.ExecuteStageRequest{
			StageName:               "K8S_SYNC",
			StageConfig:             []byte(``),
			RunningDeploymentSource: sdk.DeploymentSource{},
			TargetDeploymentSource: sdk.DeploymentSource{
				ApplicationDirectory:      filepath.Join(examplesDir(), "kubernetes", "simple"),
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

	// initialize deploy target config and dynamic client for assertions with envtest
	dtConfig, dynamicClient := setupTestDeployTargetConfigAndDynamicClient(t)

	plugin := &Plugin{}

	status := plugin.executeK8sSyncStage(ctx, input, []*sdk.DeployTarget[kubeConfigPkg.KubernetesDeployTargetConfig]{
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

func TestPlugin_executeK8sSyncStage_withInputNamespace(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	// read the application config from the example file
	cfg, err := os.ReadFile(filepath.Join(examplesDir(), "kubernetes", "simple", "app.pipecd.yaml"))
	require.NoError(t, err)

	// decode and override the autoCreateNamespace and namespace
	spec, err := config.DecodeYAML[*kubeConfigPkg.KubernetesApplicationSpec](cfg)
	require.NoError(t, err)
	spec.Spec.Input.AutoCreateNamespace = true
	spec.Spec.Input.Namespace = "test-namespace"
	cfg, err = yaml.Marshal(spec)
	require.NoError(t, err)

	// initialize tool registry
	testRegistry := toolregistrytest.NewTestToolRegistry(t)

	// prepare the input
	input := &sdk.ExecuteStageInput{
		Request: sdk.ExecuteStageRequest{
			StageName:               "K8S_SYNC",
			StageConfig:             []byte(``),
			RunningDeploymentSource: sdk.DeploymentSource{},
			TargetDeploymentSource: sdk.DeploymentSource{
				ApplicationDirectory:      filepath.Join(examplesDir(), "kubernetes", "simple"),
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

	// initialize deploy target config and dynamic client for assertions with envtest
	dtConfig, dynamicClient := setupTestDeployTargetConfigAndDynamicClient(t)

	plugin := &Plugin{}

	status := plugin.executeK8sSyncStage(ctx, input, []*sdk.DeployTarget[kubeConfigPkg.KubernetesDeployTargetConfig]{
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

func TestPlugin_executeK8sSyncStage_withPrune(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	// initialize tool registry
	testRegistry := toolregistrytest.NewTestToolRegistry(t)

	// initialize deploy target config and dynamic client for assertions with envtest
	dtConfig, dynamicClient := setupTestDeployTargetConfigAndDynamicClient(t)

	running := filepath.Join("./", "testdata", "prune", "running")

	// read the running application config from the testdata file
	runningCfg, err := os.ReadFile(filepath.Join(running, "app.pipecd.yaml"))
	require.NoError(t, err)

	ok := t.Run("prepare", func(t *testing.T) {
		// prepare the input to ensure the running deployment exists
		runningInput := &sdk.ExecuteStageInput{
			Request: sdk.ExecuteStageRequest{
				StageName:               "K8S_SYNC",
				StageConfig:             []byte(``),
				RunningDeploymentSource: sdk.DeploymentSource{},
				TargetDeploymentSource: sdk.DeploymentSource{
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
		status := plugin.executeK8sSyncStage(ctx, runningInput, []*sdk.DeployTarget[kubeConfigPkg.KubernetesDeployTargetConfig]{
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
		targetCfg, err := os.ReadFile(filepath.Join(target, "app.pipecd.yaml"))
		require.NoError(t, err)

		// prepare the input to ensure the running deployment exists
		targetInput := &sdk.ExecuteStageInput{
			Request: sdk.ExecuteStageRequest{
				StageName:   "K8S_SYNC",
				StageConfig: []byte(``),
				RunningDeploymentSource: sdk.DeploymentSource{
					ApplicationDirectory:      running,
					CommitHash:                "0123456789",
					ApplicationConfig:         runningCfg,
					ApplicationConfigFilename: "app.pipecd.yaml",
				},
				TargetDeploymentSource: sdk.DeploymentSource{
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
		status := plugin.executeK8sSyncStage(ctx, targetInput, []*sdk.DeployTarget[kubeConfigPkg.KubernetesDeployTargetConfig]{
			{
				Name:   "default",
				Config: *dtConfig,
			},
		})
		assert.Equal(t, sdk.StageStatusSuccess, status)

		_, err = dynamicClient.Resource(schema.GroupVersionResource{Group: "", Version: "v1", Resource: "services"}).Namespace("default").Get(context.Background(), "simple", metav1.GetOptions{})
		require.Error(t, err)
		require.Truef(t, apierrors.IsNotFound(err), "expected error to be NotFound, but got %v", err)
	})
}

func TestPlugin_executeK8sSyncStage_withPrune_changesNamespace(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	// initialize tool registry
	testRegistry := toolregistrytest.NewTestToolRegistry(t)

	// initialize deploy target config and dynamic client for assertions with envtest
	dtConfig, dynamicClient := setupTestDeployTargetConfigAndDynamicClient(t)

	running := filepath.Join("./", "testdata", "prune_with_change_namespace", "running")

	// read the running application config from the example file
	runningCfg, err := os.ReadFile(filepath.Join(running, "app.pipecd.yaml"))
	require.NoError(t, err)

	ok := t.Run("prepare", func(t *testing.T) {
		// prepare the input to ensure the running deployment exists
		runningInput := &sdk.ExecuteStageInput{
			Request: sdk.ExecuteStageRequest{
				StageName:               "K8S_SYNC",
				StageConfig:             []byte(``),
				RunningDeploymentSource: sdk.DeploymentSource{},
				TargetDeploymentSource: sdk.DeploymentSource{
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
		status := plugin.executeK8sSyncStage(ctx, runningInput, []*sdk.DeployTarget[kubeConfigPkg.KubernetesDeployTargetConfig]{
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
		targetCfg, err := os.ReadFile(filepath.Join(target, "app.pipecd.yaml"))
		require.NoError(t, err)

		// prepare the input to ensure the running deployment exists
		targetInput := &sdk.ExecuteStageInput{
			Request: sdk.ExecuteStageRequest{
				StageName:   "K8S_SYNC",
				StageConfig: []byte(``),
				RunningDeploymentSource: sdk.DeploymentSource{
					ApplicationDirectory:      running,
					CommitHash:                "0123456789",
					ApplicationConfig:         runningCfg,
					ApplicationConfigFilename: "app.pipecd.yaml",
				},
				TargetDeploymentSource: sdk.DeploymentSource{
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
		status := plugin.executeK8sSyncStage(ctx, targetInput, []*sdk.DeployTarget[kubeConfigPkg.KubernetesDeployTargetConfig]{
			{
				Name:   "default",
				Config: *dtConfig,
			},
		})
		require.Equal(t, sdk.StageStatusSuccess, status)

		// The service should be removed from the previous namespace
		_, err = dynamicClient.Resource(schema.GroupVersionResource{Group: "", Version: "v1", Resource: "services"}).Namespace("test-1").Get(context.Background(), "simple", metav1.GetOptions{})
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
