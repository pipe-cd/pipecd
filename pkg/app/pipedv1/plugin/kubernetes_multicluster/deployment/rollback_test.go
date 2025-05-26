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
	"github.com/pipe-cd/pipecd/pkg/plugin/toolregistry/toolregistrytest"
	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
)

func TestPlugin_executeK8sMultiRollbackStage_compatibility_k8sPlugin_NoPreviousDeployment(t *testing.T) {
	t.Parallel()

	// Initialize tool registry
	testRegistry := toolregistrytest.NewTestToolRegistry(t)

	// Initialize deploy target config and dynamic client for assertions with envtest
	dtConfig, _ := setupTestDeployTargetConfigAndDynamicClient(t)

	// Read the application config from the example file
	appCfg := sdk.LoadApplicationConfigForTest[kubeconfig.KubernetesApplicationSpec](t, filepath.Join("testdata", "simple", "app.pipecd.yaml"), "kubernetes_multicluster")

	// Prepare the input
	input := &sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec]{
		Request: sdk.ExecuteStageRequest[kubeconfig.KubernetesApplicationSpec]{
			StageName:   "K8S_MULTI_ROLLBACK",
			StageConfig: []byte(``),
			RunningDeploymentSource: sdk.DeploymentSource[kubeconfig.KubernetesApplicationSpec]{
				CommitHash: "", // Empty commit hash indicates no previous deployment
			},
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
		Client: sdk.NewClient(nil, "kubernetes", "", "", logpersistertest.NewTestLogPersister(t), testRegistry),
		Logger: zaptest.NewLogger(t),
	}

	plugin := &Plugin{}
	status := plugin.executeK8sMultiRollbackStage(t.Context(), input, []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]{
		{
			Name:   "default",
			Config: *dtConfig,
		},
	})

	assert.Equal(t, sdk.StageStatusFailure, status)
}

