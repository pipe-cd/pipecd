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

package controller

import (
	"context"
	"strings"

	"github.com/pipe-cd/pipecd/pkg/app/piped/executor"
	"github.com/pipe-cd/pipecd/pkg/config"
	"github.com/pipe-cd/pipecd/pkg/filematcher"
	"github.com/pipe-cd/pipecd/pkg/git"
)

// checkSkipStage checks whether the stage should be skipped or not.
func checkSkipStage(ctx context.Context, in executor.Input, opt config.SkipOptions) (skip bool, err error) {
	if len(opt.Paths) == 0 && len(opt.CommitMessagePrefixes) == 0 {
		// When no condition is specified.
		return false, nil
	}

	repoCfg := in.Application.GitPath.Repo
	repo, err := in.GitClient.Clone(ctx, repoCfg.Id, repoCfg.Remote, repoCfg.Branch, "")
	if err != nil {
		return false, err
	}

	// Check by path pattern
	skip, err = skipByPathPattern(ctx, opt, repo, in.RunningDSP.Revision(), in.TargetDSP.Revision())
	if err != nil {
		return false, err
	}
	if skip {
		return true, nil
	}

	// Check by prefix of commit message
	skip, err = skipByCommitMessagePrefixes(ctx, opt, repo, in.TargetDSP.Revision())
	return skip, err
}

func skipByCommitMessagePrefixes(ctx context.Context, opt config.SkipOptions, repo git.Repo, targetRev string) (skip bool, err error) {
	if len(opt.CommitMessagePrefixes) == 0 {
		return false, nil
	}

	commit, err := repo.GetCommitForRev(ctx, targetRev)
	if err != nil {
		return false, err
	}

	return commitMessageHasAnyPrefix(commit.Message, opt.CommitMessagePrefixes), nil
}

func commitMessageHasAnyPrefix(commitMessage string, prefixes []string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(commitMessage, prefix) {
			return true
		}
	}
	return false
}

func skipByPathPattern(ctx context.Context, opt config.SkipOptions, repo git.Repo, runningRev, targetRev string) (skip bool, err error) {
	if len(opt.Paths) == 0 {
		return false, nil
	}

	changedFiles, err := repo.ChangedFiles(ctx, runningRev, targetRev)
	if err != nil {
		return false, err
	}
	return hasOnlyPathsToSkip(opt.Paths, changedFiles)
}

// hasOnlyPathsToSkip returns true if and only if all changed files are included in `skipPatterns`.
// If any changed file does not match all `skipPatterns`, it returns false.
func hasOnlyPathsToSkip(skipPatterns []string, changedFiles []string) (bool, error) {
	matcher, err := filematcher.NewPatternMatcher(skipPatterns)
	if err != nil {
		return false, err
	}

	for _, changedFile := range changedFiles {
		if !matcher.Matches(changedFile) {
			return false, nil
		}
	}

	return true, nil
}
