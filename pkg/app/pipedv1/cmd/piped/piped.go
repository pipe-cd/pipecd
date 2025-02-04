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

package piped

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	awssecretsmanager "github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/credentials"
	"sigs.k8s.io/yaml"

	"github.com/pipe-cd/pipecd/pkg/admin"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/apistore/applicationstore"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/apistore/commandstore"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/apistore/deploymentstore"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/apistore/eventstore"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/appconfigreporter"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/cmd/piped/grpcapi"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/controller"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/controller/controllermetrics"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/eventwatcher"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/metadatastore"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/notifier"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/statsreporter"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/trigger"
	"github.com/pipe-cd/pipecd/pkg/app/server/service/pipedservice"
	"github.com/pipe-cd/pipecd/pkg/cli"
	config "github.com/pipe-cd/pipecd/pkg/configv1"
	"github.com/pipe-cd/pipecd/pkg/crypto"
	"github.com/pipe-cd/pipecd/pkg/git"
	"github.com/pipe-cd/pipecd/pkg/lifecycle"
	"github.com/pipe-cd/pipecd/pkg/model"
	pluginapi "github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1"
	"github.com/pipe-cd/pipecd/pkg/rpc"
	"github.com/pipe-cd/pipecd/pkg/rpc/rpcauth"
	"github.com/pipe-cd/pipecd/pkg/rpc/rpcclient"
	"github.com/pipe-cd/pipecd/pkg/version"
)

const (
	commandCheckPeriod time.Duration = 30 * time.Second
)

type piped struct {
	configFile      string
	configData      string
	configGCPSecret string
	configAWSSecret string

	insecure             bool
	certFile             string
	adminPort            int
	pluginServicePort    int
	toolsDir             string
	pluginsDir           string
	gracePeriod          time.Duration
	addLoginUserToPasswd bool
	launcherVersion      string
	maxRecvMsgSize       int
}

func NewCommand() *cobra.Command {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Sprintf("failed to detect the current user's home directory: %v", err))
	}
	p := &piped{
		adminPort:         9085,
		pluginServicePort: 9087,
		toolsDir:          path.Join(home, ".piped", "tools"),
		pluginsDir:        path.Join(home, ".piped", "plugins"),
		gracePeriod:       30 * time.Second,
		maxRecvMsgSize:    1024 * 1024 * 10, // 10MB
	}
	cmd := &cobra.Command{
		Use:   "piped",
		Short: "Start running piped.",
		RunE:  cli.WithContext(p.run),
	}

	cmd.Flags().StringVar(&p.configFile, "config-file", p.configFile, "The path to the configuration file.")
	cmd.Flags().StringVar(&p.configData, "config-data", p.configData, "The base64 encoded string of the configuration data.")
	cmd.Flags().StringVar(&p.configGCPSecret, "config-gcp-secret", p.configGCPSecret, "The resource ID of secret that contains Piped config and be stored in GCP SecretManager.")
	cmd.Flags().StringVar(&p.configAWSSecret, "config-aws-secret", p.configAWSSecret, "The ARN of secret that contains Piped config and be stored in AWS Secrets Manager.")

	cmd.Flags().BoolVar(&p.insecure, "insecure", p.insecure, "Whether disabling transport security while connecting to control-plane.")
	cmd.Flags().StringVar(&p.certFile, "cert-file", p.certFile, "The path to the TLS certificate file.")
	cmd.Flags().IntVar(&p.adminPort, "admin-port", p.adminPort, "The port number used to run a HTTP server for admin tasks such as metrics, healthz.")
	cmd.Flags().IntVar(&p.pluginServicePort, "plugin-service-port", p.pluginServicePort, "The port number used to run a gRPC server for plugin services.")

	cmd.Flags().StringVar(&p.toolsDir, "tools-dir", p.toolsDir, "The path to directory where to install needed tools such as kubectl, helm, kustomize.")
	cmd.Flags().BoolVar(&p.addLoginUserToPasswd, "add-login-user-to-passwd", p.addLoginUserToPasswd, "Whether to add login user to $HOME/passwd. This is typically for applications running as a random user ID.")
	cmd.Flags().DurationVar(&p.gracePeriod, "grace-period", p.gracePeriod, "How long to wait for graceful shutdown.")

	cmd.Flags().StringVar(&p.launcherVersion, "launcher-version", p.launcherVersion, "The version of launcher which initialized this Piped.")

	return cmd
}

