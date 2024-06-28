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
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pipe-cd/pipecd/pkg/config"
	"github.com/pipe-cd/pipecd/pkg/git"
	"github.com/pipe-cd/pipecd/pkg/git/gittest"
	"github.com/stretchr/testify/assert"
)

func TestSkipByCommitMessagePrefixes(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		name          string
		commitMessage string
		prefixes      []string
		skip          bool
	}{
		{
			name:          "no prefixes",
			commitMessage: "test message",
			prefixes:      []string{},
			skip:          false,
		},
		{
			name:          "no commit message",
			commitMessage: "",
			prefixes:      []string{"to-skip"},
			skip:          false,
		},
		{
			name:          "prefix matches",
			commitMessage: "to-skip: test message",
			prefixes:      []string{"to-skip"},
			skip:          true,
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repoMock := gittest.NewMockRepo(ctrl)
			repoMock.EXPECT().GetCommitForRev(gomock.Any(), gomock.Any()).Return(git.Commit{
				Message: tc.commitMessage,
			}, nil).AnyTimes()

			opt := config.SkipOptions{
				CommitMessagePrefixes: tc.prefixes,
			}
			skip, err := skipByCommitMessagePrefixes(context.Background(), opt, repoMock, "test-rev")
			assert.Equal(t, tc.skip, skip)
			assert.NoError(t, err)
		})
	}
}

func TestSkipByPathPattern(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		name         string
		skipPatterns []string
		changedFiles []string
		skip         bool
	}{
		{
			name:         "no skip patterns",
			skipPatterns: nil,
			changedFiles: []string{"file1"},
			skip:         false,
		},
		{
			name:         "no changed files",
			skipPatterns: []string{"file1"},
			changedFiles: nil,
			skip:         true,
		},
		{
			name:         "no skip patterns and no changed files",
			skipPatterns: nil,
			changedFiles: nil,
			skip:         false,
		},
		{
			name:         "skip pattern matches all changed files",
			skipPatterns: []string{"file1", "file2"},
			changedFiles: []string{"file1", "file2"},
			skip:         true,
		},
		{
			name:         "skip pattern does not match changed files",
			skipPatterns: []string{"file1", "file2"},
			changedFiles: []string{"file1", "file3"},
			skip:         false,
		},
		{
			name:         "skip files of a directory",
			skipPatterns: []string{"dir1/*"},
			changedFiles: []string{"dir1/file1", "dir1/file2"},
			skip:         true,
		},
		{
			name:         "skip files recursively",
			skipPatterns: []string{"dir1/**"},
			changedFiles: []string{"dir1/file1", "dir1/sub/file2"},
			skip:         true,
		},
		{
			name:         "skip files with the extension recursively",
			skipPatterns: []string{"dir1/**/*.yaml"},
			changedFiles: []string{"dir1/file1.yaml", "dir1/sub1/file2.yaml", "dir1/sub1/sub2/file3.yaml"},
			skip:         true,
		},
		{
			name:         "skip files not recursively",
			skipPatterns: []string{"*.yaml"},
			changedFiles: []string{"dir1/file1.yaml"},
			skip:         false,
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// We do not use t.Parallel() here due to https://pkg.go.dev/github.com/pipe-cd/pipecd/pkg/filematcher#PatternMatcher.Matches.
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repoMock := gittest.NewMockRepo(ctrl)
			repoMock.EXPECT().ChangedFiles(gomock.Any(), gomock.Any(), gomock.Any()).Return(tc.changedFiles, nil).AnyTimes()

			opt := config.SkipOptions{
				Paths: tc.skipPatterns,
			}
			actual, err := skipByPathPattern(context.Background(), opt, repoMock, "running-rev", "target-rev")
			assert.NoError(t, err)
			assert.Equal(t, tc.skip, actual)
		})
	}
}
