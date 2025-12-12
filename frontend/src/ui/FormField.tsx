/**
 * Form field wrapper component
 */

import type { ReactNode } from "react";

export interface FormFieldProps {
  /** Field label */
  label?: string;
  /** Field name/id */
  name: string;
  /** Error message */
  error?: string;
  /** Whether the field is required */
  required?: boolean;
  /** Helper text */
  helperText?: string;
  /** Field content */
  children: ReactNode;
  /** Additional class name */
  className?: string;
}

/**
 * A wrapper component for form fields with label, error, and helper text
 */
export function FormField({
  label,
  name,
  error,
  required = false,
  helperText,
  children,
  className = "",
}: FormFieldProps) {
  return (
    <div className={`space-y-1 ${className}`}>
      {label && (
        <label
          htmlFor={name}
          className="block text-sm font-medium text-gray-700 dark:text-gray-300"
        >
          {label}
          {required && <span className="text-red-500 ml-1">*</span>}
        </label>
      )}
      {children}
      {helperText && !error && (
        <p className="text-sm text-gray-500 dark:text-gray-400">{helperText}</p>
      )}
      {error && (
        <p id={`${name}-error`} className="text-sm text-red-500" role="alert">
          {error}
        </p>
      )}
    </div>
  );
}
