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

package kubernetes

import (
	"context"

	"github.com/kapetaniosci/pipe/pkg/app/runner/executor"
	"github.com/kapetaniosci/pipe/pkg/model"
)

const (
	subsetLabel    = "pipecd.dev/subset"
	managedByLabel = "pipecd.dev/managed-by"
)

type Executor struct {
	executor.Input
}

func init() {
	var (
		f = func(in executor.Input) executor.Executor {
			return &Executor{
				Input: in,
			}
		}
		r = executor.DefaultRegistry()
	)

	r.Register(model.StageK8sPrimaryOut, f)
	r.Register(model.StageK8sStageOut, f)
	r.Register(model.StageK8sStageIn, f)
	r.Register(model.StageK8sBaselineOut, f)
	r.Register(model.StageK8sBaselineIn, f)
	r.Register(model.StageK8sPrimaryOut, f)
	r.Register(model.StageK8sTrafficRoute, f)
}

func (e *Executor) Execute(ctx context.Context) (model.StageStatus, error) {
	return model.StageStatus_STAGE_SUCCESS, nil
}

func (e *Executor) ensureStageRollOut() error {
	return nil
}

func (e *Executor) ensureStageRemove() error {
	return nil
}

func (e *Executor) ensurePrimaryUpdate() error {
	return nil
}

func (e *Executor) ensureBaselineRollout() error {
	return nil
}

func (e *Executor) ensureBaselineRemove() error {
	return nil
}

func (e *Executor) ensureTrafficRoute() error {
	return nil
}
