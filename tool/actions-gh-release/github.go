// Copyright 2023 The PipeCD Authors.
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
	"net/http"
	"regexp"

	"github.com/google/go-github/v39/github"
	"golang.org/x/oauth2"
)

var (
	defaultMergeCommitRegex = regexp.MustCompile(`Merge pull request #([0-9]+) from .+`)
)

type (
	PullRequestState     string
	PullRequestSort      string
	PullRequestDirection string
)

const (
	// PullRequestState
	PullRequestStateOpen   PullRequestState = "open"
	PullRequestStateClosed PullRequestState = "closed"
	PullRequestStateAll    PullRequestState = "all"

	// PullRequestSort
	PullRequestSortCreated     PullRequestSort = "created"
	PullRequestSortUpdated     PullRequestSort = "updated"
	PullRequestSortPopularity  PullRequestSort = "popularity"
	PullRequestSortLongRunning PullRequestSort = "long-running"

	// PullRequestDirection
	PullRequestDirectionAsc  PullRequestDirection = "asc"
	PullRequestDirectionDesc PullRequestDirection = "desc"
)

type githubEvent struct {
	Name string

	Owner string
	Repo  string

	HeadCommit string
	BaseCommit string

	PRNumber   int
	IsComment  bool
	CommentURL string
}

const (
	eventPush         = "push"
	eventPullRequest  = "pull_request"
	eventIssueComment = "issue_comment"
)

type githubClient struct {
	restClient *github.Client
}

type githubClientConfig struct {
	Token string
}

func (p PullRequestState) String() string {
	return string(p)
}

func (p PullRequestSort) String() string {
	return string(p)
}

func (p PullRequestDirection) String() string {
	return string(p)
}

func newGitHubClient(ctx context.Context, cfg *githubClientConfig) *githubClient {
	var httpClient *http.Client
	if cfg.Token != "" {
		t := &oauth2.Token{AccessToken: cfg.Token}
		ts := oauth2.StaticTokenSource(t)
		httpClient = oauth2.NewClient(ctx, ts)
	}

	return &githubClient{
		restClient: github.NewClient(httpClient),
	}
}

// parsePullRequestEvent uses the given environment variables
// to parse and build githubEvent struct.
func (g *githubClient) parseGitHubEvent(ctx context.Context) (*githubEvent, error) {
	var parseEvents = map[string]struct{}{
		eventPush:         {},
		eventPullRequest:  {},
		eventIssueComment: {},
	}

	eventName := os.Getenv("GITHUB_EVENT_NAME")
	if _, ok := parseEvents[eventName]; !ok {
		return &githubEvent{
			Name: eventName,
		}, nil
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
	case *github.PushEvent:
		return &githubEvent{
			Name:       eventPush,
			Owner:      e.Repo.Owner.GetLogin(),
			Repo:       e.Repo.GetName(),
			HeadCommit: e.GetAfter(),
			BaseCommit: e.GetBefore(),
		}, nil

	case *github.PullRequestEvent:
		return &githubEvent{
			Name:       eventPullRequest,
			Owner:      e.Repo.Owner.GetLogin(),
			Repo:       e.Repo.GetName(),
			HeadCommit: e.PullRequest.Head.GetSHA(),
			BaseCommit: e.PullRequest.Base.GetSHA(),
			PRNumber:   e.GetNumber(),
		}, nil

	case *github.IssueCommentEvent:
		var (
			owner = e.Repo.Owner.GetLogin()
			repo  = e.Repo.GetName()
			prNum = e.Issue.GetNumber()
		)

		pr, _, err := g.restClient.PullRequests.Get(ctx, owner, repo, prNum)
		if err != nil {
			return nil, err
		}

		return &githubEvent{
			Name:       eventIssueComment,
			Owner:      owner,
			Repo:       repo,
			HeadCommit: pr.Head.GetSHA(),
			BaseCommit: pr.Base.GetSHA(),
			PRNumber:   prNum,
			IsComment:  true,
			CommentURL: e.Comment.GetHTMLURL(),
		}, nil

	default:
		return nil, fmt.Errorf("got an unexpected event type, got: %t", e)
	}
}

func (g *githubClient) sendComment(ctx context.Context, owner, repo string, prNum int, body string) (*github.IssueComment, error) {
	c, _, err := g.restClient.Issues.CreateComment(ctx, owner, repo, prNum, &github.IssueComment{
		Body: &body,
	})
	return c, err
}

func (g *githubClient) createRelease(ctx context.Context, owner, repo string, p ReleaseProposal) (*github.RepositoryRelease, error) {
	release, _, err := g.restClient.Repositories.CreateRelease(ctx, owner, repo, &github.RepositoryRelease{
		TagName:         &p.Tag,
		Name:            &p.Title,
		TargetCommitish: &p.TargetCommitish,
		Body:            &p.ReleaseNote,
		Prerelease:      &p.Prerelease,
	})
	return release, err
}

func (g *githubClient) existRelease(ctx context.Context, owner, repo, tag string) (bool, error) {
	_, resp, err := g.restClient.Repositories.GetReleaseByTag(ctx, owner, repo, tag)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return false, nil
		}
		return false, err
	}
	return resp.StatusCode == http.StatusOK, nil
}

func (g *githubClient) getPullRequest(ctx context.Context, owner, repo string, number int) (*github.PullRequest, error) {
	pr, _, err := g.restClient.PullRequests.Get(ctx, owner, repo, number)
	return pr, err
}

type ListPullRequestOptions struct {
	State     PullRequestState
	Sort      PullRequestSort
	Direction PullRequestDirection
	Limit     int
}

func (g *githubClient) listPullRequests(ctx context.Context, owner, repo string, opt *ListPullRequestOptions) ([]*github.PullRequest, error) {
	const perPage = 100
	listOpts := github.ListOptions{PerPage: perPage}
	opts := &github.PullRequestListOptions{
		State:       opt.State.String(),
		Sort:        opt.Sort.String(),
		Direction:   opt.Direction.String(),
		ListOptions: listOpts,
	}
	ret := make([]*github.PullRequest, 0, opt.Limit)
	count := opt.Limit / perPage
	for i := 0; i <= count; i++ {
		prs, resp, err := g.restClient.PullRequests.List(ctx, owner, repo, opts)
		if err != nil {
			return nil, err
		}
		for _, pr := range prs {
			if len(ret) == opt.Limit {
				break
			}
			ret = append(ret, pr)
		}
		if resp.NextPage == 0 || len(ret) == opt.Limit {
			break
		}
		opts.Page = resp.NextPage
	}
	return ret, nil
}
