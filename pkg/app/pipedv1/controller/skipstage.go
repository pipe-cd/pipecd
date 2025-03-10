// Copyright 2025 The PipeCD Authors.
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

	config "github.com/pipe-cd/pipecd/pkg/configv1"
	"github.com/pipe-cd/pipecd/pkg/filematcher"
	"github.com/pipe-cd/pipecd/pkg/git"
)

// determineSkipStage checks whether the stage should be skipped or not.
func (s *scheduler) determineSkipStage(ctx context.Context, skipOpts config.SkipOptions) (skip bool, err error) {
	if len(skipOpts.Paths) == 0 && len(skipOpts.CommitMessagePrefixes) == 0 {
		// When no condition is specified.
		return false, nil
	}

	// TODO: Do not use clone here in order to avoid unnecessary cloning.
	repoCfg := s.deployment.GetGitPath().Repo
	repo, err := s.gitClient.Clone(ctx, repoCfg.Id, repoCfg.Remote, repoCfg.Branch, "")
	if err != nil {
		return false, err
	}

	// Check by path pattern
	skip, err = skipByPathPattern(ctx, skipOpts, repo, s.runningDSP.Revision(), s.targetDSP.Revision())
	if err != nil {
		return false, err
	}
	if skip {
		return true, nil
	}

	// Check by prefix of commit message
	skip, err = skipByCommitMessagePrefixes(ctx, skipOpts, repo, s.targetDSP.Revision())
	return skip, err
}

// skipByCommitMessagePrefixes returns true if the commit message has ANY one of the specified prefixes.
func skipByCommitMessagePrefixes(ctx context.Context, opt config.SkipOptions, repo git.Repo, targetRev string) (skip bool, err error) {
	if len(opt.CommitMessagePrefixes) == 0 {
		return false, nil
	}

	commit, err := repo.GetCommitForRev(ctx, targetRev)
	if err != nil {
		return false, err
	}

	for _, prefix := range opt.CommitMessagePrefixes {
		if strings.HasPrefix(commit.Message, prefix) {
			return true, nil
		}
	}
	return false, nil
}

// skipByPathPattern returns true if and only if ALL changed files are included in `opt.Paths`.
// If ANY changed file does not match all `skipPatterns`, it returns false.
func skipByPathPattern(ctx context.Context, opt config.SkipOptions, repo git.Repo, runningRev, targetRev string) (skip bool, err error) {
	if len(opt.Paths) == 0 {
		return false, nil
	}

	changedFiles, err := repo.ChangedFiles(ctx, runningRev, targetRev)
	if err != nil {
		return false, err
	}

	matcher, err := filematcher.NewPatternMatcher(opt.Paths)
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
