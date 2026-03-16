// Copyright 2025 The PipeCD Authors.
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
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/config"
	"github.com/pipe-cd/pipecd/pkg/model"
)

func newTestWebhook(t *testing.T, serverURL, signatureKey, signatureValue string) *webhook {
	t.Helper()
	cfg := config.NotificationReceiverWebhook{
		URL:            serverURL,
		SignatureKey:   signatureKey,
		SignatureValue: signatureValue,
	}
	return newWebhookSender("test", cfg, "https://pipecd.example.com/", zap.NewNop())
}

func testEvent() model.NotificationEvent {
	return model.NotificationEvent{
		Type: model.NotificationEventType_EVENT_DEPLOYMENT_TRIGGERED,
		Metadata: &model.NotificationEventDeploymentTriggered{
			Deployment: &model.Deployment{
				Id:              "deploy-1",
				ApplicationName: "app-1",
				ProjectId:       "proj-1",
			},
		},
	}
}

func TestNewWebhookSender(t *testing.T) {
	t.Parallel()

	w := newTestWebhook(t, "https://example.com/hook", "X-Sig", "secret")

	assert.Equal(t, "test", w.name)
	assert.Equal(t, "https://pipecd.example.com", w.webURL)
	assert.NotNil(t, w.httpClient)
	assert.NotNil(t, w.eventCh)
}

func TestWebhook_SendEvent_Success(t *testing.T) {
	t.Parallel()

	var (
		receivedHeader string
		receivedBody   model.NotificationEvent
		callCount      int32
	)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&callCount, 1)
		receivedHeader = r.Header.Get("X-Signature")
		if err := json.NewDecoder(r.Body).Decode(&receivedBody); err != nil {
			http.Error(w, "bad body", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	w := newTestWebhook(t, srv.URL, "X-Signature", "my-secret")
	w.sendEvent(context.Background(), testEvent())

	assert.EqualValues(t, 1, atomic.LoadInt32(&callCount))
	assert.Equal(t, "my-secret", receivedHeader)
}

func TestWebhook_SendEvent_Non2xxLogsWarning(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	w := newTestWebhook(t, srv.URL, "X-Sig", "val")
	// Should not panic or return error — just logs a warning.
	w.sendEvent(context.Background(), testEvent())
}

func TestWebhook_SendEvent_InvalidURL(t *testing.T) {
	t.Parallel()

	w := newTestWebhook(t, "http://127.0.0.1:0/no-server", "X-Sig", "val")
	// Should not panic — logs an error internally.
	w.sendEvent(context.Background(), testEvent())
}

func TestWebhook_Notify_BuffersEvent(t *testing.T) {
	t.Parallel()

	w := newTestWebhook(t, "https://example.com", "X-Sig", "val")
	event := testEvent()

	w.Notify(event)

	require.Len(t, w.eventCh, 1)
	got := <-w.eventCh
	assert.Equal(t, event.Type, got.Type)
}

func TestWebhook_Run_ProcessesAndStopsOnContextCancel(t *testing.T) {
	t.Parallel()

	var callCount int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&callCount, 1)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	w := newTestWebhook(t, srv.URL, "X-Sig", "val")

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		done <- w.Run(ctx)
	}()

	w.Notify(testEvent())
	time.Sleep(50 * time.Millisecond)
	cancel()

	err := <-done
	assert.NoError(t, err)
	assert.EqualValues(t, 1, atomic.LoadInt32(&callCount))
}

func TestWebhook_Close_DrainsRemainingEvents(t *testing.T) {
	t.Parallel()

	var callCount int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&callCount, 1)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	w := newTestWebhook(t, srv.URL, "X-Sig", "val")

	// Buffer two events before closing.
	w.Notify(testEvent())
	w.Notify(testEvent())

	w.Close(context.Background())

	assert.EqualValues(t, 2, atomic.LoadInt32(&callCount))
}
