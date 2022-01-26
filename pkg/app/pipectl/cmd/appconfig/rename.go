// Copyright 2022 The PipeCD Authors.
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

package appconfig

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/pipe-cd/pipecd/pkg/app/server/service/apiservice"
	"github.com/pipe-cd/pipecd/pkg/cli"
	"github.com/pipe-cd/pipecd/pkg/model"
)

type rename struct {
	root *command

	repoRootPath         string
	repoID               string
	envName              string
	before               string
	after                string
	updateAtLocal        bool
	updateOnControlPlane bool
	stdout               io.Writer
}

func newRenameCommand(root *command) *cobra.Command {
	c := &rename{
		root:   root,
		stdout: os.Stdout,
	}
	cmd := &cobra.Command{
		Use:   "rename",
		Short: "Finds all applications that has the given configuration file name to change.",
		RunE:  cli.WithContext(c.run),
	}

	cmd.Flags().StringVar(&c.repoRootPath, "repo-root-path", c.repoRootPath, "The absolute path to the root directory of the Git repository.")
	cmd.Flags().StringVar(&c.repoID, "repo-id", c.repoID, "The repository ID that is being registered in Piped config.")
	cmd.Flags().StringVar(&c.envName, "env-name", c.envName, "The environment name.")
	cmd.Flags().StringVar(&c.before, "before", c.before, "The current name of configuration file.")
	cmd.Flags().StringVar(&c.after, "after", c.after, "The new name of configuration file.")
	cmd.Flags().BoolVar(&c.updateAtLocal, "update-at-local", c.updateAtLocal, "Whether to rename files in Git locally.")
	cmd.Flags().BoolVar(&c.updateOnControlPlane, "update-on-control-plane", c.updateOnControlPlane, "Whether to update application information on control plane to use the new name.")

	cmd.MarkFlagRequired("repo-root-path")
	cmd.MarkFlagRequired("repo-id")
	cmd.MarkFlagRequired("before")
	cmd.MarkFlagRequired("after")

	return cmd
}

func (c *rename) run(ctx context.Context, _ cli.Input) error {
	if !c.updateAtLocal && !c.updateOnControlPlane {
		fmt.Fprintln(c.stdout, "Nothing to do since both --update-at-local and --update-on-control-plane were not set.")
		return nil
	}

	cli, err := c.root.clientOptions.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize client: %w", err)
	}
	defer cli.Close()

	fmt.Fprintln(c.stdout, "Start finding applications to rename their configuration files...")

	var cursor string
	var targets = make([]*model.Application, 0, 0)

	for {
		req := &apiservice.ListApplicationsRequest{
			EnvName: c.envName,
			Cursor:  cursor,
		}
		resp, err := cli.ListApplications(ctx, req)
		if err != nil {
			return fmt.Errorf("failed to list application: %w", err)
		}
		for _, app := range resp.Applications {
			if app.GitPath.Repo.Id != c.repoID {
				continue
			}
			if app.GitPath.GetApplicationConfigFilename() != c.before {
				continue
			}
			targets = append(targets, app)
		}
		if resp.Cursor == "" {
			break
		}
		cursor = resp.Cursor
	}

	fmt.Fprintf(c.stdout, "Found %d applications to update\n", len(targets))
	if len(targets) == 0 {
		return nil
	}

	if c.updateAtLocal {
		for i, app := range targets {
			// Ensure the existence of the current configuration file.
			var (
				oldRelPath = app.GitPath.GetApplicationConfigFilePath()
				oldPath    = filepath.Join(c.repoRootPath, oldRelPath)
			)
			if _, err := os.Stat(oldPath); err != nil {
				return err
			}

			// Ensure that the new name is not conflicting with any existing files.
			var (
				newRelPath = filepath.Join(app.GitPath.Path, c.after)
				newPath    = filepath.Join(c.repoRootPath, newRelPath)
			)
			_, err := os.Stat(newPath)
			if err == nil {
				return fmt.Errorf("unable to use the new name %q since it will override %s", c.after, newPath)
			}
			if !errors.Is(err, os.ErrNotExist) {
				return err
			}

			// Rename file locally.
			if err := os.Rename(oldPath, newPath); err != nil {
				return err
			}
			fmt.Fprintf(c.stdout, "%d. renamed %s to %s\n", i+1, oldRelPath, newRelPath)
		}
		fmt.Fprintf(c.stdout, "Successfully renamed %d applications locally\n", len(targets))
	}

	if c.updateOnControlPlane {
		limit := 20
		for i := 0; i < len(targets); i += limit {
			req := &apiservice.RenameApplicationConfigFileRequest{
				ApplicationIds: make([]string, 0, limit),
				NewFilename:    c.after,
			}
			for j := i; j < len(targets) && j < i+limit; j++ {
				req.ApplicationIds = append(req.ApplicationIds, targets[j].Id)
			}
			_, err := cli.RenameApplicationConfigFile(ctx, req)
			if err != nil {
				return fmt.Errorf("failed to update the configuration file name on the control plane: %w", err)
			}
		}
		fmt.Fprintf(c.stdout, "Successfully renamed %d applications on the control plane\n", len(targets))
	}

	return nil
}
