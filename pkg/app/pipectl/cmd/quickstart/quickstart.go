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

package quickstart

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/manifoldco/promptui"
	"github.com/pipe-cd/pipecd/pkg/backoff"
	"github.com/pipe-cd/pipecd/pkg/cli"
	"github.com/pipe-cd/pipecd/pkg/version"
	"github.com/spf13/cobra"
)

const (
	defaultHelmVersion = "3.8.2"

	helmControlPlaneReleaseName = "pipecd"
	helmPipedReleaseName        = "piped"

	helmChartControlPlaneRepoName = "oci://ghcr.io/pipe-cd/chart/pipecd"
	helmChartPipedRepoName        = "oci://ghcr.io/pipe-cd/chart/piped"

	helmQuickstartValueRemotePath = "https://raw.githubusercontent.com/pipe-cd/pipecd/%s/quickstart/control-plane-values.yaml"

	pipecdDefaultNamespace = "pipecd"

	controlPlaneLocalhost = "http://localhost:8080/settings/piped?project=quickstart"

	pipedIDLabel       = "ID"
	pipedKeyLabel      = "Key"
	pipedGitRemoteRepo = "GitRemoteRepo"

	deploymentReadyRetryTime     = 3
	deploymentReadyRetryDuration = time.Minute
	deploymentReadyCheckDuration = 5 * time.Second
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
		panic(fmt.Sprintf("failed to detect the current user's home directory: %v", err))
	}

	defaultToolsDir := path.Join(home, ".pipectl", "tools")
	if err = os.MkdirAll(defaultToolsDir, 0755); err != nil {
		panic(fmt.Sprintf("failed to prepare tools dir: %v", err))
	}

	c := &command{
		version:   version.Get().Version,
		toolsDir:  defaultToolsDir,
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
	helm, err := c.getHelm(ctx)
	if err != nil {
		return fmt.Errorf("failed to prepare required tools (helm) for installation: %v", err)
	}

	if c.uninstall {
		return c.uninstallAll(ctx, helm, input)
	}

	if err = c.installControlPlane(ctx, helm, input); err != nil {
		input.Logger.Error("Failed to install PipeCD control plane!!")
		return err
	}

	var wg sync.WaitGroup
	if err = c.exposeService(ctx, &wg, input); err != nil {
		input.Logger.Error("Failed to expose PipeCD control plane service!!")
		return err
	}

	if err = c.installPiped(ctx, helm, input); err != nil {
		input.Logger.Error("Failed to install piped!!")
		return err
	}

	input.Logger.Info("\nPipeCD console is ready at http://localhost:8080/")

	// Wait until users hit SIG_KILL.
	wg.Wait()

	return nil
}

func (c *command) installControlPlane(ctx context.Context, helm string, input cli.Input) error {
	input.Logger.Info("Installing the controlplane in quickstart mode...")

	args := []string{
		"upgrade",
		"--install",
		helmControlPlaneReleaseName,
		helmChartControlPlaneRepoName,
		"--version",
		c.version,
		"--namespace",
		c.namespace,
		"--create-namespace",
		"--values",
		fmt.Sprintf(helmQuickstartValueRemotePath, c.version),
		"--set",
		fmt.Sprintf("mysql.image=%s", selectMySQLImage()),
	}

	var stderr, stdout bytes.Buffer
	cmd := exec.CommandContext(ctx, helm, args...)
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%w: %s", err, stderr.String())
	}

	input.Logger.Info(stdout.String())
	input.Logger.Info("Intalled the controlplane successfully!")

	return nil
}

