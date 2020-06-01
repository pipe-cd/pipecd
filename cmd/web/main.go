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

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/NYTimes/gziphandler"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/kapetaniosci/pipe/pkg/admin"
	"github.com/kapetaniosci/pipe/pkg/cli"
)

func main() {
	app := cli.NewApp(
		"web",
		"A service for serving static assets.",
	)
	app.AddCommands(newCommand())
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

type server struct {
	httpPort    int
	adminPort   int
	staticDir   string
	gracePeriod time.Duration
}

func newCommand() *cobra.Command {
	s := &server{
		httpPort:    9082,
		adminPort:   9085,
		staticDir:   "pkg/app/web/public_files",
		gracePeriod: 15 * time.Second,
	}
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Start running http server.",
		RunE:  cli.WithContext(s.run),
	}

	cmd.Flags().IntVar(&s.httpPort, "http-port", s.httpPort, "The port number used to run http server.")
	cmd.Flags().IntVar(&s.adminPort, "admin-port", s.adminPort, "The port number used to run a HTTP server for admin tasks such as metrics, healthz.")
	cmd.Flags().StringVar(&s.staticDir, "static-dir", s.staticDir, "The directory where contains static assets.")
	cmd.Flags().DurationVar(&s.gracePeriod, "grace-period", s.gracePeriod, "How long to wait for graceful shutdown.")

	return cmd
}

func (s *server) run(ctx context.Context, t cli.Telemetry) error {
	group, ctx := errgroup.WithContext(ctx)

	// Start running http server.
	{
		var (
			mux    = http.NewServeMux()
			server = &http.Server{
				Addr:    fmt.Sprintf(":%d", s.httpPort),
				Handler: mux,
			}
			fs            = http.FileServer(http.Dir(filepath.Join(s.staticDir, "assets")))
			assetsHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Cache-Control", "public, max-age=31536000")
				http.StripPrefix("/assets/", fs).ServeHTTP(w, r)
			})
		)

		mux.Handle("/assets/", gziphandler.GzipHandler(assetsHandler))
		mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, filepath.Join(s.staticDir, "favicon.ico"))
		})
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, filepath.Join(s.staticDir, "/index.html"))
		})

		group.Go(func() error {
			defer func() {
				ctx, cancel := context.WithTimeout(context.Background(), s.gracePeriod)
				defer cancel()
				t.Logger.Info("stopping http server")
				if err := server.Shutdown(ctx); err != nil {
					t.Logger.Error("failed to shutdown http server", zap.Error(err))
					return
				}
				t.Logger.Info("http server is stopped")
			}()

			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				t.Logger.Error("failed to listen and serve http server", zap.Error(err))
				return err
			}
			return nil
		})
	}

	// Start running admin server.
	{
		admin := admin.NewAdmin(s.adminPort, s.gracePeriod, t.Logger)
		if exporter, ok := t.PrometheusMetricsExporter(); ok {
			admin.Handle("/metrics", exporter)
		}
		admin.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("ok"))
		})
		group.Go(func() error {
			return admin.Run(ctx)
		})
	}

	// Wait for completions of all components.
	// A terminating sinal or a finish of any component could trigger the finish of application.
	if err := group.Wait(); err != nil {
		t.Logger.Error("failed while running", zap.Error(err))
		return err
	}

	return nil
}
