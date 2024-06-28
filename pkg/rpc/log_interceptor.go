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

package rpc

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// LogUnaryServerInterceptor logs handled unary gRPC requests.
func LogUnaryServerInterceptor(logger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		code := status.Code(err)

		switch code {
		case codes.Internal:
			logger.Error(fmt.Sprintf("failed to handle an unary gRPC request: %s", info.FullMethod),
				zap.Error(err),
				zap.Duration("duration", time.Since(start)),
			)
		default:
			logger.Info(fmt.Sprintf("handled an unary gRPC request: %s", info.FullMethod),
				zap.String("code", code.String()),
				zap.Error(err),
				zap.Duration("duration", time.Since(start)),
			)
		}
		return resp, err
	}
}
