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
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipe/pkg/datastore"
	"github.com/pipe-cd/pipe/pkg/jwt"
)

const (
	// loginPath is the path to login to pipecd projects.
	loginPath = "/auth/login"
	// passwordLoginPath is the path to login to pipecd projects with password.
	passwordLoginPath = loginPath + "/password"
	// callbackPath is the path configured in the GitHub oauth application settings.
	callbackPath = "/auth/callback"
	// logoutPath is the path for logging out from current session.
	logoutPath = "/auth/logout"
)

const (
	stateCookieKey = "state"
	errorCookieKey = "error"

	defaultTokenTTL          = 7 * 24 * time.Hour
	defaultStateCookieMaxAge = 30 * 60
	defaultErrorCookieMaxAge = 10 * 60
	defaultTokenCookieMaxAge = 7 * 24 * 60 * 60
)

// Handler handles all imcoming requests about authentication.
type Handler struct {
	signer       jwt.Signer
	projectStore datastore.ProjectStore
	logger       *zap.Logger
}

// NewHandler returns a handler that will used for authentication.
func NewHandler(
	signer jwt.Signer,
	projectStore datastore.ProjectStore,
	logger *zap.Logger,
) *Handler {
	return &Handler{
		signer:       signer,
		projectStore: projectStore,
		logger:       logger,
	}
}

// Register registers all handler into the specified registry.
func (h *Handler) Register(reg func(string, func(http.ResponseWriter, *http.Request))) {
	reg(loginPath, h.handleLogin)
	reg(logoutPath, h.handleLogout)
}

// handleLogout cleans current cookies and redirects to login page.
func (h *Handler) handleLogout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	http.SetCookie(w, makeExpiredTokenCookie())
	http.SetCookie(w, makeExpiredStateCookie())

	http.Redirect(w, r, "/", http.StatusFound)
}

func makeTokenCookie(value string) *http.Cookie {
	return &http.Cookie{
		Name:     jwt.SignedTokenKey,
		Value:    value,
		MaxAge:   defaultTokenCookieMaxAge,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
}

func makeExpiredTokenCookie() *http.Cookie {
	return &http.Cookie{
		Name:     jwt.SignedTokenKey,
		Value:    "",
		MaxAge:   -1,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
}

func makeStateCookie(value string) *http.Cookie {
	return &http.Cookie{
		Name:     stateCookieKey,
		Value:    value,
		MaxAge:   defaultStateCookieMaxAge,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
}

func makeExpiredStateCookie() *http.Cookie {
	return &http.Cookie{
		Name:     stateCookieKey,
		Value:    "",
		MaxAge:   -1,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
}

func makeErrorCookie(value string) *http.Cookie {
	return &http.Cookie{
		Name:     errorCookieKey,
		Value:    value,
		MaxAge:   defaultErrorCookieMaxAge,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
}

func serverError(w http.ResponseWriter, r *http.Request, url, msg string, logger *zap.Logger, err error) {
	logger.Error(fmt.Sprintf("auth-handler: %s", msg), zap.Error(err))
	http.SetCookie(w, makeErrorCookie(msg))
	http.Redirect(w, r, url, http.StatusSeeOther)
}
