// Copyright 2024 The PipeCD Authors.
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
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/yaml"

	kubeConfigPkg "github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/config"
	config "github.com/pipe-cd/pipecd/pkg/configv1"
	"github.com/pipe-cd/pipecd/pkg/model"
	"github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1/deployment"
	"github.com/pipe-cd/pipecd/pkg/plugin/logpersister/logpersistertest"
	"github.com/pipe-cd/pipecd/pkg/plugin/toolregistry/toolregistrytest"
)

// TODO: move to a common package
func examplesDir() string {
	d, _ := os.Getwd()
	for {
		if _, err := os.Stat(filepath.Join(d, "examples")); err == nil {
			return filepath.Join(d, "examples")
		}
		d = filepath.Dir(d)
	}
}

func kubeconfigFromRestConfig(restConfig *rest.Config) (string, error) {
	clusters := make(map[string]*clientcmdapi.Cluster)
	clusters["default-cluster"] = &clientcmdapi.Cluster{
		Server:                   restConfig.Host,
		CertificateAuthorityData: restConfig.CAData,
	}
	contexts := make(map[string]*clientcmdapi.Context)
	contexts["default-context"] = &clientcmdapi.Context{
		Cluster:  "default-cluster",
		AuthInfo: "default-user",
	}
	authinfos := make(map[string]*clientcmdapi.AuthInfo)
	authinfos["default-user"] = &clientcmdapi.AuthInfo{
		ClientCertificateData: restConfig.CertData,
		ClientKeyData:         restConfig.KeyData,
	}
	clientConfig := clientcmdapi.Config{
		Kind:           "Config",
		APIVersion:     "v1",
		Clusters:       clusters,
		Contexts:       contexts,
		CurrentContext: "default-context",
		AuthInfos:      authinfos,
	}
	b, err := clientcmd.Write(clientConfig)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func setupEnvTest(t *testing.T) *rest.Config {
	t.Helper()

	tEnv := new(envtest.Environment)
	kubeCfg, err := tEnv.Start()
	require.NoError(t, err)
	t.Cleanup(func() { tEnv.Stop() })

	return kubeCfg
}

func setupTestPluginConfig(t *testing.T, kubeCfg *rest.Config) *config.PipedPlugin {
	t.Helper()

	kubeconfig, err := kubeconfigFromRestConfig(kubeCfg)
	require.NoError(t, err)

	workDir := t.TempDir()
	kubeconfigPath := path.Join(workDir, "kubeconfig")
	err = os.WriteFile(kubeconfigPath, []byte(kubeconfig), 0755)
	require.NoError(t, err)

	deployTarget, err := json.Marshal(kubeConfigPkg.KubernetesDeployTargetConfig{KubeConfigPath: kubeconfigPath})
	require.NoError(t, err)

	// prepare the piped plugin config
	return &config.PipedPlugin{
		Name: "kubernetes",
		URL:  "file:///path/to/kubernetes/plugin", // dummy for testing
		Port: 0,                                   // dummy for testing
		DeployTargets: []config.PipedDeployTarget{{
			Name:   "default",
			Labels: map[string]string{},
			Config: json.RawMessage(deployTarget),
		}},
	}
}

func setupTestPluginConfigAndDynamicClient(t *testing.T) (*config.PipedPlugin, dynamic.Interface) {
	t.Helper()

	kubeCfg := setupEnvTest(t)
	pluginCfg := setupTestPluginConfig(t, kubeCfg)

	dynamicClient, err := dynamic.NewForConfig(kubeCfg)
	require.NoError(t, err)

	return pluginCfg, dynamicClient
}

func TestDeploymentService_executeK8sSyncStage(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	// read the application config from the example file
	cfg, err := os.ReadFile(filepath.Join(examplesDir(), "kubernetes", "simple", "app.pipecd.yaml"))
	require.NoError(t, err)

	// prepare the request
	req := &deployment.ExecuteStageRequest{
		Input: &deployment.ExecutePluginInput{
			Deployment: &model.Deployment{
				PipedId:       "piped-id",
				ApplicationId: "app-id",
				DeployTargets: []string{"default"},
			},
			Stage: &model.PipelineStage{
				Id:   "stage-id",
				Name: "K8S_SYNC",
			},
			StageConfig:             []byte(``),
			RunningDeploymentSource: nil,
			TargetDeploymentSource: &deployment.DeploymentSource{
				ApplicationDirectory:      filepath.Join(examplesDir(), "kubernetes", "simple"),
				CommitHash:                "0123456789",
				ApplicationConfig:         cfg,
				ApplicationConfigFilename: "app.pipecd.yaml",
			},
		},
	}

	// initialize tool registry
	testRegistry, err := toolregistrytest.NewToolRegistry(t)
	require.NoError(t, err)

	// initialize plugin config and dynamic client for assertions with envtest
	pluginCfg, dynamicClient := setupTestPluginConfigAndDynamicClient(t)

	svc := NewDeploymentService(pluginCfg, zaptest.NewLogger(t), testRegistry, logpersistertest.NewTestLogPersister(t))
	resp, err := svc.ExecuteStage(ctx, req)

	require.NoError(t, err)
	assert.Equal(t, model.StageStatus_STAGE_SUCCESS.String(), resp.GetStatus().String())

	deployment, err := dynamicClient.Resource(schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}).Namespace("default").Get(context.Background(), "simple", metav1.GetOptions{})
	require.NoError(t, err)

	assert.Equal(t, "simple", deployment.GetName())
	assert.Equal(t, "simple", deployment.GetLabels()["app"])

	assert.Equal(t, "piped", deployment.GetLabels()["pipecd.dev/managed-by"])
	assert.Equal(t, "piped-id", deployment.GetLabels()["pipecd.dev/piped"])
	assert.Equal(t, "app-id", deployment.GetLabels()["pipecd.dev/application"])
	assert.Equal(t, "0123456789", deployment.GetLabels()["pipecd.dev/commit-hash"])

	assert.Equal(t, "piped", deployment.GetAnnotations()["pipecd.dev/managed-by"])
	assert.Equal(t, "piped-id", deployment.GetAnnotations()["pipecd.dev/piped"])
	assert.Equal(t, "app-id", deployment.GetAnnotations()["pipecd.dev/application"])
	assert.Equal(t, "apps/v1", deployment.GetAnnotations()["pipecd.dev/original-api-version"])
	assert.Equal(t, "apps:Deployment::simple", deployment.GetAnnotations()["pipecd.dev/resource-key"]) // This assertion differs from the non-plugin-arched piped's Kubernetes platform provider, but we decided to change this behavior.
	assert.Equal(t, "0123456789", deployment.GetAnnotations()["pipecd.dev/commit-hash"])

}

