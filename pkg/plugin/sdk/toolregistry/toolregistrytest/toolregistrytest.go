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

	"github.com/pipe-cd/piped-plugin-sdk-go/toolregistry"
	"google.golang.org/grpc"

	"github.com/pipe-cd/pipecd/pkg/plugin/pipedservice"
)

type templateValues struct {
	Name    string
	Version string
	OutPath string
	TmpDir  string
	Arch    string
	Os      string
}

type fakeClient struct {
	pipedservice.PluginServiceClient
	testingT *testing.T
}

func (c *fakeClient) binDir() (string, error) {
	target := c.testingT.TempDir() + "/bin"
	if err := os.MkdirAll(target, 0o755); err != nil {
		return "", err
	}
	return target, nil
}

func (c *fakeClient) outPath() string {
	return c.testingT.TempDir() + "/out"
}

func (c *fakeClient) InstallTool(ctx context.Context, in *pipedservice.InstallToolRequest, opts ...grpc.CallOption) (*pipedservice.InstallToolResponse, error) {
	outPath := c.outPath()

	binDir, err := c.binDir()
	if err != nil {
		return nil, err
	}

	t, err := template.New("install script").Parse(in.GetInstallScript())
	if err != nil {
		return nil, err
	}

	vars := templateValues{
		Name:    in.GetName(),
		Version: in.GetVersion(),
		OutPath: outPath,
		TmpDir:  c.testingT.TempDir(),
		Arch:    runtime.GOARCH,
		Os:      runtime.GOOS,
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, vars); err != nil {
		return nil, err
	}

	cmd := exec.CommandContext(ctx, "/bin/sh", "-c", buf.String())
	if out, err := cmd.CombinedOutput(); err != nil {
		c.testingT.Log(string(out))
		return nil, err
	}

	if err := os.Chmod(outPath, 0o755); err != nil {
		return nil, err
	}

	target := binDir + "/" + in.GetName() + "-" + in.GetVersion()
	if out, err := exec.CommandContext(ctx, "/bin/sh", "-c", "mv "+outPath+" "+target).CombinedOutput(); err != nil {
		c.testingT.Log(string(out))
		return nil, err
	}

	return &pipedservice.InstallToolResponse{
		InstalledPath: target,
	}, nil
}

// NewTestToolRegistry returns a new instance of ToolRegistry for testing purpose.
func NewTestToolRegistry(t *testing.T) *toolregistry.ToolRegistry {
	return toolregistry.NewToolRegistry(&fakeClient{
		testingT: t,
	})
}
