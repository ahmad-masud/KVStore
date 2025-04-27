package kvstore

import (
	"os"
	"testing"
	"time"
)

func TestPersistentKVStore_SetGetDelete(t *testing.T) {
	// Setup temp file for testing
	tmpfile, err := os.CreateTemp("", "kvstore_test_log")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	store, err := NewPersistentKVStore(tmpfile.Name())
	if err != nil {
		t.Fatalf("failed to create PersistentKVStore: %v", err)
	}

	// Set key
	store.Set("foo", "bar")
	val, found := store.Get("foo")
	if !found || val != "bar" {
		t.Fatalf("expected to find key 'foo' with value 'bar', got found=%v val=%s", found, val)
	}

	// Delete key
	ok := store.Delete("foo")
	if !ok {
		t.Fatalf("expected delete to succeed")
	}
	_, found = store.Get("foo")
	if found {
		t.Fatalf("expected key to be deleted")
	}
}

func TestPersistentKVStore_Recovery(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "kvstore_test_log")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	// First instance: write data
	store, err := NewPersistentKVStore(tmpfile.Name())
	if err != nil {
		t.Fatalf("failed to create PersistentKVStore: %v", err)
	}
	store.Set("foo", "bar")
	store.SetWithTTL("baz", "qux", 2*time.Second)

	// Simulate server restart
	store.logFile.Close()

	store2, err := NewPersistentKVStore(tmpfile.Name())
	if err != nil {
		t.Fatalf("failed to recover PersistentKVStore: %v", err)
	}

	val, found := store2.Get("foo")
	if !found || val != "bar" {
		t.Fatalf("expected to recover key 'foo', got found=%v val=%s", found, val)
	}

	val, found = store2.Get("baz")
	if !found || val != "qux" {
		t.Fatalf("expected to recover key 'baz', got found=%v val=%s", found, val)
	}
}
