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

package insightcollector

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/pipe-cd/pipe/pkg/filestore"

	"github.com/pipe-cd/pipe/pkg/insightstore"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipe/pkg/datastore"
	"github.com/pipe-cd/pipe/pkg/model"
)

var aggregateKinds = []model.InsightMetricsKind{
	model.InsightMetricsKind_DEPLOYMENT_FREQUENCY,
	model.InsightMetricsKind_CHANGE_FAILURE_RATE,
}

// InsightCollector implements the behaviors for the gRPC definitions of InsightCollector.
type InsightCollector struct {
	applicationStore datastore.ApplicationStore
	deploymentStore  datastore.DeploymentStore
	insightstore     insightstore.Store
	logger           *zap.Logger
}

// NewInsightCollector creates a new InsightCollector instance.
func NewInsightCollector(
	ds datastore.DataStore,
	fs filestore.Store,
	logger *zap.Logger) *InsightCollector {
	a := &InsightCollector{
		applicationStore: datastore.NewApplicationStore(ds),
		deploymentStore:  datastore.NewDeploymentStore(ds),
		insightstore:     insightstore.NewStore(fs),
		logger:           logger.Named("insight-collector"),
	}
	return a
}

func (a *InsightCollector) run(ctx context.Context) {
	now := time.Now().UTC()
	apps, err := a.applicationStore.ListApplications(ctx, datastore.ListOptions{
		Page:     0,
		PageSize: 0,
		Filters:  nil,
		Orders:   nil,
		Cursor:   "",
	})
	if err != nil {
		return
	}

	for _, app := range apps {
		for _, k := range aggregateKinds {
			yearsFiles, err := a.insightstore.LoadChunks(ctx, app.ProjectId, app.Id, k, model.InsightStep_YEARLY, now, 1)
			if err != nil {
				return
			}
			years := yearsFiles[0]
			yearsAccumulatedTo := time.Unix(years.GetAccumulatedTo(), 0).UTC()

			chunkFiles, err := a.insightstore.LoadChunks(ctx, app.ProjectId, app.Id, k, model.InsightStep_MONTHLY, now, 1)
			if err != nil {
				return
			}
			chunk := chunkFiles[0]
			chunkAccumulatedTo := time.Unix(chunk.GetAccumulatedTo(), 0).UTC()

			for _, s := range model.InsightStep_value {
				step := model.InsightStep(s)
				if step == model.InsightStep_YEARLY {
					a.getInsightData(ctx, app.Id, k, step, yearsAccumulatedTo, now, years)
				} else {
					a.getInsightData(ctx, app.Id, k, step, chunkAccumulatedTo, now, chunk)
				}
			}
		}

	}

}

func (a *InsightCollector) getInsightData(
	ctx context.Context,
	appID string,
	kind model.InsightMetricsKind,
	step model.InsightStep,
	rangeFrom time.Time,
	rangeTo time.Time,
	chunk insightstore.Chunk,
) (insightstore.Chunk, error) {
	dps, err := chunk.GetDataPoints(step)
	if err != nil {
		return nil, err
	}

	var movePoint func(time.Time, int) time.Time
	switch step {
	case model.InsightStep_DAILY:
		movePoint = func(from time.Time, i int) time.Time {
			from = insightstore.NormalizeTime(from, step)
			return from.AddDate(0, 0, i)
		}
	case model.InsightStep_WEEKLY:
		movePoint = func(from time.Time, i int) time.Time {
			from = insightstore.NormalizeTime(from, step)
			return from.AddDate(0, 0, i*7)
		}
	case model.InsightStep_MONTHLY:
		movePoint = func(from time.Time, i int) time.Time {
			from = insightstore.NormalizeTime(from, step)
			return from.AddDate(0, i, 0)
		}
	case model.InsightStep_YEARLY:
		movePoint = func(from time.Time, i int) time.Time {
			from = insightstore.NormalizeTime(from, step)
			return from.AddDate(i, 0, 0)
		}
	default:
		return nil, fmt.Errorf("invalid step: %v", step)
	}

	to := movePoint(rangeFrom, 1)
	until := movePoint(rangeTo, 1)
	for {
		targetTimestamp := insightstore.NormalizeTime(rangeFrom, step).Unix()

		var data insightstore.DataPoint
		var accumulatedTo time.Time
		var err error
		switch kind {
		case model.InsightMetricsKind_DEPLOYMENT_FREQUENCY:
			data, accumulatedTo, err = a.getInsightDataForDeployFrequency(ctx, appID, targetTimestamp, rangeFrom, to)
		case model.InsightMetricsKind_CHANGE_FAILURE_RATE:
			data, accumulatedTo, err = a.getInsightDataForChangeFailureRate(ctx, appID, targetTimestamp, rangeFrom, to)
		default:
			return nil, fmt.Errorf("invalid step: %v", kind)
		}
		if err != nil {
			if err == ErrDeploymentNotFound {
				if to.Equal(until) {
					break
				}
				rangeFrom = to
				to = movePoint(to, 1)
				continue
			}
			return nil, err
		}

		// update data point and set it into chunk
		dp, err := insightstore.GetDataPoint(dps, targetTimestamp)
		if err != nil && err != insightstore.ErrNotFound {
			return nil, err
		}
		new, err := insightstore.Merge(dp, data, kind)
		if err != nil {
			return nil, err
		}
		dps = insightstore.SetDataPoint(dps, new, targetTimestamp)

		chunk.SetAccumulatedTo(accumulatedTo.Unix())
		rangeFrom = accumulatedTo
	}
	err = chunk.SetDataPoints(step, dps)
	if err != nil {
		return nil, err
	}
	return chunk, nil
}

