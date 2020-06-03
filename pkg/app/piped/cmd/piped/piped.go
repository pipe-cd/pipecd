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

package piped

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/kapetaniosci/pipe/pkg/admin"
	"github.com/kapetaniosci/pipe/pkg/app/api/service/pipedservice"
	"github.com/kapetaniosci/pipe/pkg/app/api/service/pipedservice/pipedclientfake"
	"github.com/kapetaniosci/pipe/pkg/app/piped/apistore/applicationstore"
	"github.com/kapetaniosci/pipe/pkg/app/piped/apistore/commandstore"
	"github.com/kapetaniosci/pipe/pkg/app/piped/apistore/deploymentstore"
	"github.com/kapetaniosci/pipe/pkg/app/piped/controller"
	"github.com/kapetaniosci/pipe/pkg/app/piped/livestatereporter"
	"github.com/kapetaniosci/pipe/pkg/app/piped/livestatestore"
	"github.com/kapetaniosci/pipe/pkg/app/piped/toolregistry"
	"github.com/kapetaniosci/pipe/pkg/app/piped/trigger"
	"github.com/kapetaniosci/pipe/pkg/cache/memorycache"
	"github.com/kapetaniosci/pipe/pkg/cli"
	"github.com/kapetaniosci/pipe/pkg/config"
	"github.com/kapetaniosci/pipe/pkg/git"
	"github.com/kapetaniosci/pipe/pkg/model"
	"github.com/kapetaniosci/pipe/pkg/rpc/rpcauth"
	"github.com/kapetaniosci/pipe/pkg/rpc/rpcclient"
	"github.com/kapetaniosci/pipe/pkg/version"

	// Import to preload all built-in executors to the default registry.
	_ "github.com/kapetaniosci/pipe/pkg/app/piped/executor/registry"
	// Import to preload all planners to the default registry.
	_ "github.com/kapetaniosci/pipe/pkg/app/piped/planner/registry"
)

type piped struct {
	projectID    string
	pipedID      string
	pipedKeyFile string
	configFile   string

	tls                 bool
	certFile            string
	controlPlaneAddress string
	adminPort           int

	binDir                               string
	enableDefaultKubernetesCloudProvider bool
	useFakeAPIClient                     bool
	gracePeriod                          time.Duration
}

func NewCommand() *cobra.Command {
	p := &piped{
		controlPlaneAddress: "pipecd:443",
		adminPort:           9085,
		binDir:              "/usr/local/piped",
		gracePeriod:         30 * time.Second,
	}
	cmd := &cobra.Command{
		Use:   "piped",
		Short: "Start running piped.",
		RunE:  cli.WithContext(p.run),
	}

	cmd.Flags().StringVar(&p.projectID, "project-id", p.projectID, "The identifier of the project which this piped belongs to.")
	cmd.Flags().StringVar(&p.pipedID, "piped-id", p.pipedID, "The unique identifier generated for this piped.")
	cmd.Flags().StringVar(&p.pipedKeyFile, "piped-key-file", p.pipedKeyFile, "The path to the key generated for this piped.")
	cmd.Flags().StringVar(&p.configFile, "config-file", p.configFile, "The path to the configuration file.")

	cmd.Flags().BoolVar(&p.tls, "tls", p.tls, "Whether running the gRPC server with TLS or not.")
	cmd.Flags().StringVar(&p.certFile, "cert-file", p.certFile, "The path to the TLS certificate file.")
	cmd.Flags().StringVar(&p.controlPlaneAddress, "control-plane-address", p.controlPlaneAddress, "The address used to connect to control plane.")
	cmd.Flags().IntVar(&p.adminPort, "admin-port", p.adminPort, "The port number used to run a HTTP server for admin tasks such as metrics, healthz.")

	cmd.Flags().StringVar(&p.binDir, "bin-dir", p.binDir, "The path to directory where to install needed tools such as kubectl, helm, kustomize.")
	cmd.Flags().BoolVar(&p.useFakeAPIClient, "use-fake-api-client", p.useFakeAPIClient, "Whether the fake api client should be used instead of the real one or not.")
	cmd.Flags().BoolVar(&p.enableDefaultKubernetesCloudProvider, "enable-default-kubernetes-cloud-provider", p.enableDefaultKubernetesCloudProvider, "Whether the default kubernetes provider is enabled or not.")
	cmd.Flags().DurationVar(&p.gracePeriod, "grace-period", p.gracePeriod, "How long to wait for graceful shutdown.")

	cmd.MarkFlagRequired("project-id")
	cmd.MarkFlagRequired("piped-id")
	cmd.MarkFlagRequired("piped-key-file")
	cmd.MarkFlagRequired("config-file")

	return cmd
}

