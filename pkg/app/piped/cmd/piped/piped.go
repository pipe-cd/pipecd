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
	"k8s.io/client-go/tools/clientcmd"

	// The following line to load the gcp plugin (only required to authenticate against GKE clusters).
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	"github.com/kapetaniosci/pipe/pkg/admin"
	"github.com/kapetaniosci/pipe/pkg/app/api/service/pipedservice"
	"github.com/kapetaniosci/pipe/pkg/app/api/service/pipedservice/pipedclientfake"
	"github.com/kapetaniosci/pipe/pkg/app/piped/applicationstore"
	"github.com/kapetaniosci/pipe/pkg/app/piped/appstatereporter"
	"github.com/kapetaniosci/pipe/pkg/app/piped/appstatestore"
	"github.com/kapetaniosci/pipe/pkg/app/piped/commandstore"
	"github.com/kapetaniosci/pipe/pkg/app/piped/deploymentcontroller"
	"github.com/kapetaniosci/pipe/pkg/app/piped/deploymenttrigger"
	"github.com/kapetaniosci/pipe/pkg/cli"
	"github.com/kapetaniosci/pipe/pkg/config"
	"github.com/kapetaniosci/pipe/pkg/git"
	"github.com/kapetaniosci/pipe/pkg/rpc/rpcauth"
	"github.com/kapetaniosci/pipe/pkg/rpc/rpcclient"
)

type piped struct {
	projectID    string
	pipedID      string
	pipedKeyFile string
	configFile   string
	tls          bool
	certFile     string
	apiAddress   string
	adminPort    int

	gitUserName string
	gitEmail    string

	kubeconfig string
	masterURL  string
	namespace  string

	useFakeAPIClient bool
	gracePeriod      time.Duration
}

func NewCommand() *cobra.Command {
	p := &piped{
		apiAddress:  "pipecd-api:9091",
		adminPort:   9085,
		gracePeriod: 30 * time.Second,
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
	cmd.Flags().StringVar(&p.apiAddress, "api-address", p.apiAddress, "The address used to connect to API server.")
	cmd.Flags().IntVar(&p.adminPort, "admin-port", p.adminPort, "The port number used to run a HTTP server for admin tasks such as metrics, healthz.")

	cmd.Flags().StringVar(&p.gitUserName, "git-username", p.gitUserName, "The username used to be configured fot git commands.")
	cmd.Flags().StringVar(&p.gitEmail, "git-email", p.gitEmail, "The email used to be configured fot git commands.")

	cmd.Flags().StringVar(&p.kubeconfig, "kube-config", p.kubeconfig, "Path to a kubeconfig. Only required if out-of-cluster.")
	cmd.Flags().StringVar(&p.masterURL, "master", p.masterURL, "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
	cmd.Flags().StringVar(&p.namespace, "namespace", p.namespace, "The namespace where this piped is running.")

	cmd.Flags().BoolVar(&p.useFakeAPIClient, "use-fake-api-client", p.useFakeAPIClient, "Whether the fake api client should be used instead of the real one or not.")
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
	pipedConfig, err := p.loadConfig()
	if err != nil {
		t.Logger.Error("failed to load piped configuration", zap.Error(err))
		return err
	}

	// Make gRPC client and connect to the API.
	apiClient, err := p.createAPIClient(ctx, t.Logger)
	if err != nil {
		t.Logger.Error("failed to create gRPC client to control plane", zap.Error(err))
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

	// Build kubeconfig for initialing kubernetes clients later.
	kubeConfig, err := clientcmd.BuildConfigFromFlags(p.masterURL, p.kubeconfig)
	if err != nil {
		t.Logger.Error("failed to build kube config", zap.Error(err))
		return err
	}

	// Initialize git client.
	gitClient, err := git.NewClient(p.gitUserName, p.gitEmail, t.Logger)
	if err != nil {
		t.Logger.Error("failed to initialize git client", zap.Error(err))
		return err
	}
	defer gitClient.Clean()

	// Start running application state store.
	{
		s := appstatestore.NewStore(kubeConfig, p.gracePeriod, t.Logger)
		group.Go(func() error {
			return s.Run(ctx)
		})
		// TODO: Do not block other components until this component becomes ready.
		if err := s.WaitForReady(ctx, time.Minute); err != nil {
			return err
		}
	}

	// Start running application state reporter.
	{
		r := appstatereporter.NewReporter(p.gracePeriod)
		group.Go(func() error {
			return r.Run(ctx)
		})
	}

	// Start running application store.
	var applicationStore applicationstore.Store
	{
		applicationStore = applicationstore.NewStore(apiClient, p.gracePeriod, t.Logger)
		group.Go(func() error {
			return applicationStore.Run(ctx)
		})
	}

	// Start running command store.
	var commandStore commandstore.Store
	{
		commandStore = commandstore.NewStore(apiClient, p.gracePeriod, t.Logger)
		group.Go(func() error {
			return commandStore.Run(ctx)
		})
	}

	// Start running deployment controller.
	{
		c := deploymentcontroller.NewController(apiClient, commandStore, pipedConfig, p.gracePeriod, t.Logger)
		group.Go(func() error {
			return c.Run(ctx)
		})
	}

	// Start running deployment trigger.
	{
		t := deploymenttrigger.NewTrigger(apiClient, gitClient, applicationStore, commandStore, pipedConfig, p.gracePeriod, t.Logger)
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

	client, err := pipedservice.NewClient(ctx, p.apiAddress, options...)
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
	return cfg.PipedSpec, nil
}
