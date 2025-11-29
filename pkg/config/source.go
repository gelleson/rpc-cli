package config

import (
	"jsonrpc/pkg/types"
)

const (
	// DefaultConfigName is the name for default configuration source
	DefaultConfigName = "default"
)

// Source represents a configuration source with a specific priority
type Source interface {
	// Name returns a human-readable name for this source
	Name() string

	// GetConfig returns the configuration from this source, or nil if no config
	GetConfig() *types.Config

	// Priority returns the priority of this source (higher number = higher priority)
	// Default config: 10, Named config: 20, Request overrides: 30, CLI overrides: 40
	Priority() int
}

// DefaultConfigSource provides the default configuration
type DefaultConfigSource struct {
	config *types.Config
}

// NewDefaultConfigSource creates a new default configuration source
func NewDefaultConfigSource(config *types.Config) *DefaultConfigSource {
	return &DefaultConfigSource{config: config}
}

func (s *DefaultConfigSource) Name() string {
	return DefaultConfigName
}

func (s *DefaultConfigSource) GetConfig() *types.Config {
	return s.config
}

func (s *DefaultConfigSource) Priority() int {
	return 10
}

// NamedConfigSource provides a named configuration profile
type NamedConfigSource struct {
	name   string
	config *types.Config
}

// NewNamedConfigSource creates a new named configuration source
func NewNamedConfigSource(name string, config *types.Config) *NamedConfigSource {
	return &NamedConfigSource{
		name:   name,
		config: config,
	}
}

func (s *NamedConfigSource) Name() string {
	return s.name
}

func (s *NamedConfigSource) GetConfig() *types.Config {
	return s.config
}

func (s *NamedConfigSource) Priority() int {
	return 20
}

// RequestConfigSource provides request-level configuration overrides
type RequestConfigSource struct {
	request *types.Request
}

// NewRequestConfigSource creates a new request configuration source
func NewRequestConfigSource(request *types.Request) *RequestConfigSource {
	return &RequestConfigSource{request: request}
}

func (s *RequestConfigSource) Name() string {
	return s.request.Name
}

func (s *RequestConfigSource) GetConfig() *types.Config {
	return &types.Config{
		URL:     s.request.URL,
		Headers: s.request.Headers,
		Timeout: s.request.Timeout,
	}
}

func (s *RequestConfigSource) Priority() int {
	return 30
}

// CLIConfigSource provides CLI flag overrides
type CLIConfigSource struct {
	overrides *types.CLIOverrides
}

// NewCLIConfigSource creates a new CLI configuration source
func NewCLIConfigSource(overrides *types.CLIOverrides) *CLIConfigSource {
	return &CLIConfigSource{overrides: overrides}
}

func (s *CLIConfigSource) Name() string {
	return "cli"
}

func (s *CLIConfigSource) GetConfig() *types.Config {
	return &types.Config{
		URL:     s.overrides.URL,
		Headers: s.overrides.Headers,
		Timeout: s.overrides.Timeout,
	}
}

func (s *CLIConfigSource) Priority() int {
	return 40
}