func (p *piped) run(ctx context.Context, input cli.Input) (runErr error) {
	// Make a cancellable context.
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	group, ctx := errgroup.WithContext(ctx)
	if p.addLoginUserToPasswd {
		if err := p.insertLoginUserToPasswd(ctx); err != nil {
			return fmt.Errorf("failed to insert logged-in user to passwd: %w", err)
		}
	}

	// Load piped configuration from the specified source.
	cfg, err := p.loadConfig(ctx)
	if err != nil {
		input.Logger.Error("failed to load piped configuration", zap.Error(err))
		return err
	}

	// Register all metrics.
	registry := registerMetrics(cfg.PipedID, cfg.ProjectID, p.launcherVersion)

	// // Configure SSH config if needed.
	// if cfg.Git.ShouldConfigureSSHConfig() {
	// 	if err := git.AddSSHConfig(cfg.Git); err != nil {
	// 		input.Logger.Error("failed to configure ssh-config", zap.Error(err))
	// 		return err
	// 	}
	// 	input.Logger.Info("successfully configured ssh-config")
	// }

	pipedKey, err := cfg.LoadPipedKey()
	if err != nil {
		input.Logger.Error("failed to load piped key", zap.Error(err))
		return err
	}

	// Make gRPC client and connect to the Control Plane API.
	apiClient, err := p.createAPIClient(ctx, cfg.APIAddress, cfg.ProjectID, cfg.PipedID, pipedKey, input.Logger)
	if err != nil {
		input.Logger.Error("failed to create gRPC client to control plane", zap.Error(err))
		return err
	}

	// Setup the tracer provider.
	// We don't set the global tracer provider because 3rd-party library may use the global one.
	tracerProvider, err := p.createTracerProvider(ctx, cfg.APIAddress, cfg.ProjectID, cfg.PipedID, pipedKey)
	if err != nil {
		input.Logger.Error("failed to create tracer provider", zap.Error(err))
		return err
	}

	// Send the newest piped meta to the control-plane.
	if err := p.sendPipedMeta(ctx, apiClient, cfg, input.Logger); err != nil {
		input.Logger.Error("failed to report piped meta to control-plane", zap.Error(err))
		return err
	}

	// Initialize notifier and add piped events.
	notifier, err := notifier.NewNotifier(cfg, input.Logger)
	if err != nil {
		input.Logger.Error("failed to initialize notifier", zap.Error(err))
		return err
	}
	group.Go(func() error {
		return notifier.Run(ctx)
	})

	// Start running admin server.
	{
		var (
			ver   = []byte(version.Get().Version)
			admin = admin.NewAdmin(p.adminPort, p.gracePeriod, input.Logger)
		)

		admin.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
			w.Write(ver)
		})
		admin.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("ok"))
		})
		admin.Handle("/metrics", input.PrometheusMetricsHandlerFor(registry))
		admin.HandleFunc("/debug/pprof/", pprof.Index)
		admin.HandleFunc("/debug/pprof/profile", pprof.Profile)
		admin.HandleFunc("/debug/pprof/trace", pprof.Trace)

		group.Go(func() error {
			return admin.Run(ctx)
		})
	}

	// Start running stats reporter.
	{
		url := fmt.Sprintf("http://localhost:%d/metrics", p.adminPort)
		r := statsreporter.NewReporter(url, apiClient, input.Logger)
		group.Go(func() error {
			return r.Run(ctx)
		})
	}

	// Initialize git client.
	gitOptions := []git.Option{
		git.WithUserName(cfg.Git.Username),
		git.WithEmail(cfg.Git.Email),
		git.WithLogger(input.Logger),
	}
	gitClient, err := git.NewClient(gitOptions...)
	if err != nil {
		input.Logger.Error("failed to initialize git client", zap.Error(err))
		return err
	}
	defer func() {
		if err := gitClient.Clean(); err != nil {
			input.Logger.Error("had an error while cleaning gitClient", zap.Error(err))
			return
		}
		input.Logger.Info("successfully cleaned gitClient")
	}()

	// Start running application store.
	var applicationLister applicationstore.Lister
	{
		store := applicationstore.NewStore(apiClient, p.gracePeriod, input.Logger)
		group.Go(func() error {
			return store.Run(ctx)
		})
		applicationLister = store.Lister()
	}

	// Start running deployment store.
	var deploymentLister deploymentstore.Lister
	{
		store := deploymentstore.NewStore(apiClient, p.gracePeriod, input.Logger)
		group.Go(func() error {
			return store.Run(ctx)
		})
		deploymentLister = store.Lister()
	}

	// Start running command store.
	var commandLister commandstore.Lister
	{
		store := commandstore.NewStore(apiClient, p.gracePeriod, input.Logger)
		group.Go(func() error {
			return store.Run(ctx)
		})
		commandLister = store.Lister()
	}

	// Start running event store.
	var eventLister eventstore.Lister
	{
		store := eventstore.NewStore(apiClient, p.gracePeriod, input.Logger)
		group.Go(func() error {
			return store.Run(ctx)
		})
		eventLister = store.Lister()
	}

	metadataStoreRegistry := metadatastore.NewMetadataStoreRegistry(apiClient)
	// Start running plugin service server.
	{
		var (
			service, err = grpcapi.NewPluginAPI(cfg, apiClient, p.toolsDir, input.Logger, metadataStoreRegistry, notifier)
			opts         = []rpc.Option{
				rpc.WithPort(p.pluginServicePort),
				rpc.WithGracePeriod(p.gracePeriod),
				rpc.WithLogger(input.Logger),
				rpc.WithLogUnaryInterceptor(input.Logger),
				rpc.WithRequestValidationUnaryInterceptor(),
			}
		)
		if err != nil {
			input.Logger.Error("failed to create plugin service", zap.Error(err))
			return err
		}
		// TODO: Ensure piped <-> plugin communication is secure.
		server := rpc.NewServer(service, opts...)
		group.Go(func() error {
			return server.Run(ctx)
		})
	}

	// Start plugins that registered in the configuration.
	{
		// Start all plugins and keep their commands to stop them later.
		plugins, err := p.runPlugins(ctx, cfg.Plugins, input.Logger)
		if err != nil {
			input.Logger.Error("failed to run plugins", zap.Error(err))
			return err
		}

		group.Go(func() error {
			<-ctx.Done()
			wg := &sync.WaitGroup{}
			for _, plg := range plugins {
				wg.Add(1)
				go func() {
					defer wg.Done()
					if err := plg.GracefulStop(p.gracePeriod); err != nil {
						input.Logger.Error("failed to stop plugin", zap.Error(err))
					}
				}()
			}
			wg.Wait()
			return nil
		})
	}

	// Make grpc clients to connect to plugins.
	plugins := make([]plugin.Plugin, 0, len(cfg.Plugins))
	options := []rpcclient.DialOption{
		rpcclient.WithBlock(),
		rpcclient.WithInsecure(),
	}
	for _, plg := range cfg.Plugins {
		cli, err := pluginapi.NewClient(ctx, net.JoinHostPort("localhost", strconv.Itoa(plg.Port)), options...)
		if err != nil {
			input.Logger.Error("failed to create client to connect plugin", zap.String("plugin", plg.Name), zap.Error(err))
			return err
		}

		plugins = append(plugins, plugin.Plugin{
			Name: plg.Name,
			Cli:  cli,
		})
	}

	pluginRegistry, err := plugin.NewPluginRegistry(ctx, plugins)
	if err != nil {
		input.Logger.Error("failed to create plugin registry", zap.Error(err))
		return err
	}

	// Initialize secret decrypter.
	decrypter, err := p.initializeSecretDecrypter(cfg)
	if err != nil {
		input.Logger.Error("failed to initialize secret decrypter", zap.Error(err))
		return err
	}

	// Start running application live state reporter.
	// Currently, this feature is disabled beucause many errors are showed up if the app.pipecd.yaml is not migrated.
	// {
	// 	r, err := livestatereporter.NewReporter(applicationLister, apiClient, gitClient, pluginRegistry, cfg, decrypter, input.Logger)
	// 	if err != nil {
	// 		input.Logger.Error("failed to create live state reporter", zap.Error(err))
	// 	}
	// 	group.Go(func() error {
	// 		return r.Run(ctx)
	// 	})
	// }

	// Start running deployment controller.
	{
		c := controller.NewController(
			apiClient,
			gitClient,
			pluginRegistry,
			deploymentLister,
			commandLister,
			notifier,
			decrypter,
			*metadataStoreRegistry,
			p.gracePeriod,
			input.Logger,
			tracerProvider,
		)

		group.Go(func() error {
			return c.Run(ctx)
		})
	}

	// Start running deployment trigger.
	{
		tr, err := trigger.NewTrigger(
			apiClient,
			gitClient,
			applicationLister,
			commandLister,
			notifier,
			cfg,
			p.gracePeriod,
			input.Logger,
		)
		if err != nil {
			input.Logger.Error("failed to initialize trigger", zap.Error(err))
			return err
		}

		group.Go(func() error {
			return tr.Run(ctx)
		})
	}

	// Start running event watcher.
	{
		w := eventwatcher.NewWatcher(
			cfg,
			eventLister,
			gitClient,
			apiClient,
			input.Logger,
		)
		group.Go(func() error {
			return w.Run(ctx)
		})
	}

	// Start running planpreview handler.
	{
		// TODO: Implement planpreview controller.
	}

	// Start running app-config-reporter.
	{
		r := appconfigreporter.NewReporter(
			apiClient,
			gitClient,
			applicationLister,
			cfg,
			p.gracePeriod,
			input.Logger,
		)

		group.Go(func() error {
			return r.Run(ctx)
		})
	}

	// Check for stop command.
	{
		group.Go(func() error {
			input.Logger.Info("start running piped stop checker")
			ticker := time.NewTicker(commandCheckPeriod)
			for {
				select {
				case <-ticker.C:
					shouldStop, err := stopCommandHandler(ctx, commandLister, input.Logger)
					// Don't return an error to continue this goroutine execution.
					if err != nil {
						input.Logger.Error("failed to check/handle piped stop command", zap.Error(err))
					}
					if shouldStop {
						input.Logger.Info("stop piped due to restart piped requested")
						cancel()
					}
				case <-ctx.Done():
					input.Logger.Info("piped stop checker has been stopped")
					return nil
				}
			}
		})
	}

	// Wait until all piped components have finished.
	// A terminating signal or a finish of any components
	// could trigger the finish of piped.
	// This ensures that all components are good or no one.
	if err := group.Wait(); err != nil {
		input.Logger.Error("failed while running", zap.Error(err))
		return err
	}
	return nil
}

