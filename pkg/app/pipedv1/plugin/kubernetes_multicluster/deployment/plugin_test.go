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
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"sigs.k8s.io/controller-runtime/pkg/envtest"

	kubeConfigPkg "github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes_multicluster/config"
)

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

func setupTestDeployTargetConfig(t *testing.T, kubeCfg *rest.Config) *kubeConfigPkg.KubernetesDeployTargetConfig {
	t.Helper()

	kubeconfig, err := kubeconfigFromRestConfig(kubeCfg)
	require.NoError(t, err)

	workDir := t.TempDir()
	kubeconfigPath := path.Join(workDir, "kubeconfig")
	err = os.WriteFile(kubeconfigPath, []byte(kubeconfig), 0755)
	require.NoError(t, err)

	return &kubeConfigPkg.KubernetesDeployTargetConfig{
		KubeConfigPath: kubeconfigPath,
	}
}

func setupTestDeployTargetConfigAndDynamicClient(t *testing.T) (*kubeConfigPkg.KubernetesDeployTargetConfig, dynamic.Interface) {
	t.Helper()

	kubeCfg := setupEnvTest(t)
	deployTargetCfg := setupTestDeployTargetConfig(t, kubeCfg)

	dynamicClient, err := dynamic.NewForConfig(kubeCfg)
	require.NoError(t, err)

	return deployTargetCfg, dynamicClient
}

type cluster struct {
	name string
	cli  dynamic.Interface
	dtc  *kubeConfigPkg.KubernetesDeployTargetConfig
}

func setupCluster(t *testing.T, name string) *cluster {
	t.Helper()

	clusterCfg := setupEnvTest(t)
	dtc := setupTestDeployTargetConfig(t, clusterCfg)

	cli, err := dynamic.NewForConfig(clusterCfg)
	require.NoError(t, err)

	return &cluster{
		name: name,
		cli:  cli,
		dtc:  dtc,
	}
}
