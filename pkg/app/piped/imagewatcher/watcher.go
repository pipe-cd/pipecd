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
// that periodically checks the image registry and updates
// the image if there are differences with Git.
package imagewatcher

import (
	"context"
	"time"

	"go.uber.org/zap"
)

type Watcher interface {
	Run(context.Context) error
}

type watcher struct {
	timer  *time.Timer
	logger *zap.Logger
}

type imageRepos map[string]imageRepo
type imageRepo struct {
}

func NewWatcher(interval time.Duration, logger *zap.Logger) Watcher {
	return &watcher{
		timer:  time.NewTimer(interval),
		logger: logger,
	}
}

func (w *watcher) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-w.timer.C:
			reposInReg := w.fetchFromRegistry()
			reposInGit := w.fetchFromGit()
			outdated := calculateChanges(reposInReg, reposInGit)
			if err := w.update(outdated); err != nil {
				w.logger.Error("failed to update image", zap.Error(err))
			}
		}
	}
	return nil
}

func (w *watcher) fetchFromRegistry() imageRepos {
	return nil
}

func (w *watcher) fetchFromGit() imageRepos {
	return nil
}

func (w *watcher) update(targets imageRepos) error {
	return nil
}

// calculateChanges compares between image repos in the image registry and
// image repos in git. And then gives back image repos to be updated.
func calculateChanges(x, y imageRepos) imageRepos {
	return nil
}
