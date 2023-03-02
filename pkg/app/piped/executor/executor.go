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

package executor

import (
	"context"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/app/piped/deploysource"
	"github.com/pipe-cd/pipecd/pkg/app/piped/metadatastore"
	provider "github.com/pipe-cd/pipecd/pkg/app/piped/platformprovider/kubernetes"
	"github.com/pipe-cd/pipecd/pkg/cache"
	"github.com/pipe-cd/pipecd/pkg/config"
	"github.com/pipe-cd/pipecd/pkg/git"
	"github.com/pipe-cd/pipecd/pkg/model"
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

type CommandLister interface {
	ListCommands() []model.ReportableCommand
}

type AppLiveResourceLister interface {
	ListKubernetesResources() ([]provider.Manifest, bool)
}

type AnalysisResultStore interface {
	GetLatestAnalysisResult(ctx context.Context) (*model.AnalysisResult, error)
	PutLatestAnalysisResult(ctx context.Context, analysisResult *model.AnalysisResult) error
}

type Notifier interface {
	Notify(event model.NotificationEvent)
}

type GitClient interface {
	Clone(ctx context.Context, repoID, remote, branch, destination string) (git.Repo, error)
}

type Input struct {
	Stage       *model.PipelineStage
	StageConfig config.PipelineStage
	// Readonly deployment model.
	Deployment  *model.Deployment
	Application *model.Application
	PipedConfig *config.PipedSpec
	// Deploy source at target commit
	TargetDSP deploysource.Provider
	// Deploy source at running commit
	RunningDSP               deploysource.Provider
	GitClient                GitClient
	CommandLister            CommandLister
	LogPersister             LogPersister
	MetadataStore            metadatastore.MetadataStore
	AppManifestsCache        cache.Cache
	AppLiveResourceLister    AppLiveResourceLister
	AnalysisResultStore      AnalysisResultStore
	RollbackCustomStageStack []*model.PipelineStage
	Logger                   *zap.Logger
	Notifier                 Notifier
}

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
	}
	return model.StageStatus_STAGE_FAILURE
}
