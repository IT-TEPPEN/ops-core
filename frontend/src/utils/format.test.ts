import { describe, it, expect } from "vitest";
import {
  truncate,
  capitalize,
  toTitleCase,
  toKebabCase,
  toCamelCase,
  formatNumber,
  formatBytes,
  formatPercentage,
  pluralize,
  joinWithAnd,
  toSlug,
  getInitials,
} from "./format";

describe("format utilities", () => {
  describe("truncate", () => {
    it("returns original string if shorter than maxLength", () => {
      expect(truncate("hello", 10)).toBe("hello");
    });

    it("truncates string with suffix", () => {
      expect(truncate("hello world", 8)).toBe("hello...");
    });

    it("uses custom suffix", () => {
      expect(truncate("hello world", 8, "…")).toBe("hello w…");
    });
  });

  describe("capitalize", () => {
    it("capitalizes first letter", () => {
      expect(capitalize("hello")).toBe("Hello");
    });

    it("returns empty string for empty input", () => {
      expect(capitalize("")).toBe("");
    });
  });

  describe("toTitleCase", () => {
    it("converts string to title case", () => {
      expect(toTitleCase("hello world")).toBe("Hello World");
    });
  });

  describe("toKebabCase", () => {
    it("converts camelCase to kebab-case", () => {
      expect(toKebabCase("camelCase")).toBe("camel-case");
      expect(toKebabCase("PascalCase")).toBe("pascal-case");
    });
  });

  describe("toCamelCase", () => {
    it("converts kebab-case to camelCase", () => {
      expect(toCamelCase("kebab-case")).toBe("kebabCase");
    });
  });

  describe("formatNumber", () => {
    it("formats number with thousands separators", () => {
      // Note: exact format depends on locale
      const result = formatNumber(1000000, "en-US");
      expect(result).toBe("1,000,000");
    });
  });

  describe("formatBytes", () => {
    it("formats 0 bytes", () => {
      expect(formatBytes(0)).toBe("0 Bytes");
    });

    it("formats bytes in KB", () => {
      expect(formatBytes(1024)).toBe("1 KB");
    });

    it("formats bytes in MB", () => {
      expect(formatBytes(1024 * 1024)).toBe("1 MB");
    });

    it("respects decimals parameter", () => {
      expect(formatBytes(1536, 1)).toBe("1.5 KB");
    });
  });

  describe("formatPercentage", () => {
    it("formats decimal as percentage", () => {
      expect(formatPercentage(0.75)).toBe("75.0%");
    });

    it("respects decimals parameter", () => {
      expect(formatPercentage(0.756, 2)).toBe("75.60%");
    });
  });

  describe("pluralize", () => {
    it("returns singular for count of 1", () => {
      expect(pluralize(1, "item")).toBe("item");
    });

    it("returns plural for count other than 1", () => {
      expect(pluralize(0, "item")).toBe("items");
      expect(pluralize(2, "item")).toBe("items");
    });

    it("uses custom plural form", () => {
      expect(pluralize(2, "person", "people")).toBe("people");
    });
  });

  describe("joinWithAnd", () => {
    it("returns empty string for empty array", () => {
      expect(joinWithAnd([])).toBe("");
    });

    it("returns single item for array of one", () => {
      expect(joinWithAnd(["apple"])).toBe("apple");
    });

    it("joins two items with and", () => {
      expect(joinWithAnd(["apple", "banana"])).toBe("apple and banana");
    });

    it("joins multiple items with commas and and", () => {
      expect(joinWithAnd(["apple", "banana", "cherry"])).toBe("apple, banana, and cherry");
    });
  });

  describe("toSlug", () => {
    it("converts string to URL-safe slug", () => {
      expect(toSlug("Hello World!")).toBe("hello-world");
    });

    it("handles multiple spaces and special characters", () => {
      expect(toSlug("  Hello   World  ")).toBe("hello-world");
    });
  });

  describe("getInitials", () => {
    it("extracts initials from name", () => {
      expect(getInitials("John Doe")).toBe("JD");
    });

    it("respects maxLength", () => {
      expect(getInitials("John Robert Doe", 2)).toBe("JR");
    });

    it("handles single name", () => {
      expect(getInitials("John")).toBe("J");
    });
  });
});
