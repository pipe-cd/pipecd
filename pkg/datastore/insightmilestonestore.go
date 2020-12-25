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

const insightModelKind = "InsightMilestone"
const insightMilestone = "InsightMilestone"

type MilestornStore interface {
	GetInsightMilestone(ctx context.Context) (*insight.InsightMilestone, error)
	PutInsightMilestone(ctx context.Context) error
}

func (s *InsightMilestoneStore) GetInsightMilestone(ctx context.Context) (*insight.InsightMilestone, error) {
	var entity insight.InsightMilestone
	if err := s.ds.Get(ctx, insightModelKind, insightMilestone, &entity); err != nil {
		return nil, err
	}
	return &entity, nil
}

func (s *InsightMilestoneStore) PutInsightMilestone(ctx context.Context, m *insight.InsightMilestone) error {
	return s.ds.Put(ctx, insightModelKind, insightMilestone, m)
}
