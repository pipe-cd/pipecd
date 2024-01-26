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
	"os"
	"strings"
	"testing"

	"github.com/goccy/go-yaml"
	"github.com/stretchr/testify/assert"

	"github.com/pipe-cd/pipecd/pkg/app/pipectl/cmd/initialize/prompt"
	"github.com/pipe-cd/pipecd/pkg/config"
)

func TestGenerateKubernetesConfig(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name         string
		inputs       string // mock for user's input
		expectedFile string
		expectedErr  bool
	}{
		{
			name: "valid inputs for Kustomize",
			inputs: `myApp
				0
				5.3.0
				`,
			expectedFile: "testdata/k8s-app-kustomize.yaml",
			expectedErr:  false,
		},
		{
			name: "Kustomize specific fields are all empty",
			inputs: `myApp
				0

				`,
			expectedFile: "testdata/k8s-app-kustomize-empty.yaml",
			expectedErr:  false,
		},
		{
			name:         "missing required",
			inputs:       `myApp`, // missing selecting kustomize or helm
			expectedFile: "",
			expectedErr:  true,
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			reader := strings.NewReader(tc.inputs)
			prompt := prompt.NewPrompt(reader)

			// Generate the config
			cfg, err := generateKubernetesConfig(prompt)
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
			}
		})
	}
}
