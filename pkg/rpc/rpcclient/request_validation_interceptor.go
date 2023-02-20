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

package rpcclient

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

// RequestValidationUnaryClientInterceptor validates the request payload if
// the request implements requestValidator interface.
// An InvalidArgument with the detail message will be returned to client if
// the validation was not passed.
func RequestValidationUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if v, ok := req.(requestValidator); ok {
			if err := v.Validate(); err != nil {
				return status.Error(codes.InvalidArgument, fmt.Sprintf("invalid request: %v", err))
			}
		}
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
