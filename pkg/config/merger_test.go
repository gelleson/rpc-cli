package config

import (
	"testing"

	"jsonrpc/pkg/types"
)

func TestMerger_AddSource(t *testing.T) {
	merger := NewMerger()

	// Add sources in random priority order
	merger.AddSource(NewCLIConfigSource(&types.CLIOverrides{Config: "cli"}))
	merger.AddSource(NewDefaultConfigSource(&types.Config{URL: "default"}))
	merger.AddSource(NewRequestConfigSource(&types.Request{Name: "test-request", URL: "request"}))
	merger.AddSource(NewNamedConfigSource("named", &types.Config{URL: "named"}))

	sources := merger.GetSources()

	// Should be sorted by priority (ascending)
	if len(sources) != 4 {
		t.Fatalf("Expected 4 sources, got %d", len(sources))
	}

	expectedOrder := []string{"default", "named", "test-request", "cli"}
	for i, expectedName := range expectedOrder {
		if sources[i].Name() != expectedName {
			t.Errorf("Source %d should be %s, got %s", i, expectedName, sources[i].Name())
		}
	}
}

func TestMerger_BuildEffective(t *testing.T) {
	merger := NewMerger()

	// Add sources in priority order
	merger.AddSource(NewDefaultConfigSource(&types.Config{
		URL:     "https://default.example.com",
		Timeout: 30,
		Headers: map[string]string{
			"Default-Header": "default-value",
		},
	}))

	merger.AddSource(NewNamedConfigSource("production", &types.Config{
		URL:     "https://prod.example.com",
		Timeout: 60,
		Headers: map[string]string{
			"Prod-Header": "prod-value",
		},
	}))

	merger.AddSource(NewRequestConfigSource(&types.Request{
		URL:     "https://request.example.com",
		Timeout: 45,
		Headers: map[string]string{
			"Request-Header": "request-value",
		},
	}))

	merger.AddSource(NewCLIConfigSource(&types.CLIOverrides{
		URL:     "https://cli.example.com",
		Timeout: 90,
		Headers: map[string]string{
			"CLI-Header": "cli-value",
		},
	}))

	config := merger.BuildEffective()

	// CLI overrides should take highest priority
	if config.URL != "https://cli.example.com" {
		t.Errorf("Expected URL from CLI, got %s", config.URL)
	}

	if config.Timeout != 90 {
		t.Errorf("Expected timeout from CLI, got %d", config.Timeout)
	}

	// Headers should be merged from all sources
	expectedHeaders := map[string]string{
		"Default-Header": "default-value",
		"Prod-Header":    "prod-value",
		"Request-Header": "request-value",
		"CLI-Header":     "cli-value",
	}

	for key, expectedValue := range expectedHeaders {
		if config.Headers[key] != expectedValue {
			t.Errorf("Expected header %s=%s, got %s", key, expectedValue, config.Headers[key])
		}
	}
}

func TestMerger_GetConfigName(t *testing.T) {
	merger := NewMerger()

	// Test with no sources
	if name := merger.GetConfigName(); name != "" {
		t.Errorf("Expected empty config name with no sources, got %s", name)
	}

	// Add sources
	merger.AddSource(NewDefaultConfigSource(&types.Config{URL: "default"}))
	merger.AddSource(NewRequestConfigSource(&types.Request{
		Name:   "test-request",
		Config: "production",
	}))
	merger.AddSource(NewCLIConfigSource(&types.CLIOverrides{Config: "staging"}))

	// Request config should be returned (CLI is skipped for config name determination)
	if name := merger.GetConfigName(); name != "production" {
		t.Errorf("Expected config name 'production', got %s", name)
	}
}

func TestDefaultConfigSource(t *testing.T) {
	config := &types.Config{URL: "https://default.example.com"}
	source := NewDefaultConfigSource(config)

	if source.Name() != DefaultConfigName {
		t.Errorf("Expected name '%s', got %s", DefaultConfigName, source.Name())
	}

	if source.Priority() != 10 {
		t.Errorf("Expected priority 10, got %d", source.Priority())
	}

	if source.GetConfig() != config {
		t.Error("GetConfig() should return the original config")
	}
}

func TestNamedConfigSource(t *testing.T) {
	config := &types.Config{URL: "https://prod.example.com"}
	source := NewNamedConfigSource("production", config)

	if source.Name() != "production" {
		t.Errorf("Expected name 'production', got %s", source.Name())
	}

	if source.Priority() != 20 {
		t.Errorf("Expected priority 20, got %d", source.Priority())
	}

	if source.GetConfig() != config {
		t.Error("GetConfig() should return the original config")
	}
}

func TestRequestConfigSource(t *testing.T) {
	request := &types.Request{
		Name:    "test-request",
		URL:     "https://request.example.com",
		Timeout: 45,
		Headers: map[string]string{"Request-Header": "request-value"},
		Config:  "production",
	}
	source := NewRequestConfigSource(request)

	if source.Name() != "test-request" {
		t.Errorf("Expected name 'test-request', got %s", source.Name())
	}

	if source.Priority() != 30 {
		t.Errorf("Expected priority 30, got %d", source.Priority())
	}

	config := source.GetConfig()
	if config.URL != request.URL {
		t.Errorf("Expected URL %s, got %s", request.URL, config.URL)
	}

	if config.Timeout != request.Timeout {
		t.Errorf("Expected timeout %d, got %d", request.Timeout, config.Timeout)
	}

	if config.Headers["Request-Header"] != request.Headers["Request-Header"] {
		t.Errorf("Expected header Request-Header=%s, got %s",
			request.Headers["Request-Header"], config.Headers["Request-Header"])
	}
}

func TestCLIConfigSource(t *testing.T) {
	overrides := &types.CLIOverrides{
		URL:     "https://cli.example.com",
		Timeout: 90,
		Headers: map[string]string{"CLI-Header": "cli-value"},
		Config:  "staging",
	}
	source := NewCLIConfigSource(overrides)

	if source.Name() != "cli" {
		t.Errorf("Expected name 'cli', got %s", source.Name())
	}

	if source.Priority() != 40 {
		t.Errorf("Expected priority 40, got %d", source.Priority())
	}

	config := source.GetConfig()
	if config.URL != overrides.URL {
		t.Errorf("Expected URL %s, got %s", overrides.URL, config.URL)
	}

	if config.Timeout != overrides.Timeout {
		t.Errorf("Expected timeout %d, got %d", overrides.Timeout, config.Timeout)
	}

	if config.Headers["CLI-Header"] != overrides.Headers["CLI-Header"] {
		t.Errorf("Expected header CLI-Header=%s, got %s",
			overrides.Headers["CLI-Header"], config.Headers["CLI-Header"])
	}
}

func TestMerger_ClearSources(t *testing.T) {
	merger := NewMerger()
	merger.AddSource(NewDefaultConfigSource(&types.Config{}))
	merger.AddSource(NewCLIConfigSource(&types.CLIOverrides{}))

	if len(merger.GetSources()) != 2 {
		t.Errorf("Expected 2 sources before clearing, got %d", len(merger.GetSources()))
	}

	merger.ClearSources()

	if len(merger.GetSources()) != 0 {
		t.Errorf("Expected 0 sources after clearing, got %d", len(merger.GetSources()))
	}

	// Test that we can still build effective config after clearing
	config := merger.BuildEffective()
	if config == nil {
		t.Error("BuildEffective() should return a valid config even with no sources")
	}
}
