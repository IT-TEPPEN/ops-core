package cache

import (
	"context"
	"time"
)

// Cache defines the interface for caching operations.
// Implementations should be thread-safe.
type Cache interface {
	// Get retrieves a value from the cache.
	// Returns nil if the key does not exist or has expired.
	Get(ctx context.Context, key string) (interface{}, error)

	// Set stores a value in the cache with the specified TTL.
	// A TTL of 0 means no expiration.
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error

	// Delete removes a value from the cache.
	Delete(ctx context.Context, key string) error

	// Clear removes all values from the cache.
	Clear(ctx context.Context) error

	// Exists checks if a key exists in the cache.
	Exists(ctx context.Context, key string) (bool, error)

	// GetMultiple retrieves multiple values from the cache.
	// Returns a map of key-value pairs for existing keys.
	GetMultiple(ctx context.Context, keys []string) (map[string]interface{}, error)

	// SetMultiple stores multiple values in the cache with the specified TTL.
	SetMultiple(ctx context.Context, items map[string]interface{}, ttl time.Duration) error

	// DeleteMultiple removes multiple values from the cache.
	DeleteMultiple(ctx context.Context, keys []string) error
}

// CacheConfig holds configuration for cache implementations.
type CacheConfig struct {
	// Type specifies the cache type: "memory" or "redis"
	Type string

	// DefaultTTL is the default time-to-live for cache entries
	DefaultTTL time.Duration

	// Redis-specific configuration
	RedisURL      string
	RedisPassword string
	RedisTLS      bool

	// Memory cache-specific configuration
	MaxEntries      int // Maximum number of entries (0 = unlimited)
	CleanupInterval time.Duration
}

// TTLConfig holds TTL settings for different entity types.
type TTLConfig struct {
	Document       time.Duration
	DocumentVersion time.Duration
	ViewStatistics time.Duration
	User           time.Duration
	Group          time.Duration
}

// DefaultTTLConfig returns the default TTL configuration.
func DefaultTTLConfig() TTLConfig {
	return TTLConfig{
		Document:       1 * time.Hour,      // 1 hour
		DocumentVersion: 3 * time.Hour,     // 3 hours
		ViewStatistics: 30 * time.Minute,   // 30 minutes
		User:           2 * time.Hour,      // 2 hours
		Group:          2 * time.Hour,      // 2 hours
	}
}
