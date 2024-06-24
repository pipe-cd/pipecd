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
	"fmt"
	"net/http"
	"time"

	jwtgo "github.com/golang-jwt/jwt"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/pipe-cd/pipecd/pkg/admin"
	"github.com/pipe-cd/pipecd/pkg/app/server/analysisresultstore"
	"github.com/pipe-cd/pipecd/pkg/app/server/apikeyverifier"
	"github.com/pipe-cd/pipecd/pkg/app/server/applicationlivestatestore"
	"github.com/pipe-cd/pipecd/pkg/app/server/commandoutputstore"
	"github.com/pipe-cd/pipecd/pkg/app/server/grpcapi"
	"github.com/pipe-cd/pipecd/pkg/app/server/grpcapi/grpcapimetrics"
	"github.com/pipe-cd/pipecd/pkg/app/server/httpapi"
	"github.com/pipe-cd/pipecd/pkg/app/server/httpapi/httpapimetrics"
	"github.com/pipe-cd/pipecd/pkg/app/server/pipedverifier"
	"github.com/pipe-cd/pipecd/pkg/app/server/service/webservice"
	"github.com/pipe-cd/pipecd/pkg/app/server/stagelogstore"
	"github.com/pipe-cd/pipecd/pkg/app/server/unregisteredappstore"
	"github.com/pipe-cd/pipecd/pkg/cache"
	"github.com/pipe-cd/pipecd/pkg/cache/cachemetrics"
	"github.com/pipe-cd/pipecd/pkg/cache/rediscache"
	"github.com/pipe-cd/pipecd/pkg/cli"
	"github.com/pipe-cd/pipecd/pkg/config"
	"github.com/pipe-cd/pipecd/pkg/crypto"
	"github.com/pipe-cd/pipecd/pkg/datastore"
	"github.com/pipe-cd/pipecd/pkg/datastore/filedb"
	"github.com/pipe-cd/pipecd/pkg/datastore/firestore"
	"github.com/pipe-cd/pipecd/pkg/datastore/mysql"
	"github.com/pipe-cd/pipecd/pkg/filestore"
	"github.com/pipe-cd/pipecd/pkg/filestore/gcs"
	"github.com/pipe-cd/pipecd/pkg/filestore/minio"
	"github.com/pipe-cd/pipecd/pkg/filestore/s3"
	"github.com/pipe-cd/pipecd/pkg/insight"
	"github.com/pipe-cd/pipecd/pkg/insight/insightstore"
	"github.com/pipe-cd/pipecd/pkg/jwt"
	"github.com/pipe-cd/pipecd/pkg/model"
	"github.com/pipe-cd/pipecd/pkg/redis"
	"github.com/pipe-cd/pipecd/pkg/rpc"
	"github.com/pipe-cd/pipecd/pkg/version"
)

var (
	defaultSigningMethod = jwtgo.SigningMethodHS256
)

const (
	defaultPipedStatHashKey    = "HASHKEY:PIPED:STATS"
	apiKeyLastUsedCacheHashKey = "HASHKEY:PIPED:API_KEYS" //nolint:gosec
)

type server struct {
	pipedAPIPort   int
	webAPIPort     int
	httpPort       int
	apiPort        int
	adminPort      int
	envoyAuthzPort int
	staticDir      string
	cacheAddress   string
	gracePeriod    time.Duration

	tls            bool
	certFile       string
	keyFile        string
	insecureCookie bool

	encryptionKeyFile string
	configFile        string

	enableGRPCReflection bool
}

