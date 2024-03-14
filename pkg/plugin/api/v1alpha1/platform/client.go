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

// Package planner provides a piped component
// that decides the deployment pipeline of a deployment.
// The planner bases on the changes from git commits
// then builds the deployment manifests to know the behavior of the deployment.
// From that behavior the planner can decides which pipeline should be applied.
package platform

import (
	"context"

	"google.golang.org/grpc"

	"github.com/pipe-cd/pipecd/pkg/rpc/rpcclient"
)

type PlatformPluginClient interface {
	PlannerServiceClient
	ExecutorServiceClient
	Close() error
}

type client struct {
	PlannerServiceClient
	ExecutorServiceClient
	conn *grpc.ClientConn
}

func NewClient(ctx context.Context, address string, opts ...rpcclient.DialOption) (PlatformPluginClient, error) {
	conn, err := rpcclient.DialContext(ctx, address, opts...)
	if err != nil {
		return nil, err
	}

	return &client{
		PlannerServiceClient:  NewPlannerServiceClient(conn),
		ExecutorServiceClient: NewExecutorServiceClient(conn),
		conn:                  conn,
	}, nil
}

func (c *client) Close() error {
	return c.conn.Close()
}
