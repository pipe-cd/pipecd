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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestClone(t *testing.T) {
	faker, err := newFaker()
	require.NoError(t, err)
	defer faker.clean()

	c, err := NewClient()
	require.NoError(t, err)
	require.NotNil(t, c)
	defer c.Clean()

	err = faker.makeRepo("test-clone-org", "repo-1")
	require.NoError(t, err)
	err = faker.makeRepo("test-clone-org", "repo-2")
	require.NoError(t, err)

	ctx := context.Background()

	repo1Path, err := os.MkdirTemp("", "repo1path")
	require.NoError(t, err)
	repo1, err := c.Clone(ctx, "repo-1", filepath.Join(faker.dir, "test-clone-org/repo-1"), "", repo1Path)
	require.NoError(t, err)
	require.NotNil(t, repo1)
	defer func() {
		assert.NoError(t, repo1.Clean())
	}()
	commits1, err := repo1.ListCommits(ctx, "")
	require.NoError(t, err)
	assert.Equal(t, 1, len(commits1))

	repo2Path, err := os.MkdirTemp("", "repo2path")
	require.NoError(t, err)
	repo2, err := c.Clone(ctx, "repo-2", filepath.Join(faker.dir, "test-clone-org/repo-2"), "", repo2Path)
	require.NoError(t, err)
	require.NotNil(t, repo2)
	defer func() {
		assert.NoError(t, repo2.Clean())
	}()
	commits2, err := repo1.ListCommits(ctx, "")
	require.NoError(t, err)
	assert.Equal(t, 1, len(commits2))

	// Make sure client fetches the update.
	commander := gitCommander{
		gitPath: c.(*client).gitPath,
		dir:     faker.dir,
		org:     "test-clone-org",
		repo:    "repo-1",
	}
	err = commander.addCommit("note.txt", "note.text context")
	require.NoError(t, err)
	repo12Path, err := os.MkdirTemp("", "repo12path")
	require.NoError(t, err)
	repo12, err := c.Clone(ctx, "repo-1", filepath.Join(faker.dir, "test-clone-org/repo-1"), "master", repo12Path)
	require.NoError(t, err)
	require.NotNil(t, repo12)
	defer func() {
		assert.NoError(t, repo12.Clean())
	}()
	commits12, err := repo12.ListCommits(ctx, "")
	require.NoError(t, err)
	require.Equal(t, 2, len(commits12))
	assert.Equal(t, "Added note.txt", commits12[0].Message)
}

type faker struct {
	dir     string
	gitPath string
}

func newFaker() (*faker, error) {
	gitPath, err := exec.LookPath("git")
	if err != nil {
		return nil, err
	}
	remoteDir, err := os.MkdirTemp("", "remote")
	if err != nil {
		return nil, err
	}
	return &faker{
		dir:     remoteDir,
		gitPath: gitPath,
	}, nil
}

func (f *faker) clean() {
	os.RemoveAll(f.dir)
}

//nolint:unparam
func (f *faker) repoDir(org, repo string) string {
	return filepath.Join(f.dir, org, repo)
}

func (f *faker) makeRepo(org, repo string) error {
	rdir := filepath.Join(f.dir, org, repo)
	if err := os.MkdirAll(rdir, os.ModePerm); err != nil {
		return err
	}
	commander := gitCommander{
		gitPath: f.gitPath,
		dir:     f.dir,
		org:     org,
		repo:    repo,
	}

	err := commander.runGitCommands([][]string{
		{"init", "--initial-branch", "master"},
	})
	if err != nil {
		return err
	}

	content := fmt.Sprintf("Hello, %s/%s.\n", org, repo)
	return commander.addCommit("README.md", content)
}

type gitCommander struct {
	gitPath string
	dir     string
	org     string
	repo    string
}

func (g gitCommander) runGitCommands(commands [][]string) error {
	rdir := filepath.Join(g.dir, g.org, g.repo)
	for _, cmds := range commands {
		c := exec.Command(g.gitPath, cmds...)
		c.Dir = rdir
		if b, err := c.CombinedOutput(); err != nil {
			return fmt.Errorf("%s %v: %v, %s", g.gitPath, cmds, err, string(b))
		}
	}
	return nil
}

func (g gitCommander) addCommit(filename string, content string) error {
	rdir := filepath.Join(g.dir, g.org, g.repo)
	path := filepath.Join(rdir, filename)
	if err := os.WriteFile(path, []byte(content), os.ModePerm); err != nil {
		return err
	}
	return g.runGitCommands([][]string{
		{"add", "."},
		{"config", "user.email", "test@gmail.com"},
		{"config", "user.name", "test-user"},
		{"commit", "-m", fmt.Sprintf("Added %s", filename)},
	})
}

func TestCloneUsingPasswordAuth(t *testing.T) {
	url, err := includePasswordAuthRemote("https://example.com/org/repo", "test-user", "test-password")
	require.NoError(t, err)
	assert.Equal(t, "https://test-user:test-password@example.com/org/repo", url)
}

func TestRetryCommand(t *testing.T) {
	var (
		ranCount   = 0
		commandOut = []byte("hello")
		commandErr = fmt.Errorf("test-error")
		logger     = zap.NewNop()
	)
	testcases := []struct {
		name             string
		commandSuccessAt int
		expectedError    error
	}{
		{
			name:             "success at the first time",
			commandSuccessAt: 1,
			expectedError:    nil,
		},
		{
			name:             "success at the second time",
			commandSuccessAt: 2,
			expectedError:    nil,
		},
		{
			name:             "failure at all",
			commandSuccessAt: 5,
			expectedError:    commandErr,
		},
	}
	for _, tc := range testcases {
		ranCount = 0
		out, err := retryCommand(3, time.Millisecond, logger, func() ([]byte, error) {
			ranCount++
			if tc.commandSuccessAt == ranCount {
				return commandOut, nil
			}
			return commandOut, commandErr
		})
		assert.Equal(t, commandOut, out)
		assert.Equal(t, tc.expectedError, err)
	}
}
