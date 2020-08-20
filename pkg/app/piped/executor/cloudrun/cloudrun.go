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

package cloudrun

import (
	"context"
	"path/filepath"

	provider "github.com/pipe-cd/pipe/pkg/app/piped/cloudprovider/cloudrun"
	"github.com/pipe-cd/pipe/pkg/app/piped/executor"
	"github.com/pipe-cd/pipe/pkg/config"
	"github.com/pipe-cd/pipe/pkg/model"
)

type Executor struct {
	executor.Input

	config              *config.CloudRunDeploymentSpec
	cloudProviderName   string
	cloudProviderConfig *config.CloudProviderCloudRunConfig
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

	r.Register(model.StageCloudRunSync, f)
	r.Register(model.StageCloudRunPromote, f)

	r.RegisterRollback(model.ApplicationKind_CLOUDRUN, f)
}

func (e *Executor) Execute(sig executor.StopSignal) model.StageStatus {
	e.config = e.DeploymentConfig.CloudRunDeploymentSpec
	if e.config == nil {
		e.LogPersister.Error("Malformed deployment configuration: missing CloudRunDeploymentSpec")
		return model.StageStatus_STAGE_FAILURE
	}

	e.cloudProviderName = e.Application.CloudProvider
	if e.cloudProviderName == "" {
		e.LogPersister.Error("This application configuration was missing CloudProvider name")
		return model.StageStatus_STAGE_FAILURE
	}

	cpConfig, ok := e.PipedConfig.FindCloudProvider(e.cloudProviderName, model.CloudProviderCloudRun)
	if !ok {
		e.LogPersister.Errorf("The specified cloud provider %q was not found in piped configuration", e.cloudProviderName)
		return model.StageStatus_STAGE_FAILURE
	}
	e.cloudProviderConfig = cpConfig.CloudRunConfig

	var (
		ctx            = sig.Context()
		originalStatus = e.Stage.Status
		status         model.StageStatus
	)

	switch model.Stage(e.Stage.Name) {
	case model.StageCloudRunSync:
		status = e.ensureSync(ctx)

	case model.StageCloudRunPromote:
		status = e.ensurePromote(ctx)

	case model.StageRollback:
		status = e.ensureRollback(ctx)

	default:
		e.LogPersister.Errorf("Unsupported stage %s for cloudrun application", e.Stage.Name)
		return model.StageStatus_STAGE_FAILURE
	}

	return executor.DetermineStageStatus(sig.Signal(), originalStatus, status)
}

func (e *Executor) ensureSync(ctx context.Context) model.StageStatus {
	commit := e.Deployment.Trigger.Commit.Hash
	sm, ok := e.loadServiceManifest()
	if !ok {
		return model.StageStatus_STAGE_FAILURE
	}

	e.LogPersister.Info("Generate a service manifest that configures all traffic to the revision specified at the triggered commit")
	revision, err := provider.DecideRevisionName(sm, commit)
	if err != nil {
		e.LogPersister.Errorf("Unable to decide revision name for the commit %s (%v)", commit, err)
		return model.StageStatus_STAGE_FAILURE
	}

	if err := sm.SetRevision(revision); err != nil {
		e.LogPersister.Errorf("Unable to set revision name to service manifest (%v)", err)
		return model.StageStatus_STAGE_FAILURE
	}

	if err := sm.UpdateAllTraffic(revision); err != nil {
		e.LogPersister.Errorf("Unable to configure all traffic to revision %s (%v)", revision, err)
		return model.StageStatus_STAGE_FAILURE
	}
	e.LogPersister.Info("Successfully generated the appropriate service manifest")

	e.LogPersister.Info("Start applying the service manifest")
	client, err := provider.DefaultRegistry().Client(ctx, e.cloudProviderName, e.cloudProviderConfig, e.Logger)
	if err != nil {
		e.LogPersister.Errorf("Unable to create ClourRun client for the provider (%v)", err)
		return model.StageStatus_STAGE_FAILURE
	}
	if _, err := client.Apply(ctx, sm); err != nil {
		e.LogPersister.Errorf("Failed to apply the service manifest (%v)", err)
		return model.StageStatus_STAGE_FAILURE
	}
	e.LogPersister.Info("Successfully applied the service manifest")

	return model.StageStatus_STAGE_SUCCESS
}

