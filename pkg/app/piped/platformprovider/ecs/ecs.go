// Copyright 2023 The PipeCD Authors.
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

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"

	"github.com/pipe-cd/pipecd/pkg/config"
)

const (
	LabelManagedBy   string = "pipecd-dev-managed-by"  // Always be piped.
	LabelPiped       string = "pipecd-dev-piped"       // The id of piped handling this application.
	LabelApplication string = "pipecd-dev-application" // The application this resource belongs to.
	LabelCommitHash  string = "pipecd-dev-commit-hash" // Hash value of the deployed commit.
	ManagedByPiped   string = "piped"
)

// Client is wrapper of ECS client.
type Client interface {
	ECS
	ELB
}

type ECS interface {
	ServiceExists(ctx context.Context, clusterName string, servicesName string) (bool, error)
	CreateService(ctx context.Context, service types.Service) (*types.Service, error)
	UpdateService(ctx context.Context, service types.Service) (*types.Service, error)
	WaitServiceStable(ctx context.Context, service types.Service) error
	RegisterTaskDefinition(ctx context.Context, taskDefinition types.TaskDefinition) (*types.TaskDefinition, error)
	RunTask(ctx context.Context, taskDefinition types.TaskDefinition, clusterArn string, launchType string, awsVpcConfiguration *config.ECSVpcConfiguration, tags []types.Tag) error
	GetPrimaryTaskSet(ctx context.Context, service types.Service) (*types.TaskSet, error)
	CreateTaskSet(ctx context.Context, service types.Service, taskDefinition types.TaskDefinition, targetGroup *types.LoadBalancer, scale int) (*types.TaskSet, error)
	DeleteTaskSet(ctx context.Context, service types.Service, taskSetArn string) error
	UpdateServicePrimaryTaskSet(ctx context.Context, service types.Service, taskSet types.TaskSet) (*types.TaskSet, error)
	TagResource(ctx context.Context, resourceArn string, tags []types.Tag) error
}

type ELB interface {
	GetListenerArns(ctx context.Context, targetGroup types.LoadBalancer) ([]string, error)
	GetListenerRuleArns(ctx context.Context, listenerArns []string) ([]string, error)
	ModifyListenerOrRule(ctx context.Context, listenerArns []string, listenerRuleArns []string, routingTrafficCfg RoutingTrafficConfig) error
}

// Registry holds a pool of aws client wrappers.
type Registry interface {
	Client(name string, cfg *config.PlatformProviderECSConfig, logger *zap.Logger) (Client, error)
}

// LoadServiceDefinition returns ServiceDefinition object from a given service definition file.
func LoadServiceDefinition(appDir, serviceDefinitionFilename string) (types.Service, error) {
	path := filepath.Join(appDir, serviceDefinitionFilename)
	return loadServiceDefinition(path)
}

// LoadTaskDefinition returns TaskDefinition object from a given task definition file.
func LoadTaskDefinition(appDir, taskDefinition string) (types.TaskDefinition, error) {
	path := filepath.Join(appDir, taskDefinition)
	return loadTaskDefinition(path)
}

// LoadTargetGroups returns primary & canary target groups according to the defined in pipe definition file.
func LoadTargetGroups(targetGroups config.ECSTargetGroups) (*types.LoadBalancer, *types.LoadBalancer, error) {
	return loadTargetGroups(targetGroups)
}

type registry struct {
	clients  map[string]Client
	mu       sync.RWMutex
	newGroup *singleflight.Group
}

func (r *registry) Client(name string, cfg *config.PlatformProviderECSConfig, logger *zap.Logger) (Client, error) {
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

func MakeTags(tags map[string]string) []types.Tag {
	resourceTags := make([]types.Tag, 0, len(tags))
	for key, value := range tags {
		resourceTags = append(resourceTags, types.Tag{Key: aws.String(key), Value: aws.String(value)})
	}
	return resourceTags
}