// NewServerCommand creates a new cobra command for executing api server.
func NewServerCommand() *cobra.Command {
	s := &server{
		pipedAPIPort:   9080,
		webAPIPort:     9081,
		httpPort:       9082,
		apiPort:        9083,
		adminPort:      9085,
		envoyAuthzPort: 9086,
		staticDir:      "web/static",
		cacheAddress:   "cache:6379",
		gracePeriod:    30 * time.Second,
	}
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Start running server.",
		RunE:  cli.WithContext(s.run),
	}

	cmd.Flags().IntVar(&s.pipedAPIPort, "piped-api-port", s.pipedAPIPort, "The port number used to run a grpc server that serving serves incoming piped requests.")
	cmd.Flags().IntVar(&s.webAPIPort, "web-api-port", s.webAPIPort, "The port number used to run a grpc server that serves incoming web requests.")
	cmd.Flags().IntVar(&s.httpPort, "http-port", s.httpPort, "The port number used to run a http server that serves incoming http requests such as auth callbacks or webhook events.")
	cmd.Flags().IntVar(&s.apiPort, "api-port", s.apiPort, "The port number used to run a grpc server for external apis.")
	cmd.Flags().IntVar(&s.adminPort, "admin-port", s.adminPort, "The port number used to run a HTTP server for admin tasks such as metrics, healthz.")
	cmd.Flags().IntVar(&s.envoyAuthzPort, "envoy-authz-port", s.envoyAuthzPort, "The port number used to run a gRPC server that serves envoy ExtAuthz service.")
	cmd.Flags().StringVar(&s.staticDir, "static-dir", s.staticDir, "The directory where contains static assets.")
	cmd.Flags().StringVar(&s.cacheAddress, "cache-address", s.cacheAddress, "The address to cache service.")
	cmd.Flags().DurationVar(&s.gracePeriod, "grace-period", s.gracePeriod, "How long to wait for graceful shutdown.")

	cmd.Flags().BoolVar(&s.tls, "tls", s.tls, "Whether running the gRPC server with TLS or not.")
	cmd.Flags().StringVar(&s.certFile, "cert-file", s.certFile, "The path to the TLS certificate file.")
	cmd.Flags().StringVar(&s.keyFile, "key-file", s.keyFile, "The path to the TLS key file.")
	cmd.Flags().BoolVar(&s.insecureCookie, "insecure-cookie", s.insecureCookie, "Allow cookie to be sent over an unsecured HTTP connection.")

	cmd.Flags().StringVar(&s.encryptionKeyFile, "encryption-key-file", s.encryptionKeyFile, "The path to file containing a random string of bits used to encrypt sensitive data.")
	cmd.MarkFlagRequired("encryption-key-file")
	cmd.Flags().StringVar(&s.configFile, "config-file", s.configFile, "The path to the configuration file.")
	cmd.MarkFlagRequired("config-file")

	// For debugging early in development
	cmd.Flags().BoolVar(&s.enableGRPCReflection, "enable-grpc-reflection", s.enableGRPCReflection, "Whether to enable the reflection service or not.")

	return cmd
}

