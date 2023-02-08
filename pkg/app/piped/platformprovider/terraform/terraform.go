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

	"github.com/golang-collections/collections/stack"
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

	PlanOutput string
}

func (r PlanResult) NoChanges() bool {
	return r.Adds == 0 && r.Changes == 0 && r.Destroys == 0
}

func (r PlanResult) Render() string {
	TERRAFORM_DIFF_START := "Terraform will perform the following actions:"
	TERRAFORM_DIFF_END := "─"

	startIndex := strings.Index(r.PlanOutput, TERRAFORM_DIFF_START) + len(TERRAFORM_DIFF_START)
	endIndex := strings.Index(r.PlanOutput, TERRAFORM_DIFF_END)
	out := r.PlanOutput[startIndex:endIndex]

	rendered := ""
	curlyBracketStack := stack.New()
	squareBracketStack := stack.New()

	scanner := bufio.NewScanner(strings.NewReader(out))
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}

		r := []rune(line)
		tail := r[len(r)-1]

		// The outermost nest does not have a sign.
		if tail == '{' && curlyBracketStack.Len() == 0 {
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
			curlyBracketStack.Push(r[0])
		}
		if head == '}' {
			r[0] = signMatchPair(*curlyBracketStack, r[0])
		}
		if tail == '[' {
			squareBracketStack.Push(r[0])
		}
		if head == ']' {
			r[0] = signMatchPair(*squareBracketStack, r[0])
		}

		rendered += string(r)
		rendered += "\n"
	}

	return rendered
}

// Return rune at the top of the stack, or r in case of error.
func signMatchPair(st stack.Stack, r rune) rune {
	pop := st.Pop()
	if pop == nil {
		return r
	}

	c, ok := pop.(rune)
	if ok {
		return c
	}
	return r
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
	planHasChangeRegex = regexp.MustCompile(`(?m)^Plan: (\d+) to add, (\d+) to change, (\d+) to destroy.$`)
	planNoChangesRegex = regexp.MustCompile(`(?m)^No changes. Infrastructure is up-to-date.$`)
)

// Borrowed from https://github.com/acarl005/stripansi
const ansi = "[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))"

var ansiRegex = regexp.MustCompile(ansi)

func stripAnsiCodes(str string) string {
	return ansiRegex.ReplaceAllString(str, "")
}

func parsePlanResult(out string, ansiIncluded bool) (PlanResult, error) {
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

	if ansiIncluded {
		out = stripAnsiCodes(out)
	}

	if s := planHasChangeRegex.FindStringSubmatch(out); len(s) == 4 {
		adds, changes, destroys, err := parseNums(s[1], s[2], s[3])
		if err == nil {
			return PlanResult{
				Adds:       adds,
				Changes:    changes,
				Destroys:   destroys,
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
