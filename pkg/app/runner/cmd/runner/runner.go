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

package runner

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
	apiservice "github.com/kapetaniosci/pipe/pkg/app/api/service"
	"github.com/kapetaniosci/pipe/pkg/app/runner/appstatereporter"
	"github.com/kapetaniosci/pipe/pkg/app/runner/appstatestore"
	"github.com/kapetaniosci/pipe/pkg/app/runner/deploymentcontroller"
	"github.com/kapetaniosci/pipe/pkg/app/runner/deploymenttrigger"
	"github.com/kapetaniosci/pipe/pkg/cli"
	"github.com/kapetaniosci/pipe/pkg/config"
	"github.com/kapetaniosci/pipe/pkg/rpc/rpcauth"
	"github.com/kapetaniosci/pipe/pkg/rpc/rpcclient"
)

type runner struct {
	projectID     string
	runnerID      string
	runnerKeyFile string
	configFile    string
	tls           bool
	certFile      string
	apiAddress    string
	adminPort     int

	kubeconfig string
	masterURL  string
	namespace  string

	gracePeriod time.Duration
}

func NewCommand() *cobra.Command {
	r := &runner{
		apiAddress:  "pipecd-api:9091",
		adminPort:   9085,
		gracePeriod: 30 * time.Second,
	}
	cmd := &cobra.Command{
		Use:   "runner",
		Short: "Start running Runner.",
		RunE:  cli.WithContext(r.run),
	}

	cmd.Flags().StringVar(&r.projectID, "project-id", r.projectID, "The identifier of the project which this runner belongs to.")
	cmd.Flags().StringVar(&r.runnerID, "runner-id", r.runnerID, "The unique identifier generated for this runner.")
	cmd.Flags().StringVar(&r.runnerKeyFile, "runner-key-file", r.runnerKeyFile, "The path to the key generated for this runner.")
	cmd.Flags().StringVar(&r.configFile, "config-file", r.configFile, "The path to the configuration file.")
	cmd.Flags().BoolVar(&r.tls, "tls", r.tls, "Whether running the gRPC server with TLS or not.")
	cmd.Flags().StringVar(&r.certFile, "cert-file", r.certFile, "The path to the TLS certificate file.")
	cmd.Flags().StringVar(&r.apiAddress, "api-address", r.apiAddress, "The address used to connect to API server.")
	cmd.Flags().IntVar(&r.adminPort, "admin-port", r.adminPort, "The port number used to run a HTTP server for admin tasks such as metrics, healthz.")

	cmd.Flags().StringVar(&r.kubeconfig, "kube-config", r.kubeconfig, "Path to a kubeconfig. Only required if out-of-cluster.")
	cmd.Flags().StringVar(&r.masterURL, "master", r.masterURL, "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
	cmd.Flags().StringVar(&r.namespace, "namespace", r.namespace, "The namespace where this runner is running.")

	cmd.Flags().DurationVar(&r.gracePeriod, "grace-period", r.gracePeriod, "How long to wait for graceful shutdown.")

	cmd.MarkFlagRequired("project-id")
	cmd.MarkFlagRequired("runner-id")
	cmd.MarkFlagRequired("runner-key-file")
	cmd.MarkFlagRequired("config-file")

	return cmd
}

func (r *runner) run(ctx context.Context, t cli.Telemetry) error {
	group, ctx := errgroup.WithContext(ctx)

	// Load runner configuration from specified file.
	_, err := r.loadConfig()
	if err != nil {
		t.Logger.Error("failed to load runner configuration", zap.Error(err))
		return err
	}

	// Make gRPC client and connect to the API.
	_, err = r.createAPIClient(ctx, t.Logger)
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
		c := deploymentcontroller.NewController(r.gracePeriod)
		group.Go(func() error {
			return c.Run(ctx)
		})
	}

	// Start running deployment trigger.
	{
		t := deploymenttrigger.NewTrigger(r.gracePeriod)
		group.Go(func() error {
			return t.Run(ctx)
		})
	}

	// Wait until all runner components have finished.
	// A terminating signal or a finish of any components
	// could trigger the finish of runner.
	// This ensures that all components are good or no one.
	if err := group.Wait(); err != nil {
		t.Logger.Error("failed while running", zap.Error(err))
		return err
	}
	return nil
}

// createAPIClient makes a gRPC client to connect to the API.
func (r *runner) createAPIClient(ctx context.Context, logger *zap.Logger) (apiservice.Client, error) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	runnerKey, err := ioutil.ReadFile(r.runnerKeyFile)
	if err != nil {
		logger.Error("failed to read runner key file", zap.Error(err))
		return nil, err
	}

	var (
		token   = rpcauth.MakeRunnerToken(r.projectID, r.runnerID, string(runnerKey))
		tls     = r.certFile != ""
		creds   = rpcclient.NewPerRPCCredentials(token, rpcauth.RunnerTokenCredentials, tls)
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

	client, err := apiservice.NewClient(ctx, r.apiAddress, options...)
	if err != nil {
		logger.Error("failed to create api client", zap.Error(err))
		return nil, err
	}
	return client, nil
}

// loadConfig reads the Runner configuration data from the specified file.
func (r *runner) loadConfig() (*config.RunnerSpec, error) {
	cfg, err := config.LoadFromYAML(r.configFile)
	if err != nil {
		return nil, err
	}
	if cfg.Kind != config.KindRunner {
		return nil, fmt.Errorf("wrong configuration kind for runner: %v", cfg.Kind)
	}
	return cfg.RunnerSpec, nil
}
