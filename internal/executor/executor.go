package executor

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"jsonrpc/pkg/types"
)

// Executor handles JSON-RPC request execution
type Executor struct {
	client *http.Client
}

// New creates a new Executor instance
func New() *Executor {
	return &Executor{
		client: &http.Client{},
	}
}

// Execute executes a single JSON-RPC request
func (e *Executor) Execute(hclFile *types.HCLFile, req *types.Request, overrides *types.CLIOverrides, requestID int) (*types.ExecutionResult, error) {
	startTime := time.Now()

	// Build effective configuration
	config := e.buildEffectiveConfig(hclFile, req, overrides)

	// Validate URL
	if config.URL == "" {
		return &types.ExecutionResult{
			Request:  req,
			Duration: time.Since(startTime),
			Error:    fmt.Errorf("no URL configured for request '%s'", req.Name),
		}, nil
	}

	// Create and execute JSON-RPC request
	response, err := e.executeJSONRPC(config, req, requestID)
	if err != nil {
		return &types.ExecutionResult{
			Request:  req,
			Duration: time.Since(startTime),
			Error:    err,
		}, nil
	}

	return &types.ExecutionResult{
		Request:  req,
		Response: response,
		Duration: time.Since(startTime),
	}, nil
}

// ExecuteAll executes multiple requests
func (e *Executor) ExecuteAll(hclFile *types.HCLFile, requests []*types.Request, overrides *types.CLIOverrides) ([]*types.ExecutionResult, error) {
	results := make([]*types.ExecutionResult, 0, len(requests))

	for i, req := range requests {
		result, err := e.Execute(hclFile, req, overrides, i+1)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	return results, nil
}

// buildEffectiveConfig merges configurations following priority order
func (e *Executor) buildEffectiveConfig(hclFile *types.HCLFile, req *types.Request, overrides *types.CLIOverrides) *types.EffectiveConfig {
	config := types.NewEffectiveConfig()
	merger := NewConfigMerger()

	// 1. Start with default config
	if defaultConfig, exists := hclFile.Configs["default"]; exists {
		merger.MergeFromConfig(config, defaultConfig)
	}

	// 2. Apply named config if specified
	configName := req.Config
	if overrides != nil && overrides.Config != "" {
		configName = overrides.Config
	}
	if configName != "" && configName != "default" {
		if namedConfig, exists := hclFile.Configs[configName]; exists {
			merger.MergeFromConfig(config, namedConfig)
		}
	}

	// 3. Apply request-level overrides
	merger.MergeFromRequest(config, req)

	// 4. Apply CLI overrides (highest priority)
	if overrides != nil {
		merger.MergeFromCLI(config, overrides)
	}

	return config
}

// executeJSONRPC executes a JSON-RPC request and returns the response
func (e *Executor) executeJSONRPC(config *types.EffectiveConfig, req *types.Request, requestID int) (*types.JSONRPCResponse, error) {
	// Create JSON-RPC request
	rpcReq := types.NewJSONRPCRequest(req.Method, req.ProcessedParams, requestID)

	// Marshal to JSON
	reqBody, err := json.Marshal(rpcReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequest("POST", config.URL, bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Add headers
	httpReq.Header.Set("Content-Type", "application/json")
	for k, v := range config.Headers {
		httpReq.Header.Set(k, v)
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.Timeout)*time.Second)
	defer cancel()
	httpReq = httpReq.WithContext(ctx)

	// Execute request
	httpResp, err := e.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer func() {
		_ = httpResp.Body.Close()
	}()

	// Read response body
	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check HTTP status
	if httpResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP error: %d %s - %s", httpResp.StatusCode, httpResp.Status, string(respBody))
	}

	// Parse JSON-RPC response
	var rpcResp types.JSONRPCResponse
	if err := json.Unmarshal(respBody, &rpcResp); err != nil {
		return nil, fmt.Errorf("failed to parse JSON-RPC response: %w", err)
	}

	return &rpcResp, nil
}
