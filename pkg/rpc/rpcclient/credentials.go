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

package rpcclient

import (
	"context"
	"fmt"
	"os"
	"strings"

	"google.golang.org/grpc/credentials"

	"github.com/pipe-cd/pipecd/pkg/rpc/rpcauth"
)

type perRPCCredentials struct {
	credentials              string
	credentialsType          rpcauth.CredentialsType
	requireTransportSecurity bool
}

func NewPerRPCCredentials(credentials string, t rpcauth.CredentialsType, requireTransportSecurity bool) credentials.PerRPCCredentials {
	return perRPCCredentials{
		credentials:              strings.TrimSpace(credentials),
		credentialsType:          t,
		requireTransportSecurity: requireTransportSecurity,
	}
}

func NewPerRPCCredentialsFromFile(credentialsFile string, t rpcauth.CredentialsType, requireTransportSecurity bool) (credentials.PerRPCCredentials, error) {
	credentials, err := os.ReadFile(credentialsFile)
	if err != nil {
		return nil, err
	}
	return NewPerRPCCredentials(string(credentials), t, requireTransportSecurity), nil
}

func (c perRPCCredentials) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": fmt.Sprintf("%s %s", string(c.credentialsType), c.credentials),
	}, nil
}

func (c perRPCCredentials) RequireTransportSecurity() bool {
	return c.requireTransportSecurity
}
