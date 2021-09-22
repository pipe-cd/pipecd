// Copyright 2021 The PipeCD Authors.
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

package analysis

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"text/template"
	"time"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipe/pkg/app/piped/analysisprovider/metrics"
	"github.com/pipe-cd/pipe/pkg/app/piped/apistore/analysisresultstore"
	"github.com/pipe-cd/pipe/pkg/app/piped/executor"
	"github.com/pipe-cd/pipe/pkg/config"
)

const (
	canaryVariantName   = "canary"
	baselineVariantName = "baseline"
	primaryVariantName  = "primary"
)

type metricsAnalyzer struct {
	id                  string
	cfg                 config.AnalysisMetrics
	stageStartTime      time.Time
	provider            metrics.Provider
	analysisResultStore executor.AnalysisResultStore
	logger              *zap.Logger
	logPersister        executor.LogPersister
}

func newMetricsAnalyzer(id string, cfg config.AnalysisMetrics, stageStartTime time.Time, provider metrics.Provider, analysisResultStore executor.AnalysisResultStore, logger *zap.Logger, logPersister executor.LogPersister) *metricsAnalyzer {
	return &metricsAnalyzer{
		id:                  id,
		cfg:                 cfg,
		stageStartTime:      stageStartTime,
		provider:            provider,
		analysisResultStore: analysisResultStore,
		logPersister:        logPersister,
		logger: logger.With(
			zap.String("analyzer-id", id),
		),
	}
}

// run starts an analysis which runs the query at the given interval, until the context is done.
// It returns an error when the number of failures exceeds the the failureLimit.
func (a *metricsAnalyzer) run(ctx context.Context) error {
	ticker := time.NewTicker(a.cfg.Interval.Duration())
	defer ticker.Stop()

	failureCount := 0
	for {
		select {
		case <-ticker.C:
			var (
				expected bool
				err      error
			)
			switch a.cfg.Strategy {
			case config.AnalysisStrategyThreshold:
				expected, err = a.analyzeWithThreshold(ctx)
			case config.AnalysisStrategyPrevious:
				var firstDeploy bool
				expected, firstDeploy, err = a.analyzeWithPrevious(ctx)
				if firstDeploy {
					a.logPersister.Info("[%s] PreviousAnalysis cannot be executed because this seems to be the first deployment, so it is considered as a success")
					return nil
				}
			case config.AnalysisStrategyCanaryBaseline:
				expected, err = a.analyzeWithCanaryBaseline(ctx)
			case config.AnalysisStrategyCanaryPrimary:
				expected, err = a.analyzeWithCanaryPrimary(ctx)
			default:
				return fmt.Errorf("unknown strategy %q given", a.cfg.Strategy)
			}
			// Ignore parent's context deadline exceeded error, and return immediately.
			if errors.Is(err, context.DeadlineExceeded) && ctx.Err() == context.DeadlineExceeded {
				return nil
			}
			if errors.Is(err, metrics.ErrNoDataFound) && a.cfg.SkipOnNoData {
				a.logPersister.Infof("[%s] The query result evaluation was skipped because \"skipOnNoData\" is true even though no data returned. Reason: %v. Performed query: %q", a.id, err, a.cfg.Query)
				continue
			}
			if err != nil {
				a.logPersister.Errorf("[%s] Unexpected error: %v. Performed query: %q", a.id, err, a.cfg.Query)
			}
			if expected {
				a.logPersister.Successf("[%s] The query result is expected one. Performed query: %q", a.id, a.cfg.Query)
				continue
			}
			failureCount++
			if failureCount > a.cfg.FailureLimit {
				return fmt.Errorf("analysis '%s' failed because the failure number exceeded the failure limit (%d)", a.id, a.cfg.FailureLimit)
			}
		case <-ctx.Done():
			return nil
		}
	}
}

// analyzeWithThreshold returns false if any data point is out of the prediction range.
// Return an error if the evaluation could not be executed normally.
func (a *metricsAnalyzer) analyzeWithThreshold(ctx context.Context) (bool, error) {
	if err := a.cfg.Expected.Validate(); err != nil {
		return false, fmt.Errorf("\"expected\" is required to analyze with the THRESHOLD strategy")
	}

	now := time.Now()
	queryRange := metrics.QueryRange{
		From: now.Add(-a.cfg.Interval.Duration()),
		To:   now,
	}
	points, err := a.provider.QueryPoints(ctx, a.cfg.Query, queryRange)
	if err != nil {
		return false, fmt.Errorf("failed to run query: %w", err)
	}

	var outiler metrics.DataPoint
	expected := true
	for i := range points {
		if a.cfg.Expected.InRange(points[i].Value) {
			continue
		}
		expected = false
		outiler = points[i]
		break
	}
	if !expected {
		a.logPersister.Errorf("[%s] Failed because it found a data point (%s) that is outside the expected range (%s). Performed query: %q", a.id, outiler, a.cfg.Expected, a.cfg.Query)
		return false, nil
	}

	return true, nil
}

