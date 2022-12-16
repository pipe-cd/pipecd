package planpreview

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/pipe-cd/pipecd/pkg/app/piped/deploysource"
	provider "github.com/pipe-cd/pipecd/pkg/app/piped/platformprovider/cloudrun"
	"github.com/pipe-cd/pipecd/pkg/diff"
	"github.com/pipe-cd/pipecd/pkg/model"
)

func (b *builder) cloudrundiff(
	ctx context.Context,
	app *model.Application,
	targetDSP deploysource.Provider,
	lastCommit string,
	buf *bytes.Buffer,
) (*diffResult, error) {
	var (
		oldManifest, newManifest provider.ServiceManifest
		err                      error
	)

	newManifest, err = b.loadCloudRunManifests(ctx, *app, targetDSP)
	if err != nil {
		fmt.Fprintf(buf, "failed to load cloud run manifest at the head commit (%v)\n", err)
		return nil, err
	}

	if lastCommit != "" {
		runningDSP := deploysource.NewProvider(
			b.workingDir,
			deploysource.NewGitSourceCloner(b.gitClient, b.repoCfg, "running", lastCommit),
			*app.GitPath,
			b.secretDecrypter,
		)
		oldManifest, err = b.loadCloudRunManifests(ctx, *app, runningDSP)
		if err != nil {
			fmt.Fprintf(buf, "failed to load cloud run manifest at the running commit (%v)\n", err)
			return nil, err
		}
	}

	result, err := provider.Diff(
		oldManifest,
		newManifest,
		diff.WithEquateEmpty(),
		diff.WithCompareNumberAndNumericString(),
	)
	if err != nil {
		fmt.Fprintf(buf, "failed to compare manifests (%v)\n", err)
		return nil, err
	}

	summary := fmt.Sprintf("%d changes were detected", len(result.Diff.Nodes()))
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
		summary: summary,
	}, nil

}

func (b *builder) loadCloudRunManifests(ctx context.Context, app model.Application, dsp deploysource.Provider) (manifest provider.ServiceManifest, err error) {
	commit := dsp.Revision()
	cache := provider.ServiceManifestCache{
		AppID:  app.Id,
		Cache:  b.appManifestsCache,
		Logger: b.logger,
	}

	manifest, ok := cache.Get(commit)
	if ok {
		return
	}

	ds, err := dsp.Get(ctx, io.Discard)
	if err != nil {
		return
	}

	appCfg := ds.ApplicationConfig.CloudRunApplicationSpec
	if appCfg == nil {
		err = fmt.Errorf("malformed application configuration file")
		return
	}

	manifest, err = provider.LoadServiceManifest(ds.AppDir, appCfg.Input.ServiceManifestFile)
	if err != nil {
		return
	}

	cache.Put(commit, manifest)
	return
}
