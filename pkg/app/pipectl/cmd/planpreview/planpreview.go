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

package planpreview

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/pipe-cd/pipe/pkg/app/api/service/apiservice"
	"github.com/pipe-cd/pipe/pkg/app/pipectl/client"
	"github.com/pipe-cd/pipe/pkg/cli"
	"github.com/pipe-cd/pipe/pkg/model"
)

type command struct {
	repoRemoteURL string
	headBranch    string
	headCommit    string
	baseBranch    string
	out           string
	timeout       time.Duration
	checkInterval time.Duration

	clientOptions *client.Options
}

func NewCommand() *cobra.Command {
	c := &command{
		clientOptions: &client.Options{},
		timeout:       10 * time.Minute,
		checkInterval: 10 * time.Second,
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

	cmd.MarkFlagRequired("repo-remote-url")
	cmd.MarkFlagRequired("head-branch")
	cmd.MarkFlagRequired("head-commit")
	cmd.MarkFlagRequired("base-branch")
	cmd.MarkFlagRequired("out")

	return cmd
}

func (c *command) run(ctx context.Context, _ cli.Telemetry) error {
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
	}

	resp, err := cli.RequestPlanPreview(ctx, req)
	if err != nil {
		fmt.Printf("Failed to request plan preview: %v\n", err)
		return err
	}

	getResults := func(commands []string) ([]*model.PlanPreviewCommandResult, error) {
		req := &apiservice.GetPlanPreviewResultsRequest{
			Commands: commands,
		}

		resp, err := cli.GetPlanPreviewResults(ctx, req)
		if err != nil {
			fmt.Printf("Failed to get plan preview results: %v", err)
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
				if status.Code(err) == codes.NotFound {
					break
				}
				return err
			}
			return c.printResults(results)
		}
	}

	return nil
}

func (c *command) printResults(results []*model.PlanPreviewCommandResult) error {
	// TODO: Format preview results and support writing the result into file.
	fmt.Println(results)
	return nil
}
