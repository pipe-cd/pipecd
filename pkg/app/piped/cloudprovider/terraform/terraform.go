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

package terraform

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type Terraform struct {
	execPath string
	dir      string
	vars     []string
	varFiles []string
}

func NewTerraform(execPath, dir string, vars, varFiles []string) *Terraform {
	return &Terraform{
		execPath: execPath,
		dir:      dir,
		vars:     vars,
		varFiles: varFiles,
	}
}

func (t *Terraform) Version(ctx context.Context) (string, error) {
	args := []string{"version"}
	cmd := exec.CommandContext(ctx, t.execPath, args...)
	cmd.Dir = t.dir

	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), err
	}

	return strings.TrimSpace(string(out)), nil
}

func (t *Terraform) Init(ctx context.Context, w io.Writer) error {
	args := []string{
		"init",
	}
	for _, v := range t.vars {
		args = append(args, fmt.Sprintf("-var=%s", v))
	}
	for _, f := range t.varFiles {
		args = append(args, fmt.Sprintf("-var-file=%s", f))
	}

	cmd := exec.CommandContext(ctx, t.execPath, args...)
	cmd.Dir = t.dir
	cmd.Stdout = w
	cmd.Stderr = w

	io.WriteString(w, fmt.Sprintf("terraform %s", strings.Join(args, " ")))
	return cmd.Run()
}

func (t *Terraform) SelectWorkspace(ctx context.Context, workspace string) error {
	args := []string{
		"workspace",
		"select",
		workspace,
	}
	cmd := exec.CommandContext(ctx, t.execPath, args...)
	cmd.Dir = t.dir

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to select workspace: %s (%w)", string(out), err)
	}

	return nil
}

type PlanResult struct {
	Adds     int
	Changes  int
	Destroys int
}

func (r PlanResult) NoChanges() bool {
	return r.Adds == 0 && r.Changes == 0 && r.Destroys == 0
}

func (t *Terraform) Plan(ctx context.Context, w io.Writer) (PlanResult, error) {
	args := []string{
		"plan",
		// TODO: Remove this -no-color flag after parsePlanResult supports parsing the message containing color codes.
		"-no-color",
	}
	for _, v := range t.vars {
		args = append(args, fmt.Sprintf("-var=%s", v))
	}
	for _, f := range t.varFiles {
		args = append(args, fmt.Sprintf("-var-file=%s", f))
	}

	var buf bytes.Buffer
	stdout := io.MultiWriter(w, &buf)

	cmd := exec.CommandContext(ctx, t.execPath, args...)
	cmd.Dir = t.dir
	cmd.Stdout = stdout
	cmd.Stderr = stdout

	io.WriteString(w, fmt.Sprintf("terraform %s", strings.Join(args, " ")))
	if err := cmd.Run(); err != nil {
		return PlanResult{}, err
	}

	return parsePlanResult(buf.String())
}

var (
	planHasChangeRegex = regexp.MustCompile(`(?m)^Plan: (\d+) to add, (\d+) to change, (\d+) to destroy.$`)
	planNoChangesRegex = regexp.MustCompile(`(?m)^No changes. Infrastructure is up-to-date.$`)
)

func parsePlanResult(out string) (PlanResult, error) {
	parseNums := func(add, change, destroy string) (adds int, changes int, destroys int, err error) {
		adds, err = strconv.Atoi(add)
		if err != nil {
			return
		}
		changes, err = strconv.Atoi(change)
		if err != nil {
			return
		}
		destroys, err = strconv.Atoi(destroy)
		if err != nil {
			return
		}
		return
	}

	if s := planHasChangeRegex.FindStringSubmatch(out); len(s) == 4 {
		adds, changes, destroys, err := parseNums(s[1], s[2], s[3])
		if err == nil {
			return PlanResult{
				Adds:     adds,
				Changes:  changes,
				Destroys: destroys,
			}, nil
		}
	}

	if s := planNoChangesRegex.FindStringSubmatch(out); len(s) > 0 {
		return PlanResult{}, nil
	}

	return PlanResult{}, fmt.Errorf("unable to parse plan output")
}

func (t *Terraform) Apply(ctx context.Context, w io.Writer) error {
	args := []string{
		"apply",
		"-auto-approve",
		"-input=false",
	}
	for _, v := range t.vars {
		args = append(args, fmt.Sprintf("-var=%s", v))
	}
	for _, f := range t.varFiles {
		args = append(args, fmt.Sprintf("-var-file=%s", f))
	}

	cmd := exec.CommandContext(ctx, t.execPath, args...)
	cmd.Dir = t.dir
	cmd.Stdout = w
	cmd.Stderr = w

	io.WriteString(w, fmt.Sprintf("terraform %s", strings.Join(args, " ")))
	return cmd.Run()
}
