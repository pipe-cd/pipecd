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

	api "github.com/pipe-cd/pipecd/pkg/plugin/pipedservice"
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

func (r *MetadataStoreRegistry) GetStageMetadata(ctx context.Context, req *api.GetStageMetadataRequest) (*api.GetStageMetadataResponse, error) {
	mds, ok := r.stores[req.DeploymentId]
	if !ok {
		return &api.GetStageMetadataResponse{Found: false}, fmt.Errorf("metadata store not found for deployment %s", req.DeploymentId)
	}

	value, found := mds.Stage(req.StageId).Get(req.Key)
	return &api.GetStageMetadataResponse{
		Value: value,
		Found: found,
	}, nil
}

func (r *MetadataStoreRegistry) PutStageMetadata(ctx context.Context, req *api.PutStageMetadataRequest) (*api.PutStageMetadataResponse, error) {
	mds, ok := r.stores[req.DeploymentId]
	if !ok {
		return &api.PutStageMetadataResponse{}, fmt.Errorf("metadata store not found for deployment %s", req.DeploymentId)
	}

	err := mds.Stage(req.StageId).Put(ctx, req.Key, req.Value)
	if err != nil {
		return &api.PutStageMetadataResponse{}, err
	}

	return &api.PutStageMetadataResponse{}, nil
}

func (r *MetadataStoreRegistry) PutStageMetadataMulti(ctx context.Context, req *api.PutStageMetadataMultiRequest) (*api.PutStageMetadataMultiResponse, error) {
	mds, ok := r.stores[req.DeploymentId]
	if !ok {
		return &api.PutStageMetadataMultiResponse{}, fmt.Errorf("metadata store not found for deployment %s", req.DeploymentId)
	}

	err := mds.Stage(req.StageId).PutMulti(ctx, req.Metadata)
	if err != nil {
		return &api.PutStageMetadataMultiResponse{}, err
	}

	return &api.PutStageMetadataMultiResponse{}, nil
}

func (r *MetadataStoreRegistry) GetDeploymentMetadata(ctx context.Context, req *api.GetDeploymentMetadataRequest) (*api.GetDeploymentMetadataResponse, error) {
	mds, ok := r.stores[req.DeploymentId]
	if !ok {
		return &api.GetDeploymentMetadataResponse{Found: false}, fmt.Errorf("metadata store not found for deployment %s", req.DeploymentId)
	}

	value, found := mds.Shared().Get(req.Key)
	return &api.GetDeploymentMetadataResponse{
		Value: value,
		Found: found,
	}, nil
}

func (r *MetadataStoreRegistry) PutDeploymentMetadata(ctx context.Context, req *api.PutDeploymentMetadataRequest) (*api.PutDeploymentMetadataResponse, error) {
	mds, ok := r.stores[req.DeploymentId]
	if !ok {
		return &api.PutDeploymentMetadataResponse{}, fmt.Errorf("metadata store not found for deployment %s", req.DeploymentId)
	}

	err := mds.Shared().Put(ctx, req.Key, req.Value)
	if err != nil {
		return &api.PutDeploymentMetadataResponse{}, err
	}

	return &api.PutDeploymentMetadataResponse{}, nil
}

func (r *MetadataStoreRegistry) PutDeploymentMetadataMulti(ctx context.Context, req *api.PutDeploymentMetadataMultiRequest) (*api.PutDeploymentMetadataMultiResponse, error) {
	mds, ok := r.stores[req.DeploymentId]
	if !ok {
		return &api.PutDeploymentMetadataMultiResponse{}, fmt.Errorf("metadata store not found for deployment %s", req.DeploymentId)
	}

	err := mds.Shared().PutMulti(ctx, req.Metadata)
	if err != nil {
		return &api.PutDeploymentMetadataMultiResponse{}, err
	}

	return &api.PutDeploymentMetadataMultiResponse{}, nil
}
