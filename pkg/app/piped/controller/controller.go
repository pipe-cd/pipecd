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

// Package controller provides a piped component
// that handles all of the not completed deployments by managing a pool of planners and schedulers.
// Whenever a new PENDING deployment is detected, controller spawns a new planner for deciding
// the deployment pipeline and update the deployment status to PLANNED.
// Whenever a new PLANNED deployment is detected, controller spawns a new scheduler
// for scheduling and running its pipeline executors.
package controller

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes/wrappers"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/pipe-cd/pipe/pkg/app/api/service/pipedservice"
	"github.com/pipe-cd/pipe/pkg/app/piped/logpersister"
	"github.com/pipe-cd/pipe/pkg/cache"
	"github.com/pipe-cd/pipe/pkg/config"
	"github.com/pipe-cd/pipe/pkg/git"
	"github.com/pipe-cd/pipe/pkg/model"
)

type apiClient interface {
	GetMostRecentDeployment(ctx context.Context, req *pipedservice.GetMostRecentDeploymentRequest, opts ...grpc.CallOption) (*pipedservice.GetMostRecentDeploymentResponse, error)
	ReportDeploymentPlanned(ctx context.Context, req *pipedservice.ReportDeploymentPlannedRequest, opts ...grpc.CallOption) (*pipedservice.ReportDeploymentPlannedResponse, error)
	ReportDeploymentStatusChanged(ctx context.Context, req *pipedservice.ReportDeploymentStatusChangedRequest, opts ...grpc.CallOption) (*pipedservice.ReportDeploymentStatusChangedResponse, error)
	ReportDeploymentCompleted(ctx context.Context, req *pipedservice.ReportDeploymentCompletedRequest, opts ...grpc.CallOption) (*pipedservice.ReportDeploymentCompletedResponse, error)
	SaveDeploymentMetadata(ctx context.Context, req *pipedservice.SaveDeploymentMetadataRequest, opts ...grpc.CallOption) (*pipedservice.SaveDeploymentMetadataResponse, error)
	ReportApplicationMostRecentDeployment(ctx context.Context, req *pipedservice.ReportApplicationMostRecentDeploymentRequest, opts ...grpc.CallOption) (*pipedservice.ReportApplicationMostRecentDeploymentResponse, error)

	ReportStageStatusChanged(ctx context.Context, req *pipedservice.ReportStageStatusChangedRequest, opts ...grpc.CallOption) (*pipedservice.ReportStageStatusChangedResponse, error)
	SaveStageMetadata(ctx context.Context, req *pipedservice.SaveStageMetadataRequest, opts ...grpc.CallOption) (*pipedservice.SaveStageMetadataResponse, error)
	ReportStageLogs(ctx context.Context, req *pipedservice.ReportStageLogsRequest, opts ...grpc.CallOption) (*pipedservice.ReportStageLogsResponse, error)
	ReportStageLogsFromLastCheckpoint(ctx context.Context, in *pipedservice.ReportStageLogsFromLastCheckpointRequest, opts ...grpc.CallOption) (*pipedservice.ReportStageLogsFromLastCheckpointResponse, error)
}

type gitClient interface {
	Clone(ctx context.Context, repoID, remote, branch, destination string) (git.Repo, error)
}

type deploymentLister interface {
	ListPendings() []*model.Deployment
	ListPlanneds() []*model.Deployment
	ListRunnings() []*model.Deployment
}

type commandLister interface {
	ListDeploymentCommands() []model.ReportableCommand
	ListStageCommands(deploymentID, stageID string) []model.ReportableCommand
}

type applicationLister interface {
	Get(id string) (*model.Application, bool)
}

type DeploymentController interface {
	Run(ctx context.Context) error
}

var (
	plannerStaleDuration   = time.Hour
	schedulerStaleDuration = time.Hour
)

type controller struct {
	apiClient         apiClient
	gitClient         gitClient
	deploymentLister  deploymentLister
	commandLister     commandLister
	applicationLister applicationLister
	pipedConfig       *config.PipedSpec
	appManifestsCache cache.Cache
	logPersister      logpersister.Persister

	// Map from application ID to the planner
	// of a pending deployment of that application.
	planners map[string]*planner
	// Map from deployment ID to the completion time
	// of the done planners.
	// Because when the deployment lister returns not fresh data
	// we use this to ignore the ones that have been handled previously.
	donePlanners map[string]time.Time
	// Map from application ID to the scheduler
	// of a running deployment of that application.
	schedulers map[string]*scheduler
	// Map from deployment ID to the completion time
	// of the done schedulers.
	doneSchedulers map[string]time.Time
	// Map from application ID to its most recent successful commit hash.
	mostRecentSuccessfulCommits map[string]string
	// WaitGroup for waiting the completions of all planners, schedulers.
	wg sync.WaitGroup

	workspaceDir string
	syncInternal time.Duration
	gracePeriod  time.Duration
	logger       *zap.Logger
}

