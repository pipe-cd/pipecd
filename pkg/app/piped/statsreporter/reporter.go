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

// Package statsreporter provides a piped component
// that periodically reports local metrics to control-plane.
package statsreporter

import (
	"context"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/pipe-cd/pipecd/pkg/app/server/service/pipedservice"
)

type apiClient interface {
	ReportStat(ctx context.Context, req *pipedservice.ReportStatRequest, opts ...grpc.CallOption) (*pipedservice.ReportStatResponse, error)
}

type Reporter interface {
	Run(ctx context.Context) error
}

type reporter struct {
	metricsURL string
	httpClient *http.Client
	apiClient  apiClient
	interval   time.Duration
	logger     *zap.Logger
}

func NewReporter(metricsURL string, apiClient apiClient, logger *zap.Logger) *reporter {
	return &reporter{
		metricsURL: metricsURL,
		httpClient: &http.Client{Timeout: 10 * time.Second},
		apiClient:  apiClient,
		interval:   time.Minute,
		logger:     logger.Named("stats-reporter"),
	}
}

func (r *reporter) Run(ctx context.Context) error {
	r.logger.Info("start running stats reporter")

	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			r.logger.Info("stats reporter has been stopped")
			return nil

		case now := <-ticker.C:
			if err := r.report(ctx); err != nil {
				continue
			}
			r.logger.Info("successfully collected and reported metrics",
				zap.Duration("duration", time.Since(now)),
			)
		}
	}
}

func (r *reporter) report(ctx context.Context) error {
	resp, err := r.httpClient.Get(r.metricsURL)
	if err != nil {
		r.logger.Error("failed to fetch prometheus metrics", zap.Error(err))
		return err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		r.logger.Error("failed to load prometheus metrics", zap.Error(err))
		return err
	}

	req := &pipedservice.ReportStatRequest{
		PipedStats: b,
	}
	if _, err := r.apiClient.ReportStat(ctx, req); err != nil {
		r.logger.Error("failed to report stats", zap.Error(err))
		return err
	}
	return nil
}