func (s *server) run(ctx context.Context, input cli.Input) error {
	// Register all metrics.
	reg := registerMetrics()

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
	input.Logger.Info("successfully loaded control-plane configuration")

	// Connect to the cache server.
	rd := redis.NewRedis(s.cacheAddress, "")
	defer func() {
		if err := rd.Close(); err != nil {
			input.Logger.Error("failed to close redis client", zap.Error(err))
		}
	}()

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
	input.Logger.Info("successfully connected to file store")

	dbCache := rediscache.NewTTLCache(rd, 3*time.Hour)
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
	input.Logger.Info("successfully connected to data store")

	var (
		cache                = rediscache.NewTTLCache(rd, cfg.Cache.TTLDuration())
		sls                  = stagelogstore.NewStore(fs, cache, input.Logger)
		alss                 = applicationlivestatestore.NewStore(fs, cache, input.Logger)
		las                  = analysisresultstore.NewStore(fs, input.Logger)
		insightStore         = insightstore.NewStore(fs, cfg.InsightCollector.Deployment.ChunkMaxCount, rd, input.Logger)
		insightProvider      = insight.NewProvider(insightStore)
		cmdOutputStore       = commandoutputstore.NewStore(fs, input.Logger)
		statCache            = rediscache.NewHashCache(rd, defaultPipedStatHashKey)
		unregisteredAppStore = unregisteredappstore.NewStore(rd, input.Logger)
		apiKeyLastUsedCache  = rediscache.NewHashCache(rd, apiKeyLastUsedCacheHashKey)
	)

	// Start a gRPC server for handling PipedAPI requests.
	{
		var (
			verifier = pipedverifier.NewVerifier(
				ctx,
				cfg,
				// These stores are used to handle PipedAPI request, thus the writer should be PipedWriter.
				datastore.NewProjectStore(ds, datastore.PipedCommander),
				datastore.NewPipedStore(ds, datastore.PipedCommander),
				input.Logger,
			)
			service = grpcapi.NewPipedAPI(ctx, ds, cache, sls, alss, las, statCache, cmdOutputStore, unregisteredAppStore, cfg.Address, input.Logger)
			opts    = []rpc.Option{
				rpc.WithPort(s.pipedAPIPort),
				rpc.WithGracePeriod(s.gracePeriod),
				rpc.WithLogger(input.Logger),
				rpc.WithLogUnaryInterceptor(input.Logger),
				rpc.WithPipedTokenAuthUnaryInterceptor(verifier, input.Logger),
				rpc.WithRequestValidationUnaryInterceptor(),
			}
		)
		if s.tls {
			opts = append(opts, rpc.WithTLS(s.certFile, s.keyFile))
		}
		if s.enableGRPCReflection {
			opts = append(opts, rpc.WithGRPCReflection())
		}
		if input.Flags.Metrics {
			opts = append(opts, rpc.WithPrometheusUnaryInterceptor())
		}

		server := rpc.NewServer(service, opts...)
		group.Go(func() error {
			return server.Run(ctx)
		})
	}

	// Start a gRPC server for handling external API requests.
	{
		var (
			verifier = apikeyverifier.NewVerifier(
				ctx,
				datastore.NewAPIKeyStore(ds, datastore.PipectlCommander),
				apiKeyLastUsedCache,
				input.Logger,
			)

			service = grpcapi.NewAPI(ctx, ds, fs, cache, cmdOutputStore, statCache, cfg.Address, input.Logger)
			opts    = []rpc.Option{
				rpc.WithPort(s.apiPort),
				rpc.WithGracePeriod(s.gracePeriod),
				rpc.WithLogger(input.Logger),
				rpc.WithLogUnaryInterceptor(input.Logger),
				rpc.WithAPIKeyAuthUnaryInterceptor(verifier, input.Logger),
				rpc.WithRequestValidationUnaryInterceptor(),
			}
		)
		if s.tls {
			opts = append(opts, rpc.WithTLS(s.certFile, s.keyFile))
		}
		if input.Flags.Metrics {
			opts = append(opts, rpc.WithPrometheusUnaryInterceptor())
		}

		server := rpc.NewServer(service, opts...)
		group.Go(func() error {
			return server.Run(ctx)
		})
	}

	encryptDecrypter, err := crypto.NewAESEncryptDecrypter(s.encryptionKeyFile)
	if err != nil {
		input.Logger.Error("failed to create a new AES EncryptDecrypter", zap.Error(err))
		return err
	}

	// Start a gRPC server for handling WebAPI requests.
	{
		verifier, err := jwt.NewVerifier(defaultSigningMethod, s.encryptionKeyFile)
		if err != nil {
			input.Logger.Error("failed to create a new JWT verifier", zap.Error(err))
			return err
		}

		service := grpcapi.NewWebAPI(
			ctx,
			ds,
			cache,
			sls,
			alss,
			unregisteredAppStore,
			apiKeyLastUsedCache,
			insightProvider,
			statCache,
			cfg.ProjectMap(),
			encryptDecrypter,
			input.Logger,
		)
		opts := []rpc.Option{
			rpc.WithPort(s.webAPIPort),
			rpc.WithGracePeriod(s.gracePeriod),
			rpc.WithLogger(input.Logger),
			rpc.WithLogUnaryInterceptor(input.Logger),
			rpc.WithJWTAuthUnaryInterceptor(verifier, webservice.NewRBACAuthorizer(ctx, ds, cfg.ProjectMap(), input.Logger), input.Logger),
			rpc.WithRequestValidationUnaryInterceptor(),
		}
		if s.tls {
			opts = append(opts, rpc.WithTLS(s.certFile, s.keyFile))
		}
		if s.enableGRPCReflection {
			opts = append(opts, rpc.WithGRPCReflection())
		}
		if input.Flags.Metrics {
			opts = append(opts, rpc.WithPrometheusUnaryInterceptor())
		}

		server := rpc.NewServer(service, opts...)
		group.Go(func() error {
			return server.Run(ctx)
		})
	}

	// Start an http server for handling incoming http requests
	// such as auth callbacks, webhook events and
	// serving static assets for web.
	{
		signer, err := jwt.NewSigner(defaultSigningMethod, s.encryptionKeyFile)
		if err != nil {
			input.Logger.Error("failed to create a new signer", zap.Error(err))
			return err
		}

		h := httpapi.NewHandler(
			signer,
			s.staticDir,
			encryptDecrypter,
			cfg.Address,
			cfg.StateKey,
			cfg.ProjectMap(),
			cfg.SharedSSOConfigMap(),
			datastore.NewProjectStore(ds, datastore.WebCommander),
			!s.insecureCookie,
			input.Logger,
		)
		httpServer := &http.Server{
			Addr:    fmt.Sprintf(":%d", s.httpPort),
			Handler: h,
		}

		group.Go(func() error {
			return runHTTPServer(ctx, httpServer, s.gracePeriod, input.Logger)
		})
	}

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
		admin.Handle("/metrics", input.PrometheusMetricsHandlerFor(reg))

		group.Go(func() error {
			return admin.Run(ctx)
		})
	}

	// Start a gRPC server for handling envoy ext_authz requests.
	{
		var (
			verifier = pipedverifier.NewVerifier(
				ctx,
				cfg,
				// These stores are used to handle request over envoy ext_authz, thus the writer should be PipedWriter.
				datastore.NewProjectStore(ds, datastore.PipedCommander),
				datastore.NewPipedStore(ds, datastore.PipedCommander),
				input.Logger,
			)
			service = grpcapi.NewEnvoyAuthorizationServer(verifier)
			opts    = []rpc.Option{
				rpc.WithPort(s.envoyAuthzPort),
				rpc.WithGracePeriod(s.gracePeriod),
				rpc.WithLogger(input.Logger),
				rpc.WithLogUnaryInterceptor(input.Logger),
			}
		)

		if s.tls {
			opts = append(opts, rpc.WithTLS(s.certFile, s.keyFile))
		}
		if s.enableGRPCReflection {
			opts = append(opts, rpc.WithGRPCReflection())
		}
		if input.Flags.Metrics {
			opts = append(opts, rpc.WithPrometheusUnaryInterceptor())
		}

		server := rpc.NewServer(service, opts...)
		group.Go(func() error {
			return server.Run(ctx)
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

func runHTTPServer(ctx context.Context, httpServer *http.Server, gracePeriod time.Duration, logger *zap.Logger) error {
	doneCh := make(chan error, 1)
	ctx, cancel := context.WithCancel(ctx)

	go func() {
		defer cancel()
		logger.Info(fmt.Sprintf("start running http server on %s", httpServer.Addr))
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("failed to listen and http server", zap.Error(err))
			doneCh <- err
		}
		doneCh <- nil
	}()

	<-ctx.Done()

	ctx, cancel = context.WithTimeout(context.Background(), gracePeriod)
	defer cancel()
	logger.Info("stopping http server")
	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Error("failed to shutdown http server", zap.Error(err))
	}

	return <-doneCh
}

func loadConfig(file string) (*config.ControlPlaneSpec, error) {
	cfg, err := config.LoadFromYAML(file)
	if err != nil {
		return nil, err
	}
	if cfg.Kind != config.KindControlPlane {
		return nil, fmt.Errorf("wrong configuration kind for control-plane: %v", cfg.Kind)
	}
	return cfg.ControlPlaneSpec, nil
}

func createDatastore(ctx context.Context, cfg *config.ControlPlaneSpec, fs filestore.Store, c cache.Cache, logger *zap.Logger) (datastore.DataStore, error) {
	switch cfg.Datastore.Type {
	case model.DataStoreFirestore:
		fsConfig := cfg.Datastore.FirestoreConfig
		options := []firestore.Option{
			firestore.WithCredentialsFile(fsConfig.CredentialsFile),
			firestore.WithLogger(logger),
		}
		if p := fsConfig.CollectionNamePrefix; p != "" {
			options = append(options, firestore.WithCollectionNamePrefix(p))
		}
		return firestore.NewFireStore(ctx, fsConfig.Project, fsConfig.Namespace, fsConfig.Environment, options...)

	case model.DataStoreMySQL:
		mqConfig := cfg.Datastore.MySQLConfig
		options := []mysql.Option{
			mysql.WithLogger(logger),
		}
		if mqConfig.UsernameFile != "" || mqConfig.PasswordFile != "" {
			options = append(options, mysql.WithAuthenticationFile(mqConfig.UsernameFile, mqConfig.PasswordFile))
		}
		return mysql.NewMySQL(mqConfig.URL, mqConfig.Database, options...)
	case model.DataStoreFileDB:
		options := []filedb.Option{
			filedb.WithLogger(logger),
		}
		return filedb.NewFileDB(fs, c, options...)
	default:
		return nil, fmt.Errorf("unknown datastore type %q", cfg.Datastore.Type)
	}
}

func createFilestore(ctx context.Context, cfg *config.ControlPlaneSpec, logger *zap.Logger) (filestore.Store, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	switch cfg.Filestore.Type {
	case model.FileStoreGCS:
		gcsCfg := cfg.Filestore.GCSConfig
		options := []gcs.Option{
			gcs.WithLogger(logger),
		}
		if gcsCfg.CredentialsFile != "" {
			options = append(options, gcs.WithCredentialsFile(gcsCfg.CredentialsFile))
		}
		return gcs.NewStore(ctx, gcsCfg.Bucket, options...)

	case model.FileStoreS3:
		s3Cfg := cfg.Filestore.S3Config
		options := []s3.Option{
			s3.WithLogger(logger),
		}
		if s3Cfg.CredentialsFile != "" {
			options = append(options, s3.WithCredentialsFile(s3Cfg.CredentialsFile, s3Cfg.Profile))
		}
		if s3Cfg.RoleARN != "" && s3Cfg.TokenFile != "" {
			options = append(options, s3.WithTokenFile(s3Cfg.RoleARN, s3Cfg.TokenFile))
		}
		return s3.NewStore(ctx, s3Cfg.Region, s3Cfg.Bucket, options...)

	case model.FileStoreMINIO:
		minioCfg := cfg.Filestore.MinioConfig
		options := []minio.Option{
			minio.WithLogger(logger),
		}
		s, err := minio.NewStore(minioCfg.Endpoint, minioCfg.Bucket, minioCfg.AccessKeyFile, minioCfg.SecretKeyFile, options...)
		if err != nil {
			return nil, fmt.Errorf("failed to generate minio store: %w", err)
		}
		if minioCfg.AutoCreateBucket {
			if err := s.EnsureBucket(ctx); err != nil {
				return nil, fmt.Errorf("failed to ensure bucket: %w", err)
			}
		}
		return s, nil

	default:
		return nil, fmt.Errorf("unknown filestore type %q", cfg.Filestore.Type)
	}
}

func registerMetrics() *prometheus.Registry {
	r := prometheus.NewRegistry()
	wrapped := prometheus.WrapRegistererWith(map[string]string{
		"pipecd_component": "server",
	}, r)

	wrapped.Register(collectors.NewGoCollector())
	wrapped.Register(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	wrapped.Register(grpc_prometheus.DefaultServerMetrics)

	cachemetrics.Register(wrapped)
	httpapimetrics.Register(wrapped)
	grpcapimetrics.Register(wrapped)

	return r
}
