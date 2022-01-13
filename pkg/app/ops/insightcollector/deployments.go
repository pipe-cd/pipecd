// Copyright 2021 The PipeCD Authors.
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
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/datastore"
	"github.com/pipe-cd/pipecd/pkg/filestore"
	"github.com/pipe-cd/pipecd/pkg/insight"
	"github.com/pipe-cd/pipecd/pkg/model"
)

const limit = 50

func (c *Collector) collectDeploymentChangeFailureRate(ctx context.Context, ds []*model.Deployment, target time.Time) error {
	apps, projects := groupDeployments(ds)

	var updateErr error
	for id, ds := range apps {
		if err := c.updateApplicationChunks(ctx, ds[0].ProjectId, id, ds, model.InsightMetricsKind_CHANGE_FAILURE_RATE, target); err != nil {
			c.logger.Error("failed to update application chunks", zap.Error(err))
			updateErr = err
		}
	}
	for id, ds := range projects {
		if err := c.updateApplicationChunks(ctx, id, ds[0].ApplicationId, ds, model.InsightMetricsKind_CHANGE_FAILURE_RATE, target); err != nil {
			c.logger.Error("failed to update application chunks", zap.Error(err))
			updateErr = err
		}
	}

	return updateErr
}

func (c *Collector) collectDevelopmentFrequency(ctx context.Context, ds []*model.Deployment, target time.Time) error {
	apps, projects := groupDeployments(ds)

	var updateErr error
	for id, ds := range apps {
		if err := c.updateApplicationChunks(ctx, ds[0].ProjectId, id, ds, model.InsightMetricsKind_DEPLOYMENT_FREQUENCY, target); err != nil {
			c.logger.Error("failed to update application chunks", zap.Error(err))
			updateErr = err
		}
	}
	for id, ds := range projects {
		if err := c.updateApplicationChunks(ctx, id, ds[0].ApplicationId, ds, model.InsightMetricsKind_DEPLOYMENT_FREQUENCY, target); err != nil {
			c.logger.Error("failed to update project chunks", zap.Error(err))
			updateErr = err
		}
	}

	return updateErr
}

func (c *Collector) findDeploymentsCreatedInRange(ctx context.Context, from, to int64) ([]*model.Deployment, error) {
	filters := []datastore.ListFilter{
		{
			Field:    "CreatedAt",
			Operator: datastore.OperatorGreaterThanOrEqual,
			Value:    from,
		},
	}

	var deployments []*model.Deployment
	maxCreatedAt := to
	for {
		d, _, err := c.deploymentStore.ListDeployments(ctx, datastore.ListOptions{
			Limit: limit,
			Filters: append(filters, datastore.ListFilter{
				Field:    "CreatedAt",
				Operator: datastore.OperatorLessThan,
				Value:    maxCreatedAt,
			}),
			Orders: []datastore.Order{
				{
					Field:     "CreatedAt",
					Direction: datastore.Desc,
				},
			},
		})
		if err != nil {
			return nil, err
		}
		if len(d) == 0 {
			// get all deployments in range
			break
		}

		deployments = append(deployments, d...)
		maxCreatedAt = d[len(d)-1].CreatedAt
	}
	return deployments, nil
}

func (c *Collector) findDeploymentsCompletedInRange(ctx context.Context, from, to int64) ([]*model.Deployment, error) {
	filters := []datastore.ListFilter{
		{
			Field:    "CompletedAt",
			Operator: datastore.OperatorGreaterThanOrEqual,
			Value:    from,
		},
	}

	var deployments []*model.Deployment
	maxCompletedAt := to
	for {
		d, _, err := c.deploymentStore.ListDeployments(ctx, datastore.ListOptions{
			Limit: limit,
			Filters: append(filters, datastore.ListFilter{
				Field:    "CompletedAt",
				Operator: datastore.OperatorLessThan,
				Value:    maxCompletedAt,
			}),
			Orders: []datastore.Order{
				{
					Field:     "CompletedAt",
					Direction: datastore.Desc,
				},
			},
		})
		if err != nil {
			return nil, err
		}
		if len(d) == 0 {
			// get all deployments in range
			break
		}

		deployments = append(deployments, d...)
		maxCompletedAt = d[len(d)-1].CompletedAt
	}
	return deployments, nil
}

// updateApplicationChunks updates chunk in filestore
func (c *Collector) updateApplicationChunks(ctx context.Context, projectID, appID string, deployments []*model.Deployment, kind model.InsightMetricsKind, targetDate time.Time) error {
	chunkFiles, err := c.insightstore.LoadChunks(ctx, projectID, appID, kind, model.InsightStep_MONTHLY, targetDate, 1)
	var chunk insight.Chunk
	if err == filestore.ErrNotFound {
		chunk = insight.NewChunk(projectID, kind, model.InsightStep_MONTHLY, appID, targetDate)
	} else if err != nil {
		return err
	} else {
		chunk = chunkFiles[0]
	}

	yearsFiles, err := c.insightstore.LoadChunks(ctx, projectID, appID, kind, model.InsightStep_YEARLY, targetDate, 1)
	var years insight.Chunk
	if err == filestore.ErrNotFound {
		years = insight.NewChunk(projectID, kind, model.InsightStep_YEARLY, appID, targetDate)
	} else if err != nil {
		return err
	} else {
		years = yearsFiles[0]
	}

	chunk, years, err = updateChunk(deployments, chunk, years, kind, targetDate)
	if err != nil {
		return err
	}

	err = c.insightstore.PutChunk(ctx, chunk)
	if err != nil {
		return err
	}

	err = c.insightstore.PutChunk(ctx, years)
	if err != nil {
		return err
	}

	return nil
}

