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

package planner

import (
	"context"
	"fmt"

	"github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1/deployment"
	"github.com/pipe-cd/pipecd/pkg/regexpool"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type secretDecrypter interface {
	Decrypt(string) (string, error)
}

type PlannerService struct {
	deployment.UnimplementedDeploymentServiceServer

	Decrypter secretDecrypter
	RegexPool *regexpool.Pool
	Logger    *zap.Logger
}

// Register registers all handling of this service into the specified gRPC server.
func (a *PlannerService) Register(server *grpc.Server) {
	deployment.RegisterDeploymentServiceServer(server, a)
}

// NewPlannerService creates a new planService.
func NewPlannerService(
	decrypter secretDecrypter,
	logger *zap.Logger,
) *PlannerService {
	return &PlannerService{
		Decrypter: decrypter,
		RegexPool: regexpool.DefaultPool(),
		Logger:    logger.Named("planner"),
	}
}

func (ps *PlannerService) DetermineStrategy(ctx context.Context, in *deployment.DetermineStrategyRequest) (*deployment.DetermineStrategyResponse, error) {
	return nil, fmt.Errorf("not implemented yet")
}

func (ps *PlannerService) DetermineVersions(ctx context.Context, in *deployment.DetermineVersionsRequest) (*deployment.DetermineVersionsResponse, error) {
	return nil, fmt.Errorf("not implemented yet")
}
