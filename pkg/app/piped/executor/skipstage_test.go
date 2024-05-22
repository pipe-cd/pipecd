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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSkipByCommitMessagePrefixes(t *testing.T) {

}

func TestHasOnlyPathsToSkip(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		name         string
		skipPatterns []string
		changedFiles []string
		expected     bool
	}{
		{
			name:         "no skip patterns",
			skipPatterns: nil,
			changedFiles: []string{"file1"},
			expected:     false,
		},
		{
			name:         "no changed files",
			skipPatterns: []string{"file1"},
			changedFiles: nil,
			expected:     true,
		},
		{
			name:         "no skip patterns and no changed files",
			skipPatterns: nil,
			changedFiles: nil,
			expected:     true,
		},
		{
			name:         "skip pattern matches all changed files",
			skipPatterns: []string{"file1", "file2"},
			changedFiles: []string{"file1", "file2"},
			expected:     true,
		},
		{
			name:         "skip pattern does not match changed files",
			skipPatterns: []string{"file1", "file2"},
			changedFiles: []string{"file1", "file3"},
			expected:     false,
		},
		{
			name:         "skip files of a directory",
			skipPatterns: []string{"dir1/*"},
			changedFiles: []string{"dir1/file1", "dir1/file2"},
			expected:     true,
		},
		{
			name:         "skip files recursively",
			skipPatterns: []string{"dir1/**"},
			changedFiles: []string{"dir1/file1", "dir1/sub/file2"},
			expected:     true,
		},
		{
			name:         "skip files not recursively",
			skipPatterns: []string{"dir1/*"},
			changedFiles: []string{"dir1/sub/file2"},
			expected:     false,
		},
		{
			name:         "skip files with the extension recursively",
			skipPatterns: []string{"dir1/**.yaml"},
			changedFiles: []string{"dir1/file1.yaml", "dir1/sub/file2.yaml"},
			expected:     true,
		},
		{
			name:         "skip files with the extension not recursively",
			skipPatterns: []string{"*.yaml"},
			changedFiles: []string{"dir1/file1.yaml"},
			expected:     false,
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// We do not use t.Parallel() here due to https://pkg.go.dev/github.com/pipe-cd/pipecd/pkg/filematcher#PatternMatcher.Matches.
			actual, err := hasOnlyPathsToSkip(tc.skipPatterns, tc.changedFiles)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
