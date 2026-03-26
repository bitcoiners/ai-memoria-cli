#!/bin/bash
set -e

# Check if version is provided
if [ -z "$1" ]; then
    echo "❌ Error: Version number required"
    echo ""
    echo "Usage: ./release.sh <version>"
    echo "Example: ./release.sh v0.2.0"
    exit 1
fi

VERSION=$1
BINARY_NAME="mem"
RELEASE_DIR="releases/$VERSION"
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${GREEN}🚀 Creating release $VERSION...${NC}"
echo ""

# Validate version format (should start with v)
if [[ ! $VERSION =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo -e "${RED}❌ Error: Version must be in format vX.Y.Z (e.g., v0.2.0)${NC}"
    exit 1
fi

# Check if we're on main branch
BRANCH=$(git branch --show-current)
if [ "$BRANCH" != "main" ] && [ "$BRANCH" != "master" ]; then
    echo -e "${YELLOW}⚠️  You are on branch '$BRANCH', not 'main' or 'master'${NC}"
    read -p "Continue anyway? (y/n) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

# Check for uncommitted changes
if ! git diff-index --quiet HEAD --; then
    echo -e "${YELLOW}⚠️  You have uncommitted changes.${NC}"
    read -p "Commit them now? (y/n) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        git add .
        git commit -m "Pre-release cleanup for $VERSION"
    else
        echo "Please commit or stash changes first."
        exit 1
    fi
fi

# Build binaries
echo "📦 Building binaries..."
mkdir -p "$RELEASE_DIR"

echo "  Building for Linux (amd64)..."
GOOS=linux GOARCH=amd64 go build -o "$RELEASE_DIR/${BINARY_NAME}-linux-amd64" main.go

echo "  Building for Linux (arm64)..."
GOOS=linux GOARCH=arm64 go build -o "$RELEASE_DIR/${BINARY_NAME}-linux-arm64" main.go

echo "  Building for macOS (amd64)..."
GOOS=darwin GOARCH=amd64 go build -o "$RELEASE_DIR/${BINARY_NAME}-darwin-amd64" main.go

echo "  Building for macOS (arm64)..."
GOOS=darwin GOARCH=arm64 go build -o "$RELEASE_DIR/${BINARY_NAME}-darwin-arm64" main.go

echo "  Building for Windows (amd64)..."
GOOS=windows GOARCH=amd64 go build -o "$RELEASE_DIR/${BINARY_NAME}-windows-amd64.exe" main.go

echo "  Building for Windows (arm64)..."
GOOS=windows GOARCH=arm64 go build -o "$RELEASE_DIR/${BINARY_NAME}-windows-arm64.exe" main.go

# Create checksums
echo "🔐 Creating checksums..."
cd "$RELEASE_DIR"
sha256sum * > checksums.txt
cd ../..

echo -e "\n✅ Binaries built in $RELEASE_DIR/"
ls -lh "$RELEASE_DIR/"
echo ""

# Handle git tag
if git rev-parse "$VERSION" >/dev/null 2>&1; then
    echo -e "${YELLOW}⚠️  Tag $VERSION already exists${NC}"
    read -p "Delete existing tag and create new one? (y/n) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        echo "Deleting existing tag..."
        git tag -d "$VERSION"
        git push origin ":refs/tags/$VERSION" 2>/dev/null || true
        echo "Creating new tag..."
        git tag -a "$VERSION" -m "Release $VERSION"
        git push origin "$VERSION"
    else
        echo "Keeping existing tag."
    fi
else
    echo "🏷️  Creating and pushing git tag $VERSION..."
    git tag -a "$VERSION" -m "Release $VERSION"
    git push origin "$VERSION"
fi

# Push any committed changes
git push origin "$BRANCH" 2>/dev/null || true

# Check if gh CLI is installed
if command -v gh &> /dev/null; then
    echo "📝 Creating GitHub release..."
    
    # Check if release already exists
    if gh release view "$VERSION" &>/dev/null; then
        echo -e "${YELLOW}⚠️  Release $VERSION already exists on GitHub${NC}"
        read -p "Delete existing release and create new one? (y/n) " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            echo "Deleting existing release..."
            gh release delete "$VERSION" -y
        else
            echo "Skipping release creation."
            exit 0
        fi
    fi
    
    # Generate release notes
    RELEASE_NOTES=$(cat <<-END
## What's New in $VERSION

### Added
- Comprehensive unit tests with 81%+ coverage
- Integration tests against real Rails API
- Uninstall command (\`mem uninstall\`)
- Mock server for testing
- Test documentation

### Improved
- Better error handling in API client
- Binary path detection in tests
- Makefile with coverage target
- Test isolation with temporary config directories

### Fixed
- Integration test binary path issues
- Config file handling in tests
- Revoke token error handling

## Installation

\`\`\`bash
curl -sSL https://raw.githubusercontent.com/bitcoiners/ai-memoria-cli/main/get.sh | bash
\`\`\`

## Checksums

\`\`\`
$(cat releases/$VERSION/checksums.txt)
\`\`\`
END
)
    
    # Create release
    gh release create "$VERSION" \
        --title "AI Memoria CLI $VERSION" \
        --notes "$RELEASE_NOTES" \
        releases/"$VERSION"/*
    
    echo -e "\n${GREEN}✅ Release created: https://github.com/bitcoiners/ai-memoria-cli/releases/tag/$VERSION${NC}"
else
    echo -e "${YELLOW}⚠️  GitHub CLI (gh) not installed${NC}"
    echo ""
    echo "Manual steps:"
    echo "1. Push changes: git push origin $BRANCH"
    echo "2. Push tags: git push --tags"
    echo "3. Go to: https://github.com/bitcoiners/ai-memoria-cli/releases"
    echo "4. Click 'Create a new release'"
    echo "5. Use tag: $VERSION"
    echo "6. Upload files from releases/$VERSION/"
    echo "7. Publish release"
fi

echo ""
echo -e "${GREEN}🎉 Release $VERSION complete!${NC}"
