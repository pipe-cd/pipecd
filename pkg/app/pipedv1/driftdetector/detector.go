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

// Package driftdetector provides a piped component
// that continuously checks configuration drift between the current live state
// and the state defined at the latest commit of all applications.
package driftdetector

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/driftdetector/cloudrun"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/driftdetector/kubernetes"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/driftdetector/terraform"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/livestatestore"
	"github.com/pipe-cd/pipecd/pkg/app/server/service/pipedservice"
	"github.com/pipe-cd/pipecd/pkg/cache"
	"github.com/pipe-cd/pipecd/pkg/config"
	"github.com/pipe-cd/pipecd/pkg/git"
	"github.com/pipe-cd/pipecd/pkg/model"
)

type applicationLister interface {
	ListByPlatformProvider(name string) []*model.Application
}

type deploymentLister interface {
	ListAppHeadDeployments() map[string]*model.Deployment
}

type gitClient interface {
	Clone(ctx context.Context, repoID, remote, branch, destination string) (git.Repo, error)
}

type apiClient interface {
	ReportApplicationSyncState(ctx context.Context, req *pipedservice.ReportApplicationSyncStateRequest, opts ...grpc.CallOption) (*pipedservice.ReportApplicationSyncStateResponse, error)
}

type secretDecrypter interface {
	Decrypt(string) (string, error)
}

type Detector interface {
	Run(ctx context.Context) error
}

type detector struct {
	apiClient  apiClient
	detectors  []providerDetector
	syncStates map[string]model.ApplicationSyncState
	mu         sync.RWMutex
	logger     *zap.Logger
}

type providerDetector interface {
	Run(ctx context.Context) error
	ProviderName() string
}

func NewDetector(
	appLister applicationLister,
	gitClient gitClient,
	stateGetter livestatestore.Getter,
	apiClient apiClient,
	appManifestsCache cache.Cache,
	cfg *config.PipedSpec,
	sd secretDecrypter,
	logger *zap.Logger,
) (Detector, error) {

	d := &detector{
		apiClient:  apiClient,
		detectors:  make([]providerDetector, 0, len(cfg.PlatformProviders)),
		syncStates: make(map[string]model.ApplicationSyncState),
		logger:     logger.Named("drift-detector"),
	}

	const format = "unable to find live state getter for platform provider: %s"

	for _, cp := range cfg.PlatformProviders {
		switch cp.Type {
		case model.PlatformProviderKubernetes:
			sg, ok := stateGetter.KubernetesGetter(cp.Name)
			if !ok {
				return nil, fmt.Errorf(format, cp.Name)
			}
			d.detectors = append(d.detectors, kubernetes.NewDetector(
				cp,
				appLister,
				gitClient,
				sg,
				d,
				appManifestsCache,
				cfg,
				sd,
				logger,
			))

		case model.PlatformProviderCloudRun:
			sg, ok := stateGetter.CloudRunGetter(cp.Name)
			if !ok {
				return nil, fmt.Errorf(format, cp.Name)
			}
			d.detectors = append(d.detectors, cloudrun.NewDetector(
				cp,
				appLister,
				gitClient,
				sg,
				d,
				appManifestsCache,
				cfg,
				sd,
				logger,
			))

		case model.PlatformProviderTerraform:
			if !*cp.TerraformConfig.DriftDetectionEnabled {
				continue
			}
			sg, ok := stateGetter.TerraformGetter(cp.Name)
			if !ok {
				return nil, fmt.Errorf(format, cp.Name)
			}
			d.detectors = append(d.detectors, terraform.NewDetector(
				cp,
				appLister,
				gitClient,
				sg,
				d,
				appManifestsCache,
				cfg,
				sd,
				logger,
			))

		default:
		}
	}

	return d, nil
}

func (d *detector) Run(ctx context.Context) error {
	group, ctx := errgroup.WithContext(ctx)

	for i, detector := range d.detectors {
		detector := detector
		// Avoid starting all detectors at the same time to reduce the API call burst.
		time.Sleep(time.Duration(i) * 10 * time.Second)
		d.logger.Info(fmt.Sprintf("starting drift detector for cloud provider: %s", detector.ProviderName()))

		group.Go(func() error {
			return detector.Run(ctx)
		})
	}

	d.logger.Info(fmt.Sprintf("all drift detectors of %d providers have been started", len(d.detectors)))

	if err := group.Wait(); err != nil {
		d.logger.Error("failed while running", zap.Error(err))
		return err
	}

	d.logger.Info(fmt.Sprintf("all drift detectors of %d providers have been stopped", len(d.detectors)))
	return nil
}

func (d *detector) ReportApplicationSyncState(ctx context.Context, appID string, state model.ApplicationSyncState) error {
	d.mu.RLock()
	curState, ok := d.syncStates[appID]
	d.mu.RUnlock()

	if ok && !curState.HasChanged(state) {
		return nil
	}

	_, err := d.apiClient.ReportApplicationSyncState(ctx, &pipedservice.ReportApplicationSyncStateRequest{
		ApplicationId: appID,
		State:         &state,
	})
	if err != nil {
		d.logger.Error("failed to report application sync state",
			zap.String("application-id", appID),
			zap.Any("state", state),
			zap.Error(err),
		)
		return err
	}

	d.mu.Lock()
	d.syncStates[appID] = state
	d.mu.Unlock()

	return nil
}
