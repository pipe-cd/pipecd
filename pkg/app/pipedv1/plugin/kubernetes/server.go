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
	"errors"
	"net"
	"strconv"
	"time"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/deployment"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/toolregistry"
	"github.com/pipe-cd/pipecd/pkg/cli"
	"github.com/pipe-cd/pipecd/pkg/plugin/logpersister"
	"github.com/pipe-cd/pipecd/pkg/plugin/pipedapi"
	"github.com/pipe-cd/pipecd/pkg/rpc"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type server struct {
	apiPort                int
	pipedPluginServicePort int
	gracePeriod            time.Duration
	tls                    bool
	certFile               string
	keyFile                string

	enableGRPCReflection bool
}

// NewServerCommand creates a new cobra command for executing api server.
func NewServerCommand() *cobra.Command {
	s := &server{
		apiPort:                10000,
		pipedPluginServicePort: -1, // default as error value
		gracePeriod:            30 * time.Second,
	}
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Start running server.",
		RunE:  cli.WithContext(s.run),
	}

	cmd.Flags().IntVar(&s.apiPort, "api-port", s.apiPort, "The port number used to run a grpc server for external apis.")
	cmd.Flags().IntVar(&s.pipedPluginServicePort, "piped-plugin-service-port", s.pipedPluginServicePort, "The port number used to connect to the piped plugin service.") // TODO: we should discuss about the name of this flag, or we should use environment variable instead.
	cmd.Flags().DurationVar(&s.gracePeriod, "grace-period", s.gracePeriod, "How long to wait for graceful shutdown.")

	cmd.Flags().BoolVar(&s.tls, "tls", s.tls, "Whether running the gRPC server with TLS or not.")
	cmd.Flags().StringVar(&s.certFile, "cert-file", s.certFile, "The path to the TLS certificate file.")
	cmd.Flags().StringVar(&s.keyFile, "key-file", s.keyFile, "The path to the TLS key file.")

	// For debugging early in development
	cmd.Flags().BoolVar(&s.enableGRPCReflection, "enable-grpc-reflection", s.enableGRPCReflection, "Whether to enable the reflection service or not.")

	return cmd
}

func (s *server) run(ctx context.Context, input cli.Input) (runErr error) {
	// Make a cancellable context.
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	group, ctx := errgroup.WithContext(ctx)

	if s.pipedPluginServicePort == -1 {
		input.Logger.Error("piped-plugin-service-port is required")
		return errors.New("piped-plugin-service-port is required")
	}

	pipedapiClient, err := pipedapi.NewClient(ctx, net.JoinHostPort("localhost", strconv.Itoa(s.pipedPluginServicePort)), nil)
	if err != nil {
		input.Logger.Error("failed to create piped plugin service client", zap.Error(err))
		return err
	}

	// Start a gRPC server for handling external API requests.
	{
		var (
			service = deployment.NewDeploymentService(
				input.Logger,
				toolregistry.NewToolRegistry(pipedapiClient),
				logpersister.NewPersister(pipedapiClient, input.Logger),
			)
			opts = []rpc.Option{
				rpc.WithPort(s.apiPort),
				rpc.WithGracePeriod(s.gracePeriod),
				rpc.WithLogger(input.Logger),
				rpc.WithLogUnaryInterceptor(input.Logger),
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

	if err := group.Wait(); err != nil {
		input.Logger.Error("failed while running", zap.Error(err))
		return err
	}
	return nil
}
