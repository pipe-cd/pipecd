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

package kubernetes

import (
	"context"
	"fmt"
	"time"

	provider "github.com/kapetaniosci/pipe/pkg/app/piped/cloudprovider/kubernetes"
	"github.com/kapetaniosci/pipe/pkg/app/piped/planner"
	"github.com/kapetaniosci/pipe/pkg/model"
)

type Planner struct {
}

type registerer interface {
	Register(k model.ApplicationKind, p planner.Planner) error
}

// Register registers this planner into the given registerer.
func Register(r registerer) {
	r.Register(model.ApplicationKind_KUBERNETES, &Planner{})
}

func (p *Planner) Plan(ctx context.Context, in planner.Input) (out planner.Output, err error) {
	cfg := in.DeploymentConfig.KubernetesDeploymentSpec
	if cfg == nil {
		err = fmt.Errorf("malfored deployment configuration: missing KubernetesDeploymentSpec")
		return
	}

	// This is the first time to deployment
	// or it was unabled to retrieve that value
	// We uses the specified pipeline or the default one.
	if in.LastSuccessfulCommitHash == "" {
		out.Stages = buildPipeline(cfg.Pipeline, time.Now())
		out.Description = fmt.Sprintf("Apply all manifests at commit %s", in.Deployment.CommitHash())
		return
	}

	pv := provider.NewProvider(in.RepoDir, in.AppDir, cfg.Input, in.Logger)
	if err = pv.Init(ctx); err != nil {
		return
	}

	return
}
