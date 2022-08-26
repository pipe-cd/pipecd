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

// Code generated by protoc-gen-auth. DO NOT EDIT.
// source: pkg/app/server/service/webservice/service.proto

package webservice

import (
	"context"
	"time"

	"github.com/pipe-cd/pipecd/pkg/cache"
	"github.com/pipe-cd/pipecd/pkg/cache/memorycache"
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
}

// NewRBACAuthorizer returns an RBACAuthorizer object for checking requested method based on RBAC.
func NewRBACAuthorizer(ctx context.Context, ds datastore.DataStore) rpcauth.RBACAuthorizer {
	w := datastore.WebCommander
	return &authorizer{
		projectStore: datastore.NewProjectStore(ds, w),
		rbacCache:    memorycache.NewTTLCache(ctx, 10*time.Minute, 5*time.Minute),
	}
}

func (a *authorizer) getRBAC(ctx context.Context, projectID string) (*rbac, error) {
	r, err := a.rbacCache.Get(projectID)
	if err == nil {
		return r.(*rbac), nil
	}
	project, err := a.projectStore.Get(ctx, projectID)
	if err != nil {
		return nil, err
	}
	v := &rbac{project.RbacRoles}
	a.rbacCache.Put(projectID, v)
	return v, nil
}

type rbac struct {
	Roles []*model.ProjectRBACRole
}

func (r *rbac) HasPermission(typ model.ProjectRBACResource_ResourceType, action model.ProjectRBACPolicy_Action) bool {
	for _, v := range r.Roles {
		if v.HasPermission(typ, action) {
			return true
		}
	}
	return false
}

func (r *rbac) FilterByNames(names []string) *rbac {
	roles := make([]*model.ProjectRBACRole, 0, len(names))
	rs := make(map[string]*model.ProjectRBACRole, len(r.Roles))
	for _, v := range r.Roles {
		rs[v.Name] = v
	}
	for _, n := range names {
		if v, ok := rs[n]; ok {
			roles = append(roles, v)
		}
	}
	r.Roles = roles
	return r
}

