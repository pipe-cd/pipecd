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

package api

import (
	"context"
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/pipe-cd/pipecd/pkg/app/helloworld/service"
)

var (
	helloCalls = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "hello_calls_total",
			Help: "Number of hello calls.",
		},
		[]string{"gender"},
	)
)

func init() {
	prometheus.MustRegister(helloCalls)
}

type api struct {
	service.UnimplementedHelloWorldServer

	logger *zap.Logger
}

type Option func(*api)

func WithLogger(logger *zap.Logger) Option {
	return func(a *api) {
		a.logger = logger.Named("api")
	}
}

// NewHelloWorldAPI creates new instance for api.
func NewHelloWorldAPI(opts ...Option) *api {
	a := &api{
		logger: zap.NewNop(),
	}
	for _, opt := range opts {
		opt(a)
	}
	return a
}

func (a *api) Register(server *grpc.Server) {
	service.RegisterHelloWorldServer(server, a)
}

func (a *api) Hello(ctx context.Context, in *service.HelloRequest) (*service.HelloResponse, error) {
	m := "mr"
	if in.Gender == service.HelloRequest_GENDER_FEMALE {
		m = "mrs"
	}

	helloCalls.With(prometheus.Labels{"gender": m}).Inc()

	return &service.HelloResponse{
		Message: fmt.Sprintf("Hello, %s %s", m, in.Name),
	}, nil
}
