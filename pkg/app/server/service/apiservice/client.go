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

package apiservice

import (
	"context"

	"google.golang.org/grpc"

	"github.com/pipe-cd/pipecd/pkg/rpc/rpcclient"
)

type Client interface {
	APIServiceClient
	Close() error
}

type client struct {
	APIServiceClient
	conn *grpc.ClientConn
}

func NewClient(ctx context.Context, addr string, opts ...rpcclient.DialOption) (Client, error) {
	conn, err := rpcclient.DialContext(ctx, addr, opts...)
	if err != nil {
		return nil, err
	}
	return &client{
		APIServiceClient: NewAPIServiceClient(conn),
		conn:             conn,
	}, nil
}

func (c *client) Close() error {
	return c.conn.Close()
}
