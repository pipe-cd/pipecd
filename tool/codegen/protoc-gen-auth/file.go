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

package main

const fileTpl = `// Code generated by protoc-gen-auth. DO NOT EDIT.
// source: {{ .InputPath }}

package webservice

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/cache"
	"github.com/pipe-cd/pipecd/pkg/cache/memorycache"
	"github.com/pipe-cd/pipecd/pkg/config"
	"github.com/pipe-cd/pipecd/pkg/datastore"
	"github.com/pipe-cd/pipecd/pkg/model"
	"github.com/pipe-cd/pipecd/pkg/rpc/rpcauth"
)

type webApiProjectStore interface {
	Get(ctx context.Context, id string) (*model.Project, error)
}

type authorizer struct {
	projectStore webApiProjectStore
	rbacCache    cache.Cache
	// List of debugging/quickstart projects.
	projectsInConfig map[string]config.ControlPlaneProject
	logger           *zap.Logger
}

// NewRBACAuthorizer returns an RBACAuthorizer object for checking requested method based on RBAC.
func NewRBACAuthorizer(
	ctx context.Context,
	ds datastore.DataStore,
	projects map[string]config.ControlPlaneProject,
	logger *zap.Logger,
) rpcauth.RBACAuthorizer {
	return &authorizer{
		projectStore:     datastore.NewProjectStore(ds),
		rbacCache:        memorycache.NewTTLCache(ctx, 10*time.Minute, 5*time.Minute),
		projectsInConfig: projects,
		logger:           logger.Named("authorizer"),
	}
}

func (a *authorizer) getProjectRBACRoles(ctx context.Context, projectID string, roles []string) ([]*model.ProjectRBACRole, error) {
	all, err := a.getAllProjectRBACRoles(ctx, projectID)
	if err != nil {
		return nil, err
	}
	rs := make(map[string]*model.ProjectRBACRole, len(all))
	for _, v := range all {
		rs[v.Name] = v
	}
	ret := make([]*model.ProjectRBACRole, 0, len(roles))
	for _, r := range roles {
		if v, ok := rs[r]; ok {
			ret = append(ret, v)
		}
	}
	return ret, nil
}

func (a *authorizer) getAllProjectRBACRoles(ctx context.Context, projectID string) ([]*model.ProjectRBACRole, error) {
	if _, ok := a.projectsInConfig[projectID]; ok {
		p := &model.Project{Id: projectID}
		p.SetBuiltinRBACRoles()
		return p.RbacRoles, nil
	}

	if v, err := a.rbacCache.Get(projectID); err == nil {
		return v.([]*model.ProjectRBACRole), nil
	}

	p, err := a.projectStore.Get(ctx, projectID)
	if err != nil {
		a.logger.Error("failed to get project",
			zap.String("project", projectID),
			zap.Error(err),
		)
		return nil, err
	}

	if err = a.rbacCache.Put(projectID, p.RbacRoles); err != nil {
		a.logger.Warn("unable to store the rbac in memory cache",
			zap.String("project", projectID),
			zap.Error(err),
		)
	}
	return p.RbacRoles, nil
}

// Authorize checks whether a role is enough for given gRPC method or not.
func (a *authorizer) Authorize(ctx context.Context, method string, r model.Role) bool {
	roles, err := a.getProjectRBACRoles(ctx, r.ProjectId, r.ProjectRbacRoles)
	if err != nil {
		a.logger.Error("failed to get rbac roles",
			zap.String("project", r.ProjectId),
			zap.Error(err),
		)
		return false
	}

	if len(roles) == 0 {
		a.logger.Error("the user does not have existing RBAC roles",
			zap.String("project", r.ProjectId),
			zap.Strings("user-roles", r.ProjectRbacRoles),
		)
		return false
	}

	verify := func(typ model.ProjectRBACResource_ResourceType, action model.ProjectRBACPolicy_Action) bool {
		for _, r := range roles {
			if r.HasPermission(typ, action) {
				return true
			}
		}
		return false
	}

	switch method {
	{{- range .Methods }}
	case "/grpc.service.webservice.WebService/{{ .Name }}":
	        {{- if .Ignored }}
			return true
		{{- else }}
			return verify(model.ProjectRBACResource_{{ .Resource }}, model.ProjectRBACPolicy_{{ .Action }})
		{{- end }}
	{{- end }}
	}
	return false
}
`
