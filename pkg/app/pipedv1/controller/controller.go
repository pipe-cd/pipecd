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
	"os"
	"sync"
	"time"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/controller/controllermetrics"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/metadatastore"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin"
	"github.com/pipe-cd/pipecd/pkg/app/server/service/pipedservice"
	"github.com/pipe-cd/pipecd/pkg/git"
	"github.com/pipe-cd/pipecd/pkg/model"
)

type apiClient interface {
	GetApplicationMostRecentDeployment(ctx context.Context, req *pipedservice.GetApplicationMostRecentDeploymentRequest, opts ...grpc.CallOption) (*pipedservice.GetApplicationMostRecentDeploymentResponse, error)
	ReportApplicationDeployingStatus(ctx context.Context, req *pipedservice.ReportApplicationDeployingStatusRequest, opts ...grpc.CallOption) (*pipedservice.ReportApplicationDeployingStatusResponse, error)
	ReportDeploymentPlanned(ctx context.Context, req *pipedservice.ReportDeploymentPlannedRequest, opts ...grpc.CallOption) (*pipedservice.ReportDeploymentPlannedResponse, error)
	ReportDeploymentStatusChanged(ctx context.Context, req *pipedservice.ReportDeploymentStatusChangedRequest, opts ...grpc.CallOption) (*pipedservice.ReportDeploymentStatusChangedResponse, error)
	ReportDeploymentCompleted(ctx context.Context, req *pipedservice.ReportDeploymentCompletedRequest, opts ...grpc.CallOption) (*pipedservice.ReportDeploymentCompletedResponse, error)
	SaveDeploymentMetadata(ctx context.Context, req *pipedservice.SaveDeploymentMetadataRequest, opts ...grpc.CallOption) (*pipedservice.SaveDeploymentMetadataResponse, error)
	ReportApplicationMostRecentDeployment(ctx context.Context, req *pipedservice.ReportApplicationMostRecentDeploymentRequest, opts ...grpc.CallOption) (*pipedservice.ReportApplicationMostRecentDeploymentResponse, error)

	ReportStageStatusChanged(ctx context.Context, req *pipedservice.ReportStageStatusChangedRequest, opts ...grpc.CallOption) (*pipedservice.ReportStageStatusChangedResponse, error)
	SaveStageMetadata(ctx context.Context, req *pipedservice.SaveStageMetadataRequest, opts ...grpc.CallOption) (*pipedservice.SaveStageMetadataResponse, error)

	InChainDeploymentPlannable(ctx context.Context, in *pipedservice.InChainDeploymentPlannableRequest, opts ...grpc.CallOption) (*pipedservice.InChainDeploymentPlannableResponse, error)
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

type notifier interface {
	Notify(event model.NotificationEvent)
}

type secretDecrypter interface {
	Decrypt(string) (string, error)
}

type DeploymentController interface {
	Run(ctx context.Context) error
}

var (
	plannerStaleDuration   = time.Hour
	schedulerStaleDuration = time.Hour
)

type controller struct {
	apiClient             apiClient
	gitClient             gitClient
	deploymentLister      deploymentLister
	commandLister         commandLister
	notifier              notifier
	secretDecrypter       secretDecrypter
	metadataStoreRegistry metadatastore.MetadataStoreRegistry

	// The registry of all plugins.
	pluginRegistry plugin.PluginRegistry

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
	// Map from application ID to its most recently successful commit hash.
	mostRecentlySuccessfulCommits         map[string]string
	mostRecentlySuccessfulConfigFilenames map[string]string
	// WaitGroup for waiting the completions of all planners, schedulers.
	wg sync.WaitGroup

	workspaceDir   string
	syncInternal   time.Duration
	gracePeriod    time.Duration
	logger         *zap.Logger
	tracerProvider trace.TracerProvider
}

// NewController creates a new instance for DeploymentController.
func NewController(
	apiClient apiClient,
	gitClient gitClient,
	pluginRegistry plugin.PluginRegistry,
	deploymentLister deploymentLister,
	commandLister commandLister,
	notifier notifier,
	secretDecrypter secretDecrypter,
	metadataStoreRegistry metadatastore.MetadataStoreRegistry,
	gracePeriod time.Duration,
	logger *zap.Logger,
	tracerProvider trace.TracerProvider,
) DeploymentController {

	return &controller{
		apiClient:             apiClient,
		gitClient:             gitClient,
		pluginRegistry:        pluginRegistry,
		deploymentLister:      deploymentLister,
		commandLister:         commandLister,
		notifier:              notifier,
		secretDecrypter:       secretDecrypter,
		metadataStoreRegistry: metadataStoreRegistry,

		planners:                              make(map[string]*planner),
		donePlanners:                          make(map[string]time.Time),
		schedulers:                            make(map[string]*scheduler),
		doneSchedulers:                        make(map[string]time.Time),
		mostRecentlySuccessfulCommits:         make(map[string]string),
		mostRecentlySuccessfulConfigFilenames: make(map[string]string),

		syncInternal:   10 * time.Second,
		gracePeriod:    gracePeriod,
		logger:         logger.Named("controller"),
		tracerProvider: tracerProvider,
	}
}

// Run starts running controller until the specified context has done.
// This also waits for its cleaning up before returning.
func (c *controller) Run(ctx context.Context) error {
	c.logger.Info("start running controller")

	// Make sure the existence of the workspace directory.
	// Each planner/scheduler will have a working directory inside this workspace.
	dir, err := os.MkdirTemp("", "workspace")
	if err != nil {
		c.logger.Error("failed to create workspace directory", zap.Error(err))
		return err
	}
	c.workspaceDir = dir
	c.logger.Info(fmt.Sprintf("workspace directory was configured to %s", c.workspaceDir))

	ticker := time.NewTicker(c.syncInternal)
	defer ticker.Stop()
	c.logger.Info("start syncing planners and schedulers")

	for {
		select {
		case <-ctx.Done():
			return c.shutdown()

		case <-ticker.C:
			// syncSchedulers must be called before syncPlanners because
			// after piped is restarted all running deployments need to be loaded firstly.
			c.syncSchedulers(ctx)
			c.syncPlanners(ctx)
			c.checkCommands()
		}
	}
}

// shutdown waits for stopping all planners and schedulers.
//
//nolint:unparam // returns error for future compatibility.
func (c *controller) shutdown() error {
	c.logger.Info("waiting for stopping all planners and schedulers")
	c.wg.Wait()
	c.logger.Info("controller has been stopped")
	return nil
}

// checkCommands lists all unhandled commands for running deployments
// and forwards them to their planners and schedulers.
func (c *controller) checkCommands() {
	commands := c.commandLister.ListDeploymentCommands()
	for _, cmd := range commands {
		if cmd.GetCancelDeployment() == nil {
			continue
		}

		var handled bool
		if planner, ok := c.planners[cmd.ApplicationId]; ok && planner.ID() == cmd.DeploymentId {
			handled = true
			planner.Cancel(cmd)
			c.logger.Info("a command CancelDeployment was forwarded to its planner",
				zap.String("app", cmd.ApplicationId),
				zap.String("deployment", cmd.DeploymentId),
			)
		}

		if scheduler, ok := c.schedulers[cmd.ApplicationId]; ok && scheduler.ID() == cmd.DeploymentId {
			handled = true
			scheduler.Cancel(cmd)
			c.logger.Info("a command CancelDeployment was forwarded to its scheduler",
				zap.String("app", cmd.ApplicationId),
				zap.String("deployment", cmd.DeploymentId),
			)
		}

		if !handled {
			c.logger.Info("a command CancelDeployment is still not handled",
				zap.String("app", cmd.ApplicationId),
				zap.String("deployment", cmd.DeploymentId),
			)
		}
	}
}

// syncPlanners adds new planner for newly PENDING deployments.
func (c *controller) syncPlanners(ctx context.Context) {
	// Remove stale planners from the recently completed list.
	for id, t := range c.donePlanners {
		if time.Since(t) >= plannerStaleDuration {
			delete(c.donePlanners, id)
		}
	}

	// Find all completed ones and add them to donePlaners list.
	for id, p := range c.planners {
		if !p.IsDone() {
			continue
		}
		c.logger.Info("deleted done planner",
			zap.String("deployment", p.ID()),
			zap.String("app", id),
			zap.Int("count", len(c.planners)),
		)
		c.donePlanners[p.ID()] = p.DoneTimestamp()
		delete(c.planners, id)

		// Application will be marked as NOT deploying when planner's deployment was completed.
		if p.DoneDeploymentStatus().IsCompleted() {
			if err := reportApplicationDeployingStatus(ctx, c.apiClient, id, false); err != nil {
				c.logger.Error("failed to mark application as NOT deploying",
					zap.String("deployment", p.ID()),
					zap.String("app", id),
					zap.Error(err),
				)
			}
		}
	}

	// Add missing planners.
	pendings := c.deploymentLister.ListPendings()
	if len(pendings) == 0 {
		return
	}

	c.logger.Info(fmt.Sprintf("there are %d pending deployments for planning", len(pendings)),
		zap.Int("count", len(c.planners)),
	)

	pendingByApp := make(map[string]*model.Deployment, len(pendings))
	for _, d := range pendings {
		appID := d.ApplicationId
		// Ignore already processed one.
		if _, ok := c.donePlanners[d.Id]; ok {
			c.logger.Info("ignore planning because it was already processed",
				zap.String("deployment", d.Id),
				zap.String("app", d.ApplicationId),
			)
			continue
		}
		// For each application, only one deployment can be planned at the same time.
		if p, ok := c.planners[appID]; ok {
			c.logger.Info("temporarily skip planning because another deployment is planning",
				zap.String("deployment", d.Id),
				zap.String("app", d.ApplicationId),
				zap.String("executing-deployment", p.deployment.Id),
			)
			continue
		}
		// If this application is deploying, no other deployments can be added to plan.
		if s, ok := c.schedulers[appID]; ok {
			c.logger.Info("temporarily skip planning because another deployment is running",
				zap.String("deployment", d.Id),
				zap.String("app", d.ApplicationId),
				zap.String("handling-deployment", s.deployment.Id),
			)
			continue
		}
		// Choose the oldest PENDING deployment of the application to plan.
		if pre, ok := pendingByApp[appID]; ok && !d.TriggerBefore(pre) {
			continue
		}
		controllermetrics.UpdateDeploymentStatus(d, d.Status)
		pendingByApp[appID] = d
	}

	for appID, d := range pendingByApp {
		plannable, cancel, cancelReason, err := c.shouldStartPlanningDeployment(ctx, d)
		if err != nil {
			c.logger.Error("failed to check deployment plannability",
				zap.String("deployment", d.Id),
				zap.String("app", d.ApplicationId),
				zap.Error(err),
			)
			continue
		}

		if cancel {
			if err = c.cancelDeployment(ctx, d, cancelReason); err != nil {
				c.logger.Error("failed to cancel deployment",
					zap.String("deployment", d.Id),
					zap.String("app", d.ApplicationId),
					zap.Error(err),
				)
			}
			continue
		}

		if !plannable {
			if d.IsInChainDeployment() {
				c.logger.Info("unable to start planning deployment, probably locked by the previous block in its deployment chain",
					zap.String("deployment_chain", d.DeploymentChainId),
					zap.String("deployment", d.Id),
					zap.String("app", d.ApplicationId),
				)
			} else {
				c.logger.Info("unable to start planning deployment, try again next sync interval",
					zap.String("deployment", d.Id),
					zap.String("app", d.ApplicationId),
				)
			}
			continue
		}

		planner, err := c.startNewPlanner(ctx, d)
		if err != nil {
			c.logger.Error("failed to start a new planner",
				zap.String("deployment", d.Id),
				zap.String("app", d.ApplicationId),
				zap.Error(err),
			)
			continue
		}
		c.planners[appID] = planner

		// Application will be marked as DEPLOYING after its planner was successfully created.
		if err := reportApplicationDeployingStatus(ctx, c.apiClient, d.ApplicationId, true); err != nil {
			c.logger.Error("failed to mark application as deploying",
				zap.String("deployment", d.Id),
				zap.String("app", d.ApplicationId),
				zap.Error(err),
			)
		}
	}
}

func (c *controller) startNewPlanner(ctx context.Context, d *model.Deployment) (*planner, error) {
	logger := c.logger.With(
		zap.String("deployment", d.Id),
		zap.String("app", d.ApplicationId),
	)
	logger.Info("a new planner will be started")

	// Ensure the existence of the working directory for the deployment.
	workingDir, err := os.MkdirTemp(c.workspaceDir, d.Id+"-planner-*")
	if err != nil {
		logger.Error("failed to create working directory for planner", zap.Error(err))
		return nil, err
	}

	logger = logger.With(zap.String("working-dir", workingDir))

	// The most recent successful commit is saved in memory.
	// But when the piped is restarted that data will be cleared too.
	// So in that case, we have to use the API to check.
	var (
		commitHash     = c.mostRecentlySuccessfulCommits[d.ApplicationId]
		configFilename = c.mostRecentlySuccessfulConfigFilenames[d.ApplicationId]
	)
	if commitHash == "" {
		dref, err := c.getMostRecentlySuccessfulDeployment(ctx, d.ApplicationId)
		switch {
		case err == nil:
			commitHash = dref.Trigger.Commit.Hash
			configFilename = dref.ConfigFilename
			c.mostRecentlySuccessfulCommits[d.ApplicationId] = commitHash
			c.mostRecentlySuccessfulConfigFilenames[d.ApplicationId] = configFilename

		case status.Code(err) == codes.NotFound:
			logger.Info("there is no previous successful commit for this application")

		default:
			return nil, fmt.Errorf("failed to get the most recently successful deployment (%w)", err)
		}
	}

	planner := newPlanner(
		d,
		commitHash,
		configFilename,
		workingDir,
		c.pluginRegistry,
		c.apiClient,
		c.gitClient,
		c.notifier,
		c.secretDecrypter,
		c.logger,
		c.tracerProvider,
	)

	cleanup := func() {
		logger.Info("cleaning up working directory for planner")
		if err := os.RemoveAll(workingDir); err != nil {
			logger.Warn("failed to clean working directory", zap.Error(err))
		}
	}

	// Start running planner.
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		defer cleanup()
		if err := planner.Run(ctx); err != nil {
			logger.Error("failed to run planner", zap.Error(err))
		}
	}()

	return planner, nil
}

