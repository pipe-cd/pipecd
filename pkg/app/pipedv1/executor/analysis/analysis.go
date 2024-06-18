// Copyright 2024 The PipeCD Authors.
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
	"golang.org/x/sync/errgroup"

	httpprovider "github.com/pipe-cd/pipecd/pkg/app/pipedv1/analysisprovider/http"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/analysisprovider/log"
	logfactory "github.com/pipe-cd/pipecd/pkg/app/pipedv1/analysisprovider/log/factory"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/analysisprovider/metrics"
	metricsfactory "github.com/pipe-cd/pipecd/pkg/app/pipedv1/analysisprovider/metrics/factory"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/executor"
	"github.com/pipe-cd/pipecd/pkg/config"
	"github.com/pipe-cd/pipecd/pkg/model"
)

const (
	skippedByKey = "SkippedBy"
)

type Executor struct {
	executor.Input

	repoDir             string
	config              *config.Config
	startTime           time.Time
	previousElapsedTime time.Duration
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

// Execute spawns and runs multiple analyzer that run a query at the regular time.
// Any of those fail then the stage ends with failure.
func (e *Executor) Execute(sig executor.StopSignal) model.StageStatus {
	e.startTime = time.Now()
	ctx := sig.Context()
	options := e.StageConfig.AnalysisStageOptions
	if options == nil {
		e.Logger.Error("missing analysis configuration for ANALYSIS stage")
		return model.StageStatus_STAGE_FAILURE
	}

	ds, err := e.TargetDSP.Get(ctx, e.LogPersister)
	if err != nil {
		e.LogPersister.Errorf("Failed to prepare running deploy source data (%v)", err)
		return model.StageStatus_STAGE_FAILURE
	}
	e.repoDir = ds.RepoDir
	e.config = ds.ApplicationConfig

	templateCfg, err := config.LoadAnalysisTemplate(e.repoDir)
	if errors.Is(err, config.ErrNotFound) {
		e.Logger.Info("config file for AnalysisTemplate not found")
		templateCfg = &config.AnalysisTemplateSpec{}
	} else if err != nil {
		e.LogPersister.Error(err.Error())
		return model.StageStatus_STAGE_FAILURE
	}

	timeout := time.Duration(options.Duration)
	e.previousElapsedTime = e.retrievePreviousElapsedTime()
	if e.previousElapsedTime > 0 {
		// Restart from the middle.
		timeout -= e.previousElapsedTime
	}
	defer e.saveElapsedTime(ctx)

	ctxWithTimeout, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	eg, ctxWithTimeout := errgroup.WithContext(ctxWithTimeout)

	// Sync the skip command.
	var (
		status = model.StageStatus_STAGE_SUCCESS
		doneCh = make(chan struct{})
	)
	defer close(doneCh)
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if !e.checkSkipped(ctx) {
					continue
				}
				status = model.StageStatus_STAGE_SKIPPED
				// Stop the context to cancel all running analyses.
				cancel()
				return
			case <-doneCh:
				return
			}
		}
	}()

	// Run analyses with metrics providers.
	for i := range options.Metrics {
		cfg, err := e.getMetricsConfig(options.Metrics[i], templateCfg)
		if err != nil {
			e.LogPersister.Errorf("Failed to get metrics config: %v", err)
			return model.StageStatus_STAGE_FAILURE
		}
		provider, err := e.newMetricsProvider(cfg.Provider, options.Metrics[i])
		if err != nil {
			e.LogPersister.Errorf("Failed to generate metrics provider: %v", err)
			return model.StageStatus_STAGE_FAILURE
		}

		id := fmt.Sprintf("metrics-%d", i)
		args := e.buildAppArgs(options.Metrics[i].Template.AppArgs)
		analyzer := newMetricsAnalyzer(id, *cfg, e.startTime, provider, e.AnalysisResultStore, args, e.Logger, e.LogPersister)

		eg.Go(func() error {
			e.LogPersister.Infof("[%s] Start metrics analyzer every %s with query template: %q", analyzer.id, cfg.Interval.Duration(), cfg.Query)
			return analyzer.run(ctxWithTimeout)
		})
	}
	// Run analyses with logging providers.
	for i := range options.Logs {
		analyzer, err := e.newAnalyzerForLog(i, &options.Logs[i], templateCfg)
		if err != nil {
			e.LogPersister.Errorf("Failed to spawn analyzer for %s: %v", options.Logs[i].Provider, err)
			return model.StageStatus_STAGE_FAILURE
		}
		eg.Go(func() error {
			e.LogPersister.Infof("[%s] Start log analyzer", analyzer.id)
			return analyzer.run(ctxWithTimeout)
		})
	}
	// Run analyses with http providers.
	for i := range options.HTTPS {
		analyzer, err := e.newAnalyzerForHTTP(i, &options.HTTPS[i], templateCfg)
		if err != nil {
			e.LogPersister.Errorf("Failed to spawn analyzer for HTTP: %v", err)
			return model.StageStatus_STAGE_FAILURE
		}
		eg.Go(func() error {
			e.LogPersister.Infof("[%s] Start http analyzer", analyzer.id)
			return analyzer.run(ctxWithTimeout)
		})
	}

	if err := eg.Wait(); err != nil {
		e.LogPersister.Errorf("Analysis failed: %s", err.Error())
		return model.StageStatus_STAGE_FAILURE
	}

	status = executor.DetermineStageStatus(sig.Signal(), e.Stage.Status, status)
	if status != model.StageStatus_STAGE_SUCCESS {
		return status
	}

	e.LogPersister.Success("All analyses were successful")
	err = e.AnalysisResultStore.PutLatestAnalysisResult(ctx, &model.AnalysisResult{
		StartTime: e.startTime.Unix(),
	})
	if err != nil {
		e.Logger.Error("failed to send the analysis result", zap.Error(err))
	}
	return status
}

