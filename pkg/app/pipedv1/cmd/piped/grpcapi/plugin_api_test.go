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

package grpcapi

import (
	"testing"

	"github.com/pipe-cd/pipecd/pkg/config"
	"github.com/pipe-cd/pipecd/pkg/model"
)

// Test for initializeSecretDecrypter function.
func TestInitializeSecretDecrypter(t *testing.T) {
	testcases := []struct {
		name        string
		cfg         *config.SecretManagement
		expected    bool
		expectedErr bool
	}{
		{
			name:        "no secret management config",
			cfg:         nil,
			expected:    false,
			expectedErr: false,
		},
		{
			name:        "secret management type none",
			cfg:         &config.SecretManagement{Type: model.SecretManagementTypeNone},
			expected:    false,
			expectedErr: false,
		},
		{
			name:        "unsupported secret management type",
			cfg:         &config.SecretManagement{Type: "unsupported"},
			expected:    false,
			expectedErr: true,
		},
		{
			name:        "unspoerted secret management type GCPKMS",
			cfg:         &config.SecretManagement{Type: model.SecretManagementTypeGCPKMS},
			expected:    false,
			expectedErr: true,
		},
		{
			name:        "unsupported secret mamagement type AWSKMS",
			cfg:         &config.SecretManagement{Type: model.SecretManagementTypeAWSKMS},
			expected:    false,
			expectedErr: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			decrypter, err := initializeSecretDecrypter(tc.cfg)
			if (err != nil) != tc.expectedErr {
				t.Errorf("unexpected error: %v", err)
			}
			if (decrypter != nil) != tc.expected {
				t.Errorf("unexpected result: %v", decrypter)
			}
		})
	}
}
