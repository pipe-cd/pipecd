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

package deploymentchaincontroller

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/datastore"
	"github.com/pipe-cd/pipecd/pkg/model"
)

// updater watches for a specified deployment model object
// and updates the state of that object based on states of
// its deployments.
type updater struct {
	deploymentChainID string
	// applicationRefs contains list of all applications of
	// the current handling deployment chain.
	applicationRefs []*model.ChainApplicationRef
	// deploymentRefs is a map with key is the application id and
	// value is the last state of deployment ref to the deployment
	// of that in chain application.
	deploymentRefs map[string]*model.ChainDeploymentRef

	deploymentStore      deploymentStore
	deploymentChainStore deploymentChainStore

	done          bool
	doneTimestamp time.Time

	logger  *zap.Logger
	nowFunc func() time.Time
}

func newUpdater(
	dc *model.DeploymentChain,
	ds deploymentStore,
	dcs deploymentChainStore,
	lg *zap.Logger,
) *updater {
	return &updater{
		deploymentChainID:    dc.Id,
		applicationRefs:      dc.ListAllInChainApplications(),
		deploymentRefs:       dc.ListAllInChainApplicationDeploymentsMap(),
		deploymentStore:      ds,
		deploymentChainStore: dcs,
		logger:               lg,
		nowFunc:              time.Now,
	}
}

func (u *updater) IsDone() bool {
	return u.done
}

func (u *updater) Run(ctx context.Context) error {
	// In case all in chain applications' deployments are completed
	// mark the updater as done and return immediately.
	if u.isAllInChainDeploymentsCompleted() {
		if !u.done {
			u.done = true
			u.doneTimestamp = u.nowFunc()
		}
		return nil
	}

	// Check state of all deployment which existed from previous interval.
	if len(u.deploymentRefs) != 0 {
		for _, deploymentRef := range u.deploymentRefs {
			// Ignore finished deployment.
			if deploymentRef.Status.IsCompleted() {
				continue
			}
			deployment, err := u.deploymentStore.Get(ctx, deploymentRef.DeploymentId)
			if err != nil {
				u.logger.Error("failed while update deployment chain: can not get deployment",
					zap.String("deploymentChainId", u.deploymentChainID),
					zap.String("deploymentId", deploymentRef.DeploymentId),
					zap.Error(err),
				)
				return err
			}
			// Update the deployment state in deployment chain model in case
			// the deployment state is changed.
			if deployment.Status != deploymentRef.Status {
				if err = u.deploymentChainStore.UpdateNodeDeploymentStatus(
					ctx,
					u.deploymentChainID,
					deployment.DeploymentChainBlockIndex,
					deployment.Id,
					deployment.Status,
					deployment.StatusReason,
				); err != nil {
					u.logger.Error("failed to update state of deployment in chain",
						zap.String("deploymentChainId", u.deploymentChainID),
						zap.String("deploymentId", deployment.Id),
						zap.Error(err),
					)
					return err
				}
				// Update updater's deploymentRefs.
				deploymentRef.Status = deployment.Status
				deploymentRef.StatusReason = deployment.StatusReason
			}
		}
	}

	// If not all applications in chain has its deployment ref
	// fetch all missing deployments and link the deployment ref to them.
	if len(u.deploymentRefs) != len(u.applicationRefs) {
		deployments, err := u.listAllMissingDeployments(ctx)
		if err != nil {
			return err
		}

		// Note: Only the deployment of the first block (the deployment which trigger this chain)
		// is added at the time we created this deployment chain; thus all other deployments
		// which be counted as "missng deployments" here are triggered and going to be added under
		// PENDING state, since those deployments are of blocks which have to wait for its previous
		// blocks completion. We don't need to care about updating those deployment ref status after
		// they been added here in this interval.
		for _, deployment := range deployments {
			// Connect missing deployment to the chain.
			if err := u.deploymentChainStore.AddNodeDeployment(ctx, u.deploymentChainID, deployment); err != nil {
				u.logger.Error("failed to link deployments to its chain",
					zap.String("deploymentChainId", u.deploymentChainID),
					zap.String("deploymentId", deployment.Id),
					zap.Error(err),
				)
				return err
			}
			// Store deployment ref to the updater's deploymentRefs.
			u.deploymentRefs[deployment.ApplicationId] = &model.ChainDeploymentRef{
				DeploymentId: deployment.Id,
				Status:       deployment.Status,
				StatusReason: deployment.StatusReason,
			}
		}
	}

	return nil
}

func (u *updater) listAllMissingDeployments(ctx context.Context) ([]*model.Deployment, error) {
	noDeploymentApps := make(map[string]interface{}, len(u.applicationRefs))
	for _, appRef := range u.applicationRefs {
		if _, ok := u.deploymentRefs[appRef.ApplicationId]; !ok {
			noDeploymentApps[appRef.ApplicationId] = nil
		}
	}

	// Fetch all available deployments in chain.
	options := datastore.ListOptions{
		Filters: []datastore.ListFilter{
			{
				Field:    "DeploymentChainId",
				Operator: datastore.OperatorEqual,
				Value:    u.deploymentChainID,
			},
		},
	}
	deployments, _, err := u.deploymentStore.List(ctx, options)
	if err != nil {
		u.logger.Error("failed to fetch all deployments in chain",
			zap.String("deploymentChainId", u.deploymentChainID),
			zap.Error(err),
		)
		return nil, err
	}

	missingDeployments := make([]*model.Deployment, 0, len(noDeploymentApps))
	for _, deployment := range deployments {
		if _, ok := noDeploymentApps[deployment.ApplicationId]; ok {
			missingDeployments = append(missingDeployments, deployment)
		}
	}
	return missingDeployments, nil
}

func (u *updater) isAllInChainDeploymentsCompleted() bool {
	if len(u.applicationRefs) != len(u.deploymentRefs) {
		return false
	}
	for _, dr := range u.deploymentRefs {
		if !dr.Status.IsCompleted() {
			return false
		}
	}
	return true
}
