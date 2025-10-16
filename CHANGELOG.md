# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.0] - 2025-10-16

### Added
- Initial release of rpc-cli
- HCL-based JSON-RPC request configuration and execution
- Multiple output formats: table, detailed, and JSON
- Configuration management with environment-specific profiles
- Automatic masking of sensitive headers in output
- Request filtering by name
- `ls` command to list available requests
- `run` command to execute JSON-RPC requests
- `r` alias for `run` command
- `validate` command to validate HCL syntax
- `v` alias for `validate` command
- `version` command to display version information
- Configuration override support via CLI flags (`--url`, `--header`, `--timeout`, `--config`)
- Request-level configuration overrides in HCL
- Duration tracking for all executed requests
- Timeout protection for all requests
- Input validation for HCL files
- Comprehensive error handling and validation messages
- Support for multiple architectures: amd64, arm64
- Support for multiple platforms: Linux, macOS, Windows
- Installation via binary download, install script, go install, and Docker

### Security
- Sensitive header detection and automatic masking (Authorization, Token, API-Key, Secret, Password, Bearer)
- Timeout protection on all HTTP requests
- Input validation on all HCL configurations

### Documentation
- README with installation instructions and quick start guide
- ARCHITECTURE.md with design principles and extension points
- Example configuration file (requests.hcl)

### Infrastructure
- GoReleaser configuration for automated multi-platform builds
- GitHub Actions CI/CD workflow for testing and linting
- golangci-lint integration for code quality

## Unreleased Features and Improvements

### Planned
- Extended output formatting options
- Configuration file validation enhancements
- Batch request execution with parallel processing
- Request result caching
- Request scheduling and automation

---

[Unreleased]: https://github.com/gelleson/rpc-cli/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/gelleson/rpc-cli/releases/tag/v0.1.0
