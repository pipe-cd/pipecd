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

package image

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/pipe-cd/pipe/pkg/app/api/service/apiservice"
	"github.com/pipe-cd/pipe/pkg/cli"
)

type pushReference struct {
	root *command

	repoName string
	tags     []string
	digest   string
}

func newPushReferenceCommand(root *command) *cobra.Command {
	c := &pushReference{
		root: root,
	}
	cmd := &cobra.Command{
		Use:   "push-reference",
		Short: "Push a container image reference.",
		RunE:  cli.WithContext(c.run),
	}

	cmd.Flags().StringVar(&c.repoName, "repo-name", c.repoName, "The repository name of container image. e.g. gcr.io/pipecd/pipecd, envoyproxy/envoy-alpine...")
	cmd.Flags().StringSliceVar(&c.tags, "tag", c.tags, "The image tag.")
	cmd.Flags().StringVar(&c.digest, "digest", c.digest, "The image digest.")

	cmd.MarkFlagRequired("repo-name")
	cmd.MarkFlagRequired("tag")

	return cmd
}

func (c *pushReference) run(ctx context.Context, t cli.Telemetry) error {
	cli, err := c.root.clientOptions.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize client: %w", err)
	}
	defer cli.Close()

	req := &apiservice.PushImageReferenceRequest{
		RepoName: c.repoName,
		Tags:     c.tags,
		Digest:   c.digest,
	}

	if _, err := cli.PushImageReference(ctx, req); err != nil {
		return fmt.Errorf("failed to push image reference: %w", err)
	}

	t.Logger.Info("Successfully pushed image reference")
	return nil
}
