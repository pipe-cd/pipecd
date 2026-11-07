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
	"testing"
	"time"

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
	"github.com/pipe-cd/piped-plugin-sdk-go/logpersister/logpersistertest"
	"github.com/stretchr/testify/assert"

	kubeconfig "github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes_multicluster/config"
)

func TestRunOnTargets_failureDoesNotCancelOtherTargets(t *testing.T) {
	t.Parallel()

	targets := []stageTarget{
		{deployTarget: &sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]{Name: "failing"}},
		{deployTarget: &sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]{Name: "healthy"}},
	}

	lp := sdk.NewClient(nil, "kubernetes_multicluster", "app-id", "stage-id", logpersistertest.NewTestLogPersister(t), nil).LogPersister()
	status, dtStatuses := runOnTargets(t.Context(), lp, targets,
		func(ctx context.Context, dt *sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig], _ *kubeconfig.KubernetesMultiTarget) sdk.StageStatus {
			if dt.Name == "failing" {
				return sdk.StageStatusFailure
			}
			// The healthy target takes longer than the failing one, and only
			// succeeds if its context is still alive when it finishes.
			select {
			case <-ctx.Done():
				return sdk.StageStatusFailure
			case <-time.After(500 * time.Millisecond):
				return sdk.StageStatusSuccess
			}
		})

	assert.Equal(t, sdk.StageStatusFailure, status)
	assert.ElementsMatch(t, []sdk.DeployTargetStatus{
		{Name: "failing", Status: sdk.StageStatusFailure},
		{Name: "healthy", Status: sdk.StageStatusSuccess},
	}, dtStatuses)
}
