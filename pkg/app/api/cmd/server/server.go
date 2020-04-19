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
	"fmt"
	"net/http"
	"time"

	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/kapetaniosci/pipe/pkg/admin"
	"github.com/kapetaniosci/pipe/pkg/app/api/api"
	apiservice "github.com/kapetaniosci/pipe/pkg/app/api/service"
	"github.com/kapetaniosci/pipe/pkg/cli"
	"github.com/kapetaniosci/pipe/pkg/jwt"
	"github.com/kapetaniosci/pipe/pkg/rpc"
)

var (
	defaultSigningMethod = jwtgo.SigningMethodHS256
)

type httpHandler interface {
	Register(func(pattern string, handler func(http.ResponseWriter, *http.Request)))
}

type server struct {
	runnerAPIPort int
	webAPIPort    int
	webhookPort   int
	adminPort     int
	gracePeriod   time.Duration

	tls                 bool
	certFile            string
	keyFile             string
	tokenSigningKeyFile string
}

func NewCommand() *cobra.Command {
	s := &server{
		runnerAPIPort: 9080,
		webAPIPort:    9081,
		webhookPort:   9082,
		adminPort:     9085,
		gracePeriod:   15 * time.Second,
	}
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Start running API server.",
		RunE:  cli.WithContext(s.run),
	}

	cmd.Flags().IntVar(&s.runnerAPIPort, "runner-api-port", s.runnerAPIPort, "The port number used to run a grpc server that serving serves incoming runner requests.")
	cmd.Flags().IntVar(&s.webAPIPort, "web-api-port", s.webAPIPort, "The port number used to run a grpc server that serves incoming web requests.")
	cmd.Flags().IntVar(&s.webhookPort, "webhook-port", s.webhookPort, "The port number used to run a http server that serves incoming webhook events.")
	cmd.Flags().IntVar(&s.adminPort, "admin-port", s.adminPort, "The port number used to run a HTTP server for admin tasks such as metrics, healthz.")
	cmd.Flags().DurationVar(&s.gracePeriod, "grace-period", s.gracePeriod, "How long to wait for graceful shutdown.")

	cmd.Flags().BoolVar(&s.tls, "tls", s.tls, "Whether running the gRPC server with TLS or not.")
	cmd.Flags().StringVar(&s.certFile, "cert-file", s.certFile, "The path to the TLS certificate file.")
	cmd.Flags().StringVar(&s.keyFile, "key-file", s.keyFile, "The path to the TLS key file.")
	cmd.Flags().StringVar(&s.tokenSigningKeyFile, "token-signing-key-file", s.tokenSigningKeyFile, "The path to key file used to sign ID token.")

	return cmd
}

func (s *server) run(ctx context.Context, t cli.Telemetry) error {
	group, ctx := errgroup.WithContext(ctx)

	// signer, err := jwt.NewSigner(defaultSigningMethod, s.tokenSigningKeyFile)
	// if err != nil {
	// 	t.Logger.Error("failed to create a new signer", zap.Error(err))
	// 	return err
	// }
	verifier, err := jwt.NewVerifier(defaultSigningMethod, s.tokenSigningKeyFile)
	if err != nil {
		t.Logger.Error("failed to create a new verifier", zap.Error(err))
		return err
	}

	var (
		runnerAPIServer *rpc.Server
		webAPIServer    *rpc.Server
	)

	// Start a gRPC server for handling RunnerAPI requests.
	{
		service := api.NewRunnerAPIService(t.Logger)
		opts := []rpc.Option{
			rpc.WithPort(s.runnerAPIPort),
			rpc.WithGracePeriod(s.gracePeriod),
			rpc.WithLogger(t.Logger),
			// rpc.WithRunnerTokenAuthUnaryInterceptor(verifier, t.Logger),
			rpc.WithRequestValidationUnaryInterceptor(),
		}
		if s.tls {
			opts = append(opts, rpc.WithTLS(s.certFile, s.keyFile))
		}
		runnerAPIServer = rpc.NewServer(service, opts...)
		group.Go(func() error {
			return runnerAPIServer.Run(ctx)
		})
	}

	// Start a gRPC server for handling WebAPI requests.
	{
		service := api.NewWebAPIService(t.Logger)
		opts := []rpc.Option{
			rpc.WithPort(s.runnerAPIPort),
			rpc.WithGracePeriod(s.gracePeriod),
			rpc.WithLogger(t.Logger),
			rpc.WithJWTAuthUnaryInterceptor(verifier, apiservice.NewRBACAuthorizer(), t.Logger),
			rpc.WithRequestValidationUnaryInterceptor(),
		}
		if s.tls {
			opts = append(opts, rpc.WithTLS(s.certFile, s.keyFile))
		}
		webAPIServer = rpc.NewServer(service, opts...)
		group.Go(func() error {
			return webAPIServer.Run(ctx)
		})
	}

	// Start an  http server for handling incoming webhook events.
	{
		mux := http.NewServeMux()
		httpServer := &http.Server{
			Addr:    fmt.Sprintf(":%d", s.webhookPort),
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
