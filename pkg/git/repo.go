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
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	ErrNoChange       = errors.New("no change")
	ErrBranchNotFresh = errors.New("some refs were not updated")
)

// Repo provides functions to get and handle git data.
type Repo interface {
	GetPath() string
	GetClonedBranch() string
	Copy(dest string) (Repo, error)

	ListCommits(ctx context.Context, visionRange string) ([]Commit, error)
	GetLatestCommit(ctx context.Context) (Commit, error)
	GetCommitHashForRev(ctx context.Context, rev string) (string, error)
	ChangedFiles(ctx context.Context, from, to string) ([]string, error)
	Checkout(ctx context.Context, commitish string) error
	CheckoutPullRequest(ctx context.Context, number int, branch string) error
	Clean() error

	Pull(ctx context.Context, branch string) error
	MergeRemoteBranch(ctx context.Context, branch, commit, mergeCommitMessage string) error
	Push(ctx context.Context, branch string) error
	CommitChanges(ctx context.Context, branch, message string, newBranch bool, changes map[string][]byte) error
}

type repo struct {
	dir          string
	gitPath      string
	remote       string
	clonedBranch string
	gitEnvs      []string
}

// NewRepo creates a new Repo instance.
func NewRepo(dir, gitPath, remote, clonedBranch string, gitEnvs []string) *repo {
	return &repo{
		dir:          dir,
		gitPath:      gitPath,
		remote:       remote,
		clonedBranch: clonedBranch,
		gitEnvs:      gitEnvs,
	}
}

// GetPath returns the path to the local git directory.
func (r *repo) GetPath() string {
	return r.dir
}

// GetClonedBranch returns the name of cloned branch.
func (r *repo) GetClonedBranch() string {
	return r.clonedBranch
}

// Copy does copying the repository to the given destination.
// NOTE: the given “dest” must be a path that doesn’t exist yet.
// If you don't, it copies the repo root itself to the given dest as a subdirectory.
func (r *repo) Copy(dest string) (Repo, error) {
	cmd := exec.Command("cp", "-rf", r.dir, dest)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, formatCommandError(err, out)
	}

	return &repo{
		dir:          dest,
		gitPath:      r.gitPath,
		remote:       r.remote,
		clonedBranch: r.clonedBranch,
	}, nil
}

// ListCommits returns a list of commits in a given revision range.
func (r *repo) ListCommits(ctx context.Context, revisionRange string) ([]Commit, error) {
	args := []string{
		"log",
		"--no-decorate",
		fmt.Sprintf("--pretty=format:%s", commitLogFormat),
	}
	if revisionRange != "" {
		args = append(args, revisionRange)
	}

	out, err := r.runGitCommand(ctx, args...)
	if err != nil {
		return nil, formatCommandError(err, out)
	}

	return parseCommits(string(out))
}

// GetLatestCommit returns the most recent commit of current branch.
func (r *repo) GetLatestCommit(ctx context.Context) (Commit, error) {
	commits, err := r.ListCommits(ctx, "-1")
	if err != nil {
		return Commit{}, err
	}

	if len(commits) != 1 {
		return Commit{}, fmt.Errorf("commits must contain one item, got: %d", len(commits))
	}

	return commits[0], nil
}

// GetCommitHashForRev returns the hash value of the commit for a given rev.
func (r *repo) GetCommitHashForRev(ctx context.Context, rev string) (string, error) {
	out, err := r.runGitCommand(ctx, "rev-parse", rev)
	if err != nil {
		return "", formatCommandError(err, out)
	}

	return strings.TrimSpace(string(out)), nil
}

// ChangedFiles returns a list of files those were touched between two commits.
func (r *repo) ChangedFiles(ctx context.Context, from, to string) ([]string, error) {
	out, err := r.runGitCommand(ctx, "diff", "--name-only", from, to)
	if err != nil {
		return nil, formatCommandError(err, out)
	}

	var (
		lines = strings.Split(string(out), "\n")
		files = make([]string, 0, len(lines))
	)
	// The result may include some empty lines
	// so we need to remove all of them.
	for _, f := range lines {
		if f != "" {
			files = append(files, f)
		}
	}
	return files, nil
}

// Checkout checkouts to a given commitish.
func (r *repo) Checkout(ctx context.Context, commitish string) error {
	out, err := r.runGitCommand(ctx, "checkout", commitish)
	if err != nil {
		return formatCommandError(err, out)
	}
	return nil
}

