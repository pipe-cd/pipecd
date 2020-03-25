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
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestTag(t *testing.T) {
	faker, err := newFaker()
	require.NoError(t, err)
	defer faker.clean()

	var (
		org      = "test-repo-org"
		repoName = "repo-1"
	)

	err = faker.makeRepo(org, repoName)
	require.NoError(t, err)
	r := &repo{
		dir:     faker.repoDir(org, repoName),
		gitPath: faker.gitPath,
		logger:  zap.NewNop(),
	}

	ctx := context.Background()
	tag, err := r.GetLatestTag(ctx)
	assert.Equal(t, ErrNoTag, err)
	assert.Equal(t, "", tag)

	err = r.CreateTag(ctx, "v0.0.1", "version v0.0.1")
	require.NoError(t, err)

	tag, err = r.GetLatestTag(ctx)
	require.NoError(t, err)
	assert.Equal(t, "v0.0.1", tag)
}

func TestGetLatestCommitID(t *testing.T) {
	faker, err := newFaker()
	require.NoError(t, err)
	defer faker.clean()

	var (
		org      = "test-repo-org"
		repoName = "repo-2"
	)

	err = faker.makeRepo(org, repoName)
	require.NoError(t, err)
	r := &repo{
		dir:     faker.repoDir(org, repoName),
		gitPath: faker.gitPath,
		logger:  zap.NewNop(),
	}

	ctx := context.Background()
	commits, err := r.ListCommits(ctx, "")
	require.NoError(t, err)
	assert.Equal(t, 1, len(commits))

	latestCommitID, err := r.GetLatestCommitID(ctx)
	require.NoError(t, err)
	assert.Equal(t, commits[0].Hash, latestCommitID)
}

func TestAddCommit(t *testing.T) {
	faker, err := newFaker()
	require.NoError(t, err)
	defer faker.clean()

	var (
		org      = "test-repo-org"
		repoName = "repo-3"
	)

	err = faker.makeRepo(org, repoName)
	require.NoError(t, err)
	r := &repo{
		dir:     faker.repoDir(org, repoName),
		gitPath: faker.gitPath,
		logger:  zap.NewNop(),
	}

	ctx := context.Background()
	commits, err := r.ListCommits(ctx, "")
	require.NoError(t, err)
	require.Equal(t, 1, len(commits))

	path := filepath.Join(r.dir, "new-file.txt")
	err = ioutil.WriteFile(path, []byte("content"), os.ModePerm)
	require.NoError(t, err)

	err = r.AddCommit(ctx, "Added new file")
	require.NoError(t, err)

	err = r.AddCommit(ctx, "No change")
	require.Equal(t, ErrNoChange, err)

	commits, err = r.ListCommits(ctx, "")
	require.NoError(t, err)
	require.Equal(t, 2, len(commits))
	assert.Equal(t, "Added new file", commits[0].Message)
}

func TestBranch(t *testing.T) {
	faker, err := newFaker()
	require.NoError(t, err)
	defer faker.clean()

	var (
		org      = "test-repo-org"
		repoName = "repo-4"
	)

	err = faker.makeRepo(org, repoName)
	require.NoError(t, err)
	r := &repo{
		dir:     faker.repoDir(org, repoName),
		gitPath: faker.gitPath,
		logger:  zap.NewNop(),
	}

	ctx := context.Background()
	branch, err := r.GetBranch(ctx)
	require.NoError(t, err)
	require.Equal(t, "master", branch)

	err = r.CheckoutNewBranch(ctx, "new-branch")
	require.NoError(t, err)

	branch, err = r.GetBranch(ctx)
	require.NoError(t, err)
	require.Equal(t, "new-branch", branch)
}

func TestCommitChanges(t *testing.T) {
	faker, err := newFaker()
	require.NoError(t, err)
	defer faker.clean()

	var (
		org      = "test-repo-org"
		repoName = "repo-commit-changes"
	)
	err = faker.makeRepo(org, repoName)
	require.NoError(t, err)
	r := &repo{
		dir:     faker.repoDir(org, repoName),
		gitPath: faker.gitPath,
		logger:  zap.NewNop(),
	}

	ctx := context.Background()
	commits, err := r.ListCommits(ctx, "")
	require.NoError(t, err)
	require.Equal(t, 1, len(commits))

	changes := map[string][]byte{
		"README.md":     []byte("new-readme"),
		"a/b/c/new.txt": []byte("new-hello"),
	}
	err = r.CommitChanges(ctx, "new-branch", "New commit with changes", true, changes)
	require.NoError(t, err)

	commits, err = r.ListCommits(ctx, "")
	require.NoError(t, err)
	require.Equal(t, 2, len(commits))
	assert.Equal(t, "New commit with changes", commits[0].Message)

	bytes, err := ioutil.ReadFile(filepath.Join(r.dir, "README.md"))
	require.NoError(t, err)
	assert.Equal(t, string(changes["README.md"]), string(bytes))

	bytes, err = ioutil.ReadFile(filepath.Join(r.dir, "a/b/c/new.txt"))
	require.NoError(t, err)
	assert.Equal(t, string(changes["a/b/c/new.txt"]), string(bytes))
}
