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

package kubernetes

import (
	"context"
	"time"

	"go.uber.org/zap"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	// Import to load the needs plugins such as gcp, azure, oidc, openstack.
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"github.com/kapetaniosci/pipe/pkg/config"
	"github.com/kapetaniosci/pipe/pkg/model"
)

type Store struct {
	config        *config.CloudProviderKubernetesConfig
	kubeConfig    *restclient.Config
	store         *store
	firstSyncedCh chan error
	logger        *zap.Logger
}

type Getter interface {
	GetKubernetesAppLiveResources(appID string) []model.KubernetesResource
	NewEventIterator() EventIterator
}

type EventIterator struct {
	id    int
	store *store
}

func (it EventIterator) Next(maxNum int) []model.KubernetesResourceEvent {
	return it.store.nextEvents(it.id, maxNum)
}

func NewStore(cfg *config.CloudProviderKubernetesConfig, cloudProvider string, logger *zap.Logger) *Store {
	logger = logger.Named("kubernetes").
		With(zap.String("cloud-provider", cloudProvider))

	return &Store{
		config: cfg,
		store: &store{
			apps:      make(map[string]*appNodes),
			resources: make(map[string]appResource),
			iterators: make(map[int]int, 1),
		},
		firstSyncedCh: make(chan error, 1),
		logger:        logger,
	}
}

func (s *Store) Run(ctx context.Context) error {
	s.logger.Info("start running kubernetes app state store")

	// Build kubeconfig for initialing kubernetes clients later.
	var err error
	s.kubeConfig, err = clientcmd.BuildConfigFromFlags(s.config.MasterURL, s.config.KubeConfigPath)
	if err != nil {
		s.logger.Error("failed to build kube config", zap.Error(err))
		return err
	}

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

	<-ctx.Done()
	close(stopCh)

	s.logger.Info("kubernetes app state store has been stopped")
	return nil
}

func (s *Store) WaitForReady(ctx context.Context, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	select {
	case <-ctx.Done():
		return nil
	case err := <-s.firstSyncedCh:
		return err
	}
}

func (s *Store) GetKubernetesAppLiveResources(appID string) []model.KubernetesResource {
	nodes := s.store.getAppNodes(appID)
	resources := make([]model.KubernetesResource, 0, len(nodes))
	for _, n := range nodes {
		resources = append(resources, nodeToResource(n))
	}
	return resources
}

func (s *Store) NewEventIterator() EventIterator {
	return s.store.newEventIterator()
}
