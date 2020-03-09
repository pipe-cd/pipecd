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

package rpcauth

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/nghialv/dianomi/pkg/jwt"
	"github.com/nghialv/dianomi/pkg/role"
)

var (
	errUnauthenticated  = status.Error(codes.Unauthenticated, "unauthenticated")
	errPermissionDenied = status.Error(codes.PermissionDenied, "permission denied")
)

// Authorizer defines a function to check required role for a specific method.
type Authorizer interface {
	Authorize(string, role.Role) bool
}

type (
	credentialsContextKey struct{}
	claimsContextKey      struct{}
)

var (
	credentialsKey = credentialsContextKey{}
	claimsKey      = claimsContextKey{}
)

// UnaryServerInterceptor extracts credentials from gRPC metadata
// and set the extracted credentials to the context with a fixed key.
// This interceptor will returns a gPRC error when the credentials
// was not set or was malformed.
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		creds, err := extractCredentials(ctx)
		if err != nil {
			return nil, err
		}
		ctx = context.WithValue(ctx, credentialsKey, creds)
		return handler(ctx, req)
	}
}

// StreamServerInterceptor extracts credentials from gRPC metadata
// and set the extracted credentials to the context with a fixed key.
// This interceptor will returns a gPRC error when the credentials
// was not set or was malformed.
func StreamServerInterceptor() grpc.StreamServerInterceptor {
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

// ServiceKeyUnaryServerInterceptor ensures that the credentials included in the context
// must be same with the give service key.
// This interceptor must be appended after UnaryServerInterceptor.
func ServiceKeyUnaryServerInterceptor(key string, logger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if msg, err := RequireServiceKey(ctx, key); err != nil {
			logger.Warn(msg, zap.Error(err))
			return nil, err
		}
		return handler(ctx, req)
	}
}

// RequireServiceKey is a helper to check service key from gRPC request context.
func RequireServiceKey(ctx context.Context, expected string) (msg string, err error) {
	err = errUnauthenticated
	creds, e := ExtractCredentials(ctx)
	if e != nil {
		msg = "service key credentials has not been set in context by interceptor"
		return
	}
	if creds.Type != ServiceKeyCredentials {
		msg = fmt.Sprintf("expected %v but got %v type", ServiceKeyCredentials, creds.Type)
		return
	}
	if creds.Data != expected {
		msg = fmt.Sprintf("invalid service key: %s", creds.Data)
		return
	}
	err = nil
	return
}

// JwtUnaryServerInterceptor ensures that the credentials included in the context
// must be verified by verifier.
func JwtUnaryServerInterceptor(verifier jwt.Verifier, authorizer Authorizer, logger *zap.Logger) grpc.UnaryServerInterceptor {
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
