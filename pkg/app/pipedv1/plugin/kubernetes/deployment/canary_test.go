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
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	schema "k8s.io/apimachinery/pkg/runtime/schema"

	kubeconfig "github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/config"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/provider"
	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
	"github.com/pipe-cd/piped-plugin-sdk-go/logpersister/logpersistertest"
	"github.com/pipe-cd/piped-plugin-sdk-go/toolregistry/toolregistrytest"
)

func TestPlugin_executeK8sCanaryRolloutStage_withCreateService(t *testing.T) {
	t.Parallel()

	ctx := t.Context()

	// initialize tool registry
	testRegistry := toolregistrytest.NewTestToolRegistry(t)

	configDir := filepath.Join("testdata", "canary_rollout_with_create_service")

	// read the application config from the example file
	appCfg := sdk.LoadApplicationConfigForTest[kubeconfig.KubernetesApplicationSpec](t, filepath.Join(configDir, "app.pipecd.yaml"), "kubernetes")

	input := &sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec]{
		Request: sdk.ExecuteStageRequest[kubeconfig.KubernetesApplicationSpec]{
			StageName:   "K8S_CANARY_ROLLOUT",
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
		Client: sdk.NewClient(nil, "kubernetes", "", "", logpersistertest.NewTestLogPersister(t), testRegistry),
		Logger: zaptest.NewLogger(t),
	}

	// initialize deploy target config and dynamic client for assertions with envtest
	dtConfig, dynamicClient := setupTestDeployTargetConfigAndDynamicClient(t)

	plugin := &Plugin{}

	status := plugin.executeK8sCanaryRolloutStage(ctx, input, []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]{
		{
			Name:   "default",
			Config: *dtConfig,
		},
	})

	assert.Equal(t, sdk.StageStatusSuccess, status)

	// Assert that Deployment and Service resources are created and have expected labels/annotations.
	deploymentRes := schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}
	deployment, err := dynamicClient.Resource(deploymentRes).Namespace("default").Get(ctx, "simple-canary", metav1.GetOptions{})
	require.NoError(t, err)
	assert.Equal(t, "simple-canary", deployment.GetName())
	assert.Equal(t, "simple", deployment.GetLabels()["app"])
	assert.Equal(t, "canary", deployment.GetLabels()["pipecd.dev/variant"])
	assert.Equal(t, "canary", deployment.GetAnnotations()["pipecd.dev/variant"])

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
	assert.Equal(t, "apps:Deployment::simple-canary", deployment.GetAnnotations()["pipecd.dev/resource-key"])

	serviceRes := schema.GroupVersionResource{Group: "", Version: "v1", Resource: "services"}
	service, err := dynamicClient.Resource(serviceRes).Namespace("default").Get(ctx, "simple-canary", metav1.GetOptions{})
	require.NoError(t, err)
	assert.Equal(t, "simple-canary", service.GetName())

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
	assert.Equal(t, ":Service::simple-canary", service.GetAnnotations()["pipecd.dev/resource-key"])

	// Check Service selector and ports
	selector, found, err := unstructured.NestedStringMap(service.Object, "spec", "selector")
	require.NoError(t, err)
	require.True(t, found)
	assert.Equal(t, map[string]string{"app": "simple", "pipecd.dev/variant": "canary"}, selector)
	ports, found, err := unstructured.NestedSlice(service.Object, "spec", "ports")
	require.NoError(t, err)
	require.True(t, found)
	require.Len(t, ports, 1)
	port := ports[0].(map[string]any)
	assert.Equal(t, int64(9085), port["port"])
	assert.Equal(t, int64(9085), port["targetPort"])
	assert.Equal(t, "TCP", port["protocol"])
}

func TestPlugin_executeK8sCanaryRolloutStage_withoutCreateService(t *testing.T) {
	t.Parallel()

	ctx := t.Context()

	// initialize tool registry
	testRegistry := toolregistrytest.NewTestToolRegistry(t)

	configDir := filepath.Join("testdata", "canary_rollout_without_create_service")

	// read the application config from the example file
	appCfg := sdk.LoadApplicationConfigForTest[kubeconfig.KubernetesApplicationSpec](t, filepath.Join(configDir, "app.pipecd.yaml"), "kubernetes")

	input := &sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec]{
		Request: sdk.ExecuteStageRequest[kubeconfig.KubernetesApplicationSpec]{
			StageName:   "K8S_CANARY_ROLLOUT",
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
		Client: sdk.NewClient(nil, "kubernetes", "", "", logpersistertest.NewTestLogPersister(t), testRegistry),
		Logger: zaptest.NewLogger(t),
	}

	// initialize deploy target config and dynamic client for assertions with envtest
	dtConfig, dynamicClient := setupTestDeployTargetConfigAndDynamicClient(t)

	plugin := &Plugin{}

	status := plugin.executeK8sCanaryRolloutStage(ctx, input, []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]{
		{
			Name:   "default",
			Config: *dtConfig,
		},
	})

	assert.Equal(t, sdk.StageStatusSuccess, status)

	// Assert that Deployment and Service resources are created and have expected labels/annotations.
	deploymentRes := schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}
	deployment, err := dynamicClient.Resource(deploymentRes).Namespace("default").Get(ctx, "simple-canary", metav1.GetOptions{})
	require.NoError(t, err)
	assert.Equal(t, "simple-canary", deployment.GetName())
	assert.Equal(t, "simple", deployment.GetLabels()["app"])
	assert.Equal(t, "canary", deployment.GetLabels()["pipecd.dev/variant"])
	assert.Equal(t, "canary", deployment.GetAnnotations()["pipecd.dev/variant"])

	serviceRes := schema.GroupVersionResource{Group: "", Version: "v1", Resource: "services"}
	_, err = dynamicClient.Resource(serviceRes).Namespace("default").Get(ctx, "simple-canary", metav1.GetOptions{})
	require.Error(t, err)
	assert.True(t, k8serrors.IsNotFound(err))
}

