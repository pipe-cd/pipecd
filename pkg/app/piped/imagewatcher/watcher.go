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
	mu        sync.Mutex

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
	// Pre-clone to cache the registered git repositories.
	for _, r := range w.config.Repositories {
		repo, err := w.gitClient.Clone(ctx, r.RepoID, r.Remote, r.Branch, "")
		if err != nil {
			w.logger.Error("failed to clone repository",
				zap.String("repo-id", r.RepoID),
				zap.Error(err),
			)
			return err
		}
		w.gitRepos[r.RepoID] = repo
	}

	for _, cfg := range w.config.ImageProviders {
		p, err := imageprovider.NewProvider(&cfg, w.logger)
		if err != nil {
			return err
		}

		w.wg.Add(1)
		go w.run(ctx, p, cfg.PullInterval.Duration())
	}
	w.wg.Wait()
	return nil
}

// run periodically compares the image stored in the given provider and one stored in git.
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
			updates := make([]config.ImageWatcherTarget, 0)
			for id, repo := range w.gitRepos {
				u, err := w.determineUpdates(ctx, id, repo, provider)
				if err != nil {
					w.logger.Error("failed to determine images to be updated",
						zap.String("repo-id", id),
						zap.Error(err),
					)
					continue
				}
				updates = append(updates, u...)
			}
			if len(updates) == 0 {
				w.logger.Info("no image to be updated",
					zap.String("image-provider", provider.Name()),
				)
				continue
			}
			if err := update(updates); err != nil {
				w.logger.Error("failed to update image",
					zap.String("image-provider", provider.Name()),
					zap.Error(err),
				)
				continue
			}
		}
	}
}

// determineUpdates gives back target images to be updated for a given repo.
func (w *watcher) determineUpdates(ctx context.Context, repoID string, repo git.Repo, provider imageprovider.Provider) ([]config.ImageWatcherTarget, error) {
	branch := repo.GetClonedBranch()
	w.mu.Lock()
	err := repo.Pull(ctx, branch)
	w.mu.Unlock()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch from and integrate with a local branch: %w", err)
	}

	// Load Image Watcher Config for the given repo.
	var includes, excludes []string
	for _, target := range w.config.ImageWatcher.Repos {
		if target.RepoID == repoID {
			includes = target.Includes
			excludes = target.Excludes
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

	updates := make([]config.ImageWatcherTarget, 0)
	for _, target := range cfg.Targets {
		if provider.Name() != target.Provider {
			continue
		}
		outdated, err := checkOutdated(ctx, target, repo, provider)
		if err != nil {
			return nil, fmt.Errorf("failed to check the image is outdated: %w", err)
		}
		if outdated {
			updates = append(updates, target)
		}
	}
	return updates, nil
}

// checkOutdated checks if the image defined in the given target is identical to the one in image provider.
func checkOutdated(ctx context.Context, target config.ImageWatcherTarget, repo git.Repo, provider imageprovider.Provider) (bool, error) {
	i, err := provider.ParseImage(target.Image)
	if err != nil {
		return false, err
	}
	// TODO: Control not to reach the rate limit
	imageRef, err := provider.GetLatestImage(ctx, i)
	if err != nil {
		return false, err
	}

	yml, err := ioutil.ReadFile(filepath.Join(repo.GetPath(), target.FilePath))
	if err != nil {
		return false, err
	}
	value, err := yamlprocessor.GetValue(yml, target.Field)
	if err != nil {
		return false, err
	}
	v, ok := value.(string)
	if !ok {
		return false, fmt.Errorf("unknown value is defined at %s in %s", target.FilePath, target.Field)
	}
	return imageRef.String() != v, nil
}

func update(targets []config.ImageWatcherTarget) error {
	// TODO: Make it possible to push outdated images to Git
	return nil
}
