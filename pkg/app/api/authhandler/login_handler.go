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

package authhandler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipe/pkg/jwt"
	"github.com/pipe-cd/pipe/pkg/model"
	"github.com/pipe-cd/pipe/pkg/role"
)

// handleLogin is called when user request to login PipeCD.
func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	projectID := r.FormValue("projectID")
	if projectID == "" {
		msg := "project id must be specified"
		serverError(w, r, "/", msg, h.logger, fmt.Errorf(msg))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	proj, err := h.projectStore.GetProject(ctx, projectID)
	if err != nil {
		serverError(w, r, "/", fmt.Sprintf("unabled to get project: %s", projectID), h.logger, err)
		return
	}

	if !proj.StaticAdminDisabled && r.URL.Path == passwordLoginPath {
		h.handleStaticAdminLogin(proj, w, r)
	}

	serverError(w, r, "/", fmt.Sprintf("login failed"), h.logger, fmt.Errorf("login failed"))
}

func (h *Handler) handleStaticAdminLogin(proj *model.Project, w http.ResponseWriter, r *http.Request) {
	if err := proj.StaticAdmin.Auth(r.FormValue("username"), r.FormValue("password")); err != nil {
		serverError(w, r, "/", fmt.Sprintf("login failed"), h.logger, err)
	}
	claims := jwt.NewClaims(
		proj.StaticAdmin.Username,
		"",
		defaultTokenTTL,
		role.Role{
			ProjectId:   proj.Id,
			ProjectRole: role.Role_ADMIN,
		},
	)
	signedToken, err := h.signer.Sign(claims)
	if err != nil {
		serverError(w, r, "/", err.Error(), h.logger, err)
		return
	}
	http.SetCookie(w, makeTokenCookie(signedToken))

	h.logger.Info("user logged in",
		zap.String("user", proj.StaticAdmin.Username),
		zap.String("project-id", proj.Id),
		zap.String("project-role", role.Role_ADMIN.String()),
	)
	http.Redirect(w, r, "/", http.StatusFound)
}
