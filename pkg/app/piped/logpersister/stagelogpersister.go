// Copyright 2023 The PipeCD Authors.
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
	"fmt"
	"sync"
	"time"

	"go.uber.org/atomic"
	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/model"
)

// stageLogPersister represents a log persister for a specific stage.
type stageLogPersister struct {
	key         key
	blocks      []*model.LogBlock
	curLogIndex int64
	completed   bool
	completedAt time.Time
	// Mutex to protect the fields above.
	mu sync.RWMutex

	sentIndex               int
	checkpointSentTimestamp time.Time
	done                    atomic.Bool
	doneCh                  chan struct{}

	checkpointFlushInterval time.Duration
	persister               *persister
	logger                  *zap.Logger
}

// append appends a new log block.
func (sp *stageLogPersister) append(log string, s model.LogSeverity) {
	now := time.Now()

	// We also send the error logs to the local logger.
	if s == model.LogSeverity_ERROR {
		sp.logger.Warn(fmt.Sprintf("STAGE ERROR LOG: %s", log))
	}

	sp.mu.Lock()
	defer sp.mu.Unlock()

	sp.curLogIndex++
	sp.blocks = append(sp.blocks, &model.LogBlock{
		Index:     sp.curLogIndex,
		Log:       log,
		Severity:  s,
		CreatedAt: now.Unix(),
	})
}

// Write appends a new INFO log block.
func (sp *stageLogPersister) Write(log []byte) (int, error) {
	sp.Info(string(log))
	return len(log), nil
}

// Info appends a new INFO log block.
func (sp *stageLogPersister) Info(log string) {
	sp.append(log, model.LogSeverity_INFO)
}

// Infof formats and appends a new INFO log block.
func (sp *stageLogPersister) Infof(format string, a ...interface{}) {
	sp.append(fmt.Sprintf(format, a...), model.LogSeverity_INFO)
}

// Success appends a new SUCCESS log block.
func (sp *stageLogPersister) Success(log string) {
	sp.append(log, model.LogSeverity_SUCCESS)
}

// Successf formats and appends a new SUCCESS log block.
func (sp *stageLogPersister) Successf(format string, a ...interface{}) {
	sp.append(fmt.Sprintf(format, a...), model.LogSeverity_SUCCESS)
}

// Error appends a new ERROR log block.
func (sp *stageLogPersister) Error(log string) {
	sp.append(log, model.LogSeverity_ERROR)
}

// Errorf formats and appends a new ERROR log block.
func (sp *stageLogPersister) Errorf(format string, a ...interface{}) {
	sp.append(fmt.Sprintf(format, a...), model.LogSeverity_ERROR)
}

// Complete marks the completion of logging for this stage.
// This means no more log for this stage will be added into this persister.
func (sp *stageLogPersister) Complete(timeout time.Duration) error {
	sp.mu.Lock()
	sp.completed = true
	sp.completedAt = time.Now()
	sp.mu.Unlock()

	timer := time.NewTimer(timeout)
	defer timer.Stop()

	select {
	case <-timer.C:
		return fmt.Errorf("timed out")

	case <-sp.doneCh:
		return nil
	}
}

func (sp *stageLogPersister) isStale(period time.Duration) bool {
	if sp.done.Load() {
		return true
	}

	sp.mu.RLock()
	defer sp.mu.RUnlock()

	if sp.completed && time.Since(sp.completedAt) > period {
		return true
	}
	return false
}

// flush sends the new log blocks or all of the blocks from the checkpoint
// based on the elapsed time.
// By design this flush function for a specific stageLogPersister will not be called concurrently.
func (sp *stageLogPersister) flush(ctx context.Context) error {
	// Do nothing when this persister is already done.
	if sp.done.Load() {
		return nil
	}

	sp.mu.RLock()
	completed := sp.completed
	sp.mu.RUnlock()

	if completed || time.Since(sp.checkpointSentTimestamp) > sp.checkpointFlushInterval {
		sp.checkpointSentTimestamp = time.Now()
		return sp.flushFromLastCheckpoint(ctx)
	}

	return sp.flushNewLogs(ctx)
}

func (sp *stageLogPersister) flushNewLogs(ctx context.Context) error {
	sp.mu.RLock()
	blocks := sp.blocks[sp.sentIndex:]
	sp.mu.RUnlock()

	numBlocks := len(blocks)
	if numBlocks == 0 {
		return nil
	}

	if err := sp.persister.reportStageLogs(ctx, sp.key, blocks); err != nil {
		return err
	}

	// Update sentIndex.
	sp.sentIndex += numBlocks
	return nil
}

func (sp *stageLogPersister) flushFromLastCheckpoint(ctx context.Context) (err error) {
	sp.mu.RLock()
	blocks := sp.blocks
	completed := sp.completed
	sp.mu.RUnlock()

	defer func() {
		if err == nil && completed {
			sp.done.Store(true)
			close(sp.doneCh)
		}
	}()

	numBlocks := len(blocks)
	if numBlocks == 0 {
		return nil
	}

	if err := sp.persister.reportStageLogsFromLastCheckpoint(ctx, sp.key, blocks, completed); err != nil {
		return err
	}

	// Remove all sent blocks and update checkpointSentIndex.
	sp.mu.Lock()
	sp.blocks = sp.blocks[numBlocks:]
	sp.mu.Unlock()

	sp.sentIndex = 0
	return nil
}
