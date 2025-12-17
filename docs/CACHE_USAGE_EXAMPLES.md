# Cache Usage Examples

This document provides practical examples of using the caching features in OpsCore.

## Backend Examples

### Basic Cache Usage

```go
import (
    "context"
    "time"
    "opscore/backend/internal/shared/cache"
)

// Initialize cache
config := cache.CacheConfig{
    Type:       "memory",
    MaxEntries: 10000,
    CleanupInterval: 10 * time.Minute,
}
cacheInstance, err := cache.NewCache(config)
if err != nil {
    // handle error
}

// Use cache key builder for consistency
builder := cache.NewCacheKeyBuilder("opscore:")
ttlConfig := cache.DefaultTTLConfig()

// Store a document in cache
ctx := context.Background()
key := builder.DocumentKey(documentID)
err = cacheInstance.Set(ctx, key, document, ttlConfig.Document)

// Retrieve from cache
cached, err := cacheInstance.Get(ctx, key)
if err != nil {
    // handle error
}
if cached != nil {
    doc := cached.(Document) // Type assertion
    // use document
}
```

### Caching in Repository Layer

```go
type CachedDocumentRepository struct {
    repo  repository.DocumentRepository
    cache cache.Cache
    builder *cache.CacheKeyBuilder
    ttl   time.Duration
}

func NewCachedDocumentRepository(
    repo repository.DocumentRepository,
    cache cache.Cache,
) *CachedDocumentRepository {
    return &CachedDocumentRepository{
        repo:    repo,
        cache:   cache,
        builder: cache.NewCacheKeyBuilder("opscore:"),
        ttl:     1 * time.Hour,
    }
}

func (r *CachedDocumentRepository) FindByID(
    ctx context.Context,
    id value_object.DocumentID,
) (entity.Document, error) {
    // Try to get from cache
    key := r.builder.DocumentKey(id.String())
    cached, err := r.cache.Get(ctx, key)
    if err == nil && cached != nil {
        return cached.(entity.Document), nil
    }

    // Cache miss - fetch from database
    doc, err := r.repo.FindByID(ctx, id)
    if err != nil {
        return entity.Document{}, err
    }

    // Store in cache
    _ = r.cache.Set(ctx, key, doc, r.ttl)

    return doc, nil
}

func (r *CachedDocumentRepository) Update(
    ctx context.Context,
    document entity.Document,
) error {
    // Update in database
    err := r.repo.Update(ctx, document)
    if err != nil {
        return err
    }

    // Invalidate cache
    key := r.builder.DocumentKey(document.ID().String())
    _ = r.cache.Delete(ctx, key)

    // Also invalidate related caches
    repoKey := r.builder.RepositoryDocumentsKey(document.RepositoryID().String())
    _ = r.cache.Delete(ctx, repoKey)

    return nil
}
```

### Batch Caching

```go
// Cache multiple documents at once
func (r *CachedDocumentRepository) FindByRepositoryID(
    ctx context.Context,
    repoID value_object.RepositoryID,
) ([]entity.Document, error) {
    // Try cache first
    key := r.builder.RepositoryDocumentsKey(repoID.String())
    cached, err := r.cache.Get(ctx, key)
    if err == nil && cached != nil {
        return cached.([]entity.Document), nil
    }

    // Fetch from database
    docs, err := r.repo.FindByRepositoryID(ctx, repoID)
    if err != nil {
        return nil, err
    }

    // Cache the list
    _ = r.cache.Set(ctx, key, docs, r.ttl)

    // Also cache individual documents
    items := make(map[string]interface{})
    for _, doc := range docs {
        docKey := r.builder.DocumentKey(doc.ID().String())
        items[docKey] = doc
    }
    _ = r.cache.SetMultiple(ctx, items, r.ttl)

    return docs, nil
}
```

### Using Pagination

```go
import "opscore/backend/internal/shared/pagination"

// In your handler or use case
func (h *DocumentHandler) ListDocuments(c *gin.Context) {
    // Parse pagination parameters
    page := pagination.ParseInt(c.Query("page"), 1)
    pageSize := pagination.ParseInt(c.Query("pageSize"), 20)
    
    // Create page object
    pageObj := pagination.NewPageFromPageNumber(page, pageSize)
    
    // Fetch data with pagination
    documents, total, err := h.useCase.ListDocuments(ctx, pageObj)
    if err != nil {
        // handle error
    }
    
    // Create paginated result
    result := pagination.NewPageResult(documents, total, pageObj)
    
    c.JSON(200, result)
}

// With cursor-based pagination
func (h *DocumentHandler) ListDocumentsCursor(c *gin.Context) {
    cursor := c.Query("cursor")
    limit := pagination.ParseInt(c.Query("limit"), 20)
    
    // Decode cursor
    decoded, _ := pagination.DecodeCursor(cursor)
    
    // Fetch data
    documents, hasNext, err := h.useCase.ListDocumentsCursor(
        ctx, 
        decoded.ID, 
        decoded.Timestamp,
        limit,
    )
    if err != nil {
        // handle error
    }
    
    // Generate next cursor
    var nextCursor string
    if hasNext && len(documents) > 0 {
        lastDoc := documents[len(documents)-1]
        nextCursor, _ = pagination.EncodeCursor(
            lastDoc.ID,
            lastDoc.CreatedAt.Unix(),
        )
    }
    
    result := pagination.NewCursorPageResult(documents, hasNext, nextCursor, limit)
    c.JSON(200, result)
}
```

