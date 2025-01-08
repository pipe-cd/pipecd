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

package metadatastore

import (
	"context"
	"fmt"

	"github.com/pipe-cd/pipecd/pkg/model"
	service "github.com/pipe-cd/pipecd/pkg/plugin/pipedservice"
)

// MetadataStoreRegistry is a registry of metadata stores for deployments.
type MetadataStoreRegistry struct {
	apiClient apiClient

	// stores is a map of metadata store for each deployment.
	// The key is the deployment ID.
	stores map[string]MetadataStore
}

// NewMetadataStoreRegistry creates a new MetadataStoreRegistry.
func NewMetadataStoreRegistry(apiClient apiClient) *MetadataStoreRegistry {
	return &MetadataStoreRegistry{apiClient: apiClient, stores: make(map[string]MetadataStore, 0)}
}

// Register creates a new metadata store for the given deployment.
// This must be called before other Get/Put methods are called for the deployment.
// If the metadata store already exists, the new one will replace the existing one.
func (r *MetadataStoreRegistry) Register(d *model.Deployment) {
	store := NewMetadataStore(r.apiClient, d)
	r.stores[d.Id] = store
}

// Delete deletes the metadata store for the given deployment in order to release the resources.
// If the metadata store is not found, it is a no-op.
func (r *MetadataStoreRegistry) Delete(deploymentID string) {
	delete(r.stores, deploymentID)
}

// GetStageMetadata implements the backend of PluginService.GetStageMetadata().
func (r *MetadataStoreRegistry) GetStageMetadata(ctx context.Context, req *service.GetStageMetadataRequest) (*service.GetStageMetadataResponse, error) {
	mds, ok := r.stores[req.DeploymentId]
	if !ok {
		return &service.GetStageMetadataResponse{Found: false}, fmt.Errorf("metadata store not found for deployment %s", req.DeploymentId)
	}

	value, found := mds.Stage(req.StageId).Get(req.Key)
	return &service.GetStageMetadataResponse{
		Value: value,
		Found: found,
	}, nil
}

// PutStageMetadata implements the backend of PluginService.PutStageMetadata().
func (r *MetadataStoreRegistry) PutStageMetadata(ctx context.Context, req *service.PutStageMetadataRequest) (*service.PutStageMetadataResponse, error) {
	mds, ok := r.stores[req.DeploymentId]
	if !ok {
		return &service.PutStageMetadataResponse{}, fmt.Errorf("metadata store not found for deployment %s", req.DeploymentId)
	}

	err := mds.Stage(req.StageId).Put(ctx, req.Key, req.Value)
	if err != nil {
		return &service.PutStageMetadataResponse{}, err
	}

	return &service.PutStageMetadataResponse{}, nil
}

// PutStageMetadataMulti implements the backend of PluginService.PutStageMetadataMulti().
func (r *MetadataStoreRegistry) PutStageMetadataMulti(ctx context.Context, req *service.PutStageMetadataMultiRequest) (*service.PutStageMetadataMultiResponse, error) {
	mds, ok := r.stores[req.DeploymentId]
	if !ok {
		return &service.PutStageMetadataMultiResponse{}, fmt.Errorf("metadata store not found for deployment %s", req.DeploymentId)
	}

	err := mds.Stage(req.StageId).PutMulti(ctx, req.Metadata)
	if err != nil {
		return &service.PutStageMetadataMultiResponse{}, err
	}

	return &service.PutStageMetadataMultiResponse{}, nil
}

// GetDeploymentMetadata implements the backend of PluginService.GetDeploymentMetadata().
func (r *MetadataStoreRegistry) GetDeploymentMetadata(ctx context.Context, req *service.GetDeploymentMetadataRequest) (*service.GetDeploymentMetadataResponse, error) {
	mds, ok := r.stores[req.DeploymentId]
	if !ok {
		return &service.GetDeploymentMetadataResponse{Found: false}, fmt.Errorf("metadata store not found for deployment %s", req.DeploymentId)
	}

	value, found := mds.Shared().Get(req.Key)
	return &service.GetDeploymentMetadataResponse{
		Value: value,
		Found: found,
	}, nil
}

// PutDeploymentMetadata implements the backend of PluginService.PutDeploymentMetadata().
func (r *MetadataStoreRegistry) PutDeploymentMetadata(ctx context.Context, req *service.PutDeploymentMetadataRequest) (*service.PutDeploymentMetadataResponse, error) {
	mds, ok := r.stores[req.DeploymentId]
	if !ok {
		return &service.PutDeploymentMetadataResponse{}, fmt.Errorf("metadata store not found for deployment %s", req.DeploymentId)
	}

	err := mds.Shared().Put(ctx, req.Key, req.Value)
	if err != nil {
		return &service.PutDeploymentMetadataResponse{}, err
	}

	return &service.PutDeploymentMetadataResponse{}, nil
}

// PutDeploymentMetadataMulti implements the backend of PluginService.PutDeploymentMetadataMulti().
func (r *MetadataStoreRegistry) PutDeploymentMetadataMulti(ctx context.Context, req *service.PutDeploymentMetadataMultiRequest) (*service.PutDeploymentMetadataMultiResponse, error) {
	mds, ok := r.stores[req.DeploymentId]
	if !ok {
		return &service.PutDeploymentMetadataMultiResponse{}, fmt.Errorf("metadata store not found for deployment %s", req.DeploymentId)
	}

	err := mds.Shared().PutMulti(ctx, req.Metadata)
	if err != nil {
		return &service.PutDeploymentMetadataMultiResponse{}, err
	}

	return &service.PutDeploymentMetadataMultiResponse{}, nil
}
