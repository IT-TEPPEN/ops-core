package cache

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemoryCache_GetSet(t *testing.T) {
	ctx := context.Background()
	cache := NewMemoryCache(100, 1*time.Minute)
	defer cache.Close()

	// Test setting and getting a value
	err := cache.Set(ctx, "key1", "value1", 5*time.Minute)
	require.NoError(t, err)

	val, err := cache.Get(ctx, "key1")
	require.NoError(t, err)
	assert.Equal(t, "value1", val)
}

func TestMemoryCache_GetNonExistent(t *testing.T) {
	ctx := context.Background()
	cache := NewMemoryCache(100, 1*time.Minute)
	defer cache.Close()

	val, err := cache.Get(ctx, "nonexistent")
	require.NoError(t, err)
	assert.Nil(t, val)
}

func TestMemoryCache_TTLExpiration(t *testing.T) {
	ctx := context.Background()
	cache := NewMemoryCache(100, 100*time.Millisecond)
	defer cache.Close()

	// Set a value with short TTL
	err := cache.Set(ctx, "key1", "value1", 200*time.Millisecond)
	require.NoError(t, err)

	// Value should exist immediately
	val, err := cache.Get(ctx, "key1")
	require.NoError(t, err)
	assert.Equal(t, "value1", val)

	// Wait for expiration
	time.Sleep(300 * time.Millisecond)

	// Value should be expired
	val, err = cache.Get(ctx, "key1")
	require.NoError(t, err)
	assert.Nil(t, val)
}

func TestMemoryCache_Delete(t *testing.T) {
	ctx := context.Background()
	cache := NewMemoryCache(100, 1*time.Minute)
	defer cache.Close()

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

func TestMemoryCache_Clear(t *testing.T) {
	ctx := context.Background()
	cache := NewMemoryCache(100, 1*time.Minute)
	defer cache.Close()

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

func TestMemoryCache_Exists(t *testing.T) {
	ctx := context.Background()
	cache := NewMemoryCache(100, 1*time.Minute)
	defer cache.Close()

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

func TestMemoryCache_GetMultiple(t *testing.T) {
	ctx := context.Background()
	cache := NewMemoryCache(100, 1*time.Minute)
	defer cache.Close()

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

func TestMemoryCache_SetMultiple(t *testing.T) {
	ctx := context.Background()
	cache := NewMemoryCache(100, 1*time.Minute)
	defer cache.Close()

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

func TestMemoryCache_DeleteMultiple(t *testing.T) {
	ctx := context.Background()
	cache := NewMemoryCache(100, 1*time.Minute)
	defer cache.Close()

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

func TestMemoryCache_MaxEntries(t *testing.T) {
	ctx := context.Background()
	cache := NewMemoryCache(3, 1*time.Minute) // Max 3 entries
	defer cache.Close()

	// Set 3 entries
	cache.Set(ctx, "key1", "value1", 5*time.Minute)
	cache.Set(ctx, "key2", "value2", 5*time.Minute)
	cache.Set(ctx, "key3", "value3", 5*time.Minute)

	// Adding a 4th entry should cause one old entry to be evicted
	cache.Set(ctx, "key4", "value4", 5*time.Minute)

	// Cache should have exactly 3 entries
	cache.mu.RLock()
	entryCount := len(cache.entries)
	cache.mu.RUnlock()
	assert.Equal(t, 3, entryCount)
}

func TestMemoryCache_ComplexValues(t *testing.T) {
	ctx := context.Background()
	cache := NewMemoryCache(100, 1*time.Minute)
	defer cache.Close()

	// Test with map
	mapVal := map[string]interface{}{
		"name": "test",
		"age":  30,
	}
	err := cache.Set(ctx, "map", mapVal, 5*time.Minute)
	require.NoError(t, err)

	val, err := cache.Get(ctx, "map")
	require.NoError(t, err)
	assert.Equal(t, mapVal, val)

	// Test with slice
	sliceVal := []string{"a", "b", "c"}
	err = cache.Set(ctx, "slice", sliceVal, 5*time.Minute)
	require.NoError(t, err)

	val, err = cache.Get(ctx, "slice")
	require.NoError(t, err)
	assert.Equal(t, sliceVal, val)
}

func TestMemoryCache_ZeroTTL(t *testing.T) {
	ctx := context.Background()
	cache := NewMemoryCache(100, 1*time.Minute)
	defer cache.Close()

	// Set a value with zero TTL (no expiration)
	err := cache.Set(ctx, "key1", "value1", 0)
	require.NoError(t, err)

	// Wait a bit
	time.Sleep(100 * time.Millisecond)

	// Value should still exist
	val, err := cache.Get(ctx, "key1")
	require.NoError(t, err)
	assert.Equal(t, "value1", val)
}
