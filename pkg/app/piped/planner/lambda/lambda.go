// Copyright 2023 The PipeCD Authors.
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

package lambda

import (
	"context"
	"fmt"
	"io"
	"time"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/app/piped/planner"
	provider "github.com/pipe-cd/pipecd/pkg/app/piped/platformprovider/lambda"
	"github.com/pipe-cd/pipecd/pkg/model"
)

// Planner plans the deployment pipeline for Lambda application.
type Planner struct {
}

type registerer interface {
	Register(k model.ApplicationKind, p planner.Planner) error
}

// Register registers this planner into the given registerer.
func Register(r registerer) {
	r.Register(model.ApplicationKind_LAMBDA, &Planner{})
}

// Plan decides which pipeline should be used for the given input.
func (p *Planner) Plan(ctx context.Context, in planner.Input) (out planner.Output, err error) {
	ds, err := in.TargetDSP.Get(ctx, io.Discard)
	if err != nil {
		err = fmt.Errorf("error while preparing deploy source data (%v)", err)
		return
	}

	cfg := ds.ApplicationConfig.LambdaApplicationSpec
	if cfg == nil {
		err = fmt.Errorf("missing LambdaApplicationSpec in application configuration")
		return
	}

	// Determine application version from the manifest
	if version, e := determineVersion(ds.AppDir, cfg.Input.FunctionManifestFile); e != nil {
		out.Version = "unknown"
		in.Logger.Warn("unable to determine target version", zap.Error(e))
	} else {
		out.Version = version
	}

	if versions, e := determineVersions(ds.AppDir, cfg.Input.FunctionManifestFile); e != nil || len(versions) == 0 {
		in.Logger.Warn("unable to determine target versions", zap.Error(e))
		out.Versions = []*model.ArtifactVersion{
			{
				Kind:    model.ArtifactVersion_UNKNOWN,
				Version: "unknown",
			},
		}
	} else {
		out.Versions = versions
	}

	autoRollback := *cfg.Input.AutoRollback

	// In case the strategy has been decided by trigger.
	// For example: user triggered the deployment via web console.
	switch in.Trigger.SyncStrategy {
	case model.SyncStrategy_QUICK_SYNC:
		out.SyncStrategy = model.SyncStrategy_QUICK_SYNC
		out.Stages = buildQuickSyncPipeline(autoRollback, time.Now())
		out.Summary = in.Trigger.StrategySummary
		return
	case model.SyncStrategy_PIPELINE:
		if cfg.Pipeline == nil {
			err = fmt.Errorf("unable to force sync with pipeline because no pipeline was specified")
			return
		}
		out.SyncStrategy = model.SyncStrategy_PIPELINE
		out.Stages = buildProgressivePipeline(cfg.Pipeline, autoRollback, time.Now())
		out.Summary = in.Trigger.StrategySummary
		return
	}

	// When no pipeline was configured, perform the quick sync.
	if cfg.Pipeline == nil || len(cfg.Pipeline.Stages) == 0 {
		out.SyncStrategy = model.SyncStrategy_QUICK_SYNC
		out.Stages = buildQuickSyncPipeline(autoRollback, time.Now())
		out.Summary = fmt.Sprintf("Quick sync to deploy version %s and configure all traffic to it (pipeline was not configured)", out.Version)
		return
	}

	// Force to use pipeline when the alwaysUsePipeline field was configured.
	if cfg.Planner.AlwaysUsePipeline {
		out.SyncStrategy = model.SyncStrategy_PIPELINE
		out.Stages = buildProgressivePipeline(cfg.Pipeline, autoRollback, time.Now())
		out.Summary = "Sync with the specified pipeline (alwaysUsePipeline was set)"
		return
	}

	// If this is the first time to deploy this application or it was unable to retrieve last successful commit,
	// we perform the quick sync strategy.
	if in.MostRecentSuccessfulCommitHash == "" {
		out.SyncStrategy = model.SyncStrategy_QUICK_SYNC
		out.Stages = buildQuickSyncPipeline(autoRollback, time.Now())
		out.Summary = fmt.Sprintf("Quick sync to deploy version %s and configure all traffic to it (it seems this is the first deployment)", out.Version)
		return
	}

	// Load service manifest at the last deployed commit to decide running version.
	ds, err = in.RunningDSP.Get(ctx, io.Discard)
	if err == nil {
		if lastVersion, e := determineVersion(ds.AppDir, cfg.Input.FunctionManifestFile); e == nil {
			out.SyncStrategy = model.SyncStrategy_PIPELINE
			out.Stages = buildProgressivePipeline(cfg.Pipeline, autoRollback, time.Now())
			out.Summary = fmt.Sprintf("Sync with pipeline to update version from %s to %s", lastVersion, out.Version)
			return
		}
	}

	out.SyncStrategy = model.SyncStrategy_PIPELINE
	out.Stages = buildProgressivePipeline(cfg.Pipeline, autoRollback, time.Now())
	out.Summary = "Sync with the specified pipeline"
	return
}

func determineVersion(appDir, functionManifestFile string) (string, error) {
	fm, err := provider.LoadFunctionManifest(appDir, functionManifestFile)
	if err != nil {
		return "", err
	}

	// Extract container image tag as application version.
	if fm.Spec.ImageURI != "" {
		return provider.FindImageTag(fm)
	}

	// Extract s3 object version as application version.
	if fm.Spec.S3ObjectVersion != "" {
		return fm.Spec.S3ObjectVersion, nil
	}

	// Extract source code commitish as application version.
	if fm.Spec.SourceCode.Ref != "" {
		return fm.Spec.SourceCode.Ref, nil
	}

	return "", fmt.Errorf("unable to determine version from manifest")
}

func determineVersions(appDir, functionManifestFile string) ([]*model.ArtifactVersion, error) {
	fm, err := provider.LoadFunctionManifest(appDir, functionManifestFile)
	if err != nil {
		return nil, err
	}

	return provider.FindArtifactVersions(fm)
}
