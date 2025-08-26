// Copyright 2025 The PipeCD Authors.
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
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"sigs.k8s.io/yaml"

	"github.com/pipe-cd/pipecd/pkg/cli"
	"github.com/pipe-cd/pipecd/pkg/config"
)

type applicationConfig struct {
	root *command

	configFiles []string
}

func newApplicationConfigCommand(root *command) *cobra.Command {
	c := &applicationConfig{
		root: root,
	}
	cmd := &cobra.Command{
		Use:   "application-config",
		Short: "Do migration tasks for application config.",
		Long:  "Make existing applications compatible with plugin-architectured piped. Once you execute this command for an application, it can be deployed using plugin-architectured piped.",
		RunE:  cli.WithContext(c.run),
	}

	cmd.Flags().StringSliceVar(&c.configFiles, "config-files", c.configFiles, "The list of application config files to migrate.")
	cmd.MarkFlagRequired("config-files")
	return cmd
}

func (c *applicationConfig) run(ctx context.Context, input cli.Input) error {
	for _, configFile := range c.configFiles {
		input.Logger.Info("migrating application config", zap.String("config-file", configFile))
		if err := c.migrateApplicationConfig(ctx, configFile, input.Logger); err != nil {
			input.Logger.Error("failed to migrate application config", zap.String("config-file", configFile), zap.Error(err))
			return err
		}
		input.Logger.Info("successfully migrated application config", zap.String("config-file", configFile))
	}

	return nil
}

func (c *applicationConfig) migrateApplicationConfig(ctx context.Context, configFile string, logger *zap.Logger) error {
	cfg, err := config.LoadFromYAML(configFile)
	if err != nil {
		logger.Error("failed to load application config", zap.String("config-file", configFile), zap.Error(err))
		return err
	}

	migrated := make(map[string]any)

	genericSpec, ok := cfg.GetGenericApplication()
	if !ok {
		logger.Error("failed to get generic application spec", zap.String("config-file", configFile))
		return fmt.Errorf("failed to get generic application spec")
	}

	switch cfg.Kind {
	case config.KindCloudRunApp:
		genericSpec.Plugins = map[string]any{
			"cloudrun": (*migratedCloudRunApplicationSpec)(cfg.CloudRunApplicationSpec),
		}
	case config.KindECSApp:
		genericSpec.Plugins = map[string]any{
			"ecs": (*migratedECSApplicationSpec)(cfg.ECSApplicationSpec),
		}
	case config.KindKubernetesApp:
		genericSpec.Plugins = map[string]any{
			"kubernetes": (*migratedKubernetesApplicationSpec)(cfg.KubernetesApplicationSpec),
		}
	case config.KindLambdaApp:
		genericSpec.Plugins = map[string]any{
			"lambda": (*migratedLambdaApplicationSpec)(cfg.LambdaApplicationSpec),
		}
	case config.KindTerraformApp:
		genericSpec.Plugins = map[string]any{
			"terraform": (*migratedTerraformApplicationSpec)(cfg.TerraformApplicationSpec),
		}
	}

	migrated["kind"] = config.KindApplication
	migrated["apiVersion"] = config.VersionV1Beta1
	migrated["spec"] = (*migratedApplicationSpec)(&genericSpec)

	b, err := json.Marshal(migrated)
	if err != nil {
		logger.Error("failed to marshal migrated application config", zap.String("config-file", configFile), zap.Error(err))
		return err
	}

	y, err := yaml.JSONToYAML(b)
	if err != nil {
		logger.Error("failed to convert migrated application config to YAML", zap.String("config-file", configFile), zap.Error(err))
		return err
	}

	if err := os.Rename(configFile, configFile+".old"); err != nil {
		logger.Error("failed to rename application config file", zap.String("config-file", configFile), zap.Error(err))
		return err
	}

	if err := os.WriteFile(configFile, y, 0644); err != nil {
		logger.Error("failed to write migrated application config", zap.String("config-file", configFile), zap.Error(err))

		// If the write failed, we need to restore the old config file.
		if e := os.Rename(configFile+".old", configFile); e != nil {
			logger.Error("failed to rename application config file", zap.String("config-file", configFile), zap.Error(e))
			return errors.Join(err, e)
		}

		return err
	}

	logger.Info("successfully migrated application config", zap.String("config-file", configFile))
	return nil
}

