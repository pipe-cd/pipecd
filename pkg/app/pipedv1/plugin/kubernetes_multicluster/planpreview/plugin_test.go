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

package planpreview

import (
	"context"
	"path/filepath"
	"testing"

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
	"github.com/pipe-cd/piped-plugin-sdk-go/logpersister/logpersistertest"
	"github.com/pipe-cd/piped-plugin-sdk-go/toolregistry/toolregistrytest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	kubeconfig "github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes_multicluster/config"
)

const (
	pluginName    = "kubernetes_multicluster"
	runningCommit = "abc000"
	targetCommit  = "abc123"
)

func makeDeployTargets(names ...string) []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig] {
	dts := make([]*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig], 0, len(names))
	for _, name := range names {
		dts = append(dts, &sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]{Name: name})
	}
	return dts
}

func makeInput(
	t *testing.T,
	appConfigFile string,
	targetAppDir string,
	runningAppDir string,
) *sdk.GetPlanPreviewInput[kubeconfig.KubernetesApplicationSpec] {
	t.Helper()

	appCfg := sdk.LoadApplicationConfigForTest[kubeconfig.KubernetesApplicationSpec](t, appConfigFile, pluginName)
	testRegistry := toolregistrytest.NewTestToolRegistry(t)

	targetDS := sdk.DeploymentSource[kubeconfig.KubernetesApplicationSpec]{
		ApplicationDirectory:      targetAppDir,
		CommitHash:                targetCommit,
		ApplicationConfig:         appCfg,
		ApplicationConfigFilename: "app.pipecd.yaml",
	}

	input := &sdk.GetPlanPreviewInput[kubeconfig.KubernetesApplicationSpec]{
		Request: sdk.GetPlanPreviewRequest[kubeconfig.KubernetesApplicationSpec]{
			ApplicationID:          "app-id",
			ApplicationName:        "simple",
			TargetDeploymentSource: targetDS,
		},
		Client: sdk.NewClient(nil, pluginName, "app-id", "", logpersistertest.NewTestLogPersister(t), testRegistry),
		Logger: zaptest.NewLogger(t),
	}

	if runningAppDir != "" {
		input.Request.RunningDeploymentSource = sdk.DeploymentSource[kubeconfig.KubernetesApplicationSpec]{
			ApplicationDirectory:      runningAppDir,
			CommitHash:                runningCommit,
			ApplicationConfig:         appCfg,
			ApplicationConfigFilename: "app.pipecd.yaml",
		}
	}

	return input
}

func TestPlugin_GetPlanPreview_SingleTarget(t *testing.T) {
	t.Parallel()

	appCfgFile := filepath.Join("testdata", "single", "app.pipecd.yaml")
	runningDir := filepath.Join("testdata", "single", "running")
	targetDir := filepath.Join("testdata", "single", "target")
	p := &Plugin{}

	tests := []struct {
		name          string
		targetDir     string
		runningDir    string
		wantNoChange  bool
		wantSummary   string
		wantInDetails []string
	}{
		{
			name:         "first deployment (no running source)",
			targetDir:    targetDir,
			runningDir:   "",
			wantNoChange: false,
			wantSummary:  "1 added manifests, 0 changed manifests, 0 deleted manifests",
			wantInDetails: []string{
				"simple",
				"Deployment",
			},
		},
		{
			name:         "no change (same files)",
			targetDir:    runningDir,
			runningDir:   runningDir,
			wantNoChange: true,
			wantSummary:  "No changes were detected",
		},
		{
			name:         "image tag changed",
			targetDir:    targetDir,
			runningDir:   runningDir,
			wantNoChange: false,
			wantSummary:  "0 added manifests, 1 changed manifests, 0 deleted manifests",
			wantInDetails: []string{
				"v0.1.0",
				"v0.2.0",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			input := makeInput(t, appCfgFile, tc.targetDir, tc.runningDir)

			resp, err := p.GetPlanPreview(context.Background(), nil, makeDeployTargets("default"), input)
			require.NoError(t, err)
			require.Len(t, resp.Results, 1)

			result := resp.Results[0]
			assert.Equal(t, tc.wantNoChange, result.NoChange)
			assert.Equal(t, tc.wantSummary, result.Summary)
			assert.Equal(t, "diff", result.DiffLanguage)

			for _, s := range tc.wantInDetails {
				assert.Contains(t, string(result.Details), s)
			}
			if tc.wantNoChange {
				assert.Nil(t, result.Details)
			}
		})
	}
}

