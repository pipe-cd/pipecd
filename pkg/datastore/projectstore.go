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

package datastore

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pipe-cd/pipecd/pkg/model"
)

type projectCollection struct {
	requestedBy Commander
}

func (p *projectCollection) Kind() string {
	return "Project"
}

func (p *projectCollection) Factory() Factory {
	return func() interface{} {
		return &model.Project{}
	}
}

func (p *projectCollection) ListInUsedShards() []Shard {
	return []Shard{
		ClientShard,
	}
}

func (p *projectCollection) GetUpdatableShard() (Shard, error) {
	switch p.requestedBy {
	case WebCommander:
		return ClientShard, nil
	default:
		return "", ErrUnsupported
	}
}

func (p *projectCollection) Encode(e interface{}) (map[Shard][]byte, error) {
	const errFmt = "failed while encode Project object: %s"

	me, ok := e.(*model.Project)
	if !ok {
		return nil, fmt.Errorf(errFmt, "type not matched")
	}

	data, err := json.Marshal(me)
	if err != nil {
		return nil, fmt.Errorf(errFmt, "unable to marshal entity data")
	}
	return map[Shard][]byte{
		ClientShard: data,
	}, nil
}

type ProjectStore interface {
	Add(ctx context.Context, proj *model.Project) error
	Get(ctx context.Context, id string) (*model.Project, error)
	List(ctx context.Context, opts ListOptions) ([]model.Project, error)
	UpdateProjectStaticAdmin(ctx context.Context, id, username, password string) error
	EnableStaticAdmin(ctx context.Context, id string) error
	DisableStaticAdmin(ctx context.Context, id string) error
	UpdateProjectSSOConfig(ctx context.Context, id string, sso *model.ProjectSSOConfig) error
	UpdateProjectRBACConfig(ctx context.Context, id string, rbac *model.ProjectRBACConfig) error
	AddProjectRBACRole(ctx context.Context, id, name string, policies []*model.ProjectRBACPolicy) error
	UpdateProjectRBACRole(ctx context.Context, id, name string, policies []*model.ProjectRBACPolicy) error
	DeleteProjectRBACRole(ctx context.Context, id, name string) error
	AddProjectUserGroup(ctx context.Context, id, sso, role string) error
	DeleteProjectUserGroup(ctx context.Context, id, sso string) error
}

type projectStore struct {
	backend
	commander Commander
	nowFunc   func() time.Time
}

func NewProjectStore(ds DataStore, c Commander) ProjectStore {
	return &projectStore{
		backend: backend{
			ds:  ds,
			col: &projectCollection{requestedBy: c},
		},
		commander: c,
		nowFunc:   time.Now,
	}
}

func (s *projectStore) Add(ctx context.Context, proj *model.Project) error {
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
	return s.ds.Create(ctx, s.col, proj.Id, proj)
}

func (s *projectStore) Get(ctx context.Context, id string) (*model.Project, error) {
	var entity model.Project
	if err := s.ds.Get(ctx, s.col, id, &entity); err != nil {
		return nil, err
	}
	entity.SetLegacyUserGroups()
	entity.SetBuiltinRBACRoles()
	return &entity, nil
}

func (s *projectStore) List(ctx context.Context, opts ListOptions) ([]model.Project, error) {
	it, err := s.ds.Find(ctx, s.col, opts)
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

func (s *projectStore) update(ctx context.Context, id string, updater func(project *model.Project) error) error {
	now := s.nowFunc().Unix()
	return s.ds.Update(ctx, s.col, id, func(e interface{}) error {
		p := e.(*model.Project)
		if err := updater(p); err != nil {
			return err
		}
		p.UpdatedAt = now
		return p.Validate()
	})
}

// UpdateProjectStaticAdmin updates the static admin user settings.
func (s *projectStore) UpdateProjectStaticAdmin(ctx context.Context, id, username, password string) error {
	return s.update(ctx, id, func(p *model.Project) error {
		if p.StaticAdmin == nil {
			p.StaticAdmin = &model.ProjectStaticUser{}
		}
		return p.StaticAdmin.Update(username, password)
	})
}

// EnableStaticAdmin enables static admin login.
func (s *projectStore) EnableStaticAdmin(ctx context.Context, id string) error {
	return s.update(ctx, id, func(p *model.Project) error {
		p.StaticAdminDisabled = false
		return nil
	})
}

// DisableStaticAdmin disables static admin login.
func (s *projectStore) DisableStaticAdmin(ctx context.Context, id string) error {
	return s.update(ctx, id, func(p *model.Project) error {
		p.StaticAdminDisabled = true
		return nil
	})
}

// UpdateProjectSSOConfig updates project single sign on settings.
func (s *projectStore) UpdateProjectSSOConfig(ctx context.Context, id string, sso *model.ProjectSSOConfig) error {
	return s.update(ctx, id, func(p *model.Project) error {
		if p.Sso == nil {
			p.Sso = &model.ProjectSSOConfig{}
		}
		p.Sso.Update(sso)
		return nil
	})
}

// UpdateProjectRBACConfig updates project single sign on settings.
func (s *projectStore) UpdateProjectRBACConfig(ctx context.Context, id string, rbac *model.ProjectRBACConfig) error {
	return s.update(ctx, id, func(p *model.Project) error {
		p.Rbac = rbac
		return nil
	})
}

// AddProjectRBACRole adds the custom rbac role.
func (s *projectStore) AddProjectRBACRole(ctx context.Context, id, name string, policies []*model.ProjectRBACPolicy) error {
	return s.update(ctx, id, func(p *model.Project) error {
		return p.AddRBACRole(name, policies)
	})
}

// UpdateProjectRBACRole updates the custom rbac role.
func (s *projectStore) UpdateProjectRBACRole(ctx context.Context, id, name string, policies []*model.ProjectRBACPolicy) error {
	return s.update(ctx, id, func(p *model.Project) error {
		return p.UpdateRBACRole(name, policies)
	})
}

// DeleteProjectRBACRole deletes the custom rbac role.
func (s *projectStore) DeleteProjectRBACRole(ctx context.Context, id, name string) error {
	return s.update(ctx, id, func(p *model.Project) error {
		return p.DeleteRBACRole(name)
	})
}

// AddProjectUserGroup adds the user group.
func (s *projectStore) AddProjectUserGroup(ctx context.Context, id, sso, role string) error {
	return s.update(ctx, id, func(p *model.Project) error {
		return p.AddUserGroup(sso, role)
	})
}

// DeleteProjectUserGroup deletes the user group.
func (s *projectStore) DeleteProjectUserGroup(ctx context.Context, id, sso string) error {
	return s.update(ctx, id, func(p *model.Project) error {
		return p.DeleteUserGroup(sso)
	})
}
