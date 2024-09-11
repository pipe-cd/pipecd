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

package httpapi

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseProjectAndState(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		formValues    url.Values
		expectedState string
		expectedProj  string
		expectErr     bool
	}{
		{
			name:          "missing state",
			formValues:    url.Values{},
			expectedState: "",
			expectedProj:  "",
			expectErr:     true,
		},
		{
			name: "state without project id",
			formValues: url.Values{
				stateFormKey: {"state-token"},
			},
			expectedState: "state-token",
			expectedProj:  "",
			expectErr:     true,
		},
		{
			name: "state with project id",
			formValues: url.Values{
				stateFormKey: {"state-token:project-id"},
			},
			expectedState: "state-token",
			expectedProj:  "project-id",
			expectErr:     false,
		},
		{
			name: "state without colon and project id in form",
			formValues: url.Values{
				stateFormKey:   {"state-token"},
				projectFormKey: {"project-id"},
			},
			expectedState: "state-token",
			expectedProj:  "project-id",
			expectErr:     false,
		},
		{
			name: "state with colon but missing project id in form",
			formValues: url.Values{
				stateFormKey: {"state-token:"},
			},
			expectedState: "state-token",
			expectedProj:  "",
			expectErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Form = tt.formValues

			state, project, err := parseProjectAndState(req)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedState, state)
			assert.Equal(t, tt.expectedProj, project)
		})
	}
}
