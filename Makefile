BINARY_NAME=ai-memoria-cli

.PHONY: build clean run test deps

build:
	go build -o bin/$(BINARY_NAME) main.go

build-all:
	GOOS=linux GOARCH=amd64 go build -o bin/$(BINARY_NAME)-linux-amd64
	GOOS=darwin GOARCH=amd64 go build -o bin/$(BINARY_NAME)-darwin-amd64
	GOOS=darwin GOARCH=arm64 go build -o bin/$(BINARY_NAME)-darwin-arm64
	GOOS=windows GOARCH=amd64 go build -o bin/$(BINARY_NAME)-windows-amd64.exe

clean:
	rm -rf bin/

run:
	go run main.go $(ARGS)

test:
	go test -v ./...

deps:
	go mod tidy
	go mod download

install:
	go install
