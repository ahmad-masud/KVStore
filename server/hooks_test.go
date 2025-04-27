package server

import (
	"context"
	"errors"
	"testing"
)

func TestPreHook_AllowsOperation(t *testing.T) {
	var called bool

	hook := func(ctx context.Context, method string, req interface{}) error {
		called = true
		return nil
	}

	s := &Server{preHook: hook}

	err := s.preHook(context.Background(), "Set", nil)
	if err != nil {
		t.Fatalf("expected no error from preHook, got: %v", err)
	}
	if !called {
		t.Fatalf("expected preHook to be called")
	}
}

func TestPreHook_BlocksOperation(t *testing.T) {
	hook := func(ctx context.Context, method string, req interface{}) error {
		return errors.New("blocked by preHook")
	}

	s := &Server{preHook: hook}

	err := s.preHook(context.Background(), "Set", nil)
	if err == nil {
		t.Fatalf("expected error from preHook, got nil")
	}
}

func TestPostHook_Called(t *testing.T) {
	var called bool

	hook := func(ctx context.Context, method string, req interface{}, resp interface{}) error {
		called = true
		return nil
	}

	s := &Server{postHook: hook}

	err := s.postHook(context.Background(), "Set", nil, nil)
	if err != nil {
		t.Fatalf("expected no error from postHook, got: %v", err)
	}
	if !called {
		t.Fatalf("expected postHook to be called")
	}
}
