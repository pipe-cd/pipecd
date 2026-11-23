// Copyright 2026 The PipeCD Authors.
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

package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pipe-cd/pipecd/pkg/config"
)

func TestProviderRunSendsHeaders(t *testing.T) {
	t.Parallel()

	type received struct {
		host   string
		values []string
	}

	tests := []struct {
		name       string
		headers    []config.AnalysisHTTPHeader
		wantHost   string
		wantValues []string
	}{
		{
			name: "every value of a repeated header is sent",
			headers: []config.AnalysisHTTPHeader{
				{Key: "X-Custom", Value: "first"},
				{Key: "X-Custom", Value: "second"},
			},
			wantValues: []string{"first", "second"},
		},
		{
			name: "a Host header overrides the host of the URL",
			headers: []config.AnalysisHTTPHeader{
				{Key: "Host", Value: "app.example.com"},
				{Key: "X-Custom", Value: "value"},
			},
			wantHost:   "app.example.com",
			wantValues: []string{"value"},
		},
		{
			name: "the Host header name is case-insensitive",
			headers: []config.AnalysisHTTPHeader{
				{Key: "host", Value: "app.example.com"},
			},
			wantHost: "app.example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ch := make(chan received, 1)
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ch <- received{host: r.Host, values: r.Header.Values("X-Custom")}
				w.WriteHeader(http.StatusOK)
			}))
			t.Cleanup(server.Close)

			ok, _, err := NewProvider(0).Run(context.Background(), &config.AnalysisHTTP{
				URL:          server.URL,
				Method:       http.MethodGet,
				Headers:      tt.headers,
				ExpectedCode: http.StatusOK,
			})
			require.NoError(t, err)
			assert.True(t, ok)

			got := <-ch
			wantHost := tt.wantHost
			if wantHost == "" {
				wantHost = server.Listener.Addr().String()
			}
			assert.Equal(t, wantHost, got.host)
			assert.Equal(t, tt.wantValues, got.values)
		})
	}
}
