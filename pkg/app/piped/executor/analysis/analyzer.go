package analysis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipe/pkg/app/piped/executor"
)

// analyzer contains a query for an analysis provider.
type analyzer struct {
	id           string
	providerType string
	runQuery     queryRunner
	query        string
	interval     time.Duration
	// The analysis will fail, if this value is exceeded,
	failureLimit int

	logger       *zap.Logger
	logPersister executor.LogPersister
}

type queryRunner func(ctx context.Context, query string) (expected bool, reason string, err error)

func newAnalyzer(
	id string,
	providerType string,
	query string,
	runQuery queryRunner,
	interval time.Duration,
	failureLimit int,
	logger *zap.Logger,
	logPersister executor.LogPersister,
) *analyzer {
	return &analyzer{
		id:           id,
		providerType: providerType,
		runQuery:     runQuery,
		query:        query,
		interval:     interval,
		failureLimit: failureLimit,
		logPersister: logPersister,
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
			expected, reason, err := a.runQuery(ctx, a.query)
			if errors.Is(err, context.DeadlineExceeded) && ctx.Err() == context.DeadlineExceeded {
				// Ignore parent's context deadline exceeded error, and return immediately.
				return nil
			}
			// TODO: Consider how to handle the case of analysisprovider.ErrNoValuesFound
			if err != nil {
				// The failure of the query itself is treated as an unexpected result.
				reason = fmt.Sprintf("failed to run query: %s", err.Error())
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
