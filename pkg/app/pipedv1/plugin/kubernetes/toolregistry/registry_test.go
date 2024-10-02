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

package toolregistry

import (
	"context"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/toolregistry/toolregistrytest"
)

func TestRegistry_Kubectl(t *testing.T) {
	t.Parallel()

	c, err := toolregistrytest.NewToolRegistry(t)
	require.NoError(t, err)

	r := NewRegistry(c)

	t.Cleanup(func() { c.Close() })

	p, err := r.Kubectl(context.Background(), "1.30.2")
	require.NoError(t, err)
	require.NotEmpty(t, p)

	out, err := exec.CommandContext(context.Background(), p, "version", "--client=true").CombinedOutput()
	require.NoError(t, err)

	expected := "Client Version: v1.30.2\nKustomize Version: v5.0.4-0.20230601165947-6ce0bf390ce3"

	assert.Equal(t, strings.TrimSpace(expected), strings.TrimSpace(string(out)))
}

func TestRegistry_Kustomize(t *testing.T) {
	t.Parallel()

	c, err := toolregistrytest.NewToolRegistry(t)
	require.NoError(t, err)

	r := NewRegistry(c)

	t.Cleanup(func() { c.Close() })

	p, err := r.Kustomize(context.Background(), "5.4.3")
	require.NoError(t, err)
	require.NotEmpty(t, p)

	out, err := exec.CommandContext(context.Background(), p, "version").CombinedOutput()
	require.NoError(t, err)

	expected := "v5.4.3"

	assert.Equal(t, strings.TrimSpace(expected), strings.TrimSpace(string(out)))
}

func TestRegistry_Helm(t *testing.T) {
	t.Parallel()

	c, err := toolregistrytest.NewToolRegistry(t)
	require.NoError(t, err)

	r := NewRegistry(c)

	t.Cleanup(func() { c.Close() })

	p, err := r.Helm(context.Background(), "3.16.1")
	require.NoError(t, err)
	require.NotEmpty(t, p)

	out, err := exec.CommandContext(context.Background(), p, "version").CombinedOutput()
	require.NoError(t, err)

	expected := `version.BuildInfo{Version:"v3.16.1", GitCommit:"5a5449dc42be07001fd5771d56429132984ab3ab", GitTreeState:"clean", GoVersion:"go1.22.7"}`

	assert.Equal(t, strings.TrimSpace(expected), strings.TrimSpace(string(out)))
}
