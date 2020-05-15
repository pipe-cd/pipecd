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

package executor

import (
	"context"

	"go.uber.org/zap"

	"github.com/kapetaniosci/pipe/pkg/config"
	"github.com/kapetaniosci/pipe/pkg/model"
)

type Executor interface {
	// Execute starts running executor until completion
	// or the context has been cancelled.
	Execute(ctx context.Context) model.StageStatus
}

type Factory func(in Input) Executor

type LogPersister interface {
	Append(log string, s model.LogSeverity)
	AppendInfo(log string)
	AppendSuccess(log string)
	AppendError(log string)
}

type MetadataPersister interface {
	Save(ctx context.Context, metadata []byte) error
}

type CommandStore interface {
	ListDeploymentCommands(deploymentID string) []*model.Command
	ReportCommandHandled(ctx context.Context, c *model.Command, status model.CommandStatus, metadata map[string]string) error
}

type Input struct {
	Stage             *model.PipelineStage
	Deployment        *model.Deployment
	DeploymentConfig  *config.Config
	PipedConfig       *config.PipedSpec
	WorkingDir        string
	CommandStore      CommandStore
	LogPersister      LogPersister
	MetadataPersister MetadataPersister
	Logger            *zap.Logger
}
