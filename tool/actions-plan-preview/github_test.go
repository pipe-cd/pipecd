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
	"testing"
	"time"

	"github.com/google/go-github/v36/github"
	"github.com/shurcooL/githubv4"
	"github.com/stretchr/testify/assert"
)

type dummyPullRequestsService struct {
	mergeable  bool
	createdAt  time.Time
	headBranch string
	headCommit string
	baseBranch string
}

func (d dummyPullRequestsService) Get(_ context.Context, _ string, _ string, _ int) (*github.PullRequest, *github.Response, error) {
	return &github.PullRequest{
		Mergeable: &d.mergeable,
		CreatedAt: &d.createdAt,
		Head: &github.PullRequestBranch{
			Ref: &d.headBranch,
			SHA: &d.headCommit,
		},
		Base: &github.PullRequestBranch{Ref: &d.baseBranch},
	}, nil, nil
}

func TestParseGithubEvent(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		name        string
		eventName   string
		payload     []byte
		argPRNum    int
		prService   dummyPullRequestsService
		expected    *githubEvent
		expectedErr error
	}{
		{
			name:      "successfully parsed PR event",
			eventName: "pull_request",
			payload:   readTestdataFile(t, "testdata/pull-request-payload.json"),
			expected: &githubEvent{
				Owner:       "Codertocat",
				Repo:        "Hello-World",
				RepoRemote:  "git@github.com:Codertocat/Hello-World.git",
				PRNumber:    2,
				PRMergeable: nil,
				HeadBranch:  "changes",
				HeadCommit:  "ec26c3e57ca3a959ca5aad62de7213c562f8c821",
				BaseBranch:  "master",
				SenderLogin: "Codertocat",
			},
		},
		{
			name:      "successfully parsed issue commit event",
			eventName: "issue_comment",
			payload:   readTestdataFile(t, "testdata/issue-comment-payload.json"),
			prService: dummyPullRequestsService{
				mergeable:  true,
				createdAt:  time.Unix(0, 0),
				headBranch: "head-branch",
				headCommit: "head-commit",
				baseBranch: "base-branch",
			},
			expected: &githubEvent{
				Owner:       "Codertocat",
				Repo:        "Hello-World",
				RepoRemote:  "git@github.com:Codertocat/Hello-World.git",
				PRNumber:    1,
				PRMergeable: boolPointer(true),
				PRClosed:    false,
				HeadBranch:  "head-branch",
				HeadCommit:  "head-commit",
				BaseBranch:  "base-branch",
				SenderLogin: "Codertocat",
				IsComment:   true,
				CommentURL:  "https://github.com/Codertocat/Hello-World/issues/1#issuecomment-492700400",
			},
		},
		{
			name:      "successfully parsed push event",
			eventName: "push",
			payload:   readTestdataFile(t, "testdata/push-event-payload.json"),
			argPRNum:  1,
			prService: dummyPullRequestsService{
				mergeable:  true,
				createdAt:  time.Unix(0, 0),
				headBranch: "head-branch",
				headCommit: "head-commit",
				baseBranch: "base-branch",
			},
			expected: &githubEvent{
				Owner:       "Codertocat",
				Repo:        "Hello-World",
				RepoRemote:  "git@github.com:Codertocat/Hello-World.git",
				PRNumber:    1,
				PRMergeable: boolPointer(true),
				PRClosed:    false,
				HeadBranch:  "head-branch",
				HeadCommit:  "head-commit",
				BaseBranch:  "base-branch",
				SenderLogin: "Codertocat",
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, err := parseGitHubEvent(context.Background(), tc.prService, tc.eventName, tc.payload, tc.argPRNum)
			assert.Equal(t, tc.expected, got)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestFilterLatestPlanPreviewComment(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		name     string
		comments []issueCommentQuery
		key      string
		expected *issueCommentQuery
	}{
		{
			name:     "no comments",
			comments: []issueCommentQuery{},
			key:      "",
			expected: nil,
		},
		{
			name: "no plan preview comment",
			comments: []issueCommentQuery{
				{Body: githubv4.String("no-planpreview-comment")},
			},
			key:      "",
			expected: nil,
		},
		{
			name: "latest comment with no key is returned",
			comments: []issueCommentQuery{
				{Body: githubv4.String("no-planpreview-comment")},
				{Body: githubv4.String("<!-- pipecd-plan-preview -->"), ID: 1},
				{Body: githubv4.String("<!-- pipecd-plan-preview -->"), ID: 2},
				{Body: githubv4.String("no-planpreview-comment")},
			},
			key: "",
			expected: &issueCommentQuery{
				Body: githubv4.String("<!-- pipecd-plan-preview -->"),
				ID:   2,
			},
		},
		{
			name: "multiple keys",
			comments: []issueCommentQuery{
				{Body: githubv4.String("<!-- pipecd-plan-preview foo -->")},
				{Body: githubv4.String("<!-- pipecd-plan-preview fcde2b2edba56bf408601fb721fe9b5c338d10ee429ea04fae5511b68fbf8fb9-->")}, // hashed 'bar'
				{Body: githubv4.String("<!-- pipecd-plan-preview baz -->")},
			},
			key:      "bar",
			expected: &issueCommentQuery{Body: githubv4.String("<!-- pipecd-plan-preview fcde2b2edba56bf408601fb721fe9b5c338d10ee429ea04fae5511b68fbf8fb9-->")},
		},
		{
			name: "multibyte key",
			comments: []issueCommentQuery{
				{Body: githubv4.String("<!-- pipecd-plan-preview foo -->")},
				{Body: githubv4.String("<!-- pipecd-plan-preview 1bef6bca1c45e2e0b482c46e0ba2c7b1bc711ab8aea17cbd4af275f02e651982-->")}, // hashed 'αβ'
				{Body: githubv4.String("<!-- pipecd-plan-preview baz -->")},
			},
			key:      "αβ",
			expected: &issueCommentQuery{Body: githubv4.String("<!-- pipecd-plan-preview 1bef6bca1c45e2e0b482c46e0ba2c7b1bc711ab8aea17cbd4af275f02e651982-->")},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := filterLatestPlanPreviewComment(tc.comments, tc.key)
			assert.Equal(t, tc.expected, got)
		})
	}
}
