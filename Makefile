# Makefile for kvstore project

.PHONY: test build clean

# Build the project (optional, if you add a real app later)
build:
	protoc --go_out=paths=source_relative:proto --go-grpc_out=paths=source_relative:proto --proto_path=proto proto/kvstore.proto
	go mod tidy
	go build ./...
	
# Default test target
test:
	go test -v -cover ./...

# Clean up build artifacts (optional)
clean:
	go clean
