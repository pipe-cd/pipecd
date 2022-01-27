// Copyright 2021 The PipeCD Authors.
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

package orphancommandcleaner

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/datastore"
	"github.com/pipe-cd/pipecd/pkg/model"
)

var (
	commandTimeOut = 24 * time.Hour
	interval       = 6 * time.Hour
)

type OrphanCommandCleaner struct {
	commandstore datastore.CommandStore
	logger       *zap.Logger
}

func NewOrphanCommandCleaner(
	ds datastore.DataStore,
	logger *zap.Logger,
) *OrphanCommandCleaner {
	return &OrphanCommandCleaner{
		commandstore: datastore.NewCommandStore(ds),
		logger:       logger.Named("orphan-command-cleaner"),
	}
}

func (c *OrphanCommandCleaner) Run(ctx context.Context) error {
	c.logger.Info("start running OrphanCommandCleaner")

	t := time.NewTicker(interval)
	for {
		select {
		case <-ctx.Done():
			break

		case <-t.C:
			start := time.Now()
			if err := c.updateOrphanCommandsStatus(ctx); err == nil {
				c.logger.Info("successfully cleaned orphan commands", zap.Duration("duration", time.Since(start)))
			}
		}
	}

	c.logger.Info("orphanCommandCleaner has been stopped")
	return nil
}

func (c *OrphanCommandCleaner) updateOrphanCommandsStatus(ctx context.Context) error {
	timeout := time.Now().Add(-commandTimeOut).Unix()
	opts := datastore.ListOptions{
		Filters: []datastore.ListFilter{
			{
				Field:    "Status",
				Operator: datastore.OperatorEqual,
				Value:    model.CommandStatus_COMMAND_NOT_HANDLED_YET,
			},
			{
				Field:    "CreatedAt",
				Operator: datastore.OperatorLessThanOrEqual,
				Value:    timeout,
			},
		},
	}
	commands, err := c.commandstore.ListCommands(ctx, opts)
	if err != nil {
		c.logger.Error("failed to list not-handled commands", zap.Error(err))
		return err
	}

	c.logger.Info(fmt.Sprintf("there are %d orphan commands to clean", len(commands)))
	if len(commands) == 0 {
		return nil
	}

	for _, command := range commands {
		err := c.commandstore.UpdateCommand(ctx, datastore.OpsWriter, command.Id, func(cmd *model.Command) error {
			cmd.Status = model.CommandStatus_COMMAND_TIMEOUT
			return nil
		})
		if err != nil {
			c.logger.Error("failed to mark orphan command as timed out",
				zap.String("id", command.Id),
				zap.String("type", command.Type.String()),
				zap.Error(err),
			)
		}
	}

	return nil
}
