package output

import (
	"encoding/json"
	"fmt"
	"strings"

	"jsonrpc/pkg/config"
	"jsonrpc/pkg/constants"
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
func (f *Formatter) FormatRequestList(
	hclFile *types.HCLFile,
	requests []*types.Request,
	overrides *types.CLIOverrides,
) {
	// Print header
	fmt.Printf("%-25s %-30s %-15s %-10s\n", "NAME", "METHOD", "CONFIG", "PARAMS")
	fmt.Println(strings.Repeat("-", 85))

	// Print each request
	for _, req := range requests {
		configName := config.GetConfigName(req, overrides)
		paramCount := CountParams(req.ProcessedParams)

		fmt.Printf("%-25s %-30s %-15s %-10d\n",
			truncate(req.Name, constants.MaxNameLength),
			truncate(req.Method, constants.MaxMethodLength),
			truncate(configName, constants.MaxConfigLength),
			paramCount,
		)
	}
}

// FormatRequestDetailed formats requests in detailed boxed format
func (f *Formatter) FormatRequestDetailed(
	hclFile *types.HCLFile,
	requests []*types.Request,
	overrides *types.CLIOverrides,
) {
	for i, req := range requests {
		if i > 0 {
			fmt.Println()
		}
		f.formatSingleRequestDetailed(hclFile, req, overrides)
	}
}

// formatSingleRequestDetailed formats a single request in detailed format
func (f *Formatter) formatSingleRequestDetailed(
	hclFile *types.HCLFile,
	req *types.Request,
	overrides *types.CLIOverrides,
) {
	configMgr := config.NewManager()
	config := configMgr.BuildForRequest(hclFile, req, overrides)
	configName := configMgr.GetConfigNameForRequest(hclFile, req, overrides)

	// Top border
	fmt.Println("┌" + strings.Repeat("─", constants.BoxWidth) + "┐")

	// Title
	fmt.Printf("│ %-76s │\n", req.Name)
	fmt.Println("├" + strings.Repeat("─", constants.BoxWidth) + "┤")

	// Method
	fmt.Printf("│ Method:  %-67s │\n", req.Method)

	// URL
	fmt.Printf("│ URL:     %-67s │\n", truncate(config.URL, constants.BoxContentWidth-9))

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
			fmt.Printf("│   %-74s │\n", truncate(headerLine, constants.BoxContentWidth-2))
		}
	}

	// Params
	if req.ProcessedParams != nil {
		fmt.Printf("│ Params:%-69s │\n", "")
		paramsJSON, _ := json.MarshalIndent(req.ProcessedParams, "  ", "  ")
		paramsLines := strings.Split(string(paramsJSON), "\n")
		for _, line := range paramsLines {
			fmt.Printf("│   %-74s │\n", truncate(line, constants.BoxContentWidth-2))
		}
	} else {
		fmt.Printf("│ Params:  %-67s │\n", "[]")
	}

	// Bottom border
	fmt.Println("└" + strings.Repeat("─", constants.BoxWidth) + "┘")
}

// FormatRequestJSON formats requests in JSON format
func (f *Formatter) FormatRequestJSON(requests []*types.Request) error {
	output := make([]map[string]any, 0, len(requests))
	for _, req := range requests {
		output = append(output, map[string]any{
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
		var resultObj any
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
	output := make([]map[string]any, 0, len(results))

	for _, result := range results {
		resultMap := map[string]any{
			"request":  result.Request.Name,
			"method":   result.Request.Method,
			"duration": result.Duration.Milliseconds(),
		}

		if result.Error != nil {
			resultMap["success"] = false
			resultMap["error"] = result.Error.Error()
		} else if result.Response.IsError() {
			resultMap["success"] = false
			resultMap["rpc_error"] = map[string]any{
				"code":    result.Response.Error.Code,
				"message": result.Response.Error.Message,
				"data":    result.Response.Error.Data,
			}
		} else {
			resultMap["success"] = true
			var resultObj any
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

// CountParams returns the number of parameters in a request
func CountParams(params any) int {
	if params == nil {
		return 0
	}

	switch v := params.(type) {
	case []any:
		return len(v)
	case map[string]any:
		return len(v)
	default:
		return 1
	}
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