func (p *piped) run(ctx context.Context, t cli.Telemetry) error {
	group, ctx := errgroup.WithContext(ctx)

	// Load piped configuration from specified file.
	cfg, err := p.loadConfig()
	if err != nil {
		t.Logger.Error("failed to load piped configuration", zap.Error(err))
		return err
	}

	// Configure SSH config if needed.
	if cfg.Git.ShouldConfigureSSHConfig() {
		if err := git.AddSSHConfig(cfg.Git); err != nil {
			t.Logger.Error("failed to configure ssh-config", zap.Error(err))
			return err
		}
		t.Logger.Info("successfully configured ssh-config")
	}

	// Initialize default tool registry.
	if err := toolregistry.InitDefaultRegistry(p.binDir, t.Logger); err != nil {
		t.Logger.Error("failed to initialize default tool registry", zap.Error(err))
		return err
	}

	// Make gRPC client and connect to the API.
	apiClient, err := p.createAPIClient(ctx, t.Logger)
	if err != nil {
		t.Logger.Error("failed to create gRPC client to control plane", zap.Error(err))
		return err
	}

	// Send the newest piped meta to the control-plane.
	if err := p.sendPipedMeta(ctx, apiClient, cfg, t.Logger); err != nil {
		t.Logger.Error("failed to report piped meta to control-plane", zap.Error(err))
		return err
	}

	// Start running admin server.
	{
		admin := admin.NewAdmin(p.adminPort, p.gracePeriod, t.Logger)
		if exporter, ok := t.PrometheusMetricsExporter(); ok {
			admin.Handle("/metrics", exporter)
		}
		admin.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("ok"))
		})
		group.Go(func() error {
			return admin.Run(ctx)
		})
	}

	// Initialize git client.
	gitClient, err := git.NewClient(cfg.Git.Username, cfg.Git.Email, t.Logger)
	if err != nil {
		t.Logger.Error("failed to initialize git client", zap.Error(err))
		return err
	}
	defer func() {
		if err := gitClient.Clean(); err != nil {
			t.Logger.Error("had an error while cleaning gitClient", zap.Error(err))
		} else {
			t.Logger.Info("successfully cleaned gitClient")
		}
	}()

	// Start running application store.
	var applicationLister applicationstore.Lister
	{
		store := applicationstore.NewStore(apiClient, p.gracePeriod, t.Logger)
		group.Go(func() error {
			return store.Run(ctx)
		})
		applicationLister = store.Lister()
	}

	// Start running deployment store.
	var deploymentLister deploymentstore.Lister
	{
		store := deploymentstore.NewStore(apiClient, p.gracePeriod, t.Logger)
		group.Go(func() error {
			return store.Run(ctx)
		})
		deploymentLister = store.Lister()
	}

	// Start running command store.
	var commandLister commandstore.Lister
	{
		store := commandstore.NewStore(apiClient, p.gracePeriod, t.Logger)
		group.Go(func() error {
			return store.Run(ctx)
		})
		commandLister = store.Lister()
	}

	// Create memory caches.
	appManifestsCache := memorycache.NewTTLCache(ctx, time.Hour, time.Minute)

	// Start running application live state store.
	{
		s := livestatestore.NewStore(cfg, applicationLister, p.gracePeriod, t.Logger)
		group.Go(func() error {
			return s.Run(ctx)
		})
	}

	// Start running application live state reporter.
	{
		r := livestatereporter.NewReporter(p.gracePeriod)
		group.Go(func() error {
			return r.Run(ctx)
		})
	}

	// Start running deployment controller.
	{
		c := controller.NewController(
			apiClient,
			gitClient,
			deploymentLister,
			commandLister,
			applicationLister,
			cfg,
			appManifestsCache,
			p.gracePeriod,
			t.Logger,
		)

		group.Go(func() error {
			return c.Run(ctx)
		})
	}

	// Start running deployment trigger.
	{
		t := trigger.NewTrigger(apiClient, gitClient, applicationLister, commandLister, cfg, p.gracePeriod, t.Logger)
		group.Go(func() error {
			return t.Run(ctx)
		})
	}

	// Wait until all piped components have finished.
	// A terminating signal or a finish of any components
	// could trigger the finish of piped.
	// This ensures that all components are good or no one.
	if err := group.Wait(); err != nil {
		t.Logger.Error("failed while running", zap.Error(err))
		return err
	}
	return nil
}

