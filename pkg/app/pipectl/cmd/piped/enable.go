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

package piped

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/pipe-cd/pipecd/pkg/app/server/service/apiservice"
	"github.com/pipe-cd/pipecd/pkg/cli"
)

type enable struct {
	root *command

	pipedID string
	stdout  io.Writer
}

func newEnableCommand(root *command) *cobra.Command {
	c := &enable{
		root:   root,
		stdout: os.Stdout,
	}
	cmd := &cobra.Command{
		Use:   "enable",
		Short: "Enable a given Piped.",
		RunE:  cli.WithContext(c.run),
	}

	cmd.Flags().StringVar(&c.pipedID, "piped-id", c.pipedID, "The Piped ID.")
	cmd.MarkFlagRequired("piped-id")

	return cmd
}

func (c *enable) run(ctx context.Context, _ cli.Input) error {
	cli, err := c.root.clientOptions.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize client: %w", err)
	}
	defer cli.Close()

	req := &apiservice.EnablePipedRequest{
		PipedId: c.pipedID,
	}
	if _, err := cli.EnablePiped(ctx, req); err != nil {
		fmt.Fprintf(c.stdout, "Failed to enable Piped %s (%v)\n", c.pipedID, err)
		return err
	}

	fmt.Fprintf(c.stdout, "Successfully enabled Piped %s\n", c.pipedID)
	return nil
}
