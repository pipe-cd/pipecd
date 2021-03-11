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

package ecs

import (
	"context"
	"path/filepath"
	"sync"

	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"

	"github.com/pipe-cd/pipe/pkg/config"
)

const (
	defaultTaskDefinitionFilename    = "taskdef.yaml"
	defaultserviceDefinitionFilename = "servicedef.yaml"
)

// Client is wrapper of ECS client.
type Client interface {
	ServiceExists(ctx context.Context, clusterName string, services []string) (bool, error)
	CreateService(ctx context.Context, service types.Service) (*types.Service, error)
	UpdateService(ctx context.Context, service types.Service) (*types.Service, error)
	RegisterTaskDefinition(ctx context.Context, taskDefinition types.TaskDefinition) (*types.TaskDefinition, error)
	DeregisterTaskDefinition(ctx context.Context, taskDefinition types.TaskDefinition) (*types.TaskDefinition, error)
	CreateTaskSet(ctx context.Context, service types.Service, taskDefinition types.TaskDefinition) (*types.TaskSet, error)
	DeleteTaskSet(ctx context.Context, service types.Service, taskSet types.TaskSet) (*types.TaskSet, error)
}

// Registry holds a pool of aws client wrappers.
type Registry interface {
	Client(name string, cfg *config.CloudProviderLambdaConfig, logger *zap.Logger) (Client, error)
}

// LoadServiceDefinition returns ServiceDefinition object from a given service definition file.
func LoadServiceDefinition(appDir, serviceDefinitionFilename string) (types.Service, error) {
	if serviceDefinitionFilename == "" {
		serviceDefinitionFilename = defaultserviceDefinitionFilename
	}
	path := filepath.Join(appDir, serviceDefinitionFilename)
	return loadServiceDefinition(path)
}

// LoadTaskDefinition returns TaskDefinition object from a given task definition file.
func LoadTaskDefinition(appDir, serviceDefinitionFilename string) (types.TaskDefinition, error) {
	if serviceDefinitionFilename == "" {
		serviceDefinitionFilename = defaultserviceDefinitionFilename
	}
	path := filepath.Join(appDir, serviceDefinitionFilename)
	return loadTaskDefinition(path)
}

type registry struct {
	clients  map[string]Client
	mu       sync.RWMutex
	newGroup *singleflight.Group
}

func (r *registry) Client(name string, cfg *config.CloudProviderLambdaConfig, logger *zap.Logger) (Client, error) {
	r.mu.RLock()
	client, ok := r.clients[name]
	r.mu.RUnlock()
	if ok {
		return client, nil
	}

	c, err, _ := r.newGroup.Do(name, func() (interface{}, error) {
		return newClient(cfg.Region, cfg.Profile, cfg.CredentialsFile, cfg.RoleARN, cfg.TokenFile, logger)
	})
	if err != nil {
		return nil, err
	}

	client = c.(Client)
	r.mu.Lock()
	r.clients[name] = client
	r.mu.Unlock()

	return client, nil
}

var defaultRegistry = &registry{
	clients:  make(map[string]Client),
	newGroup: &singleflight.Group{},
}

// DefaultRegistry returns a pool of aws clients and a mutex associated with it.
func DefaultRegistry() Registry {
	return defaultRegistry
}
