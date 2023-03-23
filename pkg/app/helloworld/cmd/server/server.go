// Copyright 2023 The PipeCD Authors.
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
	"net/http"
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/pipe-cd/pipecd/pkg/admin"
	"github.com/pipe-cd/pipecd/pkg/app/helloworld/api"
	"github.com/pipe-cd/pipecd/pkg/cli"
	"github.com/pipe-cd/pipecd/pkg/rpc"
	"github.com/pipe-cd/pipecd/pkg/version"
)

type server struct {
	grpcPort    int
	adminPort   int
	gracePeriod time.Duration
}

func NewCommand() *cobra.Command {
	s := &server{
		grpcPort:    9080,
		adminPort:   9085,
		gracePeriod: 15 * time.Second,
	}
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Start running HelloWorld server.",
		RunE:  cli.WithContext(s.run),
	}
	cmd.Flags().IntVar(&s.grpcPort, "grpc-port", s.grpcPort, "The port number used to run grpc server.")
	cmd.Flags().IntVar(&s.adminPort, "admin-port", s.adminPort, "The port number used to run a HTTP server for admin tasks such as metrics, healthz.")
	cmd.Flags().DurationVar(&s.gracePeriod, "grace-period", s.gracePeriod, "How long to wait for graceful shutdown.")
	return cmd
}

func (s *server) run(ctx context.Context, input cli.Input) error {
	group, ctx := errgroup.WithContext(ctx)

	// Start running gRPC server.
	{
		service := api.NewHelloWorldAPI(
			api.WithLogger(input.Logger),
		)
		server := rpc.NewServer(service,
			rpc.WithPort(s.grpcPort),
			rpc.WithGracePeriod(s.gracePeriod),
			rpc.WithLogger(input.Logger),
		)
		group.Go(func() error {
			return server.Run(ctx)
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
		admin.Handle("/metrics", input.PrometheusMetricsHandler())

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