// createAPIClient makes a gRPC client to connect to the API.
func (p *piped) createAPIClient(ctx context.Context, logger *zap.Logger) (pipedservice.Client, error) {
	if p.useFakeAPIClient {
		return pipedclientfake.NewClient(logger), nil
	}
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	pipedKey, err := ioutil.ReadFile(p.pipedKeyFile)
	if err != nil {
		logger.Error("failed to read piped key file", zap.Error(err))
		return nil, err
	}

	var (
		token   = rpcauth.MakePipedToken(p.projectID, p.pipedID, string(pipedKey))
		tls     = p.certFile != ""
		creds   = rpcclient.NewPerRPCCredentials(token, rpcauth.PipedTokenCredentials, tls)
		options = []rpcclient.DialOption{
			rpcclient.WithBlock(),
			rpcclient.WithStatsHandler(),
			rpcclient.WithPerRPCCredentials(creds),
		}
	)
	if tls {
		options = append(options, rpcclient.WithTLS(p.certFile))
	} else {
		options = append(options, rpcclient.WithInsecure())
	}

	client, err := pipedservice.NewClient(ctx, p.controlPlaneAddress, options...)
	if err != nil {
		logger.Error("failed to create api client", zap.Error(err))
		return nil, err
	}
	return client, nil
}

// loadConfig reads the Piped configuration data from the specified file.
func (p *piped) loadConfig() (*config.PipedSpec, error) {
	cfg, err := config.LoadFromYAML(p.configFile)
	if err != nil {
		return nil, err
	}
	if cfg.Kind != config.KindPiped {
		return nil, fmt.Errorf("wrong configuration kind for piped: %v", cfg.Kind)
	}
	if p.enableDefaultKubernetesCloudProvider {
		cfg.PipedSpec.EnableDefaultKubernetesCloudProvider()
	}
	return cfg.PipedSpec, nil
}

func (p *piped) sendPipedMeta(ctx context.Context, client pipedservice.Client, cfg *config.PipedSpec, logger *zap.Logger) error {
	var (
		req = &pipedservice.ReportPipedMetaRequest{
			Version:        version.Get().Version,
			CloudProviders: make([]*model.Piped_CloudProvider, 0, len(cfg.CloudProviders)),
		}
		retry = pipedservice.NewRetry(10)
		err   error
	)

	for _, cp := range cfg.CloudProviders {
		req.CloudProviders = append(req.CloudProviders, &model.Piped_CloudProvider{
			Name: cp.Name,
			Type: cp.Type.String(),
		})
	}

	for retry.WaitNext(ctx) {
		if _, err = client.ReportPipedMeta(ctx, req); err == nil {
			return nil
		}
		logger.Warn("failed to report piped meta to control-plane, wait to the next retry", zap.Int("calls", retry.Calls()))
	}

	return err
}
