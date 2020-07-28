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
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
	"golang.org/x/net/xsrftoken"

	"github.com/pipe-cd/pipe/pkg/jwt"
	"github.com/pipe-cd/pipe/pkg/model"
)

func (h *Handler) handleCallback(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	err := checkState(r, h.stateKey)
	if err != nil {
		handleError(w, r, rootPath, "unauthorized access", h.logger, err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	proj, err := h.getProject(ctx, r)
	if err != nil {
		handleError(w, r, rootPath, "wrong project", h.logger, err)
		return
	}

	user, err := proj.Sso.GenerateUserInfo(ctx, r.FormValue(authCodeFormKey))
	if err != nil {
		handleError(w, r, rootPath, "failed to generate the user information", h.logger, err)
		return
	}
	claims := jwt.NewClaims(
		user.Username,
		user.AvatarURL,
		defaultTokenTTL,
		model.Role{
			ProjectId:   proj.Id,
			ProjectRole: user.Role,
		},
	)
	signedToken, err := h.signer.Sign(claims)
	if err != nil {
		handleError(w, r, rootPath, "internal error", h.logger, err)
		return
	}
	http.SetCookie(w, makeTokenCookie(signedToken))
	http.SetCookie(w, makeExpiredStateCookie())

	h.logger.Info("user logged in",
		zap.String("user", user.Username),
		zap.String("project-id", proj.Id),
		zap.String("project-role", user.Role.String()),
	)

	http.Redirect(w, r, rootPath, http.StatusFound)
}

func checkState(r *http.Request, key string) error {
	state := r.FormValue(stateFormKey)
	rawStateToken, err := hex.DecodeString(state)
	if err != nil {
		return err
	}

	stateToken := string(rawStateToken)
	if !xsrftoken.Valid(stateToken, key, "", "") {
		return fmt.Errorf("invalid state")
	}

	c, err := r.Cookie(stateCookieKey)
	if err != nil {
		return err
	}

	secretState := c.Value
	if state == "" || subtle.ConstantTimeCompare([]byte(state), []byte(secretState)) != 1 {
		return fmt.Errorf("wrong state")
	}

	return nil
}