// NewController creates a new instance for DeploymentController.
func NewController(
	apiClient apiClient,
	gitClient gitClient,
	deploymentLister deploymentLister,
	commandLister commandLister,
	applicationLister applicationLister,
	pipedConfig *config.PipedSpec,
	appManifestsCache cache.Cache,
	gracePeriod time.Duration,
	logger *zap.Logger,
) DeploymentController {

	var (
		lp = logpersister.NewPersister(apiClient, logger)
		lg = logger.Named("controller")
	)
	return &controller{
		apiClient:         apiClient,
		gitClient:         gitClient,
		deploymentLister:  deploymentLister,
		commandLister:     commandLister,
		applicationLister: applicationLister,
		appManifestsCache: appManifestsCache,
		pipedConfig:       pipedConfig,
		logPersister:      lp,

		planners:                    make(map[string]*planner),
		donePlanners:                make(map[string]time.Time),
		schedulers:                  make(map[string]*scheduler),
		doneSchedulers:              make(map[string]time.Time),
		mostRecentSuccessfulCommits: make(map[string]string),

		syncInternal: 10 * time.Second,
		gracePeriod:  gracePeriod,
		logger:       lg,
	}
}

// Run starts running controller until the specified context has done.
// This also waits for its cleaning up before returning.
func (c *controller) Run(ctx context.Context) error {
	c.logger.Info("start running controller")

	// Make sure the existence of the workspace directory.
	// Each scheduler/deployment will have an working directory inside this workspace.
	dir, err := ioutil.TempDir("", "workspace")
	if err != nil {
		c.logger.Error("failed to create workspace directory", zap.Error(err))
		return err
	}
	c.workspaceDir = dir
	c.logger.Info(fmt.Sprintf("workspace directory was configured to %s", c.workspaceDir))

	// Start running log persister to buffer and flush the log blocks.
	// We do not use the passed ctx directly because we want log persister
	// component to be stopped at the last order to avoid lossing log from other components.
	var (
		lpStoppedCh     = make(chan error, 1)
		lpCtx, lpCancel = context.WithCancel(context.Background())
	)
	go func() {
		lpStoppedCh <- c.logPersister.Run(lpCtx)
		close(lpStoppedCh)
	}()

	ticker := time.NewTicker(c.syncInternal)
	defer ticker.Stop()
	c.logger.Info("start syncing planners and schedulers")

L:
	for {
		select {
		case <-ctx.Done():
			break L
		case <-ticker.C:
			// This must be called before syncPlanner because
			// after piped is restarted all running deployments need to be loaded firstly.
			c.syncScheduler(ctx)
			c.syncPlanner(ctx)
			c.checkCommands()
		}
	}

	c.logger.Info("waiting for stopping all planners and schedulers")
	c.wg.Wait()

	// Stop log persiter and wait for its stopping.
	lpCancel()
	err = <-lpStoppedCh

	c.logger.Info("controller has been stopped")
	return err
}

// checkCommands lists all unhandled commands for running deployments
// and forwards them to their planners and schedulers.
func (c *controller) checkCommands() {
	commands := c.commandLister.ListDeploymentCommands()
	for _, cmd := range commands {
		if cmd.GetCancelDeployment() == nil {
			continue
		}
		if planner, ok := c.planners[cmd.ApplicationId]; ok && planner.ID() == cmd.DeploymentId {
			planner.Cancel(cmd)
		}
		if scheduler, ok := c.schedulers[cmd.ApplicationId]; ok && scheduler.ID() == cmd.DeploymentId {
			scheduler.Cancel(cmd)
		}
	}
}

