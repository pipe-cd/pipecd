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

package httpapi

import (
	"net/http"
	"path/filepath"

	"github.com/NYTimes/gziphandler"
	"go.uber.org/zap"

	"github.com/pipe-cd/pipe/pkg/app/api/httpapi/metricsmiddleware"
	"github.com/pipe-cd/pipe/pkg/config"
	"github.com/pipe-cd/pipe/pkg/jwt"
	"github.com/pipe-cd/pipe/pkg/model"
)

// NewHandler gives back an HTTP handler for serving PipeCD SPA.
func NewHandler(
	signer jwt.Signer,
	staticDir string,
	decrypter decrypter,
	address string,
	stateKey string,
	projectsInConfig map[string]config.ControlPlaneProject,
	sharedSSOConfigs map[string]*model.ProjectSSOConfig,
	projectGetter projectGetter,
	secureCookie bool,
	logger *zap.Logger,
) http.Handler {
	mux := http.NewServeMux()
	a := newAuthHandler(
		signer,
		decrypter,
		address,
		stateKey,
		projectsInConfig,
		sharedSSOConfigs,
		projectGetter,
		secureCookie,
		logger,
	)

	fs := http.FileServer(http.Dir(filepath.Join(staticDir, "assets")))
	assetsHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=31536000")
		http.StripPrefix("/assets/", fs).ServeHTTP(w, r)
	})
	mux.Handle("/assets/", gziphandler.GzipHandler(assetsHandler))

	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(staticDir, "favicon.ico"))
	})

	rootHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(staticDir, "/index.html"))
	})
	mux.Handle(rootPath, metricsmiddleware.InstrumentHandler(rootPath, rootHandler))

	mux.Handle(loginPath, metricsmiddleware.InstrumentHandler(loginPath, http.HandlerFunc(a.handleSSOLogin)))
	mux.Handle(staticLoginPath, metricsmiddleware.InstrumentHandler(staticLoginPath, http.HandlerFunc(a.handleStaticAdminLogin)))
	mux.Handle(callbackPath, metricsmiddleware.InstrumentHandler(callbackPath, http.HandlerFunc(a.handleCallback)))
	mux.Handle(logoutPath, metricsmiddleware.InstrumentHandler(logoutPath, http.HandlerFunc(a.handleLogout)))

	return mux
}