// Authorize checks whether a role is enough for given gRPC method or not.
func (a *authorizer) Authorize(ctx context.Context, method string, r model.Role) bool {
	rbac, err := a.getRBAC(ctx, r.ProjectId)
	if err != nil {
		return false
	}
	rbac.FilterByNames(r.ProjectRbacRoles)

	switch method {
	case "/grpc.service.webservice.WebService/GetCommand":
		return true
	case "/grpc.service.webservice.WebService/GenerateAPIKey":
		return rbac.HasPermission(model.ProjectRBACResource_API_KEY, model.ProjectRBACPolicy_CREATE)
	case "/grpc.service.webservice.WebService/DisableAPIKey":
		return rbac.HasPermission(model.ProjectRBACResource_API_KEY, model.ProjectRBACPolicy_UPDATE)
	case "/grpc.service.webservice.WebService/ListAPIKeys":
		return rbac.HasPermission(model.ProjectRBACResource_API_KEY, model.ProjectRBACPolicy_LIST)
	case "/grpc.service.webservice.WebService/AddApplication":
		return rbac.HasPermission(model.ProjectRBACResource_APPLICATION, model.ProjectRBACPolicy_CREATE)
	case "/grpc.service.webservice.WebService/UpdateApplication":
		return rbac.HasPermission(model.ProjectRBACResource_APPLICATION, model.ProjectRBACPolicy_UPDATE)
	case "/grpc.service.webservice.WebService/EnableApplication":
		return rbac.HasPermission(model.ProjectRBACResource_APPLICATION, model.ProjectRBACPolicy_UPDATE)
	case "/grpc.service.webservice.WebService/DisableApplication":
		return rbac.HasPermission(model.ProjectRBACResource_APPLICATION, model.ProjectRBACPolicy_UPDATE)
	case "/grpc.service.webservice.WebService/DeleteApplication":
		return rbac.HasPermission(model.ProjectRBACResource_APPLICATION, model.ProjectRBACPolicy_DELETE)
	case "/grpc.service.webservice.WebService/ListApplications":
		return rbac.HasPermission(model.ProjectRBACResource_APPLICATION, model.ProjectRBACPolicy_LIST)
	case "/grpc.service.webservice.WebService/SyncApplication":
		return rbac.HasPermission(model.ProjectRBACResource_APPLICATION, model.ProjectRBACPolicy_UPDATE)
	case "/grpc.service.webservice.WebService/GetApplication":
		return rbac.HasPermission(model.ProjectRBACResource_APPLICATION, model.ProjectRBACPolicy_GET)
	case "/grpc.service.webservice.WebService/GenerateApplicationSealedSecret":
		return rbac.HasPermission(model.ProjectRBACResource_APPLICATION, model.ProjectRBACPolicy_UPDATE)
	case "/grpc.service.webservice.WebService/ListUnregisteredApplications":
		return rbac.HasPermission(model.ProjectRBACResource_APPLICATION, model.ProjectRBACPolicy_LIST)
	case "/grpc.service.webservice.WebService/GetApplicationLiveState":
		return rbac.HasPermission(model.ProjectRBACResource_APPLICATION, model.ProjectRBACPolicy_GET)
	case "/grpc.service.webservice.WebService/ListDeployments":
		return rbac.HasPermission(model.ProjectRBACResource_DEPLOYMENT, model.ProjectRBACPolicy_LIST)
	case "/grpc.service.webservice.WebService/GetDeployment":
		return rbac.HasPermission(model.ProjectRBACResource_DEPLOYMENT, model.ProjectRBACPolicy_GET)
	case "/grpc.service.webservice.WebService/GetStageLog":
		return rbac.HasPermission(model.ProjectRBACResource_DEPLOYMENT, model.ProjectRBACPolicy_GET)
	case "/grpc.service.webservice.WebService/CancelDeployment":
		return rbac.HasPermission(model.ProjectRBACResource_DEPLOYMENT, model.ProjectRBACPolicy_UPDATE)
	case "/grpc.service.webservice.WebService/SkipStage":
		return rbac.HasPermission(model.ProjectRBACResource_DEPLOYMENT, model.ProjectRBACPolicy_UPDATE)
	case "/grpc.service.webservice.WebService/ApproveStage":
		return rbac.HasPermission(model.ProjectRBACResource_DEPLOYMENT, model.ProjectRBACPolicy_UPDATE)
	case "/grpc.service.webservice.WebService/ListEvents":
		return rbac.HasPermission(model.ProjectRBACResource_EVENT, model.ProjectRBACPolicy_LIST)
	case "/grpc.service.webservice.WebService/GetInsightData":
		return rbac.HasPermission(model.ProjectRBACResource_INSIGHT, model.ProjectRBACPolicy_GET)
	case "/grpc.service.webservice.WebService/GetInsightApplicationCount":
		return rbac.HasPermission(model.ProjectRBACResource_INSIGHT, model.ProjectRBACPolicy_GET)
	case "/grpc.service.webservice.WebService/RegisterPiped":
		return rbac.HasPermission(model.ProjectRBACResource_PIPED, model.ProjectRBACPolicy_CREATE)
	case "/grpc.service.webservice.WebService/UpdatePiped":
		return rbac.HasPermission(model.ProjectRBACResource_PIPED, model.ProjectRBACPolicy_UPDATE)
	case "/grpc.service.webservice.WebService/RecreatePipedKey":
		return rbac.HasPermission(model.ProjectRBACResource_PIPED, model.ProjectRBACPolicy_UPDATE)
	case "/grpc.service.webservice.WebService/DeleteOldPipedKeys":
		return rbac.HasPermission(model.ProjectRBACResource_PIPED, model.ProjectRBACPolicy_UPDATE)
	case "/grpc.service.webservice.WebService/EnablePiped":
		return rbac.HasPermission(model.ProjectRBACResource_PIPED, model.ProjectRBACPolicy_UPDATE)
	case "/grpc.service.webservice.WebService/DisablePiped":
		return rbac.HasPermission(model.ProjectRBACResource_PIPED, model.ProjectRBACPolicy_UPDATE)
	case "/grpc.service.webservice.WebService/ListPipeds":
		return rbac.HasPermission(model.ProjectRBACResource_PIPED, model.ProjectRBACPolicy_LIST)
	case "/grpc.service.webservice.WebService/GetPiped":
		return rbac.HasPermission(model.ProjectRBACResource_PIPED, model.ProjectRBACPolicy_GET)
	case "/grpc.service.webservice.WebService/UpdatePipedDesiredVersion":
		return rbac.HasPermission(model.ProjectRBACResource_PIPED, model.ProjectRBACPolicy_UPDATE)
	case "/grpc.service.webservice.WebService/RestartPiped":
		return rbac.HasPermission(model.ProjectRBACResource_PIPED, model.ProjectRBACPolicy_UPDATE)
	case "/grpc.service.webservice.WebService/ListReleasedVersions":
		return rbac.HasPermission(model.ProjectRBACResource_PIPED, model.ProjectRBACPolicy_LIST)
	case "/grpc.service.webservice.WebService/GetProject":
		return rbac.HasPermission(model.ProjectRBACResource_PROJECT, model.ProjectRBACPolicy_GET)
	case "/grpc.service.webservice.WebService/UpdateProjectStaticAdmin":
		return rbac.HasPermission(model.ProjectRBACResource_PROJECT, model.ProjectRBACPolicy_UPDATE)
	case "/grpc.service.webservice.WebService/EnableStaticAdmin":
		return rbac.HasPermission(model.ProjectRBACResource_PROJECT, model.ProjectRBACPolicy_UPDATE)
	case "/grpc.service.webservice.WebService/DisableStaticAdmin":
		return rbac.HasPermission(model.ProjectRBACResource_PROJECT, model.ProjectRBACPolicy_UPDATE)
	case "/grpc.service.webservice.WebService/UpdateProjectSSOConfig":
		return rbac.HasPermission(model.ProjectRBACResource_PROJECT, model.ProjectRBACPolicy_UPDATE)
	case "/grpc.service.webservice.WebService/UpdateProjectRBACConfig":
		return rbac.HasPermission(model.ProjectRBACResource_PROJECT, model.ProjectRBACPolicy_UPDATE)
	case "/grpc.service.webservice.WebService/GetMe":
		return rbac.HasPermission(model.ProjectRBACResource_PROJECT, model.ProjectRBACPolicy_GET)
	}
	return false
}
