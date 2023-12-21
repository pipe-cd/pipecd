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
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/google/go-github/v36/github"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

const (
	defaultTimeout       = 5 * time.Minute
	pullRequestEventName = "pull_request"
	commentEventName     = "issue_comment"
	pushEventName        = "push"

	argAddress            = "address"
	argAPIKey             = "api-key"
	argToken              = "token"
	argTimeout            = "timeout"
	argPipedHandleTimeout = "piped-handle-timeout"
	argPRNum              = "pull-request-number"
)

func main() {
	log.Println("Start running plan-preview")

	eventName := os.Getenv("GITHUB_EVENT_NAME")
	if !isSupportedGitHubEvent(eventName) {
		log.Fatal(fmt.Errorf(
			"unexpected event %s, only %q, %q and %q event are supported",
			eventName,
			pullRequestEventName,
			commentEventName,
			pushEventName,
		))
	}

	args, err := parseArgs(os.Args)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Successfully parsed arguments")

	payload, err := os.ReadFile(os.Getenv("GITHUB_EVENT_PATH"))
	if err != nil {
		log.Fatal(fmt.Errorf("failed to read event payload: %v", err))
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: args.Token},
	)
	tc := oauth2.NewClient(ctx, ts)
	ghClient := github.NewClient(tc)
	ghGraphQLClient := githubv4.NewClient(tc)

	event, err := parseGitHubEvent(ctx, ghClient.PullRequests, eventName, payload, args.PRNum)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Successfully parsed GitHub event\n\tbase-branch %s\n\thead-branch %s\n\thead-commit %s\n", event.BaseBranch, event.HeadBranch, event.HeadCommit)

	doComment := func(body string) {
		comment, err := sendComment(
			ctx,
			ghClient.Issues,
			event.Owner,
			event.Repo,
			event.PRNumber,
			body,
		)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Successfully commented plan-preview result on pull request\n%s\n", *comment.HTMLURL)
	}

	if event.PRClosed {
		doComment(failureBadgeURL + "\nUnable to run plan-preview for a closed pull request.")
		return
	}

	// TODO: When PR opened, `Mergeable` is nil for calculation.
	// Here it is not considered for now, but needs to be handled.
	if event.PRMergeable != nil && *event.PRMergeable == false {
		doComment(failureBadgeURL + "\nUnable to run plan-preview for an un-mergeable pull request. Please resolve the conficts and try again.")
		return
	}

	result, err := retrievePlanPreview(
		ctx,
		event.RepoRemote,
		event.BaseBranch,
		event.HeadBranch,
		event.HeadCommit,
		args.Address,
		args.APIKey,
		args.Timeout,
		args.PipedHandleTimeout,
	)
	if err != nil {
		doComment(failureBadgeURL + "\nUnable to run plan-preview. \ncause: " + err.Error())
		log.Fatal(err)
	}
	log.Println("Successfully retrieved plan-preview result")

	// Maybe the PR is already closed so Piped could not clone the source code.
	if result.HasError() {
		pr, err := getPullRequest(ctx, ghClient.PullRequests, event.Owner, event.Repo, event.PRNumber)
		if err != nil {
			doComment(failureBadgeURL + "\nUnable to run plan-preview. \ncause: " + err.Error())
			log.Fatal(err)
		}
		if !pr.GetClosedAt().IsZero() {
			doComment(failureBadgeURL + "\nUnable to run plan-preview for a closed pull request.")
			return
		}

		minimizePreviousComment(ctx, ghGraphQLClient, event)

		body := makeCommentBody(event, result)
		doComment(body)

		log.Println("Successfully minimized last plan-preview result on pull request")
		log.Println("plan-preview result has error")
		os.Exit(1)
	}

	minimizePreviousComment(ctx, ghGraphQLClient, event)

	body := makeCommentBody(event, result)
	doComment(body)
}

type arguments struct {
	Address            string
	APIKey             string
	Token              string
	Timeout            time.Duration
	PipedHandleTimeout time.Duration
	PRNum              int
}

func parseArgs(args []string) (arguments, error) {
	var out arguments

	for _, arg := range args {
		ps := strings.SplitN(arg, "=", 2)
		if len(ps) != 2 {
			continue
		}
		switch ps[0] {
		case argAddress:
			out.Address = ps[1]
		case argAPIKey:
			out.APIKey = ps[1]
		case argToken:
			out.Token = ps[1]
		case argTimeout:
			d, err := time.ParseDuration(ps[1])
			if err != nil {
				return arguments{}, err
			}
			out.Timeout = d
		case argPipedHandleTimeout:
			d, err := time.ParseDuration(ps[1])
			if err != nil {
				return arguments{}, err
			}
			out.PipedHandleTimeout = d
		case argPRNum:
			if ps[1] == "" {
				continue
			}
			i, err := strconv.Atoi(ps[1])
			if err != nil {
				return out, err
			}
			if i <= 0 {
				return out, fmt.Errorf("invalid %s: %d", argPRNum, i)
			}
			out.PRNum = i
		}
	}

	if out.Address == "" {
		return out, fmt.Errorf("missing address argument")
	}
	if out.APIKey == "" {
		return out, fmt.Errorf("missing api-key argument")
	}
	if out.Token == "" {
		return out, fmt.Errorf("missing token argument")
	}
	if out.Timeout == 0 {
		out.Timeout = defaultTimeout
	}
	if out.PipedHandleTimeout == 0 {
		out.PipedHandleTimeout = defaultTimeout
	}

	return out, nil
}

func isSupportedGitHubEvent(event string) bool {
	return event == pullRequestEventName || event == commentEventName || event == pushEventName
}
