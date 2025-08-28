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

package executestage

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/analysis/config"
)

func TestExecutor_buildAppArgs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                 string
		appName              string
		analysisAppSpec      *config.AnalysisApplicationSpec
		customArgs           map[string]string
		expectedArgsTemplate argsTemplate
	}{
		{
			name:    "basic functionality with app name only",
			appName: "test-app",
			analysisAppSpec: &config.AnalysisApplicationSpec{
				AppCustomArgs: map[string]string{},
			},
			customArgs: nil,
			expectedArgsTemplate: argsTemplate{
				App:           appArgs{Name: "test-app"},
				AppCustomArgs: map[string]string{},
				K8s:           map[string]string{"Namespace": ""},
			},
		},
		{
			name:    "with app custom args only",
			appName: "test-app",
			analysisAppSpec: &config.AnalysisApplicationSpec{
				AppCustomArgs: map[string]string{
					"env":          "production",
					"k8sNamespace": "default",
				},
			},
			customArgs: nil,
			expectedArgsTemplate: argsTemplate{
				App: appArgs{Name: "test-app"},
				AppCustomArgs: map[string]string{
					"env":          "production",
					"k8sNamespace": "default",
				},
				K8s: map[string]string{"Namespace": "default"},
			},
		},
		{
			name:    "with stage custom args only",
			appName: "test-app",
			analysisAppSpec: &config.AnalysisApplicationSpec{
				AppCustomArgs: map[string]string{},
			},
			customArgs: map[string]string{
				"region":       "us-west-2",
				"k8sNamespace": "staging",
			},
			expectedArgsTemplate: argsTemplate{
				App: appArgs{Name: "test-app"},
				AppCustomArgs: map[string]string{
					"region":       "us-west-2",
					"k8sNamespace": "staging",
				},
				K8s: map[string]string{"Namespace": "staging"},
			},
		},
		{
			name:    "stage config args override app spec args",
			appName: "test-app",
			analysisAppSpec: &config.AnalysisApplicationSpec{
				AppCustomArgs: map[string]string{
					"env":          "development",
					"k8sNamespace": "default",
					"region":       "us-east-1",
				},
			},
			customArgs: map[string]string{
				"env":          "production",   // should override
				"k8sNamespace": "prod",         // should override
				"cluster":      "prod-cluster", // new key
			},
			expectedArgsTemplate: argsTemplate{
				App: appArgs{Name: "test-app"},
				AppCustomArgs: map[string]string{
					"env":          "production",   // overridden
					"k8sNamespace": "prod",         // overridden
					"region":       "us-east-1",    // preserved from app spec
					"cluster":      "prod-cluster", // added from stage config
				},
				K8s: map[string]string{"Namespace": "prod"}, // uses overridden value
			},
		},
		{
			name:    "k8s namespace backward compatibility without k8sNamespace",
			appName: "test-app",
			analysisAppSpec: &config.AnalysisApplicationSpec{
				AppCustomArgs: map[string]string{
					"env": "production",
				},
			},
			customArgs: map[string]string{
				"region": "us-west-2",
			},
			expectedArgsTemplate: argsTemplate{
				App: appArgs{Name: "test-app"},
				AppCustomArgs: map[string]string{
					"env":    "production",
					"region": "us-west-2",
				},
				K8s: map[string]string{"Namespace": ""}, // empty when k8sNamespace not present
			},
		},
		{
			name:    "empty app name",
			appName: "",
			analysisAppSpec: &config.AnalysisApplicationSpec{
				AppCustomArgs: map[string]string{
					"k8sNamespace": "test-ns",
				},
			},
			customArgs: nil,
			expectedArgsTemplate: argsTemplate{
				App: appArgs{Name: ""},
				AppCustomArgs: map[string]string{
					"k8sNamespace": "test-ns",
				},
				K8s: map[string]string{"Namespace": "test-ns"},
			},
		},
		{
			name:    "both app custom args and custom args are nil",
			appName: "test-app",
			analysisAppSpec: &config.AnalysisApplicationSpec{
				AppCustomArgs: nil,
			},
			customArgs: nil,
			expectedArgsTemplate: argsTemplate{
				App:           appArgs{Name: "test-app"},
				AppCustomArgs: map[string]string{},
				K8s:           map[string]string{"Namespace": ""},
			},
		},
		{
			name:    "app custom args is nil and custom args is not nil",
			appName: "test-app",
			analysisAppSpec: &config.AnalysisApplicationSpec{
				AppCustomArgs: map[string]string{
					"k8sNamespace": "test-ns",
				},
			},
			customArgs: nil,
			expectedArgsTemplate: argsTemplate{
				App: appArgs{Name: "test-app"},
				AppCustomArgs: map[string]string{
					"k8sNamespace": "test-ns",
				},
				K8s: map[string]string{"Namespace": "test-ns"},
			},
		},
		{
			name:    "app custom args is not nil and custom args is nil",
			appName: "test-app",
			analysisAppSpec: &config.AnalysisApplicationSpec{
				AppCustomArgs: nil,
			},
			customArgs: map[string]string{
				"k8sNamespace": "test-ns",
			},
			expectedArgsTemplate: argsTemplate{
				App:           appArgs{Name: "test-app"},
				AppCustomArgs: map[string]string{"k8sNamespace": "test-ns"},
				K8s:           map[string]string{"Namespace": "test-ns"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			e := &executor{
				appName:         tt.appName,
				analysisAppSpec: tt.analysisAppSpec,
			}

			result := e.buildAppArgs(tt.customArgs)

			assert.Equal(t, tt.expectedArgsTemplate.App, result.App)
			assert.Equal(t, tt.expectedArgsTemplate.AppCustomArgs, result.AppCustomArgs)
			assert.Equal(t, tt.expectedArgsTemplate.K8s, result.K8s)
		})
	}
}

func TestExecutor_buildAppArgs_MapsCloning(t *testing.T) {
	t.Parallel()

	// Test that the original maps are not modified
	originalAppCustomArgs := map[string]string{
		"env":          "development",
		"k8sNamespace": "default",
	}
	originalCustomArgs := map[string]string{
		"env": "production",
	}

	e := &executor{
		appName: "test-app",
		analysisAppSpec: &config.AnalysisApplicationSpec{
			AppCustomArgs: originalAppCustomArgs,
		},
	}

	result := e.buildAppArgs(originalCustomArgs)

	// Verify the result has the merged values
	assert.Equal(t, "production", result.AppCustomArgs["env"])       // overridden value
	assert.Equal(t, "default", result.AppCustomArgs["k8sNamespace"]) // preserved value

	// Verify original maps are unchanged
	assert.Equal(t, "development", originalAppCustomArgs["env"])
	assert.Equal(t, "production", originalCustomArgs["env"])

	// Modify the result to ensure it doesn't affect original maps
	result.AppCustomArgs["env"] = "modified"
	assert.Equal(t, "development", originalAppCustomArgs["env"]) // should remain unchanged
	assert.Equal(t, "production", originalCustomArgs["env"])     // should remain unchanged
}
