// Copyright 2023 The PipeCD Authors.
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

package handler

import (
	"context"
	"embed"
	"fmt"
	"html"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/config"
	"github.com/pipe-cd/pipecd/pkg/datastore"
	"github.com/pipe-cd/pipecd/pkg/model"
)

//go:embed templates/*
var templateFS embed.FS

var (
	topPageTmpl           = template.Must(template.ParseFS(templateFS, "templates/Top"))
	listProjectsTmpl      = template.Must(template.ParseFS(templateFS, "templates/ListProjects"))
	applicationCountsTmpl = template.Must(template.ParseFS(templateFS, "templates/ApplicationCounts"))
	addProjectTmpl        = template.Must(template.ParseFS(templateFS, "templates/AddProject"))
	addedProjectTmpl      = template.Must(template.ParseFS(templateFS, "templates/AddedProject"))
)

type projectStore interface {
	Add(ctx context.Context, proj *model.Project) error
	List(ctx context.Context, opts datastore.ListOptions) ([]model.Project, error)
}

type Handler struct {
	port             int
	projectStore     projectStore
	sharedSSOConfigs []config.SharedSSOConfig
	server           *http.Server
	gracePeriod      time.Duration
	logger           *zap.Logger
}

func NewHandler(port int, ps projectStore, sharedSSOConfigs []config.SharedSSOConfig, gracePeriod time.Duration, logger *zap.Logger) *Handler {
	mux := http.NewServeMux()
	h := &Handler{
		projectStore:     ps,
		sharedSSOConfigs: sharedSSOConfigs,
		server: &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: mux,
		},
		gracePeriod: gracePeriod,
		logger:      logger.Named("handler"),
	}

	mux.HandleFunc("/", h.handleTop)
	mux.HandleFunc("/projects", h.handleListProjects)
	mux.HandleFunc("/projects/add", h.handleAddProject)

	return h
}

func (h *Handler) Run(ctx context.Context) error {
	doneCh := make(chan error, 1)
	ctx, cancel := context.WithCancel(ctx)

	go func() {
		defer cancel()
		doneCh <- h.run()
	}()

	<-ctx.Done()
	h.stop()
	return <-doneCh
}

func (h *Handler) run() error {
	h.logger.Info(fmt.Sprintf("handler server is running on %d", h.port))
	if err := h.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		h.logger.Error("failed to listen and serve handler server", zap.Error(err))
		return err
	}
	return nil
}

func (h *Handler) stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), h.gracePeriod)
	defer cancel()
	h.logger.Info("stopping handler server")
	if err := h.server.Shutdown(ctx); err != nil {
		h.logger.Error("failed to shutdown handler server", zap.Error(err))
		return err
	}
	return nil
}

func (h *Handler) handleTop(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	if err := topPageTmpl.Execute(w, nil); err != nil {
		h.logger.Error("failed to render Top page template", zap.Error(err))
	}
}

func (h *Handler) handleListProjects(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	projects, err := h.projectStore.List(ctx, datastore.ListOptions{})
	if err != nil {
		h.logger.Error("failed to retrieve the list of projects", zap.Error(err))
		http.Error(w, "Unable to retrieve projects", http.StatusInternalServerError)
		return
	}

	data := make([]map[string]string, 0, len(projects))
	for i := range projects {
		data = append(data, map[string]string{
			"ID":                  projects[i].Id,
			"Description":         projects[i].Desc,
			"StaticAdminDisabled": strconv.FormatBool(projects[i].StaticAdminDisabled),
			"SharedSSOName":       projects[i].SharedSsoName,
			"CreatedAt":           time.Unix(projects[i].CreatedAt, 0).String(),
		})
	}
	if err := listProjectsTmpl.Execute(w, data); err != nil {
		h.logger.Error("failed to render ListProjects page template", zap.Error(err))
	}
}

func (h *Handler) handleAddProject(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	if r.Method == http.MethodGet {
		if err := addProjectTmpl.Execute(w, nil); err != nil {
			h.logger.Error("failed to render AddProject page template", zap.Error(err))
		}
		return
	}

	var (
		id                 = html.EscapeString(r.FormValue("ID"))
		description        = html.EscapeString(r.FormValue("Description"))
		sharedSSOName      = html.EscapeString(r.FormValue("SharedSSO"))
		allowStrayAsViewer = r.FormValue("AllowStrayAsViewer") == "on"
	)
	if id == "" {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	if sharedSSOName != "" {
		found := false
		for i := range h.sharedSSOConfigs {
			if h.sharedSSOConfigs[i].Name == sharedSSOName {
				found = true
				break
			}
		}
		if !found {
			http.Error(w, fmt.Sprintf("SharedSSOConfig %q was not found in Control Plane configuration", sharedSSOName), http.StatusBadRequest)
			return
		}
	}

	var (
		project = &model.Project{
			Id:                 id,
			Desc:               description,
			SharedSsoName:      sharedSSOName,
			AllowStrayAsViewer: allowStrayAsViewer,
		}
		username = model.GenerateRandomString(10)
		password = model.GenerateRandomString(30)
	)

	if err := project.SetStaticAdmin(username, password); err != nil {
		h.logger.Error("failed to set static admin",
			zap.String("id", id),
			zap.Error(err),
		)
		http.Error(w, fmt.Sprintf("Unable to add the project (%v)", err), http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := h.projectStore.Add(ctx, project); err != nil {
		h.logger.Error("failed to add a new project",
			zap.String("id", id),
			zap.Error(err),
		)
		http.Error(w, fmt.Sprintf("Unable to add the project (%v)", err), http.StatusInternalServerError)
		return
	}
	h.logger.Info("successfully added a new project", zap.String("id", id))

	data := map[string]string{
		"ID":                  id,
		"Description":         description,
		"SharedSSOName":       sharedSSOName,
		"StaticAdminUsername": username,
		"StaticAdminPassword": password,
	}
	if err := addedProjectTmpl.Execute(w, data); err != nil {
		h.logger.Error("failed to render AddedProject page template", zap.Error(err))
	}
}

func groupApplicationCounts(counts []model.InsightApplicationCount) (total int, groups map[string]int) {
	groups = make(map[string]int)
	for _, c := range counts {
		total += int(c.Count)
		kind := c.Labels[model.InsightApplicationCountLabelKey_KIND.String()]
		groups[kind] = groups[kind] + int(c.Count)
	}
	return
}
