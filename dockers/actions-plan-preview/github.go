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
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/v36/github"
	"github.com/shurcooL/githubv4"
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

type issueCommentQuery struct {
	ID     githubv4.ID
	Author struct {
		Login githubv4.String
	}
	Body        githubv4.String
	IsMinimized githubv4.Boolean
}

type issueCommentsQuery struct {
	Nodes []issueCommentQuery
}

type pullRequestCommentQuery struct {
	Repository struct {
		PullRequest struct {
			Comments issueCommentsQuery `graphql:"comments(last: 100)"`
		} `graphql:"pullRequest(number: $prNumber)"`
	} `graphql:"repository(owner: $repositoryOwner, name: $repositoryName)"`
}

var errNotFound = errors.New("not found")

// find the latest plan preview comment in the specified issue
// if there is no plan preview comment, return errNotFound err
func findLatestPlanPreviewComment(ctx context.Context, client *githubv4.Client, owner, repo string, prNumber int) (*issueCommentQuery, error) {
	variables := map[string]interface{}{
		"repositoryOwner": githubv4.String(owner),
		"repositoryName":  githubv4.String(repo),
		"prNumber":        githubv4.Int(prNumber),
	}

	var q pullRequestCommentQuery
	if err := client.Query(ctx, &q, variables); err != nil {
		return nil, err
	}

	comment := filterLatestPlanPreviewComment(q.Repository.PullRequest.Comments.Nodes)
	if comment == nil {
		return nil, errNotFound
	}
	return comment, nil
}

// Expect comments to be sorted in ascending order by created_at
func filterLatestPlanPreviewComment(comments []issueCommentQuery) *issueCommentQuery {
	const planPreviewCommentStart = "<!-- pipecd-plan-preview-->"
	const commentLogin = "github-actions"

	for i := range comments {
		comment := comments[len(comments)-i-1]
		if strings.HasPrefix(string(comment.Body), planPreviewCommentStart) && comment.Author.Login == commentLogin {
			return &comment
		}
	}

	return nil
}

type minimizeCommentMutation struct {
	MinimizeComment struct {
		MinimizedComment struct {
			IsMinimized bool
		}
	} `graphql:"minimizeComment(input: $input)"`
}

func minimizeComment(ctx context.Context, client *githubv4.Client, id githubv4.ID, classifier string) error {
	var m minimizeCommentMutation
	input := githubv4.MinimizeCommentInput{
		SubjectID:        id,
		Classifier:       githubv4.ReportedContentClassifiers(classifier),
		ClientMutationID: nil,
	}
	if err := client.Mutate(ctx, &m, input, nil); err != nil {
		return err
	}

	if !m.MinimizeComment.MinimizedComment.IsMinimized {
		return fmt.Errorf("cannot minimize comment. id: %s, classifier: %s", id, classifier)
	}

	return nil
}
