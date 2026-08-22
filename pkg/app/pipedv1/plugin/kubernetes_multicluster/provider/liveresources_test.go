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

package provider

import (
	"context"
	"testing"

	"github.com/pipe-cd/piped-plugin-sdk-go/toolregistry/toolregistrytest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes_multicluster/toolregistry"
)

// TestGetLiveResources_InvalidKubeconfig verifies that GetLiveResources
// returns an error when the kubeconfig path does not exist / is unreachable.
func TestGetLiveResources_InvalidKubeconfig(t *testing.T) {
	t.Parallel()

	registry := toolregistry.NewRegistry(toolregistrytest.NewTestToolRegistry(t))
	kubectlPath, err := registry.Kubectl(context.Background(), "")
	require.NoError(t, err)

	kubectl := NewKubectl(kubectlPath)
	_, _, err = GetLiveResources(context.Background(), kubectl, "/nonexistent/kubeconfig", "app-id-123")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed while listing all namespace-scoped resources")
}