func TestPlugin_executeK8sMultiRollbackStage_compatibility_k8sPlugin_SuccessfulRollback(t *testing.T) {
	t.Parallel()

	// Initialize tool registry
	testRegistry := toolregistrytest.NewTestToolRegistry(t)

	// Initialize deploy target config and dynamic client for assertions with envtest
	dtConfig, dynamicClient := setupTestDeployTargetConfigAndDynamicClient(t)

	// Read the application config from the example file
	appCfg := sdk.LoadApplicationConfigForTest[kubeconfig.KubernetesApplicationSpec](t, filepath.Join("testdata", "simple", "app.pipecd.yaml"), "kubernetes_multicluster")

	// Prepare the input
	input := &sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec]{
		Request: sdk.ExecuteStageRequest[kubeconfig.KubernetesApplicationSpec]{
			StageName:   "K8S_MULTI_ROLLBACK",
			StageConfig: []byte(``),
			RunningDeploymentSource: sdk.DeploymentSource[kubeconfig.KubernetesApplicationSpec]{
				ApplicationDirectory:      filepath.Join("testdata", "simple"),
				CommitHash:                "previous-hash",
				ApplicationConfig:         appCfg,
				ApplicationConfigFilename: "app.pipecd.yaml",
			},
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
		Client: sdk.NewClient(nil, "kubernetes", "", "", logpersistertest.NewTestLogPersister(t), testRegistry),
		Logger: zaptest.NewLogger(t),
	}

	plugin := &Plugin{}
	status := plugin.executeK8sMultiRollbackStage(t.Context(), input, []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]{
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

func TestPlugin_executeK8sMultiRollbackStage_compatibility_k8sPlugin_WithVariantLabels(t *testing.T) {
	t.Parallel()

	// Initialize tool registry
	testRegistry := toolregistrytest.NewTestToolRegistry(t)

	// Initialize deploy target config and dynamic client for assertions with envtest
	dtConfig, dynamicClient := setupTestDeployTargetConfigAndDynamicClient(t)

	// Read the application config and modify it to include variant labels
	appCfg := sdk.LoadApplicationConfigForTest[kubeconfig.KubernetesApplicationSpec](t, filepath.Join("testdata", "simple", "app.pipecd.yaml"), "kubernetes_multicluster")

	// Prepare the input
	input := &sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec]{
		Request: sdk.ExecuteStageRequest[kubeconfig.KubernetesApplicationSpec]{
			StageName:   "K8S_MULTI_ROLLBACK",
			StageConfig: []byte(``),
			RunningDeploymentSource: sdk.DeploymentSource[kubeconfig.KubernetesApplicationSpec]{
				ApplicationDirectory:      filepath.Join("testdata", "simple"),
				CommitHash:                "previous-hash",
				ApplicationConfig:         appCfg,
				ApplicationConfigFilename: "app.pipecd.yaml",
			},
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
		Client: sdk.NewClient(nil, "kubernetes", "", "", logpersistertest.NewTestLogPersister(t), testRegistry),
		Logger: zaptest.NewLogger(t),
	}

	plugin := &Plugin{}
	status := plugin.executeK8sMultiRollbackStage(t.Context(), input, []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]{
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

func TestPlugin_executeK8sMultiRollbackStage_SuccessfulRollback(t *testing.T) {
	t.Parallel()

	// Initialize tool registry
	testRegistry := toolregistrytest.NewTestToolRegistry(t)

	// Read the application config from the example file
	appCfg := sdk.LoadApplicationConfigForTest[kubeconfig.KubernetesApplicationSpec](t, filepath.Join("testdata", "simple", "app.pipecd.yaml"), "kubernetes_multicluster")

	// Prepare the input
	input := &sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec]{
		Request: sdk.ExecuteStageRequest[kubeconfig.KubernetesApplicationSpec]{
			StageName:   "K8S_MULTI_ROLLBACK",
			StageConfig: []byte(``),
			RunningDeploymentSource: sdk.DeploymentSource[kubeconfig.KubernetesApplicationSpec]{
				ApplicationDirectory:      filepath.Join("testdata", "simple"),
				CommitHash:                "previous-hash",
				ApplicationConfig:         appCfg,
				ApplicationConfigFilename: "app.pipecd.yaml",
			},
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
		Client: sdk.NewClient(nil, "kubernete_multicluster", "", "", logpersistertest.NewTestLogPersister(t), testRegistry),
		Logger: zaptest.NewLogger(t),
	}

	// prepare the cluster
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
	status := plugin.executeK8sMultiRollbackStage(t.Context(), input, dts)

	assert.Equal(t, sdk.StageStatusSuccess, status)

	// Verify the deployment was rolled back
	for _, cluster := range []*cluster{cluster1, cluster2} {
		deployment, err := cluster.cli.Resource(schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}).Namespace("default").Get(t.Context(), "simple", metav1.GetOptions{})
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
}

func TestPlugin_executeK8sMultiRollbackStage_SuccessfulRollback_when_adding_multiTarget(t *testing.T) {
	t.Parallel()

	// Initialize tool registry
	testRegistry := toolregistrytest.NewTestToolRegistry(t)

	// prepare the cluster
	cluster1 := setupCluster(t, "cluster1")
	cluster2 := setupCluster(t, "cluster2")

	ok := t.Run("prepare", func(t *testing.T) {
		// prepare the input to ensure the running deployment exists
		runningInput := &sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec]{
			Request: sdk.ExecuteStageRequest[kubeconfig.KubernetesApplicationSpec]{
				StageName:               "K8S_MULTI_SYNC",
				StageConfig:             []byte(``),
				RunningDeploymentSource: sdk.DeploymentSource[kubeconfig.KubernetesApplicationSpec]{},
				TargetDeploymentSource: sdk.DeploymentSource[kubeconfig.KubernetesApplicationSpec]{
					ApplicationDirectory:      filepath.Join("testdata", "add_deploy_target", "running"),
					CommitHash:                "previous-hash",
					ApplicationConfig:         sdk.LoadApplicationConfigForTest[kubeconfig.KubernetesApplicationSpec](t, filepath.Join("testdata", "add_deploy_target", "running", "app.pipecd.yaml"), "kubernetes_multicluster"),
					ApplicationConfigFilename: "app.pipecd.yaml",
				},
				Deployment: sdk.Deployment{
					PipedID:       "piped-id",
					ApplicationID: "app-id",
				},
			},
			Client: sdk.NewClient(nil, "kubernete_multicluster", "", "", logpersistertest.NewTestLogPersister(t), testRegistry),
			Logger: zaptest.NewLogger(t),
		}

		plugin := &Plugin{}
		status := plugin.executeK8sMultiSyncStage(t.Context(), runningInput, []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]{
			{
				Name:   "cluster1",
				Config: *cluster1.dtc,
			},
		})
		require.Equal(t, sdk.StageStatusSuccess, status)

		deployment, err := cluster1.cli.Resource(schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}).Namespace("default").Get(t.Context(), fmt.Sprintf("simple-%s", cluster1.name), metav1.GetOptions{})
		require.NoError(t, err)

		assert.Equal(t, "piped", deployment.GetLabels()["pipecd.dev/managed-by"])
		assert.Equal(t, "piped-id", deployment.GetLabels()["pipecd.dev/piped"])
		assert.Equal(t, "app-id", deployment.GetLabels()["pipecd.dev/application"])
		assert.Equal(t, "previous-hash", deployment.GetLabels()["pipecd.dev/commit-hash"])

		assert.Equal(t, "piped", deployment.GetAnnotations()["pipecd.dev/managed-by"])
		assert.Equal(t, "piped-id", deployment.GetAnnotations()["pipecd.dev/piped"])
		assert.Equal(t, "app-id", deployment.GetAnnotations()["pipecd.dev/application"])
		assert.Equal(t, "apps/v1", deployment.GetAnnotations()["pipecd.dev/original-api-version"])
		assert.Equal(t, "apps:Deployment::"+fmt.Sprintf("simple-%s", cluster1.name), deployment.GetAnnotations()["pipecd.dev/resource-key"])
		assert.Equal(t, "previous-hash", deployment.GetAnnotations()["pipecd.dev/commit-hash"])
	})
	require.Truef(t, ok, "expected prepare to succeed")

	// Prepare the input
	input := &sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec]{
		Request: sdk.ExecuteStageRequest[kubeconfig.KubernetesApplicationSpec]{
			StageName:   "K8S_MULTI_ROLLBACK",
			StageConfig: []byte(``),
			RunningDeploymentSource: sdk.DeploymentSource[kubeconfig.KubernetesApplicationSpec]{
				ApplicationDirectory:      filepath.Join("testdata", "add_deploy_target", "running"),
				CommitHash:                "previous-hash",
				ApplicationConfig:         sdk.LoadApplicationConfigForTest[kubeconfig.KubernetesApplicationSpec](t, filepath.Join("testdata", "add_deploy_target", "running", "app.pipecd.yaml"), "kubernetes_multicluster"),
				ApplicationConfigFilename: "app.pipecd.yaml",
			},
			TargetDeploymentSource: sdk.DeploymentSource[kubeconfig.KubernetesApplicationSpec]{
				ApplicationDirectory:      filepath.Join("testdata", "add_deploy_target", "target"),
				CommitHash:                "0123456789",
				ApplicationConfig:         sdk.LoadApplicationConfigForTest[kubeconfig.KubernetesApplicationSpec](t, filepath.Join("testdata", "add_deploy_target", "target", "app.pipecd.yaml"), "kubernetes_multicluster"),
				ApplicationConfigFilename: "app.pipecd.yaml",
			},
			Deployment: sdk.Deployment{
				PipedID:       "piped-id",
				ApplicationID: "app-id",
			},
		},
		Client: sdk.NewClient(nil, "kubernete_multicluster", "", "", logpersistertest.NewTestLogPersister(t), testRegistry),
		Logger: zaptest.NewLogger(t),
	}

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
	status := plugin.executeK8sMultiRollbackStage(t.Context(), input, dts)

	assert.Equal(t, sdk.StageStatusSuccess, status)

	{
		// for cluster1
		deployment, err := cluster1.cli.Resource(schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}).Namespace("default").Get(t.Context(), fmt.Sprintf("simple-%s", cluster1.name), metav1.GetOptions{})
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
		assert.Equal(t, "apps:Deployment::"+fmt.Sprintf("simple-%s", cluster1.name), deployment.GetAnnotations()["pipecd.dev/resource-key"])
		assert.Equal(t, "previous-hash", deployment.GetAnnotations()["pipecd.dev/commit-hash"])
	}

	{
		// for cluster2
		_, err := cluster2.cli.Resource(schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}).Namespace("default").Get(t.Context(), fmt.Sprintf("simple-%s", cluster2.name), metav1.GetOptions{})
		assert.True(t, apierrors.IsNotFound(err))
	}
}

