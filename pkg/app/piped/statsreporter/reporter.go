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
	"io"
	"net/http"
	"time"

	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/pipe-cd/pipe/pkg/app/api/service/pipedservice"
	"github.com/pipe-cd/pipe/pkg/model"
	"github.com/pipe-cd/pipe/pkg/version"
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

L:
	for {
		select {
		case <-ctx.Done():
			break L

		case now := <-ticker.C:
			metrics, err := r.collect()
			if err != nil {
				continue
			}
			if len(metrics) == 0 {
				r.logger.Info("there are no metrics to report")
				continue
			}
			if r.report(ctx, metrics, now); err != nil {
				continue
			}
			r.logger.Info("successfully collected and reported metrics",
				zap.Int("num", len(metrics)),
				zap.Duration("duration", time.Since(now)),
			)
		}
	}

	r.logger.Info("stats reporter has been stopped")
	return nil
}

func (r *reporter) collect() ([]*model.PrometheusMetrics, error) {
	resp, err := r.httpClient.Get(r.metricsURL)
	if err != nil {
		r.logger.Error("failed to collect prometheus metrics", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	metrics, err := parsePrometheusMetrics(resp.Body)
	if err != nil {
		r.logger.Error("failed to parse prometheus metrics", zap.Error(err))
		return nil, err
	}

	return metrics, nil
}

func (r *reporter) report(ctx context.Context, metrics []*model.PrometheusMetrics, now time.Time) error {
	req := &pipedservice.PingRequest{
		PipedStats: &model.PipedStats{
			Version:           version.Get().Version,
			Timestamp:         now.Unix(),
			PrometheusMetrics: metrics,
		},
	}
	if _, err := r.apiClient.Ping(ctx, req); err != nil {
		r.logger.Error("failed to report stats", zap.Error(err))
		return err
	}
	return nil
}

var parser expfmt.TextParser

// TODO: Add a metrics whitelist and fiter out not needed ones.
func parsePrometheusMetrics(reader io.Reader) ([]*model.PrometheusMetrics, error) {
	metricFamily, err := parser.TextToMetricFamilies(reader)
	if err != nil {
		return nil, err
	}

	metrics := make([]*model.PrometheusMetrics, 0, len(metricFamily))

L:
	for _, mf := range metricFamily {
		var metricType model.PrometheusMetrics_Type

		switch mf.GetType() {
		case dto.MetricType_COUNTER:
			metricType = model.PrometheusMetrics_COUNTER
		case dto.MetricType_GAUGE:
			metricType = model.PrometheusMetrics_GAUGE
		default:
			continue L
		}

		metric := &model.PrometheusMetrics{
			Name: *mf.Name,
			Type: metricType,
		}

		for _, m := range mf.Metric {
			sample := &model.PrometheusMetrics_Sample{
				Labels: make([]*model.PrometheusMetrics_LabelPair, 0, len(m.Label)),
			}
			metric.Samples = append(metric.Samples, sample)

			for _, l := range m.Label {
				sample.Labels = append(sample.Labels, &model.PrometheusMetrics_LabelPair{
					Name:  l.GetName(),
					Value: l.GetValue(),
				})
			}

			switch metric.Type {
			case model.PrometheusMetrics_COUNTER:
				sample.Value = m.Counter.GetValue()
			case model.PrometheusMetrics_GAUGE:
				sample.Value = m.Gauge.GetValue()
			}
		}

		if len(metric.Samples) > 0 {
			metrics = append(metrics, metric)
		}
	}

	return metrics, nil
}
