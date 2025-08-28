// Copyright 2025 The PipeCD Authors.
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

package executestage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"time"

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	httpprovider "github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/analysis/analysisprovider/http"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/analysis/analysisprovider/log"
	logfactory "github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/analysis/analysisprovider/log/factory"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/analysis/analysisprovider/metrics"
	metricsfactory "github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/analysis/analysisprovider/metrics/factory"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/analysis/analysisresultstore"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/analysis/config"
)

type executor struct {
	targetDS            *sdk.DeploymentSource[any] // TODO: do not use any
	stageConfig         *config.AnalysisStageOptions
	pluginConfig        *config.PluginConfig
	analysisAppSpec     *config.AnalysisApplicationSpec
	analysisResultStore analysisresultstore.Store
	appName             string
	sharedConfigDir     string

	logger       *zap.Logger
	logPersister sdk.StageLogPersister
	client       *sdk.Client

	startTime           time.Time
	previousElapsedTime time.Duration
}

func ExecuteAnalysisStage(ctx context.Context, input *sdk.ExecuteStageInput[config.AnalysisApplicationSpec], pluginCfg *config.PluginConfig) sdk.StageStatus {
	stageCfg := &config.AnalysisStageOptions{}
	if err := json.Unmarshal(input.Request.StageConfig, stageCfg); err != nil {
		return sdk.StageStatusFailure
	}
	resultStore := analysisresultstore.NewStore(input.Client, input.Logger)

	e := &executor{
		stageConfig:         stageCfg,
		pluginConfig:        pluginCfg,
		analysisAppSpec:     input.Request.TargetDeploymentSource.ApplicationConfig.Spec,
		analysisResultStore: resultStore,
		appName:             input.Request.Deployment.ApplicationName,
		sharedConfigDir:     input.Request.TargetDeploymentSource.SharedConfigDirectory,
		logger:              zap.NewNop(),
		logPersister:        input.Client.LogPersister(),
		client:              input.Client,
	}
	return e.execute(ctx)
}

// Execute spawns and runs multiple analyzer that run a query at the regular time.
// Any of those fail then the stage ends with failure.
func (e *executor) execute(ctx context.Context) sdk.StageStatus {
	e.startTime = time.Now()
	options := e.stageConfig
	if options == nil {
		e.logger.Error("missing analysis configuration for ANALYSIS stage")
		return sdk.StageStatusFailure
	}

	templateCfg, err := config.LoadAnalysisTemplate(e.targetDS.SharedConfigDirectory)
	if errors.Is(err, config.ErrNotFound) {
		e.logger.Info("config file for AnalysisTemplate not found")
		templateCfg = &config.AnalysisTemplateSpec{}
	} else if err != nil {
		e.logPersister.Error(err.Error())
		return sdk.StageStatusFailure
	}

	timeout := time.Duration(options.Duration)
	e.previousElapsedTime, err = e.retrievePreviousElapsedTime(ctx)
	if err != nil {
		e.logPersister.Errorf("Failed to retrieve previous elapsed time: %v", err)
		return sdk.StageStatusFailure
	}
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
		status = sdk.StageStatusSuccess
		doneCh = make(chan struct{})
	)
	defer close(doneCh)
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if !e.checkSkippedByCmd(ctx) {
					continue
				}
				status = sdk.StageStatusSkipped
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
			e.logPersister.Errorf("Failed to get metrics config: %v", err)
			return sdk.StageStatusFailure
		}
		provider, err := e.newMetricsProvider(cfg.Provider, options.Metrics[i])
		if err != nil {
			e.logPersister.Errorf("Failed to generate metrics provider: %v", err)
			return sdk.StageStatusFailure
		}

		id := fmt.Sprintf("metrics-%d", i)
		args := e.buildAppArgs(options.Metrics[i].Template.AppArgs)
		analyzer := newMetricsAnalyzer(id, *cfg, e.startTime, provider, e.analysisResultStore, args, e.logger, e.logPersister)

		eg.Go(func() error {
			e.logPersister.Infof("[%s] Start metrics analyzer every %s with query template: %q", analyzer.id, cfg.Interval.Duration(), cfg.Query)
			return analyzer.run(ctxWithTimeout)
		})
	}
	// Run analyses with logging providers.
	for i := range options.Logs {
		analyzer, err := e.newAnalyzerForLog(i, &options.Logs[i], templateCfg)
		if err != nil {
			e.logPersister.Errorf("Failed to spawn analyzer for %s: %v", options.Logs[i].Provider, err)
			return sdk.StageStatusFailure
		}
		eg.Go(func() error {
			e.logPersister.Infof("[%s] Start log analyzer", analyzer.id)
			return analyzer.run(ctxWithTimeout)
		})
	}
	// Run analyses with http providers.
	for i := range options.HTTPS {
		analyzer, err := e.newAnalyzerForHTTP(i, &options.HTTPS[i], templateCfg)
		if err != nil {
			e.logPersister.Errorf("Failed to spawn analyzer for HTTP: %v", err)
			return sdk.StageStatusFailure
		}
		eg.Go(func() error {
			e.logPersister.Infof("[%s] Start http analyzer", analyzer.id)
			return analyzer.run(ctxWithTimeout)
		})
	}

	if err := eg.Wait(); err != nil {
		e.logPersister.Errorf("Analysis failed: %s", err.Error())
		return sdk.StageStatusFailure
	}

	// TODO: OK?
	if status != sdk.StageStatusSuccess {
		return status
	}

	e.logPersister.Success("All analyses were successful")
	err = e.analysisResultStore.PutLatestAnalysisResult(ctx, &analysisresultstore.AnalysisResult{
		StartTime: e.startTime.Unix(),
	})
	if err != nil {
		e.logger.Error("failed to send the analysis result", zap.Error(err))
	}
	return status
}

