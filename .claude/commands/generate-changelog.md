# Generate Changelog

Generate a changelog based on git commits between a previous tag and the current state (HEAD). Uses Conventional Commits format to categorize changes.

## Usage

```bash
/generate-changelog [previous-tag]
```

## Parameters

- `previous-tag` (optional): The git tag to compare against. Defaults to the most recent tag if not specified.

## Description

This command:
1. Fetches git log between the specified previous tag and HEAD
2. Parses commits using Conventional Commits format (feat:, fix:, docs:, etc.)
3. Groups commits by type (Features, Bug Fixes, Enhancements, etc.)
4. Generates a formatted changelog entry ready for CHANGELOG.md
5. Shows the generated entry for review before updating

The changelog entry includes:
- Version number (based on semantic versioning)
- Release date
- Categorized commits with descriptions
- Reference to full changelog on GitHub

## Example

```bash
# Generate changelog since last tag
/generate-changelog

# Generate changelog since specific tag
/generate-changelog v0.1.0
```

## Conventional Commits Format

Supported commit types:
- `feat:` - New features
- `fix:` - Bug fixes
- `docs:` - Documentation changes
- `refactor:` - Code refactoring
- `perf:` - Performance improvements
- `test:` - Test additions/changes
- `chore:` - Build/dependency changes
- `ci:` - CI/CD changes

Scopes are optional: `feat(parser): add new feature`

## Implementation Notes

The command uses `git log --pretty=format` to extract commit information and parses the raw output to build the changelog. It respects the GoReleaser configuration for commit filtering and grouping.