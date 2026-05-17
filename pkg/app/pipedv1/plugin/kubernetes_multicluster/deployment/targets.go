// Copyright 2025 The PipeCD Authors.
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
	"fmt"

	"golang.org/x/sync/errgroup"

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"

	kubeconfig "github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes_multicluster/config"
)

// stageTarget pairs a deploy target with its optional per-target multiTarget config.
type stageTarget struct {
	deployTarget *sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]
	multiTarget  *kubeconfig.KubernetesMultiTarget
}

// buildStageTargets resolves which deploy targets a stage should run on.
// If multiTargets is empty, all dts are included. Otherwise only the named targets.
func buildStageTargets(
	lp sdk.StageLogPersister,
	dts []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig],
	multiTargets []kubeconfig.KubernetesMultiTarget,
) []stageTarget {
	dtMap := make(map[string]*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig], len(dts))
	for _, dt := range dts {
		dtMap[dt.Name] = dt
	}
	if len(multiTargets) == 0 {
		targets := make([]stageTarget, len(dts))
		for i, dt := range dts {
			targets[i] = stageTarget{deployTarget: dt}
		}
		return targets
	}
	targets := make([]stageTarget, 0, len(multiTargets))
	for _, mt := range multiTargets {
		dt, ok := dtMap[mt.Target.Name]
		if !ok {
			lp.Infof("Ignore multi target '%s': not matched any deployTarget", mt.Target.Name)
			continue
		}
		mt := mt // capture loop var
		targets = append(targets, stageTarget{deployTarget: dt, multiTarget: &mt})
	}
	return targets
}

// runOnTargets fans out fn across targets in parallel using errgroup.
// fn receives each target's deploy target and multiTarget (may be nil).
// Returns StageStatusSuccess if all succeed, StageStatusFailure if any fail.
func runOnTargets(
	ctx context.Context,
	lp sdk.StageLogPersister,
	targets []stageTarget,
	fn func(ctx context.Context, dt *sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig], mt *kubeconfig.KubernetesMultiTarget) sdk.StageStatus,
) sdk.StageStatus {
	eg, ctx := errgroup.WithContext(ctx)
	for _, tc := range targets {
		eg.Go(func() error {
			if status := fn(ctx, tc.deployTarget, tc.multiTarget); status == sdk.StageStatusFailure {
				return fmt.Errorf("stage failed for target %s", tc.deployTarget.Name)
			}
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		lp.Errorf("Stage failed on one or more targets: %v", err)
		return sdk.StageStatusFailure
	}
	return sdk.StageStatusSuccess
}
