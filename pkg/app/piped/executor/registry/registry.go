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

	"github.com/pipe-cd/pipecd/pkg/app/piped/executor"
	"github.com/pipe-cd/pipecd/pkg/app/piped/executor/analysis"
	"github.com/pipe-cd/pipecd/pkg/app/piped/executor/cloudrun"
	"github.com/pipe-cd/pipecd/pkg/app/piped/executor/customsync"
	"github.com/pipe-cd/pipecd/pkg/app/piped/executor/ecs"
	"github.com/pipe-cd/pipecd/pkg/app/piped/executor/kubernetes"
	"github.com/pipe-cd/pipecd/pkg/app/piped/executor/lambda"
	"github.com/pipe-cd/pipecd/pkg/app/piped/executor/scriptrun"
	"github.com/pipe-cd/pipecd/pkg/app/piped/executor/terraform"
	"github.com/pipe-cd/pipecd/pkg/app/piped/executor/wait"
	"github.com/pipe-cd/pipecd/pkg/app/piped/executor/waitapproval"
	"github.com/pipe-cd/pipecd/pkg/model"
)

type Registry interface {
	Executor(stage model.Stage, in executor.Input) (executor.Executor, bool)
	RollbackExecutor(kind model.ApplicationKind, in executor.Input) (executor.Executor, bool)
}

type registry struct {
	factories         map[model.Stage]executor.Factory
	rollbackFactories map[model.RollbackKind]executor.Factory
	mu                sync.RWMutex
}

func (r *registry) Register(stage model.Stage, f executor.Factory) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.factories[stage]; ok {
		return fmt.Errorf("executor for %s stage has already been registered", stage)
	}
	r.factories[stage] = f
	return nil
}

func (r *registry) RegisterRollback(kind model.RollbackKind, f executor.Factory) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.rollbackFactories[kind]; ok {
		return fmt.Errorf("rollback executor for %s application kind has already been registered", kind.String())
	}
	r.rollbackFactories[kind] = f
	return nil
}

func (r *registry) Executor(stage model.Stage, in executor.Input) (executor.Executor, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	f, ok := r.factories[stage]
	if !ok {
		return nil, false
	}
	return f(in), true
}

func (r *registry) RollbackExecutor(kind model.ApplicationKind, in executor.Input) (executor.Executor, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var rollbackKind model.RollbackKind
	if in.Stage.Name == model.StageCustomSyncRollback.String() {
		rollbackKind = model.RollbackKind_Rollback_CUSTOM_SYNC
	} else {
		rollbackKind = kind.ToRollbackKind()
	}

	f, ok := r.rollbackFactories[rollbackKind]
	if !ok {
		return nil, false
	}
	return f(in), true
}

var defaultRegistry = &registry{
	factories:         make(map[model.Stage]executor.Factory),
	rollbackFactories: make(map[model.RollbackKind]executor.Factory),
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
	ecs.Register(defaultRegistry)
	wait.Register(defaultRegistry)
	waitapproval.Register(defaultRegistry)
	customsync.Register(defaultRegistry)
	scriptrun.Register(defaultRegistry)
}
