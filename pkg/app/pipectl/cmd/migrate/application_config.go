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
	data, err := os.ReadFile(configFile)
	if err != nil {
		return err
	}

	var cfg map[string]any
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return err
	}

	migrated := make(map[string]any)
	migrated["kind"] = config.KindApplication
	migrated["apiVersion"] = config.VersionV1Beta1

	spec := make(map[string]any)

	// generic spec
	keys := []string{
		"name",
		"labels",
		"description",
		"planner",
		"commitMatcher",
		"pipeline",
		"trigger",
		"postSync",
		"timeout",
		"encryption",
		"attachment",
		"notification",
		"eventWatcher",
		"driftDetection",
	}

	oldSpec := cfg["spec"].(map[string]any)
	for _, key := range keys {
		if _, ok := oldSpec[key]; ok {
			spec[key] = oldSpec[key]
		}
	}

	switch config.Kind(cfg["kind"].(string)) {
	case config.KindKubernetesApp:
		keys := []string{
			"input",
			"quickSync",
			"service",
			"workloads",
			"trafficRouting",
			"variantLabel",
			"resourceRoutes",
		}
		pluginCfg := make(map[string]map[string]any)
		pluginCfg["kubernetes"] = make(map[string]any)
		for _, key := range keys {
			if _, ok := oldSpec[key]; ok {
				pluginCfg["kubernetes"][key] = oldSpec[key]
			}
		}
		spec["plugins"] = pluginCfg
	case config.KindTerraformApp:
		keys := []string{
			"input",
			"quickSync",
		}
		pluginCfg := make(map[string]map[string]any)
		pluginCfg["terraform"] = make(map[string]any)
		for _, key := range keys {
			if _, ok := oldSpec[key]; ok {
				pluginCfg["terraform"][key] = oldSpec[key]
			}
		}
		spec["plugins"] = pluginCfg
	case config.KindECSApp:
		keys := []string{
			"input",
			"quickSync",
		}
		pluginCfg := make(map[string]map[string]any)
		pluginCfg["ecs"] = make(map[string]any)
		for _, key := range keys {
			if _, ok := oldSpec[key]; ok {
				pluginCfg["ecs"][key] = oldSpec[key]
			}
		}
		spec["plugins"] = pluginCfg
	case config.KindLambdaApp:
		keys := []string{
			"input",
			"quickSync",
		}
		pluginCfg := make(map[string]map[string]any)
		pluginCfg["lambda"] = make(map[string]any)
		for _, key := range keys {
			if _, ok := oldSpec[key]; ok {
				pluginCfg["lambda"][key] = oldSpec[key]
			}
		}
		spec["plugins"] = pluginCfg
	case config.KindCloudRunApp:
		keys := []string{
			"input",
			"quickSync",
		}
		pluginCfg := make(map[string]map[string]any)
		pluginCfg["cloudrun"] = make(map[string]any)
		for _, key := range keys {
			if _, ok := oldSpec[key]; ok {
				pluginCfg["cloudrun"][key] = oldSpec[key]
			}
		}
		spec["plugins"] = pluginCfg
	default:
		return fmt.Errorf("unsupported application kind: %s", cfg["kind"])
	}

	migrated["spec"] = spec

	yamlData, err := yaml.Marshal(migrated)
	if err != nil {
		return err
	}

	if err := os.Rename(configFile, configFile+".old"); err != nil {
		return err
	}

	if err := os.WriteFile(configFile, yamlData, 0644); err != nil {
		if e := os.Rename(configFile+".old", configFile); e != nil {
			return errors.Join(err, e)
		}
		return err
	}

	return nil
}
