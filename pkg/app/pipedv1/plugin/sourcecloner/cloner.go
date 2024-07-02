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

package sourcecloner

import (
	"context"
	"os"
	"os/exec"
	"sync"

	"github.com/pipe-cd/pipecd/pkg/app/piped/deploysource"
	"github.com/pipe-cd/pipecd/pkg/git"
)

var _ deploysource.SourceCloner = (*Cloner)(nil)

var gitPath = sync.OnceValues(func() (string, error) {
	path, err := exec.LookPath("git")
	if err != nil {
		return "", err
	}
	return path, nil
})

// Cloner is a source cloner for git repositories.
type Cloner struct {
	gitPath      string
	remotePath   string
	revision     string
	revisionName string
}

// NewCloner creates a new Cloner instance.
// if revision is empty, it will checkout defualt branch.
func NewCloner(remotePath, revision, revisionName string) (*Cloner, error) {
	gitPath, err := gitPath()
	if err != nil {
		return nil, err
	}

	return &Cloner{
		gitPath:      gitPath,
		remotePath:   remotePath,
		revision:     revision,
		revisionName: revisionName,
	}, nil
}

// Clone implements deploysource.SourceCloner.
func (c *Cloner) Clone(ctx context.Context, dest string) error {
	if _, err := git.RunGitCommand(ctx, c.gitPath, "", os.Environ(), "clone", c.remotePath, dest); err != nil {
		return err
	}
	if c.revision == "" {
		if _, err := git.RunGitCommand(ctx, c.gitPath, "", os.Environ(), "-C", dest, "checkout", c.revision); err != nil {
			return err
		}
	}
	return nil
}

// Revision implements deploysource.SourceCloner.
func (c *Cloner) Revision() string {
	return c.revision
}

// RevisionName implements deploysource.SourceCloner.
func (c *Cloner) RevisionName() string {
	return c.revisionName
}
