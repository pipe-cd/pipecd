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

package handler

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/config"
	"github.com/pipe-cd/pipecd/pkg/model"
)

func createMockHandler(ctrl *gomock.Controller) (*MockprojectStore, *Handler) {
	m := NewMockprojectStore(ctrl)
	logger, _ := zap.NewProduction()

	h := NewHandler(
		10101,
		m,
		[]config.SharedSSOConfig{},
		0,
		logger,
	)

	return m, h
}

func setupMockHandler(ctrl *gomock.Controller, errToReturn error) (*Handler, string, *model.Project) {
	m, h := createMockHandler(ctrl)

	// fake details
	id := "test_id"
	project := &model.Project{
		Id: id,
	}

	if errToReturn == nil {
		// mock the call to the project store
		m.EXPECT().Get(gomock.Any(), gomock.Eq(id)).Return(project, nil)
	} else {
		// mock the call to the project store
		m.EXPECT().Get(gomock.Any(), gomock.Eq(id)).Return(nil, errToReturn)
	}

	return h, id, project
}

func TestGetProjectByIDOrReturnError(t *testing.T) {
	// Create a new controller
	ctrl := gomock.NewController(t)

	h, id, project := setupMockHandler(ctrl, nil)
	w := httptest.NewRecorder()

	actualProject := h.getProjectByIDOrReturnError(id, w)

	assert.Equal(t, project, actualProject)
}

func TestGetProjectByIDOrReturnErrorError(t *testing.T) {
	// Create a new controller
	ctrl := gomock.NewController(t)

	h, id, _ := setupMockHandler(ctrl, errors.New("example error"))
	w := httptest.NewRecorder()

	actualProject := h.getProjectByIDOrReturnError(id, w)

	assert.Nil(t, actualProject)

	res := w.Result()
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)

	defer res.Body.Close()
	data, _ := io.ReadAll(res.Body)

	assert.Equal(t, "Unable to retrieve existing project (example error)\n", string(data))
}

func TestRun(t *testing.T) {
	// Create a new controller
	ctrl := gomock.NewController(t)
	_, h := createMockHandler(ctrl)

	go func() {
		h.Run(context.TODO())
	}()
	err := h.stop()

	assert.Nil(t, err)
}

func buildHandleResetPasswordRequest(method string, id string, confirmationID string) *http.Request {
	req, _ := http.NewRequest(method, "/projects/reset-password", nil)

	q := req.URL.Query()
	if id != "" {
		q.Add("ID", id)
	}
	if confirmationID != "" {
		q.Add("confirmationID", confirmationID)
	}
	req.URL.RawQuery = q.Encode()

	return req
}

