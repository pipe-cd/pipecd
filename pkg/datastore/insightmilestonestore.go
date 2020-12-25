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
