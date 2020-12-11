// Copyright 2020 The PipeCD Authors.
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

	"github.com/spf13/cobra"

	"github.com/pipe-cd/pipe/pkg/app/api/service/apiservice"
	"github.com/pipe-cd/pipe/pkg/cli"
	"github.com/pipe-cd/pipe/pkg/model"
)

type add struct {
	root *command

	appName       string
	appKind       string
	envID         string
	pipedID       string
	cloudProvider string

	gitPathRepoID         string
	gitPathAppDir         string
	gitPathConfigFileName string
}

func newAddCommand(root *command) *cobra.Command {
	c := &add{
		root:                  root,
		gitPathConfigFileName: ".pipe.yaml",
	}
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a new application.",
		RunE:  cli.WithContext(c.run),
	}

	cmd.Flags().StringVar(&c.appName, "app-name", c.appName, "The application name.")
	cmd.Flags().StringVar(&c.appKind, "app-kind", c.appKind, "The kind of application. (KUBERNETES|TERRAFORM|LAMBDA|CLOUDRUN)")
	cmd.Flags().StringVar(&c.envID, "env-id", c.envID, "The ID of environment where this application should belong to.")
	cmd.Flags().StringVar(&c.pipedID, "piped-id", c.pipedID, "The ID of piped that should handle this applicaiton.")
	cmd.Flags().StringVar(&c.cloudProvider, "cloud-provider", c.cloudProvider, "The cloud provider name. One of the registered providers in the piped configuration.")

	cmd.Flags().StringVar(&c.gitPathRepoID, "gitpath-repo-id", c.gitPathRepoID, "The repository ID. One the registered repositories in the piped configuration.")
	cmd.Flags().StringVar(&c.gitPathAppDir, "gitpath-app-dir", c.gitPathAppDir, "The relative path from the root of repository to the application directory.")
	cmd.Flags().StringVar(&c.gitPathConfigFileName, "gitpath-config-file-name", c.gitPathConfigFileName, "The configuration file name. Default is .pipe.yaml")

	cmd.MarkFlagRequired("app-name")
	cmd.MarkFlagRequired("app-kind")
	cmd.MarkFlagRequired("env-id")
	cmd.MarkFlagRequired("piped-id")
	cmd.MarkFlagRequired("cloud-provider")
	cmd.MarkFlagRequired("gitpath-repo-id")
	cmd.MarkFlagRequired("gitpath-app-dir")

	return cmd
}

func (c *add) run(ctx context.Context, _ cli.Telemetry) error {
	logger := c.root.logOptions.NewLogger()
	cli, err := c.root.clientOptions.NewClient(ctx)
	if err != nil {
		logger.Fatal("Failed to initialize client: %v", err)
	}
	defer cli.Close()

	appKind, ok := model.ApplicationKind_value[c.appKind]
	if !ok {
		logger.Fatal("Unsupported application kind %s", c.appKind)
	}

	req := &apiservice.AddApplicationRequest{
		Name:    c.appName,
		EnvId:   c.envID,
		PipedId: c.pipedID,
		GitPath: &model.ApplicationGitPath{
			Repo: &model.ApplicationGitRepository{
				Id: c.gitPathRepoID,
			},
			Path:           c.gitPathAppDir,
			ConfigFilename: c.gitPathConfigFileName,
		},
		Kind:          model.ApplicationKind(appKind),
		CloudProvider: c.cloudProvider,
	}

	resp, err := cli.AddApplication(ctx, req)
	if err != nil {
		logger.Fatal("Failed to add application: %v", err)
	}

	logger.Info("Successfully added application id = %s", resp.ApplicationId)
	return nil
}
