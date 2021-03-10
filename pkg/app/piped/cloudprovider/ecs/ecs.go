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

const defaultTaskDefinitionFilename = "taskdef.json"
const defaultserviceDefinitionFilename = "servicedef.json"

// Client is wrapper of AWS client.
type Client interface {
	IsServiceExist(ctx context.Context, name string) (bool, error)
	CreateService(ctx context.Context) error
	UpdateService(ctx context.Context) error
	PublishService(ctx context.Context) (version string, err error)
	GetTrafficConfig(ctx context.Context) (err error)
	CreateTrafficConfig(ctx context.Context, version string) error
	UpdateTrafficConfig(ctx context.Context) error
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
