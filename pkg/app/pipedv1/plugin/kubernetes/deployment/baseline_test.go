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
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	schema "k8s.io/apimachinery/pkg/runtime/schema"

	kubeConfigPkg "github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/config"
	"github.com/pipe-cd/pipecd/pkg/plugin/logpersister/logpersistertest"
	"github.com/pipe-cd/pipecd/pkg/plugin/sdk"
	"github.com/pipe-cd/pipecd/pkg/plugin/toolregistry/toolregistrytest"
)

func TestPlugin_executeK8sBaselineRolloutStage_withCreateService(t *testing.T) {
	t.Parallel()

	ctx := t.Context()

	// initialize tool registry
	testRegistry := toolregistrytest.NewTestToolRegistry(t)

	configDir := filepath.Join("testdata", "baseline_rollout_with_create_service")

	// read the application config from the example file
	appCfg := sdk.LoadApplicationConfigForTest[kubeConfigPkg.KubernetesApplicationSpec](t, filepath.Join(configDir, "app.pipecd.yaml"), "kubernetes")

	input := &sdk.ExecuteStageInput[kubeConfigPkg.KubernetesApplicationSpec]{
		Request: sdk.ExecuteStageRequest[kubeConfigPkg.KubernetesApplicationSpec]{
			StageName:   "K8S_BASELINE_ROLLOUT",
			StageConfig: []byte(`{"createService": true}`),
			RunningDeploymentSource: sdk.DeploymentSource[kubeConfigPkg.KubernetesApplicationSpec]{
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
		Client: sdk.NewClient(nil, "kubernetes", "", "", logpersistertest.NewTestLogPersister(t), testRegistry),
		Logger: zaptest.NewLogger(t),
	}

	// initialize deploy target config and dynamic client for assertions with envtest
	dtConfig, dynamicClient := setupTestDeployTargetConfigAndDynamicClient(t)

	plugin := &Plugin{}

	status := plugin.executeK8sBaselineRolloutStage(ctx, input, []*sdk.DeployTarget[kubeConfigPkg.KubernetesDeployTargetConfig]{
		{
			Name:   "default",
			Config: *dtConfig,
		},
	})

	assert.Equal(t, sdk.StageStatusSuccess, status)

	// Assert that Deployment and Service resources are created and have expected labels/annotations.
	deploymentRes := schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}
	deployment, err := dynamicClient.Resource(deploymentRes).Namespace("default").Get(ctx, "simple-baseline", metav1.GetOptions{})
	require.NoError(t, err)
	assert.Equal(t, "simple-baseline", deployment.GetName())
	assert.Equal(t, "simple", deployment.GetLabels()["app"])
	assert.Equal(t, "baseline", deployment.GetLabels()["pipecd.dev/variant"])
	assert.Equal(t, "baseline", deployment.GetAnnotations()["pipecd.dev/variant"])

	// Additional assertions for builtin labels and annotations
	assert.Equal(t, "piped", deployment.GetLabels()["pipecd.dev/managed-by"])
	assert.Equal(t, "piped-id", deployment.GetLabels()["pipecd.dev/piped"])
	assert.Equal(t, "app-id", deployment.GetLabels()["pipecd.dev/application"])
	assert.Equal(t, "0123456789", deployment.GetLabels()["pipecd.dev/commit-hash"])
	assert.Equal(t, "piped", deployment.GetAnnotations()["pipecd.dev/managed-by"])
	assert.Equal(t, "piped-id", deployment.GetAnnotations()["pipecd.dev/piped"])
	assert.Equal(t, "app-id", deployment.GetAnnotations()["pipecd.dev/application"])
	assert.Equal(t, "0123456789", deployment.GetAnnotations()["pipecd.dev/commit-hash"])
	assert.Equal(t, "apps/v1", deployment.GetAnnotations()["pipecd.dev/original-api-version"])
	assert.Equal(t, "apps:Deployment::simple-baseline", deployment.GetAnnotations()["pipecd.dev/resource-key"])

	serviceRes := schema.GroupVersionResource{Group: "", Version: "v1", Resource: "services"}
	service, err := dynamicClient.Resource(serviceRes).Namespace("default").Get(ctx, "simple-baseline", metav1.GetOptions{})
	require.NoError(t, err)
	assert.Equal(t, "simple-baseline", service.GetName())

	// Additional assertions for Service labels, annotations, selector, and ports
	assert.Equal(t, "piped", service.GetLabels()["pipecd.dev/managed-by"])
	assert.Equal(t, "piped-id", service.GetLabels()["pipecd.dev/piped"])
	assert.Equal(t, "app-id", service.GetLabels()["pipecd.dev/application"])
	assert.Equal(t, "0123456789", service.GetLabels()["pipecd.dev/commit-hash"])
	assert.Equal(t, "piped", service.GetAnnotations()["pipecd.dev/managed-by"])
	assert.Equal(t, "piped-id", service.GetAnnotations()["pipecd.dev/piped"])
	assert.Equal(t, "app-id", service.GetAnnotations()["pipecd.dev/application"])
	assert.Equal(t, "0123456789", service.GetAnnotations()["pipecd.dev/commit-hash"])
	assert.Equal(t, "v1", service.GetAnnotations()["pipecd.dev/original-api-version"])
	assert.Equal(t, ":Service::simple-baseline", service.GetAnnotations()["pipecd.dev/resource-key"])

	// Check Service selector and ports
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
	assert.Equal(t, "TCP", port["protocol"])
}

func TestPlugin_executeK8sBaselineRolloutStage_withoutCreateService(t *testing.T) {
	t.Parallel()

	ctx := t.Context()

	// initialize tool registry
	testRegistry := toolregistrytest.NewTestToolRegistry(t)

	configDir := filepath.Join("testdata", "baseline_rollout_without_create_service")

	// read the application config from the example file
	appCfg := sdk.LoadApplicationConfigForTest[kubeConfigPkg.KubernetesApplicationSpec](t, filepath.Join(configDir, "app.pipecd.yaml"), "kubernetes")

	input := &sdk.ExecuteStageInput[kubeConfigPkg.KubernetesApplicationSpec]{
		Request: sdk.ExecuteStageRequest[kubeConfigPkg.KubernetesApplicationSpec]{
			StageName:   "K8S_BASELINE_ROLLOUT",
			StageConfig: []byte(`{}`),
			RunningDeploymentSource: sdk.DeploymentSource[kubeConfigPkg.KubernetesApplicationSpec]{
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
		Client: sdk.NewClient(nil, "kubernetes", "", "", logpersistertest.NewTestLogPersister(t), testRegistry),
		Logger: zaptest.NewLogger(t),
	}

	// initialize deploy target config and dynamic client for assertions with envtest
	dtConfig, dynamicClient := setupTestDeployTargetConfigAndDynamicClient(t)

	plugin := &Plugin{}

	status := plugin.executeK8sBaselineRolloutStage(ctx, input, []*sdk.DeployTarget[kubeConfigPkg.KubernetesDeployTargetConfig]{
		{
			Name:   "default",
			Config: *dtConfig,
		},
	})

	assert.Equal(t, sdk.StageStatusSuccess, status)

	// Assert that Deployment and Service resources are created and have expected labels/annotations.
	deploymentRes := schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}
	deployment, err := dynamicClient.Resource(deploymentRes).Namespace("default").Get(ctx, "simple-baseline", metav1.GetOptions{})
	require.NoError(t, err)
	assert.Equal(t, "simple-baseline", deployment.GetName())
	assert.Equal(t, "simple", deployment.GetLabels()["app"])
	assert.Equal(t, "baseline", deployment.GetLabels()["pipecd.dev/variant"])
	assert.Equal(t, "baseline", deployment.GetAnnotations()["pipecd.dev/variant"])

	serviceRes := schema.GroupVersionResource{Group: "", Version: "v1", Resource: "services"}
	_, err = dynamicClient.Resource(serviceRes).Namespace("default").Get(ctx, "simple-baseline", metav1.GetOptions{})
	require.Error(t, err)
	assert.True(t, errors.IsNotFound(err))
}
