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

package cli

import (
	"context"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"cloud.google.com/go/profiler"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"google.golang.org/api/option"

	"github.com/pipe-cd/pipe/pkg/log"
	"github.com/pipe-cd/pipe/pkg/version"
)

type Telemetry struct {
	Logger *zap.Logger
	Flags  TelemetryFlags
}

type Runner func(ctx context.Context, telemetry Telemetry) error

func WithContext(runner Runner) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		defer signal.Stop(ch)
		return runWithContext(cmd, ch, runner)
	}
}

func runWithContext(cmd *cobra.Command, signalCh <-chan os.Signal, runner Runner) error {
	flags, err := parseTelemetryFlags(cmd.Flags())
	if err != nil {
		return err
	}
	telemetry := Telemetry{
		Flags: flags,
	}
	service := extractServiceName(cmd)
	version := version.Get()

	// Initialize logger.
	logger, err := newLogger(service, version.Version, flags.LogLevel, flags.LogEncoding)
	if err != nil {
		return err
	}
	defer logger.Sync()
	telemetry.Logger = logger

	// Start running profiler.
	if flags.Profile {
		if err := startProfiler(service, version.Version, flags.ProfilerCredentialsFile, flags.ProfileDebugLogging, logger); err != nil {
			logger.Error("failed to run profiler", zap.Error(err))
			return err
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		select {
		case s := <-signalCh:
			logger.Info("stopping due to signal", zap.Any("signal", s))
			cancel()
		case <-ctx.Done():
		}
	}()

	return runner(ctx, telemetry)
}

func newLogger(service, version, level, encoding string) (*zap.Logger, error) {
	configs := log.DefaultConfigs
	configs.ServiceContext = &log.ServiceContext{
		Service: service,
		Version: version,
	}
	configs.Level = level
	configs.Encoding = log.EncodingType(encoding)
	return log.NewLogger(configs)
}

func startProfiler(service, version, credentialsFile string, debugLogging bool, logger *zap.Logger) error {
	var options []option.ClientOption
	if credentialsFile != "" {
		options = append(options, option.WithCredentialsFile(credentialsFile))
	}
	config := profiler.Config{
		Service:        service,
		ServiceVersion: version,
		DebugLogging:   debugLogging,
	}

	logger.Info("start running profiler", zap.String("service", service))
	return profiler.Start(config, options...)
}

func (t Telemetry) PrometheusMetricsHandler() http.Handler {
	if t.Flags.Metrics {
		return promhttp.Handler()
	}
	var empty http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(""))
	}
	return empty
}

func (t Telemetry) PrometheusMetricsHandlerFor(r *prometheus.Registry) http.Handler {
	if t.Flags.Metrics {
		return promhttp.HandlerFor(r, promhttp.HandlerOpts{})
	}
	var empty http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(""))
	}
	return empty
}

type MetricsBuilder interface {
	Build() (io.ReadCloser, error)
}

func (t Telemetry) CustomedMetricsHandlerFor(mb MetricsBuilder) http.Handler {
	if t.Flags.Metrics {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rc, err := mb.Build()
			if err != nil {
				http.NotFound(w, r)
			}
			_, err = io.Copy(w, rc)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		})
	}
	var empty http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(""))
	}
	return empty
}

func extractServiceName(cmd *cobra.Command) string {
	return strings.Replace(cmd.CommandPath(), " ", ".", -1)
}
