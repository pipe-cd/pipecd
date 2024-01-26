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

package git

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
)

const (
	defaultUsername = "piped"
	defaultEmail    = "pipecd.dev@gmail.com"
)

// Client is a git client for cloning/fetching git repo.
// It keeps a local cache for faster future cloning.
type Client interface {
	// Clone clones a specific git repository to the given destination.
	Clone(ctx context.Context, repoID, remote, branch, destination string) (Repo, error)
	// Clean removes all cache data.
	Clean() error
}

type client struct {
	username  string
	email     string
	gitPath   string
	cacheDir  string
	mu        sync.Mutex
	repoLocks map[string]*sync.Mutex

	gitEnvs         []string
	gitEnvsByRepo   map[string][]string
	gitGCAutoDetach bool
	logger          *zap.Logger
}

type Option func(*client)

func WithGitEnv(env string) Option {
	return func(c *client) {
		c.gitEnvs = append(c.gitEnvs, env)
	}
}

func WithGitEnvForRepo(remote string, env string) Option {
	return func(c *client) {
		c.gitEnvsByRepo[remote] = append(c.gitEnvsByRepo[remote], env)
	}
}

func WithLogger(logger *zap.Logger) Option {
	return func(c *client) {
		c.logger = logger
	}
}

func WithUserName(n string) Option {
	return func(c *client) {
		if n != "" {
			c.username = n
		}
	}
}

func WithEmail(e string) Option {
	return func(c *client) {
		if e != "" {
			c.email = e
		}
	}
}

func WithAutoDetach(a bool) Option {
	return func(c *client) {
		c.gitGCAutoDetach = a
	}
}

// NewClient creates a new CLient instance for cloning git repositories.
// After using Clean should be called to delete cache data.
func NewClient(opts ...Option) (Client, error) {
	gitPath, err := exec.LookPath("git")
	if err != nil {
		return nil, fmt.Errorf("unable to find the path of git: %v", err)
	}

	cacheDir, err := os.MkdirTemp("", "gitcache")
	if err != nil {
		return nil, fmt.Errorf("unable to create a temporary directory for git cache: %v", err)
	}

	c := &client{
		username:        defaultUsername,
		email:           defaultEmail,
		gitPath:         gitPath,
		cacheDir:        cacheDir,
		repoLocks:       make(map[string]*sync.Mutex),
		gitEnvsByRepo:   make(map[string][]string, 0),
		gitGCAutoDetach: true,
		logger:          zap.NewNop(),
	}

	for _, opt := range opts {
		opt(c)
	}

	return c, nil
}

// Clone clones a specific git repository to the given destination.
func (c *client) Clone(ctx context.Context, repoID, remote, branch, destination string) (Repo, error) {
	var (
		repoCachePath = filepath.Join(c.cacheDir, repoID)
		logger        = c.logger.With(
			zap.String("repo-id", repoID),
			zap.String("remote", remote),
			zap.String("repo-cache-path", repoCachePath),
		)
	)

	c.lockRepo(repoID)
	defer c.unlockRepo(repoID)

	_, err := os.Stat(repoCachePath)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	if os.IsNotExist(err) {
		// Cache miss, clone for the first time.
		logger.Info(fmt.Sprintf("cloning %s for the first time", repoID))
		if err := os.MkdirAll(filepath.Dir(repoCachePath), os.ModePerm); err != nil && !os.IsExist(err) {
			return nil, err
		}
		out, err := retryCommand(3, time.Second, logger, func() ([]byte, error) {
			return runGitCommand(ctx, c.gitPath, "", c.envsForRepo(remote), "clone", "--mirror", remote, repoCachePath)
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
		c.logger.Info(fmt.Sprintf("fetching %s to update the cache", repoID))
		out, err := retryCommand(3, time.Second, c.logger, func() ([]byte, error) {
			return runGitCommand(ctx, c.gitPath, repoCachePath, c.envsForRepo(remote), "fetch")
		})
		if err != nil {
			logger.Error("failed to fetch from remote",
				zap.String("out", string(out)),
				zap.Error(err),
			)
			return nil, fmt.Errorf("failed to fetch: %v", err)
		}
	}

	if destination != "" {
		err = os.MkdirAll(destination, os.ModePerm)
		if err != nil {
			return nil, err
		}
	} else {
		destination, err = os.MkdirTemp("", "git")
		if err != nil {
			return nil, err
		}
	}

	args := []string{"clone"}
	if branch != "" {
		args = append(args, "-b", branch)
	}
	args = append(args, repoCachePath, destination)
	if out, err := runGitCommand(ctx, c.gitPath, "", c.envsForRepo(remote), args...); err != nil {
		logger.Error("failed to clone from local",
			zap.String("out", string(out)),
			zap.String("branch", branch),
			zap.String("repo-path", destination),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to clone from local: %v", err)
	}

	r := NewRepo(destination, c.gitPath, remote, branch, c.envsForRepo(remote))
	if c.username != "" || c.email != "" {
		if err := r.setUser(ctx, c.username, c.email); err != nil {
			return nil, fmt.Errorf("failed to set user: %v", err)
		}
	}

	if err := r.setAutoDetach(ctx, c.gitGCAutoDetach); err != nil {
		return nil, fmt.Errorf("failed to set auto detach: %v", err)
	}

	// Because we did a local cloning so the remote url of origin
	// is the path to the cache directory.
	// We do this change to correct it.
	if err := r.setRemote(ctx, remote); err != nil {
		return nil, fmt.Errorf("failed to set remote: %v", err)
	}

	return r, nil
}

// Clean removes all cache data.
func (c *client) Clean() error {
	return os.RemoveAll(c.cacheDir)
}

// getLatestRemoteHashForBranch returns the hash of the latest commit of a remote branch.
func (c *client) getLatestRemoteHashForBranch(ctx context.Context, remote, branch string) (string, error) {
	ref := "refs/heads/" + branch
	out, err := retryCommand(3, time.Second, c.logger, func() ([]byte, error) {
		return runGitCommand(ctx, c.gitPath, "", c.envsForRepo(remote), "ls-remote", ref)
	})
	if err != nil {
		c.logger.Error("failed to get latest remote hash for branch",
			zap.String("remote", remote),
			zap.String("branch", branch),
			zap.String("out", string(out)),
			zap.Error(err),
		)
		return "", err
	}
	parts := strings.Split(string(out), "\t")
	return parts[0], nil
}

func (c *client) lockRepo(repoID string) {
	c.mu.Lock()
	if _, ok := c.repoLocks[repoID]; !ok {
		c.repoLocks[repoID] = &sync.Mutex{}
	}
	mu := c.repoLocks[repoID]
	c.mu.Unlock()

	mu.Lock()
}

func (c *client) unlockRepo(repoID string) {
	c.mu.Lock()
	c.repoLocks[repoID].Unlock()
	c.mu.Unlock()
}

func (c *client) envsForRepo(remote string) []string {
	envs := c.gitEnvsByRepo[remote]
	return append(envs, c.gitEnvs...)
}

func runGitCommand(ctx context.Context, execPath, dir string, envs []string, args ...string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, execPath, args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), envs...)
	return cmd.CombinedOutput()
}

// retryCommand retries a command a few times with a constant backoff.
//
//nolint:unparam
func retryCommand(retries int, interval time.Duration, logger *zap.Logger, commander func() ([]byte, error)) (out []byte, err error) {
	for i := 0; i < retries; i++ {
		out, err = commander()
		if err == nil {
			return
		}
		logger.Warn(fmt.Sprintf("command was failed %d times, sleep %v before retrying command", i+1, interval))
		time.Sleep(interval)
	}
	return
}