func TestHandleResetPassword(t *testing.T) {
	// Create a new controller
	ctrl := gomock.NewController(t)
	m, h := createMockHandler(ctrl)

	// fake details
	id := "test_id"
	project := &model.Project{
		Id: id,
	}

	testcases := []struct {
		name           string
		req            *http.Request
		expectedStatus int
		expectedBody   string
		extraMocks     func()
	}{
		{
			name:           "wrong method",
			req:            buildHandleResetPasswordRequest("PUT", "", ""),
			expectedStatus: http.StatusNotFound,
			expectedBody:   "not found",
		},
		{
			name:           "get returns confirmation page",
			req:            buildHandleResetPasswordRequest("GET", id, ""),
			expectedStatus: http.StatusOK,
			expectedBody:   fmt.Sprintf("Confirm you want to reset the static admin password for %s", id),
			extraMocks: func() {
				m.EXPECT().Get(gomock.Any(), gomock.Eq(id)).Return(project, nil)
			},
		},
		{
			name:           "missing-id-from-query-and-post",
			req:            buildHandleResetPasswordRequest("POST", "", ""),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid id",
		},
		{
			name:           "missing-confirmation-id",
			req:            buildHandleResetPasswordRequest("POST", id, ""),
			expectedStatus: http.StatusOK,
			expectedBody:   "Missing confirmation ID",
			extraMocks: func() {
				m.EXPECT().Get(gomock.Any(), gomock.Eq(id)).Return(project, nil)
			},
		},
		{
			name:           "wrong-confirmation-id",
			req:            buildHandleResetPasswordRequest("POST", id, fmt.Sprintf("%s mis match", id)),
			expectedStatus: http.StatusOK,
			expectedBody:   "Confirmation ID doesn&#39;t match",
			extraMocks: func() {
				m.EXPECT().Get(gomock.Any(), gomock.Eq(id)).Return(project, nil)
			},
		},
		{
			name:           "valid-reset-post",
			req:            buildHandleResetPasswordRequest("POST", id, id),
			expectedStatus: http.StatusOK,
			expectedBody:   "Successfully reset password for project",
			extraMocks: func() {
				m.EXPECT().Get(gomock.Any(), gomock.Eq(id)).Return(project, nil)
				m.EXPECT().UpdateProjectStaticAdmin(gomock.Any(), id, gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		{
			name:           "unable to update project static admin",
			req:            buildHandleResetPasswordRequest("POST", id, id),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "Unable to reset the password for project",
			extraMocks: func() {
				m.EXPECT().Get(gomock.Any(), gomock.Eq(id)).Return(project, nil)
				m.EXPECT().UpdateProjectStaticAdmin(gomock.Any(), id, gomock.Any(), gomock.Any()).Return(errors.New("error updating admin"))
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			if tc.extraMocks != nil {
				tc.extraMocks()
			}

			h.handleResetPassword(w, tc.req)

			res := w.Result()
			assert.Equal(t, tc.expectedStatus, res.StatusCode)

			defer res.Body.Close()
			data, _ := io.ReadAll(res.Body)

			assert.True(t, strings.Contains(string(data), tc.expectedBody))
		})
	}
}

func TestApplicationCountsTmpl(t *testing.T) {
	testcases := []struct {
		name        string
		data        []map[string]interface{}
		expected    string
		expectedErr error
	}{
		{
			name: "ok",
			data: []map[string]interface{}{
				{
					"Project": "one-count",
					"Total":   5,
					"Counts": map[string]int{
						"KUBERNETES": 5,
					},
				},
				{
					"Project": "not-found",
					"Error":   "No data for this project",
				},
				{
					"Project": "multi-counts",
					"Total":   20,
					"Counts": map[string]int{
						"KUBERNETES": 10,
						"CLOUD_RUN":  2,
						"LAMBDA":     8,
					},
				},
			},
			expected: `<!DOCTYPE html>
<html>
<head>
<style>
table {
  font-family: arial, sans-serif;
  border-collapse: collapse;
  width: 100%;
}

td, th {
  border: 1px solid #dddddd;
  text-align: left;
  padding: 8px;
}

tr:nth-child(1) {
  background-color: #dddddd;
}
</style>
</head>
<body>

<h2 style="text-align: center;"><a href="/">Welcome to PipeCD Owner Page!</a></h2>

<h3>There are 3 registered projects</h3>
<h4>0. one-count</h4>

<table>
  <tr>
    <th>Application Kind</th>
    <th>Count</th>
  </tr>
  <tr>
    <td>KUBERNETES</td>
    <td>5</td>
  </tr>
  <tr>
    <td>TOTAL</td>
    <td>5</td>
  </tr>
</table>

<h4>1. not-found</h4>

Unable to fetch application counts (No data for this project).

<h4>2. multi-counts</h4>

<table>
  <tr>
    <th>Application Kind</th>
    <th>Count</th>
  </tr>
  <tr>
    <td>CLOUD_RUN</td>
    <td>2</td>
  </tr>
  <tr>
    <td>KUBERNETES</td>
    <td>10</td>
  </tr>
  <tr>
    <td>LAMBDA</td>
    <td>8</td>
  </tr>
  <tr>
    <td>TOTAL</td>
    <td>20</td>
  </tr>
</table>

</body>
</html>
`,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := applicationCountsTmpl.Execute(&buf, tc.data)
			assert.Equal(t, tc.expected, buf.String())
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}