type migratedApplicationSpec config.GenericApplicationSpec

func (c *migratedApplicationSpec) MarshalJSON() ([]byte, error) {

	type migratedPipelineStage struct {
		ID      string          `json:"id,omitempty"`
		Name    string          `json:"name"`
		Desc    string          `json:"desc,omitempty"`
		Timeout config.Duration `json:"timeout,omitempty"`
		With    json.RawMessage `json:"with,omitempty"`
	}

	type migratedDeploymentPipeline struct {
		Stages []migratedPipelineStage `json:"stages"`
	}

	type spec struct {
		// The application name.
		// This is required if you set the application through the application configuration file.
		Name string `json:"name,omitempty"`
		// Additional attributes to identify applications.
		Labels map[string]string `json:"labels,omitempty"`
		// Notes on the Application.
		Description string `json:"description,omitempty"`

		// Configuration used while planning deployment.
		Planner config.DeploymentPlanner `json:"planner,omitempty"`
		// Forcibly use QuickSync or Pipeline when commit message matched the specified pattern.
		CommitMatcher config.DeploymentCommitMatcher `json:"commitMatcher,omitempty"`
		// Pipeline for deploying progressively.
		Pipeline *migratedDeploymentPipeline `json:"pipeline"`
		// The trigger configuration use to determine trigger logic.
		Trigger config.Trigger `json:"trigger,omitempty"`
		// Configuration to be used once the deployment is triggered successfully.
		PostSync *config.PostSync `json:"postSync,omitempty"`
		// The maximum length of time to execute deployment before giving up.
		// Default is 6h.
		Timeout config.Duration `json:"timeout,omitempty" default:"6h"`
		// List of encrypted secrets and targets that should be decoded before using.
		Encryption *config.SecretEncryption `json:"encryption,omitempty"`
		// List of files that should be attached to application manifests before using.
		Attachment *config.Attachment `json:"attachment,omitempty"`
		// Additional configuration used while sending notification to external services.
		DeploymentNotification *config.DeploymentNotification `json:"notification,omitempty"`
		// List of the configuration for event watcher.
		EventWatcher []config.EventWatcherConfig `json:"eventWatcher,omitempty"`
		// Configuration for drift detection
		DriftDetection *config.DriftDetection `json:"driftDetection,omitempty"`

		// This is a workaround not to raise unknown-field error when the application config file contains the plugins field.
		Plugins any `json:"plugins"`
	}

	pipeline := &migratedDeploymentPipeline{}
	for _, stage := range c.Pipeline.Stages {
		pipeline.Stages = append(pipeline.Stages, migratedPipelineStage{
			ID:      stage.ID,
			Name:    string(stage.Name),
			Desc:    stage.Desc,
			Timeout: stage.Timeout,
			With:    stage.With,
		})
	}

	return json.Marshal(spec{
		Name:                   c.Name,
		Labels:                 c.Labels,
		Description:            c.Description,
		Planner:                c.Planner,
		CommitMatcher:          c.CommitMatcher,
		Pipeline:               pipeline,
		Trigger:                c.Trigger,
		PostSync:               c.PostSync,
		Timeout:                c.Timeout,
		Encryption:             c.Encryption,
		Attachment:             c.Attachment,
		DeploymentNotification: c.DeploymentNotification,
		EventWatcher:           c.EventWatcher,
		DriftDetection:         c.DriftDetection,
		Plugins:                c.Plugins,
	})
}

type migratedCloudRunApplicationSpec config.CloudRunApplicationSpec

func (c *migratedCloudRunApplicationSpec) MarshalJSON() ([]byte, error) {
	// spec is copied from config.CloudRunApplicationSpec, but we need to remove the GenericApplicationSpec field.
	type spec struct {
		// Input for CloudRun deployment such as docker image...
		Input config.CloudRunDeploymentInput `json:"input"`
		// Configuration for quick sync.
		QuickSync config.CloudRunSyncStageOptions `json:"quickSync"`
	}
	return json.Marshal(spec{
		Input:     c.Input,
		QuickSync: c.QuickSync,
	})
}

