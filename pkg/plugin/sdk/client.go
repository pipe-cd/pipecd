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

	"github.com/pipe-cd/piped-plugin-sdk-go/toolregistry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/pipe-cd/pipecd/pkg/model"
	"github.com/pipe-cd/pipecd/pkg/plugin/pipedservice"
	"github.com/pipe-cd/pipecd/pkg/rpc/rpcclient"
)

const (
	// MetadataKeyStageDisplay is the key of the stage metadata to be displayed on the deployment detail UI.
	MetadataKeyStageDisplay = model.MetadataKeyStageDisplay
	// MetadataKeyStageApprovedUsers is the key of the metadata of who approved the stage.
	// It will be displayed in the DEPLOYMENT_APPROVED notification.
	// e.g. user-1,user-2
	MetadataKeyStageApprovedUsers = model.MetadataKeyStageApprovedUsers

	listStageCommandsInterval = 5 * time.Second
)

type pluginServiceClient struct {
	pipedservice.PluginServiceClient
	conn *grpc.ClientConn
}

func newPluginServiceClient(ctx context.Context, address string, opts ...rpcclient.DialOption) (*pluginServiceClient, error) {
	// Clone the opts to avoid modifying the original opts slice.
	opts = slices.Clone(opts)

	// Append the required options.
	// The WithBlock option is required to make the client wait until the connection is up.
	// The WithInsecure option is required to disable the transport security.
	// The piped service does not require transport security because it is only used in localhost.
	opts = append(opts, rpcclient.WithBlock(), rpcclient.WithInsecure())

	conn, err := rpcclient.DialContext(ctx, address, opts...)
	if err != nil {
		return nil, err
	}

	return &pluginServiceClient{
		PluginServiceClient: pipedservice.NewPluginServiceClient(conn),
		conn:                conn,
	}, nil
}

func (c *pluginServiceClient) Close() error {
	return c.conn.Close()
}

// Client is a toolkit for interacting with the piped service.
// It provides methods to call the piped service APIs.
// It's a wrapper around the raw piped service client.
type Client struct {
	base *pluginServiceClient

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

	// stageLogPersister is used to persist the stage logs.
	// This field exists only when the client is working with a specific stage; for example, when this client is passed as the ExecuteStage method's argument.
	stageLogPersister StageLogPersister

	// toolRegistry is used to install and get the path of the tools used in the plugin.
	// TODO: We should consider installing the tools in other way.
	toolRegistry *toolregistry.ToolRegistry
}

// NewClient creates a new client.
// DO NOT USE this function except in tests.
// FIXME: Remove this function and make a better way for tests.
func NewClient(base *pluginServiceClient, pluginName, applicationID, stageID string, slp StageLogPersister, tr *toolregistry.ToolRegistry) *Client {
	return &Client{
		base:              base,
		pluginName:        pluginName,
		applicationID:     applicationID,
		stageID:           stageID,
		stageLogPersister: slp,
		toolRegistry:      tr,
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
func (c *Client) GetStageMetadata(ctx context.Context, key string) (string, bool, error) {
	resp, err := c.base.GetStageMetadata(ctx, &pipedservice.GetStageMetadataRequest{
		DeploymentId: c.deploymentID,
		StageId:      c.stageID,
		Key:          key,
	})
	if err != nil {
		return "", false, err
	}
	return resp.Value, resp.Found, nil
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
func (c *Client) GetDeploymentPluginMetadata(ctx context.Context, key string) (string, bool, error) {
	resp, err := c.base.GetDeploymentPluginMetadata(ctx, &pipedservice.GetDeploymentPluginMetadataRequest{
		DeploymentId: c.deploymentID,
		PluginName:   c.pluginName,
		Key:          key,
	})
	if err != nil {
		return "", false, err
	}
	return resp.Value, resp.Found, err
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
func (c *Client) GetDeploymentSharedMetadata(ctx context.Context, key string) (string, bool, error) {
	resp, err := c.base.GetDeploymentSharedMetadata(ctx, &pipedservice.GetDeploymentSharedMetadataRequest{
		DeploymentId: c.deploymentID,
		Key:          key,
	})
	if err != nil {
		return "", false, err
	}
	return resp.Value, resp.Found, err
}

// GetApplicationSharedObject gets the application object which is shared across deployments.
func (c *Client) GetApplicationSharedObject(ctx context.Context, key string) (obj []byte, found bool, err error) {
	resp, err := c.base.GetApplicationSharedObject(ctx, &pipedservice.GetApplicationSharedObjectRequest{
		ApplicationId: c.applicationID,
		PluginName:    c.pluginName,
		Key:           key,
	})
	if status.Code(err) == codes.NotFound {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	return resp.Object, true, nil
}

// PutApplicationSharedObject stores the application object which is shared across deployments.
func (c *Client) PutApplicationSharedObject(ctx context.Context, key string, object []byte) error {
	_, err := c.base.PutApplicationSharedObject(ctx, &pipedservice.PutApplicationSharedObjectRequest{
		ApplicationId: c.applicationID,
		PluginName:    c.pluginName,
		Key:           key,
		Object:        object,
	})
	return err
}

// StageLogPersister returns the stage log persister.
// Use this to persist the stage logs and make it viewable on the UI.
// This method should be called only when the client is working with a specific stage, for example, when this client is passed as the ExecuteStage method's argument.
// Otherwise, it will return nil.
// TODO: we should consider returning an error instead of nil, or return logger which prints to stdout.
func (c *Client) StageLogPersister() StageLogPersister {
	return c.stageLogPersister
}

// ToolRegistry returns the tool registry.
// Use this to install and get the path of the tools used in the plugin.
func (c *Client) ToolRegistry() *toolregistry.ToolRegistry {
	return c.toolRegistry
}

// ListStageCommands returns the list of stage commands of the given command types.
func (c Client) ListStageCommands(ctx context.Context, commandTypes ...CommandType) iter.Seq2[*StageCommand, error] {
	return func(yield func(*StageCommand, error) bool) {
		returned := map[string]struct{}{}

		modelCommandTypes := make([]model.Command_Type, 0, len(commandTypes))
		for _, cmdType := range commandTypes {
			modelType, err := cmdType.toModelEnum()
			if err != nil {
				if !yield(nil, err) {
					return
				}
				continue
			}
			modelCommandTypes = append(modelCommandTypes, modelType)
		}

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
				if !slices.Contains(modelCommandTypes, command.Type) {
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