// syncPlanner adds new planner for newly PENDING deployments.
func (c *controller) syncPlanner(ctx context.Context) error {
	// Remove stale planners.
	for id, t := range c.donePlanners {
		if time.Since(t) < plannerStaleDuration {
			continue
		}
		delete(c.donePlanners, id)
	}

	for id, p := range c.planners {
		if !p.IsDone() {
			continue
		}
		c.logger.Info("deleted done planner",
			zap.String("deployment-id", p.ID()),
			zap.String("application-id", id),
			zap.Int("planner-count", len(c.planners)),
		)
		c.donePlanners[p.ID()] = p.DoneTimestamp()
		delete(c.planners, id)
	}

	// Add missing planners.
	pendings := c.deploymentLister.ListPendings()
	if len(pendings) == 0 {
		return nil
	}

	c.logger.Info(fmt.Sprintf("there are %d pending deployments for planning", len(pendings)),
		zap.Int("planner-count", len(c.planners)),
	)

	pendingByApp := make(map[string]*model.Deployment, len(pendings))
	for _, d := range pendings {
		appID := d.ApplicationId
		// Ignore already processed one.
		if _, ok := c.donePlanners[d.Id]; ok {
			continue
		}
		// For each application, only one deployment can be planned at the same time.
		if _, ok := c.planners[appID]; ok {
			continue
		}
		// If this application is deploying, no other deployments can be added to plan.
		if _, ok := c.schedulers[appID]; ok {
			continue
		}
		// Choose the oldest PENDING deployment of the application to plan.
		if pre, ok := pendingByApp[appID]; ok {
			if !d.TriggerBefore(pre) {
				continue
			}
		}
		pendingByApp[appID] = d
	}

	for appID, d := range pendingByApp {
		planner, err := c.startNewPlanner(ctx, d)
		if err != nil {
			continue
		}
		c.planners[appID] = planner
	}

	return nil
}

func (c *controller) startNewPlanner(ctx context.Context, d *model.Deployment) (*planner, error) {
	logger := c.logger.With(
		zap.String("deployment-id", d.Id),
		zap.String("application-id", d.ApplicationId),
	)
	logger.Info("will add a new planner")

	// Ensure the existence of the working directory for the deployment.
	workingDir, err := ioutil.TempDir(c.workspaceDir, d.Id+"-planner-*")
	if err != nil {
		logger.Error("failed to create working directory for planner", zap.Error(err))
		return nil, err
	}
	logger.Info("created working directory for planner", zap.String("working-dir", workingDir))

	// The most recent successful commit is saved in memory.
	// But when the piped is restarted that data will be cleared too.
	// So in that case, we have to use the API to check.
	commit := c.mostRecentSuccessfulCommits[d.ApplicationId]
	if commit == "" {
		mostRecent, err := c.getMostRecentlySuccessfulDeployment(ctx, d.ApplicationId)
		if err == nil {
			commit = mostRecent.CommitHash()
			c.mostRecentSuccessfulCommits[d.ApplicationId] = commit
		} else if status.Code(err) == codes.NotFound {
			logger.Info("there is no previous successful commit for this application")
		} else {
			logger.Error("unabled to get the most recent successful deployment", zap.Error(err))
		}
	}

	planner := newPlanner(
		d,
		commit,
		workingDir,
		c.apiClient,
		c.gitClient,
		c.pipedConfig,
		c.appManifestsCache,
		c.logger,
	)

	cleanup := func() {
		logger.Info("cleaning up working directory for planner", zap.String("working-dir", workingDir))
		if err := os.RemoveAll(workingDir); err != nil {
			logger.Warn("failed to clean working directory",
				zap.String("working-dir", workingDir),
				zap.Error(err),
			)
		}
	}

	// Start running planner.
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		defer cleanup()
		planner.Run(ctx)
	}()

	return planner, nil
}

// syncScheduler adds new scheduler for newly PLANNED/RUNNING deployments
// as well as removes the schedulers for the completed deployments.
func (c *controller) syncScheduler(ctx context.Context) error {
	// Update the most recent successful commit hashes.
	for id, s := range c.schedulers {
		if !s.IsDone() {
			continue
		}
		if s.DoneDeploymentStatus() != model.DeploymentStatus_DEPLOYMENT_SUCCESS {
			continue
		}
		c.mostRecentSuccessfulCommits[id] = s.CommitHash()
	}

	// Remove done schedulers.
	for id, t := range c.doneSchedulers {
		if time.Since(t) < schedulerStaleDuration {
			continue
		}
		delete(c.doneSchedulers, id)
	}

	for id, s := range c.schedulers {
		if !s.IsDone() {
			continue
		}
		c.logger.Info("deleted done scheduler",
			zap.String("deployment-id", s.ID()),
			zap.String("application-id", id),
			zap.Int("scheduler-count", len(c.schedulers)),
		)
		c.doneSchedulers[s.ID()] = s.DoneTimestamp()
		delete(c.schedulers, id)
	}

	// Add missing schedulers.
	planneds := c.deploymentLister.ListPlanneds()
	runnings := c.deploymentLister.ListRunnings()
	runnings = append(runnings, planneds...)

	if len(runnings) == 0 {
		return nil
	}

	c.logger.Info(fmt.Sprintf("there are %d planned/running deployments for scheduling", len(runnings)),
		zap.Int("scheduler-count", len(c.schedulers)),
	)

	for _, d := range runnings {
		// Ignore already processed one.
		if _, ok := c.doneSchedulers[d.Id]; ok {
			continue
		}
		if s, ok := c.schedulers[d.ApplicationId]; ok {
			if s.ID() != d.Id {
				c.logger.Warn("detected an application that has more than one running deployments",
					zap.String("application-id", d.ApplicationId),
					zap.String("handling-deployment-id", s.ID()),
					zap.String("deployment-id", d.Id),
				)
			}
			continue
		}
		s, err := c.startNewScheduler(ctx, d)
		if err != nil {
			continue
		}
		c.schedulers[d.ApplicationId] = s
		c.logger.Info("added a new scheduler",
			zap.String("deployment-id", d.Id),
			zap.String("application-id", d.ApplicationId),
			zap.Int("scheduler-count", len(c.schedulers)),
		)
	}

	return nil
}

