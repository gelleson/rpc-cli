package config

import (
	"jsonrpc/pkg/types"
)

// Manager manages configuration sources and building effective configurations
type Manager struct {
	merger *Merger
}

// NewManager creates a new configuration manager
func NewManager() *Manager {
	return &Manager{
		merger: NewMerger(),
	}
}

// BuildForRequest builds an effective configuration for a specific request
// using the provided HCL file and CLI overrides
func (m *Manager) BuildForRequest(
	hclFile *types.HCLFile,
	request *types.Request,
	cliOverrides *types.CLIOverrides,
) *types.EffectiveConfig {
	m.merger.ClearSources()

	// Add sources in priority order (they will be auto-sorted by the merger)

	// 1. Default config (priority: 10)
	if defaultConfig, exists := hclFile.Configs[DefaultConfigName]; exists {
		m.merger.AddSource(NewDefaultConfigSource(defaultConfig))
	}

	// 2. Named config if specified (priority: 20)
	configName := GetConfigName(request, cliOverrides)
	if configName != "" && configName != DefaultConfigName {
		if namedConfig, exists := hclFile.Configs[configName]; exists {
			m.merger.AddSource(NewNamedConfigSource(configName, namedConfig))
		}
	}

	// 3. Request overrides (priority: 30)
	m.merger.AddSource(NewRequestConfigSource(request))

	// 4. CLI overrides (priority: 40)
	if cliOverrides != nil {
		m.merger.AddSource(NewCLIConfigSource(cliOverrides))
	}

	return m.merger.BuildEffective()
}

// BuildForCLI builds an effective configuration using only CLI overrides
func (m *Manager) BuildForCLI(cliOverrides *types.CLIOverrides) *types.EffectiveConfig {
	m.merger.ClearSources()

	if cliOverrides != nil {
		m.merger.AddSource(NewCLIConfigSource(cliOverrides))
	}

	return m.merger.BuildEffective()
}

// GetConfigNameForRequest returns the effective configuration name for a request
func (m *Manager) GetConfigNameForRequest(
	hclFile *types.HCLFile,
	request *types.Request,
	cliOverrides *types.CLIOverrides,
) string {
	m.merger.ClearSources()

	// Build sources just for config name determination
	configName := GetConfigName(request, cliOverrides)

	// Validate that the config exists in HCL file
	if configName != "" && configName != DefaultConfigName {
		if _, exists := hclFile.Configs[configName]; !exists {
			// Fallback to default if named config doesn't exist
			if _, hasDefault := hclFile.Configs["default"]; hasDefault {
				configName = "default"
			}
		}
	}

	return configName
}

// GetConfigName is a utility function to determine the config name for a request
// This maintains backward compatibility with the existing helper function
func GetConfigName(request *types.Request, cliOverrides *types.CLIOverrides) string {
	// CLI config override takes highest priority
	if cliOverrides != nil && cliOverrides.Config != "" {
		return cliOverrides.Config
	}

	// Request-level config override
	if request.Config != "" {
		return request.Config
	}

	// Default fallback
	return "default"
}