func Test_findConfigMapManifests(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		manifests []provider.Manifest
		want      []provider.Manifest
	}{
		{
			name: "found ConfigMap",
			manifests: mustParseManifests(t, strings.TrimSpace(`
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
spec:
  template:
    spec:
      containers:
      - name: nginx
        image: nginx:1.19.3
---
apiVersion: v1
kind: Service
metadata:
  name: nginx-service
spec:
  selector:
    app: nginx
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: nginx-configmap
data:
  conf: hoge
`)),
			want: mustParseManifests(t, strings.TrimSpace(`
apiVersion: v1
kind: ConfigMap
metadata:
  name: nginx-configmap
data:
  conf: hoge
`)),
		},
		{
			name: "no match",
			manifests: mustParseManifests(t, strings.TrimSpace(`
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
spec:
  template:
    spec:
      containers:
      - name: nginx
        image: nginx:1.19.3
---
apiVersion: v1
kind: Service
metadata:
  name: nginx-service
spec:
  selector:
    app: nginx
`)),
			want: []provider.Manifest{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := findConfigMapManifests(tt.manifests)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_findSecretManifests(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		manifests []provider.Manifest
		want      []provider.Manifest
	}{
		{
			name: "found Secret",
			manifests: mustParseManifests(t, strings.TrimSpace(`
apiVersion: v1
kind: Secret
metadata:
  name: nginx-secret
data:
  password: dGVzdA==
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: nginx-configmap
data:
  conf: hoge
`)),
			want: mustParseManifests(t, strings.TrimSpace(`
apiVersion: v1
kind: Secret
metadata:
  name: nginx-secret
data:
  password: dGVzdA==
`)),
		},
		{
			name: "no match",
			manifests: mustParseManifests(t, strings.TrimSpace(`
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
spec:
  template:
    spec:
      containers:
      - name: nginx
        image: nginx:1.19.3
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: nginx-configmap
data:
  conf: hoge
`)),
			want: []provider.Manifest{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := findSecretManifests(tt.manifests)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_patchManifest(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name          string
		manifests     string
		patch         kubeconfig.K8sResourcePatch
		expectedError error
	}{
		{
			name:      "one op",
			manifests: "testdata/patch_manifest/patch_configmap.yaml",
			patch: kubeconfig.K8sResourcePatch{
				Ops: []kubeconfig.K8sResourcePatchOp{
					{
						Op:    kubeconfig.K8sResourcePatchOpYAMLReplace,
						Path:  "$.data.key1",
						Value: "value-1",
					},
				},
			},
		},
		{
			name:      "multi ops",
			manifests: "testdata/patch_manifest/patch_configmap_multi_ops.yaml",
			patch: kubeconfig.K8sResourcePatch{
				Ops: []kubeconfig.K8sResourcePatchOp{
					{
						Op:    kubeconfig.K8sResourcePatchOpYAMLReplace,
						Path:  "$.data.key1",
						Value: "value-1",
					},
					{
						Op:    kubeconfig.K8sResourcePatchOpYAMLReplace,
						Path:  "$.data.key2",
						Value: "value-2",
					},
				},
			},
		},
		{
			name:      "one op with a given field",
			manifests: "testdata/patch_manifest/patch_configmap_field.yaml",
			patch: kubeconfig.K8sResourcePatch{
				Target: kubeconfig.K8sResourcePatchTarget{
					DocumentRoot: "$.data.envoy-config",
				},
				Ops: []kubeconfig.K8sResourcePatchOp{
					{
						Op:    kubeconfig.K8sResourcePatchOpYAMLReplace,
						Path:  "$.admin.address.socket_address.port_value",
						Value: "9096",
					},
				},
			},
		},
		{
			name:      "multi ops with a given field",
			manifests: "testdata/patch_manifest/patch_configmap_field_multi_ops.yaml",
			patch: kubeconfig.K8sResourcePatch{
				Target: kubeconfig.K8sResourcePatchTarget{
					DocumentRoot: "$.data.envoy-config",
				},
				Ops: []kubeconfig.K8sResourcePatchOp{
					{
						Op:    kubeconfig.K8sResourcePatchOpYAMLReplace,
						Path:  "$.admin.address.socket_address.port_value",
						Value: "19095",
					},
					{
						Op:    kubeconfig.K8sResourcePatchOpYAMLReplace,
						Path:  "$.static_resources.clusters[1].load_assignment.endpoints[0].lb_endpoints[0].endpoint.address.socket_address.port_value",
						Value: "19081",
					},
					{
						Op:    kubeconfig.K8sResourcePatchOpYAMLReplace,
						Path:  "$.static_resources.clusters[1].type",
						Value: "DNS",
					},
				},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			manifests, err := provider.LoadManifestsFromYAMLFile(tc.manifests)
			require.NoError(t, err)

			if tc.expectedError == nil {
				require.Equal(t, 2, len(manifests))
			} else {
				require.Equal(t, 1, len(manifests))
			}

			got, err := patchManifest(manifests[0], tc.patch)
			require.Equal(t, tc.expectedError, err)

			expectedBytes, err := manifests[1].YamlBytes()
			require.NoError(t, err)

			gotBytes, err := got.YamlBytes()
			require.NoError(t, err)

			if tc.expectedError == nil {
				assert.Equal(t, string(expectedBytes), string(gotBytes))
			}
		})
	}
}

func Test_patchManifests(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name          string
		manifests     []provider.Manifest
		patches       []kubeconfig.K8sResourcePatch
		expected      []provider.Manifest
		expectedError error
	}{
		{
			name: "no patches",
			manifests: mustParseManifests(t, strings.TrimSpace(`
apiVersion: apps/v1
kind: Deployment
metadata:
  name: deployment-1
spec:
  template:
    spec:
      containers:
      - name: nginx
        image: nginx:1.19.3
        env:
        - name: VALUE
          value: none
`)),
			patches: []kubeconfig.K8sResourcePatch{},
			expected: mustParseManifests(t, strings.TrimSpace(`
apiVersion: apps/v1
kind: Deployment
metadata:
  name: deployment-1
spec:
  template:
    spec:
      containers:
      - name: nginx
        image: nginx:1.19.3
        env:
        - name: VALUE
          value: none
`)),
		},
		{
			name: "no manifest for the given patch",
			manifests: mustParseManifests(t, strings.TrimSpace(`
apiVersion: apps/v1
kind: Deployment
metadata:
  name: deployment-1
spec:
  template:
    spec:
      containers:
      - name: nginx
        image: nginx:1.19.3
        env:
        - name: VALUE
          value: none
`)),
			patches: []kubeconfig.K8sResourcePatch{
				{
					Target: kubeconfig.K8sResourcePatchTarget{
						K8sResourceReference: kubeconfig.K8sResourceReference{
							Kind: "Deployment",
							Name: "deployment-2",
						},
					},
				},
			},
			expectedError: errors.New("no manifest matches the given patch: kind=Deployment, name=deployment-2"),
		},
		{
			name: "multiple patches",
			manifests: mustParseManifests(t, strings.TrimSpace(`
apiVersion: apps/v1
kind: Deployment
metadata:
  name: deployment-1
spec:
  template:
    spec:
      containers:
      - name: nginx
        image: nginx:1.19.3
        env:
        - name: VALUE
          value: none
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: config
data:
  key1: value-1
  key2: value-2
`)),
			patches: []kubeconfig.K8sResourcePatch{
				{
					Target: kubeconfig.K8sResourcePatchTarget{
						K8sResourceReference: kubeconfig.K8sResourceReference{
							Kind: "Deployment",
							Name: "deployment-1",
						},
					},
					Ops: []kubeconfig.K8sResourcePatchOp{
						{
							Op:    kubeconfig.K8sResourcePatchOpYAMLReplace,
							Path:  "$.spec.template.spec.containers[0].env[0].value",
							Value: "patched",
						},
					},
				},
				{
					Target: kubeconfig.K8sResourcePatchTarget{
						K8sResourceReference: kubeconfig.K8sResourceReference{
							Kind: "ConfigMap",
							Name: "config",
						},
					},
					Ops: []kubeconfig.K8sResourcePatchOp{
						{
							Op:    kubeconfig.K8sResourcePatchOpYAMLReplace,
							Path:  "$.data.key1",
							Value: "patched",
						},
						{
							Op:    kubeconfig.K8sResourcePatchOpYAMLReplace,
							Path:  "$.data.key2",
							Value: "patched",
						},
					},
				},
			},
			expected: mustParseManifests(t, strings.TrimSpace(`
apiVersion: apps/v1
kind: Deployment
metadata:
  name: deployment-1
spec:
  template:
    spec:
      containers:
      - name: nginx
        image: nginx:1.19.3
        env:
        - name: VALUE
          value: patched
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: config
data:
  key1: patched
  key2: patched
`)),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := patchManifests(tc.manifests, tc.patches, patchManifest)
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expected, got)
		})
	}
}
