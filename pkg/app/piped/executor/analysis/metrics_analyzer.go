package analysis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipe/pkg/app/piped/analysisprovider/metrics"
	"github.com/pipe-cd/pipe/pkg/app/piped/executor"
	"github.com/pipe-cd/pipe/pkg/config"
)

type metricsAnalyzer struct {
	id           string
	cfg          config.AnalysisMetrics
	provider     metrics.Provider
	logger       *zap.Logger
	logPersister executor.LogPersister
}

func newMetricsAnalyzer(id string, cfg config.AnalysisMetrics, provider metrics.Provider, logger *zap.Logger, logPersister executor.LogPersister) *metricsAnalyzer {
	return &metricsAnalyzer{
		id:           id,
		cfg:          cfg,
		provider:     provider,
		logPersister: logPersister,
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
			// FIXME: Implement ADA strategies other than THRESHOLD
			switch a.cfg.Strategy {
			case config.StrategyThreshold:
				expected, err = a.analyzeWithThreshold(ctx)
			case config.StrategyPrevious:
			case config.StrategyCanaryBaseline:
			case config.StrategyCanaryPrimary:
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
			a.logPersister.Errorf("[%s] The query result is unexpected. Performed query: %q", a.id, a.cfg.Query)
			failureCount++
			if failureCount > a.cfg.FailureLimit {
				return fmt.Errorf("analysis '%s' failed because the failure number exceeded the failure limit (%d)", a.id, a.cfg.FailureLimit)
			}
		case <-ctx.Done():
			return nil
		}
	}
}

// Return false if any data point is out of the prediction range.
func (a *metricsAnalyzer) analyzeWithThreshold(ctx context.Context) (expected bool, err error) {
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
	expected = true
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

	a.logPersister.Successf("[%s] The query result is expected one. Performed query: %q", a.id, a.cfg.Query)
	return true, nil
}
