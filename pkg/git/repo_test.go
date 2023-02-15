// Copyright 2023 The PipeCD Authors.
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
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetCommitHashForRev(t *testing.T) {
	faker, err := newFaker()
	require.NoError(t, err)
	defer faker.clean()

	var (
		org      = "test-repo-org"
		repoName = "repo-get-commit-hash-for-rev"
		ctx      = context.Background()
	)

	err = faker.makeRepo(org, repoName)
	require.NoError(t, err)
	r := &repo{
		dir:     faker.repoDir(org, repoName),
		gitPath: faker.gitPath,
	}

	commits, err := r.ListCommits(ctx, "")
	require.NoError(t, err)
	assert.Equal(t, 1, len(commits))

	latestCommitHash, err := r.GetCommitHashForRev(ctx, "HEAD")
	require.NoError(t, err)
	assert.Equal(t, commits[0].Hash, latestCommitHash)
}

func TestChangedFiles(t *testing.T) {
	faker, err := newFaker()
	require.NoError(t, err)
	defer faker.clean()

	var (
		org      = "test-repo-org"
		repoName = "repo-changed-files"
		ctx      = context.Background()
	)

	err = faker.makeRepo(org, repoName)
	require.NoError(t, err)
	r := &repo{
		dir:     faker.repoDir(org, repoName),
		gitPath: faker.gitPath,
	}

	previousCommitHash, err := r.GetCommitHashForRev(ctx, "HEAD")
	require.NoError(t, err)
	require.NotEqual(t, "", previousCommitHash)

	err = os.MkdirAll(filepath.Join(r.dir, "new-dir"), os.ModePerm)
	require.NoError(t, err)
	path := filepath.Join(r.dir, "new-dir", "new-file.txt")
	err = os.WriteFile(path, []byte("content"), os.ModePerm)
	require.NoError(t, err)

	readmeFilePath := filepath.Join(r.dir, "README.md")
	err = os.WriteFile(readmeFilePath, []byte("new content"), os.ModePerm)
	require.NoError(t, err)

	err = r.addCommit(ctx, "Added new file")
	require.NoError(t, err)

	headCommitHash, err := r.GetCommitHashForRev(ctx, "HEAD")
	require.NoError(t, err)
	require.NotEqual(t, "", headCommitHash)

	changedFiles, err := r.ChangedFiles(ctx, previousCommitHash, headCommitHash)
	sort.Strings(changedFiles)
	expectedChangedFiles := []string{
		"new-dir/new-file.txt",
		"README.md",
	}
	sort.Strings(expectedChangedFiles)

	require.NoError(t, err)
	assert.Equal(t, expectedChangedFiles, changedFiles)
}

func TestAddCommit(t *testing.T) {
	faker, err := newFaker()
	require.NoError(t, err)
	defer faker.clean()

	var (
		org      = "test-repo-org"
		repoName = "repo-add-commit"
		ctx      = context.Background()
	)

	err = faker.makeRepo(org, repoName)
	require.NoError(t, err)
	r := &repo{
		dir:     faker.repoDir(org, repoName),
		gitPath: faker.gitPath,
	}

	commits, err := r.ListCommits(ctx, "")
	require.NoError(t, err)
	require.Equal(t, 1, len(commits))

	path := filepath.Join(r.dir, "new-file.txt")
	err = os.WriteFile(path, []byte("content"), os.ModePerm)
	require.NoError(t, err)

	err = r.addCommit(ctx, "Added new file")
	require.NoError(t, err)

	err = r.addCommit(ctx, "No change")
	require.Equal(t, ErrNoChange, err)

	commits, err = r.ListCommits(ctx, "")
	require.NoError(t, err)
	require.Equal(t, 2, len(commits))
	assert.Equal(t, "Added new file", commits[0].Message)
}

func TestCommitChanges(t *testing.T) {
	faker, err := newFaker()
	require.NoError(t, err)
	defer faker.clean()

	var (
		org      = "test-repo-org"
		repoName = "repo-commit-changes"
		ctx      = context.Background()
	)

	err = faker.makeRepo(org, repoName)
	require.NoError(t, err)
	r := &repo{
		dir:     faker.repoDir(org, repoName),
		gitPath: faker.gitPath,
	}

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

	bytes, err := os.ReadFile(filepath.Join(r.dir, "README.md"))
	require.NoError(t, err)
	assert.Equal(t, string(changes["README.md"]), string(bytes))

	bytes, err = os.ReadFile(filepath.Join(r.dir, "a/b/c/new.txt"))
	require.NoError(t, err)
	assert.Equal(t, string(changes["a/b/c/new.txt"]), string(bytes))
}