// createAPIClient makes a gRPC client to connect to the API.
func (p *piped) createAPIClient(ctx context.Context, address, projectID, pipedID string, pipedKey []byte, logger *zap.Logger) (pipedservice.Client, error) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	var (
		token   = rpcauth.MakePipedToken(projectID, pipedID, string(pipedKey))
		creds   = rpcclient.NewPerRPCCredentials(token, rpcauth.PipedTokenCredentials, !p.insecure)
		options = []rpcclient.DialOption{
			rpcclient.WithBlock(),
			rpcclient.WithPerRPCCredentials(creds),
			rpcclient.WithMaxRecvMsgSize(p.maxRecvMsgSize),
		}
	)

	if !p.insecure {
		if p.certFile != "" {
			options = append(options, rpcclient.WithTLS(p.certFile))
		} else {
			config := &tls.Config{}
			options = append(options, rpcclient.WithTransportCredentials(credentials.NewTLS(config)))
		}
	} else {
		options = append(options, rpcclient.WithInsecure())
	}

	client, err := pipedservice.NewClient(ctx, address, options...)
	if err != nil {
		logger.Error("failed to create api client", zap.Error(err))
		return nil, err
	}
	return client, nil
}

// createTracerProvider makes a OpenTelemetry Trace's TracerProvider.
func (p *piped) createTracerProvider(ctx context.Context, address, projectID, pipedID string, pipedKey []byte) (trace.TracerProvider, error) {
	options := []otlptracegrpc.Option{
		otlptracegrpc.WithEndpoint(address),
		otlptracegrpc.WithHeaders(map[string]string{
			"authorization": "Bearer " + rpcauth.MakePipedToken(projectID, pipedID, string(pipedKey)),
		}),
	}

	if !p.insecure {
		if p.certFile != "" {
			creds, err := credentials.NewClientTLSFromFile(p.certFile, "")
			if err != nil {
				return nil, fmt.Errorf("failed to load client TLS credentials: %w", err)
			}
			options = append(options, otlptracegrpc.WithTLSCredentials(creds))
		} else {
			config := &tls.Config{}
			options = append(options, otlptracegrpc.WithTLSCredentials(credentials.NewTLS(config)))
		}
	} else {
		options = append(options, otlptracegrpc.WithInsecure())
	}

	otlpTraceExporter, err := otlptracegrpc.New(ctx, options...)
	if err != nil {
		return nil, err
	}

	otlpResource, err := resource.New(ctx, resource.WithAttributes(
		// Set common attributes for all spans.
		attribute.String("service.name", "piped"),
		attribute.String("service.version", version.Get().Version),
		attribute.String("service.namespace", projectID),
		attribute.String("service.instance.id", pipedID),

		// Set the project and piped IDs as attributes.
		attribute.String("project-id", projectID),
		attribute.String("piped-id", pipedID),
	))
	if err != nil {
		return nil, err
	}

	otlpResource, err = resource.Merge(resource.Default(), otlpResource) // the later one has higher priority
	if err != nil {
		return nil, err
	}

	return sdktrace.NewTracerProvider(
		sdktrace.WithResource(otlpResource),
		sdktrace.WithBatcher(otlpTraceExporter),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	), nil
}

