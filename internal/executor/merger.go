package executor

import "jsonrpc/pkg/types"

// ConfigMerger handles merging configurations from different sources
type ConfigMerger struct{}

// NewConfigMerger creates a new ConfigMerger
func NewConfigMerger() *ConfigMerger {
	return &ConfigMerger{}
}

// MergeFromConfig merges settings from a Config into EffectiveConfig
func (m *ConfigMerger) MergeFromConfig(effective *types.EffectiveConfig, config *types.Config) {
	if config.URL != "" {
		effective.URL = config.URL
	}

	for k, v := range config.Headers {
		effective.Headers[k] = v
	}

	if config.Timeout > 0 {
		effective.Timeout = config.Timeout
	}
}

// MergeFromRequest merges request-level overrides into EffectiveConfig
func (m *ConfigMerger) MergeFromRequest(effective *types.EffectiveConfig, req *types.Request) {
	if req.URL != "" {
		effective.URL = req.URL
	}

	for k, v := range req.Headers {
		effective.Headers[k] = v
	}

	if req.Timeout > 0 {
		effective.Timeout = req.Timeout
	}
}

// MergeFromCLI merges CLI overrides into EffectiveConfig
func (m *ConfigMerger) MergeFromCLI(effective *types.EffectiveConfig, overrides *types.CLIOverrides) {
	if overrides.URL != "" {
		effective.URL = overrides.URL
	}

	for k, v := range overrides.Headers {
		effective.Headers[k] = v
	}

	if overrides.Timeout > 0 {
		effective.Timeout = overrides.Timeout
	}
}
