/**
 * Validation utility functions
 */

/** Validation result */
export interface ValidationResult {
  isValid: boolean;
  error?: string;
}

/** Validation function type */
export type ValidationFn<T = string> = (value: T) => ValidationResult;

/**
 * Check if a value is required and not empty
 */
export const required = (message = "This field is required"): ValidationFn => {
  return (value: string): ValidationResult => {
    const isValid = value !== undefined && value !== null && value.trim() !== "";
    return {
      isValid,
      error: isValid ? undefined : message,
    };
  };
};

/**
 * Check minimum length
 */
export const minLength = (min: number, message?: string): ValidationFn => {
  return (value: string): ValidationResult => {
    const isValid = value.length >= min;
    return {
      isValid,
      error: isValid ? undefined : message || `Minimum ${min} characters required`,
    };
  };
};

/**
 * Check maximum length
 */
export const maxLength = (max: number, message?: string): ValidationFn => {
  return (value: string): ValidationResult => {
    const isValid = value.length <= max;
    return {
      isValid,
      error: isValid ? undefined : message || `Maximum ${max} characters allowed`,
    };
  };
};

/**
 * Check if value matches a pattern
 */
export const pattern = (regex: RegExp, message = "Invalid format"): ValidationFn => {
  return (value: string): ValidationResult => {
    const isValid = regex.test(value);
    return {
      isValid,
      error: isValid ? undefined : message,
    };
  };
};

/**
 * Check if value is a valid email
 */
export const email = (message = "Invalid email address"): ValidationFn => {
  const emailPattern = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  return pattern(emailPattern, message);
};

/**
 * Check if value is a valid URL
 */
export const url = (message = "Invalid URL"): ValidationFn => {
  return (value: string): ValidationResult => {
    try {
      new URL(value);
      return { isValid: true };
    } catch {
      return { isValid: false, error: message };
    }
  };
};

/**
 * Check if value is a number within range
 */
export const numberRange = (
  min?: number,
  max?: number,
  message?: string
): ValidationFn<number> => {
  return (value: number): ValidationResult => {
    let isValid = true;
    if (min !== undefined && value < min) isValid = false;
    if (max !== undefined && value > max) isValid = false;
    return {
      isValid,
      error: isValid ? undefined : message || `Value must be between ${min} and ${max}`,
    };
  };
};

/**
 * Combine multiple validators
 */
export const compose = (...validators: ValidationFn[]): ValidationFn => {
  return (value: string): ValidationResult => {
    for (const validator of validators) {
      const result = validator(value);
      if (!result.isValid) {
        return result;
      }
    }
    return { isValid: true };
  };
};

/**
 * Validate a form object
 */
export const validateForm = <T extends Record<string, unknown>>(
  values: T,
  rules: Partial<Record<keyof T, ValidationFn>>
): Partial<Record<keyof T, string>> => {
  const errors: Partial<Record<keyof T, string>> = {};

  for (const key in rules) {
    const validator = rules[key];
    if (validator) {
      const value = values[key];
      const result = validator(String(value ?? ""));
      if (!result.isValid && result.error) {
        errors[key] = result.error;
      }
    }
  }

  return errors;
};