func TestDeploymentService_executeK8sSyncStage_withInputNamespace(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	// read the application config from the example file
	cfg, err := os.ReadFile(filepath.Join(examplesDir(), "kubernetes", "simple", "app.pipecd.yaml"))
	require.NoError(t, err)

	// decode and override the autoCreateNamespace and namespace
	spec, err := config.DecodeYAML[*kubeConfigPkg.KubernetesApplicationSpec](cfg)
	require.NoError(t, err)
	spec.Spec.Input.AutoCreateNamespace = true
	spec.Spec.Input.Namespace = "test-namespace"
	cfg, err = yaml.Marshal(spec)
	require.NoError(t, err)

	// prepare the request
	req := &deployment.ExecuteStageRequest{
		Input: &deployment.ExecutePluginInput{
			Deployment: &model.Deployment{
				PipedId:       "piped-id",
				ApplicationId: "app-id",
				DeployTargets: []string{"default"},
			},
			Stage: &model.PipelineStage{
				Id:   "stage-id",
				Name: "K8S_SYNC",
			},
			StageConfig:             []byte(``),
			RunningDeploymentSource: nil,
			TargetDeploymentSource: &deployment.DeploymentSource{
				ApplicationDirectory:      filepath.Join(examplesDir(), "kubernetes", "simple"),
				CommitHash:                "0123456789",
				ApplicationConfig:         cfg,
				ApplicationConfigFilename: "app.pipecd.yaml",
			},
		},
	}

	// initialize tool registry
	testRegistry, err := toolregistrytest.NewToolRegistry(t)
	require.NoError(t, err)

	// initialize plugin config and dynamic client for assertions with envtest
	pluginCfg, dynamicClient := setupTestPluginConfigAndDynamicClient(t)

	svc := NewDeploymentService(pluginCfg, zaptest.NewLogger(t), testRegistry, logpersistertest.NewTestLogPersister(t))
	resp, err := svc.ExecuteStage(ctx, req)

	require.NoError(t, err)
	assert.Equal(t, model.StageStatus_STAGE_SUCCESS.String(), resp.GetStatus().String())

	deployment, err := dynamicClient.Resource(schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}).Namespace("test-namespace").Get(context.Background(), "simple", metav1.GetOptions{})
	require.NoError(t, err)

	assert.Equal(t, "piped", deployment.GetLabels()["pipecd.dev/managed-by"])
	assert.Equal(t, "piped-id", deployment.GetLabels()["pipecd.dev/piped"])
	assert.Equal(t, "app-id", deployment.GetLabels()["pipecd.dev/application"])
	assert.Equal(t, "0123456789", deployment.GetLabels()["pipecd.dev/commit-hash"])

	assert.Equal(t, "simple", deployment.GetName())
	assert.Equal(t, "simple", deployment.GetLabels()["app"])
	assert.Equal(t, "piped", deployment.GetAnnotations()["pipecd.dev/managed-by"])
	assert.Equal(t, "piped-id", deployment.GetAnnotations()["pipecd.dev/piped"])
	assert.Equal(t, "app-id", deployment.GetAnnotations()["pipecd.dev/application"])
	assert.Equal(t, "apps/v1", deployment.GetAnnotations()["pipecd.dev/original-api-version"])
	assert.Equal(t, "apps:Deployment:test-namespace:simple", deployment.GetAnnotations()["pipecd.dev/resource-key"]) // This assertion differs from the non-plugin-arched piped's Kubernetes platform provider, but we decided to change this behavior.
	assert.Equal(t, "0123456789", deployment.GetAnnotations()["pipecd.dev/commit-hash"])
}

