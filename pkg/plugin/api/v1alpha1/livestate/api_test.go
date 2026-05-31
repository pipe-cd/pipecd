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

package livestate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLivestateRequestValidate(t *testing.T) {
	tests := []struct {
		name    string
		req     *GetLivestateRequest
		wantErr bool
	}{
		{
			name:    "nil struct is valid",
			req:     nil,
			wantErr: false,
		},
		{
			name:    "empty struct is valid",
			req:     &GetLivestateRequest{},
			wantErr: false,
		},
		{
			name: "full struct is valid",
			req: &GetLivestateRequest{
				PipedId:         "piped-id",
				ApplicationId:   "app-id",
				ApplicationName: "app-name",
				DeployTargets:   []string{"target-1"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
