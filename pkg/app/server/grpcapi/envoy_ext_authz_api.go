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

package grpcapi

import (
	"context"
	"errors"
	"strings"

	"github.com/pipe-cd/pipecd/pkg/app/server/pipedverifier"
	"github.com/pipe-cd/pipecd/pkg/rpc"
	"github.com/pipe-cd/pipecd/pkg/rpc/rpcauth"

	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	_ authv3.AuthorizationServer = (*EnvoyAuthorizationServer)(nil)
	_ rpc.Service                = (*EnvoyAuthorizationServer)(nil)
)

type EnvoyAuthorizationServer struct {
	authv3.UnimplementedAuthorizationServer

	verifier *pipedverifier.Verifier
}

func NewEnvoyAuthorizationServer(verifier *pipedverifier.Verifier) *EnvoyAuthorizationServer {
	return &EnvoyAuthorizationServer{
		verifier: verifier,
	}
}

// Register implements rpc.Service.
func (e *EnvoyAuthorizationServer) Register(server *grpc.Server) {
	authv3.RegisterAuthorizationServer(server, e)
}

// Check implements authv3.AuthorizationServer.
func (e *EnvoyAuthorizationServer) Check(ctx context.Context, request *authv3.CheckRequest) (*authv3.CheckResponse, error) {
	a, ok := request.GetAttributes().GetRequest().GetHttp().GetHeaders()["authorization"]
	if !ok {
		return &authv3.CheckResponse{Status: status.New(codes.Unauthenticated, "missing authorization header").Proto()}, nil
	}

	projectID, pipedID, pipedKey, err := e.parseAuthorizationHeader(a)
	if err != nil {
		return &authv3.CheckResponse{Status: status.New(codes.PermissionDenied, err.Error()).Proto()}, nil
	}

	if err := e.verifier.Verify(ctx, projectID, pipedID, pipedKey); err != nil {
		return &authv3.CheckResponse{Status: status.New(codes.PermissionDenied, err.Error()).Proto()}, nil
	}

	return &authv3.CheckResponse{Status: status.New(codes.OK, "OK").Proto()}, nil
}

func (e *EnvoyAuthorizationServer) parseAuthorizationHeader(header string) (string, string, string, error) {
	if !strings.HasPrefix(header, "Bearer ") {
		return "", "", "", errors.New("invalid authorization header")
	}

	return rpcauth.ParsePipedToken(strings.TrimPrefix(header, "Bearer "))
}
