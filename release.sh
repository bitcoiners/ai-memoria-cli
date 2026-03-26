#!/bin/bash
set -e

VERSION=${1:-"v0.1.0"}

echo "📦 Creating release $VERSION..."

# Check if we're on main branch
BRANCH=$(git branch --show-current)
if [ "$BRANCH" != "main" ]; then
    echo "⚠️  You are on branch '$BRANCH', not 'main'"
    read -p "Continue anyway? (y/n) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

# Check for uncommitted changes
if ! git diff-index --quiet HEAD --; then
    echo "❌ You have uncommitted changes. Commit or stash them first."
    exit 1
fi

# Build binaries
echo "Building binaries..."
./build-release.sh "$VERSION"

# Create tag if it doesn't exist
if ! git rev-parse "$VERSION" >/dev/null 2>&1; then
    echo "Creating tag $VERSION..."
    git tag "$VERSION"
    git push origin "$VERSION"
else
    echo "Tag $VERSION already exists"
fi

# Check if gh CLI is installed
if command -v gh &> /dev/null; then
    echo "Creating release with GitHub CLI..."
    
    # Check if release already exists
    if gh release view "$VERSION" &>/dev/null; then
        echo "Release $VERSION already exists. Deleting..."
        gh release delete "$VERSION" -y
    fi
    
    # Create release
    gh release create "$VERSION" \
        --title "AI Memoria CLI $VERSION" \
        --notes "AI Memoria CLI $VERSION

## Installation

### One-liner
\`\`\`bash
curl -sSL https://raw.githubusercontent.com/bitcoiners/ai-memoria-cli/main/get.sh | bash
\`\`\`

### Download
Download the binary for your platform below.

## Features
- Token-based authentication
- User management
- Status checks
- JSON output
- Profile support (development/production)

## Checksums
\`\`\`
$(cat releases/$VERSION/checksums.txt)
\`\`\`" \
        releases/"$VERSION"/*
    
    echo "✅ Release created: https://github.com/bitcoiners/ai-memoria-cli/releases/tag/$VERSION"
else
    echo "⚠️  GitHub CLI (gh) not installed"
    echo "Manual steps:"
    echo "1. Go to: https://github.com/bitcoiners/ai-memoria-cli/releases"
    echo "2. Click 'Create a new release'"
    echo "3. Use tag: $VERSION"
    echo "4. Upload files from releases/$VERSION/"
    echo "5. Publish release"
fi
