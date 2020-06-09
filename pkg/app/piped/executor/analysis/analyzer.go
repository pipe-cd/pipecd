package analysis

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/kapetaniosci/pipe/pkg/app/piped/executor"
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
	a.logger.Info("start the analysis")
	ticker := time.NewTicker(a.interval)
	defer ticker.Stop()

	failureCount := 0
	for {
		select {
		case <-ticker.C:
			a.logPersister.AppendInfo(fmt.Sprintf("Start running query against %s", a.providerType))
			success, err := a.runQuery(ctx)
			if err != nil {
				// The failure of the query itself is treated as a failure.
				a.logPersister.AppendError(fmt.Sprintf("Failed to run query: %s", err.Error()))
				success = false
			}
			if success {
				a.logPersister.AppendSuccess(fmt.Sprintf("The result of the query for %s by analysis '%s' is a success.", a.providerType, a.id))
			} else {
				failureCount++
				a.logPersister.AppendError(fmt.Sprintf("The result of the query for %s by analysis '%s' is a failure. This analysis will fail if it fails %d more times.", a.providerType, a.id, a.failureLimit+1-failureCount))
			}

			if failureCount > a.failureLimit {
				return fmt.Errorf("anslysis '%s' by %s failed because the failure number exceeded the failure limit %d", a.id, a.providerType, a.failureLimit)
			}
		case <-ctx.Done():
			return nil
		}
	}
}