// syncSchedulers adds new scheduler for newly PLANNED/RUNNING deployments
// as well as removes the schedulers for the completed deployments.
func (c *controller) syncSchedulers(ctx context.Context) {
	// Update the most recent successful commit hashes.
	for id, s := range c.schedulers {
		if !s.IsDone() {
			continue
		}
		if s.DoneDeploymentStatus() != model.DeploymentStatus_DEPLOYMENT_SUCCESS {
			continue
		}
		c.mostRecentlySuccessfulCommits[id] = s.CommitHash()
		c.mostRecentlySuccessfulConfigFilenames[id] = s.ConfigFilename()
	}

	// Remove done schedulers.
	for id, t := range c.doneSchedulers {
		if time.Since(t) >= schedulerStaleDuration {
			delete(c.doneSchedulers, id)
		}
	}

	for id, s := range c.schedulers {
		if !s.IsDone() {
			continue
		}
		c.logger.Info("deleted done scheduler",
			zap.String("deployment", s.ID()),
			zap.String("app", id),
			zap.Int("count", len(c.schedulers)),
		)
		c.doneSchedulers[s.ID()] = s.DoneTimestamp()
		delete(c.schedulers, id)

		// Application will be marked as NOT deploying when scheduler's deployment was completed.
		if s.DoneDeploymentStatus().IsCompleted() {
			if err := reportApplicationDeployingStatus(ctx, c.apiClient, id, false); err != nil {
				c.logger.Error("failed to mark application as NOT deploying",
					zap.String("deployment", s.ID()),
					zap.String("app", id),
					zap.Error(err),
				)
			}
		}
	}

	// Add missing schedulers.
	planneds := c.deploymentLister.ListPlanneds()
	runnings := c.deploymentLister.ListRunnings()
	targets := append(runnings, planneds...)

	if len(targets) == 0 {
		return
	}

	c.logger.Info(fmt.Sprintf("there are %d planned/running deployments for scheduling", len(targets)),
		zap.Int("count", len(c.schedulers)),
	)

	for _, d := range targets {
		// Ignore already processed one.
		if _, ok := c.doneSchedulers[d.Id]; ok {
			continue
		}
		if s, ok := c.schedulers[d.ApplicationId]; ok {
			if s.ID() != d.Id {
				c.logger.Warn("detected an application that has more than one running deployments",
					zap.String("app", d.ApplicationId),
					zap.String("handling-deployment", s.ID()),
					zap.String("deployment", d.Id),
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
			zap.String("deployment", d.Id),
			zap.String("app", d.ApplicationId),
			zap.Int("count", len(c.schedulers)),
		)
	}
}

// startNewScheduler creates and starts running a new scheduler
// for a specific PLANNED deployment.
// This adds the newly created one to the scheduler list
// for tracking its lifetime periodically later.
func (c *controller) startNewScheduler(ctx context.Context, d *model.Deployment) (*scheduler, error) {
	logger := c.logger.With(
		zap.String("deployment", d.Id),
		zap.String("app", d.ApplicationId),
	)
	logger.Info("will add a new scheduler")

	// Ensure the existence of the working directory for the deployment.
	workingDir, err := os.MkdirTemp(c.workspaceDir, d.Id+"-scheduler-*")
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
		c.pluginRegistry,
		c.notifier,
		c.secretDecrypter,
		c.logger,
		c.tracerProvider,
	)

	c.metadataStoreRegistry.Register(d)

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
		defer c.metadataStoreRegistry.Delete(d.Id)
		if err := scheduler.Run(ctx); err != nil {
			logger.Error("failed to run scheduler", zap.Error(err))
		}
	}()

	return scheduler, nil
}

