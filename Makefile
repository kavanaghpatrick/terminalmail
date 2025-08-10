# Makefile for TerminalMail

# Build configuration
BINARY_NAME=tmail
MAIN_PACKAGE=./cmd/tmail

# Default target
all: build

# Build the application
build:
	go build -tags sqlite_fts5 -o $(BINARY_NAME) $(MAIN_PACKAGE)

# Build for different platforms
build-mac:
	GOOS=darwin GOARCH=arm64 go build -tags sqlite_fts5 -o $(BINARY_NAME)-mac $(MAIN_PACKAGE)

build-linux:
	GOOS=linux GOARCH=amd64 go build -tags sqlite_fts5 -o $(BINARY_NAME)-linux $(MAIN_PACKAGE)

build-windows:
	GOOS=windows GOARCH=amd64 go build -tags sqlite_fts5 -o $(BINARY_NAME)-windows.exe $(MAIN_PACKAGE)

# Clean build artifacts
clean:
	rm -f $(BINARY_NAME) $(BINARY_NAME)-*

# Install dependencies
deps:
	go mod tidy

# Run tests
test:
	go test ./...

.PHONY: all build clean deps test