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

package samplecli

import (
	"context"
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	helloworldservice "github.com/pipe-cd/pipecd/pkg/app/helloworld/service"
	"github.com/pipe-cd/pipecd/pkg/cli"
	"github.com/pipe-cd/pipecd/pkg/rpc/rpcclient"
)

type samplecli struct {
	address string
	name    string
}

func NewCommand() *cobra.Command {
	s := &samplecli{
		address: "localhost:9080",
		name:    "samplecli",
	}
	cmd := &cobra.Command{
		Use:   "samplecli",
		Short: "Start running sample client to HelloWorld service",
		RunE:  cli.WithContext(s.run),
	}
	cmd.Flags().StringVar(&s.address, "address", s.address, "The address to HelloWorld service.")
	cmd.Flags().StringVar(&s.name, "name", s.name, "The name to be sent.")
	return cmd
}

func (s *samplecli) run(ctx context.Context, input cli.Input) error {
	cli, err := s.createHelloWorldClient(ctx, input.Logger)
	if err != nil {
		input.Logger.Error("failed to create client", zap.Error(err))
		return err
	}
	defer cli.Close()

	req := &helloworldservice.HelloRequest{
		Name:   s.name,
		Gender: helloworldservice.HelloRequest_GENDER_MALE,
	}
	resp, err := cli.Hello(ctx, req)
	if err != nil {
		input.Logger.Error("failed to send Hello", zap.Error(err))
		return err
	}
	input.Logger.Info("succeeded to send Hello", zap.String("message", resp.Message))
	return nil
}

func (s *samplecli) createHelloWorldClient(ctx context.Context, logger *zap.Logger) (helloworldservice.Client, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	options := []rpcclient.DialOption{
		rpcclient.WithBlock(),
		rpcclient.WithInsecure(),
	}
	client, err := helloworldservice.NewClient(ctx, s.address, options...)
	if err != nil {
		logger.Error("failed to create HelloWorld client", zap.Error(err))
		return nil, err
	}
	return client, nil
}
