#!/bin/bash

# Changelog Generator for rpc-cli
# Generates changelog entries based on git commits between tags using Conventional Commits format

REPO_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$REPO_DIR"

# Get the previous tag (default to most recent tag if not specified)
PREV_TAG=""
if [ -n "$1" ]; then
    PREV_TAG="$1"
else
    # Try to get the most recent tag
    PREV_TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo "")
fi

CURRENT_DATE=$(date +%Y-%m-%d)
COMMIT_RANGE="HEAD"

# If we have a previous tag, use it; otherwise use all history
if [ -n "$PREV_TAG" ]; then
    COMMIT_RANGE="$PREV_TAG..HEAD"
    echo "Generating changelog between $PREV_TAG and HEAD..."
else
    echo "Generating changelog for all commits (no previous tag)..."
fi
echo ""

# Function to get commits of a specific type
get_commits_by_type() {
    local regex="$1"
    if [ -n "$PREV_TAG" ]; then
        git log "$PREV_TAG"..HEAD --pretty=format:"%s" 2>/dev/null | grep -E "$regex" || true
    else
        git log HEAD --pretty=format:"%s" 2>/dev/null | grep -E "$regex" || true
    fi
}

# Collect commits by type
features=$(get_commits_by_type "^feat(\([^)]*\))?:")
fixes=$(get_commits_by_type "^fix(\([^)]*\))?:")
enhancements=$(get_commits_by_type "^enhance(\([^)]*\))?:")
refactors=$(get_commits_by_type "^refactor(\([^)]*\))?:")
perf=$(get_commits_by_type "^perf(\([^)]*\))?:")
docs=$(get_commits_by_type "^docs(\([^)]*\))?:")

# Generate changelog entry
echo "## [Unreleased] - $CURRENT_DATE"
echo ""

if [ -n "$features" ]; then
    echo "### Added"
    echo "$features" | while IFS= read -r line; do
        # Clean up commit message: remove type prefix and scope
        clean_msg=$(echo "$line" | sed -E 's/^[a-z]+(\([^)]*\))?:[ ]*//')
        [ -n "$clean_msg" ] && echo "- $clean_msg"
    done
    echo ""
fi

if [ -n "$fixes" ]; then
    echo "### Fixed"
    echo "$fixes" | while IFS= read -r line; do
        clean_msg=$(echo "$line" | sed -E 's/^[a-z]+(\([^)]*\))?:[ ]*//')
        [ -n "$clean_msg" ] && echo "- $clean_msg"
    done
    echo ""
fi

if [ -n "$enhancements" ]; then
    echo "### Enhanced"
    echo "$enhancements" | while IFS= read -r line; do
        clean_msg=$(echo "$line" | sed -E 's/^[a-z]+(\([^)]*\))?:[ ]*//')
        [ -n "$clean_msg" ] && echo "- $clean_msg"
    done
    echo ""
fi

if [ -n "$refactors" ]; then
    echo "### Changed"
    echo "$refactors" | while IFS= read -r line; do
        clean_msg=$(echo "$line" | sed -E 's/^[a-z]+(\([^)]*\))?:[ ]*//')
        [ -n "$clean_msg" ] && echo "- $clean_msg"
    done
    echo ""
fi

if [ -n "$perf" ]; then
    echo "### Performance"
    echo "$perf" | while IFS= read -r line; do
        clean_msg=$(echo "$line" | sed -E 's/^[a-z]+(\([^)]*\))?:[ ]*//')
        [ -n "$clean_msg" ] && echo "- $clean_msg"
    done
    echo ""
fi

if [ -n "$docs" ]; then
    echo "### Documentation"
    echo "$docs" | while IFS= read -r line; do
        clean_msg=$(echo "$line" | sed -E 's/^[a-z]+(\([^)]*\))?:[ ]*//')
        [ -n "$clean_msg" ] && echo "- $clean_msg"
    done
    echo ""
fi

echo "---"
echo ""
if [ -n "$PREV_TAG" ]; then
    echo "**Full Changelog**: https://github.com/gelleson/rpc-cli/compare/$PREV_TAG...HEAD"
else
    echo "**Full Changelog**: https://github.com/gelleson/rpc-cli"
fi
