package terraform

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type DiffListResult struct {
	Adds     []string
	Changes  []DiffListChanges
	Destroys []string
}

type DiffListChanges struct {
	Old string
	New string
}

func (t *Terraform) Diff(ctx context.Context, w io.Writer) (DiffListResult, error) {
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
		return DiffListResult{}, nil
	case 2:
		return parseDiffListResult(buf.String(), !t.options.noColor)
	default:
		return DiffListResult{}, err
	}
}

func parseDiffListResult(out string, ansiIncluded bool) (DiffListResult, error) {
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
			return DiffListResult{
				Adds:     adds,
				Changes:  changes,
				Destroys: destroys,
			}, nil
		}
	}

	if s := planNoChangesRegex.FindStringSubmatch(out); len(s) > 0 {
		return DiffListResult{}, nil
	}

	return DiffListResult{}, fmt.Errorf("unable to parse plan output")
}