## Frontend Examples

### Basic Cache Usage

```typescript
import { cache, TTL, CacheKeys } from '@/utils/cache';

// Cache a document
const document = await fetchDocument(docId);
cache.set(CacheKeys.document(docId), document, {
    ttl: TTL.ONE_HOUR,
    storage: 'memory'
});

// Retrieve from cache
const cachedDoc = cache.get(CacheKeys.document(docId));
if (cachedDoc) {
    // Use cached document
} else {
    // Fetch fresh data
}
```

### React Hook with Cache

```typescript
import { useState, useEffect } from 'react';
import { cache, TTL, CacheKeys } from '@/utils/cache';

function useDocument(docId: string) {
    const [document, setDocument] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    useEffect(() => {
        async function loadDocument() {
            try {
                setLoading(true);

                // Try to get from cache first
                const cached = cache.get(CacheKeys.document(docId));
                if (cached) {
                    setDocument(cached);
                    setLoading(false);
                    return;
                }

                // Cache miss - fetch from API
                const response = await fetch(`/api/v1/documents/${docId}`);
                if (!response.ok) throw new Error('Failed to fetch');
                
                const data = await response.json();
                
                // Store in cache
                cache.set(CacheKeys.document(docId), data, {
                    ttl: TTL.ONE_HOUR,
                    storage: 'memory'
                });
                
                setDocument(data);
            } catch (err) {
                setError(err);
            } finally {
                setLoading(false);
            }
        }

        loadDocument();
    }, [docId]);

    const invalidateCache = () => {
        cache.delete(CacheKeys.document(docId));
    };

    return { document, loading, error, invalidateCache };
}

// Usage in component
function DocumentViewer({ docId }) {
    const { document, loading, error, invalidateCache } = useDocument(docId);

    if (loading) return <div>Loading...</div>;
    if (error) return <div>Error: {error.message}</div>;

    return (
        <div>
            <h1>{document.title}</h1>
            <button onClick={invalidateCache}>Refresh</button>
        </div>
    );
}
```

### getOrSet Pattern

```typescript
import { cache, TTL, CacheKeys } from '@/utils/cache';

async function getDocument(docId: string) {
    return await cache.getOrSet(
        CacheKeys.document(docId),
        async () => {
            // This function only runs on cache miss
            const response = await fetch(`/api/v1/documents/${docId}`);
            if (!response.ok) throw new Error('Failed to fetch');
            return await response.json();
        },
        { ttl: TTL.ONE_HOUR, storage: 'memory' }
    );
}
```

### Custom Cache Instance with Configuration

```typescript
import { Cache, TTL } from '@/utils/cache';

// Create a cache instance with custom cleanup interval
const shortLivedCache = new Cache({
    prefix: 'temp_',
    cleanupInterval: TTL.ONE_MINUTE // Clean up every minute
});

// Use for temporary data
shortLivedCache.set('temp_key', data, {
    ttl: TTL.FIVE_MINUTES,
    storage: 'memory'
});
```

### localStorage for Persistent Cache

```typescript
import { cache, TTL, CacheKeys } from '@/utils/cache';

// Cache user preferences persistently
function saveUserPreferences(prefs: UserPreferences) {
    cache.set('user:preferences', prefs, {
        ttl: TTL.ONE_DAY,
        storage: 'localStorage' // Survives page refresh
    });
}

function loadUserPreferences(): UserPreferences | null {
    return cache.get('user:preferences', 'localStorage');
}
```

### Cache Invalidation on Mutation

```typescript
import { cache, CacheKeys } from '@/utils/cache';

async function updateDocument(docId: string, updates: Partial<Document>) {
    // Make the update
    const response = await fetch(`/api/v1/documents/${docId}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(updates)
    });

    if (!response.ok) throw new Error('Update failed');

    // Invalidate all related caches
    cache.delete(CacheKeys.document(docId));
    cache.delete(CacheKeys.documentVersions(docId));
    cache.delete(CacheKeys.repositories()); // If document list changed

    return await response.json();
}
```

### Clearing Specific Cache Types

```typescript
import { cache } from '@/utils/cache';

