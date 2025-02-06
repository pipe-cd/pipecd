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

package migrate

import (
	"context"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/pipe-cd/pipecd/pkg/app/server/service/apiservice"
	"github.com/pipe-cd/pipecd/pkg/cli"
	"github.com/pipe-cd/pipecd/pkg/model"
)

type database struct {
	root *command

	applications []string
}

func newDatabaseCommand(root *command) *cobra.Command {
	c := &database{
		root: root,
	}
	cmd := &cobra.Command{
		Use:   "database",
		Short: "Do migration tasks for database.",
		Long:  "Make existing applications compatible with plugin-architectured piped. Once you execute this command for an application, it can be deployed using plugin-architectured piped.",
		RunE:  cli.WithContext(c.run),
	}

	cmd.Flags().StringSliceVar(&c.applications, "applications", c.applications, "The list of application to migrate database.")
	cmd.MarkFlagRequired("applications")
	return cmd
}

func (c *database) run(ctx context.Context, input cli.Input) error {
	client, err := c.root.clientOptions.NewClient(ctx)
	if err != nil {
		input.Logger.Error("failed to create client", zap.Error(err))
		return err
	}

	for _, appID := range c.applications {
		input.Logger.Info("migrating database", zap.String("application", appID))
		if err := c.migrateApplication(ctx, client, appID, input.Logger); err != nil {
			input.Logger.Error("failed to migrate database", zap.String("application", appID), zap.Error(err))
			return err
		}
		input.Logger.Info("successfully migrated database", zap.String("application", appID))
	}

	return nil
}

func (c *database) migrateApplication(ctx context.Context, client apiservice.Client, appID string, logger *zap.Logger) error {
	app, err := client.GetApplication(ctx, &apiservice.GetApplicationRequest{ApplicationId: appID})
	if err != nil {
		logger.Error("failed to get application", zap.Error(err), zap.String("application", appID))
		return err
	}

	if len(app.GetApplication().GetDeployTargets()) > 0 {
		logger.Info("skip migrating database because the deploy target is already set", zap.String("application", appID))
		return nil
	}

	provider := app.GetApplication().GetPlatformProvider()

	if provider == "" {
		logger.Info("skip migrating database because the platform provider is not set", zap.String("application", appID))
		return nil
	}

	deployTargets, err := structpb.NewList([]any{provider})
	if err != nil {
		logger.Error("error while determining the deploy targets from previous application platform provider", zap.String("application", appID), zap.Error(err))
		return err
	}
	// Migrate database for the application.
	if _, err := client.UpdateApplicationDeployTargets(ctx, &apiservice.UpdateApplicationDeployTargetsRequest{
		ApplicationId:         appID,
		DeployTargetsByPlugin: map[string]*structpb.ListValue{convertApplicationKindToPluginName(app.Application.Kind): deployTargets},
	}); err != nil {
		logger.Error("failed to update application deploy targets", zap.Error(err), zap.String("application", appID))
		return err
	}

	return nil
}

// NOTE: Convention for Application plugins migration
// The plugins name for this migration task are defined based on the Application Kind
// Eg: KubernetesApp -> kubernetes | ECSApp -> ecs | ...
func convertApplicationKindToPluginName(k model.ApplicationKind) string {
	switch k {
	case model.ApplicationKind_KUBERNETES:
		return "kubernetes"
	case model.ApplicationKind_CLOUDRUN:
		return "cloudrun"
	case model.ApplicationKind_ECS:
		return "ecs"
	case model.ApplicationKind_LAMBDA:
		return "lambda"
	case model.ApplicationKind_TERRAFORM:
		return "terraform"
	}
	return "" // Unexpected
}
