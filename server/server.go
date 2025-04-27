package server

import (
	"context"
	"net"
	"time"

	"github.com/ahmad-masud/KVStore/kvstore"
	"github.com/ahmad-masud/KVStore/proto"

	"google.golang.org/grpc"
)

// Server is a gRPC server that handles key-value store operations.
// It wraps a Storage backend and supports optional hooks for customization.
type Server struct {
	proto.UnimplementedKVStoreServer

	storage    kvstore.Storage
	preHook    PreHookFunc
	postHook   PostHookFunc
	defaultTTL time.Duration
}

// NewServer creates a new Server instance with optional functional configuration.
// By default, it uses an in-memory storage backend.
func NewServer(opts ...Option) *Server {
	s := &Server{
		storage: kvstore.New(),
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// Set stores a key-value pair into the storage backend, optionally applying a TTL (time-to-live).
// If a PreHookFunc is set, it runs before the operation.
// If a PostHookFunc is set, it runs after a successful operation.
func (s *Server) Set(ctx context.Context, req *proto.SetRequest) (*proto.SetResponse, error) {
	if s.preHook != nil {
		if err := s.preHook(ctx, "Set", req); err != nil {
			return nil, err
		}
	}

	var ttl time.Duration
	if req.TtlSeconds > 0 {
		ttl = time.Duration(req.TtlSeconds) * time.Second
	} else if s.defaultTTL > 0 {
		ttl = s.defaultTTL
	}

	s.storage.SetWithTTL(req.Key, req.Value, ttl)

	resp := &proto.SetResponse{Success: true}

	if s.postHook != nil {
		_ = s.postHook(ctx, "Set", req, resp)
	}

	return resp, nil
}

// Get retrieves the value for a given key from the storage backend.
// If a PreHookFunc is set, it runs before the operation.
// If a PostHookFunc is set, it runs after retrieving the value.
func (s *Server) Get(ctx context.Context, req *proto.GetRequest) (*proto.GetResponse, error) {
	if s.preHook != nil {
		if err := s.preHook(ctx, "Get", req); err != nil {
			return nil, err
		}
	}

	value, found := s.storage.Get(req.Key)

	resp := &proto.GetResponse{
		Value: value,
		Found: found,
	}

	if s.postHook != nil {
		_ = s.postHook(ctx, "Get", req, resp)
	}

	return resp, nil
}

// Delete removes a key-value pair from the storage backend.
// If a PreHookFunc is set, it runs before the operation.
// If a PostHookFunc is set, it runs after a successful deletion.
func (s *Server) Delete(ctx context.Context, req *proto.DeleteRequest) (*proto.DeleteResponse, error) {
	if s.preHook != nil {
		if err := s.preHook(ctx, "Delete", req); err != nil {
			return nil, err
		}
	}

	success := s.storage.Delete(req.Key)

	resp := &proto.DeleteResponse{
		Success: success,
	}

	if s.postHook != nil {
		_ = s.postHook(ctx, "Delete", req, resp)
	}

	return resp, nil
}

// Listen starts the gRPC server on the specified TCP address (e.g., ":50051").
// It registers the KVStore service and begins serving incoming requests.
func (s *Server) Listen(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	proto.RegisterKVStoreServer(grpcServer, s)

	return grpcServer.Serve(lis)
}