func TestPlugin_executeK8sMultiRollbackStage_SuccessfulRollback_when_removing_multiTarget(t *testing.T) {
	t.Parallel()

	// Initialize tool registry
	testRegistry := toolregistrytest.NewTestToolRegistry(t)

	// prepare the cluster
	cluster1 := setupCluster(t, "cluster1")
	cluster2 := setupCluster(t, "cluster2")

	ok := t.Run("prepare", func(t *testing.T) {
		// prepare the input to ensure the running deployment exists
		runningInput := &sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec]{
			Request: sdk.ExecuteStageRequest[kubeconfig.KubernetesApplicationSpec]{
				StageName:               "K8S_MULTI_SYNC",
				StageConfig:             []byte(``),
				RunningDeploymentSource: sdk.DeploymentSource[kubeconfig.KubernetesApplicationSpec]{},
				TargetDeploymentSource: sdk.DeploymentSource[kubeconfig.KubernetesApplicationSpec]{
					ApplicationDirectory:      filepath.Join("testdata", "remove_deploy_target", "running"),
					CommitHash:                "previous-hash",
					ApplicationConfig:         sdk.LoadApplicationConfigForTest[kubeconfig.KubernetesApplicationSpec](t, filepath.Join("testdata", "remove_deploy_target", "running", "app.pipecd.yaml"), "kubernetes_multicluster"),
					ApplicationConfigFilename: "app.pipecd.yaml",
				},
				Deployment: sdk.Deployment{
					PipedID:       "piped-id",
					ApplicationID: "app-id",
				},
			},
			Client: sdk.NewClient(nil, "kubernete_multicluster", "", "", logpersistertest.NewTestLogPersister(t), testRegistry),
			Logger: zaptest.NewLogger(t),
		}

		plugin := &Plugin{}
		status := plugin.executeK8sMultiSyncStage(t.Context(), runningInput, []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]{
			{
				Name:   "cluster1",
				Config: *cluster1.dtc,
			},
			{
				Name:   "cluster2",
				Config: *cluster2.dtc,
			},
		})
		require.Equal(t, sdk.StageStatusSuccess, status)
	})
	require.Truef(t, ok, "expected prepare to succeed")

	// Prepare the input
	input := &sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec]{
		Request: sdk.ExecuteStageRequest[kubeconfig.KubernetesApplicationSpec]{
			StageName:   "K8S_MULTI_ROLLBACK",
			StageConfig: []byte(``),
			RunningDeploymentSource: sdk.DeploymentSource[kubeconfig.KubernetesApplicationSpec]{
				ApplicationDirectory:      filepath.Join("testdata", "remove_deploy_target", "running"),
				CommitHash:                "previous-hash",
				ApplicationConfig:         sdk.LoadApplicationConfigForTest[kubeconfig.KubernetesApplicationSpec](t, filepath.Join("testdata", "remove_deploy_target", "running", "app.pipecd.yaml"), "kubernetes_multicluster"),
				ApplicationConfigFilename: "app.pipecd.yaml",
			},
			TargetDeploymentSource: sdk.DeploymentSource[kubeconfig.KubernetesApplicationSpec]{
				ApplicationDirectory:      filepath.Join("testdata", "remove_deploy_target", "target"),
				CommitHash:                "0123456789",
				ApplicationConfig:         sdk.LoadApplicationConfigForTest[kubeconfig.KubernetesApplicationSpec](t, filepath.Join("testdata", "remove_deploy_target", "target", "app.pipecd.yaml"), "kubernetes_multicluster"),
				ApplicationConfigFilename: "app.pipecd.yaml",
			},
			Deployment: sdk.Deployment{
				PipedID:       "piped-id",
				ApplicationID: "app-id",
			},
		},
		Client: sdk.NewClient(nil, "kubernete_multicluster", "", "", logpersistertest.NewTestLogPersister(t), testRegistry),
		Logger: zaptest.NewLogger(t),
	}

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
	status := plugin.executeK8sMultiRollbackStage(t.Context(), input, dts)

	assert.Equal(t, sdk.StageStatusSuccess, status)

	for _, cluster := range []*cluster{cluster1, cluster2} {
		deployment, err := cluster.cli.Resource(schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}).Namespace("default").Get(t.Context(), fmt.Sprintf("simple-%s", cluster.name), metav1.GetOptions{})
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
		assert.Equal(t, "apps:Deployment::"+fmt.Sprintf("simple-%s", cluster.name), deployment.GetAnnotations()["pipecd.dev/resource-key"])
		assert.Equal(t, "previous-hash", deployment.GetAnnotations()["pipecd.dev/commit-hash"])
	}
}

