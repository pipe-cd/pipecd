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
	"fmt"
	"path/filepath"

	"go.uber.org/zap"

	provider "github.com/kapetaniosci/pipe/pkg/app/piped/cloudprovider/kubernetes"
	"github.com/kapetaniosci/pipe/pkg/app/piped/executor"
	"github.com/kapetaniosci/pipe/pkg/config"
	"github.com/kapetaniosci/pipe/pkg/model"
)

const (
	primaryVariant  = "primary"
	canaryVariant   = "canary"
	baselineVariant = "baseline"

	metadataKeyAddedStageResources    = "k8s-canary-resources"
	metadataKeyAddedBaselineResources = "k8s-baseline-resources"
)

type Executor struct {
	executor.Input

	provider provider.Provider
	config   *config.KubernetesDeploymentSpec
}

type registerer interface {
	Register(stage model.Stage, f executor.Factory) error
}

// Register registers this executor factory into a given registerer.
func Register(r registerer) {
	f := func(in executor.Input) executor.Executor {
		return &Executor{
			Input: in,
		}
	}
	r.Register(model.StageK8sPrimaryUpdate, f)
	r.Register(model.StageK8sCanaryRollout, f)
	r.Register(model.StageK8sCanaryClean, f)
	r.Register(model.StageK8sBaselineRollout, f)
	r.Register(model.StageK8sBaselineClean, f)
	r.Register(model.StageK8sTrafficSplit, f)
}

func (e *Executor) Execute(sig executor.StopSignal) model.StageStatus {
	e.config = e.DeploymentConfig.KubernetesDeploymentSpec
	if e.config == nil {
		e.LogPersister.AppendError("Malformed deployment configuration: missing KubernetesDeploymentSpec")
		return model.StageStatus_STAGE_FAILURE
	}

	var (
		ctx    = sig.Context()
		appDir = filepath.Join(e.RepoDir, e.Deployment.GitPath.Path)
	)

	e.provider = provider.NewProvider(e.RepoDir, appDir, e.config.Input, e.Logger)
	if err := e.provider.Init(ctx); err != nil {
		e.LogPersister.AppendError(fmt.Sprintf("Failed while initializing kubernetes provider (%v)", err))
		return model.StageStatus_STAGE_FAILURE
	}

	e.Logger.Info("start executing kubernetes stage",
		zap.String("stage-name", e.Stage.Name),
		zap.String("app-dir", appDir),
	)

	var (
		originalStatus = e.Stage.Status
		status         model.StageStatus
	)

	switch model.Stage(e.Stage.Name) {
	case model.StageK8sPrimaryUpdate:
		status = e.ensurePrimaryUpdate(ctx)
	case model.StageK8sCanaryRollout:
		status = e.ensureCanaryRollout(ctx)
	case model.StageK8sCanaryClean:
		status = e.ensureCanaryClean(ctx)
	case model.StageK8sBaselineRollout:
		status = e.ensureBaselineRollout(ctx)
	case model.StageK8sBaselineClean:
		status = e.ensureBaselineClean(ctx)
	case model.StageK8sTrafficSplit:
		status = e.ensureTrafficSplit(ctx)
	default:
		e.LogPersister.AppendError(fmt.Sprintf("Unsupported stage %s for kubernetes application", e.Stage.Name))
		return model.StageStatus_STAGE_FAILURE
	}

	return executor.DetermineStageStatus(sig.Signal(), originalStatus, status)
}

func (e *Executor) ensureTrafficSplit(ctx context.Context) model.StageStatus {
	return model.StageStatus_STAGE_SUCCESS
}

func (e *Executor) ensureRollback(ctx context.Context) model.StageStatus {
	return model.StageStatus_STAGE_SUCCESS
}
