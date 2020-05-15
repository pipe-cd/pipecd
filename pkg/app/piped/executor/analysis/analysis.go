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
	"encoding/json"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/kapetaniosci/pipe/pkg/app/piped/executor"
	"github.com/kapetaniosci/pipe/pkg/app/piped/executor/analysis/metricsprovider"
	"github.com/kapetaniosci/pipe/pkg/app/piped/executor/analysis/metricsprovider/datadog"
	"github.com/kapetaniosci/pipe/pkg/app/piped/executor/analysis/metricsprovider/prometheus"
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

func (e *Executor) Execute(ctx context.Context) model.StageStatus {
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
	// TODO: Use provider's identifier defined in piped-config instead of AnalysisMetrics.Provider.
	//   Be sure to replace all `conf.Provider` with that.
	//   https://github.com/kapetaniosci/pipe/blob/master/pkg/config/testdata/piped/piped-config.yaml#L33
	for _, m := range options.Metrics {
		var provider metricsprovider.Provider
		var err error
		switch m.Provider {
		case "prometheus":
			provider, err = prometheus.NewProvider()
			if err != nil {
				return model.StageStatus_STAGE_FAILURE
			}
		case "datadog":
			provider, err = datadog.NewProvider()
			if err != nil {
				return model.StageStatus_STAGE_FAILURE
			}
		default:
			e.Logger.Error("unknown provider given")
			continue
		}
		go e.runMetricsQuery(ctx, &m, provider, resultCh)
	}
	// TODO: Support metrics provider for log and http.
	/*	for _, _ = range options.Logs {

		}
		for _, _ = range options.Https {

		}
	*/

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
			}
			if failureCount > options.Threshold {
				// Stop running all queries.
				cancel()
			}
		}
	}
	return model.StageStatus_STAGE_SUCCESS
}

func (e *Executor) runMetricsQuery(ctx context.Context, cfg *config.AnalysisMetrics, provider metricsprovider.Provider, resultCh chan<- providerResult) {
	ticker := time.NewTicker(time.Duration(cfg.Interval))
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			res, err := provider.RunQuery(cfg.Query)
			if err != nil {
				e.Logger.Error("failed to run query", zap.Error(err))
				// TODO: Decide how to handle query failures.
				continue
			}
			success := e.evaluate(cfg.Expected, res)
			resultCh <- providerResult{
				success:  success,
				provider: cfg.Provider,
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

// TODO: Implement to evaluate metrics provider result
func (e *Executor) evaluate(expected string, result float64) bool {
	return true
}

// saveQueryCount stores query count into metadata persister in json format.
// The analysis stage can be restarted from the beginning even if it ends unexpectedly,
// that's why count should be stored.
func (e *Executor) saveQueryCount(ctx context.Context, queryCount map[string]int) {
	b, err := json.Marshal(queryCount)
	if err != nil {
		e.Logger.Warn("failed to convert query cont to json")
		return
	}
	e.MetadataPersister.Save(ctx, b)
}

// getQueryCount decodes metadata and populates query count to own field.
// The returned value is the number of queries executed per provider.
func (e *Executor) getQueryCount() map[string]int {
	var m map[string]int
	err := json.Unmarshal(e.Stage.Metadata, &m)
	if err != nil {
		e.Logger.Warn("failed to decode query cont")
		return nil
	}
	return m
}
