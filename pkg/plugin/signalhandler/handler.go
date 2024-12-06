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

// Package planner provides a piped component
// that decides the deployment pipeline of a deployment.
// The planner bases on the changes from git commits
// then builds the deployment manifests to know the behavior of the deployment.
// From that behavior the planner can decides which pipeline should be applied.

package signalhandler

import (
	"context"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"

	"google.golang.org/grpc"
)

var (
	terminated atomic.Bool

	signals = []os.Signal{syscall.SIGINT, syscall.SIGTERM}
)

func init() {
	ctx, cancel := signal.NotifyContext(context.Background(), signals...)
	go func() {
		defer cancel()
		<-ctx.Done()
		terminated.Store(true)
	}()
}

// Terminated returns true if the signal handler has received a termination signal.
func Terminated() bool {
	return terminated.Load()
}

// CancelInterceptor is a gRPC interceptor that cancels the context when a termination signal is received.
func CancelInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	ctx, cancel := signal.NotifyContext(ctx, signals...)
	defer cancel()

	return handler(ctx, req)
}
