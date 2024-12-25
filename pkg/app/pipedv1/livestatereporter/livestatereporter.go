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
	"time"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	"github.com/pipe-cd/pipecd/pkg/app/server/service/pipedservice"
	"github.com/pipe-cd/pipecd/pkg/model"
	pluginapi "github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1"
	"github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1/livestate"
)

type applicationLister interface {
	ListByPluginName(name string) []*model.Application
}

type apiClient interface {
	ReportApplicationLiveState(ctx context.Context, req *pipedservice.ReportApplicationLiveStateRequest, opts ...grpc.CallOption) (*pipedservice.ReportApplicationLiveStateResponse, error)
	ReportApplicationSyncState(ctx context.Context, req *pipedservice.ReportApplicationSyncStateRequest, opts ...grpc.CallOption) (*pipedservice.ReportApplicationSyncStateResponse, error)
}

// Reporter represents a component that reports the snapshot about live state of registered applications.
type Reporter interface {
	Run(ctx context.Context) error
}

type reporter struct {
	reporters []pluginReporter
	logger    *zap.Logger
}

// NewReporter creates a new reporter.
func NewReporter(appLister applicationLister, apiClient apiClient, plugins map[string]pluginapi.PluginClient, logger *zap.Logger) Reporter {
	rlogger := logger.Named("live-state-reporter")
	r := &reporter{
		reporters: make([]pluginReporter, 0, len(plugins)),
		logger:    rlogger,
	}

	for name, p := range plugins {
		r.reporters = append(r.reporters, pluginReporter{
			pluginName:            name,
			snapshotFlushInterval: time.Minute,
			appLister:             appLister,
			apiClient:             apiClient,
			pluginClient:          p,
			logger:                rlogger.With(zap.String("plugin-name", name)),
		})
	}
	return r
}

func (r *reporter) Run(ctx context.Context) error {
	group, ctx := errgroup.WithContext(ctx)

	r.logger.Info(fmt.Sprintf("starting %d live state reporters of plugins", len(r.reporters)))

	for _, reporter := range r.reporters {
		// Avoid starting all reporters at the same time to reduce the API call burst.
		time.Sleep(10 * time.Second)

		group.Go(func() error {
			reporter.logger.Info("starting app live state reporter for plugin")
			return reporter.Run(ctx)
		})
	}

	r.logger.Info(fmt.Sprintf("all live state reporters of %d plugins have been started", len(r.reporters)))

	if err := group.Wait(); err != nil {
		r.logger.Error("failed while running", zap.Error(err))
		return err
	}

	r.logger.Info(fmt.Sprintf("all live state reporters of %d plugins have been stopped", len(r.reporters)))
	return nil
}

type pluginReporter struct {
	pluginName            string
	snapshotFlushInterval time.Duration
	appLister             applicationLister
	apiClient             apiClient
	pluginClient          pluginapi.PluginClient
	logger                *zap.Logger
}

func (pr *pluginReporter) Run(ctx context.Context) error {
	pr.logger.Info("start running app live state reporter", zap.Duration("snapshot-flush-interval", pr.snapshotFlushInterval))

	snapshotTicker := time.NewTicker(pr.snapshotFlushInterval)
	defer snapshotTicker.Stop()

	for {
		select {
		case <-snapshotTicker.C:
			pr.flushSnapshots(ctx)

		case <-ctx.Done():
			pr.logger.Info("app live state reporter has been stopped")
			return nil
		}
	}
}

func (pr *pluginReporter) flushSnapshots(ctx context.Context) {
	// TODO: Implement appLister.ListByPluginName.
	apps := pr.appLister.ListByPluginName(pr.pluginName)
	for _, app := range apps {
		res, err := pr.pluginClient.GetLivestate(ctx, &livestate.GetLivestateRequest{ApplicationId: app.Id})
		if err != nil {
			pr.logger.Info(fmt.Sprintf("no app state of application %s to report", app.Id))
			continue
		}

		// Report the application live state to the control plane.
		snapshot := &model.ApplicationLiveStateSnapshot{
			ApplicationId:        app.Id,
			PipedId:              app.PipedId,
			ProjectId:            app.ProjectId,
			Kind:                 app.Kind,
			ApplicationLiveState: res.GetApplicationLiveState(),
		}
		// TODO: Implement DetermineAppHealthStatus for the case of plugin architecture.
		snapshot.DetermineAppHealthStatus()

		// TODO: Fix ReportApplicationLiveState to store ApplicationLiveState.
		if _, err := pr.apiClient.ReportApplicationLiveState(ctx, &pipedservice.ReportApplicationLiveStateRequest{
			Snapshot: snapshot,
		}); err != nil {
			pr.logger.Error("failed to report application live state",
				zap.String("application-id", app.Id),
				zap.Error(err),
			)
			continue
		}

		// Report the application sync state to the control plane.
		if _, err := pr.apiClient.ReportApplicationSyncState(ctx, &pipedservice.ReportApplicationSyncStateRequest{
			ApplicationId: app.Id,
			State:         res.GetSyncState(),
		}); err != nil {
			pr.logger.Error("failed to report application live state",
				zap.String("application-id", app.Id),
				zap.Error(err),
			)
			continue
		}

		pr.logger.Info(fmt.Sprintf("successfully reported application live state for application: %s", app.Id))
	}
}
