// Copyright 2020 The PipeCD Authors.
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

package rpcauth

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/kapetaniosci/pipe/pkg/jwt"
	"github.com/kapetaniosci/pipe/pkg/role"
)

var (
	errUnauthenticated  = status.Error(codes.Unauthenticated, "unauthenticated")
	errPermissionDenied = status.Error(codes.PermissionDenied, "permission denied")
)

// RBACAuthorizer defines a function to check required role for a specific method.
type RBACAuthorizer interface {
	Authorize(string, role.Role) bool
}

type RunnerKeyVerifier interface {
	Verify(projectID, runnerID, runnerKey string) error
}

type (
	credentialsContextKey struct{}
	claimsContextKey      struct{}
)

var (
	credentialsKey = credentialsContextKey{}
	claimsKey      = claimsContextKey{}
)

// RunnerKeyUnaryServerInterceptor extracts credentials from gRPC metadata
// and set the extracted credentials to the context with a fixed key.
// This interceptor will returns a gPRC error when the credentials
// was not set or was malformed.
func RunnerKeyUnaryServerInterceptor(verifier RunnerKeyVerifier, logger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		creds, err := extractCredentials(ctx)
		if err != nil {
			return nil, err
		}
		ctx = context.WithValue(ctx, credentialsKey, creds)
		return handler(ctx, req)
	}
}

// RunnerKeyStreamServerInterceptor extracts credentials from gRPC metadata
// and set the extracted credentials to the context with a fixed key.
// This interceptor will returns a gPRC error when the credentials
// was not set or was malformed.
func RunnerKeyStreamServerInterceptor(verifier RunnerKeyVerifier, logger *zap.Logger) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, _ *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := stream.Context()
		creds, err := extractCredentials(ctx)
		if err != nil {
			return err
		}
		ctx = context.WithValue(ctx, credentialsKey, creds)
		wrappedStream := &wrappedServerStream{
			ServerStream: stream,
			ctx:          ctx,
		}
		return handler(srv, wrappedStream)
	}
}

// JWTUnaryServerInterceptor ensures that the JWT credentials included in the context
// must be verified by verifier.
func JWTUnaryServerInterceptor(verifier jwt.Verifier, authorizer RBACAuthorizer, logger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		cookie, err := extractCookie(ctx)
		if err != nil {
			logger.Warn("failed to extract cookie", zap.Error(err))
			return nil, errUnauthenticated
		}
		token, ok := cookie[jwt.SignedTokenKey]
		if !ok {
			logger.Warn("token does not exist in cookie")
			return nil, errUnauthenticated
		}
		claims, err := verifier.Verify(token)
		if err != nil {
			logger.Warn("unable to verify token", zap.Error(err))
			return nil, errUnauthenticated
		}
		if !authorizer.Authorize(info.FullMethod, claims.Role) {
			logger.Warn(fmt.Sprintf("unsufficient permission for method: %s", info.FullMethod),
				zap.Any("claims", claims),
			)
			return nil, errPermissionDenied
		}
		ctx = context.WithValue(ctx, claimsKey, *claims)
		return handler(ctx, req)
	}
}

// ExtractCredentials returns the credentials inside a given context.
func ExtractCredentials(ctx context.Context) (Credentials, error) {
	creds, ok := ctx.Value(credentialsKey).(Credentials)
	if !ok {
		return creds, errUnauthenticated
	}
	return creds, nil
}

// ExtractClaims returns the claims inside a given context.
func ExtractClaims(ctx context.Context) (jwt.Claims, error) {
	claims, ok := ctx.Value(claimsKey).(jwt.Claims)
	if !ok {
		return claims, errUnauthenticated
	}
	return claims, nil
}