const elapsedTimeKey = "elapsedTime"

// saveElapsedTime stores the elapsed time of analysis stage into metadata persister.
// The analysis stage can be restarted from the middle even if it ends unexpectedly,
// that's why count should be stored.
func (e *executor) saveElapsedTime(ctx context.Context) {
	elapsedTime := time.Since(e.startTime) + e.previousElapsedTime
	metadata := map[string]string{
		elapsedTimeKey: elapsedTime.String(),
	}
	if err := e.client.PutStageMetadataMulti(ctx, metadata); err != nil {
		e.logger.Error("failed to store metadata", zap.Error(err))
	}
}

// retrievePreviousElapsedTime sets the elapsed time of analysis stage by decoding metadata.
func (e *executor) retrievePreviousElapsedTime(ctx context.Context) (time.Duration, error) {
	s, ok, err := e.client.GetStageMetadata(ctx, elapsedTimeKey)
	if err != nil {
		return 0, err
	}
	if !ok {
		return 0, nil
	}
	et, err := time.ParseDuration(s)
	if err != nil {
		e.logger.Error("unexpected elapsed time is stored", zap.String("stored-value", s), zap.Error(err))
		return 0, err
	}
	return et, nil
}

func (e *executor) newAnalyzerForLog(i int, templatable *config.TemplatableAnalysisLog, templateCfg *config.AnalysisTemplateSpec) (*analyzer, error) {
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
	return newAnalyzer(id, provider.Type(), cfg.Query, runner, time.Duration(cfg.Interval), cfg.FailureLimit, cfg.SkipOnNoData, e.logger, e.logPersister), nil
}

func (e *executor) newAnalyzerForHTTP(i int, templatable *config.TemplatableAnalysisHTTP, templateCfg *config.AnalysisTemplateSpec) (*analyzer, error) {
	cfg, err := e.getHTTPConfig(templatable, templateCfg)
	if err != nil {
		return nil, err
	}
	provider := httpprovider.NewProvider(time.Duration(cfg.Timeout))
	id := fmt.Sprintf("http-%d", i)
	runner := func(ctx context.Context, query string) (bool, string, error) {
		return provider.Run(ctx, cfg)
	}
	return newAnalyzer(id, provider.Type(), "", runner, time.Duration(cfg.Interval), cfg.FailureLimit, cfg.SkipOnNoData, e.logger, e.logPersister), nil
}

func (e *executor) newMetricsProvider(providerName string, templatable config.TemplatableAnalysisMetrics) (metrics.Provider, error) {
	cfg, ok := e.pluginConfig.GetAnalysisProvider(providerName)
	if !ok {
		return nil, fmt.Errorf("unknown provider name %s", providerName)
	}
	provider, err := metricsfactory.NewProvider(&templatable, &cfg, e.logger)
	if err != nil {
		return nil, err
	}
	return provider, nil
}

func (e *executor) newLogProvider(providerName string) (log.Provider, error) {
	cfg, ok := e.pluginConfig.GetAnalysisProvider(providerName)
	if !ok {
		return nil, fmt.Errorf("unknown provider name %s", providerName)
	}
	provider, err := logfactory.NewProvider(&cfg, e.logger)
	if err != nil {
		return nil, err
	}
	return provider, nil
}

// getMetricsConfig renders the given template and returns the metrics config.
// Just returns metrics config if no template specified.
func (e *executor) getMetricsConfig(templatableCfg config.TemplatableAnalysisMetrics, templateCfg *config.AnalysisTemplateSpec) (*config.AnalysisMetrics, error) {
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
func (e *executor) getLogConfig(templatableCfg *config.TemplatableAnalysisLog, templateCfg *config.AnalysisTemplateSpec) (*config.AnalysisLog, error) {
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
func (e *executor) getHTTPConfig(templatableCfg *config.TemplatableAnalysisHTTP, templateCfg *config.AnalysisTemplateSpec) (*config.AnalysisHTTP, error) {
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

func (e *executor) buildAppArgs(customArgs map[string]string) argsTemplate {
	// merge custom args specified under stage config and application plugin spec
	// the values under stage config has higher priority
	appCustomArgs := maps.Clone(e.analysisAppSpec.AppCustomArgs)
	maps.Copy(appCustomArgs, customArgs)
	args := argsTemplate{
		App:           appArgs{Name: e.appName},
		AppCustomArgs: appCustomArgs,

		// This is for temporary support for the `{{ .K8s.Namespace }}` syntax in the query template.
		// Please use `{{ .AppCustomArgs.k8sNamespace }}` instead.
		K8s: map[string]string{"Namespace": appCustomArgs["k8sNamespace"]},
	}
	return args
}

func (e *executor) checkSkippedByCmd(ctx context.Context) bool {
	for cmd, err := range e.client.ListStageCommands(ctx, sdk.CommandTypeSkipStage) {
		if err != nil {
			e.logger.Error("failed to list stage skip commands in analysis stage", zap.Error(err))
			return false
		}

		md := fmt.Sprintf("SkippedBy: %s", cmd.Commander)
		if err := e.client.PutStageMetadata(ctx, sdk.MetadataKeyStageDisplay, md); err != nil {
			e.logPersister.Errorf("Unable to save the commander who skipped the stage information to deployment, %v", err)
		}
		e.logPersister.Infof("Got the skip command from %q", cmd.Commander)
		e.logPersister.Infof("This stage has been skipped by user (%s)", cmd.Commander)

		return true
	}
	return false
}
