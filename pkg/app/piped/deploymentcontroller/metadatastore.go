// Copyright 2020 The PipeCD Authors.
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

package deploymentcontroller

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/kapetaniosci/pipe/pkg/app/api/service/pipedservice"
	"github.com/kapetaniosci/pipe/pkg/model"
)

type metadataStore struct {
	apiClient     apiClient
	deployment    *model.Deployment
	metadata      sync.Map // map[key-string]value-string
	stageMetadata sync.Map // map[stage-id-string]json-string
}

func NewMetadataStore(apiClient apiClient, d *model.Deployment) *metadataStore {
	s := &metadataStore{
		apiClient:     apiClient,
		deployment:    d,
		metadata:      sync.Map{},
		stageMetadata: sync.Map{},
	}
	// Store shared metadata of deployment.
	for k, v := range d.Metadata {
		s.metadata.Store(k, v)
	}
	// Store metadata of all stages.
	for _, stage := range d.Stages {
		s.stageMetadata.Store(stage.Id, stage.JsonMetadata)
	}
	return s
}

func (s *metadataStore) Set(ctx context.Context, key, value string) error {
	s.metadata.Store(key, value)

	metadata := make(map[string]string)
	s.metadata.Range(func(key, value interface{}) bool {
		var (
			k = key.(string)
			v = value.(string)
		)
		metadata[k] = v
		return true
	})

	_, err := s.apiClient.SaveDeploymentMetadata(ctx, &pipedservice.SaveDeploymentMetadataRequest{
		DeploymentId: s.deployment.Id,
		Metadata:     metadata,
	})
	return err
}

func (s *metadataStore) Get(key string) (string, bool) {
	if value, ok := s.metadata.Load(key); ok {
		return value.(string), ok
	}
	return "", false
}

func (s *metadataStore) SetStageMetadata(ctx context.Context, stageID string, metadata interface{}) error {
	data, err := json.Marshal(metadata)
	if err != nil {
		return err
	}
	s.stageMetadata.Store(stageID, string(data))

	_, err = s.apiClient.SaveStageMetadata(ctx, &pipedservice.SaveStageMetadataRequest{
		DeploymentId: s.deployment.Id,
		StageId:      stageID,
		JsonMetadata: string(data),
	})
	return err
}

func (s *metadataStore) GetStageMetadata(stageID string, metadata interface{}) error {
	var data string
	if v, ok := s.stageMetadata.Load(stageID); ok {
		data = v.(string)
	}
	return json.Unmarshal([]byte(data), metadata)
}
