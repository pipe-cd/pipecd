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

	kubeConfigPkg "github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes_multicluster/config"
	"github.com/pipe-cd/pipecd/pkg/plugin/logpersister/logpersistertest"
	"github.com/pipe-cd/pipecd/pkg/plugin/sdk"
	"github.com/pipe-cd/pipecd/pkg/plugin/sdk/sdktest"
	"github.com/pipe-cd/pipecd/pkg/plugin/toolregistry/toolregistrytest"
)

func TestPlugin_executeK8sMultiRollbackStage_NoPreviousDeployment(t *testing.T) {
	t.Parallel()

	// Initialize tool registry
	testRegistry := toolregistrytest.NewTestToolRegistry(t)

	// Initialize deploy target config and dynamic client for assertions with envtest
	dtConfig, _ := setupTestDeployTargetConfigAndDynamicClient(t)

	// Read the application config from the example file
	appCfg := sdktest.LoadApplicationConfig[kubeConfigPkg.KubernetesApplicationSpec](t, filepath.Join(examplesDir(), "kubernetes", "simple", "app.pipecd.yaml"))

	// Prepare the input
	input := &sdk.ExecuteStageInput[kubeConfigPkg.KubernetesApplicationSpec]{
		Request: sdk.ExecuteStageRequest[kubeConfigPkg.KubernetesApplicationSpec]{
			StageName:   "K8S_MULTI_ROLLBACK",
			StageConfig: []byte(``),
			RunningDeploymentSource: sdk.DeploymentSource[kubeConfigPkg.KubernetesApplicationSpec]{
				CommitHash: "", // Empty commit hash indicates no previous deployment
			},
			TargetDeploymentSource: sdk.DeploymentSource[kubeConfigPkg.KubernetesApplicationSpec]{
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

	plugin := &Plugin{}
	status := plugin.executeK8sMultiRollbackStage(t.Context(), input, []*sdk.DeployTarget[kubeConfigPkg.KubernetesDeployTargetConfig]{
		{
			Name:   "default",
			Config: *dtConfig,
		},
	})

	assert.Equal(t, sdk.StageStatusFailure, status)
}

func TestPlugin_executeK8sMultiRollbackStage_SuccessfulRollback(t *testing.T) {
	t.Parallel()

	// Initialize tool registry
	testRegistry := toolregistrytest.NewTestToolRegistry(t)

	// Initialize deploy target config and dynamic client for assertions with envtest
	dtConfig, dynamicClient := setupTestDeployTargetConfigAndDynamicClient(t)

	// Read the application config from the example file
	appCfg := sdktest.LoadApplicationConfig[kubeConfigPkg.KubernetesApplicationSpec](t, filepath.Join(examplesDir(), "kubernetes", "simple", "app.pipecd.yaml"))

	// Prepare the input
	input := &sdk.ExecuteStageInput[kubeConfigPkg.KubernetesApplicationSpec]{
		Request: sdk.ExecuteStageRequest[kubeConfigPkg.KubernetesApplicationSpec]{
			StageName:   "K8S_MULTI_ROLLBACK",
			StageConfig: []byte(``),
			RunningDeploymentSource: sdk.DeploymentSource[kubeConfigPkg.KubernetesApplicationSpec]{
				ApplicationDirectory:      filepath.Join(examplesDir(), "kubernetes", "simple"),
				CommitHash:                "previous-hash",
				ApplicationConfig:         appCfg,
				ApplicationConfigFilename: "app.pipecd.yaml",
			},
			TargetDeploymentSource: sdk.DeploymentSource[kubeConfigPkg.KubernetesApplicationSpec]{
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

	plugin := &Plugin{}
	status := plugin.executeK8sMultiRollbackStage(t.Context(), input, []*sdk.DeployTarget[kubeConfigPkg.KubernetesDeployTargetConfig]{
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

func TestPlugin_executeK8sMultiRollbackStage_WithVariantLabels(t *testing.T) {
	t.Parallel()

	// Initialize tool registry
	testRegistry := toolregistrytest.NewTestToolRegistry(t)

	// Initialize deploy target config and dynamic client for assertions with envtest
	dtConfig, dynamicClient := setupTestDeployTargetConfigAndDynamicClient(t)

	// Read the application config and modify it to include variant labels
	appCfg := sdktest.LoadApplicationConfig[kubeConfigPkg.KubernetesApplicationSpec](t, filepath.Join(examplesDir(), "kubernetes", "simple", "app.pipecd.yaml"))

	// Prepare the input
	input := &sdk.ExecuteStageInput[kubeConfigPkg.KubernetesApplicationSpec]{
		Request: sdk.ExecuteStageRequest[kubeConfigPkg.KubernetesApplicationSpec]{
			StageName:   "K8S_MULTI_ROLLBACK",
			StageConfig: []byte(``),
			RunningDeploymentSource: sdk.DeploymentSource[kubeConfigPkg.KubernetesApplicationSpec]{
				ApplicationDirectory:      filepath.Join(examplesDir(), "kubernetes", "simple"),
				CommitHash:                "previous-hash",
				ApplicationConfig:         appCfg,
				ApplicationConfigFilename: "app.pipecd.yaml",
			},
			TargetDeploymentSource: sdk.DeploymentSource[kubeConfigPkg.KubernetesApplicationSpec]{
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

	plugin := &Plugin{}
	status := plugin.executeK8sMultiRollbackStage(t.Context(), input, []*sdk.DeployTarget[kubeConfigPkg.KubernetesDeployTargetConfig]{
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
