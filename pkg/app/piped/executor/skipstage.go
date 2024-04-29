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

package executor

import (
	"context"
	"strings"

	"github.com/pipe-cd/pipecd/pkg/config"
	"github.com/pipe-cd/pipecd/pkg/filematcher"
	"github.com/pipe-cd/pipecd/pkg/git"
)

// based on stage's config.
func checkSkipStage(ctx context.Context, in Input, opt config.SkipStageOptions) (skip bool, err error) {
	if opt.Paths == nil && len(opt.CommitMessagePrefixes) == 0 {
		// When no condition is specified for skipping.
		return false, nil
	}

	appRepo := in.Application.GitPath.Repo
	clonedRepo, err := in.GitClient.Clone(ctx, appRepo.Id, appRepo.Remote, appRepo.Branch, "")
	if err != nil {
		return false, err
	}

	// (1)と(2)はOR。どちらか一方でも満たせばスキップする。
	// (1)ファイルパスで判定する場合
	skip, err = skipByPathPattern(ctx, in, opt, clonedRepo)
	if err != nil {
		return false, err
	}
	if skip {
		return true, nil
	}

	// (2)Gitのコミットメッセージで判定する場合
	skip, err = skipByCommitMessagePrefixes(ctx, in, opt, clonedRepo)
	return skip, err
}

func skipByCommitMessagePrefixes(ctx context.Context, in Input, opt config.SkipStageOptions, repo git.Repo) (skip bool, err error) {
	if len(opt.CommitMessagePrefixes) > 0 {
		return false, nil
	}

	commit, err := repo.GetCommitFromRev(ctx, in.TargetDSP.Revision())
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

func skipByPathPattern(ctx context.Context, in Input, opt config.SkipStageOptions, repo git.Repo) (skip bool, err error) {
	if opt.Paths == nil {
		return false, nil
	}

	changedFiles, err := repo.ChangedFiles(ctx, in.RunningDSP.Revision(), in.TargetDSP.Revision())
	if err != nil {
		return false, err
	}
	// TODO: Which should we choose when no file is changed? (e.g. when SYNC command is executed from Console)
	//       hasOnlyPathToSkip() returns true in this case.
	// if len(changedFiles) == 0 {
	// 	return false, err
	// }

	return hasOnlyPathsToSkip(opt.Paths, changedFiles)
}

// hasOnlyPathsToSkip returns true if and only if all changed files are included in `skipPatterns`.
// If any changed file is not included in `skipPatterns`, it returns false.
// If `changedFiles` is empty, it returns true.
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
