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
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
	"github.com/pipe-cd/piped-plugin-sdk-go/logpersister/logpersistertest"
	"github.com/pipe-cd/piped-plugin-sdk-go/toolregistry/toolregistrytest"
	"github.com/pipe-cd/piped-plugin-sdk-go/unit"

	kubeconfig "github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/config"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/provider"
)

type trafficRoutingTestCase struct {
	name            string
	testdataDir     string
	stageCfg        kubeconfig.K8sTrafficRoutingStageOptions
	shouldApplySync bool
	expectedStatus  sdk.StageStatus
	verifyFunc      func(t *testing.T, dynamicClient dynamic.Interface)
}

// setupTrafficRoutingTest initializes common test components
func setupTrafficRoutingTest(t *testing.T, tc trafficRoutingTestCase) (
	input *sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec],
	dtConfig *kubeconfig.KubernetesDeployTargetConfig,
	dynamicClient dynamic.Interface,
) {
	t.Helper()

	// Initialize tool registry
	testRegistry := toolregistrytest.NewTestToolRegistry(t)

	// Read the application config from the testdata file
	appCfg := sdk.LoadApplicationConfigForTest[kubeconfig.KubernetesApplicationSpec](t, filepath.Join("testdata", tc.testdataDir, "app.pipecd.yaml"), "kubernetes")

	// Prepare stage config
	stageCfgBytes, err := json.Marshal(tc.stageCfg)
	require.NoError(t, err)

	// Prepare the input
	input = &sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec]{
		Request: sdk.ExecuteStageRequest[kubeconfig.KubernetesApplicationSpec]{
			StageName:   "K8S_TRAFFIC_ROUTING",
			StageConfig: stageCfgBytes,
			TargetDeploymentSource: sdk.DeploymentSource[kubeconfig.KubernetesApplicationSpec]{
				ApplicationDirectory:      filepath.Join("testdata", tc.testdataDir),
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

	// Initialize deploy target config and dynamic client
	dtConfig, dynamicClient = setupTestDeployTargetConfigAndDynamicClient(t)

	return input, dtConfig, dynamicClient
}

// applyServiceByK8sSync executes K8S_SYNC stage to apply the service
func applyServiceByK8sSync(t *testing.T, ctx context.Context, testdataDir string, dtConfig *kubeconfig.KubernetesDeployTargetConfig) {
	t.Helper()

	testRegistry := toolregistrytest.NewTestToolRegistry(t)
	applyCfg := sdk.LoadApplicationConfigForTest[kubeconfig.KubernetesApplicationSpec](t, filepath.Join("testdata", testdataDir, "app.pipecd.yaml"), "kubernetes")
	applyInput := &sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec]{
		Request: sdk.ExecuteStageRequest[kubeconfig.KubernetesApplicationSpec]{
			StageName:   "K8S_SYNC",
			StageConfig: []byte(``),
			TargetDeploymentSource: sdk.DeploymentSource[kubeconfig.KubernetesApplicationSpec]{
				ApplicationDirectory:      filepath.Join("testdata", testdataDir),
				CommitHash:                "0123456789",
				ApplicationConfig:         applyCfg,
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
	status := plugin.executeK8sSyncStage(ctx, applyInput, []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]{
		{
			Name:   "default",
			Config: *dtConfig,
		},
	})
	require.Equal(t, sdk.StageStatusSuccess, status)
}

// verifyServiceSelector checks if the service selector has the expected variant
func verifyServiceSelector(t *testing.T, dynamicClient dynamic.Interface, serviceName, expectedVariant, variantLabel string) {
	t.Helper()

	service, err := dynamicClient.Resource(schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "services",
	}).Namespace("default").Get(t.Context(), serviceName, metav1.GetOptions{})
	require.NoError(t, err)

	selector := service.Object["spec"].(map[string]interface{})["selector"].(map[string]interface{})
	assert.Equal(t, expectedVariant, selector[variantLabel])
}

func TestPlugin_executeK8sTrafficRoutingStagePodSelector(t *testing.T) {
	t.Parallel()

	testCases := []trafficRoutingTestCase{
		{
			name:        "route to primary",
			testdataDir: "traffic_routing_pod_selector",
			stageCfg: kubeconfig.K8sTrafficRoutingStageOptions{
				All: "primary",
			},
			shouldApplySync: true,
			expectedStatus:  sdk.StageStatusSuccess,
			verifyFunc: func(t *testing.T, dynamicClient dynamic.Interface) {
				verifyServiceSelector(t, dynamicClient, "traffic-test", "primary", "pipecd.dev/variant")
			},
		},
		{
			name:        "route to canary",
			testdataDir: "traffic_routing_pod_selector",
			stageCfg: kubeconfig.K8sTrafficRoutingStageOptions{
				All: "canary",
			},
			shouldApplySync: true,
			expectedStatus:  sdk.StageStatusSuccess,
			verifyFunc: func(t *testing.T, dynamicClient dynamic.Interface) {
				verifyServiceSelector(t, dynamicClient, "traffic-test", "canary", "pipecd.dev/variant")
			},
		},
		{
			name:        "50-50 split not supported",
			testdataDir: "traffic_routing_pod_selector",
			stageCfg: kubeconfig.K8sTrafficRoutingStageOptions{
				Primary: unit.Percentage{Number: 50},
				Canary:  unit.Percentage{Number: 50},
			},
			shouldApplySync: false,
			expectedStatus:  sdk.StageStatusFailure,
		},
		{
			name:        "0-0 split not supported",
			testdataDir: "traffic_routing_pod_selector",
			stageCfg: kubeconfig.K8sTrafficRoutingStageOptions{
				Primary: unit.Percentage{Number: 0},
				Canary:  unit.Percentage{Number: 0},
			},
			shouldApplySync: false,
			expectedStatus:  sdk.StageStatusFailure,
		},
		{
			name:        "baseline not supported",
			testdataDir: "traffic_routing_pod_selector",
			stageCfg: kubeconfig.K8sTrafficRoutingStageOptions{
				Baseline: unit.Percentage{Number: 100},
			},
			shouldApplySync: false,
			expectedStatus:  sdk.StageStatusFailure,
		},
		{
			name:        "no service",
			testdataDir: "traffic_routing_no_service",
			stageCfg: kubeconfig.K8sTrafficRoutingStageOptions{
				All: "primary",
			},
			shouldApplySync: false,
			expectedStatus:  sdk.StageStatusFailure,
		},
		{
			name:        "missing variant label",
			testdataDir: "traffic_routing_missing_variant",
			stageCfg: kubeconfig.K8sTrafficRoutingStageOptions{
				All: "primary",
			},
			shouldApplySync: true,
			expectedStatus:  sdk.StageStatusFailure,
			verifyFunc: func(t *testing.T, dynamicClient dynamic.Interface) {
				// Verify service was created by K8S_SYNC stage
				_, err := dynamicClient.Resource(schema.GroupVersionResource{
					Group:    "",
					Version:  "v1",
					Resource: "services",
				}).Namespace("default").Get(t.Context(), "traffic-test", metav1.GetOptions{})
				require.NoError(t, err)
			},
		},
		{
			name:        "wrong variant value",
			testdataDir: "traffic_routing_wrong_variant",
			stageCfg: kubeconfig.K8sTrafficRoutingStageOptions{
				All: "primary",
			},
			shouldApplySync: true,
			expectedStatus:  sdk.StageStatusFailure,
			verifyFunc: func(t *testing.T, dynamicClient dynamic.Interface) {
				// Verify service was created by K8S_SYNC stage
				_, err := dynamicClient.Resource(schema.GroupVersionResource{
					Group:    "",
					Version:  "v1",
					Resource: "services",
				}).Namespace("default").Get(t.Context(), "traffic-test", metav1.GetOptions{})
				require.NoError(t, err)
			},
		},
		{
			name:        "custom variant label",
			testdataDir: "traffic_routing_custom_variant",
			stageCfg: kubeconfig.K8sTrafficRoutingStageOptions{
				All: "primary",
			},
			shouldApplySync: true,
			expectedStatus:  sdk.StageStatusSuccess,
			verifyFunc: func(t *testing.T, dynamicClient dynamic.Interface) {
				verifyServiceSelector(t, dynamicClient, "traffic-test", "main", "my-custom/variant")
			},
		},
		{
			name:        "multiple services",
			testdataDir: "traffic_routing_multiple_services",
			stageCfg: kubeconfig.K8sTrafficRoutingStageOptions{
				All: "primary",
			},
			shouldApplySync: true,
			expectedStatus:  sdk.StageStatusSuccess,
			verifyFunc: func(t *testing.T, dynamicClient dynamic.Interface) {
				// Verify only the first Service selector was updated
				verifyServiceSelector(t, dynamicClient, "traffic-test-1", "primary", "pipecd.dev/variant")

				// Second service should remain unchanged
				service2, err := dynamicClient.Resource(schema.GroupVersionResource{
					Group:    "",
					Version:  "v1",
					Resource: "services",
				}).Namespace("default").Get(t.Context(), "traffic-test-2", metav1.GetOptions{})
				require.NoError(t, err)
				selector2 := service2.Object["spec"].(map[string]interface{})["selector"].(map[string]interface{})
				assert.Equal(t, "canary", selector2["pipecd.dev/variant"])
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			// Setup test components
			input, dtConfig, dynamicClient := setupTrafficRoutingTest(t, tc)

			// Apply service if needed
			if tc.shouldApplySync {
				applyServiceByK8sSync(t, ctx, tc.testdataDir, dtConfig)
			}

			// Execute traffic routing
			plugin := &Plugin{}
			appCfg := input.Request.TargetDeploymentSource.ApplicationConfig

			deployTargets := []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]{
				{
					Name:   "default",
					Config: *dtConfig,
				},
			}

			status := plugin.executeK8sTrafficRoutingStagePodSelector(ctx, input, deployTargets, appCfg)
			assert.Equal(t, tc.expectedStatus, status)

			// Run verification if provided
			if tc.verifyFunc != nil {
				tc.verifyFunc(t, dynamicClient)
			}
		})
	}
}

// This test assumes that the parsing of the stage config is done before the assertion of the deploy target.
// If the order is changed, this test will not work.
func TestPlugin_executeK8sTrafficRoutingStagePodSelector_InvalidInputs(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		stageCfg []byte
	}{
		{
			name:     "empty stage config",
			stageCfg: []byte(``),
		},
		{
			name:     "invalid stage config",
			stageCfg: []byte(`invalid`),
		},
		{
			name:     "valid stage config but no deploy target",
			stageCfg: []byte(`{"all": "primary"}`),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			appCfg := sdk.LoadApplicationConfigForTest[kubeconfig.KubernetesApplicationSpec](t, filepath.Join("testdata", "traffic_routing_pod_selector", "app.pipecd.yaml"), "kubernetes")

			plugin := &Plugin{}
			status := plugin.executeK8sTrafficRoutingStagePodSelector(t.Context(), &sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec]{
				Request: sdk.ExecuteStageRequest[kubeconfig.KubernetesApplicationSpec]{
					StageConfig: tc.stageCfg,
					TargetDeploymentSource: sdk.DeploymentSource[kubeconfig.KubernetesApplicationSpec]{
						ApplicationDirectory:      filepath.Join("testdata", "traffic_routing_pod_selector"),
						CommitHash:                "0123456789",
						ApplicationConfig:         appCfg,
						ApplicationConfigFilename: "app.pipecd.yaml",
					},
				},
				Client: sdk.NewClient(nil, "kubernetes", "app-id", "stage-id", logpersistertest.NewTestLogPersister(t), toolregistrytest.NewTestToolRegistry(t)),
				Logger: zaptest.NewLogger(t),
			}, []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]{}, appCfg)
			assert.Equal(t, sdk.StageStatusFailure, status)
		})
	}
}

