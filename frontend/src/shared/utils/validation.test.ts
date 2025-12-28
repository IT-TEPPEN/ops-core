import { describe, it, expect } from "vitest";
import {
  required,
  minLength,
  maxLength,
  pattern,
  email,
  url,
  numberRange,
  compose,
  validateForm,
} from "./validation";

describe("validation utilities", () => {
  describe("required", () => {
    it("returns valid for non-empty string", () => {
      const result = required()("hello");
      expect(result.isValid).toBe(true);
      expect(result.error).toBeUndefined();
    });

    it("returns invalid for empty string", () => {
      const result = required()("");
      expect(result.isValid).toBe(false);
      expect(result.error).toBe("This field is required");
    });

    it("returns invalid for whitespace-only string", () => {
      const result = required()("   ");
      expect(result.isValid).toBe(false);
    });

    it("uses custom error message", () => {
      const result = required("Custom error")("");
      expect(result.error).toBe("Custom error");
    });
  });

  describe("minLength", () => {
    it("returns valid when length meets minimum", () => {
      const result = minLength(5)("hello");
      expect(result.isValid).toBe(true);
    });

    it("returns invalid when length is below minimum", () => {
      const result = minLength(5)("hi");
      expect(result.isValid).toBe(false);
      expect(result.error).toBe("Minimum 5 characters required");
    });
  });

  describe("maxLength", () => {
    it("returns valid when length is within maximum", () => {
      const result = maxLength(10)("hello");
      expect(result.isValid).toBe(true);
    });

    it("returns invalid when length exceeds maximum", () => {
      const result = maxLength(3)("hello");
      expect(result.isValid).toBe(false);
      expect(result.error).toBe("Maximum 3 characters allowed");
    });
  });

  describe("pattern", () => {
    it("returns valid when pattern matches", () => {
      const result = pattern(/^[a-z]+$/)("hello");
      expect(result.isValid).toBe(true);
    });

    it("returns invalid when pattern does not match", () => {
      const result = pattern(/^[a-z]+$/, "Only lowercase letters allowed")("Hello123");
      expect(result.isValid).toBe(false);
      expect(result.error).toBe("Only lowercase letters allowed");
    });
  });

  describe("email", () => {
    it("returns valid for valid email", () => {
      const result = email()("test@example.com");
      expect(result.isValid).toBe(true);
    });

    it("returns invalid for invalid email", () => {
      const result = email()("invalid-email");
      expect(result.isValid).toBe(false);
      expect(result.error).toBe("Invalid email address");
    });
  });

  describe("url", () => {
    it("returns valid for valid URL", () => {
      const result = url()("https://example.com");
      expect(result.isValid).toBe(true);
    });

    it("returns invalid for invalid URL", () => {
      const result = url()("not-a-url");
      expect(result.isValid).toBe(false);
      expect(result.error).toBe("Invalid URL");
    });
  });

  describe("numberRange", () => {
    it("returns valid when number is within range", () => {
      const result = numberRange(1, 10)(5);
      expect(result.isValid).toBe(true);
    });

    it("returns invalid when number is below minimum", () => {
      const result = numberRange(1, 10)(0);
      expect(result.isValid).toBe(false);
    });

    it("returns invalid when number is above maximum", () => {
      const result = numberRange(1, 10)(11);
      expect(result.isValid).toBe(false);
    });
  });

  describe("compose", () => {
    it("returns valid when all validators pass", () => {
      const validator = compose(required(), minLength(3), maxLength(10));
      const result = validator("hello");
      expect(result.isValid).toBe(true);
    });

    it("returns first error when a validator fails", () => {
      const validator = compose(required(), minLength(10));
      const result = validator("hi");
      expect(result.isValid).toBe(false);
      expect(result.error).toBe("Minimum 10 characters required");
    });
  });

  describe("validateForm", () => {
    it("validates all fields and returns errors", () => {
      const values = { name: "", email: "invalid" };
      const rules = {
        name: required(),
        email: email(),
      };

      const errors = validateForm(values, rules);
      expect(errors.name).toBe("This field is required");
      expect(errors.email).toBe("Invalid email address");
    });

    it("returns empty object when all fields are valid", () => {
      const values = { name: "John", email: "john@example.com" };
      const rules = {
        name: required(),
        email: email(),
      };

      const errors = validateForm(values, rules);
      expect(Object.keys(errors)).toHaveLength(0);
    });
  });
});