// loadConfig reads the Piped configuration data from the specified source.
func (p *piped) loadConfig(ctx context.Context) (*config.PipedSpec, error) {
	// HACK: When the version of cobra is updated to >=v1.8.0, this should be replaced with https://pkg.go.dev/github.com/spf13/cobra#Command.MarkFlagsMutuallyExclusive.
	if err := p.hasTooManyConfigFlags(); err != nil {
		return nil, err
	}

	extract := func(cfg *config.Config[*config.PipedSpec, config.PipedSpec]) (*config.PipedSpec, error) {
		if cfg.Kind != config.KindPiped {
			return nil, fmt.Errorf("wrong configuration kind for piped: %v", cfg.Kind)
		}
		return cfg.Spec, nil
	}

	if p.configFile != "" {
		cfg, err := config.LoadFromYAML[*config.PipedSpec](p.configFile)
		if err != nil {
			return nil, err
		}
		return extract(cfg)
	}

	if p.configData != "" {
		data, err := base64.StdEncoding.DecodeString(p.configData)
		if err != nil {
			return nil, fmt.Errorf("the given config-data isn't base64 encoded: %w", err)
		}

		cfg, err := config.DecodeYAML[*config.PipedSpec](data)
		if err != nil {
			return nil, err
		}
		return extract(cfg)
	}

	if p.configGCPSecret != "" {
		data, err := p.getConfigDataFromSecretManager(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to load config from SecretManager (%w)", err)
		}
		cfg, err := config.DecodeYAML[*config.PipedSpec](data)
		if err != nil {
			return nil, err
		}
		return extract(cfg)
	}

	if p.configAWSSecret != "" {
		data, err := p.getConfigDataFromAWSSecretsManager(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to load config from AWS Secrets Manager (%w)", err)
		}
		cfg, err := config.DecodeYAML[*config.PipedSpec](data)
		if err != nil {
			return nil, err
		}
		return extract(cfg)
	}

	return nil, fmt.Errorf("one of config-file, config-gcp-secret or config-aws-secret must be set")
}

