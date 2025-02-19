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

	"github.com/pipe-cd/pipecd/pkg/plugin/logpersister"
	"github.com/pipe-cd/pipecd/pkg/plugin/pipedapi"
	"github.com/pipe-cd/pipecd/pkg/plugin/pipedservice"
)

// Client is a toolkit for interacting with the piped service.
// It provides methods to call the piped service APIs.
// It's a wrapper around the raw piped service client.
type Client struct {
	base *pipedapi.PipedServiceClient

	// pluginName is used to identify which plugin sends requests to piped.
	pluginName string

	// applicationID is used to identify the application that the client is working with.
	applicationID string
	// deploymentID is used to identify the deployment that the client is working with.
	// This field exists only when the client is working with a specific deployment; for example, when this client is passed as the deployoment plugin's argument.
	deploymentID string
	// stageID is used to identify the stage that the client is working with.
	// This field exists only when the client is working with a specific stage; for example, when this client is passed as the ExecuteStage method's argument.
	stageID string

	// TODO: Define another interface (e.g. Remove Complete() method)
	LogPersister logpersister.StageLogPersister
}

// GetStageMetadata gets the metadata of the current stage.
func (c *Client) GetStageMetadata(ctx context.Context, key string) (string, error) {
	resp, err := c.base.GetStageMetadata(ctx, &pipedservice.GetStageMetadataRequest{
		DeploymentId: c.deploymentID,
		StageId:      c.stageID,
		Key:          key,
	})
	if err != nil {
		return "", err
	}
	return resp.Value, nil
}

// PutStageMetadata stores the metadata of the current stage.
func (c *Client) PutStageMetadata(ctx context.Context, key, value string) error {
	_, err := c.base.PutStageMetadata(ctx, &pipedservice.PutStageMetadataRequest{
		DeploymentId: c.deploymentID,
		StageId:      c.stageID,
		Key:          key,
		Value:        value,
	})
	return err
}

// PutStageMetadataMulti stores the multiple metadata of the current stage.
func (c *Client) PutStageMetadataMulti(ctx context.Context, metadata map[string]string) error {
	_, err := c.base.PutStageMetadataMulti(ctx, &pipedservice.PutStageMetadataMultiRequest{
		DeploymentId: c.deploymentID,
		StageId:      c.stageID,
		Metadata:     metadata,
	})
	return err
}

// GetDeploymentPluginMetadata gets the metadata of the current deployment and plugin.
func (c *Client) GetDeploymentPluginMetadata(ctx context.Context, key string) (string, error) {
	resp, err := c.base.GetDeploymentPluginMetadata(ctx, &pipedservice.GetDeploymentPluginMetadataRequest{
		DeploymentId: c.deploymentID,
		PluginName:   c.pluginName,
		Key:          key,
	})
	return resp.Value, err
}

// PutDeploymentPluginMetadata stores the metadata of the current deployment and plugin.
func (c *Client) PutDeploymentPluginMetadata(ctx context.Context, key, value string) error {
	_, err := c.base.PutDeploymentPluginMetadata(ctx, &pipedservice.PutDeploymentPluginMetadataRequest{
		DeploymentId: c.deploymentID,
		PluginName:   c.pluginName,
		Key:          key,
		Value:        value,
	})
	return err
}

// PutDeploymentPluginMetadataMulti stores the multiple metadata of the current deployment and plugin.
func (c *Client) PutDeploymentPluginMetadataMulti(ctx context.Context, metadata map[string]string) error {
	_, err := c.base.PutDeploymentPluginMetadataMulti(ctx, &pipedservice.PutDeploymentPluginMetadataMultiRequest{
		DeploymentId: c.deploymentID,
		PluginName:   c.pluginName,
		Metadata:     metadata,
	})
	return err
}

// GetDeploymentSharedMetadata gets the metadata of the current deployment
// which is shared among piped and plugins.
func (c *Client) GetDeploymentSharedMetadata(ctx context.Context, key string) (string, error) {
	resp, err := c.base.GetDeploymentSharedMetadata(ctx, &pipedservice.GetDeploymentSharedMetadataRequest{
		DeploymentId: c.deploymentID,
		Key:          key,
	})
	return resp.Value, err
}
