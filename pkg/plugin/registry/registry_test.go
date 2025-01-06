package registry

import (
	"testing"

	config "github.com/pipe-cd/pipecd/pkg/configv1"
	pluginapi "github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1"
	"github.com/stretchr/testify/assert"
)

type mockPluginClient struct {
	pluginapi.PluginClient
	name string
}

func createMockPluginClient(name string) pluginapi.PluginClient {
	return &mockPluginClient{
		name: name,
	}
}
func TestPluginDeterminer_GetPluginClientsByAppConfig(t *testing.T) {
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
						"stage1": createMockPluginClient("stage1"),
						"stage2": createMockPluginClient("stage2"),
					},
				}
			},
			expected: []pluginapi.PluginClient{
				createMockPluginClient("stage1"),
				createMockPluginClient("stage2"),
			},
			wantErr: false,
		},
		{
			name: "get plugins by plugin names",
			cfg: &config.GenericApplicationSpec{
				Pipeline: nil,
				Plugins:  []string{"plugin1", "plugin2"},
			},
			setup: func() *pluginRegistry {
				return &pluginRegistry{
					nameBasedPlugins: map[string]pluginapi.PluginClient{
						"plugin1": createMockPluginClient("plugin1"),
						"plugin2": createMockPluginClient("plugin2"),
					},
				}
			},
			expected: []pluginapi.PluginClient{
				createMockPluginClient("plugin1"),
				createMockPluginClient("plugin2"),
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
				Plugins: []string{"plugin1", "plugin2"},
			},
			setup: func() *pluginRegistry {
				return &pluginRegistry{
					stageBasedPlugins: map[string]pluginapi.PluginClient{
						"stage1": createMockPluginClient("stage1"),
						"stage2": createMockPluginClient("stage2"),
					},
					nameBasedPlugins: map[string]pluginapi.PluginClient{
						"plugin1": createMockPluginClient("plugin1"),
						"plugin2": createMockPluginClient("plugin2"),
					},
				}
			},
			expected: []pluginapi.PluginClient{
				createMockPluginClient("stage1"),
				createMockPluginClient("stage2"),
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
			pr := tt.setup()
			plugins, err := pr.GetPluginClientsByAppConfig(tt.cfg)
			assert.Equal(t, tt.expected, plugins)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
func TestPluginRegistry_getPluginClientsByPipeline(t *testing.T) {
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
						"stage1": createMockPluginClient("stage1"),
						"stage2": createMockPluginClient("stage2"),
					},
				}
			},
			expected: []pluginapi.PluginClient{
				createMockPluginClient("stage1"),
				createMockPluginClient("stage2"),
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
			pr := tt.setup()
			plugins, err := pr.getPluginClientsByPipeline(tt.pipeline)
			assert.Equal(t, tt.expected, plugins)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
func TestPluginRegistry_getPluginsByPluginNames(t *testing.T) {
	tests := []struct {
		name        string
		pluginNames []string
		setup       func() *pluginRegistry
		expected    []pluginapi.PluginClient
		wantErr     bool
	}{
		{
			name:        "get plugins by valid plugin names",
			pluginNames: []string{"plugin1", "plugin2"},
			setup: func() *pluginRegistry {
				return &pluginRegistry{
					nameBasedPlugins: map[string]pluginapi.PluginClient{
						"plugin1": createMockPluginClient("plugin1"),
						"plugin2": createMockPluginClient("plugin2"),
					},
				}
			},
			expected: []pluginapi.PluginClient{
				createMockPluginClient("plugin1"),
				createMockPluginClient("plugin2"),
			},
			wantErr: false,
		},
		{
			name:        "no plugins found for empty plugin names",
			pluginNames: []string{},
			setup: func() *pluginRegistry {
				return &pluginRegistry{
					nameBasedPlugins: map[string]pluginapi.PluginClient{
						"plugin1": createMockPluginClient("plugin1"),
						"plugin2": createMockPluginClient("plugin2"),
					},
				}
			},
			wantErr: true,
		},
		{
			name:        "no plugins found for non-existent plugin names",
			pluginNames: []string{"plugin1", "plugin2"},
			setup: func() *pluginRegistry {
				return &pluginRegistry{
					nameBasedPlugins: map[string]pluginapi.PluginClient{
						"plugin3": createMockPluginClient("plugin3"),
						"plugin4": createMockPluginClient("plugin4"),
					},
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pr := tt.setup()
			plugins, err := pr.getPluginsByPluginNames(tt.pluginNames)
			assert.Equal(t, tt.expected, plugins)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
func TestPluginRegistry_GetPluginClientByStageName(t *testing.T) {
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
						"stage1": createMockPluginClient("stage1"),
					},
				}
			},
			expected: createMockPluginClient("stage1"),
			wantErr:  false,
		},
		{
			name:  "no plugin found for non-existent stage name",
			stage: "stage2",
			setup: func() *pluginRegistry {
				return &pluginRegistry{
					stageBasedPlugins: map[string]pluginapi.PluginClient{
						"stage1": createMockPluginClient("stage1"),
					},
				}
			},
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pr := tt.setup()
			plugin, err := pr.GetPluginClientByStageName(tt.stage)
			assert.Equal(t, tt.expected, plugin)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
