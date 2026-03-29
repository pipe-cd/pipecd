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
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes_multicluster/provider"
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

/*
// installIstioCRDs applies the Istio CRD bundle (via Helm/kustomize) to the envtest cluster
// pointed at by dtConfig. It reuses the K8S_MULTI_SYNC stage so the same code path that
// handles Helm kustomize overlays is exercised.
func installIstioCRDs(t *testing.T, dtConfig *kubeconfig.KubernetesDeployTargetConfig) {
	t.Helper()

	// Copy the crds dir to a temp dir — kustomize writes intermediate files alongside the
	// kustomization.yaml and will fail if the source is read-only or shared between tests.
	istioCrdsDir := filepath.Join("testdata", "istio_crds")
	istioTempDir := t.TempDir()
	require.NoError(t, os.CopyFS(istioTempDir, os.DirFS(istioCrdsDir)))

	testRegistry := toolregistrytest.NewTestToolRegistry(t)
	crdCfg := sdk.LoadApplicationConfigForTest[kubeconfig.KubernetesApplicationSpec](
		t, filepath.Join(istioTempDir, "app.pipecd.yaml"), "kubernetes_multicluster")

	input := &sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec]{
		Request: sdk.ExecuteStageRequest[kubeconfig.KubernetesApplicationSpec]{
			StageName:   StageK8sMultiSync,
			StageConfig: []byte(``),
			TargetDeploymentSource: sdk.DeploymentSource[kubeconfig.KubernetesApplicationSpec]{
				ApplicationDirectory:      istioTempDir,
				CommitHash:                "0123456789",
				ApplicationConfig:         crdCfg,
				ApplicationConfigFilename: "app.pipecd.yaml",
			},
			Deployment: sdk.Deployment{
				PipedID:       "piped-id",
				ApplicationID: "istio-crds-app-id",
			},
		},
		Client: sdk.NewClient(nil, "kubernetes_multicluster", "istio-crds-app-id", "stage-id",
			logpersistertest.NewTestLogPersister(t), testRegistry),
		Logger: zaptest.NewLogger(t),
	}

	plugin := &Plugin{}
	status := plugin.executeK8sMultiSyncStage(t.Context(), input,
		[]*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]{
			{Name: "default", Config: *dtConfig},
		})
	require.Equal(t, sdk.StageStatusSuccess, status)
}

// applyManifestsByMultiSync runs K8S_MULTI_SYNC to pre-apply the manifests in testdataDir
// into the cluster — used to seed the VirtualService before the traffic routing stage runs.
func applyManifestsByMultiSync(t *testing.T, ctx context.Context, testdataDir string, dtConfig *kubeconfig.KubernetesDeployTargetConfig) {
	t.Helper()

	testRegistry := toolregistrytest.NewTestToolRegistry(t)
	cfg := sdk.LoadApplicationConfigForTest[kubeconfig.KubernetesApplicationSpec](
		t, filepath.Join("testdata", testdataDir, "app.pipecd.yaml"), "kubernetes_multicluster")

	input := &sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec]{
		Request: sdk.ExecuteStageRequest[kubeconfig.KubernetesApplicationSpec]{
			StageName:   StageK8sMultiSync,
			StageConfig: []byte(``),
			TargetDeploymentSource: sdk.DeploymentSource[kubeconfig.KubernetesApplicationSpec]{
				ApplicationDirectory:      filepath.Join("testdata", testdataDir),
				CommitHash:                "0123456789",
				ApplicationConfig:         cfg,
				ApplicationConfigFilename: "app.pipecd.yaml",
			},
			Deployment: sdk.Deployment{
				PipedID:       "piped-id",
				ApplicationID: "app-id",
			},
		},
		Client: sdk.NewClient(nil, "kubernetes_multicluster", "app-id", "stage-id",
			logpersistertest.NewTestLogPersister(t), testRegistry),
		Logger: zaptest.NewLogger(t),
	}

	plugin := &Plugin{}
	status := plugin.executeK8sMultiSyncStage(ctx, input,
		[]*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]{
			{Name: "default", Config: *dtConfig},
		})
	require.Equal(t, sdk.StageStatusSuccess, status)
}

// expectedRoute holds the expected destination for one route entry in a VirtualService.
type expectedRoute struct {
	host   string
	subset string
	weight int32
}

// verifyVirtualServiceRouting reads the named VirtualService from the cluster and asserts
// that the first HTTP route has exactly the given destinations in order.
func verifyVirtualServiceRouting(t *testing.T, dynamicClient dynamic.Interface, vsName string, expected []expectedRoute) {
	t.Helper()

	vs, err := dynamicClient.Resource(schema.GroupVersionResource{
		Group:   "networking.istio.io",
		Version: "v1",
		Resource: "virtualservices",
	}).Namespace("default").Get(t.Context(), vsName, metav1.GetOptions{})
	require.NoError(t, err)

	spec := vs.Object["spec"].(map[string]any)
	httpRoutes := spec["http"].([]any)
	require.NotEmpty(t, httpRoutes)

	routes := httpRoutes[0].(map[string]any)["route"].([]any)
	require.Len(t, routes, len(expected), "number of route destinations")

	for i, exp := range expected {
		dest := routes[i].(map[string]any)["destination"].(map[string]any)
		assert.Equal(t, exp.host, dest["host"], "host mismatch for route %d", i)
		assert.Equal(t, exp.subset, dest["subset"], "subset mismatch for route %d", i)

		rawWeight := routes[i].(map[string]any)["weight"]
		switch w := rawWeight.(type) {
		case int64:
			assert.Equal(t, int64(exp.weight), w, "weight mismatch for route %d", i)
		case float64:
			assert.Equal(t, float64(exp.weight), w, "weight mismatch for route %d", i)
		default:
			t.Errorf("unexpected weight type %T for route %d", rawWeight, i)
		}
	}
}

// verifyVirtualServiceEditableRoutes checks that only the named editable route was
// modified and the non-editable route still has a single primary-only destination.
func verifyVirtualServiceEditableRoutes(t *testing.T, dynamicClient dynamic.Interface, vsName string) {
	t.Helper()

	vs, err := dynamicClient.Resource(schema.GroupVersionResource{
		Group:   "networking.istio.io",
		Version: "v1",
		Resource: "virtualservices",
	}).Namespace("default").Get(t.Context(), vsName, metav1.GetOptions{})
	require.NoError(t, err)

	spec := vs.Object["spec"].(map[string]any)
	httpRoutes := spec["http"].([]any)
	require.Len(t, httpRoutes, 2, "expected exactly two HTTP routes")

	// api-route should have been modified (primary + canary)
	apiRoute := httpRoutes[0].(map[string]any)
	assert.Equal(t, "api-route", apiRoute["name"])
	assert.Len(t, apiRoute["route"].([]any), 2, "api-route should have 2 destinations")

	// web-route should remain unchanged (primary only)
	webRoute := httpRoutes[1].(map[string]any)
	assert.Equal(t, "web-route", webRoute["name"])
	webDests := webRoute["route"].([]any)
	assert.Len(t, webDests, 1, "web-route should remain with 1 destination")
	dest := webDests[0].(map[string]any)["destination"].(map[string]any)
	assert.Equal(t, "primary", dest["subset"])
}

func TestPlugin_executeK8sMultiTrafficRoutingStageIstio(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name        string
		testdataDir string
		stageCfg    kubeconfig.K8sTrafficRoutingStageOptions
		wantStatus  sdk.StageStatus
		verify      func(t *testing.T, dynamicClient dynamic.Interface)
	}

	tests := []testCase{
		{
			name:        "canary 30%",
			testdataDir: "traffic_routing_istio",
			stageCfg:    kubeconfig.K8sTrafficRoutingStageOptions{Canary: unit.Percentage{Number: 30}},
			wantStatus:  sdk.StageStatusSuccess,
			verify: func(t *testing.T, dynamicClient dynamic.Interface) {
				verifyVirtualServiceRouting(t, dynamicClient, "traffic-test-vs", []expectedRoute{
					{host: "traffic-test", subset: "primary", weight: 70},
					{host: "traffic-test", subset: "canary", weight: 30},
				})
			},
		},
		{
			name:        "primary canary baseline split",
			testdataDir: "traffic_routing_istio_baseline",
			stageCfg: kubeconfig.K8sTrafficRoutingStageOptions{
				Primary:  unit.Percentage{Number: 50},
				Canary:   unit.Percentage{Number: 30},
				Baseline: unit.Percentage{Number: 20},
			},
			wantStatus: sdk.StageStatusSuccess,
			verify: func(t *testing.T, dynamicClient dynamic.Interface) {
				verifyVirtualServiceRouting(t, dynamicClient, "traffic-test-vs", []expectedRoute{
					{host: "traffic-test", subset: "primary", weight: 50},
					{host: "traffic-test", subset: "canary", weight: 30},
					{host: "traffic-test", subset: "baseline", weight: 20},
				})
			},
		},
		{
			name:        "restore to primary 100%",
			testdataDir: "traffic_routing_istio",
			stageCfg:    kubeconfig.K8sTrafficRoutingStageOptions{Primary: unit.Percentage{Number: 100}},
			wantStatus:  sdk.StageStatusSuccess,
			verify: func(t *testing.T, dynamicClient dynamic.Interface) {
				verifyVirtualServiceRouting(t, dynamicClient, "traffic-test-vs", []expectedRoute{
					{host: "traffic-test", subset: "primary", weight: 100},
				})
			},
		},
		{
			name:        "editable routes filter",
			testdataDir: "traffic_routing_istio_editable_routes",
			stageCfg:    kubeconfig.K8sTrafficRoutingStageOptions{Canary: unit.Percentage{Number: 40}},
			wantStatus:  sdk.StageStatusSuccess,
			verify: func(t *testing.T, dynamicClient dynamic.Interface) {
				verifyVirtualServiceEditableRoutes(t, dynamicClient, "traffic-test-vs")
			},
		},
		{
			name:        "no virtualservice manifest",
			testdataDir: "traffic_routing_istio_no_virtualservice",
			stageCfg:    kubeconfig.K8sTrafficRoutingStageOptions{Canary: unit.Percentage{Number: 30}},
			wantStatus:  sdk.StageStatusFailure,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()

			dtConfig, dynamicClient := setupTestDeployTargetConfigAndDynamicClient(t)

			// Install Istio CRDs so the envtest cluster accepts VirtualService objects.
			installIstioCRDs(t, dtConfig)

			// Pre-apply the manifests (including the VirtualService) so the traffic routing
			// stage has something to read back and update.
			applyManifestsByMultiSync(t, ctx, tt.testdataDir, dtConfig)

			stageCfgBytes, err := json.Marshal(tt.stageCfg)
			require.NoError(t, err)

			appCfg := sdk.LoadApplicationConfigForTest[kubeconfig.KubernetesApplicationSpec](
				t, filepath.Join("testdata", tt.testdataDir, "app.pipecd.yaml"), "kubernetes_multicluster")

			input := &sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec]{
				Request: sdk.ExecuteStageRequest[kubeconfig.KubernetesApplicationSpec]{
					StageName:   StageK8sMultiTrafficRouting,
					StageConfig: stageCfgBytes,
					TargetDeploymentSource: sdk.DeploymentSource[kubeconfig.KubernetesApplicationSpec]{
						ApplicationDirectory:      filepath.Join("testdata", tt.testdataDir),
						CommitHash:                "0123456789",
						ApplicationConfig:         appCfg,
						ApplicationConfigFilename: "app.pipecd.yaml",
					},
					Deployment: sdk.Deployment{
						PipedID:       "piped-id",
						ApplicationID: "app-id",
					},
				},
				Client: sdk.NewClient(nil, "kubernetes_multicluster", "app-id", "stage-id",
					logpersistertest.NewTestLogPersister(t), toolregistrytest.NewTestToolRegistry(t)),
				Logger: zaptest.NewLogger(t),
			}

			plugin := &Plugin{}
			status := plugin.executeK8sMultiTrafficRoutingStage(ctx, input,
				[]*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]{
					{Name: "default", Config: *dtConfig},
				})

			assert.Equal(t, tt.wantStatus, status)

			if tt.verify != nil {
				tt.verify(t, dynamicClient)
			}
		})
	}
}

func TestPlugin_executeK8sMultiTrafficRoutingStageIstio_MultiCluster(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	clusterUS := setupCluster(t, "cluster-us")
	clusterEU := setupCluster(t, "cluster-eu")

	for _, c := range []*cluster{clusterUS, clusterEU} {
		installIstioCRDs(t, c.dtc)
		applyManifestsByMultiSync(t, ctx, "traffic_routing_istio", c.dtc)
	}

	stageCfgBytes, err := json.Marshal(kubeconfig.K8sTrafficRoutingStageOptions{
		Canary: unit.Percentage{Number: 30},
	})
	require.NoError(t, err)

	appCfg := sdk.LoadApplicationConfigForTest[kubeconfig.KubernetesApplicationSpec](
		t, filepath.Join("testdata", "traffic_routing_istio", "app.pipecd.yaml"), "kubernetes_multicluster")

	input := &sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec]{
		Request: sdk.ExecuteStageRequest[kubeconfig.KubernetesApplicationSpec]{
			StageName:   StageK8sMultiTrafficRouting,
			StageConfig: stageCfgBytes,
			TargetDeploymentSource: sdk.DeploymentSource[kubeconfig.KubernetesApplicationSpec]{
				ApplicationDirectory:      filepath.Join("testdata", "traffic_routing_istio"),
				CommitHash:                "0123456789",
				ApplicationConfig:         appCfg,
				ApplicationConfigFilename: "app.pipecd.yaml",
			},
			Deployment: sdk.Deployment{
				PipedID:       "piped-id",
				ApplicationID: "app-id",
			},
		},
		Client: sdk.NewClient(nil, "kubernetes_multicluster", "app-id", "stage-id",
			logpersistertest.NewTestLogPersister(t), toolregistrytest.NewTestToolRegistry(t)),
		Logger: zaptest.NewLogger(t),
	}

	plugin := &Plugin{}
	status := plugin.executeK8sMultiTrafficRoutingStage(ctx, input,
		[]*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]{
			{Name: clusterUS.name, Config: *clusterUS.dtc},
			{Name: clusterEU.name, Config: *clusterEU.dtc},
		})

	assert.Equal(t, sdk.StageStatusSuccess, status)

	for _, c := range []*cluster{clusterUS, clusterEU} {
		verifyVirtualServiceRouting(t, c.cli, "traffic-test-vs", []expectedRoute{
			{host: "traffic-test", subset: "primary", weight: 70},
			{host: "traffic-test", subset: "canary", weight: 30},
		})
	}
}

*/

