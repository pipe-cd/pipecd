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

package provider

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/config"
)

var (
	allowedURLSchemes = []string{"http", "https"}

	updateGroup = &singleflight.Group{}
)

type Helm struct {
	execPath string
	logger   *zap.Logger
}

func NewHelm(path string, logger *zap.Logger) *Helm {
	return &Helm{
		execPath: path,
		logger:   logger,
	}
}

func (h *Helm) LoginToOCIRegistry(ctx context.Context, address, username, password string) error {
	args := []string{
		"registry",
		"login",
		"-u",
		username,
		"-p",
		password,
		address,
	}

	var stderr bytes.Buffer
	cmd := exec.CommandContext(ctx, h.execPath, args...)
	cmd.Stderr = &stderr

	h.logger.Info("login to oci registry", zap.String("address", address))
	if err := cmd.Run(); err != nil {
		h.logger.Error("failed to login to oci registry", zap.String("address", address), zap.Error(err))
		return fmt.Errorf("%w: %s", err, stderr.String())
	}
	return nil
}

func (h *Helm) TemplateLocalChart(ctx context.Context, appName, appDir, namespace, chartPath string, opts *config.InputHelmOptions) (string, error) {
	releaseName := appName
	if opts != nil && opts.ReleaseName != "" {
		releaseName = opts.ReleaseName
	}

	args := []string{
		"template",
		"--no-hooks",
		"--include-crds",
		"--dependency-update",
		releaseName,
		chartPath,
	}

	if namespace != "" {
		args = append(args, fmt.Sprintf("--namespace=%s", namespace))
	}

	if opts != nil {
		for k, v := range opts.SetValues {
			args = append(args, "--set", fmt.Sprintf("%s=%s", k, v))
		}
		for _, v := range opts.ValueFiles {
			if err := verifyHelmValueFilePath(appDir, v); err != nil {
				h.logger.Error("failed to verify values file path", zap.Error(err))
				return "", err
			}
			args = append(args, "-f", v)
		}
		for k, v := range opts.SetFiles {
			args = append(args, "--set-file", fmt.Sprintf("%s=%s", k, v))
		}
		for _, v := range opts.APIVersions {
			args = append(args, "--api-versions", v)
		}
		if opts.KubeVersion != "" {
			args = append(args, "--kube-version", opts.KubeVersion)
		}
	}

	var stdout, stderr bytes.Buffer
	cmd := exec.CommandContext(ctx, h.execPath, args...)
	cmd.Dir = appDir
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	h.logger.Info(fmt.Sprintf("start templating a local chart (or cloned remote git chart) for application %s", appName),
		zap.Any("args", args),
	)

	if err := cmd.Run(); err != nil {
		return stdout.String(), fmt.Errorf("%w: %s", err, stderr.String())
	}
	return stdout.String(), nil
}

// Add installs all specified Helm Chart repositories.
// https://helm.sh/docs/topics/chart_repository/
// helm repo add fantastic-charts https://fantastic-charts.storage.googleapis.com
// helm repo add fantastic-charts https://fantastic-charts.storage.googleapis.com --username my-username --password my-password
func (h *Helm) AddRepository(ctx context.Context, repo config.HelmChartRepository) error {
	args := []string{"repo", "add", repo.Name, repo.Address}
	if repo.Insecure {
		args = append(args, "--insecure-skip-tls-verify")
	}
	if repo.Username != "" || repo.Password != "" {
		args = append(args, "--username", repo.Username, "--password", repo.Password)
	}
	cmd := exec.CommandContext(ctx, h.execPath, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		h.logger.Error("failed to add chart repository", zap.String("name", repo.Name), zap.Error(err))
		return fmt.Errorf("failed to add chart repository %s: %s (%w)", repo.Name, string(out), err)
	}
	h.logger.Info("successfully added chart repository", zap.String("name", repo.Name))
	return nil
}

func (h *Helm) UpdateRepositories(ctx context.Context) error {
	_, err, _ := updateGroup.Do("update", func() (interface{}, error) {
		args := []string{"repo", "update"}
		cmd := exec.CommandContext(ctx, h.execPath, args...)
		out, err := cmd.CombinedOutput()
		if err != nil {
			h.logger.Error("failed to update chart repositories", zap.Error(err))
			return nil, fmt.Errorf("failed to update chart repositories: %s (%w)", string(out), err)
		}
		h.logger.Info("successfully updated chart repositories")
		return nil, nil
	})
	return err
}

