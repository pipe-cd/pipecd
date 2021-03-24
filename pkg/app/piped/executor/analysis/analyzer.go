package analysis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipe/pkg/app/piped/analysisprovider/metrics"
	"github.com/pipe-cd/pipe/pkg/app/piped/executor"
)

// analyzer contains a query for an analysis provider.
type analyzer struct {
	id           string
	providerType string
	evaluate     evaluator
	query        string
	interval     time.Duration
	// The analysis will fail, if this value is exceeded,
	failureLimit    int
	nodataAsSuccess bool

	logger       *zap.Logger
	logPersister executor.LogPersister
}

type evaluator func(ctx context.Context, query string) (expected bool, reason string, err error)

func newAnalyzer(
	id string,
	providerType string,
	query string,
	evaluate evaluator,
	interval time.Duration,
	failureLimit int,
	noDataAsSuccess bool,
	logger *zap.Logger,
	logPersister executor.LogPersister,
) *analyzer {
	return &analyzer{
		id:              id,
		providerType:    providerType,
		evaluate:        evaluate,
		query:           query,
		interval:        interval,
		failureLimit:    failureLimit,
		nodataAsSuccess: noDataAsSuccess,
		logPersister:    logPersister,
		logger: logger.With(
			zap.String("analyzer-id", id),
			zap.String("provider-type", providerType),
		),
	}
}

// run starts an analysis which runs the query at the given interval, until the context is done.
// It returns an error when the number of failures exceeds the the failureLimit.
func (a *analyzer) run(ctx context.Context) error {
	ticker := time.NewTicker(a.interval)
	defer ticker.Stop()

	failureCount := 0
	for {
		select {
		case <-ticker.C:
			expected, reason, err := a.evaluate(ctx, a.query)
			switch {
			case errors.Is(err, context.DeadlineExceeded) && ctx.Err() == context.DeadlineExceeded:
				// Ignore parent's context deadline exceeded error, and return immediately.
				return nil
			case errors.Is(err, metrics.ErrNoDataFound) && a.nodataAsSuccess:
				reason = "no data returned but \"nodataAsSuccess\" is true"
				expected = true
			case err != nil:
				reason = fmt.Sprintf("failed to run query: %s", err.Error())
			default:
			}

			if expected {
				a.logPersister.Successf("[%s] The query result is expected one. Reason: %s. Performed query: %s", a.id, reason, a.query)
			} else {
				failureCount++
				a.logPersister.Errorf("[%s] The query result is unexpected. Reason: %s. Performed query: %s", a.id, reason, a.query)
			}

			if failureCount > a.failureLimit {
				return fmt.Errorf("anslysis '%s' failed because the failure number exceeded the failure limit (%d)", a.id, a.failureLimit)
			}
		case <-ctx.Done():
			return nil
		}
	}
}
