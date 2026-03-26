BINARY_NAME=mem
PROJECT_NAME=ai-memoria-cli

.PHONY: build clean run test deps test-unit test-integration install uninstall

build:
	go build -o bin/$(BINARY_NAME) main.go

build-all:
	GOOS=linux GOARCH=amd64 go build -o bin/$(BINARY_NAME)-linux-amd64
	GOOS=linux GOARCH=arm64 go build -o bin/$(BINARY_NAME)-linux-arm64
	GOOS=darwin GOARCH=amd64 go build -o bin/$(BINARY_NAME)-darwin-amd64
	GOOS=darwin GOARCH=arm64 go build -o bin/$(BINARY_NAME)-darwin-arm64
	GOOS=windows GOARCH=amd64 go build -o bin/$(BINARY_NAME)-windows-amd64.exe
	GOOS=windows GOARCH=arm64 go build -o bin/$(BINARY_NAME)-windows-arm64.exe

install: build
	@mkdir -p ~/.local/bin
	@cp bin/$(BINARY_NAME) ~/.local/bin/
	@echo "✅ Installed to ~/.local/bin/$(BINARY_NAME)"
	@echo ""
	@echo "Make sure ~/.local/bin is in your PATH. Add this to your ~/.bashrc or ~/.zshrc:"
	@echo "  export PATH=\$$PATH:~/.local/bin"
	@echo ""
	@echo "Try it out: $(BINARY_NAME) --help"

uninstall:
	@rm -f ~/.local/bin/$(BINARY_NAME)
	@rm -rf ~/.ai-memoria
	@echo "✅ Uninstalled $(BINARY_NAME)"

clean:
	rm -rf bin/
	rm -f coverage.out

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
	go test -coverprofile=coverage.out ./tests/unit/...
	go tool cover -html=coverage.out

dev: build
	@echo "Running development build..."
	@./bin/$(BINARY_NAME) --help

release: build-release
build-release:
	./build-release.sh $(VERSION)
