package executor

import (
	"testing"

	"jsonrpc/pkg/types"
)

func TestGetConfigName(t *testing.T) {
	tests := []struct {
		name      string
		req       *types.Request
		overrides *types.CLIOverrides
		want      string
	}{
		{
			name:      "no overrides, no request config",
			req:       &types.Request{},
			overrides: nil,
			want:      "default",
		},
		{
			name:      "request config set",
			req:       &types.Request{Config: "staging"},
			overrides: nil,
			want:      "staging",
		},
		{
			name:      "CLI override takes precedence",
			req:       &types.Request{Config: "staging"},
			overrides: &types.CLIOverrides{Config: "production"},
			want:      "production",
		},
		{
			name:      "CLI override set, no request config",
			req:       &types.Request{},
			overrides: &types.CLIOverrides{Config: "production"},
			want:      "production",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetConfigName(tt.req, tt.overrides)
			if got != tt.want {
				t.Errorf("GetConfigName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCountParams(t *testing.T) {
	tests := []struct {
		name   string
		params interface{}
		want   int
	}{
		{
			name:   "nil params",
			params: nil,
			want:   0,
		},
		{
			name:   "empty slice",
			params: []interface{}{},
			want:   0,
		},
		{
			name:   "slice with items",
			params: []interface{}{"a", "b", "c"},
			want:   3,
		},
		{
			name:   "empty map",
			params: map[string]interface{}{},
			want:   0,
		},
		{
			name: "map with items",
			params: map[string]interface{}{
				"key1": "value1",
				"key2": "value2",
			},
			want: 2,
		},
		{
			name:   "single value",
			params: "single",
			want:   1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CountParams(tt.params)
			if got != tt.want {
				t.Errorf("CountParams() = %v, want %v", got, tt.want)
			}
		})
	}
}
