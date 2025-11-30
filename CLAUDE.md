# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**rpc-cli** is a Go CLI tool for executing JSON-RPC requests defined in HCL (HashiCorp Configuration Language) configuration files. It supports flexible configuration management, multiple output formats, and automatic masking of sensitive headers.

**Key Technologies**:
- Go 1.24
- Cobra (CLI framework)
- HCL/v2 (configuration parsing)
- cty (type conversion)
- Bubble Tea (TUI framework)
- Lipgloss (terminal styling)

## Common Development Commands

### Building and Running
```bash
# Build the binary
go build -o rpc-cli ./cmd/rpc-cli

# Run a command directly
go run ./cmd/rpc-cli ls requests.hcl
go run ./cmd/rpc-cli run requests.hcl --url https://example.com
go run ./cmd/rpc-cli tui requests.hcl

# Run TUI (interactive mode - cannot be tested in non-interactive shell)
./rpc-cli tui requests.hcl
./rpc-cli tui  # Auto-discovers HCL files
```

### Testing
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...

# Run tests with race detection
go test -race ./...

# Run specific package tests
go test ./internal/parser -v
go test ./internal/executor -v
go test ./internal/output -v
go test ./pkg/types -v

# Run a single test
go test ./internal/parser -run TestParseHCLFile -v
```

### Code Quality
```bash
# Run linter (requires golangci-lint)
golangci-lint run

# Format code
go fmt ./...

# Tidy dependencies
go mod tidy
```

### Release and Distribution
```bash
# Test release locally (doesn't publish)
goreleaser release --snapshot --clean

# View generated artifacts
ls -la dist/

# Create a release tag (triggers GitHub Actions release)
git tag -a v0.2.0 -m "Release v0.2.0"
git push origin v0.2.0
```

### Claude Commands
```bash
# Generate changelog based on current state and previous tag
/generate-changelog

# Generate changelog since specific tag
/generate-changelog v0.1.0
```

The `/generate-changelog` command parses git commits using Conventional Commits format and generates a formatted changelog entry ready for CHANGELOG.md. It automatically categorizes commits by type (Features, Bug Fixes, Enhancements, etc.).

## Architecture Overview

### Package Structure

**cmd/rpc-cli** - CLI entry point using Cobra framework
- Defines 5 commands: `ls`, `run`/`r`, `validate`/`v`, `tui`, `version`
- Handles CLI flags and argument parsing
- Orchestrates calls to internal packages
- Auto-discovers HCL files when no file argument provided

**internal/parser** - HCL file parsing and validation (25.5% coverage)
- `parser.go`: Main parsing logic, extracts config and request blocks
- `decoder.go`: Converts HCL attributes to Go types
- `converter.go`: Transforms cty.Value to native Go types
- `validator.go`: Validates parsed files, checks required fields

**internal/executor** - JSON-RPC execution (34.9% coverage)
- `executor.go`: HTTP client management, JSON-RPC protocol implementation
- `merger.go`: Configuration merging with priority order (CLI flags > request overrides > named config > default config)
- `helpers.go`: Utility functions for config resolution

**internal/output** - Result formatting and masking (10.1% coverage)
- `formatter.go`: Multiple output formats (table, detailed, JSON)
- `masker.go`: Sensitive header detection and masking

**internal/tui** - Interactive terminal UI using Bubble Tea
- `model.go`: Core TUI state management, event handling, Bubble Tea integration
- `view.go`: View rendering with viewport, JSON syntax highlighting
- `styles.go`: Visual styling and color schemes (default, dark, light, high-contrast)
- File selection screen for browsing HCL files
- Search/filter functionality for requests
- Multi-select with keyboard navigation (vim-style hjkl + arrows)
- Real-time execution results with colored output

**pkg/types** - Shared types (100% coverage)
- Core data structures: `Config`, `Request`, `HCLFile`, `JSONRPCRequest`, `ExecutionResult`, `EffectiveConfig`

### Configuration Priority (Highest to Lowest)
1. CLI flags (`--url`, `--header`, `--timeout`, `--config`)
2. Request-level overrides in HCL
3. Named config profile
4. Default config

### Data Flow
```
CLI command → Parser → Validator → ConfigMerger → Executor → Formatter → Output
```

## Key Design Patterns

### Configuration Merging
The `internal/executor/merger.go` implements a cascading configuration system where CLI flags override request-level overrides, which override named config profiles, which override the default config. This enables flexible per-request customization.

### Type Safety
The parser uses HCL/v2's cty type system and manually converts cty.Value to Go types via `converter.go`, ensuring type safety while maintaining flexibility for complex nested structures.

### Output Masking
Sensitive headers (Authorization, Token, API-Key, Secret, Password, Bearer) are automatically masked in output via case-insensitive keyword matching in `masker.go`.

### Error Handling
Errors are categorized at multiple levels: parse errors (HCL syntax), validation errors (missing fields), network errors (timeouts, connection failures), and RPC errors (JSON-RPC error responses).

### TUI Architecture
The TUI follows the Elm Architecture (Model-View-Update) via Bubble Tea:
- **Model**: Holds all application state (files, requests, selection, viewport, search)
- **Update**: Processes messages (key presses, window resize, file load, execution results)
- **View**: Renders the UI based on current model state

Key TUI patterns:
- **View States**: `ViewFileSelect`, `ViewList`, `ViewDetail`, `ViewResults`, `ViewHelp`
- **Viewport Management**: Uses Bubble Tea viewport for scrollable content areas
- **Search/Filter**: Real-time filtering with separate filtered request list and index mapping
- **State Transitions**: File selection → Request list → Details → Results
- **Content Separation**: Status lines rendered outside viewport to prevent spacing issues

## Testing Strategy

- **Unit tests** in `*_test.go` files use table-driven test patterns
- **Coverage targets**: Critical packages >80% (types at 100%), lower for output formatting
- **Test files** are co-located with implementation files
- Focus on error paths and edge cases, especially for public APIs

### Writing Tests
Follow the existing table-driven test pattern:
```go
tests := []struct {
    name    string
    input   interface{}
    want    interface{}
    wantErr bool
}{
    // test cases
}

