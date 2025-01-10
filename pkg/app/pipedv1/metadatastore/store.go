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
	"sync"

	"google.golang.org/grpc"

	"github.com/pipe-cd/pipecd/pkg/app/server/service/pipedservice"
	"github.com/pipe-cd/pipecd/pkg/model"
)

type Getter interface {
	// Get finds and returns value of a given key.
	Get(key string) (string, bool)
}

type Putter interface {
	// Put adds a single key, value into store.
	// If the key is already existing, it overwrite the old value by the new one.
	Put(ctx context.Context, key, value string) error
	// PutMulti adds multiple (key, value) into store.
	// If any key is already existing, it overwrite the old value by the new one.
	PutMulti(ctx context.Context, md map[string]string) error
}

type Store interface {
	Getter
	Putter
}

type MetadataStore interface {
	Shared() Store
	Stage(stageID string) Store
}

type apiClient interface {
	SaveDeploymentSharedMetadata(ctx context.Context, req *pipedservice.SaveDeploymentSharedMetadataRequest, opts ...grpc.CallOption) (*pipedservice.SaveDeploymentSharedMetadataResponse, error)
	SaveDeploymentPluginMetadata(ctx context.Context, req *pipedservice.SaveDeploymentPluginMetadataRequest, opts ...grpc.CallOption) (*pipedservice.SaveDeploymentPluginMetadataResponse, error)
	SaveStageMetadata(ctx context.Context, req *pipedservice.SaveStageMetadataRequest, opts ...grpc.CallOption) (*pipedservice.SaveStageMetadataResponse, error)
}

type metadata map[string]string

type metadataStore struct {
	apiClient  apiClient
	deployment *model.Deployment

	shared  metadata
	plugins map[string]metadata
	stages  map[string]metadata

	sharedMu  sync.RWMutex
	pluginsMu sync.RWMutex
	stagesMu  sync.RWMutex
}

// newMetadataStore builds a metadata store for the given deployment.
// It keeps local data and makes sure that they are synced with the remote store.
func newMetadataStore(apiClient apiClient, d *model.Deployment) *metadataStore {
	s := &metadataStore{
		apiClient:  apiClient,
		deployment: d,
		shared:     make(map[string]string, 0),
		plugins:    make(map[string]metadata, 0),
		stages:     make(map[string]metadata, 0),
	}

	// Initialize shared metadata of deployment.
	for k, v := range d.GetMetadataV2().GetShared().GetKeyValues() {
		s.shared[k] = v
	}

	// Initialize metadata of plugins of the deployment.
	for _, pluginKV := range d.GetMetadataV2().GetPlugins() {
		s.plugins[pluginKV.String()] = pluginKV.GetKeyValues()
	}

	// Initialize metadata of all stages.
	for _, stage := range d.GetStages() {
		if md := stage.Metadata; md != nil {
			s.stages[stage.Id] = md
		}
	}
	return s
}

func (s *metadataStore) Shared() Store {
	return s
}

func (s *metadataStore) Get(key string) (value string, found bool) {
	s.sharedMu.RLock()
	defer s.sharedMu.RUnlock()

	value, found = s.shared[key]
	return
}

func (s *metadataStore) Put(ctx context.Context, key, value string) error {
	s.sharedMu.Lock()
	s.shared[key] = value
	s.sharedMu.Unlock()

	return s.syncSharedMetadata(ctx)
}

func (s *metadataStore) PutMulti(ctx context.Context, md map[string]string) error {
	s.sharedMu.Lock()
	for key, value := range md {
		s.shared[key] = value
	}
	s.sharedMu.Unlock()

	return s.syncSharedMetadata(ctx)
}

func (s *metadataStore) syncSharedMetadata(ctx context.Context) error {
	s.sharedMu.RLock()
	md := make(map[string]string, len(s.shared))
	for k, v := range s.shared {
		md[k] = v
	}
	s.sharedMu.RUnlock()

	// Send full list of metadata to ensure that they will be synced.
	_, err := s.apiClient.SaveDeploymentSharedMetadata(ctx, &pipedservice.SaveDeploymentSharedMetadataRequest{
		DeploymentId: s.deployment.Id,
		Metadata:     md,
	})
	return err
}

func (s *metadataStore) stagePutMulti(ctx context.Context, stageID string, md map[string]string) error {
	s.stagesMu.Lock()
	merged := make(map[string]string, len(md)+len(s.stages[stageID]))
	for k, v := range s.stages[stageID] {
		merged[k] = v
	}
	for k, v := range md {
		merged[k] = v
	}
	s.stages[stageID] = merged
	s.stagesMu.Unlock()

	// Send full list of metadata to ensure that they will be synced.
	_, err := s.apiClient.SaveStageMetadata(ctx, &pipedservice.SaveStageMetadataRequest{
		DeploymentId: s.deployment.Id,
		StageId:      stageID,
		Metadata:     merged,
	})
	return err
}

func (s *metadataStore) pluginGet(pluginName, key string) (value string, found bool) {
	s.pluginsMu.RLock()
	defer s.pluginsMu.RUnlock()

	md, ok := s.plugins[pluginName]
	if !ok {
		return "", false
	}

	value, found = md[key]
	return
}

func (s *metadataStore) pluginPutMulti(ctx context.Context, pluginName string, md map[string]string) error {
	s.pluginsMu.Lock()
	merged := make(map[string]string, len(md)+len(s.plugins[pluginName]))
	for k, v := range s.plugins[pluginName] {
		merged[k] = v
	}
	for k, v := range md {
		merged[k] = v
	}
	s.plugins[pluginName] = merged
	s.pluginsMu.Unlock()

	// Persist to the remote store.
	_, err := s.apiClient.SaveDeploymentPluginMetadata(ctx, &pipedservice.SaveDeploymentPluginMetadataRequest{
		DeploymentId: s.deployment.Id,
		PluginName:   pluginName,
		Metadata:     merged,
	})
	return err
}

func (s *metadataStore) stageGet(stageID, key string) (value string, found bool) {
	s.stagesMu.RLock()
	defer s.stagesMu.RUnlock()

	md, ok := s.stages[stageID]
	if !ok {
		return "", false
	}

	value, found = md[key]
	return
}

func (s *metadataStore) Stage(stageID string) Store {
	return &stageMetadataStore{backend: s, stageID: stageID}
}

type stageMetadataStore struct {
	stageID string
	backend *metadataStore
}

func (s *stageMetadataStore) PutMulti(ctx context.Context, md map[string]string) error {
	return s.backend.stagePutMulti(ctx, s.stageID, md)
}

func (s *stageMetadataStore) Put(ctx context.Context, key, value string) error {
	return s.backend.stagePutMulti(ctx, s.stageID, map[string]string{key: value})
}

func (s *stageMetadataStore) Get(key string) (string, bool) {
	return s.backend.stageGet(s.stageID, key)
}
