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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/app/piped/toolregistry"
)

func TestKustomizeTemplate(t *testing.T) {
	t.Parallel()

	var (
		ctx     = context.TODO()
		appName = "testapp"
		appDir  = "testdata/testkustomize"
	)

	kustomizePath, _, err := toolregistry.DefaultRegistry().Kustomize(ctx, "")
	require.NoError(t, err)
	helmPath, _, err := toolregistry.DefaultRegistry().Helm(ctx, "")
	require.NoError(t, err)

	kustomize := NewKustomize("", kustomizePath, zap.NewNop())
	helm := NewHelm("", helmPath, zap.NewNop())
	out, err := kustomize.Template(ctx, appName, appDir, map[string]string{
		"load_restrictor": "LoadRestrictionsNone",
	}, helm)
	require.NoError(t, err)
	assert.True(t, len(out) > 0)
}

func TestKustomizeTemplate_WithHelm(t *testing.T) {
	t.Parallel()

	var (
		ctx     = context.TODO()
		appName = "testapp"
		appDir  = "testdata/testkustomize-with-helm"
	)

	kustomizePath, _, err := toolregistry.DefaultRegistry().Kustomize(ctx, "5.6.0")
	require.NoError(t, err)
	helmPath, _, err := toolregistry.DefaultRegistry().Helm(ctx, "3.17.0")
	require.NoError(t, err)

	kustomize := NewKustomize("5.6.0", kustomizePath, zap.NewNop())
	helm := NewHelm("3.17.0", helmPath, zap.NewNop())
	out, err := kustomize.Template(ctx, appName, appDir, map[string]string{
		"enable-helm": "",
	}, helm)
	require.NoError(t, err)
	assert.True(t, len(out) > 0)
}
