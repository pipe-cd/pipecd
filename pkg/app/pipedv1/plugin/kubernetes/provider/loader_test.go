// Copyright 2024 The PipeCD Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package provider

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/config"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/toolregistry"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/toolregistry/toolregistrytest"
)

func mustParseManifests(t *testing.T, data string) []Manifest {
	t.Helper()

	manifests, err := ParseManifests(data)
	require.NoError(t, err)

	return manifests
}

func TestParseManifests(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		want    []Manifest
		wantErr bool
	}{
		{
			name: "single manifest",
			data: `
apiVersion: v1
kind: ConfigMap
metadata:
  name: test-config
`,
			want: []Manifest{
				{
					Key: ResourceKey{
						APIVersion: "v1",
						Kind:       "ConfigMap",
						Name:       "test-config",
					},
					Body: &unstructured.Unstructured{
						Object: map[string]interface{}{
							"apiVersion": "v1",
							"kind":       "ConfigMap",
							"metadata": map[string]interface{}{
								"name": "test-config",
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "multiple manifests",
			data: `
apiVersion: v1
kind: ConfigMap
metadata:
  name: test-config
---
apiVersion: v1
kind: Service
metadata:
  name: test-service
`,
			want: []Manifest{
				{
					Key: ResourceKey{
						APIVersion: "v1",
						Kind:       "ConfigMap",
						Name:       "test-config",
					},
					Body: &unstructured.Unstructured{
						Object: map[string]interface{}{
							"apiVersion": "v1",
							"kind":       "ConfigMap",
							"metadata": map[string]interface{}{
								"name": "test-config",
							},
						},
					},
				},
				{
					Key: ResourceKey{
						APIVersion: "v1",
						Kind:       "Service",
						Name:       "test-service",
					},
					Body: &unstructured.Unstructured{
						Object: map[string]interface{}{
							"apiVersion": "v1",
							"kind":       "Service",
							"metadata": map[string]interface{}{
								"name": "test-service",
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid manifest",
			data: `
apiVersion: v1
kind: ConfigMap
metadata:
  name: test-config
---
invalid yaml
`,
			want:    nil,
			wantErr: true,
		},
		{
			name: "empty manifest",
			data: `
---
`,
			want:    []Manifest{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseManifests(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseManifests() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseManifests() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoadPlainYAMLManifests(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		dir            string
		names          []string
		configFilename string
		setup          func(dir string) error
		want           []Manifest
		wantErr        bool
	}{
		{
			name:           "load single manifest",
			dir:            "testdata/single",
			names:          []string{"configmap.yaml"},
			configFilename: "pipecd-config.yaml",
			setup: func(dir string) error {
				return os.WriteFile(filepath.Join(dir, "configmap.yaml"), []byte(`
apiVersion: v1
kind: ConfigMap
metadata:
  name: test-config
`), 0644)
			},
			want: []Manifest{
				{
					Key: ResourceKey{
						APIVersion: "v1",
						Kind:       "ConfigMap",
						Name:       "test-config",
					},
					Body: &unstructured.Unstructured{
						Object: map[string]interface{}{
							"apiVersion": "v1",
							"kind":       "ConfigMap",
							"metadata": map[string]interface{}{
								"name": "test-config",
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:           "ignore config file",
			dir:            "testdata/ignore-config",
			names:          []string{},
			configFilename: "pipecd-config.yaml",
			setup: func(dir string) error {
				// Place dummy files to ensure the loader ignores them.
				if err := os.WriteFile(filepath.Join(dir, "pipecd-config.yaml"), []byte(`
apiVersion: v1
kind: ConfigMap
metadata:
  name: pipecd-config
`), 0644); err != nil {
					return err
				}
				return os.WriteFile(filepath.Join(dir, "service.yaml"), []byte(`
apiVersion: v1
kind: Service
metadata:
  name: test-service
`), 0644)
			},
			want: []Manifest{
				{
					Key: ResourceKey{
						APIVersion: "v1",
						Kind:       "Service",
						Name:       "test-service",
					},
					Body: &unstructured.Unstructured{
						Object: map[string]interface{}{
							"apiVersion": "v1",
							"kind":       "Service",
							"metadata": map[string]interface{}{
								"name": "test-service",
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:           "load multiple manifests",
			dir:            "testdata/multiple",
			names:          []string{"configmap.yaml", "service.yaml"},
			configFilename: "pipecd-config.yaml",
			setup: func(dir string) error {
				if err := os.WriteFile(filepath.Join(dir, "configmap.yaml"), []byte(`
apiVersion: v1
kind: ConfigMap
metadata:
  name: test-config
`), 0644); err != nil {
					return err
				}
				return os.WriteFile(filepath.Join(dir, "service.yaml"), []byte(`
apiVersion: v1
kind: Service
metadata:
  name: test-service
`), 0644)
			},
			want: []Manifest{
				{
					Key: ResourceKey{
						APIVersion: "v1",
						Kind:       "ConfigMap",
						Name:       "test-config",
					},
					Body: &unstructured.Unstructured{
						Object: map[string]interface{}{
							"apiVersion": "v1",
							"kind":       "ConfigMap",
							"metadata": map[string]interface{}{
								"name": "test-config",
							},
						},
					},
				},
				{
					Key: ResourceKey{
						APIVersion: "v1",
						Kind:       "Service",
						Name:       "test-service",
					},
					Body: &unstructured.Unstructured{
						Object: map[string]interface{}{
							"apiVersion": "v1",
							"kind":       "Service",
							"metadata": map[string]interface{}{
								"name": "test-service",
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:           "invalid manifest",
			dir:            "testdata/invalid",
			names:          []string{"invalid.yaml"},
			configFilename: "pipecd-config.yaml",
			setup: func(dir string) error {
				return os.WriteFile(filepath.Join(dir, "invalid.yaml"), []byte(`
invalid yaml content
`), 0644)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:           "no manifests",
			dir:            "testdata/empty",
			names:          []string{},
			configFilename: "pipecd-config.yaml",
			setup: func(dir string) error {
				return nil
			},
			want:    []Manifest{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			dir := filepath.Join(t.TempDir(), tt.dir)
			require.NoError(t, os.MkdirAll(dir, 0755))

			if tt.setup != nil {
				require.NoError(t, tt.setup(dir))
			}

			got, err := LoadPlainYAMLManifests(dir, tt.names, tt.configFilename)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			assert.ElementsMatch(t, tt.want, got)
		})
	}
}

func TestLoader_templateHelmChart(t *testing.T) {
	c, err := toolregistrytest.NewToolRegistry(t)
	require.NoError(t, err)
	t.Cleanup(func() { c.Close() })

	loader := &Loader{
		toolRegistry: toolregistry.NewRegistry(c),
	}

	tests := []struct {
		name    string
		input   LoaderInput
		wantErr bool
	}{
		{
			name: "local chart",
			input: LoaderInput{
				AppName:     "test-app",
				AppDir:      "testdata/testhelm/appconfdir",
				Namespace:   "default",
				HelmVersion: "3.16.1",
				HelmChart:   &config.InputHelmChart{Path: "../../testchart"},
				HelmOptions: &config.InputHelmOptions{},
				Logger:      zap.NewNop(),
			},
			wantErr: false,
		},
		{
			name: "helm chart from git remote",
			input: LoaderInput{
				AppName:     "test-app",
				AppDir:      "testdata/testhelm/appconfdir",
				Namespace:   "default",
				HelmVersion: "3.16.1",
				HelmChart:   &config.InputHelmChart{GitRemote: "https://github.com/test/repo.git"},
				HelmOptions: &config.InputHelmOptions{},
				Logger:      zap.NewNop(),
			},
			wantErr: true, // it's not implemented yet
		},
		{
			name: "helm chart from repository",
			input: LoaderInput{
				AppName:     "test-app",
				AppDir:      "testdata/testhelm/appconfdir",
				Namespace:   "default",
				HelmVersion: "3.16.1",
				HelmChart:   &config.InputHelmChart{Repository: "https://charts.helm.sh/stable"},
				HelmOptions: &config.InputHelmOptions{},
				Logger:      zap.NewNop(),
			},
			wantErr: true, // it's not implemented yet
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := loader.templateHelmChart(context.Background(), tt.input)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
