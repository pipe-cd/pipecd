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
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipe/pkg/datastore"
	"github.com/pipe-cd/pipe/pkg/model"
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

type DeploymentChainController struct {
	applicationStore     datastore.ApplicationStore
	deploymentStore      datastore.DeploymentStore
	deploymentChainStore datastore.DeploymentChainStore
	// Map from deployment chain ID to the updater
	// who in charge for the deployment chain updating.
	updaters map[string]*updater
	// WaitGroup for waiting the fetcher and updaters to stop.
	wg sync.WaitGroup

	logger *zap.Logger
}

func NewDeploymentChainController(
	ds datastore.DataStore,
	logger *zap.Logger,
) *DeploymentChainController {
	return &DeploymentChainController{
		applicationStore:     datastore.NewApplicationStore(ds),
		deploymentStore:      datastore.NewDeploymentStore(ds),
		deploymentChainStore: datastore.NewDeploymentChainStore(ds),
		updaters:             make(map[string]*updater),
		logger:               logger.Named("deployment-chain-controller"),
	}
}

func (d *DeploymentChainController) Run(ctx context.Context) error {
	syncDeploymentChainsTicker := time.NewTicker(syncDeploymentChainsInterval)
	defer syncDeploymentChainsTicker.Stop()
	syncUpdatersTicker := time.NewTicker(syncUpdatersInterval)
	defer syncUpdatersTicker.Stop()

	d.logger.Info("start running deployment chain controller")

	d.wg.Add(1)
	go func() {
		defer d.wg.Done()
		for {
			select {
			case <-ctx.Done():
				return

			case <-syncUpdatersTicker.C:
				d.syncUpdaters(ctx)
			}
		}
	}()

	d.wg.Add(1)
	go func() {
		defer d.wg.Done()
		for {
			select {
			case <-ctx.Done():
				return

			case <-syncDeploymentChainsTicker.C:
				d.syncDeploymentChains(ctx)
			}
		}
	}()

	d.wg.Wait()
	d.logger.Info("deployment chain controller has been stopped")

	return nil
}

func (d *DeploymentChainController) syncUpdaters(ctx context.Context) error {
	// Remove done updater of completed deployment chain.
	for id, u := range d.updaters {
		if u.IsDone() {
			d.logger.Info("remove done updater of deployment chain",
				zap.String("id", id),
				zap.Time("completed_at", u.doneTimestamp),
			)
			delete(d.updaters, id)
		}
	}

	// Find all not completed deployment chains and create updater if does not exist.
	notCompletedChains, err := listNotCompletedDeploymentChain(ctx, d.deploymentChainStore)
	if err != nil {
		d.logger.Error("failed to fetch all not completed deployment chain", zap.Error(err))
		return err
	}
	for _, c := range notCompletedChains {
		if _, ok := d.updaters[c.Id]; !ok {
			d.updaters[c.Id] = newUpdater(
				c,
				d.applicationStore,
				d.deploymentStore,
				d.deploymentChainStore,
				d.logger,
			)
		}
	}

	return nil
}

func (d *DeploymentChainController) syncDeploymentChains(ctx context.Context) error {
	var (
		dcUpdatersNum = len(d.updaters)
		updatersCh    = make(chan *updater, dcUpdatersNum)
		resultCh      = make(chan error, dcUpdatersNum)
	)
	updaterWorkerNum := maxUpdaterWorkerNum
	if updaterWorkerNum > dcUpdatersNum {
		updaterWorkerNum = dcUpdatersNum
	}

	d.logger.Info(fmt.Sprintf("there are %d running deployment chain updaters", dcUpdatersNum))
	for w := 0; w < updaterWorkerNum; w++ {
		go func(wid int) {
			d.logger.Info(fmt.Sprintf("worker id (%d) is handling deployment chain updaters", wid))
			for updater := range updatersCh {
				resultCh <- updater.Run(ctx)
			}
			d.logger.Info(fmt.Sprintf("worker id (%d) has stopped", wid))
		}(w)
	}

	for chainID := range d.updaters {
		// Ignore updater which be marked as done.
		if d.updaters[chainID].IsDone() {
			continue
		}
		updatersCh <- d.updaters[chainID]
	}
	close(updatersCh)

	d.logger.Info("waiting for all updaters to finish")
	for i := 0; i < dcUpdatersNum; i++ {
		<-resultCh
	}

	return nil
}

func listNotCompletedDeploymentChain(ctx context.Context, dcs datastore.DeploymentChainStore) ([]*model.DeploymentChain, error) {
	opts := datastore.ListOptions{
		Filters: []datastore.ListFilter{
			{
				Field:    "Status",
				Operator: datastore.OperatorIn,
				Value:    model.GetNotCompletedDeploymentChainStatuses(),
			},
		},
	}

	chains, _, err := dcs.ListDeploymentChains(ctx, opts)
	return chains, err
}