func (p *piped) runPlugins(ctx context.Context, pluginsCfg []config.PipedPlugin, logger *zap.Logger) ([]*lifecycle.Command, error) {
	plugins := make([]*lifecycle.Command, 0, len(pluginsCfg))
	for _, pCfg := range pluginsCfg {
		// Download plugin binary to piped's pluginsDir.
		pPath, err := lifecycle.DownloadBinary(pCfg.URL, p.pluginsDir, pCfg.Name, logger)
		if err != nil {
			return nil, fmt.Errorf("failed to download plugin %s: %w", pCfg.Name, err)
		}

		// Build plugin's args.
		args := make([]string, 0, 4)
		args = append(args, "start", "--piped-plugin-service", net.JoinHostPort("localhost", strconv.Itoa(p.pluginServicePort)))
		b, err := json.Marshal(pCfg)
		if err != nil {
			return nil, fmt.Errorf("failed to prepare plugin %s config: %w", pCfg.Name, err)
		}
		args = append(args, "--config", string(b))

		// Run the plugin binary.
		cmd, err := lifecycle.RunBinary(ctx, pPath, args)
		if err != nil {
			return nil, fmt.Errorf("failed to run plugin %s: %w", pCfg.Name, err)
		}

		plugins = append(plugins, cmd)
	}
	return plugins, nil
}

