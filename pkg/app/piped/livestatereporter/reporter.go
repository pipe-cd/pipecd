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

	"github.com/pipe-cd/pipecd/pkg/app/piped/livestatereporter/cloudrun"
	"github.com/pipe-cd/pipecd/pkg/app/piped/livestatereporter/ecs"
	"github.com/pipe-cd/pipecd/pkg/app/piped/livestatereporter/kubernetes"
	"github.com/pipe-cd/pipecd/pkg/app/piped/livestatestore"
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
}

type reporter struct {
	reporters []providerReporter
	logger    *zap.Logger
}

type providerReporter interface {
	Run(ctx context.Context) error
	ProviderName() string
}

func NewReporter(appLister applicationLister, stateGetter livestatestore.Getter, apiClient apiClient, cfg *config.PipedSpec, logger *zap.Logger) Reporter {
	r := &reporter{
		reporters: make([]providerReporter, 0, len(cfg.PlatformProviders)),
		logger:    logger.Named("live-state-reporter"),
	}

	const errFmt = "unable to find live state getter for platform provider: %s"
	for _, cp := range cfg.PlatformProviders {
		switch cp.Type {
		case model.PlatformProviderKubernetes:
			sg, ok := stateGetter.KubernetesGetter(cp.Name)
			if !ok {
				r.logger.Error(fmt.Sprintf(errFmt, cp.Name))
				continue
			}
			r.reporters = append(r.reporters, kubernetes.NewReporter(cp, appLister, sg, apiClient, logger))
		case model.PlatformProviderCloudRun:
			sg, ok := stateGetter.CloudRunGetter(cp.Name)
			if !ok {
				r.logger.Error(fmt.Sprintf(errFmt, cp.Name))
				continue
			}
			r.reporters = append(r.reporters, cloudrun.NewReporter(cp, appLister, sg, apiClient, logger))
		case model.PlatformProviderECS:
			sg, ok := stateGetter.ECSGetter(cp.Name)
			if !ok {
				r.logger.Error(fmt.Sprintf(errFmt, cp.Name))
				continue
			}
			r.reporters = append(r.reporters, ecs.NewReporter(cp, appLister, sg, apiClient, logger))
		}
	}

	return r
}

func (r *reporter) Run(ctx context.Context) error {
	group, ctx := errgroup.WithContext(ctx)

	for i, reporter := range r.reporters {
		reporter := reporter
		// Avoid starting all reporters at the same time to reduce the API call burst.
		time.Sleep(time.Duration(i) * 10 * time.Second)
		r.logger.Info(fmt.Sprintf("starting app live state reporter for cloud provider: %s", reporter.ProviderName()))

		group.Go(func() error {
			return reporter.Run(ctx)
		})
	}

	r.logger.Info(fmt.Sprintf("all live state reporters of %d providers have been started", len(r.reporters)))

	if err := group.Wait(); err != nil {
		r.logger.Error("failed while running", zap.Error(err))
		return err
	}

	r.logger.Info(fmt.Sprintf("all live state reporters of %d providers have been stopped", len(r.reporters)))
	return nil
}
