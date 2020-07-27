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

package webservice

import (
	"github.com/pipe-cd/pipe/pkg/model"
	"github.com/pipe-cd/pipe/pkg/rpc/rpcauth"
)

type authorizer struct{}

// NewRBACAuthorizer returns an RBACAuthorizer object for checking requested method based on RBAC.
func NewRBACAuthorizer() rpcauth.RBACAuthorizer {
	return &authorizer{}
}

func isAdmin(r model.Role) bool {
	return r.ProjectRole == model.Role_ADMIN
}

func isEditor(r model.Role) bool {
	return r.ProjectRole == model.Role_EDITOR
}

func isViewer(r model.Role) bool {
	return r.ProjectRole == model.Role_VIEWER
}

// Authorize checks whether a role is enough for given gRPC method or not.
// Todo: Auto generate this file from protobuf.
func (a *authorizer) Authorize(method string, r model.Role) bool {
	switch method {
	case "/pipe.api.service.WebAPI/AddEnvironment":
		return isAdmin(r)
	case "/pipe.api.service.WebAPI/UpdateEnvironmentDesc":
		return isAdmin(r)
	case "/pipe.api.service.WebAPI/RegisterPiped":
		return isAdmin(r)
	case "/pipe.api.service.WebAPI/RecreatePipedKey":
		return isAdmin(r)
	case "/pipe.api.service.WebAPI/EnablePiped":
		return isAdmin(r)
	case "/pipe.api.service.WebAPI/DisablePiped":
		return isAdmin(r)
	case "/pipe.api.service.WebAPI/AddApplication":
		return isAdmin(r)
	case "/pipe.api.service.WebAPI/EnableApplication":
		return isAdmin(r)
	case "/pipe.api.service.WebAPI/DisableApplication":
		return isAdmin(r)
	case "/pipe.api.service.WebAPI/ListEnvironments":
		return isAdmin(r) || isEditor(r)
	case "/pipe.api.service.WebAPI/ListPipeds":
		return isAdmin(r) || isEditor(r)
	case "/pipe.api.service.WebAPI/GetPiped":
		return isAdmin(r) || isEditor(r)
	case "/pipe.api.service.WebAPI/ListApplications":
		return isAdmin(r) || isEditor(r)
	case "/pipe.api.service.WebAPI/SyncApplication":
		return isAdmin(r) || isEditor(r)
	case "/pipe.api.service.WebAPI/GetApplication":
		return isAdmin(r) || isEditor(r)
	case "/pipe.api.service.WebAPI/ListDeployments":
		return isAdmin(r) || isEditor(r)
	case "/pipe.api.service.WebAPI/GetDeployment":
		return isAdmin(r) || isEditor(r)
	case "/pipe.api.service.WebAPI/GetStageLog":
		return isAdmin(r) || isEditor(r)
	case "/pipe.api.service.WebAPI/CancelDeployment":
		return isAdmin(r) || isEditor(r)
	case "/pipe.api.service.WebAPI/ApproveStage":
		return isAdmin(r) || isEditor(r)
	case "/pipe.api.service.WebAPI/GetApplicationLiveState":
		return isAdmin(r) || isEditor(r)
	case "/pipe.api.service.WebAPI/GetProject":
		return isAdmin(r) || isEditor(r)
	case "/pipe.api.service.WebAPI/GetCommand":
		return isAdmin(r) || isEditor(r)
	case "/pipe.api.service.WebAPI/ListDeploymentConfigTemplates":
		return isAdmin(r) || isEditor(r)
	case "/pipe.api.service.WebAPI/GetMe":
		return isAdmin(r) || isEditor(r) || isViewer(r)
	}
	return false
}
