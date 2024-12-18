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

package lambda

import (
	"context"
	"io"
	"path/filepath"
	"sync"

	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"

	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"

	"github.com/pipe-cd/pipecd/pkg/config"
)

const (
	LabelManagedBy   string = "pipecd-dev-managed-by"  // Always be piped.
	LabelPiped       string = "pipecd-dev-piped"       // The id of piped handling this application.
	LabelApplication string = "pipecd-dev-application" // The application this resource belongs to.
	LabelCommitHash  string = "pipecd-dev-commit-hash" // Hash value of the deployed commit.
	ManagedByPiped   string = "piped"
)

// Client is wrapper of AWS client.
type Client interface {
	IsFunctionExist(ctx context.Context, name string) (bool, error)
	CreateFunction(ctx context.Context, fm FunctionManifest) error
	CreateFunctionFromSource(ctx context.Context, fm FunctionManifest, zip io.Reader) error
	UpdateFunction(ctx context.Context, fm FunctionManifest) error
	UpdateFunctionFromSource(ctx context.Context, fm FunctionManifest, zip io.Reader) error
	PublishFunction(ctx context.Context, fm FunctionManifest) (version string, err error)
	ListFunctions(ctx context.Context) ([]types.FunctionConfiguration, error)
	GetFunction(ctx context.Context, functionName string) (*lambda.GetFunctionOutput, error)
	GetTrafficConfig(ctx context.Context, fm FunctionManifest) (routingTrafficCfg RoutingTrafficConfig, err error)
	CreateTrafficConfig(ctx context.Context, fm FunctionManifest, version string) error
	UpdateTrafficConfig(ctx context.Context, fm FunctionManifest, routingTraffic RoutingTrafficConfig) error
}

// Registry holds a pool of aws client wrappers.
type Registry interface {
	Client(name string, cfg *config.PlatformProviderLambdaConfig, logger *zap.Logger) (Client, error)
}

// LoadFunctionManifest returns FunctionManifest object from a given Function config manifest file.
func LoadFunctionManifest(appDir, functionManifestFilename string) (FunctionManifest, error) {
	path := filepath.Join(appDir, functionManifestFilename)
	return loadFunctionManifest(path)
}

type registry struct {
	clients  map[string]Client
	mu       sync.RWMutex
	newGroup *singleflight.Group
}

func (r *registry) Client(name string, cfg *config.PlatformProviderLambdaConfig, logger *zap.Logger) (Client, error) {
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
