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

	"go.uber.org/zap"

	"github.com/pipe-cd/pipe/pkg/datastore"
	"github.com/pipe-cd/pipe/pkg/filestore"
	"github.com/pipe-cd/pipe/pkg/insight"
	"github.com/pipe-cd/pipe/pkg/model"
)

var aggregateKinds = []model.InsightMetricsKind{
	model.InsightMetricsKind_DEPLOYMENT_FREQUENCY,
	model.InsightMetricsKind_CHANGE_FAILURE_RATE,
}

// InsightCollector implements the behaviors for the gRPC definitions of InsightCollector.
type InsightCollector struct {
	projectStore     datastore.ProjectStore
	applicationStore datastore.ApplicationStore
	deploymentStore  datastore.DeploymentStore
	insightstore     insight.Store
	logger           *zap.Logger
}

// NewInsightCollector creates a new InsightCollector instance.
func NewInsightCollector(
	ds datastore.DataStore,
	fs filestore.Store,
	logger *zap.Logger,
) *InsightCollector {
	a := &InsightCollector{
		projectStore:     datastore.NewProjectStore(ds),
		applicationStore: datastore.NewApplicationStore(ds),
		deploymentStore:  datastore.NewDeploymentStore(ds),
		insightstore:     insight.NewStore(fs),
		logger:           logger.Named("insight-collector"),
	}
	return a
}

var (
	pageSize = 50
)

func (i *InsightCollector) CollectProjectsInsight(ctx context.Context) error {
	now := time.Now().UTC()
	maxCreatedAt := now.Unix()

	for {
		projects, err := i.projectStore.ListProjects(ctx, datastore.ListOptions{
			PageSize: pageSize,
			Filters: []datastore.ListFilter{
				{
					Field:    "CreatedAt",
					Operator: "<",
					Value:    maxCreatedAt,
				},
			},
			Orders: []datastore.Order{
				{
					Field:     "CreatedAt",
					Direction: datastore.Desc,
				},
			},
		})
		if err != nil {
			return err
		}
		if len(projects) == 0 {
			// updated all project's insights completely
			break
		}

		for _, p := range projects {
			for _, k := range aggregateKinds {
				if err := i.updateApplicationChunks(ctx, p.Id, "", k, now); err != nil {
					i.logger.Error("failed to update application chunks", zap.Error(err))
				}
			}
		}
		maxCreatedAt = projects[len(projects)-1].CreatedAt
	}
	return nil
}

func (i *InsightCollector) CollectApplicationInsight(ctx context.Context) error {
	now := time.Now().UTC()
	maxCreatedAt := now.Unix()

	for {
		apps, err := i.applicationStore.ListApplications(ctx, datastore.ListOptions{
			PageSize: pageSize,
			Filters: []datastore.ListFilter{
				{
					Field:    "CreatedAt",
					Operator: "<",
					Value:    maxCreatedAt,
				},
			},
			Orders: []datastore.Order{
				{
					Field:     "CreatedAt",
					Direction: datastore.Desc,
				},
			},
		})
		if err != nil {
			return err
		}
		if len(apps) == 0 {
			// updated all application's insights completely
			break
		}

		for _, app := range apps {
			if app.Deleted {
				continue
			}
			for _, k := range aggregateKinds {

				if err := i.updateApplicationChunks(ctx, app.ProjectId, app.Id, k, now); err != nil {
					i.logger.Error("failed to update application chunks", zap.Error(err))
				}
			}
		}
		maxCreatedAt = apps[len(apps)-1].CreatedAt
	}
	return nil
}

func (i *InsightCollector) updateApplicationChunks(
	ctx context.Context,
	projectID, appID string,
	kind model.InsightMetricsKind,
	to time.Time,
) error {
	chunkFiles, err := i.insightstore.LoadChunks(ctx, projectID, appID, kind, model.InsightStep_MONTHLY, to, 1)
	var chunk insight.Chunk
	if err == filestore.ErrNotFound {
		chunk = insight.NewChunk(projectID, kind, model.InsightStep_MONTHLY, appID, to)
	} else if err != nil {
		return err
	} else {
		chunk = chunkFiles[0]
	}

	yearsFiles, err := i.insightstore.LoadChunks(ctx, projectID, appID, kind, model.InsightStep_YEARLY, to, 1)
	var years insight.Chunk
	if err == filestore.ErrNotFound {
		years = insight.NewChunk(projectID, kind, model.InsightStep_YEARLY, appID, to)
	} else if err != nil {
		return err
	} else {
		years = yearsFiles[0]
	}

	chunk, years, err = i.updateChunk(ctx, chunk, years, projectID, appID, kind, to)
	if err != nil {
		return err
	}

	err = i.insightstore.PutChunk(ctx, chunk)
	if err != nil {
		return err
	}

	err = i.insightstore.PutChunk(ctx, years)
	if err != nil {
		return err
	}

	return nil
}

