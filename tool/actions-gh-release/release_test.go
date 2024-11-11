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

package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/go-github/v39/github"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseReleaseConfig(t *testing.T) {
	fakeShowCommitter := true
	testcases := []struct {
		name        string
		configFile  string
		expected    *ReleaseConfig
		expectedErr error
	}{
		{
			name:        "empty config",
			configFile:  "testdata/empty-config.txt",
			expectedErr: fmt.Errorf("tag must be specified"),
		},
		{
			name:       "valid config",
			configFile: "testdata/valid-config.txt",
			expected: &ReleaseConfig{
				Tag:  "v1.1.0",
				Name: "hello",
				CommitInclude: ReleaseCommitMatcherConfig{
					Contains: []string{
						"app/hello",
					},
				},
				CommitExclude: ReleaseCommitMatcherConfig{
					Prefixes: []string{
						"Merge pull request #",
					},
				},
				CommitCategories: []ReleaseCommitCategoryConfig{
					ReleaseCommitCategoryConfig{
						ID:    "_category_0",
						Title: "Breaking Changes",
						ReleaseCommitMatcherConfig: ReleaseCommitMatcherConfig{
							Contains: []string{"change-category/breaking-change"},
						},
					},
					ReleaseCommitCategoryConfig{
						ID:    "_category_1",
						Title: "New Features",
						ReleaseCommitMatcherConfig: ReleaseCommitMatcherConfig{
							Contains: []string{"change-category/new-feature"},
						},
					},
					ReleaseCommitCategoryConfig{
						ID:    "_category_2",
						Title: "Notable Changes",
						ReleaseCommitMatcherConfig: ReleaseCommitMatcherConfig{
							Contains: []string{"change-category/notable-change"},
						},
					},
					ReleaseCommitCategoryConfig{
						ID:                         "_category_3",
						Title:                      "Internal Changes",
						ReleaseCommitMatcherConfig: ReleaseCommitMatcherConfig{},
					},
				},
				ReleaseNoteGenerator: ReleaseNoteGeneratorConfig{
					ShowCommitter:       &fakeShowCommitter,
					UseReleaseNoteBlock: true,
				},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			data, err := testdata.ReadFile(tc.configFile)
			require.NoError(t, err)

			cfg, err := parseReleaseConfig(data)
			assert.Equal(t, tc.expected, cfg)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestBuildReleaseCommits(t *testing.T) {
	ctx := context.Background()
	fakeShowCommitter := true
	config := ReleaseConfig{
		Tag:  "v1.1.0",
		Name: "hello",
		CommitInclude: ReleaseCommitMatcherConfig{
			Contains: []string{
				"app/hello",
			},
		},
		CommitExclude: ReleaseCommitMatcherConfig{
			Prefixes: []string{
				"Merge pull request #",
			},
		},
		CommitCategories: []ReleaseCommitCategoryConfig{
			ReleaseCommitCategoryConfig{
				ID:    "breaking-change",
				Title: "Breaking Changes",
				ReleaseCommitMatcherConfig: ReleaseCommitMatcherConfig{
					Contains: []string{"change-category/breaking-change"},
				},
			},
			ReleaseCommitCategoryConfig{
				ID:    "new-feature",
				Title: "New Features",
				ReleaseCommitMatcherConfig: ReleaseCommitMatcherConfig{
					Contains: []string{"change-category/new-feature"},
				},
			},
			ReleaseCommitCategoryConfig{
				ID:    "notable-change",
				Title: "Notable Changes",
				ReleaseCommitMatcherConfig: ReleaseCommitMatcherConfig{
					Contains: []string{"change-category/notable-change"},
				},
			},
			ReleaseCommitCategoryConfig{
				ID:                         "internal-change",
				Title:                      "Internal Changes",
				ReleaseCommitMatcherConfig: ReleaseCommitMatcherConfig{},
			},
		},
		ReleaseNoteGenerator: ReleaseNoteGeneratorConfig{
			ShowCommitter:       &fakeShowCommitter,
			UseReleaseNoteBlock: true,
		},
	}

	testcases := []struct {
		name     string
		commits  []Commit
		config   ReleaseConfig
		expected []ReleaseCommit
		wantErr  bool
	}{
		{
			name:     "empty",
			expected: []ReleaseCommit{},
			wantErr:  false,
		},
		{
			name: "ok",
			commits: []Commit{
				Commit{
					Subject: "Commit 1 message",
					Body:    "commit 1\napp/hello\n- change-category/breaking-change",
				},
				Commit{
					Subject: "Commit 2 message",
					Body:    "commit 2\napp/hello",
				},
				Commit{
					Subject: "Commit 3 message",
					Body:    "commit 3\napp/hello\n- change-category/notable-change",
				},
				Commit{
					Subject: "Commit 4 message",
					Body:    "commit 4\napp/hello\n```release-note\nCommit 4 release note\n```\n- change-category/notable-change\n",
				},
				Commit{
					Subject: "Commit 5 message",
					Body:    "commit 5",
				},
			},
			config: config,
			expected: []ReleaseCommit{
				ReleaseCommit{
					Commit: Commit{
						Subject: "Commit 1 message",
						Body:    "commit 1\napp/hello\n- change-category/breaking-change",
					},
					CategoryName: "breaking-change",
					ReleaseNote:  "Commit 1 message",
				},
				ReleaseCommit{
					Commit: Commit{
						Subject: "Commit 2 message",
						Body:    "commit 2\napp/hello",
					},
					CategoryName: "internal-change",
					ReleaseNote:  "Commit 2 message",
				},
				ReleaseCommit{
					Commit: Commit{
						Subject: "Commit 3 message",
						Body:    "commit 3\napp/hello\n- change-category/notable-change",
					},
					CategoryName: "notable-change",
					ReleaseNote:  "Commit 3 message",
				},
				ReleaseCommit{
					Commit: Commit{
						Subject: "Commit 4 message",
						Body:    "commit 4\napp/hello\n```release-note\nCommit 4 release note\n```\n- change-category/notable-change\n",
					},
					CategoryName: "notable-change",
					ReleaseNote:  "Commit 4 release note",
				},
			},
			wantErr: false,
		},
		{
			name: "Add include condition: parent of merge commit",
			commits: []Commit{
				{
					Hash:         "a",
					ParentHashes: []string{"z"},
					Subject:      "Commit 1 message",
					Body:         "commit 1",
				},
				{
					Hash:         "b",
					ParentHashes: []string{"a"},
					Subject:      "Commit 2 message",
					Body:         "commit 2",
				},
				{
					Hash:         "c",
					ParentHashes: []string{"z", "b"},
					Subject:      "Commit 3 message",
					Body:         "commit 3\napp/hello\n- change-category/notable-change",
				},
				{
					Hash:         "d",
					ParentHashes: []string{"c"},
					Subject:      "Commit 4 message",
					Body:         "commit 4",
				},
				{
					Hash:         "e",
					ParentHashes: []string{"c", "d"},
					Subject:      "Commit 5 message",
					Body:         "commit 5",
				},
			},
			config: func(base ReleaseConfig) ReleaseConfig {
				base.CommitInclude.ParentOfMergeCommit = true
				return base
			}(config),
			expected: []ReleaseCommit{
				{
					Commit: Commit{
						Hash:         "a",
						ParentHashes: []string{"z"},
						Subject:      "Commit 1 message",
						Body:         "commit 1",
					},
					CategoryName: "internal-change",
					ReleaseNote:  "Commit 1 message",
				},
				{
					Commit: Commit{
						Hash:         "b",
						ParentHashes: []string{"a"},
						Subject:      "Commit 2 message",
						Body:         "commit 2",
					},
					CategoryName: "internal-change",
					ReleaseNote:  "Commit 2 message",
				},
				{
					Commit: Commit{
						Hash:         "c",
						ParentHashes: []string{"z", "b"},
						Subject:      "Commit 3 message",
						Body:         "commit 3\napp/hello\n- change-category/notable-change",
					},
					CategoryName: "notable-change",
					ReleaseNote:  "Commit 3 message",
				},
			},
			wantErr: false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := buildReleaseCommits(ctx, nil, tc.commits, tc.config, nil)
			assert.Equal(t, tc.wantErr, err != nil)
			assert.Equal(t, tc.expected, got)
		})
	}
}

func TestBuildReleaseCommitsForCategoryLabel(t *testing.T) {
	t.Parallel()

	cfgs := []ReleaseCommitCategoryConfig{
		{
			ID:    "breaking-change",
			Title: "Breaking Changes",
			ReleaseCommitMatcherConfig: ReleaseCommitMatcherConfig{
				Labels: []string{"change-category/breaking-change"},
			},
		},
		{
			ID:    "notable-change",
			Title: "Notable Changes",
			ReleaseCommitMatcherConfig: ReleaseCommitMatcherConfig{
				Labels: []string{"change-category/notable-change"},
			},
		},
		{
			ID:                         "internal-change",
			Title:                      "Internal Changes",
			ReleaseCommitMatcherConfig: ReleaseCommitMatcherConfig{},
		},
	}

	testcases := []struct {
		title            string
		labels           []*github.Label
		categories       []ReleaseCommitCategoryConfig
		expectedCategory string
	}{
		{
			title:            "no categories config provided",
			labels:           []*github.Label{},
			categories:       []ReleaseCommitCategoryConfig{},
			expectedCategory: "",
		},
		{
			title:            "breaking-change",
			labels:           []*github.Label{{Name: github.String("change-category/breaking-change")}},
			categories:       cfgs,
			expectedCategory: "breaking-change",
		},
		{
			title:            "not match any category",
			labels:           []*github.Label{{Name: github.String("foo")}},
			categories:       cfgs,
			expectedCategory: "",
		},
		{
			title: "matching multiple labels results in the first one defined in config",
			labels: []*github.Label{
				{Name: github.String("change-category/notable-change")},
				{Name: github.String("change-category/breaking-change")},
			},
			categories:       cfgs,
			expectedCategory: "breaking-change",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.title, func(t *testing.T) {
			t.Parallel()

			pr := &github.PullRequest{
				Labels: tc.labels,
			}
			got := determineCommitCategoryOfPR(pr, tc.categories)
			assert.Equal(t, tc.expectedCategory, got)
		})
	}
}

func TestRenderReleaseNote(t *testing.T) {
	testcases := []struct {
		name     string
		proposal ReleaseProposal
		config   ReleaseConfig
		expected string
	}{
		{
			name: "no category",
			proposal: ReleaseProposal{
				Tag:    "v0.2.0",
				PreTag: "v0.1.0",
				Commits: []ReleaseCommit{
					ReleaseCommit{
						Commit: Commit{
							Subject: "Commit 1 message",
							Body:    "commit 1\n- change-category/breaking-change",
						},
						ReleaseNote: "Commit 1 message",
					},
					ReleaseCommit{
						Commit: Commit{
							Subject: "Commit 2 message",
							Body:    "commit 2",
						},
						ReleaseNote: "Commit 2 message",
					},
					ReleaseCommit{
						Commit: Commit{
							Subject: "Commit 3 message",
							Body:    "commit 3\n- change-category/notable-change",
						},
						ReleaseNote: "Commit 3 message",
					},
					ReleaseCommit{
						Commit: Commit{
							Subject: "Commit 4 message",
							Body:    "commit 4\n```release-note\nCommit 4 release note\n```\n- change-category/notable-change",
						},
						ReleaseNote: "Commit 4 release note",
					},
				},
			},
			config:   ReleaseConfig{},
			expected: "testdata/no-category-release-note.txt",
		},
		{
			name: "has category",
			proposal: ReleaseProposal{
				Tag:    "v0.2.0",
				PreTag: "v0.1.0",
				Commits: []ReleaseCommit{
					ReleaseCommit{
						Commit: Commit{
							Subject: "Commit 1 message",
							Body:    "commit 1\n- change-category/breaking-change",
						},
						CategoryName: "breaking-change",
						ReleaseNote:  "Commit 1 message",
					},
					ReleaseCommit{
						Commit: Commit{
							Subject: "Commit 2 message",
							Body:    "commit 2",
						},
						CategoryName: "internal-change",
						ReleaseNote:  "Commit 2 message",
					},
					ReleaseCommit{
						Commit: Commit{
							Subject: "Commit 3 message",
							Body:    "commit 3\n- change-category/notable-change",
						},
						CategoryName: "notable-change",
						ReleaseNote:  "Commit 3 message",
					},
					ReleaseCommit{
						Commit: Commit{
							Subject: "Commit 4 message",
							Body:    "commit 4\n```release-note\nCommit 4 release note\n```\n- change-category/notable-change",
						},
						CategoryName: "notable-change",
						ReleaseNote:  "Commit 4 release note",
					},
				},
			},
			config: ReleaseConfig{
				CommitCategories: []ReleaseCommitCategoryConfig{
					ReleaseCommitCategoryConfig{
						ID:    "breaking-change",
						Title: "Breaking Changes",
					},
					ReleaseCommitCategoryConfig{
						ID:    "new-feature",
						Title: "New Features",
					},
					ReleaseCommitCategoryConfig{
						ID:    "notable-change",
						Title: "Notable Changes",
					},
					ReleaseCommitCategoryConfig{
						ID:    "internal-change",
						Title: "Internal Changes",
					},
				},
			},
			expected: "testdata/has-category-release-note.txt",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := renderReleaseNote(tc.proposal, tc.config)

			expected, err := testdata.ReadFile(tc.expected)
			require.NoError(t, err)

			assert.Equal(t, string(expected), string(got))
		})
	}
}
