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

package httpapi

import (
	"context"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
	"golang.org/x/net/xsrftoken"

	"github.com/pipe-cd/pipecd/pkg/jwt"
	"github.com/pipe-cd/pipecd/pkg/model"
)

// handleSSOLogin is called when an user requested to login via SSO.
func (h *authHandler) handleSSOLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	// Validate request's payload.
	if r.Method != http.MethodPost {
		h.handleError(w, r, "Method not allowed", nil)
		return
	}
	projectID := r.FormValue(projectFormKey)
	if projectID == "" {
		h.handleError(w, r, "Missing project id", nil)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	proj, err := h.projectGetter.Get(ctx, projectID)
	if err != nil {
		h.handleError(w, r, fmt.Sprintf("Unable to find project %s", projectID), err)
		return
	}

	sso, shared, err := h.findSSOConfig(proj)
	if err != nil {
		h.handleError(w, r, fmt.Sprintf("Invalid SSO configuration: %v", err), nil)
		return
	}

	if !shared {
		if err := sso.Decrypt(h.decrypter); err != nil {
			h.handleError(w, r, "Failed to decrypt SSO configuration", err)
			return
		}
	}

	var (
		stateToken = xsrftoken.Generate(h.stateKey, "", "")
		state      = hex.EncodeToString([]byte(stateToken))
	)
	authURL, err := sso.GenerateAuthCodeURL(proj.Id, h.callbackURL, state)
	if err != nil {
		h.handleError(w, r, "Internal error", err)
		return
	}

	http.SetCookie(w, makeStateCookie(state, h.secureCookie))
	http.Redirect(w, r, authURL, http.StatusFound)
}

// handleStaticAdminLogin is called when an user requested to login as a static admin.
func (h *authHandler) handleStaticAdminLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	// Validate request's payload.
	if r.Method != http.MethodPost {
		h.handleError(w, r, "Method not allowed", nil)
		return
	}
	projectID := r.FormValue(projectFormKey)
	if projectID == "" {
		h.handleError(w, r, "Missing project id", nil)
		return
	}
	username := r.FormValue(usernameFormKey)
	if username == "" {
		h.handleError(w, r, "Missing username", nil)
		return
	}
	password := r.FormValue(passwordFormKey)
	if password == "" {
		h.handleError(w, r, "Missing password", nil)
		return
	}

	var admin *model.ProjectStaticUser
	if p, ok := h.projectsInConfig[projectID]; ok {
		admin = &model.ProjectStaticUser{
			Username:     p.StaticAdmin.Username,
			PasswordHash: p.StaticAdmin.PasswordHash,
		}
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		proj, err := h.projectGetter.Get(ctx, projectID)
		if err != nil {
			h.handleError(w, r, fmt.Sprintf("Unable to find project: %s", projectID), err)
			return
		}
		if proj.StaticAdminDisabled {
			h.handleError(w, r, "Static admin is disabling", nil)
			return
		}
		admin = proj.StaticAdmin
	}

	if err := admin.Auth(username, password); err != nil {
		h.handleError(w, r, "Unable to login", err)
		return
	}

	claims := jwt.NewClaims(
		admin.Username,
		"",
		defaultTokenTTL,
		model.Role{
			ProjectId:        projectID,
			ProjectRbacRoles: []string{model.BuiltinRBACRoleAdmin.String()},
		},
	)
	signedToken, err := h.signer.Sign(claims)
	if err != nil {
		h.handleError(w, r, "Internal error", err)
		return
	}

	h.logger.Info("a new user has been logged in",
		zap.String("user", admin.Username),
		zap.String("project-id", projectID),
		zap.String("project-role", model.BuiltinRBACRoleAdmin.String()),
	)
	http.SetCookie(w, makeTokenCookie(signedToken, h.secureCookie))
	http.Redirect(w, r, rootPath, http.StatusFound)
}