// Clear all memory cache (lost on refresh anyway)
cache.clear('memory');

// Clear all localStorage cache (persistent)
cache.clear('localStorage');

// Clear both
cache.clear('all');

// Selective clearing - delete by pattern
function clearDocumentCaches() {
    // Since we can't pattern-match easily, delete known keys
    // This is why cache key builders are useful
    const docIds = getKnownDocumentIds();
    docIds.forEach(id => {
        cache.delete(CacheKeys.document(id));
        cache.delete(CacheKeys.documentVersions(id));
    });
}
```

## Best Practices

### 1. Always Use Key Builders

```typescript
// ❌ Bad - inconsistent keys
cache.set('doc-' + id, data);
cache.get('document-' + id); // Won't find it!

// ✅ Good - use key builders
cache.set(CacheKeys.document(id), data);
cache.get(CacheKeys.document(id));
```

### 2. Set Appropriate TTLs

```typescript
// ❌ Bad - too short for stable data
cache.set(key, document, { ttl: TTL.ONE_MINUTE }); // Document rarely changes

// ✅ Good - match TTL to data volatility
cache.set(CacheKeys.document(id), document, { ttl: TTL.ONE_HOUR });
cache.set(CacheKeys.viewStatistics(id), stats, { ttl: TTL.THIRTY_MINUTES });
```

### 3. Invalidate on Updates

```typescript
// ❌ Bad - cache becomes stale
await updateDocument(id, changes);
// Cache still has old data!

// ✅ Good - invalidate after updates
await updateDocument(id, changes);
cache.delete(CacheKeys.document(id));
```

### 4. Handle Cache Failures Gracefully

```typescript
// ✅ Good - don't let cache errors break functionality
async function getDocument(id: string) {
    try {
        const cached = cache.get(CacheKeys.document(id));
        if (cached) return cached;
    } catch (error) {
        console.warn('Cache error, falling back to API', error);
    }

    // Always fetch if cache fails or misses
    return await fetchFromApi(id);
}
```

### 5. Use Memory for Volatile, localStorage for Stable

```typescript
// Volatile data - memory cache
cache.set(CacheKeys.viewStatistics(id), stats, {
    ttl: TTL.THIRTY_MINUTES,
    storage: 'memory'
});

// Stable data - localStorage
cache.set(CacheKeys.document(id), document, {
    ttl: TTL.ONE_HOUR,
    storage: 'localStorage'
});
```

## Performance Monitoring

### Backend

```go
// Track cache hit/miss rates
type CacheMetrics struct {
    hits   int64
    misses int64
}

func (r *CachedRepository) FindByID(ctx context.Context, id string) (interface{}, error) {
    cached, err := r.cache.Get(ctx, key)
    if err == nil && cached != nil {
        atomic.AddInt64(&r.metrics.hits, 1)
        return cached, nil
    }
    
    atomic.AddInt64(&r.metrics.misses, 1)
    // ... fetch from database
}

func (r *CachedRepository) GetHitRate() float64 {
    hits := atomic.LoadInt64(&r.metrics.hits)
    misses := atomic.LoadInt64(&r.metrics.misses)
    total := hits + misses
    if total == 0 {
        return 0
    }
    return float64(hits) / float64(total)
}
```

### Frontend

```typescript
class CacheMetrics {
    private hits = 0;
    private misses = 0;

    recordHit() {
        this.hits++;
    }

    recordMiss() {
        this.misses++;
    }

    getHitRate(): number {
        const total = this.hits + this.misses;
        return total === 0 ? 0 : this.hits / total;
    }

    reset() {
        this.hits = 0;
        this.misses = 0;
    }
}

export const cacheMetrics = new CacheMetrics();

// Track in your cache wrapper
async function getCached<T>(key: string): Promise<T | null> {
    const cached = cache.get<T>(key);
    if (cached) {
        cacheMetrics.recordHit();
        return cached;
    }
    cacheMetrics.recordMiss();
    return null;
}

// Log metrics periodically
setInterval(() => {
    console.log('Cache hit rate:', cacheMetrics.getHitRate());
}, 60000);
```

## Troubleshooting

### Cache Not Working

1. Check that keys are consistent (use key builders)
2. Verify TTL is not too short
3. Check that cache is initialized
4. Look for cache clearing calls

### Stale Data

1. Verify cache invalidation on updates
2. Reduce TTL for frequently changing data
3. Check for race conditions between update and read

### Memory Issues

1. Monitor cache size
2. Set MaxEntries for memory cache
3. Use shorter TTLs
4. Implement cache eviction
5. Use localStorage for non-critical data

### Performance Not Improving

1. Check cache hit rate (should be > 80% for stable data)
2. Verify indexes are being used (backend)
3. Profile application to find bottlenecks
4. Consider adding more cache layers
