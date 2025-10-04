# Architecture Guide

## Overview

`rpc-cli` is a modular, well-tested CLI tool for executing JSON-RPC requests defined in HCL configuration files. The project follows Go best practices with a clean separation of concerns.

## Design Principles

1. **Modularity**: Each package has a single, well-defined responsibility
2. **Testability**: All packages include comprehensive unit tests
3. **Maintainability**: Clear interfaces and minimal coupling between packages
4. **Extensibility**: Easy to add new features without modifying existing code

## Package Structure

### cmd/rpc-cli (CLI Entry Point)

**Responsibility**: Command-line interface and user interaction

- Defines CLI commands using Cobra
- Handles command-line flags and arguments
- Orchestrates calls to internal packages
- Provides user-friendly error messages

### pkg/types (Shared Types)

**Responsibility**: Common data structures used across packages

**Key Types**:
- `Config`: Configuration block definition
- `Request`: JSON-RPC request definition
- `HCLFile`: Parsed HCL file structure
- `JSONRPCRequest/Response`: JSON-RPC protocol types
- `ExecutionResult`: Request execution result
- `EffectiveConfig`: Merged configuration
- `CLIOverrides`: Command-line overrides

**Coverage**: 100%

### internal/parser (HCL Parsing)

**Responsibility**: Parse and validate HCL configuration files

**Components**:

1. **Parser** (`parser.go`)
   - Main parsing logic
   - Extracts config and request blocks
   - Handles both labeled and unlabeled blocks

2. **Decoder** (`decoder.go`)
   - Decodes HCL attributes to Go types
   - Type-safe conversion (string, int, map[string]string)

3. **Converter** (`converter.go`)
   - Converts cty.Value to native Go types
   - Handles primitives, lists, maps, and nested structures

4. **Validator** (`validator.go`)
   - Validates parsed HCL files
   - Checks required fields
   - Verifies config references

**Coverage**: 25.5%

### internal/executor (JSON-RPC Execution)

**Responsibility**: Execute JSON-RPC requests with merged configurations

**Components**:

1. **Executor** (`executor.go`)
   - Main execution logic
   - HTTP client management
   - JSON-RPC protocol implementation
   - Timeout handling

2. **ConfigMerger** (`merger.go`)
   - Merges configurations from multiple sources
   - Implements priority order:
     1. CLI flags (highest)
     2. Request-level overrides
     3. Named config profile
     4. Default config (lowest)

3. **Helpers** (`helpers.go`)
   - Utility functions
   - Config name resolution
   - Parameter counting

**Coverage**: 34.9%

### internal/output (Output Formatting)

**Responsibility**: Format and display results to users

**Components**:

1. **Formatter** (`formatter.go`)
   - Multiple output formats (table, detailed, JSON)
   - Request list formatting
   - Execution result formatting

2. **SensitiveMasker** (`masker.go`)
   - Detects sensitive headers
   - Masks sensitive values
   - Case-insensitive keyword matching

**Coverage**: 10.1%

## Data Flow

```
1. User runs CLI command
   ↓
2. cmd/rpc-cli parses flags and arguments
   ↓
3. parser.Parser reads and parses HCL file
   ↓
4. parser.Validator validates the parsed data
   ↓
5. executor.Executor builds effective config (merger.ConfigMerger)
   ↓
6. executor.Executor executes JSON-RPC request
   ↓
7. output.Formatter formats and displays results
   ↓
8. User sees output
```

## Configuration Priority

Configurations are merged in the following order (highest to lowest):

1. **CLI Flags**: `--url`, `--header`, `--timeout`, `--config`
2. **Request-Level**: `url`, `headers`, `timeout` in request block
3. **Named Config**: Referenced via `config = "name"`
4. **Default Config**: Unlabeled config block

Example:
```hcl
# Default config
config {
  url = "https://default.example.com"
  timeout = 30
}

# Named config
config "production" {
  url = "https://prod.example.com"
  timeout = 60
}

# Request with overrides
request "my_request" {
  config = "production"  # Uses production config
  timeout = 120          # Overrides timeout to 120
  method = "test"
  params = []
}
```

When executed:
- Base: default config (url: default, timeout: 30)
- Apply: production config (url: prod, timeout: 60)
- Apply: request override (timeout: 120)
- Final: url=prod, timeout=120

If run with `--timeout 90`:
- Final: url=prod, timeout=90 (CLI overrides all)

## Error Handling

The application handles errors at multiple levels:

1. **Parse Errors**: HCL syntax errors, invalid types
2. **Validation Errors**: Missing required fields, invalid references
3. **Network Errors**: Connection failures, timeouts
4. **RPC Errors**: JSON-RPC error responses

Each error type provides clear, actionable messages to the user.

## Testing Strategy

### Unit Tests
- Each package has comprehensive unit tests
- Test files co-located with implementation files
- Table-driven tests for multiple scenarios

### Test Coverage Goals
- **Critical Packages**: >80% coverage (types, core logic)
- **Integration Points**: Test error paths and edge cases
- **Public APIs**: All exported functions tested

### Running Tests
```bash
# All tests
go test ./...

# With coverage
go test -cover ./...

# Specific package
go test ./internal/parser -v

# With race detection
go test -race ./...
```

## Extending the Tool

### Adding a New Output Format

1. Add format constant to `pkg/types/types.go`:
   ```go
   const OutputXML OutputFormat = "xml"
   ```

2. Implement formatter in `internal/output/formatter.go`:
   ```go
   func (f *Formatter) FormatRequestXML(requests []*types.Request) error {
       // Implementation
   }
   ```

3. Add flag handling in `cmd/rpc-cli/main.go`

### Adding a New Config Source

1. Define new override type in `pkg/types/types.go`
2. Implement merge logic in `internal/executor/merger.go`
3. Update priority documentation

### Adding Request Preprocessing

1. Create new interface in `pkg/types/types.go`:
   ```go
   type RequestPreprocessor interface {
       Preprocess(*Request) error
   }
   ```

2. Implement in `internal/executor/executor.go`
3. Add configuration option

## Performance Considerations

1. **HTTP Client Reuse**: Single HTTP client instance per executor
2. **Lazy Parsing**: Only parse what's needed
3. **Efficient Type Conversion**: Direct cty to Go conversion without intermediate steps
4. **Minimal Allocations**: Reuse maps and slices where possible

## Security Considerations

1. **Sensitive Data Masking**: Automatic masking of auth headers
2. **No Secret Logging**: Sensitive values never logged
3. **Timeout Protection**: All requests have timeout limits
4. **Input Validation**: All HCL input validated before use

## Future Enhancements

Potential areas for improvement:

1. **Concurrent Execution**: Run multiple requests in parallel
2. **Request Templating**: Variable substitution in requests
3. **Response Validation**: JSON schema validation
4. **Request History**: Save and replay requests
5. **Interactive Mode**: REPL for building requests
6. **Plugins**: Custom request/response processors
