// Copyright 2020 The PipeCD Authors.
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

package terraform

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/pipe-cd/pipe/pkg/app/piped/planner"
	"github.com/pipe-cd/pipe/pkg/model"
)

// Planner plans the deployment pipeline for terraform application.
type Planner struct {
}

type registerer interface {
	Register(k model.ApplicationKind, p planner.Planner) error
}

// Register registers this planner into the given registerer.
func Register(r registerer) {
	r.Register(model.ApplicationKind_TERRAFORM, &Planner{})
}

// Plan decides which pipeline should be used for the given input.
func (p *Planner) Plan(ctx context.Context, in planner.Input) (out planner.Output, err error) {
	ds, err := in.TargetDSP.Get(ctx, io.Discard)
	if err != nil {
		err = fmt.Errorf("error while preparing deploy source data (%v)", err)
		return
	}

	cfg := ds.DeploymentConfig.TerraformDeploymentSpec
	if cfg == nil {
		err = fmt.Errorf("missing TerraformDeploymentSpec in deployment configuration")
		return
	}

	// If the deployment was triggered by forcing via web UI,
	// we rely on the user's decision.
	switch in.Trigger.SyncStrategy {
	case model.SyncStrategy_QUICK_SYNC:
		out.SyncStrategy = model.SyncStrategy_QUICK_SYNC
		out.Stages = buildQuickSyncPipeline(cfg.Input.AutoRollback, time.Now())
		out.Summary = "Quick sync by automatically applying any detected changes because no pipeline was configured (forced via web)"
		return
	case model.SyncStrategy_PIPELINE:
		if cfg.Pipeline == nil {
			err = fmt.Errorf("unable to force sync with pipeline because no pipeline was specified")
			return
		}
		out.SyncStrategy = model.SyncStrategy_PIPELINE
		out.Stages = buildProgressivePipeline(cfg.Pipeline, cfg.Input.AutoRollback, time.Now())
		out.Summary = "Sync with the specified progressive pipeline (forced via web)"
		return
	}

	now := time.Now()
	out.Version = "N/A"

	if cfg.Pipeline == nil || len(cfg.Pipeline.Stages) == 0 {
		out.SyncStrategy = model.SyncStrategy_QUICK_SYNC
		out.Stages = buildQuickSyncPipeline(cfg.Input.AutoRollback, now)
		out.Summary = "Quick sync by automatically applying any detected changes because no pipeline was configured"
		return
	}

	// Force to use pipeline when the alwaysUsePipeline field was configured.
	if cfg.Planner.AlwaysUsePipeline {
		out.SyncStrategy = model.SyncStrategy_PIPELINE
		out.Stages = buildProgressivePipeline(cfg.Pipeline, cfg.Input.AutoRollback, time.Now())
		out.Summary = "Sync with the specified pipeline (alwaysUsePipeline was set)"
		return
	}

	out.SyncStrategy = model.SyncStrategy_PIPELINE
	out.Stages = buildProgressivePipeline(cfg.Pipeline, cfg.Input.AutoRollback, now)
	out.Summary = "Sync with the specified progressive pipeline"
	return
}
