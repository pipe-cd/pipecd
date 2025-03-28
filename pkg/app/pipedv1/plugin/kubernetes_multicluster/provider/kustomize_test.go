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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes_multicluster/toolregistry"
	"github.com/pipe-cd/pipecd/pkg/plugin/toolregistry/toolregistrytest"
)

func TestKustomizeTemplate(t *testing.T) {
	t.Parallel()

	var (
		ctx     = context.TODO()
		appName = "testapp"
		appDir  = "testdata/testkustomize"
	)

	c := toolregistrytest.NewTestToolRegistry(t)
	r := toolregistry.NewRegistry(c)

	kustomizePath, err := r.Kustomize(context.Background(), "5.4.3")
	require.NoError(t, err)
	require.NotEmpty(t, kustomizePath)

	kustomize := NewKustomize("5.4.3", kustomizePath, zaptest.NewLogger(t))
	out, err := kustomize.Template(ctx, appName, appDir, map[string]string{
		"load-restrictor": "LoadRestrictionsNone",
	}, nil)
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

	c := toolregistrytest.NewTestToolRegistry(t)
	r := toolregistry.NewRegistry(c)

	kustomizePath, err := r.Kustomize(context.Background(), "5.6.0")
	require.NoError(t, err)
	helmPath, err := r.Helm(context.Background(), "3.17.0")
	require.NoError(t, err)

	kustomize := NewKustomize("5.6.0", kustomizePath, zaptest.NewLogger(t))
	helm := NewHelm("3.17.0", helmPath, zaptest.NewLogger(t))
	out, err := kustomize.Template(ctx, appName, appDir, map[string]string{
		"enable-helm": "",
	}, helm)
	require.NoError(t, err)
	assert.True(t, len(out) > 0)
}

func TestKustomizeIsHelmCommandFlagAvailable(t *testing.T) {
	t.Parallel()

	kustomize := NewKustomize("4.0.1337", "", zaptest.NewLogger(t))
	assert.False(t, kustomize.isHelmCommandFlagAvailable())
	kustomize = NewKustomize("4.1.0", "", zaptest.NewLogger(t))
	assert.True(t, kustomize.isHelmCommandFlagAvailable())
	kustomize = NewKustomize("10.0.0", "", zaptest.NewLogger(t))
	assert.True(t, kustomize.isHelmCommandFlagAvailable())
}
