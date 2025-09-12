// Copyright 2025 The PipeCD Authors.
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

package sdk

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1/planpreview"
)

// PlanPreviewPlugin is the interface that must be implemented by a PlanPreview plugin.
// In addition to the Plugin interface, it provides a method to get the plan preview result.
// The Config and DeployTargetConfig are the plugin's config defined in piped's config.
type PlanPreviewPlugin[Config, DeployTargetConfig, ApplicationConfigSpec any] interface {
	// GetPlanPreview returns the plan preview result of the given application.
	GetPlanPreview(context.Context, *Config, []*DeployTarget[DeployTargetConfig], *GetPlanPreviewInput[ApplicationConfigSpec]) (*GetPlanPreviewResponse, error)
}

// PlanPreviewPluginServer is a wrapper for PlanPreviewPlugin to satisfy the PlanPreviewServiceServer interface.
// It is used to register the plugin to the gRPC server.
type PlanPreviewPluginServer[Config, DeployTargetConfig, ApplicationConfigSpec any] struct {
	planpreview.UnimplementedPlanPreviewServiceServer
	commonFields[Config, DeployTargetConfig]

	base PlanPreviewPlugin[Config, DeployTargetConfig, ApplicationConfigSpec]
}

// Register registers the plugin to the gRPC server.
func (s *PlanPreviewPluginServer[Config, DeployTargetConfig, ApplicationConfigSpec]) Register(server *grpc.Server) {
	planpreview.RegisterPlanPreviewServiceServer(server, s)
}

// GetPlanPreview returns the plan preview of the resources in the given application.
func (s *PlanPreviewPluginServer[Config, DeployTargetConfig, ApplicationConfigSpec]) GetPlanPreview(ctx context.Context, request *planpreview.GetPlanPreviewRequest) (*planpreview.GetPlanPreviewResponse, error) {
	// Get the deploy targets set on the deployment from the piped plugin config.
	deployTargets := make([]*DeployTarget[DeployTargetConfig], 0, len(request.GetDeployTargets()))
	for _, name := range request.GetDeployTargets() {
		dt, ok := s.deployTargets[name]
		if !ok {
			return nil, status.Errorf(codes.Internal, "the deploy target %s is not found in the piped plugin config", name)
		}

		deployTargets = append(deployTargets, dt)
	}

	client := &Client{
		base:          s.client,
		pluginName:    s.name,
		applicationID: request.GetApplicationId(),
		toolRegistry:  s.toolRegistry,
	}

	targetDS, err := newDeploymentSource[ApplicationConfigSpec](s.name, request.GetTargetDeploymentSource())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to parse target deployment source: %v", err)
	}

	runningDS, err := newDeploymentSource[ApplicationConfigSpec](s.name, request.GetRunningDeploymentSource())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to parse running deployment source: %v", err)
	}

	response, err := s.base.GetPlanPreview(ctx, s.pluginConfig, deployTargets, &GetPlanPreviewInput[ApplicationConfigSpec]{
		Request: GetPlanPreviewRequest[ApplicationConfigSpec]{
			ApplicationID:           request.GetApplicationId(),
			ApplicationName:         request.GetApplicationName(),
			PipedID:                 request.GetPipedId(),
			DeployTargets:           request.GetDeployTargets(),
			TargetDeploymentSource:  targetDS,
			RunningDeploymentSource: runningDS,
		},
		Client: client,
		Logger: s.logger,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get the plan preview: %v", err)
	}

	return response.toProto(), nil
}

// GetPlanPreviewInput is the input for the GetPlanPreview method.
type GetPlanPreviewInput[ApplicationConfigSpec any] struct {
	// Request is the request for getting the plan preview.
	Request GetPlanPreviewRequest[ApplicationConfigSpec]
	// Client is the client for accessing the piped API.
	Client *Client
	// Logger is the logger for logging.
	Logger *zap.Logger
}

// GetPlanPreviewRequest is the request for the GetPlanPreview method.
type GetPlanPreviewRequest[ApplicationConfigSpec any] struct {
	// ApplicationID is the ID of the application.
	ApplicationID string
	// ApplicationName is the name of the application.
	ApplicationName string
	// PipedID is the ID of the piped.
	PipedID string
	// DeployTargets is the names of the deploy targets.
	DeployTargets []string
	// TargetDeploymentSource is the target source of the deployment.
	TargetDeploymentSource DeploymentSource[ApplicationConfigSpec]
	// RunningDeploymentSource is the running source of the deployment.
	RunningDeploymentSource DeploymentSource[ApplicationConfigSpec]
}

// GetPlanPreviewResponse is the response for the GetPlanPreview method.
type GetPlanPreviewResponse struct {
	// Results is the results of the plan preview.
	Results []PlanPreviewResult
}

// PlanPreviewResult is the result of the plan preview.
type PlanPreviewResult struct {
	// DeployTarget is the name of the deploy target.
	DeployTarget string
	// Summary is a human-readable summary of the plan preview.
	Summary string
	// NoChange indicates whether any changes were detected.
	NoChange bool
	// Details contains the detailed plan preview information.
	Details []byte
	// DiffLanguage is the language to render the details like "diff","hcl".
	// If this is empty, "diff" will be used by default.
	DiffLanguage string
}

// toProto converts the GetPlanPreviewResponse to the planpreview.GetPlanPreviewResponse.
func (r *GetPlanPreviewResponse) toProto() *planpreview.GetPlanPreviewResponse {
	results := make([]*planpreview.PlanPreviewResult, 0, len(r.Results))
	for _, result := range r.Results {
		results = append(results, &planpreview.PlanPreviewResult{
			DeployTarget: result.DeployTarget,
			Summary:      result.Summary,
			NoChange:     result.NoChange,
			Details:      result.Details,
			DiffLanguage: result.DiffLanguage,
		})
	}

	return &planpreview.GetPlanPreviewResponse{
		Results: results,
	}
}
