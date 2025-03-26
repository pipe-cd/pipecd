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
	"github.com/pipe-cd/pipecd/pkg/plugin/toolregistry"
	"github.com/pipe-cd/pipecd/pkg/rpc"
)

// PluginOption is a function that configures the plugin.
type PluginOption[Config, DeployTargetConfig any] func(*Plugin[Config, DeployTargetConfig])

// WithStagePlugin is a function that sets the stage plugin.
// This is mutually exclusive with WithDeploymentPlugin.
func WithStagePlugin[Config, DeployTargetConfig any](stagePlugin StagePlugin[Config, DeployTargetConfig]) PluginOption[Config, DeployTargetConfig] {
	return func(plugin *Plugin[Config, DeployTargetConfig]) {
		plugin.stagePlugin = stagePlugin
	}
}

// WithDeploymentPlugin is a function that sets the deployment plugin.
// This is mutually exclusive with WithStagePlugin.
func WithDeploymentPlugin[Config, DeployTargetConfig any](deploymentPlugin DeploymentPlugin[Config, DeployTargetConfig]) PluginOption[Config, DeployTargetConfig] {
	return func(plugin *Plugin[Config, DeployTargetConfig]) {
		plugin.deploymentPlugin = deploymentPlugin
	}
}

// WithLivestatePlugin is a function that sets the livestate plugin.
func WithLivestatePlugin[Config, DeployTargetConfig any](livestatePlugin LivestatePlugin[Config, DeployTargetConfig]) PluginOption[Config, DeployTargetConfig] {
	return func(plugin *Plugin[Config, DeployTargetConfig]) {
		plugin.livestatePlugin = livestatePlugin
	}
}

// Plugin is a wrapper for the plugin.
// It provides a way to run the plugin with the given config and deploy target config.
type Plugin[Config, DeployTargetConfig any] struct {
	// plugin info
	name    string
	version string

	// plugin implementations
	stagePlugin      StagePlugin[Config, DeployTargetConfig]
	deploymentPlugin DeploymentPlugin[Config, DeployTargetConfig]
	livestatePlugin  LivestatePlugin[Config, DeployTargetConfig]

	// command line options
	pipedPluginService   string
	gracePeriod          time.Duration
	tls                  bool
	certFile             string
	keyFile              string
	config               string
	enableGRPCReflection bool
}

// NewPlugin creates a new plugin.
func NewPlugin[Config, DeployTargetConfig any](name, version string, options ...PluginOption[Config, DeployTargetConfig]) (*Plugin[Config, DeployTargetConfig], error) {
	plugin := &Plugin[Config, DeployTargetConfig]{
		name:    name,
		version: version,

		// Default values of command line options
		gracePeriod: 30 * time.Second,
	}

	for _, option := range options {
		option(plugin)
	}

	if plugin.stagePlugin == nil && plugin.deploymentPlugin == nil && plugin.livestatePlugin == nil {
		return nil, fmt.Errorf("at least one plugin must be registered")
	}

	if _, ok := plugin.stagePlugin.(DeploymentPlugin[Config, DeployTargetConfig]); ok {
		return nil, fmt.Errorf("stage plugin cannot be a deployment plugin, you must use WithDeploymentPlugin instead")
	}

	if plugin.stagePlugin != nil && plugin.deploymentPlugin != nil {
		return nil, fmt.Errorf("stage plugin and deployment plugin cannot be registered at the same time")
	}

	return plugin, nil
}

// Run runs the plugin.
func (p *Plugin[Config, DeployTargetConfig]) Run() error {
	app := cli.NewApp(
		fmt.Sprintf("pipecd-plugin-%s", p.name),
		"Plugin component for Piped.",
	)

	app.AddCommands(
		p.command(),
	)

	if err := app.Run(); err != nil {
		return err
	}

	return nil
}

// command returns the cobra command for the plugin.
func (p *Plugin[Config, DeployTargetConfig]) command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: fmt.Sprintf("Start running a %s plugin.", p.name),
		RunE:  cli.WithContext(p.run),
	}

	cmd.Flags().StringVar(&p.pipedPluginService, "piped-plugin-service", p.pipedPluginService, "The address used to connect to the piped plugin service.")
	cmd.Flags().StringVar(&p.config, "config", p.config, "The configuration for the plugin.")
	cmd.Flags().DurationVar(&p.gracePeriod, "grace-period", p.gracePeriod, "How long to wait for graceful shutdown.")

	cmd.Flags().BoolVar(&p.tls, "tls", p.tls, "Whether running the gRPC server with TLS or not.")
	cmd.Flags().StringVar(&p.certFile, "cert-file", p.certFile, "The path to the TLS certificate file.")
	cmd.Flags().StringVar(&p.keyFile, "key-file", p.keyFile, "The path to the TLS key file.")

	// For debugging early in development
	cmd.Flags().BoolVar(&p.enableGRPCReflection, "enable-grpc-reflection", p.enableGRPCReflection, "Whether to enable the reflection service or not.")

	cmd.MarkFlagRequired("piped-plugin-service")
	cmd.MarkFlagRequired("config")

	return cmd
}

