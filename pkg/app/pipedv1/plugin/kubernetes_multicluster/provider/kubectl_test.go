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

package provider

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
)

func setupEnvTest(t *testing.T) *rest.Config {
	t.Helper()

	tEnv := new(envtest.Environment)
	kubeCfg, err := tEnv.Start()
	require.NoError(t, err)
	t.Cleanup(func() { tEnv.Stop() })

	return kubeCfg
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

func TestKubectl_GetAll(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name      string
		namespace string
		selectors []string
		manifests []Manifest
		want      []ResourceKey
		wantErr   bool
	}{
		{
			name:      "get all namespace-scoped resources in defau namespace",
			namespace: "default",
			selectors: []string{"env=test"},
			manifests: mustParseManifests(t, `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  labels:
    env: test
spec:
  selector:
    matchLabels:
      env: test
  template:
    metadata:
      labels:
        env: test
    spec:
      containers:
      - name: nginx
        image: nginx:1.19.3
---
apiVersion: v1
kind: Service
metadata:
  name: nginx-service
  labels:
    env: test
spec:
  selector:
    env: test
  ports:
    - port: 80
      targetPort: 8080
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: test-config
  labels:
    env: test
data:
  key: value
`),
			want: []ResourceKey{
				{
					groupKind: schema.GroupKind{
						Group: "apps",
						Kind:  "Deployment",
					},
					namespace: "default",
					name:      "nginx-deployment",
				},
				{
					groupKind: schema.GroupKind{
						Group: "",
						Kind:  "Service",
					},
					namespace: "default",
					name:      "nginx-service",
				},
				{
					groupKind: schema.GroupKind{
						Group: "",
						Kind:  "ConfigMap",
					},
					namespace: "default",
					name:      "test-config",
				},
			},
			wantErr: false,
		},
		{
			name:      "get all namespace-scoped resources for all namespaces when namespace is empty",
			namespace: "",
			selectors: []string{"env=test"},
			manifests: mustParseManifests(t, `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  namespace: test
  labels:
    env: test
spec:
  selector:
    matchLabels:
      env: test
  template:
    metadata:
      labels:
        env: test
    spec:
      containers:
      - name: nginx
        image: nginx:1.19.3
---
apiVersion: v1
kind: Service
metadata:
  name: nginx-service
  namespace: test
  labels:
    env: test
spec:
  selector:
    env: test
  ports:
    - port: 80
      targetPort: 8080
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: test-config
  labels:
    env: test
data:
  key: value
`),
			want: []ResourceKey{
				{
					groupKind: schema.GroupKind{
						Group: "apps",
						Kind:  "Deployment",
					},
					namespace: "test",
					name:      "nginx-deployment",
				},
				{
					groupKind: schema.GroupKind{
						Group: "",
						Kind:  "Service",
					},
					namespace: "test",
					name:      "nginx-service",
				},
				{
					groupKind: schema.GroupKind{
						Group: "",
						Kind:  "ConfigMap",
					},
					namespace: "default",
					name:      "test-config",
				},
			},
			wantErr: false,
		},
		{
			name:      "get all namespace-scoped resources in a specific namespace",
			namespace: "test",
			selectors: []string{"env=test"},
			manifests: mustParseManifests(t, `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  namespace: test
  labels:
    env: test
spec:
  selector:
    matchLabels:
      env: test
  template:
    metadata:
      labels:
        env: test
    spec:
      containers:
      - name: nginx
        image: nginx:1.19.3
---
apiVersion: v1
kind: Service
metadata:
  name: nginx-service
  namespace: test
  labels:
    env: test
spec:
  selector:
    env: test
  ports:
    - port: 80
      targetPort: 8080
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: test-config
  labels:
    env: test
data:
  key: value			
`),
			want: []ResourceKey{
				{
					groupKind: schema.GroupKind{
						Group: "apps",
						Kind:  "Deployment",
					},
					namespace: "test",
					name:      "nginx-deployment",
				},
				{
					groupKind: schema.GroupKind{
						Group: "",
						Kind:  "Service",
					},
					namespace: "test",
					name:      "nginx-service",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// prepare k8s cluster
			restConfig := setupEnvTest(t)
			kubeconfig, err := kubeconfigFromRestConfig(restConfig)
			require.NoError(t, err)

			kubeconfigPath := path.Join(t.TempDir(), "kubeconfig")
			err = os.WriteFile(kubeconfigPath, []byte(kubeconfig), 0755)
			require.NoError(t, err)

			kubectl := NewKubectl(path.Join(os.Getenv("KUBEBUILDER_ASSETS"), "kubectl"))

			// create namespace for testing
			err = kubectl.CreateNamespace(t.Context(), kubeconfigPath, "test")
			require.NoError(t, err)

			// apply resources before testing
			for _, m := range tt.manifests {
				// set empty namespace to use defined one in the manifest
				err := kubectl.Apply(t.Context(), kubeconfigPath, "", m)
				require.NoError(t, err)
			}

			got, err := kubectl.GetAll(t.Context(), kubeconfigPath, tt.namespace, tt.selectors...)
			assert.Equal(t, tt.wantErr, err != nil)

			keys := make([]ResourceKey, 0, len(tt.want))
			for _, m := range got {
				keys = append(keys, m.Key())
			}

			assert.ElementsMatch(t, tt.want, keys)
		})
	}

}
