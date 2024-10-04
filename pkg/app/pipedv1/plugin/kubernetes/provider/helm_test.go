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

package provider

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/toolregistry"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/toolregistry/toolregistrytest"
)

func TestTemplateLocalChart(t *testing.T) {
	t.Parallel()

	var (
		ctx       = context.Background()
		appName   = "testapp"
		appDir    = "testdata"
		chartPath = "testchart"
	)

	r, err := toolregistrytest.NewToolRegistry(t)
	require.NoError(t, err)
	t.Cleanup(func() { r.Close() })

	registry := toolregistry.NewRegistry(r)
	helmPath, err := registry.Helm(ctx, "3.8.2")
	require.NoError(t, err)

	helm := NewHelm("", helmPath, zap.NewNop(), registry)
	out, err := helm.TemplateLocalChart(ctx, appName, appDir, "", chartPath, nil)
	require.NoError(t, err)

	out = strings.TrimPrefix(out, "---")
	manifests := strings.Split(out, "---")
	assert.Equal(t, 3, len(manifests))
}

func TestTemplateLocalChart_WithNamespace(t *testing.T) {
	t.Parallel()

	var (
		ctx       = context.Background()
		appName   = "testapp"
		appDir    = "testdata"
		chartPath = "testchart"
		namespace = "testnamespace"
	)

	r, err := toolregistrytest.NewToolRegistry(t)
	require.NoError(t, err)
	t.Cleanup(func() { r.Close() })

	registry := toolregistry.NewRegistry(r)
	helmPath, err := registry.Helm(ctx, "3.8.2")
	require.NoError(t, err)

	helm := NewHelm("", helmPath, zap.NewNop(), registry)
	out, err := helm.TemplateLocalChart(ctx, appName, appDir, namespace, chartPath, nil)
	require.NoError(t, err)

	out = strings.TrimPrefix(out, "---")

	manifests, _ := ParseManifests(out)
	for _, manifest := range manifests {
		metadata, err := manifest.GetNestedMap("metadata")
		require.NoError(t, err)
		require.Equal(t, namespace, metadata["namespace"])
	}
}

func TestVerifyHelmValueFilePath(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name          string
		appDir        string
		valueFilePath string
		wantErr       bool
	}{
		{
			name:          "Values file locates inside the app dir",
			appDir:        "testdata/testhelm/appconfdir",
			valueFilePath: "values.yaml",
			wantErr:       false,
		},
		{
			name:          "Values file locates inside the app dir (with ..)",
			appDir:        "testdata/testhelm/appconfdir",
			valueFilePath: "../../../testdata/testhelm/appconfdir/values.yaml",
			wantErr:       false,
		},
		{
			name:          "Values file locates under the app dir",
			appDir:        "testdata/testhelm/appconfdir",
			valueFilePath: "dir/values.yaml",
			wantErr:       false,
		},
		{
			name:          "Values file locates under the app dir (with ..)",
			appDir:        "testdata/testhelm/appconfdir",
			valueFilePath: "../../../testdata/testhelm/appconfdir/dir/values.yaml",
			wantErr:       false,
		},
		{
			name:          "arbitrary file locates outside the app dir",
			appDir:        "testdata/testhelm/appconfdir",
			valueFilePath: "/etc/hosts",
			wantErr:       true,
		},
		{
			name:          "arbitrary file locates outside the app dir (with ..)",
			appDir:        "testdata/testhelm/appconfdir",
			valueFilePath: "../../../../../../../../../../../../etc/hosts",
			wantErr:       true,
		},
		{
			name:          "Values file locates allowed remote URL (http)",
			appDir:        "testdata/testhelm/appconfdir",
			valueFilePath: "http://exmaple.com/values.yaml",
			wantErr:       false,
		},
		{
			name:          "Values file locates allowed remote URL (https)",
			appDir:        "testdata/testhelm/appconfdir",
			valueFilePath: "https://exmaple.com/values.yaml",
			wantErr:       false,
		},
		{
			name:          "Values file locates disallowed remote URL (ftp)",
			appDir:        "testdata/testhelm/appconfdir",
			valueFilePath: "ftp://exmaple.com/values.yaml",
			wantErr:       true,
		},
		{
			name:          "Values file is symlink targeting valid values file",
			appDir:        "testdata/testhelm/appconfdir",
			valueFilePath: "valid-symlink",
			wantErr:       false,
		},
		{
			name:          "Values file is symlink targeting invalid values file",
			appDir:        "testdata/testhelm/appconfdir",
			valueFilePath: "invalid-symlink",
			wantErr:       true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			err := verifyHelmValueFilePath(tc.appDir, tc.valueFilePath)
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
