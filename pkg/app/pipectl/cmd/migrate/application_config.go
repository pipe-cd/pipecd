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
	"io/fs"
	"maps"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"sigs.k8s.io/yaml"

	"github.com/pipe-cd/pipecd/pkg/cli"
	"github.com/pipe-cd/pipecd/pkg/config"
	"github.com/pipe-cd/pipecd/pkg/model"
)

type applicationConfig struct {
	root *command

	configFiles []string
	directories []string
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
	cmd.Flags().StringSliceVar(&c.directories, "dirs", c.directories, "The list of application config directories to migrate.")

	cmd.MarkFlagsOneRequired("config-files", "dirs")
	cmd.MarkFlagsMutuallyExclusive("config-files", "dirs")
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

	for _, directory := range c.directories {
		input.Logger.Info("migrating application configs in directory", zap.String("directory", directory))

		fileSystem := os.DirFS(directory)
		// Scan all files under the repository.
		err := fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}
			if !model.IsApplicationConfigFile(d.Name()) {
				return nil
			}

			input.Logger.Info("migrating application config", zap.String("config-file", path))
			if err := c.migrateApplicationConfig(ctx, filepath.Join(directory, path), input.Logger); err != nil {
				input.Logger.Error("failed to migrate application config", zap.String("config-file", path), zap.Error(err))
				// Continue to migrate other application configs.
				return nil
			}
			input.Logger.Info("successfully migrated application config", zap.String("config-file", path))
			return nil
		})
		if err != nil {
			input.Logger.Error("failed to migrate application configs in directory", zap.String("directory", directory), zap.Error(err))
			return err
		}
		input.Logger.Info("successfully migrated application configs in directory", zap.String("directory", directory))
	}
	return nil
}

func (c *applicationConfig) migrateApplicationConfig(_ context.Context, configFile string, logger *zap.Logger) error {
	data, err := os.ReadFile(configFile)
	if err != nil {
		logger.Error("failed to read application config", zap.String("config-file", configFile), zap.Error(err))
		return err
	}

	var cfg map[string]any
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		logger.Error("failed to unmarshal application config", zap.String("config-file", configFile), zap.Error(err))
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

	// Copy STAGE `timeout` and `skipOn` config under pipeline.stages[].with to pipeline.stages[]
	// NOTE: We keep the original `timeout` and `skipOn` config under pipeline.stages[].with. for backward compatibility,
	//  in case user want to downgrade pipedv1 to pipedv0.
	// `pipeline.stages[].with.{timeout, skipOn}` will be marked as deprecated in v1.
	if oldPipelineCfg, ok := oldSpec["pipeline"]; ok {
		pipelineCfg := make(map[string][]any)
		for _, oldStage := range oldPipelineCfg.(map[string]any)["stages"].([]any) {
			if oldStageCfg, ok := oldStage.(map[string]any); ok {
				stageCfg := maps.Clone(oldStageCfg)
				if withCfg, ok := stageCfg["with"].(map[string]any); ok {
					if _, ok := withCfg["timeout"]; ok {
						stageCfg["timeout"] = withCfg["timeout"]
					}
					if _, ok := withCfg["skipOn"]; ok {
						stageCfg["skipOn"] = withCfg["skipOn"]
					}
				}
				pipelineCfg["stages"] = append(pipelineCfg["stages"], stageCfg)
			}
		}
		spec["pipeline"] = pipelineCfg
	}

	switch config.Kind(cfg["kind"].(string)) {
	case config.KindKubernetesApp:
		logger.Info("migrating kubernetes application config", zap.String("config-file", configFile))
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
		logger.Info("migrating terraform application config", zap.String("config-file", configFile))
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
		logger.Info("migrating ecs application config", zap.String("config-file", configFile))
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
		logger.Info("migrating lambda application config", zap.String("config-file", configFile))
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
		logger.Info("migrating cloudrun application config", zap.String("config-file", configFile))
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
		logger.Error("unsupported application kind", zap.String("config-file", configFile), zap.String("kind", cfg["kind"].(string)))
		return fmt.Errorf("unsupported application kind: %s", cfg["kind"])
	}

	migrated["spec"] = spec

	yamlData, err := yaml.Marshal(migrated)
	if err != nil {
		logger.Error("failed to marshal migrated application config", zap.String("config-file", configFile), zap.Error(err))
		return err
	}

	if err := os.Rename(configFile, configFile+".old"); err != nil {
		logger.Error("failed to rename application config", zap.String("config-file", configFile), zap.Error(err))
		return err
	}

	if err := os.WriteFile(configFile, yamlData, 0644); err != nil {
		logger.Error("failed to write migrated application config", zap.String("config-file", configFile), zap.Error(err))
		if e := os.Rename(configFile+".old", configFile); e != nil {
			logger.Error("failed to restore application config", zap.String("config-file", configFile), zap.Error(e))
			return errors.Join(err, e)
		}
		return err
	}

	return nil
}
