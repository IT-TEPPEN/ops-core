package cache

import (
	"fmt"
	"time"
)

// NewCache creates a new cache instance based on the provided configuration.
func NewCache(config CacheConfig) (Cache, error) {
	switch config.Type {
	case "memory", "":
		// Default to memory cache
		maxEntries := config.MaxEntries
		if maxEntries == 0 {
			maxEntries = 10000 // Default max entries
		}

		cleanupInterval := config.CleanupInterval
		if cleanupInterval == 0 {
			cleanupInterval = 10 * time.Minute
		}

		return NewMemoryCache(maxEntries, cleanupInterval), nil

	case "redis":
		// Redis cache requires a client to be provided separately
		// This is because the Redis client setup may vary based on the application
		return nil, fmt.Errorf("redis cache requires a RedisClient to be provided via NewRedisCache()")

	default:
		return nil, fmt.Errorf("unsupported cache type: %s", config.Type)
	}
}

// CacheKeyBuilder provides helper methods for building cache keys.
type CacheKeyBuilder struct {
	prefix string
}

// NewCacheKeyBuilder creates a new cache key builder.
func NewCacheKeyBuilder(prefix string) *CacheKeyBuilder {
	return &CacheKeyBuilder{prefix: prefix}
}

// DocumentKey builds a cache key for a document.
func (b *CacheKeyBuilder) DocumentKey(id string) string {
	return fmt.Sprintf("%sdocument:%s", b.prefix, id)
}

// DocumentVersionKey builds a cache key for a document version.
func (b *CacheKeyBuilder) DocumentVersionKey(documentID string, version string) string {
	return fmt.Sprintf("%sdocument_version:%s:%s", b.prefix, documentID, version)
}

// DocumentVersionsKey builds a cache key for all versions of a document.
func (b *CacheKeyBuilder) DocumentVersionsKey(documentID string) string {
	return fmt.Sprintf("%sdocument_versions:%s", b.prefix, documentID)
}

// ViewStatisticsKey builds a cache key for view statistics.
func (b *CacheKeyBuilder) ViewStatisticsKey(documentID string) string {
	return fmt.Sprintf("%sview_stats:%s", b.prefix, documentID)
}

// UserKey builds a cache key for a user.
func (b *CacheKeyBuilder) UserKey(id string) string {
	return fmt.Sprintf("%suser:%s", b.prefix, id)
}

// GroupKey builds a cache key for a group.
func (b *CacheKeyBuilder) GroupKey(id string) string {
	return fmt.Sprintf("%sgroup:%s", b.prefix, id)
}

// UserGroupsKey builds a cache key for user's groups.
func (b *CacheKeyBuilder) UserGroupsKey(userID string) string {
	return fmt.Sprintf("%suser_groups:%s", b.prefix, userID)
}

// RepositoryDocumentsKey builds a cache key for repository's documents.
func (b *CacheKeyBuilder) RepositoryDocumentsKey(repositoryID string) string {
	return fmt.Sprintf("%srepo_docs:%s", b.prefix, repositoryID)
}

// PopularDocumentsKey builds a cache key for popular documents.
func (b *CacheKeyBuilder) PopularDocumentsKey(limit int) string {
	return fmt.Sprintf("%spopular_docs:%d", b.prefix, limit)
}

// RecentDocumentsKey builds a cache key for recently viewed documents.
func (b *CacheKeyBuilder) RecentDocumentsKey(userID string, limit int) string {
	return fmt.Sprintf("%srecent_docs:%s:%d", b.prefix, userID, limit)
}
