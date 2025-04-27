package kvstore

import "time"

// Storage is an interface that defines the methods for a key-value store.
// It allows setting, getting, and deleting key-value pairs, as well as setting a value with a time-to-live (TTL).
// The interface is designed to be implemented by different storage backends, such as in-memory, Redis, or any other key-value store.
// The methods are designed to be simple and efficient, allowing for easy integration with various storage solutions.
type Storage interface {
	Set(key, value string)
	SetWithTTL(key, value string, ttl time.Duration)
	Get(key string) (string, bool)
	Delete(key string) bool
}