func TestPlugin_GetPlanPreview_MultiTarget(t *testing.T) {
	t.Parallel()

	appCfgFile := filepath.Join("testdata", "multi", "app.pipecd.yaml")
	runningDir := filepath.Join("testdata", "multi", "running")
	targetDir := filepath.Join("testdata", "multi", "target")
	p := &Plugin{}
	dts := makeDeployTargets("cluster1", "cluster2")

	tests := []struct {
		name   string
		checks []struct {
			deployTarget  string
			wantNoChange  bool
			wantSummary   string
			wantInDetails []string
		}
	}{
		{
			name: "cluster1 changed, cluster2 unchanged",
			checks: []struct {
				deployTarget  string
				wantNoChange  bool
				wantSummary   string
				wantInDetails []string
			}{
				{
					deployTarget:  "cluster1",
					wantNoChange:  false,
					wantSummary:   "0 added manifests, 1 changed manifests, 0 deleted manifests",
					wantInDetails: []string{"v0.1.0", "v0.2.0"},
				},
				{
					deployTarget: "cluster2",
					wantNoChange: true,
					wantSummary:  "No changes were detected",
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			input := makeInput(t, appCfgFile, targetDir, runningDir)

			resp, err := p.GetPlanPreview(context.Background(), nil, dts, input)
			require.NoError(t, err)
			require.Len(t, resp.Results, 2)

			for _, check := range tc.checks {
				var result *sdk.PlanPreviewResult
				for i := range resp.Results {
					if resp.Results[i].DeployTarget == check.deployTarget {
						result = &resp.Results[i]
						break
					}
				}
				require.NotNil(t, result, "result for deploy target %q not found", check.deployTarget)

				assert.Equal(t, check.wantNoChange, result.NoChange)
				assert.Equal(t, check.wantSummary, result.Summary)
				assert.Equal(t, "diff", result.DiffLanguage)

				for _, s := range check.wantInDetails {
					assert.Contains(t, string(result.Details), s)
				}
				if check.wantNoChange {
					assert.Nil(t, result.Details)
				}
			}
		})
	}
}

func TestPlugin_GetPlanPreview_MultiTarget_KustomizeDir(t *testing.T) {
	t.Parallel()

	appCfgFile := filepath.Join("testdata", "kustomize", "app.pipecd.yaml")
	runningDir := filepath.Join("testdata", "kustomize", "running")
	targetDir := filepath.Join("testdata", "kustomize", "target")
	p := &Plugin{}
	dts := makeDeployTargets("cluster1", "cluster2")

	input := makeInput(t, appCfgFile, targetDir, runningDir)

	resp, err := p.GetPlanPreview(context.Background(), nil, dts, input)
	require.NoError(t, err)
	require.Len(t, resp.Results, 2)

	results := make(map[string]sdk.PlanPreviewResult, len(resp.Results))
	for _, r := range resp.Results {
		results[r.DeployTarget] = r
	}

	// Only the overlay of cluster1 changes the image tag.
	cluster1, ok := results["cluster1"]
	require.True(t, ok, "result for deploy target cluster1 not found")
	assert.False(t, cluster1.NoChange)
	assert.Equal(t, "0 added manifests, 1 changed manifests, 0 deleted manifests", cluster1.Summary)
	assert.Contains(t, string(cluster1.Details), "v0.1.0")
	assert.Contains(t, string(cluster1.Details), "v0.2.0")

	cluster2, ok := results["cluster2"]
	require.True(t, ok, "result for deploy target cluster2 not found")
	assert.True(t, cluster2.NoChange)
	assert.Equal(t, "No changes were detected", cluster2.Summary)
	assert.Nil(t, cluster2.Details)
}

func TestBuildLoaderInput_KustomizeOverrides(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		input       kubeconfig.KubernetesDeploymentInput
		multiTarget *kubeconfig.KubernetesMultiTarget
		wantDir     string
		wantVersion string
		wantOptions map[string]string
	}{
		{
			name: "no multiTarget uses the top-level values",
			input: kubeconfig.KubernetesDeploymentInput{
				KustomizeVersion: "5.3.0",
				KustomizeOptions: map[string]string{"flag": "val"},
			},
			wantVersion: "5.3.0",
			wantOptions: map[string]string{"flag": "val"},
		},
		{
			name: "multiTarget without overrides falls back to the top-level values",
			input: kubeconfig.KubernetesDeploymentInput{
				KustomizeVersion: "5.3.0",
				KustomizeOptions: map[string]string{"flag": "val"},
			},
			multiTarget: &kubeconfig.KubernetesMultiTarget{},
			wantVersion: "5.3.0",
			wantOptions: map[string]string{"flag": "val"},
		},
		{
			name: "multiTarget overrides dir, version and options",
			input: kubeconfig.KubernetesDeploymentInput{
				KustomizeVersion: "5.3.0",
				KustomizeOptions: map[string]string{"flag": "val"},
			},
			multiTarget: &kubeconfig.KubernetesMultiTarget{
				KustomizeDir:     "overlays/cluster1",
				KustomizeVersion: "5.4.3",
				KustomizeOptions: map[string]string{"load-restrictor": "LoadRestrictionsNone"},
			},
			wantDir:     "overlays/cluster1",
			wantVersion: "5.4.3",
			wantOptions: map[string]string{"load-restrictor": "LoadRestrictionsNone"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			input := &sdk.GetPlanPreviewInput[kubeconfig.KubernetesApplicationSpec]{
				Logger: zaptest.NewLogger(t),
			}
			ds := &sdk.DeploymentSource[kubeconfig.KubernetesApplicationSpec]{}
			spec := &kubeconfig.KubernetesApplicationSpec{Input: tc.input}

			got := buildLoaderInput(input, ds, spec, tc.multiTarget)
			assert.Equal(t, tc.wantDir, got.KustomizeDir)
			assert.Equal(t, tc.wantVersion, got.KustomizeVersion)
			assert.Equal(t, tc.wantOptions, got.KustomizeOptions)
		})
	}
}
