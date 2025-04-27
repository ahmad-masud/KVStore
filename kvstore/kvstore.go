package kvstore

import (
	"sync"
	"time"
)

type item struct {
	value     string
	expiresAt time.Time
}

type KVStore struct {
	mu    sync.RWMutex
	store map[string]item
}

func New() *KVStore {
	return &KVStore{
		store: make(map[string]item),
	}
}

func (kv *KVStore) Set(key, value string) {
	kv.SetWithTTL(key, value, 0)
}

func (kv *KVStore) SetWithTTL(key, value string, ttl time.Duration) {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	var expiresAt time.Time
	if ttl > 0 {
		expiresAt = time.Now().Add(ttl)
	}

	kv.store[key] = item{
		value:     value,
		expiresAt: expiresAt,
	}
}

func (kv *KVStore) Get(key string) (string, bool) {
	kv.mu.RLock()
	defer kv.mu.RUnlock()

	it, ok := kv.store[key]
	if !ok {
		return "", false
	}
	if !it.expiresAt.IsZero() && time.Now().After(it.expiresAt) {
		go kv.deleteKeyAsync(key)
		return "", false
	}
	return it.value, true
}

func (kv *KVStore) Delete(key string) bool {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	if _, ok := kv.store[key]; ok {
		delete(kv.store, key)
		return true
	}
	return false
}

func (kv *KVStore) deleteKeyAsync(key string) {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	delete(kv.store, key)
}
