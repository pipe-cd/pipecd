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

package initialize

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"

	"github.com/pipe-cd/pipecd/pkg/app/pipectl/exporter"
	"github.com/pipe-cd/pipecd/pkg/app/pipectl/prompt"
	"github.com/pipe-cd/pipecd/pkg/cli"
	"github.com/pipe-cd/pipecd/pkg/config"
)

type command struct {
	// Add flags if needed.
}

var (
	platform   int
	exportPath string
)

const (
	platformKubernetes = 0
	platformECS        = 1
)

// Use genericConfigs in order to simplify using the spec.
type genericConfig struct {
	APIVersion      string      `json:"apiVersion"`
	Kind            config.Kind `json:"kind"`
	ApplicationSpec interface{} `json:"spec"`
}

func NewCommand() *cobra.Command {
	c := &command{}
	cmd := &cobra.Command{
		Use:     "init",
		Short:   "Generate a app.pipecd.yaml easily and interactively",
		Example: `  pipectl init`,
		Long:    "Generate a app.pipecd.yaml easily, interactively selecting options.",
		RunE:    cli.WithContext(c.run),
	}

	return cmd
}

func (c *command) run(ctx context.Context, input cli.Input) error {
	// Enable interrupt signal.
	ctx, cancel := context.WithCancel(ctx)
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGHUP, syscall.SIGTERM)

	defer func() {
		signal.Stop(signals)
		cancel()
	}()

	go func() {
		select {
		case s := <-signals:
			fmt.Printf("Interrupted by signal: %v\n", s)
			cancel()
			os.Exit(1)
		case <-ctx.Done():
		}
	}()

	p := prompt.NewPrompt(os.Stdin)
	return generateConfig(ctx, input, p)
}

func generateConfig(ctx context.Context, input cli.Input, p prompt.Prompt) error {
	platformInput := prompt.Input{
		Message:       fmt.Sprintf("Which platform? Enter the number [%d]Kubernetes [%d]ECS", platformKubernetes, platformECS),
		TargetPointer: &platform,
		Required:      true,
	}
	exportPathInput := prompt.Input{
		Message:       "Path to save the config (if not specified, it goes to stdout)",
		TargetPointer: &exportPath,
		Required:      false,
	}

	err := p.Run([]prompt.Input{platformInput})
	if err != nil {
		return fmt.Errorf("invalid platform number: %v", err)
	}

	var cfg *genericConfig
	switch platform {
	case platformKubernetes:
		// cfg, err = generateKubernetesConfig(...)
		panic("not implemented")
	case platformECS:
		cfg, err = generateECSConfig(p)
	default:
		return fmt.Errorf("invalid platform number: %d", platform)
	}

	if err != nil {
		return err
	}

	cfgBytes, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	fmt.Println("### The config model was successfully prepared. Move on to exporting. ###")
	err = p.RunOne(exportPathInput)
	if err != nil {
		return err
	}
	if len(exportPath) == 0 {
		printConfig(cfgBytes)
		return nil
	}
	err = exporter.Export(cfgBytes, exportPath)
	if err != nil {
		// fmt.Printf("Failed to export to %s: %v\n", exportPath, err)
		printConfig(cfgBytes)
		return err
	}

	return nil
}

// Print the config to stdout.
func printConfig(configBytes []byte) {
	fmt.Printf("\n### Generated Config is below ###\n%s\n", string(configBytes))
}
