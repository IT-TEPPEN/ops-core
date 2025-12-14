import { describe, it, expect, beforeEach, afterEach, vi } from "vitest";
import { Cache, TTL, CacheKeys } from "./cache";

describe("Cache", () => {
  let cache: Cache;

  beforeEach(() => {
    cache = new Cache("test_");
    localStorage.clear();
  });

  afterEach(() => {
    cache.stopCleanup();
    localStorage.clear();
  });

  describe("Memory cache", () => {
    it("should set and get a value", () => {
      cache.set("key1", "value1", { storage: "memory" });
      const value = cache.get<string>("key1", "memory");
      expect(value).toBe("value1");
    });

    it("should return null for non-existent key", () => {
      const value = cache.get<string>("nonexistent", "memory");
      expect(value).toBeNull();
    });

    it("should delete a value", () => {
      cache.set("key1", "value1", { storage: "memory" });
      cache.delete("key1", "memory");
      const value = cache.get<string>("key1", "memory");
      expect(value).toBeNull();
    });

    it("should check if key exists", () => {
      cache.set("key1", "value1", { storage: "memory" });
      expect(cache.has("key1", "memory")).toBe(true);
      expect(cache.has("nonexistent", "memory")).toBe(false);
    });

    it("should handle TTL expiration", async () => {
      vi.useFakeTimers();
      cache.set("key1", "value1", { storage: "memory", ttl: 1000 });

      // Value should exist immediately
      expect(cache.get<string>("key1", "memory")).toBe("value1");

      // Advance time past TTL
      vi.advanceTimersByTime(1001);

      // Value should be expired
      expect(cache.get<string>("key1", "memory")).toBeNull();

      vi.useRealTimers();
    });

    it("should handle zero TTL (no expiration)", () => {
      cache.set("key1", "value1", { storage: "memory", ttl: 0 });
      const value = cache.get<string>("key1", "memory");
      expect(value).toBe("value1");
    });

    it("should clear all memory cache", () => {
      cache.set("key1", "value1", { storage: "memory" });
      cache.set("key2", "value2", { storage: "memory" });
      cache.clear("memory");
      expect(cache.get<string>("key1", "memory")).toBeNull();
      expect(cache.get<string>("key2", "memory")).toBeNull();
    });
  });

  describe("localStorage cache", () => {
    it("should set and get a value", () => {
      cache.set("key1", "value1", { storage: "localStorage" });
      const value = cache.get<string>("key1", "localStorage");
      expect(value).toBe("value1");
    });

    it("should return null for non-existent key", () => {
      const value = cache.get<string>("nonexistent", "localStorage");
      expect(value).toBeNull();
    });

    it("should delete a value", () => {
      cache.set("key1", "value1", { storage: "localStorage" });
      cache.delete("key1", "localStorage");
      const value = cache.get<string>("key1", "localStorage");
      expect(value).toBeNull();
    });

    it("should check if key exists", () => {
      cache.set("key1", "value1", { storage: "localStorage" });
      expect(cache.has("key1", "localStorage")).toBe(true);
      expect(cache.has("nonexistent", "localStorage")).toBe(false);
    });

    it("should handle TTL expiration", () => {
      vi.useFakeTimers();
      cache.set("key1", "value1", { storage: "localStorage", ttl: 1000 });

      // Value should exist immediately
      expect(cache.get<string>("key1", "localStorage")).toBe("value1");

      // Advance time past TTL
      vi.advanceTimersByTime(1001);

      // Value should be expired
      expect(cache.get<string>("key1", "localStorage")).toBeNull();

      vi.useRealTimers();
    });

    it("should clear all localStorage cache", () => {
      cache.set("key1", "value1", { storage: "localStorage" });
      cache.set("key2", "value2", { storage: "localStorage" });
      cache.clear("localStorage");
      expect(cache.get<string>("key1", "localStorage")).toBeNull();
      expect(cache.get<string>("key2", "localStorage")).toBeNull();
    });
  });

  describe("Complex values", () => {
    it("should cache objects", () => {
      const obj = { name: "test", age: 30 };
      cache.set("obj", obj, { storage: "memory" });
      const cached = cache.get<typeof obj>("obj", "memory");
      expect(cached).toEqual(obj);
    });

    it("should cache arrays", () => {
      const arr = [1, 2, 3, 4, 5];
      cache.set("arr", arr, { storage: "memory" });
      const cached = cache.get<typeof arr>("arr", "memory");
      expect(cached).toEqual(arr);
    });

    it("should cache objects in localStorage", () => {
      const obj = { name: "test", age: 30 };
      cache.set("obj", obj, { storage: "localStorage" });
      const cached = cache.get<typeof obj>("obj", "localStorage");
      expect(cached).toEqual(obj);
    });
  });

  describe("getOrSet", () => {
    it("should return cached value if exists", async () => {
      cache.set("key1", "cached", { storage: "memory" });
      const factory = vi.fn(() => "new");
      const value = await cache.getOrSet("key1", factory, { storage: "memory" });
      expect(value).toBe("cached");
      expect(factory).not.toHaveBeenCalled();
    });

    it("should call factory and cache result if not exists", async () => {
      const factory = vi.fn(() => "new");
      const value = await cache.getOrSet("key1", factory, { storage: "memory" });
      expect(value).toBe("new");
      expect(factory).toHaveBeenCalledOnce();
      expect(cache.get<string>("key1", "memory")).toBe("new");
    });

    it("should work with async factory", async () => {
      const factory = vi.fn(async () => {
        return "async value";
      });
      const value = await cache.getOrSet("key1", factory, { storage: "memory" });
      expect(value).toBe("async value");
      expect(cache.get<string>("key1", "memory")).toBe("async value");
    });
  });

  describe("Clear all", () => {
    it("should clear both memory and localStorage", () => {
      cache.set("key1", "value1", { storage: "memory" });
      cache.set("key2", "value2", { storage: "localStorage" });
      cache.clear("all");
      expect(cache.get<string>("key1", "memory")).toBeNull();
      expect(cache.get<string>("key2", "localStorage")).toBeNull();
    });
  });

  describe("Prefix isolation", () => {
    it("should only clear keys with matching prefix", () => {
      const cache1 = new Cache("prefix1_");
      const cache2 = new Cache("prefix2_");

      cache1.set("key", "value1", { storage: "localStorage" });
      cache2.set("key", "value2", { storage: "localStorage" });

      cache1.clear("localStorage");

      expect(cache1.get<string>("key", "localStorage")).toBeNull();
      expect(cache2.get<string>("key", "localStorage")).toBe("value2");

      cache1.stopCleanup();
      cache2.stopCleanup();
    });
  });
});

