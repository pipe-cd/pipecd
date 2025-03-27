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
	"fmt"
	"os/exec"

	"go.uber.org/zap"
	"golang.org/x/mod/semver"
)

type Kustomize struct {
	version  string
	execPath string
	logger   *zap.Logger
}

func NewKustomize(version, path string, logger *zap.Logger) *Kustomize {
	return &Kustomize{
		version:  version,
		execPath: path,
		logger:   logger,
	}
}

func (c *Kustomize) Template(ctx context.Context, appName, appDir string, opts map[string]string, helm *Helm) (string, error) {
	args := []string{
		"build",
		".",
	}

	// Pass the Helm command path to kustomize to use the specified version of Helm.
	// Unconditionally adding this flag as it's unharmful when Helm is not used.
	if c.isHelmCommandFlagAvailable() && helm != nil {
		args = append(args, "--helm-command", helm.execPath)
	}

	for k, v := range opts {
		args = append(args, fmt.Sprintf("--%s", k))
		if v != "" {
			args = append(args, v)
		}
	}

	var stdout, stderr bytes.Buffer
	cmd := exec.CommandContext(ctx, c.execPath, args...)
	cmd.Dir = appDir
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	c.logger.Info(fmt.Sprintf("start templating a Kustomize application %s", appName),
		zap.Any("args", args),
	)

	if err := cmd.Run(); err != nil {
		return stdout.String(), fmt.Errorf("%w: %s", err, stderr.String())
	}
	return stdout.String(), nil
}

// isHelmCommandFlagAvailable returns true if the `--helm-command` flag is available
// on the installed Kustomize version.
func (c *Kustomize) isHelmCommandFlagAvailable() bool {
	// It's only available on Kustomize v4.1.0 and higher.
	return semver.Compare("v"+c.version, "v4.1.0") >= 0
}
