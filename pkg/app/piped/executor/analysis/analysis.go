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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"text/template"
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

	startTime             time.Time
	previouslyElapsedTime time.Duration
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

// templateArgs allows deployment-specific data to be embedded in the analysis template.
// NOTE: Changing its fields will force users to change the template definition.
type templateArgs struct {
	App struct {
		Name string
		Env  string
	}
	// User-defined custom args.
	Args map[string]string
}

// Execute runs multiple analyses that execute queries against analysis providers at regular intervals.
// An executor runs multiple analyses, an analysis may run a query multiple times.
func (e *Executor) Execute(sig executor.StopSignal) model.StageStatus {
	e.startTime = time.Now()
	ctx := sig.Context()
	options := e.StageConfig.AnalysisStageOptions
	if options == nil {
		e.Logger.Error("missing analysis configuration for ANALYSIS stage")
		return model.StageStatus_STAGE_FAILURE
	}

	templateCfg, ok, err := config.LoadAnalysisTemplate(e.RepoDir)
	if err != nil {
		e.LogPersister.AppendError(err.Error())
		return model.StageStatus_STAGE_FAILURE
	}
	if !ok {
		e.Logger.Info("config file for AnalysisTemplate not found")
		templateCfg = &config.AnalysisTemplateSpec{}
	}

	timeout := time.Duration(options.Duration)
	e.previouslyElapsedTime = e.retrievePreviouslyElapsedTime()
	if e.previouslyElapsedTime > 0 {
		// Restart from the middle.
		timeout -= e.previouslyElapsedTime
	}
	defer e.saveElapsedTime(ctx)

	ctx, cancel := context.WithTimeout(sig.Context(), timeout)
	defer cancel()

	eg, ctx := errgroup.WithContext(ctx)

	// Run analyses with metrics providers.
	mf := metrics.NewFactory(e.Logger)
	for i, m := range options.Metrics {
		// TODO: Encapsulate implementation of analysis as an Analyzer
		templateCfg, err = e.render(*templateCfg, m.Template.Args)
		if err != nil {
			e.LogPersister.AppendError(err.Error())
			return model.StageStatus_STAGE_FAILURE
		}
		cfg, err := e.getMetricsConfig(&m, templateCfg)
		if err != nil {
			e.LogPersister.AppendError(err.Error())
			continue
		}
		provider, err := e.newMetricsProvider(cfg.Provider, mf)
		if err != nil {
			e.LogPersister.AppendError(err.Error())
			continue
		}
		id := fmt.Sprintf("metrics-%d", i)
		eg.Go(func() error {
			runner := func(ctx context.Context) (bool, error) {
				return provider.RunQuery(ctx, cfg.Query, cfg.Expected)
			}
			return e.runAnalysis(ctx, id, provider.Type(), time.Duration(cfg.Interval), runner, cfg.FailureLimit)
		})
	}
	// Run analyses with logging providers.
	lf := log.NewFactory(e.Logger)
	for i, l := range options.Logs {
		templateCfg, err = e.render(*templateCfg, l.Template.Args)
		if err != nil {
			e.LogPersister.AppendError(err.Error())
			return model.StageStatus_STAGE_FAILURE
		}
		cfg, err := e.getLogConfig(&l, templateCfg)
		if err != nil {
			e.LogPersister.AppendError(err.Error())
			continue
		}
		provider, err := e.newLogProvider(cfg.Provider, lf)
		if err != nil {
			e.LogPersister.AppendError(err.Error())
			continue
		}
		id := fmt.Sprintf("log-%d", i)
		eg.Go(func() error {
			runner := func(ctx context.Context) (bool, error) {
				return provider.RunQuery(ctx, cfg.Query)
			}
			return e.runAnalysis(ctx, id, provider.Type(), time.Duration(cfg.Interval), runner, cfg.FailureLimit)
		})
	}
	// Run analyses with http providers.
	for i, h := range options.Https {
		templateCfg, err = e.render(*templateCfg, h.Template.Args)
		if err != nil {
			e.LogPersister.AppendError(err.Error())
			return model.StageStatus_STAGE_FAILURE
		}
		cfg, err := e.getHTTPConfig(&h, templateCfg)
		if err != nil {
			e.LogPersister.AppendError(err.Error())
			continue
		}
		provider := httpprovider.NewProvider(time.Duration(cfg.Timeout))
		id := fmt.Sprintf("http-%d", i)
		eg.Go(func() error {
			runner := func(ctx context.Context) (bool, error) {
				return provider.Run(ctx, cfg)
			}
			return e.runAnalysis(ctx, id, provider.Type(), time.Duration(cfg.Interval), runner, cfg.FailureLimit)
		})
	}

	if err := eg.Wait(); err != nil {
		e.LogPersister.AppendError(fmt.Sprintf("Analysis failed: %s", err.Error()))
		return model.StageStatus_STAGE_FAILURE
	}

	e.LogPersister.AppendSuccess("All analyses were successful.")
	return model.StageStatus_STAGE_SUCCESS
}

