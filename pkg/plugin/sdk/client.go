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
	"iter"
	"slices"
	"time"

	"github.com/pipe-cd/pipecd/pkg/model"
	"github.com/pipe-cd/pipecd/pkg/plugin/pipedapi"
	"github.com/pipe-cd/pipecd/pkg/plugin/pipedservice"
	"github.com/pipe-cd/pipecd/pkg/plugin/toolregistry"
)

const (
	listStageCommandsInterval = 5 * time.Second
)

// Client is a toolkit for interacting with the piped service.
// It provides methods to call the piped service APIs.
// It's a wrapper around the raw piped service client.
type Client struct {
	base *pipedapi.PluginServiceClient

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

	// logPersister is used to persist the stage logs.
	// This field exists only when the client is working with a specific stage; for example, when this client is passed as the ExecuteStage method's argument.
	logPersister StageLogPersister

	// toolRegistry is used to install and get the path of the tools used in the plugin.
	// TODO: We should consider installing the tools in other way.
	toolRegistry *toolregistry.ToolRegistry
}

// NewClient creates a new client.
// DO NOT USE this function except in tests.
// FIXME: Remove this function and make a better way for tests.
func NewClient(base *pipedapi.PluginServiceClient, pluginName, applicationID, stageID string, lp StageLogPersister, tr *toolregistry.ToolRegistry) *Client {
	return &Client{
		base:          base,
		pluginName:    pluginName,
		applicationID: applicationID,
		stageID:       stageID,
		logPersister:  lp,
		toolRegistry:  tr,
	}
}

// StageLogPersister is a interface for persisting the stage logs.
// Use this to persist the stage logs and make it viewable on the UI.
type StageLogPersister interface {
	Write(log []byte) (int, error)
	Info(log string)
	Infof(format string, a ...interface{})
	Success(log string)
	Successf(format string, a ...interface{})
	Error(log string)
	Errorf(format string, a ...interface{})
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

// LogPersister returns the stage log persister.
// Use this to persist the stage logs and make it viewable on the UI.
// This method should be called only when the client is working with a specific stage, for example, when this client is passed as the ExecuteStage method's argument.
// Otherwise, it will return nil.
// TODO: we should consider returning an error instead of nil, or return logger which prints to stdout.
func (c *Client) LogPersister() StageLogPersister {
	return c.logPersister
}

// ToolRegistry returns the tool registry.
// Use this to install and get the path of the tools used in the plugin.
func (c *Client) ToolRegistry() *toolregistry.ToolRegistry {
	return c.toolRegistry
}

// ListStageCommands returns the list of stage commands of the given command types.
func (c Client) ListStageCommands(ctx context.Context, commandTypes ...model.Command_Type) iter.Seq2[*StageCommand, error] {
	return func(yield func(*StageCommand, error) bool) {
		returned := map[string]struct{}{}

		for {
			resp, err := c.base.ListStageCommands(ctx, &pipedservice.ListStageCommandsRequest{
				DeploymentId: c.deploymentID,
				StageId:      c.stageID,
			})
			if err != nil {
				if !yield(nil, err) {
					return
				}
				continue
			}

			for _, command := range resp.Commands {
				if !slices.Contains(commandTypes, command.Type) {
					continue
				}

				if _, ok := returned[command.Id]; ok {
					continue
				}
				returned[command.Id] = struct{}{}

				stageCommand, err := newStageCommand(command)
				if err != nil {
					if !yield(nil, err) {
						return
					}
					continue
				}

				if !yield(&stageCommand, nil) {
					return
				}
			}

			time.Sleep(listStageCommandsInterval)
		}
	}
}
