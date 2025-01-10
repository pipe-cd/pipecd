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
	plugins map[string]metadata
	stages  map[string]metadata
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
		plugins: make(map[string]metadata, 0),
		stages:  make(map[string]metadata, 0),
	}
	d := &model.Deployment{
		Metadata: map[string]string{
			"key-1": "value-1",
		},
		PluginMetadata: map[string]*model.KeyValues{
			"plugin-1": {
				KeyValues: map[string]string{
					"plugin-1-key-1": "plugin-1-value-1",
				},
			},
		},
		Stages: []*model.PipelineStage{
			{
				Id: "stage-1",
				Metadata: map[string]string{
					"stage-1-key-1": "stage-1-value-1",
				},
			},
		},
	}

	ctx := context.Background()
	store := newMetadataStore(ac, d)

	// Shared metadata.
	{
		// existing key
		value, found := store.sharedGet("key-1")
		assert.Equal(t, "value-1", value)
		assert.Equal(t, true, found)

		// nonexistent key
		value, found = store.sharedGet("key-2")
		assert.Equal(t, "", value)
		assert.Equal(t, false, found)
	}

	// Plugin metadata.
	{
		// existing key
		value, found := store.pluginGet("plugin-1", "plugin-1-key-1")
		assert.Equal(t, "plugin-1-value-1", value)
		assert.Equal(t, true, found)

		// nonexistent key
		value, found = store.pluginGet("plugin-1", "plugin-1-key-2")
		assert.Equal(t, "", value)
		assert.Equal(t, false, found)

		// put new and existing keys
		err := store.pluginPutMulti(ctx, "plugin-1", map[string]string{
			"plugin-1-key-2": "plugin-1-value-2",
			"plugin-1-key-1": "plugin-1-value-1-new",
			"plugin-1-key-3": "plugin-1-value-3",
		})
		assert.Equal(t, nil, err)
		value, found = store.pluginGet("plugin-1", "plugin-1-key-1")
		assert.Equal(t, "plugin-1-value-1-new", value)
		assert.Equal(t, true, found)
		value, found = store.pluginGet("plugin-1", "plugin-1-key-2")
		assert.Equal(t, "plugin-1-value-2", value)
		assert.Equal(t, true, found)
		value, found = store.pluginGet("plugin-1", "plugin-1-key-3")
		assert.Equal(t, "plugin-1-value-3", value)
		assert.Equal(t, true, found)

		assert.Equal(t, metadata{
			"plugin-1-key-1": "plugin-1-value-1-new",
			"plugin-1-key-2": "plugin-1-value-2",
			"plugin-1-key-3": "plugin-1-value-3",
		}, ac.plugins["plugin-1"])
	}

	// Stage metadata.
	{
		// existing key
		value, found := store.stageGet("stage-1", "stage-1-key-1")
		assert.Equal(t, "stage-1-value-1", value)
		assert.Equal(t, true, found)

		// nonexistent key
		value, found = store.stageGet("stage-1", "nonexistent-key")
		assert.Equal(t, "", value)
		assert.Equal(t, false, found)

		// put new and existing keys
		err := store.stagePutMulti(ctx, "stage-1", map[string]string{
			"stage-1-key-1": "stage-1-value-1-new",
			"stage-1-key-2": "stage-1-value-2",
		})
		assert.Equal(t, nil, err)
		value, found = store.stageGet("stage-1", "stage-1-key-1")
		assert.Equal(t, "stage-1-value-1-new", value)
		assert.Equal(t, true, found)
		value, found = store.stageGet("stage-1", "stage-1-key-2")
		assert.Equal(t, "stage-1-value-2", value)
		assert.Equal(t, true, found)

		assert.Equal(t, map[string]metadata{
			"stage-1": {
				"stage-1-key-1": "stage-1-value-1-new",
				"stage-1-key-2": "stage-1-value-2",
			},
		}, ac.stages)
	}
}
