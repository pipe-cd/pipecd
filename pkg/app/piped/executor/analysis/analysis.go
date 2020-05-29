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

package analysis

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/kapetaniosci/pipe/pkg/app/piped/analysisprovider/log"
	"github.com/kapetaniosci/pipe/pkg/app/piped/analysisprovider/metrics"
	"github.com/kapetaniosci/pipe/pkg/app/piped/executor"
	"github.com/kapetaniosci/pipe/pkg/model"
)

type Executor struct {
	executor.Input
	// The number of queries executed per provider.
	queryCount map[string]int
	mu         sync.RWMutex
}

type registerer interface {
	Register(stage model.Stage, f executor.Factory) error
}

func Register(r registerer) {
	f := func(in executor.Input) executor.Executor {
		return &Executor{
			Input: in,
		}
	}
	r.Register(model.StageAnalysis, f)
}

// Execute runs multiple analyses that execute queries against analysis providers at regular intervals.
// An executor runs multiple analyses, an analysis may run a query multiple times.
func (e *Executor) Execute(sig executor.StopSignal) model.StageStatus {
	ctx := sig.Context()
	e.setQueryCount()
	defer e.saveQueryCount(ctx)

	options := e.StageConfig.AnalysisStageOptions
	if options == nil {
		e.Logger.Error("missing analysis configuration for ANALYSIS stage")
		return model.StageStatus_STAGE_FAILURE
	}

	ctx, cancel := context.WithTimeout(ctx, time.Duration(options.Duration))
	defer cancel()

	eg, ctx := errgroup.WithContext(ctx)
	// Run analyses with metrics providers.
	mf := metrics.NewFactory(e.Logger)
	for _, m := range options.Metrics {
		provider, err := e.newMetricsProvider(m.Provider, mf)
		if err != nil {
			e.LogPersister.AppendError(err.Error())
			continue
		}
		eg.Go(func() error {
			return e.runAnalysis(ctx, time.Duration(m.Interval), provider.Type(), func(ctx context.Context) (bool, error) {
				return provider.RunQuery(ctx, m.Query, m.Expected)
			}, m.FailureLimit)
		})
	}
	// Run analyses with logging providers.
	lf := log.NewFactory(e.Logger)
	for _, l := range options.Logs {
		provider, err := e.newLogProvider(l.Provider, lf)
		if err != nil {
			e.LogPersister.AppendError(err.Error())
			continue
		}
		eg.Go(func() error {
			return e.runAnalysis(ctx, time.Duration(l.Interval), provider.Type(), func(ctx context.Context) (bool, error) {
				return provider.RunQuery(l.Query, l.FailureLimit)
			}, l.FailureLimit)
		})
	}
	// TODO: Make HTTP analysis part of metrics provider.
	for range options.Https {

	}

	if err := eg.Wait(); err != nil {
		return model.StageStatus_STAGE_FAILURE
	}

	return model.StageStatus_STAGE_SUCCESS
}

// runAnalysis calls `runQuery` function at the given interval and reports back to failureCh
// when the number of failures exceeds the failureLimit.
func (e *Executor) runAnalysis(ctx context.Context, interval time.Duration, providerType string, runQuery func(context.Context) (bool, error), failureLimit int) error {
	e.Logger.Info("start the analysis", zap.String("provider", providerType))
	// TODO: Address the case when using template
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	failureCount := 0
	for {
		select {
		case <-ticker.C:
			success, err := runQuery(ctx)
			if err != nil {
				e.Logger.Error("failed to run query", zap.Error(err))
				// TODO: Decide how to handle query failures.
				success = false
			}
			if !success {
				failureCount++
			}

			e.mu.Lock()
			e.queryCount[providerType]++
			e.mu.Unlock()
			e.saveQueryCount(ctx)

			if failureCount > failureLimit {
				return fmt.Errorf("anslysis by %s failed", providerType)
			}
		case <-ctx.Done():
			return nil
		}
	}
}

// saveQueryCount stores metadata into metadata persister.
// The analysis stage can be restarted from the middle even if it ends unexpectedly,
// that's why count should be stored.
func (e *Executor) saveQueryCount(ctx context.Context) {
	// Copy to local variable to avoid to lock in a long time.
	e.mu.RLock()
	qc := make(map[string]int, len(e.queryCount))
	for k, v := range e.queryCount {
		qc[k] = v
	}
	e.mu.RUnlock()

	if err := e.MetadataStore.SetStageMetadata(ctx, e.Stage.Id, qc); err != nil {
		e.Logger.Error("failed to store metadata", zap.Error(err))
	}
}

// setQueryCount decodes metadata and populates query count to own field.
func (e *Executor) setQueryCount() {
	err := e.MetadataStore.GetStageMetadata(e.Stage.Id, &e.queryCount)
	if err != nil {
		e.Logger.Error("failed to get stage metadata", zap.Error(err))
		e.queryCount = make(map[string]int)
	}
}

func (e *Executor) newMetricsProvider(providerName string, factory *metrics.Factory) (metrics.Provider, error) {
	// TODO: Address the case when using template
	cfg, ok := e.PipedConfig.GetProvider(providerName)
	if !ok {
		return nil, fmt.Errorf("unknown provider name %s", providerName)
	}
	provider, err := factory.NewProvider(&cfg)
	if err != nil {
		return nil, err
	}
	return provider, nil
}

func (e *Executor) newLogProvider(providerName string, factory *log.Factory) (log.Provider, error) {
	// TODO: Address the case when using template
	cfg, ok := e.PipedConfig.GetProvider(providerName)
	if !ok {
		return nil, fmt.Errorf("unknown provider name %s", providerName)
	}
	provider, err := factory.NewProvider(&cfg)
	if err != nil {
		return nil, err
	}
	return provider, nil
}
