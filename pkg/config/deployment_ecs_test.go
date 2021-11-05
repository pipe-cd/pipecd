// Copyright 2020 The PipeCD Authors.
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

package config

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestECSDeploymentConfig(t *testing.T) {
	testcases := []struct {
		fileName           string
		expectedKind       Kind
		expectedAPIVersion string
		expectedSpec       interface{}
		expectedError      error
	}{
		{
			fileName:           "testdata/application/ecs-app.yaml",
			expectedKind:       KindECSApp,
			expectedAPIVersion: "pipecd.dev/v1beta1",
			expectedSpec: &ECSDeploymentSpec{
				GenericDeploymentSpec: GenericDeploymentSpec{
					Timeout: Duration(6 * time.Hour),
				},
				Input: ECSDeploymentInput{
					ServiceDefinitionFile: "/path/to/servicedef.yaml",
					TaskDefinitionFile:    "/path/to/taskdef.yaml",
					TargetGroups: ECSTargetGroups{
						Primary: json.RawMessage(`{"containerName":"web","containerPort":80,"targetGroupArn":"arn:aws:elasticloadbalancing:xyz"}`),
						Canary:  nil,
					},
					AutoRollback: true,
				},
			},
			expectedError: nil,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.fileName, func(t *testing.T) {
			cfg, err := LoadFromYAML(tc.fileName)
			require.Equal(t, tc.expectedError, err)
			if err == nil {
				assert.Equal(t, tc.expectedKind, cfg.Kind)
				assert.Equal(t, tc.expectedAPIVersion, cfg.APIVersion)
				assert.Equal(t, tc.expectedSpec, cfg.spec)
			}
		})
	}
}
