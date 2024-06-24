package grpcapi

import (
	"context"
	"errors"
	"strings"

	"github.com/pipe-cd/pipecd/pkg/app/server/pipedverifier"
	"github.com/pipe-cd/pipecd/pkg/rpc"

	"github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
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

	projectID, pipedID, pipedKey, err := e.parsePipedToken(a)
	if err != nil {
		return &authv3.CheckResponse{Status: status.New(codes.PermissionDenied, err.Error()).Proto()}, nil
	}

	if err := e.verifier.Verify(ctx, projectID, pipedID, pipedKey); err != nil {
		return &authv3.CheckResponse{Status: status.New(codes.PermissionDenied, err.Error()).Proto()}, nil
	}

	return &authv3.CheckResponse{Status: status.New(codes.OK, "OK").Proto()}, nil
}

func (e *EnvoyAuthorizationServer) parsePipedToken(a string) (string, string, string, error) {
	if !strings.HasPrefix(a, "Bearer ") {
		return "", "", "", errors.New("invalid authorization header")
	}

	parts := strings.Split(strings.TrimPrefix(a, "Bearer "), ",")
	if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
		return "", "", "", errors.New("malformed piped token")
	}
	return parts[0], parts[1], parts[2], nil
}
