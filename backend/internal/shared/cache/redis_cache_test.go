package cache

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockRedisClient is a mock implementation of RedisClient for testing
type MockRedisClient struct {
	data      map[string]string
	expiryMap map[string]time.Time
}

func NewMockRedisClient() *MockRedisClient {
	return &MockRedisClient{
		data:      make(map[string]string),
		expiryMap: make(map[string]time.Time),
	}
}

func (m *MockRedisClient) Get(ctx context.Context, key string) (string, error) {
	if expiry, ok := m.expiryMap[key]; ok && time.Now().After(expiry) {
		delete(m.data, key)
		delete(m.expiryMap, key)
		return "", ErrKeyNotFound
	}
	val, ok := m.data[key]
	if !ok {
		return "", ErrKeyNotFound
	}
	return val, nil
}

func (m *MockRedisClient) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	m.data[key] = value
	if expiration > 0 {
		m.expiryMap[key] = time.Now().Add(expiration)
	}
	return nil
}

func (m *MockRedisClient) Del(ctx context.Context, keys ...string) error {
	for _, key := range keys {
		delete(m.data, key)
		delete(m.expiryMap, key)
	}
	return nil
}

func (m *MockRedisClient) FlushAll(ctx context.Context) error {
	m.data = make(map[string]string)
	m.expiryMap = make(map[string]time.Time)
	return nil
}

func (m *MockRedisClient) Exists(ctx context.Context, key string) (bool, error) {
	if expiry, ok := m.expiryMap[key]; ok && time.Now().After(expiry) {
		delete(m.data, key)
		delete(m.expiryMap, key)
		return false, nil
	}
	_, ok := m.data[key]
	return ok, nil
}

func (m *MockRedisClient) MGet(ctx context.Context, keys ...string) ([]interface{}, error) {
	result := make([]interface{}, len(keys))
	for i, key := range keys {
		if expiry, ok := m.expiryMap[key]; ok && time.Now().After(expiry) {
			delete(m.data, key)
			delete(m.expiryMap, key)
			result[i] = nil
			continue
		}
		if val, ok := m.data[key]; ok {
			result[i] = val
		} else {
			result[i] = nil
		}
	}
	return result, nil
}

func (m *MockRedisClient) MSet(ctx context.Context, pairs ...interface{}) error {
	for i := 0; i < len(pairs); i += 2 {
		key := pairs[i].(string)
		value := pairs[i+1].(string)
		m.data[key] = value
	}
	return nil
}

func (m *MockRedisClient) Expire(ctx context.Context, key string, expiration time.Duration) error {
	if _, ok := m.data[key]; !ok {
		return errors.New("key not found")
	}
	m.expiryMap[key] = time.Now().Add(expiration)
	return nil
}

func (m *MockRedisClient) Close() error {
	return nil
}

func TestRedisCache_GetSet(t *testing.T) {
	ctx := context.Background()
	client := NewMockRedisClient()
	cache := NewRedisCache(client, "test:")

	// Test setting and getting a value
	err := cache.Set(ctx, "key1", "value1", 5*time.Minute)
	require.NoError(t, err)

	val, err := cache.Get(ctx, "key1")
	require.NoError(t, err)
	assert.Equal(t, "value1", val)
}

func TestRedisCache_GetNonExistent(t *testing.T) {
	ctx := context.Background()
	client := NewMockRedisClient()
	cache := NewRedisCache(client, "test:")

	val, err := cache.Get(ctx, "nonexistent")
	require.NoError(t, err)
	assert.Nil(t, val)
}

func TestRedisCache_Delete(t *testing.T) {
	ctx := context.Background()
	client := NewMockRedisClient()
	cache := NewRedisCache(client, "test:")

	// Set a value
	err := cache.Set(ctx, "key1", "value1", 5*time.Minute)
	require.NoError(t, err)

	// Delete the value
	err = cache.Delete(ctx, "key1")
	require.NoError(t, err)

	// Value should not exist
	val, err := cache.Get(ctx, "key1")
	require.NoError(t, err)
	assert.Nil(t, val)
}

func TestRedisCache_Clear(t *testing.T) {
	ctx := context.Background()
	client := NewMockRedisClient()
	cache := NewRedisCache(client, "test:")

	// Set multiple values
	cache.Set(ctx, "key1", "value1", 5*time.Minute)
	cache.Set(ctx, "key2", "value2", 5*time.Minute)
	cache.Set(ctx, "key3", "value3", 5*time.Minute)

	// Clear cache
	err := cache.Clear(ctx)
	require.NoError(t, err)

	// All values should be gone
	val, _ := cache.Get(ctx, "key1")
	assert.Nil(t, val)
	val, _ = cache.Get(ctx, "key2")
	assert.Nil(t, val)
	val, _ = cache.Get(ctx, "key3")
	assert.Nil(t, val)
}

