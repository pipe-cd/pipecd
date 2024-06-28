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

package cloudrun

import (
	"context"
	"strconv"
	"time"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/deploysource"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/executor"
	provider "github.com/pipe-cd/pipecd/pkg/app/pipedv1/platformprovider/cloudrun"
	"github.com/pipe-cd/pipecd/pkg/config"
	"github.com/pipe-cd/pipecd/pkg/model"

	"go.uber.org/zap"
)

const (
	promotePercentageMetadataKey = "promote-percentage"
	revisionCheckDuration        = 10 * time.Second
	revisionCheckTimeout         = 2 * time.Minute
)

type deployExecutor struct {
	executor.Input

	deploySource *deploysource.DeploySource
	appCfg       *config.CloudRunApplicationSpec
	client       provider.Client
}

func (e *deployExecutor) Execute(sig executor.StopSignal) model.StageStatus {
	ctx := sig.Context()
	ds, err := e.TargetDSP.GetReadOnly(ctx, e.LogPersister)
	if err != nil {
		e.LogPersister.Errorf("Failed to prepare target deploy source data (%v)", err)
		return model.StageStatus_STAGE_FAILURE
	}

	e.deploySource = ds
	e.appCfg = ds.ApplicationConfig.CloudRunApplicationSpec
	if e.appCfg == nil {
		e.LogPersister.Error("Malformed application configuration: missing CloudRunApplicationSpec")
		return model.StageStatus_STAGE_FAILURE
	}

	cpName, cpCfg, found := findPlatformProvider(&e.Input)
	if !found {
		return model.StageStatus_STAGE_FAILURE
	}

	e.client, err = provider.DefaultRegistry().Client(ctx, cpName, cpCfg, e.Logger)
	if err != nil {
		e.LogPersister.Errorf("Unable to create ClourRun client for the provider (%v)", err)
		return model.StageStatus_STAGE_FAILURE
	}

	var (
		originalStatus = e.Stage.Status
		status         model.StageStatus
	)

	switch model.Stage(e.Stage.Name) {
	case model.StageCloudRunSync:
		status = e.ensureSync(ctx)

	case model.StageCloudRunPromote:
		status = e.ensurePromote(ctx)

	default:
		e.LogPersister.Errorf("Unsupported stage %s for cloudrun application", e.Stage.Name)
		return model.StageStatus_STAGE_FAILURE
	}

	return executor.DetermineStageStatus(sig.Signal(), originalStatus, status)
}

func (e *deployExecutor) ensureSync(ctx context.Context) model.StageStatus {
	sm, ok := loadServiceManifest(&e.Input, e.appCfg.Input.ServiceManifestFile, e.deploySource)
	if !ok {
		return model.StageStatus_STAGE_FAILURE
	}

	revision, ok := decideRevisionName(sm, e.Deployment.Trigger.Commit.Hash, e.LogPersister)
	if !ok {
		return model.StageStatus_STAGE_FAILURE
	}

	traffics := []provider.RevisionTraffic{
		{
			RevisionName: revision,
			Percent:      100,
		},
	}
	if !configureServiceManifest(sm, revision, traffics, e.LogPersister) {
		return model.StageStatus_STAGE_FAILURE
	}

	// Add builtin labels for tracking application live state
	commit := e.Deployment.CommitHash()
	if !addBuiltinLabels(sm, commit, e.PipedConfig.PipedID, e.Deployment.ApplicationId, revision, e.LogPersister) {
		return model.StageStatus_STAGE_FAILURE
	}

	if !apply(ctx, e.client, sm, e.LogPersister) {
		return model.StageStatus_STAGE_FAILURE
	}

	if err := waitRevisionReady(
		ctx,
		e.client,
		revision,
		revisionCheckDuration,
		revisionCheckTimeout,
		e.LogPersister,
	); err != nil {
		return model.StageStatus_STAGE_FAILURE
	}

	return model.StageStatus_STAGE_SUCCESS
}

func (e *deployExecutor) ensurePromote(ctx context.Context) model.StageStatus {
	options := e.StageConfig.CloudRunPromoteStageOptions
	if options == nil {
		e.LogPersister.Errorf("Malformed configuration for stage %s", e.Stage.Name)
		return model.StageStatus_STAGE_FAILURE
	}
	metadata := map[string]string{
		promotePercentageMetadataKey: strconv.FormatInt(int64(options.Percent.Int()), 10),
	}
	if err := e.MetadataStore.Stage(e.Stage.Id).PutMulti(ctx, metadata); err != nil {
		e.Logger.Error("failed to save routing percentages to metadata", zap.Error(err))
	}

	// Loaded the last deployed data.
	if e.Deployment.RunningCommitHash == "" {
		e.LogPersister.Errorf("Unable to determine the last deployed commit")
		return model.StageStatus_STAGE_FAILURE
	}

	runningDS, err := e.RunningDSP.GetReadOnly(ctx, e.LogPersister)
	if err != nil {
		e.LogPersister.Errorf("Failed to prepare running deploy source data (%v)", err)
		return model.StageStatus_STAGE_FAILURE
	}

	runningAppCfg := runningDS.ApplicationConfig.CloudRunApplicationSpec
	if runningAppCfg == nil {
		e.LogPersister.Error("Malformed application configuration in running commit: missing CloudRunApplicationSpec")
		return model.StageStatus_STAGE_FAILURE
	}

	lastDeployedSM, ok := loadServiceManifest(&e.Input, runningAppCfg.Input.ServiceManifestFile, runningDS)
	if !ok {
		return model.StageStatus_STAGE_FAILURE
	}

	lastDeployedRevision, ok := decideRevisionName(lastDeployedSM, e.Deployment.RunningCommitHash, e.LogPersister)
	if !ok {
		return model.StageStatus_STAGE_FAILURE
	}

	// Load the service manifest at the target commit.
	sm, ok := loadServiceManifest(&e.Input, e.appCfg.Input.ServiceManifestFile, e.deploySource)
	if !ok {
		return model.StageStatus_STAGE_FAILURE
	}

	revision, ok := decideRevisionName(sm, e.Deployment.Trigger.Commit.Hash, e.LogPersister)
	if !ok {
		return model.StageStatus_STAGE_FAILURE
	}

	traffics := []provider.RevisionTraffic{
		{
			RevisionName: revision,
			Percent:      options.Percent.Int(),
		},
		{
			RevisionName: lastDeployedRevision,
			Percent:      100 - options.Percent.Int(),
		},
	}

	exist, err := revisionExists(ctx, e.client, revision, e.LogPersister)
	if err != nil {
		return model.StageStatus_STAGE_FAILURE
	}

	newRevision := revision
	if exist {
		newRevision = ""
		e.LogPersister.Infof("Revision %s was already registered", revision)
	}

	if !configureServiceManifest(sm, newRevision, traffics, e.LogPersister) {
		return model.StageStatus_STAGE_FAILURE
	}

	commit := e.Deployment.CommitHash()
	if !addBuiltinLabels(sm, commit, e.PipedConfig.PipedID, e.Deployment.ApplicationId, newRevision, e.LogPersister) {
		return model.StageStatus_STAGE_FAILURE
	}

	if !apply(ctx, e.client, sm, e.LogPersister) {
		return model.StageStatus_STAGE_FAILURE
	}

	if err := waitRevisionReady(
		ctx,
		e.client,
		revision,
		revisionCheckDuration,
		revisionCheckTimeout,
		e.LogPersister,
	); err != nil {
		return model.StageStatus_STAGE_FAILURE
	}

	return model.StageStatus_STAGE_SUCCESS
}
