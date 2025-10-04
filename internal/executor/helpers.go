package executor

import "jsonrpc/pkg/types"

// GetConfigName returns the effective config name for a request
func GetConfigName(req *types.Request, overrides *types.CLIOverrides) string {
	if overrides != nil && overrides.Config != "" {
		return overrides.Config
	}
	if req.Config != "" {
		return req.Config
	}
	return "default"
}

// CountParams returns the number of parameters in a request
func CountParams(params interface{}) int {
	if params == nil {
		return 0
	}

	switch v := params.(type) {
	case []interface{}:
		return len(v)
	case map[string]interface{}:
		return len(v)
	default:
		return 1
	}
}
