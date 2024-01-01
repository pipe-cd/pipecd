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

package admin

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// Admin is a http server for exposing private information, e.g.
// - application metrics
// - prom metrics: go, process
// - service health check
// - runtime configuration
type Admin struct {
	port        int
	mux         *http.ServeMux
	server      *http.Server
	patterns    []string
	gracePeriod time.Duration
	logger      *zap.Logger
}

func NewAdmin(port int, gracePeriod time.Duration, logger *zap.Logger) *Admin {
	mux := http.NewServeMux()
	a := &Admin{
		port: port,
		mux:  mux,
		server: &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: mux,
		},
		gracePeriod: gracePeriod,
		logger:      logger.Named("admin"),
	}
	mux.HandleFunc("/", a.handleTop)
	return a
}

func (a *Admin) Handle(pattern string, handler http.Handler) {
	a.patterns = append(a.patterns, pattern)
	a.mux.Handle(pattern, handler)
}

func (a *Admin) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	a.patterns = append(a.patterns, pattern)
	a.mux.HandleFunc(pattern, handler)
}

func (a *Admin) handleTop(w http.ResponseWriter, r *http.Request) {
	buf := new(bytes.Buffer)
	if err := topPageTmpl.Execute(buf, a.patterns); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(buf.Bytes())
}

func (a *Admin) Run(ctx context.Context) error {
	doneCh := make(chan error, 1)
	ctx, cancel := context.WithCancel(ctx)

	go func() {
		defer cancel()
		doneCh <- a.run()
	}()

	<-ctx.Done()
	a.stop()
	return <-doneCh
}

func (a *Admin) run() error {
	a.logger.Info(fmt.Sprintf("admin server is running on %d", a.port))
	if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		a.logger.Error("failed to listen and serve admin server", zap.Error(err))
		return err
	}
	return nil
}

func (a *Admin) stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), a.gracePeriod)
	defer cancel()
	a.logger.Info("stopping admin server")
	if err := a.server.Shutdown(ctx); err != nil {
		a.logger.Error("failed to shutdown admin server", zap.Error(err))
		return err
	}
	return nil
}

const topPageTemplate = `
<!DOCTYPE html>
<html>
<body>

<h3>Admin Page</h3>
{{- range . }}
<p><a href="{{ . }}">{{ . }}</a></p>
{{- end }}
</body>
</html>
`

var topPageTmpl = template.Must(template.New("toppage").Parse(topPageTemplate))
