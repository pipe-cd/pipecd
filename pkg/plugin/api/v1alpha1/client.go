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

package pluginapi

import (
	"context"

	"google.golang.org/grpc"

	"github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1/deployment"
	"github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1/livestate"
	"github.com/pipe-cd/pipecd/pkg/rpc/rpcclient"
)

type PluginClient interface {
	deployment.DeploymentServiceClient
	livestate.LivestateServiceClient
	Close() error
}

type client struct {
	deployment.DeploymentServiceClient
	livestate.LivestateServiceClient
	conn *grpc.ClientConn
}

func NewClient(ctx context.Context, address string, opts ...rpcclient.DialOption) (PluginClient, error) {
	conn, err := rpcclient.DialContext(ctx, address, opts...)
	if err != nil {
		return nil, err
	}

	return &client{
		DeploymentServiceClient: deployment.NewDeploymentServiceClient(conn),
		LivestateServiceClient:  livestate.NewLivestateServiceClient(conn),
		conn:                    conn,
	}, nil
}

func (c *client) Close() error {
	return c.conn.Close()
}
