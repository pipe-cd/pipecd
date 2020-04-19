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
	"io/ioutil"
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	apiservice "github.com/kapetaniosci/pipe/pkg/app/api/service"
	"github.com/kapetaniosci/pipe/pkg/cli"
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
	gracePeriod   time.Duration
}

func NewCommand() *cobra.Command {
	r := &runner{
		apiAddress:  "pipecd-api:9091",
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
	cmd.Flags().DurationVar(&r.gracePeriod, "grace-period", r.gracePeriod, "How long to wait for graceful shutdown.")

	cmd.MarkFlagRequired("project-id")
	cmd.MarkFlagRequired("runner-id")
	cmd.MarkFlagRequired("runner-key-file")
	cmd.MarkFlagRequired("config-file")

	return cmd
}

func (r *runner) run(ctx context.Context, t cli.Telemetry) error {
	return nil
}

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