func (c *controller) getMostRecentlySuccessfulDeployment(ctx context.Context, applicationID string) (*model.ApplicationDeploymentReference, error) {
	req := &pipedservice.GetApplicationMostRecentDeploymentRequest{
		ApplicationId: applicationID,
		Status:        model.DeploymentStatus_DEPLOYMENT_SUCCESS,
	}

	d, err := pipedservice.NewRetry(3).Do(ctx, func() (interface{}, error) {
		resp, err := c.apiClient.GetApplicationMostRecentDeployment(ctx, req)
		if err == nil {
			return resp.Deployment, nil
		}
		return nil, pipedservice.NewRetriableErr(err)
	})
	if err != nil {
		return nil, err
	}
	return d.(*model.ApplicationDeploymentReference), nil
}

func (c *controller) shouldStartPlanningDeployment(ctx context.Context, d *model.Deployment) (plannable, cancel bool, cancelReason string, err error) {
	if !d.IsInChainDeployment() {
		plannable = true
		return
	}
	resp, err := c.apiClient.InChainDeploymentPlannable(ctx, &pipedservice.InChainDeploymentPlannableRequest{
		DeploymentId:              d.Id,
		DeploymentChainId:         d.DeploymentChainId,
		DeploymentChainBlockIndex: d.DeploymentChainBlockIndex,
	})
	if err != nil {
		return
	}
	plannable = resp.Plannable
	cancel = resp.Cancel
	cancelReason = resp.CancelReason
	return
}

