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

// Package logpersister provides a piped component
// that enqueues all log blocks from running stages
// and then periodically sends to the control plane.
package logpersister

import (
	"context"
	"sync"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/kapetaniosci/pipe/pkg/app/api/service/pipedservice"
	"github.com/kapetaniosci/pipe/pkg/model"
)

type apiClient interface {
	ReportStageLogs(ctx context.Context, in *pipedservice.ReportStageLogsRequest, opts ...grpc.CallOption) (*pipedservice.ReportStageLogsResponse, error)
	ReportStageLogsFromLastCheckpoint(ctx context.Context, in *pipedservice.ReportStageLogsFromLastCheckpointRequest, opts ...grpc.CallOption) (*pipedservice.ReportStageLogsFromLastCheckpointResponse, error)
}

type Persister interface {
	Run(ctx context.Context) error
	StageLogPersister(deploymentID, stageID string) StageLogPersister
}

type StageLogPersister interface {
	Append(log string, s model.LogSeverity)
	AppendInfo(log string)
	AppendSuccess(log string)
	AppendError(log string)
	Complete(timeout time.Duration) error
}

type key struct {
	DeploymentID string
	StageID      string
}

type persister struct {
	apiClient               apiClient
	stagePersisters         sync.Map
	flushInterval           time.Duration
	checkpointFlushInterval time.Duration
	stalePeriod             time.Duration
	gracePeriod             time.Duration
	logger                  *zap.Logger
}

// NewPersister creates a new persister instance for saving the stage logs into server's storage.
// This controls how many concurent api calls should be executed and when to flush the logs.
func NewPersister(apiClient apiClient, logger *zap.Logger) Persister {
	return &persister{
		apiClient:               apiClient,
		flushInterval:           5 * time.Second,
		checkpointFlushInterval: 2 * time.Minute,
		stalePeriod:             time.Minute,
		gracePeriod:             30 * time.Second,
		logger:                  logger.Named("logger-persister"),
	}
}

// Run starts running workers to flush logs to server.
func (p *persister) Run(ctx context.Context) error {
	p.logger.Info("start running log persister")
	ticker := time.NewTicker(p.flushInterval)
	defer ticker.Stop()

L:
	for {
		select {
		case <-ticker.C:
			p.flush(ctx)

		case <-ctx.Done():
			break L
		}
	}

	p.logger.Info("flush all logs before stopping")
	ctx, cancel := context.WithTimeout(context.Background(), p.gracePeriod)
	defer cancel()
	p.flushAll(ctx)

	p.logger.Info("log persister has been stopped")
	return nil
}

// StageLogPersister creates a child persister instance for a specific stage.
func (p *persister) StageLogPersister(deploymentID, stageID string) StageLogPersister {
	var (
		k = key{
			DeploymentID: deploymentID,
			StageID:      stageID,
		}
		logger = p.logger.With(
			zap.String("deployment-id", deploymentID),
			zap.String("stage-id", stageID),
		)
		sp = &stageLogPersister{
			key:                     k,
			curLogIndex:             time.Now().Unix(),
			doneCh:                  make(chan struct{}),
			checkpointFlushInterval: p.checkpointFlushInterval,
			persister:               p,
			logger:                  logger,
		}
	)
	p.stagePersisters.Store(k, sp)
	return sp
}

func (p *persister) flush(ctx context.Context) {
	completedKeys := make([]key, 0)

	// Check for new log entries and flush if needed.
	p.stagePersisters.Range(func(_, v interface{}) bool {
		sp := v.(*stageLogPersister)
		if sp.isStale(p.stalePeriod) {
			completedKeys = append(completedKeys, sp.key)
		} else {
			sp.flush(ctx)
		}
		return false
	})

	// Clean up all completed stages.
	for _, k := range completedKeys {
		p.stagePersisters.Delete(k)
	}
}

func (p *persister) flushAll(ctx context.Context) {
	p.stagePersisters.Range(func(_, v interface{}) bool {
		sp := v.(*stageLogPersister)
		go sp.flushFromLastCheckpoint(ctx)
		return false
	})
}

func (p *persister) reportStageLogs(ctx context.Context, k key, blocks []*model.LogBlock) error {
	req := &pipedservice.ReportStageLogsRequest{
		DeploymentId: k.DeploymentID,
		StageId:      k.StageID,
		Blocks:       blocks,
	}
	if _, err := p.apiClient.ReportStageLogs(ctx, req); err != nil {
		p.logger.Error("failed to report stage logs",
			zap.Any("key", k),
			zap.Error(err),
		)
		return err
	}
	return nil
}

func (p *persister) reportStageLogsFromLastCheckpoint(ctx context.Context, k key, blocks []*model.LogBlock, completed bool) error {
	req := &pipedservice.ReportStageLogsFromLastCheckpointRequest{
		DeploymentId: k.DeploymentID,
		StageId:      k.StageID,
		Blocks:       blocks,
		Completed:    completed,
	}
	if _, err := p.apiClient.ReportStageLogsFromLastCheckpoint(ctx, req); err != nil {
		p.logger.Error("failed to report stage logs from last checkpoint",
			zap.Any("key", k),
			zap.Error(err),
		)
		return err
	}
	return nil
}
