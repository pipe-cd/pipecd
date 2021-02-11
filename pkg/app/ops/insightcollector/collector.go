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

	"github.com/pipe-cd/pipe/pkg/model"

	"github.com/pipe-cd/pipe/pkg/datastore"
	"github.com/pipe-cd/pipe/pkg/filestore"
	"github.com/pipe-cd/pipe/pkg/insight/insightstore"
	"go.uber.org/zap"
)

// InsightCollector implements the behaviors for the gRPC definitions of InsightCollector.
type InsightCollector struct {
	projectStore     datastore.ProjectStore
	applicationStore datastore.ApplicationStore
	deploymentStore  datastore.DeploymentStore
	insightstore     insightstore.Store
	collectFuncs     collectFunctions
	logger           *zap.Logger
}

type collectFunctions struct {
	applicationCollectFns []func(ctx context.Context, applications []*model.Application, target time.Time) error
	developmentCollectFns developmentCollectFunctions
}

type developmentCollectFunctions struct {
	newlyCreatedDeploymentsFns   []func(ctx context.Context, developments []*model.Deployment, target time.Time) error
	newlyCompletedDeploymentsFns []func(ctx context.Context, developments []*model.Deployment, target time.Time) error
}

// NewInsightCollector creates a new InsightCollector instance.
func NewInsightCollector(
	ds datastore.DataStore,
	fs filestore.Store,
	logger *zap.Logger,
	mode CollectorMode,
) *InsightCollector {
	i := &InsightCollector{
		projectStore:     datastore.NewProjectStore(ds),
		applicationStore: datastore.NewApplicationStore(ds),
		deploymentStore:  datastore.NewDeploymentStore(ds),
		insightstore:     insightstore.NewStore(fs),
		logger:           logger.Named("insight-collector"),
	}
	i.setCollectFunctions(mode)
	return i
}

func (i *InsightCollector) setCollectFunctions(mode CollectorMode) {
	cf := collectFunctions{}
	if mode.EnableApplicationCount() {
		cf.applicationCollectFns = append(cf.applicationCollectFns, i.collectApplicationCount)
	}
	if mode.EnableDevelopmentFrequency() {
		cf.developmentCollectFns.newlyCreatedDeploymentsFns = append(cf.developmentCollectFns.newlyCreatedDeploymentsFns, i.collectDevelopmentFrequency)
	}
	if mode.EnableChangeFailureRate() {
		cf.developmentCollectFns.newlyCompletedDeploymentsFns = append(cf.developmentCollectFns.newlyCompletedDeploymentsFns, i.collectChangeFailureRate)
	}
	i.collectFuncs = cf
}