var (
	ErrDeploymentNotFound = errors.New("deployments not found")
)

// getInsightDataForDeployFrequency accumulate insight data in target range for deploy frequency.
func (a *InsightCollector) getInsightDataForDeployFrequency(
	ctx context.Context,
	applicationID string,
	targetTimestamp int64,
	from time.Time,
	to time.Time) (insightstore.DeployFrequency, time.Time, error) {
	filters := []datastore.ListFilter{
		{
			Field:    "CreatedAt",
			Operator: ">=",
			Value:    from.Unix(),
		},
		{
			Field:    "CreatedAt",
			Operator: "<",
			Value:    to.Unix(),
		},
	}

	if applicationID != "" {
		filters = append(filters, datastore.ListFilter{
			Field:    "ApplicationId",
			Operator: "==",
			Value:    applicationID,
		})
	}

	pageSize := 50
	deployments, err := a.deploymentStore.ListDeployments(ctx, datastore.ListOptions{
		PageSize: pageSize,
		Filters:  filters,
	})
	if err != nil {
		a.logger.Error("failed to get deployments", zap.Error(err))
		return insightstore.DeployFrequency{}, time.Time{}, fmt.Errorf("failed to get deployments")
	}
	if len(deployments) == 0 {
		return insightstore.DeployFrequency{}, time.Time{}, ErrDeploymentNotFound
	}

	accumulatedTo := from.Unix()
	for _, d := range deployments {
		if d.CreatedAt > accumulatedTo {
			accumulatedTo = d.CreatedAt
		}
	}

	return insightstore.DeployFrequency{
		Timestamp:   targetTimestamp,
		DeployCount: float32(len(deployments)),
	}, time.Unix(accumulatedTo, 0).UTC(), nil
}

// getInsightDataForChangeFailureRate accumulate insight data in target range for change failure rate
func (a *InsightCollector) getInsightDataForChangeFailureRate(
	ctx context.Context,
	applicationID string,
	targetTimestamp int64,
	from time.Time,
	to time.Time) (insightstore.ChangeFailureRate, time.Time, error) {

	filters := []datastore.ListFilter{
		{
			Field:    "CreatedAt",
			Operator: ">=",
			Value:    from.Unix(),
		},
		{
			Field:    "CreatedAt",
			Operator: "<",
			Value:    to.Unix(),
		},
		{
			Field:    "Status",
			Operator: "in",
			Value:    []model.DeploymentStatus{model.DeploymentStatus_DEPLOYMENT_FAILURE, model.DeploymentStatus_DEPLOYMENT_SUCCESS},
		},
	}

	if applicationID != "" {
		filters = append(filters, datastore.ListFilter{
			Field:    "ApplicationId",
			Operator: "==",
			Value:    applicationID,
		})
	}

	pageSize := 50
	deployments, err := a.deploymentStore.ListDeployments(ctx, datastore.ListOptions{
		PageSize: pageSize,
		Filters:  filters,
	})
	if err != nil {
		a.logger.Error("failed to get deployments", zap.Error(err))
		return insightstore.ChangeFailureRate{}, time.Time{}, fmt.Errorf("failed to get deployments")
	}

	if len(deployments) == 0 {
		return insightstore.ChangeFailureRate{}, time.Time{}, ErrDeploymentNotFound
	}

	var successCount int64 = 0
	var failureCount int64 = 0
	for _, d := range deployments {
		switch d.Status {
		case model.DeploymentStatus_DEPLOYMENT_SUCCESS:
			successCount++
		case model.DeploymentStatus_DEPLOYMENT_FAILURE:
			failureCount++
		}
	}

	var changeFailureRate float32
	if successCount+failureCount != 0 {
		changeFailureRate = float32(failureCount) / float32(successCount+failureCount)
	} else {
		changeFailureRate = 0
	}

	accumulatedTo := from.Unix()
	for _, d := range deployments {
		if d.CreatedAt > accumulatedTo {
			accumulatedTo = d.CreatedAt
		}
	}

	return insightstore.ChangeFailureRate{
		Timestamp:    targetTimestamp,
		Rate:         changeFailureRate,
		SuccessCount: successCount,
		FailureCount: failureCount,
	}, time.Unix(accumulatedTo, 0).UTC(), nil
}
