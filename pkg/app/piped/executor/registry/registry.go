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

package registry

import (
	"fmt"
	"sync"

	"github.com/kapetaniosci/pipe/pkg/app/piped/executor"
	"github.com/kapetaniosci/pipe/pkg/app/piped/executor/analysis"
	"github.com/kapetaniosci/pipe/pkg/app/piped/executor/cloudrun"
	"github.com/kapetaniosci/pipe/pkg/app/piped/executor/kubernetes"
	"github.com/kapetaniosci/pipe/pkg/app/piped/executor/lambda"
	"github.com/kapetaniosci/pipe/pkg/app/piped/executor/terraform"
	"github.com/kapetaniosci/pipe/pkg/app/piped/executor/wait"
	"github.com/kapetaniosci/pipe/pkg/app/piped/executor/waitapproval"
	"github.com/kapetaniosci/pipe/pkg/model"
)

type Registry interface {
	Executor(stage model.Stage, in executor.Input) (executor.Executor, error)
}

type registry struct {
	factories map[model.Stage]executor.Factory
	mu        sync.RWMutex
}

func (r *registry) Register(stage model.Stage, f executor.Factory) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.factories[stage]; ok {
		return fmt.Errorf("executor for %s stage has already registered", stage)
	}
	r.factories[stage] = f
	return nil
}

func (r *registry) Executor(stage model.Stage, in executor.Input) (executor.Executor, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	f, ok := r.factories[stage]
	if !ok {
		return nil, fmt.Errorf("no registered executor for stage %s", stage)
	}
	return f(in), nil
}

var defaultRegistry = &registry{
	factories: make(map[model.Stage]executor.Factory),
}

func DefaultRegistry() Registry {
	return defaultRegistry
}

// init registers all built-in executors to the default registry.
func init() {
	analysis.Register(defaultRegistry)
	cloudrun.Register(defaultRegistry)
	kubernetes.Register(defaultRegistry)
	lambda.Register(defaultRegistry)
	terraform.Register(defaultRegistry)
	wait.Register(defaultRegistry)
	waitapproval.Register(defaultRegistry)
}