func TestPlugin_executeK8sMultiRollbackStage_FailureRollback_when_all_rollback_of_targets_are_failed(t *testing.T) {
	t.Parallel()

	// Initialize tool registry
	testRegistry := toolregistrytest.NewTestToolRegistry(t)

	// prepare the cluster
	cluster1 := setupCluster(t, "cluster1")
	cluster2 := setupCluster(t, "cluster2")

	// Prepare the input
	input := &sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec]{
		Request: sdk.ExecuteStageRequest[kubeconfig.KubernetesApplicationSpec]{
			StageName:   "K8S_MULTI_ROLLBACK",
			StageConfig: []byte(``),
			RunningDeploymentSource: sdk.DeploymentSource[kubeconfig.KubernetesApplicationSpec]{
				ApplicationDirectory:      filepath.Join("testdata", "failed_all_of_rollback", "running"),
				CommitHash:                "previous-hash",
				ApplicationConfig:         sdk.LoadApplicationConfigForTest[kubeconfig.KubernetesApplicationSpec](t, filepath.Join("testdata", "failed_all_of_rollback", "running", "app.pipecd.yaml"), "kubernetes_multicluster"),
				ApplicationConfigFilename: "app.pipecd.yaml",
			},
			TargetDeploymentSource: sdk.DeploymentSource[kubeconfig.KubernetesApplicationSpec]{},
			Deployment: sdk.Deployment{
				PipedID:       "piped-id",
				ApplicationID: "app-id",
			},
		},
		Client: sdk.NewClient(nil, "kubernete_multicluster", "", "", logpersistertest.NewTestLogPersister(t), testRegistry),
		Logger: zaptest.NewLogger(t),
	}

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
	status := plugin.executeK8sMultiRollbackStage(t.Context(), input, dts)

	assert.Equal(t, sdk.StageStatusFailure, status)
}

