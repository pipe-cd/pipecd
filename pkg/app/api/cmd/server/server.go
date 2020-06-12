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

	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/pipe-cd/pipe/pkg/admin"
	"github.com/pipe-cd/pipe/pkg/app/api/api"
	"github.com/pipe-cd/pipe/pkg/app/api/applicationlivestatestore"
	"github.com/pipe-cd/pipe/pkg/app/api/stagelogstore"
	"github.com/pipe-cd/pipe/pkg/cache/rediscache"
	"github.com/pipe-cd/pipe/pkg/cli"
	"github.com/pipe-cd/pipe/pkg/config"
	"github.com/pipe-cd/pipe/pkg/datastore"
	"github.com/pipe-cd/pipe/pkg/datastore/firestore"
	"github.com/pipe-cd/pipe/pkg/filestore"
	"github.com/pipe-cd/pipe/pkg/filestore/gcs"
	"github.com/pipe-cd/pipe/pkg/redis"
	"github.com/pipe-cd/pipe/pkg/rpc"
)

var (
	defaultSigningMethod = jwtgo.SigningMethodHS256
)

type httpHandler interface {
	Register(func(pattern string, handler func(http.ResponseWriter, *http.Request)))
}

type server struct {
	pipedAPIPort int
	webAPIPort   int
	httpPort     int
	adminPort    int
	gracePeriod  time.Duration

	tls                 bool
	certFile            string
	keyFile             string
	tokenSigningKeyFile string

	configFile string

	useFakeResponse      bool
	enableGRPCReflection bool
}

