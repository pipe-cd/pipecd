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

package logpersister

import (
	"context"
	"testing"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/cmd/piped/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/atomic"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type fakeAPIClient struct {
	reportStageLogsCount                   atomic.Uint32
	reportStageLogsFromLastCheckpointCount atomic.Uint32
}

func (c *fakeAPIClient) ReportStageLogs(ctx context.Context, in *service.ReportStageLogsRequest, opts ...grpc.CallOption) (*service.ReportStageLogsResponse, error) {
	c.reportStageLogsCount.Inc()
	return &service.ReportStageLogsResponse{}, nil
}

func (c *fakeAPIClient) ReportStageLogsFromLastCheckpoint(ctx context.Context, in *service.ReportStageLogsFromLastCheckpointRequest, opts ...grpc.CallOption) (*service.ReportStageLogsFromLastCheckpointResponse, error) {
	c.reportStageLogsFromLastCheckpointCount.Inc()
	return &service.ReportStageLogsFromLastCheckpointResponse{}, nil
}

func (c *fakeAPIClient) NumberOfReportStageLogs() int {
	return int(c.reportStageLogsCount.Load())
}

func (c *fakeAPIClient) NumberOfReportStageLogsFromLastCheckpoint() int {
	return int(c.reportStageLogsFromLastCheckpointCount.Load())
}

func TestPersister(t *testing.T) {
	t.Parallel()

	apiClient := &fakeAPIClient{}
	p := NewPersister(apiClient, zap.NewNop())
	p.stalePeriod = 0

	flushes, deletes := p.flush(context.TODO())
	require.Equal(t, 0, apiClient.NumberOfReportStageLogs())
	require.Equal(t, 0, apiClient.NumberOfReportStageLogsFromLastCheckpoint())
	assert.Equal(t, 0, flushes)
	assert.Equal(t, 0, deletes)

	num := p.flushAll(context.TODO())
	require.Equal(t, 0, apiClient.NumberOfReportStageLogs())
	require.Equal(t, 0, apiClient.NumberOfReportStageLogsFromLastCheckpoint())
	assert.Equal(t, 0, num)

	sp1 := p.StageLogPersister("deployment-1", "stage-1")
	p.StageLogPersister("deployment-2", "stage-2")

	num = p.flushAll(context.TODO())
	require.Equal(t, 0, apiClient.NumberOfReportStageLogs())
	require.Equal(t, 0, apiClient.NumberOfReportStageLogsFromLastCheckpoint())
	assert.Equal(t, 2, num)

	sp1.Complete(0)

	flushes, deletes = p.flush(context.TODO())
	require.Equal(t, 0, apiClient.NumberOfReportStageLogs())
	require.Equal(t, 0, apiClient.NumberOfReportStageLogsFromLastCheckpoint())
	assert.Equal(t, 1, flushes)
	assert.Equal(t, 1, deletes)

	num = p.flushAll(context.TODO())
	require.Equal(t, 0, apiClient.NumberOfReportStageLogs())
	require.Equal(t, 0, apiClient.NumberOfReportStageLogsFromLastCheckpoint())
	assert.Equal(t, 1, num)
}
