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

	"github.com/pipe-cd/pipecd/pkg/model"
	service "github.com/pipe-cd/pipecd/pkg/plugin/pipedservice"
)

func TestRegistry(t *testing.T) {
	t.Parallel()

	ac := &fakeAPIClient{
		shared:  make(map[string]string, 0),
		plugins: make(map[string]metadata, 0),
		stages:  make(map[string]metadata, 0),
	}
	d := &model.Deployment{
		MetadataV2: &model.DeploymentMetadata{
			Shared: &model.DeploymentMetadata_KeyValues{
				KeyValues: map[string]string{
					"key-1": "value-1",
				},
			},
			Plugins: map[string]*model.DeploymentMetadata_KeyValues{
				"plugin-1": {
					KeyValues: map[string]string{
						"plugin-1-key-1": "plugin-1-value-1",
					},
				},
			},
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

	r := NewMetadataStoreRegistry(ac)
	r.Register(d)

	ctx := context.Background()

	// DeploymentShared metadata.
	{
		// Get
		{
			// Existing key
			resp, err := r.GetDeploymentSharedMetadata(ctx, &service.GetDeploymentSharedMetadataRequest{
				DeploymentId: d.Id,
				Key:          "key-1",
			})
			assert.NoError(t, err)
			assert.Equal(t, true, resp.Found)
			assert.Equal(t, "value-1", resp.Value)

			// Nonexistent key
			resp, err = r.GetDeploymentSharedMetadata(ctx, &service.GetDeploymentSharedMetadataRequest{
				DeploymentId: d.Id,
				Key:          "nonexistent-key",
			})
			assert.NoError(t, err)
			assert.Equal(t, false, resp.Found)
			assert.Equal(t, "", resp.Value)
		}
	}

	// DeploymentPlugin metadata.
	{
		// Get
		{
			// Existing key
			resp, err := r.GetDeploymentPluginMetadata(ctx, &service.GetDeploymentPluginMetadataRequest{
				DeploymentId: d.Id,
				PluginName:   "plugin-1",
				Key:          "plugin-1-key-1",
			})
			assert.NoError(t, err)
			assert.Equal(t, true, resp.Found)
			assert.Equal(t, "plugin-1-value-1", resp.Value)

			// Nonexistent key
			resp, err = r.GetDeploymentPluginMetadata(ctx, &service.GetDeploymentPluginMetadataRequest{
				DeploymentId: d.Id,
				PluginName:   "plugin-1",
				Key:          "nonexistent-key",
			})
			assert.NoError(t, err)
			assert.Equal(t, false, resp.Found)
			assert.Equal(t, "", resp.Value)

			// Nonexistent plugin
			resp, err = r.GetDeploymentPluginMetadata(ctx, &service.GetDeploymentPluginMetadataRequest{
				DeploymentId: d.Id,
				PluginName:   "nonexistent-plugin",
				Key:          "plugin-1-key-1",
			})
			assert.NoError(t, err)
			assert.Equal(t, false, resp.Found)
			assert.Equal(t, "", resp.Value)

			// Nonexistent deployment
			resp, err = r.GetDeploymentPluginMetadata(ctx, &service.GetDeploymentPluginMetadataRequest{
				DeploymentId: "not-exist-id",
				PluginName:   "plugin-1",
				Key:          "plugin-1-key-1",
			})
			assert.Error(t, err)
			assert.Equal(t, false, resp.Found)
			assert.Equal(t, "", resp.Value)
		}
		// Put
		{
			// New key
			_, err := r.PutDeploymentPluginMetadata(ctx, &service.PutDeploymentPluginMetadataRequest{
				DeploymentId: d.Id,
				PluginName:   "plugin-1",
				Key:          "plugin-1-key-2",
				Value:        "plugin-1-value-2",
			})
			assert.NoError(t, err)
			assert.Equal(t, metadata{
				"plugin-1-key-1": "plugin-1-value-1",
				"plugin-1-key-2": "plugin-1-value-2",
			}, ac.plugins["plugin-1"])

			// Nonexistent deployment
			_, err = r.PutDeploymentPluginMetadata(ctx, &service.PutDeploymentPluginMetadataRequest{
				DeploymentId: "not-exist-id",
				PluginName:   "plugin-1",
				Key:          "plugin-1-key-2",
				Value:        "plugin-1-value-2",
			})
			assert.Error(t, err)
		}
		// PutMulti
		{
			// New keys(3,4) with one existing key(1)
			_, err := r.PutDeploymentPluginMetadataMulti(ctx, &service.PutDeploymentPluginMetadataMultiRequest{
				DeploymentId: d.Id,
				PluginName:   "plugin-1",
				Metadata: map[string]string{
					"plugin-1-key-3": "plugin-1-value-3",
					"plugin-1-key-1": "plugin-1-value-1-new",
					"plugin-1-key-4": "plugin-1-value-4",
				},
			})
			assert.NoError(t, err)
			assert.Equal(t, metadata{
				"plugin-1-key-1": "plugin-1-value-1-new",
				"plugin-1-key-2": "plugin-1-value-2",
				"plugin-1-key-3": "plugin-1-value-3",
				"plugin-1-key-4": "plugin-1-value-4",
			}, ac.plugins["plugin-1"])

			// Nonexistent deployment
			_, err = r.PutDeploymentPluginMetadataMulti(ctx, &service.PutDeploymentPluginMetadataMultiRequest{
				DeploymentId: "nonexistent-id",
				PluginName:   "plugin-1",
				Metadata: map[string]string{
					"plugin-1-key-3": "plugin-1-value-3",
					"plugin-1-key-4": "plugin-1-value-4",
				},
			})
			assert.Error(t, err)
		}
	}

	// Stage metadata.
	{
		// Get
		{
			// Existing key
			resp, err := r.GetStageMetadata(ctx, &service.GetStageMetadataRequest{
				DeploymentId: d.Id,
				StageId:      "stage-2",
				Key:          "stage-2-key-1",
			})
			assert.NoError(t, err)
			assert.Equal(t, true, resp.Found)
			assert.Equal(t, "stage-2-value-1", resp.Value)

			// Nonexistent key
			resp, err = r.GetStageMetadata(ctx, &service.GetStageMetadataRequest{
				DeploymentId: d.Id,
				StageId:      "stage-1",
				Key:          "not-exist-key",
			})
			assert.NoError(t, err)
			assert.Equal(t, false, resp.Found)
			assert.Equal(t, "", resp.Value)

			// Nonexistent stage
			resp, err = r.GetStageMetadata(ctx, &service.GetStageMetadataRequest{
				DeploymentId: d.Id,
				StageId:      "not-exist-stage",
				Key:          "key-1",
			})
			assert.NoError(t, err)
			assert.Equal(t, false, resp.Found)
			assert.Equal(t, "", resp.Value)

			// Nonexistent deployment
			resp, err = r.GetStageMetadata(ctx, &service.GetStageMetadataRequest{
				DeploymentId: "not-exist-id",
				StageId:      "stage-1",
				Key:          "key-1",
			})
			assert.Error(t, err)
			assert.Equal(t, false, resp.Found)
			assert.Equal(t, "", resp.Value)
		}
		// Put
		{
			// New key
			_, err := r.PutStageMetadata(ctx, &service.PutStageMetadataRequest{
				DeploymentId: d.Id,
				StageId:      "stage-1",
				Key:          "stage-1-key-1",
				Value:        "stage-1-value-1",
			})
			assert.NoError(t, err)
			assert.Equal(t, metadata{
				"stage-1-key-1": "stage-1-value-1",
			}, ac.stages["stage-1"])

			// Nonexistent deployment
			_, err = r.PutStageMetadata(ctx, &service.PutStageMetadataRequest{
				DeploymentId: "not-exist-id",
				StageId:      "stage-1",
				Key:          "stage-1-key-1",
				Value:        "stage-1-value-1",
			})
			assert.Error(t, err)
		}
		// PutMulti
		{
			// New key(2) with one existing key(1)
			_, err := r.PutStageMetadataMulti(ctx, &service.PutStageMetadataMultiRequest{
				DeploymentId: d.Id,
				StageId:      "stage-2",
				Metadata: map[string]string{
					"stage-2-key-1": "stage-2-value-12",
					"stage-2-key-2": "stage-2-value-2",
				},
			})
			assert.NoError(t, err)
			assert.Equal(t, map[string]metadata{
				"stage-1": {
					"stage-1-key-1": "stage-1-value-1",
				},
				"stage-2": {
					"stage-2-key-1": "stage-2-value-12",
					"stage-2-key-2": "stage-2-value-2",
				},
			}, ac.stages)

			// Nonexistent deployment
			_, err = r.PutStageMetadataMulti(ctx, &service.PutStageMetadataMultiRequest{
				DeploymentId: "not-exist-id",
				StageId:      "stage-1",
				Metadata: map[string]string{
					"stage-1-key-1": "stage-1-value-1",
				},
			})
			assert.Error(t, err)
		}
	}
}
