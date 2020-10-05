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
)

type Watcher interface {
	Run(context.Context) error
}

type watcher struct {
	timer *time.Timer
}

type imageRepos map[string]imageRepo
type imageRepo struct {
}

func NewWatcher(interval time.Duration) Watcher {
	return &watcher{
		timer: time.NewTimer(interval),
	}
}

func (w *watcher) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-w.timer.C:
			reposInReg := w.fetchImageReposFromRegistry()
			reposInGit := w.fetchImageReposFromGit()
			outdated := calculateChanges(reposInReg, reposInGit)
			if err := w.update(outdated); err != nil {
				// FIXME: Emit log
			}
		}
	}
	return nil
}

func (w *watcher) fetchImageReposFromRegistry() imageRepos {
	return nil
}

func (w *watcher) fetchImageReposFromGit() imageRepos {
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
