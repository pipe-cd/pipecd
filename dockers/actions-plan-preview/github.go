// Copyright 2021 The PipeCD Authors.
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
	"os"

	"github.com/google/go-github/v36/github"
)

type githubEvent struct {
	Owner       string
	Repo        string
	RepoRemote  string
	PRNumber    int
	PRMergeable *bool
	PRClosed    bool
	HeadBranch  string
	HeadCommit  string
	BaseBranch  string
	SenderLogin string
	IsComment   bool
	CommentURL  string
}

// parsePullRequestEvent uses the given environment variables
// to parse and build githubEvent struct.
// Currently, we support 2 kinds of event as below:
// - PullRequestEvent
//   https://pkg.go.dev/github.com/google/go-github/v36/github#PullRequestEvent
// - IssueCommentEvent
//   https://pkg.go.dev/github.com/google/go-github/v36/github#IssueCommentEvent
func parseGitHubEvent(ctx context.Context, client *github.Client) (*githubEvent, error) {
	const (
		pullRequestEventName = "pull_request"
		commentEventName     = "issue_comment"
	)

	eventName := os.Getenv("GITHUB_EVENT_NAME")
	if eventName != pullRequestEventName && eventName != commentEventName {
		return nil, fmt.Errorf("unexpected event %s, only %q and %q event are supported", eventName, pullRequestEventName, commentEventName)
	}

	eventPath := os.Getenv("GITHUB_EVENT_PATH")
	payload, err := os.ReadFile(eventPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read event payload: %v", err)
	}

	event, err := github.ParseWebHook(eventName, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to parse event payload: %v", err)
	}

	switch e := event.(type) {
	case *github.PullRequestEvent:
		return &githubEvent{
			Owner:       e.Repo.Owner.GetLogin(),
			Repo:        e.Repo.GetName(),
			RepoRemote:  e.Repo.GetSSHURL(),
			PRNumber:    e.GetNumber(),
			PRMergeable: e.PullRequest.Mergeable,
			PRClosed:    !e.PullRequest.GetClosedAt().IsZero(),
			HeadBranch:  e.PullRequest.Head.GetRef(),
			HeadCommit:  e.PullRequest.Head.GetSHA(),
			BaseBranch:  e.PullRequest.Base.GetRef(),
			SenderLogin: e.Sender.GetLogin(),
		}, nil

	case *github.IssueCommentEvent:
		var (
			owner = e.Repo.Owner.GetLogin()
			repo  = e.Repo.GetName()
			prNum = e.Issue.GetNumber()
		)
		pr, err := getPullRequest(ctx, client, owner, repo, prNum)
		if err != nil {
			return nil, err
		}

		return &githubEvent{
			Owner:       owner,
			Repo:        repo,
			RepoRemote:  e.Repo.GetSSHURL(),
			PRNumber:    prNum,
			PRMergeable: pr.Mergeable,
			PRClosed:    !pr.GetClosedAt().IsZero(),
			HeadBranch:  pr.Head.GetRef(),
			HeadCommit:  pr.Head.GetSHA(),
			BaseBranch:  pr.Base.GetRef(),
			SenderLogin: e.Sender.GetLogin(),
			IsComment:   true,
			CommentURL:  e.Comment.GetHTMLURL(),
		}, nil

	default:
		return nil, fmt.Errorf("got an unexpected event type, got: %t", e)
	}
}

func sendComment(ctx context.Context, client *github.Client, owner, repo string, prNum int, body string) (*github.IssueComment, error) {
	c, _, err := client.Issues.CreateComment(ctx, owner, repo, prNum, &github.IssueComment{
		Body: &body,
	})
	return c, err
}

func getPullRequest(ctx context.Context, client *github.Client, owner, repo string, prNum int) (*github.PullRequest, error) {
	pr, _, err := client.PullRequests.Get(ctx, owner, repo, prNum)
	return pr, err
}
