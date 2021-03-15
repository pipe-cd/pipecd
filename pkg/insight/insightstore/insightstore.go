package insightstore

import (
	"context"
	"time"

	"github.com/pipe-cd/pipe/pkg/insight"
	"github.com/pipe-cd/pipe/pkg/model"
)

type IStore interface {
	LoadChunks(
		ctx context.Context,
		projectID, appID string,
		kind model.InsightMetricsKind,
		step model.InsightStep,
		from time.Time,
		count int,
	) (insight.Chunks, error)
	PutChunk(ctx context.Context, chunk insight.Chunk) error
	LoadMilestone(ctx context.Context) (*insight.Milestone, error)
	PutMilestone(ctx context.Context, m *insight.Milestone) error
	LoadApplicationCount(ctx context.Context, projectID string) (*insight.ApplicationCount, error)
	PutApplicationCount(ctx context.Context, ac *insight.ApplicationCount, projectID string) error
}
