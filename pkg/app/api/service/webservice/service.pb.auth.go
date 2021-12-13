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
	case "/pipe.api.service.webservice.WebService/AddEnvironment":
		return isAdmin(r)
	case "/pipe.api.service.webservice.WebService/UpdateEnvironmentDesc":
		return isAdmin(r)
	case "/pipe.api.service.webservice.WebService/DeleteEnvironment":
		return isAdmin(r)
	case "/pipe.api.service.webservice.WebService/RegisterPiped":
		return isAdmin(r)
	case "/pipe.api.service.webservice.WebService/UpdatePiped":
		return isAdmin(r)
	case "/pipe.api.service.webservice.WebService/RecreatePipedKey":
		return isAdmin(r)
	case "/pipe.api.service.webservice.WebService/DeleteOldPipedKeys":
		return isAdmin(r)
	case "/pipe.api.service.webservice.WebService/EnablePiped":
		return isAdmin(r)
	case "/pipe.api.service.webservice.WebService/DisablePiped":
		return isAdmin(r)
	case "/pipe.api.service.webservice.WebService/UpdatePipedDesiredVersion":
		return isAdmin(r)
	case "/pipe.api.service.webservice.WebService/UpdateProjectStaticAdmin":
		return isAdmin(r)
	case "/pipe.api.service.webservice.WebService/EnableStaticAdmin":
		return isAdmin(r)
	case "/pipe.api.service.webservice.WebService/DisableStaticAdmin":
		return isAdmin(r)
	case "/pipe.api.service.webservice.WebService/UpdateProjectSSOConfig":
		return isAdmin(r)
	case "/pipe.api.service.webservice.WebService/UpdateProjectRBACConfig":
		return isAdmin(r)
	case "/pipe.api.service.webservice.WebService/GenerateAPIKey":
		return isAdmin(r)
	case "/pipe.api.service.webservice.WebService/DisableAPIKey":
		return isAdmin(r)
	case "/pipe.api.service.webservice.WebService/ListAPIKeys":
		return isAdmin(r)

	case "/pipe.api.service.webservice.WebService/AddApplication":
		return isAdmin(r) || isEditor(r)
	case "/pipe.api.service.webservice.WebService/UpdateApplication":
		return isAdmin(r) || isEditor(r)
	case "/pipe.api.service.webservice.WebService/UpdateApplicationDescription":
		return isAdmin(r) || isEditor(r)
	case "/pipe.api.service.webservice.WebService/EnableApplication":
		return isAdmin(r) || isEditor(r)
	case "/pipe.api.service.webservice.WebService/DisableApplication":
		return isAdmin(r) || isEditor(r)
	case "/pipe.api.service.webservice.WebService/DeleteApplication":
		return isAdmin(r) || isEditor(r)
	case "/pipe.api.service.webservice.WebService/SyncApplication":
		return isAdmin(r) || isEditor(r)
	case "/pipe.api.service.webservice.WebService/CancelDeployment":
		return isAdmin(r) || isEditor(r)
	case "/pipe.api.service.webservice.WebService/ApproveStage":
		return isAdmin(r) || isEditor(r)
	case "/pipe.api.service.webservice.WebService/GenerateApplicationSealedSecret":
		return isAdmin(r) || isEditor(r)

	case "/pipe.api.service.webservice.WebService/GetApplicationLiveState":
		return isAdmin(r) || isEditor(r) || isViewer(r)
	case "/pipe.api.service.webservice.WebService/GetProject":
		return isAdmin(r) || isEditor(r) || isViewer(r)
	case "/pipe.api.service.webservice.WebService/GetCommand":
		return isAdmin(r) || isEditor(r) || isViewer(r)
	case "/pipe.api.service.webservice.WebService/ListEnvironments":
		return isAdmin(r) || isEditor(r) || isViewer(r)
	case "/pipe.api.service.webservice.WebService/ListPipeds":
		return isAdmin(r) || isEditor(r) || isViewer(r)
	case "/pipe.api.service.webservice.WebService/GetPiped":
		return isAdmin(r) || isEditor(r) || isViewer(r)
	case "/pipe.api.service.webservice.WebService/ListApplications":
		return isAdmin(r) || isEditor(r) || isViewer(r)
	case "/pipe.api.service.webservice.WebService/GetApplication":
		return isAdmin(r) || isEditor(r) || isViewer(r)
	case "/pipe.api.service.webservice.WebService/ListDeployments":
		return isAdmin(r) || isEditor(r) || isViewer(r)
	case "/pipe.api.service.webservice.WebService/GetDeployment":
		return isAdmin(r) || isEditor(r) || isViewer(r)
	case "/pipe.api.service.webservice.WebService/GetStageLog":
		return isAdmin(r) || isEditor(r) || isViewer(r)
	case "/pipe.api.service.webservice.WebService/GetMe":
		return isAdmin(r) || isEditor(r) || isViewer(r)
	case "/pipe.api.service.webservice.WebService/GetInsightData":
		return isAdmin(r) || isEditor(r) || isViewer(r)
	case "/pipe.api.service.webservice.WebService/GetInsightApplicationCount":
		return isAdmin(r) || isEditor(r) || isViewer(r)
	case "/pipe.api.service.webservice.WebService/ListUnregisteredApplications":
		return isAdmin(r) || isEditor(r) || isViewer(r)
	}

	return false
}
