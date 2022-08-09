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

	toolsDir          string
	values            string
	configFile        string
	encryptionKeyFile string

	firestoreServiceAccount string
	gcsServiceAccount       string
	cloudSQLServiceAccount  string
	minioAccessKey          string
	minioSecretKey          string
	internalTLSKey          string
	internalTLSCert         string
}

func newInstallControlplaneCommand() *cobra.Command {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Sprintf("failed to detect the current user's home directory: %v", err))
	}

	c := &controlplane{
		version:   version.Get().Version,
		namespace: pipecdDefaultNamespace,
		toolsDir:  path.Join(home, ".pipectl", "tools"),
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
	cmd.Flags().StringVar(&c.encryptionKeyFile, "encryption-key-file", c.encryptionKeyFile, "The path to the Control Plane encryption key file.")
	cmd.Flags().StringVar(&c.values, "values", c.values, "The Helm chart '--values' flag, which specify values in a YAML file or a URL (can specify multiple).")

	cmd.Flags().StringVar(&c.firestoreServiceAccount, "firestore-service-account-file", c.firestoreServiceAccount, "The path to service account which used to access controlplane Firestore database (if using).")
	cmd.Flags().StringVar(&c.cloudSQLServiceAccount, "cloud-sql-service-account-file", c.cloudSQLServiceAccount, "The path to service account which used to access controlplane Google cloud SQL database (if using).")
	cmd.Flags().StringVar(&c.gcsServiceAccount, "gcs-service-account-file", c.gcsServiceAccount, "The path to service account which used to access controlplane Google Cloud Storage service (if using).")
	cmd.Flags().StringVar(&c.minioAccessKey, "minio-access-key-file", c.minioAccessKey, "The path to access key which used to connect Minio filestore (if using).")
	cmd.Flags().StringVar(&c.minioSecretKey, "minio-secret-key-file", c.minioSecretKey, "The path to secret key which used to connect Minio filestore (if using).")
	cmd.Flags().StringVar(&c.internalTLSKey, "internal-tls-key-file", c.internalTLSKey, "The path to internal TLS key file (if using).")
	cmd.Flags().StringVar(&c.internalTLSCert, "internal-tls-cert-file", c.internalTLSCert, "The path to internal TLS certificate file (if using).")

	return cmd
}

func (c *controlplane) run(ctx context.Context, input cli.Input) error {
	helm, err := c.findHelm()
	if err != nil {
		return fmt.Errorf("failed to check required tools before install: %v", err)
	}

	input.Logger.Info("Installing the controlplane's components...")

	args, err := c.buildHelmArgs()
	if err != nil {
		return fmt.Errorf("failed to install controlplane: %v", err)
	}

	var stderr, stdout bytes.Buffer
	cmd := exec.CommandContext(ctx, helm, args...)
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%w: %s", err, stderr.String())
	}

	input.Logger.Info(stdout.String())
	input.Logger.Info("Installed the controlplane's components successfully.")

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

func (c *controlplane) buildHelmArgs() ([]string, error) {
	args := []string{
		"upgrade",
		"--install",
		helmReleaseName,
		helmChartRepoName,
		"--version",
		c.version,
		"--namespace",
		c.namespace,
		"--create-namespace",
	}

	if c.quickstart {
		args = append(args,
			"--values",
			fmt.Sprintf(helmQuickstartValueRemotePath, c.version),
		)
		return args, nil
	}

	if c.configFile == "" || c.encryptionKeyFile == "" {
		return nil, fmt.Errorf("missing required fields: config-file or encryption-key-file")
	}

	args = append(args,
		"--set-file",
		fmt.Sprintf("config.data=%s", c.configFile),
		"--set-file",
		fmt.Sprintf("secret.encryptionKey.data=%s", c.encryptionKeyFile),
	)

	if c.firestoreServiceAccount != "" {
		args = append(args,
			"--set-file",
			fmt.Sprintf("secret.firestoreServiceAccount.data=%s", c.firestoreServiceAccount),
		)
	}

	if c.cloudSQLServiceAccount != "" {
		args = append(args,
			"--set-file",
			fmt.Sprintf("secret.cloudSQLServiceAccount.data=%s", c.cloudSQLServiceAccount),
		)
	}

	if c.gcsServiceAccount != "" {
		args = append(args,
			"--set-file",
			fmt.Sprintf("secret.gcsServiceAccount.data=%s", c.gcsServiceAccount),
		)
	}

	if c.minioAccessKey != "" {
		args = append(args,
			"--set-file",
			fmt.Sprintf("secret.minioAccessKey.data=%s", c.minioAccessKey),
		)
	}

	if c.minioSecretKey != "" {
		args = append(args,
			"--set-file",
			fmt.Sprintf("secret.minioSecretKey.data=%s", c.minioSecretKey),
		)
	}

	if c.internalTLSKey != "" {
		args = append(args,
			"--set-file",
			fmt.Sprintf("secret.internalTLSKey.data=%s", c.internalTLSKey),
		)
	}

	if c.internalTLSCert != "" {
		args = append(args,
			"--set-file",
			fmt.Sprintf("secret.internalTLSCert.data=%s", c.internalTLSCert),
		)
	}

	if c.values != "" {
		args = append(args, "--values", c.values)
	}

	return args, nil
}
