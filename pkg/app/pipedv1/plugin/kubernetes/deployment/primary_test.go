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
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	schema "k8s.io/apimachinery/pkg/runtime/schema"

	kubeConfigPkg "github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/config"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/provider"
	"github.com/pipe-cd/pipecd/pkg/plugin/logpersister/logpersistertest"
	"github.com/pipe-cd/pipecd/pkg/plugin/sdk"
	"github.com/pipe-cd/pipecd/pkg/plugin/toolregistry/toolregistrytest"
)

func TestPlugin_executeK8sPrimaryRolloutStage(t *testing.T) {
	t.Parallel()

	ctx := t.Context()

	// initialize tool registry
	testRegistry := toolregistrytest.NewTestToolRegistry(t)

	// read the application config from the example file
	appCfg := sdk.LoadApplicationConfigForTest[kubeConfigPkg.KubernetesApplicationSpec](t, filepath.Join("testdata", "primary_rollout", "app.pipecd.yaml"), "kubernetes")

	input := &sdk.ExecuteStageInput[kubeConfigPkg.KubernetesApplicationSpec]{
		Request: sdk.ExecuteStageRequest[kubeConfigPkg.KubernetesApplicationSpec]{
			StageName:   "K8S_PRIMARY_ROLLOUT",
			StageConfig: []byte(`{}`),
			TargetDeploymentSource: sdk.DeploymentSource[kubeConfigPkg.KubernetesApplicationSpec]{
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

	// initialize deploy target config and dynamic client for assertions with envtest
	dtConfig, dynamicClient := setupTestDeployTargetConfigAndDynamicClient(t)

	plugin := &Plugin{}

	status := plugin.executeK8sPrimaryRolloutStage(ctx, input, []*sdk.DeployTarget[kubeConfigPkg.KubernetesDeployTargetConfig]{
		{
			Name:   "default",
			Config: *dtConfig,
		},
	})

	assert.Equal(t, sdk.StageStatusSuccess, status)

	// Assert that Deployment and Service resources are created and have expected labels/annotations.
	deploymentRes := schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}
	deployment, err := dynamicClient.Resource(deploymentRes).Namespace("default").Get(ctx, "simple", metav1.GetOptions{})
	if err != nil {
		t.Fatalf("failed to get deployment: %v", err)
	}
	assert.Equal(t, "simple", deployment.GetName())
	assert.Equal(t, "simple", deployment.GetLabels()["app"])
	assert.Equal(t, "primary", deployment.GetLabels()["pipecd.dev/variant"])
	assert.Equal(t, "primary", deployment.GetAnnotations()["pipecd.dev/variant"])

	serviceRes := schema.GroupVersionResource{Group: "", Version: "v1", Resource: "services"}
	service, err := dynamicClient.Resource(serviceRes).Namespace("default").Get(ctx, "simple", metav1.GetOptions{})
	if err != nil {
		t.Fatalf("failed to get service: %v", err)
	}
	assert.Equal(t, "simple", service.GetName())
}

func TestPlugin_executeK8sPrimaryRolloutStage_withPrune(t *testing.T) {
	t.Parallel()

	ctx := t.Context()

	// initialize tool registry
	testRegistry := toolregistrytest.NewTestToolRegistry(t)

	// initialize deploy target config and dynamic client for assertions with envtest
	dtConfig, dynamicClient := setupTestDeployTargetConfigAndDynamicClient(t)

	runningOk := t.Run("prepare running state", func(t *testing.T) {
		running := filepath.Join("testdata", "primary_rollout_prune", "running")
		runningCfg := sdk.LoadApplicationConfigForTest[kubeConfigPkg.KubernetesApplicationSpec](t, filepath.Join(running, "app.pipecd.yaml"), "kubernetes")

		runningInput := &sdk.ExecuteStageInput[kubeConfigPkg.KubernetesApplicationSpec]{
			Request: sdk.ExecuteStageRequest[kubeConfigPkg.KubernetesApplicationSpec]{
				StageName:   "K8S_PRIMARY_ROLLOUT",
				StageConfig: []byte(`{"prune":true}`),
				TargetDeploymentSource: sdk.DeploymentSource[kubeConfigPkg.KubernetesApplicationSpec]{
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
			Client: sdk.NewClient(nil, "kubernetes", "", "", logpersistertest.NewTestLogPersister(t), testRegistry),
			Logger: zaptest.NewLogger(t),
		}

		plugin := &Plugin{}
		status := plugin.executeK8sPrimaryRolloutStage(ctx, runningInput, []*sdk.DeployTarget[kubeConfigPkg.KubernetesDeployTargetConfig]{
			{
				Name:   "default",
				Config: *dtConfig,
			},
		})
		assert.Equal(t, sdk.StageStatusSuccess, status)

		deploymentRes := schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}
		_, err := dynamicClient.Resource(deploymentRes).Namespace("default").Get(ctx, "simple", metav1.GetOptions{})
		assert.NoError(t, err)
		serviceRes := schema.GroupVersionResource{Group: "", Version: "v1", Resource: "services"}
		_, err = dynamicClient.Resource(serviceRes).Namespace("default").Get(ctx, "simple", metav1.GetOptions{})
		assert.NoError(t, err)
	})
	require.True(t, runningOk, "prepare running state subtest failed, aborting")

	t.Run("prune with target state", func(t *testing.T) {
		target := filepath.Join("testdata", "primary_rollout_prune", "target")
		targetCfg := sdk.LoadApplicationConfigForTest[kubeConfigPkg.KubernetesApplicationSpec](t, filepath.Join(target, "app.pipecd.yaml"), "kubernetes")

		running := filepath.Join("testdata", "primary_rollout_prune", "running")
		runningCfg := sdk.LoadApplicationConfigForTest[kubeConfigPkg.KubernetesApplicationSpec](t, filepath.Join(running, "app.pipecd.yaml"), "kubernetes")
		targetInput := &sdk.ExecuteStageInput[kubeConfigPkg.KubernetesApplicationSpec]{
			Request: sdk.ExecuteStageRequest[kubeConfigPkg.KubernetesApplicationSpec]{
				StageName:   "K8S_PRIMARY_ROLLOUT",
				StageConfig: []byte(`{"prune":true}`),
				RunningDeploymentSource: sdk.DeploymentSource[kubeConfigPkg.KubernetesApplicationSpec]{
					ApplicationDirectory:      running,
					CommitHash:                "0123456789",
					ApplicationConfig:         runningCfg,
					ApplicationConfigFilename: "app.pipecd.yaml",
				},
				TargetDeploymentSource: sdk.DeploymentSource[kubeConfigPkg.KubernetesApplicationSpec]{
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
			Client: sdk.NewClient(nil, "kubernetes", "", "", logpersistertest.NewTestLogPersister(t), testRegistry),
			Logger: zaptest.NewLogger(t),
		}

		plugin := &Plugin{}
		status := plugin.executeK8sPrimaryRolloutStage(ctx, targetInput, []*sdk.DeployTarget[kubeConfigPkg.KubernetesDeployTargetConfig]{
			{
				Name:   "default",
				Config: *dtConfig,
			},
		})
		assert.Equal(t, sdk.StageStatusSuccess, status)

		deploymentRes := schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}
		_, err := dynamicClient.Resource(deploymentRes).Namespace("default").Get(ctx, "simple", metav1.GetOptions{})
		assert.NoError(t, err)
		serviceRes := schema.GroupVersionResource{Group: "", Version: "v1", Resource: "services"}
		_, err = dynamicClient.Resource(serviceRes).Namespace("default").Get(ctx, "simple", metav1.GetOptions{})
		assert.Error(t, err)
		require.Truef(t, apierrors.IsNotFound(err), "expected error to be NotFound, but got %v", err)
	})
}

func TestGeneratePrimaryManifests(t *testing.T) {
	t.Parallel()

	appCfg := &kubeConfigPkg.KubernetesApplicationSpec{
		Service: kubeConfigPkg.K8sResourceReference{
			Kind: "Service",
			Name: "my-service",
		},
	}

	const variantLabel = "pipecd.dev/variant"
	const variant = "primary"

	testcases := []struct {
		name       string
		inputYAML  string
		stageCfg   kubeConfigPkg.K8sPrimaryRolloutStageOptions
		expectYAML string
		expectErr  bool
	}{
		{
			name: "add variant label to selector",
			inputYAML: `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: simple
spec:
  selector:
    matchLabels:
      app: simple
  template:
    metadata:
      labels:
        app: simple
`,
			stageCfg: kubeConfigPkg.K8sPrimaryRolloutStageOptions{
				AddVariantLabelToSelector: true,
			},
			expectYAML: `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: simple
spec:
  selector:
    matchLabels:
      app: simple
      pipecd.dev/variant: primary
  template:
    metadata:
      labels:
        app: simple
        pipecd.dev/variant: primary
`,
		},
		{
			name: "create service with suffix",
			inputYAML: `
apiVersion: v1
kind: Service
metadata:
  name: my-service
spec:
  selector:
    app: my-app
  type: NodePort
  ports:
    - port: 80
      targetPort: 8080
`,
			stageCfg: kubeConfigPkg.K8sPrimaryRolloutStageOptions{
				CreateService: true,
				Suffix:        "custom",
			},
			expectYAML: `
apiVersion: v1
kind: Service
metadata:
  name: my-service
spec:
  selector:
    app: my-app
  type: NodePort
  ports:
    - port: 80
      targetPort: 8080
---
apiVersion: v1
kind: Service
metadata:
  name: my-service-custom
spec:
  selector:
    app: my-app
    pipecd.dev/variant: primary
  type: ClusterIP
  ports:
    - port: 80
      targetPort: 8080
`,
		},
		{
			name: "error when service is not matched with config",
			inputYAML: `
apiVersion: v1
kind: Service
metadata:
  name: other-service
spec:
  selector:
    app: my-app
  type: NodePort
  ports:
    - port: 80
      targetPort: 8080
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: simple
spec:
  selector:
    matchLabels:
      app: simple
  template:
    metadata:
      labels:
        app: simple
`,
			stageCfg: kubeConfigPkg.K8sPrimaryRolloutStageOptions{
				CreateService: true,
			},
			expectErr: true,
		},
		{
			name: "error when no service found for CreateService",
			inputYAML: `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: simple
spec:
  selector:
    matchLabels:
      app: simple
  template:
    metadata:
      labels:
        app: simple
`,
			stageCfg: kubeConfigPkg.K8sPrimaryRolloutStageOptions{
				CreateService: true,
			},
			expectErr: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			manifests, err := provider.ParseManifests(tc.inputYAML)
			require.NoError(t, err)

			result, err := generatePrimaryManifests(appCfg, manifests, tc.stageCfg, variantLabel, variant)
			if tc.expectErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			expects, err := provider.ParseManifests(tc.expectYAML)
			require.NoError(t, err)
			require.Equal(t, len(expects), len(result))

			for i := range expects {
				// Compare manifests by converting to structured objects.
				switch expects[i].Kind() {
				case "Deployment":
					var want, got appsv1.Deployment
					err := expects[i].ConvertToStructuredObject(&want)
					require.NoError(t, err)
					err = result[i].ConvertToStructuredObject(&got)
					require.NoError(t, err)
					assert.Equal(t, want, got)
				case "Service":
					var want, got corev1.Service
					err := expects[i].ConvertToStructuredObject(&want)
					require.NoError(t, err)
					err = result[i].ConvertToStructuredObject(&got)
					require.NoError(t, err)
					assert.Equal(t, want, got)
				}
			}
		})
	}
}
