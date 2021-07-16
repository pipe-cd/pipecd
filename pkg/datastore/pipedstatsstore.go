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

package datastore

import (
	"context"
	"time"

	"github.com/pipe-cd/pipe/pkg/model"
)

const PipedStatsModelKind = "PipedStats"

var pipedStatsFactory = func() interface{} {
	return &model.PipedStats{}
}

// Deprecated: PipedStats model is deprecated, along with its store interface.
type PipedStatsStore interface {
	AddPipedStats(ctx context.Context, ps *model.PipedStats) error
	ListPipedStatss(ctx context.Context, opts ListOptions) ([]model.PipedStats, error)
}

type pipedStatsStore struct {
	backend
	nowFunc func() time.Time
}

// Deprecated: PipedStats model is deprecated, along with its store interface.
func NewPipedStatsStore(ds DataStore) PipedStatsStore {
	return &pipedStatsStore{
		backend: backend{
			ds: ds,
		},
		nowFunc: time.Now,
	}
}

func (s *pipedStatsStore) AddPipedStats(ctx context.Context, ps *model.PipedStats) error {
	now := s.nowFunc().Unix()
	if ps.Timestamp == 0 {
		ps.Timestamp = now
	}
	if err := ps.Validate(); err != nil {
		return err
	}
	return s.ds.Create(ctx, PipedStatsModelKind, "", ps)
}

func (s *pipedStatsStore) ListPipedStatss(ctx context.Context, opts ListOptions) ([]model.PipedStats, error) {
	it, err := s.ds.Find(ctx, PipedStatsModelKind, opts)
	if err != nil {
		return nil, err
	}
	pss := make([]model.PipedStats, 0)
	for {
		var ps model.PipedStats
		err := it.Next(&ps)
		if err == ErrIteratorDone {
			break
		}
		if err != nil {
			return nil, err
		}
		pss = append(pss, ps)
	}
	return pss, nil
}
