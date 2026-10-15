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

package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	analysisconfig "github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/analysis/config"
)

func TestProviderRun(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		statusCode       int
		responseBody     string
		expectedCode     int
		expectedResponse string
		wantOK           bool
		wantErr          string
	}{
		{
			name:             "matching expected response",
			statusCode:       http.StatusOK,
			responseBody:     `{"status":"healthy"}`,
			expectedCode:     http.StatusOK,
			expectedResponse: `{"status":"healthy"}`,
			wantOK:           true,
		},
		{
			name:             "unexpected response body",
			statusCode:       http.StatusOK,
			responseBody:     `{"status":"unhealthy"}`,
			expectedCode:     http.StatusOK,
			expectedResponse: `{"status":"healthy"}`,
			wantErr:          "unexpected response body",
		},
		{
			name:         "empty expected response skips body validation",
			statusCode:   http.StatusOK,
			responseBody: "any response",
			expectedCode: http.StatusOK,
			wantOK:       true,
		},
		{
			name:             "unexpected status code takes precedence",
			statusCode:       http.StatusServiceUnavailable,
			responseBody:     `{"status":"healthy"}`,
			expectedCode:     http.StatusOK,
			expectedResponse: `{"status":"healthy"}`,
			wantErr:          "unexpected status code 503",
		},
		{
			name:             "response body exceeds limit",
			statusCode:       http.StatusOK,
			responseBody:     strings.Repeat("a", maxResponseBodySize+1),
			expectedCode:     http.StatusOK,
			expectedResponse: "expected",
			wantErr:          "response body exceeds maximum size of 1048576 bytes",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				_, _ = w.Write([]byte(tt.responseBody))
			}))
			t.Cleanup(server.Close)

			ok, _, err := NewProvider(0).Run(context.Background(), &analysisconfig.AnalysisHTTP{
				URL:              server.URL,
				Method:           http.MethodGet,
				ExpectedCode:     tt.expectedCode,
				ExpectedResponse: tt.expectedResponse,
			})

			assert.Equal(t, tt.wantOK, ok)
			if tt.wantErr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.wantErr)
			}
		})
	}
}

func TestProviderRunExpectedResponseExceedsLimit(t *testing.T) {
	t.Parallel()

	ok, _, err := NewProvider(0).Run(context.Background(), &analysisconfig.AnalysisHTTP{
		URL:              "://invalid-url",
		Method:           http.MethodGet,
		ExpectedResponse: strings.Repeat("a", maxResponseBodySize+1),
	})

	assert.False(t, ok)
	assert.EqualError(t, err, "expected response exceeds maximum size of 1048576 bytes")
}
