package server

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/ahmad-masud/KVStore/kvstore"
)

func TestWithStorage(t *testing.T) {
	customStore := kvstore.New()

	s := &Server{}
	opt := WithStorage(customStore)
	opt(s)

	if s.storage != customStore {
		t.Fatalf("expected storage to be set")
	}
}

func TestWithPreHook(t *testing.T) {
	hook := func(ctx context.Context, method string, req interface{}) error {
		return nil
	}

	s := &Server{}
	opt := WithPreHook(hook)
	opt(s)

	if s.preHook == nil {
		t.Fatalf("expected preHook to be set")
	}
}

func TestWithPostHook(t *testing.T) {
	hook := func(ctx context.Context, method string, req, resp interface{}) error {
		return nil
	}

	s := &Server{}
	opt := WithPostHook(hook)
	opt(s)

	if s.postHook == nil {
		t.Fatalf("expected postHook to be set")
	}
}

func TestWithDefaultTTL(t *testing.T) {
	ttl := 5 * time.Minute

	s := &Server{}
	opt := WithDefaultTTL(ttl)
	opt(s)

	if s.defaultTTL != ttl {
		t.Fatalf("expected defaultTTL to be set")
	}
}

func TestWithDiskPersistence(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "kvstore_test_log")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	s := NewServer(
		WithDiskPersistence(tmpfile.Name(), true),
	)

	if s.storage == nil {
		t.Fatalf("expected storage to be initialized")
	}
}