func (p *piped) initializeSecretDecrypter(cfg *config.PipedSpec) (crypto.Decrypter, error) {
	sm := cfg.SecretManagement
	if sm == nil {
		return nil, nil
	}

	switch sm.Type {
	case model.SecretManagementTypeNone:
		return nil, nil

	case model.SecretManagementTypeKeyPair:
		key, err := sm.KeyPair.LoadPrivateKey()
		if err != nil {
			return nil, err
		}
		decrypter, err := crypto.NewHybridDecrypter(key)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize decrypter (%w)", err)
		}
		return decrypter, nil

	case model.SecretManagementTypeGCPKMS:
		return nil, fmt.Errorf("type %q is not implemented yet", sm.Type.String())

	case model.SecretManagementTypeAWSKMS:
		return nil, fmt.Errorf("type %q is not implemented yet", sm.Type.String())

	default:
		return nil, fmt.Errorf("unsupported secret management type: %s", sm.Type.String())
	}
}

func (p *piped) sendPipedMeta(ctx context.Context, client pipedservice.Client, cfg *config.PipedSpec, logger *zap.Logger) error {
	repos := make([]*model.ApplicationGitRepository, 0, len(cfg.Repositories))
	for _, r := range cfg.Repositories {
		repos = append(repos, &model.ApplicationGitRepository{
			Id:     r.RepoID,
			Remote: r.Remote,
			Branch: r.Branch,
		})
	}

	cloneCfg, err := cfg.Clone()
	if err != nil {
		return err
	}

	cloneCfg.Mask()
	maskedCfg, err := yaml.Marshal(cloneCfg)
	if err != nil {
		return err
	}

	req := &pipedservice.ReportPipedMetaRequest{
		Version:      version.Get().Version,
		Config:       string(maskedCfg),
		Repositories: repos,
	}

	// Configure secret management.
	if sm := cfg.SecretManagement; sm != nil && sm.Type == model.SecretManagementTypeKeyPair {
		publicKey, err := sm.KeyPair.LoadPublicKey()
		if err != nil {
			return fmt.Errorf("failed to read public key for secret management (%w)", err)
		}
		req.SecretEncryption = &model.Piped_SecretEncryption{
			Type:      sm.Type.String(),
			PublicKey: string(publicKey),
		}
	}
	if req.SecretEncryption == nil {
		req.SecretEncryption = &model.Piped_SecretEncryption{
			Type: model.SecretManagementTypeNone.String(),
		}
	}

	retry := pipedservice.NewRetry(5)
	_, err = retry.Do(ctx, func() (interface{}, error) {
		if res, err := client.ReportPipedMeta(ctx, req); err == nil {
			cfg.Name = res.Name
			if cfg.WebAddress == "" {
				cfg.WebAddress = res.WebBaseUrl
			}
			return nil, nil
		}
		logger.Warn("failed to report piped meta to control-plane, wait to the next retry",
			zap.Int("calls", retry.Calls()),
			zap.Error(err),
		)
		return nil, err
	})

	return err
}