// runAnalysis calls `runQuery` function at the given interval and reports back to failureCh
// when the number of failures exceeds the failureLimit.
func (e *Executor) runAnalysis(ctx context.Context, id, providerType string, interval time.Duration, runQuery func(context.Context) (bool, error), failureLimit int) error {
	e.Logger.Info("start the analysis", zap.String("analysis-id", id), zap.String("provider", providerType))
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	failureCount := 0
	for {
		select {
		case <-ticker.C:
			e.LogPersister.AppendInfo(fmt.Sprintf("Start running query against %s", providerType))
			success, err := runQuery(ctx)
			if err != nil {
				// The failure of the query itself is treated as a failure.
				e.LogPersister.AppendError(fmt.Sprintf("Failed to run query: %s", err.Error()))
				success = false
			}
			if success {
				e.LogPersister.AppendSuccess(fmt.Sprintf("The result of the query for %s by analysis '%s' is a success.", providerType, id))
			} else {
				failureCount++
				e.LogPersister.AppendError(fmt.Sprintf("The result of the query for %s by analysis '%s' is a failure. This analysis will fail if it fails %d more times.", providerType, id, failureLimit+1-failureCount))
			}

			if failureCount > failureLimit {
				return fmt.Errorf("anslysis '%s' by %s failed because the failure number exceeded the failure limit %d", id, providerType, failureLimit)
			}
		case <-ctx.Done():
			return nil
		}
	}
}

const elapsedTimeKey = "elapsedTime"

// saveElapsedTime stores the elapsed time of analysis stage into metadata persister.
// The analysis stage can be restarted from the middle even if it ends unexpectedly,
// that's why count should be stored.
func (e *Executor) saveElapsedTime(ctx context.Context) {
	elapsedTime := time.Since(e.startTime) + e.previouslyElapsedTime
	metadata := map[string]string{
		elapsedTimeKey: elapsedTime.String(),
	}
	if err := e.MetadataStore.SetStageMetadata(ctx, e.Stage.Id, metadata); err != nil {
		e.Logger.Error("failed to store metadata", zap.Error(err))
	}
}

// retrievePreviouslyElapsedTime sets the elapsed time of analysis stage by decoding metadata.
func (e *Executor) retrievePreviouslyElapsedTime() time.Duration {
	metadata, ok := e.MetadataStore.GetStageMetadata(e.Stage.Id)
	if !ok {
		return 0
	}
	s, ok := metadata[elapsedTimeKey]
	if !ok {
		return 0
	}
	et, err := time.ParseDuration(s)
	if err != nil {
		e.Logger.Error("unexpected elapsed time is stored", zap.String("stored-value", s), zap.Error(err))
		return 0
	}
	return et
}

func (e *Executor) newMetricsProvider(providerName string, factory *metrics.Factory) (metrics.Provider, error) {
	cfg, ok := e.PipedConfig.GetAnalysisProvider(providerName)
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
	cfg, ok := e.PipedConfig.GetAnalysisProvider(providerName)
	if !ok {
		return nil, fmt.Errorf("unknown provider name %s", providerName)
	}
	provider, err := factory.NewProvider(&cfg)
	if err != nil {
		return nil, err
	}
	return provider, nil
}

// getMetricsConfig renders the given template and returns the metrics config.
// Just returns metrics config if no template specified.
func (e *Executor) getMetricsConfig(templatableCfg *config.TemplatableAnalysisMetrics, templateSpec *config.AnalysisTemplateSpec) (*config.AnalysisMetrics, error) {
	name := templatableCfg.Template.Name
	if name == "" {
		return &templatableCfg.AnalysisMetrics, nil
	}
	cfg, ok := templateSpec.Metrics[name]
	if !ok {
		return nil, fmt.Errorf("analysis template %s not found despite useTemplate specified", name)
	}
	return &cfg, nil
}

// getLogConfig renders the given template and returns the log config.
// Just returns log config if no template specified.
func (e *Executor) getLogConfig(templatableCfg *config.TemplatableAnalysisLog, templateSpec *config.AnalysisTemplateSpec) (*config.AnalysisLog, error) {
	name := templatableCfg.Template.Name
	if name == "" {
		return &templatableCfg.AnalysisLog, nil
	}
	cfg, ok := templateSpec.Logs[name]
	if !ok {
		return nil, fmt.Errorf("analysis template %s not found despite useTemplate specified", name)
	}
	return &cfg, nil
}

// getHTTPConfig renders the given template and returns the http config.
// Just returns http config if no template specified.
func (e *Executor) getHTTPConfig(templatableCfg *config.TemplatableAnalysisHTTP, templateSpec *config.AnalysisTemplateSpec) (*config.AnalysisHTTP, error) {
	name := templatableCfg.Template.Name
	if name == "" {
		return &templatableCfg.AnalysisHTTP, nil
	}
	cfg, ok := templateSpec.HTTPs[name]
	if !ok {
		return nil, fmt.Errorf("analysis template %s not found despite useTemplate specified", name)
	}
	return &cfg, nil
}

// render returns a new AnalysisTemplateSpec, where deployment-specific arguments populated.
func (e *Executor) render(templateCfg config.AnalysisTemplateSpec, customArgs map[string]string) (*config.AnalysisTemplateSpec, error) {
	args := templateArgs{
		Args: customArgs,
		App: struct {
			Name string
			Env  string
		}{
			Name: e.Application.Name,
			// TODO: Populate Environment.
		},
	}

	cfg, err := json.Marshal(templateCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal json: %w", err)
	}
	t, err := template.New("AnalysisTemplate").Parse(string(cfg))
	if err != nil {
		return nil, fmt.Errorf("failed to parse text: %w", err)
	}
	b := new(bytes.Buffer)
	if err := t.Execute(b, args); err != nil {
		return nil, fmt.Errorf("failed to apply template: %w", err)
	}
	newCfg := &config.AnalysisTemplateSpec{}
	err = json.Unmarshal(b.Bytes(), newCfg)
	return newCfg, err
}
