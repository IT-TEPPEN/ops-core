# Performance Optimization and Caching Strategy

This document describes the performance optimization features implemented in OpsCore, including caching strategies, database indexing, and pagination support.

## Overview

Performance optimization is critical for providing a responsive user experience, especially as the system scales to support more users and documents. The implementation includes:

- **Backend caching layer** with multiple storage backends
- **Database indexing** for frequently queried tables
- **Pagination support** with both offset-based and cursor-based strategies
- **Frontend client-side caching** for API responses

## Backend Caching

### Cache Interface

The caching layer is defined by the `Cache` interface in `backend/internal/shared/cache/cache.go`:

```go
type Cache interface {
    Get(ctx context.Context, key string) (interface{}, error)
    Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
    Delete(ctx context.Context, key string) error
    Clear(ctx context.Context) error
    Exists(ctx context.Context, key string) (bool, error)
    GetMultiple(ctx context.Context, keys []string) (map[string]interface{}, error)
    SetMultiple(ctx context.Context, items map[string]interface{}, ttl time.Duration) error
    DeleteMultiple(ctx context.Context, keys []string) error
}
```

### Cache Implementations

#### Memory Cache

- **Location**: `backend/internal/shared/cache/memory_cache.go`
- **Use case**: Single-instance deployments, development, testing
- **Features**:
  - Thread-safe in-memory storage
  - Automatic cleanup of expired entries
  - Configurable max entries limit
  - No external dependencies

**Configuration**:
```bash
CACHE_TYPE=memory
```

#### Redis Cache

- **Location**: `backend/internal/shared/cache/redis_cache.go`
- **Use case**: Distributed deployments, production environments
- **Features**:
  - Distributed caching across multiple instances
  - Persistence support
  - High performance
  - Scalable

**Configuration**:
```bash
CACHE_TYPE=redis
REDIS_URL=redis://localhost:6379/0
REDIS_PASSWORD=your_password
REDIS_TLS=true
```

### Cache TTL Configuration

Default TTL values for different entity types:

| Entity Type | TTL | Rationale |
|------------|-----|-----------|
| Document | 1 hour | Documents change infrequently |
| DocumentVersion | 3 hours | Versions are immutable once created |
| ViewStatistics | 30 minutes | Statistics update frequently |
| User | 2 hours | User data changes occasionally |
| Group | 2 hours | Group data changes occasionally |

**Configuration**:
```bash
CACHE_TTL_DOCUMENT=3600        # 1 hour in seconds
CACHE_TTL_DOCUMENT_VERSION=10800  # 3 hours
CACHE_TTL_STATISTICS=1800      # 30 minutes
CACHE_TTL_USER=7200           # 2 hours
CACHE_TTL_GROUP=7200          # 2 hours
```

### Cache Key Builder

Use the `CacheKeyBuilder` to generate consistent cache keys:

```go
builder := cache.NewCacheKeyBuilder("opscore:")

// Generate keys
docKey := builder.DocumentKey("doc-123")
versionKey := builder.DocumentVersionKey("doc-123", "v1")
statsKey := builder.ViewStatisticsKey("doc-123")
```

### Cache Invalidation Strategy

- **Update-triggered**: Clear cache immediately when data is modified
- **Time-based**: Automatic expiration based on TTL
- **Manual**: Explicit invalidation when needed

## Database Indexing

### Index Strategy

Database indexes are defined in `backend/internal/shared/db/indexes.go` and applied via migration `000008_add_performance_indexes.up.sql`. A total of 32 indexes cover frequently queried tables.

### Index Summary

| Table | Index Count | Purpose |
|-------|-------------|---------|
| documents | 5 | Filter by repository, owner, publication status |
| document_versions | 3 | Version lookup and commit tracking |
| execution_records | 6 | Search by document, user, status, dates |
| execution_steps | 2 | Step ordering and lookup |
| attachments | 2 | Link to records and steps |
| view_history | 4 | User views and aggregation |
| view_statistics | 4 | Popular and recent documents |
| users | 2 | Email and username uniqueness |
| groups | 1 | Group name search |
| group_members | 3 | Membership queries |

### Key Indexes

#### Document Indexes
- `idx_documents_repository_id`: Filter documents by repository
- `idx_documents_owner`: Filter documents by owner
- `idx_documents_is_published`: Filter published/unpublished documents
- `idx_documents_published`: Partial index for published documents (optimized for common queries)
- `idx_documents_created_at`: Sort and filter by creation date

#### Execution Record Indexes
- `idx_execution_records_document_id`: Find executions for a document
- `idx_execution_records_executor_id`: Find executions by user
- `idx_execution_records_status`: Filter by execution status
- `idx_execution_records_started_at`: Sort by start time
- `idx_execution_records_search`: Composite index for complex searches

#### View History Indexes
- `idx_view_history_document_id`: Document view history
- `idx_view_history_user_id`: User view history
- `idx_view_history_viewed_at`: Time-based queries
- `idx_view_history_aggregation`: Optimize statistics aggregation

#### View Statistics Indexes
- `idx_view_statistics_document_id`: Unique index per document
- `idx_view_statistics_total_views`: Sort by popularity
- `idx_view_statistics_last_viewed_at`: Recent documents
- `idx_view_statistics_popularity`: Composite index for popular documents

### Index Maintenance

Indexes are automatically created during database migration. To apply indexes:

```bash
cd backend
go run cmd/migrate/main.go up
```

## Pagination

### Pagination Utilities

Location: `backend/internal/shared/pagination/pagination.go`

### Offset-Based Pagination

Traditional page number-based pagination:

