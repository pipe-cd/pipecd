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

package lambda

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/pipe-cd/pipecd/pkg/app/piped/livestatestore/lambda"
	"github.com/pipe-cd/pipecd/pkg/app/server/service/pipedservice"
	"github.com/pipe-cd/pipecd/pkg/config"
	"github.com/pipe-cd/pipecd/pkg/model"
)

type applicationLister interface {
	ListByPlatformProvider(name string) []*model.Application
}

type apiClient interface {
	ReportApplicationLiveState(ctx context.Context, req *pipedservice.ReportApplicationLiveStateRequest, opts ...grpc.CallOption) (*pipedservice.ReportApplicationLiveStateResponse, error)
	ReportApplicationLiveStateEvents(ctx context.Context, req *pipedservice.ReportApplicationLiveStateEventsRequest, opts ...grpc.CallOption) (*pipedservice.ReportApplicationLiveStateEventsResponse, error)
}

type Reporter interface {
	Run(ctx context.Context) error
	ProviderName() string
}

type reporter struct {
	provider              config.PipedPlatformProvider
	appLister             applicationLister
	stateGetter           lambda.Getter
	apiClient             apiClient
	snapshotFlushInterval time.Duration
	logger                *zap.Logger

	snapshotVersions map[string]model.ApplicationLiveStateVersion
}

func NewReporter(cp config.PipedPlatformProvider, appLister applicationLister, stateGetter lambda.Getter, apiClient apiClient, logger *zap.Logger) Reporter {
	logger = logger.Named("lambda-reporter").With(
		zap.String("platform-provider", cp.Name),
	)
	return &reporter{
		provider:              cp,
		appLister:             appLister,
		stateGetter:           stateGetter,
		apiClient:             apiClient,
		snapshotFlushInterval: time.Minute,
		logger:                logger,
		snapshotVersions:      make(map[string]model.ApplicationLiveStateVersion),
	}
}

func (r *reporter) Run(ctx context.Context) error {
	r.logger.Info("start running app live state reporter")

	r.logger.Info("waiting for livestatestore to be ready")
	if err := r.stateGetter.WaitForReady(ctx, 10*time.Minute); err != nil {
		r.logger.Error("livestatestore was unable to be ready in time", zap.Error(err))
		return err
	}

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

func (r *reporter) ProviderName() string {
	return r.provider.Name
}

func (r *reporter) flushSnapshots(ctx context.Context) {
	apps := r.appLister.ListByPlatformProvider(r.provider.Name)
	for _, app := range apps {
		state, ok := r.stateGetter.GetState(app.Id)
		if !ok {
			r.logger.Info(fmt.Sprintf("no app state of lambda application %s to report", app.Id))
			continue
		}

		snapshot := &model.ApplicationLiveStateSnapshot{
			ApplicationId: app.Id,
			PipedId:       app.PipedId,
			ProjectId:     app.ProjectId,
			Kind:          app.Kind,
			Lambda: &model.LambdaApplicationLiveState{
				Resources: state.Resources,
			},
			Version: &state.Version,
		}
		snapshot.DetermineAppHealthStatus()
		req := &pipedservice.ReportApplicationLiveStateRequest{
			Snapshot: snapshot,
		}

		if _, err := r.apiClient.ReportApplicationLiveState(ctx, req); err != nil {
			r.logger.Error("failed to report application live state",
				zap.String("application-id", app.Id),
				zap.Error(err),
			)
			continue
		}
		r.snapshotVersions[app.Id] = state.Version
		r.logger.Info(fmt.Sprintf("successfully reported application live state for application: %s", app.Id))
	}
}
