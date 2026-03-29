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

func TestPlugin_executeK8sMultiTrafficRoutingStage_PodSelector_RouteToCanary(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	configDir := filepath.Join("testdata", "traffic_routing_pod_selector")
	appCfg := sdk.LoadApplicationConfigForTest[kubeconfig.KubernetesApplicationSpec](t, filepath.Join(configDir, "app.pipecd.yaml"), "kubernetes_multicluster")

	testRegistry := toolregistrytest.NewTestToolRegistry(t)

	input := &sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec]{
		Request: sdk.ExecuteStageRequest[kubeconfig.KubernetesApplicationSpec]{
			StageName:   StageK8sMultiTrafficRouting,
			StageConfig: []byte(`{"canary": 100}`),
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
		Client: sdk.NewClient(nil, "kubernetes_multicluster", "app-id", "stage-id", logpersistertest.NewTestLogPersister(t), testRegistry),
		Logger: zaptest.NewLogger(t),
	}

	dtConfig, dynamicClient := setupTestDeployTargetConfigAndDynamicClient(t)

	plugin := &Plugin{}
	status := plugin.executeK8sMultiTrafficRoutingStage(ctx, input, []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]{
		{Name: "default", Config: *dtConfig},
	})

	assert.Equal(t, sdk.StageStatusSuccess, status)

	// The Service selector should now point to canary.
	serviceRes := schema.GroupVersionResource{Group: "", Version: "v1", Resource: "services"}
	svc, err := dynamicClient.Resource(serviceRes).Namespace("default").Get(ctx, "simple", metav1.GetOptions{})
	require.NoError(t, err)

	selector, _, err := unstructuredNestedStringMap(svc.Object, "spec", "selector")
	require.NoError(t, err)
	assert.Equal(t, "canary", selector["pipecd.dev/variant"])
}

func TestPlugin_executeK8sMultiTrafficRoutingStage_PodSelector_RouteToCanary_MultiCluster(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	configDir := filepath.Join("testdata", "traffic_routing_pod_selector")
	appCfg := sdk.LoadApplicationConfigForTest[kubeconfig.KubernetesApplicationSpec](t, filepath.Join(configDir, "app.pipecd.yaml"), "kubernetes_multicluster")

	testRegistry := toolregistrytest.NewTestToolRegistry(t)

	input := &sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec]{
		Request: sdk.ExecuteStageRequest[kubeconfig.KubernetesApplicationSpec]{
			StageName:   StageK8sMultiTrafficRouting,
			StageConfig: []byte(`{"canary": 100}`),
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
		Client: sdk.NewClient(nil, "kubernetes_multicluster", "app-id", "stage-id", logpersistertest.NewTestLogPersister(t), testRegistry),
		Logger: zaptest.NewLogger(t),
	}

	clusterUS := setupCluster(t, "cluster-us")
	clusterEU := setupCluster(t, "cluster-eu")

	plugin := &Plugin{}
	status := plugin.executeK8sMultiTrafficRoutingStage(ctx, input, []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]{
		{Name: clusterUS.name, Config: *clusterUS.dtc},
		{Name: clusterEU.name, Config: *clusterEU.dtc},
	})

	assert.Equal(t, sdk.StageStatusSuccess, status)

	serviceRes := schema.GroupVersionResource{Group: "", Version: "v1", Resource: "services"}
	for _, c := range []*cluster{clusterUS, clusterEU} {
		svc, err := c.cli.Resource(serviceRes).Namespace("default").Get(ctx, "simple", metav1.GetOptions{})
		require.NoError(t, err)
		selector, _, err := unstructuredNestedStringMap(svc.Object, "spec", "selector")
		require.NoError(t, err)
		assert.Equal(t, "canary", selector["pipecd.dev/variant"], "cluster %s should have canary selector", c.name)
	}
}

