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
				Revision:                  "0123456789",
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
	assert.Equal(t, "piped", deployment.GetAnnotations()["pipecd.dev/managed-by"])
	assert.Equal(t, "piped-id", deployment.GetAnnotations()["pipecd.dev/piped"])
	assert.Equal(t, "app-id", deployment.GetAnnotations()["pipecd.dev/application"])
	assert.Equal(t, "apps/v1", deployment.GetAnnotations()["pipecd.dev/original-api-version"])
	assert.Equal(t, "apps/v1:Deployment::simple", deployment.GetAnnotations()["pipecd.dev/resource-key"]) // This assertion differs from the non-plugin-arched piped's Kubernetes platform provider, but we decided to change this behavior.
	assert.Equal(t, "0123456789", deployment.GetAnnotations()["pipecd.dev/commit-hash"])

}

func TestDeploymentService_executeK8sSyncStage_withInputNamespace(t *testing.T) {
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
				Revision:                  "0123456789",
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

	assert.Equal(t, "simple", deployment.GetName())
	assert.Equal(t, "simple", deployment.GetLabels()["app"])
	assert.Equal(t, "piped", deployment.GetAnnotations()["pipecd.dev/managed-by"])
	assert.Equal(t, "piped-id", deployment.GetAnnotations()["pipecd.dev/piped"])
	assert.Equal(t, "app-id", deployment.GetAnnotations()["pipecd.dev/application"])
	assert.Equal(t, "apps/v1", deployment.GetAnnotations()["pipecd.dev/original-api-version"])
	assert.Equal(t, "apps/v1:Deployment::simple", deployment.GetAnnotations()["pipecd.dev/resource-key"]) // This assertion differs from the non-plugin-arched piped's Kubernetes platform provider, but we decided to change this behavior.
	assert.Equal(t, "0123456789", deployment.GetAnnotations()["pipecd.dev/commit-hash"])
}
