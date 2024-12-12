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
	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/toolregistry"
	"github.com/pipe-cd/pipecd/pkg/plugin/toolregistry/toolregistrytest"
)

func TestKustomizeTemplate(t *testing.T) {
	t.Parallel()

	var (
		ctx     = context.TODO()
		appName = "testapp"
		appDir  = "testdata/testkustomize"
	)

	c, err := toolregistrytest.NewToolRegistry(t)
	require.NoError(t, err)

	r := toolregistry.NewRegistry(c)

	t.Cleanup(func() { c.Close() })

	kustomizePath, err := r.Kustomize(context.Background(), "5.4.3")
	require.NoError(t, err)
	require.NotEmpty(t, kustomizePath)

	kustomize := NewKustomize(kustomizePath, zap.NewNop())
	out, err := kustomize.Template(ctx, appName, appDir, map[string]string{
		"load-restrictor": "LoadRestrictionsNone",
	})
	require.NoError(t, err)
	assert.True(t, len(out) > 0)
}
