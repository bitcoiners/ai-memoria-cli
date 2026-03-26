#!/bin/bash
set -e

BINARY_NAME="mem"
INSTALL_DIR="$HOME/.local/bin"
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "🚀 Installing AI Memoria CLI..."

# Create install directory if it doesn't exist
mkdir -p "$INSTALL_DIR"

# Check if we're in a source directory with a binary in bin/
if [ -f "bin/$BINARY_NAME" ]; then
    echo "Found built binary in bin/$BINARY_NAME"
    cp "bin/$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"
    chmod +x "$INSTALL_DIR/$BINARY_NAME"
    
# Check if we're in a downloaded release with the binary in current directory
elif [ -f "./$BINARY_NAME" ]; then
    echo "Found binary in current directory"
    cp "./$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"
    chmod +x "$INSTALL_DIR/$BINARY_NAME"
    
else
    echo -e "${RED}Error: Binary not found!${NC}"
    echo ""
    echo "Usage:"
    echo "  From source:    make install"
    echo "  From release:   ./install.sh (run in directory with the binary)"
    echo ""
    echo "Or download the binary manually and copy it to ~/.local/bin/mem"
    exit 1
fi

echo -e "${GREEN}✅ Installed to $INSTALL_DIR/$BINARY_NAME${NC}"
echo ""

# Check if INSTALL_DIR is in PATH
if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
    echo -e "${YELLOW}⚠️  $INSTALL_DIR is not in your PATH${NC}"
    echo ""
    echo "Add this to your ~/.bashrc or ~/.zshrc:"
    echo "  export PATH=\"\$PATH:$INSTALL_DIR\""
    echo ""
    echo "Then reload your shell: source ~/.bashrc (or restart terminal)"
    echo ""
fi

echo -e "${GREEN}Try it out: $BINARY_NAME --help${NC}"
