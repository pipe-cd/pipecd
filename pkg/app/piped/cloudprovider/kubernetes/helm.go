// Copyright 2020 The PipeCD Authors.
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
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipe/pkg/config"
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

func (c *Helm) TemplateLocalChart(ctx context.Context, appName, appDir, chartPath string, opts *config.InputHelmOptions) (string, error) {
	releaseName := appName
	if opts != nil && opts.ReleaseName != "" {
		releaseName = opts.ReleaseName
	}

	args := []string{
		"template",
		"--no-hooks",
		releaseName,
		chartPath,
	}
	if opts != nil {
		for _, v := range opts.ValueFiles {
			args = append(args, "-f", v)
		}
		for k, v := range opts.SetFiles {
			args = append(args, "--set-file", fmt.Sprintf("%s=%s", k, v))
		}
	}

	var stdout, stderr bytes.Buffer
	cmd := exec.CommandContext(ctx, c.execPath, args...)
	cmd.Dir = appDir
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	c.logger.Info(fmt.Sprintf("start templating a local chart (or cloned remote git chart) for application %s", appName),
		zap.Any("args", args),
	)
	if err := cmd.Start(); err != nil {
		return "", err
	}

	if err := cmd.Wait(); err != nil {
		return stdout.String(), fmt.Errorf("%w: %s", err, stderr.String())
	}

	return stdout.String(), nil
}

type helmRemoteGitChart struct {
	GitRemote string
	Ref       string
	Path      string
}

func (c *Helm) TemplateRemoteGitChart(ctx context.Context, appName, appDir string, chart helmRemoteGitChart, gitClient gitClient, opts *config.InputHelmOptions) (string, error) {
	// Firstly, we need to download the remote repositoy.
	repoDir, err := ioutil.TempDir("", "helm-remote-chart")
	if err != nil {
		return "", fmt.Errorf("unabled to created temporary directory for storing remote helm chart: %w", err)
	}
	defer os.RemoveAll(repoDir)

	repo, err := gitClient.Clone(ctx, chart.GitRemote, chart.GitRemote, "", repoDir)
	if err != nil {
		return "", fmt.Errorf("unabled to clone git repository containing remote helm chart: %w", err)
	}

	if chart.Ref != "" {
		if err := repo.Checkout(ctx, chart.Ref); err != nil {
			return "", fmt.Errorf("unabled to checkout to specified ref %s: %w", chart.Ref, err)
		}
	}
	chartPath := filepath.Join(repoDir, chart.Path)

	// After that handle it as a local chart.
	return c.TemplateLocalChart(ctx, appName, appDir, chartPath, opts)
}

type helmRemoteChart struct {
	Repository string
	Name       string
	Version    string
}

func (c *Helm) TemplateRemoteChart(ctx context.Context, appName, appDir string, chart helmRemoteChart, opts *config.InputHelmOptions) (string, error) {
	return "", nil
}
