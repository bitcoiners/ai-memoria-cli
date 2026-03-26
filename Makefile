BINARY_NAME=mem
PROJECT_NAME=ai-memoria-cli
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-X main.Version=${VERSION}"

.PHONY: build clean run test deps test-unit test-integration install uninstall coverage release

build:
	go build $(LDFLAGS) -o bin/$(BINARY_NAME) main.go

build-all:
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-linux-amd64 main.go
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-linux-arm64 main.go
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-darwin-amd64 main.go
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-darwin-arm64 main.go
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-windows-amd64.exe main.go
	GOOS=windows GOARCH=arm64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-windows-arm64.exe main.go

install: build
	@mkdir -p ~/.local/bin
	@cp bin/$(BINARY_NAME) ~/.local/bin/
	@echo "✅ Installed to ~/.local/bin/$(BINARY_NAME) (version $(VERSION))"
	@echo ""
	@echo "Make sure ~/.local/bin is in your PATH. Add this to your ~/.bashrc or ~/.zshrc:"
	@echo "  export PATH=\$$PATH:~/.local/bin"
	@echo ""
	@echo "Try it out: $(BINARY_NAME) --version"

uninstall:
	@rm -f ~/.local/bin/$(BINARY_NAME)
	@rm -rf ~/.ai-memoria
	@echo "✅ Uninstalled $(BINARY_NAME)"

clean:
	rm -rf bin/
	rm -rf releases/
	rm -f coverage.out coverage.html

run:
	go run main.go $(ARGS)

test: test-unit test-integration

test-unit:
	go test -v ./tests/unit/...

test-integration:
	@echo "🔧 Checking Rails API server..."
	@curl -s http://localhost:3000/up > /dev/null && echo "✅ Rails API is running" || (echo "❌ Rails API not running. Start with: cd ../api && rails server" && exit 1)
	@echo ""
	@echo "🧪 Running integration tests..."
	@go build -o bin/$(BINARY_NAME) main.go
	@go test -v ./tests/integration/...

deps:
	go mod tidy
	go mod download

coverage:
	@echo "📊 Generating coverage report..."
	@go test -coverprofile=coverage.out -coverpkg=./internal/... ./tests/unit/...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report generated: coverage.html"
	@xdg-open coverage.html 2>/dev/null || open coverage.html 2>/dev/null || echo "Open coverage.html manually"

release:
	./release.sh $(VERSION)

dev: build
	@echo "Running development build (version $(VERSION))..."
	@./bin/$(BINARY_NAME) --version
