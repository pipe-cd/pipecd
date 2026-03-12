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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
	"github.com/pipe-cd/piped-plugin-sdk-go/logpersister/logpersistertest"
	"github.com/pipe-cd/piped-plugin-sdk-go/toolregistry/toolregistrytest"

	kubeconfig "github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes_multicluster/config"
)

func TestPlugin_executeK8sMultiCanaryRolloutStage_SingleCluster(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	// Load the application config from testdata.
	appCfg := sdk.LoadApplicationConfigForTest[kubeconfig.KubernetesApplicationSpec](t, filepath.Join("testdata", "canary", "app.pipecd.yaml"), "kubernetes_multicluster")

	testRegistry := toolregistrytest.NewTestToolRegistry(t)

	stageConfig := []byte(`{"replicas": "50%", "suffix": "canary"}`)

	input := &sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec]{
		Request: sdk.ExecuteStageRequest[kubeconfig.KubernetesApplicationSpec]{
			StageName:   StageK8sMultiCanaryRollout,
			StageConfig: stageConfig,
			TargetDeploymentSource: sdk.DeploymentSource[kubeconfig.KubernetesApplicationSpec]{
				ApplicationDirectory:      filepath.Join("testdata", "canary"),
				CommitHash:                "0123456789",
				ApplicationConfig:         appCfg,
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

	dtConfig, dynamicClient := setupTestDeployTargetConfigAndDynamicClient(t)

	plugin := &Plugin{}
	status := plugin.executeK8sMultiCanaryRolloutStage(ctx, input, []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]{
		{
			Name:   "default",
			Config: *dtConfig,
		},
	})

	assert.Equal(t, sdk.StageStatusSuccess, status)

	// The canary deployment should be created with "-canary" suffix.
	deployment, err := dynamicClient.Resource(schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}).Namespace("default").Get(ctx, "simple-canary", metav1.GetOptions{})
	require.NoError(t, err)

	assert.Equal(t, "simple-canary", deployment.GetName())

	// Verify variant label is set to "canary".
	assert.Equal(t, "canary", deployment.GetLabels()["pipecd.dev/variant"])
	assert.Equal(t, "canary", deployment.GetAnnotations()["pipecd.dev/variant"])

	// Verify replica count is 1 (50% of 2 = 1).
	spec, ok := deployment.Object["spec"].(map[string]interface{})
	require.True(t, ok)
	replicas, ok := spec["replicas"].(int64)
	require.True(t, ok)
	assert.Equal(t, int64(1), replicas)
}

func TestPlugin_executeK8sMultiCanaryRolloutStage_MultiCluster(t *testing.T) {
	t.Parallel()

	appCfg := sdk.LoadApplicationConfigForTest[kubeconfig.KubernetesApplicationSpec](t, filepath.Join("testdata", "canary", "app.pipecd.yaml"), "kubernetes_multicluster")

	testRegistry := toolregistrytest.NewTestToolRegistry(t)

	stageConfig := []byte(`{"replicas": 1, "suffix": "canary"}`)

	input := &sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec]{
		Request: sdk.ExecuteStageRequest[kubeconfig.KubernetesApplicationSpec]{
			StageName:   StageK8sMultiCanaryRollout,
			StageConfig: stageConfig,
			TargetDeploymentSource: sdk.DeploymentSource[kubeconfig.KubernetesApplicationSpec]{
				ApplicationDirectory:      filepath.Join("testdata", "canary"),
				CommitHash:                "0123456789",
				ApplicationConfig:         appCfg,
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
	status := plugin.executeK8sMultiCanaryRolloutStage(t.Context(), input, dts)

	require.Equal(t, sdk.StageStatusSuccess, status)

	// Both clusters should have a canary deployment.
	for _, cl := range []*cluster{cluster1, cluster2} {
		deployment, err := cl.cli.Resource(schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}).Namespace("default").Get(context.Background(), "simple-canary", metav1.GetOptions{})
		require.NoError(t, err)

		assert.Equal(t, "simple-canary", deployment.GetName())
		assert.Equal(t, "canary", deployment.GetLabels()["pipecd.dev/variant"])
		assert.Equal(t, "piped-id", deployment.GetLabels()["pipecd.dev/piped"])
		assert.Equal(t, "app-id", deployment.GetLabels()["pipecd.dev/application"])
	}
}

func TestPlugin_executeK8sMultiCanaryRolloutStage_Failure(t *testing.T) {
	t.Parallel()

	// Use an invalid kubeconfig path to force failure.
	appCfg := sdk.LoadApplicationConfigForTest[kubeconfig.KubernetesApplicationSpec](t, filepath.Join("testdata", "canary", "app.pipecd.yaml"), "kubernetes_multicluster")

	testRegistry := toolregistrytest.NewTestToolRegistry(t)

	stageConfig := []byte(`{"replicas": 1}`)

	input := &sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec]{
		Request: sdk.ExecuteStageRequest[kubeconfig.KubernetesApplicationSpec]{
			StageName:   StageK8sMultiCanaryRollout,
			StageConfig: stageConfig,
			TargetDeploymentSource: sdk.DeploymentSource[kubeconfig.KubernetesApplicationSpec]{
				ApplicationDirectory:      filepath.Join("testdata", "canary"),
				CommitHash:                "0123456789",
				ApplicationConfig:         appCfg,
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
	status := plugin.executeK8sMultiCanaryRolloutStage(t.Context(), input, dts)

	assert.Equal(t, sdk.StageStatusFailure, status)
}
