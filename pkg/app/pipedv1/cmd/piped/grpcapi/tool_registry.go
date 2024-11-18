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

package grpcapi

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"text/template"

	"golang.org/x/sync/singleflight"
)

type toolRegistry struct {
	toolsDir string
	tmpDir   string
	group    singleflight.Group
}

type templateValues struct {
	Name    string
	Version string
	OutPath string
	TmpDir  string
	Arch    string
	Os      string
}

func newToolRegistry(toolsDir, tmpDir string) (*toolRegistry, error) {
	tmpDir, err := os.MkdirTemp(tmpDir, "tool-registry")
	if err != nil {
		return nil, fmt.Errorf("failed to create a temporary directory: %w", err)
	}
	return &toolRegistry{
		toolsDir: toolsDir,
		tmpDir:   tmpDir,
	}, nil
}

func (r *toolRegistry) newTmpDir() (string, error) {
	return os.MkdirTemp(r.tmpDir, "")
}

func (r *toolRegistry) outPath() (string, error) {
	target, err := r.newTmpDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(target, "out"), nil
}

func (r *toolRegistry) InstallTool(ctx context.Context, name, version, script string) (string, error) {
	out, err, _ := r.group.Do(fmt.Sprintf("%s-%s", name, version), func() (interface{}, error) {
		return r.installTool(ctx, name, version, script)
	})
	if err != nil {
		return "", fmt.Errorf("failed to install the tool %s-%s: %w", name, version, err)
	}
	return out.(string), nil // the result is always string
}

func (r *toolRegistry) installTool(ctx context.Context, name, version, script string) (path string, err error) {
	outPath, err := r.outPath()
	if err != nil {
		return "", err
	}

	tmpDir, err := r.newTmpDir()
	if err != nil {
		return "", err
	}

	t, err := template.New("install script").Parse(script)
	if err != nil {
		return "", err
	}

	vars := templateValues{
		Name:    name,
		Version: version,
		OutPath: outPath,
		TmpDir:  tmpDir,
		Arch:    runtime.GOARCH,
		Os:      runtime.GOOS,
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, vars); err != nil {
		return "", err
	}

	cmd := exec.CommandContext(ctx, "/bin/sh", "-c", buf.String())
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to execute the install script: %w, output: %s", err, string(out))
	}

	if err := os.Chmod(outPath, 0o755); err != nil {
		return "", err
	}

	target := filepath.Join(r.toolsDir, fmt.Sprintf("%s-%s", name, version))
	if out, err := exec.CommandContext(ctx, "/bin/sh", "-c", "mv "+outPath+" "+target).CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to move the installed binary: %w, output: %s", err, string(out))
	}

	if err := os.RemoveAll(tmpDir); err != nil {
		return "", err
	}

	return target, nil
}

func (r *toolRegistry) Close() error {
	if err := os.RemoveAll(r.tmpDir); err != nil {
		return err
	}
	return nil
}
