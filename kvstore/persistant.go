package kvstore

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

// PersistentKVStore wraps a KVStore and adds disk persistence.
// It writes every Set and Delete operation to a log file and replays the log on startup.
type PersistentKVStore struct {
	memStore *KVStore   // in-memory store
	logFile  *os.File   // append-only log file
	mu       sync.Mutex // protects logFile writes
}

// NewPersistentKVStore creates a new PersistentKVStore, replaying any existing log to rebuild the in-memory store.
// The logPath specifies the file to be used for persistence.
func NewPersistentKVStore(logPath string) (*PersistentKVStore, error) {
	dir := filepath.Dir(logPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory for persistence: %w", err)
	}

	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open persistence file: %w", err)
	}

	store := New()
	p := &PersistentKVStore{
		memStore: store,
		logFile:  file,
	}

	// Replay the existing log to rebuild memory state
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		p.replayLine(line)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading persistence file: %w", err)
	}

	return p, nil
}

// replayLine processes a single log line and applies it to the in-memory store.
func (p *PersistentKVStore) replayLine(line string) {
	parts := strings.SplitN(line, " ", 4)
	if len(parts) < 2 {
		return
	}

	switch parts[0] {
	case "SET":
		if len(parts) < 3 {
			return
		}
		key, value := parts[1], parts[2]
		p.memStore.Set(key, value)
	case "SETTTL":
		if len(parts) < 4 {
			return
		}
		key, value := parts[1], parts[2]
		ttlMillis, err := strconv.ParseInt(parts[3], 10, 64)
		if err != nil {
			return
		}
		ttl := time.Duration(ttlMillis) * time.Millisecond
		p.memStore.SetWithTTL(key, value, ttl)
	case "DEL":
		if len(parts) < 2 {
			return
		}
		key := parts[1]
		p.memStore.Delete(key)
	}
}

// Set stores a key-value pair in the in-memory store and appends the operation to the log file.
func (p *PersistentKVStore) Set(key, value string) {
	p.memStore.Set(key, value)
	p.appendLog(fmt.Sprintf("SET %s %s\n", key, value))
}

// SetWithTTL stores a key-value pair with a TTL and appends the operation to the log file.
func (p *PersistentKVStore) SetWithTTL(key, value string, ttl time.Duration) {
	p.memStore.SetWithTTL(key, value, ttl)
	p.appendLog(fmt.Sprintf("SETTTL %s %s %d\n", key, value, ttl.Milliseconds()))
}

// Get retrieves the value associated with the key from the in-memory store.
func (p *PersistentKVStore) Get(key string) (string, bool) {
	return p.memStore.Get(key)
}

// Delete removes the key-value pair from the in-memory store and appends the operation to the log file.
func (p *PersistentKVStore) Delete(key string) bool {
	ok := p.memStore.Delete(key)
	if ok {
		p.appendLog(fmt.Sprintf("DEL %s\n", key))
	}
	return ok
}

// appendLog safely appends an operation to the log file and ensures it is flushed to disk.
func (p *PersistentKVStore) appendLog(entry string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.logFile.WriteString(entry)
	p.logFile.Sync() // ensure durability
}