func TestPlugin_executeK8sMultiRollbackStage_FailureRollback_when_at_least_one_of_targets_are_success(t *testing.T) {
	t.Parallel()

	// Initialize tool registry
	testRegistry := toolregistrytest.NewTestToolRegistry(t)

	// prepare the cluster
	cluster1 := setupCluster(t, "cluster1")
	cluster2 := setupCluster(t, "cluster2")

	// Prepare the input
	input := &sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec]{
		Request: sdk.ExecuteStageRequest[kubeconfig.KubernetesApplicationSpec]{
			StageName:   "K8S_MULTI_ROLLBACK",
			StageConfig: []byte(``),
			RunningDeploymentSource: sdk.DeploymentSource[kubeconfig.KubernetesApplicationSpec]{
				ApplicationDirectory:      filepath.Join("testdata", "succes_one_of_rollback", "running"),
				CommitHash:                "previous-hash",
				ApplicationConfig:         sdk.LoadApplicationConfigForTest[kubeconfig.KubernetesApplicationSpec](t, filepath.Join("testdata", "succes_one_of_rollback", "running", "app.pipecd.yaml"), "kubernetes_multicluster"),
				ApplicationConfigFilename: "app.pipecd.yaml",
			},
			TargetDeploymentSource: sdk.DeploymentSource[kubeconfig.KubernetesApplicationSpec]{},
			Deployment: sdk.Deployment{
				PipedID:       "piped-id",
				ApplicationID: "app-id",
			},
		},
		Client: sdk.NewClient(nil, "kubernete_multicluster", "", "", logpersistertest.NewTestLogPersister(t), testRegistry),
		Logger: zaptest.NewLogger(t),
	}

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
	status := plugin.executeK8sMultiRollbackStage(t.Context(), input, dts)

	assert.Equal(t, sdk.StageStatusSuccess, status)
}
