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
	"sort"
	"time"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipe/pkg/datastore"
	"github.com/pipe-cd/pipe/pkg/filestore"
	"github.com/pipe-cd/pipe/pkg/insightstore"
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
	logger *zap.Logger,
) *InsightCollector {
	a := &InsightCollector{
		applicationStore: datastore.NewApplicationStore(ds),
		deploymentStore:  datastore.NewDeploymentStore(ds),
		insightstore:     insightstore.NewStore(fs),
		logger:           logger.Named("insight-collector"),
	}
	return a
}

func (i *InsightCollector) Run(ctx context.Context) error {
	now := time.Now().UTC()
	maxUpdateAt := now.Unix()

	for {
		apps, err := i.applicationStore.ListApplications(ctx, datastore.ListOptions{
			PageSize: 50,
			Filters: []datastore.ListFilter{
				{
					Field:    "UpdatedAt",
					Operator: "<",
					Value:    maxUpdateAt,
				},
			},
			Orders: []datastore.Order{
				{
					Field:     "UpdatedAt",
					Direction: datastore.Desc,
				},
			},
		})
		if err != nil {
			return err
		}
		if len(apps) == 0 {
			// updated all application completely
			break
		}

		for _, app := range apps {
			for _, k := range aggregateKinds {
				// Update chunk insight file
				{
					chunkFiles, err := i.insightstore.LoadChunks(ctx, app.ProjectId, app.Id, k, model.InsightStep_MONTHLY, now, 1)
					if err != nil {
						return err
					}
					chunk := chunkFiles[0]
					accumulatedTo := time.Unix(chunk.GetAccumulatedTo(), 0).UTC()

					updatedps, accumulateTo, err := i.getDailyInsightData(ctx, app.Id, k, accumulatedTo, now)
					if err != nil {
						return err
					}

					for _, s := range model.InsightStep_value {
						step := model.InsightStep(s)
						chunk, err = i.updateChunk(chunk, step, k, updatedps, accumulateTo)
						if err != nil {
							return err
						}
					}
					err = i.insightstore.PutChunk(ctx, chunk)
					if err != nil {
						return err
					}
				}
				// Update years insight file
				{
					yearsFiles, err := i.insightstore.LoadChunks(ctx, app.ProjectId, app.Id, k, model.InsightStep_YEARLY, now, 1)
					if err != nil {
						return err
					}
					years := yearsFiles[0]
					accumulatedTo := time.Unix(yearsFiles[0].GetAccumulatedTo(), 0).UTC()

					updatedps, accumulateTo, err := i.getDailyInsightData(ctx, app.Id, k, accumulatedTo, now)
					if err != nil {
						return err
					}

					for _, s := range model.InsightStep_value {
						step := model.InsightStep(s)
						years, err = i.updateChunk(years, step, k, updatedps, accumulateTo)
						if err != nil {
							return err
						}
					}
					err = i.insightstore.PutChunk(ctx, years)
					if err != nil {
						return err
					}
				}
			}
		}
		maxUpdateAt = apps[len(apps)-1].UpdatedAt
	}
	return nil
}

func (i *InsightCollector) updateChunk(chunk insightstore.Chunk, step model.InsightStep, kind model.InsightMetricsKind, updatedps map[int64]insightstore.DataPoint, accumulatedTo int64) (insightstore.Chunk, error) {
	dps, err := chunk.GetDataPoints(step)
	if err != nil {
		return nil, err
	}

	for k, d := range updatedps {
		key := insightstore.NormalizeTime(time.Unix(k, 0).UTC(), step)
		dp, err := insightstore.GetDataPoint(dps, key.Unix())
		if err != nil && err != insightstore.ErrNotFound {
			return nil, err
		}
		new, err := insightstore.Merge(dp, d, kind)
		if err != nil {
			return nil, err
		}
		dps = insightstore.SetDataPoint(dps, new, k)
	}
	sort.SliceStable(dps, func(i, j int) bool { return dps[i].GetTimestamp() < dps[j].GetTimestamp() })

	chunk.SetAccumulatedTo(accumulatedTo)
	err = chunk.SetDataPoints(step, dps)
	if err != nil {
		return nil, err
	}

	return chunk, nil
}

func (i *InsightCollector) getDailyInsightData(
	ctx context.Context,
	appID string,
	kind model.InsightMetricsKind,
	rangeFrom time.Time,
	rangeTo time.Time,
) (map[int64]insightstore.DataPoint, int64, error) {
	step := model.InsightStep_DAILY

	var movePoint func(time.Time, int) time.Time
	movePoint = func(from time.Time, i int) time.Time {
		from = insightstore.NormalizeTime(from, step)
		return from.AddDate(0, 0, i)
	}

	updatedps := map[int64]insightstore.DataPoint{}

	to := movePoint(rangeFrom, 1)
	until := movePoint(rangeTo, 1)
	var accumulatedTo time.Time
	for {
		targetTimestamp := insightstore.NormalizeTime(rangeFrom, step).Unix()

		var data insightstore.DataPoint
		var a time.Time
		var err error
		switch kind {
		case model.InsightMetricsKind_DEPLOYMENT_FREQUENCY:
			data, a, err = i.getInsightDataForDeployFrequency(ctx, appID, targetTimestamp, rangeFrom, to)
		case model.InsightMetricsKind_CHANGE_FAILURE_RATE:
			data, a, err = i.getInsightDataForChangeFailureRate(ctx, appID, targetTimestamp, rangeFrom, to)
		default:
			return nil, 0, fmt.Errorf("invalid step: %v", kind)
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
			return nil, 0, err
		}

		updatedps[targetTimestamp] = data
		rangeFrom = a
		accumulatedTo = a
	}

	return updatedps, accumulatedTo.Unix(), nil
}

var (
	ErrDeploymentNotFound = errors.New("deployments not found")
)

// getInsightDataForDeployFrequency accumulate insight data in target range for deploy frequency.
func (i *InsightCollector) getInsightDataForDeployFrequency(
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
	deployments, err := i.deploymentStore.ListDeployments(ctx, datastore.ListOptions{
		PageSize: pageSize,
		Filters:  filters,
	})
	if err != nil {
		i.logger.Error("failed to get deployments", zap.Error(err))
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
func (i *InsightCollector) getInsightDataForChangeFailureRate(
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
	deployments, err := i.deploymentStore.ListDeployments(ctx, datastore.ListOptions{
		PageSize: pageSize,
		Filters:  filters,
	})
	if err != nil {
		i.logger.Error("failed to get deployments", zap.Error(err))
		return insightstore.ChangeFailureRate{}, time.Time{}, fmt.Errorf("failed to get deployments")
	}

	if len(deployments) == 0 {
		return insightstore.ChangeFailureRate{}, time.Time{}, ErrDeploymentNotFound
	}

	var successCount int64
	var failureCount int64
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
