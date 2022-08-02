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

// Package livestatestore provides a piped component
// that watches the live state of applications in the cluster
// to construct it cache data that will be used to provide
// data to another components quickly.
package livestatestore

import (
	"context"
	"time"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/pipe-cd/pipecd/pkg/app/piped/livestatestore/cloudrun"
	"github.com/pipe-cd/pipecd/pkg/app/piped/livestatestore/ecs"
	"github.com/pipe-cd/pipecd/pkg/app/piped/livestatestore/kubernetes"
	"github.com/pipe-cd/pipecd/pkg/app/piped/livestatestore/lambda"
	"github.com/pipe-cd/pipecd/pkg/app/piped/livestatestore/terraform"
	provider "github.com/pipe-cd/pipecd/pkg/app/piped/platformprovider/kubernetes"
	"github.com/pipe-cd/pipecd/pkg/config"
	"github.com/pipe-cd/pipecd/pkg/model"
)

type applicationLister interface {
	List() []*model.Application
}

type Getter interface {
	CloudRunGetter(cloudProvider string) (cloudrun.Getter, bool)
	ECSRunGetter(cloudProvider string) (ecs.Getter, bool)
	KubernetesGetter(cloudProvider string) (kubernetes.Getter, bool)
	LambdaGetter(cloudProvider string) (lambda.Getter, bool)
	TerraformGetter(cloudProvider string) (terraform.Getter, bool)
}

type Store interface {
	Run(ctx context.Context) error
	Getter() Getter
}

type kubernetesStore interface {
	Run(ctx context.Context) error
	kubernetes.Getter
}

type terraformStore interface {
	Run(ctx context.Context) error
}

type cloudRunStore interface {
	Run(ctx context.Context) error
	cloudrun.Getter
}

type lambdaStore interface {
	Run(ctx context.Context) error
}

type ecsStore interface {
	Run(ctx context.Context) error
}

// store manages a list of particular stores for all cloud providers.
type store struct {
	// Map thats contains a list of kubernetesStore where key is the cloud provider name.
	kubernetesStores map[string]kubernetesStore
	// Map thats contains a list of terraformStore where key is the cloud provider name.
	terraformStores map[string]terraformStore
	// Map thats contains a list of cloudRunStore where key is the cloud provider name.
	cloudrunStores map[string]cloudRunStore
	// Map thats contains a list of lambdaStore where key is the cloud provider name.
	lambdaStores map[string]lambdaStore
	// Map thats contains a list of ecsStore where key is the cloud provider name.
	ecsStores map[string]ecsStore

	gracePeriod time.Duration
	logger      *zap.Logger
}

func NewStore(ctx context.Context, cfg *config.PipedSpec, appLister applicationLister, gracePeriod time.Duration, logger *zap.Logger) Store {
	logger = logger.Named("livestatestore")

	s := &store{
		kubernetesStores: make(map[string]kubernetesStore),
		terraformStores:  make(map[string]terraformStore),
		cloudrunStores:   make(map[string]cloudRunStore),
		lambdaStores:     make(map[string]lambdaStore),
		ecsStores:        make(map[string]ecsStore),
		gracePeriod:      gracePeriod,
		logger:           logger,
	}
	for _, cp := range cfg.PlatformProviders {
		switch cp.Type {
		case model.PlatformProviderKubernetes:
			store := kubernetes.NewStore(cp.KubernetesConfig, cfg, cp.Name, logger)
			s.kubernetesStores[cp.Name] = store

		case model.PlatformProviderTerraform:
			store := terraform.NewStore(cp.TerraformConfig, cp.Name, appLister, logger)
			s.terraformStores[cp.Name] = store

		case model.PlatformProviderCloudRun:
			store, err := cloudrun.NewStore(ctx, cp.CloudRunConfig, cp.Name, logger)
			if err != nil {
				logger.Error("failed to create a new cloudrun's livestatestore", zap.Error(err))
				continue
			}
			s.cloudrunStores[cp.Name] = store

		case model.PlatformProviderLambda:
			store := lambda.NewStore(cp.LambdaConfig, cp.Name, appLister, logger)
			s.lambdaStores[cp.Name] = store

		case model.PlatformProviderECS:
			store := ecs.NewStore(cp.ECSConfig, cp.Name, appLister, logger)
			s.ecsStores[cp.Name] = store
		}
	}

	return s
}

func (s *store) Run(ctx context.Context) error {
	s.logger.Info("start running appsatestore")

	group, ctx := errgroup.WithContext(ctx)

	for i := range s.kubernetesStores {
		cpName := i
		group.Go(func() error {
			return s.kubernetesStores[cpName].Run(ctx)
		})
	}

	for i := range s.terraformStores {
		cpName := i
		group.Go(func() error {
			return s.terraformStores[cpName].Run(ctx)
		})
	}

	for i := range s.cloudrunStores {
		cpName := i
		group.Go(func() error {
			return s.cloudrunStores[cpName].Run(ctx)
		})
	}

	for i := range s.lambdaStores {
		cpName := i
		group.Go(func() error {
			return s.lambdaStores[cpName].Run(ctx)
		})
	}

	for i := range s.ecsStores {
		cpName := i
		group.Go(func() error {
			return s.ecsStores[cpName].Run(ctx)
		})
	}

	err := group.Wait()
	if err == nil {
		s.logger.Info("all state stores have been stopped")
	} else {
		s.logger.Error("all state stores have been stopped", zap.Error(err))
	}
	return err
}

func (s *store) Getter() Getter {
	return s
}

func (s *store) CloudRunGetter(cloudProvider string) (cloudrun.Getter, bool) {
	ks, ok := s.cloudrunStores[cloudProvider]
	return ks, ok
}

func (s *store) ECSRunGetter(cloudProvider string) (ecs.Getter, bool) {
	ks, ok := s.ecsStores[cloudProvider]
	return ks, ok
}

func (s *store) KubernetesGetter(cloudProvider string) (kubernetes.Getter, bool) {
	ks, ok := s.kubernetesStores[cloudProvider]
	return ks, ok
}

func (s *store) LambdaGetter(cloudProvider string) (lambda.Getter, bool) {
	ks, ok := s.lambdaStores[cloudProvider]
	return ks, ok
}

func (s *store) TerraformGetter(cloudProvider string) (terraform.Getter, bool) {
	ks, ok := s.terraformStores[cloudProvider]
	return ks, ok
}

type LiveResourceLister struct {
	Getter
}

func (g LiveResourceLister) ListKubernetesAppLiveResources(cloudProvider, appID string) ([]provider.Manifest, bool) {
	kg, ok := g.KubernetesGetter(cloudProvider)
	if !ok {
		return nil, false
	}
	return kg.GetAppLiveManifests(appID), true
}
