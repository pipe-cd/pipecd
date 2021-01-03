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

package lambda

import (
	"github.com/pipe-cd/pipe/pkg/app/piped/executor"
	"github.com/pipe-cd/pipe/pkg/config"
	"github.com/pipe-cd/pipe/pkg/model"
)

type Executor struct {
	executor.Input
}

type registerer interface {
	Register(stage model.Stage, f executor.Factory) error
	RegisterRollback(kind model.ApplicationKind, f executor.Factory) error
}

func Register(r registerer) {
	f := func(in executor.Input) executor.Executor {
		return &Executor{
			Input: in,
		}
	}
	r.Register(model.StageLambdaSync, f)
	r.Register(model.StageLambdaPromote, f)
	// TODO registerRollback
}

func (e *Executor) Execute(sig executor.StopSignal) model.StageStatus {
	return model.StageStatus_STAGE_SUCCESS
}

func findCloudProvider(in *executor.Input) (name string, cfg *config.CloudProviderLambdaConfig, found bool) {
	name = in.Application.CloudProvider
	if name == "" {
		in.LogPersister.Error("Missing the CloudProvider name in the application configuration")
		return
	}

	cp, ok := in.PipedConfig.FindCloudProvider(name, model.CloudProviderLambda)
	if !ok {
		in.LogPersister.Errorf("The specified cloud provider %q was not found in piped configuration", name)
		return
	}

	cfg = cp.LambdaConfig
	found = true
	return
}
