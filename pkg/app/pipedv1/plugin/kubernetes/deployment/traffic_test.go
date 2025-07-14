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
