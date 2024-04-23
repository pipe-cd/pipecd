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
	"strings"

	"github.com/pipe-cd/pipecd/pkg/diff"
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
	d, err := diff.DiffStructs(old, new, opts...)
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
	taskDiff, err := diff.DiffByCommand(command, old.TaskDefinition, new.TaskDefinition)
	if err != nil {
		return nil, err
	}

	serviceDiff, err := diff.DiffByCommand(command, old.ServiceDefinition, new.ServiceDefinition)
	if err != nil {
		return nil, err
	}

	return bytes.Join([][]byte{
		[]byte("# 1. ServiceDefinition"),
		serviceDiff,
		[]byte("\n# 2. TaskDefinition"),
		taskDiff,
	}, []byte("\n")), nil
}
