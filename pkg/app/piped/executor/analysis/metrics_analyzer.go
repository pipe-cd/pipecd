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
	"github.com/pipe-cd/pipe/pkg/app/piped/executor/analysis/mannwhitney"
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
	// Application-specific arguments using when rendering the query.
	argsTemplate argsTemplate
	logger       *zap.Logger
	logPersister executor.LogPersister
}

func newMetricsAnalyzer(id string, cfg config.AnalysisMetrics, stageStartTime time.Time, provider metrics.Provider, analysisResultStore executor.AnalysisResultStore, argsTemplate argsTemplate, logger *zap.Logger, logPersister executor.LogPersister) *metricsAnalyzer {
	return &metricsAnalyzer{
		id:                  id,
		cfg:                 cfg,
		stageStartTime:      stageStartTime,
		provider:            provider,
		analysisResultStore: analysisResultStore,
		argsTemplate:        argsTemplate,
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
				a.logPersister.Infof("[%s] The query result evaluation was skipped because \"skipOnNoData\" is true though no data returned. Reason: %v", a.id, err)
				continue
			}
			if err != nil {
				a.logPersister.Errorf("[%s] Unexpected error: %v", a.id, err)
			}
			if expected {
				a.logPersister.Successf("[%s] The query result is expected one", a.id)
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
		a.logPersister.Errorf("[%s] Failed because it found a data point (%s) that is outside the expected range (%s). Performed query: %q", a.id, &outiler, &a.cfg.Expected, a.cfg.Query)
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
		return false, false, fmt.Errorf("failed to run query: %w: performed query: %q", err, a.cfg.Query)
	}
	values := make([]float64, 0, len(points))
	for i := range points {
		values = append(values, points[i].Value)
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
		return false, false, fmt.Errorf("failed to run query to fetch metrics for the previous deployment: %w: performed query: %q", err, a.cfg.Query)
	}
	prevValues := make([]float64, 0, len(prevPoints))
	for i := range prevPoints {
		prevValues = append(prevValues, prevPoints[i].Value)
	}
	if err := compare(values, prevValues, a.cfg.Deviation); err != nil {
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
	canaryQuery, err := a.renderQuery(a.cfg.Query, a.cfg.CanaryArgs, canaryVariantName)
	if err != nil {
		return false, fmt.Errorf("failed to render query template for Canary: %w", err)
	}
	baselineQuery, err := a.renderQuery(a.cfg.Query, a.cfg.BaselineArgs, baselineVariantName)
	if err != nil {
		return false, fmt.Errorf("failed to render query template for Baseline: %w", err)
	}

	// Fetch data points from Canary.
	canaryPoints, err := a.provider.QueryPoints(ctx, canaryQuery, queryRange)
	if err != nil {
		return false, fmt.Errorf("failed to run query to fetch metrics for the Canary variant: %w: performed query: %q", err, canaryQuery)
	}
	canaryPointsCtn := len(canaryPoints)
	a.logPersister.Infof("[%s] Got %d data points for Canary from the query: %q", a.id, canaryPointsCtn, canaryQuery)
	canaryValues := make([]float64, 0, canaryPointsCtn)
	for i := range canaryPoints {
		canaryValues = append(canaryValues, canaryPoints[i].Value)
	}

	// Fetch data points from Baseline.
	baselinePoints, err := a.provider.QueryPoints(ctx, baselineQuery, queryRange)
	if err != nil {
		return false, fmt.Errorf("failed to run query to fetch metrics for the Baseline variant: %w: performed query: %q", err, baselineQuery)
	}
	baselinePointsCtn := len(baselinePoints)
	a.logPersister.Infof("[%s] Got %d data points for Baseline from the query: %q", a.id, baselinePointsCtn, baselineQuery)
	baselineValues := make([]float64, 0, baselinePointsCtn)
	for i := range baselinePoints {
		baselineValues = append(baselineValues, baselinePoints[i].Value)
	}

	if err := compare(canaryValues, baselineValues, a.cfg.Deviation); err != nil {
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
	canaryQuery, err := a.renderQuery(a.cfg.Query, a.cfg.CanaryArgs, canaryVariantName)
	if err != nil {
		return false, fmt.Errorf("failed to render query template for Canary: %w", err)
	}
	primaryQuery, err := a.renderQuery(a.cfg.Query, a.cfg.PrimaryArgs, primaryVariantName)
	if err != nil {
		return false, fmt.Errorf("failed to render query template for Primary: %w", err)
	}

	canaryPoints, err := a.provider.QueryPoints(ctx, canaryQuery, queryRange)
	if err != nil {
		return false, fmt.Errorf("failed to run query to fetch metrics for the Canary variant: %w: performed query: %q", err, canaryQuery)
	}
	canaryValues := make([]float64, 0, len(canaryPoints))
	for i := range canaryPoints {
		canaryValues = append(canaryValues, canaryPoints[i].Value)
	}
	primaryPoints, err := a.provider.QueryPoints(ctx, primaryQuery, queryRange)
	if err != nil {
		return false, fmt.Errorf("failed to run query to fetch metrics for the Primary variant: %w: performed query: %q", err, primaryQuery)
	}
	primaryValues := make([]float64, 0, len(primaryPoints))
	for i := range primaryPoints {
		primaryValues = append(primaryValues, primaryPoints[i].Value)
	}
	if err := compare(canaryValues, primaryValues, a.cfg.Deviation); err != nil {
		a.logPersister.Errorf("[%s] Failed because %v. Performed query for canary: %q. Performed query for primary: %q", a.id, err, canaryQuery, primaryQuery)
		return false, nil
	}
	return true, nil
}

// compare compares the given two samples using Mann-Whitney U test.
// Considered as failure if it deviates in the specified direction as the third argument.
// No error means that the result is expected.
func compare(experiment, control []float64, deviation string) (err error) {
	if len(experiment) == 0 {
		return fmt.Errorf("no data points of Experiment found")
	}
	if len(control) == 0 {
		return fmt.Errorf("no data points of Control found")
	}
	var alternativeHypothesis mannwhitney.LocationHypothesis
	switch deviation {
	case config.AnalysisDeviationEither:
		alternativeHypothesis = mannwhitney.LocationDiffers
	case config.AnalysisDeviationLow:
		alternativeHypothesis = mannwhitney.LocationLess
	case config.AnalysisDeviationHigh:
		alternativeHypothesis = mannwhitney.LocationGreater
	default:
		return fmt.Errorf("unknown deviation %q given", deviation)
	}
	res, err := mannwhitney.MannWhitneyUTest(experiment, control, alternativeHypothesis)
	if errors.Is(err, mannwhitney.ErrSamplesEqual) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to perform the Mann-Whitney U test: %w", err)
	}

	// alpha is the significance level. Typically 5% is used.
	const alpha = 0.05
	// If the p-value is greater than the significance level,
	// we cannot say that the distributions in the two groups differed significantly.
	// See: https://support.minitab.com/en-us/minitab-express/1/help-and-how-to/basic-statistics/inference/how-to/two-samples/mann-whitney-test/interpret-the-results/key-results/
	if res.P > alpha {
		return nil
	}
	return fmt.Errorf("the difference between the medians is statistically significant")
}

// argsTemplate is a collection of available template arguments.
// NOTE: Changing its fields will force users to change the template definition.
type argsTemplate struct {
	// The args that are automatically populated.
	App     appArgs
	K8s     k8sArgs
	Variant variantArgs

	// User-defined custom args.
	VariantCustomArgs map[string]string
	AppCustomArgs     map[string]string
}

// appArgs allows application-specific data to be embedded in the query.
type appArgs struct {
	Name string
	Env  string
}

type k8sArgs struct {
	Namespace string
}

// variantArgs allows variant-specific data to be embedded in the query.
type variantArgs struct {
	// One of "primary", "canary", or "baseline" will be populated.
	Name string
}

// renderQuery applies the given variant args to the query template.
func (a *metricsAnalyzer) renderQuery(queryTemplate string, variantCustomArgs map[string]string, variant string) (string, error) {
	args := argsTemplate{
		Variant:           variantArgs{Name: variant},
		VariantCustomArgs: variantCustomArgs,
		App:               a.argsTemplate.App,
		K8s:               a.argsTemplate.K8s,
		AppCustomArgs:     a.argsTemplate.AppCustomArgs,
	}

	t, err := template.New("AnalysisTemplate").Parse(queryTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse query template: %w", err)
	}

	b := new(bytes.Buffer)
	if err := t.Execute(b, args); err != nil {
		return "", fmt.Errorf("failed to apply template: %w", err)
	}
	return b.String(), err
}
