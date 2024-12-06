package kubernetes

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
)

func TestKubectl_Apply(t *testing.T) {
	t.Setenv("KUBEBUILDER_ASSETS", "/Users/s14218/Library/Application Support/io.kubebuilder.envtest/k8s/1.30.0-darwin-arm64")

	tEnv := new(envtest.Environment)
	cfg, err := tEnv.Start()
	if err != nil {
		panic(err)
	}
	defer tEnv.Stop()

	workspace := t.TempDir()
	kubeconfig := filepath.Join(workspace, "kubeconfig")
	kubecfgBytes, err := kubeconfigFromRestConfig(cfg)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(kubeconfig, []byte(kubecfgBytes), 0755))

	kubectl := Kubectl{execPath: filepath.Join(tEnv.BinaryAssetsDirectory, "kubectl")}
	manifest, err := ParseManifests(`
apiVersion: apps/v1
kind: Deployment
metadata:
  name: simple
  labels:
    app: simple
spec:
  replicas: 2
  selector:
    matchLabels:
      app: simple
      pipecd.dev/variant: primary
  template:
    metadata:
      labels:
        app: simple
        pipecd.dev/variant: primary
      annotations:
        sidecar.istio.io/inject: "false"
    spec:
      containers:
      - name: helloworld
        image: ghcr.io/pipe-cd/helloworld:v0.32.0
        args:
          - server
        ports:
        - containerPort: 9085
`)
	require.NoError(t, err)

	err = kubectl.Apply(context.Background(), kubeconfig, "default", manifest[0])
	require.NoError(t, err)

	// kubectlのApply結果を確認する client-goを使って
	dynamicClient, err := dynamic.NewForConfig(cfg)
	require.NoError(t, err)
	obj, err := dynamicClient.Resource(manifest[0].u.GroupVersionKind().GroupVersion().WithResource("deployments")).Namespace("default").Get(context.Background(), "simple", metav1.GetOptions{})
	require.NoError(t, err)

	assert.Equal(t, "Deployment", obj.GetKind())
	assert.Equal(t, "simple", obj.GetName())
	assert.Equal(t, "simple", obj.GetLabels()["app"])
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
