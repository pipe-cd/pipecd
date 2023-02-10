// Copyright 2023 The PipeCD Authors.
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

package client

import (
	"context"
	"crypto/tls"
	"errors"
	"time"

	"github.com/spf13/cobra"
	"google.golang.org/grpc/credentials"

	"github.com/pipe-cd/pipecd/pkg/app/server/service/apiservice"
	"github.com/pipe-cd/pipecd/pkg/model"
	"github.com/pipe-cd/pipecd/pkg/rpc/rpcauth"
	"github.com/pipe-cd/pipecd/pkg/rpc/rpcclient"
)

type Options struct {
	Address    string
	APIKey     string
	APIKeyFile string
	Insecure   bool
	CertFile   string
}

func (o *Options) RegisterPersistentFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&o.Address, "address", o.Address, "The address to control-plane api.")
	cmd.PersistentFlags().StringVar(&o.APIKey, "api-key", o.APIKey, "The API key used while authenticating with control-plane.")
	cmd.PersistentFlags().StringVar(&o.APIKeyFile, "api-key-file", o.APIKeyFile, "Path to the file containing API key used while authenticating with control-plane.")
	cmd.PersistentFlags().BoolVar(&o.Insecure, "insecure", o.Insecure, "Whether disabling transport security while connecting to control-plane.")
	cmd.PersistentFlags().StringVar(&o.CertFile, "cert-file", o.CertFile, "The path to the TLS certificate file.")
}

func (o *Options) Validate() error {
	if o.Address == "" {
		return errors.New("address must be set")
	}
	if o.APIKey == "" && o.APIKeyFile == "" {
		return errors.New("either api-key or api-key-file must be set")
	}
	return nil
}

func (o *Options) NewClient(ctx context.Context) (apiservice.Client, error) {
	if err := o.Validate(); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var creds credentials.PerRPCCredentials
	var err error

	if o.APIKey != "" {
		creds = rpcclient.NewPerRPCCredentials(o.APIKey, rpcauth.APIKeyCredentials, !o.Insecure)
	} else {
		creds, err = rpcclient.NewPerRPCCredentialsFromFile(o.APIKeyFile, rpcauth.APIKeyCredentials, !o.Insecure)
		if err != nil {
			return nil, err
		}
	}

	options := []rpcclient.DialOption{
		rpcclient.WithBlock(),
		rpcclient.WithPerRPCCredentials(creds),
	}

	if !o.Insecure {
		if o.CertFile != "" {
			options = append(options, rpcclient.WithTLS(o.CertFile))
		} else {
			config := &tls.Config{}
			options = append(options, rpcclient.WithTransportCredentials(credentials.NewTLS(config)))
		}
	} else {
		options = append(options, rpcclient.WithInsecure())
	}

	client, err := apiservice.NewClient(ctx, o.Address, options...)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func getCommand(ctx context.Context, cli apiservice.Client, cmdID string) (*model.Command, error) {
	req := &apiservice.GetCommandRequest{
		CommandId: cmdID,
	}
	resp, err := cli.GetCommand(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Command, nil
}
