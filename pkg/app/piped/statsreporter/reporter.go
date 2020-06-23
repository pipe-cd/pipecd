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
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

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
			stats, err := r.collect()
			if err != nil {
				continue
			}
			if len(stats) == 0 {
				r.logger.Info("there are no stats to report")
				continue
			}
			if r.report(ctx, stats, now); err != nil {
				continue
			}
			r.logger.Info("successfully collected and reported stats",
				zap.Int("num", len(stats)),
				zap.Duration("duration", time.Since(now)),
			)
		}
	}

	r.logger.Info("stats reporter has been stopped")
	return nil
}

func (r *reporter) collect() ([]*model.PipedStats_PrometheusMetrics, error) {
	resp, err := r.httpClient.Get(r.metricsURL)
	if err != nil {
		r.logger.Error("failed to collect prometheus metrics", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	stats, err := parsePrometheusMetrics(resp.Body)
	if err != nil {
		r.logger.Error("failed to parse prometheus metrics", zap.Error(err))
		return nil, err
	}

	return stats, nil
}

func (r *reporter) report(ctx context.Context, stats []*model.PipedStats_PrometheusMetrics, now time.Time) error {
	req := &pipedservice.PingRequest{
		PipedStats: &model.PipedStats{
			Version:         version.Get().Version,
			Timestamp:       now.Unix(),
			PrometheusStats: stats,
		},
	}
	if _, err := r.apiClient.Ping(ctx, req); err != nil {
		r.logger.Error("failed to report stats", zap.Error(err))
		return err
	}
	return nil
}

var helpPrefix = []byte("# HELP")

const typePrefix = "# TYPE"

// TODO: Add a metrics whitelist and fiter out not needed ones.
func parsePrometheusMetrics(reader io.Reader) ([]*model.PipedStats_PrometheusMetrics, error) {
	var (
		curType = model.PipedStats_PrometheusMetrics_UNKNOWN
		metrics []*model.PipedStats_PrometheusMetrics
	)
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		lb := scanner.Bytes()
		if len(lb) == 0 {
			continue
		}

		// Ignore all HELP line.
		if bytes.HasPrefix(lb, helpPrefix) {
			continue
		}

		// Extract current type from TYPE line.
		line := string(lb)
		if strings.HasPrefix(line, typePrefix) {
			parts := strings.Split(line, " ")
			if len(parts) < 3 {
				return nil, fmt.Errorf("malformed TYPE line %s", line)
			}
			switch parts[len(parts)-1] {
			case "gauge":
				curType = model.PipedStats_PrometheusMetrics_GAUGE
			case "counter":
				curType = model.PipedStats_PrometheusMetrics_COUNTER
			default:
				curType = model.PipedStats_PrometheusMetrics_UNKNOWN
			}
			continue
		}

		if curType == model.PipedStats_PrometheusMetrics_UNKNOWN {
			continue
		}

		// Extract metrics data.
		lastSpaceIndex := strings.LastIndexByte(line, ' ')
		if lastSpaceIndex < 0 {
			continue
		}
		value, err := strconv.ParseFloat(line[lastSpaceIndex+1:], 64)
		if err != nil {
			return nil, err
		}
		metrics = append(metrics, &model.PipedStats_PrometheusMetrics{
			Type:  curType,
			Name:  line[:lastSpaceIndex],
			Value: value,
		})
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return metrics, nil
}
