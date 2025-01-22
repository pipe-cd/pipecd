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

// Package notifier provides a piped component
// that sends notifications to the configured targets.
package notifier

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/atomic"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	config "github.com/pipe-cd/pipecd/pkg/configv1"
	"github.com/pipe-cd/pipecd/pkg/model"
	"github.com/pipe-cd/pipecd/pkg/version"
)

type Notifier struct {
	config      *config.PipedSpec
	handlers    []handler
	gracePeriod time.Duration
	closed      atomic.Bool
	logger      *zap.Logger
}

type handler struct {
	matcher *matcher
	sender  sender
}

type sender interface {
	Run(ctx context.Context) error
	Notify(event model.NotificationEvent)
	Close(ctx context.Context)
}

func NewNotifier(cfg *config.PipedSpec, logger *zap.Logger) (*Notifier, error) {
	logger = logger.Named("notifier")
	receivers := make(map[string]config.NotificationReceiver, len(cfg.Notifications.Receivers))
	for _, r := range cfg.Notifications.Receivers {
		receivers[r.Name] = r
	}

	handlers := make([]handler, 0, len(cfg.Notifications.Routes))
	// for _, route := range cfg.Notifications.Routes {
	// 	// receiver, ok := receivers[route.Receiver]
	// 	// if !ok {
	// 	// 	return nil, fmt.Errorf("missing receiver %s that is used in route %s", route.Receiver, route.Name)
	// 	// }

	// 	var sd sender
	// 	switch {
	// 	case receiver.Slack != nil:
	// 		slacksender, err := newSlackSender(receiver.Name, *receiver.Slack, cfg.WebAddress, logger)
	// 		if err != nil {
	// 			return nil, fmt.Errorf("failed to create slack sender: %w", err)
	// 		}
	// 		sd = slacksender
	// 	case receiver.Webhook != nil:
	// 		sd = newWebhookSender(receiver.Name, *receiver.Webhook, cfg.WebAddress, logger)
	// 	default:
	// 		continue
	// 	}

	// 	handlers = append(handlers, handler{
	// 		matcher: newMatcher(route),
	// 		sender:  sd,
	// 	})
	// }

	return &Notifier{
		config:      cfg,
		handlers:    handlers,
		gracePeriod: 10 * time.Second,
		logger:      logger,
	}, nil
}

func (n *Notifier) Run(ctx context.Context) error {
	group, ctx := errgroup.WithContext(ctx)

	// Start running all senders.
	for i := range n.handlers {
		sender := n.handlers[i].sender
		group.Go(func() error {
			return sender.Run(ctx)
		})
	}

	// Send the PIPED_STARTED event.
	n.Notify(model.NotificationEvent{
		Type: model.NotificationEventType_EVENT_PIPED_STARTED,
		Metadata: &model.NotificationEventPipedStarted{
			Id:        n.config.PipedID,
			Name:      n.config.Name,
			Version:   version.Get().Version,
			ProjectId: n.config.ProjectID,
		},
	})

	n.logger.Info(fmt.Sprintf("all %d notifiers have been started", len(n.handlers)))
	if err := group.Wait(); err != nil {
		n.logger.Error("failed while running", zap.Error(err))
		return err
	}

	// Send the PIPED_STOPPED event.
	n.Notify(model.NotificationEvent{
		Type: model.NotificationEventType_EVENT_PIPED_STOPPED,
		Metadata: &model.NotificationEventPipedStopped{
			Id:        n.config.PipedID,
			Name:      n.config.Name,
			Version:   version.Get().Version,
			ProjectId: n.config.ProjectID,
		},
	})

	// Mark to ignore all incoming events from this time and close all senders.
	n.closed.Store(true)
	ctx, cancel := context.WithTimeout(context.Background(), n.gracePeriod)
	defer cancel()

	for i := range n.handlers {
		sender := n.handlers[i].sender
		sender.Close(ctx)
	}

	n.logger.Info(fmt.Sprintf("all %d notifiers have been stopped", len(n.handlers)))
	return nil
}

func (n *Notifier) Notify(event model.NotificationEvent) {
	if n.closed.Load() {
		n.logger.Warn("ignore an event because notifier is already closed", zap.String("type", event.Type.String()))
		return
	}
	for _, h := range n.handlers {
		if !h.matcher.Match(event) {
			continue
		}
		h.sender.Notify(event)
	}
}
