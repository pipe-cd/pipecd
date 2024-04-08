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

// Package controller provides a piped component
// that handles all of the not completed deployments by managing a pool of planners and schedulers.
// Whenever a new PENDING deployment is detected, controller spawns a new planner for deciding
// the deployment pipeline and update the deployment status to PLANNED.
// Whenever a new PLANNED deployment is detected, controller spawns a new scheduler
// for scheduling and running its pipeline executors.
package controller

import (
	"sync"

	"github.com/pipe-cd/pipecd/pkg/model"
	"github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1/platform"
)

type PluginRegistry interface {
	Plugin(k model.ApplicationKind) (platform.PlatformPluginClient, bool)
}

type pluginRegistry struct {
	plugins map[model.ApplicationKind]platform.PlatformPluginClient
	mu      sync.RWMutex
}

func (r *pluginRegistry) Plugin(k model.ApplicationKind) (platform.PlatformPluginClient, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	e, ok := r.plugins[k]
	if !ok {
		return nil, false
	}

	return e, true
}

var defaultPluginRegistry = &pluginRegistry{
	plugins: make(map[model.ApplicationKind]platform.PlatformPluginClient),
}

func DefaultPluginRegistry() PluginRegistry {
	return defaultPluginRegistry
}

func init() {
	// TODO: Register all available built-in plugins.

	// NOTE: If you want to directry test the plugin, you can use the following code.

	// defaultPluginRegistry.mu.Lock()
	// defer defaultPluginRegistry.mu.Unlock()

	// options := []rpcclient.DialOption{
	// 	rpcclient.WithBlock(),
	// 	rpcclient.WithInsecure(),
	// }

	// cli, err := platform.NewClient(context.Background(), "localhost:10000", options...)
	// if err != nil {
	// 	panic(err)
	// }

	// defaultPluginRegistry.plugins[model.ApplicationKind_KUBERNETES] = cli
}
