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

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
	"github.com/pipe-cd/piped-plugin-sdk-go/logpersister/logpersistertest"
	"github.com/pipe-cd/piped-plugin-sdk-go/toolregistry/toolregistrytest"
	"github.com/pipe-cd/piped-plugin-sdk-go/unit"

	kubeconfig "github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/config"
)

func TestPlugin_executeK8sTrafficRoutingStagePodSelector_toPrimary(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	// Initialize tool registry
	testRegistry := toolregistrytest.NewTestToolRegistry(t)

	// Read the application config from the testdata file
	appCfg := sdk.LoadApplicationConfigForTest[kubeconfig.KubernetesApplicationSpec](t, filepath.Join("testdata", "traffic_routing_pod_selector_primary", "app.pipecd.yaml"), "kubernetes")

	// Prepare stage config to route 100% traffic to primary
	stageCfg := kubeconfig.K8sTrafficRoutingStageOptions{
		All: "primary",
	}
	stageCfgBytes, err := json.Marshal(stageCfg)
	require.NoError(t, err)

	// Prepare the input
	input := &sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec]{
		Request: sdk.ExecuteStageRequest[kubeconfig.KubernetesApplicationSpec]{
			StageName:   "K8S_TRAFFIC_ROUTING",
			StageConfig: stageCfgBytes,
			TargetDeploymentSource: sdk.DeploymentSource[kubeconfig.KubernetesApplicationSpec]{
				ApplicationDirectory:      filepath.Join("testdata", "traffic_routing_pod_selector_primary"),
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

	// Initialize deploy target config and dynamic client for assertions with envtest
	dtConfig, dynamicClient := setupTestDeployTargetConfigAndDynamicClient(t)

	// First apply the service with initial canary selector
	applyCfg := sdk.LoadApplicationConfigForTest[kubeconfig.KubernetesApplicationSpec](t, filepath.Join("testdata", "traffic_routing_pod_selector_primary", "app.pipecd.yaml"), "kubernetes")
	applyInput := &sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec]{
		Request: sdk.ExecuteStageRequest[kubeconfig.KubernetesApplicationSpec]{
			StageName:   "K8S_SYNC",
			StageConfig: []byte(``),
			TargetDeploymentSource: sdk.DeploymentSource[kubeconfig.KubernetesApplicationSpec]{
				ApplicationDirectory:      filepath.Join("testdata", "traffic_routing_pod_selector_primary"),
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

	// Now execute traffic routing
	status = plugin.executeK8sTrafficRoutingStagePodSelector(ctx, input, []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]{
		{
			Name:   "default",
			Config: *dtConfig,
		},
	}, appCfg)

	assert.Equal(t, sdk.StageStatusSuccess, status)

	// Verify the Service selector was updated to primary variant
	service, err := dynamicClient.Resource(schema.GroupVersionResource{Group: "", Version: "v1", Resource: "services"}).Namespace("default").Get(context.Background(), "traffic-test", metav1.GetOptions{})
	require.NoError(t, err)

	selector := service.Object["spec"].(map[string]interface{})["selector"].(map[string]interface{})
	assert.Equal(t, "primary", selector["pipecd.dev/variant"])
}

func TestPlugin_executeK8sTrafficRoutingStagePodSelector_toCanary(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	// Initialize tool registry
	testRegistry := toolregistrytest.NewTestToolRegistry(t)

	// Read the application config from the testdata file
	appCfg := sdk.LoadApplicationConfigForTest[kubeconfig.KubernetesApplicationSpec](t, filepath.Join("testdata", "traffic_routing_pod_selector_canary", "app.pipecd.yaml"), "kubernetes")

	// Prepare stage config to route 100% traffic to canary
	stageCfg := kubeconfig.K8sTrafficRoutingStageOptions{
		All: "canary",
	}
	stageCfgBytes, err := json.Marshal(stageCfg)
	require.NoError(t, err)

	// Prepare the input
	input := &sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec]{
		Request: sdk.ExecuteStageRequest[kubeconfig.KubernetesApplicationSpec]{
			StageName:   "K8S_TRAFFIC_ROUTING",
			StageConfig: stageCfgBytes,
			TargetDeploymentSource: sdk.DeploymentSource[kubeconfig.KubernetesApplicationSpec]{
				ApplicationDirectory:      filepath.Join("testdata", "traffic_routing_pod_selector_canary"),
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

	// Initialize deploy target config and dynamic client for assertions with envtest
	dtConfig, dynamicClient := setupTestDeployTargetConfigAndDynamicClient(t)

	// First apply the service with initial primary selector
	applyCfg := sdk.LoadApplicationConfigForTest[kubeconfig.KubernetesApplicationSpec](t, filepath.Join("testdata", "traffic_routing_pod_selector_canary", "app.pipecd.yaml"), "kubernetes")
	applyInput := &sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec]{
		Request: sdk.ExecuteStageRequest[kubeconfig.KubernetesApplicationSpec]{
			StageName:   "K8S_SYNC",
			StageConfig: []byte(``),
			TargetDeploymentSource: sdk.DeploymentSource[kubeconfig.KubernetesApplicationSpec]{
				ApplicationDirectory:      filepath.Join("testdata", "traffic_routing_pod_selector_canary"),
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

	// Now execute traffic routing
	status = plugin.executeK8sTrafficRoutingStagePodSelector(ctx, input, []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]{
		{
			Name:   "default",
			Config: *dtConfig,
		},
	}, appCfg)

	assert.Equal(t, sdk.StageStatusSuccess, status)

	// Verify the Service selector was updated to canary variant
	service, err := dynamicClient.Resource(schema.GroupVersionResource{Group: "", Version: "v1", Resource: "services"}).Namespace("default").Get(context.Background(), "traffic-test", metav1.GetOptions{})
	require.NoError(t, err)

	selector := service.Object["spec"].(map[string]interface{})["selector"].(map[string]interface{})
	assert.Equal(t, "canary", selector["pipecd.dev/variant"])
}

func TestPlugin_executeK8sTrafficRoutingStagePodSelector_invalidPercentages(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	// Initialize tool registry
	testRegistry := toolregistrytest.NewTestToolRegistry(t)

	// Read the application config from the testdata file
	appCfg := sdk.LoadApplicationConfigForTest[kubeconfig.KubernetesApplicationSpec](t, filepath.Join("testdata", "traffic_routing_pod_selector_primary", "app.pipecd.yaml"), "kubernetes")

	testCases := []struct {
		name          string
		stageCfg      kubeconfig.K8sTrafficRoutingStageOptions
		expectedError string
	}{
		{
			name: "50-50 split not supported",
			stageCfg: kubeconfig.K8sTrafficRoutingStageOptions{
				Primary: unit.Percentage{Number: 50},
				Canary:  unit.Percentage{Number: 50},
			},
			expectedError: "PodSelector requires either primary or canary to be 100%",
		},
		{
			name: "80-20 split not supported",
			stageCfg: kubeconfig.K8sTrafficRoutingStageOptions{
				Primary: unit.Percentage{Number: 80},
				Canary:  unit.Percentage{Number: 20},
			},
			expectedError: "PodSelector requires either primary or canary to be 100%",
		},
		{
			name: "0-0 split not supported",
			stageCfg: kubeconfig.K8sTrafficRoutingStageOptions{
				Primary: unit.Percentage{Number: 0},
				Canary:  unit.Percentage{Number: 0},
			},
			expectedError: "PodSelector requires either primary or canary to be 100%",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			stageCfgBytes, err := json.Marshal(tc.stageCfg)
			require.NoError(t, err)

			// Prepare the input
			input := &sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec]{
				Request: sdk.ExecuteStageRequest[kubeconfig.KubernetesApplicationSpec]{
					StageName:   "K8S_TRAFFIC_ROUTING",
					StageConfig: stageCfgBytes,
					TargetDeploymentSource: sdk.DeploymentSource[kubeconfig.KubernetesApplicationSpec]{
						ApplicationDirectory:      filepath.Join("testdata", "traffic_routing_pod_selector_primary"),
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

			// Initialize deploy target config
			dtConfig, _ := setupTestDeployTargetConfigAndDynamicClient(t)

			plugin := &Plugin{}
			status := plugin.executeK8sTrafficRoutingStagePodSelector(ctx, input, []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]{
				{
					Name:   "default",
					Config: *dtConfig,
				},
			}, appCfg)

			assert.Equal(t, sdk.StageStatusFailure, status)
		})
	}
}

func TestPlugin_executeK8sTrafficRoutingStagePodSelector_withBaseline(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	// Initialize tool registry
	testRegistry := toolregistrytest.NewTestToolRegistry(t)

	// Read the application config from the testdata file
	appCfg := sdk.LoadApplicationConfigForTest[kubeconfig.KubernetesApplicationSpec](t, filepath.Join("testdata", "traffic_routing_pod_selector_primary", "app.pipecd.yaml"), "kubernetes")

	// Prepare stage config with baseline percentage
	stageCfg := kubeconfig.K8sTrafficRoutingStageOptions{
		Primary:  unit.Percentage{Number: 50},
		Canary:   unit.Percentage{Number: 30},
		Baseline: unit.Percentage{Number: 20},
	}
	stageCfgBytes, err := json.Marshal(stageCfg)
	require.NoError(t, err)

	// Prepare the input
	input := &sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec]{
		Request: sdk.ExecuteStageRequest[kubeconfig.KubernetesApplicationSpec]{
			StageName:   "K8S_TRAFFIC_ROUTING",
			StageConfig: stageCfgBytes,
			TargetDeploymentSource: sdk.DeploymentSource[kubeconfig.KubernetesApplicationSpec]{
				ApplicationDirectory:      filepath.Join("testdata", "traffic_routing_pod_selector_primary"),
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

	// Initialize deploy target config
	dtConfig, _ := setupTestDeployTargetConfigAndDynamicClient(t)

	plugin := &Plugin{}
	status := plugin.executeK8sTrafficRoutingStagePodSelector(ctx, input, []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]{
		{
			Name:   "default",
			Config: *dtConfig,
		},
	}, appCfg)

	assert.Equal(t, sdk.StageStatusFailure, status)
}

func TestPlugin_executeK8sTrafficRoutingStagePodSelector_noService(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	// Initialize tool registry
	testRegistry := toolregistrytest.NewTestToolRegistry(t)

	// Read the application config from the testdata file
	appCfg := sdk.LoadApplicationConfigForTest[kubeconfig.KubernetesApplicationSpec](t, filepath.Join("testdata", "traffic_routing_no_service", "app.pipecd.yaml"), "kubernetes")

	// Prepare stage config to route 100% traffic to primary
	stageCfg := kubeconfig.K8sTrafficRoutingStageOptions{
		All: "primary",
	}
	stageCfgBytes, err := json.Marshal(stageCfg)
	require.NoError(t, err)

	// Prepare the input
	input := &sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec]{
		Request: sdk.ExecuteStageRequest[kubeconfig.KubernetesApplicationSpec]{
			StageName:   "K8S_TRAFFIC_ROUTING",
			StageConfig: stageCfgBytes,
			TargetDeploymentSource: sdk.DeploymentSource[kubeconfig.KubernetesApplicationSpec]{
				ApplicationDirectory:      filepath.Join("testdata", "traffic_routing_no_service"),
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

	// Initialize deploy target config
	dtConfig, _ := setupTestDeployTargetConfigAndDynamicClient(t)

	plugin := &Plugin{}
	status := plugin.executeK8sTrafficRoutingStagePodSelector(ctx, input, []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]{
		{
			Name:   "default",
			Config: *dtConfig,
		},
	}, appCfg)

	assert.Equal(t, sdk.StageStatusFailure, status)
}

func TestPlugin_executeK8sTrafficRoutingStagePodSelector_missingVariantLabel(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	// Initialize tool registry
	testRegistry := toolregistrytest.NewTestToolRegistry(t)

	// Read the application config from the testdata file
	appCfg := sdk.LoadApplicationConfigForTest[kubeconfig.KubernetesApplicationSpec](t, filepath.Join("testdata", "traffic_routing_missing_variant", "app.pipecd.yaml"), "kubernetes")

	// Prepare stage config to route 100% traffic to primary
	stageCfg := kubeconfig.K8sTrafficRoutingStageOptions{
		All: "primary",
	}
	stageCfgBytes, err := json.Marshal(stageCfg)
	require.NoError(t, err)

	// Prepare the input
	input := &sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec]{
		Request: sdk.ExecuteStageRequest[kubeconfig.KubernetesApplicationSpec]{
			StageName:   "K8S_TRAFFIC_ROUTING",
			StageConfig: stageCfgBytes,
			TargetDeploymentSource: sdk.DeploymentSource[kubeconfig.KubernetesApplicationSpec]{
				ApplicationDirectory:      filepath.Join("testdata", "traffic_routing_missing_variant"),
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

	// Initialize deploy target config and dynamic client for assertions with envtest
	dtConfig, dynamicClient := setupTestDeployTargetConfigAndDynamicClient(t)

	// First apply the service
	applyCfg := sdk.LoadApplicationConfigForTest[kubeconfig.KubernetesApplicationSpec](t, filepath.Join("testdata", "traffic_routing_missing_variant", "app.pipecd.yaml"), "kubernetes")
	applyInput := &sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec]{
		Request: sdk.ExecuteStageRequest[kubeconfig.KubernetesApplicationSpec]{
			StageName:   "K8S_SYNC",
			StageConfig: []byte(``),
			TargetDeploymentSource: sdk.DeploymentSource[kubeconfig.KubernetesApplicationSpec]{
				ApplicationDirectory:      filepath.Join("testdata", "traffic_routing_missing_variant"),
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

	// Verify service was created
	_, err = dynamicClient.Resource(schema.GroupVersionResource{Group: "", Version: "v1", Resource: "services"}).Namespace("default").Get(context.Background(), "traffic-test", metav1.GetOptions{})
	require.NoError(t, err)

	// Now execute traffic routing
	status = plugin.executeK8sTrafficRoutingStagePodSelector(ctx, input, []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]{
		{
			Name:   "default",
			Config: *dtConfig,
		},
	}, appCfg)

	assert.Equal(t, sdk.StageStatusFailure, status)
}

func TestPlugin_executeK8sTrafficRoutingStagePodSelector_wrongVariantValue(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	// Initialize tool registry
	testRegistry := toolregistrytest.NewTestToolRegistry(t)

	// Read the application config from the testdata file
	appCfg := sdk.LoadApplicationConfigForTest[kubeconfig.KubernetesApplicationSpec](t, filepath.Join("testdata", "traffic_routing_wrong_variant", "app.pipecd.yaml"), "kubernetes")

	// Prepare stage config to route 100% traffic to primary
	stageCfg := kubeconfig.K8sTrafficRoutingStageOptions{
		All: "primary",
	}
	stageCfgBytes, err := json.Marshal(stageCfg)
	require.NoError(t, err)

	// Prepare the input
	input := &sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec]{
		Request: sdk.ExecuteStageRequest[kubeconfig.KubernetesApplicationSpec]{
			StageName:   "K8S_TRAFFIC_ROUTING",
			StageConfig: stageCfgBytes,
			TargetDeploymentSource: sdk.DeploymentSource[kubeconfig.KubernetesApplicationSpec]{
				ApplicationDirectory:      filepath.Join("testdata", "traffic_routing_wrong_variant"),
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

	// Initialize deploy target config and dynamic client for assertions with envtest
	dtConfig, dynamicClient := setupTestDeployTargetConfigAndDynamicClient(t)

	// First apply the service
	applyCfg := sdk.LoadApplicationConfigForTest[kubeconfig.KubernetesApplicationSpec](t, filepath.Join("testdata", "traffic_routing_wrong_variant", "app.pipecd.yaml"), "kubernetes")
	applyInput := &sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec]{
		Request: sdk.ExecuteStageRequest[kubeconfig.KubernetesApplicationSpec]{
			StageName:   "K8S_SYNC",
			StageConfig: []byte(``),
			TargetDeploymentSource: sdk.DeploymentSource[kubeconfig.KubernetesApplicationSpec]{
				ApplicationDirectory:      filepath.Join("testdata", "traffic_routing_wrong_variant"),
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

	// Verify service was created
	_, err = dynamicClient.Resource(schema.GroupVersionResource{Group: "", Version: "v1", Resource: "services"}).Namespace("default").Get(context.Background(), "traffic-test", metav1.GetOptions{})
	require.NoError(t, err)

	// Now execute traffic routing
	status = plugin.executeK8sTrafficRoutingStagePodSelector(ctx, input, []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]{
		{
			Name:   "default",
			Config: *dtConfig,
		},
	}, appCfg)

	assert.Equal(t, sdk.StageStatusFailure, status)
}

func TestPlugin_executeK8sTrafficRoutingStagePodSelector_customVariantLabel(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	// Initialize tool registry
	testRegistry := toolregistrytest.NewTestToolRegistry(t)

	// Read the application config from the testdata file
	appCfg := sdk.LoadApplicationConfigForTest[kubeconfig.KubernetesApplicationSpec](t, filepath.Join("testdata", "traffic_routing_custom_variant", "app.pipecd.yaml"), "kubernetes")

	// Prepare stage config to route 100% traffic to custom primary
	stageCfg := kubeconfig.K8sTrafficRoutingStageOptions{
		All: "primary",
	}
	stageCfgBytes, err := json.Marshal(stageCfg)
	require.NoError(t, err)

	// Prepare the input
	input := &sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec]{
		Request: sdk.ExecuteStageRequest[kubeconfig.KubernetesApplicationSpec]{
			StageName:   "K8S_TRAFFIC_ROUTING",
			StageConfig: stageCfgBytes,
			TargetDeploymentSource: sdk.DeploymentSource[kubeconfig.KubernetesApplicationSpec]{
				ApplicationDirectory:      filepath.Join("testdata", "traffic_routing_custom_variant"),
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

	// Initialize deploy target config and dynamic client for assertions with envtest
	dtConfig, dynamicClient := setupTestDeployTargetConfigAndDynamicClient(t)

	// First apply the service with initial variant
	applyCfg := sdk.LoadApplicationConfigForTest[kubeconfig.KubernetesApplicationSpec](t, filepath.Join("testdata", "traffic_routing_custom_variant", "app.pipecd.yaml"), "kubernetes")
	applyInput := &sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec]{
		Request: sdk.ExecuteStageRequest[kubeconfig.KubernetesApplicationSpec]{
			StageName:   "K8S_SYNC",
			StageConfig: []byte(``),
			TargetDeploymentSource: sdk.DeploymentSource[kubeconfig.KubernetesApplicationSpec]{
				ApplicationDirectory:      filepath.Join("testdata", "traffic_routing_custom_variant"),
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

	// Now execute traffic routing
	status = plugin.executeK8sTrafficRoutingStagePodSelector(ctx, input, []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]{
		{
			Name:   "default",
			Config: *dtConfig,
		},
	}, appCfg)

	assert.Equal(t, sdk.StageStatusSuccess, status)

	// Verify the Service selector was updated to custom primary variant
	service, err := dynamicClient.Resource(schema.GroupVersionResource{Group: "", Version: "v1", Resource: "services"}).Namespace("default").Get(context.Background(), "traffic-test", metav1.GetOptions{})
	require.NoError(t, err)

	selector := service.Object["spec"].(map[string]interface{})["selector"].(map[string]interface{})
	assert.Equal(t, "main", selector["my-custom/variant"])
}

func TestPlugin_executeK8sTrafficRoutingStagePodSelector_noDeployTarget(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	// Initialize tool registry
	testRegistry := toolregistrytest.NewTestToolRegistry(t)

	// Read the application config from the testdata file
	appCfg := sdk.LoadApplicationConfigForTest[kubeconfig.KubernetesApplicationSpec](t, filepath.Join("testdata", "traffic_routing_no_deploy_target", "app.pipecd.yaml"), "kubernetes")

	// Prepare stage config to route 100% traffic to primary
	stageCfg := kubeconfig.K8sTrafficRoutingStageOptions{
		All: "primary",
	}
	stageCfgBytes, err := json.Marshal(stageCfg)
	require.NoError(t, err)

	// Prepare the input
	input := &sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec]{
		Request: sdk.ExecuteStageRequest[kubeconfig.KubernetesApplicationSpec]{
			StageName:   "K8S_TRAFFIC_ROUTING",
			StageConfig: stageCfgBytes,
			TargetDeploymentSource: sdk.DeploymentSource[kubeconfig.KubernetesApplicationSpec]{
				ApplicationDirectory:      filepath.Join("testdata", "traffic_routing_no_deploy_target"),
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

	plugin := &Plugin{}
	// Execute with empty deploy targets
	status := plugin.executeK8sTrafficRoutingStagePodSelector(ctx, input, []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]{}, appCfg)

	assert.Equal(t, sdk.StageStatusFailure, status)
}

func TestPlugin_executeK8sTrafficRoutingStagePodSelector_multipleServices(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	// Initialize tool registry
	testRegistry := toolregistrytest.NewTestToolRegistry(t)

	// Read the application config from the testdata file
	appCfg := sdk.LoadApplicationConfigForTest[kubeconfig.KubernetesApplicationSpec](t, filepath.Join("testdata", "traffic_routing_multiple_services", "app.pipecd.yaml"), "kubernetes")

	// Prepare stage config to route 100% traffic to primary
	stageCfg := kubeconfig.K8sTrafficRoutingStageOptions{
		All: "primary",
	}
	stageCfgBytes, err := json.Marshal(stageCfg)
	require.NoError(t, err)

	// Prepare the input
	input := &sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec]{
		Request: sdk.ExecuteStageRequest[kubeconfig.KubernetesApplicationSpec]{
			StageName:   "K8S_TRAFFIC_ROUTING",
			StageConfig: stageCfgBytes,
			TargetDeploymentSource: sdk.DeploymentSource[kubeconfig.KubernetesApplicationSpec]{
				ApplicationDirectory:      filepath.Join("testdata", "traffic_routing_multiple_services"),
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

	// Initialize deploy target config and dynamic client for assertions with envtest
	dtConfig, dynamicClient := setupTestDeployTargetConfigAndDynamicClient(t)

	// First apply the services
	applyCfg := sdk.LoadApplicationConfigForTest[kubeconfig.KubernetesApplicationSpec](t, filepath.Join("testdata", "traffic_routing_multiple_services", "app.pipecd.yaml"), "kubernetes")
	applyInput := &sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec]{
		Request: sdk.ExecuteStageRequest[kubeconfig.KubernetesApplicationSpec]{
			StageName:   "K8S_SYNC",
			StageConfig: []byte(``),
			TargetDeploymentSource: sdk.DeploymentSource[kubeconfig.KubernetesApplicationSpec]{
				ApplicationDirectory:      filepath.Join("testdata", "traffic_routing_multiple_services"),
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

	// Now execute traffic routing
	status = plugin.executeK8sTrafficRoutingStagePodSelector(ctx, input, []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]{
		{
			Name:   "default",
			Config: *dtConfig,
		},
	}, appCfg)

	assert.Equal(t, sdk.StageStatusSuccess, status)

	// Verify only the first Service selector was updated
	service1, err := dynamicClient.Resource(schema.GroupVersionResource{Group: "", Version: "v1", Resource: "services"}).Namespace("default").Get(context.Background(), "traffic-test-1", metav1.GetOptions{})
	require.NoError(t, err)
	selector1 := service1.Object["spec"].(map[string]interface{})["selector"].(map[string]interface{})
	assert.Equal(t, "primary", selector1["pipecd.dev/variant"])

	// Second service should remain unchanged
	service2, err := dynamicClient.Resource(schema.GroupVersionResource{Group: "", Version: "v1", Resource: "services"}).Namespace("default").Get(context.Background(), "traffic-test-2", metav1.GetOptions{})
	require.NoError(t, err)
	selector2 := service2.Object["spec"].(map[string]interface{})["selector"].(map[string]interface{})
	assert.Equal(t, "canary", selector2["pipecd.dev/variant"])

	// Multiple services test verified by checking that only the first service was updated
}
