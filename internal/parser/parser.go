package parser

import (
	"fmt"
	"os"

	"jsonrpc/pkg/types"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

// Parser handles HCL file parsing
type Parser struct {
	hclParser *hclparse.Parser
}

// New creates a new Parser instance
func New() *Parser {
	return &Parser{
		hclParser: hclparse.NewParser(),
	}
}

// ParseFile parses an HCL file and returns the structured data
func (p *Parser) ParseFile(filename string) (*types.HCLFile, error) {
	src, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	file, diags := p.hclParser.ParseHCL(src, filename)
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to parse HCL file: %s", diags.Error())
	}

	result := types.NewHCLFile()

	// Extract all blocks from the HCL body
	blocks, diags := p.extractBlocks(file.Body)
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to extract blocks: %s", diags.Error())
	}

	// Parse config blocks
	for _, block := range blocks {
		if block.Type == "config" {
			config, err := p.parseConfigBlock(block)
			if err != nil {
				return nil, err
			}

			configName := p.getConfigName(block)
			result.Configs[configName] = config
		}
	}

	// Parse request blocks
	for _, block := range blocks {
		if block.Type == "request" {
			request, err := p.parseRequestBlock(block)
			if err != nil {
				return nil, err
			}
			result.Requests = append(result.Requests, request)
		}
	}

	return result, nil
}

// extractBlocks extracts all blocks from an HCL body
func (p *Parser) extractBlocks(body hcl.Body) (hcl.Blocks, hcl.Diagnostics) {
	// Check if this is an hclsyntax.Body (native syntax)
	syntaxBody, ok := body.(*hclsyntax.Body)
	if !ok {
		return nil, hcl.Diagnostics{
			&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Unsupported body type",
				Detail:   "Only HCL native syntax is supported",
			},
		}
	}

	// Directly access blocks without schema restrictions
	var blocks hcl.Blocks
	for _, block := range syntaxBody.Blocks {
		blocks = append(blocks, block.AsHCLBlock())
	}

	return blocks, nil
}

// getConfigName returns the config name from block labels
func (p *Parser) getConfigName(block *hcl.Block) string {
	if len(block.Labels) > 0 {
		return block.Labels[0]
	}
	return "default"
}

// parseConfigBlock parses a config block
func (p *Parser) parseConfigBlock(block *hcl.Block) (*types.Config, error) {
	config := types.NewConfig()
	decoder := NewAttributeDecoder()

	schema := &hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{Name: "url"},
			{Name: "headers"},
			{Name: "timeout"},
		},
	}

	content, diags := block.Body.Content(schema)
	if diags.HasErrors() {
		configName := p.getConfigName(block)
		return nil, fmt.Errorf("failed to decode config '%s': %s", configName, diags.Error())
	}

	// Decode URL
	if attr, exists := content.Attributes["url"]; exists {
		if err := decoder.DecodeString(attr, &config.URL); err != nil {
			return nil, err
		}
	}

	// Decode headers
	if attr, exists := content.Attributes["headers"]; exists {
		if err := decoder.DecodeStringMap(attr, &config.Headers); err != nil {
			return nil, err
		}
	}

	// Decode timeout
	if attr, exists := content.Attributes["timeout"]; exists {
		if err := decoder.DecodeInt(attr, &config.Timeout); err != nil {
			return nil, err
		}
	}

	return config, nil
}

// parseRequestBlock parses a request block
func (p *Parser) parseRequestBlock(block *hcl.Block) (*types.Request, error) {
	if len(block.Labels) == 0 {
		return nil, fmt.Errorf("request block must have a name label")
	}

	request := types.NewRequest(block.Labels[0])
	decoder := NewAttributeDecoder()

	schema := &hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{Name: "method", Required: true},
			{Name: "params"},
			{Name: "url"},
			{Name: "headers"},
			{Name: "timeout"},
			{Name: "config"},
		},
	}

	content, diags := block.Body.Content(schema)
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to decode request '%s': %s", request.Name, diags.Error())
	}

	// Decode method (required)
	if attr, exists := content.Attributes["method"]; exists {
		if err := decoder.DecodeString(attr, &request.Method); err != nil {
			return nil, err
		}
	}

	// Decode params (complex type)
	if attr, exists := content.Attributes["params"]; exists {
		val, attrDiags := attr.Expr.Value(nil)
		if attrDiags.HasErrors() {
			return nil, fmt.Errorf("failed to decode params: %s", attrDiags.Error())
		}
		request.Params = val
		request.ProcessedParams = ConvertCtyToGo(val)
	}

	// Decode URL
	if attr, exists := content.Attributes["url"]; exists {
		if err := decoder.DecodeString(attr, &request.URL); err != nil {
			return nil, err
		}
	}

	// Decode headers
	if attr, exists := content.Attributes["headers"]; exists {
		if err := decoder.DecodeStringMap(attr, &request.Headers); err != nil {
			return nil, err
		}
	}

	// Decode timeout
	if attr, exists := content.Attributes["timeout"]; exists {
		if err := decoder.DecodeInt(attr, &request.Timeout); err != nil {
			return nil, err
		}
	}

	// Decode config reference
	if attr, exists := content.Attributes["config"]; exists {
		if err := decoder.DecodeString(attr, &request.Config); err != nil {
			return nil, err
		}
	}

	return request, nil
}
