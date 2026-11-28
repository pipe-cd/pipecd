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

package kubernetes

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

	"github.com/pipe-cd/pipecd/pkg/app/piped/chartrepo"
	"github.com/pipe-cd/pipecd/pkg/app/piped/toolregistry"
	"github.com/pipe-cd/pipecd/pkg/config"
)

var (
	allowedURLSchemes = []string{"http", "https"}
)

type Helm struct {
	version  string
	execPath string
	logger   *zap.Logger
}

func NewHelm(version, path string, logger *zap.Logger) *Helm {
	return &Helm{
		version:  version,
		execPath: path,
		logger:   logger,
	}
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

type helmRemoteGitChart struct {
	GitRemote string
	Ref       string
	Path      string
}

func (h *Helm) TemplateRemoteGitChart(ctx context.Context, appName, appDir, namespace string, chart helmRemoteGitChart, gitClient gitClient, opts *config.InputHelmOptions) (string, error) {
	// Firstly, we need to download the remote repository.
	repoDir, err := os.MkdirTemp("", "helm-remote-chart")
	if err != nil {
		return "", fmt.Errorf("unable to created temporary directory for storing remote helm chart: %w", err)
	}
	defer os.RemoveAll(repoDir)

	repo, err := gitClient.Clone(ctx, chart.GitRemote, chart.GitRemote, "", repoDir)
	if err != nil {
		return "", fmt.Errorf("unable to clone git repository containing remote helm chart: %w", err)
	}

	if chart.Ref != "" {
		if err := repo.Checkout(ctx, chart.Ref); err != nil {
			return "", fmt.Errorf("unable to checkout to specified ref %s: %w", chart.Ref, err)
		}
	}
	chartPath := filepath.Join(repoDir, chart.Path)

	// After that handle it as a local chart.
	return h.TemplateLocalChart(ctx, appName, appDir, namespace, chartPath, opts)
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
	if e := chartrepo.Update(ctx, toolregistry.DefaultRegistry(), h.logger); e != nil {
		h.logger.Error("failed to update Helm chart repositories", zap.Error(e))
		return "", err
	}
	return executor()
}

// verifyHelmValueFilePath verifies if the path of the values file references
// a remote URL or inside the path where the application configuration file (i.e. *.pipecd.yaml) is located.
func verifyHelmValueFilePath(appDir, valueFilePath string) error {
	isRooted := filepath.IsAbs(valueFilePath) || strings.HasPrefix(valueFilePath, "/") || strings.HasPrefix(valueFilePath, "\\")
	if !isRooted {
		u, err := url.Parse(valueFilePath)
		if err == nil && u.Scheme != "" {
			for _, s := range allowedURLSchemes {
				if strings.EqualFold(u.Scheme, s) {
					return nil
				}
			}

			return fmt.Errorf("scheme %s is not allowed to load values file", u.Scheme)
		}
	}

	absAppDir, err := filepath.Abs(appDir)
	if err != nil {
		return err
	}

	// valueFilePath is a path where non-default Helm values file is located.
	if !isRooted {
		valueFilePath = filepath.Join(absAppDir, valueFilePath)
	}
	absValueFilePath, err := filepath.Abs(valueFilePath)
	if err != nil {
		return err
	}

	if isOutsideDir(absAppDir, absValueFilePath) {
		return fmt.Errorf("values file %s references outside the application configuration directory", absValueFilePath)
	}

	resolvedAppDir, err := resolveSymlinks(absAppDir)
	if err != nil {
		return err
	}
	resolvedValueFilePath, err := resolveSymlinks(absValueFilePath)
	if err != nil {
		return err
	}

	// If a path outside of appDir is specified as the path for the values file,
	// it may indicate that someone trying to illegally read a file as values file that
	// exists in the environment where Piped is running.
	if isOutsideDir(resolvedAppDir, resolvedValueFilePath) {
		return fmt.Errorf("values file %s references outside the application configuration directory", resolvedValueFilePath)
	}

	return nil
}

func isOutsideDir(baseDir, targetPath string) bool {
	rel, err := filepath.Rel(baseDir, targetPath)
	if err != nil {
		return true
	}
	return rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator))
}

