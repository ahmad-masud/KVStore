package kvstore

import (
	"testing"
	"time"
)

func TestKVStore_SetAndGet(t *testing.T) {
	store := New()

	store.Set("foo", "bar")
	val, ok := store.Get("foo")
	if !ok {
		t.Fatalf("expected key to exist")
	}
	if val != "bar" {
		t.Fatalf("expected value 'bar', got '%s'", val)
	}
}

func TestKVStore_Delete(t *testing.T) {
	store := New()

	store.Set("foo", "bar")
	ok := store.Delete("foo")
	if !ok {
		t.Fatalf("expected delete to succeed")
	}

	_, exists := store.Get("foo")
	if exists {
		t.Fatalf("expected key to be deleted")
	}
}

func TestKVStore_SetWithTTL(t *testing.T) {
	store := New()

	store.SetWithTTL("foo", "bar", 100*time.Millisecond)

	val, ok := store.Get("foo")
	if !ok || val != "bar" {
		t.Fatalf("expected key to exist immediately")
	}

	time.Sleep(150 * time.Millisecond)

	val, ok = store.Get("foo")
	if ok {
		t.Fatalf("expected key to expire, got value '%s'", val)
	}
}
