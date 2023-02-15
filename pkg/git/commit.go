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

package git

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	separator       = "__GIT_LOG_SEPARATOR__"
	delimiter       = "__GIT_LOG_DELIMITER__"
	fieldNum        = 7
	commitLogFormat = separator +
		"%an" + delimiter +
		"%cn" + delimiter +
		"%at" + delimiter +
		"%H" + delimiter +
		"%h" + delimiter +
		"%s" + delimiter +
		"%b"
)

type Commit struct {
	Author          string
	Committer       string
	CreatedAt       int
	Hash            string
	AbbreviatedHash string
	Message         string
	Body            string
}

// We was using json encoding to parse commit log,
// but the commit message may contain various escape chars,
// so I think reading each log line and map to Commit field is a good way.
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
		Author:          fields[0],
		Committer:       fields[1],
		CreatedAt:       createdAt,
		Hash:            fields[3],
		AbbreviatedHash: fields[4],
		Message:         fields[5],
		Body:            strings.TrimSpace(fields[6]),
	}, nil
}
