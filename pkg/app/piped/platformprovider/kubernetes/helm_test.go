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

package kubernetes

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/app/piped/toolregistry"
)

func TestTemplateLocalChart(t *testing.T) {
	t.Parallel()

	var (
		ctx       = context.Background()
		appName   = "testapp"
		appDir    = "testdata"
		chartPath = "testchart"
	)

	// TODO: Preinstall a helm version inside CI runner to avoid installing.
	helmPath, _, err := toolregistry.DefaultRegistry().Helm(ctx, "")
	require.NoError(t, err)

	helm := NewHelm("", helmPath, zap.NewNop())
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

	// TODO: Preinstall a helm version inside CI runner to avoid installing.
	helmPath, _, err := toolregistry.DefaultRegistry().Helm(ctx, "")
	require.NoError(t, err)

	helm := NewHelm("", helmPath, zap.NewNop())
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
			valueFilePath: "http://example.com/values.yaml",
			wantErr:       false,
		},
		{
			name:          "Values file locates allowed remote URL (https)",
			appDir:        "testdata/testhelm/appconfdir",
			valueFilePath: "https://example.com/values.yaml",
			wantErr:       false,
		},
		{
			name:          "Values file locates disallowed remote URL (ftp)",
			appDir:        "testdata/testhelm/appconfdir",
			valueFilePath: "ftp://example.com/values.yaml",
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
		{
			name:          "Values file in sibling directory with matching prefix",
			appDir:        "testdata/testhelm/appconfdir",
			valueFilePath: "../appconfdir-sibling/values.yaml",
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

func TestVerifyHelmValueFilePath_TraversalAndSymlinks(t *testing.T) {
	t.Parallel()

	rootDir := t.TempDir()
	appDir := filepath.Join(rootDir, "app")
	subDir := filepath.Join(appDir, "sub")
	siblingDir := filepath.Join(rootDir, "app-secrets")
	outsideDir := filepath.Join(rootDir, "outside")

	require.NoError(t, os.MkdirAll(subDir, 0o755))
	require.NoError(t, os.MkdirAll(siblingDir, 0o755))
	require.NoError(t, os.MkdirAll(outsideDir, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(appDir, "values.yaml"), []byte("ok: true\n"), 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(siblingDir, "values.yaml"), []byte("secret: true\n"), 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(outsideDir, "secret.yaml"), []byte("secret: true\n"), 0o600))

	t.Run("absolute path with parent traversal escaping appDir", func(t *testing.T) {
		t.Parallel()
		uncleaned := appDir + string(filepath.Separator) + ".." + string(filepath.Separator) + "outside" + string(filepath.Separator) + "secret.yaml"
		require.Error(t, verifyHelmValueFilePath(appDir, uncleaned))
	})

	t.Run("sibling directory sharing appDir prefix via relative path", func(t *testing.T) {
		t.Parallel()
		require.Error(t, verifyHelmValueFilePath(appDir, filepath.Join("..", "app-secrets", "values.yaml")))
	})

	t.Run("sibling directory sharing appDir prefix via absolute path", func(t *testing.T) {
		t.Parallel()
		require.Error(t, verifyHelmValueFilePath(appDir, filepath.Join(siblingDir, "values.yaml")))
	})

	// Test symlink scenarios when the OS/filesystem allows creating symlinks.
	validSubSymlink := filepath.Join(subDir, "valid-link.yaml")
	if err := os.Symlink(filepath.Join("..", "values.yaml"), validSubSymlink); err != nil {
		t.Skipf("symlinks unavailable in this environment: %v", err)
	}
	escapeSubSymlink := filepath.Join(subDir, "escape-link.yaml")
	require.NoError(t, os.Symlink(filepath.Join("..", "..", "outside", "secret.yaml"), escapeSubSymlink))

	chainedTarget := filepath.Join(appDir, "chain-mid.yaml")
	chainedStart := filepath.Join(appDir, "chain-start.yaml")
	require.NoError(t, os.Symlink(filepath.Join("..", "outside", "secret.yaml"), chainedTarget))
	require.NoError(t, os.Symlink("chain-mid.yaml", chainedStart))

	symlinkedSubDir := filepath.Join(appDir, "linked-dir")
	require.NoError(t, os.Symlink(outsideDir, symlinkedSubDir))

	danglingEscape := filepath.Join(subDir, "dangling-escape.yaml")
	require.NoError(t, os.Symlink(filepath.Join("..", "..", "outside", "missing.yaml"), danglingEscape))

	t.Run("subdirectory relative symlink inside appDir", func(t *testing.T) {
		t.Parallel()
		require.NoError(t, verifyHelmValueFilePath(appDir, filepath.Join("sub", "valid-link.yaml")))
	})

	t.Run("subdirectory relative symlink escaping appDir", func(t *testing.T) {
		t.Parallel()
		require.Error(t, verifyHelmValueFilePath(appDir, filepath.Join("sub", "escape-link.yaml")))
	})

	t.Run("chained symlink escaping appDir", func(t *testing.T) {
		t.Parallel()
		require.Error(t, verifyHelmValueFilePath(appDir, "chain-start.yaml"))
	})

	t.Run("symlinked subdirectory escaping appDir", func(t *testing.T) {
		t.Parallel()
		require.Error(t, verifyHelmValueFilePath(appDir, filepath.Join("linked-dir", "secret.yaml")))
	})

	t.Run("dangling symlink in subdirectory escaping appDir", func(t *testing.T) {
		t.Parallel()
		require.Error(t, verifyHelmValueFilePath(appDir, filepath.Join("sub", "dangling-escape.yaml")))
	})
}
