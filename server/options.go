package server

import (
	"time"

	"github.com/ahmadmasud/KVStore/kvstore"
)

// Option configures the Server.
type Option func(*Server)

// WithStorage allows injecting a custom storage backend.
func WithStorage(storage kvstore.Storage) Option {
	return func(s *Server) {
		s.storage = storage
	}
}

// WithPreHook sets a hook that runs before every operation.
func WithPreHook(hook PreHookFunc) Option {
	return func(s *Server) {
		s.preHook = hook
	}
}

// WithPostHook sets a hook that runs after every successful operation.
func WithPostHook(hook PostHookFunc) Option {
	return func(s *Server) {
		s.postHook = hook
	}
}

// WithDefaultTTL sets a default TTL for keys if none is specified.
func WithDefaultTTL(ttl time.Duration) Option {
	return func(s *Server) {
		s.defaultTTL = ttl
	}
}
