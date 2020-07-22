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

// Package notifier provides a piped component
// that sends notifications to the configured targets.
package notifier

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/pipe-cd/pipe/pkg/config"
	"github.com/pipe-cd/pipe/pkg/model"
)

type Notifier struct {
	config   *config.PipedSpec
	handlers []handler
	logger   *zap.Logger
}

type handler struct {
	matcher *matcher
	sender  sender
}

type sender interface {
	Run(ctx context.Context) error
	Notify(event model.Event)
}

func NewNotifier(cfg *config.PipedSpec, logger *zap.Logger) (*Notifier, error) {
	logger = logger.Named("notifier")
	receivers := make(map[string]config.NotificationReceiver, len(cfg.Notifications.Receivers))
	for _, r := range cfg.Notifications.Receivers {
		receivers[r.Name] = r
	}

	handlers := make([]handler, 0, len(cfg.Notifications.Routes))
	for _, route := range cfg.Notifications.Routes {
		receiver, ok := receivers[route.Receiver]
		if !ok {
			return nil, fmt.Errorf("missing receiver %s that is used in route %s", route.Receiver, route.Name)
		}

		var sd sender
		switch {
		case receiver.Slack != nil:
			sd = newSlackSender(receiver.Name, *receiver.Slack, cfg.WebURL, logger)
		case receiver.Webhook != nil:
			sd = newWebhookSender(receiver.Name, *receiver.Webhook, logger)
		default:
			continue
		}

		handlers = append(handlers, handler{
			matcher: newMatcher(route),
			sender:  sd,
		})
	}

	return &Notifier{config: cfg, handlers: handlers, logger: logger}, nil
}

func (n *Notifier) Run(ctx context.Context) error {
	group, ctx := errgroup.WithContext(ctx)

	for i := range n.handlers {
		sender := n.handlers[i].sender
		group.Go(func() error {
			return sender.Run(ctx)
		})
	}

	n.logger.Info(fmt.Sprintf("all %d notifiers have been started", len(n.handlers)))
	if err := group.Wait(); err != nil {
		n.logger.Error("failed while running", zap.Error(err))
		return err
	}

	n.logger.Info(fmt.Sprintf("all %d notifiers have been stopped", len(n.handlers)))
	return nil
}

func (n *Notifier) Notify(event model.Event) {
	for _, h := range n.handlers {
		if !h.matcher.Match(event) {
			continue
		}
		h.sender.Notify(event)
	}
}
