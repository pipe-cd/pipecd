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
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/pipe-cd/pipe/pkg/admin"
	"github.com/pipe-cd/pipe/pkg/app/ops/firestoreindexensurer"
	"github.com/pipe-cd/pipe/pkg/app/ops/handler"
	"github.com/pipe-cd/pipe/pkg/app/ops/insightcollector"
	"github.com/pipe-cd/pipe/pkg/app/ops/mysqlensurer"
	"github.com/pipe-cd/pipe/pkg/app/ops/orphancommandcleaner"
	"github.com/pipe-cd/pipe/pkg/app/ops/pipedstatsbuilder"
	"github.com/pipe-cd/pipe/pkg/cache/rediscache"
	"github.com/pipe-cd/pipe/pkg/cli"
	"github.com/pipe-cd/pipe/pkg/config"
	"github.com/pipe-cd/pipe/pkg/datastore"
	"github.com/pipe-cd/pipe/pkg/insight/insightmetrics"
	"github.com/pipe-cd/pipe/pkg/insight/insightstore"
	"github.com/pipe-cd/pipe/pkg/model"
	"github.com/pipe-cd/pipe/pkg/redis"
	"github.com/pipe-cd/pipe/pkg/version"
)

type ops struct {
	httpPort               int
	adminPort              int
	gracePeriod            time.Duration
	enableInsightCollector bool
	configFile             string
	gcloudPath             string
	cacheAddress           string
}

func NewOpsCommand() *cobra.Command {
	s := &ops{
		httpPort:     9082,
		adminPort:    9085,
		cacheAddress: "cache:6379",
		gracePeriod:  15 * time.Second,
	}
	cmd := &cobra.Command{
		Use:   "ops",
		Short: "Start running ops server.",
		RunE:  cli.WithContext(s.run),
	}
	cmd.Flags().IntVar(&s.httpPort, "http-port", s.httpPort, "The port number used to run http server.")
	cmd.Flags().IntVar(&s.adminPort, "admin-port", s.adminPort, "The port number used to run a HTTP server for admin tasks such as metrics, healthz.")
	cmd.Flags().DurationVar(&s.gracePeriod, "grace-period", s.gracePeriod, "How long to wait for graceful shutdown.")
	cmd.Flags().StringVar(&s.configFile, "config-file", s.configFile, "The path to the configuration file.")
	cmd.Flags().StringVar(&s.gcloudPath, "gcloud-path", s.gcloudPath, "The path to the gcloud command executable.")
	cmd.Flags().StringVar(&s.cacheAddress, "cache-address", s.cacheAddress, "The address to cache service.")
	return cmd
}

func (s *ops) run(ctx context.Context, t cli.Telemetry) error {
	group, ctx := errgroup.WithContext(ctx)

	// Load control plane configuration from the specified file.
	cfg, err := loadConfig(s.configFile)
	if err != nil {
		t.Logger.Error("failed to load control-plane configuration",
			zap.String("config-file", s.configFile),
			zap.Error(err),
		)
		return err
	}

	// Prepare sql database.
	if cfg.Datastore.Type == model.DataStoreMySQL {
		if err := ensureSQLDatabase(ctx, cfg, t.Logger); err != nil {
			t.Logger.Error("failed to ensure prepare SQL database", zap.Error(err))
			return err
		}
	}

	if cfg.Datastore.Type == model.DataStoreFirestore {
		// Create needed composite indexes for Firestore.
		ensurer := firestoreindexensurer.NewIndexEnsurer(
			s.gcloudPath,
			cfg.Datastore.FirestoreConfig.Project,
			cfg.Datastore.FirestoreConfig.CredentialsFile,
			cfg.Datastore.FirestoreConfig.CollectionNamePrefix,
			t.Logger,
		)
		group.Go(func() error {
			return ensurer.CreateIndexes(ctx)
		})
	}

	// Connect to the data store.
	ds, err := createDatastore(ctx, cfg, t.Logger)
	if err != nil {
		t.Logger.Error("failed to create datastore", zap.Error(err))
		return err
	}
	defer func() {
		if err := ds.Close(); err != nil {
			t.Logger.Error("failed to close datastore client", zap.Error(err))
		}
	}()

	// Connect to the file store.
	fs, err := createFilestore(ctx, cfg, t.Logger)
	if err != nil {
		t.Logger.Error("failed to create filestore", zap.Error(err))
		return err
	}
	defer func() {
		if err := fs.Close(); err != nil {
			t.Logger.Error("failed to close filestore client", zap.Error(err))
		}
	}()

	// Start running command cleaner.
	cleaner := orphancommandcleaner.NewOrphanCommandCleaner(ds, t.Logger)
	group.Go(func() error {
		return cleaner.Run(ctx)
	})

	// Start running insight collector.
	ic := insightcollector.NewCollector(ds, fs, cfg.InsightCollector, t.Logger)
	group.Go(func() error {
		return ic.Run(ctx)
	})
	insightMetricsCollector := insightmetrics.NewInsightMetricsCollector(insightstore.NewStore(fs), datastore.NewProjectStore(ds))

	// Start running HTTP server.
	{
		handler := handler.NewHandler(s.httpPort, datastore.NewProjectStore(ds), insightstore.NewStore(fs), cfg.SharedSSOConfigs, s.gracePeriod, t.Logger)
		group.Go(func() error {
			return handler.Run(ctx)
		})
	}

	rd := redis.NewRedis(s.cacheAddress, "")
	statCache := rediscache.NewHashCache(rd, defaultPipedStatHashKey)
	psb := pipedstatsbuilder.NewPipedStatsBuilder(statCache, t.Logger)

	// Register all pipecd ops metrics collectors.
	reg := registerOpsMetrics(insightMetricsCollector)
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
		admin.Handle("/metrics", t.CustomMetricsHandlerFor(reg, psb))

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

func ensureSQLDatabase(ctx context.Context, cfg *config.ControlPlaneSpec, logger *zap.Logger) error {
	mysqlEnsurer, err := mysqlensurer.NewMySQLEnsurer(
		cfg.Datastore.MySQLConfig.URL,
		cfg.Datastore.MySQLConfig.Database,
		cfg.Datastore.MySQLConfig.UsernameFile,
		cfg.Datastore.MySQLConfig.PasswordFile,
		logger,
	)
	if err != nil {
		logger.Error("failed to create SQL ensurer instance", zap.Error(err))
		return err
	}
	defer func() {
		// Close connection held by the client.
		if err := mysqlEnsurer.Close(); err != nil {
			logger.Error("failed to close database ensurer connection", zap.Error(err))
		}
	}()

	if err = mysqlEnsurer.Run(ctx); err != nil {
		logger.Error("failed to ensure SQL schema and indexes", zap.Error(err))
		return err
	}

	logger.Info("prepare SQL schema and indexes successfully")
	return nil
}

func registerOpsMetrics(col ...prometheus.Collector) *prometheus.Registry {
	r := prometheus.NewRegistry()
	wrapped := prometheus.WrapRegistererWithPrefix("pipecd_ops_", r)

	wrapped.Register(prometheus.NewGoCollector())
	wrapped.Register(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))
	wrapped.MustRegister(col...)

	return r
}
