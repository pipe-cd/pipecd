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

// Package livestatereporter provides a piped component
// that reports the changes as well as full snapshot about live state of registered applications.
package livestatereporter

import (
	"context"
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/deploysource"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin"
	"github.com/pipe-cd/pipecd/pkg/app/server/service/pipedservice"
	config "github.com/pipe-cd/pipecd/pkg/configv1"
	"github.com/pipe-cd/pipecd/pkg/git"
	"github.com/pipe-cd/pipecd/pkg/model"
	"github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1/livestate"
)

type applicationLister interface {
	List() []*model.Application
}

type apiClient interface {
	ReportApplicationLiveState(ctx context.Context, req *pipedservice.ReportApplicationLiveStateRequest, opts ...grpc.CallOption) (*pipedservice.ReportApplicationLiveStateResponse, error)
	ReportApplicationSyncState(ctx context.Context, req *pipedservice.ReportApplicationSyncStateRequest, opts ...grpc.CallOption) (*pipedservice.ReportApplicationSyncStateResponse, error)
}

// Reporter represents a component that reports the snapshot about live state of registered applications.
type Reporter interface {
	Run(ctx context.Context) error
}

type gitClient interface {
	Clone(ctx context.Context, repoID, remote, branch, destination string) (git.Repo, error)
}

type secretDecrypter interface {
	Decrypt(string) (string, error)
}

type reporter struct {
	snapshotFlushInterval time.Duration
	appLister             applicationLister
	apiClient             apiClient
	gitClient             gitClient
	repoMap               map[string]git.Repo
	pluginRegistry        plugin.PluginRegistry
	pipedConfig           *config.PipedSpec
	secretDecrypter       secretDecrypter
	workingDir            string
	logger                *zap.Logger
}

// NewReporter creates a new reporter.
func NewReporter(appLister applicationLister, apiClient apiClient, gitClient gitClient, pluginRegistry plugin.PluginRegistry, pipedConfig *config.PipedSpec, secretDecrypter secretDecrypter, logger *zap.Logger) (Reporter, error) {
	rlogger := logger.Named("live-state-reporter")

	workingDir, err := os.MkdirTemp("", "livestate-reporter-*")
	if err != nil {
		rlogger.Error("failed to create working directory", zap.Error(err))
		return nil, err
	}

	r := &reporter{
		snapshotFlushInterval: 1 * time.Minute,
		appLister:             appLister,
		apiClient:             apiClient,
		gitClient:             gitClient,
		repoMap:               make(map[string]git.Repo, 0),
		pluginRegistry:        pluginRegistry,
		pipedConfig:           pipedConfig,
		secretDecrypter:       secretDecrypter,
		workingDir:            workingDir,
		logger:                rlogger,
	}

	return r, nil
}

func (r *reporter) Run(ctx context.Context) error {
	r.logger.Info("start running app live state reporter", zap.Duration("snapshot-flush-interval", r.snapshotFlushInterval))

	snapshotTicker := time.NewTicker(r.snapshotFlushInterval)
	defer snapshotTicker.Stop()

	for {
		select {
		case <-snapshotTicker.C:
			r.flushSnapshots(ctx)

		case <-ctx.Done():
			r.logger.Info("app live state reporter has been stopped")
			return nil
		}
	}
}

func (r *reporter) flushSnapshots(ctx context.Context) {
	r.logger.Info("start flushing snapshots")

	r.logger.Info("updating git repositories")
	for id, cfgRepo := range r.pipedConfig.GetRepositoryMap() {
		if _, ok := r.repoMap[id]; ok {
			r.logger.Info("the repo is cloned, so pulling repository", zap.String("repo-id", id))
			err := r.repoMap[id].Pull(ctx, cfgRepo.Branch)
			if err != nil {
				r.logger.Error("failed to pull repository", zap.String("repo-id", id), zap.Error(err))
			}
			continue
		}

		r.logger.Info("the repo is not cloned, so cloning repository", zap.String("repo-id", id))
		repo, err := r.gitClient.Clone(ctx, id, cfgRepo.Remote, cfgRepo.Branch, fmt.Sprintf("%s/%s", r.workingDir, id))
		if err != nil {
			r.logger.Error("failed to clone repository", zap.String("repo-id", id), zap.Error(err))
			continue
		}
		r.repoMap[id] = repo
	}

	apps := r.appLister.List()

	// TODO: Set the limit based on the piped config
	limit := max(runtime.GOMAXPROCS(0)/2, 1)

	appsPerReporter := max(len(apps)/limit, 1)

	appsGrouped := make([][]*model.Application, 0)
	for i := 0; i < len(apps); i += appsPerReporter {
		end := i + appsPerReporter
		if end > len(apps) {
			end = len(apps)
		}
		appsGrouped = append(appsGrouped, apps[i:end])
	}

	var wg sync.WaitGroup

	r.logger.Info("flushing snapshots", zap.Int("total-applications", len(apps)), zap.Int("parallel-count", len(appsGrouped)))
	for _, apps := range appsGrouped {
		wg.Add(1)
		go func() {
			r.flushAll(ctx, apps, r.repoMap)
			wg.Done()
		}()
	}

	wg.Wait()

	r.logger.Info("finished flushing snapshots")
}

