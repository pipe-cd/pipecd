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

package ecs

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/pipe-cd/pipecd/pkg/diff"
	"sigs.k8s.io/yaml"
)

const (
	diffCommand = "diff"
)

type DiffResult struct {
	Diff *diff.Result
	Old  ECSManifest
	New  ECSManifest
}

func (d *DiffResult) NoChange() bool {
	return len(d.Diff.Nodes()) == 0
}

func Diff(old, new ECSManifest, opts ...diff.Option) (*DiffResult, error) {
	key := old.ServiceDefinition.ServiceName
	if key == nil {
		return nil, fmt.Errorf("service name is required in the old service definition")
	}
	old_u, err := old.Unstructured()
	if err != nil {
		return nil, fmt.Errorf("failed to convert old service definition to unstructured: %w", err)
	}
	new_u, err := new.Unstructured()
	if err != nil {
		return nil, fmt.Errorf("failed to convert new service definition to unstructured: %w", err)
	}

	d, err := diff.DiffUnstructureds(old_u, new_u, *key, opts...)
	if err != nil {
		return nil, err
	}

	if !d.HasDiff() {
		return &DiffResult{Diff: d}, nil
	}

	ret := &DiffResult{
		Old:  old,
		New:  new,
		Diff: d,
	}
	return ret, nil
}

type DiffRenderOptions struct {
	// If true, use "diff" command to render.
	UseDiffCommand bool
}

func (d *DiffResult) Render(opt DiffRenderOptions) string {
	var b strings.Builder
	opts := []diff.RenderOption{
		diff.WithLeftPadding(1),
	}
	renderer := diff.NewRenderer(opts...)
	if !opt.UseDiffCommand {
		b.WriteString(renderer.Render(d.Diff.Nodes()))
	} else {
		d, err := diffByCommand(diffCommand, d.Old, d.New)
		if err != nil {
			b.WriteString(fmt.Sprintf("An error occurred while rendering diff (%v)", err))
		} else {
			b.Write(d)
		}
	}
	b.WriteString("\n")

	return b.String()
}

func diffByCommand(command string, old, new ECSManifest) ([]byte, error) {
	taskDiff, err := diffYamlByCommand(command, old.TaskDefinition, new.TaskDefinition)
	if err != nil {
		return nil, err
	}
	serviceDiff, err := diffYamlByCommand(command, old.ServiceDefinition, new.ServiceDefinition)
	if err != nil {
		return nil, err
	}

	// TODO merge? or just return both?
	return bytes.Join([][]byte{taskDiff, serviceDiff}, []byte("\n")), nil
}

func diffYamlByCommand(command string, old, new interface{}) ([]byte, error) {
	oldBytes, err := yaml.Marshal(old)
	if err != nil {
		return nil, err
	}
	newBytes, err := yaml.Marshal(new)
	if err != nil {
		return nil, err
	}

	oldFile, err := os.CreateTemp("", "old")
	if err != nil {
		return nil, err
	}
	defer os.Remove(oldFile.Name())
	if _, err := oldFile.Write(oldBytes); err != nil {
		return nil, err
	}

	newFile, err := os.CreateTemp("", "new")
	if err != nil {
		return nil, err
	}
	defer os.Remove(newFile.Name())
	if _, err := newFile.Write(newBytes); err != nil {
		return nil, err
	}

	var stdout, stderr bytes.Buffer
	cmd := exec.Command(command, "-u", "-N", oldFile.Name(), newFile.Name())
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	if stdout.Len() > 0 {
		// diff exits with a non-zero status when the files don't match.
		// Ignore that failure as long as we get output.
		err = nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to run diff, err = %w, %s", err, stderr.String())
	}

	// Remove two-line header from output.
	data := bytes.TrimSpace(stdout.Bytes())
	rows := bytes.SplitN(data, []byte("\n"), 3)
	if len(rows) == 3 {
		return rows[2], nil
	}
	return data, nil
}
