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
		inputConfig map[string]any
		expected    map[string]any
		expectError bool
	}{
		{
			name: "kubernetes application migration",
			inputConfig: map[string]any{
				"kind":       "KubernetesApp",
				"apiVersion": "pipecd.dev/v1beta1",
				"spec": map[string]any{
					"name":        "test-app",
					"description": "Test application",
					"labels": map[string]any{
						"env": "test",
					},
					"input": map[string]any{
						"namespace": "default",
					},
					"quickSync": map[string]any{
						"prune": true,
					},
					"service": map[string]any{
						"name": "test-service",
					},
					"pipeline": map[string]any{
						"stages": []any{
							map[string]any{
								"id":   "stage1",
								"name": "Stage 1",
								"with": map[string]any{
									"timeout": "10m",
									"skipOn":  "failure",
								},
							},
						},
					},
				},
			},
			expected: map[string]any{
				"kind":       "Application",
				"apiVersion": "pipecd.dev/v1beta1",
				"spec": map[string]any{
					"name":        "test-app",
					"description": "Test application",
					"labels": map[string]any{
						"env": "test",
					},
					"plugins": map[string]any{
						"kubernetes": map[string]any{
							"input": map[string]any{
								"namespace": "default",
							},
							"quickSync": map[string]any{
								"prune": true,
							},
							"service": map[string]any{
								"name": "test-service",
							},
						},
					},
					"pipeline": map[string]any{
						"stages": []any{
							map[string]any{
								"id":      "stage1",
								"name":    "Stage 1",
								"timeout": "10m",
								"skipOn":  "failure",
								"with": map[string]any{
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
			inputConfig: map[string]any{
				"kind":       "TerraformApp",
				"apiVersion": "pipecd.dev/v1beta1",
				"spec": map[string]any{
					"name":        "terraform-app",
					"description": "Terraform application",
					"input": map[string]any{
						"workspace": "default",
					},
					"quickSync": map[string]any{
						"prune": true,
					},
				},
			},
			expected: map[string]any{
				"kind":       "Application",
				"apiVersion": "pipecd.dev/v1beta1",
				"spec": map[string]any{
					"name":        "terraform-app",
					"description": "Terraform application",
					"plugins": map[string]any{
						"terraform": map[string]any{
							"input": map[string]any{
								"workspace": "default",
							},
							"quickSync": map[string]any{
								"prune": true,
							},
						},
					},
				},
			},
		},
		{
			name: "ecs application migration",
			inputConfig: map[string]any{
				"kind":       "ECSApp",
				"apiVersion": "pipecd.dev/v1beta1",
				"spec": map[string]any{
					"name":        "ecs-app",
					"description": "ECS application",
					"input": map[string]any{
						"cluster": "test-cluster",
					},
					"quickSync": map[string]any{
						"prune": true,
					},
				},
			},
			expected: map[string]any{
				"kind":       "Application",
				"apiVersion": "pipecd.dev/v1beta1",
				"spec": map[string]any{
					"name":        "ecs-app",
					"description": "ECS application",
					"plugins": map[string]any{
						"ecs": map[string]any{
							"input": map[string]any{
								"cluster": "test-cluster",
							},
							"quickSync": map[string]any{
								"prune": true,
							},
						},
					},
				},
			},
		},
		{
			name: "lambda application migration",
			inputConfig: map[string]any{
				"kind":       "LambdaApp",
				"apiVersion": "pipecd.dev/v1beta1",
				"spec": map[string]any{
					"name":        "lambda-app",
					"description": "Lambda application",
					"input": map[string]any{
						"region": "us-west-2",
					},
					"quickSync": map[string]any{
						"prune": true,
					},
				},
			},
			expected: map[string]any{
				"kind":       "Application",
				"apiVersion": "pipecd.dev/v1beta1",
				"spec": map[string]any{
					"name":        "lambda-app",
					"description": "Lambda application",
					"plugins": map[string]any{
						"lambda": map[string]any{
							"input": map[string]any{
								"region": "us-west-2",
							},
							"quickSync": map[string]any{
								"prune": true,
							},
						},
					},
				},
			},
		},
		{
			name: "cloudrun application migration",
			inputConfig: map[string]any{
				"kind":       "CloudRunApp",
				"apiVersion": "pipecd.dev/v1beta1",
				"spec": map[string]any{
					"name":        "cloudrun-app",
					"description": "Cloud Run application",
					"input": map[string]any{
						"project": "test-project",
					},
					"quickSync": map[string]any{
						"prune": true,
					},
				},
			},
			expected: map[string]any{
				"kind":       "Application",
				"apiVersion": "pipecd.dev/v1beta1",
				"spec": map[string]any{
					"name":        "cloudrun-app",
					"description": "Cloud Run application",
					"plugins": map[string]any{
						"cloudrun": map[string]any{
							"input": map[string]any{
								"project": "test-project",
							},
							"quickSync": map[string]any{
								"prune": true,
							},
						},
					},
				},
			},
		},
		{
			name: "unsupported application kind",
			inputConfig: map[string]any{
				"kind":       "Application",
				"apiVersion": "pipecd.dev/v1beta1",
				"spec": map[string]any{
					"name": "test-app",
				},
			},
			expectError: true,
		},
		{
			name: "pipeline with timeout and skipOn migration",
			inputConfig: map[string]any{
				"kind":       "KubernetesApp",
				"apiVersion": "pipecd.dev/v1beta1",
				"spec": map[string]any{
					"name": "test-app",
					"pipeline": map[string]any{
						"stages": []any{
							map[string]any{
								"id":   "stage1",
								"name": "Stage 1",
								"with": map[string]any{
									"timeout": "5m",
									"skipOn":  "success",
								},
							},
							map[string]any{
								"id":   "stage2",
								"name": "Stage 2",
								"with": map[string]any{
									"timeout": "15m",
								},
							},
							map[string]any{
								"id":   "stage3",
								"name": "Stage 3",
								"with": map[string]any{
									"skipOn": "failure",
								},
							},
						},
					},
				},
			},
			expected: map[string]any{
				"kind":       "Application",
				"apiVersion": "pipecd.dev/v1beta1",
				"spec": map[string]any{
					"name": "test-app",
					"plugins": map[string]any{
						"kubernetes": map[string]any{},
					},
					"pipeline": map[string]any{
						"stages": []any{
							map[string]any{
								"id":      "stage1",
								"name":    "Stage 1",
								"timeout": "5m",
								"skipOn":  "success",
								"with": map[string]any{
									"timeout": "5m",
									"skipOn":  "success",
								},
							},
							map[string]any{
								"id":      "stage2",
								"name":    "Stage 2",
								"timeout": "15m",
								"with": map[string]any{
									"timeout": "15m",
								},
							},
							map[string]any{
								"id":     "stage3",
								"name":   "Stage 3",
								"skipOn": "failure",
								"with": map[string]any{
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
			inputConfig: map[string]any{
				"kind":       "KubernetesApp",
				"apiVersion": "pipecd.dev/v1beta1",
				"spec": map[string]any{
					"name":        "test-app",
					"description": "Test application",
					"labels": map[string]any{
						"env": "test",
					},
					"input": map[string]any{
						"namespace": "test",
					},
					"quickSync": map[string]any{
						"prune": true,
					},
					"service": map[string]any{
						"name": "test-service",
					},
					"pipeline": map[string]any{
						"stages": []any{
							map[string]any{
								"id":   "stage1",
								"name": "ANALYSIS",
							},
						},
					},
				},
			},
			expected: map[string]any{
				"kind":       "Application",
				"apiVersion": "pipecd.dev/v1beta1",
				"spec": map[string]any{
					"name":        "test-app",
					"description": "Test application",
					"labels": map[string]any{
						"env": "test",
					},
					"plugins": map[string]any{
						"kubernetes": map[string]any{
							"input": map[string]any{
								"namespace": "test",
							},
							"quickSync": map[string]any{
								"prune": true,
							},
							"service": map[string]any{
								"name": "test-service",
							},
						},
						"analysis": map[string]any{
							"appCustomArgs": map[string]any{
								"k8sNamespace": "test",
							},
						},
					},
					"pipeline": map[string]any{
						"stages": []any{
							map[string]any{
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

			var migratedConfig map[string]any
			err = yaml.Unmarshal(migratedData, &migratedConfig)
			require.NoError(t, err)

			// Verify migration
			assert.Equal(t, tt.expected["kind"], migratedConfig["kind"])
			assert.Equal(t, tt.expected["apiVersion"], migratedConfig["apiVersion"])

			// Verify spec structure
			expectedSpec := tt.expected["spec"].(map[string]any)
			migratedSpec := migratedConfig["spec"].(map[string]any)

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
		config := map[string]any{
			"kind":       "KubernetesApp",
			"apiVersion": "pipecd.dev/v1beta1",
			"spec": map[string]any{
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
	inputConfig := map[string]any{
		"kind":       "KubernetesApp",
		"apiVersion": "pipecd.dev/v1beta1",
		"spec": map[string]any{
			"name":        "test-app",
			"labels":      map[string]any{"env": "test", "team": "backend"},
			"description": "Test application with all generic fields",
			"planner": map[string]any{
				"alwaysUsePipeline": true,
			},
			"commitMatcher": map[string]any{
				"branches": []any{"main", "develop"},
			},
			"trigger": map[string]any{
				"onCommit": map[string]any{
					"branches": []any{"main"},
				},
			},
			"postSync": map[string]any{
				"analysis": map[string]any{
					"enabled": true,
				},
			},
			"timeout": "30m",
			"encryption": map[string]any{
				"enabled": true,
			},
			"attachment": map[string]any{
				"enabled": true,
			},
			"notification": map[string]any{
				"slack": map[string]any{
					"enabled": true,
				},
			},
			"eventWatcher": map[string]any{
				"enabled": true,
			},
			"driftDetection": map[string]any{
				"enabled": true,
			},
			"input": map[string]any{
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

	var migratedConfig map[string]any
	err = yaml.Unmarshal(migratedData, &migratedConfig)
	require.NoError(t, err)

	// Verify all generic fields are preserved
	migratedSpec := migratedConfig["spec"].(map[string]any)
	expectedFields := []string{
		"name", "labels", "description", "planner", "commitMatcher",
		"trigger", "postSync", "timeout", "encryption", "attachment",
		"notification", "eventWatcher", "driftDetection",
	}

	for _, field := range expectedFields {
		if inputConfig["spec"].(map[string]any)[field] != nil {
			assert.Equal(t, inputConfig["spec"].(map[string]any)[field], migratedSpec[field])
		}
	}

	// Verify plugin-specific fields are moved to plugins.kubernetes
	assert.NotNil(t, migratedSpec["plugins"])
	plugins := migratedSpec["plugins"].(map[string]any)
	assert.NotNil(t, plugins["kubernetes"])
	kubernetesPlugin := plugins["kubernetes"].(map[string]any)
	assert.Equal(t, inputConfig["spec"].(map[string]any)["input"], kubernetesPlugin["input"])
}
