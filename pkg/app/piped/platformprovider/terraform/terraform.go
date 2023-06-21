// Copyright 2023 The PipeCD Authors.
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
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type options struct {
	noColor  bool
	vars     []string
	varFiles []string

	sharedFlags []string
	initFlags   []string
	planFlags   []string
	applyFlags  []string

	sharedEnvs []string
	initEnvs   []string
	planEnvs   []string
	applyEnvs  []string
}

type Option func(*options)

func WithoutColor() Option {
	return func(opts *options) {
		opts.noColor = true
	}
}

func WithVars(vars []string) Option {
	return func(opts *options) {
		opts.vars = vars
	}
}

func WithVarFiles(files []string) Option {
	return func(opts *options) {
		opts.varFiles = files
	}
}

func WithAdditionalFlags(shared, init, plan, apply []string) Option {
	return func(opts *options) {
		opts.sharedFlags = append(opts.sharedFlags, shared...)
		opts.initFlags = append(opts.initFlags, init...)
		opts.planFlags = append(opts.planFlags, plan...)
		opts.applyFlags = append(opts.applyFlags, apply...)
	}
}

func WithAdditionalEnvs(shared, init, plan, apply []string) Option {
	return func(opts *options) {
		opts.sharedEnvs = append(opts.sharedEnvs, shared...)
		opts.initEnvs = append(opts.initEnvs, init...)
		opts.planEnvs = append(opts.planEnvs, plan...)
		opts.applyEnvs = append(opts.applyEnvs, apply...)
	}
}

type Terraform struct {
	execPath string
	dir      string

	options options
}

func NewTerraform(execPath, dir string, opts ...Option) *Terraform {
	opt := options{}
	for _, o := range opts {
		o(&opt)
	}

	return &Terraform{
		execPath: execPath,
		dir:      dir,
		options:  opt,
	}
}

func (t *Terraform) Version(ctx context.Context) (string, error) {
	args := []string{"version"}
	cmd := exec.CommandContext(ctx, t.execPath, args...)
	cmd.Dir = t.dir
	cmd.Env = append(os.Environ(), t.options.sharedEnvs...)

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
	args = append(args, t.makeCommonCommandArgs()...)
	for _, f := range t.options.initFlags {
		args = append(args, f)
	}

	cmd := exec.CommandContext(ctx, t.execPath, args...)
	cmd.Dir = t.dir
	cmd.Stdout = w
	cmd.Stderr = w

	env := append(os.Environ(), t.options.sharedEnvs...)
	env = append(env, t.options.initEnvs...)
	cmd.Env = env

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
	cmd.Env = append(os.Environ(), t.options.sharedEnvs...)

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
	Imports  int

	PlanOutput string
}

func (r PlanResult) NoChanges() bool {
	return r.Adds == 0 && r.Changes == 0 && r.Destroys == 0 && r.Imports == 0
}

func (r PlanResult) Render() string {
	terraformDiffStart := "Terraform will perform the following actions:"
	terraformDiffEnd := fmt.Sprintf("Plan: %d to import, %d to add, %d to change, %d to destroy.", r.Imports, r.Adds, r.Changes, r.Destroys)

	startIndex := strings.Index(r.PlanOutput, terraformDiffStart) + len(terraformDiffStart)
	endIndex := strings.Index(r.PlanOutput, terraformDiffEnd) + len(terraformDiffEnd)
	out := r.PlanOutput[startIndex:endIndex]

	rendered := ""
	var curlyBracketStack []rune
	var squareBracketStack []rune

	scanner := bufio.NewScanner(strings.NewReader(out))
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}

		r := []rune(line)
		tail := r[len(r)-1]

		// The outermost nest does not have a sign.
		if tail == '{' && len(curlyBracketStack) == 0 {
			// Terraform's outermost block would be resource block.
			deadline := strings.Index(string(r), "resource")
			for i := 0; i < deadline; i++ {
				r[i] = ' '
			}
		}

		// Get head rune without tab and space.
		head, pos := headRuneWithoutWhiteSpace(r)
		if pos < 0 {
			continue
		}

		// Move sign to the beginning.
		if head == '+' || head == '-' || head == '~' {
			r[0], r[pos] = r[pos], r[0]
		}

		// Corresponding pairs with corresponding sign.
		if tail == '{' {
			curlyBracketStack = append(curlyBracketStack, r[0])
		}
		if head == '}' {
			r[0] = signMatchBracket(&curlyBracketStack, r[0])
		}
		if tail == '[' {
			squareBracketStack = append(squareBracketStack, r[0])
		}
		if head == ']' {
			r[0] = signMatchBracket(&squareBracketStack, r[0])
		}

		rendered += string(r)
		rendered += "\n"
	}

	return rendered
}

