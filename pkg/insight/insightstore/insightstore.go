package insightstore

import (
	"context"

	"github.com/pipe-cd/pipecd/pkg/insight"
)

type Store interface {
	// LoadChunks(
	// 	ctx context.Context,
	// 	projectID, appID string,
	// 	kind model.InsightMetricsKind,
	// 	step model.InsightStep,
	// 	from time.Time,
	// 	count int,
	// ) (insight.Chunks, error)
	// PutChunk(ctx context.Context, chunk insight.Chunk) error
	// LoadMilestone(ctx context.Context) (*insight.Milestone, error)
	// PutMilestone(ctx context.Context, m *insight.Milestone) error

	LoadApplicationCounts(ctx context.Context, projectID string) (*insight.ApplicationCounts, error)
	PutApplicationCounts(ctx context.Context, projectID string, counts insight.ApplicationCounts) error
}