func TestRedisCache_Exists(t *testing.T) {
	ctx := context.Background()
	client := NewMockRedisClient()
	cache := NewRedisCache(client, "test:")

	// Key should not exist initially
	exists, err := cache.Exists(ctx, "key1")
	require.NoError(t, err)
	assert.False(t, exists)

	// Set a value
	cache.Set(ctx, "key1", "value1", 5*time.Minute)

	// Key should exist now
	exists, err = cache.Exists(ctx, "key1")
	require.NoError(t, err)
	assert.True(t, exists)
}

func TestRedisCache_GetMultiple(t *testing.T) {
	ctx := context.Background()
	client := NewMockRedisClient()
	cache := NewRedisCache(client, "test:")

	// Set multiple values
	cache.Set(ctx, "key1", "value1", 5*time.Minute)
	cache.Set(ctx, "key2", "value2", 5*time.Minute)
	cache.Set(ctx, "key3", "value3", 5*time.Minute)

	// Get multiple values
	values, err := cache.GetMultiple(ctx, []string{"key1", "key2", "key3", "nonexistent"})
	require.NoError(t, err)
	assert.Len(t, values, 3)
	assert.Equal(t, "value1", values["key1"])
	assert.Equal(t, "value2", values["key2"])
	assert.Equal(t, "value3", values["key3"])
	assert.Nil(t, values["nonexistent"])
}

func TestRedisCache_SetMultiple(t *testing.T) {
	ctx := context.Background()
	client := NewMockRedisClient()
	cache := NewRedisCache(client, "test:")

	// Set multiple values at once
	items := map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}
	err := cache.SetMultiple(ctx, items, 5*time.Minute)
	require.NoError(t, err)

	// Verify all values are set
	val, _ := cache.Get(ctx, "key1")
	assert.Equal(t, "value1", val)
	val, _ = cache.Get(ctx, "key2")
	assert.Equal(t, "value2", val)
	val, _ = cache.Get(ctx, "key3")
	assert.Equal(t, "value3", val)
}

func TestRedisCache_DeleteMultiple(t *testing.T) {
	ctx := context.Background()
	client := NewMockRedisClient()
	cache := NewRedisCache(client, "test:")

	// Set multiple values
	cache.Set(ctx, "key1", "value1", 5*time.Minute)
	cache.Set(ctx, "key2", "value2", 5*time.Minute)
	cache.Set(ctx, "key3", "value3", 5*time.Minute)

	// Delete multiple values
	err := cache.DeleteMultiple(ctx, []string{"key1", "key3"})
	require.NoError(t, err)

	// key1 and key3 should be gone, key2 should remain
	val, _ := cache.Get(ctx, "key1")
	assert.Nil(t, val)
	val, _ = cache.Get(ctx, "key2")
	assert.Equal(t, "value2", val)
	val, _ = cache.Get(ctx, "key3")
	assert.Nil(t, val)
}

func TestRedisCache_ComplexValues(t *testing.T) {
	ctx := context.Background()
	client := NewMockRedisClient()
	cache := NewRedisCache(client, "test:")

	// Test with map
	mapVal := map[string]interface{}{
		"name": "test",
		"age":  float64(30), // JSON unmarshals numbers as float64
	}
	err := cache.Set(ctx, "map", mapVal, 5*time.Minute)
	require.NoError(t, err)

	val, err := cache.Get(ctx, "map")
	require.NoError(t, err)
	assert.Equal(t, mapVal, val)
}

func TestRedisCache_PrefixIsolation(t *testing.T) {
	ctx := context.Background()
	client := NewMockRedisClient()
	cache1 := NewRedisCache(client, "prefix1:")
	cache2 := NewRedisCache(client, "prefix2:")

	// Set values with different prefixes
	cache1.Set(ctx, "key", "value1", 5*time.Minute)
	cache2.Set(ctx, "key", "value2", 5*time.Minute)

	// Both values should coexist
	val1, _ := cache1.Get(ctx, "key")
	val2, _ := cache2.Get(ctx, "key")
	assert.Equal(t, "value1", val1)
	assert.Equal(t, "value2", val2)
}

func TestRedisCache_Close(t *testing.T) {
	client := NewMockRedisClient()
	cache := NewRedisCache(client, "test:")

	err := cache.Close()
	assert.NoError(t, err)
}