func Test_generateVirtualServiceManifest(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		inputYAML       string
		host            string
		editableRoutes  []string
		canaryPercent   int32
		baselinePercent int32
		variantLabel    kubeconfig.KubernetesVariantLabel
		checkFunc       func(t *testing.T, result provider.Manifest)
	}{
		{
			name: "basic canary traffic routing",
			inputYAML: `
apiVersion: networking.istio.io/v1
kind: VirtualService
metadata:
  name: test-vs
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
			host:            "test-service",
			editableRoutes:  []string{},
			canaryPercent:   30,
			baselinePercent: 0,
			variantLabel: kubeconfig.KubernetesVariantLabel{
				Key:           "pipecd.dev/variant",
				PrimaryValue:  "primary",
				CanaryValue:   "canary",
				BaselineValue: "baseline",
			},
			checkFunc: func(t *testing.T, result provider.Manifest) {
				vs, err := convertVirtualService(result)
				require.NoError(t, err)
				require.Len(t, vs.Spec.Http, 1)
				require.Len(t, vs.Spec.Http[0].Route, 2)

				primaryRoute := vs.Spec.Http[0].Route[0]
				assert.Equal(t, "test-service", primaryRoute.Destination.Host)
				assert.Equal(t, "primary", primaryRoute.Destination.Subset)
				assert.Equal(t, int32(70), primaryRoute.Weight)

				canaryRoute := vs.Spec.Http[0].Route[1]
				assert.Equal(t, "test-service", canaryRoute.Destination.Host)
				assert.Equal(t, "canary", canaryRoute.Destination.Subset)
				assert.Equal(t, int32(30), canaryRoute.Weight)
			},
		},
		{
			name: "canary and baseline traffic routing",
			inputYAML: `
apiVersion: networking.istio.io/v1
kind: VirtualService
metadata:
  name: test-vs
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
			host:            "test-service",
			editableRoutes:  []string{},
			canaryPercent:   20,
			baselinePercent: 30,
			variantLabel: kubeconfig.KubernetesVariantLabel{
				Key:           "pipecd.dev/variant",
				PrimaryValue:  "primary",
				CanaryValue:   "canary",
				BaselineValue: "baseline",
			},
			checkFunc: func(t *testing.T, result provider.Manifest) {
				vs, err := convertVirtualService(result)
				require.NoError(t, err)
				require.Len(t, vs.Spec.Http, 1)
				require.Len(t, vs.Spec.Http[0].Route, 3)

				primaryRoute := vs.Spec.Http[0].Route[0]
				assert.Equal(t, "primary", primaryRoute.Destination.Subset)
				assert.Equal(t, int32(50), primaryRoute.Weight)

				canaryRoute := vs.Spec.Http[0].Route[1]
				assert.Equal(t, "canary", canaryRoute.Destination.Subset)
				assert.Equal(t, int32(20), canaryRoute.Weight)

				baselineRoute := vs.Spec.Http[0].Route[2]
				assert.Equal(t, "baseline", baselineRoute.Destination.Subset)
				assert.Equal(t, int32(30), baselineRoute.Weight)
			},
		},
		{
			name: "preserve other host routes",
			inputYAML: `
apiVersion: networking.istio.io/v1
kind: VirtualService
metadata:
  name: test-vs
spec:
  hosts:
  - test-service
  http:
  - route:
    - destination:
        host: test-service
        subset: primary
      weight: 60
    - destination:
        host: other-service
        subset: v1
      weight: 40
`,
			host:            "test-service",
			editableRoutes:  []string{},
			canaryPercent:   50,
			baselinePercent: 0,
			variantLabel: kubeconfig.KubernetesVariantLabel{
				Key:           "pipecd.dev/variant",
				PrimaryValue:  "primary",
				CanaryValue:   "canary",
				BaselineValue: "baseline",
			},
			checkFunc: func(t *testing.T, result provider.Manifest) {
				vs, err := convertVirtualService(result)
				require.NoError(t, err)
				require.Len(t, vs.Spec.Http, 1)
				require.Len(t, vs.Spec.Http[0].Route, 3)

				// primary and canary each get 50% of the 60% variant budget = 30 each
				primaryRoute := vs.Spec.Http[0].Route[0]
				assert.Equal(t, "primary", primaryRoute.Destination.Subset)
				assert.Equal(t, int32(30), primaryRoute.Weight)

				canaryRoute := vs.Spec.Http[0].Route[1]
				assert.Equal(t, "canary", canaryRoute.Destination.Subset)
				assert.Equal(t, int32(30), canaryRoute.Weight)

				// other-service route preserved unchanged
				otherRoute := vs.Spec.Http[0].Route[2]
				assert.Equal(t, "other-service", otherRoute.Destination.Host)
				assert.Equal(t, "v1", otherRoute.Destination.Subset)
				assert.Equal(t, int32(40), otherRoute.Weight)
			},
		},
		{
			name: "editable routes filter",
			inputYAML: `
apiVersion: networking.istio.io/v1
kind: VirtualService
metadata:
  name: test-vs
spec:
  hosts:
  - test-service
  http:
  - name: editable-route
    route:
    - destination:
        host: test-service
        subset: primary
      weight: 100
  - name: non-editable-route
    route:
    - destination:
        host: test-service
        subset: primary
      weight: 100
`,
			host:            "test-service",
			editableRoutes:  []string{"editable-route"},
			canaryPercent:   40,
			baselinePercent: 0,
			variantLabel: kubeconfig.KubernetesVariantLabel{
				Key:           "pipecd.dev/variant",
				PrimaryValue:  "primary",
				CanaryValue:   "canary",
				BaselineValue: "baseline",
			},
			checkFunc: func(t *testing.T, result provider.Manifest) {
				vs, err := convertVirtualService(result)
				require.NoError(t, err)
				require.Len(t, vs.Spec.Http, 2)

				editableHTTP := vs.Spec.Http[0]
				assert.Equal(t, "editable-route", editableHTTP.Name)
				require.Len(t, editableHTTP.Route, 2)
				assert.Equal(t, int32(60), editableHTTP.Route[0].Weight) // primary
				assert.Equal(t, int32(40), editableHTTP.Route[1].Weight) // canary

				nonEditableHTTP := vs.Spec.Http[1]
				assert.Equal(t, "non-editable-route", nonEditableHTTP.Name)
				require.Len(t, nonEditableHTTP.Route, 1)
				assert.Equal(t, int32(100), nonEditableHTTP.Route[0].Weight)
				assert.Equal(t, "primary", nonEditableHTTP.Route[0].Destination.Subset)
			},
		},
		{
			name: "restore to primary 100%",
			inputYAML: `
apiVersion: networking.istio.io/v1
kind: VirtualService
metadata:
  name: test-vs
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
			host:            "test-service",
			editableRoutes:  []string{},
			canaryPercent:   0,
			baselinePercent: 0,
			variantLabel: kubeconfig.KubernetesVariantLabel{
				Key:           "pipecd.dev/variant",
				PrimaryValue:  "primary",
				CanaryValue:   "canary",
				BaselineValue: "baseline",
			},
			checkFunc: func(t *testing.T, result provider.Manifest) {
				vs, err := convertVirtualService(result)
				require.NoError(t, err)
				require.Len(t, vs.Spec.Http, 1)
				require.Len(t, vs.Spec.Http[0].Route, 1)

				primaryRoute := vs.Spec.Http[0].Route[0]
				assert.Equal(t, "primary", primaryRoute.Destination.Subset)
				assert.Equal(t, int32(100), primaryRoute.Weight)
			},
		},
		{
			name: "custom variant labels",
			inputYAML: `
apiVersion: networking.istio.io/v1
kind: VirtualService
metadata:
  name: test-vs
spec:
  hosts:
  - test-service
  http:
  - route:
    - destination:
        host: test-service
        subset: stable
      weight: 100
`,
			host:            "test-service",
			editableRoutes:  []string{},
			canaryPercent:   25,
			baselinePercent: 25,
			variantLabel: kubeconfig.KubernetesVariantLabel{
				Key:           "custom/variant",
				PrimaryValue:  "stable",
				CanaryValue:   "preview",
				BaselineValue: "test",
			},
			checkFunc: func(t *testing.T, result provider.Manifest) {
				vs, err := convertVirtualService(result)
				require.NoError(t, err)
				require.Len(t, vs.Spec.Http, 1)
				require.Len(t, vs.Spec.Http[0].Route, 3)

				assert.Equal(t, "stable", vs.Spec.Http[0].Route[0].Destination.Subset)
				assert.Equal(t, int32(50), vs.Spec.Http[0].Route[0].Weight)

				assert.Equal(t, "preview", vs.Spec.Http[0].Route[1].Destination.Subset)
				assert.Equal(t, int32(25), vs.Spec.Http[0].Route[1].Weight)

				assert.Equal(t, "test", vs.Spec.Http[0].Route[2].Destination.Subset)
				assert.Equal(t, int32(25), vs.Spec.Http[0].Route[2].Weight)
			},
		},
		{
			name: "multiple http routes all editable",
			inputYAML: `
apiVersion: networking.istio.io/v1
kind: VirtualService
metadata:
  name: test-vs
spec:
  hosts:
  - test-service
  http:
  - name: route1
    route:
    - destination:
        host: test-service
        subset: primary
      weight: 100
  - name: route2
    route:
    - destination:
        host: test-service
        subset: primary
      weight: 100
`,
			host:            "test-service",
			editableRoutes:  []string{},
			canaryPercent:   50,
			baselinePercent: 0,
			variantLabel: kubeconfig.KubernetesVariantLabel{
				Key:           "pipecd.dev/variant",
				PrimaryValue:  "primary",
				CanaryValue:   "canary",
				BaselineValue: "baseline",
			},
			checkFunc: func(t *testing.T, result provider.Manifest) {
				vs, err := convertVirtualService(result)
				require.NoError(t, err)
				require.Len(t, vs.Spec.Http, 2)

				for i, httpRoute := range vs.Spec.Http {
					require.Len(t, httpRoute.Route, 2, "route %d", i)
					assert.Equal(t, "primary", httpRoute.Route[0].Destination.Subset)
					assert.Equal(t, int32(50), httpRoute.Route[0].Weight)
					assert.Equal(t, "canary", httpRoute.Route[1].Destination.Subset)
					assert.Equal(t, int32(50), httpRoute.Route[1].Weight)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			manifests := mustParseManifests(t, tt.inputYAML)
			require.Len(t, manifests, 1)

			result, err := generateVirtualServiceManifest(
				manifests[0], tt.host, tt.editableRoutes, tt.variantLabel,
				tt.canaryPercent, tt.baselinePercent,
			)
			require.NoError(t, err)

			if tt.checkFunc != nil {
				tt.checkFunc(t, result)
			}
		})
	}
}

