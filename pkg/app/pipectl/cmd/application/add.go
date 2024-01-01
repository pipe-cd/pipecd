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

package application

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/pipe-cd/pipecd/pkg/app/server/service/apiservice"
	"github.com/pipe-cd/pipecd/pkg/cli"
	"github.com/pipe-cd/pipecd/pkg/model"
)

type add struct {
	root *command

	appName          string
	appKind          string
	pipedID          string
	platformProvider string
	description      string

	repoID         string
	appDir         string
	configFileName string
}

func newAddCommand(root *command) *cobra.Command {
	c := &add{
		root:           root,
		configFileName: model.DefaultApplicationConfigFilename,
	}
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a new application.",
		RunE:  cli.WithContext(c.run),
	}

	cmd.Flags().StringVar(&c.appName, "app-name", c.appName, "The application name.")
	cmd.Flags().StringVar(&c.appKind, "app-kind", c.appKind, "The kind of application. (KUBERNETES|TERRAFORM|LAMBDA|CLOUDRUN)")
	cmd.Flags().StringVar(&c.pipedID, "piped-id", c.pipedID, "The ID of piped that should handle this application.")
	cmd.Flags().StringVar(&c.platformProvider, "platform-provider", c.platformProvider, "The platform provider name. One of the registered providers in the piped configuration. Previous name of this field is cloud-provider.")

	cmd.Flags().StringVar(&c.repoID, "repo-id", c.repoID, "The repository ID. One the registered repositories in the piped configuration.")
	cmd.Flags().StringVar(&c.appDir, "app-dir", c.appDir, "The relative path from the root of repository to the application directory.")
	cmd.Flags().StringVar(&c.configFileName, "config-file-name", c.configFileName, "The configuration file name")
	cmd.Flags().StringVar(&c.description, "description", c.description, "The description of the application.")

	cmd.MarkFlagRequired("app-name")
	cmd.MarkFlagRequired("app-kind")
	cmd.MarkFlagRequired("piped-id")
	cmd.MarkFlagRequired("platform-provider")
	cmd.MarkFlagRequired("repo-id")
	cmd.MarkFlagRequired("app-dir")

	return cmd
}

func (c *add) run(ctx context.Context, input cli.Input) error {
	cli, err := c.root.clientOptions.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize client: %w", err)
	}
	defer cli.Close()

	appKind, ok := model.ApplicationKind_value[c.appKind]
	if !ok {
		return fmt.Errorf("unsupported application kind %s", c.appKind)
	}

	req := &apiservice.AddApplicationRequest{
		Name:    c.appName,
		PipedId: c.pipedID,
		GitPath: &model.ApplicationGitPath{
			Repo: &model.ApplicationGitRepository{
				Id: c.repoID,
			},
			Path:           c.appDir,
			ConfigFilename: c.configFileName,
		},
		Kind:             model.ApplicationKind(appKind),
		PlatformProvider: c.platformProvider,
		Description:      c.description,
	}

	resp, err := cli.AddApplication(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to add application: %w", err)
	}

	input.Logger.Info(fmt.Sprintf("Successfully added application id = %s", resp.ApplicationId))
	return nil
}