// verifyHelmValueFilePath verifies if the path of the values file references
// a remote URL or inside the path where the application configuration file (i.e. *.pipecd.yaml) is located.
func verifyHelmValueFilePath(appDir, valueFilePath string) error {
	url, err := url.Parse(valueFilePath)
	if err == nil && url.Scheme != "" {
		for _, s := range allowedURLSchemes {
			if strings.EqualFold(url.Scheme, s) {
				return nil
			}
		}

		return fmt.Errorf("scheme %s is not allowed to load values file", url.Scheme)
	}

	// valueFilePath is a path where non-default Helm values file is located.
	if !filepath.IsAbs(valueFilePath) {
		valueFilePath = filepath.Join(appDir, valueFilePath)
	}

	if isSymlink(valueFilePath) {
		if valueFilePath, err = resolveSymlinkToAbsPath(valueFilePath, appDir); err != nil {
			return err
		}
	}

	// If a path outside of appDir is specified as the path for the values file,
	// it may indicate that someone trying to illegally read a file as values file that
	// exists in the environment where Piped is running.
	if !strings.HasPrefix(valueFilePath, appDir) {
		return fmt.Errorf("values file %s references outside the application configuration directory", valueFilePath)
	}

	return nil
}

// isSymlink returns the path is whether symbolic link or not.
func isSymlink(path string) bool {
	lstat, err := os.Lstat(path)
	if err != nil {
		return false
	}

	return lstat.Mode()&os.ModeSymlink == os.ModeSymlink
}

// resolveSymlinkToAbsPath resolves symbolic link to an absolute path.
func resolveSymlinkToAbsPath(path, absParentDir string) (string, error) {
	resolved, err := os.Readlink(path)
	if err != nil {
		return "", err
	}

	if !filepath.IsAbs(resolved) {
		resolved = filepath.Join(absParentDir, resolved)
	}

	return resolved, nil
}

type helmRemoteChart struct {
	Repository string
	Name       string
	Version    string
	Insecure   bool
}

func (h *Helm) TemplateRemoteChart(ctx context.Context, appName, appDir, namespace string, chart helmRemoteChart, opts *config.InputHelmOptions) (string, error) {
	releaseName := appName
	if opts != nil && opts.ReleaseName != "" {
		releaseName = opts.ReleaseName
	}

	args := []string{
		"template",
		"--no-hooks",
		"--include-crds",
		"--dependency-update",
		releaseName,
		fmt.Sprintf("%s/%s", chart.Repository, chart.Name),
		fmt.Sprintf("--version=%s", chart.Version),
	}

	if chart.Insecure {
		args = append(args, "--insecure-skip-tls-verify")
	}

	if namespace != "" {
		args = append(args, fmt.Sprintf("--namespace=%s", namespace))
	}

	if opts != nil {
		for k, v := range opts.SetValues {
			args = append(args, "--set", fmt.Sprintf("%s=%s", k, v))
		}
		for _, v := range opts.ValueFiles {
			if err := verifyHelmValueFilePath(appDir, v); err != nil {
				h.logger.Error("failed to verify values file path", zap.Error(err))
				return "", err
			}
			args = append(args, "-f", v)
		}
		for k, v := range opts.SetFiles {
			args = append(args, "--set-file", fmt.Sprintf("%s=%s", k, v))
		}
		for _, v := range opts.APIVersions {
			args = append(args, "--api-versions", v)
		}
		if opts.KubeVersion != "" {
			args = append(args, "--kube-version", opts.KubeVersion)
		}
	}

	h.logger.Info(fmt.Sprintf("start templating a chart from Helm repository for application %s", appName),
		zap.Any("args", args),
	)

	executor := func() (string, error) {
		var stdout, stderr bytes.Buffer
		cmd := exec.CommandContext(ctx, h.execPath, args...)
		cmd.Dir = appDir
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		if err := cmd.Run(); err != nil {
			return stdout.String(), fmt.Errorf("%w: %s", err, stderr.String())
		}
		return stdout.String(), nil
	}

	out, err := executor()
	if err == nil {
		return out, nil
	}

	if !strings.Contains(err.Error(), "helm repo update") {
		return "", err
	}

	// If the error is a "Not Found", we update the repositories and try again.
	if e := h.UpdateRepositories(ctx); e != nil {
		h.logger.Error("failed to update Helm chart repositories", zap.Error(e))
		return "", errors.Join(err, e)
	}
	return executor()
}