```go
// Create page from page number
page := pagination.NewPageFromPageNumber(2, 20) // Page 2, 20 items per page

// Create result
result := pagination.NewPageResult(items, total, page)
```

**Configuration**:
```bash
DEFAULT_PAGE_SIZE=20
MAX_PAGE_SIZE=100
```

### Cursor-Based Pagination

More efficient for large datasets:

```go
// Encode cursor
cursor, _ := pagination.EncodeCursor("last-item-id", timestamp)

// Decode cursor
decoded, _ := pagination.DecodeCursor(cursor)

// Create result
result := pagination.NewCursorPageResult(items, hasNext, nextCursor, pageSize)
```

**Configuration**:
```bash
CURSOR_PAGINATION_ENABLED=true
```

## Frontend Caching

### Client Cache Utility

Location: `frontend/src/utils/cache.ts`

### Usage

```typescript
import { cache, TTL, CacheKeys } from '@/utils/cache';

// Set a value
cache.set(CacheKeys.document('doc-123'), documentData, {
  ttl: TTL.ONE_HOUR,
  storage: 'memory'
});

// Get a value
const doc = cache.get(CacheKeys.document('doc-123'));

// Get or set with factory
const doc = await cache.getOrSet(
  CacheKeys.document('doc-123'),
  () => fetchDocument('doc-123'),
  { ttl: TTL.ONE_HOUR }
);

// Delete a value
cache.delete(CacheKeys.document('doc-123'));

// Clear all cache
cache.clear('all');
```

### Storage Options

- **memory**: In-memory cache (lost on page refresh)
- **localStorage**: Persistent cache (survives page refresh)

### TTL Constants

```typescript
TTL.ONE_MINUTE      // 1 minute
TTL.FIVE_MINUTES    // 5 minutes
TTL.FIFTEEN_MINUTES // 15 minutes
TTL.THIRTY_MINUTES  // 30 minutes
TTL.ONE_HOUR        // 1 hour
TTL.ONE_DAY         // 24 hours
```

### Predefined Cache Keys

```typescript
CacheKeys.document(id)
CacheKeys.documentVersion(docId, version)
CacheKeys.documentVersions(docId)
CacheKeys.user(id)
CacheKeys.repository(id)
CacheKeys.repositories()
CacheKeys.executionRecord(id)
CacheKeys.viewStatistics(docId)
```

## Performance Metrics

### Backend Metrics

Monitor these metrics to assess performance:

- **Database query time**: Average query execution time
- **Cache hit rate**: Percentage of requests served from cache
- **API response time**: Time to complete API requests

### Frontend Metrics

Use browser tools and Lighthouse to measure:

- **First Contentful Paint (FCP)**: Time to first content render
- **Largest Contentful Paint (LCP)**: Time to main content render
- **Cumulative Layout Shift (CLS)**: Visual stability
- **Time to Interactive (TTI)**: Time until page is interactive

## Best Practices

### When to Use Caching

✅ **Use cache for**:
- Frequently accessed, rarely changing data (documents, users)
- Expensive computations or queries
- External API responses
- Aggregated statistics

❌ **Don't cache**:
- Frequently changing data
- User-specific sensitive data (unless properly secured)
- Data that must be real-time

### Cache Invalidation

Always invalidate cache when:
- Data is updated
- Data is deleted
- Related data changes

Example:
```go
// After updating a document
cache.Delete(ctx, builder.DocumentKey(docID))
cache.Delete(ctx, builder.RepositoryDocumentsKey(repoID))
```

### Database Indexing

✅ **Add indexes for**:
- Foreign key columns
- Frequently used WHERE clauses
- ORDER BY columns
- JOIN columns

❌ **Avoid excessive indexes**:
- Rarely queried columns
- High-write tables (indexes slow down writes)
- Small tables (full scan may be faster)

## Testing

### Backend Tests

```bash
cd backend
go test ./internal/shared/cache/...
go test ./internal/shared/pagination/...
```

### Frontend Tests

```bash
cd frontend
npm test -- cache.test.ts
```

## Monitoring and Tuning

### Cache Hit Rate Monitoring

Track cache hit rate to optimize TTL values:

```go
hits := cacheHits.Load()
misses := cacheMisses.Load()
hitRate := float64(hits) / float64(hits + misses)
```

### Query Performance

Use PostgreSQL's `EXPLAIN ANALYZE` to verify index usage:

```sql
EXPLAIN ANALYZE
SELECT * FROM documents
WHERE repository_id = 'repo-123'
AND is_published = true;
```

### Adjusting TTL

If cache hit rate is low:
- Increase TTL for stable data
- Decrease TTL for frequently changing data

If data becomes stale:
- Decrease TTL
- Improve cache invalidation strategy

## Troubleshooting

### Cache Not Working

1. Check cache type configuration
2. Verify Redis connection (if using Redis)
3. Check cache key generation
4. Verify TTL is not too short

### Slow Queries

1. Check if indexes are created (`\di` in psql)
2. Use `EXPLAIN ANALYZE` to verify index usage
3. Consider adding composite indexes
4. Check table statistics are up to date

### Memory Issues

1. Monitor memory cache size
2. Set `MaxEntries` limit for memory cache
3. Consider using Redis for large datasets
4. Implement cache eviction policies

## Future Enhancements

Potential improvements for future phases:

- **LRU cache eviction**: Implement least-recently-used eviction
- **Cache warming**: Pre-populate cache on startup
- **Distributed cache invalidation**: Pub/sub for multi-instance invalidation
- **Query result caching**: Cache at database query level
- **Service Worker**: Offline support and background sync
- **CDN integration**: Cache static assets
- **GraphQL DataLoader**: Batch and cache data fetching
