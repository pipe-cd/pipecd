// Copyright 2024 The PipeCD Authors.
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

package planpreview

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/pipe-cd/pipecd/pkg/app/piped/deploysource"
	provider "github.com/pipe-cd/pipecd/pkg/app/piped/platformprovider/ecs"
	"github.com/pipe-cd/pipecd/pkg/diff"
	"github.com/pipe-cd/pipecd/pkg/model"
)

func (b *builder) ecsdiff(
	ctx context.Context,
	app *model.Application,
	targetDSP deploysource.Provider,
	lastCommit string,
	buf *bytes.Buffer,
) (*diffResult, error) {
	var (
		oldManifests, newManifests provider.ECSManifests
		err                        error
	)

	newManifests, err = b.loadECSManifests(ctx, *app, targetDSP)
	if err != nil {
		fmt.Fprintf(buf, "failed to load ecs manifests at the head commit (%v)\n", err)
		return nil, err
	}

	if lastCommit == "" {
		fmt.Fprintf(buf, "failed to find the commit of the last successful deployment")
		return nil, fmt.Errorf("cannot get the old manifests without the last successful deployment")
	}

	runningDSP := deploysource.NewProvider(
		b.workingDir,
		deploysource.NewGitSourceCloner(b.gitClient, b.repoCfg, "running", lastCommit),
		*app.GitPath,
		b.secretDecrypter,
	)

	oldManifests, err = b.loadECSManifests(ctx, *app, runningDSP)
	if err != nil {
		fmt.Fprintf(buf, "failed to load ecs manifests at the running commit (%v)\n", err)
		return nil, err
	}

	result, err := provider.Diff(
		oldManifests,
		newManifests,
		diff.WithEquateEmpty(),
		diff.WithCompareNumberAndNumericString(),
	)
	if err != nil {
		fmt.Fprintf(buf, "failed to compare manifests (%v)\n", err)
		return nil, err
	}

	if result.NoChange() {
		fmt.Fprintln(buf, "No changes were detected")
		return &diffResult{
			summary:  "No changes were detected",
			noChange: true,
		}, nil
	}

	details := result.Render(provider.DiffRenderOptions{
		UseDiffCommand: true,
	})
	fmt.Fprintf(buf, "--- Last Deploy\n+++ Head Commit\n\n%s\n", details)

	return &diffResult{
		summary: fmt.Sprintf("%d changes were detected", len(result.Diff.Nodes())),
	}, nil
}

func (b *builder) loadECSManifests(ctx context.Context, app model.Application, dsp deploysource.Provider) (provider.ECSManifests, error) {
	commit := dsp.Revision()
	cache := provider.ECSManifestsCache{
		AppID:  app.Id,
		Cache:  b.appManifestsCache,
		Logger: b.logger,
	}

	manifests, ok := cache.Get(commit)
	if ok {
		return manifests, nil
	}

	ds, err := dsp.Get(ctx, io.Discard)
	if err != nil {
		return provider.ECSManifests{}, err
	}

	appCfg := ds.ApplicationConfig.ECSApplicationSpec
	if appCfg == nil {
		return provider.ECSManifests{}, fmt.Errorf("malformed application configuration file")
	}

	taskDef, err := provider.LoadTaskDefinition(ds.AppDir, appCfg.Input.TaskDefinitionFile)
	if err != nil {
		return provider.ECSManifests{}, err
	}
	serviceDef, err := provider.LoadServiceDefinition(ds.AppDir, appCfg.Input.ServiceDefinitionFile)
	if err != nil {
		return provider.ECSManifests{}, err
	}

	manifests = provider.ECSManifests{
		TaskDefinition:    &taskDef,
		ServiceDefinition: &serviceDef,
	}

	cache.Put(commit, manifests)
	return manifests, nil
}
