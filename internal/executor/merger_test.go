package executor

import (
	"reflect"
	"testing"

	"jsonrpc/pkg/types"
)

func TestConfigMerger_MergeFromConfig(t *testing.T) {
	merger := NewConfigMerger()

	tests := []struct {
		name      string
		effective *types.EffectiveConfig
		config    *types.Config
		want      *types.EffectiveConfig
	}{
		{
			name:      "merge URL",
			effective: types.NewEffectiveConfig(),
			config:    &types.Config{URL: "https://api.example.com"},
			want: &types.EffectiveConfig{
				URL:     "https://api.example.com",
				Headers: map[string]string{},
				Timeout: 30,
			},
		},
		{
			name:      "merge headers",
			effective: types.NewEffectiveConfig(),
			config: &types.Config{
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
			},
			want: &types.EffectiveConfig{
				URL: "",
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
				Timeout: 30,
			},
		},
		{
			name:      "merge timeout",
			effective: types.NewEffectiveConfig(),
			config:    &types.Config{Timeout: 60},
			want: &types.EffectiveConfig{
				URL:     "",
				Headers: map[string]string{},
				Timeout: 60,
			},
		},
		{
			name:      "merge all fields",
			effective: types.NewEffectiveConfig(),
			config: &types.Config{
				URL: "https://api.example.com",
				Headers: map[string]string{
					"Authorization": "Bearer token",
				},
				Timeout: 90,
			},
			want: &types.EffectiveConfig{
				URL: "https://api.example.com",
				Headers: map[string]string{
					"Authorization": "Bearer token",
				},
				Timeout: 90,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			merger.MergeFromConfig(tt.effective, tt.config)
			if !reflect.DeepEqual(tt.effective, tt.want) {
				t.Errorf("MergeFromConfig() got = %+v, want %+v", tt.effective, tt.want)
			}
		})
	}
}

func TestConfigMerger_MergeFromRequest(t *testing.T) {
	merger := NewConfigMerger()

	req := &types.Request{
		URL: "https://custom.example.com",
		Headers: map[string]string{
			"X-Custom": "value",
		},
		Timeout: 120,
	}

	effective := types.NewEffectiveConfig()
	merger.MergeFromRequest(effective, req)

	if effective.URL != req.URL {
		t.Errorf("URL not merged correctly")
	}

	if effective.Headers["X-Custom"] != "value" {
		t.Errorf("Headers not merged correctly")
	}

	if effective.Timeout != req.Timeout {
		t.Errorf("Timeout not merged correctly")
	}
}

func TestConfigMerger_MergeFromCLI(t *testing.T) {
	merger := NewConfigMerger()

	overrides := &types.CLIOverrides{
		URL: "https://cli.example.com",
		Headers: map[string]string{
			"X-CLI": "cli-value",
		},
		Timeout: 45,
	}

	effective := types.NewEffectiveConfig()
	merger.MergeFromCLI(effective, overrides)

	if effective.URL != overrides.URL {
		t.Errorf("URL not merged correctly")
	}

	if effective.Headers["X-CLI"] != "cli-value" {
		t.Errorf("Headers not merged correctly")
	}

	if effective.Timeout != overrides.Timeout {
		t.Errorf("Timeout not merged correctly")
	}
}

func TestConfigMerger_MergePriority(t *testing.T) {
	merger := NewConfigMerger()
	effective := types.NewEffectiveConfig()

	// Start with base config
	baseConfig := &types.Config{
		URL:     "https://default.example.com",
		Timeout: 30,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
	merger.MergeFromConfig(effective, baseConfig)

	// Apply request override
	req := &types.Request{
		URL: "https://request.example.com",
		Headers: map[string]string{
			"Authorization": "Bearer token",
		},
	}
	merger.MergeFromRequest(effective, req)

	// Apply CLI override (highest priority)
	overrides := &types.CLIOverrides{
		Timeout: 60,
		Headers: map[string]string{
			"X-CLI-Override": "true",
		},
	}
	merger.MergeFromCLI(effective, overrides)

	// Verify final state
	if effective.URL != "https://request.example.com" {
		t.Errorf("URL should be from request, got %s", effective.URL)
	}

	if effective.Timeout != 60 {
		t.Errorf("Timeout should be from CLI, got %d", effective.Timeout)
	}

	if effective.Headers["Content-Type"] != "application/json" {
		t.Error("Base config header should be preserved")
	}

	if effective.Headers["Authorization"] != "Bearer token" {
		t.Error("Request header should be present")
	}

	if effective.Headers["X-CLI-Override"] != "true" {
		t.Error("CLI header should be present")
	}
}
