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
	runQuery     func(ctx context.Context) (bool, error)
	interval     time.Duration
	// The analysis will fail, ff this value is exceeded,
	failureLimit int

	logger       *zap.Logger
	logPersister executor.LogPersister
}

func newAnalyzer(id string, providerType string, runQuery func(ctx context.Context) (bool, error), interval time.Duration, failureLimit int, logger *zap.Logger, logPersister executor.LogPersister) *analyzer {
	l := logger.With(
		zap.String("analyzer-id", id),
		zap.String("provider-type", providerType),
	)
	return &analyzer{
		id:           id,
		providerType: providerType,
		runQuery:     runQuery,
		interval:     interval,
		failureLimit: failureLimit,
		logger:       l,
		logPersister: logPersister,
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
			reason := ""
			success, err := a.runQuery(ctx)
			if errors.Is(err, context.DeadlineExceeded) && ctx.Err() == context.DeadlineExceeded {
				// Ignore parent's context deadline exceeded error, and return immediately.
				return nil
			}
			// TODO: Consider how to handle the case of analysisprovider.ErrNoValuesFound
			if err != nil {
				// The failure of the query itself is treated as a failure.
				reason = fmt.Sprintf("failed to run query: %s", err.Error())
				success = false
			}
			if success {
				a.logPersister.Successf("[%s] The query result is a success.", a.id)
			} else {
				failureCount++
				if reason == "" {
					reason = "the response is not expected value"
				}
				a.logPersister.Errorf("[%s] The query result is a failure. Reason: %s",
					a.id,
					reason,
				)
			}

			if failureCount > a.failureLimit {
				return fmt.Errorf("anslysis '%s' failed because the failure number exceeded the failure limit (%d)", a.id, a.failureLimit)
			}
		case <-ctx.Done():
			return nil
		}
	}
}
