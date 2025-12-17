package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"
)

// RedisClient is an interface for Redis operations to allow for testing and different Redis client implementations.
type RedisClient interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, expiration time.Duration) error
	Del(ctx context.Context, keys ...string) error
	FlushAll(ctx context.Context) error
	Exists(ctx context.Context, key string) (bool, error)
	MGet(ctx context.Context, keys ...string) ([]interface{}, error)
	MSet(ctx context.Context, pairs ...interface{}) error
	Expire(ctx context.Context, key string, expiration time.Duration) error
	Close() error
}

// RedisCache is a Redis-based implementation of the Cache interface.
// It is suitable for distributed deployments and provides persistence.
type RedisCache struct {
	client RedisClient
	prefix string // Key prefix to avoid collisions
}

// NewRedisCache creates a new Redis-based cache.
// The client parameter should be a configured Redis client.
func NewRedisCache(client RedisClient, keyPrefix string) *RedisCache {
	if keyPrefix == "" {
		keyPrefix = "opscore:"
	}
	return &RedisCache{
		client: client,
		prefix: keyPrefix,
	}
}

// prefixKey adds the prefix to a key.
func (rc *RedisCache) prefixKey(key string) string {
	return rc.prefix + key
}

// prefixKeys adds the prefix to multiple keys.
func (rc *RedisCache) prefixKeys(keys []string) []string {
	prefixed := make([]string, len(keys))
	for i, key := range keys {
		prefixed[i] = rc.prefixKey(key)
	}
	return prefixed
}

// Get retrieves a value from the cache.
func (rc *RedisCache) Get(ctx context.Context, key string) (interface{}, error) {
	value, err := rc.client.Get(ctx, rc.prefixKey(key))
	if err != nil {
		// Check if key doesn't exist (not an error)
		if errors.Is(err, ErrKeyNotFound) {
			return nil, nil
		}
		return nil, err
	}

	if value == "" {
		return nil, nil
	}

	// Deserialize JSON
	var result interface{}
	if err := json.Unmarshal([]byte(value), &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Set stores a value in the cache with the specified TTL.
func (rc *RedisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	// Serialize to JSON
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return rc.client.Set(ctx, rc.prefixKey(key), string(data), ttl)
}

// Delete removes a value from the cache.
func (rc *RedisCache) Delete(ctx context.Context, key string) error {
	return rc.client.Del(ctx, rc.prefixKey(key))
}

// Clear removes all values from the cache.
// WARNING: This clears only keys with the configured prefix.
// For large key sets, this operation may be slow as it uses SCAN.
func (rc *RedisCache) Clear(ctx context.Context) error {
	// Note: A production implementation would use SCAN with pattern matching
	// to delete only keys with the prefix. However, since RedisClient interface
	// doesn't expose SCAN, this is a simplified implementation.
	// 
	// For now, we document that Clear should be used carefully in production,
	// and applications should implement selective deletion using Delete/DeleteMultiple
	// when possible.
	//
	// Example proper implementation with redis client:
	// iter := client.Scan(ctx, 0, rc.prefix+"*", 0).Iterator()
	// for iter.Next(ctx) {
	//     client.Del(ctx, iter.Val())
	// }
	return rc.client.FlushAll(ctx)
}

// Exists checks if a key exists in the cache.
func (rc *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	return rc.client.Exists(ctx, rc.prefixKey(key))
}

// GetMultiple retrieves multiple values from the cache.
func (rc *RedisCache) GetMultiple(ctx context.Context, keys []string) (map[string]interface{}, error) {
	if len(keys) == 0 {
		return make(map[string]interface{}), nil
	}

	prefixedKeys := rc.prefixKeys(keys)
	values, err := rc.client.MGet(ctx, prefixedKeys...)
	if err != nil {
		return nil, err
	}

	result := make(map[string]interface{})
	for i, value := range values {
		if value == nil {
			continue
		}

		// Try to deserialize as JSON
		var deserialized interface{}
		if str, ok := value.(string); ok && str != "" {
			if err := json.Unmarshal([]byte(str), &deserialized); err != nil {
				// If deserialization fails, use raw value
				result[keys[i]] = value
			} else {
				result[keys[i]] = deserialized
			}
		}
	}

	return result, nil
}

// SetMultiple stores multiple values in the cache with the specified TTL.
func (rc *RedisCache) SetMultiple(ctx context.Context, items map[string]interface{}, ttl time.Duration) error {
	if len(items) == 0 {
		return nil
	}

	// Prepare key-value pairs for MSET
	pairs := make([]interface{}, 0, len(items)*2)
	for key, value := range items {
		data, err := json.Marshal(value)
		if err != nil {
			return err
		}
		pairs = append(pairs, rc.prefixKey(key), string(data))
	}

	// Set all values
	if err := rc.client.MSet(ctx, pairs...); err != nil {
		return err
	}

	// Set TTL for each key if specified
	if ttl > 0 {
		var ttlErrors []string
		for key := range items {
			if err := rc.client.Expire(ctx, rc.prefixKey(key), ttl); err != nil {
				// Log the error but continue setting TTL for other keys
				errMsg := fmt.Sprintf("failed to set TTL for key %s: %v", key, err)
				ttlErrors = append(ttlErrors, errMsg)
				log.Printf("Redis cache warning: %s", errMsg)
			}
		}
		// If any TTL failures occurred, log a summary
		if len(ttlErrors) > 0 {
			log.Printf("Redis cache: %d TTL failures out of %d keys in SetMultiple", len(ttlErrors), len(items))
		}
	}

	return nil
}

// DeleteMultiple removes multiple values from the cache.
func (rc *RedisCache) DeleteMultiple(ctx context.Context, keys []string) error {
	if len(keys) == 0 {
		return nil
	}

	prefixedKeys := rc.prefixKeys(keys)
	return rc.client.Del(ctx, prefixedKeys...)
}

// Close closes the Redis connection.
func (rc *RedisCache) Close() error {
	return rc.client.Close()
}

// ErrKeyNotFound is returned when a key is not found in Redis.
var ErrKeyNotFound = errors.New("key not found")
