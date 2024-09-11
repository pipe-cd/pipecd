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
	"hash/fnv"

	"github.com/pipe-cd/pipecd/pkg/model"
	"go.opentelemetry.io/otel/trace"
)

func newContextWithDeploymentSpan(ctx context.Context, deployment *model.Deployment) context.Context {
	return trace.ContextWithSpanContext(ctx, trace.NewSpanContext(trace.SpanContextConfig{TraceID: deploymentTraceID(deployment)}))
}

func deploymentTraceID(deployment *model.Deployment) trace.TraceID {
	w := fnv.New128a()
	w.Write([]byte(deployment.Id))

	var id trace.TraceID
	copy(id[:], w.Sum(nil))
	return id
}
