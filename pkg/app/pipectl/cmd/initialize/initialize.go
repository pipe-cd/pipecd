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
	"io"
	"os"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"

	"github.com/pipe-cd/pipecd/pkg/cli"
	"github.com/pipe-cd/pipecd/pkg/config"
)

type command struct {
	someTextOption string
}

// Use genericConfigs in order to simplify using the spec.
type genericConfig struct {
	APIVersion      string      `json:"apiVersion"`
	Kind            config.Kind `json:"kind"`
	ApplicationSpec interface{} `json:"spec"`
}

func NewCommand() *cobra.Command {
	c := &command{
		someTextOption: "default-value",
	}
	cmd := &cobra.Command{
		Use:     "init",
		Short:   "Generate a app.pipecd.yaml easily and interactively",
		Example: `  pipectl init`,
		Long:    "Generate a app.pipecd.yaml easily, interactively selecting options.",
		RunE:    cli.WithContext(c.run),
	}

	cmd.Flags().StringVar(&c.someTextOption, "some-text-option", c.someTextOption, "Some text option")

	return cmd
}

func (c *command) run(ctx context.Context, input cli.Input) error {
	return generateConfig(ctx, input, os.Stdin)
}

func generateConfig(ctx context.Context, input cli.Input, in io.Reader) error {
	platform := promptString("Which platform? Enter the number [0]Kubernetes [1]ECS : ", in)

	var cfg *genericConfig
	var e error
	switch platform {
	case "0": // Kubernetes
		panic("not implemented")
		// cfg := createKubernetesConfig()
	case "1": // ECS
		cfg, e = generateECSConfig(in)
	default:
		return fmt.Errorf("invalid platform number: %s", platform)
	}

	if e != nil {
		return e
	}

	cfgBytes, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	fmt.Println("### The config model was successfully generated. Move on to exporting. ###")

	targetPath := promptString("Path to save the generated config (if not specified, it goes to stdout) : ", in)
	if len(targetPath) == 0 {
		// If the target path is not specified, print to stdout.
		printConfig(cfgBytes)
	} else {
		exportConfig(cfgBytes, targetPath, in)
	}

	return nil
}

// Write the config to the specified path.
func exportConfig(configBytes []byte, path string, in io.Reader) {
	if fInfo, err := os.Stat(path); err == nil {
		if fInfo.IsDir() {
			fmt.Printf("The path %s is a directory. Please specify a file path.\n", path)
			printConfig(configBytes)
			return
		}

		// If the file exists, ask if overwrite it.
		overwrite := promptStringRequired(fmt.Sprintf("The file %s already exists. Overwrite it? [y/n] : ", path), in)
		if overwrite != "y" && overwrite != "Y" {
			fmt.Println("Cancelled exporting the config.")
			printConfig(configBytes)
			return
		}
	}

	// If the file does not exist or overwrite, write to the path, including validating.
	fmt.Printf("Start exporting the config to %s\n", path)
	err := os.WriteFile(path, configBytes, 0644)
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

// Read a string value from stdin.
func promptString(message string, in io.Reader) string {
	var s string
	fmt.Printf("%s ", message)
	fmt.Fscanln(in, &s)
	return s
}

// Read a string value from stdin, and validate int.
func promptInt(message string, in io.Reader) (int, error) {
	var s int
	fmt.Printf("%s ", message)
	_, e := fmt.Fscanln(in, &s)
	if e != nil {
		return 0, e
	}
	return s, nil
}

// Read a string value from stdin, and validate it is not empty.
func promptStringRequired(message string, in io.Reader) string {
	for {
		in := promptString(message, in)
		if in != "" {
			return in
		}
		fmt.Printf("[WARN] This field is required. \n")
	}
}
