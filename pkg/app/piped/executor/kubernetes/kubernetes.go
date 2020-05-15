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
	"os"
	"path/filepath"

	"go.uber.org/zap"

	"github.com/kapetaniosci/pipe/pkg/app/piped/executor"
	"github.com/kapetaniosci/pipe/pkg/app/piped/toolregistry"
	"github.com/kapetaniosci/pipe/pkg/config"
	"github.com/kapetaniosci/pipe/pkg/model"
)

const (
	variantLabel    = "pipecd.dev/variant"
	managedByLabel  = "pipecd.dev/managed-by"
	commitHashLabel = "pipecd.dev/commit-hash"

	kustomizationFileName = "kustomization.yaml"
)

type TemplatingMethod string

const (
	TemplatingMethodHelm      TemplatingMethod = "helm"
	TemplatingMethodKustomize TemplatingMethod = "kustomize"
	TemplatingMethodNone      TemplatingMethod = "none"
)

type Executor struct {
	executor.Input
}

type registerer interface {
	Register(stage model.Stage, f executor.Factory) error
}

func Register(r registerer) {
	f := func(in executor.Input) executor.Executor {
		return &Executor{
			Input: in,
		}
	}
	r.Register(model.StageK8sPrimaryUpdate, f)
	r.Register(model.StageK8sStageRollout, f)
	r.Register(model.StageK8sStageClean, f)
	r.Register(model.StageK8sBaselineRollout, f)
	r.Register(model.StageK8sBaselineClean, f)
	r.Register(model.StageK8sTrafficSplit, f)
}

func (e *Executor) Execute(ctx context.Context) model.StageStatus {
	var (
		appDirPath       = filepath.Join(e.RepoDir, e.Deployment.GitPath.Path)
		templatingMethod = determineTemplatingMethod(e.DeploymentConfig, appDirPath)
	)
	e.Logger.Info("start executing kubernetes stage",
		zap.String("app-dir-path", appDirPath),
		zap.String("templating-method", string(templatingMethod)),
	)

	_, _, _ = toolregistry.DefaultRegistry().Kubectl(ctx, "1.8.0")

	return model.StageStatus_STAGE_SUCCESS
}

func (e *Executor) ensureStageRollout() error {
	return nil
}

func (e *Executor) ensureStageClean() error {
	return nil
}

func (e *Executor) ensurePrimaryUpdate() error {
	return nil
}

func (e *Executor) ensureBaselineRollout() error {
	return nil
}

func (e *Executor) ensureBaselineClean() error {
	return nil
}

func (e *Executor) ensureTrafficSplit() error {
	return nil
}

func (e *Executor) generateStageManifests() error {
	return nil
}

func (e *Executor) generateBaselineManifests() error {
	return nil
}

func determineTemplatingMethod(deploymentConfig *config.Config, appDirPath string) TemplatingMethod {
	if input := deploymentConfig.KubernetesDeploymentSpec.Input; input != nil {
		if input.HelmChart != nil {
			return TemplatingMethodHelm
		}
		if len(input.HelmValueFiles) > 0 {
			return TemplatingMethodHelm
		}
		if input.HelmVersion != "" {
			return TemplatingMethodHelm
		}
	}
	if _, err := os.Stat(filepath.Join(appDirPath, kustomizationFileName)); err == nil {
		return TemplatingMethodKustomize
	}
	return TemplatingMethodNone
}
