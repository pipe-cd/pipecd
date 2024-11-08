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

package executor

import (
	"context"

	"go.uber.org/zap"

	config "github.com/pipe-cd/pipecd/pkg/configv1"
	"github.com/pipe-cd/pipecd/pkg/git"
	"github.com/pipe-cd/pipecd/pkg/model"
	"github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1/deployment"
)

type Executor interface {
	// Execute starts running executor until completion
	// or the StopSignal has emitted.
	Execute(sig StopSignal) model.StageStatus
}

type Factory func(in Input) Executor

type LogPersister interface {
	Write(log []byte) (int, error)
	Info(log string)
	Infof(format string, a ...interface{})
	Success(log string)
	Successf(format string, a ...interface{})
	Error(log string)
	Errorf(format string, a ...interface{})
}

type GitClient interface {
	Clone(ctx context.Context, repoID, remote, branch, destination string) (git.Repo, error)
}

type Input struct {
	Stage       *model.PipelineStage
	Deployment  *model.Deployment
	StageConfig *config.PipelineStage
	// Deploy source at target commit
	TargetDS *deployment.DeploymentSource
	// Deploy source at running commit
	RunningDS    *deployment.DeploymentSource
	LogPersister LogPersister
	Logger       *zap.Logger
}

// DetermineStageStatus determines the final status of the stage based on the given stop signal.
// Normal is the case when the stop signal is StopSignalNone.
func DetermineStageStatus(sig StopSignalType, ori, got model.StageStatus) model.StageStatus {
	switch sig {
	case StopSignalNone:
		return got
	case StopSignalTerminate:
		return ori
	case StopSignalCancel:
		return model.StageStatus_STAGE_CANCELLED
	case StopSignalTimeout:
		return model.StageStatus_STAGE_FAILURE
	default:
		return model.StageStatus_STAGE_FAILURE
	}
}
