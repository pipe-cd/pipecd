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
	"time"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipe/pkg/datastore"
	"github.com/pipe-cd/pipe/pkg/model"
)

const (
	// interval represents time to sync state for all deployment chain.
	interval = time.Minute
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
	t := time.NewTicker(interval)
	defer t.Stop()

	d.logger.Info("start running deployment chain controller")
	defer d.logger.Info("deployment chain controller has been stopped")
	for {
		select {
		case <-ctx.Done():
			return nil

		case <-t.C:
			d.syncUpdaters(ctx)
		}
	}
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

	var (
		dcUpdatersCtn = len(d.updaters)
		updatersCh    = make(chan *updater, dcUpdatersCtn)
		resultCh      = make(chan error, dcUpdatersCtn)
	)
	updaterWorkerNum := maxUpdaterWorkerNum
	if updaterWorkerNum > dcUpdatersCtn {
		updaterWorkerNum = dcUpdatersCtn
	}

	d.logger.Info(fmt.Sprintf("there are %d running deployment chain updaters", dcUpdatersCtn))
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
		updatersCh <- d.updaters[chainID]
	}
	close(updatersCh)

	d.logger.Info("waiting for all updaters to finish")
	for i := 0; i < dcUpdatersCtn; i++ {
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
