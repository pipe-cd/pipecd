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
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"

	"github.com/pipe-cd/pipecd/pkg/app/server/service/pipedservice"
	"github.com/pipe-cd/pipecd/pkg/model"
)

type fakeAPIClient struct {
	shared  metadata
	plugins map[string]metadata
	stages  map[string]metadata
}

func (c *fakeAPIClient) SaveDeploymentSharedMetadata(ctx context.Context, req *pipedservice.SaveDeploymentSharedMetadataRequest, opts ...grpc.CallOption) (*pipedservice.SaveDeploymentSharedMetadataResponse, error) {
	md := make(map[string]string, len(c.shared)+len(req.Metadata))
	for k, v := range c.shared {
		md[k] = v
	}
	for k, v := range req.Metadata {
		md[k] = v
	}
	c.shared = md
	return &pipedservice.SaveDeploymentSharedMetadataResponse{}, nil
}

func (c *fakeAPIClient) SaveDeploymentPluginMetadata(ctx context.Context, req *pipedservice.SaveDeploymentPluginMetadataRequest, opts ...grpc.CallOption) (*pipedservice.SaveDeploymentPluginMetadataResponse, error) {
	ori := c.plugins[req.PluginName]
	md := make(map[string]string, len(ori)+len(req.Metadata))
	for k, v := range ori {
		md[k] = v
	}
	for k, v := range req.Metadata {
		md[k] = v
	}
	c.plugins[req.PluginName] = md
	return &pipedservice.SaveDeploymentPluginMetadataResponse{}, nil
}

func (c *fakeAPIClient) SaveStageMetadata(ctx context.Context, req *pipedservice.SaveStageMetadataRequest, opts ...grpc.CallOption) (*pipedservice.SaveStageMetadataResponse, error) {
	ori := c.stages[req.StageId]
	md := make(map[string]string, len(ori)+len(req.Metadata))
	for k, v := range ori {
		md[k] = v
	}
	for k, v := range req.Metadata {
		md[k] = v
	}
	c.stages[req.StageId] = md
	return &pipedservice.SaveStageMetadataResponse{}, nil
}

func TestStore(t *testing.T) {
	t.Parallel()

	ac := &fakeAPIClient{
		shared: make(map[string]string, 0),
		stages: make(map[string]metadata, 0),
	}
	d := &model.Deployment{
		Metadata: map[string]string{
			"key-1": "value-1",
		},
		Stages: []*model.PipelineStage{
			{
				Id: "stage-1",
			},
			{
				Id: "stage-2",
				Metadata: map[string]string{
					"stage-2-key-1": "stage-2-value-1",
				},
			},
		},
	}

	ctx := context.Background()
	store := newMetadataStore(ac, d)

	// Shared metadata.
	value, found := store.Shared().Get("key-1")
	assert.Equal(t, "value-1", value)
	assert.Equal(t, true, found)

	value, found = store.Shared().Get("key-2")
	assert.Equal(t, "", value)
	assert.Equal(t, false, found)

	err := store.Shared().Put(ctx, "key-2", "value-2")
	assert.Equal(t, nil, err)

	assert.Equal(t, metadata{
		"key-1": "value-1",
		"key-2": "value-2",
	}, ac.shared)

	err = store.Shared().PutMulti(ctx, map[string]string{
		"key-3": "value-3",
		"key-1": "value-12",
		"key-4": "value-4",
	})
	assert.Equal(t, nil, err)

	assert.Equal(t, metadata{
		"key-1": "value-12",
		"key-2": "value-2",
		"key-3": "value-3",
		"key-4": "value-4",
	}, ac.shared)

	// Stage metadata.
	value, found = store.Stage("stage-1").Get("key-1")
	assert.Equal(t, "", value)
	assert.Equal(t, false, found)

	value, found = store.Stage("stage-2").Get("stage-2-key-1")
	assert.Equal(t, "stage-2-value-1", value)
	assert.Equal(t, true, found)

	err = store.Stage("stage-1").Put(ctx, "stage-1-key-1", "stage-1-value-1")
	assert.Equal(t, nil, err)

	value, found = store.Stage("stage-1").Get("stage-1-key-1")
	assert.Equal(t, "stage-1-value-1", value)
	assert.Equal(t, true, found)

	err = store.Stage("stage-2").PutMulti(ctx, map[string]string{
		"stage-2-key-1": "stage-2-value-12",
		"stage-2-key-2": "stage-2-value-2",
	})
	assert.Equal(t, nil, err)

	value, found = store.Stage("stage-2").Get("stage-2-key-1")
	assert.Equal(t, "stage-2-value-12", value)
	assert.Equal(t, true, found)

	value, found = store.Stage("stage-2").Get("stage-2-key-2")
	assert.Equal(t, "stage-2-value-2", value)
	assert.Equal(t, true, found)

	assert.Equal(t, map[string]metadata{
		"stage-1": {
			"stage-1-key-1": "stage-1-value-1",
		},
		"stage-2": {
			"stage-2-key-1": "stage-2-value-12",
			"stage-2-key-2": "stage-2-value-2",
		},
	}, ac.stages)
}
