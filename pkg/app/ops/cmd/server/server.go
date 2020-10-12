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

package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/pipe-cd/pipe/pkg/admin"
	"github.com/pipe-cd/pipe/pkg/app/ops/handler"
	"github.com/pipe-cd/pipe/pkg/cli"
	"github.com/pipe-cd/pipe/pkg/config"
	"github.com/pipe-cd/pipe/pkg/datastore"
	"github.com/pipe-cd/pipe/pkg/datastore/firestore"
	"github.com/pipe-cd/pipe/pkg/datastore/mongodb"
	"github.com/pipe-cd/pipe/pkg/model"
	"github.com/pipe-cd/pipe/pkg/version"
)

type server struct {
	httpPort    int
	adminPort   int
	gracePeriod time.Duration
	configFile  string
}

func NewCommand() *cobra.Command {
	s := &server{
		httpPort:    9082,
		adminPort:   9085,
		gracePeriod: 15 * time.Second,
	}
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Start running ops server.",
		RunE:  cli.WithContext(s.run),
	}
	cmd.Flags().IntVar(&s.httpPort, "http-port", s.httpPort, "The port number used to run http server.")
	cmd.Flags().IntVar(&s.adminPort, "admin-port", s.adminPort, "The port number used to run a HTTP server for admin tasks such as metrics, healthz.")
	cmd.Flags().DurationVar(&s.gracePeriod, "grace-period", s.gracePeriod, "How long to wait for graceful shutdown.")

	cmd.Flags().StringVar(&s.configFile, "config-file", s.configFile, "The path to the configuration file.")
	return cmd
}

func (s *server) run(ctx context.Context, t cli.Telemetry) error {
	group, ctx := errgroup.WithContext(ctx)

	// Load control plane configuration from the specified file.
	cfg, err := s.loadConfig()
	if err != nil {
		t.Logger.Error("failed to load control-plane configuration",
			zap.String("config-file", s.configFile),
			zap.Error(err),
		)
		return err
	}

	// Connect to the data store.
	ds, err := s.createDatastore(ctx, cfg, t.Logger)
	if err != nil {
		t.Logger.Error("failed to create datastore", zap.Error(err))
		return err
	}
	defer func() {
		if err := ds.Close(); err != nil {
			t.Logger.Error("failed to close datastore client", zap.Error(err))

		}
	}()

	// Start running HTTP server.
	{
		handler := handler.NewHandler(s.httpPort, datastore.NewProjectStore(ds), cfg.SharedSSOConfigs, s.gracePeriod, t.Logger)
		group.Go(func() error {
			return handler.Run(ctx)
		})
	}

	// Start running admin server.
	{
		var (
			ver   = []byte(version.Get().Version)
			admin = admin.NewAdmin(s.adminPort, s.gracePeriod, t.Logger)
		)

		admin.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
			w.Write(ver)
		})
		admin.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("ok"))
		})
		admin.Handle("/metrics", t.PrometheusMetricsHandler())

		group.Go(func() error {
			return admin.Run(ctx)
		})
	}

	// Wait until all components have finished.
	// A terminating signal or a finish of any components
	// could trigger the finish of server.
	// This ensures that all components are good or no one.
	if err := group.Wait(); err != nil {
		t.Logger.Error("failed while running", zap.Error(err))
		return err
	}
	return nil
}

func (s *server) loadConfig() (*config.ControlPlaneSpec, error) {
	cfg, err := config.LoadFromYAML(s.configFile)
	if err != nil {
		return nil, err
	}
	if cfg.Kind != config.KindControlPlane {
		return nil, fmt.Errorf("wrong configuration kind for control-plane: %v", cfg.Kind)
	}
	return cfg.ControlPlaneSpec, nil
}

func (s *server) createDatastore(ctx context.Context, cfg *config.ControlPlaneSpec, logger *zap.Logger) (datastore.DataStore, error) {
	switch cfg.Datastore.Type {
	case model.DataStoreFirestore:
		fsConfig := cfg.Datastore.FirestoreConfig
		options := []firestore.Option{
			firestore.WithCredentialsFile(fsConfig.CredentialsFile),
			firestore.WithLogger(logger),
		}
		return firestore.NewFireStore(ctx, fsConfig.Project, fsConfig.Namespace, fsConfig.Environment, options...)

	case model.DataStoreDynamoDB:
		return nil, errors.New("dynamodb is unimplemented yet")

	case model.DataStoreMongoDB:
		mdConfig := cfg.Datastore.MongoDBConfig
		options := []mongodb.Option{
			mongodb.WithLogger(logger),
		}
		return mongodb.NewMongoDB(ctx, mdConfig.URL, mdConfig.Database, options...)

	default:
		return nil, fmt.Errorf("unknown datastore type %q", cfg.Datastore.Type)
	}
}
