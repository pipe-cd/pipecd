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

// Package eventwatcher provides facilities to update config files when new
// event found. It can be done by periodically comparing the latest value user
// registered and the value in the files placed at Git repositories.
package eventwatcher

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipe/pkg/config"
	"github.com/pipe-cd/pipe/pkg/git"
	"github.com/pipe-cd/pipe/pkg/model"
	"github.com/pipe-cd/pipe/pkg/yamlprocessor"
)

const (
	// The latest value and Event name are supposed.
	defaultCommitMessageFormat = "Replace values with %q set by Event %q"
	defaultCheckInterval       = 5 * time.Minute
)

type Watcher interface {
	Run(context.Context) error
}

type eventGetter interface {
	GetLatest(ctx context.Context, name string, labels map[string]string) (*model.Event, bool)
}

type gitClient interface {
	Clone(ctx context.Context, repoID, remote, branch, destination string) (git.Repo, error)
}

type commit struct {
	changes map[string][]byte
	message string
}

type watcher struct {
	config      *config.PipedSpec
	eventGetter eventGetter
	gitClient   gitClient
	logger      *zap.Logger
	wg          sync.WaitGroup
}

func NewWatcher(cfg *config.PipedSpec, eventGetter eventGetter, gitClient gitClient, logger *zap.Logger) Watcher {
	return &watcher{
		config:      cfg,
		eventGetter: eventGetter,
		gitClient:   gitClient,
		logger:      logger.Named("event-watcher"),
	}
}

// Run spawns goroutines for each git repository. They periodically fetch the latest Event
// from the control-plane to compare the value with one in the git repository.
func (w *watcher) Run(ctx context.Context) error {
	w.logger.Info("start running event watcher")

	for _, repoCfg := range w.config.Repositories {
		repo, err := w.gitClient.Clone(ctx, repoCfg.RepoID, repoCfg.Remote, repoCfg.Branch, "")
		if err != nil {
			w.logger.Error("failed to clone repository",
				zap.String("repo-id", repoCfg.RepoID),
				zap.Error(err),
			)
			return fmt.Errorf("failed to clone repository %s: %w", repoCfg.RepoID, err)
		}
		defer os.RemoveAll(repo.GetPath())

		w.wg.Add(1)
		go w.run(ctx, repo, &repoCfg)
	}

	w.wg.Wait()
	return nil
}

// run works against a single git repo. It periodically compares the value in the given
// git repository and one in the control-plane. And then pushes those with differences.
func (w *watcher) run(ctx context.Context, repo git.Repo, repoCfg *config.PipedRepository) {
	defer w.wg.Done()

	var (
		commitMsg                  string
		includedCfgs, excludedCfgs []string
	)
	// Use user-defined settings if there is.
	for _, r := range w.config.EventWatcher.GitRepos {
		if r.RepoID != repoCfg.RepoID {
			continue
		}
		commitMsg = r.CommitMessage
		includedCfgs = r.Includes
		excludedCfgs = r.Excludes
		break
	}
	checkInterval := time.Duration(w.config.EventWatcher.CheckInterval)
	if checkInterval == 0 {
		checkInterval = defaultCheckInterval
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
			cfg, err := config.LoadEventWatcher(repo.GetPath(), includedCfgs, excludedCfgs)
			if errors.Is(err, config.ErrNotFound) {
				w.logger.Info("configuration file for Event Watcher not found",
					zap.String("repo-id", repoCfg.RepoID),
					zap.Error(err),
				)
				continue
			}
			if err != nil {
				w.logger.Error("failed to load configuration file for Event Watcher",
					zap.String("repo-id", repoCfg.RepoID),
					zap.Error(err),
				)
				continue
			}
			if err := w.updateValues(ctx, repo, cfg.Events, commitMsg); err != nil {
				w.logger.Error("failed to update the values",
					zap.String("repo-id", repoCfg.RepoID),
					zap.Error(err),
				)
			}
		}
	}
}

