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

package main

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

const (
	separator       = "__GIT_LOG_SEPARATOR__"
	delimiter       = "__GIT_LOG_DELIMITER__"
	fieldNum        = 9
	commitLogFormat = separator +
		"%an" + delimiter + // author name
		"%cn" + delimiter + // committer name
		"%at" + delimiter + // author date (UNIX timestamp)
		"%H" + delimiter + // commit hash
		"%h" + delimiter + // abbreviated commit hash
		"%P" + delimiter + // parent hashes
		"%p" + delimiter + // abbreviated parent hashes
		"%s" + delimiter + // subject
		"%b" // body
)

type Commit struct {
	Author                  string   `json:"author,omitempty"`
	Committer               string   `json:"committer,omitempty"`
	CreatedAt               int      `json:"createdAt,omitempty"`
	Hash                    string   `json:"hash,omitempty"`
	AbbreviatedHash         string   `json:"abbreviatedHash,omitempty"`
	ParentHashes            []string `json:"parentHashes,omitempty"`
	AbbreviatedParentHashes []string `json:"abbreviatedParentHashes,omitempty"`
	Subject                 string   `json:"subject,omitempty"`
	Body                    string   `json:"body,omitempty"`
}

func (c Commit) IsMerge() bool {
	return len(c.ParentHashes) == 2
}

func (c Commit) PullRequestNumber() (int, bool) {
	if !c.IsMerge() {
		return 0, false
	}

	subs := defaultMergeCommitRegex.FindStringSubmatch(c.Subject)
	if len(subs) != 2 {
		return 0, false
	}

	prNumber, err := strconv.Atoi(subs[1])
	if err != nil {
		return 0, false
	}

	return prNumber, true
}

func parseCommits(log string) ([]Commit, error) {
	lines := strings.Split(log, separator)
	if len(lines) < 1 {
		return nil, fmt.Errorf("invalid log")
	}
	commits := make([]Commit, 0, len(lines))
	for _, line := range lines[1:] {
		commit, err := parseCommit(line)
		if err != nil {
			return nil, err
		}
		commits = append(commits, commit)
	}
	return commits, nil
}

func parseCommit(log string) (Commit, error) {
	fields := strings.Split(log, delimiter)
	if len(fields) != fieldNum {
		return Commit{}, fmt.Errorf("invalid log: log line should contain %d fields but got %d", fieldNum, len(fields))
	}
	createdAt, err := strconv.Atoi(fields[2])
	if err != nil {
		return Commit{}, err
	}
	return Commit{
		Author:                  fields[0],
		Committer:               fields[1],
		CreatedAt:               createdAt,
		Hash:                    fields[3],
		AbbreviatedHash:         fields[4],
		ParentHashes:            strings.Split(fields[5], " "),
		AbbreviatedParentHashes: strings.Split(fields[6], " "),
		Subject:                 fields[7],
		Body:                    strings.TrimSpace(fields[8]),
	}, nil
}

// listCommits returns a list of commits between the given revision range.
func listCommits(ctx context.Context, gitExecPath, repoDir string, revisionRange string) ([]Commit, error) {
	args := []string{
		"log",
		"--no-decorate",
		fmt.Sprintf("--pretty=format:%s", commitLogFormat),
	}
	if revisionRange != "" {
		args = append(args, revisionRange)
	}

	cmd := exec.CommandContext(ctx, gitExecPath, args...)
	cmd.Dir = repoDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("err: %w, out: %s", err, string(out))
	}

	return parseCommits(string(out))
}

// changedFiles returns a list of files those were touched between two commits.
func changedFiles(ctx context.Context, gitExecPath, repoDir, from, to string) ([]string, error) {
	cmd := exec.CommandContext(ctx, gitExecPath, "diff", "--name-only", from, to)
	cmd.Dir = repoDir
	out, err := cmd.CombinedOutput()

	if err != nil {
		return nil, fmt.Errorf("err: %w, out: %s", err, string(out))
	}

	var (
		lines = strings.Split(string(out), "\n")
		files = make([]string, 0, len(lines))
	)
	// We need to remove all empty lines since the result may include them.
	for _, f := range lines {
		if f != "" {
			files = append(files, f)
		}
	}
	return files, nil
}

// readFileAtCommit reads the content of a specific file at the given commit.
func readFileAtCommit(ctx context.Context, gitExecPath, repoDir, filePath, commit string) ([]byte, error) {
	args := []string{
		"show",
		fmt.Sprintf("%s:%s", commit, filePath),
	}

	cmd := exec.CommandContext(ctx, gitExecPath, args...)
	cmd.Dir = repoDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("err: %w, out: %s", err, string(out))
	}

	return bytes.TrimSpace(out), nil
}

func addSafeDirectory(ctx context.Context, gitExecPath, repoDir string) error {
	args := []string{
		"config",
		"--global",
		"--add",
		"safe.directory",
		repoDir,
	}

	cmd := exec.CommandContext(ctx, gitExecPath, args...)
	cmd.Dir = repoDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("err: %w, out: %s", err, string(out))
	}
	return nil
}