// run is the entrypoint of the plugin.
func (p *Plugin[Config, DeployTargetConfig]) run(ctx context.Context, input cli.Input) error {
	var (
		stagePluginServiceServer      *StagePluginServiceServer[Config, DeployTargetConfig]
		deploymentPluginServiceServer *DeploymentPluginServiceServer[Config, DeployTargetConfig]
		livestatePluginServiceServer  *LivestatePluginServer[Config, DeployTargetConfig]
	)

	if p.stagePlugin != nil && p.deploymentPlugin != nil {
		return fmt.Errorf("stage plugin and deployment plugin cannot be registered at the same time")
	}

	if p.stagePlugin != nil {
		stagePluginServiceServer = &StagePluginServiceServer[Config, DeployTargetConfig]{base: p.stagePlugin}
	}

	if p.deploymentPlugin != nil {
		deploymentPluginServiceServer = &DeploymentPluginServiceServer[Config, DeployTargetConfig]{base: p.deploymentPlugin}
	}

	if p.livestatePlugin != nil {
		livestatePluginServiceServer = &LivestatePluginServer[Config, DeployTargetConfig]{base: p.livestatePlugin}
	}

	// Make a cancellable context.
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	group, ctx := errgroup.WithContext(ctx)

	pipedapiClient, err := pipedapi.NewClient(ctx, p.pipedPluginService)
	if err != nil {
		input.Logger.Error("failed to create piped plugin service client", zap.Error(err))
		return err
	}

	// Load the configuration.
	cfg, err := config.ParsePluginConfig(p.config)
	if err != nil {
		input.Logger.Error("failed to parse the configuration", zap.Error(err))
		return err
	}

	// Start running admin server.
	{
		var (
			ver   = []byte(p.version)
			admin = admin.NewAdmin(0, p.gracePeriod, input.Logger) // TODO: add config for admin port
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
		commonFields := commonFields{
			name:         p.name,
			version:      p.version,
			config:       cfg,
			logPersister: persister,
			client:       pipedapiClient,
			toolRegistry: toolregistry.NewToolRegistry(pipedapiClient),
		}

		var services []rpc.Service

		if p.stagePlugin != nil {
			if err := stagePluginServiceServer.setFields(
				commonFields.withLogger(input.Logger.Named("stage-service")),
			); err != nil {
				input.Logger.Error("failed to set fields", zap.Error(err))
				return err
			}
			services = append(services, stagePluginServiceServer)
		}

		if p.deploymentPlugin != nil {
			if err := deploymentPluginServiceServer.setFields(
				commonFields.withLogger(input.Logger.Named("deployment-service")),
			); err != nil {
				input.Logger.Error("failed to set fields", zap.Error(err))
				return err
			}
			services = append(services, deploymentPluginServiceServer)
		}

		if p.livestatePlugin != nil {
			if err := livestatePluginServiceServer.setFields(
				commonFields.withLogger(input.Logger.Named("livestate-service")),
			); err != nil {
				input.Logger.Error("failed to set fields", zap.Error(err))
				return err
			}
			services = append(services, livestatePluginServiceServer)
		}

		if len(services) == 0 {
			// This is promised in the NewPlugin function.
			// When this happens, it means that *Plugin was initialized without using NewPlugin.
			return fmt.Errorf("no plugin is registered, you must use NewPlugin to initialize the plugin")
		}

		var (
			opts = []rpc.Option{
				rpc.WithPort(cfg.Port),
				rpc.WithGracePeriod(p.gracePeriod),
				rpc.WithLogger(input.Logger),
				rpc.WithLogUnaryInterceptor(input.Logger),
				rpc.WithRequestValidationUnaryInterceptor(),
				rpc.WithSignalHandlingUnaryInterceptor(),
			}
		)
		if p.tls {
			opts = append(opts, rpc.WithTLS(p.certFile, p.keyFile))
		}
		if p.enableGRPCReflection {
			opts = append(opts, rpc.WithGRPCReflection())
		}
		if input.Flags.Metrics {
			opts = append(opts, rpc.WithPrometheusUnaryInterceptor())
		}
		if len(services) > 1 {
			for _, service := range services[1:] {
				opts = append(opts, rpc.WithService(service))
			}
		}

		server := rpc.NewServer(services[0], opts...)

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
