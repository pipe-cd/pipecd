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

	"github.com/pipe-cd/pipe/pkg/insight"
)

type InsightMilestoneStore struct {
	backend
}

func NewInsightMilestoneStore(ds DataStore) *InsightMilestoneStore {
	return &InsightMilestoneStore{
		backend: backend{
			ds: ds,
		},
	}
}

const insightModelKind = "insight"
const insightMilestone = "milestone"

type MilestornStore interface {
	GetInsightMilestone(ctx context.Context) (*insight.Milestone, error)
	PutInsightMilestone(ctx context.Context) error
}

func (s *InsightMilestoneStore) GetInsightMilestone(ctx context.Context) (*insight.Milestone, error) {
	var entity insight.Milestone
	if err := s.ds.Get(ctx, insightModelKind, insightMilestone, &entity); err != nil {
		return nil, err
	}
	return &entity, nil
}

func (s *InsightMilestoneStore) PutInsightMilestone(ctx context.Context, m *insight.Milestone) error {
	return s.ds.Put(ctx, insightModelKind, insightMilestone, m)
}
