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
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"cloud.google.com/go/profiler"
	"contrib.go.opencensus.io/exporter/prometheus"
	"contrib.go.opencensus.io/exporter/stackdriver"
	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/cobra"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
	"go.uber.org/zap"
	"google.golang.org/api/option"

	"github.com/kapetaniosci/pipe/pkg/log"
	"github.com/kapetaniosci/pipe/pkg/version"
)

const (
	metricsNamespace    = "pipe"
	PrometheusExporter  = "prometheus"
	StackdriverExporter = "stackdriver"
)

type Telemetry struct {
	Logger          *zap.Logger
	MetricsExporter view.Exporter
	TracingExporter trace.Exporter
	Flags           TelemetryFlags
}

type Piped func(ctx context.Context, telemetry Telemetry) error

func WithContext(piped Piped) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		defer signal.Stop(ch)
		return runWithContext(cmd, ch, piped)
	}
}

func runWithContext(cmd *cobra.Command, signalCh <-chan os.Signal, piped Piped) error {
	flags, err := parseTelemetryFlags(cmd.Flags())
	if err != nil {
		return err
	}
	telemetry := Telemetry{
		Flags: flags,
	}
	service := extractServiceName(cmd)
	version := version.Get().Version

	// Initialize logger.
	logger, err := newLogger(service, version, flags.LogLevel, flags.LogEncoding)
	if err != nil {
		return err
	}
	defer logger.Sync()
	telemetry.Logger = logger
	// Start profiler.
	if flags.Profile {
		if err := startProfiler(service, version, flags.ProfilerCredentialsFile, flags.ProfileDebugLogging, logger); err != nil {
			logger.Error("failed to run profiler", zap.Error(err))
			return err
		}
	}
	// Initialize metrics exporter.
	if flags.Metrics {
		exporter, err := newMetricsExporter(flags.MetricsExporter)
		if err != nil {
			logger.Error("failed to create metrics exporter", zap.Error(err))
			return err
		}
		telemetry.MetricsExporter = exporter
		// Ensure that we register it as a stats exporter.
		view.RegisterExporter(exporter)
	}
	// Initialize tracing exporter.
	if flags.Tracing {
		exporter, err := newTracingExporter(flags)
		if err != nil {
			logger.Error("failed to create tracing exporter", zap.Error(err))
			return err
		}
		telemetry.TracingExporter = exporter
		// Ensure that we register it as a trace exporter.
		trace.RegisterExporter(exporter)
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
	logger.Info(fmt.Sprintf("start running %s %s", service, version))
	return piped(ctx, telemetry)
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
	logger.Info("start running profiler", zap.String("service", service))
	config := profiler.Config{
		Service:        service,
		ServiceVersion: version,
		DebugLogging:   debugLogging,
	}
	return profiler.Start(config, options...)
}

func newMetricsExporter(exporter string) (view.Exporter, error) {
	if exporter != PrometheusExporter {
		return nil, fmt.Errorf("unsupported metrics exporter: %s", exporter)
	}
	r := prom.NewRegistry()
	r.MustRegister(
		prom.NewGoCollector(),
		prom.NewProcessCollector(prom.ProcessCollectorOpts{}),
	)
	return prometheus.NewExporter(prometheus.Options{
		Namespace: metricsNamespace,
		Registry:  r,
	})
}

func newTracingExporter(f TelemetryFlags) (trace.Exporter, error) {
	if f.TracingExporter != StackdriverExporter {
		return nil, fmt.Errorf("unsupported tracing exporter: %s", f.TracingExporter)
	}
	if f.StackdriverProjectID == "" {
		return nil, fmt.Errorf("missing stackdriver-project-id")
	}
	var options []option.ClientOption
	if f.StackdriverCredentialsFile != "" {
		options = append(options, option.WithCredentialsFile(f.StackdriverCredentialsFile))
	}
	return stackdriver.NewExporter(stackdriver.Options{
		ProjectID:          f.StackdriverProjectID,
		MetricPrefix:       metricsNamespace,
		TraceClientOptions: options,
	})
}

func (t Telemetry) PrometheusMetricsExporter() (*prometheus.Exporter, bool) {
	exporter, ok := t.MetricsExporter.(*prometheus.Exporter)
	return exporter, ok
}

func extractServiceName(cmd *cobra.Command) string {
	return strings.Replace(cmd.CommandPath(), " ", ".", -1)
}
