package executor

import (
	"jsonrpc/pkg/config"
	"jsonrpc/pkg/types"
)

// GetConfigName returns the effective config name for a request
// This function is deprecated - use config.GetConfigName directly
// Maintained for backward compatibility
func GetConfigName(req *types.Request, overrides *types.CLIOverrides) string {
	return config.GetConfigName(req, overrides)
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
