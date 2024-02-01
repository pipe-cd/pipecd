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

	"github.com/pipe-cd/pipecd/pkg/app/pipectl/cmd/initialize/exporter"
	"github.com/pipe-cd/pipecd/pkg/app/pipectl/cmd/initialize/prompt"
	"github.com/pipe-cd/pipecd/pkg/cli"
	"github.com/pipe-cd/pipecd/pkg/config"
)

type command struct {
	// Add flags if needed.
}

const (
	// platform numbers to select which platform to use.
	platformKubernetes string = "0" // for KubernetesApp
	platformECS        string = "1" // for ECSApp
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
		Short:   "Generate an application config (app.pipecd.yaml) easily and interactively.",
		Example: `  pipectl init`,
		Long:    "Generate an application config (app.pipecd.yaml) easily, interactively selecting options.",
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
	// user's inputs
	var (
		platform   string
		exportPath string
	)

	platformInput := prompt.Input{
		Message:       fmt.Sprintf("Which platform? Enter the number [%s]Kubernetes [%s]ECS", platformKubernetes, platformECS),
		TargetPointer: &platform,
		Required:      true,
	}
	exportPathInput := prompt.Input{
		Message:       "Path to save the config (if not specified, it goes to stdout)",
		TargetPointer: &exportPath,
		Required:      false,
	}

	err := p.Run(platformInput)
	if err != nil {
		return fmt.Errorf("invalid platform number: %v", err)
	}

	var cfg *genericConfig
	switch platform {
	case platformKubernetes:
		panic("not implemented")
	case platformECS:
		cfg, err = generateECSConfig(p)
	default:
		return fmt.Errorf("invalid platform number: %s", platform)
	}

	if err != nil {
		return err
	}

	cfgBytes, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	fmt.Println("### The config model was successfully prepared. Move on to exporting. ###")
	err = p.Run(exportPathInput)
	if err != nil {
		printConfig(cfgBytes)
		return err
	}
	err = export(cfgBytes, exportPath)
	if err != nil {
		return nil
	}

	return nil
}

func export(cfgBytes []byte, exportPath string) error {
	if len(exportPath) == 0 {
		// if the path is not specified, print to stdout
		printConfig(cfgBytes)
		return nil
	}
	err := exporter.Export(cfgBytes, exportPath)
	if err != nil {
		printConfig(cfgBytes)
		return err
	}
	return nil
}

// Print the config to stdout.
func printConfig(configBytes []byte) {
	fmt.Printf("\n### Generated Config is below ###\n%s\n", string(configBytes))
}