func NewCommand() *cobra.Command {
	s := &server{
		pipedAPIPort: 9080,
		webAPIPort:   9081,
		httpPort:     9082,
		adminPort:    9085,
		gracePeriod:  30 * time.Second,
	}
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Start running API server.",
		RunE:  cli.WithContext(s.run),
	}

	cmd.Flags().IntVar(&s.pipedAPIPort, "piped-api-port", s.pipedAPIPort, "The port number used to run a grpc server that serving serves incoming piped requests.")
	cmd.Flags().IntVar(&s.webAPIPort, "web-api-port", s.webAPIPort, "The port number used to run a grpc server that serves incoming web requests.")
	cmd.Flags().IntVar(&s.httpPort, "http-port", s.httpPort, "The port number used to run a http server that serves incoming http requests such as auth callbacks or webhook events.")
	cmd.Flags().IntVar(&s.adminPort, "admin-port", s.adminPort, "The port number used to run a HTTP server for admin tasks such as metrics, healthz.")
	cmd.Flags().DurationVar(&s.gracePeriod, "grace-period", s.gracePeriod, "How long to wait for graceful shutdown.")

	cmd.Flags().BoolVar(&s.tls, "tls", s.tls, "Whether running the gRPC server with TLS or not.")
	cmd.Flags().StringVar(&s.certFile, "cert-file", s.certFile, "The path to the TLS certificate file.")
	cmd.Flags().StringVar(&s.keyFile, "key-file", s.keyFile, "The path to the TLS key file.")
	cmd.Flags().StringVar(&s.tokenSigningKeyFile, "token-signing-key-file", s.tokenSigningKeyFile, "The path to key file used to sign ID token.")
	cmd.Flags().StringVar(&s.configFile, "config-file", s.configFile, "The path to the configuration file.")

	// For debugging early in development
	cmd.Flags().BoolVar(&s.useFakeResponse, "use-fake-response", s.useFakeResponse, "Whether the server responds fake response or not.")
	cmd.Flags().BoolVar(&s.enableGRPCReflection, "enable-grpc-reflection", s.enableGRPCReflection, "Whether to enable the reflection service or not.")

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

	// signer, err := jwt.NewSigner(defaultSigningMethod, s.tokenSigningKeyFile)
	// if err != nil {
	// 	t.Logger.Error("failed to create a new signer", zap.Error(err))
	// 	return err
	// }

	// Left comment out until authentication is ready
	//verifier, err := jwt.NewVerifier(defaultSigningMethod, s.tokenSigningKeyFile)
	//if err != nil {
	//	t.Logger.Error("failed to create a new verifier", zap.Error(err))
	//	return err
	//}

	var (
		pipedAPIServer *rpc.Server
		webAPIServer   *rpc.Server
	)

	ds, err := s.createDatastore(ctx, cfg, t.Logger)
	if err != nil {
		t.Logger.Error("failed creating datastore", zap.Error(err))
		return err
	}
	defer func() {
		if err := ds.Close(); err != nil {
			t.Logger.Error("failed closing datastore client", zap.Error(err))
		}
	}()

	fs, err := s.createFilestore(ctx, cfg, t.Logger)
	if err != nil {
		t.Logger.Error("failed creating firestore", zap.Error(err))
		return err
	}
	defer func() {
		if err := fs.Close(); err != nil {
			t.Logger.Error("failed closing firestore client", zap.Error(err))
		}
	}()

	rd := redis.NewRedis(cfg.Cache.RedisAddress, "")
	defer func() {
		if err := rd.Close(); err != nil {
			t.Logger.Error("failed closing redis client", zap.Error(err))
		}
	}()
	cache := rediscache.NewTTLCache(rd, cfg.Cache.TTL.Duration())
	sls := stagelogstore.NewStore(fs, cache, t.Logger)
	alss := applicationlivestatestore.NewStore(fs, cache, t.Logger)

	// Start a gRPC server for handling PipedAPI requests.
	{
		service := api.NewPipedAPI(ds, sls, t.Logger)
		opts := []rpc.Option{
			rpc.WithPort(s.pipedAPIPort),
			rpc.WithGracePeriod(s.gracePeriod),
			rpc.WithLogger(t.Logger),
			// rpc.WithPipedTokenAuthUnaryInterceptor(verifier, t.Logger),
			rpc.WithRequestValidationUnaryInterceptor(),
		}
		if s.tls {
			opts = append(opts, rpc.WithTLS(s.certFile, s.keyFile))
		}
		if s.enableGRPCReflection {
			opts = append(opts, rpc.WithGRPCReflection())
		}
		pipedAPIServer = rpc.NewServer(service, opts...)
		group.Go(func() error {
			return pipedAPIServer.Run(ctx)
		})
	}

	// Start a gRPC server for handling WebAPI requests.
	{
		service := api.NewWebAPI(ds, sls, alss, s.useFakeResponse, t.Logger)
		opts := []rpc.Option{
			rpc.WithPort(s.webAPIPort),
			rpc.WithGracePeriod(s.gracePeriod),
			rpc.WithLogger(t.Logger),
			// Left comment out until authentication is ready
			// rpc.WithJWTAuthUnaryInterceptor(verifier, webservice.NewRBACAuthorizer(), t.Logger),
			rpc.WithRequestValidationUnaryInterceptor(),
		}
		if s.tls {
			opts = append(opts, rpc.WithTLS(s.certFile, s.keyFile))
		}
		if s.enableGRPCReflection {
			opts = append(opts, rpc.WithGRPCReflection())
		}
		webAPIServer = rpc.NewServer(service, opts...)
		group.Go(func() error {
			return webAPIServer.Run(ctx)
		})
	}

	// Start an http server for handling incoming http requests such as auth callbacks or webhook events.
	{
		mux := http.NewServeMux()
		httpServer := &http.Server{
			Addr:    fmt.Sprintf(":%d", s.httpPort),
			Handler: mux,
		}
		handlers := []httpHandler{
			//authhandler.NewHandler(signer, verifier, config, ac, t.Logger),
		}
		for _, h := range handlers {
			h.Register(mux.HandleFunc)
		}
		group.Go(func() error {
			return runHttpServer(ctx, httpServer, s.gracePeriod, t.Logger)
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

func runHttpServer(ctx context.Context, httpServer *http.Server, gracePeriod time.Duration, logger *zap.Logger) error {
	doneCh := make(chan error, 1)
	ctx, cancel := context.WithCancel(ctx)

	go func() {
		defer cancel()
		logger.Info("start running http server")
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("failed to listen and http server", zap.Error(err))
			doneCh <- err
		}
		doneCh <- nil
	}()

	<-ctx.Done()

	ctx, _ = context.WithTimeout(context.Background(), gracePeriod)
	logger.Info("stopping http server")
	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Error("failed to shutdown http server", zap.Error(err))
	}

	return <-doneCh
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
	if cfg.Datastore.FirestoreConfig != nil {
		fsConfig := cfg.Datastore.FirestoreConfig
		options := []firestore.Option{
			firestore.WithCredentialsFile(fsConfig.CredentialsFile),
			firestore.WithLogger(logger),
		}
		return firestore.NewFireStore(ctx, fsConfig.Project, fsConfig.Namespace, options...)
	}

	if cfg.Datastore.DynamoDBConfig != nil {
		return nil, errors.New("dynamodb is unimplemented now")
	}

	if cfg.Datastore.MongoDBConfig != nil {
		return nil, errors.New("mongodb is unimplemented now")
	}

	//return nil, errors.New("datastore configuration is invalid")
	return nil, nil
}

func (s *server) createFilestore(ctx context.Context, cfg *config.ControlPlaneSpec, logger *zap.Logger) (filestore.Store, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if cfg.Filestore.GCSConfig != nil {
		gcsConfig := cfg.Filestore.GCSConfig
		options := []gcs.Option{
			gcs.WithLogger(logger),
		}
		if gcsConfig.CredentialsFile != "" {
			options = append(options, gcs.WithCredentialsFile(gcsConfig.CredentialsFile))
		}
		client, err := gcs.NewStore(ctx, gcsConfig.Bucket, options...)
		if err != nil {
			return nil, err
		}
		return client, nil
	}

	if cfg.Filestore.S3Config != nil {
		return nil, errors.New("s3 is unimplemented now")
	}

	//return nil, errors.New("filestore configuration is invalid")
	return nil, nil
}
