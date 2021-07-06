// Copyright 2021 The PipeCD Authors.
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

	"go.uber.org/zap"

	provider "github.com/pipe-cd/pipe/pkg/app/piped/cloudprovider/kubernetes"
	"github.com/pipe-cd/pipe/pkg/app/piped/deploysource"
	"github.com/pipe-cd/pipe/pkg/cache"
	"github.com/pipe-cd/pipe/pkg/config"
	"github.com/pipe-cd/pipe/pkg/diff"
	"github.com/pipe-cd/pipe/pkg/model"
)

func (b *builder) kubernetesDiff(
	ctx context.Context,
	app *model.Application,
	cmd model.Command_BuildPlanPreview,
	lastSuccessfulCommit string,
	buf *bytes.Buffer,
) (string, error) {

	var oldManifests, newManifests []provider.Manifest
	var err error

	repoCfg := config.PipedRepository{
		RepoID: b.repoCfg.RepoID,
		Remote: b.repoCfg.Remote,
		Branch: cmd.HeadBranch,
	}

	targetDSP := deploysource.NewProvider(
		b.workingDir,
		repoCfg,
		"target",
		cmd.HeadCommit,
		b.gitClient,
		app.GitPath,
		b.secretDecrypter,
	)
	newManifests, err = loadKubernetesManifests(ctx, *app, cmd.HeadCommit, targetDSP, b.appManifestsCache, b.logger)
	if err != nil {
		fmt.Fprintf(buf, "failed to load kubernetes manifests at the head commit (%v)\n", err)
		return "", err
	}

	if lastSuccessfulCommit != "" {
		runningDSP := deploysource.NewProvider(
			b.workingDir,
			repoCfg,
			"running",
			lastSuccessfulCommit,
			b.gitClient,
			app.GitPath,
			b.secretDecrypter,
		)
		oldManifests, err = loadKubernetesManifests(ctx, *app, lastSuccessfulCommit, runningDSP, b.appManifestsCache, b.logger)
		if err != nil {
			fmt.Fprintf(buf, "failed to load kubernetes manifests at the running commit (%v)\n", err)
			return "", err
		}
	}

	result, err := provider.DiffList(oldManifests, newManifests,
		diff.WithEquateEmpty(),
		diff.WithIgnoreAddingMapKeys(),
		diff.WithCompareNumberAndNumericString(),
	)
	if err != nil {
		fmt.Fprintf(buf, "failed to compare manifests (%v)\n", err)
		return "", err
	}

	if result.NoChange() {
		fmt.Fprintln(buf, "No changes were detected")
		return "No changes were detected", nil
	}

	summary := fmt.Sprintf("%d added manifests, %d changed manifests, %d deleted manifests", len(result.Adds), len(result.Changes), len(result.Deletes))
	fmt.Fprintf(buf, "--- Last Deploy\n+++ Head Commit\n\n%s\n", result.DiffString())

	return summary, nil
}

func loadKubernetesManifests(ctx context.Context, app model.Application, commit string, dsp deploysource.Provider, manifestsCache cache.Cache, logger *zap.Logger) (manifests []provider.Manifest, err error) {
	cache := provider.AppManifestsCache{
		AppID:  app.Id,
		Cache:  manifestsCache,
		Logger: logger,
	}
	manifests, ok := cache.Get(commit)
	if ok {
		return manifests, nil
	}

	// When the manifests were not in the cache we have to load them.
	ds, err := dsp.Get(ctx, io.Discard)
	if err != nil {
		return nil, err
	}

	deployCfg := ds.DeploymentConfig.KubernetesDeploymentSpec
	if deployCfg == nil {
		return nil, fmt.Errorf("malformed deployment configuration file")
	}

	loader := provider.NewManifestLoader(
		app.Name,
		ds.AppDir,
		ds.RepoDir,
		app.GitPath.ConfigFilename,
		deployCfg.Input,
		logger,
	)
	manifests, err = loader.LoadManifests(ctx)
	if err != nil {
		return nil, err
	}

	cache.Put(commit, manifests)
	return manifests, nil
}
