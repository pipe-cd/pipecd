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
	w := datastore.WebCommander
	return &authorizer{
		projectStore:     datastore.NewProjectStore(ds, w),
		rbacCache:        memorycache.NewTTLCache(ctx, 10*time.Minute, 5*time.Minute),
		projectsInConfig: projects,
		logger:           logger.Named("authorizer"),
	}
}

func (a *authorizer) getRBACRoles(ctx context.Context, projectID string) ([]*model.ProjectRBACRole, error) {
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
	roles, err := a.getRBACRoles(ctx, r.ProjectId)
	if err != nil {
		a.logger.Error("failed to get the rbac",
			zap.String("project", r.ProjectId),
			zap.Error(err),
		)
		return false
	}

	switch method {
	case "/grpc.service.webservice.WebService/GetCommand":
		return true
	case "/grpc.service.webservice.WebService/GenerateAPIKey":
		return model.VerifyRBACPermission(model.ProjectRBACResource_API_KEY, model.ProjectRBACPolicy_CREATE, r.ProjectRbacRoles, roles)
	case "/grpc.service.webservice.WebService/DisableAPIKey":
		return model.VerifyRBACPermission(model.ProjectRBACResource_API_KEY, model.ProjectRBACPolicy_UPDATE, r.ProjectRbacRoles, roles)
	case "/grpc.service.webservice.WebService/ListAPIKeys":
		return model.VerifyRBACPermission(model.ProjectRBACResource_API_KEY, model.ProjectRBACPolicy_LIST, r.ProjectRbacRoles, roles)
	case "/grpc.service.webservice.WebService/AddApplication":
		return model.VerifyRBACPermission(model.ProjectRBACResource_APPLICATION, model.ProjectRBACPolicy_CREATE, r.ProjectRbacRoles, roles)
	case "/grpc.service.webservice.WebService/UpdateApplication":
		return model.VerifyRBACPermission(model.ProjectRBACResource_APPLICATION, model.ProjectRBACPolicy_UPDATE, r.ProjectRbacRoles, roles)
	case "/grpc.service.webservice.WebService/EnableApplication":
		return model.VerifyRBACPermission(model.ProjectRBACResource_APPLICATION, model.ProjectRBACPolicy_UPDATE, r.ProjectRbacRoles, roles)
	case "/grpc.service.webservice.WebService/DisableApplication":
		return model.VerifyRBACPermission(model.ProjectRBACResource_APPLICATION, model.ProjectRBACPolicy_UPDATE, r.ProjectRbacRoles, roles)
	case "/grpc.service.webservice.WebService/DeleteApplication":
		return model.VerifyRBACPermission(model.ProjectRBACResource_APPLICATION, model.ProjectRBACPolicy_DELETE, r.ProjectRbacRoles, roles)
	case "/grpc.service.webservice.WebService/ListApplications":
		return model.VerifyRBACPermission(model.ProjectRBACResource_APPLICATION, model.ProjectRBACPolicy_LIST, r.ProjectRbacRoles, roles)
	case "/grpc.service.webservice.WebService/SyncApplication":
		return model.VerifyRBACPermission(model.ProjectRBACResource_APPLICATION, model.ProjectRBACPolicy_UPDATE, r.ProjectRbacRoles, roles)
	case "/grpc.service.webservice.WebService/GetApplication":
		return model.VerifyRBACPermission(model.ProjectRBACResource_APPLICATION, model.ProjectRBACPolicy_GET, r.ProjectRbacRoles, roles)
	case "/grpc.service.webservice.WebService/GenerateApplicationSealedSecret":
		return model.VerifyRBACPermission(model.ProjectRBACResource_APPLICATION, model.ProjectRBACPolicy_UPDATE, r.ProjectRbacRoles, roles)
	case "/grpc.service.webservice.WebService/ListUnregisteredApplications":
		return model.VerifyRBACPermission(model.ProjectRBACResource_APPLICATION, model.ProjectRBACPolicy_LIST, r.ProjectRbacRoles, roles)
	case "/grpc.service.webservice.WebService/GetApplicationLiveState":
		return model.VerifyRBACPermission(model.ProjectRBACResource_APPLICATION, model.ProjectRBACPolicy_GET, r.ProjectRbacRoles, roles)
	case "/grpc.service.webservice.WebService/ListDeployments":
		return model.VerifyRBACPermission(model.ProjectRBACResource_DEPLOYMENT, model.ProjectRBACPolicy_LIST, r.ProjectRbacRoles, roles)
	case "/grpc.service.webservice.WebService/GetDeployment":
		return model.VerifyRBACPermission(model.ProjectRBACResource_DEPLOYMENT, model.ProjectRBACPolicy_GET, r.ProjectRbacRoles, roles)
	case "/grpc.service.webservice.WebService/GetStageLog":
		return model.VerifyRBACPermission(model.ProjectRBACResource_DEPLOYMENT, model.ProjectRBACPolicy_GET, r.ProjectRbacRoles, roles)
	case "/grpc.service.webservice.WebService/CancelDeployment":
		return model.VerifyRBACPermission(model.ProjectRBACResource_DEPLOYMENT, model.ProjectRBACPolicy_UPDATE, r.ProjectRbacRoles, roles)
	case "/grpc.service.webservice.WebService/SkipStage":
		return model.VerifyRBACPermission(model.ProjectRBACResource_DEPLOYMENT, model.ProjectRBACPolicy_UPDATE, r.ProjectRbacRoles, roles)
	case "/grpc.service.webservice.WebService/ApproveStage":
		return model.VerifyRBACPermission(model.ProjectRBACResource_DEPLOYMENT, model.ProjectRBACPolicy_UPDATE, r.ProjectRbacRoles, roles)
	case "/grpc.service.webservice.WebService/ListEvents":
		return model.VerifyRBACPermission(model.ProjectRBACResource_EVENT, model.ProjectRBACPolicy_LIST, r.ProjectRbacRoles, roles)
	case "/grpc.service.webservice.WebService/GetInsightData":
		return model.VerifyRBACPermission(model.ProjectRBACResource_INSIGHT, model.ProjectRBACPolicy_GET, r.ProjectRbacRoles, roles)
	case "/grpc.service.webservice.WebService/GetInsightApplicationCount":
		return model.VerifyRBACPermission(model.ProjectRBACResource_INSIGHT, model.ProjectRBACPolicy_GET, r.ProjectRbacRoles, roles)
	case "/grpc.service.webservice.WebService/RegisterPiped":
		return model.VerifyRBACPermission(model.ProjectRBACResource_PIPED, model.ProjectRBACPolicy_CREATE, r.ProjectRbacRoles, roles)
	case "/grpc.service.webservice.WebService/UpdatePiped":
		return model.VerifyRBACPermission(model.ProjectRBACResource_PIPED, model.ProjectRBACPolicy_UPDATE, r.ProjectRbacRoles, roles)
	case "/grpc.service.webservice.WebService/RecreatePipedKey":
		return model.VerifyRBACPermission(model.ProjectRBACResource_PIPED, model.ProjectRBACPolicy_UPDATE, r.ProjectRbacRoles, roles)
	case "/grpc.service.webservice.WebService/DeleteOldPipedKeys":
		return model.VerifyRBACPermission(model.ProjectRBACResource_PIPED, model.ProjectRBACPolicy_UPDATE, r.ProjectRbacRoles, roles)
	case "/grpc.service.webservice.WebService/EnablePiped":
		return model.VerifyRBACPermission(model.ProjectRBACResource_PIPED, model.ProjectRBACPolicy_UPDATE, r.ProjectRbacRoles, roles)
	case "/grpc.service.webservice.WebService/DisablePiped":
		return model.VerifyRBACPermission(model.ProjectRBACResource_PIPED, model.ProjectRBACPolicy_UPDATE, r.ProjectRbacRoles, roles)
	case "/grpc.service.webservice.WebService/ListPipeds":
		return model.VerifyRBACPermission(model.ProjectRBACResource_PIPED, model.ProjectRBACPolicy_LIST, r.ProjectRbacRoles, roles)
	case "/grpc.service.webservice.WebService/GetPiped":
		return model.VerifyRBACPermission(model.ProjectRBACResource_PIPED, model.ProjectRBACPolicy_GET, r.ProjectRbacRoles, roles)
	case "/grpc.service.webservice.WebService/UpdatePipedDesiredVersion":
		return model.VerifyRBACPermission(model.ProjectRBACResource_PIPED, model.ProjectRBACPolicy_UPDATE, r.ProjectRbacRoles, roles)
	case "/grpc.service.webservice.WebService/RestartPiped":
		return model.VerifyRBACPermission(model.ProjectRBACResource_PIPED, model.ProjectRBACPolicy_UPDATE, r.ProjectRbacRoles, roles)
	case "/grpc.service.webservice.WebService/ListReleasedVersions":
		return model.VerifyRBACPermission(model.ProjectRBACResource_PIPED, model.ProjectRBACPolicy_LIST, r.ProjectRbacRoles, roles)
	case "/grpc.service.webservice.WebService/GetProject":
		return model.VerifyRBACPermission(model.ProjectRBACResource_PROJECT, model.ProjectRBACPolicy_GET, r.ProjectRbacRoles, roles)
	case "/grpc.service.webservice.WebService/UpdateProjectStaticAdmin":
		return model.VerifyRBACPermission(model.ProjectRBACResource_PROJECT, model.ProjectRBACPolicy_UPDATE, r.ProjectRbacRoles, roles)
	case "/grpc.service.webservice.WebService/EnableStaticAdmin":
		return model.VerifyRBACPermission(model.ProjectRBACResource_PROJECT, model.ProjectRBACPolicy_UPDATE, r.ProjectRbacRoles, roles)
	case "/grpc.service.webservice.WebService/DisableStaticAdmin":
		return model.VerifyRBACPermission(model.ProjectRBACResource_PROJECT, model.ProjectRBACPolicy_UPDATE, r.ProjectRbacRoles, roles)
	case "/grpc.service.webservice.WebService/UpdateProjectSSOConfig":
		return model.VerifyRBACPermission(model.ProjectRBACResource_PROJECT, model.ProjectRBACPolicy_UPDATE, r.ProjectRbacRoles, roles)
	case "/grpc.service.webservice.WebService/UpdateProjectRBACConfig":
		return model.VerifyRBACPermission(model.ProjectRBACResource_PROJECT, model.ProjectRBACPolicy_UPDATE, r.ProjectRbacRoles, roles)
	case "/grpc.service.webservice.WebService/GetMe":
		return model.VerifyRBACPermission(model.ProjectRBACResource_PROJECT, model.ProjectRBACPolicy_GET, r.ProjectRbacRoles, roles)
	case "/grpc.service.webservice.WebService/AddProjectRBACRole":
		return model.VerifyRBACPermission(model.ProjectRBACResource_PROJECT, model.ProjectRBACPolicy_UPDATE, r.ProjectRbacRoles, roles)
	case "/grpc.service.webservice.WebService/UpdateProjectRBACRole":
		return model.VerifyRBACPermission(model.ProjectRBACResource_PROJECT, model.ProjectRBACPolicy_UPDATE, r.ProjectRbacRoles, roles)
	case "/grpc.service.webservice.WebService/DeleteProjectRBACRole":
		return model.VerifyRBACPermission(model.ProjectRBACResource_PROJECT, model.ProjectRBACPolicy_UPDATE, r.ProjectRbacRoles, roles)
	case "/grpc.service.webservice.WebService/AddProjectUserGroup":
		return model.VerifyRBACPermission(model.ProjectRBACResource_PROJECT, model.ProjectRBACPolicy_UPDATE, r.ProjectRbacRoles, roles)
	case "/grpc.service.webservice.WebService/DeleteProjectUserGroup":
		return model.VerifyRBACPermission(model.ProjectRBACResource_PROJECT, model.ProjectRBACPolicy_UPDATE, r.ProjectRbacRoles, roles)
	}
	return false
}
