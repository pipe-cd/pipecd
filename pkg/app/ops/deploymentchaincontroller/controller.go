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
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/datastore"
	"github.com/pipe-cd/pipecd/pkg/model"
)

const (
	// syncDeploymentChainsInterval represents time to sync state for all controlling deployment chains.
	syncDeploymentChainsInterval = 15 * time.Second
	// syncUpdatersInterval represents time to update list of controlling deployment chain updaters.
	syncUpdatersInterval = time.Minute
	// updaterWorkerNum represents the maximum number of updaters which can be
	// run at the same time.
	maxUpdaterWorkerNum = 10
)

type deploymentStore interface {
	Get(ctx context.Context, id string) (*model.Deployment, error)
	List(ctx context.Context, opts datastore.ListOptions) ([]*model.Deployment, string, error)
}

type deploymentChainStore interface {
	List(ctx context.Context, opts datastore.ListOptions) ([]*model.DeploymentChain, string, error)
	AddNodeDeployment(ctx context.Context, chainID string, deployment *model.Deployment) error
	UpdateNodeDeploymentStatus(ctx context.Context, chainID string, blockIndex uint32, deploymentID string, status model.DeploymentStatus, reason string) error
}

type DeploymentChainController struct {
	deploymentStore      deploymentStore
	deploymentChainStore deploymentChainStore
	// Map from deployment chain ID to the updater
	// who in charge for the deployment chain updating.
	updaters map[string]*updater

	logger *zap.Logger
}

func NewDeploymentChainController(
	ds datastore.DataStore,
	logger *zap.Logger,
) *DeploymentChainController {
	w := datastore.OpsCommander
	return &DeploymentChainController{
		deploymentStore:      datastore.NewDeploymentStore(ds, w),
		deploymentChainStore: datastore.NewDeploymentChainStore(ds, w),
		updaters:             make(map[string]*updater),
		logger:               logger.Named("deployment-chain-controller"),
	}
}

func (d *DeploymentChainController) Run(ctx context.Context) error {
	syncUpdatersTicker := time.NewTicker(syncUpdatersInterval)
	defer syncUpdatersTicker.Stop()
	syncDeploymentChainsTicker := time.NewTicker(syncDeploymentChainsInterval)
	defer syncDeploymentChainsTicker.Stop()

	d.logger.Info("start running deployment chain controller")
	defer d.logger.Info("deployment chain controller has been stopped")

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-syncUpdatersTicker.C:
			if err := d.syncUpdaters(ctx); err != nil {
				d.logger.Error("failed while sync controller updaters", zap.Error(err))
			}
		case <-syncDeploymentChainsTicker.C:
			d.syncDeploymentChains(ctx)
		}
	}
}

func (d *DeploymentChainController) syncUpdaters(ctx context.Context) error {
	// Remove done updater of completed deployment chain.
	for id, u := range d.updaters {
		if !u.IsDone() {
			continue
		}
		d.logger.Info("remove done updater of deployment chain",
			zap.String("id", id),
			zap.Time("completed_at", u.doneTimestamp),
		)
		delete(d.updaters, id)
	}

	// Find all not completed deployment chains and create updater if does not exist.
	notCompletedChains, err := listNotCompletedDeploymentChain(ctx, d.deploymentChainStore)
	if err != nil {
		d.logger.Error("failed to fetch all not completed deployment chain", zap.Error(err))
		return err
	}
	for _, c := range notCompletedChains {
		// Ignore in case there is updater for that deployment chain existed.
		if _, ok := d.updaters[c.Id]; ok {
			continue
		}
		d.updaters[c.Id] = newUpdater(
			c,
			d.deploymentStore,
			d.deploymentChainStore,
			d.logger,
		)
	}

	return nil
}

func (d *DeploymentChainController) syncDeploymentChains(ctx context.Context) {
	var (
		updatersNum = len(d.updaters)
		updaterCh   = make(chan *updater, updatersNum)
		wg          sync.WaitGroup
	)
	updaterWorkerNum := maxUpdaterWorkerNum
	if updaterWorkerNum > updatersNum {
		updaterWorkerNum = updatersNum
	}

	d.logger.Info(fmt.Sprintf("there are %d running deployment chain updaters", updatersNum))
	for w := 0; w < updaterWorkerNum; w++ {
		wg.Add(1)
		go func(wid int) {
			d.logger.Info(fmt.Sprintf("worker id (%d) is handling deployment chain updaters", wid))
			defer wg.Done()

			for updater := range updaterCh {
				updater.Run(ctx)
			}
			d.logger.Info(fmt.Sprintf("worker id (%d) has stopped", wid))
		}(w)
	}

	for chainID := range d.updaters {
		updaterCh <- d.updaters[chainID]
	}
	close(updaterCh)

	d.logger.Info("waiting for all updaters to finish")
	wg.Wait()
}

func listNotCompletedDeploymentChain(ctx context.Context, dcs deploymentChainStore) ([]*model.DeploymentChain, error) {
	opts := datastore.ListOptions{
		Filters: []datastore.ListFilter{
			{
				Field:    "Status",
				Operator: datastore.OperatorIn,
				Value:    model.GetNotCompletedDeploymentChainStatuses(),
			},
		},
	}

	chains, _, err := dcs.List(ctx, opts)
	return chains, err
}
