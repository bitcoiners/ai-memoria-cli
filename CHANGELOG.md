# Changelog

## [0.1.0] - 2024-03-26

### Added
- Initial release
- Token-based authentication (`auth login`, `auth logout`, `auth whoami`)
- User management (`users create`)
- Status check (`status`)
- JSON output support (`--json`)
- Profile support for development/production
- Configuration stored in `~/.ai-memoria/config.json`
- Environment variable support (`AI_MEMORIA_API_KEY`, `AI_MEMORIA_API_URL`, `AI_MEMORIA_PROFILE`)
- Cross-platform builds (Linux, macOS, Windows)
- Installation scripts (`install.sh`, `get.sh`)
- Makefile with build and test targets
- Unit and integration tests
- Mock server for testing

### Installation
```bash
curl -sSL https://raw.githubusercontent.com/bitcoiners/ai-memoria-cli/main/get.sh | bash
```
### Usage
```bash
mem auth login --email user@example.com --password secret
mem auth whoami
mem users create --email new@example.com --username newuser --name "New User" --password secret
mem status
```
