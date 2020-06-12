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

package datastore

import (
	"context"
	"time"

	"github.com/pipe-cd/pipe/pkg/model"
)

const projectModelKind = "Project"

var projectFactory = func() interface{} {
	return &model.Project{}
}

type ProjectStore interface {
	AddProject(ctx context.Context, proj *model.Project) error
	ListProjects(ctx context.Context, opts ListOptions) ([]model.Project, error)
}

type projectStore struct {
	backend
	nowFunc func() time.Time
}

func NewProjectStore(ds DataStore) ProjectStore {
	return &projectStore{
		backend: backend{
			ds: ds,
		},
		nowFunc: time.Now,
	}
}

func (s *projectStore) AddProject(ctx context.Context, proj *model.Project) error {
	now := s.nowFunc().Unix()
	if proj.CreatedAt == 0 {
		proj.CreatedAt = now
	}
	if proj.UpdatedAt == 0 {
		proj.UpdatedAt = now
	}
	if err := proj.Validate(); err != nil {
		return err
	}
	return s.ds.Create(ctx, projectModelKind, proj.Id, proj)
}

func (s *projectStore) ListProjects(ctx context.Context, opts ListOptions) ([]model.Project, error) {
	it, err := s.ds.Find(ctx, projectModelKind, opts)
	if err != nil {
		return nil, err
	}
	ps := make([]model.Project, 0)
	for {
		var p model.Project
		err := it.Next(&p)
		if err == ErrIteratorDone {
			break
		}
		if err != nil {
			return nil, err
		}
		ps = append(ps, p)
	}
	return ps, nil
}
