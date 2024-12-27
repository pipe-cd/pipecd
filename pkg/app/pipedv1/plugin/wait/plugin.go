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
	"net/http/pprof"
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/pipe-cd/pipecd/pkg/admin"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/wait/deployment"
	"github.com/pipe-cd/pipecd/pkg/cli"
	config "github.com/pipe-cd/pipecd/pkg/configv1"
	"github.com/pipe-cd/pipecd/pkg/plugin/logpersister"
	"github.com/pipe-cd/pipecd/pkg/plugin/pipedapi"
	"github.com/pipe-cd/pipecd/pkg/rpc"
	"github.com/pipe-cd/pipecd/pkg/version"
)

type plugin struct {
	pipedPluginService string
	gracePeriod        time.Duration
	tls                bool
	certFile           string
	keyFile            string
	config             string

	enableGRPCReflection bool
}

// newPluginCommand creates a new cobra command for executing api server.
func newPluginCommand() *cobra.Command {
	s := &plugin{
		gracePeriod: 30 * time.Second,
	}
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start running the wait-stage-plugin.",
		RunE:  cli.WithContext(s.run),
	}

	cmd.Flags().StringVar(&s.pipedPluginService, "piped-plugin-service", s.pipedPluginService, "The port number used to connect to the piped plugin service.") // TODO: we should discuss about the name of this flag, or we should use environment variable instead.
	cmd.Flags().StringVar(&s.config, "config", s.config, "The configuration for the plugin.")
	cmd.Flags().DurationVar(&s.gracePeriod, "grace-period", s.gracePeriod, "How long to wait for graceful shutdown.")

	cmd.Flags().BoolVar(&s.tls, "tls", s.tls, "Whether running the gRPC server with TLS or not.")
	cmd.Flags().StringVar(&s.certFile, "cert-file", s.certFile, "The path to the TLS certificate file.")
	cmd.Flags().StringVar(&s.keyFile, "key-file", s.keyFile, "The path to the TLS key file.")

	// For debugging early in development
	cmd.Flags().BoolVar(&s.enableGRPCReflection, "enable-grpc-reflection", s.enableGRPCReflection, "Whether to enable the reflection service or not.")

	cmd.MarkFlagRequired("piped-plugin-service")
	cmd.MarkFlagRequired("config")

	return cmd
}

func (s *plugin) run(ctx context.Context, input cli.Input) (runErr error) {
	// Make a cancellable context.
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	group, ctx := errgroup.WithContext(ctx)

	pipedapiClient, err := pipedapi.NewClient(ctx, s.pipedPluginService)
	if err != nil {
		input.Logger.Error("failed to create piped plugin service client", zap.Error(err))
		return err
	}

	// Load the configuration.
	cfg, err := config.ParsePluginConfig(s.config)
	if err != nil {
		input.Logger.Error("failed to parse the configuration", zap.Error(err))
		return err
	}

	// Start running admin server.
	{
		var (
			ver   = []byte(version.Get().Version)                  // TODO: get the plugin's version
			admin = admin.NewAdmin(0, s.gracePeriod, input.Logger) // TODO: add config for admin port
		)

		admin.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
			w.Write(ver)
		})
		admin.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("ok"))
		})
		admin.HandleFunc("/debug/pprof/", pprof.Index)
		admin.HandleFunc("/debug/pprof/profile", pprof.Profile)
		admin.HandleFunc("/debug/pprof/trace", pprof.Trace)

		group.Go(func() error {
			return admin.Run(ctx)
		})
	}

	// Start log persister
	persister := logpersister.NewPersister(pipedapiClient, input.Logger)
	group.Go(func() error {
		return persister.Run(ctx)
	})

	// Start a gRPC server for handling external API requests.
	{
		var (
			service = deployment.NewDeploymentService(
				cfg,
				input.Logger,
				persister,
			)
			opts = []rpc.Option{
				rpc.WithPort(cfg.Port),
				rpc.WithGracePeriod(s.gracePeriod),
				rpc.WithLogger(input.Logger),
				rpc.WithLogUnaryInterceptor(input.Logger),
				rpc.WithRequestValidationUnaryInterceptor(),
				rpc.WithSignalHandlingUnaryInterceptor(),
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

	if err := group.Wait(); err != nil {
		input.Logger.Error("failed while running", zap.Error(err))
		return err
	}
	return nil
}
