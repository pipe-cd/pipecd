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

package cleaner

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
)

const (
	logFileName            = "step.log"
	completionMarkFileName = "completion"
)

type Option func(*Cleaner)

func WithPeriod(period time.Duration) Option {
	return func(c *Cleaner) {
		c.period = period
	}
}

func WithIterationTimeout(timeout time.Duration) Option {
	return func(c *Cleaner) {
		c.iterationTimeout = timeout
	}
}

func WithMaxTTL(ttl time.Duration) Option {
	return func(c *Cleaner) {
		c.maxTTL = ttl
	}
}

func WithLogger(logger *zap.Logger) Option {
	return func(c *Cleaner) {
		c.logger = logger.Named("cleaner")
	}
}

type Cleaner struct {
	dataDir          string
	period           time.Duration
	iterationTimeout time.Duration
	maxTTL           time.Duration
	logger           *zap.Logger
	ctx              context.Context
	cancel           context.CancelFunc
	doneCh           chan struct{}
}

func NewCleaner(dataDir string, opts ...Option) *Cleaner {
	ctx, cancel := context.WithCancel(context.Background())
	c := &Cleaner{
		dataDir:          dataDir,
		period:           time.Hour,
		iterationTimeout: 30 * time.Minute,
		maxTTL:           6 * time.Hour,
		logger:           zap.NewNop(),
		ctx:              ctx,
		cancel:           cancel,
		doneCh:           make(chan struct{}),
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *Cleaner) Run() error {
	ticker := time.NewTicker(c.period)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ctx, _ := context.WithTimeout(c.ctx, c.iterationTimeout)
			if err := c.clean(ctx); err != nil {
				c.logger.Error("failed while cleaning old logs", zap.Error(err))
			} else {
				c.logger.Info("successfully cleaned old logs")
			}
		case <-c.ctx.Done():
			c.logger.Info("cleaner is stopped")
			return nil
		}
	}
}

func (c *Cleaner) Stop() {
	c.cancel()
	<-c.doneCh
}

func (c *Cleaner) clean(ctx context.Context) error {
	var total, cleaned int64
	remover := func(path string) {
		if err := os.RemoveAll(path); err != nil {
			c.logger.Error(fmt.Sprintf("failed to remove old data: %s", path), zap.Error(err))
			return
		}
		cleaned++
		c.logger.Info(fmt.Sprintf("removed old data: %s", path))
	}
	err := filepath.Walk(c.dataDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if path == c.dataDir {
			return nil
		}
		if !info.IsDir() {
			return nil
		}
		total++
		if time.Since(info.ModTime()) > c.maxTTL {
			remover(path)
			return filepath.SkipDir
		}
		markFilePath := filepath.Join(path, completionMarkFileName)
		if _, err := os.Stat(markFilePath); os.IsNotExist(err) {
			remover(path)
		}
		return filepath.SkipDir
	})
	c.logger.Info(fmt.Sprintf("cleaned %d from %d logs", cleaned, total))
	return err
}
