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
	"encoding/json"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/pipe-cd/pipecd/pkg/model"
	"github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1/livestate"
)

var (
	livestateServiceServer interface {
		Plugin

		Register(server *grpc.Server)
		setFields(commonFields) error
		livestate.LivestateServiceServer
	}
)

// LivestatePlugin is the interface that must be implemented by a Livestate plugin.
// In addition to the Plugin interface, it provides a method to get the live state of the resources.
// The Config and DeployTargetConfig are the plugin's config defined in piped's config.
type LivestatePlugin[Config, DeployTargetConfig any] interface {
	Plugin

	// GetLivestate returns the live state of the resources in the given application.
	// It returns the resources' live state and the difference between the desired state and the live state.
	// It's allowed to return only the resources' live state if the difference is not available, or only the difference if the live state is not available.
	GetLivestate(context.Context, *Config, []*DeployTarget[DeployTargetConfig], *GetLivestateInput) (*GetLivestateResponse, error)
}

// LivestatePluginServer is a wrapper for LivestatePlugin to satisfy the LivestateServiceServer interface.
// It is used to register the plugin to the gRPC server.
type LivestatePluginServer[Config, DeployTargetConfig any] struct {
	livestate.UnimplementedLivestateServiceServer
	commonFields

	base          LivestatePlugin[Config, DeployTargetConfig]
	config        Config
	deployTargets map[string]*DeployTarget[DeployTargetConfig]
}

// RegisterLivestatePlugin registers the given LivestatePlugin to the sdk.
func RegisterLivestatePlugin[Config, DeployTargetConfig any](plugin LivestatePlugin[Config, DeployTargetConfig]) {
	livestateServiceServer = &LivestatePluginServer[Config, DeployTargetConfig]{base: plugin}
}

// Name returns the name of the plugin.
func (s *LivestatePluginServer[Config, DeployTargetConfig]) Name() string {
	return s.base.Name()
}

// Version returns the version of the plugin.
func (s *LivestatePluginServer[Config, DeployTargetConfig]) Version() string {
	return s.base.Version()
}

// Register registers the plugin to the gRPC server.
func (s *LivestatePluginServer[Config, DeployTargetConfig]) Register(server *grpc.Server) {
	livestate.RegisterLivestateServiceServer(server, s)
}

// setFields sets the common fields and configs to the server.
func (s *LivestatePluginServer[Config, DeployTargetConfig]) setFields(fields commonFields) error {
	s.commonFields = fields

	cfg := fields.config
	if cfg.Config != nil {
		if err := json.Unmarshal(cfg.Config, &s.config); err != nil {
			s.logger.Fatal("failed to unmarshal the plugin config", zap.Error(err))
			return err
		}
	}

	s.deployTargets = make(map[string]*DeployTarget[DeployTargetConfig], len(cfg.DeployTargets))
	for _, dt := range cfg.DeployTargets {
		var sdkDt DeployTargetConfig
		if err := json.Unmarshal(dt.Config, &sdkDt); err != nil {
			s.logger.Fatal("failed to unmarshal deploy target config", zap.Error(err))
			return err
		}
		s.deployTargets[dt.Name] = &DeployTarget[DeployTargetConfig]{
			Name:   dt.Name,
			Labels: dt.Labels,
			Config: sdkDt,
		}
	}

	return nil
}

// GetLivestate returns the live state of the resources in the given application.
func (s *LivestatePluginServer[Config, DeployTargetConfig]) GetLivestate(ctx context.Context, request *livestate.GetLivestateRequest) (*livestate.GetLivestateResponse, error) {
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
		pluginName:    s.Name(),
		applicationID: request.GetApplicationId(),
		toolRegistry:  s.toolRegistry,
	}

	response, err := s.base.GetLivestate(ctx, &s.config, deployTargets, &GetLivestateInput{
		Request: GetLivestateRequest{
			ApplicationID:    request.ApplicationId,
			DeploymentSource: newDeploymentSource(request.GetDeploySource()),
		},
		Client: client,
		Logger: s.logger,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get the live state: %v", err)
	}

	return response.toModel(time.Now()), nil
}

type GetLivestateInput struct {
	Request GetLivestateRequest
	Client  *Client
	Logger  *zap.Logger
}

type GetLivestateRequest struct {
	// ApplicationID is the ID of the application.
	ApplicationID string
	// DeploymentSource is the source of the deployment.
	DeploymentSource DeploymentSource
}

