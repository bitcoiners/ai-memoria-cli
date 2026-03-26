# AI Memoria CLI

Command line interface for AI Memoria API.

## Installation

### From Source

```bash
git clone git@github.com:bitcoiners/ai-memoria-cli.git
cd ai-memoria-cli
make build
```

### Download Binary
Download the latest binary from https://github.com/bitcoiners/ai-memoria-cli/releases

### Configuration
The CLI stores configuration in ~/.ai-memoria/config.json

#### Create Config File

```bash
mkdir -p ~/.ai-memoria

cat > ~/.ai-memoria/config.json << 'EOF'
{
  "development": {
    "base_url": "http://localhost:3000",
    "profile": "development"
  },
  "production": {
    "base_url": "https://api.ai-memoria.com",
    "profile": "production"
  },
  "default_profile": "development"
}
EOF
```

### Environment Variables
- AI_MEMORIA_API_KEY - API key for authentication
- AI_MEMORIA_API_URL - API base URL (default: http://localhost:3000)
- AI_MEMORIA_PROFILE - Configuration profile (development/production)

### Usage
#### Authentication
```bash
# Login
ai-memoria-cli auth login --email user@example.com --password secret

# Show current user
ai-memoria-cli auth whoami

# Logout
ai-memoria-cli auth logout
```

### User Management
```bash
# Create a new user (public signup)
ai-memoria-cli users create --email new@example.com --username newuser --name "New User" --password secret
```

###Status
```bash
# Check API health
ai-memoria-cli status
```

### JSON Output
Use --json flag for machine-readable output:

```bash
ai-memoria-cli --json auth whoami
```
### Profiles
Switch between development and production:

```bash
# Use production profile
ai-memoria-cli --profile production auth whoami

# Or set environment variable
export AI_MEMORIA_PROFILE=production
ai-memoria-cli auth whoami
```
Development
```bash
# Build
make build

# Run tests
make test

# Build for all platforms
make build-all

# Clean
make clean
```

### License
MIT

