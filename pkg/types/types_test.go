package types

import (
	"encoding/json"
	"testing"
)

func TestNewConfig(t *testing.T) {
	config := NewConfig()

	if config == nil {
		t.Fatal("NewConfig() returned nil")
		return
	}

	if config.Headers == nil {
		t.Error("Headers map should be initialized")
	}

	if config.Timeout != 30 {
		t.Errorf("Default timeout should be 30, got %d", config.Timeout)
	}
}

func TestNewRequest(t *testing.T) {
	name := "test_request"
	req := NewRequest(name)

	if req == nil {
		t.Fatal("NewRequest() returned nil")
		return
	}

	if req.Name != name {
		t.Errorf("Expected name %s, got %s", name, req.Name)
	}

	if req.Headers == nil {
		t.Error("Headers map should be initialized")
	}
}

func TestNewHCLFile(t *testing.T) {
	hclFile := NewHCLFile()

	if hclFile == nil {
		t.Fatal("NewHCLFile() returned nil")
		return
	}

	if hclFile.Configs == nil {
		t.Error("Configs map should be initialized")
	}

	if hclFile.Requests == nil {
		t.Error("Requests slice should be initialized")
	}
}

func TestNewJSONRPCRequest(t *testing.T) {
	method := "eth_blockNumber"
	params := []any{}
	id := 1

	req := NewJSONRPCRequest(method, params, id)

	if req == nil {
		t.Fatal("NewJSONRPCRequest() returned nil")
		return
	}

	if req.JSONRPC != "2.0" {
		t.Errorf("Expected JSONRPC version 2.0, got %s", req.JSONRPC)
	}

	if req.Method != method {
		t.Errorf("Expected method %s, got %s", method, req.Method)
	}

	if req.ID != id {
		t.Errorf("Expected ID %d, got %d", id, req.ID)
	}
}

func TestJSONRPCResponse_IsError(t *testing.T) {
	tests := []struct {
		name     string
		response *JSONRPCResponse
		want     bool
	}{
		{
			name: "response with error",
			response: &JSONRPCResponse{
				Error: &RPCError{Code: -32600, Message: "Invalid Request"},
			},
			want: true,
		},
		{
			name: "response without error",
			response: &JSONRPCResponse{
				Result: json.RawMessage(`"0x123"`),
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.response.IsError(); got != tt.want {
				t.Errorf("IsError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRPCError_Error(t *testing.T) {
	msg := "Invalid params"
	rpcErr := &RPCError{
		Code:    -32602,
		Message: msg,
	}

	if rpcErr.Error() != msg {
		t.Errorf("Error() = %s, want %s", rpcErr.Error(), msg)
	}
}

func TestExecutionResult_IsSuccess(t *testing.T) {
	tests := []struct {
		name   string
		result *ExecutionResult
		want   bool
	}{
		{
			name: "successful execution",
			result: &ExecutionResult{
				Response: &JSONRPCResponse{
					Result: json.RawMessage(`"0x123"`),
				},
			},
			want: true,
		},
		{
			name: "execution with error",
			result: &ExecutionResult{
				Error: &RPCError{Code: -32600, Message: "Invalid Request"},
			},
			want: false,
		},
		{
			name: "execution with RPC error",
			result: &ExecutionResult{
				Response: &JSONRPCResponse{
					Error: &RPCError{Code: -32600, Message: "Invalid Request"},
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.result.IsSuccess(); got != tt.want {
				t.Errorf("IsSuccess() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewEffectiveConfig(t *testing.T) {
	config := NewEffectiveConfig()

	if config == nil {
		t.Fatal("NewEffectiveConfig() returned nil")
		return
	}

	if config.Headers == nil {
		t.Error("Headers map should be initialized")
	}

	if config.Timeout != 30 {
		t.Errorf("Default timeout should be 30, got %d", config.Timeout)
	}
}

func TestNewCLIOverrides(t *testing.T) {
	overrides := NewCLIOverrides()

	if overrides == nil {
		t.Fatal("NewCLIOverrides() returned nil")
		return
	}

	if overrides.Headers == nil {
		t.Error("Headers map should be initialized")
	}
}
