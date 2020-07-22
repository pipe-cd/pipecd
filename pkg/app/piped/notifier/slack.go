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

package notifier

import (
	"context"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipe/pkg/config"
	"github.com/pipe-cd/pipe/pkg/model"
)

type slack struct {
	name   string
	config config.NotificationReceiverSlack
	logger *zap.Logger
}

func newSlackSender(name string, cfg config.NotificationReceiverSlack, logger *zap.Logger) *slack {
	return &slack{
		name:   name,
		config: cfg,
		logger: logger.Named("slack"),
	}
}

func (s *slack) Run(ctx context.Context) error {
	return nil
}

func (s *slack) Notify(event model.Event) {
}
