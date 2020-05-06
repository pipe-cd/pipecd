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

	"github.com/kapetaniosci/pipe/pkg/app/api/service/runnerservice"
)

type metadataPersister struct {
	apiClient apiClient
}

func (p metadataPersister) StageMetadataPersister(deploymentID, stageID string) stageMetadataPersister {
	return stageMetadataPersister{
		deploymentID: deploymentID,
		stageID:      stageID,
		apiClient:    p.apiClient,
	}
}

type stageMetadataPersister struct {
	deploymentID string
	stageID      string
	apiClient    apiClient
}

func (p stageMetadataPersister) Save(ctx context.Context, metadata []byte) error {
	_, err := p.apiClient.SaveStageMetadata(ctx, &runnerservice.SaveStageMetadataRequest{
		Id:           p.stageID,
		DeploymentId: p.deploymentID,
		Metadata:     metadata,
	})
	return err
}
