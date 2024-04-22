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
)

// based on stage's config.
func checkSkipStage(ctx context.Context, in Input, opt config.SkipStageOptions) (skip bool, err error) {
	if opt.Paths == nil && opt.CommitMessagePrefix == "" {
		// When no condition is specified for skipping.
		return false, nil
	}

	repo := in.Application.GitPath.Repo
	clonedRepo, err := in.GitClient.Clone(ctx, repo.Id, repo.Remote, repo.Branch, "")
	if err != nil {
		return false, err
	}

	// (1)と(2)はOR。どちらか一方でも満たせばスキップする。
	// (1)ファイルパスで判定する場合
	if opt.Paths != nil {
		changedFiles, err := clonedRepo.ChangedFiles(ctx, in.RunningDSP.Revision(), in.TargetDSP.Revision())
		if err != nil {
			return false, err
		}

		// check whether changed files are included in opt.Paths.
		// if any file is included, return true.
		if opt.Paths != nil {
			for _, path := range changedFiles {
				// TODO use regex
				if opt.Paths.Match(path) {
					return true, nil
				}
			}
		}
	}

	// (2)Gitのコミットメッセージで判定する場合
	if opt.CommitMessagePrefix != "" {
		commit, err := clonedRepo.GetCommitFromHash(ctx, in.TargetDSP.Revision())
		if err != nil {
			return false, err
		}

		if strings.HasPrefix(commit.Message, opt.CommitMessagePrefix) {
			return true, nil
		}
	}

	return false, nil
}
