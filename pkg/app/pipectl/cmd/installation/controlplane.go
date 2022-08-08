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

package installation

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/spf13/cobra"

	"github.com/pipe-cd/pipecd/pkg/cli"
	"github.com/pipe-cd/pipecd/pkg/version"
)

const (
	pipecdDefaultNamespace = "default"

	helmBinaryName    = "helm"
	helmReleaseName   = "pipecd"
	helmChartRepoName = "oci://ghcr.io/pipe-cd/chart/pipecd"

	helmQuickstartValueRemotePath = "https://raw.githubusercontent.com/pipe-cd/pipecd/%s/quickstart/control-plane-values.yaml"
)

type controlplane struct {
	quickstart bool
	version    string
	namespace  string

	toolsDir   string
	configFile string
}

func newInstallControlplaneCommand() *cobra.Command {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Sprintf("failed to detect the current user's home directory: %v", err))
	}

	c := &controlplane{
		version:   version.Get().Version,
		namespace: pipecdDefaultNamespace,
		toolsDir:  path.Join(home, ".piped", "tools"),
	}

	cmd := &cobra.Command{
		Use:   "controlplane",
		Short: "Install PipeCD Control Plane.",
		RunE:  cli.WithContext(c.run),
	}

	cmd.Flags().BoolVar(&c.quickstart, "quickstart", c.quickstart, "Whether installing the Control Plane in quickstart mode.")
	cmd.Flags().StringVar(&c.version, "version", c.version, "The Control Plane version. Default is the version of pipectl.")
	cmd.Flags().StringVar(&c.namespace, "namespace", c.namespace, "The cluster namespace where to install Control Plane. Default is 'default'.")

	cmd.Flags().StringVar(&c.toolsDir, "tools-dir", c.toolsDir, "The path to directory where to install needed tools such as helm.")
	cmd.Flags().StringVar(&c.configFile, "config-file", c.configFile, "The path to the Control Plane configuration file.")

	return cmd
}

func (c *controlplane) run(ctx context.Context, input cli.Input) error {
	helm, err := c.findHelm()
	if err != nil {
		return fmt.Errorf("failed to check required tools before install: %v", err)
	}

	input.Logger.Info("Installing the controlplane's components")

	var args []string
	if c.quickstart {
		args = c.buildHelmArgsForQuickstart()
	} else {
		args = []string{}
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

func (c *controlplane) findHelm() (string, error) {
	fi, err := os.Stat(path.Join(c.toolsDir, helmBinaryName))
	if err != nil && !os.IsNotExist(err) {
		return "", err
	}

	// If the Helm executable binary exists in tools dir, use it.
	if fi != nil {
		return path.Join(c.toolsDir, helmBinaryName), nil
	}

	// If the Helm executable binary exists in $PATH, use it.
	path, err := exec.LookPath("helm")
	if err != nil {
		return "", err
	}
	return path, nil
}

func (c *controlplane) buildHelmArgsForQuickstart() []string {
	return []string{
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
