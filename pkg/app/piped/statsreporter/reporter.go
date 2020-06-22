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

// Package statsreporter provides a piped component
// that periodically reports local metrics to control-plane.
package statsreporter

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/pipe-cd/pipe/pkg/app/api/service/pipedservice"
	"github.com/pipe-cd/pipe/pkg/model"
)

type apiClient interface {
	Ping(ctx context.Context, req *pipedservice.PingRequest, opts ...grpc.CallOption) (*pipedservice.PingResponse, error)
}

type Reporter interface {
	Run(ctx context.Context) error
}

type reporter struct {
	metricsURL string
	httpClient *http.Client
	interval   time.Duration
	logger     *zap.Logger
}

func NewReporter(metricsURL string, apiClient apiClient, logger *zap.Logger) *reporter {
	return &reporter{
		metricsURL: metricsURL,
		httpClient: &http.Client{Timeout: 10 * time.Second},
		interval:   time.Minute,
		logger:     logger.Named("stats-reporter"),
	}
}

func (r *reporter) Run(ctx context.Context) error {
	r.logger.Info("start running stats reporter")

	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()

	{
		stats, err := r.collect(ctx)
		if err != nil {
			r.logger.Error("failed while collecting stats", zap.Error(err))
		}
		if r.report(ctx, stats); err != nil {
			r.logger.Error("failed while reporting stats", zap.Error(err))
		}
		r.logger.Info("successfully collected and reported stats", zap.Int("num", len(stats)))
	}

L:
	for {
		select {
		case <-ctx.Done():
			break L

		case <-ticker.C:
			stats, err := r.collect(ctx)
			if err != nil {
				r.logger.Error("failed while collecting stats", zap.Error(err))
				continue
			}
			if r.report(ctx, stats); err != nil {
				r.logger.Error("failed while reporting stats", zap.Error(err))
				continue
			}
			r.logger.Info("successfully collected and reported stats", zap.Int("num", len(stats)))
		}
	}

	r.logger.Info("stats reporter has been stopped")
	return nil
}

func (r *reporter) collect(ctx context.Context) ([]*model.PipedStats_PrometheusMetrics, error) {
	resp, err := r.httpClient.Get(r.metricsURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	metrics, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	r.logger.Info(fmt.Sprintf("collected %d bytes of metrics", len(metrics)))

	return nil, nil
}

func (r *reporter) report(ctx context.Context, stats []*model.PipedStats_PrometheusMetrics) error {
	return nil
}
