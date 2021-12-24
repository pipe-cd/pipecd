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
	interval = 1 * time.Minute
)

type DeploymentChainController struct {
	applicationStore     datastore.ApplicationStore
	deploymentStore      datastore.DeploymentStore
	deploymentChainStore datastore.DeploymentChainStore
	// Map from deployment chain ID to the updater
	// who in charge for the deployment chain updating.
	updaters map[string]*updater
	// WaitGroup for waiting all deployment chain updaters to be completed.
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
	d.logger.Info("start running DeploymentChainController")

	t := time.NewTicker(interval)
	defer t.Stop()
	d.logger.Info("start syncing updaters")

L:
	for {
		select {
		case <-ctx.Done():
			break L

		case <-t.C:
			d.syncUpdaters(ctx)
		}
	}

	d.logger.Info("deploymentChainController has been stopped")
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

	d.logger.Info(fmt.Sprintf("there are %d running deployment chain updaters", len(d.updaters)),
		zap.Int("count", len(d.updaters)),
	)

	for chainID := range d.updaters {
		updater := d.updaters[chainID]
		d.wg.Add(1)
		go func() {
			defer d.wg.Done()
			if err := updater.Run(ctx); err != nil {
				d.logger.Error("failed to update deployment chain",
					zap.String("deploymentChainId", updater.deploymentChainID),
					zap.Error(err),
				)
			}
		}()
	}

	d.logger.Info("waiting for all updaters to finish")
	d.wg.Wait()

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
