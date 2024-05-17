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

	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/pipe-cd/pipecd/pkg/app/piped/deploysource"
	provider "github.com/pipe-cd/pipecd/pkg/app/piped/platformprovider/kubernetes"
	"github.com/pipe-cd/pipecd/pkg/cache"
	"github.com/pipe-cd/pipecd/pkg/diff"
	"github.com/pipe-cd/pipecd/pkg/model"
)

func (b *builder) kubernetesDiff(
	ctx context.Context,
	app *model.Application,
	targetDSP deploysource.Provider,
	lastSuccessfulCommit string,
	buf *bytes.Buffer,
) (*diffResult, error) {

	var oldManifests, newManifests []provider.Manifest
	var err error

	// Use discovery to discover APIs supported by the Kubernetes API server.
	// This should be run periodically with a low rate because the APIs are not added frequently.
	// https://godoc.org/k8s.io/client-go/discovery
	cp, ok := b.pipedCfg.FindPlatformProvider(app.PlatformProvider, model.ApplicationKind_KUBERNETES)
	if !ok {
		err = fmt.Errorf("provider %s was not found", app.PlatformProvider)
		return nil, err
	}
	kubeConfig, err := clientcmd.BuildConfigFromFlags(cp.KubernetesConfig.MasterURL, cp.KubernetesConfig.KubeConfigPath)
	if err != nil {
		err = fmt.Errorf("failed to build kube config: %w", err)
		return nil, err
	}
	discoveryClient, err := discovery.NewDiscoveryClientForConfig(kubeConfig)
	if err != nil {
		err = fmt.Errorf("failed to create discovery client: %w", err)
		return nil, err
	}
	groupResources, err := discoveryClient.ServerPreferredResources()
	if err != nil {
		err = fmt.Errorf("failed to fetch preferred resources: %w", err)
		return nil, err
	}

	isNamespacedResources := make(map[schema.GroupVersionKind]bool)
	for _, gr := range groupResources {
		for _, resource := range gr.APIResources {
			gvk := schema.FromAPIVersionAndKind(gr.GroupVersion, resource.Kind)
			isNamespacedResources[gvk] = resource.Namespaced
		}
	}

	newManifests, err = loadKubernetesManifests(ctx, *app, targetDSP, b.appManifestsCache, isNamespacedResources, b.gitClient, b.logger)
	if err != nil {
		fmt.Fprintf(buf, "failed to load kubernetes manifests at the head commit (%v)\n", err)
		return nil, err
	}

	if lastSuccessfulCommit != "" {
		runningDSP := deploysource.NewProvider(
			b.workingDir,
			deploysource.NewGitSourceCloner(b.gitClient, b.repoCfg, "running", lastSuccessfulCommit),
			*app.GitPath,
			b.secretDecrypter,
		)
		oldManifests, err = loadKubernetesManifests(ctx, *app, runningDSP, b.appManifestsCache, isNamespacedResources, b.gitClient, b.logger)
		if err != nil {
			fmt.Fprintf(buf, "failed to load kubernetes manifests at the running commit (%v)\n", err)
			return nil, err
		}
	}

	result, err := provider.DiffList(
		oldManifests,
		newManifests,
		b.logger,
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

	summary := fmt.Sprintf("%d added manifests, %d changed manifests, %d deleted manifests", len(result.Adds), len(result.Changes), len(result.Deletes))
	details := result.Render(provider.DiffRenderOptions{
		MaskSecret:     true,
		UseDiffCommand: true,
	})
	fmt.Fprintf(buf, "--- Last Deploy\n+++ Head Commit\n\n%s\n", details)

	return &diffResult{
		summary: summary,
	}, nil
}

func loadKubernetesManifests(ctx context.Context, app model.Application, dsp deploysource.Provider, manifestsCache cache.Cache, isNamespacedResources map[schema.GroupVersionKind]bool, gc gitClient, logger *zap.Logger) (manifests []provider.Manifest, err error) {
	commit := dsp.Revision()
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

	appCfg := ds.ApplicationConfig.KubernetesApplicationSpec
	if appCfg == nil {
		return nil, fmt.Errorf("malformed application configuration file")
	}

	loader := provider.NewLoader(
		app.Name,
		ds.AppDir,
		ds.RepoDir,
		app.GitPath.ConfigFilename,
		appCfg.Input,
		isNamespacedResources,
		gc,
		logger,
	)
	manifests, err = loader.LoadManifests(ctx)
	if err != nil {
		return nil, err
	}

	cache.Put(commit, manifests)
	return manifests, nil
}
