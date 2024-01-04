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

package notifier

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/config"
	"github.com/pipe-cd/pipecd/pkg/model"
)

const eventChannelBufferSize = 1000

type webhook struct {
	name       string
	config     config.NotificationReceiverWebhook
	webURL     string
	httpClient *http.Client
	eventCh    chan model.NotificationEvent
	logger     *zap.Logger
}

func newWebhookSender(name string, cfg config.NotificationReceiverWebhook, webURL string, logger *zap.Logger) *webhook {
	return &webhook{
		name:   name,
		config: cfg,
		webURL: strings.TrimRight(webURL, "/"),
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
		eventCh: make(chan model.NotificationEvent, eventChannelBufferSize),
		logger:  logger.Named("webhook").With(zap.String("name", name)),
	}
}

func (w *webhook) Run(ctx context.Context) error {
	for {
		select {
		case event, ok := <-w.eventCh:
			if ok {
				w.sendEvent(ctx, event)
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func (w *webhook) Notify(event model.NotificationEvent) {
	w.eventCh <- event
}

func (w *webhook) sendEvent(ctx context.Context, event model.NotificationEvent) {
	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(event); err != nil {
		w.logger.Error("unable to send data to webhook url", zap.Error(err))
		return
	}

	req, err := http.NewRequestWithContext(ctx, "POST", w.config.URL, buf)
	if err != nil {
		w.logger.Error("unable to send data to webhook url", zap.Error(err))
		return
	}

	signature, err := w.config.LoadSignatureValue()
	if err != nil {
		w.logger.Error("unable to load webhook signature value", zap.Error(err))
		return
	}

	req.Header.Add(w.config.SignatureKey, signature)

	resp, err := w.httpClient.Do(req)
	if err != nil {
		w.logger.Error("unable to send data to webhook url", zap.Error(err))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		w.logger.Warn("unexpected status was returned from the destination of webhook", zap.String("status", resp.Status))
	}
}

func (w *webhook) Close(ctx context.Context) {
	close(w.eventCh)

	// Send all remaining events.
	for {
		select {
		case event, ok := <-w.eventCh:
			if !ok {
				return
			}
			w.sendEvent(ctx, event)
		case <-ctx.Done():
			return
		}
	}
}
