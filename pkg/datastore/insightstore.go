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

const insightModelKind = "Insight"

var insightFactory = func() interface{} {
	return &model.InsightDataPoint{}
}

type InsightStore interface {
	AddInsight(ctx context.Context, idp *model.InsightDataPoint) error
	ListInsights(ctx context.Context, opts ListOptions) ([]*model.InsightDataPoint, error)
}

type insightStore struct {
	backend
	nowFunc func() time.Time
}

func NewInsightStore(ds DataStore) InsightStore {
	return &insightStore{
		backend: backend{
			ds: ds,
		},
		nowFunc: time.Now,
	}
}

func (s *insightStore) AddInsight(ctx context.Context, idp *model.InsightDataPoint) error {
	now := s.nowFunc().Unix()
	if idp.CreatedAt == 0 {
		idp.CreatedAt = now
	}
	if err := idp.Validate(); err != nil {
		return err
	}
	return s.ds.Create(ctx, insightModelKind, idp.Id, idp)
}

func (s *insightStore) ListInsights(ctx context.Context, opts ListOptions) ([]*model.InsightDataPoint, error) {
	it, err := s.ds.Find(ctx, insightModelKind, opts)
	if err != nil {
		return nil, err
	}
	idps := make([]*model.InsightDataPoint, 0)
	for {
		var idp model.InsightDataPoint
		err := it.Next(&idp)
		if err == ErrIteratorDone {
			break
		}
		if err != nil {
			return nil, err
		}
		idps = append(idps, &idp)
	}
	return idps, nil
}