func TestDeploymentService_executeK8sSyncStage_withPrune(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	// initialize tool registry
	testRegistry, err := toolregistrytest.NewToolRegistry(t)
	require.NoError(t, err)

	// initialize plugin config and dynamic client for assertions with envtest
	pluginCfg, dynamicClient := setupTestPluginConfigAndDynamicClient(t)

	svc := NewDeploymentService(pluginCfg, zaptest.NewLogger(t), testRegistry, logpersistertest.NewTestLogPersister(t))

	running := filepath.Join("./", "testdata", "prune", "running")

	// read the running application config from the testdata file
	runningCfg, err := os.ReadFile(filepath.Join(running, "app.pipecd.yaml"))
	require.NoError(t, err)

	ok := t.Run("prepare", func(t *testing.T) {
		runningRequest := &deployment.ExecuteStageRequest{
			Input: &deployment.ExecutePluginInput{
				Deployment: &model.Deployment{
					PipedId:       "piped-id",
					ApplicationId: "app-id",
					DeployTargets: []string{"default"},
				},
				Stage: &model.PipelineStage{
					Id:   "stage-id",
					Name: "K8S_SYNC",
				},
				StageConfig:             []byte(``),
				RunningDeploymentSource: nil,
				TargetDeploymentSource: &deployment.DeploymentSource{
					ApplicationDirectory:      running,
					CommitHash:                "0123456789",
					ApplicationConfig:         runningCfg,
					ApplicationConfigFilename: "app.pipecd.yaml",
				},
			},
		}

		resp, err := svc.ExecuteStage(ctx, runningRequest)

		require.NoError(t, err)
		require.Equal(t, model.StageStatus_STAGE_SUCCESS.String(), resp.GetStatus().String())

		service, err := dynamicClient.Resource(schema.GroupVersionResource{Group: "", Version: "v1", Resource: "services"}).Namespace("default").Get(context.Background(), "simple", metav1.GetOptions{})
		require.NoError(t, err)

		require.Equal(t, "piped", service.GetLabels()["pipecd.dev/managed-by"])
		require.Equal(t, "piped-id", service.GetLabels()["pipecd.dev/piped"])
		require.Equal(t, "app-id", service.GetLabels()["pipecd.dev/application"])
		require.Equal(t, "0123456789", service.GetLabels()["pipecd.dev/commit-hash"])

		require.Equal(t, "simple", service.GetName())
		require.Equal(t, "piped", service.GetAnnotations()["pipecd.dev/managed-by"])
		require.Equal(t, "piped-id", service.GetAnnotations()["pipecd.dev/piped"])
		require.Equal(t, "app-id", service.GetAnnotations()["pipecd.dev/application"])
		require.Equal(t, "v1", service.GetAnnotations()["pipecd.dev/original-api-version"])
		require.Equal(t, ":Service::simple", service.GetAnnotations()["pipecd.dev/resource-key"]) // This assertion differs from the non-plugin-arched piped's Kubernetes platform provider, but we decided to change this behavior.
		require.Equal(t, "0123456789", service.GetAnnotations()["pipecd.dev/commit-hash"])
	})
	require.Truef(t, ok, "expected prepare to succeed")

	t.Run("run with prune", func(t *testing.T) {

		// prepare the request to ensure the running deployment exists

		target := filepath.Join("./", "testdata", "prune", "target")

		// read the running application config from the testdata file
		targetCfg, err := os.ReadFile(filepath.Join(target, "app.pipecd.yaml"))
		require.NoError(t, err)

		// prepare the request to ensure the running deployment exists
		targetRequest := &deployment.ExecuteStageRequest{
			Input: &deployment.ExecutePluginInput{
				Deployment: &model.Deployment{
					PipedId:       "piped-id",
					ApplicationId: "app-id",
					DeployTargets: []string{"default"},
				},
				Stage: &model.PipelineStage{
					Id:   "stage-id",
					Name: "K8S_SYNC",
				},
				StageConfig: []byte(``),
				RunningDeploymentSource: &deployment.DeploymentSource{
					ApplicationDirectory:      running,
					CommitHash:                "0123456789",
					ApplicationConfig:         runningCfg,
					ApplicationConfigFilename: "app.pipecd.yaml",
				},
				TargetDeploymentSource: &deployment.DeploymentSource{
					ApplicationDirectory:      target,
					CommitHash:                "0012345678",
					ApplicationConfig:         targetCfg,
					ApplicationConfigFilename: "app.pipecd.yaml",
				},
			},
		}

		resp, err := svc.ExecuteStage(ctx, targetRequest)
		require.NoError(t, err)
		require.Equal(t, model.StageStatus_STAGE_SUCCESS.String(), resp.GetStatus().String())

		_, err = dynamicClient.Resource(schema.GroupVersionResource{Group: "", Version: "v1", Resource: "services"}).Namespace("default").Get(context.Background(), "simple", metav1.GetOptions{})
		require.Error(t, err)
		require.Truef(t, apierrors.IsNotFound(err), "expected error to be NotFound, but got %v", err)
	})
}