func Test_findIstioVirtualServiceManifests(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		manifestsYAML string
		ref           kubeconfig.K8sResourceReference
		wantCount     int
		wantNames     []string
		wantErr       bool
		errMsg        string
	}{
		{
			name: "finds matching VirtualService by name",
			manifestsYAML: `
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: test-vs
spec:
  hosts:
  - test.example.com
---
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: other-vs
spec:
  hosts:
  - other.example.com
---
apiVersion: v1
kind: Service
metadata:
  name: test-service
`,
			ref: kubeconfig.K8sResourceReference{
				Kind: "VirtualService",
				Name: "test-vs",
			},
			wantCount: 1,
			wantNames: []string{"test-vs"},
			wantErr:   false,
		},
		{
			name: "finds all VirtualServices when name is empty",
			manifestsYAML: `
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: vs1
spec:
  hosts:
  - vs1.example.com
---
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: vs2
spec:
  hosts:
  - vs2.example.com
---
apiVersion: v1
kind: Service
metadata:
  name: test-service
`,
			ref: kubeconfig.K8sResourceReference{
				Kind: "VirtualService",
				Name: "",
			},
			wantCount: 2,
			wantNames: []string{"vs1", "vs2"},
			wantErr:   false,
		},
		{
			name: "finds all VirtualServices when kind is empty",
			manifestsYAML: `
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: vs1
spec:
  hosts:
  - vs1.example.com
---
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: vs2
spec:
  hosts:
  - vs2.example.com
---
apiVersion: v1
kind: Service
metadata:
  name: test-service
`,
			ref: kubeconfig.K8sResourceReference{
				Kind: "",
				Name: "",
			},
			wantCount: 2,
			wantNames: []string{"vs1", "vs2"},
			wantErr:   false,
		},
		{
			name: "returns empty when no VirtualServices found",
			manifestsYAML: `
apiVersion: v1
kind: Service
metadata:
  name: test-service
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: test-deployment
spec:
  replicas: 1
  selector:
    matchLabels:
      app: test
  template:
    metadata:
      labels:
        app: test
    spec:
      containers:
      - name: test
        image: nginx
`,
			ref: kubeconfig.K8sResourceReference{
				Kind: "VirtualService",
				Name: "test-vs",
			},
			wantCount: 0,
			wantNames: []string{},
			wantErr:   false,
		},
		{
			name: "returns empty when no matching name found",
			manifestsYAML: `
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: other-vs
spec:
  hosts:
  - other.example.com
`,
			ref: kubeconfig.K8sResourceReference{
				Kind: "VirtualService",
				Name: "test-vs",
			},
			wantCount: 0,
			wantNames: []string{},
			wantErr:   false,
		},
		{
			name: "filters out manifests with wrong group",
			manifestsYAML: `
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: istio-vs
spec:
  hosts:
  - istio.example.com
---
apiVersion: networking.k8s.io/v1
kind: VirtualService
metadata:
  name: k8s-ingress
spec:
  rules:
  - host: k8s.example.com
`,
			ref: kubeconfig.K8sResourceReference{
				Kind: "VirtualService",
				Name: "",
			},
			wantCount: 1,
			wantNames: []string{"istio-vs"},
			wantErr:   false,
		},
		{
			name: "filters out manifests with wrong kind",
			manifestsYAML: `
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: test-vs
spec:
  hosts:
  - test.example.com
---
apiVersion: networking.istio.io/v1beta1
kind: DestinationRule
metadata:
  name: test-dr
spec:
  host: test.example.com
`,
			ref: kubeconfig.K8sResourceReference{
				Kind: "VirtualService",
				Name: "",
			},
			wantCount: 1,
			wantNames: []string{"test-vs"},
			wantErr:   false,
		},
		{
			name: "returns error for invalid kind",
			manifestsYAML: `
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: test-vs
spec:
  hosts:
  - test.example.com
`,
			ref: kubeconfig.K8sResourceReference{
				Kind: "DestinationRule",
				Name: "test-dr",
			},
			wantCount: 0,
			wantNames: []string{},
			wantErr:   true,
			errMsg:    `support only "VirtualService" kind for VirtualService reference`,
		},
		{
			name:          "returns empty for empty manifest list",
			manifestsYAML: "",
			ref: kubeconfig.K8sResourceReference{
				Kind: "VirtualService",
				Name: "test-vs",
			},
			wantCount: 0,
			wantNames: []string{},
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var manifests []provider.Manifest
			if tt.manifestsYAML != "" {
				manifests = mustParseManifests(t, tt.manifestsYAML)
			}

			got, err := findIstioVirtualServiceManifests(manifests, tt.ref)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.wantCount, len(got), "Expected %d manifests, got %d", tt.wantCount, len(got))

			// Verify the names of the returned manifests
			gotNames := make([]string, len(got))
			for i, manifest := range got {
				gotNames[i] = manifest.Name()
			}
			assert.ElementsMatch(t, tt.wantNames, gotNames, "Expected manifest names to match")

			// Verify that all returned manifests are VirtualServices
			for _, manifest := range got {
				assert.Equal(t, "VirtualService", manifest.Kind())
				assert.Equal(t, "networking.istio.io", manifest.GroupVersionKind().Group)
			}
		})
	}
}

