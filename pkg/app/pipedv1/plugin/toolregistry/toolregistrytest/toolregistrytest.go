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

package toolregistrytest

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"runtime"
	"testing"
	"text/template"
)

type ToolRegistry struct {
	testingT *testing.T
	tmpDir   string
}

type templateValues struct {
	Name    string
	Version string
	OutPath string
	TmpDir  string
	Arch    string
	Os      string
}

func NewToolRegistry(t *testing.T) (*ToolRegistry, error) {
	tmpDir, err := os.MkdirTemp("", "tool-registry-test")
	if err != nil {
		return nil, err
	}
	return &ToolRegistry{
		testingT: t,
		tmpDir:   tmpDir,
	}, nil
}

func (r *ToolRegistry) newTmpDir() (string, error) {
	return os.MkdirTemp(r.tmpDir, "")
}

func (r *ToolRegistry) binDir() (string, error) {
	target := r.tmpDir + "/bin"
	if err := os.MkdirAll(target, 0o755); err != nil {
		return "", err
	}
	return target, nil
}

func (r *ToolRegistry) outPath() (string, error) {
	target, err := r.newTmpDir()
	if err != nil {
		return "", err
	}
	return target + "/out", nil
}

func (r *ToolRegistry) InstallTool(ctx context.Context, name, version, script string) (path string, err error) {
	outPath, err := r.outPath()
	if err != nil {
		return "", err
	}

	tmpDir, err := r.newTmpDir()
	if err != nil {
		return "", err
	}

	binDir, err := r.binDir()
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
		r.testingT.Log(string(out))
		return "", err
	}

	if err := os.Chmod(outPath, 0o755); err != nil {
		return "", err
	}

	target := binDir + "/" + name + "-" + version
	if out, err := exec.CommandContext(ctx, "/bin/sh", "-c", "mv "+outPath+" "+target).CombinedOutput(); err != nil {
		r.testingT.Log(string(out))
		return "", err
	}

	if err := os.RemoveAll(tmpDir); err != nil {
		return "", err
	}

	return target, nil
}

func (r *ToolRegistry) Close() error {
	if err := os.RemoveAll(r.tmpDir); err != nil {
		return err
	}
	return nil
}
