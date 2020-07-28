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
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
	"golang.org/x/net/xsrftoken"

	"github.com/pipe-cd/pipe/pkg/jwt"
	"github.com/pipe-cd/pipe/pkg/model"
)

// handleLogin is called when user request to login PipeCD.
func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	if r.Method != http.MethodPost {
		handleError(w, r, rootPath, "method not allowed", h.logger, fmt.Errorf("method not allowed: %v", r.Method))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	proj, err := h.getProject(ctx, r)
	if err != nil {
		handleError(w, r, rootPath, "wrong project", h.logger, err)
		return
	}

	stateToken := xsrftoken.Generate(h.stateSeed, "", "")
	state := hex.EncodeToString([]byte(stateToken))
	authURL, err := proj.Sso.GenerateAuthCodeURL(proj.Id, h.apiURL, callbackPath, state)
	if err != nil {
		handleError(w, r, rootPath, "internal error", h.logger, err)
		return
	}
	http.SetCookie(w, makeStateCookie(state))
	http.Redirect(w, r, authURL, http.StatusFound)
}

// handleStaticLogin is called when user request to login PipeCD as a static user.
func (h *Handler) handleStaticLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	if r.Method != http.MethodPost {
		handleError(w, r, rootPath, "method not allowed", h.logger, fmt.Errorf("method not allowed: %v", r.Method))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	proj, err := h.getProject(ctx, r)
	if err != nil {
		handleError(w, r, rootPath, "wrong project", h.logger, err)
		return
	}
	if proj.StaticAdminDisabled {
		msg := "static login is disabled"
		handleError(w, r, rootPath, msg, h.logger, fmt.Errorf(msg))
		return
	}

	if err := proj.StaticAdmin.Auth(r.FormValue(usernameFormKey), r.FormValue(passwordFormKey)); err != nil {
		handleError(w, r, rootPath, "login failed", h.logger, err)
		return
	}
	claims := jwt.NewClaims(
		proj.StaticAdmin.Username,
		"",
		defaultTokenTTL,
		model.Role{
			ProjectId:   proj.Id,
			ProjectRole: model.Role_ADMIN,
		},
	)
	signedToken, err := h.signer.Sign(claims)
	if err != nil {
		handleError(w, r, rootPath, "internal error", h.logger, err)
		return
	}
	http.SetCookie(w, makeTokenCookie(signedToken))

	h.logger.Info("a new user has been logged in",
		zap.String("user", proj.StaticAdmin.Username),
		zap.String("project-id", proj.Id),
		zap.String("project-role", model.Role_ADMIN.String()),
	)
	http.Redirect(w, r, rootPath, http.StatusFound)
}
