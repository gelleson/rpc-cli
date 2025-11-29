package types

import (
	"encoding/json"
	"time"

	"github.com/zclconf/go-cty/cty"
	"jsonrpc/pkg/constants"
)

// Config represents a configuration block for JSON-RPC requests
type Config struct {
	URL     string            `hcl:"url,optional" json:"url,omitempty"`
	Headers map[string]string `hcl:"headers,optional" json:"headers,omitempty"`
	Timeout int               `hcl:"timeout,optional" json:"timeout,omitempty"` // in seconds
}

// NewConfig creates a new Config with sensible defaults
func NewConfig() *Config {
	return &Config{
		Headers: make(map[string]string),
		Timeout: constants.DefaultTimeoutSeconds,
	}
}

// Request represents a JSON-RPC request definition
type Request struct {
	Name            string            `hcl:"name,label" json:"name"`
	Method          string            `hcl:"method" json:"method"`
	Params          cty.Value         `hcl:"params,optional" json:"-"`
	URL             string            `hcl:"url,optional" json:"url,omitempty"`
	Headers         map[string]string `hcl:"headers,optional" json:"headers,omitempty"`
	Timeout         int               `hcl:"timeout,optional" json:"timeout,omitempty"`
	Config          string            `hcl:"config,optional" json:"config,omitempty"`
	ProcessedParams any               `hcl:"-" json:"params,omitempty"`
}

// NewRequest creates a new Request with initialized maps
func NewRequest(name string) *Request {
	return &Request{
		Name:    name,
		Headers: make(map[string]string),
	}
}

// HCLFile represents the entire parsed HCL file structure
type HCLFile struct {
	Configs  map[string]*Config
	Requests []*Request
}

// NewHCLFile creates a new HCLFile with initialized maps
func NewHCLFile() *HCLFile {
	return &HCLFile{
		Configs:  make(map[string]*Config),
		Requests: make([]*Request, 0),
	}
}

// JSONRPCRequest represents a JSON-RPC 2.0 request
type JSONRPCRequest struct {
	JSONRPC string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  any    `json:"params"`
	ID      int    `json:"id"`
}

// NewJSONRPCRequest creates a new JSON-RPC request
func NewJSONRPCRequest(method string, params any, id int) *JSONRPCRequest {
	return &JSONRPCRequest{
		JSONRPC: constants.DefaultJSONRPCVersion,
		Method:  method,
		Params:  params,
		ID:      id,
	}
}

// JSONRPCResponse represents a JSON-RPC 2.0 response
type JSONRPCResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *RPCError       `json:"error,omitempty"`
	ID      int             `json:"id"`
}

// IsError returns true if the response contains an error
func (r *JSONRPCResponse) IsError() bool {
	return r.Error != nil
}

// RPCError represents a JSON-RPC error object
type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// Error implements the error interface
func (e *RPCError) Error() string {
	return e.Message
}

// ExecutionResult contains the result of executing a request
type ExecutionResult struct {
	Request  *Request
	Response *JSONRPCResponse
	Duration time.Duration
	Error    error
}

// IsSuccess returns true if the execution was successful
func (r *ExecutionResult) IsSuccess() bool {
	return r.Error == nil && (r.Response == nil || !r.Response.IsError())
}

// EffectiveConfig holds the final merged configuration for a request
type EffectiveConfig struct {
	URL     string
	Headers map[string]string
	Timeout int
}

// NewEffectiveConfig creates a new EffectiveConfig with defaults
func NewEffectiveConfig() *EffectiveConfig {
	return &EffectiveConfig{
		Headers: make(map[string]string),
		Timeout: constants.DefaultTimeoutSeconds,
	}
}

// CLIOverrides holds configuration overrides from CLI flags
type CLIOverrides struct {
	URL     string
	Headers map[string]string
	Timeout int
	Config  string
}

// NewCLIOverrides creates a new CLIOverrides with initialized maps
func NewCLIOverrides() *CLIOverrides {
	return &CLIOverrides{
		Headers: make(map[string]string),
	}
}

// OutputFormat represents the output format type
type OutputFormat string

const (
	OutputTable    OutputFormat = "table"
	OutputDetailed OutputFormat = "detailed"
	OutputJSON     OutputFormat = "json"
)
