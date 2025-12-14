/**
 * Client-side cache utility for storing and retrieving data.
 * Supports both memory and localStorage-based caching with TTL.
 */

export interface CacheEntry<T> {
  value: T;
  expiration: number; // Timestamp in milliseconds
}

export interface CacheOptions {
  ttl?: number; // Time to live in milliseconds
  storage?: "memory" | "localStorage";
}

/**
 * Cache class for client-side data caching.
 */
export class Cache {
  private memoryCache: Map<string, CacheEntry<any>> = new Map();
  private readonly prefix: string;
  private cleanupInterval: number | null = null;

  constructor(prefix: string = "opscore_cache_") {
    this.prefix = prefix;
    this.startCleanup();
  }

  /**
   * Get a value from the cache.
   * @param key Cache key
   * @param storage Storage type (memory or localStorage)
   * @returns Cached value or null if not found or expired
   */
  get<T>(key: string, storage: "memory" | "localStorage" = "memory"): T | null {
    const fullKey = this.prefix + key;

    if (storage === "memory") {
      const entry = this.memoryCache.get(fullKey);
      if (!entry) return null;

      if (this.isExpired(entry)) {
        this.memoryCache.delete(fullKey);
        return null;
      }

      return entry.value as T;
    } else {
      try {
        const item = localStorage.getItem(fullKey);
        if (!item) return null;

        const entry: CacheEntry<T> = JSON.parse(item);
        if (this.isExpired(entry)) {
          localStorage.removeItem(fullKey);
          return null;
        }

        return entry.value;
      } catch (error) {
        console.error("Cache get error:", error);
        return null;
      }
    }
  }

  /**
   * Set a value in the cache.
   * @param key Cache key
   * @param value Value to cache
   * @param options Cache options (ttl, storage)
   */
  set<T>(
    key: string,
    value: T,
    options: CacheOptions = {}
  ): void {
    const { ttl = 3600000, storage = "memory" } = options; // Default 1 hour TTL
    const fullKey = this.prefix + key;
    const expiration = ttl > 0 ? Date.now() + ttl : 0;

    const entry: CacheEntry<T> = {
      value,
      expiration,
    };

    if (storage === "memory") {
      this.memoryCache.set(fullKey, entry);
    } else {
      try {
        localStorage.setItem(fullKey, JSON.stringify(entry));
      } catch (error) {
        console.error("Cache set error:", error);
        // Fall back to memory if localStorage is full or unavailable
        this.memoryCache.set(fullKey, entry);
      }
    }
  }

  /**
   * Delete a value from the cache.
   * @param key Cache key
   * @param storage Storage type
   */
  delete(key: string, storage: "memory" | "localStorage" = "memory"): void {
    const fullKey = this.prefix + key;

    if (storage === "memory") {
      this.memoryCache.delete(fullKey);
    } else {
      try {
        localStorage.removeItem(fullKey);
      } catch (error) {
        console.error("Cache delete error:", error);
      }
    }
  }

  /**
   * Clear all cache entries.
   * @param storage Storage type to clear (or "all" for both)
   */
  clear(storage: "memory" | "localStorage" | "all" = "all"): void {
    if (storage === "memory" || storage === "all") {
      this.memoryCache.clear();
    }

    if (storage === "localStorage" || storage === "all") {
      try {
        // Remove all keys with our prefix
        const keysToRemove: string[] = [];
        for (let i = 0; i < localStorage.length; i++) {
          const key = localStorage.key(i);
          if (key && key.startsWith(this.prefix)) {
            keysToRemove.push(key);
          }
        }
        keysToRemove.forEach((key) => localStorage.removeItem(key));
      } catch (error) {
        console.error("Cache clear error:", error);
      }
    }
  }

  /**
   * Check if a key exists in the cache.
   * @param key Cache key
   * @param storage Storage type
   */
  has(key: string, storage: "memory" | "localStorage" = "memory"): boolean {
    const fullKey = this.prefix + key;

    if (storage === "memory") {
      const entry = this.memoryCache.get(fullKey);
      return entry !== undefined && !this.isExpired(entry);
    } else {
      try {
        const item = localStorage.getItem(fullKey);
        if (!item) return false;

        const entry: CacheEntry<any> = JSON.parse(item);
        return !this.isExpired(entry);
      } catch (error) {
        return false;
      }
    }
  }

  /**
   * Get or set a value with a factory function.
   * If the value exists and is not expired, return it.
   * Otherwise, call the factory function and cache the result.
   * @param key Cache key
   * @param factory Function to generate the value if not cached
   * @param options Cache options
   */
  async getOrSet<T>(
    key: string,
    factory: () => Promise<T> | T,
    options: CacheOptions = {}
  ): Promise<T> {
    const cached = this.get<T>(key, options.storage);
    if (cached !== null) {
      return cached;
    }

    const value = await factory();
    this.set(key, value, options);
    return value;
  }

  /**
   * Check if a cache entry is expired.
   */
  private isExpired(entry: CacheEntry<any>): boolean {
    if (entry.expiration === 0) return false; // No expiration
    return Date.now() > entry.expiration;
  }

  /**
   * Start background cleanup of expired entries.
   */
  private startCleanup(): void {
    // Clean up expired entries every 5 minutes
    this.cleanupInterval = window.setInterval(() => {
      this.cleanupExpired();
    }, 5 * 60 * 1000);
  }

  /**
   * Stop background cleanup.
   */
  stopCleanup(): void {
    if (this.cleanupInterval !== null) {
      clearInterval(this.cleanupInterval);
      this.cleanupInterval = null;
    }
  }

  /**
   * Remove all expired entries from the cache.
   */
  private cleanupExpired(): void {
    // Clean memory cache
    for (const [key, entry] of this.memoryCache.entries()) {
      if (this.isExpired(entry)) {
        this.memoryCache.delete(key);
      }
    }

    // Clean localStorage
    try {
      const keysToRemove: string[] = [];
      for (let i = 0; i < localStorage.length; i++) {
        const key = localStorage.key(i);
        if (key && key.startsWith(this.prefix)) {
          const item = localStorage.getItem(key);
          if (item) {
            const entry: CacheEntry<any> = JSON.parse(item);
            if (this.isExpired(entry)) {
              keysToRemove.push(key);
            }
          }
        }
      }
      keysToRemove.forEach((key) => localStorage.removeItem(key));
    } catch (error) {
      console.error("Cache cleanup error:", error);
    }
  }
}

// Default cache instance
export const cache = new Cache();

// TTL constants (in milliseconds)
export const TTL = {
  ONE_MINUTE: 60 * 1000,
  FIVE_MINUTES: 5 * 60 * 1000,
  FIFTEEN_MINUTES: 15 * 60 * 1000,
  THIRTY_MINUTES: 30 * 60 * 1000,
  ONE_HOUR: 60 * 60 * 1000,
  ONE_DAY: 24 * 60 * 60 * 1000,
};

// Predefined cache key builders
export const CacheKeys = {
  document: (id: string) => `document:${id}`,
  documentVersion: (docId: string, version: string) => `document_version:${docId}:${version}`,
  documentVersions: (docId: string) => `document_versions:${docId}`,
  user: (id: string) => `user:${id}`,
  repository: (id: string) => `repository:${id}`,
  repositories: () => `repositories:all`,
  executionRecord: (id: string) => `execution_record:${id}`,
  viewStatistics: (docId: string) => `view_statistics:${docId}`,
};
