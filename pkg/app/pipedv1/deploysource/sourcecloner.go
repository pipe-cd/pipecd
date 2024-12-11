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

package deploysource

import (
	"context"

	config "github.com/pipe-cd/pipecd/pkg/configv1"
	"github.com/pipe-cd/pipecd/pkg/git"
)

type SourceCloner interface {
	Clone(ctx context.Context, dest string) error
	Revision() string
	RevisionName() string
}

type gitClient interface {
	Clone(ctx context.Context, repoID, remote, branch, destination string) (git.Repo, error)
}

func NewGitSourceCloner(gc gitClient, cfg config.PipedRepository, revisionName, revision string) SourceCloner {
	return &gitSourceCloner{
		revision:     revision,
		revisionName: revisionName,
		gc:           gc,
		cfg:          cfg,
	}
}

type gitSourceCloner struct {
	revision     string
	revisionName string
	gc           gitClient
	cfg          config.PipedRepository
}

func (d *gitSourceCloner) Revision() string {
	return d.revision
}

func (d *gitSourceCloner) RevisionName() string {
	return d.revisionName
}

func (d *gitSourceCloner) Clone(ctx context.Context, dest string) error {
	repo, err := d.gc.Clone(ctx, d.cfg.RepoID, d.cfg.Remote, d.cfg.Branch, dest)
	if err != nil {
		return err
	}
	if err := repo.Checkout(ctx, d.revision); err != nil {
		return err
	}
	return nil
}

type localSourceCloner struct {
	revision     string
	revisionName string
	repo         git.Repo
}

func NewLocalSourceCloner(repo git.Repo, revisionName, revision string) SourceCloner {
	return &localSourceCloner{
		revision:     revision,
		revisionName: revisionName,
		repo:         repo,
	}
}

func (d *localSourceCloner) Revision() string {
	return d.revision
}

func (d *localSourceCloner) RevisionName() string {
	return d.revisionName
}

func (d *localSourceCloner) Clone(ctx context.Context, dest string) error {
	repo, err := d.repo.Copy(dest)
	if err != nil {
		return err
	}
	if err := repo.Checkout(ctx, d.revision); err != nil {
		return err
	}
	return nil
}
