// Copyright 2022 The PipeCD Authors.
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

package quickstart

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/spf13/cobra"

	"github.com/pipe-cd/pipecd/pkg/cli"
	"github.com/pipe-cd/pipecd/pkg/version"
)

const (
	pipecdDefaultNamespace = "pipecd"

	helmBinaryName    = "helm"
	helmReleaseName   = "pipecd"
	helmChartRepoName = "oci://ghcr.io/pipe-cd/chart/pipecd"

	helmQuickstartValueRemotePath = "https://raw.githubusercontent.com/pipe-cd/pipecd/%s/quickstart/control-plane-values.yaml"
)

type command struct {
	version   string
	toolsDir  string
	namespace string

	uninstall bool
}

func NewCommand() *cobra.Command {
	home, err := os.UserHomeDir()
	if err != nil {
		{
			panic(fmt.Sprintf("failed to detect the current user's home directory: %v", err))
		}
	}

	c := &command{
		version:   version.Get().Version,
		toolsDir:  path.Join(home, ".pipectl", "tools"),
		namespace: pipecdDefaultNamespace,
	}

	cmd := &cobra.Command{
		Use:   "quickstart",
		Short: "Quick prepare PipeCD control plane in quickstart mode.",
		Long:  "Quick prepare PipeCD control plane in quickstart mode.\nTo install PipeCD control plane for real-life usage, please read the docs: https://pipecd.dev/docs/installation/install-controlplane",
		RunE:  cli.WithContext(c.run),
	}

	cmd.Flags().StringVar(&c.version, "version", c.version, "The Control Plane version. Default is the version of pipectl.")
	cmd.Flags().StringVar(&c.toolsDir, "tools-dir", c.toolsDir, "The path to directory where to install tools such as helm.")
	cmd.Flags().StringVar(&c.namespace, "namespace", c.namespace, "The Kubernetes cluster namespace where to install Control Plane.")

	cmd.Flags().BoolVar(&c.uninstall, "uninstall", c.uninstall, "Uninstall the quickstart mode installed PipeCD control plane.")

	return cmd
}

func (c *command) run(ctx context.Context, input cli.Input) error {
	helm, err := c.getHelm()
	if err != nil {
		return fmt.Errorf("failed to prepare required tools (helm) for installation: %v", err)
	}

	var args []string

	if c.uninstall {
		input.Logger.Info("Uninstalling the controlplane...\n")

		args = []string{
			"uninstall",
			helmReleaseName,
			"--namespace",
			c.namespace,
		}
	} else {
		input.Logger.Info("Installing the controlplane in quickstart mode...\n")

		args = []string{
			"upgrade",
			"--install",
			helmReleaseName,
			helmChartRepoName,
			"--version",
			c.version,
			"--namespace",
			c.namespace,
			"--create-namespace",
			"--values",
			fmt.Sprintf(helmQuickstartValueRemotePath, c.version),
		}
	}

	var stderr, stdout bytes.Buffer
	cmd := exec.CommandContext(ctx, helm, args...)
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%w: %s", err, stderr.String())
	}

	input.Logger.Info(stdout.String())

	return nil
}

// getHelm finds and returns helm executable binary in the following priority:
//   1. pre-installed in command specified toolsDir (default is $HOME/.pipectl/tools)
//   2. $PATH
//   3. install new helm to command specified toolsDir
func (c *command) getHelm() (string, error) {
	fi, err := os.Stat(path.Join(c.toolsDir, helmBinaryName))
	if err != nil && !os.IsNotExist(err) {
		return "", err
	}

	// If the Helm executable binary exists in tools dir, use it.
	if fi != nil {
		return path.Join(c.toolsDir, helmBinaryName), nil
	}

	// If the Helm executable binary exists in $PATH, use it.
	path, err := exec.LookPath(helmBinaryName)
	if err != nil && !errors.Is(err, exec.ErrNotFound) {
		return "", err
	}

	if path != "" {
		return path, nil
	}

	// TODO: Install helm to command toolsDir.

	return path, nil
}