func (e *Executor) ensurePromote(ctx context.Context) model.StageStatus {
	var options = e.StageConfig.CloudRunPromoteStageOptions
	if options == nil {
		e.LogPersister.Errorf("Malformed configuration for stage %s", e.Stage.Name)
		return model.StageStatus_STAGE_FAILURE
	}

	// Determine the last deployed revision name.
	lastDeployedCommit := e.Deployment.RunningCommitHash
	if lastDeployedCommit == "" {
		e.LogPersister.Errorf("Unable to determine the last deployed commit")
	}

	lastDeployedServiceManifest, ok := e.loadLastDeployedServiceManifest()
	if !ok {
		return model.StageStatus_STAGE_FAILURE
	}

	lastDeployedRevisionName, err := provider.DecideRevisionName(lastDeployedServiceManifest, lastDeployedCommit)
	if err != nil {
		e.LogPersister.Errorf("Unable to decide the last deployed revision name for the commit %s (%v)", lastDeployedCommit, err)
		return model.StageStatus_STAGE_FAILURE
	}

	// Load triggered service manifest to apply.
	commit := e.Deployment.Trigger.Commit.Hash
	sm, ok := e.loadServiceManifest()
	if !ok {
		return model.StageStatus_STAGE_FAILURE
	}

	e.LogPersister.Infof("Generating a service manifest that configures traffic as: %d%% to new version, %d%% to old version", options.Percent, 100-options.Percent)
	revisionName, err := provider.DecideRevisionName(sm, commit)
	if err != nil {
		e.LogPersister.Errorf("Unable to decide revision name for the commit %s (%v)", commit, err)
		return model.StageStatus_STAGE_FAILURE
	}

	if err := sm.SetRevision(revisionName); err != nil {
		e.LogPersister.Errorf("Unable to set revision name to service manifest (%v)", err)
		return model.StageStatus_STAGE_FAILURE
	}

	revisions := []provider.RevisionTraffic{
		{
			RevisionName: revisionName,
			Percent:      options.Percent,
		},
		{
			RevisionName: lastDeployedRevisionName,
			Percent:      100 - options.Percent,
		},
	}
	if err := sm.UpdateTraffic(revisions); err != nil {
		e.LogPersister.Errorf("Unable to configure traffic (%v)", err)
		return model.StageStatus_STAGE_FAILURE
	}
	e.LogPersister.Info("Successfully generated the appropriate service manifest")

	e.LogPersister.Info("Start applying the service manifest")
	client, err := provider.DefaultRegistry().Client(ctx, e.cloudProviderName, e.cloudProviderConfig, e.Logger)
	if err != nil {
		e.LogPersister.Errorf("Unable to create ClourRun client for the provider (%v)", err)
		return model.StageStatus_STAGE_FAILURE
	}
	if _, err := client.Apply(ctx, sm); err != nil {
		e.LogPersister.Errorf("Failed to apply the service manifest (%v)", err)
		return model.StageStatus_STAGE_FAILURE
	}
	e.LogPersister.Info("Successfully applied the service manifest")

	// TODO: Wait to ensure the traffic was fully configured.
	return model.StageStatus_STAGE_SUCCESS
}

func (e *Executor) ensureRollback(ctx context.Context) model.StageStatus {
	commit := e.Deployment.RunningCommitHash
	if commit == "" {
		e.LogPersister.Errorf("Unable to determine the last deployed commit to rollback. It seems this is the first deployment.")
		return model.StageStatus_STAGE_FAILURE
	}

	sm, ok := e.loadLastDeployedServiceManifest()
	if !ok {
		return model.StageStatus_STAGE_FAILURE
	}

	e.LogPersister.Info("Generate a service manifest that configures all traffic to the last deployed revision")
	revision, err := provider.DecideRevisionName(sm, commit)
	if err != nil {
		e.LogPersister.Errorf("Unable to decide revision name for the commit %s (%v)", commit, err)
		return model.StageStatus_STAGE_FAILURE
	}

	if err := sm.SetRevision(revision); err != nil {
		e.LogPersister.Errorf("Unable to set revision name to service manifest (%v)", err)
		return model.StageStatus_STAGE_FAILURE
	}

	if err := sm.UpdateAllTraffic(revision); err != nil {
		e.LogPersister.Errorf("Unable to configure all traffic to revision %s (%v)", revision, err)
		return model.StageStatus_STAGE_FAILURE
	}
	e.LogPersister.Info("Successfully generated the appropriate service manifest")

	e.LogPersister.Info("Start applying the service manifest")
	client, err := provider.DefaultRegistry().Client(ctx, e.cloudProviderName, e.cloudProviderConfig, e.Logger)
	if err != nil {
		e.LogPersister.Errorf("Unable to create ClourRun client for the provider (%v)", err)
		return model.StageStatus_STAGE_FAILURE
	}
	if _, err := client.Apply(ctx, sm); err != nil {
		e.LogPersister.Errorf("Failed to apply the service manifest (%v)", err)
		return model.StageStatus_STAGE_FAILURE
	}
	e.LogPersister.Info("Successfully applied the service manifest")

	return model.StageStatus_STAGE_SUCCESS
}

func (e *Executor) loadServiceManifest() (provider.ServiceManifest, bool) {
	var (
		commit = e.Deployment.Trigger.Commit.Hash
		appDir = filepath.Join(e.RepoDir, e.Deployment.GitPath.Path)
	)

	e.LogPersister.Infof("Loading service manifest at the triggered commit %s", commit)
	sm, err := provider.LoadServiceManifest(appDir, e.config.Input.ServiceManifestFile)
	if err != nil {
		e.LogPersister.Errorf("Failed to load service manifest file (%v)", err)
		return provider.ServiceManifest{}, false
	}
	e.LogPersister.Info("Successfully loaded the service manifest")

	return sm, true
}

func (e *Executor) loadLastDeployedServiceManifest() (provider.ServiceManifest, bool) {
	var (
		commit = e.Deployment.RunningCommitHash
		appDir = filepath.Join(e.RunningRepoDir, e.Deployment.GitPath.Path)
	)

	e.LogPersister.Infof("Loading service manifest at the last deployed commit %s", commit)
	sm, err := provider.LoadServiceManifest(appDir, e.config.Input.ServiceManifestFile)
	if err != nil {
		e.LogPersister.Errorf("Failed to load service manifest file (%v)", err)
		return provider.ServiceManifest{}, false
	}
	e.LogPersister.Info("Successfully loaded the service manifest")

	return sm, true
}
