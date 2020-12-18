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
	"os"
	"path/filepath"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipe/pkg/app/piped/imageprovider"
	"github.com/pipe-cd/pipe/pkg/config"
	"github.com/pipe-cd/pipe/pkg/git"
	"github.com/pipe-cd/pipe/pkg/yamlprocessor"
)

const (
	defaultCommitMessageFormat = "Update image %s to %s defined at %s in %s"
	defaultCheckInterval       = 5 * time.Minute
)

type Watcher interface {
	Run(context.Context) error
}

type gitClient interface {
	Clone(ctx context.Context, repoID, remote, branch, destination string) (git.Repo, error)
}

type commit struct {
	changes map[string][]byte
	message string
}

type watcher struct {
	config    *config.PipedSpec
	gitClient gitClient
	logger    *zap.Logger
	wg        sync.WaitGroup

	// Indexed by the Image Provider name.
	providerCfgs map[string]config.PipedImageProvider
}

func NewWatcher(cfg *config.PipedSpec, gitClient gitClient, logger *zap.Logger) Watcher {
	return &watcher{
		config:    cfg,
		gitClient: gitClient,
		logger:    logger.Named("image-watcher"),
	}
}

// Run spawns goroutines for each git repository. They periodically pull the image
// from the container registry to compare the image with one in the git repository.
func (w *watcher) Run(ctx context.Context) error {
	w.providerCfgs = make(map[string]config.PipedImageProvider, len(w.config.ImageProviders))
	for _, cfg := range w.config.ImageProviders {
		w.providerCfgs[cfg.Name] = cfg
	}

	for _, repoCfg := range w.config.Repositories {
		repo, err := w.gitClient.Clone(ctx, repoCfg.RepoID, repoCfg.Remote, repoCfg.Branch, "")
		if err != nil {
			w.logger.Error("failed to clone repository",
				zap.String("repo-id", repoCfg.RepoID),
				zap.Error(err),
			)
			return fmt.Errorf("failed to clone repository %s: %w", repoCfg.RepoID, err)
		}

		w.wg.Add(1)
		go w.run(ctx, repo, &repoCfg)
	}

	w.wg.Wait()
	return nil
}

// run periodically compares the image in the given git repository and one in the image provider.
// And then pushes those with differences.
func (w *watcher) run(ctx context.Context, repo git.Repo, repoCfg *config.PipedRepository) {
	defer w.wg.Done()

	var (
		checkInterval              = defaultCheckInterval
		commitMsg                  string
		includedCfgs, excludedCfgs []string
	)
	// Use user-defined settings if there is.
	for _, r := range w.config.ImageWatcher.Repos {
		if r.RepoID != repoCfg.RepoID {
			continue
		}
		checkInterval = time.Duration(r.CheckInterval)
		commitMsg = r.CommitMessage
		includedCfgs = r.Includes
		excludedCfgs = r.Excludes
		break
	}

	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			err := repo.Pull(ctx, repo.GetClonedBranch())
			if err != nil {
				w.logger.Error("failed to perform git pull",
					zap.String("repo-id", repoCfg.RepoID),
					zap.String("branch", repo.GetClonedBranch()),
					zap.Error(err),
				)
				continue
			}
			cfg, ok, err := config.LoadImageWatcher(repo.GetPath(), includedCfgs, excludedCfgs)
			if err != nil {
				w.logger.Error("failed to load configuration file for Image Watcher",
					zap.String("repo-id", repoCfg.RepoID),
					zap.Error(err),
				)
				continue
			}
			if !ok {
				w.logger.Info("configuration file for Image Watcher not found",
					zap.String("repo-id", repoCfg.RepoID),
					zap.Error(err),
				)
				continue
			}
			if err := w.updateOutdatedImages(ctx, repo, cfg.Targets, commitMsg); err != nil {
				w.logger.Error("failed to update the targets",
					zap.String("repo-id", repoCfg.RepoID),
					zap.Error(err),
				)
			}
		}
	}
}

// updateOutdatedImages inspects all targets and pushes the changes to git repo if there is.
func (w *watcher) updateOutdatedImages(ctx context.Context, repo git.Repo, targets []config.ImageWatcherTarget, commitMsg string) error {
	commits := make([]*commit, 0)
	for _, t := range targets {
		c, err := w.checkOutdatedImage(ctx, &t, repo, commitMsg)
		if err != nil {
			w.logger.Error("failed to update image", zap.Error(err))
			continue
		}
		if c != nil {
			commits = append(commits, c)
		}
	}
	if len(commits) == 0 {
		return nil
	}

	// Copy the repo to another directory to avoid pull failure in the future.
	tmpDir, err := ioutil.TempDir("", "image-watcher")
	if err != nil {
		return fmt.Errorf("failed to create a new temporary directory: %w", err)
	}
	defer os.RemoveAll(tmpDir)
	tmpRepo, err := repo.Copy(tmpDir)
	if err != nil {
		return fmt.Errorf("failed to copy the repository to the temporary directory: %w", err)
	}
	for _, c := range commits {
		if err := tmpRepo.CommitChanges(ctx, tmpRepo.GetClonedBranch(), c.message, false, c.changes); err != nil {
			return fmt.Errorf("failed to perform git commit: %w", err)
		}
	}

	return tmpRepo.Push(ctx, tmpRepo.GetClonedBranch())
}

// checkOutdatedImage gives back a change content if any deviation exists
// between the image in the given git repository and one in the image provider.
func (w *watcher) checkOutdatedImage(ctx context.Context, target *config.ImageWatcherTarget, repo git.Repo, commitMsg string) (*commit, error) {
	// Retrieve the image from the image provider.
	providerCfg, ok := w.providerCfgs[target.Provider]
	if !ok {
		return nil, fmt.Errorf("unknown image provider %s is defined", target.Provider)
	}
	provider, err := imageprovider.NewProvider(&providerCfg, w.logger)
	if err != nil {
		return nil, fmt.Errorf("failed to yield image provider %s: %w", providerCfg.Name, err)
	}
	i, err := provider.ParseImage(target.Image)
	if err != nil {
		return nil, fmt.Errorf("failed to parse image string \"%s\": %w", target.Image, err)
	}
	// TODO: Control not to reach the rate limit
	imageInRegistry, err := provider.GetLatestImage(ctx, i)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest image from %s: %w", provider.Name(), err)
	}

	// Retrieve the image from the file cloned from the git repository.
	path := filepath.Join(repo.GetPath(), target.FilePath)
	yml, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	value, err := yamlprocessor.GetValue(yml, target.Field)
	if err != nil {
		return nil, fmt.Errorf("failed to get value at %s in %s: %w", target.Field, target.FilePath, err)
	}
	imageInGit, ok := value.(string)
	if !ok {
		return nil, fmt.Errorf("unknown value is defined at %s in %s", target.FilePath, target.Field)
	}

	outdated := imageInRegistry.String() != imageInGit
	if !outdated {
		return nil, nil
	}

	// Give back a change content.
	newYml, err := yamlprocessor.ReplaceValue(yml, target.Field, imageInRegistry.String())
	if err != nil {
		return nil, fmt.Errorf("failed to replace value at %s with %s: %w", target.Field, imageInRegistry, err)
	}
	if commitMsg == "" {
		commitMsg = fmt.Sprintf(defaultCommitMessageFormat, imageInGit, imageInRegistry.String(), target.Field, target.FilePath)
	}
	return &commit{
		changes: map[string][]byte{
			target.FilePath: newYml,
		},
		message: commitMsg,
	}, nil
}
