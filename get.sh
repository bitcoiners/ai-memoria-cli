#!/bin/bash
set -e

VERSION=${1:-"latest"}
BINARY_NAME="mem"
INSTALL_DIR="$HOME/.local/bin"
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case $ARCH in
    x86_64)
        ARCH="amd64"
        ;;
    aarch64|arm64)
        ARCH="arm64"
        ;;
    *)
        echo -e "${RED}Unsupported architecture: $ARCH${NC}"
        exit 1
        ;;
esac

case $OS in
    linux|darwin)
        # OK
        ;;
    *)
        echo -e "${RED}Unsupported OS: $OS${NC}"
        exit 1
        ;;
esac

if [ "$VERSION" = "latest" ]; then
    # Get latest version from GitHub
    VERSION=$(curl -s https://api.github.com/repos/bitcoiners/ai-memoria-cli/releases/latest | grep -o '"tag_name": "[^"]*"' | cut -d'"' -f4)
    if [ -z "$VERSION" ]; then
        echo -e "${RED}Failed to get latest version${NC}"
        exit 1
    fi
fi

echo "📦 Downloading AI Memoria CLI $VERSION for $OS-$ARCH..."

# Download URL
URL="https://github.com/bitcoiners/ai-memoria-cli/releases/download/$VERSION/${BINARY_NAME}-${OS}-${ARCH}"
if [ "$OS" = "windows" ]; then
    URL="${URL}.exe"
fi

# Download
curl -L -o "$BINARY_NAME" "$URL"
chmod +x "$BINARY_NAME"

# Install
mkdir -p "$INSTALL_DIR"
mv "$BINARY_NAME" "$INSTALL_DIR/"

echo -e "${GREEN}✅ Installed to $INSTALL_DIR/$BINARY_NAME${NC}"
echo ""
echo "Try it out: $BINARY_NAME --help"
