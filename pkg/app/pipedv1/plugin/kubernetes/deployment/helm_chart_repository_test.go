// Copyright 2026 The PipeCD Authors.
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
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
	"github.com/pipe-cd/piped-plugin-sdk-go/logpersister/logpersistertest"
	"github.com/pipe-cd/piped-plugin-sdk-go/toolregistry/toolregistrytest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	kubeConfigPkg "github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/config"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/provider"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/toolregistry"
)

// The chart repository is served over HTTPS with a self-signed certificate,
// so Helm can only download the chart when the repository is marked as insecure.
func TestPlugin_DetermineVersions_insecureChartRepository(t *testing.T) {
	// Not parallel: the Helm home directories are set through environment variables.
	home := t.TempDir()
	t.Setenv("HELM_CONFIG_HOME", filepath.Join(home, "config"))
	t.Setenv("HELM_CACHE_HOME", filepath.Join(home, "cache"))
	t.Setenv("HELM_DATA_HOME", filepath.Join(home, "data"))

	ctx := context.Background()
	logger := zaptest.NewLogger(t)

	testRegistry := toolregistrytest.NewTestToolRegistry(t)
	helmPath, err := toolregistry.NewRegistry(testRegistry).Helm(ctx, "")
	require.NoError(t, err)

	helm := func(args ...string) {
		out, err := exec.CommandContext(ctx, helmPath, args...).CombinedOutput()
		require.NoError(t, err, string(out))
	}

	// Package a chart and serve it from a chart repository.
	workDir := t.TempDir()
	repoDir := t.TempDir()
	helm("create", filepath.Join(workDir, "mychart"))
	helm("package", filepath.Join(workDir, "mychart"), "--destination", repoDir)
	server := httptest.NewTLSServer(http.FileServer(http.Dir(repoDir)))
	t.Cleanup(server.Close)
	helm("repo", "index", repoDir, "--url", server.URL)

	// Add the repository the same way the plugin does when it starts.
	pluginCfg := &kubeConfigPkg.KubernetesPluginConfig{
		ChartRepositories: []kubeConfigPkg.HelmChartRepository{
			{
				Type:     kubeConfigPkg.HTTPHelmChartRepository,
				Name:     "insecure-repo",
				Address:  server.URL,
				Insecure: true,
			},
		},
	}
	for _, repo := range pluginCfg.HTTPHelmChartRepositories() {
		require.NoError(t, provider.NewHelm(helmPath, logger).AddRepository(ctx, repo))
	}

	appDir := t.TempDir()
	appCfgFile := filepath.Join(appDir, "app.pipecd.yaml")
	require.NoError(t, os.WriteFile(appCfgFile, []byte(`apiVersion: pipecd.dev/v1beta1
kind: KubernetesApp
spec:
  name: mychart
  plugins:
    kubernetes:
      input:
        helmChart:
          repository: insecure-repo
          name: mychart
          version: 0.1.0
`), 0o644))
	appCfg := sdk.LoadApplicationConfigForTest[kubeConfigPkg.KubernetesApplicationSpec](t, appCfgFile, "kubernetes")

	input := &sdk.DetermineVersionsInput[kubeConfigPkg.KubernetesApplicationSpec]{
		Request: sdk.DetermineVersionsRequest[kubeConfigPkg.KubernetesApplicationSpec]{
			Deployment: sdk.Deployment{
				PipedID:         "piped-id",
				ApplicationID:   "app-id",
				ApplicationName: "mychart",
			},
			DeploymentSource: sdk.DeploymentSource[kubeConfigPkg.KubernetesApplicationSpec]{
				ApplicationDirectory:      appDir,
				CommitHash:                "0123456789",
				ApplicationConfig:         appCfg,
				ApplicationConfigFilename: "app.pipecd.yaml",
			},
		},
		Client: sdk.NewClient(nil, "kubernetes", "", "", logpersistertest.NewTestLogPersister(t), testRegistry),
		Logger: logger,
	}

	plugin := &Plugin{}
	resp, err := plugin.DetermineVersions(ctx, pluginCfg, input)
	require.NoError(t, err)
	require.Len(t, resp.Versions, 1)
	assert.Equal(t, "nginx", resp.Versions[0].Name)
}
