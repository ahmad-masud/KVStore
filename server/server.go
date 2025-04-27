package server

import (
	"context"
	"net"
	"time"

	"github.com/ahmad-masud/KVStore/kvstore"
	"github.com/ahmad-masud/KVStore/proto"

	"google.golang.org/grpc"
)

// PreHookFunc is a function type that defines the signature for pre-hook functions.
// It takes a context, method name, and request object as parameters and returns an error if any.
type Server struct {
	proto.UnimplementedKVStoreServer

	storage    kvstore.Storage
	preHook    PreHookFunc
	postHook   PostHookFunc
	defaultTTL time.Duration
}

// PreHookFunc is a function type that defines the signature for pre-hook functions.
// It takes a context, method name, and request object as parameters and returns an error if any.
func NewServer(opts ...Option) *Server {
	s := &Server{
		storage: kvstore.New(),
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *Server) Set(ctx context.Context, req *proto.SetRequest) (*proto.SetResponse, error) {
	// Run pre-hook if provided
	if s.preHook != nil {
		if err := s.preHook(ctx, "Set", req); err != nil {
			return nil, err
		}
	}

	// Determine TTL
	var ttl time.Duration
	if req.TtlSeconds > 0 {
		ttl = time.Duration(req.TtlSeconds) * time.Second
	} else if s.defaultTTL > 0 {
		ttl = s.defaultTTL
	}

	// Perform the storage operation
	s.storage.SetWithTTL(req.Key, req.Value, ttl)

	// Prepare response
	resp := &proto.SetResponse{Success: true}

	// Run post-hook if provided
	if s.postHook != nil {
		_ = s.postHook(ctx, "Set", req, resp)
	}

	return resp, nil
}

func (s *Server) Get(ctx context.Context, req *proto.GetRequest) (*proto.GetResponse, error) {
	// Run pre-hook if provided
	if s.preHook != nil {
		if err := s.preHook(ctx, "Get", req); err != nil {
			return nil, err
		}
	}

	// Perform the storage operation
	value, found := s.storage.Get(req.Key)

	// Prepare response
	resp := &proto.GetResponse{
		Value: value,
		Found: found,
	}

	// Run post-hook if provided
	if s.postHook != nil {
		_ = s.postHook(ctx, "Get", req, resp)
	}

	return resp, nil
}

func (s *Server) Delete(ctx context.Context, req *proto.DeleteRequest) (*proto.DeleteResponse, error) {
	// Run pre-hook if provided
	if s.preHook != nil {
		if err := s.preHook(ctx, "Delete", req); err != nil {
			return nil, err
		}
	}

	// Perform the storage operation
	success := s.storage.Delete(req.Key)

	// Prepare response
	resp := &proto.DeleteResponse{
		Success: success,
	}

	// Run post-hook if provided
	if s.postHook != nil {
		_ = s.postHook(ctx, "Delete", req, resp)
	}

	return resp, nil
}

func (s *Server) Listen(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	proto.RegisterKVStoreServer(grpcServer, s)

	return grpcServer.Serve(lis)
}
