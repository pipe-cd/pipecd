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

package planpreview

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/pipe-cd/pipecd/pkg/app/pipectl/client"
	"github.com/pipe-cd/pipecd/pkg/app/server/service/apiservice"
	"github.com/pipe-cd/pipecd/pkg/cli"
	"github.com/pipe-cd/pipecd/pkg/model"
)

const (
	defaultTimeout            = 10 * time.Minute
	defaultPipedHandleTimeout = 5 * time.Minute
	defaultCheckInterval      = 10 * time.Second
	labelEnvKey               = "env"
)

type command struct {
	repoRemoteURL      string
	headBranch         string
	headCommit         string
	baseBranch         string
	out                string
	timeout            time.Duration
	pipedHandleTimeout time.Duration
	checkInterval      time.Duration

	clientOptions *client.Options
}

func NewCommand() *cobra.Command {
	c := &command{
		clientOptions:      &client.Options{},
		pipedHandleTimeout: defaultPipedHandleTimeout,
		timeout:            defaultTimeout,
		checkInterval:      defaultCheckInterval,
	}
	cmd := &cobra.Command{
		Use:   "plan-preview",
		Short: "Show plan preview against the specified commit.",
		RunE:  cli.WithContext(c.run),
	}

	c.clientOptions.RegisterPersistentFlags(cmd)

	cmd.Flags().StringVar(&c.repoRemoteURL, "repo-remote-url", c.repoRemoteURL, "The remote URL of Git repository.")
	cmd.Flags().StringVar(&c.headBranch, "head-branch", c.headBranch, "The head branch of the change.")
	cmd.Flags().StringVar(&c.headCommit, "head-commit", c.headCommit, "The SHA of the head commit.")
	cmd.Flags().StringVar(&c.baseBranch, "base-branch", c.baseBranch, "The base branch of the change.")
	cmd.Flags().StringVar(&c.out, "out", c.out, "Write planpreview result to the given path.")
	cmd.Flags().DurationVar(&c.timeout, "timeout", c.timeout, "Maximum amount of time this command has to complete. Default is 10m.")
	cmd.Flags().DurationVar(&c.pipedHandleTimeout, "piped-handle-timeout", c.pipedHandleTimeout, "Maximum amount of time that a Piped can take to handle. Default is 5m.")

	cmd.MarkFlagRequired("repo-remote-url")
	cmd.MarkFlagRequired("head-branch")
	cmd.MarkFlagRequired("head-commit")
	cmd.MarkFlagRequired("base-branch")

	return cmd
}

func (c *command) run(ctx context.Context, _ cli.Input) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	cli, err := c.clientOptions.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize client: %w", err)
	}
	defer cli.Close()

	req := &apiservice.RequestPlanPreviewRequest{
		RepoRemoteUrl: c.repoRemoteURL,
		HeadBranch:    c.headBranch,
		HeadCommit:    c.headCommit,
		BaseBranch:    c.baseBranch,
		Timeout:       int64(c.pipedHandleTimeout.Seconds()),
	}

	resp, err := cli.RequestPlanPreview(ctx, req)
	if err != nil {
		fmt.Printf("Failed to request plan-preview: %v\n", err)
		return err
	}
	if len(resp.Commands) == 0 {
		fmt.Println("There is no piped that is handling the given Git repository")
		return nil
	}
	fmt.Printf("Requested plan-preview, waiting for its results (commands: %v)\n", resp.Commands)

	getResults := func(commands []string) ([]*model.PlanPreviewCommandResult, error) {
		req := &apiservice.GetPlanPreviewResultsRequest{
			Commands:             commands,
			CommandHandleTimeout: int64(c.pipedHandleTimeout.Seconds()),
		}

		resp, err := cli.GetPlanPreviewResults(ctx, req)
		if err != nil {
			return nil, err
		}

		return resp.Results, nil
	}

	ticker := time.NewTicker(c.checkInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case <-ticker.C:
			results, err := getResults(resp.Commands)
			if err != nil {
				s := status.Convert(err)
				if s.Code() == codes.NotFound {
					fmt.Println(s.Message())
					fmt.Println("waiting...")
					break
				}
				fmt.Printf("Failed to retrieve plan-preview results: %v\n", err)
				return err
			}
			return printResults(results, os.Stdout, c.out)
		}
	}
}

func printResults(results []*model.PlanPreviewCommandResult, stdout io.Writer, outFile string) error {
	r := convert(results)

	// Print out a readable format to stdout.
	fmt.Fprint(stdout, r)

	if outFile == "" {
		return nil
	}

	// Write JSON format to the given file.
	data, err := json.Marshal(r)
	if err != nil {
		fmt.Printf("Failed to encode result to JSON: %v\n", err)
		return err
	}
	return os.WriteFile(outFile, data, 0644)
}

