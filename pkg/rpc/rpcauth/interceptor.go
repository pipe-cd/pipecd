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

package rpcauth

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/pipe-cd/pipecd/pkg/jwt"
	"github.com/pipe-cd/pipecd/pkg/model"
)

var (
	errUnauthenticated  = status.Error(codes.Unauthenticated, "Unauthenticated")
	errPermissionDenied = status.Error(codes.PermissionDenied, "Permission Denied")
)

// RBACAuthorizer defines a function to check required role for a specific RPC method.
type RBACAuthorizer interface {
	Authorize(context.Context, string, model.Role) bool
}

// PipedTokenVerifier verifies the given piped token.
type PipedTokenVerifier interface {
	Verify(ctx context.Context, projectID, pipedID, pipedKey string) error
}

// APIKeyVerifier verifies the given API key.
type APIKeyVerifier interface {
	Verify(ctx context.Context, key string) (*model.APIKey, error)
}

type (
	claimsContextKey       struct{}
	pipedTokenContextKey   struct{}
	pipedTokenContextValue struct {
		ProjectID string
		PipedID   string
		PipedKey  string
	}
	apiKeyContextKey struct{}
)

var (
	claimsKey     = claimsContextKey{}
	pipedTokenKey = pipedTokenContextKey{}
	apiKeyKey     = apiKeyContextKey{}
)

// PipedTokenUnaryServerInterceptor extracts credentials from gRPC metadata
// and validates it by the specified Verifier.
// If the token was valid the parsed ProjectID, PipedID, PipedKey will be set to the context.
func PipedTokenUnaryServerInterceptor(verifier PipedTokenVerifier, logger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		creds, err := extractCredentials(ctx)
		if err != nil {
			return nil, err
		}
		if creds.Type != PipedTokenCredentials {
			logger.Warn("wrong credentials type for PipedTokenCredentials", zap.Any("credentials", creds))
			return nil, errUnauthenticated
		}
		projectID, pipedID, pipedKey, err := parsePipedToken(creds.Data)
		if err != nil {
			logger.Warn(fmt.Sprintf("malformed credentials: %s, err: %v", creds.Data, err))
			return nil, errUnauthenticated
		}
		if err := verifier.Verify(ctx, projectID, pipedID, pipedKey); err != nil {
			logger.Warn("unable to verify piped token", zap.Error(err))
			return nil, errUnauthenticated
		}
		ctx = context.WithValue(ctx, pipedTokenKey, pipedTokenContextValue{
			ProjectID: projectID,
			PipedID:   pipedID,
			PipedKey:  pipedKey,
		})
		return handler(ctx, req)
	}
}

// PipedTokenStreamServerInterceptor extracts credentials from gRPC metadata
// and set the extracted credentials to the context with a fixed key.
// This interceptor will returns a gPRC error when the credentials
// was not set or was malformed.
func PipedTokenStreamServerInterceptor(verifier PipedTokenVerifier, logger *zap.Logger) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, _ *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := stream.Context()
		creds, err := extractCredentials(ctx)
		if err != nil {
			return err
		}
		if creds.Type != PipedTokenCredentials {
			logger.Warn("wrong credentials type for PipedTokenCredentials", zap.Any("credentials", creds))
			return errUnauthenticated
		}
		projectID, pipedID, pipedKey, err := parsePipedToken(creds.Data)
		if err != nil {
			logger.Warn(fmt.Sprintf("malformed credentials: %s, err: %v", creds.Data, err))
			return errUnauthenticated
		}
		if err := verifier.Verify(ctx, projectID, pipedID, pipedKey); err != nil {
			logger.Warn("unable to verify piped token", zap.Error(err))
			return errUnauthenticated
		}
		ctx = context.WithValue(ctx, pipedTokenKey, pipedTokenContextValue{
			ProjectID: projectID,
			PipedID:   pipedID,
			PipedKey:  pipedKey,
		})
		wrappedStream := &wrappedServerStream{
			ServerStream: stream,
			ctx:          ctx,
		}
		return handler(srv, wrappedStream)
	}
}

// ExtractPipedToken returns the verified piped key inside a given context.
func ExtractPipedToken(ctx context.Context) (projectID, pipedID, pipedKey string, err error) {
	v, ok := ctx.Value(pipedTokenKey).(pipedTokenContextValue)
	if !ok {
		err = errUnauthenticated
		return
	}
	projectID = v.ProjectID
	pipedID = v.PipedID
	pipedKey = v.PipedKey
	return
}

// APIKeyUnaryServerInterceptor extracts credentials from gRPC metadata
// and validates it by the specified Verifier.
// The valid API key will be set to the context.
func APIKeyUnaryServerInterceptor(verifier APIKeyVerifier, logger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		creds, err := extractCredentials(ctx)
		if err != nil {
			return nil, err
		}
		if creds.Type != APIKeyCredentials {
			logger.Warn("wrong credentials type for APIKeyCredentials", zap.Any("credentials", creds))
			return nil, errUnauthenticated
		}
		apiKey, err := verifier.Verify(ctx, creds.Data)
		if err != nil {
			logger.Warn("unable to verify api key", zap.Error(err))
			return nil, errUnauthenticated
		}
		ctx = ContextWithAPIKey(ctx, apiKey)
		return handler(ctx, req)
	}
}

// ContextWithAPIKey returns a new context in which the given API key was attached.
func ContextWithAPIKey(ctx context.Context, k *model.APIKey) context.Context {
	return context.WithValue(ctx, apiKeyKey, k)
}

// ExtractAPIKey returns the verified API key inside the given context.
func ExtractAPIKey(ctx context.Context) (*model.APIKey, error) {
	k, ok := ctx.Value(apiKeyKey).(*model.APIKey)
	if !ok {
		return nil, errUnauthenticated
	}
	return k, nil
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
		if !authorizer.Authorize(ctx, info.FullMethod, claims.Role) {
			logger.Warn(fmt.Sprintf("unsufficient permission for method: %s", info.FullMethod),
				zap.Any("claims", claims),
			)
			return nil, errPermissionDenied
		}
		ctx = context.WithValue(ctx, claimsKey, *claims)
		return handler(ctx, req)
	}
}

// ExtractClaims returns the claims inside a given context.
func ExtractClaims(ctx context.Context) (jwt.Claims, error) {
	claims, ok := ctx.Value(claimsKey).(jwt.Claims)
	if !ok {
		return claims, errUnauthenticated
	}
	return claims, nil
}