func (c *command) installPiped(ctx context.Context, helm string, input cli.Input) error {
	input.Logger.Info("\nInstalling the piped for quickstart...")

	input.Logger.Info("\nOpenning PipeCD control plane at http://localhost:8080/\nPlease login using the following account:\n- Username: hello-pipecd\n- Password: hello-pipecd\nFor more information refer to https://pipecd.dev/docs/quickstart/\n")

	if err := openbrowser(controlPlaneLocalhost); err != nil {
		return fmt.Errorf("failed to open PipeCD control plane: %w", err)
	}

	input.Logger.Info("Fill up your registered Piped information:")

	pipedID := getPromptInput(pipedIDLabel)
	pipedKey := getPromptInput(pipedKeyLabel)
	sourceRepo := getPromptInput(pipedGitRemoteRepo)

	args := []string{
		"upgrade",
		"--install",
		helmPipedReleaseName,
		helmChartPipedRepoName,
		"--version",
		c.version,
		"--namespace",
		c.namespace,
		"--set",
		"quickstart.enabled=true",
		"--set",
		fmt.Sprintf("quickstart.pipedId=%s", pipedID),
		"--set",
		fmt.Sprintf("secret.data.piped-key=%s", pipedKey),
		"--set",
		fmt.Sprintf("quickstart.gitRepoRemote=%s", sourceRepo),
	}

	var stderr, stdout bytes.Buffer
	cmd := exec.CommandContext(ctx, helm, args...)
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%w: %s", err, stderr.String())
	}

	input.Logger.Info(stdout.String())
	input.Logger.Info("Intalled the piped successfully!")

	return nil
}

func (c *command) printExposeState(ctx context.Context, input cli.Input) {
	binName := "bash"
	epath, err := exec.LookPath(binName)
	if err != nil {
		return
	}

	if epath != "" {
		binName = epath
	}
	kubectl, _ := c.getKubectl()
	cmdText := fmt.Sprintf("%s -n %s", kubectl, c.namespace) +
		` get pods --no-headers | awk '{print $3}' | sort | uniq -c | awk '{total+=$1; statuses[$2]=$1} END {for (status in statuses) printf " %s %s", status, statuses[status]}'`
	args := []string{
		"-c",
		cmdText,
	}
	cmd := exec.CommandContext(ctx, binName, args...)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Run()
	input.Logger.Sugar().Infof("PipeCD control plane status:%s", stdout.String())
}

func (c *command) exposeService(ctx context.Context, wg *sync.WaitGroup, input cli.Input) error {
	input.Logger.Info("\nWaiting for PipeCD control plane to be ready...")
	notify := make(chan struct{})
	go func() {
		ticker := time.NewTicker(deploymentReadyCheckDuration)
		for {
			select {
			case <-ticker.C:
				c.printExposeState(ctx, input)
			case <-notify:
				return
			}
		}
	}()
	defer close(notify)
	kubectl, err := c.getKubectl()
	if err != nil {
		return fmt.Errorf("failed to prepare required tool (kubectl) for installation: %v", err)
	}

	// Wait the control plane service to be ready.
	args := []string{
		"rollout",
		"status",
		"deploy/pipecd-server",
		"-n",
		c.namespace,
	}
	var stdout bytes.Buffer
	cmd := exec.CommandContext(ctx, kubectl, args...)
	cmd.Stdout = &stdout

	retry := backoff.NewRetry(deploymentReadyRetryTime, backoff.NewConstant(deploymentReadyRetryDuration))
	var serverIsReady bool
	for retry.WaitNext(ctx) {
		cmd.Run()
		if strings.Contains(stdout.String(), "successfully rolled out") {
			notify <- struct{}{}
			serverIsReady = true
			break
		}
	}

	if !serverIsReady {
		return fmt.Errorf("failed while waiting for server to be ready")
	}

	// Expose the PipeCD control plane to localhost:8080.
	args = []string{
		"port-forward",
		"svc/pipecd",
		"8080",
		"-n",
		c.namespace,
	}
	var stderr bytes.Buffer
	cmd = exec.CommandContext(ctx, kubectl, args...)
	cmd.Stderr = &stderr

	wg.Add(1)
	go func() {
		cmd.Run()
		defer wg.Done()
	}()

	return nil
}

