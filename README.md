# KVStore

A lightweight, extensible, and customizable Key-Value Store library in Go, served over gRPC.

This project is designed to be minimal but powerful:
- In-memory key-value storage
- TTL (expiration) support
- gRPC server exposing Set, Get, and Delete operations
- Hook system for custom authentication, logging, rate-limiting, and more
- Functional options to customize server behavior
- Storage backend pluggability

Built for developers who want control without unnecessary complexity.

---

## Features

- **In-Memory Key-Value Store** with concurrency safety.
- **TTL Expiration** (keys can expire automatically).
- **gRPC Interface** (Set, Get, Delete operations).
- **Pre and Post Hooks** (inject custom logic before/after every operation).
- **Customizable Storage Backend** (swap in Redis, database, etc.).
- **Functional Options** for server customization.
- **Extensive Unit and Integration Tests**.
- **Simple Makefile** for easy building, testing, and running.
- **Disk Persistance** for easy backups

---

## Getting Started

### Prerequisites
- Go 1.21+
- `protoc` compiler installed
- `protoc-gen-go` and `protoc-gen-go-grpc` plugins installed

Install plugins if needed:
```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

Ensure `protoc` is in your PATH.

---

## Project Structure

```
kvstore/
├── kvstore/                # Core storage logic
│    ├── kvstore.go          # KV store implementation
│    └── storage.go          # Storage interface
├── server/                  # gRPC server wrapper
│    ├── server.go           # gRPC service + Listen
│    ├── hooks.go            # PreHookFunc and PostHookFunc
│    └── options.go          # Functional options for server configuration
├── proto/                   # Protobuf definitions
│    ├── kvstore.proto
│    ├── kvstore.pb.go
│    └── kvstore_grpc.pb.go
├── Makefile                 # Build, test, run automation
├── go.mod
├── go.sum
└── README.md                # You are here
```

---

## Building the Project

To generate Go files from the .proto file and build the project:

```bash
make build
```

You can also manually run:
```bash
protoc --go_out=paths=source_relative:proto --go-grpc_out=paths=source_relative:proto --proto_path=proto proto/kvstore.proto
go build ./...
```
This will:
- Regenerate `proto/kvstore.pb.go` and `proto/kvstore_grpc.pb.go`
- Build the Go project

---

## Running Tests

Run all tests (unit + integration):

```bash
make test
```

You can also manually run:
```bash
go test -v -cover ./...
```

- `-v` : verbose output
- `-race` : detect race conditions
- `-cover` : show test coverage

---

## Usage Example (Client Side)

After running the server, you can connect using a gRPC client.

```go
package main

import (
	"log"
	"time"

	"github.com/ahmad-masud/KVStore/server"
)

func main() {
	// Create the server with options
	s := server.NewServer(
		server.WithDefaultTTL(5*time.Minute),
	)

	// Start listening on a port
	log.Println("Starting KVStore server on :50051...")
	if err := s.Listen(":50051"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
```

---

## Hooks (Advanced Customization)

You can inject custom logic before and after every operation.

Example PreHook:
```go
func authHook(ctx context.Context, method string, req interface{}) error {
	if method == "Delete" {
		return status.Error(codes.PermissionDenied, "Delete not allowed")
	}
	return nil
}
```

Pass it to your server:
```go
s := server.NewServer(
	server.WithPreHook(authHook),
)
```

---

## Functional Options

Available options:
- `WithStorage(storage kvstore.Storage)` - Use a custom storage backend
- `WithPreHook(hook server.PreHookFunc)` - Inject logic before operations
- `WithPostHook(hook server.PostHookFunc)` - Inject logic after successful operations
- `WithDefaultTTL(ttl time.Duration)` - Set a default TTL for all keys

Example:
```go
s := server.NewServer(
	server.WithDefaultTTL(10*time.Minute),
)
```

---

## Contributing

Feel free to open issues or pull requests!

Future plans:
- Optional TTL background cleaning goroutine
- Metrics / Prometheus support
- Clustered / distributed version

---

## License

MIT License.