// updateValues inspects all Event-definition and pushes the changes to git repo if there is.
func (w *watcher) updateValues(ctx context.Context, repo git.Repo, events []config.EventWatcherEvent, commitMsg string) error {
	// Copy the repo to another directory to avoid pull failure in the future.
	tmpDir, err := ioutil.TempDir("", "event-watcher")
	if err != nil {
		return fmt.Errorf("failed to create a new temporary directory: %w", err)
	}
	defer os.RemoveAll(tmpDir)
	tmpRepo, err := repo.Copy(filepath.Join(tmpDir, "tmp-repo"))
	if err != nil {
		return fmt.Errorf("failed to copy the repository to the temporary directory: %w", err)
	}

	commits := make([]*commit, 0)
	for _, e := range events {
		c, err := w.modifyFiles(ctx, &e, tmpRepo, commitMsg)
		if err != nil {
			w.logger.Error("failed to check outdated value", zap.Error(err))
			continue
		}
		if c != nil {
			commits = append(commits, c)
		}
	}
	if len(commits) == 0 {
		return nil
	}

	w.logger.Info(fmt.Sprintf("there are %d outdated values", len(commits)))
	for _, c := range commits {
		if err := tmpRepo.CommitChanges(ctx, tmpRepo.GetClonedBranch(), c.message, false, c.changes); err != nil {
			return fmt.Errorf("failed to perform git commit: %w", err)
		}
	}
	return tmpRepo.Push(ctx, tmpRepo.GetClonedBranch())
}

// modifyFiles modifies files defined in a given Event if any deviation exists between the value in
// the git repository and one in the control-plane. And gives back a change contents.
func (w *watcher) modifyFiles(ctx context.Context, event *config.EventWatcherEvent, repo git.Repo, commitMsg string) (*commit, error) {
	latestEvent, ok := w.eventGetter.GetLatest(ctx, event.Name, event.Labels)
	if !ok {
		return nil, fmt.Errorf("failed to get the latest Event with the name %q", event.Name)
	}

	// Determine files to be changed.
	changes := make(map[string][]byte, 0)
	for _, r := range event.Replacements {
		path := filepath.Join(repo.GetPath(), r.File)
		yml, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to read file: %w", err)
		}
		v, err := yamlprocessor.GetValue(yml, r.YAMLField)
		if err != nil {
			return nil, fmt.Errorf("failed to get value at %s in %s: %w", r.YAMLField, r.File, err)
		}
		value, err := convertStr(v)
		if err != nil {
			return nil, fmt.Errorf("a value of unknown type is defined at %s in %s: %w", err, r.YAMLField, r.File)
		}
		if latestEvent.Data == value {
			// Already up-to-date.
			continue
		}
		// Modify the local file and put it into the change list.
		newYml, err := yamlprocessor.ReplaceValue(yml, r.YAMLField, latestEvent.Data)
		if err != nil {
			return nil, fmt.Errorf("failed to replace value at %s with %s: %w", r.YAMLField, latestEvent.Data, err)
		}
		if err := ioutil.WriteFile(path, newYml, os.ModePerm); err != nil {
			return nil, fmt.Errorf("failed to write file: %w", err)
		}
		changes[r.File] = newYml
	}

	if len(changes) == 0 {
		return nil, nil
	}

	if commitMsg == "" {
		commitMsg = fmt.Sprintf(defaultCommitMessageFormat, latestEvent.Data, event.Name)
	}
	return &commit{
		changes: changes,
		message: commitMsg,
	}, nil
}

// convertStr converts a given value into a string.
func convertStr(value interface{}) (out string, err error) {
	switch v := value.(type) {
	case string:
		out = v
	case int:
		out = strconv.Itoa(v)
	case int64:
		out = strconv.FormatInt(v, 10)
	case uint64:
		out = strconv.FormatUint(v, 10)
	case float64:
		out = strconv.FormatFloat(v, 'f', -1, 64)
	case bool:
		out = strconv.FormatBool(v)
	default:
		err = fmt.Errorf("failed to convert %T into string", v)
	}
	return
}
