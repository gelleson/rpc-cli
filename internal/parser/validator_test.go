package parser

import (
	"strings"
	"testing"

	"jsonrpc/pkg/types"
)

func TestValidator_Validate(t *testing.T) {
	tests := []struct {
		name    string
		hclFile *types.HCLFile
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid HCL file",
			hclFile: &types.HCLFile{
				Configs: map[string]*types.Config{
					"default": {URL: "https://api.example.com"},
				},
				Requests: []*types.Request{
					{Name: "test", Method: "eth_blockNumber"},
				},
			},
			wantErr: false,
		},
		{
			name: "missing method",
			hclFile: &types.HCLFile{
				Requests: []*types.Request{
					{Name: "test", Method: ""},
				},
			},
			wantErr: true,
			errMsg:  "missing required 'method' field",
		},
		{
			name: "non-existent config reference",
			hclFile: &types.HCLFile{
				Configs: map[string]*types.Config{
					"default": {URL: "https://api.example.com"},
				},
				Requests: []*types.Request{
					{Name: "test", Method: "eth_blockNumber", Config: "production"},
				},
			},
			wantErr: true,
			errMsg:  "references non-existent config",
		},
		{
			name: "valid config reference",
			hclFile: &types.HCLFile{
				Configs: map[string]*types.Config{
					"production": {URL: "https://prod.example.com"},
				},
				Requests: []*types.Request{
					{Name: "test", Method: "eth_blockNumber", Config: "production"},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := NewValidator()
			err := v.Validate(tt.hclFile)

			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("Validate() error = %v, should contain %s", err, tt.errMsg)
			}
		})
	}
}

func TestValidator_validateRequest(t *testing.T) {
	v := NewValidator()

	tests := []struct {
		name    string
		req     *types.Request
		configs map[string]*types.Config
		wantErr bool
	}{
		{
			name:    "valid request",
			req:     &types.Request{Name: "test", Method: "eth_blockNumber"},
			configs: map[string]*types.Config{},
			wantErr: false,
		},
		{
			name:    "missing method",
			req:     &types.Request{Name: "test", Method: ""},
			configs: map[string]*types.Config{},
			wantErr: true,
		},
		{
			name: "valid config reference",
			req:  &types.Request{Name: "test", Method: "test", Config: "prod"},
			configs: map[string]*types.Config{
				"prod": {URL: "https://prod.example.com"},
			},
			wantErr: false,
		},
		{
			name:    "invalid config reference",
			req:     &types.Request{Name: "test", Method: "test", Config: "missing"},
			configs: map[string]*types.Config{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.validateRequest(tt.req, tt.configs)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