func getPromptInput(label string) string {
	validate := func(input string) error {
		switch label {
		case pipedIDLabel:
			if len(input) != 36 {
				return fmt.Errorf("invalid ID")
			}
		case pipedKeyLabel:
			if len(input) != 50 {
				return fmt.Errorf("invalid Key")
			}
		default:
			if len(input) == 0 {
				return fmt.Errorf("missing value: %s", label)
			}
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    label,
		Validate: validate,
	}

	result, err := prompt.Run()
	if err != nil {
		return ""
	}
	return result
}

func openbrowser(url string) error {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
	return err
}

func selectMySQLImage() string {
	var mysqlImage string
	switch runtime.GOARCH {
	case "amd64":
		mysqlImage = "mysql"
	case "arm64":
		mysqlImage = "arm64v8/mysql"
	default:
		mysqlImage = "mysql"
	}
	return mysqlImage
}

func (c *command) uninstallAll(ctx context.Context, helm string, input cli.Input) error {
	input.Logger.Info("Uninstalling PipeCD components...")

	var stderr, stdout bytes.Buffer

	// Uninstall PipeCD control plane.
	args := []string{
		"uninstall",
		helmControlPlaneReleaseName,
		"--namespace",
		c.namespace,
	}

	cmd := exec.CommandContext(ctx, helm, args...)
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%w: %s", err, stderr.String())
	}

	// Uninstall Piped.
	args = []string{
		"uninstall",
		helmPipedReleaseName,
		"--namespace",
		c.namespace,
	}

	cmd = exec.CommandContext(ctx, helm, args...)
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%w: %s", err, stderr.String())
	}

	input.Logger.Info(stdout.String())
	input.Logger.Info("Uninstalled the PipeCD components successfully!")

	return nil
}

func (c *command) getKubectl() (string, error) {
	binName := "kubectl"

	fi, err := os.Stat(path.Join(c.toolsDir, binName))
	if err != nil && !os.IsNotExist(err) {
		return "", err
	}

	if fi != nil {
		return path.Join(c.toolsDir, binName), nil
	}

	epath, err := exec.LookPath(binName)
	if err != nil && !errors.Is(err, exec.ErrNotFound) {
		return "", err
	}

	if epath != "" {
		return epath, nil
	}

	return "", fmt.Errorf("%s not found", binName)
}

// getHelm finds and returns helm executable binary in the following priority:
//  1. pre-installed in command specified toolsDir (default is $HOME/.pipectl/tools)
//  2. $PATH
//  3. install new helm to command specified toolsDir
func (c *command) getHelm(ctx context.Context) (string, error) {
	binName := "helm"

	fi, err := os.Stat(path.Join(c.toolsDir, binName))
	if err != nil && !os.IsNotExist(err) {
		return "", err
	}

	// If the Helm executable binary exists in tools dir, use it.
	if fi != nil {
		return path.Join(c.toolsDir, binName), nil
	}

	// If the Helm executable binary exists in $PATH, use it.
	epath, err := exec.LookPath(binName)
	if err != nil && !errors.Is(err, exec.ErrNotFound) {
		return "", err
	}

	if epath != "" {
		return epath, nil
	}

	// Install helm to command toolsDir.
	helmInstallScriptTmpl := template.Must(template.New("helm").Parse(helmInstallScript))
	var (
		buf  bytes.Buffer
		data = map[string]interface{}{
			"Version": defaultHelmVersion,
			"BinDir":  c.toolsDir,
		}
	)
	if err := helmInstallScriptTmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to install helm %s (%v)", defaultHelmVersion, err)
	}

	var (
		script = buf.String()
		cmd    = exec.CommandContext(ctx, "/bin/sh", "-c", script)
	)
	if _, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to install helm %s (%v)", defaultHelmVersion, err)
	}

	return path.Join(c.toolsDir, binName), nil
}
