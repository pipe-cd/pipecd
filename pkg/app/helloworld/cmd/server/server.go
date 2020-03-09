// Copyright 2020 The Dianomi Authors.
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

	"github.com/nghialv/dianomi/pkg/admin"
	"github.com/nghialv/dianomi/pkg/app/helloworld/api"
	"github.com/nghialv/dianomi/pkg/cli"
	"github.com/nghialv/dianomi/pkg/rpc"
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

func (s *server) run(ctx context.Context, t cli.Telemetry) error {
	doneCh := make(chan error)

	// Start grpc server.
	service := api.NewHelloWorldService(
		api.WithLogger(t.Logger),
	)
	server := rpc.NewServer(service,
		rpc.WithPort(s.grpcPort),
		rpc.WithLogger(t.Logger),
	)
	go func() {
		doneCh <- server.Run()
	}()
	defer server.Stop(s.gracePeriod)

	// Start admin server.
	admin := admin.NewAdmin(s.adminPort, t.Logger)
	if exporter, ok := t.PrometheusMetricsExporter(); ok {
		admin.Handle("/metrics", exporter)
	}
	admin.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	go func() {
		doneCh <- admin.Run()
	}()
	defer admin.Stop(s.gracePeriod)

	select {
	case <-ctx.Done():
		return nil
	case err := <-doneCh:
		if err != nil {
			t.Logger.Error("failed while running", zap.Error(err))
			return err
		}
		return nil
	}
}