const elapsedTimeKey = "elapsedTime"

// saveElapsedTime stores the elapsed time of analysis stage into metadata persister.
// The analysis stage can be restarted from the middle even if it ends unexpectedly,
// that's why count should be stored.
func (e *Executor) saveElapsedTime(ctx context.Context) {
	elapsedTime := time.Since(e.startTime) + e.previousElapsedTime
	metadata := map[string]string{
		elapsedTimeKey: elapsedTime.String(),
	}
	if err := e.MetadataStore.Stage(e.Stage.Id).PutMulti(ctx, metadata); err != nil {
		e.Logger.Error("failed to store metadata", zap.Error(err))
	}
}

// retrievePreviousElapsedTime sets the elapsed time of analysis stage by decoding metadata.
func (e *Executor) retrievePreviousElapsedTime() time.Duration {
	s, ok := e.MetadataStore.Stage(e.Stage.Id).Get(elapsedTimeKey)
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

func (e *Executor) newAnalyzerForLog(i int, templatable *config.TemplatableAnalysisLog, templateCfg *config.AnalysisTemplateSpec) (*analyzer, error) {
	cfg, err := e.getLogConfig(templatable, templateCfg)
	if err != nil {
		return nil, err
	}
	provider, err := e.newLogProvider(cfg.Provider)
	if err != nil {
		return nil, err
	}
	id := fmt.Sprintf("log-%d", i)
	runner := func(ctx context.Context, query string) (bool, string, error) {
		return provider.Evaluate(ctx, query)
	}
	return newAnalyzer(id, provider.Type(), cfg.Query, runner, time.Duration(cfg.Interval), cfg.FailureLimit, cfg.SkipOnNoData, e.Logger, e.LogPersister), nil
}

func (e *Executor) newAnalyzerForHTTP(i int, templatable *config.TemplatableAnalysisHTTP, templateCfg *config.AnalysisTemplateSpec) (*analyzer, error) {
	cfg, err := e.getHTTPConfig(templatable, templateCfg)
	if err != nil {
		return nil, err
	}
	provider := httpprovider.NewProvider(time.Duration(cfg.Timeout))
	id := fmt.Sprintf("http-%d", i)
	runner := func(ctx context.Context, query string) (bool, string, error) {
		return provider.Run(ctx, cfg)
	}
	return newAnalyzer(id, provider.Type(), "", runner, time.Duration(cfg.Interval), cfg.FailureLimit, cfg.SkipOnNoData, e.Logger, e.LogPersister), nil
}

func (e *Executor) newMetricsProvider(providerName string, templatable config.TemplatableAnalysisMetrics) (metrics.Provider, error) {
	cfg, ok := e.PipedConfig.GetAnalysisProvider(providerName)
	if !ok {
		return nil, fmt.Errorf("unknown provider name %s", providerName)
	}
	provider, err := metricsfactory.NewProvider(&templatable, &cfg, e.Logger)
	if err != nil {
		return nil, err
	}
	return provider, nil
}

func (e *Executor) newLogProvider(providerName string) (log.Provider, error) {
	cfg, ok := e.PipedConfig.GetAnalysisProvider(providerName)
	if !ok {
		return nil, fmt.Errorf("unknown provider name %s", providerName)
	}
	provider, err := logfactory.NewProvider(&cfg, e.Logger)
	if err != nil {
		return nil, err
	}
	return provider, nil
}

// getMetricsConfig renders the given template and returns the metrics config.
// Just returns metrics config if no template specified.
func (e *Executor) getMetricsConfig(templatableCfg config.TemplatableAnalysisMetrics, templateCfg *config.AnalysisTemplateSpec) (*config.AnalysisMetrics, error) {
	name := templatableCfg.Template.Name
	if name == "" {
		return &templatableCfg.AnalysisMetrics, nil
	}

	cfg, ok := templateCfg.Metrics[name]
	if !ok {
		return nil, fmt.Errorf("analysis template %s not found despite template specified", name)
	}
	return &cfg, nil
}

// getLogConfig renders the given template and returns the log config.
// Just returns log config if no template specified.
func (e *Executor) getLogConfig(templatableCfg *config.TemplatableAnalysisLog, templateCfg *config.AnalysisTemplateSpec) (*config.AnalysisLog, error) {
	name := templatableCfg.Template.Name
	if name == "" {
		return &templatableCfg.AnalysisLog, nil
	}

	cfg, ok := templateCfg.Logs[name]
	if !ok {
		return nil, fmt.Errorf("analysis template %s not found despite template specified", name)
	}
	return &cfg, nil
}

// getHTTPConfig renders the given template and returns the http config.
// Just returns http config if no template specified.
func (e *Executor) getHTTPConfig(templatableCfg *config.TemplatableAnalysisHTTP, templateCfg *config.AnalysisTemplateSpec) (*config.AnalysisHTTP, error) {
	name := templatableCfg.Template.Name
	if name == "" {
		return &templatableCfg.AnalysisHTTP, nil
	}

	cfg, ok := templateCfg.HTTPS[name]
	if !ok {
		return nil, fmt.Errorf("analysis template %s not found despite template specified", name)
	}
	return &cfg, nil
}

func (e *Executor) buildAppArgs(customArgs map[string]string) argsTemplate {
	args := argsTemplate{
		App: appArgs{
			Name: e.Application.Name,
			// TODO: Populate Env
			Env: "",
		},
		AppCustomArgs: customArgs,
	}
	if e.config.Kind != config.KindKubernetesApp {
		return args
	}
	namespace := "default"
	if n := e.config.KubernetesApplicationSpec.Input.Namespace; n != "" {
		namespace = n
	}
	args.K8s.Namespace = namespace
	return args
}

func (e *Executor) checkSkipped(ctx context.Context) bool {
	var skipCmd *model.ReportableCommand
	commands := e.CommandLister.ListCommands()

	for i, cmd := range commands {
		if cmd.GetSkipStage() != nil {
			skipCmd = &commands[i]
			break
		}
	}
	if skipCmd == nil {
		return false
	}

	if err := e.MetadataStore.Stage(e.Stage.Id).Put(ctx, skippedByKey, skipCmd.Commander); err != nil {
		e.LogPersister.Errorf("Unable to save the commander who skipped the stage information to deployment, %v", err)
	}
	e.LogPersister.Infof("Got the skip command from %q", skipCmd.Commander)
	e.LogPersister.Infof("This stage has been skipped by user (%s)", skipCmd.Commander)

	if err := skipCmd.Report(ctx, model.CommandStatus_COMMAND_SUCCEEDED, nil, nil); err != nil {
		e.Logger.Error("failed to report handled command", zap.Error(err))
	}
	return true
}
