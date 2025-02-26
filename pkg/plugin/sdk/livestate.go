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
	"fmt"

	"github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1/livestate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	livestateServiceServer interface {
		Plugin

		Register(server *grpc.Server)
		setCommonFields(commonFields)
		setConfig([]byte) error
		livestate.LivestateServiceServer
	}
)

// LivestatePlugin is the interface that must be implemented by a Livestate plugin.
// In addition to the Plugin interface, it provides a method to get the live state of the resources.
// The Config and DeployTargetConfig are the plugin's config defined in piped's config.
type LivestatePlugin[Config, DeployTargetConfig any] interface {
	Plugin

	// GetLivestate returns the live state of the resources in the given deploy target.
	// It returns the resources' live state and the difference between the desired state and the live state.
	// It's allowed to return only the resources' live state if the difference is not available, or only the difference if the live state is not available.
	GetLivestate(context.Context, *Config, []*DeployTarget[DeployTargetConfig], TODO) (TODO, error)
}

// LivestatePluginServer is a wrapper for LivestatePlugin to satisfy the LivestateServiceServer interface.
// It is used to register the plugin to the gRPC server.
type LivestatePluginServer[Config, DeployTargetConfig any] struct {
	livestate.UnimplementedLivestateServiceServer
	commonFields

	base   LivestatePlugin[Config, DeployTargetConfig]
	config Config
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

// setCommonFields sets the common fields to the plugin server.
func (s *LivestatePluginServer[Config, DeployTargetConfig]) setCommonFields(c commonFields) {
	s.commonFields = c
}

// setConfig sets the configuration to the plugin server.
func (s *LivestatePluginServer[Config, DeployTargetConfig]) setConfig(bytes []byte) error {
	if bytes == nil {
		return nil
	}
	if err := json.Unmarshal(bytes, &s.config); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}
	return nil
}

// GetLivestate returns the live state of the resources in the given deploy target.
func (s *LivestatePluginServer[Config, DeployTargetConfig]) GetLivestate(context.Context, *livestate.GetLivestateRequest) (*livestate.GetLivestateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLivestate not implemented")
}
