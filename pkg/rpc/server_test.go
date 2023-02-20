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

package rpc

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/app/helloworld/api"
	"github.com/pipe-cd/pipecd/pkg/app/helloworld/service"
	"github.com/pipe-cd/pipecd/pkg/rpc/rpcauth"
	"github.com/pipe-cd/pipecd/pkg/rpc/rpcclient"
)

func TestMain(m *testing.M) {
	logger := zap.NewExample()
	server := NewServer(
		api.NewHelloWorldAPI(api.WithLogger(logger)),
		WithTLS("testdata/tls.crt", "testdata/tls.key"),
		WithPort(9090),
		WithGracePeriod(time.Second),
		WithLogger(logger),
		WithPipedTokenAuthUnaryInterceptor(testPipedTokenVerifier{"test-piped-key"}, logger),
	)
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	go server.Run(ctx)
	os.Exit(m.Run())
}

type testPipedTokenVerifier struct {
	pipedKey string
}

func (v testPipedTokenVerifier) Verify(ctx context.Context, projectID, pipedID, pipedKey string) error {
	if pipedKey != v.pipedKey {
		return fmt.Errorf("invalid piped key, want: %s, got: %s", v.pipedKey, pipedKey)
	}
	return nil
}

func TestRPCRequestOK(t *testing.T) {
	ctx := context.Background()
	pipedToken := rpcauth.MakePipedToken("test-project-id", "test-piped-id", "test-piped-key")
	creds := rpcclient.NewPerRPCCredentials(pipedToken, rpcauth.PipedTokenCredentials, true)
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
