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

package migrate

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"sigs.k8s.io/yaml"
)

func TestApplicationConfig_migrateApplicationConfig(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		inputConfig map[string]interface{}
		expected    map[string]interface{}
		expectError bool
	}{
		{
			name: "kubernetes application migration",
			inputConfig: map[string]interface{}{
				"kind":       "KubernetesApp",
				"apiVersion": "pipecd.dev/v1beta1",
				"spec": map[string]interface{}{
					"name":        "test-app",
					"description": "Test application",
					"labels": map[string]interface{}{
						"env": "test",
					},
					"input": map[string]interface{}{
						"namespace": "default",
					},
					"quickSync": map[string]interface{}{
						"prune": true,
					},
					"service": map[string]interface{}{
						"name": "test-service",
					},
					"pipeline": map[string]interface{}{
						"stages": []interface{}{
							map[string]interface{}{
								"id":   "stage1",
								"name": "Stage 1",
								"with": map[string]interface{}{
									"timeout": "10m",
									"skipOn":  "failure",
								},
							},
						},
					},
				},
			},
			expected: map[string]interface{}{
				"kind":       "Application",
				"apiVersion": "pipecd.dev/v1beta1",
				"spec": map[string]interface{}{
					"name":        "test-app",
					"description": "Test application",
					"labels": map[string]interface{}{
						"env": "test",
					},
					"plugins": map[string]interface{}{
						"kubernetes": map[string]interface{}{
							"input": map[string]interface{}{
								"namespace": "default",
							},
							"quickSync": map[string]interface{}{
								"prune": true,
							},
							"service": map[string]interface{}{
								"name": "test-service",
							},
						},
					},
					"pipeline": map[string]interface{}{
						"stages": []interface{}{
							map[string]interface{}{
								"id":      "stage1",
								"name":    "Stage 1",
								"timeout": "10m",
								"skipOn":  "failure",
								"with": map[string]interface{}{
									"timeout": "10m",
									"skipOn":  "failure",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "terraform application migration",
			inputConfig: map[string]interface{}{
				"kind":       "TerraformApp",
				"apiVersion": "pipecd.dev/v1beta1",
				"spec": map[string]interface{}{
					"name":        "terraform-app",
					"description": "Terraform application",
					"input": map[string]interface{}{
						"workspace": "default",
					},
					"quickSync": map[string]interface{}{
						"prune": true,
					},
				},
			},
			expected: map[string]interface{}{
				"kind":       "Application",
				"apiVersion": "pipecd.dev/v1beta1",
				"spec": map[string]interface{}{
					"name":        "terraform-app",
					"description": "Terraform application",
					"plugins": map[string]interface{}{
						"terraform": map[string]interface{}{
							"input": map[string]interface{}{
								"workspace": "default",
							},
							"quickSync": map[string]interface{}{
								"prune": true,
							},
						},
					},
				},
			},
		},
		{
			name: "ecs application migration",
			inputConfig: map[string]interface{}{
				"kind":       "ECSApp",
				"apiVersion": "pipecd.dev/v1beta1",
				"spec": map[string]interface{}{
					"name":        "ecs-app",
					"description": "ECS application",
					"input": map[string]interface{}{
						"cluster": "test-cluster",
					},
					"quickSync": map[string]interface{}{
						"prune": true,
					},
				},
			},
			expected: map[string]interface{}{
				"kind":       "Application",
				"apiVersion": "pipecd.dev/v1beta1",
				"spec": map[string]interface{}{
					"name":        "ecs-app",
					"description": "ECS application",
					"plugins": map[string]interface{}{
						"ecs": map[string]interface{}{
							"input": map[string]interface{}{
								"cluster": "test-cluster",
							},
							"quickSync": map[string]interface{}{
								"prune": true,
							},
						},
					},
				},
			},
		},
		{
			name: "lambda application migration",
			inputConfig: map[string]interface{}{
				"kind":       "LambdaApp",
				"apiVersion": "pipecd.dev/v1beta1",
				"spec": map[string]interface{}{
					"name":        "lambda-app",
					"description": "Lambda application",
					"input": map[string]interface{}{
						"region": "us-west-2",
					},
					"quickSync": map[string]interface{}{
						"prune": true,
					},
				},
			},
			expected: map[string]interface{}{
				"kind":       "Application",
				"apiVersion": "pipecd.dev/v1beta1",
				"spec": map[string]interface{}{
					"name":        "lambda-app",
					"description": "Lambda application",
					"plugins": map[string]interface{}{
						"lambda": map[string]interface{}{
							"input": map[string]interface{}{
								"region": "us-west-2",
							},
							"quickSync": map[string]interface{}{
								"prune": true,
							},
						},
					},
				},
			},
		},
		{
			name: "cloudrun application migration",
			inputConfig: map[string]interface{}{
				"kind":       "CloudRunApp",
				"apiVersion": "pipecd.dev/v1beta1",
				"spec": map[string]interface{}{
					"name":        "cloudrun-app",
					"description": "Cloud Run application",
					"input": map[string]interface{}{
						"project": "test-project",
					},
					"quickSync": map[string]interface{}{
						"prune": true,
					},
				},
			},
			expected: map[string]interface{}{
				"kind":       "Application",
				"apiVersion": "pipecd.dev/v1beta1",
				"spec": map[string]interface{}{
					"name":        "cloudrun-app",
					"description": "Cloud Run application",
					"plugins": map[string]interface{}{
						"cloudrun": map[string]interface{}{
							"input": map[string]interface{}{
								"project": "test-project",
							},
							"quickSync": map[string]interface{}{
								"prune": true,
							},
						},
					},
				},
			},
		},
		{
			name: "unsupported application kind",
			inputConfig: map[string]interface{}{
				"kind":       "Application",
				"apiVersion": "pipecd.dev/v1beta1",
				"spec": map[string]interface{}{
					"name": "test-app",
				},
			},
			expectError: true,
		},
		{
			name: "pipeline with timeout and skipOn migration",
			inputConfig: map[string]interface{}{
				"kind":       "KubernetesApp",
				"apiVersion": "pipecd.dev/v1beta1",
				"spec": map[string]interface{}{
					"name": "test-app",
					"pipeline": map[string]interface{}{
						"stages": []interface{}{
							map[string]interface{}{
								"id":   "stage1",
								"name": "Stage 1",
								"with": map[string]interface{}{
									"timeout": "5m",
									"skipOn":  "success",
								},
							},
							map[string]interface{}{
								"id":   "stage2",
								"name": "Stage 2",
								"with": map[string]interface{}{
									"timeout": "15m",
								},
							},
							map[string]interface{}{
								"id":   "stage3",
								"name": "Stage 3",
								"with": map[string]interface{}{
									"skipOn": "failure",
								},
							},
						},
					},
				},
			},
			expected: map[string]interface{}{
				"kind":       "Application",
				"apiVersion": "pipecd.dev/v1beta1",
				"spec": map[string]interface{}{
					"name": "test-app",
					"plugins": map[string]interface{}{
						"kubernetes": map[string]interface{}{},
					},
					"pipeline": map[string]interface{}{
						"stages": []interface{}{
							map[string]interface{}{
								"id":      "stage1",
								"name":    "Stage 1",
								"timeout": "5m",
								"skipOn":  "success",
								"with": map[string]interface{}{
									"timeout": "5m",
									"skipOn":  "success",
								},
							},
							map[string]interface{}{
								"id":      "stage2",
								"name":    "Stage 2",
								"timeout": "15m",
								"with": map[string]interface{}{
									"timeout": "15m",
								},
							},
							map[string]interface{}{
								"id":     "stage3",
								"name":   "Stage 3",
								"skipOn": "failure",
								"with": map[string]interface{}{
									"skipOn": "failure",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "KubernetesApp with analysis stage",
			inputConfig: map[string]interface{}{
				"kind":       "KubernetesApp",
				"apiVersion": "pipecd.dev/v1beta1",
				"spec": map[string]interface{}{
					"name":        "test-app",
					"description": "Test application",
					"labels": map[string]interface{}{
						"env": "test",
					},
					"input": map[string]interface{}{
						"namespace": "test",
					},
					"quickSync": map[string]interface{}{
						"prune": true,
					},
					"service": map[string]interface{}{
						"name": "test-service",
					},
					"pipeline": map[string]interface{}{
						"stages": []interface{}{
							map[string]interface{}{
								"id":   "stage1",
								"name": "ANALYSIS",
							},
						},
					},
				},
			},
			expected: map[string]interface{}{
				"kind":       "Application",
				"apiVersion": "pipecd.dev/v1beta1",
				"spec": map[string]interface{}{
					"name":        "test-app",
					"description": "Test application",
					"labels": map[string]interface{}{
						"env": "test",
					},
					"plugins": map[string]interface{}{
						"kubernetes": map[string]interface{}{
							"input": map[string]interface{}{
								"namespace": "test",
							},
							"quickSync": map[string]interface{}{
								"prune": true,
							},
							"service": map[string]interface{}{
								"name": "test-service",
							},
						},
						"analysis": map[string]interface{}{
							"appCustomArgs": map[string]interface{}{
								"k8sNamespace": "test",
							},
						},
					},
					"pipeline": map[string]interface{}{
						"stages": []interface{}{
							map[string]interface{}{
								"id":   "stage1",
								"name": "ANALYSIS",
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Create temporary directory for test
			tempDir := t.TempDir()
			configFile := filepath.Join(tempDir, "app.yaml")

			// Write input config to file
			inputData, err := yaml.Marshal(tt.inputConfig)
			require.NoError(t, err)
			err = os.WriteFile(configFile, inputData, 0644)
			require.NoError(t, err)

			// Create application config instance
			appConfig := &applicationConfig{}

			// Run migration
			logger := zap.NewNop()
			err = appConfig.migrateApplicationConfig(t.Context(), configFile, logger)

			if tt.expectError {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)

			// Read migrated config
			migratedData, err := os.ReadFile(configFile)
			require.NoError(t, err)

			var migratedConfig map[string]interface{}
			err = yaml.Unmarshal(migratedData, &migratedConfig)
			require.NoError(t, err)

			// Verify migration
			assert.Equal(t, tt.expected["kind"], migratedConfig["kind"])
			assert.Equal(t, tt.expected["apiVersion"], migratedConfig["apiVersion"])

			// Verify spec structure
			expectedSpec := tt.expected["spec"].(map[string]interface{})
			migratedSpec := migratedConfig["spec"].(map[string]interface{})

			// Check generic fields
			for _, key := range []string{"name", "description", "labels"} {
				if expectedValue, exists := expectedSpec[key]; exists {
					assert.Equal(t, expectedValue, migratedSpec[key])
				}
			}

			// Check plugins structure
			if expectedPlugins, exists := expectedSpec["plugins"]; exists {
				assert.Equal(t, expectedPlugins, migratedSpec["plugins"])
			}

			// Check pipeline structure
			if expectedPipeline, exists := expectedSpec["pipeline"]; exists {
				assert.Equal(t, expectedPipeline, migratedSpec["pipeline"])
			}

			// Verify backup file was created
			backupFile := configFile + ".old"
			_, err = os.Stat(backupFile)
			assert.NoError(t, err, "Backup file should exist")
		})
	}
}

func TestApplicationConfig_migrateApplicationConfig_WriteErrors(t *testing.T) {
	t.Parallel()
	t.Run("write permission denied", func(t *testing.T) {
		tempDir := t.TempDir()
		configFile := filepath.Join(tempDir, "app.yaml")

		// Create a valid config file
		config := map[string]interface{}{
			"kind":       "KubernetesApp",
			"apiVersion": "pipecd.dev/v1beta1",
			"spec": map[string]interface{}{
				"name": "test-app",
			},
		}
		data, _ := yaml.Marshal(config)
		err := os.WriteFile(configFile, data, 0644)
		require.NoError(t, err)

		// Make directory read-only to simulate write permission issues
		err = os.Chmod(tempDir, 0444)
		require.NoError(t, err)
		defer os.Chmod(tempDir, 0755) // Restore permissions

		appConfig := &applicationConfig{}
		logger := zap.NewNop()
		err = appConfig.migrateApplicationConfig(context.Background(), configFile, logger)

		assert.Error(t, err)
	})
}

func TestApplicationConfig_migrateApplicationConfig_GenericFields(t *testing.T) {
	t.Parallel()
	inputConfig := map[string]interface{}{
		"kind":       "KubernetesApp",
		"apiVersion": "pipecd.dev/v1beta1",
		"spec": map[string]interface{}{
			"name":        "test-app",
			"labels":      map[string]interface{}{"env": "test", "team": "backend"},
			"description": "Test application with all generic fields",
			"planner": map[string]interface{}{
				"alwaysUsePipeline": true,
			},
			"commitMatcher": map[string]interface{}{
				"branches": []interface{}{"main", "develop"},
			},
			"trigger": map[string]interface{}{
				"onCommit": map[string]interface{}{
					"branches": []interface{}{"main"},
				},
			},
			"postSync": map[string]interface{}{
				"analysis": map[string]interface{}{
					"enabled": true,
				},
			},
			"timeout": "30m",
			"encryption": map[string]interface{}{
				"enabled": true,
			},
			"attachment": map[string]interface{}{
				"enabled": true,
			},
			"notification": map[string]interface{}{
				"slack": map[string]interface{}{
					"enabled": true,
				},
			},
			"eventWatcher": map[string]interface{}{
				"enabled": true,
			},
			"driftDetection": map[string]interface{}{
				"enabled": true,
			},
			"input": map[string]interface{}{
				"namespace": "default",
			},
		},
	}

	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "app.yaml")

	// Write input config to file
	inputData, err := yaml.Marshal(inputConfig)
	require.NoError(t, err)
	err = os.WriteFile(configFile, inputData, 0644)
	require.NoError(t, err)

	// Create application config instance
	appConfig := &applicationConfig{}

	// Run migration
	logger := zap.NewNop()
	err = appConfig.migrateApplicationConfig(context.Background(), configFile, logger)
	require.NoError(t, err)

	// Read migrated config
	migratedData, err := os.ReadFile(configFile)
	require.NoError(t, err)

	var migratedConfig map[string]interface{}
	err = yaml.Unmarshal(migratedData, &migratedConfig)
	require.NoError(t, err)

	// Verify all generic fields are preserved
	migratedSpec := migratedConfig["spec"].(map[string]interface{})
	expectedFields := []string{
		"name", "labels", "description", "planner", "commitMatcher",
		"trigger", "postSync", "timeout", "encryption", "attachment",
		"notification", "eventWatcher", "driftDetection",
	}

	for _, field := range expectedFields {
		if inputConfig["spec"].(map[string]interface{})[field] != nil {
			assert.Equal(t, inputConfig["spec"].(map[string]interface{})[field], migratedSpec[field])
		}
	}

	// Verify plugin-specific fields are moved to plugins.kubernetes
	assert.NotNil(t, migratedSpec["plugins"])
	plugins := migratedSpec["plugins"].(map[string]interface{})
	assert.NotNil(t, plugins["kubernetes"])
	kubernetesPlugin := plugins["kubernetes"].(map[string]interface{})
	assert.Equal(t, inputConfig["spec"].(map[string]interface{})["input"], kubernetesPlugin["input"])
}
