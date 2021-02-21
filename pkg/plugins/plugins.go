package plugins

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/traefik/traefik/v2/pkg/log"
)

// Setup setup plugins environment.
func Setup(client *Client, plugins map[string]Descriptor, devPlugins map[string]DevPlugin) error {
	err := checkPluginsConfiguration(plugins, devPlugins)
	if err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	err = client.CleanArchives(plugins)
	if err != nil {
		return fmt.Errorf("failed to clean archives: %w", err)
	}

	ctx := context.Background()

	for pAlias, desc := range plugins {
		log.FromContext(ctx).Debugf("loading of plugin: %s: %s@%s", pAlias, desc.ModuleName, desc.Version)

		hash, err := client.Download(ctx, desc.ModuleName, desc.Version)
		if err != nil {
			_ = client.ResetAll()
			return fmt.Errorf("failed to download plugin %s: %w", desc.ModuleName, err)
		}

		err = client.Check(ctx, desc.ModuleName, desc.Version, hash)
		if err != nil {
			_ = client.ResetAll()
			return fmt.Errorf("failed to check archive integrity of the plugin %s: %w", desc.ModuleName, err)
		}
	}

	err = client.WriteState(plugins)
	if err != nil {
		_ = client.ResetAll()
		return fmt.Errorf("failed to write plugins state: %w", err)
	}

	for _, desc := range plugins {
		err = client.Unzip(desc.ModuleName, desc.Version)
		if err != nil {
			_ = client.ResetAll()
			return fmt.Errorf("failed to unzip archive: %w", err)
		}
	}

	return nil
}

func checkPluginsConfiguration(plugins map[string]Descriptor, devPlugins map[string]DevPlugin) error {
	if plugins == nil && devPlugins == nil {
		return nil
	}

	uniq := make(map[string]struct{})
	uniqAlias := make(map[string]struct{})

	var errs []string

	for pAlias, descriptor := range plugins {
		if descriptor.ModuleName == "" {
			errs = append(errs, fmt.Sprintf("%s: plugin name is missing", pAlias))
		}

		if descriptor.Version == "" {
			errs = append(errs, fmt.Sprintf("%s: plugin version is missing", pAlias))
		}

		if strings.HasPrefix(descriptor.ModuleName, "/") || strings.HasSuffix(descriptor.ModuleName, "/") {
			errs = append(errs, fmt.Sprintf("%s: plugin name should not start or end with a /", pAlias))
			continue
		}

		if _, ok := uniq[descriptor.ModuleName]; ok {
			errs = append(errs, fmt.Sprintf("only one version of a plugin is allowed, there is a duplicate of %s", descriptor.ModuleName))
			continue
		}

		if _, ok := uniqAlias[pAlias]; ok {
			errs = append(errs, fmt.Sprintf("only one alias of a plugin is allowed, there is a duplicate of %s", pAlias))
			continue
		}

		uniq[descriptor.ModuleName] = struct{}{}
		uniqAlias[pAlias] = struct{}{}
	}

	for pAlias, devDescriptor := range devPlugins {
		if devDescriptor.ModuleName == "" {
			errs = append(errs, fmt.Sprintf("%s: plugin name is missing", pAlias))
		}

		if devDescriptor.GoPath == "" {
			errs = append(errs, fmt.Sprintf("%s: plugin Go Path is missing (prefer a dedicated Go Path)", pAlias))
		}

		if strings.HasPrefix(devDescriptor.ModuleName, "/") || strings.HasSuffix(devDescriptor.ModuleName, "/") {
			errs = append(errs, fmt.Sprintf("%s: plugin name should not start or end with a /", pAlias))
			continue
		}

		if _, ok := uniq[devDescriptor.ModuleName]; ok {
			errs = append(errs, fmt.Sprintf("only one version of a plugin is allowed, there is a duplicate of %s", devDescriptor.ModuleName))
			continue
		}

		if _, ok := uniqAlias[pAlias]; ok {
			errs = append(errs, fmt.Sprintf("only one alias of a plugin is allowed, there is a duplicate of %s", pAlias))
			continue
		}

		m, err := ReadManifest(devDescriptor.GoPath, devDescriptor.ModuleName)
		if err != nil {
			errs = append(errs, err.Error())
			continue
		}

		if m.Type != "middleware" {
			errs = append(errs, fmt.Sprintf("%s: plugin type is unsupported", pAlias))
			continue
		}

		if m.Import == "" {
			errs = append(errs, fmt.Sprintf("%s: plugin is missing import", pAlias))
			continue
		}

		if !strings.HasPrefix(m.Import, devDescriptor.ModuleName) {
			errs = append(errs, fmt.Sprintf("%s: the import %q must be related to the module name %q", pAlias, m.Import, devDescriptor.ModuleName))
			continue
		}

		if m.DisplayName == "" {
			errs = append(errs, fmt.Sprintf("%s: plugin is missing display name", pAlias))
			continue
		}

		if m.Summary == "" {
			errs = append(errs, fmt.Sprintf("%s: plugin is missing summary", pAlias))
			continue
		}

		if m.TestData == nil {
			errs = append(errs, fmt.Sprintf("%s: plugin is missing test data", pAlias))
			continue
		}

		uniq[devDescriptor.ModuleName] = struct{}{}
		uniqAlias[pAlias] = struct{}{}
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, ": "))
	}

	return nil
}