func TestConvertVirtualService(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		yaml      string
		checkFunc func(t *testing.T, vs *virtualService)
	}{
		{
			name: "valid VirtualService with basic HTTP routing",
			yaml: `
apiVersion: networking.istio.io/v1
kind: VirtualService
metadata:
  name: test-virtual-service
  namespace: default
spec:
  hosts:
  - test-service
  http:
  - route:
    - destination:
        host: test-service
        subset: primary
      weight: 100
`,
			checkFunc: func(t *testing.T, vs *virtualService) {
				assert.Equal(t, "networking.istio.io/v1", vs.APIVersion)
				assert.Equal(t, "VirtualService", vs.Kind)
				assert.Equal(t, "test-virtual-service", vs.Name)
				assert.Equal(t, "default", vs.Namespace)
				assert.Equal(t, []string{"test-service"}, vs.Spec.Hosts)
				assert.Len(t, vs.Spec.Http, 1)
				assert.Len(t, vs.Spec.Http[0].Route, 1)
				assert.Equal(t, "test-service", vs.Spec.Http[0].Route[0].Destination.Host)
				assert.Equal(t, "primary", vs.Spec.Http[0].Route[0].Destination.Subset)
				assert.Equal(t, int32(100), vs.Spec.Http[0].Route[0].Weight)
			},
		},
		{
			name: "VirtualService with multiple routes and match conditions",
			yaml: `
apiVersion: networking.istio.io/v1
kind: VirtualService
metadata:
  name: multi-route-vs
  namespace: istio-system
  labels:
    app: test
spec:
  hosts:
  - example.com
  - test-service.default.svc.cluster.local
  gateways:
  - test-gateway
  http:
  - name: primary-route
    match:
    - headers:
        version:
          exact: v1
    route:
    - destination:
        host: test-service
        subset: primary
      weight: 80
    - destination:
        host: test-service
        subset: canary
      weight: 20
  - name: default-route
    route:
    - destination:
        host: test-service
        subset: primary
      weight: 100
`,
			checkFunc: func(t *testing.T, vs *virtualService) {
				assert.Equal(t, "networking.istio.io/v1", vs.APIVersion)
				assert.Equal(t, "VirtualService", vs.Kind)
				assert.Equal(t, "multi-route-vs", vs.Name)
				assert.Equal(t, "istio-system", vs.Namespace)
				assert.Equal(t, "test", vs.Labels["app"])
				assert.Equal(t, []string{"example.com", "test-service.default.svc.cluster.local"}, vs.Spec.Hosts)
				assert.Equal(t, []string{"test-gateway"}, vs.Spec.Gateways)
				assert.Len(t, vs.Spec.Http, 2)

				// Check first HTTP route
				assert.Equal(t, "primary-route", vs.Spec.Http[0].Name)
				assert.Len(t, vs.Spec.Http[0].Match, 1)
				assert.Len(t, vs.Spec.Http[0].Route, 2)
				assert.Equal(t, "test-service", vs.Spec.Http[0].Route[0].Destination.Host)
				assert.Equal(t, "primary", vs.Spec.Http[0].Route[0].Destination.Subset)
				assert.Equal(t, int32(80), vs.Spec.Http[0].Route[0].Weight)
				assert.Equal(t, "canary", vs.Spec.Http[0].Route[1].Destination.Subset)
				assert.Equal(t, int32(20), vs.Spec.Http[0].Route[1].Weight)

				// Check second HTTP route
				assert.Equal(t, "default-route", vs.Spec.Http[1].Name)
				assert.Len(t, vs.Spec.Http[1].Route, 1)
				assert.Equal(t, int32(100), vs.Spec.Http[1].Route[0].Weight)
			},
		},
		{
			name: "VirtualService with TCP routing",
			yaml: `
apiVersion: networking.istio.io/v1
kind: VirtualService
metadata:
  name: tcp-vs
spec:
  hosts:
  - tcp-service
  tcp:
  - match:
    - port: 3306
    route:
    - destination:
        host: mysql-service
        port:
          number: 3306
      weight: 100
`,
			checkFunc: func(t *testing.T, vs *virtualService) {
				assert.Equal(t, "tcp-vs", vs.Name)
				assert.Equal(t, []string{"tcp-service"}, vs.Spec.Hosts)
				assert.Len(t, vs.Spec.Tcp, 1)
				assert.Len(t, vs.Spec.Tcp[0].Route, 1)
				assert.Equal(t, "mysql-service", vs.Spec.Tcp[0].Route[0].Destination.Host)
				assert.Equal(t, uint32(3306), vs.Spec.Tcp[0].Route[0].Destination.Port.Number)
			},
		},
		{
			name: "minimal VirtualService",
			yaml: `
apiVersion: networking.istio.io/v1
kind: VirtualService
metadata:
  name: minimal-vs
spec:
  hosts:
  - minimal-service
`,
			checkFunc: func(t *testing.T, vs *virtualService) {
				assert.Equal(t, "minimal-vs", vs.Name)
				assert.Equal(t, []string{"minimal-service"}, vs.Spec.Hosts)
				assert.Empty(t, vs.Spec.Http)
				assert.Empty(t, vs.Spec.Tcp)
			},
		},
		{
			name: "invalid YAML - not a VirtualService",
			yaml: `
apiVersion: v1
kind: Service
metadata:
  name: test-service
spec:
  selector:
    app: test
  ports:
  - port: 80
    targetPort: 8080
`,
			checkFunc: func(t *testing.T, vs *virtualService) {
				assert.Equal(t, "v1", vs.APIVersion)
				assert.Equal(t, "Service", vs.Kind)
				assert.Equal(t, "test-service", vs.Name)
				// Spec will be empty since Service spec doesn't match VirtualService spec structure
			},
		},
		{
			name: "YAML with invalid structure for VirtualService conversion",
			yaml: `
apiVersion: v1
kind: VirtualService
metadata:
  name: invalid-structure
spec:
  # This will cause conversion errors because these fields don't match the expected VirtualService structure
  invalidField: "invalid value"
  anotherInvalidField:
    nested: "structure that doesn't match Istio VirtualService schema"
`,
			checkFunc: func(t *testing.T, vs *virtualService) {
				assert.Equal(t, "v1", vs.APIVersion)
				assert.Equal(t, "VirtualService", vs.Kind)
				assert.Equal(t, "invalid-structure", vs.Name)
				// The Spec will be mostly empty since the fields don't match VirtualService structure
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			manifests := mustParseManifests(t, tt.yaml)
			require.Len(t, manifests, 1, "test should provide exactly one manifest")

			vs, err := convertVirtualService(manifests[0])

			require.NoError(t, err)
			require.NotNil(t, vs)
			if tt.checkFunc != nil {
				tt.checkFunc(t, vs)
			}
		})
	}
}

