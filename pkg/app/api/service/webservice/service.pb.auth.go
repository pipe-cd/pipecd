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
	"github.com/kapetaniosci/pipe/pkg/role"
	"github.com/kapetaniosci/pipe/pkg/rpc/rpcauth"
)

type authorizer struct{}

// NewRBACAuthorizer returns an RBACAuthorizer object for checking requested method based on RBAC.
func NewRBACAuthorizer() rpcauth.RBACAuthorizer {
	return &authorizer{}
}

func isOwner(r role.Role) bool {
	return r.Owner
}

func isAdmin(r role.Role) bool {
	return r.ProjectRole == role.Role_ADMIN
}

func isEditor(r role.Role) bool {
	return r.ProjectRole == role.Role_EDITOR
}

func isViewer(r role.Role) bool {
	return r.ProjectRole == role.Role_VIEWER
}

// Authorize checks whether a role is enough for given gRPC method or not.
// Todo: Auto generate this file from protobuf.
func (a *authorizer) Authorize(method string, r role.Role) bool {
	switch method {
	case "/pipe.api.service.WebAPI/RegisterRunner":
		return isAdmin(r)
	}
	return false
}
