package kvstore

import "time"

type Storage interface {
	Set(key, value string)
	SetWithTTL(key, value string, ttlSeconds time.Duration)
	Get(key string) (string, bool)
	Delete(key string) bool
}
