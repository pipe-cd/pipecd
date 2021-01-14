package orphancommandcleaner

import (
	"context"
	"time"

	"github.com/pipe-cd/pipe/pkg/datastore"
	"github.com/pipe-cd/pipe/pkg/model"
	"go.uber.org/zap"
)

var (
	commandTimeOut = 7 * 24 * time.Hour
	interval       = 24 * time.Hour
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
		logger:       logger.Named("orphan-command-finder"),
	}
}

func (o *OrphanCommandCleaner) Run(ctx context.Context) {
	for {
		if err := o.updateOrphanCommandsStatus(ctx); err != nil {
			o.logger.Error("failed to update orphan commands", zap.Error(err))
		}
		time.Sleep(interval)
	}
}

func (o *OrphanCommandCleaner) updateOrphanCommandsStatus(ctx context.Context) error {
	timeoutedTime := time.Now().Add(-commandTimeOut).Unix()
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
				Value:    timeoutedTime,
			},
		},
	}
	commands, err := o.commandstore.ListCommands(ctx, opts)
	if err != nil {
		return err
	}

	for _, c := range commands {
		o.commandstore.UpdateCommand(ctx, c.Id, func(c *model.Command) error {
			c.Status = model.CommandStatus_COMMAND_TIMEOUT
			return nil
		})
	}

	return nil
}