func convert(results []*model.PlanPreviewCommandResult) ReadableResult {
	out := ReadableResult{}
	for _, r := range results {
		if r.Error != "" {
			out.FailurePipeds = append(out.FailurePipeds, FailurePiped{
				PipedInfo: PipedInfo{
					PipedID:  r.PipedId,
					PipedURL: r.PipedUrl,
				},
				Reason: r.Error,
			})
			continue
		}

		for _, a := range r.Results {
			appInfo := ApplicationInfo{
				ApplicationID:        a.ApplicationId,
				ApplicationName:      a.ApplicationName,
				ApplicationURL:       a.ApplicationUrl,
				ApplicationKind:      a.ApplicationKind.String(),
				ApplicationDirectory: a.ApplicationDirectory,
				Env:                  a.Labels[labelEnvKey],
			}
			if a.Error != "" {
				out.FailureApplications = append(out.FailureApplications, FailureApplication{
					ApplicationInfo: appInfo,
					Reason:          a.Error,
					PlanDetails:     string(a.PlanDetails),
				})
				continue
			}
			out.Applications = append(out.Applications, ApplicationResult{
				ApplicationInfo: appInfo,
				SyncStrategy:    a.SyncStrategy.String(),
				PlanSummary:     string(a.PlanSummary),
				PlanDetails:     string(a.PlanDetails),
				NoChange:        a.NoChange,
			})
		}
	}

	return out
}

type ReadableResult struct {
	Applications        []ApplicationResult
	FailureApplications []FailureApplication
	FailurePipeds       []FailurePiped
}

type ApplicationResult struct {
	ApplicationInfo
	SyncStrategy string // QUICK_SYNC, PIPELINE
	PlanSummary  string
	PlanDetails  string
	NoChange     bool
}

type FailurePiped struct {
	PipedInfo
	Reason string
}

type FailureApplication struct {
	ApplicationInfo
	Reason      string
	PlanDetails string
}

type PipedInfo struct {
	PipedID  string
	PipedURL string
}

type ApplicationInfo struct {
	ApplicationID        string
	ApplicationName      string
	ApplicationURL       string
	ApplicationKind      string // KUBERNETES, TERRAFORM, CLOUDRUN, LAMBDA, ECS
	ApplicationDirectory string
	Env                  string
}

func (r ReadableResult) String() string {
	var b strings.Builder
	if len(r.Applications)+len(r.FailureApplications)+len(r.FailurePipeds) == 0 {
		fmt.Fprintf(&b, "\nThere are no updated applications. It means no deployment will be triggered once this pull request got merged.\n")
		return b.String()
	}

	if len(r.Applications) > 0 {
		if len(r.Applications) > 1 {
			fmt.Fprintf(&b, "\nHere are plan-preview for %d applications:\n", len(r.Applications))
		} else {
			fmt.Fprintf(&b, "\nHere are plan-preview for 1 application:\n")
		}
		for i, app := range r.Applications {
			title := fmt.Sprintf("\n%d. app: %s, env: %s, kind: %s\n", i+1, app.ApplicationName, app.Env, app.ApplicationKind)
			if app.Env == "" {
				title = fmt.Sprintf("\n%d. app: %s, kind: %s\n", i+1, app.ApplicationName, app.ApplicationKind)
			}

			b.WriteString(title)
			fmt.Fprintf(&b, "  sync strategy: %s\n", app.SyncStrategy)
			fmt.Fprintf(&b, "  summary: %s\n", app.PlanSummary)
			fmt.Fprintf(&b, "  details:\n\n  ---DETAILS_BEGIN---\n%s\n  ---DETAILS_END---\n", app.PlanDetails)
		}
	}

	if len(r.FailureApplications) > 0 {
		if len(r.FailureApplications) > 1 {
			fmt.Fprintf(&b, "\nNOTE: An error occurred while building plan-preview for the following %d applications:\n", len(r.FailureApplications))
		} else {
			fmt.Fprintf(&b, "\nNOTE: An error occurred while building plan-preview for the following application:\n")
		}
		for i, app := range r.FailureApplications {
			title := fmt.Sprintf("\n%d. app: %s, env: %s, kind: %s\n", i+1, app.ApplicationName, app.Env, app.ApplicationKind)
			if app.Env == "" {
				title = fmt.Sprintf("\n%d. app: %s, kind: %s\n", i+1, app.ApplicationName, app.ApplicationKind)
			}

			b.WriteString(title)
			fmt.Fprintf(&b, "  reason: %s\n", app.Reason)
			if len(app.PlanDetails) > 0 {
				fmt.Fprintf(&b, "  details:\n\n  ---DETAILS_BEGIN---\n%s\n  ---DETAILS_END---\n", app.PlanDetails)
			}
		}
	}

	if len(r.FailurePipeds) > 0 {
		if len(r.FailurePipeds) > 1 {
			fmt.Fprintf(&b, "\nNOTE: An error occurred while building plan-preview for applications of the following %d Pipeds:\n", len(r.FailurePipeds))
		} else {
			fmt.Fprintf(&b, "\nNOTE: An error occurred while building plan-preview for applications of the following Piped:\n")
		}
		for i, piped := range r.FailurePipeds {
			fmt.Fprintf(&b, "\n%d. piped: %s\n", i+1, piped.PipedID)
			fmt.Fprintf(&b, "  reason: %s\n", piped.Reason)
		}
	}

	return b.String()
}
