#!/bin/bash
set -e

VERSION=${1:-"latest"}
BINARY_NAME="mem"
RELEASE_DIR="releases/$VERSION"

echo "📦 Building release $VERSION..."

# Create release directory
mkdir -p "$RELEASE_DIR"

# Build for all platforms
echo "Building for Linux (amd64)..."
GOOS=linux GOARCH=amd64 go build -o "$RELEASE_DIR/${BINARY_NAME}-linux-amd64" main.go

echo "Building for Linux (arm64)..."
GOOS=linux GOARCH=arm64 go build -o "$RELEASE_DIR/${BINARY_NAME}-linux-arm64" main.go

echo "Building for macOS (amd64)..."
GOOS=darwin GOARCH=amd64 go build -o "$RELEASE_DIR/${BINARY_NAME}-darwin-amd64" main.go

echo "Building for macOS (arm64)..."
GOOS=darwin GOARCH=arm64 go build -o "$RELEASE_DIR/${BINARY_NAME}-darwin-arm64" main.go

echo "Building for Windows (amd64)..."
GOOS=windows GOARCH=amd64 go build -o "$RELEASE_DIR/${BINARY_NAME}-windows-amd64.exe" main.go

echo "Building for Windows (arm64)..."
GOOS=windows GOARCH=arm64 go build -o "$RELEASE_DIR/${BINARY_NAME}-windows-arm64.exe" main.go

# Create checksums
echo "Creating checksums..."
cd "$RELEASE_DIR"
sha256sum * > checksums.txt
cd ../..

echo -e "\n✅ Release built in $RELEASE_DIR/"
echo ""
echo "Files:"
ls -lh "$RELEASE_DIR/"
echo ""
echo "To create GitHub release, run:"
echo "  gh release create $VERSION $RELEASE_DIR/*"
