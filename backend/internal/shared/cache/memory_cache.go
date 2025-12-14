package cache

import (
	"context"
	"sync"
	"time"
)

// cacheEntry represents a cached value with expiration time.
type cacheEntry struct {
	value      interface{}
	expiration time.Time
}

// isExpired checks if the entry has expired.
func (e *cacheEntry) isExpired() bool {
	if e.expiration.IsZero() {
		return false // No expiration
	}
	return time.Now().After(e.expiration)
}

// MemoryCache is an in-memory implementation of the Cache interface.
// It is thread-safe and suitable for single-instance deployments.
type MemoryCache struct {
	mu              sync.RWMutex
	entries         map[string]*cacheEntry
	maxEntries      int
	cleanupInterval time.Duration
	stopCleanup     chan struct{}
}

// NewMemoryCache creates a new in-memory cache.
func NewMemoryCache(maxEntries int, cleanupInterval time.Duration) *MemoryCache {
	if cleanupInterval == 0 {
		cleanupInterval = 10 * time.Minute // Default cleanup interval
	}

	mc := &MemoryCache{
		entries:         make(map[string]*cacheEntry),
		maxEntries:      maxEntries,
		cleanupInterval: cleanupInterval,
		stopCleanup:     make(chan struct{}),
	}

	// Start background cleanup goroutine
	go mc.cleanupLoop()

	return mc
}

// cleanupLoop periodically removes expired entries.
func (mc *MemoryCache) cleanupLoop() {
	ticker := time.NewTicker(mc.cleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			mc.removeExpired()
		case <-mc.stopCleanup:
			return
		}
	}
}

// removeExpired removes all expired entries from the cache.
func (mc *MemoryCache) removeExpired() {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	for key, entry := range mc.entries {
		if entry.isExpired() {
			delete(mc.entries, key)
		}
	}
}

// Get retrieves a value from the cache.
func (mc *MemoryCache) Get(ctx context.Context, key string) (interface{}, error) {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	entry, exists := mc.entries[key]
	if !exists {
		return nil, nil
	}

	if entry.isExpired() {
		// Entry expired, will be cleaned up later
		return nil, nil
	}

	return entry.value, nil
}

// Set stores a value in the cache with the specified TTL.
func (mc *MemoryCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	// Check if we need to enforce max entries limit
	if mc.maxEntries > 0 && len(mc.entries) >= mc.maxEntries {
		if _, exists := mc.entries[key]; !exists {
			// Remove oldest entry (simple strategy: remove first entry)
			for k := range mc.entries {
				delete(mc.entries, k)
				break
			}
		}
	}

	var expiration time.Time
	if ttl > 0 {
		expiration = time.Now().Add(ttl)
	}

	mc.entries[key] = &cacheEntry{
		value:      value,
		expiration: expiration,
	}

	return nil
}

// Delete removes a value from the cache.
func (mc *MemoryCache) Delete(ctx context.Context, key string) error {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	delete(mc.entries, key)
	return nil
}

// Clear removes all values from the cache.
func (mc *MemoryCache) Clear(ctx context.Context) error {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	mc.entries = make(map[string]*cacheEntry)
	return nil
}

// Exists checks if a key exists in the cache.
func (mc *MemoryCache) Exists(ctx context.Context, key string) (bool, error) {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	entry, exists := mc.entries[key]
	if !exists {
		return false, nil
	}

	if entry.isExpired() {
		return false, nil
	}

	return true, nil
}

// GetMultiple retrieves multiple values from the cache.
func (mc *MemoryCache) GetMultiple(ctx context.Context, keys []string) (map[string]interface{}, error) {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	result := make(map[string]interface{})
	for _, key := range keys {
		entry, exists := mc.entries[key]
		if exists && !entry.isExpired() {
			result[key] = entry.value
		}
	}

	return result, nil
}

// SetMultiple stores multiple values in the cache with the specified TTL.
func (mc *MemoryCache) SetMultiple(ctx context.Context, items map[string]interface{}, ttl time.Duration) error {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	var expiration time.Time
	if ttl > 0 {
		expiration = time.Now().Add(ttl)
	}

	for key, value := range items {
		// Check max entries limit
		if mc.maxEntries > 0 && len(mc.entries) >= mc.maxEntries {
			if _, exists := mc.entries[key]; !exists {
				// Remove oldest entry
				for k := range mc.entries {
					delete(mc.entries, k)
					break
				}
			}
		}

		mc.entries[key] = &cacheEntry{
			value:      value,
			expiration: expiration,
		}
	}

	return nil
}

// DeleteMultiple removes multiple values from the cache.
func (mc *MemoryCache) DeleteMultiple(ctx context.Context, keys []string) error {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	for _, key := range keys {
		delete(mc.entries, key)
	}

	return nil
}

// Close stops the cleanup goroutine and cleans up resources.
func (mc *MemoryCache) Close() error {
	close(mc.stopCleanup)
	return nil
}
