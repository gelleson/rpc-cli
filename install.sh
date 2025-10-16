#!/usr/bin/env bash

set -e

# Configuration
BINARY_NAME="rpc-cli"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"
REPO="gelleson/rpc-cli"

# Detect OS and architecture
OS="$(uname -s)"
ARCH="$(uname -m)"

case "$OS" in
    Linux*)     OS="linux";;
    Darwin*)    OS="darwin";;
    CYGWIN*|MINGW*|MSYS*) OS="windows";;
    *)          echo "Unsupported OS: $OS"; exit 1;;
esac

case "$ARCH" in
    x86_64)     ARCH="amd64";;
    aarch64|arm64) ARCH="arm64";;
    *)          echo "Unsupported architecture: $ARCH"; exit 1;;
esac

echo "Installing $BINARY_NAME for $OS/$ARCH..."

# Get latest release version
VERSION=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$VERSION" ]; then
    echo "Failed to get latest version"
    exit 1
fi

echo "Latest version: $VERSION"

# Construct download URL
ARCHIVE_NAME="${BINARY_NAME}_${VERSION#v}_${OS}_${ARCH}.tar.gz"
DOWNLOAD_URL="https://github.com/$REPO/releases/download/$VERSION/$ARCHIVE_NAME"

# Create temporary directory
TMP_DIR=$(mktemp -d)
trap "rm -rf $TMP_DIR" EXIT

# Download archive
echo "Downloading from $DOWNLOAD_URL..."
curl -sL "$DOWNLOAD_URL" -o "$TMP_DIR/$ARCHIVE_NAME"

# Extract binary
echo "Extracting binary..."
tar xz -C "$TMP_DIR" -f "$TMP_DIR/$ARCHIVE_NAME" "$BINARY_NAME"

# Make executable
chmod +x "$TMP_DIR/$BINARY_NAME"

# Install binary
echo "Installing to $INSTALL_DIR/$BINARY_NAME..."
if [ -w "$INSTALL_DIR" ]; then
    mv "$TMP_DIR/$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"
else
    sudo mv "$TMP_DIR/$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"
fi

echo "Successfully installed $BINARY_NAME $VERSION"
echo "Run '$BINARY_NAME --help' to get started"
