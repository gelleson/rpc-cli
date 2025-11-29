# Linting and Pre-commit Setup

This document describes the linting and pre-commit hooks configuration for the rpc-cli project.

## Overview

The project uses a comprehensive linting setup with:
- **golangci-lint**: Go code quality and static analysis
- **pre-commit hooks**: Automated code quality checks before commits
- **commitlint**: Conventional commits validation
- **Additional tools**: Shell check, spell check, secret detection

## Files Added

- `.golangci.yml` - Production-ready golangci-lint configuration
- `.pre-commit-config.yaml` - Pre-commit hooks configuration
- `.commitlint.yaml` - Conventional commits validation rules
- `.codespellignore` - Common words to ignore in spell checking
- `LINTING.md` - This documentation file

## Pre-commit Hooks

### Installed Hooks

1. **Basic file checks**
   - `trailing-whitespace` - Remove trailing whitespace
   - `end-of-file-fixer` - Ensure files end with newline
   - `check-yaml` - Validate YAML syntax
   - `check-added-large-files` - Prevent large files (>1MB)
   - `check-case-conflict` - Prevent case-insensitive filename conflicts
   - `check-merge-conflict` - Detect merge conflict markers
   - `check-json` - Validate JSON syntax
   - `pretty-format-json` - Auto-format JSON files
   - `mixed-line-ending` - Ensure consistent line endings
   - `no-commit-to-branch` - Prevent commits to main/master
   - `requirements-txt-fixer` - Auto-sort requirements
   - `check-toml` - Validate TOML syntax
   - `check-xml` - Validate XML syntax

2. **Go-specific hooks**
   - `go-fmt` - Format Go code using gofmt
   - `go-mod-tidy` - Clean up go.mod and go.sum
   - `go-test` - Run Go tests
   - `go-build` - Verify Go code builds
   - `golangci-lint` - Run comprehensive Go linting

3. **Commit message validation**
   - `commitlint` - Enforce conventional commit format

4. **Additional quality checks**
   - `shellcheck` - Shell script linting
   - `codespell` - Spell checking
   - `detect-secrets` - Detect potential secrets in code

### Usage

```bash
# Install hooks (already done)
pre-commit install
pre-commit install --hook-type commit-msg

# Run all hooks on all files
pre-commit run --all-files

# Run specific hook
pre-commit run golangci-lint --all-files

# Skip hooks (not recommended, use with care)
git commit --no-verify
```

## golangci-lint Configuration

### Enabled Linters

- **Core linters**: govet, errcheck, staticcheck, unused, ineffassign
- **Security**: gosec (security vulnerability detection)
- **Style**: goconst (magic strings), misspell (spelling)
- **Complexity**: gocyclo (cyclomatic complexity), dupl (code duplication)
- **Organization**: gochecknoinits (prevent init functions)
- **Formatting**: lll (line length limit - 120 chars)

### Key Settings

- **Line length**: 120 characters
- **Complexity threshold**: 15 for gocyclo
- **Duplication threshold**: 100 tokens for dupl
- **Timeout**: 5 minutes for analysis
- **Format**: Colored output with line numbers

### Exclusions

- Test files are excluded from some strict linters
- Shadow variable checking excludes common `err` patterns
- Known security false positives are excluded

## Conventional Commits

### Required Format

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

### Types

- `feat` - New feature
- `fix` - Bug fix
- `docs` - Documentation only changes
- `style` - Changes that don't affect code meaning
- ` refactor` - Code change that neither fixes bug nor adds feature
- `perf` - Performance improvement
- `test` - Adding missing tests or correcting existing tests
- `build` - Build system or external dependencies
- `ci` - CI configuration files and scripts
- `chore` - Other changes that don't modify src or test files
- `revert` - Reverts a previous commit
- `wip` - Work in progress

### Examples

- `feat(parser): add support for nested HCL blocks`
- `fix(executor): resolve timeout issue with large requests`
- `docs: update README with installation instructions`

## Current Issues Found

The linting setup identified the following issues in the codebase:

1. **Line length violations** (8 issues) - Several function signatures exceed 120 characters
2. **Magic string** (1 issue) - String "default" appears 8 times, should be a constant
3. **Security warning** (1 issue) - File reading via variable (acceptable in this context)

## Integration with Development Workflow

### Before Committing

```bash
# Run checks manually (optional, pre-commit will run automatically)
go fmt ./...
go test ./...
golangci-lint run
```

### CI/CD Integration

The pre-commit hooks are configured to run automatically before each commit, ensuring:
- Code formatting consistency
- No syntax or basic logical errors
- Conventional commit messages
- No accidental inclusion of secrets or large files
- Consistent file formatting (JSON, YAML, TOML)

### IDE Integration

Most IDEs can be configured to:
- Run gofmt on save
- Show golangci-lint warnings inline
- Integrate pre-commit hooks with Git operations

## Maintenance

### Updating Tools

```bash
# Update pre-commit hook versions
pre-commit autoupdate

# Update golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### Adding New Rules

1. Update `.golangci.yml` for new Go linters
2. Update `.pre-commit-config.yaml` for new pre-commit hooks
3. Update `.commitlint.yaml` for new commit rules
4. Test changes with `pre-commit run --all-files`

### Troubleshooting

- **Hook timeouts**: Increase timeout in `.golangci.yml`
- **False positives**: Add exclusions to the issues section
- **Version conflicts**: Ensure compatible versions in configuration files

## Benefits

1. **Code Quality**: Consistent formatting and style across the codebase
2. **Early Detection**: Catch issues before they reach CI/CD
3. **Security**: Automatic detection of potential secrets and vulnerabilities
4. **Documentation**: Enforced conventional commits for better changelog generation
5. **Team Collaboration**: Shared standards reduce code review friction
6. **Maintainability**: Automated checks prevent technical debt accumulation
