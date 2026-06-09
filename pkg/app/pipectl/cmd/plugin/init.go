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
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/app/pipectl/pluginscaffold"
	"github.com/pipe-cd/pipecd/pkg/cli"
)

type initCmd struct {
	root *command

	outputDir  string
	kind       string
	stages     []string
	modulePath string
	pluginName string
	dryRun     bool
	force      bool
}

func newInitCommand(root *command) *cobra.Command {
	c := &initCmd{root: root}
	cmd := &cobra.Command{
		Use:   "init [output-dir]",
		Short: "Scaffold a new Piped v1 plugin.",
		Long: `Generate a compile-ready plugin skeleton for stage-only or deployment-minimal plugins.

With --force, an existing output directory is removed before writing. Double-check the path to avoid deleting unrelated data.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c.outputDir = args[0]
			return cli.WithContext(c.run)(cmd, args)
		},
	}

	cmd.Flags().StringVar(&c.kind, "kind", string(pluginscaffold.KindStage), "Plugin kind: stage or deployment")
	cmd.Flags().StringSliceVar(&c.stages, "stages", nil, "Comma-separated stage names in UPPER_SNAKE_CASE (required)")
	cmd.Flags().StringVar(&c.modulePath, "module", "", "Go module path (default: github.com/example/piped-plugin-<name>)")
	cmd.Flags().StringVar(&c.pluginName, "name", "", "Plugin name (default: basename of output-dir)")
	cmd.Flags().BoolVar(&c.dryRun, "dry-run", false, "Print files without writing")
	cmd.Flags().BoolVar(&c.force, "force", false, "Remove and recreate the output directory if it already exists")

	_ = cmd.MarkFlagRequired("stages")

	return cmd
}

func (c *initCmd) run(_ context.Context, input cli.Input) error {
	outputDir, err := filepath.Abs(c.outputDir)
	if err != nil {
		return err
	}

	pluginName := c.pluginName
	if pluginName == "" {
		pluginName = filepath.Base(outputDir)
	}

	stages := normalizeStages(c.stages)
	if len(stages) == 0 {
		return fmt.Errorf("--stages is required")
	}

	kind := pluginscaffold.Kind(c.kind)
	opts := pluginscaffold.Options{
		OutputDir:  outputDir,
		PluginName: pluginName,
		ModulePath: c.modulePath,
		Kind:       kind,
		Stages:     stages,
		DryRun:     c.dryRun,
		Force:      c.force,
	}

	files, err := pluginscaffold.Generate(opts)
	if err != nil {
		return err
	}

	if c.dryRun {
		for _, f := range files {
			input.Logger.Info("would write file", zap.String("path", f.Path))
		}
		return nil
	}

	if err := pluginscaffold.Write(opts, files); err != nil {
		return err
	}

	input.Logger.Info("plugin scaffold written", zap.String("dir", outputDir), zap.Int("files", len(files)))
	return nil
}

func normalizeStages(in []string) []string {
	var out []string
	for _, s := range in {
		for _, part := range strings.Split(s, ",") {
			part = strings.TrimSpace(part)
			if part != "" {
				out = append(out, part)
			}
		}
	}
	return out
}
