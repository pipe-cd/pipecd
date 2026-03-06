// Copyright 2026 The PipeCD Authors.
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

package provider

import (
	"path/filepath"
	"sync"

	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/ecs/config"
	"golang.org/x/sync/singleflight"
)

const (
	LabelManagedBy     string = "pipecd-dev-managed-by"  // Always be piped-ecs-plugin.
	LabelPiped         string = "pipecd-dev-piped"       // The id of piped handling this application.
	LabelApplication   string = "pipecd-dev-application" // The application this resource belongs to.
	LabelCommitHash    string = "pipecd-dev-commit-hash" // Hash value of the deployed commit.
	ManagedByECSPlugin string = "piped-ecs-plugin"
)

type Client interface {
	ECS
}

// ECS defines methods for interacting with ECS resources.
type ECS interface {
}

// LoadTaskDefinition returns TaskDefinition object from a given task definition file.
func LoadTaskDefinition(appDir, taskDefinition string) (types.TaskDefinition, error) {
	path := filepath.Join(appDir, taskDefinition)
	return loadTaskDefinition(path)
}

// LoadServiceDefinition returns Service object from a given service definition file.
func LoadServiceDefinition(appDir, serviceDefinition string) (types.Service, error) {
	path := filepath.Join(appDir, serviceDefinition)
	return loadServiceDefinition(path)
}

// Registry holds a pool of aws client wrappers.
type Registry interface {
	Client(name string, cfg config.ECSDeployTargetConfig) (Client, error)
}

type registry struct {
	clients  map[string]Client
	mu       sync.RWMutex
	newGroup *singleflight.Group
}

func (r *registry) Client(name string, cfg config.ECSDeployTargetConfig) (Client, error) {
	r.mu.RLock()
	client, ok := r.clients[name]
	r.mu.RUnlock()
	if ok {
		return client, nil
	}

	c, err, _ := r.newGroup.Do(name, func() (interface{}, error) {
		return newClient(cfg.Region, cfg.Profile, cfg.CredentialsFile, cfg.RoleARN, cfg.TokenFile)
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