func (c *controller) cancelDeployment(ctx context.Context, d *model.Deployment, reason string) error {
	req := &pipedservice.ReportDeploymentCompletedRequest{
		DeploymentId:              d.Id,
		Status:                    model.DeploymentStatus_DEPLOYMENT_CANCELLED,
		StatusReason:              reason,
		StageStatuses:             nil,
		DeploymentChainId:         d.DeploymentChainId,
		DeploymentChainBlockIndex: d.DeploymentChainBlockIndex,
		CompletedAt:               time.Now().Unix(),
	}

	_, err := pipedservice.NewRetry(10).Do(ctx, func() (interface{}, error) {
		_, err := c.apiClient.ReportDeploymentCompleted(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("failed to report deployment status to control-plane: %w", err)
		}
		return nil, nil
	})
	return err
}

func reportApplicationDeployingStatus(ctx context.Context, c apiClient, appID string, deploying bool) error {
	req := &pipedservice.ReportApplicationDeployingStatusRequest{
		ApplicationId: appID,
		Deploying:     deploying,
	}

	_, err := pipedservice.NewRetry(10).Do(ctx, func() (interface{}, error) {
		_, err := c.ReportApplicationDeployingStatus(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("failed to report application deploying status to control-plane: %w", err)
		}
		return nil, nil
	})
	return err
}
