package config

import (
	"sort"

	"jsonrpc/pkg/types"
)

// Merger handles merging configurations from multiple sources
type Merger struct {
	sources []Source
}

// NewMerger creates a new configuration merger
func NewMerger() *Merger {
	return &Merger{
		sources: make([]Source, 0),
	}
}

// AddSource adds a configuration source to the merger
// Sources are automatically sorted by priority (lowest first)
func (m *Merger) AddSource(source Source) {
	m.sources = append(m.sources, source)
	// Sort by priority (ascending - lower priority sources are applied first)
	sort.Slice(m.sources, func(i, j int) bool {
		return m.sources[i].Priority() < m.sources[j].Priority()
	})
}

// BuildEffective builds the effective configuration by merging all sources
// Sources are applied in priority order, with higher priority sources overriding lower priority ones
func (m *Merger) BuildEffective() *types.EffectiveConfig {
	config := types.NewEffectiveConfig()

	for _, source := range m.sources {
		sourceConfig := source.GetConfig()
		if sourceConfig == nil {
			continue
		}
		m.mergeInto(config, sourceConfig)
	}

	return config
}

// GetConfigName returns the effective configuration name for display purposes
// Returns the name of the highest priority non-default, non-CLI config source
func (m *Merger) GetConfigName() string {
	if len(m.sources) == 0 {
		return ""
	}

	// Find the highest priority source that provides a meaningful name
	// Skip CLI and look for named config first
	configName := ""
	for i := len(m.sources) - 1; i >= 0; i-- {
		source := m.sources[i]
		switch s := source.(type) {
		case *NamedConfigSource:
			if s.config != nil {
				configName = s.name
			}
		case *RequestConfigSource:
			// If request has a config reference, use that
			if s.request.Config != "" && s.request.Config != DefaultConfigName {
				configName = s.request.Config
			}
		}
	}

	// If no named config found, default to "default"
	if configName == "" {
		// Check if we have a default config source
		for _, source := range m.sources {
			if _, ok := source.(*DefaultConfigSource); ok {
				configName = DefaultConfigName
				break
			}
		}
	}

	return configName
}

// mergeInto merges a source configuration into an effective configuration
// This is the core merging logic that eliminates duplication
func (m *Merger) mergeInto(effective *types.EffectiveConfig, source *types.Config) {
	if source.URL != "" {
		effective.URL = source.URL
	}

	// Merge headers - source headers override existing ones
	for k, v := range source.Headers {
		effective.Headers[k] = v
	}

	if source.Timeout > 0 {
		effective.Timeout = source.Timeout
	}
}

// ClearSources removes all sources from the merger
func (m *Merger) ClearSources() {
	m.sources = m.sources[:0]
}

// GetSources returns a copy of all sources in priority order
func (m *Merger) GetSources() []Source {
	result := make([]Source, len(m.sources))
	copy(result, m.sources)
	return result
}
