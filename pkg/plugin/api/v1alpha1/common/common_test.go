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

package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeploymentSourceValidate(t *testing.T) {
	tests := []struct {
		name    string
		source  *DeploymentSource
		wantErr bool
	}{
		{
			name:    "nil struct is valid",
			source:  nil,
			wantErr: false,
		},
		{
			name:    "empty struct is valid",
			source:  &DeploymentSource{},
			wantErr: false,
		},
		{
			name: "full struct is valid",
			source: &DeploymentSource{
				ApplicationDirectory:      "app-dir",
				CommitHash:                "commit-hash",
				ApplicationConfig:         []byte("config"),
				ApplicationConfigFilename: ".pipe.yaml",
				SharedConfigDirectory:     "shared-dir",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.source.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
