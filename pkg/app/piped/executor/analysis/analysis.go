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
	"io/ioutil"
	"time"

	"go.uber.org/zap"

	"github.com/kapetaniosci/pipe/pkg/app/piped/executor"
	"github.com/kapetaniosci/pipe/pkg/app/piped/executor/analysis/log"
	"github.com/kapetaniosci/pipe/pkg/app/piped/executor/analysis/log/stackdriver"
	"github.com/kapetaniosci/pipe/pkg/app/piped/executor/analysis/metric"
	"github.com/kapetaniosci/pipe/pkg/app/piped/executor/analysis/metric/datadog"
	"github.com/kapetaniosci/pipe/pkg/app/piped/executor/analysis/metric/prometheus"
	"github.com/kapetaniosci/pipe/pkg/config"
	"github.com/kapetaniosci/pipe/pkg/model"
)

type Executor struct {
	executor.Input
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

// providerResult describes a providerResult of the query for provider.
type providerResult struct {
	success  bool
	provider string
}

// Execute runs multiple analyses that execute queries against analysis providers at regular intervals.
// An executor runs multiple analyses, an analysis runs multiple queries.
func (e *Executor) Execute(sig executor.StopSignal) model.StageStatus {
	ctx := sig.Context()
	queryCount := e.getQueryCount()
	defer e.saveQueryCount(ctx, queryCount)

	options, err := e.getStageOptions()
	if err != nil {
		e.Logger.Error("failed to get analysis options", zap.Error(err))
		return model.StageStatus_STAGE_FAILURE
	}

	ctx, cancel := context.WithTimeout(ctx, time.Duration(options.Duration))
	defer cancel()

	resultCh := make(chan providerResult)
	// Run analyses with metrics providers.
	for _, m := range options.Metrics {
		provider, err := e.newMetricsProvider(&m)
		if err != nil {
			e.LogPersister.AppendError(err.Error())
			continue
		}
		go e.runAnalysis(ctx, time.Duration(m.Interval), provider.Type(), func(ctx context.Context) (bool, error) {
			return provider.RunQuery(ctx, m.Query, m.Expected)
		}, resultCh)
	}
	// Run analyses with logging providers.
	for _, l := range options.Logs {
		provider, err := e.newLogProvider(&l)
		if err != nil {
			e.LogPersister.AppendError(err.Error())
			continue
		}
		go e.runAnalysis(ctx, time.Duration(l.Interval), provider.Type(), func(ctx context.Context) (bool, error) {
			return provider.RunQuery(l.Query, l.FailureLimit)
		}, resultCh)
	}
	// TODO: Make HTTP analysis part of metrics provider.
	for range options.Https {

	}

	// Start watching the result of queries.
	var failureCount int
LOOP:
	for {
		select {
		case <-ctx.Done():
			break LOOP
		case res := <-resultCh:
			queryCount[res.provider]++
			e.saveQueryCount(ctx, queryCount)
			if !res.success {
				failureCount++
				e.Logger.Info("analysis run failed",
					zap.String("provider", res.provider),
					zap.Int("failure count", failureCount),
				)
			}
			if failureCount > options.FailureLimit {
				e.Logger.Info("stop all analysis")
				cancel()
			}
		}
	}
	if failureCount > options.FailureLimit {
		return model.StageStatus_STAGE_FAILURE
	}

	return model.StageStatus_STAGE_SUCCESS
}

// newMetricsProvider generates an appropriate metrics provider according to analysis metrics config.
func (e *Executor) newMetricsProvider(metrics *config.AnalysisMetrics) (metric.Provider, error) {
	// TODO: Address the case when using template
	providerCfg, ok := e.PipedConfig.GetProvider(metrics.Provider)
	if !ok {
		return nil, fmt.Errorf("unknown provider name %s", metrics.Provider)
	}

	var provider metric.Provider
	switch {
	case providerCfg.Prometheus != nil:
		cfg := providerCfg.Prometheus
		// TODO: Decide the way to authenticate.
		/*		username, err := ioutil.ReadFile(cfg.UsernameFile)
				if err != nil {
					return nil, err
				}
				password, err := ioutil.ReadFile(cfg.PasswordFile)
				if err != nil {
					return nil, err
				}
				provider, err = prometheus.NewProvider(cfg.Address, string(username), string(password))
		*/
		var err error
		provider, err = prometheus.NewProvider(cfg.Address, "", "", e.LogPersister)
		if err != nil {
			return nil, err
		}
	case providerCfg.Datadog != nil:
		cfg := providerCfg.Datadog
		apiKey, err := ioutil.ReadFile(cfg.APIKeyFile)
		if err != nil {
			return nil, err
		}
		applicationKey, err := ioutil.ReadFile(cfg.ApplicationKeyFile)
		if err != nil {
			return nil, err
		}
		provider, err = datadog.NewProvider(cfg.Address, string(apiKey), string(applicationKey))
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("provider config not found")
	}
	return provider, nil
}

// newLogProvider generates an appropriate log provider according to log metrics config.
func (e *Executor) newLogProvider(analysisLog *config.AnalysisLog) (log.Provider, error) {
	// TODO: Address the case when using template
	providerCfg, ok := e.PipedConfig.GetProvider(analysisLog.Provider)
	if !ok {
		return nil, fmt.Errorf("unknown provider name %s", analysisLog.Provider)
	}

	var provider log.Provider
	switch {
	case providerCfg.Stackdriver != nil:
		cfg := providerCfg.Stackdriver
		sa, err := ioutil.ReadFile(cfg.ServiceAccountFile)
		if err != nil {
			return nil, err
		}
		provider, err = stackdriver.NewProvider(sa)
		if err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("provider config not found")
	}
	return provider, nil
}

// runAnalysis calls `runQuery` function at the given interval and reports the result.
func (e *Executor) runAnalysis(ctx context.Context, interval time.Duration, providerType string, runQuery func(context.Context) (bool, error), resultCh chan<- providerResult) {
	e.Logger.Info("start the analysis", zap.String("provider", providerType))
	// TODO: Address the case when using template
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			success, err := runQuery(ctx)
			if err != nil {
				e.Logger.Error("failed to run query", zap.Error(err))
				// TODO: Decide how to handle query failures.
				success = false
			}
			resultCh <- providerResult{
				success:  success,
				provider: providerType,
			}

		case <-ctx.Done():
			return
		}
	}
}

