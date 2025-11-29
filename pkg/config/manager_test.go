package config

import (
	"testing"

	"jsonrpc/pkg/types"
)

func TestGetConfigName(t *testing.T) {
	tests := []struct {
		name      string
		request   *types.Request
		overrides *types.CLIOverrides
		want      string
	}{
		{
			name: "CLI config override takes highest priority",
			request: &types.Request{
				Config: "request-config",
			},
			overrides: &types.CLIOverrides{
				Config: "cli-config",
			},
			want: "cli-config",
		},
		{
			name: "Request config override when no CLI override",
			request: &types.Request{
				Config: "request-config",
			},
			overrides: &types.CLIOverrides{},
			want:      "request-config",
		},
		{
			name: "Default when no overrides",
			request: &types.Request{
				Config: "",
			},
			overrides: &types.CLIOverrides{},
			want:      "default",
		},
		{
			name: "CLI override when request has no config",
			request: &types.Request{
				Config: "",
			},
			overrides: &types.CLIOverrides{
				Config: "cli-only",
			},
			want: "cli-only",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetConfigName(tt.request, tt.overrides)
			if got != tt.want {
				t.Errorf("GetConfigName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_BuildForRequest(t *testing.T) {
	hclFile := &types.HCLFile{
		Configs: map[string]*types.Config{
			"default": {
				URL:     "https://default.example.com",
				Timeout: 30,
				Headers: map[string]string{
					"Default-Header": "default-value",
				},
			},
			"production": {
				URL:     "https://prod.example.com",
				Timeout: 60,
				Headers: map[string]string{
					"Prod-Header": "prod-value",
				},
			},
		},
	}

	tests := []struct {
		name        string
		request     *types.Request
		overrides   *types.CLIOverrides
		wantURL     string
		wantTimeout int
	}{
		{
			name: "Default config only",
			request: &types.Request{
				Name:   "test-request",
				Method: "test.method",
			},
			overrides:   nil,
			wantURL:     "https://default.example.com",
			wantTimeout: 30,
		},
		{
			name: "Named config overrides default",
			request: &types.Request{
				Name:   "test-request",
				Method: "test.method",
				Config: "production",
			},
			overrides:   nil,
			wantURL:     "https://prod.example.com",
			wantTimeout: 60,
		},
		{
			name: "Request overrides named config",
			request: &types.Request{
				Name:    "test-request",
				Method:  "test.method",
				Config:  "production",
				URL:     "https://request.example.com",
				Timeout: 45,
			},
			overrides:   nil,
			wantURL:     "https://request.example.com",
			wantTimeout: 45,
		},
		{
			name: "CLI overrides take highest priority",
			request: &types.Request{
				Name:    "test-request",
				Method:  "test.method",
				Config:  "production",
				URL:     "https://request.example.com",
				Timeout: 45,
			},
			overrides: &types.CLIOverrides{
				URL:     "https://cli.example.com",
				Timeout: 90,
			},
			wantURL:     "https://cli.example.com",
			wantTimeout: 90,
		},
	}

	manager := NewManager()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := manager.BuildForRequest(hclFile, tt.request, tt.overrides)

			if config.URL != tt.wantURL {
				t.Errorf("BuildForRequest().URL = %v, want %v", config.URL, tt.wantURL)
			}

			if config.Timeout != tt.wantTimeout {
				t.Errorf("BuildForRequest().Timeout = %v, want %v", config.Timeout, tt.wantTimeout)
			}
		})
	}
}

func TestManager_GetConfigNameForRequest(t *testing.T) {
	hclFile := &types.HCLFile{
		Configs: map[string]*types.Config{
			"default":    {URL: "https://default.example.com"},
			"production": {URL: "https://prod.example.com"},
		},
	}

	tests := []struct {
		name      string
		request   *types.Request
		overrides *types.CLIOverrides
		want      string
	}{
		{
			name: "CLI config override",
			request: &types.Request{
				Config: "production",
			},
			overrides: &types.CLIOverrides{
				Config: "cli-config",
			},
			want: "default", // Falls back to default since cli-config doesn't exist
		},
		{
			name: "Existing named config",
			request: &types.Request{
				Config: "production",
			},
			overrides: nil,
			want:      "production",
		},
		{
			name: "Non-existent config falls back to default",
			request: &types.Request{
				Config: "non-existent",
			},
			overrides: nil,
			want:      "default",
		},
	}

	manager := NewManager()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := manager.GetConfigNameForRequest(hclFile, tt.request, tt.overrides)
			if got != tt.want {
				t.Errorf("GetConfigNameForRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}
