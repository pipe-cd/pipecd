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

package cloudrun

import (
	"context"
	"fmt"
	"io/ioutil"
	"time"

	"go.uber.org/zap"

	provider "github.com/pipe-cd/pipe/pkg/app/piped/cloudprovider/cloudrun"
	"github.com/pipe-cd/pipe/pkg/app/piped/planner"
	"github.com/pipe-cd/pipe/pkg/model"
)

// Planner plans the deployment pipeline for CloudRun application.
type Planner struct {
}

type registerer interface {
	Register(k model.ApplicationKind, p planner.Planner) error
}

// Register registers this planner into the given registerer.
func Register(r registerer) {
	r.Register(model.ApplicationKind_CLOUDRUN, &Planner{})
}

// Plan decides which pipeline should be used for the given input.
func (p *Planner) Plan(ctx context.Context, in planner.Input) (out planner.Output, err error) {
	ds, err := in.TargetDSP.Get(ctx, ioutil.Discard)
	if err != nil {
		err = fmt.Errorf("error while preparing deploy source data (%v)", err)
		return
	}

	cfg := ds.DeploymentConfig.CloudRunDeploymentSpec
	if cfg == nil {
		err = fmt.Errorf("missing CloudRunDeploymentSpec in deployment configuration")
		return
	}

	// Determine application version from the manifest.
	if version, e := p.determineVersion(ds.AppDir, cfg.Input.ServiceManifestFile); e == nil {
		out.Version = version
	} else {
		out.Version = "unknown"
		in.Logger.Warn("unable to determine target version", zap.Error(e))
	}

	// This is the first time to deploy this application or it was unable to retrieve that value.
	// We just do the quick sync.
	if in.MostRecentSuccessfulCommitHash == "" {
		out.Stages = buildQuickSyncPipeline(cfg.Input.AutoRollback, time.Now())
		out.Summary = "Quick sync by deploying the new version and configuring all traffic to it because it seems this is the first deployment"
		return
	}

	// When no pipeline was configured, do the quick sync.
	if cfg.Pipeline == nil || len(cfg.Pipeline.Stages) == 0 {
		out.Stages = buildQuickSyncPipeline(cfg.Input.AutoRollback, time.Now())
		out.Summary = "Quick sync by deploying the new version and configuring all traffic to it because no pipeline was configured"
		return
	}

	// Load service manifest at the last deployed commit to decide running version.
	ds, err = in.RunningDSP.Get(ctx, ioutil.Discard)
	if err == nil {
		if lastVersion, e := p.determineVersion(ds.AppDir, cfg.Input.ServiceManifestFile); e == nil {
			out.Stages = buildProgressivePipeline(cfg.Pipeline, cfg.Input.AutoRollback, time.Now())
			out.Summary = fmt.Sprintf("Sync progressively because of updating image from %s to %s", lastVersion, out.Version)
			return
		}
	}

	out.Stages = buildProgressivePipeline(cfg.Pipeline, cfg.Input.AutoRollback, time.Now())
	out.Summary = "Sync progressively with the specified pipeline"
	return
}

func (p *Planner) determineVersion(appDir, serviceManifestFile string) (string, error) {
	sm, err := provider.LoadServiceManifest(appDir, serviceManifestFile)
	if err != nil {
		return "", err
	}

	return provider.FindImageTag(sm)
}