type GetLivestateResponse struct {
	LiveState ApplicationLiveState
	SyncState ApplicationSyncState
}

func (r *GetLivestateResponse) toModel(now time.Time) *livestate.GetLivestateResponse {
	return &livestate.GetLivestateResponse{
		ApplicationLiveState: r.LiveState.toModel(now),
		SyncState:            r.SyncState.toModel(now),
	}
}

// ApplicationLiveState represents the live state of an application.
type ApplicationLiveState struct {
	Resources    []ResourceState
	HealthStatus ApplicationHealthStatus
}

// toModel converts the ApplicationLiveState to the model.ApplicationLiveState.
func (s *ApplicationLiveState) toModel(now time.Time) *model.ApplicationLiveState {
	resources := make([]*model.ResourceState, 0, len(s.Resources))
	for _, rs := range s.Resources {
		resources = append(resources, rs.toModel(now))
	}
	return &model.ApplicationLiveState{
		Resources:    resources,
		HealthStatus: s.HealthStatus.toModel(),
	}
}

// ResourceState represents the live state of a resource.
type ResourceState struct {
	// ID is the unique identifier of the resource.
	ID string
	// ParentIDs is the list of the parent resource's IDs.
	ParentIDs []string
	// Name is the name of the resource.
	Name string
	// ResourceType is the type of the resource.
	ResourceType string
	// ResourceMetadata is the metadata of the resource.
	ResourceMetadata map[string]string
	// HealthStatus is the health status of the resource.
	HealthStatus ResourceHealthStatus
	// HealthDescription is the description of the health status.
	HealthDescription string
	// DeployTarget is the target where the resource is deployed.
	DeployTarget string
	// PluginName is the name of the plugin that provides the resource.
	PluginName string
	// CreatedAt is the time when the resource was created.
	CreatedAt time.Time
}

// toModel converts the ResourceState to the model.ResourceState.
func (s *ResourceState) toModel(now time.Time) *model.ResourceState {
	return &model.ResourceState{
		Id:                s.ID,
		ParentIds:         s.ParentIDs,
		Name:              s.Name,
		ResourceType:      s.ResourceType,
		ResourceMetadata:  s.ResourceMetadata,
		HealthStatus:      s.HealthStatus.toModel(),
		HealthDescription: s.HealthDescription,
		DeployTarget:      s.DeployTarget,
		PluginName:        s.PluginName,
		CreatedAt:         s.CreatedAt.Unix(),
		UpdatedAt:         now.Unix(),
	}
}

// ApplicationHealthStatus represents the health status of an application.
type ApplicationHealthStatus int

const (
	// ApplicationHealthStateUnknown represents the unknown health status of an application.
	ApplicationHealthStateUnknown ApplicationHealthStatus = iota
	// ApplicationHealthStateHealthy represents the healthy health status of an application.
	ApplicationHealthStateHealthy
	// ApplicationHealthStateOther represents the other health status of an application.
	ApplicationHealthStateOther
)

// toModel converts the ApplicationHealthStatus to the model.ApplicationLiveState_Status.
func (s ApplicationHealthStatus) toModel() model.ApplicationLiveState_Status {
	switch s {
	case ApplicationHealthStateHealthy:
		return model.ApplicationLiveState_HEALTHY
	case ApplicationHealthStateOther:
		return model.ApplicationLiveState_OTHER
	default:
		return model.ApplicationLiveState_UNKNOWN
	}
}

// ResourceHealthStatus represents the health status of a resource.
type ResourceHealthStatus int

const (
	// ResourceHealthStateUnknown represents the unknown health status of a resource.
	ResourceHealthStateUnknown ResourceHealthStatus = iota
	// ResourceHealthStateHealthy represents the healthy health status of a resource.
	ResourceHealthStateHealthy
	// ResourceHealthStateUnhealthy represents the unhealthy health status of a resource.
	ResourceHealthStateUnhealthy
)

// toModel converts the ResourceHealthStatus to the model.ResourceState_HealthStatus.
func (s ResourceHealthStatus) toModel() model.ResourceState_HealthStatus {
	switch s {
	case ResourceHealthStateHealthy:
		return model.ResourceState_HEALTHY
	case ResourceHealthStateUnhealthy:
		return model.ResourceState_UNHEALTHY
	default:
		return model.ResourceState_UNKNOWN
	}
}

type ApplicationSyncState struct {
}

func (s *ApplicationSyncState) toModel(now time.Time) *model.ApplicationSyncState {
	return nil
}
