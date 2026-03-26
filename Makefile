BINARY_NAME=mem
PROJECT_NAME=ai-memoria-cli

.PHONY: build clean run test deps test-unit test-integration install uninstall

build:
	go build -o bin/$(BINARY_NAME) main.go

build-all:
	GOOS=linux GOARCH=amd64 go build -o bin/$(BINARY_NAME)-linux-amd64
	GOOS=darwin GOARCH=amd64 go build -o bin/$(BINARY_NAME)-darwin-amd64
	GOOS=darwin GOARCH=arm64 go build -o bin/$(BINARY_NAME)-darwin-arm64
	GOOS=windows GOARCH=amd64 go build -o bin/$(BINARY_NAME)-windows-amd64.exe

# Install to user directory (default, no sudo required)
install: build
	@mkdir -p ~/.local/bin
	@cp bin/$(BINARY_NAME) ~/.local/bin/
	@echo "✅ Installed to ~/.local/bin/$(BINARY_NAME)"
	@echo ""
	@echo "Make sure ~/.local/bin is in your PATH. Add this to your ~/.bashrc or ~/.zshrc:"
	@echo "  export PATH=\$$PATH:~/.local/bin"
	@echo ""
	@echo "Then reload your shell: source ~/.bashrc (or restart terminal)"
	@echo ""
	@echo "Try it out: $(BINARY_NAME) --help"

# Uninstall
uninstall:
	@rm -f ~/.local/bin/$(BINARY_NAME)
	@echo "✅ Uninstalled $(BINARY_NAME)"

clean:
	rm -rf bin/
	rm -rf tests/tmp/

run:
	go run main.go $(ARGS)

test: test-unit test-integration

test-unit:
	go test -v ./internal/...

test-integration:
	go build -o bin/$(BINARY_NAME) main.go
	go test -v ./tests/integration/...

deps:
	go mod tidy
	go mod download

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# Development helpers
dev: build
	@echo "Running development build..."
	@./bin/$(BINARY_NAME) --help

# Create GitHub release
release:
	./release.sh $(VERSION)
