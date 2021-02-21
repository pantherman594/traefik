package main

import (
	"github.com/traefik/traefik/v2/pkg/config/static"
	"github.com/traefik/traefik/v2/pkg/plugins"
)

const outputDir = "./plugins-storage/"

func createPluginBuilder(staticConfiguration *static.Configuration) (*plugins.Builder, error) {
	client, plgs, devPlugins, err := initPlugins(staticConfiguration)
	if err != nil {
		return nil, err
	}

	return plugins.NewBuilder(client, plgs, devPlugins)
}

func initPlugins(staticCfg *static.Configuration) (*plugins.Client, map[string]plugins.Descriptor, map[string]plugins.DevPlugin, error) {
	if !isPilotEnabled(staticCfg) || !hasPlugins(staticCfg) {
		return nil, map[string]plugins.Descriptor{}, map[string]plugins.DevPlugin{}, nil
	}

	opts := plugins.ClientOptions{
		Output: outputDir,
		Token:  staticCfg.Pilot.Token,
	}

	client, err := plugins.NewClient(opts)
	if err != nil {
		return nil, nil, nil, err
	}

	err = plugins.Setup(client, staticCfg.Experimental.Plugins, staticCfg.Experimental.DevPlugins)
	if err != nil {
		return nil, nil, nil, err
	}

	return client, staticCfg.Experimental.Plugins, staticCfg.Experimental.DevPlugins, nil
}

func isPilotEnabled(staticCfg *static.Configuration) bool {
	return staticCfg.Pilot != nil && staticCfg.Pilot.Token != ""
}

func hasPlugins(staticCfg *static.Configuration) bool {
	return staticCfg.Experimental != nil &&
		(len(staticCfg.Experimental.Plugins) + len(staticCfg.Experimental.DevPlugins) > 0)
}
