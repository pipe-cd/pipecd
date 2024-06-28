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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func TestChainUnaryServerInterceptors(t *testing.T) {
	type parentKey string
	parent := parentKey("parent")
	ctx := context.WithValue(context.Background(), parent, "")
	serverInfo := &grpc.UnaryServerInfo{
		FullMethod: "service.test",
	}
	out := "out"
	var firstRun, secondRun, handlerRun bool
	first := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		require.Equal(t, serverInfo, info)
		require.Equal(t, "", ctx.Value(parent).(string))
		ctx = context.WithValue(ctx, parent, "first")
		firstRun = true
		return handler(ctx, req)
	}
	second := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		require.Equal(t, serverInfo, info)
		require.Equal(t, "first", ctx.Value(parent).(string))
		ctx = context.WithValue(ctx, parent, "second")
		secondRun = true
		return handler(ctx, req)
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		require.Equal(t, "second", ctx.Value(parent).(string))
		handlerRun = true
		return out, nil
	}
	interceptors := ChainUnaryServerInterceptors(first, second)
	result, _ := interceptors(ctx, "req", serverInfo, handler)
	assert.Equal(t, out, result)
	assert.True(t, firstRun)
	assert.True(t, secondRun)
	assert.True(t, handlerRun)
}