func TestVirtualService_toManifest(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		yaml      string
		checkFunc func(t *testing.T, original *virtualService, converted provider.Manifest)
	}{
		{
			name: "basic VirtualService conversion",
			yaml: `
apiVersion: networking.istio.io/v1
kind: VirtualService
metadata:
  name: test-virtual-service
  namespace: default
spec:
  hosts:
  - test-service
  http:
  - route:
    - destination:
        host: test-service
        subset: primary
      weight: 100
`,
			checkFunc: func(t *testing.T, original *virtualService, converted provider.Manifest) {
				// Check basic metadata
				assert.Equal(t, original.Name, converted.Name())
				assert.Equal(t, original.Kind, converted.Kind())
				assert.Equal(t, original.APIVersion, converted.APIVersion())

				// Check namespace using NestedString (metadata.namespace)
				namespace, found, err := converted.NestedString("metadata", "namespace")
				require.NoError(t, err)
				if original.Namespace != "" {
					require.True(t, found)
					assert.Equal(t, original.Namespace, namespace)
				}

				// Verify hosts are preserved using NestedMap approach
				spec, found, err := converted.NestedMap("spec")
				require.NoError(t, err)
				require.True(t, found)

				hosts, ok := spec["hosts"].([]interface{})
				require.True(t, ok, "hosts should be a slice")
				hostStrings := make([]string, len(hosts))
				for i, h := range hosts {
					hostStrings[i] = h.(string)
				}
				assert.Equal(t, original.Spec.Hosts, hostStrings)

				// Verify the converted manifest can be serialized back to the same structure
				var reconverted virtualService
				err = converted.ConvertToStructuredObject(&reconverted)
				require.NoError(t, err)
				assert.Equal(t, original.Name, reconverted.Name)
				assert.Equal(t, original.Namespace, reconverted.Namespace)
				assert.Equal(t, original.Spec.Hosts, reconverted.Spec.Hosts)
			},
		},
		{
			name: "VirtualService with multiple routes",
			yaml: `
apiVersion: networking.istio.io/v1
kind: VirtualService
metadata:
  name: multi-route-vs
  namespace: istio-system
  labels:
    app: test
spec:
  hosts:
  - example.com
  - test-service.default.svc.cluster.local
  gateways:
  - test-gateway
  http:
  - name: primary-route
    route:
    - destination:
        host: test-service
        subset: primary
      weight: 80
    - destination:
        host: test-service
        subset: canary
      weight: 20
  - name: default-route
    route:
    - destination:
        host: test-service
        subset: primary
      weight: 100
`,
			checkFunc: func(t *testing.T, original *virtualService, converted provider.Manifest) {
				// Check metadata preservation
				assert.Equal(t, original.Name, converted.Name())
				assert.Equal(t, original.Kind, converted.Kind())

				// Check namespace using NestedString (metadata.namespace)
				namespace, found, err := converted.NestedString("metadata", "namespace")
				require.NoError(t, err)
				if original.Namespace != "" {
					require.True(t, found)
					assert.Equal(t, original.Namespace, namespace)
				}

				// Check labels are preserved using NestedMap
				metadata, found, err := converted.NestedMap("metadata")
				require.NoError(t, err)
				require.True(t, found)
				if labels, ok := metadata["labels"].(map[string]interface{}); ok {
					assert.Equal(t, "test", labels["app"])
				}

				// Check hosts are preserved using NestedMap
				spec, found, err := converted.NestedMap("spec")
				require.NoError(t, err)
				require.True(t, found)

				hosts, ok := spec["hosts"].([]interface{})
				require.True(t, ok, "hosts should be a slice")
				hostStrings := make([]string, len(hosts))
				for i, h := range hosts {
					hostStrings[i] = h.(string)
				}
				assert.Equal(t, original.Spec.Hosts, hostStrings)

				// Check gateways are preserved using NestedMap
				if gateways, ok := spec["gateways"].([]interface{}); ok {
					gatewayStrings := make([]string, len(gateways))
					for i, g := range gateways {
						gatewayStrings[i] = g.(string)
					}
					assert.Equal(t, original.Spec.Gateways, gatewayStrings)
				}

				// Verify round-trip conversion
				var roundTrip virtualService
				err = converted.ConvertToStructuredObject(&roundTrip)
				require.NoError(t, err)
				assert.Equal(t, original.Spec.Hosts, roundTrip.Spec.Hosts)
				assert.Equal(t, original.Spec.Gateways, roundTrip.Spec.Gateways)
				assert.Len(t, roundTrip.Spec.Http, 2)
			},
		},
		{
			name: "VirtualService with TCP routing",
			yaml: `
apiVersion: networking.istio.io/v1
kind: VirtualService
metadata:
  name: tcp-vs
spec:
  hosts:
  - tcp-service
  tcp:
  - match:
    - port: 3306
    route:
    - destination:
        host: mysql-service
        port:
          number: 3306
      weight: 100
`,
			checkFunc: func(t *testing.T, original *virtualService, converted provider.Manifest) {
				// Check basic conversion
				assert.Equal(t, original.Name, converted.Name())
				assert.Equal(t, original.Kind, converted.Kind())

				// Check hosts preservation using NestedMap
				spec, found, err := converted.NestedMap("spec")
				require.NoError(t, err)
				require.True(t, found)

				hosts, ok := spec["hosts"].([]interface{})
				require.True(t, ok, "hosts should be a slice")
				hostStrings := make([]string, len(hosts))
				for i, h := range hosts {
					hostStrings[i] = h.(string)
				}
				assert.Equal(t, original.Spec.Hosts, hostStrings)

				// Verify TCP routes are preserved through round-trip
				var roundTrip virtualService
				err = converted.ConvertToStructuredObject(&roundTrip)
				require.NoError(t, err)
				assert.Len(t, roundTrip.Spec.Tcp, 1)
				assert.Len(t, roundTrip.Spec.Tcp[0].Route, 1)
				assert.Equal(t, "mysql-service", roundTrip.Spec.Tcp[0].Route[0].Destination.Host)
			},
		},
		{
			name: "minimal VirtualService",
			yaml: `
apiVersion: networking.istio.io/v1
kind: VirtualService
metadata:
  name: minimal-vs
spec:
  hosts:
  - minimal-service
`,
			checkFunc: func(t *testing.T, original *virtualService, converted provider.Manifest) {
				// Check basic conversion for minimal case
				assert.Equal(t, original.Name, converted.Name())
				assert.Equal(t, original.Kind, converted.Kind())

				// Check hosts are preserved using NestedMap
				spec, found, err := converted.NestedMap("spec")
				require.NoError(t, err)
				require.True(t, found)

				hosts, ok := spec["hosts"].([]interface{})
				require.True(t, ok, "hosts should be a slice")
				hostStrings := make([]string, len(hosts))
				for i, h := range hosts {
					hostStrings[i] = h.(string)
				}
				assert.Equal(t, original.Spec.Hosts, hostStrings)

				// Verify round-trip for minimal case
				var roundTrip virtualService
				err = converted.ConvertToStructuredObject(&roundTrip)
				require.NoError(t, err)
				assert.Equal(t, original.Spec.Hosts, roundTrip.Spec.Hosts)
				assert.Empty(t, roundTrip.Spec.Http)
				assert.Empty(t, roundTrip.Spec.Tcp)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Parse YAML to manifest
			manifests := mustParseManifests(t, tt.yaml)
			require.Len(t, manifests, 1, "test should provide exactly one manifest")

			// Convert to virtualService
			vs, err := convertVirtualService(manifests[0])
			require.NoError(t, err)
			require.NotNil(t, vs)

			// Test toManifest conversion
			convertedManifest, err := vs.toManifest()
			require.NoError(t, err)
			require.NotNil(t, convertedManifest)

			// Run verification
			if tt.checkFunc != nil {
				tt.checkFunc(t, vs, convertedManifest)
			}
		})
	}
}