// updateChunk updates passed chunk with deployments
func updateChunk(deployments []*model.Deployment, chunk, years insight.Chunk, kind model.InsightMetricsKind, targetDate time.Time) (insight.Chunk, insight.Chunk, error) {
	accumulatedTo := time.Unix(chunk.GetAccumulatedTo(), 0).UTC()
	yearsAccumulatedTo := time.Unix(years.GetAccumulatedTo(), 0).UTC()

	if accumulatedTo != targetDate {
		updatedps, err := extractDailyInsightDataPoints(deployments, kind, accumulatedTo, targetDate)
		if err != nil {
			return nil, nil, err
		}
		for _, s := range model.InsightStep_value {
			step := model.InsightStep(s)
			if step != model.InsightStep_YEARLY {
				chunk, err = updateDataPoints(years, step, updatedps, targetDate.Unix())
				if err != nil {
					return nil, nil, err
				}
			}
		}
	}

	if yearsAccumulatedTo != targetDate {
		updatedpsForYears, err := extractDailyInsightDataPoints(deployments, kind, yearsAccumulatedTo, targetDate)
		if err != nil {
			return nil, nil, err
		}

		chunk, err = updateDataPoints(chunk, model.InsightStep_YEARLY, updatedpsForYears, targetDate.Unix())
		if err != nil {
			return nil, nil, err
		}
	}

	return chunk, years, nil
}

// updateDataPoints updates chunk's datapoints with accumuleatedTo and datapoints for update
func updateDataPoints(chunk insight.Chunk, step model.InsightStep, updatedps []insight.DataPoint, accumulatedTo int64) (insight.Chunk, error) {
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

// extractDailyInsightDataPoints extracts the daily datapoints from deployment
func extractDailyInsightDataPoints(deployments []*model.Deployment, kind model.InsightMetricsKind, rangeFrom, rangeTo time.Time) ([]insight.DataPoint, error) {
	step := model.InsightStep_DAILY

	var movePoint func(time.Time, int) time.Time
	movePoint = func(from time.Time, i int) time.Time {
		from = insight.NormalizeTime(from, step)
		return from.AddDate(0, 0, i)
	}

	var updatedps []insight.DataPoint

	to := movePoint(rangeFrom, 1)
	for !to.After(rangeTo) {
		targetTimestamp := insight.NormalizeTime(rangeFrom, step).Unix()

		var data insight.DataPoint
		switch kind {
		case model.InsightMetricsKind_DEPLOYMENT_FREQUENCY:
			data, deployments = extractDeployFrequency(deployments, rangeFrom.Unix(), to.Unix(), targetTimestamp)
		case model.InsightMetricsKind_CHANGE_FAILURE_RATE:
			data, deployments = extractChangeFailureRate(deployments, rangeFrom.Unix(), to.Unix(), targetTimestamp)
		default:
			return nil, fmt.Errorf("invalid step: %v", kind)
		}

		updatedps = append(updatedps, data)
		rangeFrom = movePoint(rangeFrom, 1)
		to = movePoint(to, 1)
	}

	return updatedps, nil
}

// extractDeployFrequency extracts deploy frequency from deployments with specified range
func extractDeployFrequency(deployments []*model.Deployment, from, to int64, targetTimestamp int64) (*insight.DeployFrequency, []*model.Deployment) {
	var ds []*model.Deployment
	var rest []*model.Deployment
	for _, d := range deployments {
		if d.CreatedAt < to && d.CreatedAt >= from {
			ds = append(ds, d)
		} else {
			rest = append(rest, d)
		}
	}

	return &insight.DeployFrequency{
		Timestamp:   targetTimestamp,
		DeployCount: float32(len(ds)),
	}, rest
}

// extractChangeFailureRate extracts change failure rate from deployments with specified range
func extractChangeFailureRate(deployments []*model.Deployment, from, to int64, targetTimestamp int64) (*insight.ChangeFailureRate, []*model.Deployment) {
	var ds []*model.Deployment
	var rest []*model.Deployment
	for _, d := range deployments {
		if d.CompletedAt < to && d.CompletedAt >= from {
			ds = append(ds, d)
		} else {
			rest = append(rest, d)
		}
	}
	var successCount int64
	var failureCount int64
	for _, d := range ds {
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

	return &insight.ChangeFailureRate{
		Timestamp:    targetTimestamp,
		Rate:         changeFailureRate,
		SuccessCount: successCount,
		FailureCount: failureCount,
	}, rest
}

// groupDeployments groups deployments by applicationID and projectID
func groupDeployments(deployments []*model.Deployment) (apps, projects map[string][]*model.Deployment) {
	apps = make(map[string][]*model.Deployment)
	projects = make(map[string][]*model.Deployment)

	for _, d := range deployments {
		apps[d.ApplicationId] = append(apps[d.ApplicationId], d)
		projects[d.ProjectId] = append(projects[d.ProjectId], d)
	}
	return
}
