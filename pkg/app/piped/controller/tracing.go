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

package controller

import (
	"context"
	"crypto/rand"
	"hash/fnv"

	"github.com/pipe-cd/pipecd/pkg/model"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

var _ (sdktrace.IDGenerator) = (*tracingIDGenerator)(nil)

func newContextWithDeploymentSpan(ctx context.Context, deployment *model.Deployment) context.Context {
	gen := &tracingIDGenerator{deployment: deployment}
	return trace.ContextWithSpanContext(ctx, trace.NewSpanContext(trace.SpanContextConfig{TraceID: gen.deploymentTraceID()}))
}

type tracingIDGenerator struct {
	deployment *model.Deployment
}

func (gen *tracingIDGenerator) spanContextConfig() trace.SpanContextConfig {
	tid, sid := gen.NewIDs(context.Background())
	return trace.SpanContextConfig{
		TraceID:    tid,
		SpanID:     sid,
		TraceFlags: trace.FlagsSampled,
	}
}

// NewIDs implements trace.IDGenerator.
func (gen *tracingIDGenerator) NewIDs(ctx context.Context) (trace.TraceID, trace.SpanID) {
	tid := gen.deploymentTraceID()
	return tid, gen.NewSpanID(ctx, tid)
}

// NewSpanID implements trace.IDGenerator.
func (gen *tracingIDGenerator) NewSpanID(ctx context.Context, traceID trace.TraceID) trace.SpanID {
	var sid trace.SpanID
	_, _ = rand.Read(sid[:])
	return sid
}

func (gen *tracingIDGenerator) deploymentTraceID() trace.TraceID {
	w := fnv.New128a()
	w.Write([]byte(gen.deployment.Id))

	var id trace.TraceID
	copy(id[:], w.Sum(nil))
	return id
}
