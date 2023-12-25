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

	"github.com/spf13/cobra"

	// "github.com/go-yaml/yaml"
	"gopkg.in/yaml.v3"
	// "sigs.k8s.io/yaml"

	"github.com/pipe-cd/pipecd/pkg/cli"
	"github.com/pipe-cd/pipecd/pkg/config"
)

type command struct {
	someTextOption string
}

// Use genericConfigs in order to
//   - keep the order as we want
//   - use only simple fields, without attatching `omitempty` to all fields
//   - enable modifying the original configs isolately from init command
type genericConfig struct {
	APIVersion      string      `yaml:"apiVersion"`
	Kind            config.Kind `yaml:"kind"`
	ApplicationSpec interface{} `yaml:"spec"`
}

func NewCommand() *cobra.Command {
	c := &command{
		someTextOption: "default-value",
	}
	cmd := &cobra.Command{
		Use:     "init",
		Short:   "Create a app.pipecd.yaml easily (interactively)",
		Example: `  pipectl init`,
		Long:    "Create a app.pipecd.yaml easily, interactively selecting options.",
		RunE:    cli.WithContext(c.run),
	}

	cmd.Flags().StringVar(&c.someTextOption, "some-text-option", c.someTextOption, "Some text option")

	return cmd
}

func (c *command) run(ctx context.Context, input cli.Input) error {
	platform := promptString("Which platform? Enter the number [0]Kubernetes [1]ECS : ")

	var cfg *genericConfig
	var e error
	switch platform {
	case "0": // Kubernetes
		panic("not implemented")
		// cfg := createKubernetesConfig()
	case "1": // ECS
		cfg, e = generateECSConfig()
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

	fmt.Println("### The config model was successfully generated.")

	targetPath := promptString("Path to save the generated config (if not specified, it goes to stdout) : ")
	if len(targetPath) == 0 {
		// If the target path is not specified, print the config to stdout.
		printConfig(cfgBytes)
	} else {
		if _, err := os.Stat(targetPath); err == nil {
			// If the file exists, ask if overwrite it.
			overwrite := promptStringRequired(fmt.Sprintf("The file %s already exists. Overwrite it? [y/n] : ", targetPath))
			if overwrite == "y" || overwrite == "Y" {
				exportConfig(cfgBytes, targetPath)
			} else {
				fmt.Println("Cancelled exporting the config.")
				printConfig(cfgBytes)
			}
		} else {
			// If the file does not exist, simply write to the new file, including validating the path.
			exportConfig(cfgBytes, targetPath)
		}
	}

	return nil
}

// Write the config to the specified path file.
func exportConfig(configBytes []byte, path string) {
	fmt.Printf("Start exporting the config to %s\n", path)
	err := os.WriteFile(path, configBytes, 0644)
	if err != nil {
		fmt.Printf("Failed to export the config to %s: %v\n", path, err)
		// If failed, print the config to prevent losing it.
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
func promptString(message string) string {
	var in string
	fmt.Printf("%s ", message)
	fmt.Scanln(&in)
	return in
}

// Read a string value from stdin, and validate int.
func promptInt(message string) (int, error) {
	var in int
	fmt.Printf("%s ", message)
	_, e := fmt.Scanln(&in)
	if e != nil {
		return 0, e
	}
	return in, nil
}

// Read a string value from stdin, and validate it is not empty.
func promptStringRequired(message string) string {
	for {
		in := promptString(message)
		if in != "" {
			return in
		}
		fmt.Printf("[WARN] This field is required. \n")
	}
}
