// Copyright 2023 The PipeCD Authors.
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
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/app/piped/analysisprovider/metrics"
	"github.com/pipe-cd/pipecd/pkg/app/piped/executor"
)

// analyzer contains a query for an analysis provider.
type analyzer struct {
	id           string
	providerType string
	evaluate     evaluator
	query        string
	interval     time.Duration
	// The analysis will fail, if this value is exceeded,
	failureLimit int
	skipOnNoData bool

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
	skipOnNodata bool,
	logger *zap.Logger,
	logPersister executor.LogPersister,
) *analyzer {
	return &analyzer{
		id:           id,
		providerType: providerType,
		evaluate:     evaluate,
		query:        query,
		interval:     interval,
		failureLimit: failureLimit,
		skipOnNoData: skipOnNodata,
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
			expected, reason, err := a.evaluate(ctx, a.query)
			// Ignore parent's context deadline exceeded error, and return immediately.
			if errors.Is(err, context.DeadlineExceeded) && ctx.Err() == context.DeadlineExceeded {
				return nil
			}
			if errors.Is(err, metrics.ErrNoDataFound) && a.skipOnNoData {
				a.logPersister.Infof("[%s] The query result evaluation was skipped because \"skipOnNoData\" is true even though no data returned. Reason: %v. Performed query: %q", a.id, err, a.query)
				continue
			}
			if err != nil {
				reason = fmt.Sprintf("failed to run query: %s", err.Error())
			}

			if expected {
				a.logPersister.Successf("[%s] The query result is expected one. Reason: %s. Performed query: %q", a.id, reason, a.query)
				continue
			}

			a.logPersister.Errorf("[%s] The query result is unexpected. Reason: %s. Performed query: %q", a.id, reason, a.query)
			failureCount++
			if failureCount > a.failureLimit {
				return fmt.Errorf("analysis '%s' failed because the failure number exceeded the failure limit (%d)", a.id, a.failureLimit)
			}
		case <-ctx.Done():
			return nil
		}
	}
}
