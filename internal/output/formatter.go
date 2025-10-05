package output

import (
	"encoding/json"
	"fmt"
	"strings"

	"jsonrpc/internal/executor"
	"jsonrpc/pkg/types"
)

// Formatter handles output formatting
type Formatter struct {
	masker *SensitiveMasker
}

// New creates a new Formatter instance
func New() *Formatter {
	return &Formatter{
		masker: NewSensitiveMasker(),
	}
}

// FormatRequestList formats requests as a table
func (f *Formatter) FormatRequestList(hclFile *types.HCLFile, requests []*types.Request, overrides *types.CLIOverrides) {
	// Print header
	fmt.Printf("%-25s %-30s %-15s %-10s\n", "NAME", "METHOD", "CONFIG", "PARAMS")
	fmt.Println(strings.Repeat("-", 85))

	// Print each request
	for _, req := range requests {
		configName := executor.GetConfigName(req, overrides)
		paramCount := executor.CountParams(req.ProcessedParams)

		fmt.Printf("%-25s %-30s %-15s %-10d\n",
			truncate(req.Name, 25),
			truncate(req.Method, 30),
			truncate(configName, 15),
			paramCount,
		)
	}
}

// FormatRequestDetailed formats requests in detailed boxed format
func (f *Formatter) FormatRequestDetailed(hclFile *types.HCLFile, requests []*types.Request, overrides *types.CLIOverrides) {
	for i, req := range requests {
		if i > 0 {
			fmt.Println()
		}
		f.formatSingleRequestDetailed(hclFile, req, overrides)
	}
}

// formatSingleRequestDetailed formats a single request in detailed format
func (f *Formatter) formatSingleRequestDetailed(hclFile *types.HCLFile, req *types.Request, overrides *types.CLIOverrides) {
	merger := executor.NewConfigMerger()
	config := types.NewEffectiveConfig()

	// Build effective config for display
	if defaultConfig, exists := hclFile.Configs["default"]; exists {
		merger.MergeFromConfig(config, defaultConfig)
	}

	configName := executor.GetConfigName(req, overrides)
	if configName != "" && configName != "default" {
		if namedConfig, exists := hclFile.Configs[configName]; exists {
			merger.MergeFromConfig(config, namedConfig)
		}
	}

	merger.MergeFromRequest(config, req)
	if overrides != nil {
		merger.MergeFromCLI(config, overrides)
	}

	// Top border
	fmt.Println("┌" + strings.Repeat("─", 78) + "┐")

	// Title
	fmt.Printf("│ %-76s │\n", req.Name)
	fmt.Println("├" + strings.Repeat("─", 78) + "┤")

	// Method
	fmt.Printf("│ Method:  %-67s │\n", req.Method)

	// URL
	fmt.Printf("│ URL:     %-67s │\n", truncate(config.URL, 67))

	// Config
	fmt.Printf("│ Config:  %-67s │\n", configName)

	// Timeout
	fmt.Printf("│ Timeout: %-67s │\n", fmt.Sprintf("%ds", config.Timeout))

	// Headers (if any)
	if len(config.Headers) > 0 {
		fmt.Printf("│ Headers:%-68s │\n", "")
		for k, v := range config.Headers {
			value := f.masker.MaskIfSensitive(k, v)
			headerLine := fmt.Sprintf("  %s: %s", k, value)
			fmt.Printf("│   %-74s │\n", truncate(headerLine, 74))
		}
	}

	// Params
	if req.ProcessedParams != nil {
		fmt.Printf("│ Params:%-69s │\n", "")
		paramsJSON, _ := json.MarshalIndent(req.ProcessedParams, "  ", "  ")
		paramsLines := strings.Split(string(paramsJSON), "\n")
		for _, line := range paramsLines {
			fmt.Printf("│   %-74s │\n", truncate(line, 74))
		}
	} else {
		fmt.Printf("│ Params:  %-67s │\n", "[]")
	}

	// Bottom border
	fmt.Println("└" + strings.Repeat("─", 78) + "┘")
}

