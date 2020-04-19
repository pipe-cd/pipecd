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
	"time"

	"github.com/spf13/cobra"

	"github.com/kapetaniosci/pipe/pkg/cli"
)

type runner struct {
	apiAddress  string
	apiPort     int
	gracePeriod time.Duration

	tls      bool
	certFile string
}

func NewCommand() *cobra.Command {
	r := &runner{
		apiAddress:  "pipecd-api",
		apiPort:     9080,
		gracePeriod: 15 * time.Second,
	}
	cmd := &cobra.Command{
		Use:   "runner",
		Short: "Start running Runner.",
		RunE:  cli.WithContext(r.run),
	}

	cmd.Flags().StringVar(&r.apiAddress, "api-address", r.apiAddress, "The address used to connect to API server.")
	cmd.Flags().IntVar(&r.apiPort, "api-port", r.apiPort, "The port used to connect to API server.")
	cmd.Flags().DurationVar(&r.gracePeriod, "grace-period", r.gracePeriod, "How long to wait for graceful shutdown.")

	cmd.Flags().BoolVar(&r.tls, "tls", r.tls, "Whether running the gRPC server with TLS or not.")
	cmd.Flags().StringVar(&r.certFile, "cert-file", r.certFile, "The path to the TLS certificate file.")

	return cmd
}

func (r *runner) run(ctx context.Context, t cli.Telemetry) error {
	return nil
}