// isSymlink returns the path is whether symbolic link or not.
func isSymlink(path string) bool {
	lstat, err := os.Lstat(path)
	if err != nil {
		return false
	}

	return lstat.Mode()&os.ModeSymlink == os.ModeSymlink
}

// resolveSymlinks resolves symlinks in path, including dangling symlinks and symlinked ancestors.
func resolveSymlinks(path string) (string, error) {
	resolved, err := filepath.EvalSymlinks(path)
	if err == nil {
		return resolved, nil
	}
	if !errors.Is(err, os.ErrNotExist) {
		return "", err
	}

	if isSymlink(path) {
		target, err := os.Readlink(path)
		if err != nil {
			return "", err
		}
		if !filepath.IsAbs(target) {
			parent, err := resolveSymlinks(filepath.Dir(path))
			if err != nil {
				return "", err
			}
			target = filepath.Join(parent, target)
		}
		return resolveSymlinks(filepath.Clean(target))
	}

	parent := filepath.Dir(path)
	if parent == path {
		return path, nil
	}
	resolvedParent, err := resolveSymlinks(parent)
	if err != nil {
		return "", err
	}
	return filepath.Join(resolvedParent, filepath.Base(path)), nil
}

func (h *Helm) UpgradeLocalChart(ctx context.Context, appName, appDir, namespace, chartPath string, opts *config.InputHelmOptions) (string, error) {
	releaseName := appName
	if opts != nil && opts.ReleaseName != "" {
		releaseName = opts.ReleaseName
	}

	args := []string{
		"upgrade",
		"--install",
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
	}
	var stdout, stderr bytes.Buffer
	cmd := exec.CommandContext(ctx, h.execPath, args...)
	cmd.Dir = appDir
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	h.logger.Info(fmt.Sprintf("start upgrading a release (or cloned remote git chart) for application %s", appName),
		zap.Any("args", args),
	)

	if err := cmd.Run(); err != nil {
		return stdout.String(), fmt.Errorf("%w: %s", err, stderr.String())
	}
	return stdout.String(), nil
}

func (h *Helm) UpgradeRemoteGitChart(ctx context.Context, appName, appDir, namespace string, chart helmRemoteGitChart, gitClient gitClient, opts *config.InputHelmOptions) (string, error) {
	repoDir, err := os.MkdirTemp("", "helm-remote-chart")
	if err != nil {
		return "", fmt.Errorf("unable to created temporary directory for storing remote helm chart: %w", err)
	}
	defer os.RemoveAll(repoDir)

	repo, err := gitClient.Clone(ctx, chart.GitRemote, chart.GitRemote, "", repoDir)
	if err != nil {
		return "", fmt.Errorf("unable to clone git repository containing remote helm chart: %w", err)
	}

	if chart.Ref != "" {
		if err := repo.Checkout(ctx, chart.Ref); err != nil {
			return "", fmt.Errorf("unable to checkout to specified ref %s: %w", chart.Ref, err)
		}
	}
	chartPath := filepath.Join(repoDir, chart.Path)

	// After that handle it as a local chart.
	return h.UpgradeLocalChart(ctx, appName, appDir, namespace, chartPath, opts)
}

func (h *Helm) UpgradeRemoteChart(ctx context.Context, appName, appDir, namespace string, chart helmRemoteChart, gitClient gitClient, opts *config.InputHelmOptions) (string, error) {
	releaseName := appName
	if opts != nil && opts.ReleaseName != "" {
		releaseName = opts.ReleaseName
	}

	args := []string{
		"upgrade",
		"--install",
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
	}

	h.logger.Info(fmt.Sprintf("start upgrading a release from Helm repository for application %s", appName),
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
	if e := chartrepo.Update(ctx, toolregistry.DefaultRegistry(), h.logger); e != nil {
		h.logger.Error("failed to update Helm chart repositories", zap.Error(e))
		return "", err
	}
	return executor()

}
