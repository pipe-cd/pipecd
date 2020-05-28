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

	in.MostRecentSuccessfulCommitHash = "626ad85b9c6c02c6409b9aa79ee433fb9b5507d7"

	// This is the first time to deploy this application
	// or it was unabled to retrieve that value
	// We just apply all manifests.
	if in.MostRecentSuccessfulCommitHash == "" {
		out.Stages = buildPipeline(time.Now())
		out.Description = fmt.Sprintf("Apply all manifests because no most recent successful commit")
		return
	}

	// If the commit is a revert one. Let's apply primary to rollback.
	// TODO: Determine if the new commit is a revert one.
	// out.Description = fmt.Sprintf("Rollback from %s", in.MostRecentSuccessfulCommitHash)

	// Load previous deployed manifests and new manifests to compare.
	pv := provider.NewProvider(in.RepoDir, in.AppDir, cfg.Input, in.Logger)
	if err = pv.Init(ctx); err != nil {
		return
	}

	// Load manifests of the new commit.
	newManifests, err := pv.LoadManifests(ctx)
	if err != nil {
		return
	}

	// Checkout to the most recent successful commit to load its manifests.
	err = in.Repo.Checkout(ctx, in.MostRecentSuccessfulCommitHash)
	if err != nil {
		return
	}

	// Load manifests of the previously applied commit.
	oldManifests, err := pv.LoadManifests(ctx)
	if err != nil {
		return
	}

	progressive, desc := decideStrategy(oldManifests, newManifests)
	out.Description = desc

	if progressive {
		out.Stages = buildProgressivePipeline(cfg.Pipeline, time.Now())
		return
	}

	out.Stages = buildPipeline(time.Now())
	return
}

func decideStrategy(olds, news []provider.Manifest) (progressive bool, desc string) {
	oldWorkload, ok := findWorkload(olds)
	if !ok {
		return false, "Apply all manifests because there is no workload running"
	}

	newWorkload, ok := findWorkload(news)
	if !ok {
		return false, "Apply all manifests because there is no workload in the new manifests"
	}

	diff := provider.Diff(oldWorkload, newWorkload)

	// If the workload' pod template or config/secret was touched
	// let's do the specified progressive pipeline.
	// out.Description = fmt.Sprintf("Do progressive deployment because image was changed to %s", "v1.0.0")

	// out.Description = fmt.Sprintf("Do progressive deployment because configmap %s was updated", "config")

	// out.Description = fmt.Sprintf("Do progressive deployment because pod template for workload %s was changed", "config")

	// // Otherwise, just apply the primary.

	// out.Description = fmt.Sprintf("Scale workload %s from %d to %d", "deployment-name", 1, 2)

	return
}

func findWorkload(manifests []provider.Manifest) (provider.Manifest, bool) {
	for _, m := range manifests {
		if m.Key.Kind != "Deployment" {
			continue
		}
		switch m.Key.APIVersion {
		case "v1", "apps/v1":
			return m, true
		default:
			continue
		}
	}
	return provider.Manifest{}, false
}

func findConfig(manifests []provider.Manifest) []provider.Manifest {
	configs := make([]provider.Manifest, 0)
	for _, m := range manifests {
		if m.Key.Kind != "ConfigMap" && m.Key.Kind != "Secret" {
			continue
		}
		switch m.Key.APIVersion {
		case "v1", "apps/v1":
			configs = append(configs, m)
		default:
			continue
		}
	}
	return configs
}
