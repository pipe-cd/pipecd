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
	"encoding/json"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"sigs.k8s.io/controller-runtime/pkg/envtest"

	kubeConfigPkg "github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/config"
	config "github.com/pipe-cd/pipecd/pkg/configv1"
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
