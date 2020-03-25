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
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/kapetaniosci/pipe/pkg/admin"
	"github.com/kapetaniosci/pipe/pkg/app/livelog/api"
	"github.com/kapetaniosci/pipe/pkg/app/livelog/cleaner"
	"github.com/kapetaniosci/pipe/pkg/cli"
	"github.com/kapetaniosci/pipe/pkg/rpc"
)

const (
	dataPath = "/data"
)

type server struct {
	grpcPort                int
	adminPort               int
	gracePeriod             time.Duration
	cleanerPeriod           time.Duration
	cleanerIterationTimeout time.Duration
	cleanerMaxTTL           time.Duration

	tls      bool
	certFile string
	keyFile  string
}

func NewCommand() *cobra.Command {
	s := &server{
		grpcPort:                9080,
		adminPort:               9085,
		gracePeriod:             15 * time.Second,
		cleanerPeriod:           time.Hour,
		cleanerIterationTimeout: 30 * time.Minute,
		cleanerMaxTTL:           6 * time.Hour,
	}
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Start running livelog server.",
		RunE:  cli.WithContext(s.run),
	}

	cmd.Flags().IntVar(&s.grpcPort, "grpc-port", s.grpcPort, "The port number used to run grpc server.")
	cmd.Flags().IntVar(&s.adminPort, "admin-port", s.adminPort, "The port number used to run a HTTP server for admin tasks such as metrics, healthz.")
	cmd.Flags().DurationVar(&s.gracePeriod, "grace-period", s.gracePeriod, "How long to wait for graceful shutdown.")

	cmd.Flags().DurationVar(&s.cleanerPeriod, "cleaner-period", s.cleanerPeriod, "The interval between two cleaning times of cleaner.")
	cmd.Flags().DurationVar(&s.cleanerIterationTimeout, "cleaner-iteration-timeout", s.cleanerIterationTimeout, "The maximum time one iteration of cleaner will take.")
	cmd.Flags().DurationVar(&s.cleanerMaxTTL, "cleaner-max-ttl", s.cleanerMaxTTL, "The maximum TTL of a log file.")

	cmd.Flags().BoolVar(&s.tls, "tls", s.tls, "Whether running the gRPC server with TLS or not.")
	cmd.Flags().StringVar(&s.certFile, "cert-file", s.certFile, "The path to the TLS certificate file.")
	cmd.Flags().StringVar(&s.keyFile, "key-file", s.keyFile, "The path to the TLS key file.")

	return cmd
}

func (s *server) run(ctx context.Context, t cli.Telemetry) error {
	doneCh := make(chan error, 3)

	// Ensure the existence of data directory.
	if _, err := os.Stat(dataPath); os.IsNotExist(err) {
		if err := os.Mkdir(dataPath, os.ModePerm); err != nil {
			t.Logger.Error("failed to create data directory", zap.Error(err))
			return err
		}
	}

	// Start grpc server.
	service := api.NewLiveLogService(
		dataPath,
		api.WithLogger(t.Logger),
	)
	opts := []rpc.Option{
		rpc.WithPort(s.grpcPort),
		rpc.WithLogger(t.Logger),
		//rpc.WithRunnerKeyAuthUnaryInterceptor(),
		//rpc.WithRunnerKeyAuthStreamInterceptor(),
		rpc.WithRequestValidationUnaryInterceptor(),
	}
	if s.tls {
		opts = append(opts, rpc.WithTLS(s.certFile, s.keyFile))
	}
	server := rpc.NewServer(service, opts...)
	go func() {
		doneCh <- server.Run()
	}()
	defer server.Stop(s.gracePeriod)

	// Start cleaner.
	cln := cleaner.NewCleaner(
		dataPath,
		cleaner.WithPeriod(s.cleanerPeriod),
		cleaner.WithIterationTimeout(s.cleanerIterationTimeout),
		cleaner.WithMaxTTL(s.cleanerMaxTTL),
		cleaner.WithLogger(t.Logger),
	)
	go func() {
		doneCh <- cln.Run()
	}()
	defer cln.Stop()

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
