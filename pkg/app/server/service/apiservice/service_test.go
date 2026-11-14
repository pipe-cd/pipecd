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

package apiservice

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pipe-cd/pipecd/pkg/model"
)

func TestSyncApplicationRequestValidate(t *testing.T) {
	testcases := []struct {
		name    string
		req     *SyncApplicationRequest
		wantErr bool
	}{
		{
			name: "valid: strategy left unset defaults to AUTO",
			req:  &SyncApplicationRequest{ApplicationId: "app-1"},
		},
		{
			name: "valid: quick sync",
			req: &SyncApplicationRequest{
				ApplicationId: "app-1",
				SyncStrategy:  model.SyncStrategy_QUICK_SYNC,
			},
		},
		{
			name: "valid: pipeline",
			req: &SyncApplicationRequest{
				ApplicationId: "app-1",
				SyncStrategy:  model.SyncStrategy_PIPELINE,
			},
		},
		{
			name: "invalid: undefined strategy",
			req: &SyncApplicationRequest{
				ApplicationId: "app-1",
				SyncStrategy:  model.SyncStrategy(999),
			},
			wantErr: true,
		},
		{
			name:    "invalid: missing application id",
			req:     &SyncApplicationRequest{},
			wantErr: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.wantErr, tc.req.Validate() != nil)
		})
	}
}