// FormatRequestJSON formats requests in JSON format
func (f *Formatter) FormatRequestJSON(requests []*types.Request) error {
	output := make([]map[string]interface{}, 0, len(requests))
	for _, req := range requests {
		output = append(output, map[string]interface{}{
			"name":    req.Name,
			"method":  req.Method,
			"params":  req.ProcessedParams,
			"url":     req.URL,
			"headers": req.Headers,
			"timeout": req.Timeout,
			"config":  req.Config,
		})
	}

	jsonBytes, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(jsonBytes))
	return nil
}

// FormatExecutionResults formats the results of request execution
func (f *Formatter) FormatExecutionResults(results []*types.ExecutionResult, jsonOutput bool) {
	if jsonOutput {
		f.formatExecutionResultsJSON(results)
		return
	}

	totalCount := len(results)
	successCount := 0
	failedCount := 0

	for i, result := range results {
		fmt.Printf("\n[%d/%d] Executing: %s\n", i+1, totalCount, result.Request.Name)

		if result.Error != nil {
			fmt.Printf("  ✗ Failed\n")
			fmt.Printf("  Duration: %dms\n", result.Duration.Milliseconds())
			fmt.Printf("  Error: %s\n", result.Error.Error())
			failedCount++
			continue
		}

		if result.Response.IsError() {
			fmt.Printf("  ✗ RPC Error\n")
			fmt.Printf("  Duration: %dms\n", result.Duration.Milliseconds())
			fmt.Printf("  Error Code: %d\n", result.Response.Error.Code)
			fmt.Printf("  Error Message: %s\n", result.Response.Error.Message)
			if result.Response.Error.Data != nil {
				dataJSON, _ := json.MarshalIndent(result.Response.Error.Data, "  ", "  ")
				fmt.Printf("  Error Data:\n  %s\n", string(dataJSON))
			}
			failedCount++
			continue
		}

		fmt.Printf("  ✓ Success\n")
		fmt.Printf("  Duration: %dms\n", result.Duration.Milliseconds())
		fmt.Printf("  Result:\n")

		// Pretty print result
		var resultObj interface{}
		if err := json.Unmarshal(result.Response.Result, &resultObj); err == nil {
			resultJSON, _ := json.MarshalIndent(resultObj, "  ", "  ")
			fmt.Printf("  %s\n", string(resultJSON))
		} else {
			fmt.Printf("  %s\n", string(result.Response.Result))
		}

		successCount++
	}

	// Summary
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Printf("Summary: %d total, %d successful, %d failed\n", totalCount, successCount, failedCount)
}

// formatExecutionResultsJSON formats execution results in JSON format
func (f *Formatter) formatExecutionResultsJSON(results []*types.ExecutionResult) {
	output := make([]map[string]interface{}, 0, len(results))

	for _, result := range results {
		resultMap := map[string]interface{}{
			"request":  result.Request.Name,
			"method":   result.Request.Method,
			"duration": result.Duration.Milliseconds(),
		}

		if result.Error != nil {
			resultMap["success"] = false
			resultMap["error"] = result.Error.Error()
		} else if result.Response.IsError() {
			resultMap["success"] = false
			resultMap["rpc_error"] = map[string]interface{}{
				"code":    result.Response.Error.Code,
				"message": result.Response.Error.Message,
				"data":    result.Response.Error.Data,
			}
		} else {
			resultMap["success"] = true
			var resultObj interface{}
			if err := json.Unmarshal(result.Response.Result, &resultObj); err == nil {
				resultMap["result"] = resultObj
			} else {
				resultMap["result"] = string(result.Response.Result)
			}
		}

		output = append(output, resultMap)
	}

	jsonBytes, _ := json.MarshalIndent(output, "", "  ")
	fmt.Println(string(jsonBytes))
}

// truncate truncates a string to a maximum length
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}
