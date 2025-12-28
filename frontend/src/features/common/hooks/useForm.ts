/**
 * Custom hook for form management
 */

import { useState, useCallback, useMemo } from "react";
import type { FormState } from "../types/ui";
import type { ValidationFn } from "../utils/validation";

export interface UseFormOptions<T> {
  /** Initial form values */
  initialValues: T;
  /** Validation rules */
  validationRules?: Partial<Record<keyof T, ValidationFn>>;
  /** Callback when form is submitted */
  onSubmit?: (values: T) => void | Promise<void>;
}

export interface UseFormResult<T> extends FormState<T> {
  /** Handle input change */
  handleChange: (name: keyof T, value: unknown) => void;
  /** Handle input blur */
  handleBlur: (name: keyof T) => void;
  /** Handle form submission */
  handleSubmit: (e?: React.FormEvent) => void;
  /** Reset form to initial values */
  reset: () => void;
  /** Set a specific field value */
  setValue: (name: keyof T, value: unknown) => void;
  /** Set a specific field error */
  setError: (name: keyof T, error: string) => void;
  /** Get props for an input field */
  getFieldProps: (name: keyof T) => {
    name: keyof T;
    value: unknown;
    onChange: (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement | HTMLTextAreaElement>) => void;
    onBlur: () => void;
  };
}

/**
 * Hook for managing form state and validation
 */
export function useForm<T extends Record<string, unknown>>(
  options: UseFormOptions<T>
): UseFormResult<T> {
  const { initialValues, validationRules = {}, onSubmit } = options;

  const [values, setValues] = useState<T>(initialValues);
  const [errors, setErrors] = useState<Partial<Record<keyof T, string>>>({});
  const [touched, setTouched] = useState<Partial<Record<keyof T, boolean>>>({});
  const [isSubmitting, setIsSubmitting] = useState(false);

  const validateField = useCallback(
    (name: keyof T, value: unknown): string | undefined => {
      const rules = validationRules as Record<string, ValidationFn | undefined>;
      const validator = rules[name as string];
      if (validator) {
        const result = validator(String(value ?? ""));
        return result.error;
      }
      return undefined;
    },
    [validationRules]
  );

  const validateAll = useCallback((): Partial<Record<keyof T, string>> => {
    const newErrors: Partial<Record<keyof T, string>> = {};
    for (const key in validationRules) {
      const error = validateField(key as keyof T, values[key as keyof T]);
      if (error) {
        newErrors[key as keyof T] = error;
      }
    }
    return newErrors;
  }, [values, validationRules, validateField]);

  const isValid = useMemo(() => {
    return Object.keys(validateAll()).length === 0;
  }, [validateAll]);

  const handleChange = useCallback(
    (name: keyof T, value: unknown) => {
      setValues((prev) => ({ ...prev, [name]: value }));
      
      // Clear error when user starts typing
      if (errors[name]) {
        setErrors((prev) => ({ ...prev, [name]: undefined }));
      }
    },
    [errors]
  );

  const handleBlur = useCallback(
    (name: keyof T) => {
      setTouched((prev) => ({ ...prev, [name]: true }));
      
      const error = validateField(name, values[name]);
      if (error) {
        setErrors((prev) => ({ ...prev, [name]: error }));
      }
    },
    [values, validateField]
  );

  const handleSubmit = useCallback(
    async (e?: React.FormEvent) => {
      e?.preventDefault();

      // Mark all fields as touched
      const allTouched = Object.keys(values).reduce(
        (acc, key) => ({ ...acc, [key]: true }),
        {} as Partial<Record<keyof T, boolean>>
      );
      setTouched(allTouched);

      // Validate all fields
      const validationErrors = validateAll();
      setErrors(validationErrors);

      if (Object.keys(validationErrors).length > 0) {
        return;
      }

      if (onSubmit) {
        setIsSubmitting(true);
        try {
          await onSubmit(values);
        } finally {
          setIsSubmitting(false);
        }
      }
    },
    [values, validateAll, onSubmit]
  );

  const reset = useCallback(() => {
    setValues(initialValues);
    setErrors({});
    setTouched({});
    setIsSubmitting(false);
  }, [initialValues]);

  const setValue = useCallback((name: keyof T, value: unknown) => {
    setValues((prev) => ({ ...prev, [name]: value }));
  }, []);

  const setError = useCallback((name: keyof T, error: string) => {
    setErrors((prev) => ({ ...prev, [name]: error }));
  }, []);

  const getFieldProps = useCallback(
    (name: keyof T) => ({
      name,
      value: values[name],
      onChange: (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement | HTMLTextAreaElement>) => {
        const target = e.target;
        const value = target.type === "checkbox" ? (target as HTMLInputElement).checked : target.value;
        handleChange(name, value);
      },
      onBlur: () => handleBlur(name),
    }),
    [values, handleChange, handleBlur]
  );

  return {
    values,
    errors,
    touched,
    isSubmitting,
    isValid,
    handleChange,
    handleBlur,
    handleSubmit,
    reset,
    setValue,
    setError,
    getFieldProps,
  };
}
