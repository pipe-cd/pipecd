package registry

import (
	"context"
	"fmt"

	config "github.com/pipe-cd/pipecd/pkg/configv1"
	pluginapi "github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1"
	"github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1/deployment"
)

type Plugin struct {
	Name string
	Cli  pluginapi.PluginClient
}

// PluginRegistry is the interface that provides methods to get plugin clients.
type PluginRegistry interface {
	GetPluginClientByStageName(name string) (pluginapi.PluginClient, error)
	GetPluginsByAppConfig(cfg *config.GenericApplicationSpec) ([]pluginapi.PluginClient, error)
}

type pluginRegistry struct {
	nameBasedPlugins  map[string]pluginapi.PluginClient
	stageBasedPlugins map[string]pluginapi.PluginClient

	// TODO: add more fields if needed (e.g. deploymentBasedPlugins, livestateBasedPlugins)
}

// NewPluginRegistry creates a new PluginRegistry based on the given plugins.
func NewPluginRegistry(ctx context.Context, plugins []Plugin) (PluginRegistry, error) {
	nameBasedPlugins := make(map[string]pluginapi.PluginClient)
	stageBasedPlugins := make(map[string]pluginapi.PluginClient)

	for _, plg := range plugins {
		// add the plugin to the name-based plugins
		nameBasedPlugins[plg.Name] = plg.Cli

		// add the plugin to the stage-based plugins
		res, err := plg.Cli.FetchDefinedStages(ctx, &deployment.FetchDefinedStagesRequest{})
		if err != nil {
			return nil, err
		}

		for _, stage := range res.Stages {
			stageBasedPlugins[stage] = plg.Cli
		}
	}

	return &pluginRegistry{
		nameBasedPlugins:  nameBasedPlugins,
		stageBasedPlugins: stageBasedPlugins,
	}, nil
}

// GetPluginClientByStageName returns the plugin client based on the given stage name.
func (pr *pluginRegistry) GetPluginClientByStageName(name string) (pluginapi.PluginClient, error) {
	plugin, ok := pr.stageBasedPlugins[name]
	if !ok {
		return nil, fmt.Errorf("no plugin found for the specified stage")
	}

	return plugin, nil
}

// GetPluginsByAppConfig returns the plugin clients based on the given configuration.
// The priority of determining plugins is as follows:
//  1. If the pipeline is specified, it will determine the plugins based on the pipeline stages.
//  2. If the plugins are specified, it will determine the plugins based on the plugin names.
//  3. If neither the pipeline nor the plugins are specified, it will return an error.
func (pr *pluginRegistry) GetPluginsByAppConfig(cfg *config.GenericApplicationSpec) ([]pluginapi.PluginClient, error) {
	if cfg.Pipeline != nil && len(cfg.Pipeline.Stages) > 0 {
		return pr.getPluginsByPipeline(cfg.Pipeline)
	}

	if cfg.Plugins != nil {
		return pr.getPluginsByPluginNames(cfg.Plugins)
	}

	return nil, fmt.Errorf("no plugin specified")
}

func (pr *pluginRegistry) getPluginsByPipeline(pipeline *config.DeploymentPipeline) ([]pluginapi.PluginClient, error) {
	plugins := make([]pluginapi.PluginClient, 0, len(pipeline.Stages))

	if len(pipeline.Stages) == 0 {
		return nil, fmt.Errorf("no plugin found for the specified stages")
	}

	for _, stage := range pipeline.Stages {
		plugin, ok := pr.stageBasedPlugins[stage.Name.String()]
		if ok {
			plugins = append(plugins, plugin)
		}
	}

	if len(plugins) == 0 {
		return nil, fmt.Errorf("no plugin found for the specified stages")
	}

	return plugins, nil
}

func (pr *pluginRegistry) getPluginsByPluginNames(names []string) ([]pluginapi.PluginClient, error) {
	if len(names) == 0 {
		return nil, fmt.Errorf("no plugin specified")
	}

	plugins := make([]pluginapi.PluginClient, 0, len(names))
	for _, name := range names {
		plugin, ok := pr.nameBasedPlugins[name]
		if ok {
			plugins = append(plugins, plugin)
		}
	}

	if len(plugins) == 0 {
		return nil, fmt.Errorf("no plugin found for the specified stages")
	}

	return plugins, nil
}