func TestPlugin_executeK8sMultiTrafficRoutingStage_PodSelector_RestoreToPrimary(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	configDir := filepath.Join("testdata", "traffic_routing_pod_selector")
	appCfg := sdk.LoadApplicationConfigForTest[kubeconfig.KubernetesApplicationSpec](t, filepath.Join(configDir, "app.pipecd.yaml"), "kubernetes_multicluster")

	testRegistry := toolregistrytest.NewTestToolRegistry(t)

	input := &sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec]{
		Request: sdk.ExecuteStageRequest[kubeconfig.KubernetesApplicationSpec]{
			StageName:   StageK8sMultiTrafficRouting,
			StageConfig: []byte(`{"primary": 100}`),
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
		Client: sdk.NewClient(nil, "kubernetes_multicluster", "app-id", "stage-id", logpersistertest.NewTestLogPersister(t), testRegistry),
		Logger: zaptest.NewLogger(t),
	}

	dtConfig, dynamicClient := setupTestDeployTargetConfigAndDynamicClient(t)

	plugin := &Plugin{}
	status := plugin.executeK8sMultiTrafficRoutingStage(ctx, input, []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]{
		{Name: "default", Config: *dtConfig},
	})

	assert.Equal(t, sdk.StageStatusSuccess, status)

	serviceRes := schema.GroupVersionResource{Group: "", Version: "v1", Resource: "services"}
	svc, err := dynamicClient.Resource(serviceRes).Namespace("default").Get(ctx, "simple", metav1.GetOptions{})
	require.NoError(t, err)

	selector, _, err := unstructuredNestedStringMap(svc.Object, "spec", "selector")
	require.NoError(t, err)
	assert.Equal(t, "primary", selector["pipecd.dev/variant"])
}

func TestPlugin_executeK8sMultiTrafficRoutingStage_PodSelector_RejectBaseline(t *testing.T) {
	t.Parallel()

	configDir := filepath.Join("testdata", "traffic_routing_pod_selector")
	appCfg := sdk.LoadApplicationConfigForTest[kubeconfig.KubernetesApplicationSpec](t, filepath.Join(configDir, "app.pipecd.yaml"), "kubernetes_multicluster")

	testRegistry := toolregistrytest.NewTestToolRegistry(t)

	input := &sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec]{
		Request: sdk.ExecuteStageRequest[kubeconfig.KubernetesApplicationSpec]{
			StageName:   StageK8sMultiTrafficRouting,
			StageConfig: []byte(`{"baseline": 100}`),
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
		Client: sdk.NewClient(nil, "kubernetes_multicluster", "app-id", "stage-id", logpersistertest.NewTestLogPersister(t), testRegistry),
		Logger: zaptest.NewLogger(t),
	}

	dtConfig, _ := setupTestDeployTargetConfigAndDynamicClient(t)

	plugin := &Plugin{}
	status := plugin.executeK8sMultiTrafficRoutingStage(t.Context(), input, []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]{
		{Name: "default", Config: *dtConfig},
	})

	// PodSelector does not support baseline variant.
	assert.Equal(t, sdk.StageStatusFailure, status)
}

func TestPlugin_executeK8sMultiTrafficRoutingStage_PodSelector_RejectSplit(t *testing.T) {
	t.Parallel()

	configDir := filepath.Join("testdata", "traffic_routing_pod_selector")
	appCfg := sdk.LoadApplicationConfigForTest[kubeconfig.KubernetesApplicationSpec](t, filepath.Join(configDir, "app.pipecd.yaml"), "kubernetes_multicluster")

	testRegistry := toolregistrytest.NewTestToolRegistry(t)

	input := &sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec]{
		Request: sdk.ExecuteStageRequest[kubeconfig.KubernetesApplicationSpec]{
			StageName:   StageK8sMultiTrafficRouting,
			StageConfig: []byte(`{"primary": 50, "canary": 50}`),
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
		Client: sdk.NewClient(nil, "kubernetes_multicluster", "app-id", "stage-id", logpersistertest.NewTestLogPersister(t), testRegistry),
		Logger: zaptest.NewLogger(t),
	}

	dtConfig, _ := setupTestDeployTargetConfigAndDynamicClient(t)

	plugin := &Plugin{}
	status := plugin.executeK8sMultiTrafficRoutingStage(t.Context(), input, []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]{
		{Name: "default", Config: *dtConfig},
	})

	// PodSelector requires one variant to be 100%.
	assert.Equal(t, sdk.StageStatusFailure, status)
}

// unstructuredNestedStringMap extracts a map[string]string from an unstructured object.
func unstructuredNestedStringMap(obj map[string]any, fields ...string) (map[string]string, bool, error) {
	m, found, err := nestedMap(obj, fields...)
	if !found || err != nil {
		return nil, found, err
	}
	result := make(map[string]string, len(m))
	for k, v := range m {
		s, ok := v.(string)
		if !ok {
			continue
		}
		result[k] = s
	}
	return result, true, nil
}

func nestedMap(obj map[string]any, fields ...string) (map[string]any, bool, error) {
	cur := obj
	for _, f := range fields {
		v, ok := cur[f]
		if !ok {
			return nil, false, nil
		}
		m, ok := v.(map[string]any)
		if !ok {
			return nil, false, nil
		}
		cur = m
	}
	return cur, true, nil
}
