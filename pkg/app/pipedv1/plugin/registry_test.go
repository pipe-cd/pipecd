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

package plugin

import (
	"testing"

	"github.com/stretchr/testify/assert"

	config "github.com/pipe-cd/pipecd/pkg/configv1"
	pluginapi "github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1"
)

type fakePluginClient struct {
	pluginapi.PluginClient
	name string
}

func TestPluginRegistry_GetPluginClientsByAppConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		cfg      *config.GenericApplicationSpec
		setup    func() *pluginRegistry
		expected []pluginapi.PluginClient
		wantErr  bool
	}{
		{
			name: "get plugins by pipeline",
			cfg: &config.GenericApplicationSpec{
				Pipeline: &config.DeploymentPipeline{
					Stages: []config.PipelineStage{
						{Name: "stage1"},
						{Name: "stage2"},
					},
				},
				Plugins: nil,
			},
			setup: func() *pluginRegistry {
				return &pluginRegistry{
					stageBasedPlugins: map[string]pluginapi.PluginClient{
						"stage1": fakePluginClient{name: "stage1"},
						"stage2": fakePluginClient{name: "stage2"},
					},
				}
			},
			expected: []pluginapi.PluginClient{
				fakePluginClient{name: "stage1"},
				fakePluginClient{name: "stage2"},
			},
			wantErr: false,
		},
		{
			name: "get plugins by plugin names",
			cfg: &config.GenericApplicationSpec{
				Pipeline: nil,
				Plugins:  map[string]struct{}{"plugin1": {}, "plugin2": {}},
			},
			setup: func() *pluginRegistry {
				return &pluginRegistry{
					nameBasedPlugins: map[string]pluginapi.PluginClient{
						"plugin1": fakePluginClient{name: "plugin1"},
						"plugin2": fakePluginClient{name: "plugin2"},
					},
				}
			},
			expected: []pluginapi.PluginClient{
				fakePluginClient{name: "plugin1"},
				fakePluginClient{name: "plugin2"},
			},
			wantErr: false,
		},
		{
			name: "get plugins by pipeline when both pipeline and plugins are specified",
			cfg: &config.GenericApplicationSpec{
				Pipeline: &config.DeploymentPipeline{
					Stages: []config.PipelineStage{
						{Name: "stage1"},
						{Name: "stage2"},
					},
				},
				Plugins: map[string]struct{}{"plugin1": {}, "plugin2": {}},
			},
			setup: func() *pluginRegistry {
				return &pluginRegistry{
					stageBasedPlugins: map[string]pluginapi.PluginClient{
						"stage1": fakePluginClient{name: "stage1"},
						"stage2": fakePluginClient{name: "stage2"},
					},
					nameBasedPlugins: map[string]pluginapi.PluginClient{
						"plugin1": fakePluginClient{name: "plugin1"},
						"plugin2": fakePluginClient{name: "plugin2"},
					},
				}
			},
			expected: []pluginapi.PluginClient{
				fakePluginClient{name: "stage1"},
				fakePluginClient{name: "stage2"},
			},
			wantErr: false,
		},
		{
			name: "no plugins found when no pipeline or plugins specified",
			cfg:  &config.GenericApplicationSpec{},
			setup: func() *pluginRegistry {
				return &pluginRegistry{
					nameBasedPlugins: map[string]pluginapi.PluginClient{},
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			pr := tt.setup()
			plugins, err := pr.GetPluginClientsByAppConfig(tt.cfg)
			assert.ElementsMatch(t, tt.expected, plugins)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
func TestPluginRegistry_getPluginClientsByPipeline(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		pipeline *config.DeploymentPipeline
		setup    func() *pluginRegistry
		expected []pluginapi.PluginClient
		wantErr  bool
	}{
		{
			name: "get plugins by valid pipeline stages",
			pipeline: &config.DeploymentPipeline{
				Stages: []config.PipelineStage{
					{Name: "stage1"},
					{Name: "stage2"},
				},
			},
			setup: func() *pluginRegistry {
				return &pluginRegistry{
					stageBasedPlugins: map[string]pluginapi.PluginClient{
						"stage1": fakePluginClient{name: "stage1"},
						"stage2": fakePluginClient{name: "stage2"},
					},
				}
			},
			expected: []pluginapi.PluginClient{
				fakePluginClient{name: "stage1"},
				fakePluginClient{name: "stage2"},
			},
			wantErr: false,
		},
		{
			name: "no plugins found for empty pipeline stages",
			pipeline: &config.DeploymentPipeline{
				Stages: []config.PipelineStage{},
			},
			setup: func() *pluginRegistry {
				return &pluginRegistry{
					stageBasedPlugins: map[string]pluginapi.PluginClient{},
				}
			},
			expected: nil,
			wantErr:  true,
		},
		{
			name: "no plugins found for non-existent pipeline stages",
			pipeline: &config.DeploymentPipeline{
				Stages: []config.PipelineStage{
					{Name: "stage1"},
					{Name: "stage2"},
				},
			},
			setup: func() *pluginRegistry {
				return &pluginRegistry{
					stageBasedPlugins: map[string]pluginapi.PluginClient{},
				}
			},
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			pr := tt.setup()
			plugins, err := pr.getPluginClientsByPipeline(tt.pipeline)
			assert.ElementsMatch(t, tt.expected, plugins)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
func TestPluginRegistry_getPluginClientsByNames(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		pluginNames map[string]struct{}
		setup       func() *pluginRegistry
		expected    []pluginapi.PluginClient
		wantErr     bool
	}{
		{
			name:        "get plugins by valid plugin names",
			pluginNames: map[string]struct{}{"plugin1": {}, "plugin2": {}},
			setup: func() *pluginRegistry {
				return &pluginRegistry{
					nameBasedPlugins: map[string]pluginapi.PluginClient{
						"plugin1": fakePluginClient{name: "plugin1"},
						"plugin2": fakePluginClient{name: "plugin2"},
					},
				}
			},
			expected: []pluginapi.PluginClient{
				fakePluginClient{name: "plugin1"},
				fakePluginClient{name: "plugin2"},
			},
			wantErr: false,
		},
		{
			name:        "no plugins found for empty plugin names",
			pluginNames: map[string]struct{}{},
			setup: func() *pluginRegistry {
				return &pluginRegistry{
					nameBasedPlugins: map[string]pluginapi.PluginClient{
						"plugin1": fakePluginClient{name: "plugin1"},
						"plugin2": fakePluginClient{name: "plugin2"},
					},
				}
			},
			wantErr: true,
		},
		{
			name:        "no plugins found for non-existent plugin names",
			pluginNames: map[string]struct{}{"plugin1": {}, "plugin2": {}},
			setup: func() *pluginRegistry {
				return &pluginRegistry{
					nameBasedPlugins: map[string]pluginapi.PluginClient{
						"plugin3": fakePluginClient{name: "plugin3"},
						"plugin4": fakePluginClient{name: "plugin4"},
					},
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			pr := tt.setup()
			plugins, err := pr.getPluginClientsByNames(tt.pluginNames)
			assert.ElementsMatch(t, tt.expected, plugins)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
func TestPluginRegistry_GetPluginClientByStageName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		stage    string
		setup    func() *pluginRegistry
		expected pluginapi.PluginClient
		wantErr  bool
	}{
		{
			name:  "get plugin by valid stage name",
			stage: "stage1",
			setup: func() *pluginRegistry {
				return &pluginRegistry{
					stageBasedPlugins: map[string]pluginapi.PluginClient{
						"stage1": fakePluginClient{name: "stage1"},
					},
				}
			},
			expected: fakePluginClient{name: "stage1"},
			wantErr:  false,
		},
		{
			name:  "no plugin found for non-existent stage name",
			stage: "stage2",
			setup: func() *pluginRegistry {
				return &pluginRegistry{
					stageBasedPlugins: map[string]pluginapi.PluginClient{
						"stage1": fakePluginClient{name: "stage1"},
					},
				}
			},
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			pr := tt.setup()
			plugin, err := pr.GetPluginClientByStageName(tt.stage)
			assert.Equal(t, tt.expected, plugin)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
