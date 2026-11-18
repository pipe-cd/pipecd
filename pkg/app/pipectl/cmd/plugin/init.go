// Copyright 2026 The PipeCD Authors.
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

package plugin

import (
	"context"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/app/pipectl/pluginscaffold"
	"github.com/pipe-cd/pipecd/pkg/cli"
)

type initCommand struct {
	root *command

	name      string
	module    string
	stages    []string
	outputDir string
}

func newInitCommand(root *command) *cobra.Command {
	i := &initCommand{
		root: root,
	}
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Generate a new stage plugin.",
		Long:  "Generate the source of a stage plugin that compiles and passes its own tests, with one stub for each stage.",
		Example: `  pipectl plugin init --name my-plugin --module github.com/your-org/my-plugin
  pipectl plugin init --name my-plugin --module github.com/your-org/my-plugin --stages MY_PREPARE,MY_CHECK`,
		RunE: cli.WithContext(i.run),
	}

	cmd.Flags().StringVar(&i.name, "name", i.name, "The name of the plugin, as it is registered in the piped config. For example, 'my-plugin'.")
	cmd.Flags().StringVar(&i.module, "module", i.module, "The Go module path of the plugin. For example, 'github.com/your-org/my-plugin'.")
	cmd.Flags().StringSliceVar(&i.stages, "stages", i.stages, "The stages of the plugin, in UPPER_SNAKE_CASE. For example, 'MY_PREPARE,MY_CHECK'. Defaults to one stage named after the plugin.")
	cmd.Flags().StringVar(&i.outputDir, "output-dir", i.outputDir, "The directory to generate the plugin in. It must not exist or must be empty. Defaults to a directory named after the plugin.")

	cmd.MarkFlagRequired("name")
	cmd.MarkFlagRequired("module")

	return cmd
}

func (i *initCommand) run(ctx context.Context, input cli.Input) error {
	files, err := pluginscaffold.Generate(pluginscaffold.Options{
		Name:   i.name,
		Module: i.module,
		Stages: i.stages,
	})
	if err != nil {
		return err
	}

	dir := i.outputDir
	if dir == "" {
		dir = i.name
	}

	if err := pluginscaffold.Write(dir, files); err != nil {
		return err
	}

	input.Logger.Info("successfully generated the plugin", zap.String("dir", dir), zap.Int("files", len(files)))

	return nil
}
