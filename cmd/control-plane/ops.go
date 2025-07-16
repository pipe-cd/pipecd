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

package main

import (
	"context"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/pipe-cd/pipecd/pkg/admin"
	"github.com/pipe-cd/pipecd/pkg/app/ops/apikeylastusedtimeupdater"
	"github.com/pipe-cd/pipecd/pkg/app/ops/deploymentchaincontroller"
	"github.com/pipe-cd/pipecd/pkg/app/ops/firestoreindexensurer"
	"github.com/pipe-cd/pipecd/pkg/app/ops/handler"
	"github.com/pipe-cd/pipecd/pkg/app/ops/insightcollector"
	"github.com/pipe-cd/pipecd/pkg/app/ops/mysqlensurer"
	"github.com/pipe-cd/pipecd/pkg/app/ops/orphancommandcleaner"
	"github.com/pipe-cd/pipecd/pkg/app/ops/pipedstatsbuilder"
	"github.com/pipe-cd/pipecd/pkg/app/ops/planpreviewoutputcleaner"
	"github.com/pipe-cd/pipecd/pkg/app/ops/staledpipedstatcleaner"
	"github.com/pipe-cd/pipecd/pkg/cache/rediscache"
	"github.com/pipe-cd/pipecd/pkg/cli"
	"github.com/pipe-cd/pipecd/pkg/config"
	"github.com/pipe-cd/pipecd/pkg/datastore"
	"github.com/pipe-cd/pipecd/pkg/insight"
	"github.com/pipe-cd/pipecd/pkg/insight/insightmetrics"
	"github.com/pipe-cd/pipecd/pkg/insight/insightstore"
	"github.com/pipe-cd/pipecd/pkg/model"
	"github.com/pipe-cd/pipecd/pkg/redis"
	"github.com/pipe-cd/pipecd/pkg/version"
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

func (s *ops) run(ctx context.Context, input cli.Input) error {
	group, ctx := errgroup.WithContext(ctx)

	// Load control plane configuration from the specified file.
	cfg, err := loadConfig(s.configFile)
	if err != nil {
		input.Logger.Error("failed to load control-plane configuration",
			zap.String("config-file", s.configFile),
			zap.Error(err),
		)
		return err
	}

	// Connect to the cache.
	rd := redis.NewRedis(s.cacheAddress, "")
	defer func() {
		if err := rd.Close(); err != nil {
			input.Logger.Error("failed to close redis client", zap.Error(err))
		}
	}()

	// Connect to the file store.
	fs, err := createFilestore(ctx, cfg, input.Logger)
	if err != nil {
		input.Logger.Error("failed to create filestore", zap.Error(err))
		return err
	}
	defer func() {
		if err := fs.Close(); err != nil {
			input.Logger.Error("failed to close filestore client", zap.Error(err))
		}
	}()

	// Prepare sql database.
	if cfg.Datastore.Type == model.DataStoreMySQL {
		if err := ensureSQLDatabase(ctx, cfg, input.Logger); err != nil {
			input.Logger.Error("failed to ensure prepare SQL database", zap.Error(err))
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
			input.Logger,
		)
		group.Go(func() error {
			return ensurer.CreateIndexes(ctx)
		})
	}

	dbCache := rediscache.NewTTLCache(rd, 3*time.Hour)
	// Connect to the data store.
	ds, err := createDatastore(ctx, cfg, fs, dbCache, input.Logger)
	if err != nil {
		input.Logger.Error("failed to create datastore", zap.Error(err))
		return err
	}
	defer func() {
		if err := ds.Close(); err != nil {
			input.Logger.Error("failed to close datastore client", zap.Error(err))
		}
	}()

	statCache := rediscache.NewHashCache(rd, defaultPipedStatHashKey)
	// Start running staled piped stat cleaner.
	{
		cleaner := staledpipedstatcleaner.NewStaledPipedStatCleaner(statCache, input.Logger)
		group.Go(func() error {
			return cleaner.Run(ctx)
		})
	}

	// Start running command cleaner.
	{
		cleaner := orphancommandcleaner.NewOrphanCommandCleaner(ds, input.Logger)
		group.Go(func() error {
			return cleaner.Run(ctx)
		})
	}

	// Start running planpreview output cleaner.
	{
		cleaner := planpreviewoutputcleaner.NewCleaner(fs, input.Logger)
		group.Go(func() error {
			return cleaner.Run(ctx)
		})
	}

	// Start runnning apiKeyLastUsedTime updater.
	{
		updater := apikeylastusedtimeupdater.NewAPIKeyLastUsedTimeUpdater(ds, rd, input.Logger)
		group.Go(func() error {
			return updater.Run(ctx)
		})
	}

	// Start deployment chain controller.
	{
		controller := deploymentchaincontroller.NewDeploymentChainController(ds, input.Logger)
		group.Go(func() error {
			return controller.Run(ctx)
		})
	}

	insightStore := insightstore.NewStore(
		fs,
		cfg.InsightCollector.Deployment.ChunkMaxCount,
		rd,
		input.Logger,
	)
	// Start running insight collector.
	{
		ic := insightcollector.NewCollector(ds, insightStore, cfg.InsightCollector, input.Logger)
		group.Go(func() error {
			return ic.Run(ctx)
		})
	}

	insightMetricsCollector := insightmetrics.NewInsightMetricsCollector(
		insight.NewProvider(insightStore),
		datastore.NewProjectStore(ds, datastore.OpsCommander),
	)

	// Start running HTTP server.
	{
		handler := handler.NewHandler(
			s.httpPort,
			datastore.NewProjectStore(ds, datastore.OpsCommander),
			cfg.SharedSSOConfigs,
			s.gracePeriod,
			input.Logger,
		)
		group.Go(func() error {
			return handler.Run(ctx)
		})
	}

	psb := pipedstatsbuilder.NewPipedStatsBuilder(statCache, input.Logger)

	// Register all pipecd ops metrics collectors.
	reg := registerOpsMetrics(insightMetricsCollector)
	// Start running admin server.
	{
		var (
			ver   = []byte(version.Get().Version)
			admin = admin.NewAdmin(s.adminPort, s.gracePeriod, input.Logger)
		)

		admin.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
			w.Write(ver)
		})
		admin.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("ok"))
		})
		admin.Handle("/metrics", input.CustomMetricsHandlerFor(reg, psb))

		group.Go(func() error {
			return admin.Run(ctx)
		})
	}

	// Wait until all components have finished.
	// A terminating signal or a finish of any components
	// could trigger the finish of server.
	// This ensures that all components are good or no one.
	if err := group.Wait(); err != nil {
		input.Logger.Error("failed while running", zap.Error(err))
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

	logger.Info("start running SQL schema and indexes ensurer")
	if err = mysqlEnsurer.Run(ctx); err != nil {
		logger.Error("failed to ensure SQL schema and indexes", zap.Error(err))
		return err
	}

	logger.Info("prepare SQL schema and indexes successfully")
	return nil
}

func registerOpsMetrics(col ...prometheus.Collector) *prometheus.Registry {
	r := prometheus.NewRegistry()
	wrapped := prometheus.WrapRegistererWith(map[string]string{
		"pipecd_component": "ops",
	}, r)

	wrapped.Register(collectors.NewGoCollector())
	wrapped.Register(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	wrapped.MustRegister(col...)

	return r
}
