# AI Memoria CLI

Command line interface for AI Memoria API. The CLI is designed to be simple and memorable with the command `mem`.

## Quick Install

### Option 1: One-liner (Recommended)

```bash

curl -sSL https://raw.githubusercontent.com/bitcoiners/ai-memoria-cli/main/get.sh | bash

```

### Option 2: Download from GitHub Releases
1. Download the binary for your platform from [releases](https://github.com/bitcoiners/ai-memoria-cli/releases)
2. Install manually:

   **Linux/macOS:**
   
```bash
   # Download (example for Linux amd64)
   wget https://github.com/bitcoiners/ai-memoria-cli/releases/download/v0.2.0/mem-linux-amd64 -O mem
   
   # Make executable
   chmod +x mem
   
   # Move to a directory in your PATH
   mkdir -p ~/.local/bin
   mv mem ~/.local/bin/
   
   # Verify installation
   mem --version
   
```

   **Windows:**
```powershell

   # Download mem-windows-amd64.exe
   # Rename to mem.exe
   # Move to a directory in your PATH, e.g., C:\Users\YourName\bin\
   # Add that directory to your PATH environment variable
   
```

### Option 3: Build from Source
```bash

git clone git@github.com:bitcoiners/ai-memoria-cli.git
cd ai-memoria-cli
make install

```

### Option 4: Using Go Install
```bash

go install github.com/bitcoiners/ai-memoria-cli@latest

```

## Building Binaries

### Build for Current Platform
```bash
make build
```
This creates a binary at `bin/mem`

### Build for All Platforms (for releases)
```bash
make build-all
```
This creates binaries for:
- Linux (amd64, arm64)
- macOS (amd64, arm64)
- Windows (amd64, arm64)

Binaries will be in the `bin/` directory:
- `mem-linux-amd64`
- `mem-linux-arm64`
- `mem-darwin-amd64`
- `mem-darwin-arm64`
- `mem-windows-amd64.exe`
- `mem-windows-arm64.exe`

### Build Release Package
```bash
# Build with version tag
./release.sh v1.0.0

# Or use make
make release VERSION=v1.0.0
```
This creates a `releases/` directory with all platform binaries and checksums.

### Manual Build
```bash
# Build for your current platform
go build -o mem main.go

# Build for specific platform
GOOS=linux GOARCH=amd64 go build -o mem-linux-amd64 main.go
GOOS=darwin GOARCH=arm64 go build -o mem-darwin-arm64 main.go
GOOS=windows GOARCH=amd64 go build -o mem-windows-amd64.exe main.go
```

## Creating GitHub Releases

### Prerequisites

1. Install [GitHub CLI](https://cli.github.com/):
   ```bash
   # Ubuntu/Debian
   curl -fsSL https://cli.github.com/packages/githubcli-archive-keyring.gpg | sudo dd of=/usr/share/keyrings/githubcli-archive-keyring.gpg
   echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main" | sudo tee /etc/apt/sources.list.d/github-cli.list > /dev/null
   sudo apt update
   sudo apt install gh
   
   # macOS
   brew install gh
   
   # Windows (using winget)
   winget install --id GitHub.cli
   ```

2. Authenticate with GitHub:
   ```bash
   gh auth login
   ```
   Follow the prompts to authenticate with your GitHub account.

### Creating a Release

#### Step 1: Ensure all changes are committed
```bash
git status
git add .
git commit -m "Your commit message"
git push origin main
```

#### Step 2: Create and push a version tag
```bash
# Create a tag (using semantic versioning)
git tag v0.1.0

# Push the tag to GitHub
git push origin v0.1.0
```

#### Step 3: Build binaries for the release
```bash
# Build all platform binaries
./release.sh v0.1.0

# Or use make
make release VERSION=v0.1.0
```

#### Step 4: Create GitHub Release

**Using GitHub CLI (automated):**
```bash
# Create release with all binaries
gh release create v0.1.0 \
  --title "AI Memoria CLI v0.1.0" \
  --notes "Initial release

Features:
- Token-based authentication
- User management
- Status checks
- JSON output support
- Profile support

Installation:
```bash
curl -sSL https://raw.githubusercontent.com/bitcoiners/ai-memoria-cli/main/get.sh | bash
```" \
  releases/v0.1.0/*
```

**Using GitHub Web Interface (manual):**
1. Go to: https://github.com/bitcoiners/ai-memoria-cli/releases
2. Click "Create a new release"
3. Select the tag you created (e.g., `v0.1.0`)
4. Enter a title: "AI Memoria CLI v0.1.0"
5. Add release notes
6. Drag and drop all files from `releases/v0.1.0/` into the attachments area
7. Click "Publish release"

### Automated Release Script

Use the included `release.sh` script to automate the entire process:
```bash
# Make the script executable
chmod +x release.sh

# Create release (automates steps 1-4)
./release.sh v0.1.0

# Or use make
make release VERSION=v0.1.0
```

### Release Checklist

- [ ] All tests passing (`make test`)
- [ ] Code is committed and pushed
- [ ] Version tag is created and pushed
- [ ] Binaries are built for all platforms
- [ ] Checksums are generated
- [ ] Release notes are written
- [ ] Release is published on GitHub
- [ ] Installation instructions are tested

### Versioning Guidelines

Follow [Semantic Versioning](https://semver.org/):
- **MAJOR** version (v1.0.0): Incompatible API changes
- **MINOR** version (v0.1.0): Backward-compatible functionality
- **PATCH** version (v0.1.1): Backward-compatible bug fixes

## Post-Installation

Make sure `~/.local/bin` is in your PATH. Add this to your `~/.bashrc` or `~/.zshrc`:
```bash
export PATH="$PATH:$HOME/.local/bin"
```

Then reload your shell:
```bash
source ~/.bashrc  # or source ~/.zshrc
```

## Configuration

The CLI stores configuration in `~/.ai-memoria/config.json`

### Environment Variables
- `AI_MEMORIA_API_KEY` - API key for authentication
- `AI_MEMORIA_API_URL` - API base URL (default: http://localhost:3000)
- `AI_MEMORIA_PROFILE` - Configuration profile (development/production)

## Usage

### Authentication
```bash
# Login
mem auth login --email user@example.com --password secret

# Show current user
mem auth whoami

# Logout
mem auth logout
```

### User Management
```bash
# Create a new user (public signup)
mem users create --email new@example.com --username newuser --name "New User" --password secret
```

### Status
```bash
# Check API health
mem status
```

### JSON Output
Use `--json` flag for machine-readable output:
```bash
mem --json auth whoami
```

### Profiles
Switch between development and production:
```bash
# Use production profile
mem --profile production auth whoami

# Or set environment variable
export AI_MEMORIA_PROFILE=production
mem auth whoami
```

## Development

```bash
# Build
make build

# Run tests
make test

# Run unit tests only
make test-unit

# Run integration tests only
make test-integration

# Install locally
make install

# Uninstall
make uninstall

# Build for all platforms
make build-all

# Create release package
make release VERSION=v1.0.0

# Clean build artifacts
make clean

# Run with coverage
make coverage
```

## Project Structure
```
.
├── bin/                 # Compiled binaries
├── cmd/                 # Command implementations
│   ├── auth/           # Authentication commands
│   ├── users/          # User management commands
│   └── status/         # Status command
├── internal/           # Internal packages
│   ├── api/            # API client
│   ├── config/         # Configuration management
│   ├── models/         # Data models
│   └── utils/          # Utility functions
├── tests/              # Test files
│   ├── integration/    # Integration tests
│   └── mock/           # Mock server for testing
├── releases/           # Release packages
├── Makefile            # Build automation
├── go.mod              # Go module definition
├── main.go             # Entry point
├── install.sh          # Binary installer (for source builds)
├── get.sh              # One-liner installer
├── release.sh          # GitHub release automation
├── CHANGELOG.md        # Version history
└── README.md           # This file
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests: `make test`
5. Commit and push
6. Create a pull request

## Uninstall

To completely remove AI Memoria CLI:

```bash
# If you have the CLI installed
mem uninstall

# Or manually remove
rm -f ~/.local/bin/mem
rm -rf ~/.ai-memoria
```

## License

MIT

## Testing

### Unit Tests
```bash
make test-unit
```

### Integration Tests
Integration tests require a running Rails API server at `http://localhost:3000`.

First, start the Rails API:
```bash
cd ../api
rails server
```

Then run the integration tests:
```bash
make test-integration
```

### All Tests
```bash
make test
```

### Test Coverage
```bash
make coverage
```
