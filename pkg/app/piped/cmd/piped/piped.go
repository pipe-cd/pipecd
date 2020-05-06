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
	"github.com/kapetaniosci/pipe/pkg/app/piped/appstatereporter"
	"github.com/kapetaniosci/pipe/pkg/app/piped/appstatestore"
	"github.com/kapetaniosci/pipe/pkg/app/piped/deploymentcontroller"
	"github.com/kapetaniosci/pipe/pkg/app/piped/deploymenttrigger"
	"github.com/kapetaniosci/pipe/pkg/cli"
	"github.com/kapetaniosci/pipe/pkg/config"
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

	kubeconfig string
	masterURL  string
	namespace  string

	useFakeAPIClient bool
	gracePeriod      time.Duration
}

func NewCommand() *cobra.Command {
	r := &piped{
		apiAddress:  "pipecd-api:9091",
		adminPort:   9085,
		gracePeriod: 30 * time.Second,
	}
	cmd := &cobra.Command{
		Use:   "piped",
		Short: "Start running Piped.",
		RunE:  cli.WithContext(r.run),
	}

	cmd.Flags().StringVar(&r.projectID, "project-id", r.projectID, "The identifier of the project which this piped belongs to.")
	cmd.Flags().StringVar(&r.pipedID, "piped-id", r.pipedID, "The unique identifier generated for this piped.")
	cmd.Flags().StringVar(&r.pipedKeyFile, "piped-key-file", r.pipedKeyFile, "The path to the key generated for this piped.")
	cmd.Flags().StringVar(&r.configFile, "config-file", r.configFile, "The path to the configuration file.")
	cmd.Flags().BoolVar(&r.tls, "tls", r.tls, "Whether running the gRPC server with TLS or not.")
	cmd.Flags().StringVar(&r.certFile, "cert-file", r.certFile, "The path to the TLS certificate file.")
	cmd.Flags().StringVar(&r.apiAddress, "api-address", r.apiAddress, "The address used to connect to API server.")
	cmd.Flags().IntVar(&r.adminPort, "admin-port", r.adminPort, "The port number used to run a HTTP server for admin tasks such as metrics, healthz.")

	cmd.Flags().StringVar(&r.kubeconfig, "kube-config", r.kubeconfig, "Path to a kubeconfig. Only required if out-of-cluster.")
	cmd.Flags().StringVar(&r.masterURL, "master", r.masterURL, "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
	cmd.Flags().StringVar(&r.namespace, "namespace", r.namespace, "The namespace where this piped is running.")

	cmd.Flags().BoolVar(&r.useFakeAPIClient, "use-fake-api-client", r.useFakeAPIClient, "Whether the fake api client should be used instead of the real one or not.")
	cmd.Flags().DurationVar(&r.gracePeriod, "grace-period", r.gracePeriod, "How long to wait for graceful shutdown.")

	cmd.MarkFlagRequired("project-id")
	cmd.MarkFlagRequired("piped-id")
	cmd.MarkFlagRequired("piped-key-file")
	cmd.MarkFlagRequired("config-file")

	return cmd
}

func (r *piped) run(ctx context.Context, t cli.Telemetry) error {
	group, ctx := errgroup.WithContext(ctx)

	// Load piped configuration from specified file.
	pipedConfig, err := r.loadConfig()
	if err != nil {
		t.Logger.Error("failed to load piped configuration", zap.Error(err))
		return err
	}

	// Make gRPC client and connect to the API.
	apiClient, err := r.createAPIClient(ctx, t.Logger)
	if err != nil {
		t.Logger.Error("failed to create gRPC client to control plane", zap.Error(err))
		return err
	}

	// Start running admin server.
	{
		admin := admin.NewAdmin(r.adminPort, r.gracePeriod, t.Logger)
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
	kubeConfig, err := clientcmd.BuildConfigFromFlags(r.masterURL, r.kubeconfig)
	if err != nil {
		t.Logger.Error("failed to build kube config", zap.Error(err))
		return err
	}

	// Start running application state store.
	{
		s := appstatestore.NewStore(kubeConfig, r.gracePeriod, t.Logger)
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
		r := appstatereporter.NewReporter(r.gracePeriod)
		group.Go(func() error {
			return r.Run(ctx)
		})
	}

	// Start running deployment controller.
	{
		c := deploymentcontroller.NewController(apiClient, pipedConfig, r.gracePeriod, t.Logger)
		group.Go(func() error {
			return c.Run(ctx)
		})
	}

	// Start running deployment trigger.
	{
		t := deploymenttrigger.NewTrigger(apiClient, pipedConfig, r.gracePeriod, t.Logger)
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
func (r *piped) createAPIClient(ctx context.Context, logger *zap.Logger) (pipedservice.Client, error) {
	if r.useFakeAPIClient {
		return pipedclientfake.NewClient(logger), nil
	}
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	pipedKey, err := ioutil.ReadFile(r.pipedKeyFile)
	if err != nil {
		logger.Error("failed to read piped key file", zap.Error(err))
		return nil, err
	}

	var (
		token   = rpcauth.MakePipedToken(r.projectID, r.pipedID, string(pipedKey))
		tls     = r.certFile != ""
		creds   = rpcclient.NewPerRPCCredentials(token, rpcauth.PipedTokenCredentials, tls)
		options = []rpcclient.DialOption{
			rpcclient.WithBlock(),
			rpcclient.WithStatsHandler(),
			rpcclient.WithPerRPCCredentials(creds),
		}
	)
	if tls {
		options = append(options, rpcclient.WithTLS(r.certFile))
	} else {
		options = append(options, rpcclient.WithInsecure())
	}

	client, err := pipedservice.NewClient(ctx, r.apiAddress, options...)
	if err != nil {
		logger.Error("failed to create api client", zap.Error(err))
		return nil, err
	}
	return client, nil
}

// loadConfig reads the Piped configuration data from the specified file.
func (r *piped) loadConfig() (*config.PipedSpec, error) {
	cfg, err := config.LoadFromYAML(r.configFile)
	if err != nil {
		return nil, err
	}
	if cfg.Kind != config.KindPiped {
		return nil, fmt.Errorf("wrong configuration kind for piped: %v", cfg.Kind)
	}
	return cfg.PipedSpec, nil
}