for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        // test logic
    })
}
```

## Important Files and Conventions

- **requests.hcl**: Example HCL configuration file showing all supported features
- **ARCHITECTURE.md**: Detailed architecture guide with extension points
- **.goreleaser.yaml**: Multi-platform build configuration for Linux, macOS, Windows (amd64, arm64)
- **.github/workflows/test.yml**: CI/CD pipeline runs tests with race detection and coverage

## Release and Versioning

- Uses [GoReleaser](https://goreleaser.com/) for automated multi-platform builds
- Changelog is automatically generated from git commits using Conventional Commits
- Commit message format: `type(scope): description` (e.g., `feat: add new output format`, `fix(parser): handle nested objects`)
- GoReleaser groups commits: Features, Bug fixes, Enhancements, Others
- Excludes commits starting with `docs:`, `test:`, `chore:`, `ci:` from changelog

## Common Extension Points

### Adding a New Output Format
1. Add constant to `pkg/types/types.go` (OutputFormat type)
2. Implement formatter in `internal/output/formatter.go`
3. Add flag handling in `cmd/rpc-cli/main.go`

### Adding a New Config Source
1. Define override type in `pkg/types/types.go`
2. Implement merge logic in `internal/executor/merger.go`
3. Update configuration priority documentation

### Adding Request Preprocessing
1. Create interface in `pkg/types/types.go`
2. Implement in `internal/executor/executor.go`
3. Add configuration option

### Adding a New TUI View
1. Add view constant to `internal/tui/model.go` (View enum)
2. Implement `buildXContent()` method in `internal/tui/view.go` for viewport content
3. Implement `renderXView()` method in `internal/tui/view.go` for full view with header/footer
4. Add case to `Model.View()` switch statement
5. Add case to `updateViewportContent()` switch statement
6. Create `handleXKeys()` method for view-specific key handling
7. Add transition logic in relevant key handlers

### Modifying TUI Rendering
**Important**: The TUI has a specific rendering structure to avoid spacing issues:
- Status/info lines are rendered OUTSIDE the viewport (in the header area)
- Viewport content starts directly with list items (no leading newlines)
- Use `strings.TrimSpace(m.viewport.View())` when rendering viewport content
- Never add blank lines at the start of `buildXContent()` methods

## Performance Notes

- HTTP client is reused across multiple requests
- HCL parsing is lazy (only parses what's needed)
- Type conversion is direct from cty to Go without intermediate steps
- Maps and slices are reused where possible to minimize allocations

## Security Considerations

- Sensitive headers are automatically masked in output (never exposed in logs)
- All requests have configurable timeout limits
- All HCL input is validated before processing
- No hardcoded secrets in configuration examples
