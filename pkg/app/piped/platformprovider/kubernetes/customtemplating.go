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

package kubernetes

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"os/exec"
	"path/filepath"
	"strings"

	"go.uber.org/zap"
)

type CustomTemplating struct {
	execPath string
	logger   *zap.Logger
}

func NewCustomTemplating(path string, logger *zap.Logger) *CustomTemplating {
	return &CustomTemplating{
		execPath: path,
		logger:   logger,
	}
}

func (c *CustomTemplating) Template(ctx context.Context, appName, appDir string, inputArgs []string) (string, error) {
	var args []string
	for _, v := range inputArgs {
		args = append(args, strings.Split(v, " ")...)
	}

	for _, v := range args {
		if err := verifyCustomtemplatingArgs(appDir, v); err != nil {
			c.logger.Error("failed to verify args", zap.Error(err))
			return "", err
		}
	}
	var stdout, stderr bytes.Buffer
	cmd := exec.CommandContext(ctx, c.execPath, args...)
	cmd.Dir = appDir
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	c.logger.Info(fmt.Sprintf("start custom templating %s application %s", c.execPath, appName),
		zap.Any("args", args),
	)

	if err := cmd.Run(); err != nil {
		return stdout.String(), fmt.Errorf("%w: %s", err, stderr.String())
	}
	return stdout.String(), nil
}

func verifyCustomtemplatingArgs(appDir, arg string) error {
	url, err := url.Parse(arg)
	if err == nil && url.Scheme != "" {
		for _, s := range allowedURLSchemes {
			if strings.EqualFold(url.Scheme, s) {
				return nil
			}
		}

		return fmt.Errorf("scheme %s is not allowed to use as args", url.Scheme)
	}

	// arg is a path where non-default file is located.
	if !filepath.IsAbs(arg) {
		arg = filepath.Join(appDir, arg)
	}

	if isSymlink(arg) {
		if arg, err = resolveSymlinkToAbsPath(arg, appDir); err != nil {
			return err
		}
	}

	// If a path outside of appDir is specified as the path in args,
	// it may indicate that someone trying to illegally read a file that
	// exists in the environment where Piped is running.
	if !strings.HasPrefix(arg, appDir) {
		return fmt.Errorf("arg %s references outside the application configuration directory", arg)
	}

	return nil
}