// CheckoutPullRequest checkouts to the latest commit of a given pull request.
func (r *repo) CheckoutPullRequest(ctx context.Context, number int, branch string) error {
	target := fmt.Sprintf("pull/%d/head:%s", number, branch)
	out, err := r.runGitCommand(ctx, "fetch", r.remote, target)
	if err != nil {
		return formatCommandError(err, out)
	}
	return r.Checkout(ctx, branch)
}

// Pull fetches from and integrate with a local branch.
func (r *repo) Pull(ctx context.Context, branch string) error {
	out, err := r.runGitCommand(ctx, "pull", r.remote, branch)
	if err != nil {
		return formatCommandError(err, out)
	}
	return nil
}

// MergeRemoteBranch merges all commits until the given one
// from a remote branch to current local branch.
// This always adds a new merge commit into tree.
func (r *repo) MergeRemoteBranch(ctx context.Context, branch, commit, mergeCommitMessage string) error {
	out, err := r.runGitCommand(ctx, "fetch", r.remote, branch)
	if err != nil {
		return formatCommandError(err, out)
	}
	out, err = r.runGitCommand(ctx, "merge", "-q", "--no-ff", "-m", mergeCommitMessage, commit)
	if err != nil {
		return formatCommandError(err, out)
	}
	return nil
}

// Push pushes local changes of a given branch to the remote.
func (r *repo) Push(ctx context.Context, branch string) error {
	out, err := r.runGitCommand(ctx, "push", r.remote, branch)
	if err == nil {
		return nil
	}
	if strings.Contains(string(out), "failed to push some refs to") {
		return ErrBranchNotFresh
	}
	return formatCommandError(err, out)
}

// CommitChanges commits some changes into a branch.
func (r *repo) CommitChanges(ctx context.Context, branch, message string, newBranch bool, changes map[string][]byte) error {
	if newBranch {
		if err := r.checkoutNewBranch(ctx, branch); err != nil {
			return fmt.Errorf("failed to checkout new branch, branch: %v, error: %v", branch, err)
		}
	} else {
		if err := r.Checkout(ctx, branch); err != nil {
			return fmt.Errorf("failed to checkout branch, branch: %v, error: %v", branch, err)
		}
	}
	// Apply the changes.
	for p, bytes := range changes {
		filePath := filepath.Join(r.dir, p)
		dirPath := filepath.Dir(filePath)
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
				return fmt.Errorf("failed to create directory, dir: %s, err: %v", dirPath, err)
			}
		}
		if err := os.WriteFile(filePath, bytes, os.ModePerm); err != nil {
			return fmt.Errorf("failed to write file, file: %s, error: %v", filePath, err)
		}
	}
	// Commit the changes.
	if err := r.addCommit(ctx, message); err != nil {
		return fmt.Errorf("failed to commit, branch: %s, error: %v", branch, err)
	}
	return nil
}

// Clean deletes all local git data.
func (r repo) Clean() error {
	return os.RemoveAll(r.dir)
}

func (r *repo) checkoutNewBranch(ctx context.Context, branch string) error {
	out, err := r.runGitCommand(ctx, "checkout", "-b", branch)
	if err != nil {
		return formatCommandError(err, out)
	}
	return nil
}

func (r repo) addCommit(ctx context.Context, message string) error {
	out, err := r.runGitCommand(ctx, "add", ".")
	if err != nil {
		return formatCommandError(err, out)
	}
	out, err = r.runGitCommand(ctx, "commit", "-m", message)
	if err != nil {
		msg := string(out)
		if strings.Contains(msg, "nothing to commit, working tree clean") {
			return ErrNoChange
		}
		return formatCommandError(err, out)
	}
	return nil
}

// setUser configures username and email for local user of this repo.
func (r *repo) setUser(ctx context.Context, username, email string) error {
	if out, err := r.runGitCommand(ctx, "config", "user.name", username); err != nil {
		return formatCommandError(err, out)
	}
	if out, err := r.runGitCommand(ctx, "config", "user.email", email); err != nil {
		return formatCommandError(err, out)
	}
	return nil
}

func (r *repo) setRemote(ctx context.Context, remote string) error {
	out, err := r.runGitCommand(ctx, "remote", "set-url", "origin", remote)
	if err != nil {
		return formatCommandError(err, out)
	}
	return nil
}

func (r *repo) runGitCommand(ctx context.Context, args ...string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, r.gitPath, args...)
	cmd.Dir = r.dir
	cmd.Env = append(os.Environ(), r.gitEnvs...)
	return cmd.CombinedOutput()
}

func formatCommandError(err error, out []byte) error {
	return fmt.Errorf("err: %w, out: %s", err, string(out))
}
