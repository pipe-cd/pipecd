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

	"github.com/pipe-cd/pipecd/pkg/config"
	"github.com/pipe-cd/pipecd/pkg/model"
)

type applicationLister interface {
	List() []*model.Application
}

type Getter interface {
	// TODO: generic getter methods
}

type Store interface {
	Run(ctx context.Context) error
	Getter() Getter
}

// store manages a list of particular stores for all cloud providers.
type store struct {
	// TODO: generic store fields

	gracePeriod time.Duration
	logger      *zap.Logger
}

func NewStore(ctx context.Context, cfg *config.PipedSpec, appLister applicationLister, gracePeriod time.Duration, logger *zap.Logger) Store {
	logger = logger.Named("livestatestore")

	s := &store{
		gracePeriod:      gracePeriod,
		logger:           logger,
	}
	for _, cp := range cfg.PlatformProviders {
		_ = cp // TODO: general state from plugin from store fields
	}

	return s
}

func (s *store) Run(ctx context.Context) error {
	s.logger.Info("start running appsatestore")

	group, ctx := errgroup.WithContext(ctx)

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

type LiveResourceLister struct {
	Getter
}
