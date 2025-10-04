package parser

import (
	"fmt"

	"jsonrpc/pkg/types"
)

// Validator validates parsed HCL files
type Validator struct{}

// NewValidator creates a new Validator
func NewValidator() *Validator {
	return &Validator{}
}

// Validate validates the HCL file for correctness
func (v *Validator) Validate(hclFile *types.HCLFile) error {
	// Validate all requests
	for _, req := range hclFile.Requests {
		if err := v.validateRequest(req, hclFile.Configs); err != nil {
			return err
		}
	}

	return nil
}

// validateRequest validates a single request
func (v *Validator) validateRequest(req *types.Request, configs map[string]*types.Config) error {
	// Check that the request has a method
	if req.Method == "" {
		return fmt.Errorf("request '%s' is missing required 'method' field", req.Name)
	}

	// Check if config reference exists
	if req.Config != "" {
		if _, exists := configs[req.Config]; !exists {
			return fmt.Errorf("request '%s' references non-existent config '%s'", req.Name, req.Config)
		}
	}

	return nil
}
