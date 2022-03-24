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
	"time"

	"github.com/pipe-cd/pipecd/pkg/datastore"
	"github.com/pipe-cd/pipecd/pkg/model"
)

const limit = 50

var deploymentFrequencyMinimumVersion = model.InsightDeploymentVersion_V0

// func (c *Collector) collectDeploymentChangeFailureRate(ctx context.Context, ds []*model.Deployment, target time.Time) error {
// 	apps, projects := groupDeployments(ds)

// 	var updateErr error
// 	for id, ds := range apps {
// 		if err := c.updateApplicationChunks(ctx, ds[0].ProjectId, id, ds, model.InsightMetricsKind_CHANGE_FAILURE_RATE, target); err != nil {
// 			c.logger.Error("failed to update application chunks", zap.Error(err))
// 			updateErr = err
// 		}
// 	}
// 	for id, ds := range projects {
// 		if err := c.updateApplicationChunks(ctx, id, ds[0].ApplicationId, ds, model.InsightMetricsKind_CHANGE_FAILURE_RATE, target); err != nil {
// 			c.logger.Error("failed to update application chunks", zap.Error(err))
// 			updateErr = err
// 		}
// 	}

// 	return updateErr
// }

func (c *Collector) collectDevelopmentFrequency(ctx context.Context, ds []*model.Deployment, target time.Time) error {
	dailyDeployments := groupDeploymentByProjectID(ds)

	for projectID, deployments := range dailyDeployments {
		err := c.insightstore.Put(ctx, projectID, deployments, deploymentFrequencyMinimumVersion)
		if err != nil {
			return err
		}
	}

	return nil
}

// returned deployments are sorted by CompletedAt ASC
func (c *Collector) findDeploymentsCompletedInRange(ctx context.Context, from, to int64) ([]*model.Deployment, error) {
	if from > to {
		return []*model.Deployment{}, nil
	}

	filters := []datastore.ListFilter{
		{
			Field:    "CompletedAt",
			Operator: datastore.OperatorLessThanOrEqual,
			Value:    to,
		},
	}

	var deployments []*model.Deployment
	minCompletedAt := from
	for {
		d, _, err := c.deploymentStore.List(ctx, datastore.ListOptions{
			Limit: limit,
			Filters: append(filters, datastore.ListFilter{
				Field:    "CompletedAt",
				Operator: datastore.OperatorGreaterThanOrEqual,
				Value:    minCompletedAt,
			}),
			Orders: []datastore.Order{
				{
					Field:     "CompletedAt",
					Direction: datastore.Asc,
				},
				{
					Field:     "Id",
					Direction: datastore.Asc,
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
		// Add 1 for preventing infinity loop
		minCompletedAt = d[len(d)-1].CompletedAt + 1
	}
	return deployments, nil
}

func groupDeploymentByProjectID(deployments []*model.Deployment) map[string][]*model.InsightDeployment {
	groupById := make(map[string][]*model.InsightDeployment)
	for _, d := range deployments {
		var rollbackStartedAt int64
		if s, ok := d.FindRollbackStage(); ok {
			rollbackStartedAt = s.CreatedAt
		}
		groupById[d.ProjectId] = append(groupById[d.ProjectId], &model.InsightDeployment{
			Id:                d.Id,
			AppId:             d.ApplicationId,
			Labels:            d.Labels,
			StartedAt:         d.CreatedAt,
			CompletedAt:       d.CompletedAt,
			RollbackStartedAt: rollbackStartedAt,
			CompleteStatus:    d.Status,
		})
	}

	return groupById
}
