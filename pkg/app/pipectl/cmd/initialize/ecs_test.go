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

package initialize

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/goccy/go-yaml"
	"github.com/stretchr/testify/assert"

	"github.com/pipe-cd/pipecd/pkg/app/pipectl/cmd/initialize/prompt"
	"github.com/pipe-cd/pipecd/pkg/config"
)

func TestGenerateECSConfig(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name         string
		inputs       string // mock for user's input
		expectedFile string
		expectedErr  bool
	}{
		// QuickSync
		{
			name: "valid inputs for QuickSync",
			inputs: `myApp
				serviceDef.yaml
				taskDef.yaml
				0
				arn:aws:elasticloadbalancing:ap-northeast-1:123456789012:targetgroup/xxx/xxx
				web
				80
				`,
			expectedFile: "testdata/ecs-app.yaml",
			expectedErr:  false,
		},
		{
			name: "missing required deploymentStrategy",
			inputs: `myApp
				serviceDef.yaml
				taskDef.yaml
				
				`,
			expectedFile: "",
			expectedErr:  true,
		},
		// Canary
		{
			name: "valid inputs for Canary",
			inputs: `myApp
				serviceDef.yaml
				taskDef.yaml
				1 
				arn:aws:elasticloadbalancing:ap-northeast-1:123456789012:targetgroup/xxx/xxx
				web
				80
				arn:aws:elasticloadbalancing:ap-northeast-1:123456789012:targetgroup/yyy/yyy
				web
				80
				20
				`,
			expectedFile: "testdata/ecs-app-canary.yaml",
			expectedErr:  false,
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			reader := strings.NewReader(tc.inputs)
			prompt := prompt.NewPrompt(reader)

			// Generate the config
			cfg, err := generateECSConfig(prompt)
			assert.Equal(t, tc.expectedErr, err != nil)

			if err == nil {
				// Compare the YAML output
				yml, err := yaml.Marshal(cfg)
				assert.NoError(t, err)
				file, err := os.ReadFile(tc.expectedFile)
				assert.NoError(t, err)
				assert.Equal(t, string(file), string(yml))

				// Check if the YAML output is compatible with the original Config model
				_, err = config.DecodeYAML(yml)
				assert.NoError(t, err)
			} else {
				fmt.Println(err)
			}
		})
	}
}
