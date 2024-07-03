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

package lambda

import (
	"fmt"
	"strings"

	"github.com/pipe-cd/pipecd/pkg/diff"
)

const (
	diffCommand = "diff"
)

type DiffResult struct {
	Diff *diff.Result
	Old  FunctionManifest
	New  FunctionManifest
}

func (d *DiffResult) NoChange() bool {
	return len(d.Diff.Nodes()) == 0
}

func Diff(old, new FunctionManifest, opts ...diff.Option) (*DiffResult, error) {
	d, err := diff.DiffStructureds(old, new, opts...)
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
		d, err := renderByCommand(diffCommand, d.Old, d.New)
		if err != nil {
			b.WriteString(fmt.Sprintf("An error occurred while rendering diff (%v)", err))
		} else {
			b.Write(d)
		}
	}
	b.WriteString("\n")

	return b.String()
}

func renderByCommand(command string, old, new FunctionManifest) ([]byte, error) {
	diff, err := diff.RenderByCommand(command, old, new)
	if err != nil {
		return nil, err
	}

	return diff, nil
}
