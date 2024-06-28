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

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type requestValidator interface {
	Validate() error
}

// RequestValidationUnaryServerInterceptor validates the request payload if
// the request implements requestValidator interface.
// An InvalidArgument with the detail message will be returned to client if
// the validation was not passed.
func RequestValidationUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if v, ok := req.(requestValidator); ok {
			if err := v.Validate(); err != nil {
				return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid request: %v", err))
			}
		}
		return handler(ctx, req)
	}
}