func TestDeploymentService_executeK8sSyncStage_withPrune_changesNamespace(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	// initialize tool registry
	testRegistry, err := toolregistrytest.NewToolRegistry(t)
	require.NoError(t, err)

	// initialize plugin config and dynamic client for assertions with envtest
	pluginCfg, dynamicClient := setupTestPluginConfigAndDynamicClient(t)

	svc := NewDeploymentService(pluginCfg, zaptest.NewLogger(t), testRegistry, logpersistertest.NewTestLogPersister(t))

	running := filepath.Join("./", "testdata", "prune_with_change_namespace", "running")

	// read the running application config from the example file
	runningCfg, err := os.ReadFile(filepath.Join(running, "app.pipecd.yaml"))
	require.NoError(t, err)

	ok := t.Run("prepare", func(t *testing.T) {
		// prepare the request to ensure the running deployment exists
		runningRequest := &deployment.ExecuteStageRequest{
			Input: &deployment.ExecutePluginInput{
				Deployment: &model.Deployment{
					PipedId:       "piped-id",
					ApplicationId: "app-id",
					DeployTargets: []string{"default"},
				},
				Stage: &model.PipelineStage{
					Id:   "stage-id",
					Name: "K8S_SYNC",
				},
				StageConfig:             []byte(``),
				RunningDeploymentSource: nil,
				TargetDeploymentSource: &deployment.DeploymentSource{
					ApplicationDirectory:      running,
					CommitHash:                "0123456789",
					ApplicationConfig:         runningCfg,
					ApplicationConfigFilename: "app.pipecd.yaml",
				},
			},
		}

		resp, err := svc.ExecuteStage(ctx, runningRequest)

		require.NoError(t, err)
		require.Equal(t, model.StageStatus_STAGE_SUCCESS.String(), resp.GetStatus().String())

		service, err := dynamicClient.Resource(schema.GroupVersionResource{Group: "", Version: "v1", Resource: "services"}).Namespace("test-1").Get(context.Background(), "simple", metav1.GetOptions{})
		require.NoError(t, err)

		require.Equal(t, "piped", service.GetLabels()["pipecd.dev/managed-by"])
		require.Equal(t, "piped-id", service.GetLabels()["pipecd.dev/piped"])
		require.Equal(t, "app-id", service.GetLabels()["pipecd.dev/application"])
		require.Equal(t, "0123456789", service.GetLabels()["pipecd.dev/commit-hash"])

		require.Equal(t, "simple", service.GetName())
		require.Equal(t, "piped", service.GetAnnotations()["pipecd.dev/managed-by"])
		require.Equal(t, "piped-id", service.GetAnnotations()["pipecd.dev/piped"])
		require.Equal(t, "app-id", service.GetAnnotations()["pipecd.dev/application"])
		require.Equal(t, "v1", service.GetAnnotations()["pipecd.dev/original-api-version"])
		require.Equal(t, "0123456789", service.GetAnnotations()["pipecd.dev/commit-hash"])
		require.Equal(t, ":Service:test-1:simple", service.GetAnnotations()["pipecd.dev/resource-key"])
	})
	require.Truef(t, ok, "expected prepare to succeed")

	t.Run("run with prune", func(t *testing.T) {
		target := filepath.Join("./", "testdata", "prune_with_change_namespace", "target")

		// read the running application config from the example file
		targetCfg, err := os.ReadFile(filepath.Join(target, "app.pipecd.yaml"))
		require.NoError(t, err)

		// prepare the request to ensure the running deployment exists
		targetRequest := &deployment.ExecuteStageRequest{
			Input: &deployment.ExecutePluginInput{
				Deployment: &model.Deployment{
					PipedId:       "piped-id",
					ApplicationId: "app-id",
					DeployTargets: []string{"default"},
				},
				Stage: &model.PipelineStage{
					Id:   "stage-id",
					Name: "K8S_SYNC",
				},
				StageConfig: []byte(``),
				RunningDeploymentSource: &deployment.DeploymentSource{
					ApplicationDirectory:      running,
					CommitHash:                "0123456789",
					ApplicationConfig:         runningCfg,
					ApplicationConfigFilename: "app.pipecd.yaml",
				},
				TargetDeploymentSource: &deployment.DeploymentSource{
					ApplicationDirectory:      target,
					CommitHash:                "0012345678",
					ApplicationConfig:         targetCfg,
					ApplicationConfigFilename: "app.pipecd.yaml",
				},
			},
		}

		resp, err := svc.ExecuteStage(ctx, targetRequest)
		require.NoError(t, err)
		require.Equal(t, model.StageStatus_STAGE_SUCCESS.String(), resp.GetStatus().String())

		// The service should be removed from the previous namespace
		_, err = dynamicClient.Resource(schema.GroupVersionResource{Group: "", Version: "v1", Resource: "services"}).Namespace("test-1").Get(context.Background(), "simple", metav1.GetOptions{})
		require.Error(t, err)
		require.Truef(t, apierrors.IsNotFound(err), "expected error to be NotFound, but got %v", err)

		// The service should be created in the new namespace
		service, err := dynamicClient.Resource(schema.GroupVersionResource{Group: "", Version: "v1", Resource: "services"}).Namespace("test-2").Get(context.Background(), "simple", metav1.GetOptions{})
		require.NoError(t, err)

		require.Equal(t, "piped", service.GetLabels()["pipecd.dev/managed-by"])
		require.Equal(t, "piped-id", service.GetLabels()["pipecd.dev/piped"])
		require.Equal(t, "app-id", service.GetLabels()["pipecd.dev/application"])
		require.Equal(t, "0012345678", service.GetLabels()["pipecd.dev/commit-hash"])

		require.Equal(t, "simple", service.GetName())
		require.Equal(t, "piped", service.GetAnnotations()["pipecd.dev/managed-by"])
		require.Equal(t, "piped-id", service.GetAnnotations()["pipecd.dev/piped"])
		require.Equal(t, "app-id", service.GetAnnotations()["pipecd.dev/application"])
		require.Equal(t, "v1", service.GetAnnotations()["pipecd.dev/original-api-version"])
		require.Equal(t, "0012345678", service.GetAnnotations()["pipecd.dev/commit-hash"])
		require.Equal(t, ":Service:test-2:simple", service.GetAnnotations()["pipecd.dev/resource-key"])
	})
}

