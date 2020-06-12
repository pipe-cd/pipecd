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
	cfg := in.DeploymentConfig.TerraformDeploymentSpec
	if cfg == nil {
		err = fmt.Errorf("malfored deployment configuration: missing TerraformDeploymentSpec")
		return
	}

	if cfg.Pipeline != nil && len(cfg.Pipeline.Stages) > 0 {
		out.Stages = buildProgressivePipeline(cfg.Pipeline.Stages, time.Now())
		out.Description = "Deploy terraform with the specified progressive pipeline."
		return
	}

	out.Stages = builDefaultPipeline(time.Now())
	out.Description = "Deploy terraform with the default pipeline."
	return
}