type migratedECSApplicationSpec config.ECSApplicationSpec

func (c *migratedECSApplicationSpec) MarshalJSON() ([]byte, error) {
	// spec is copied from config.ECSApplicationSpec, but we need to remove the GenericApplicationSpec field.
	type spec struct {
		// Input for ECS deployment such as where to fetch source code...
		Input config.ECSDeploymentInput `json:"input"`
		// Configuration for quick sync.
		QuickSync config.ECSSyncStageOptions `json:"quickSync"`
	}
	return json.Marshal(spec{
		Input:     c.Input,
		QuickSync: c.QuickSync,
	})
}

type migratedKubernetesApplicationSpec config.KubernetesApplicationSpec

func (c *migratedKubernetesApplicationSpec) MarshalJSON() ([]byte, error) {
	// spec is copied from config.KubernetesApplicationSpec, but we need to remove the GenericApplicationSpec field.
	type spec struct {
		// Input for Kubernetes deployment such as kubectl version, helm version, manifests filter...
		Input config.KubernetesDeploymentInput `json:"input"`
		// Configuration for quick sync.
		QuickSync config.K8sSyncStageOptions `json:"quickSync"`
		// Which resource should be considered as the Service of application.
		// Empty means the first Service resource will be used.
		Service config.K8sResourceReference `json:"service"`
		// Which resources should be considered as the Workload of application.
		// Empty means all Deployments.
		// e.g.
		// - kind: Deployment
		//   name: deployment-name
		// - kind: ReplicationController
		//   name: replication-controller-name
		Workloads []config.K8sResourceReference `json:"workloads"`
		// Which method should be used for traffic routing.
		TrafficRouting *config.KubernetesTrafficRouting `json:"trafficRouting"`
		// The label will be configured to variant manifests used to distinguish them.
		VariantLabel config.KubernetesVariantLabel `json:"variantLabel"`
		// List of route configurations to resolve the platform provider for application resources.
		// Each resource will be checked over the match conditions of each route.
		// If matches, it will be applied to the route's provider,
		// otherwise, it will be fallen through the next route to check.
		// Any resource which does not match any specified route will be applied
		// to the default platform provider which had been specified while registering the application.
		ResourceRoutes []config.KubernetesResourceRoute `json:"resourceRoutes"`
	}

	return json.Marshal(spec{
		Input:          c.Input,
		QuickSync:      c.QuickSync,
		Service:        c.Service,
		Workloads:      c.Workloads,
		TrafficRouting: c.TrafficRouting,
		VariantLabel:   c.VariantLabel,
		ResourceRoutes: c.ResourceRoutes,
	})
}

type migratedLambdaApplicationSpec config.LambdaApplicationSpec

func (c *migratedLambdaApplicationSpec) MarshalJSON() ([]byte, error) {
	// spec is copied from config.LambdaApplicationSpec, but we need to remove the GenericApplicationSpec field.
	type spec struct {
		// Input for Lambda deployment such as where to fetch source code...
		Input config.LambdaDeploymentInput `json:"input"`
		// Configuration for quick sync.
		QuickSync config.LambdaSyncStageOptions `json:"quickSync"`
	}
	return json.Marshal(spec{
		Input:     c.Input,
		QuickSync: c.QuickSync,
	})
}

type migratedTerraformApplicationSpec config.TerraformApplicationSpec

func (c *migratedTerraformApplicationSpec) MarshalJSON() ([]byte, error) {
	// spec is copied from config.TerraformApplicationSpec, but we need to remove the GenericApplicationSpec field.
	type spec struct {
		// Input for Terraform deployment such as terraform version, workspace...
		Input config.TerraformDeploymentInput `json:"input"`
		// Configuration for quick sync.
		QuickSync config.TerraformApplyStageOptions `json:"quickSync"`
	}
	return json.Marshal(spec{
		Input:     c.Input,
		QuickSync: c.QuickSync,
	})
}
