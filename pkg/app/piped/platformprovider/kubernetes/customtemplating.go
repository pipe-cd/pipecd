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
	"os/exec"

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

func (c *CustomTemplating) Template(ctx context.Context, appName, appDir string, args []string) (string, error) {
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