func TestDeploymentService_executeK8sSyncStage_withPrune_clusterScoped(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	// initialize tool registry
	testRegistry, err := toolregistrytest.NewToolRegistry(t)
	require.NoError(t, err)

	// initialize plugin config and dynamic client for assertions with envtest
	pluginCfg, dynamicClient := setupTestPluginConfigAndDynamicClient(t)

	svc := NewDeploymentService(pluginCfg, zaptest.NewLogger(t), testRegistry, logpersistertest.NewTestLogPersister(t))

	// prepare the custom resource definition
	prepare := filepath.Join("./", "testdata", "prune_cluster_scoped_resource", "prepare")

	prepareCfg, err := os.ReadFile(filepath.Join(prepare, "app.pipecd.yaml"))
	require.NoError(t, err)

	ok := t.Run("prepare crd", func(t *testing.T) {
		prepareRequest := &deployment.ExecuteStageRequest{
			Input: &deployment.ExecutePluginInput{
				Deployment: &model.Deployment{
					PipedId:       "piped-id",
					ApplicationId: "prepare-app-id",
					DeployTargets: []string{"default"},
				},
				Stage: &model.PipelineStage{
					Id:   "stage-id",
					Name: "K8S_SYNC",
				},
				StageConfig:             []byte(``),
				RunningDeploymentSource: nil,
				TargetDeploymentSource: &deployment.DeploymentSource{
					ApplicationDirectory:      prepare,
					CommitHash:                "0123456789",
					ApplicationConfig:         prepareCfg,
					ApplicationConfigFilename: "app.pipecd.yaml",
				},
			},
		}

		resp, err := svc.ExecuteStage(ctx, prepareRequest)

		require.NoError(t, err)
		require.Equal(t, model.StageStatus_STAGE_SUCCESS.String(), resp.GetStatus().String())
	})
	require.Truef(t, ok, "expected prepare to succeed")

	// prepare the running resources
	running := filepath.Join("./", "testdata", "prune_cluster_scoped_resource", "running")

	// read the running application config from the example file
	runningCfg, err := os.ReadFile(filepath.Join(running, "app.pipecd.yaml"))
	require.NoError(t, err)

	ok = t.Run("prepare running", func(t *testing.T) {
		// prepare the request to ensure the running deployment exists
		runningRequest := &deployment.ExecuteStageRequest{
			Input: &deployment.ExecutePluginInput{
				Deployment: &model.Deployment{
					PipedId:       "piped-id",
					ApplicationId: "app-id",
					DeployTargets: []string{"default"},
				},
				Stage: &model.PipelineStage{
					Id:   "stage-id",
					Name: "K8S_SYNC",
				},
				StageConfig:             []byte(``),
				RunningDeploymentSource: nil,
				TargetDeploymentSource: &deployment.DeploymentSource{
					ApplicationDirectory:      running,
					CommitHash:                "0123456789",
					ApplicationConfig:         runningCfg,
					ApplicationConfigFilename: "app.pipecd.yaml",
				},
			},
		}

		resp, err := svc.ExecuteStage(ctx, runningRequest)

		require.NoError(t, err)
		require.Equal(t, model.StageStatus_STAGE_SUCCESS.String(), resp.GetStatus().String())

		// The my-new-cron-object/my-new-cron-object-2/my-new-cron-object-v1beta1 should be created
		_, err = dynamicClient.Resource(schema.GroupVersionResource{Group: "stable.example.com", Version: "v1", Resource: "crontabs"}).Get(context.Background(), "my-new-cron-object", metav1.GetOptions{})
		require.NoError(t, err)
		_, err = dynamicClient.Resource(schema.GroupVersionResource{Group: "stable.example.com", Version: "v1", Resource: "crontabs"}).Get(context.Background(), "my-new-cron-object-2", metav1.GetOptions{})
		require.NoError(t, err)
		_, err = dynamicClient.Resource(schema.GroupVersionResource{Group: "stable.example.com", Version: "v1", Resource: "crontabs"}).Get(context.Background(), "my-new-cron-object-v1beta1", metav1.GetOptions{})
		require.NoError(t, err)
	})
	require.Truef(t, ok, "expected prepare to succeed")

	t.Run("sync", func(t *testing.T) {
		// sync the target resources and assert the prune behavior
		target := filepath.Join("./", "testdata", "prune_cluster_scoped_resource", "target")

		// read the running application config from the example file
		targetCfg, err := os.ReadFile(filepath.Join(target, "app.pipecd.yaml"))
		require.NoError(t, err)

		// prepare the request to ensure the running deployment exists
		targetRequest := &deployment.ExecuteStageRequest{
			Input: &deployment.ExecutePluginInput{
				Deployment: &model.Deployment{
					PipedId:       "piped-id",
					ApplicationId: "app-id",
					DeployTargets: []string{"default"},
				},
				Stage: &model.PipelineStage{
					Id:   "stage-id",
					Name: "K8S_SYNC",
				},
				StageConfig: []byte(``),
				RunningDeploymentSource: &deployment.DeploymentSource{
					ApplicationDirectory:      running,
					CommitHash:                "0123456789",
					ApplicationConfig:         runningCfg,
					ApplicationConfigFilename: "app.pipecd.yaml",
				},
				TargetDeploymentSource: &deployment.DeploymentSource{
					ApplicationDirectory:      target,
					CommitHash:                "0012345678",
					ApplicationConfig:         targetCfg,
					ApplicationConfigFilename: "app.pipecd.yaml",
				},
			},
		}

		resp, err := svc.ExecuteStage(ctx, targetRequest)
		require.NoError(t, err)
		require.Equal(t, model.StageStatus_STAGE_SUCCESS.String(), resp.GetStatus().String())

		// The my-new-cron-object should not be removed
		_, err = dynamicClient.Resource(schema.GroupVersionResource{Group: "stable.example.com", Version: "v1", Resource: "crontabs"}).Get(context.Background(), "my-new-cron-object", metav1.GetOptions{})
		require.NoError(t, err)
		// The my-new-cron-object-2 should be removed
		_, err = dynamicClient.Resource(schema.GroupVersionResource{Group: "stable.example.com", Version: "v1", Resource: "crontabs"}).Get(context.Background(), "my-new-cron-object-2", metav1.GetOptions{})
		require.Error(t, err)
		require.Truef(t, apierrors.IsNotFound(err), "expected error to be NotFound, but got %v", err)
		// The my-new-cron-object-v1beta1 should be removed
		_, err = dynamicClient.Resource(schema.GroupVersionResource{Group: "stable.example.com", Version: "v1", Resource: "crontabs"}).Get(context.Background(), "my-new-cron-object-v1beta1", metav1.GetOptions{})
		require.Error(t, err)
		require.Truef(t, apierrors.IsNotFound(err), "expected error to be NotFound, but got %v", err)
	})
}