func (e *Executor) getStageOptions() (*config.AnalysisStageOptions, error) {
	if e.Stage == nil {
		return nil, fmt.Errorf("stage information not found")
	}
	index := e.Stage.Index

	// TODO: Make `config.Deployment.GetPipelines`
	var stageConfig config.PipelineStage
	switch e.Deployment.Kind {
	case model.ApplicationKind_KUBERNETES:
		stages := e.DeploymentConfig.KubernetesDeploymentSpec.Pipeline.Stages
		if len(stages) < int(index)+1 {
			return nil, fmt.Errorf("unexpected stage index given")
		}
		stageConfig = stages[index]
	case model.ApplicationKind_TERRAFORM:
		stages := e.DeploymentConfig.TerraformDeploymentSpec.Pipeline.Stages
		if len(stages) < int(index)+1 {
			return nil, fmt.Errorf("unexpected stage index given")
		}
		stageConfig = stages[index]
	}
	if stageConfig.AnalysisStageOptions == nil {
		return nil, fmt.Errorf("no analysis options found")
	}

	return stageConfig.AnalysisStageOptions, nil
}

// saveQueryCount stores query count into metadata persister in json format.
// The analysis stage can be restarted from the beginning even if it ends unexpectedly,
// that's why count should be stored.
func (e *Executor) saveQueryCount(ctx context.Context, queryCount map[string]int) {
	if err := e.MetadataStore.SetStageMetadata(ctx, e.Stage.Id, queryCount); err != nil {
		e.Logger.Error("failed to store query count to stage metadata", zap.Error(err))
	}
}

// getQueryCount decodes metadata and populates query count to own field.
// The returned value is the number of queries executed per provider.
func (e *Executor) getQueryCount() map[string]int {
	var m map[string]int
	err := e.MetadataStore.GetStageMetadata(e.Stage.Id, &m)
	if err != nil {
		e.Logger.Error("failed to get stage metadata", zap.Error(err))
		return make(map[string]int)
	}
	return m
}