// startNewScheduler creates and starts running a new scheduler
// for a specific PLANNED deployment.
// This adds the newly created one to the scheduler list
// for tracking its lifetime periodically later.
func (c *controller) startNewScheduler(ctx context.Context, d *model.Deployment) (*scheduler, error) {
	logger := c.logger.With(
		zap.String("deployment-id", d.Id),
		zap.String("application-id", d.ApplicationId),
	)
	logger.Info("will add a new scheduler")

	// Ensure the existence of the working directory for the deployment.
	workingDir, err := ioutil.TempDir(c.workspaceDir, d.Id+"-scheduler-*")
	if err != nil {
		logger.Error("failed to create working directory for scheduler", zap.Error(err))
		return nil, err
	}
	logger.Info("created working directory for scheduler", zap.String("working-dir", workingDir))

	// Create a new scheduler and append to the list for tracking.
	scheduler := newScheduler(
		d,
		workingDir,
		c.apiClient,
		c.gitClient,
		c.commandLister,
		c.applicationLister,
		c.logPersister,
		c.pipedConfig,
		c.appManifestsCache,
		c.logger,
	)

	cleanup := func() {
		logger.Info("cleaning up working directory for scheduler", zap.String("working-dir", workingDir))
		err := os.RemoveAll(workingDir)
		if err == nil {
			return
		}
		logger.Warn("failed to clean working directory",
			zap.String("working-dir", workingDir),
			zap.Error(err),
		)
	}

	// Start running scheduler.
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		defer cleanup()
		scheduler.Run(ctx)
	}()

	return scheduler, nil
}

func (c *controller) getMostRecentlySuccessfulDeployment(ctx context.Context, applicationID string) (*model.Deployment, error) {
	var (
		err   error
		resp  *pipedservice.GetMostRecentDeploymentResponse
		retry = pipedservice.NewRetry(3)
		req   = &pipedservice.GetMostRecentDeploymentRequest{
			ApplicationId: applicationID,
			Status: &wrappers.Int32Value{
				Value: int32(model.DeploymentStatus_DEPLOYMENT_SUCCESS),
			},
		}
	)

	for retry.WaitNext(ctx) {
		if resp, err = c.apiClient.GetMostRecentDeployment(ctx, req); err == nil {
			return resp.Deployment, nil
		}
		if !pipedservice.Retriable(err) {
			return nil, err
		}
	}
	return nil, err
}

// prepareDeployRepository clones repository and checkouts to the target revision.
func prepareDeployRepository(ctx context.Context, d *model.Deployment, gitClient gitClient, repoDirPath string, pipedConfig *config.PipedSpec) (git.Repo, error) {
	var (
		appID       = d.ApplicationId
		repoID      = d.GitPath.RepoId
		revision    = d.Trigger.Commit.Hash
		repoCfg, ok = pipedConfig.GetRepository(repoID)
	)
	if !ok {
		err := fmt.Errorf("no registered repository id %s for application %s", repoID, appID)
		return nil, err
	}

	gitRepo, err := gitClient.Clone(ctx, repoCfg.RepoID, repoCfg.Remote, repoCfg.Branch, repoDirPath)
	if err != nil {
		err = fmt.Errorf("failed to clone repository %s for application %s (%v)", repoID, appID, err)
		return nil, err
	}

	err = gitRepo.Checkout(ctx, revision)
	if err != nil {
		err = fmt.Errorf("failed to clone repository %s for application %s (%v)", repoID, appID, err)
		return nil, err
	}

	return gitRepo, nil
}

func loadDeploymentConfiguration(repoPath string, d *model.Deployment) (*config.Config, error) {
	path := filepath.Join(repoPath, d.GitPath.GetDeploymentConfigFilePath(config.DeploymentConfigurationFileName))
	cfg, err := config.LoadFromYAML(path)
	if err != nil {
		return nil, err
	}
	if appKind, ok := config.ToApplicationKind(cfg.Kind); !ok || appKind != d.Kind {
		return nil, fmt.Errorf("application in deployment configuration file is not match, got: %s, expected: %s", appKind, d.Kind)
	}
	return cfg, nil
}
