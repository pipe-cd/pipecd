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
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipe/pkg/app/piped/imageprovider"
	"github.com/pipe-cd/pipe/pkg/config"
	"github.com/pipe-cd/pipe/pkg/git"
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

	mu sync.RWMutex
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

// Run spawns goroutines for each image provider.
func (w *watcher) Run(ctx context.Context) error {
	// Pre-clone to cache the registered git repositories.
	w.gitRepos = make(map[string]git.Repo, len(w.config.Repositories))
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

		go w.run(ctx, p, cfg.PullInterval.Duration())
	}
	return nil
}

// run periodically compares the image stored in the given provider and one stored in git.
// And then pushes those with differences.
func (w *watcher) run(ctx context.Context, provider imageprovider.Provider, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			targets := make([]config.ImageWatcherTarget, 0)
			// Collect target images for each git repository.
			for id, repo := range w.gitRepos {
				w.mu.RLock()
				branch := repo.GetClonedBranch()
				w.mu.RUnlock()
				if err := repo.Pull(ctx, branch); err != nil {
					w.logger.Error("failed to update repository branch",
						zap.String("repo-id", id),
						zap.Error(err),
					)
					continue
				}

				cfg, ok, err := config.LoadImageWatcher(repo.GetPath())
				if err != nil {
					w.logger.Error("failed to load configuration file for Image Watcher", zap.Error(err))
					continue
				}
				if !ok {
					w.logger.Error("configuration file for Image Watcher not found", zap.Error(err))
					continue
				}
				t := filterTargets(provider.Name(), cfg.Targets)
				targets = append(targets, t...)
			}

			outdated, err := determineUpdates(ctx, targets, provider)
			if err != nil {
				w.logger.Error("failed to determine which one should be updated", zap.Error(err))
				continue
			}
			if len(outdated) == 0 {
				w.logger.Info("no image to be updated")
				continue
			}
			if err := update(outdated); err != nil {
				w.logger.Error("failed to update image", zap.Error(err))
				continue
			}
		}
	}
}

// filterTargets gives back the targets corresponding to the given provider.
func filterTargets(provider string, targets []config.ImageWatcherTarget) (filtered []config.ImageWatcherTarget) {
	for _, t := range targets {
		if t.Provider == provider {
			filtered = append(filtered, t)
		}
	}
	return
}

// determineUpdates gives back target images to be updated.
func determineUpdates(ctx context.Context, targets []config.ImageWatcherTarget, provider imageprovider.Provider) (outdated []config.ImageWatcherTarget, err error) {
	for _, target := range targets {
		i, err := provider.ParseImage(target.Image)
		if err != nil {
			return nil, err
		}
		// TODO: Control not to reach the rate limit
		_, err = provider.GetLatestImage(ctx, i)
		if err != nil {
			return nil, err
		}
		// TODO: Compares between image repos in the image registry and image repos in git
		//   And then gives back image repos to be updated.
	}

	return
}

func update(targets []config.ImageWatcherTarget) error {
	// TODO: Make it possible to push outdated images to Git
	return nil
}
