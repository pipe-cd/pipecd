package orphancommandcleaner

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipe/pkg/datastore"
	"github.com/pipe-cd/pipe/pkg/model"
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
	t := time.NewTicker(interval)
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-t.C:
			start := time.Now()
			if err := c.updateOrphanCommandsStatus(ctx); err == nil {
				c.logger.Info("successfully update orphan commands status", zap.Duration("duration", time.Since(start)))
			}
		}
	}
}

func (c *OrphanCommandCleaner) updateOrphanCommandsStatus(ctx context.Context) error {
	timeout := time.Now().Add(-commandTimeOut).Unix()
	opts := datastore.ListOptions{
		Filters: []datastore.ListFilter{
			{
				Field:    "Status",
				Operator: "==",
				Value:    model.CommandStatus_COMMAND_NOT_HANDLED_YET,
			},
			{
				Field:    "CreatedAt",
				Operator: "<=",
				Value:    timeout,
			},
		},
	}
	commands, err := c.commandstore.ListCommands(ctx, opts)
	if err != nil {
		return err
	}

	for _, command := range commands {
		if err := c.commandstore.UpdateCommand(ctx, command.Id, func(cmd *model.Command) error {
			cmd.Status = model.CommandStatus_COMMAND_TIMEOUT
			return nil
		}); err != nil {
			c.logger.Error("failed to update orphan command", zap.Error(err), zap.String("id", command.Id), zap.String("type", command.Type.String()))
		}
	}

	return nil
}
