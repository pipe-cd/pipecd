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

package registry

import (
	"fmt"
	"sync"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/planner"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/planner/cloudrun"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/planner/ecs"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/planner/kubernetes"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/planner/lambda"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/planner/terraform"
	"github.com/pipe-cd/pipecd/pkg/model"
)

type Registry interface {
	Planner(k model.ApplicationKind) (planner.Planner, bool)
}

type registry struct {
	planners map[model.ApplicationKind]planner.Planner
	mu       sync.RWMutex
}

func (r *registry) Register(k model.ApplicationKind, p planner.Planner) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.planners[k]; ok {
		return fmt.Errorf("planner for %v application kind has already been registered", k)
	}
	r.planners[k] = p
	return nil
}

func (r *registry) Planner(k model.ApplicationKind) (planner.Planner, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	p, ok := r.planners[k]
	if !ok {
		return nil, false
	}
	return p, true
}

var defaultRegistry = &registry{
	planners: make(map[model.ApplicationKind]planner.Planner),
}

func DefaultRegistry() Registry {
	return defaultRegistry
}

// init registers all planners to the default registry.
func init() {
	cloudrun.Register(defaultRegistry)
	kubernetes.Register(defaultRegistry)
	lambda.Register(defaultRegistry)
	terraform.Register(defaultRegistry)
	ecs.Register(defaultRegistry)
}
