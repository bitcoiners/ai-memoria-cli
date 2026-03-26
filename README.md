# AI Memoria CLI

Command line interface for AI Memoria API. The CLI is designed to be simple and memorable with the command `mem`.

## Quick Install

### Option 1: One-liner (Recommended)
\`\`\`bash
curl -sSL https://raw.githubusercontent.com/bitcoiners/ai-memoria-cli/main/get.sh | bash
\`\`\`

### Option 2: Download from GitHub Releases
1. Download the binary for your platform from [releases](https://github.com/bitcoiners/ai-memoria-cli/releases)
2. Run the installer:
   \`\`\`bash
   chmod +x install.sh
   ./install.sh
   \`\`\`

### Option 3: Build from Source
\`\`\`bash
git clone git@github.com:bitcoiners/ai-memoria-cli.git
cd ai-memoria-cli
make install
\`\`\`

### Option 4: Using Go Install
\`\`\`bash
go install github.com/bitcoiners/ai-memoria-cli@latest
\`\`\`

## Building Binaries

### Build for Current Platform
\`\`\`bash
make build
\`\`\`
This creates a binary at `bin/mem`

### Build for All Platforms (for releases)
\`\`\`bash
make build-all
\`\`\`
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
\`\`\`bash
# Build with version tag
./build-release.sh v1.0.0

# Or use make
make release VERSION=v1.0.0
\`\`\`
This creates a `releases/` directory with all platform binaries and checksums.

### Manual Build
\`\`\`bash
# Build for your current platform
go build -o mem main.go

# Build for specific platform
GOOS=linux GOARCH=amd64 go build -o mem-linux-amd64 main.go
GOOS=darwin GOARCH=arm64 go build -o mem-darwin-arm64 main.go
GOOS=windows GOARCH=amd64 go build -o mem-windows-amd64.exe main.go
\`\`\`

## Post-Installation

Make sure `~/.local/bin` is in your PATH. Add this to your `~/.bashrc` or `~/.zshrc`:
\`\`\`bash
export PATH="$PATH:$HOME/.local/bin"
\`\`\`

Then reload your shell:
\`\`\`bash
source ~/.bashrc  # or source ~/.zshrc
\`\`\`

## Configuration

The CLI stores configuration in `~/.ai-memoria/config.json`

### Environment Variables
- `AI_MEMORIA_API_KEY` - API key for authentication
- `AI_MEMORIA_API_URL` - API base URL (default: http://localhost:3000)
- `AI_MEMORIA_PROFILE` - Configuration profile (development/production)

## Usage

### Authentication
\`\`\`bash
# Login
mem auth login --email user@example.com --password secret

# Show current user
mem auth whoami

# Logout
mem auth logout
\`\`\`

### User Management
\`\`\`bash
# Create a new user (public signup)
mem users create --email new@example.com --username newuser --name "New User" --password secret
\`\`\`

### Status
\`\`\`bash
# Check API health
mem status
\`\`\`

### JSON Output
Use `--json` flag for machine-readable output:
\`\`\`bash
mem --json auth whoami
\`\`\`

### Profiles
Switch between development and production:
\`\`\`bash
# Use production profile
mem --profile production auth whoami

# Or set environment variable
export AI_MEMORIA_PROFILE=production
mem auth whoami
\`\`\`

## Development

\`\`\`bash
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
\`\`\`

## Project Structure
\`\`\`
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
├── Makefile            # Build automation
├── go.mod              # Go module definition
└── main.go             # Entry point
\`\`\`

## License

MIT
