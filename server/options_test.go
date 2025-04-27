package server

import (
	"context"
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
