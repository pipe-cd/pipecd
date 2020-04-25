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

// Package appstatestore provides a runner component
// that watches the live state of applications in the cluster
// to construct it cache data that will be used to provide
// data to another components quickly.
package appstatestore

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	restclient "k8s.io/client-go/rest"

	"github.com/kapetaniosci/pipe/pkg/model"
)

type Store interface {
	Run(ctx context.Context) error
	WaitForReady(ctx context.Context, timeout time.Duration) error
	GetAppLiveResources(appID string) ([]model.K8SResource, error)
}

// appStateStore syncs the live state of application with the cluster
// and provides some functions for other components to query those states.
type appStateStore struct {
	kubeConfig    *restclient.Config
	store         *store
	firstSyncedCh chan error
	gracePeriod   time.Duration
	logger        *zap.Logger
}

func NewStore(kubeConfig *restclient.Config, gracePeriod time.Duration, logger *zap.Logger) Store {
	return &appStateStore{
		kubeConfig: kubeConfig,
		store: &store{
			apps:      make(map[string]*appLiveNodes),
			resources: make(map[string]appResource),
		},
		firstSyncedCh: make(chan error, 1),
		gracePeriod:   gracePeriod,
		logger:        logger,
	}
}

func (s *appStateStore) Run(ctx context.Context) error {
	stopCh := make(chan struct{})
	rf := reflector{
		kubeConfig: s.kubeConfig,
		onAdd:      s.store.onAddResource,
		onUpdate:   s.store.onUpdateResource,
		onDelete:   s.store.onDeleteResource,
		stopCh:     stopCh,
		logger:     s.logger.Named("reflector"),
	}
	if err := rf.start(ctx); err != nil {
		s.firstSyncedCh <- err
		return err
	}
	s.logger.Info("the reflector has done the first sync")
	s.store.initialize()
	s.logger.Info("the store has done the initializing")
	close(s.firstSyncedCh)

	s.logger.Info("DEBUG\n\n")
	apps := []string{
		"nginx-apps-v1",
		"nginx-apps-v1beta2",
		"nginx-extensions-v1beta1",
	}
	for _, app := range apps {
		s.logger.Info(fmt.Sprintf("Application: %s", app))

		managingNodes := s.store.getManagingNodesForApp(app)
		s.logger.Info(fmt.Sprintf("\tmanaging nodes %d", len(managingNodes)))
		for k, n := range managingNodes {
			s.logger.Info(fmt.Sprintf("\t\t%s: %s, %s", k, n.firstResourceKey, n.matchResourceKey))
		}

		dependedNodes := s.store.getDependedNodesForApp(app)
		s.logger.Info(fmt.Sprintf("\tdepended nodes %d", len(dependedNodes)))
		for k, n := range dependedNodes {
			s.logger.Info(fmt.Sprintf("\t\t%s: %s, %s", k, n.firstResourceKey, n.matchResourceKey))
		}
	}

	<-ctx.Done()
	close(stopCh)
	return nil
}

func (s *appStateStore) WaitForReady(ctx context.Context, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	select {
	case <-ctx.Done():
		return nil
	case err := <-s.firstSyncedCh:
		return err
	}
}

func (s *appStateStore) GetAppLiveResources(appID string) ([]model.K8SResource, error) {
	return nil, nil
}