func (r *reporter) flushAll(ctx context.Context, apps []*model.Application, repoMap map[string]git.Repo) {
	for _, app := range apps {
		repo, ok := repoMap[app.GitPath.Repo.Id]
		if !ok {
			r.logger.Error("failed to find repository for application", zap.String("application-id", app.Id))
			continue
		}

		if err := r.flush(ctx, app, repo); err != nil {
			r.logger.Error("failed to flush application live state", zap.String("application-id", app.Id), zap.Error(err))
			continue
		}
	}
}

func (r *reporter) flush(ctx context.Context, app *model.Application, repo git.Repo) error {
	dir, err := os.MkdirTemp(r.workingDir, fmt.Sprintf("app-%s-*", app.Id))
	if err != nil {
		r.logger.Error("failed to create temporary directory", zap.Error(err))
		return err
	}

	headCommit, err := repo.GetLatestCommit(ctx)
	if err != nil {
		r.logger.Error("failed to get latest commit", zap.Error(err))
		return err
	}

	dsp := deploysource.NewProvider(dir, deploysource.NewLocalSourceCloner(repo, "target", headCommit.Hash), app.GitPath, r.secretDecrypter)
	ds, err := dsp.Get(ctx, io.Discard)
	if err != nil {
		r.logger.Error("failed to get deploy source", zap.String("application-id", app.Id), zap.Error(err))
		return err
	}

	cfg, err := config.DecodeYAML[*config.GenericApplicationSpec](ds.ApplicationConfig)
	if err != nil {
		r.logger.Error("unable to parse application config", zap.Error(err))
		return err
	}

	pluginClis, err := r.pluginRegistry.GetPluginClientsByAppConfig(cfg.Spec)
	if err != nil {
		r.logger.Error("failed to get plugin clients", zap.Error(err))
		return err
	}

	// Get the application live state from the plugins.
	resourceStates := make([]*model.ResourceState, 0)
	syncStates := make([]*model.ApplicationSyncState, 0)
	for _, pluginClient := range pluginClis {
		res, err := pluginClient.GetLivestate(ctx, &livestate.GetLivestateRequest{
			PipedId:         app.GetPipedId(),
			ApplicationId:   app.GetId(),
			ApplicationName: app.GetName(),
			DeploySource:    ds.ToPluginDeploySource(),
			DeployTargets:   app.GetDeployTargets(),
		})
		if err != nil {
			r.logger.Info(fmt.Sprintf("no app state of application %s to report", app.Id))
			return err
		}

		resourceStates = append(resourceStates, res.GetApplicationLiveState().GetResources()...)
		syncStates = append(syncStates, res.GetSyncState())
	}

	// Report the application live state to the control plane.
	snapshot := &model.ApplicationLiveStateSnapshot{
		ApplicationId: app.Id,
		PipedId:       app.PipedId,
		ProjectId:     app.ProjectId,
		Kind:          app.Kind,
		ApplicationLiveState: &model.ApplicationLiveState{
			Resources: resourceStates,
		},
		Version: &model.ApplicationLiveStateVersion{
			Timestamp: time.Now().Unix(),
		},
	}
	snapshot.DetermineApplicationHealthStatus()

	if _, err := r.apiClient.ReportApplicationLiveState(ctx, &pipedservice.ReportApplicationLiveStateRequest{
		Snapshot: snapshot,
	}); err != nil {
		r.logger.Error("failed to report application live state",
			zap.String("application-id", app.Id),
			zap.Error(err),
		)
		return err
	}

	// Report the application sync state to the control plane.
	if _, err := r.apiClient.ReportApplicationSyncState(ctx, &pipedservice.ReportApplicationSyncStateRequest{
		ApplicationId: app.Id,
		State:         model.MergeApplicationSyncState(syncStates),
	}); err != nil {
		r.logger.Error("failed to report application live state",
			zap.String("application-id", app.Id),
			zap.Error(err),
		)
		return err
	}

	r.logger.Info(fmt.Sprintf("successfully reported application live state for application: %s", app.Id))

	os.RemoveAll(dir)
	return nil
}