// insertLoginUserToPasswd adds the logged-in user to /etc/passwd.
// It requires nss_wrapper (https://cwrap.org/nss_wrapper.html)
// to get the operation done.
//
// This is a workaround to deal with OpenShift less than 4.2
// See more: https://github.com/pipe-cd/pipecd/issues/1905
func (p *piped) insertLoginUserToPasswd(ctx context.Context) error {
	var stdout, stderr bytes.Buffer

	// Use the id command so that it gets proper ids even in pure Go.
	cmd := exec.CommandContext(ctx, "id", "-u")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to get uid: %s", &stderr)
	}
	uid := strings.TrimSpace(stdout.String())

	stdout.Reset()
	stderr.Reset()

	cmd = exec.CommandContext(ctx, "id", "-g")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to get gid: %s", &stderr)
	}
	gid := strings.TrimSpace(stdout.String())

	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to detect the current user's home directory: %w", err)
	}

	// echo "default:x:${USER_ID}:${GROUP_ID}:Dynamically created user:${HOME}:/sbin/nologin" >> "$HOME/passwd"
	entry := fmt.Sprintf("\ndefault:x:%s:%s:Dynamically created user:%s:/sbin/nologin", uid, gid, home)
	nssPasswdPath := filepath.Join(home, "passwd")
	f, err := os.OpenFile(nssPasswdPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
	if err != nil {
		return fmt.Errorf("failed to open file %q: %w", nssPasswdPath, err)
	}
	defer f.Close()
	if _, err := f.WriteString(entry); err != nil {
		return fmt.Errorf("failed to append entry to %q: %w", nssPasswdPath, err)
	}

	return nil
}

func (p *piped) getConfigDataFromSecretManager(ctx context.Context) ([]byte, error) {
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: p.configGCPSecret,
	}

	resp, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.Payload.Data, nil
}

func (p *piped) getConfigDataFromAWSSecretsManager(ctx context.Context) ([]byte, error) {
	cfg, err := awsconfig.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}
	client := awssecretsmanager.NewFromConfig(cfg)

	in := &awssecretsmanager.GetSecretValueInput{
		SecretId: &p.configAWSSecret,
	}

	result, err := client.GetSecretValue(ctx, in)
	if err != nil {
		return nil, err
	}

	decoded, err := base64.StdEncoding.DecodeString(*result.SecretString)
	if err != nil {
		return nil, err
	}

	return decoded, nil
}

func registerMetrics(pipedID, projectID, launcherVersion string) *prometheus.Registry {
	r := prometheus.NewRegistry()
	wrapped := prometheus.WrapRegistererWith(
		map[string]string{
			"pipecd_component": "piped",
			"piped":            pipedID,
			"piped_version":    version.Get().Version,
			"launcher_version": launcherVersion,
			"project":          projectID,
		},
		r,
	)
	wrapped.Register(collectors.NewGoCollector())
	wrapped.Register(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))

	controllermetrics.Register(wrapped)

	return r
}

func stopCommandHandler(ctx context.Context, cmdLister commandstore.Lister, logger *zap.Logger) (bool, error) {
	logger.Debug("fetch unhandled piped commands")

	commands := cmdLister.ListPipedCommands()
	if len(commands) == 0 {
		return false, nil
	}

	stopCmds := make([]model.ReportableCommand, 0, len(commands))
	for _, command := range commands {
		if command.IsRestartPipedCmd() {
			stopCmds = append(stopCmds, command)
		}
	}

	if len(stopCmds) == 0 {
		return false, nil
	}

	for _, command := range stopCmds {
		if err := command.Report(ctx, model.CommandStatus_COMMAND_SUCCEEDED, nil, []byte(command.Id)); err != nil {
			return false, fmt.Errorf("failed to report command %s: %w", command.Id, err)
		}
	}

	return true, nil
}

func (p *piped) hasTooManyConfigFlags() error {
	cnt := 0
	for _, v := range []string{p.configFile, p.configGCPSecret, p.configAWSSecret} {
		if v != "" {
			cnt++
		}
	}
	if cnt > 1 {
		return fmt.Errorf("only one of config-file, config-gcp-secret or config-aws-secret could be set")
	}
	return nil
}
