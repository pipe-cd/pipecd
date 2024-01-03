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

package planpreviewoutputcleaner

import (
	"context"
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/filestore"
)

const (
	outputTTL    = 48 * time.Hour
	cronSchedule = "0 9 * * *" // Run at 09:00 every day.
	prefix       = "command-output/"
)

type store interface {
	filestore.Lister
	filestore.Deleter
}

type Cleaner struct {
	store  store
	logger *zap.Logger
}

func NewCleaner(s store, logger *zap.Logger) *Cleaner {
	return &Cleaner{
		store:  s,
		logger: logger.Named("planpreview-output-cleaner"),
	}
}

func (c *Cleaner) Run(ctx context.Context) error {
	c.logger.Info("start running planpreview output cleaner")

	cr := cron.New()
	if _, err := cr.AddFunc(cronSchedule, func() { c.clean(ctx) }); err != nil {
		return err
	}

	cr.Start()
	<-ctx.Done()
	cr.Stop()

	c.logger.Info("planpreview output cleaner has been stopped")
	return nil
}

func (c *Cleaner) clean(ctx context.Context) error {
	c.logger.Info("will find stale planpreview outputs to delete")

	objects, err := c.store.List(ctx, prefix)
	if err != nil {
		c.logger.Error("failed to list planpreview output objects",
			zap.String("prefix", prefix),
			zap.Error(err),
		)
		return err
	}

	ttl := outputTTL.Seconds()
	now := time.Now()
	deletes := 0

	for _, obj := range objects {
		if float64(now.Unix()-obj.UpdatedAt) <= ttl {
			continue
		}
		if err := c.store.Delete(ctx, obj.Path); err != nil {
			c.logger.Error("failed to delete planpreview output object",
				zap.String("path", obj.Path),
				zap.Error(err),
			)
			continue
		}
		c.logger.Info("successfully deleted a stale planpreview output",
			zap.String("path", obj.Path),
		)
		deletes++
	}

	c.logger.Info(fmt.Sprintf("deleted %d/%d stale planpreview outputs", deletes, len(objects)))
	return nil
}
