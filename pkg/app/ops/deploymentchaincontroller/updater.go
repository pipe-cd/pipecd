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

package deploymentchaincontroller

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipe/pkg/datastore"
	"github.com/pipe-cd/pipe/pkg/model"
)

// updater watches for a specified deployment model object
// and updates the state of that object based on states of
// its deployments.
type updater struct {
	deploymentChainID string
	// applicationsRef contains list of all applications of
	// the current handling deployment chain.
	applicationsRef []*model.ChainApplicationRef
	// deploymentsRef is a map with key is the application id and
	// value is the last state of deployment ref to the deployment
	// of that in chain application.
	deploymentsRef map[string]*model.ChainDeploymentRef

	applicationStore     datastore.ApplicationStore
	deploymentStore      datastore.DeploymentStore
	deploymentChainStore datastore.DeploymentChainStore

	done          bool
	doneTimestamp time.Time

	logger  *zap.Logger
	nowFunc func() time.Time
}

func newUpdater(
	dc *model.DeploymentChain,
	as datastore.ApplicationStore,
	ds datastore.DeploymentStore,
	dcs datastore.DeploymentChainStore,
	lg *zap.Logger,
) *updater {
	return &updater{
		deploymentChainID:    dc.Id,
		applicationsRef:      dc.ListAllInChainApplications(),
		deploymentsRef:       dc.ListAllInChainApplicationDeploymentsMap(),
		applicationStore:     as,
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
		u.done = true
		u.doneTimestamp = u.nowFunc()
		return nil
	}

	// Check state of all deployment which existed from previous interval.
	if len(u.deploymentsRef) != 0 {
		for _, deploymentRef := range u.deploymentsRef {
			// Ignore finished deployment.
			if model.IsCompletedDeployment(deploymentRef.Status) {
				continue
			}
			deployment, err := u.deploymentStore.GetDeployment(ctx, deploymentRef.DeploymentId)
			if err != nil {
				return err
			}
			// Update the deployment state in deployment chain model in case
			// the deployment state is changed.
			if deployment.Status != deploymentRef.Status {
				if err = u.deploymentChainStore.UpdateDeploymentChain(ctx, u.deploymentChainID,
					datastore.DeploymentChainNodeDeploymentStatusUpdater(
						deployment.DeploymentChainBlockIndex,
						deployment.Id,
						deployment.Status,
						deployment.StatusReason),
				); err != nil {
					return err
				}
				// Update updater's deploymentsRef.
				deploymentRef.Status = deployment.Status
				deploymentRef.StatusReason = deployment.StatusReason
			}
		}
	}

	// If not all applications in chain has its deployment ref
	// fetch all missing deployments and link the deployment ref to them.
	if len(u.deploymentsRef) != len(u.applicationsRef) {
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
			if err := u.deploymentChainStore.UpdateDeploymentChain(
				ctx,
				u.deploymentChainID,
				datastore.DeploymentChainAddDeploymentToBlock(deployment),
			); err != nil {
				return err
			}
			// Store deployment ref to the updater's deploymentsRef.
			u.deploymentsRef[deployment.ApplicationId] = &model.ChainDeploymentRef{
				DeploymentId: deployment.Id,
				Status:       deployment.Status,
				StatusReason: deployment.StatusReason,
			}
		}
	}

	return nil
}

func (u *updater) listAllMissingDeployments(ctx context.Context) ([]*model.Deployment, error) {
	noDeploymentApps := make([]string, len(u.applicationsRef))
	for _, appRef := range u.applicationsRef {
		if _, ok := u.deploymentsRef[appRef.ApplicationId]; !ok {
			noDeploymentApps = append(noDeploymentApps, appRef.ApplicationId)
		}
	}

	deployments := make([]*model.Deployment, len(noDeploymentApps))
	// TODO: Find a better way to fetch applications in batch.
	for _, appId := range noDeploymentApps {
		app, err := u.applicationStore.GetApplication(ctx, appId)
		if err != nil {
			return nil, err
		}
		// If the most recent triggered deployment does not exist, ignore it.
		if app.MostRecentlyTriggeredDeployment == nil {
			continue
		}

		deployment, err := u.deploymentStore.GetDeployment(ctx, app.MostRecentlyTriggeredDeployment.DeploymentId)
		if err != nil {
			return nil, err
		}
		deployments = append(deployments, deployment)
	}
	return deployments, nil
}

func (u *updater) isAllInChainDeploymentsCompleted() bool {
	if len(u.applicationsRef) != len(u.deploymentsRef) {
		return false
	}
	allDeploymentCompleted := true
	for _, dr := range u.deploymentsRef {
		if !model.IsCompletedDeployment(dr.Status) {
			allDeploymentCompleted = false
			break
		}
	}
	return allDeploymentCompleted
}
