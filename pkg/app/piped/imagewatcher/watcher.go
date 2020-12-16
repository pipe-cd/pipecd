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

// Package imagewatcher provides a piped component
// that periodically checks the container registry and updates
// the image if there are differences with Git.
package imagewatcher

import (
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipe/pkg/app/piped/imageprovider"
	"github.com/pipe-cd/pipe/pkg/config"
	"github.com/pipe-cd/pipe/pkg/git"
	"github.com/pipe-cd/pipe/pkg/yamlprocessor"
)

const defaultCommitMessageFormat = "Update image %s to %s defined at %s in %s"

type Watcher interface {
	Run(context.Context) error
}

type gitClient interface {
	Clone(ctx context.Context, repoID, remote, branch, destination string) (git.Repo, error)
}

type watcher struct {
	config    *config.PipedSpec
	gitClient gitClient
	logger    *zap.Logger
	wg        sync.WaitGroup
	// For file locking.
	mu sync.Mutex

	// Indexed by repo id.
	gitRepos map[string]git.Repo
}

func NewWatcher(cfg *config.PipedSpec, gitClient gitClient, logger *zap.Logger) Watcher {
	return &watcher{
		config:    cfg,
		gitClient: gitClient,
		logger:    logger.Named("image-watcher"),
	}
}

// Run spawns goroutines for each image provider. They periodically pull the image
// from the container registry to compare the image with one in the git repository.
func (w *watcher) Run(ctx context.Context) error {
	// TODO: Spawn goroutines for each repository
	// Pre-clone to cache the registered git repositories.
	for _, r := range w.config.Repositories {
		// TODO: Clone repository another temporary destination
		repo, err := w.gitClient.Clone(ctx, r.RepoID, r.Remote, r.Branch, "")
		if err != nil {
			w.logger.Error("failed to clone repository",
				zap.String("repo-id", r.RepoID),
				zap.Error(err),
			)
			return fmt.Errorf("failed to clone repository %s: %w", r.RepoID, err)
		}
		w.gitRepos[r.RepoID] = repo
	}

	for _, cfg := range w.config.ImageProviders {
		p, err := imageprovider.NewProvider(&cfg, w.logger)
		if err != nil {
			return fmt.Errorf("failed to yield image provider %s: %w", cfg.Name, err)
		}

		w.wg.Add(1)
		go w.run(ctx, p, cfg.PullInterval.Duration())
	}
	w.wg.Wait()
	return nil
}

// run periodically compares the image in the given provider and one in git repository.
// And then pushes those with differences.
func (w *watcher) run(ctx context.Context, provider imageprovider.Provider, interval time.Duration) {
	defer w.wg.Done()
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// Inspect all targets in all repos, and compare only images managed by the given provider.
			for id, repo := range w.gitRepos {
				cfg, err := w.loadImageWatcherConfig(ctx, id, repo)
				if err != nil {
					w.logger.Error("failed to load image watcher config",
						zap.String("repo-id", id),
						zap.Error(err),
					)
					continue
				}
				for _, target := range cfg.Targets {
					if target.Provider != provider.Name() {
						continue
					}
					if err := w.updateOutdatedImage(ctx, &target, repo, provider); err != nil {
						w.logger.Error("failed to update image",
							zap.String("repo-id", id),
							zap.String("image-provider", provider.Name()),
							zap.Error(err),
						)
						continue
					}
				}
			}
		}
	}
}

// loadImageWatcherConfig gives back an Image Watcher Config for the given repo.
func (w *watcher) loadImageWatcherConfig(ctx context.Context, repoID string, repo git.Repo) (*config.ImageWatcherSpec, error) {
	w.mu.Lock()
	err := repo.Pull(ctx, repo.GetClonedBranch())
	w.mu.Unlock()
	if err != nil {
		return nil, fmt.Errorf("failed to perform git pull: %w", err)
	}

	var includes, excludes []string
	for _, repos := range w.config.ImageWatcher.Repos {
		if repos.RepoID == repoID {
			includes = repos.Includes
			excludes = repos.Excludes
			break
		}
	}
	cfg, ok, err := config.LoadImageWatcher(repo.GetPath(), includes, excludes)
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration file for Image Watcher: %w", err)
	}
	if !ok {
		return nil, fmt.Errorf("configuration file for Image Watcher not found: %w", err)
	}
	return cfg, nil
}

// updateOutdatedImage first compares the image in git repository and one in image provider.
// Then pushes rewritten one to the git repository if any deviation exists.
func (w *watcher) updateOutdatedImage(ctx context.Context, target *config.ImageWatcherTarget, repo git.Repo, provider imageprovider.Provider) error {
	// Fetch from the image provider.
	i, err := provider.ParseImage(target.Image)
	if err != nil {
		return fmt.Errorf("failed to parse image string \"%s\": %w", target.Image, err)
	}
	// TODO: Control not to reach the rate limit
	imageInRegistry, err := provider.GetLatestImage(ctx, i)
	if err != nil {
		return fmt.Errorf("failed to get latest image from %s: %w", provider.Name(), err)
	}

	// Fetch from the git repository.
	path := filepath.Join(repo.GetPath(), target.FilePath)
	yml, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}
	value, err := yamlprocessor.GetValue(yml, target.Field)
	if err != nil {
		return fmt.Errorf("failed to get value at %s in %s: %w", target.Field, target.FilePath, err)
	}
	imageInGit, ok := value.(string)
	if !ok {
		return fmt.Errorf("unknown value is defined at %s in %s", target.FilePath, target.Field)
	}

	outdated := imageInRegistry.String() != imageInGit
	if !outdated {
		return nil
	}

	// Update the outdated image.
	newYml, err := yamlprocessor.ReplaceValue(yml, target.Field, imageInRegistry.String())
	if err != nil {
		return fmt.Errorf("failed to replace value at %s with %s: %w", target.Field, imageInRegistry, err)
	}
	changes := map[string][]byte{
		target.FilePath: newYml,
	}
	// TODO: Make it changeable the commit message
	msg := fmt.Sprintf(defaultCommitMessageFormat, imageInGit, imageInRegistry.String(), target.Field, target.FilePath)
	w.mu.Lock()
	if err := repo.CommitChanges(ctx, repo.GetClonedBranch(), msg, false, changes); err != nil {
		return fmt.Errorf("failed to perform git commit: %w", err)
	}
	err = repo.Push(ctx, repo.GetClonedBranch())
	w.mu.Unlock()
	if err != nil {
		return fmt.Errorf("failed to perform git push: %w", err)
	}
	return nil
}
