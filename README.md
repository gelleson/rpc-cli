# rpc-cli

A command-line tool for executing JSON-RPC requests defined in HCL configuration files.

## Features

- 📝 **HCL Configuration**: Define JSON-RPC requests in easy-to-read HCL files
- 🎨 **Interactive TUI**: Advanced terminal UI with search, filtering, and JSON syntax highlighting
- 🔧 **Flexible Configuration**: Support for multiple environments and config profiles
- 🔐 **Security**: Automatic masking of sensitive headers in output
- ⏱️ **Performance Tracking**: Built-in duration tracking for all requests
- 🎯 **Selective Execution**: Run all requests or specific ones by name
- 📊 **Multiple Output Formats**: Table, detailed, and JSON output modes
- 🔄 **Config Overrides**: Override configurations at multiple levels (CLI, request, profile, default)
- 🔍 **Auto-discovery**: Automatically finds HCL files in current directory

## Installation

### Binary Releases

Download pre-built binaries from the [releases page](https://github.com/gelleson/rpc-cli/releases).

### Install Script

```bash
curl -sSL https://raw.githubusercontent.com/gelleson/rpc-cli/main/install.sh | bash
```

### Go Install

```bash
go install github.com/gelleson/rpc-cli/cmd/rpc-cli@latest
```

### Docker

```bash
docker pull ghcr.io/gelleson/rpc-cli:latest

# Run with HCL file mounted
docker run -v $(pwd)/requests.hcl:/requests.hcl ghcr.io/gelleson/rpc-cli:latest ls /requests.hcl
```

### Build from Source

```bash
# Clone the repository
git clone https://github.com/gelleson/rpc-cli.git
cd rpc-cli

# Build
go build -o rpc-cli ./cmd/rpc-cli

# Or install system-wide
go install ./cmd/rpc-cli
```

## Development

```bash
# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...

# Run tests with race detection
go test -race ./...
```

### Release Process

This project uses [GoReleaser](https://goreleaser.com/) for automated releases.

```bash
# Create a new tag
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0

# GoReleaser will automatically:
# - Build binaries for multiple platforms (Linux, macOS, Windows)
# - Create GitHub release with changelog
# - Publish Docker images to GitHub Container Registry
```

Local release testing:

```bash
# Test release locally (without publishing)
goreleaser release --snapshot --clean

# Check generated artifacts
ls -la dist/
```

## Quick Start

1. Create an HCL file with your requests (see `requests.hcl` for examples)
2. List available requests:
   ```bash
   rpc-cli ls requests.hcl
   ```
3. Execute a request:
   ```bash
   rpc-cli run requests.hcl get_balance
   ```

## Commands

### ls - List requests

List all requests or specific requests from an HCL file.

```bash
# List all requests
rpc-cli ls requests.hcl

# List specific requests
rpc-cli ls requests.hcl get_balance admin_call

# Show detailed information
rpc-cli ls requests.hcl get_balance --detailed

# JSON output
rpc-cli ls requests.hcl --json
```

### run - Execute requests

Execute all requests or specific requests from an HCL file.

```bash
# Execute all requests
rpc-cli run requests.hcl

# Execute specific requests
rpc-cli run requests.hcl get_balance

# Execute with URL override
rpc-cli run requests.hcl get_balance --url https://custom-rpc.com

# Execute with config profile
rpc-cli run requests.hcl get_balance --config production

# Execute with header overrides
rpc-cli run requests.hcl get_balance --header "Authorization: Bearer xyz"

# Execute with timeout override (in seconds)
rpc-cli run requests.hcl get_balance --timeout 60

# JSON output for scripting
rpc-cli run requests.hcl get_balance --json
```

### validate - Validate HCL syntax

Validate HCL file syntax and check for errors.

```bash
rpc-cli validate requests.hcl
```

### tui - Interactive Terminal UI

Launch an interactive terminal user interface for browsing and executing JSON-RPC requests.

```bash
# Auto-discover and select HCL files from current directory
rpc-cli tui

# Launch TUI with specific file (skips file selection)
rpc-cli tui requests.hcl
```

**File Selection:**
When running `rpc-cli tui` without arguments, an interactive file browser appears showing all HCL files in the current directory with their file sizes. Use arrow keys or `j/k` to navigate and press Enter to select a file.

**TUI Features:**
- 🔍 **Search/Filter** - Press `/` to search requests by name or method
- ⌨️ **Vim-style Navigation** - Use `hjkl` or arrow keys
- 📋 **Multi-select** - Space to toggle, `a` to select all, `A` to deselect all
- 🎨 **JSON Syntax Highlighting** - Color-coded responses
- 📊 **Real-time Results** - Execution results with response times
- ❓ **Help Modal** - Press `?` for keyboard shortcuts
- 🔄 **Viewport Scrolling** - Smooth scrolling for large content

**Keyboard Shortcuts:**
- `?` - Show help
- `/` - Search/filter requests
- `↑/k` - Move up
- `↓/j` - Move down
- `space` - Toggle selection
- `enter/l` - View details
- `r` - Run selected requests
- `a` - Select all
- `A` - Deselect all
- `ESC/h` - Go back / Clear search
- `q` - Quit

## HCL File Structure

### Config Blocks

Define reusable configurations. You can have:
- **Default config** (no label): Used when no other config is specified
- **Named configs**: Environment-specific configurations

```hcl
# Default config
config {
  url = "https://api.example.com/rpc"
  headers = {
    Content-Type = "application/json"
  }
  timeout = 30
}

# Production config
config "production" {
  url = "https://prod-api.example.com/rpc"
  headers = {
    Content-Type  = "application/json"
    Authorization = "Bearer prod_token"
  }
  timeout = 60
}

# Staging config
config "staging" {
  url = "https://staging-api.example.com/rpc"
  headers = {
    Content-Type = "application/json"
  }
  timeout = 45
}
```

### Request Blocks

Define individual JSON-RPC requests.

```hcl
# Simple request
request "get_balance" {
  method = "eth_getBalance"
  params = ["0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0", "latest"]
}

# Request using a named config
request "get_balance_prod" {
  config = "production"
  method = "eth_getBalance"
  params = ["0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045", "latest"]
}

# Request with URL override
request "external_call" {
  url = "https://external-api.com/rpc"
  method = "getData"
  params = []
}

# Request with header overrides
request "admin_call" {
  headers = {
    Authorization = "Bearer admin_token"
    X-Admin-Key   = "secret123"
  }
  method = "admin.getStats"
  params = []
}

# Request with timeout override
request "long_query" {
  timeout = 120
  method = "heavy.computation"
  params = [{ data = "large_dataset" }]
}

# Complex nested parameters
request "batch_transfer" {
  method = "wallet.batchTransfer"
  params = [
    {
      to     = "0x123..."
      amount = 100
      token  = "USDT"
    },
    {
      to     = "0x456..."
      amount = 250
      token  = "USDC"
    }
  ]
}
```

## Configuration Override Priority

Configurations are merged in the following order (highest to lowest priority):

1. **CLI flags** (`--url`, `--header`, `--timeout`, `--config`)
2. **Request-level overrides** (`url`, `headers`, `timeout` in request block)
3. **Named config profile** (if `config = "name"` specified)
4. **Default config** (unlabeled config block)

## Output Examples

### Table Output (ls)

```
NAME                  METHOD                      CONFIG      PARAMS
---------------------------------------------------------------------------------
get_balance          eth_getBalance              default     2
admin_call           admin.getStats              custom      0
batch_transfer       wallet.batchTransfer        default     3
```

### Detailed Output (ls --detailed)

```
┌──────────────────────────────────────────────────────────────────────────────┐
│ get_balance                                                                  │
├──────────────────────────────────────────────────────────────────────────────┤
│ Method:  eth_getBalance                                                      │
│ URL:     https://api.example.com/rpc                                        │
│ Config:  default                                                             │
│ Timeout: 30s                                                                 │
│ Params:                                                                      │
│   ["0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0", "latest"]                   │
└──────────────────────────────────────────────────────────────────────────────┘
```

### Execution Output (run)

```
[1/2] Executing: get_balance
  ✓ Success
  Duration: 245ms
  Result:
  "0x1bc16d674ec80000"

[2/2] Executing: get_block
  ✓ Success
  Duration: 312ms
  Result:
  {
    "number": "0x12a4b5c",
    "hash": "0xabc123...",
    ...
  }

=============================================================
Summary: 2 total, 2 successful, 0 failed
```

## Security Features

### Sensitive Header Masking

The tool automatically masks sensitive information in output. Headers containing any of the following keywords are masked:
- authorization
- token
- api-key
- apikey
- secret
- password
- bearer

Example:
```
Authorization: Bear****
X-API-Key: abc1****
```

## Error Handling

The tool provides clear error messages for:
- **Parse errors**: Invalid HCL syntax
- **Network errors**: Connection failures, timeouts
- **RPC errors**: JSON-RPC error responses with error codes and messages
- **Validation errors**: Missing required fields, invalid config references

## Examples

### Working with Ethereum RPC

```hcl
config {
  url = "https://eth-mainnet.g.alchemy.com/v2/demo"
  headers = {
    Content-Type = "application/json"
  }
  timeout = 30
}

request "get_block_number" {
  method = "eth_blockNumber"
  params = []
}

request "get_balance" {
  method = "eth_getBalance"
  params = ["0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0", "latest"]
}

request "get_transaction" {
  method = "eth_getTransactionByHash"
  params = ["0x88df016429689c079f3b2f6ad39fa052532c56795b733da78a91ebe6a713944b"]
}
```

### Testing Different Environments

```bash
# Test with staging config
rpc-cli run requests.hcl get_balance --config staging

# Test with production config
rpc-cli run requests.hcl get_balance --config production

# Override URL for testing
rpc-cli run requests.hcl get_balance --url https://test-rpc.com
```

### Scripting with JSON Output

```bash
# Get result in JSON for processing
result=$(rpc-cli run requests.hcl get_balance --json)
echo $result | jq '.[] | .result'
```

## Project Structure

```
rpc-cli/
├── cmd/
│   └── rpc-cli/
│       └── main.go              # CLI entry point
├── internal/
│   ├── executor/
│   │   ├── executor.go          # JSON-RPC execution logic
│   │   ├── merger.go            # Configuration merging
│   │   ├── helpers.go           # Helper functions
│   │   ├── executor_test.go
│   │   ├── merger_test.go
│   │   └── helpers_test.go
│   ├── output/
│   │   ├── formatter.go         # Output formatting
│   │   ├── masker.go            # Sensitive data masking
│   │   ├── formatter_test.go
│   │   └── masker_test.go
│   └── parser/
│       ├── parser.go            # HCL file parsing
│       ├── decoder.go           # Attribute decoding
│       ├── converter.go         # Cty to Go conversion
│       ├── validator.go         # HCL validation
│       ├── parser_test.go
│       ├── converter_test.go
│       └── validator_test.go
├── pkg/
│   └── types/
│       ├── types.go             # Shared type definitions
│       └── types_test.go
├── requests.hcl                 # Example HCL file
├── go.mod
├── go.sum
└── README.md
```

### Architecture

The project follows a clean, modular architecture:

- **cmd/**: Contains the CLI application entry point
- **internal/**: Internal packages not meant for external use
  - **executor**: Handles JSON-RPC request execution and configuration merging
  - **output**: Manages all output formatting and sensitive data masking
  - **parser**: HCL file parsing, validation, and type conversion
- **pkg/types**: Shared types used across the application

Each package is:
- **Self-contained**: Minimal dependencies on other packages
- **Well-tested**: Comprehensive unit tests with >80% coverage
- **Documented**: Clear interfaces and function documentation

## Dependencies

- [github.com/hashicorp/hcl/v2](https://github.com/hashicorp/hcl) - HCL parsing
- [github.com/spf13/cobra](https://github.com/spf13/cobra) - CLI framework
- [github.com/zclconf/go-cty](https://github.com/zclconf/go-cty) - Cty value conversion

## License

This project is provided as-is for educational and development purposes.
