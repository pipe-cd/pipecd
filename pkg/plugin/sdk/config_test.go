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

package sdk

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Define a test struct and its pointer type satisfying the Spec interface.
type testSpecRuntime struct {
	Field1 string `yaml:"field1"`
	Field2 int    `yaml:"field2"`
	// Used to trigger validation error in tests
	ShouldFailValidation bool `yaml:"shouldFailValidation"`
}

func (s *testSpecRuntime) Validate() error {
	if s.ShouldFailValidation {
		return errors.New("validation failed")
	}
	return nil
}

func TestLoadConfigSpec(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name       string
		inputYAML  string
		wantErr    bool
		wantErrMsg string
		wantSpec   *testSpecRuntime
	}{
		{
			name: "success",
			inputYAML: `
apiVersion: pipecd.dev/v1beta1
kind: TestPlugin
spec:
  field1: "value1"
  field2: 123
`,
			wantErr: false,
			wantSpec: &testSpecRuntime{
				Field1: "value1",
				Field2: 123,
			},
		},
		{
			name: "invalid yaml",
			inputYAML: `
apiVersion: pipecd.dev/v1beta1
kind: TestPlugin
spec:
  field1: "value1"
  field2: 123
invalid-yaml-token
`,
			wantErr:    true,
			wantErrMsg: "yaml:",
		},
		{
			name: "validation fail",
			inputYAML: `
apiVersion: pipecd.dev/v1beta1
kind: TestPlugin
spec:
  field1: "value1"
  field2: 456
  shouldFailValidation: true
`,
			wantErr:    true,
			wantErrMsg: "validation failed",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ds := DeploymentSource{
				ApplicationConfig: []byte(tc.inputYAML),
			}

			spec, err := LoadConfigSpec[*testSpecRuntime](ds)

			if tc.wantErr {
				assert.Error(t, err)
				assert.Nil(t, spec)
				if tc.wantErrMsg != "" {
					assert.ErrorContains(t, err, tc.wantErrMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.wantSpec, spec)
			}
		})
	}
}

func TestLoadStages(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name       string
		inputYAML  string
		wantErr    bool
		wantStages Stages
	}{
		{
			name: "success with stages",
			inputYAML: `
apiVersion: pipecd.dev/v1beta1
kind: TestApp
spec:
  pipeline:
    stages:
      - name: STAGE1
        desc: Stage 1 description
      - name: STAGE2
        desc: Stage 2 description
`,
			wantErr: false,
			wantStages: Stages{
				stages: []string{"STAGE1", "STAGE2"},
			},
		},
		{
			name: "success with no pipeline",
			inputYAML: `
apiVersion: pipecd.dev/v1beta1
kind: TestApp
spec:
  field: value
`,
			wantErr: false,
			wantStages: Stages{
				stages: nil, // Expecting nil slice, which serializes to null in JSON/YAML
			},
		},
		{
			name: "success with pipeline but no stages field",
			inputYAML: `
apiVersion: pipecd.dev/v1beta1
kind: TestApp
spec:
  pipeline:
    field: value
`,
			wantErr: false,
			wantStages: Stages{
				stages: nil, // Expecting nil slice
			},
		},
		{
			name: "success with empty stages list",
			inputYAML: `
apiVersion: pipecd.dev/v1beta1
kind: TestApp
spec:
  pipeline:
    stages: []
`,
			wantErr: false,
			wantStages: Stages{
				stages: nil, // Expecting nil slice
			},
		},
		{
			name: "invalid yaml",
			inputYAML: `
apiVersion: pipecd.dev/v1beta1
kind: TestApp
spec:
  pipeline:
    stages:
      - name: STAGE1
invalid-yaml-token
`,
			wantErr:    true,
			wantStages: Stages{}, // Value doesn't matter on error
		},
		{
			name:       "empty config",
			inputYAML:  ``,
			wantErr:    true, // The yaml decoder will return an error for empty config.
			wantStages: Stages{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ds := DeploymentSource{
				ApplicationConfig: []byte(tc.inputYAML),
			}

			stages, err := LoadStages(ds)

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				// Handle nil vs empty slice explicitly for clarity if needed,
				// but assert.Equal should handle this correctly for struct comparison.
				if tc.wantStages.stages == nil {
					assert.Nil(t, stages.stages)
				} else {
					assert.Equal(t, tc.wantStages, stages)
				}
			}
		})
	}
}

func TestStages_Has(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name   string
		stages Stages
		stage  string
		want   bool
	}{
		{
			name:   "stage exists",
			stages: Stages{stages: []string{"STAGE1", "STAGE2", "STAGE3"}},
			stage:  "STAGE2",
			want:   true,
		},
		{
			name:   "stage does not exist",
			stages: Stages{stages: []string{"STAGE1", "STAGE3"}},
			stage:  "STAGE2",
			want:   false,
		},
		{
			name:   "empty stages list",
			stages: Stages{stages: []string{}},
			stage:  "STAGE1",
			want:   false,
		},
		{
			name:   "nil stages list",
			stages: Stages{stages: nil},
			stage:  "STAGE1",
			want:   false,
		},
		{
			name:   "checking for empty string stage",
			stages: Stages{stages: []string{"STAGE1", "", "STAGE3"}},
			stage:  "",
			want:   true,
		},
		{
			name:   "checking for empty string stage when not present",
			stages: Stages{stages: []string{"STAGE1", "STAGE3"}},
			stage:  "",
			want:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := tc.stages.Has(tc.stage)
			assert.Equal(t, tc.want, got)
		})
	}
}