func (i *InsightCollector) updateChunk(
	ctx context.Context,
	chunk, years insight.Chunk,
	projectID, appID string,
	kind model.InsightMetricsKind,
	to time.Time,
) (insight.Chunk, insight.Chunk, error) {
	accumulatedTo := time.Unix(chunk.GetAccumulatedTo(), 0).UTC()
	yearsAccumulatedTo := time.Unix(years.GetAccumulatedTo(), 0).UTC()

	updatedps, accumulateTo, err := i.getDailyInsightData(ctx, projectID, appID, kind, accumulatedTo, to)
	if err != nil {
		return nil, nil, err
	}

	updatedpsForYears, yearAccumulateTo, err := i.getDailyInsightData(ctx, projectID, appID, kind, yearsAccumulatedTo, to)
	if err != nil {
		return nil, nil, err
	}

	for _, s := range model.InsightStep_value {
		step := model.InsightStep(s)
		if step == model.InsightStep_YEARLY {
			chunk, err = i.updateDataPoints(years, step, updatedpsForYears, yearAccumulateTo)
		} else {
			chunk, err = i.updateDataPoints(chunk, step, updatedps, accumulateTo)
		}
		if err != nil {
			return nil, nil, err
		}
	}
	return chunk, years, nil
}

func (i *InsightCollector) updateDataPoints(chunk insight.Chunk, step model.InsightStep, updatedps []insight.DataPoint, accumulatedTo int64) (insight.Chunk, error) {
	dps, err := chunk.GetDataPoints(step)
	if err != nil {
		return nil, err
	}

	for _, d := range updatedps {
		key := insight.NormalizeTime(time.Unix(d.GetTimestamp(), 0).UTC(), step)

		dps, err = insight.UpdateDataPoint(dps, d, key.Unix())
		if err != nil {
			return nil, err
		}
	}
	chunk.SetAccumulatedTo(accumulatedTo)
	err = chunk.SetDataPoints(step, dps)
	if err != nil {
		return nil, err
	}

	return chunk, nil
}

func (i *InsightCollector) getDailyInsightData(
	ctx context.Context,
	projectID, appID string,
	kind model.InsightMetricsKind,
	rangeFrom time.Time,
	rangeTo time.Time,
) ([]insight.DataPoint, int64, error) {
	step := model.InsightStep_DAILY

	var movePoint func(time.Time, int) time.Time
	movePoint = func(from time.Time, i int) time.Time {
		from = insight.NormalizeTime(from, step)
		return from.AddDate(0, 0, i)
	}

	var updatedps []insight.DataPoint

	to := movePoint(rangeFrom, 1)
	until := movePoint(rangeTo, 1)
	var accumulatedTo time.Time
	for {
		targetTimestamp := insight.NormalizeTime(rangeFrom, step).Unix()

		var data insight.DataPoint
		var a time.Time
		var err error
		switch kind {
		case model.InsightMetricsKind_DEPLOYMENT_FREQUENCY:
			data, a, err = i.getInsightDataForDeployFrequency(ctx, projectID, appID, targetTimestamp, rangeFrom, to)
		case model.InsightMetricsKind_CHANGE_FAILURE_RATE:
			data, a, err = i.getInsightDataForChangeFailureRate(ctx, projectID, appID, targetTimestamp, rangeFrom, to)
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

		updatedps = append(updatedps, data)
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
	projectID, applicationID string,
	targetTimestamp int64,
	from time.Time,
	to time.Time) (*insight.DeployFrequency, time.Time, error) {
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

	if projectID != "" {
		filters = append(filters, datastore.ListFilter{
			Field:    "ProjectId",
			Operator: "==",
			Value:    projectID,
		})
	}

	deployments, err := i.deploymentStore.ListDeployments(ctx, datastore.ListOptions{
		PageSize: pageSize,
		Filters:  filters,
	})
	if err != nil {
		i.logger.Error("failed to get deployments", zap.Error(err))
		return &insight.DeployFrequency{}, time.Time{}, fmt.Errorf("failed to get deployments")
	}
	if len(deployments) == 0 {
		return &insight.DeployFrequency{}, time.Time{}, ErrDeploymentNotFound
	}

	accumulatedTo := from.Unix()
	for _, d := range deployments {
		if d.CreatedAt > accumulatedTo {
			accumulatedTo = d.CreatedAt
		}
	}

	return &insight.DeployFrequency{
		Timestamp:   targetTimestamp,
		DeployCount: float32(len(deployments)),
	}, time.Unix(accumulatedTo, 0).UTC(), nil
}

// getInsightDataForChangeFailureRate accumulate insight data in target range for change failure rate
func (i *InsightCollector) getInsightDataForChangeFailureRate(
	ctx context.Context,
	projectID, applicationID string,
	targetTimestamp int64,
	from time.Time,
	to time.Time) (*insight.ChangeFailureRate, time.Time, error) {

	filters := []datastore.ListFilter{
		{
			Field:    "CompletedAt",
			Operator: ">=",
			Value:    from.Unix(),
		},
		{
			Field:    "CompletedAt",
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

	if projectID != "" {
		filters = append(filters, datastore.ListFilter{
			Field:    "ProjectId",
			Operator: "==",
			Value:    projectID,
		})
	}

	deployments, err := i.deploymentStore.ListDeployments(ctx, datastore.ListOptions{
		PageSize: pageSize,
		Filters:  filters,
	})
	if err != nil {
		i.logger.Error("failed to get deployments", zap.Error(err))
		return &insight.ChangeFailureRate{}, time.Time{}, fmt.Errorf("failed to get deployments")
	}

	if len(deployments) == 0 {
		return &insight.ChangeFailureRate{}, time.Time{}, ErrDeploymentNotFound
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
		if d.CompletedAt > accumulatedTo {
			accumulatedTo = d.CompletedAt
		}
	}

	return &insight.ChangeFailureRate{
		Timestamp:    targetTimestamp,
		Rate:         changeFailureRate,
		SuccessCount: successCount,
		FailureCount: failureCount,
	}, time.Unix(accumulatedTo, 0).UTC(), nil
}
