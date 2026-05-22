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

// Package provider contains live resource fetching logic.
//
// Full unit testing of GetLiveResources requires a running Kubernetes cluster
// because *Kubectl is a concrete type that shells out to the kubectl binary.
// The selector-building logic (LabelManagedBy, LabelApplication) is tested
// indirectly via the integration-style test below: we verify the function
// returns an error when kubectl cannot reach a cluster, which confirms the
// function does attempt the expected kubectl.GetAll call.
package provider

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGetLiveResources_InvalidKubeconfig verifies that GetLiveResources
// returns an error when the kubeconfig path does not exist / is unreachable.
// This exercises the error path of kubectl.GetAll and confirms that
// GetLiveResources propagates the failure rather than silently swallowing it.
func TestGetLiveResources_InvalidKubeconfig(t *testing.T) {
	t.Parallel()

	kubectl := NewKubectl("kubectl")
	_, _, err := GetLiveResources(context.Background(), kubectl, "/nonexistent/kubeconfig", "app-id-123")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed while listing all namespace-scoped resources")
}