describe("TTL constants", () => {
  it("should have correct values", () => {
    expect(TTL.ONE_MINUTE).toBe(60 * 1000);
    expect(TTL.FIVE_MINUTES).toBe(5 * 60 * 1000);
    expect(TTL.FIFTEEN_MINUTES).toBe(15 * 60 * 1000);
    expect(TTL.THIRTY_MINUTES).toBe(30 * 60 * 1000);
    expect(TTL.ONE_HOUR).toBe(60 * 60 * 1000);
    expect(TTL.ONE_DAY).toBe(24 * 60 * 60 * 1000);
  });
});

describe("CacheKeys", () => {
  it("should generate correct keys", () => {
    expect(CacheKeys.document("123")).toBe("document:123");
    expect(CacheKeys.documentVersion("123", "v1")).toBe("document_version:123:v1");
    expect(CacheKeys.documentVersions("123")).toBe("document_versions:123");
    expect(CacheKeys.user("456")).toBe("user:456");
    expect(CacheKeys.repository("789")).toBe("repository:789");
    expect(CacheKeys.repositories()).toBe("repositories:all");
    expect(CacheKeys.executionRecord("abc")).toBe("execution_record:abc");
    expect(CacheKeys.viewStatistics("doc1")).toBe("view_statistics:doc1");
  });
});