// analyzeWithPrevious returns false if primary deviates in the specified direction compared to the previous deployment.
// Return an error if the evaluation could not be executed normally.
// elapsedTime is used to compare metrics at the same point in time after the analysis has started.
func (a *metricsAnalyzer) analyzeWithPrevious(ctx context.Context) (expected, firstDeploy bool, err error) {
	now := time.Now()
	queryRange := metrics.QueryRange{
		From: now.Add(-a.cfg.Interval.Duration()),
		To:   now,
	}
	points, err := a.provider.QueryPoints(ctx, a.cfg.Query, queryRange)
	if err != nil {
		return false, false, fmt.Errorf("failed to run query: %w", err)
	}

	prevMetadata, err := a.analysisResultStore.GetLatestAnalysisResult(ctx)
	if errors.Is(err, analysisresultstore.ErrNotFound) {
		return false, true, nil
	}
	if err != nil {
		return false, false, fmt.Errorf("failed to fetch the most recent successful analysis metadata: %w", err)
	}
	// Compare it with the previous metrics when the same amount of time as now has passed since the start of the stage.
	elapsedTime := now.Sub(a.stageStartTime)
	prevTo := time.Unix(prevMetadata.StartTime, 0).Add(elapsedTime)
	prevFrom := prevTo.Add(-a.cfg.Interval.Duration())
	prevQueryRange := metrics.QueryRange{
		From: prevFrom,
		To:   prevTo,
	}
	prevPoints, err := a.provider.QueryPoints(ctx, a.cfg.Query, prevQueryRange)
	if err != nil {
		return false, false, fmt.Errorf("failed to run query to fetch metrics for the previous deployment: %w", err)
	}
	if err := mannWhitneyUTest(points, prevPoints, a.cfg.Deviation); err != nil {
		a.logPersister.Errorf("[%s] Failed because %v. Performed query: %q", a.id, err, a.cfg.Query)
		return false, false, err
	}
	return true, false, nil
}

// analyzeWithCanaryBaseline returns false if canary deviates in the specified direction compared to baseline.
// Return an error if the evaluation could not be executed normally.
func (a *metricsAnalyzer) analyzeWithCanaryBaseline(ctx context.Context) (bool, error) {
	now := time.Now()
	queryRange := metrics.QueryRange{
		From: now.Add(-a.cfg.Interval.Duration()),
		To:   now,
	}
	canaryQuery, err := a.render(a.cfg.Query, a.cfg.CanaryArgs, canaryVariantName)
	if err != nil {
		return false, fmt.Errorf("failed to render query template for Canary: %w", err)
	}
	baselineQuery, err := a.render(a.cfg.Query, a.cfg.BaselineArgs, baselineVariantName)
	if err != nil {
		return false, fmt.Errorf("failed to render query template for Baseline: %w", err)
	}

	canaryPoints, err := a.provider.QueryPoints(ctx, canaryQuery, queryRange)
	if err != nil {
		return false, fmt.Errorf("failed to run query to fetch metrics for the Canary variant: %w", err)
	}
	baselinePoints, err := a.provider.QueryPoints(ctx, baselineQuery, queryRange)
	if err != nil {
		return false, fmt.Errorf("failed to run query to fetch metrics for the Baseline variant: %w", err)
	}

	if err := mannWhitneyUTest(canaryPoints, baselinePoints, a.cfg.Deviation); err != nil {
		a.logPersister.Errorf("[%s] Failed because %v. Performed query for canary: %q. Performed query for baseline: %q", a.id, err, canaryQuery, baselineQuery)
		return false, nil
	}
	return true, nil
}

// analyzeWithCanaryPrimary returns false if canary deviates in the specified direction compared to primary.
// Return an error if the evaluation could not be executed normally.
func (a *metricsAnalyzer) analyzeWithCanaryPrimary(ctx context.Context) (bool, error) {
	now := time.Now()
	queryRange := metrics.QueryRange{
		From: now.Add(-a.cfg.Interval.Duration()),
		To:   now,
	}
	canaryQuery, err := a.render(a.cfg.Query, a.cfg.CanaryArgs, canaryVariantName)
	if err != nil {
		return false, fmt.Errorf("failed to render query template for Canary: %w", err)
	}
	primaryQuery, err := a.render(a.cfg.Query, a.cfg.PrimaryArgs, primaryVariantName)
	if err != nil {
		return false, fmt.Errorf("failed to render query template for Primary: %w", err)
	}

	canaryPoints, err := a.provider.QueryPoints(ctx, canaryQuery, queryRange)
	if err != nil {
		return false, fmt.Errorf("failed to run query to fetch metrics for the Canary variant: %w", err)
	}
	primaryPoints, err := a.provider.QueryPoints(ctx, primaryQuery, queryRange)
	if err != nil {
		return false, fmt.Errorf("failed to run query to fetch metrics for the Primary variant: %w", err)
	}
	if err := mannWhitneyUTest(canaryPoints, primaryPoints, a.cfg.Deviation); err != nil {
		a.logPersister.Errorf("[%s] Failed because %v. Performed query for canary: %q. Performed query for primary: %q", a.id, err, canaryQuery, primaryQuery)
		return false, nil
	}
	return true, nil
}

type argsForTemplate struct {
	BuiltInArgs builtInArgs
	// User-defined custom args.
	VariantArgs map[string]string
}

type builtInArgs struct {
	Variant struct {
		Name string
	}
}

// render applies the given variant args to the query template.
func (a *metricsAnalyzer) render(queryTemplate string, variantArgs map[string]string, variant string) (string, error) {
	args := argsForTemplate{
		BuiltInArgs: builtInArgs{Variant: struct{ Name string }{Name: variant}},
		VariantArgs: variantArgs,
	}

	t, err := template.New("AnalysisVariantTemplate").Parse(queryTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse query template: %w", err)
	}

	b := new(bytes.Buffer)
	if err := t.Execute(b, args); err != nil {
		return "", fmt.Errorf("failed to apply template: %w", err)
	}
	return b.String(), err
}

// No error means that the result is expected.
func mannWhitneyUTest(target, base []metrics.DataPoint, deviation string) (err error) {
	if len(target) == 0 {
		return fmt.Errorf("no data points for target found")
	}
	if len(base) == 0 {
		return fmt.Errorf("no data points for base found")
	}
	// TODO: Implement mannWhitneyUTest
	return fmt.Errorf("mannWhitneyUTest isn't implemented yet")
}
