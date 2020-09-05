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
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipe/pkg/config"
	"github.com/pipe-cd/pipe/pkg/jwt"
	"github.com/pipe-cd/pipe/pkg/model"
)

const (
	rootPath = "/"
	// loginPath is the path to login to pipecd projects.
	loginPath = "/auth/login"
	// staticLoginPath is the path to login to pipecd projects with password.
	staticLoginPath = loginPath + "/static"
	// callbackPath is the path configured in the GitHub oauth application settings.
	callbackPath = "/auth/callback"
	// logoutPath is the path for logging out from current session.
	logoutPath = "/auth/logout"

	projectFormKey  = "project"
	usernameFormKey = "username"
	passwordFormKey = "password"
	authCodeFormKey = "code"
	stateFormKey    = "state"

	stateCookieKey = "state"
	errorCookieKey = "error"

	defaultTokenTTL          = 7 * 24 * time.Hour
	defaultStateCookieMaxAge = 30 * 60
	defaultErrorCookieMaxAge = 10 * 60
	defaultTokenCookieMaxAge = 7 * 24 * 60 * 60
)

type projectGetter interface {
	GetProject(ctx context.Context, id string) (*model.Project, error)
}

// Handler handles all imcoming requests about authentication.
type Handler struct {
	signer           jwt.Signer
	callbackURL      string
	stateKey         string
	projectsInConfig map[string]config.ControlPlaneProject
	projectGetter    projectGetter
	logger           *zap.Logger
}

// NewHandler returns a handler that will used for authentication.
func NewHandler(
	signer jwt.Signer,
	address string,
	stateKey string,
	projectsInConfig map[string]config.ControlPlaneProject,
	projectGetter projectGetter,
	logger *zap.Logger,
) *Handler {
	return &Handler{
		signer:           signer,
		callbackURL:      strings.TrimSuffix(address, "/") + "/" + callbackPath,
		stateKey:         stateKey,
		projectsInConfig: projectsInConfig,
		projectGetter:    projectGetter,
		logger:           logger,
	}
}

// Register registers all handler into the specified registry.
func (h *Handler) Register(r func(string, func(http.ResponseWriter, *http.Request))) {
	r(loginPath, h.handleLogin)
	r(staticLoginPath, h.handleStaticLogin)
	r(callbackPath, h.handleCallback)
	r(logoutPath, h.handleLogout)
}

// handleLogout cleans current cookies and redirects to login page.
func (h *Handler) handleLogout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	http.SetCookie(w, makeExpiredTokenCookie())
	http.SetCookie(w, makeExpiredStateCookie())

	http.Redirect(w, r, rootPath, http.StatusFound)
}

func (h *Handler) getProject(ctx context.Context, projectID string) (*model.Project, error) {
	if projectID == "" {
		return nil, fmt.Errorf("project id must be specified")
	}

	proj, err := h.projectGetter.GetProject(ctx, projectID)
	if err != nil {
		return nil, err
	}
	return proj, nil
}

func makeTokenCookie(value string, secure bool) *http.Cookie {
	return &http.Cookie{
		Name:     jwt.SignedTokenKey,
		Value:    value,
		MaxAge:   defaultTokenCookieMaxAge,
		Path:     rootPath,
		Secure:   secure,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
}

func makeExpiredTokenCookie() *http.Cookie {
	return &http.Cookie{
		Name:     jwt.SignedTokenKey,
		Value:    "",
		MaxAge:   -1,
		Path:     rootPath,
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
		Path:     rootPath,
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
		Path:     rootPath,
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
		Path:     rootPath,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
}

func handleError(w http.ResponseWriter, r *http.Request, redirectURL, responseMessage string, logger *zap.Logger, err error) {
	logger.Error(fmt.Sprintf("auth-handler: %s", responseMessage), zap.Error(err))
	http.SetCookie(w, makeErrorCookie(responseMessage))
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}
