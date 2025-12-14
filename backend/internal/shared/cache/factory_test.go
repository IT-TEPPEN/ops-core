package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCache_Memory(t *testing.T) {
	config := CacheConfig{
		Type:            "memory",
		MaxEntries:      100,
		CleanupInterval: 1 * time.Minute,
	}

	cache, err := NewCache(config)
	require.NoError(t, err)
	assert.NotNil(t, cache)
	assert.IsType(t, &MemoryCache{}, cache)
}

func TestNewCache_DefaultToMemory(t *testing.T) {
	config := CacheConfig{
		Type: "", // Empty should default to memory
	}

	cache, err := NewCache(config)
	require.NoError(t, err)
	assert.NotNil(t, cache)
	assert.IsType(t, &MemoryCache{}, cache)
}

func TestNewCache_Redis(t *testing.T) {
	config := CacheConfig{
		Type:       "redis",
		RedisURL:   "redis://localhost:6379",
		RedisTLS:   false,
	}

	_, err := NewCache(config)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "redis cache requires a RedisClient")
}

func TestNewCache_InvalidType(t *testing.T) {
	config := CacheConfig{
		Type: "invalid",
	}

	_, err := NewCache(config)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported cache type")
}

func TestCacheKeyBuilder(t *testing.T) {
	builder := NewCacheKeyBuilder("test:")

	tests := []struct {
		name     string
		method   string
		args     []string
		expected string
	}{
		{
			name:     "Document key",
			method:   "DocumentKey",
			args:     []string{"doc123"},
			expected: "test:document:doc123",
		},
		{
			name:     "Document version key",
			method:   "DocumentVersionKey",
			args:     []string{"doc123", "v1"},
			expected: "test:document_version:doc123:v1",
		},
		{
			name:     "Document versions key",
			method:   "DocumentVersionsKey",
			args:     []string{"doc123"},
			expected: "test:document_versions:doc123",
		},
		{
			name:     "View statistics key",
			method:   "ViewStatisticsKey",
			args:     []string{"doc123"},
			expected: "test:view_stats:doc123",
		},
		{
			name:     "User key",
			method:   "UserKey",
			args:     []string{"user123"},
			expected: "test:user:user123",
		},
		{
			name:     "Group key",
			method:   "GroupKey",
			args:     []string{"group123"},
			expected: "test:group:group123",
		},
		{
			name:     "User groups key",
			method:   "UserGroupsKey",
			args:     []string{"user123"},
			expected: "test:user_groups:user123",
		},
		{
			name:     "Repository documents key",
			method:   "RepositoryDocumentsKey",
			args:     []string{"repo123"},
			expected: "test:repo_docs:repo123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result string
			switch tt.method {
			case "DocumentKey":
				result = builder.DocumentKey(tt.args[0])
			case "DocumentVersionKey":
				result = builder.DocumentVersionKey(tt.args[0], tt.args[1])
			case "DocumentVersionsKey":
				result = builder.DocumentVersionsKey(tt.args[0])
			case "ViewStatisticsKey":
				result = builder.ViewStatisticsKey(tt.args[0])
			case "UserKey":
				result = builder.UserKey(tt.args[0])
			case "GroupKey":
				result = builder.GroupKey(tt.args[0])
			case "UserGroupsKey":
				result = builder.UserGroupsKey(tt.args[0])
			case "RepositoryDocumentsKey":
				result = builder.RepositoryDocumentsKey(tt.args[0])
			}
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCacheKeyBuilder_PopularDocuments(t *testing.T) {
	builder := NewCacheKeyBuilder("test:")
	key := builder.PopularDocumentsKey(10)
	assert.Equal(t, "test:popular_docs:10", key)
}

func TestCacheKeyBuilder_RecentDocuments(t *testing.T) {
	builder := NewCacheKeyBuilder("test:")
	key := builder.RecentDocumentsKey("user123", 5)
	assert.Equal(t, "test:recent_docs:user123:5", key)
}

func TestDefaultTTLConfig(t *testing.T) {
	config := DefaultTTLConfig()

	assert.Equal(t, 1*time.Hour, config.Document)
	assert.Equal(t, 3*time.Hour, config.DocumentVersion)
	assert.Equal(t, 30*time.Minute, config.ViewStatistics)
	assert.Equal(t, 2*time.Hour, config.User)
	assert.Equal(t, 2*time.Hour, config.Group)
}
