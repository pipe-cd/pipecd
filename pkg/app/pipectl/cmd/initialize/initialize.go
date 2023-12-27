// Copyright 2023 The PipeCD Authors.
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

	"github.com/pipe-cd/pipecd/pkg/app/pipectl/prompt"
	"github.com/pipe-cd/pipecd/pkg/cli"
	"github.com/pipe-cd/pipecd/pkg/config"
)

type command struct {
	// Add flags if needed.
}

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
			fmt.Printf(" Interrupted by signal: %v\n", s)
			cancel()
			os.Exit(1)
		case <-ctx.Done():
		}
	}()

	reader := prompt.NewStdinReader()
	return generateConfig(ctx, input, reader)
}

func generateConfig(ctx context.Context, input cli.Input, reader prompt.Reader) error {
	platform, err := reader.ReadString("Which platform? Enter the number [0]Kubernetes [1]ECS : ")
	if err != nil {
		return fmt.Errorf("invalid input: %v`", err)
	}

	var cfg *genericConfig
	switch platform {
	case "0": // Kubernetes
		// cfg, err = generateKubernetesConfig(in)
		panic("not implemented")
	case "1": // ECS
		cfg, err = generateECSConfig(reader)
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
	exportConfig(cfgBytes, reader)

	return nil
}

func exportConfig(configBytes []byte, reader prompt.Reader) {
	path, err := reader.ReadString("Path to save the config (if not specified, it goes to stdout) : ")
	if err != nil {
		fmt.Printf("Failed to read path %s \n", path)
		printConfig(configBytes)
		return
	}
	if len(path) == 0 {
		// If the target path is not specified, print to stdout.
		printConfig(configBytes)
		return
	}

	// Check if the file/directory already exists and ask if overwrite it.
	if fInfo, err := os.Stat(path); err == nil {
		if fInfo.IsDir() {
			fmt.Printf("The path %s is a directory. Please specify a file path.\n", path)
			printConfig(configBytes)
			return
		}

		// If the file exists, ask if overwrite it.
		overwrite, err := reader.ReadStringRequired(fmt.Sprintf("The file %s already exists. Overwrite it? [y/n] : ", path))
		if err != nil {
			fmt.Printf("Invalid input for overwrite(string): %v\n", err)
			printConfig(configBytes)
			return
		}
		if overwrite != "y" && overwrite != "Y" {
			fmt.Println("Cancelled exporting the config.")
			printConfig(configBytes)
			return
		}
	}

	// If the file does not exist or overwrite, write to the path, including validating.
	fmt.Printf("Start exporting the config to %s\n", path)
	err = os.WriteFile(path, configBytes, 0644)
	if err != nil {
		fmt.Printf("Failed to export the config to %s: %v\n", path, err)
		// If failed, print the config to avoid losing it.
		printConfig(configBytes)
	} else {
		fmt.Printf("Successfully exported the config to %s\n", path)
	}
}

// Print the config to stdout.
func printConfig(configBytes []byte) {
	fmt.Printf("\n### Generated Config is below ###\n%s\n", string(configBytes))
}
