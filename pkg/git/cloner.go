// Copyright 2020 The Pipe Authors.
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

package git

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"go.uber.org/zap"
)

// Cloner is a git client for cloning GitHub repository.
// It keeps a local cache for faster future cloning.
// TODO: Limit the cache size. (LRU?)
type Cloner struct {
	gitPath   string
	cacheDir  string
	mu        sync.Mutex
	repoLocks map[string]*sync.Mutex
	logger    *zap.Logger
}

// NewCloner creates a new Cloner instance for cloning GitHub repositories.
// After using Clean should be called to delete cache data.
func NewCloner(logger *zap.Logger) (*Cloner, error) {
	gitPath, err := exec.LookPath("git")
	if err != nil {
		return nil, fmt.Errorf("unabled to find the path of git: %v", err)
	}
	cacheDir, err := ioutil.TempDir("", "gitcache")
	if err != nil {
		return nil, fmt.Errorf("unabled to create a temporary directory for git cache: %v", err)
	}
	return &Cloner{
		gitPath:   gitPath,
		cacheDir:  cacheDir,
		repoLocks: make(map[string]*sync.Mutex),
		logger:    logger,
	}, nil
}

// Clone clones a specific GitHub repository.
func (c *Cloner) Clone(ctx context.Context, base, repoFullname, username, email string) (Repo, error) {
	remote := fmt.Sprintf("%s/%s", base, repoFullname)

	c.lockRepo(repoFullname)
	defer c.unlockRepo(repoFullname)

	repoCachePath := filepath.Join(c.cacheDir, repoFullname) + ".git"
	_, err := os.Stat(repoCachePath)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	logger := c.logger.With(
		zap.String("repo", repoFullname),
		zap.String("repo-cache-path", repoCachePath),
	)

	if os.IsNotExist(err) {
		// Cache miss, clone for the first time.
		logger.Info(fmt.Sprintf("cloning %s for the first time", repoFullname))
		if err := os.MkdirAll(filepath.Dir(repoCachePath), os.ModePerm); err != nil && !os.IsExist(err) {
			return nil, err
		}
		out, err := retryCommand(3, time.Second, logger, func() ([]byte, error) {
			return c.runGitCommand(ctx, "", "clone", "--mirror", remote, repoCachePath)
		})
		if err != nil {
			logger.Error("failed to clone from remote",
				zap.String("out", string(out)),
				zap.Error(err),
			)
			return nil, fmt.Errorf("failed to clone from remote: %v", err)
		}
	} else {
		// Cache hit. Do a git fetch to keep updated.
		c.logger.Info(fmt.Sprintf("fetching %s to update the cache", repoFullname))
		out, err := retryCommand(3, time.Second, c.logger, func() ([]byte, error) {
			return c.runGitCommand(ctx, repoCachePath, "fetch")
		})
		if err != nil {
			logger.Error("failed to fetch from remote",
				zap.String("out", string(out)),
				zap.Error(err),
			)
			return nil, fmt.Errorf("failed to fetch: %v", err)
		}
	}

	repoPath, err := ioutil.TempDir("", "git")
	if err != nil {
		return nil, err
	}
	if out, err := c.runGitCommand(ctx, "", "clone", repoCachePath, repoPath); err != nil {
		logger.Error("failed to clone from local",
			zap.String("out", string(out)),
			zap.String("repo-path", repoPath),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to clone from local: %v", err)
	}
	r := &repo{
		dir:     repoPath,
		gitPath: c.gitPath,
		remote:  remote,
		logger: c.logger.With(
			zap.String("repo", repoFullname),
		),
	}
	if err := r.SetUser(ctx, username, email); err != nil {
		return nil, fmt.Errorf("failed to set user: %v", err)
	}
	return r, nil
}

// Clean removes all cache data.
func (c *Cloner) Clean() error {
	return os.RemoveAll(c.cacheDir)
}

func (c *Cloner) lockRepo(repoFullname string) {
	c.mu.Lock()
	if _, ok := c.repoLocks[repoFullname]; !ok {
		c.repoLocks[repoFullname] = &sync.Mutex{}
	}
	mu := c.repoLocks[repoFullname]
	c.mu.Unlock()

	mu.Lock()
}

func (c *Cloner) unlockRepo(repoFullname string) {
	c.mu.Lock()
	c.repoLocks[repoFullname].Unlock()
	c.mu.Unlock()
}

func (c *Cloner) runGitCommand(ctx context.Context, dir string, args ...string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, c.gitPath, args...)
	cmd.Dir = dir
	return cmd.CombinedOutput()
}

// retryCommand retries a command a few times with a constant backoff.
func retryCommand(retries int, internal time.Duration, logger *zap.Logger, commander func() ([]byte, error)) (out []byte, err error) {
	for i := 0; i < retries; i++ {
		out, err = commander()
		if err == nil {
			return
		}
		logger.Warn(fmt.Sprintf("command was failed %d times, sleep %d seconds before retrying command", i+1, internal))
		time.Sleep(internal)
	}
	return
}
