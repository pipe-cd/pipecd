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

	httpprovider "github.com/kapetaniosci/pipe/pkg/app/piped/analysisprovider/http"
	"github.com/kapetaniosci/pipe/pkg/app/piped/analysisprovider/log"
	"github.com/kapetaniosci/pipe/pkg/app/piped/analysisprovider/metrics"
	"github.com/kapetaniosci/pipe/pkg/app/piped/executor"
	"github.com/kapetaniosci/pipe/pkg/config"
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

	ctx, cancel := context.WithTimeout(sig.Context(), time.Duration(options.Duration))
	defer cancel()

	eg, ctx := errgroup.WithContext(ctx)
	// Run analyses with metrics providers.
	mf := metrics.NewFactory(e.Logger)
	for _, m := range options.Metrics {
		cfg, err := e.getMetricsCfg(&m)
		if err != nil {
			e.LogPersister.AppendError(err.Error())
			continue
		}
		provider, err := e.newMetricsProvider(cfg.Provider, mf)
		if err != nil {
			e.LogPersister.AppendError(err.Error())
			continue
		}
		eg.Go(func() error {
			runner := func(ctx context.Context) (bool, error) {
				return provider.RunQuery(ctx, cfg.Query, cfg.Expected)
			}
			return e.runAnalysis(ctx, time.Duration(cfg.Interval), provider.Type(), runner, cfg.FailureLimit)
		})
	}
	// Run analyses with logging providers.
	lf := log.NewFactory(e.Logger)
	for _, l := range options.Logs {
		cfg, err := e.getLogCfg(&l)
		if err != nil {
			e.LogPersister.AppendError(err.Error())
			continue
		}
		provider, err := e.newLogProvider(cfg.Provider, lf)
		if err != nil {
			e.LogPersister.AppendError(err.Error())
			continue
		}
		eg.Go(func() error {
			runner := func(ctx context.Context) (bool, error) {
				return provider.RunQuery(cfg.Query, cfg.FailureLimit)
			}
			return e.runAnalysis(ctx, time.Duration(cfg.Interval), provider.Type(), runner, cfg.FailureLimit)
		})
	}
	// Run analyses with http providers.
	for _, h := range options.Https {
		cfg, err := e.getHTTPCfg(&h)
		if err != nil {
			e.LogPersister.AppendError(err.Error())
			continue
		}
		provider := httpprovider.NewProvider(time.Duration(cfg.Timeout))
		eg.Go(func() error {
			runner := func(ctx context.Context) (bool, error) {
				return provider.Run(ctx, cfg)
			}
			return e.runAnalysis(ctx, time.Duration(cfg.Interval), provider.Type(), runner, cfg.FailureLimit)
		})
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

// getMetricsCfg renders the given template and returns the metrics config.
// Just returns metrics config if no template specified.
func (e *Executor) getMetricsCfg(templatableCfg *config.TemplatableAnalysisMetrics) (*config.AnalysisMetrics, error) {
	name := templatableCfg.UseTemplate
	if name == "" {
		return &templatableCfg.AnalysisMetrics, nil
	}
	// TODO: Load template config from .piped underneath e.RepoDir
	templateSpec := &config.AnalysisTemplateSpec{}
	cfg, ok := templateSpec.Metrics[name]
	if !ok {
		return nil, fmt.Errorf("analysis template %s not found", name)
	}
	// TODO: Render the application specific data into a template.
	return &cfg, nil
}

// getLogCfg renders the given template and returns the log config.
// Just returns log config if no template specified.
func (e *Executor) getLogCfg(templatableCfg *config.TemplatableAnalysisLog) (*config.AnalysisLog, error) {
	name := templatableCfg.UseTemplate
	if name == "" {
		return &templatableCfg.AnalysisLog, nil
	}
	// TODO: Load template config from .piped underneath e.RepoDir
	templateSpec := &config.AnalysisTemplateSpec{}
	cfg, ok := templateSpec.Logs[name]
	if !ok {
		return nil, fmt.Errorf("analysis template %s not found", name)
	}
	// TODO: Render the application specific data into a template.
	return &cfg, nil
}

// getHTTPCfg renders the given template and returns the http config.
// Just returns http config if no template specified.
func (e *Executor) getHTTPCfg(templatableCfg *config.TemplatableAnalysisHTTP) (*config.AnalysisHTTP, error) {
	name := templatableCfg.UseTemplate
	if name == "" {
		return &templatableCfg.AnalysisHTTP, nil
	}
	// TODO: Load template config from .piped underneath e.RepoDir
	templateSpec := &config.AnalysisTemplateSpec{}
	cfg, ok := templateSpec.HTTPs[name]
	if !ok {
		return nil, fmt.Errorf("analysis template %s not found", name)
	}
	// TODO: Render the application specific data into a template.
	return &cfg, nil
}
