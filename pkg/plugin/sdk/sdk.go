// Copyright 2025 The PipeCD Authors.
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

// package pipedsdk provides software development kits for building PipeCD piped plugins.
package sdk

import (
	"context"
	"fmt"
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/pipe-cd/pipecd/pkg/admin"
	"github.com/pipe-cd/pipecd/pkg/cli"
	config "github.com/pipe-cd/pipecd/pkg/configv1"
	"github.com/pipe-cd/pipecd/pkg/plugin/logpersister"
	"github.com/pipe-cd/pipecd/pkg/plugin/pipedapi"
	"github.com/pipe-cd/pipecd/pkg/rpc"
)

// Plugin is the interface that must be implemented by a piped plugin.
type Plugin interface {
	// Name returns the name of the plugin.
	Name() string
	// Version returns the version of the plugin.
	Version() string
}

// TODO is a placeholder for the real type.
// This type will be replaced by the real type when implementing the sdk.
type TODO struct{}

// Run runs the registered plugins.
// It will listen the gRPC server and handle all requests from piped.
func Run() error {
	if deploymentServiceServer == nil { // TODO: support livestate plugin
		return fmt.Errorf("deployment service server is not registered")
	}

	app := cli.NewApp(
		fmt.Sprintf("pipecd-plugin-%s", deploymentServiceServer.Name()),
		"Plugin component for Piped.",
	)
	app.AddCommands(
		NewPluginCommand(),
	)
	if err := app.Run(); err != nil {
		return err
	}

	return nil
}

type plugin struct {
	pipedPluginService string
	gracePeriod        time.Duration
	tls                bool
	certFile           string
	keyFile            string
	config             string

	enableGRPCReflection bool
}

// NewPluginCommand creates a new cobra command for executing api server.
func NewPluginCommand() *cobra.Command {
	s := &plugin{
		gracePeriod: 30 * time.Second,
	}
	cmd := &cobra.Command{
		Use:   "start",
		Short: fmt.Sprintf("Start running a %s plugin.", deploymentServiceServer.Name()),
		RunE:  cli.WithContext(s.run),
	}

	cmd.Flags().StringVar(&s.pipedPluginService, "piped-plugin-service", s.pipedPluginService, "The address used to connect to the piped plugin service.")
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
			ver   = []byte(deploymentServiceServer.Version())      // TODO: not use the deploymentServiceServer directly
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
		deploymentServiceServer.setCommonFields(commonFields{
			config:       cfg,
			logger:       input.Logger.Named("deployment-service"),
			logPersister: persister,
			client:       pipedapiClient,
		})
		if err := deploymentServiceServer.setConfig(cfg.Config); err != nil {
			input.Logger.Error("failed to set configuration", zap.Error(err))
			return err
		}
		var (
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

		server := rpc.NewServer(deploymentServiceServer, opts...)
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