// Return rune at the top of the stack, or r in case of error.
func signMatchBracket(l *[]rune, r rune) rune {
	list := *l
	if len(list) == 0 {
		return r
	}
	n := len(list) - 1
	v := list[n]
	*l = list[:n]
	return v
}

func headRuneWithoutWhiteSpace(r []rune) (rune, int) {
	for i, ri := range r {
		if !(ri == '\t' || ri == ' ') {
			return ri, i
		}
	}
	return ' ', -1
}

func GetExitCode(err error) int {
	if err == nil {
		return 0
	}
	if exitErr, ok := err.(*exec.ExitError); ok {
		return exitErr.ExitCode()
	}
	return 1
}

func (t *Terraform) Plan(ctx context.Context, w io.Writer) (PlanResult, error) {
	args := []string{
		"plan",
		"-lock=false",
		"-detailed-exitcode",
	}
	args = append(args, t.makeCommonCommandArgs()...)
	for _, f := range t.options.planFlags {
		args = append(args, f)
	}

	var buf bytes.Buffer
	stdout := io.MultiWriter(w, &buf)

	cmd := exec.CommandContext(ctx, t.execPath, args...)
	cmd.Dir = t.dir
	cmd.Stdout = stdout
	cmd.Stderr = stdout

	env := append(os.Environ(), t.options.sharedEnvs...)
	env = append(env, t.options.planEnvs...)
	cmd.Env = env

	io.WriteString(w, fmt.Sprintf("terraform %s", strings.Join(args, " ")))
	err := cmd.Run()
	switch GetExitCode(err) {
	case 0:
		return PlanResult{}, nil
	case 2:
		return parsePlanResult(buf.String(), !t.options.noColor)
	default:
		return PlanResult{}, err
	}
}

func (t *Terraform) makeCommonCommandArgs() (args []string) {
	if t.options.noColor {
		args = append(args, "-no-color")
	}
	for _, v := range t.options.vars {
		args = append(args, fmt.Sprintf("-var=%s", v))
	}
	for _, f := range t.options.varFiles {
		args = append(args, fmt.Sprintf("-var-file=%s", f))
	}
	for _, f := range t.options.sharedFlags {
		args = append(args, f)
	}
	return
}

var (
	// Import block was introduced from Terraform v1.5.0.
	// Keep this regex for backward compatibility.
	planHasChangeRegex = regexp.MustCompile(`(?m)^Plan:(?: \d+ to import,)?? (\d+) to add, (\d+) to change, (\d+) to destroy.$`)
	planNoChangesRegex = regexp.MustCompile(`(?m)^No changes. Infrastructure is up-to-date.$`)
)

// Borrowed from https://github.com/acarl005/stripansi
const ansi = "[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))"

var ansiRegex = regexp.MustCompile(ansi)

func stripAnsiCodes(str string) string {
	return ansiRegex.ReplaceAllString(str, "")
}

func parsePlanResult(out string, ansiIncluded bool) (PlanResult, error) {
	parseNums := func(vals ...string) (imports, adds, changes, destroys int, err error) {
		if len(vals) < 3 || len(vals) > 4 {
			err = fmt.Errorf("invalid plan result: %v", vals)
			return
		}

		var impt, add, change, destroy string
		if len(vals) == 3 {
			add = vals[0]
			change = vals[1]
			destroy = vals[2]
		}

		if len(vals) == 4 {
			impt = vals[0]
			add = vals[1]
			change = vals[2]
			destroy = vals[3]
		}

		if impt != "" {
			imports, err = strconv.Atoi(impt)
			if err != nil {
				return
			}
		}

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

	if ansiIncluded {
		out = stripAnsiCodes(out)
	}

	if s := planHasChangeRegex.FindStringSubmatch(out); len(s) > 0 {
		imports, adds, changes, destroys, err := parseNums(s[1:]...)
		if err == nil {
			return PlanResult{
				Adds:       adds,
				Changes:    changes,
				Destroys:   destroys,
				Imports:    imports,
				PlanOutput: out,
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
	args = append(args, t.makeCommonCommandArgs()...)
	for _, f := range t.options.applyFlags {
		args = append(args, f)
	}

	cmd := exec.CommandContext(ctx, t.execPath, args...)
	cmd.Dir = t.dir
	cmd.Stdout = w
	cmd.Stderr = w

	env := append(os.Environ(), t.options.sharedEnvs...)
	env = append(env, t.options.applyEnvs...)
	cmd.Env = env

	io.WriteString(w, fmt.Sprintf("terraform %s", strings.Join(args, " ")))
	return cmd.Run()
}
