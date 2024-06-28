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
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

	previousCommit, err := r.GetCommitForRev(ctx, "HEAD")
	require.NoError(t, err)
	require.NotEqual(t, "", previousCommit.Hash)

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

	headCommit, err := r.GetCommitForRev(ctx, "HEAD")
	require.NoError(t, err)
	require.NotEqual(t, "", headCommit.Hash)

	changedFiles, err := r.ChangedFiles(ctx, previousCommit.Hash, headCommit.Hash)
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

func Test_setGCAutoDetach(t *testing.T) {
	faker, err := newFaker()
	require.NoError(t, err)
	defer faker.clean()

	var (
		org      = "test-repo-org"
		repoName = "repo-set-gc-auto-detach"
		ctx      = context.Background()
	)

	err = faker.makeRepo(org, repoName)
	require.NoError(t, err)
	r := &repo{
		dir:     faker.repoDir(org, repoName),
		gitPath: faker.gitPath,
	}

	getGCAutoDetach := func(ctx context.Context, repo *repo) (bool, error) {
		cmd := exec.CommandContext(ctx, repo.gitPath, "config", "--get", "gc.autoDetach")
		cmd.Dir = r.dir
		out, err := cmd.CombinedOutput()
		if err != nil {
			return false, err
		}
		v, err := strconv.ParseBool(strings.TrimSuffix(string(out), "\n"))
		if err != nil {
			return false, err
		}

		return v, nil
	}

	// set  as true firstly, and then set as false.
	// set true
	err = r.setGCAutoDetach(ctx, true)
	require.NoError(t, err)

	got, err := getGCAutoDetach(ctx, r)
	if err != nil {
		t.Fatal(err)
	}
	require.NoError(t, err)

	assert.Equal(t, true, got)

	// set false
	err = r.setGCAutoDetach(ctx, false)
	require.NoError(t, err)

	got, err = getGCAutoDetach(ctx, r)
	require.NoError(t, err)

	assert.Equal(t, false, got)
}

func TestCopy(t *testing.T) {
	faker, err := newFaker()
	require.NoError(t, err)
	defer faker.clean()

	var (
		org      = "test-repo-org"
		repoName = "repo-copy"
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

	tmpDir := filepath.Join(faker.dir, "tmp-repo")
	newRepo, err := r.Copy(tmpDir)
	require.NoError(t, err)

	assert.NotEqual(t, r, newRepo)

	newRepoCommits, err := newRepo.ListCommits(ctx, "")
	require.NoError(t, err)
	assert.Equal(t, 1, len(newRepoCommits))

	assert.Equal(t, commits, newRepoCommits)
}

func TestGetCommitForRev(t *testing.T) {
	faker, err := newFaker()
	require.NoError(t, err)
	defer faker.clean()

	var (
		org      = "test-repo-org"
		repoName = "repo-get-commit-from-rev"
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

	commit, err := r.GetCommitForRev(ctx, "HEAD")
	require.NoError(t, err)
	assert.Equal(t, commits[0].Hash, commit.Hash)
}
