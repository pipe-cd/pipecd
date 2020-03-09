// Copyright 2020 The Dianomi Authors.
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

package rpc

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/nghialv/dianomi/pkg/app/helloworld/api"
	"github.com/nghialv/dianomi/pkg/app/helloworld/service"
	"github.com/nghialv/dianomi/pkg/rpc/rpcauth"
	"github.com/nghialv/dianomi/pkg/rpc/rpcclient"
)

func TestMain(m *testing.M) {
	logger := zap.NewExample()
	server := NewServer(
		api.NewHelloWorldService(api.WithLogger(logger)),
		WithTLS("testdata/tls.crt", "testdata/tls.key"),
		WithPort(9090),
		WithLogger(logger),
		WithAuthUnaryInterceptor(),
	)
	defer server.Stop(time.Second)
	go server.Run()
	os.Exit(m.Run())
}

func TestRPCRequestOK(t *testing.T) {
	ctx := context.Background()
	creds := rpcclient.NewPerRPCCredentials("service-key", rpcauth.ServiceKeyCredentials, true)
	var cli service.Client
	var err error
	// Waiting the gRPC server.
	for i := 0; i < 3; i++ {
		ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
		defer cancel()
		cli, err = service.NewClient(
			ctx,
			"localhost:9090",
			rpcclient.WithBlock(),
			rpcclient.WithStatsHandler(),
			rpcclient.WithTLS("testdata/tls.crt"),
			rpcclient.WithPerRPCCredentials(creds),
		)
		if err == nil {
			break
		}
	}
	require.NoError(t, err)
	require.NotNil(t, cli)
	defer cli.Close()

	ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()
	resp, err := cli.Hello(ctx, &service.HelloRequest{
		Name:   "test-name",
		Gender: service.HelloRequest_GENDER_MALE,
	})
	require.NoError(t, err)
	assert.True(t, len(resp.Message) > 0)
}

func TestRPCRequestWithoutCredentials(t *testing.T) {
	ctx := context.Background()
	var cli service.Client
	var err error
	// Waiting the gRPC server.
	for i := 0; i < 3; i++ {
		ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
		defer cancel()
		cli, err = service.NewClient(
			ctx,
			"localhost:9090",
			rpcclient.WithBlock(),
			rpcclient.WithStatsHandler(),
			rpcclient.WithTLS("testdata/tls.crt"),
		)
		if err == nil {
			break
		}
	}
	require.NoError(t, err)
	require.NotNil(t, cli)
	defer cli.Close()

	ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()
	_, err = cli.Hello(ctx, &service.HelloRequest{
		Name:   "test-name",
		Gender: service.HelloRequest_GENDER_MALE,
	})
	assert.NotNil(t, err)
}
