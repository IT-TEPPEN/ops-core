import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";
import {
  formatDateISO,
  formatDateLocale,
  formatDateTime,
  formatRelativeTime,
  isToday,
  isPast,
  isFuture,
  addDays,
  startOfDay,
  endOfDay,
  parseDate,
} from "./date";

describe("date utilities", () => {
  describe("formatDateISO", () => {
    it("formats Date object to ISO string", () => {
      const date = new Date("2023-06-15T12:00:00Z");
      expect(formatDateISO(date)).toBe("2023-06-15");
    });

    it("formats date string to ISO string", () => {
      expect(formatDateISO("2023-06-15T12:00:00Z")).toBe("2023-06-15");
    });
  });

  describe("formatDateLocale", () => {
    it("formats date with default locale", () => {
      const date = new Date("2023-06-15");
      const result = formatDateLocale(date, "en-US");
      // Result format depends on locale
      expect(result).toContain("2023");
      expect(result).toContain("15");
    });
  });

  describe("formatDateTime", () => {
    it("formats date with time", () => {
      const date = new Date("2023-06-15T14:30:00");
      const result = formatDateTime(date, "en-US");
      expect(result).toBeTruthy();
    });
  });

  describe("formatRelativeTime", () => {
    beforeEach(() => {
      vi.useFakeTimers();
      vi.setSystemTime(new Date("2023-06-15T12:00:00"));
    });

    afterEach(() => {
      vi.useRealTimers();
    });

    it("returns 'just now' for recent times", () => {
      const date = new Date("2023-06-15T11:59:30");
      expect(formatRelativeTime(date)).toBe("just now");
    });

    it("returns minutes ago for recent times", () => {
      const date = new Date("2023-06-15T11:55:00");
      expect(formatRelativeTime(date)).toBe("5 minutes ago");
    });

    it("returns hours ago", () => {
      const date = new Date("2023-06-15T09:00:00");
      expect(formatRelativeTime(date)).toBe("3 hours ago");
    });

    it("returns days ago", () => {
      const date = new Date("2023-06-10T12:00:00");
      expect(formatRelativeTime(date)).toBe("5 days ago");
    });
  });

  describe("isToday", () => {
    beforeEach(() => {
      vi.useFakeTimers();
      vi.setSystemTime(new Date("2023-06-15T12:00:00"));
    });

    afterEach(() => {
      vi.useRealTimers();
    });

    it("returns true for today's date", () => {
      expect(isToday(new Date("2023-06-15T08:00:00"))).toBe(true);
    });

    it("returns false for yesterday", () => {
      expect(isToday(new Date("2023-06-14T12:00:00"))).toBe(false);
    });
  });

  describe("isPast", () => {
    it("returns true for past date", () => {
      const pastDate = new Date(Date.now() - 1000);
      expect(isPast(pastDate)).toBe(true);
    });

    it("returns false for future date", () => {
      const futureDate = new Date(Date.now() + 1000000);
      expect(isPast(futureDate)).toBe(false);
    });
  });

  describe("isFuture", () => {
    it("returns true for future date", () => {
      const futureDate = new Date(Date.now() + 1000000);
      expect(isFuture(futureDate)).toBe(true);
    });

    it("returns false for past date", () => {
      const pastDate = new Date(Date.now() - 1000);
      expect(isFuture(pastDate)).toBe(false);
    });
  });

  describe("addDays", () => {
    it("adds days to date", () => {
      const date = new Date("2023-06-15");
      const result = addDays(date, 5);
      expect(result.getDate()).toBe(20);
    });

    it("handles negative days", () => {
      const date = new Date("2023-06-15");
      const result = addDays(date, -5);
      expect(result.getDate()).toBe(10);
    });
  });

  describe("startOfDay", () => {
    it("returns start of day", () => {
      const date = new Date("2023-06-15T14:30:45");
      const result = startOfDay(date);
      expect(result.getHours()).toBe(0);
      expect(result.getMinutes()).toBe(0);
      expect(result.getSeconds()).toBe(0);
    });
  });

  describe("endOfDay", () => {
    it("returns end of day", () => {
      const date = new Date("2023-06-15T14:30:45");
      const result = endOfDay(date);
      expect(result.getHours()).toBe(23);
      expect(result.getMinutes()).toBe(59);
      expect(result.getSeconds()).toBe(59);
    });
  });

  describe("parseDate", () => {
    it("parses valid date string", () => {
      const result = parseDate("2023-06-15");
      expect(result).toBeInstanceOf(Date);
      expect(result?.getFullYear()).toBe(2023);
    });

    it("returns null for invalid date", () => {
      expect(parseDate("invalid")).toBeNull();
    });
  });
});
