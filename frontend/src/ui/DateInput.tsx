/**
 * Date input component
 */

import type { BaseInputProps } from "../types/ui";

export interface DateInputProps extends BaseInputProps {
  /** Input value (ISO date string) */
  value: string;
  /** Change handler */
  onChange: (value: string) => void;
  /** Blur handler */
  onBlur?: () => void;
  /** Minimum date */
  min?: string;
  /** Maximum date */
  max?: string;
}

/**
 * A styled date input component
 */
export function DateInput({
  id,
  name,
  label,
  value,
  onChange,
  error,
  required = false,
  disabled = false,
  className = "",
  onBlur,
  min,
  max,
}: DateInputProps) {
  const inputId = id || name;

  return (
    <div className={className}>
      {label && (
        <label
          htmlFor={inputId}
          className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
        >
          {label}
          {required && <span className="text-red-500 ml-1">*</span>}
        </label>
      )}
      <input
        id={inputId}
        name={name}
        type="date"
        value={value}
        onChange={(e) => onChange(e.target.value)}
        onBlur={onBlur}
        disabled={disabled}
        required={required}
        min={min}
        max={max}
        className={`
          w-full px-3 py-2 rounded-md shadow-sm
          border ${error ? "border-red-500" : "border-gray-300 dark:border-gray-600"}
          bg-white dark:bg-gray-700
          text-gray-900 dark:text-gray-100
          focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500
          disabled:bg-gray-100 dark:disabled:bg-gray-800 disabled:cursor-not-allowed
          transition-colors
        `}
        aria-invalid={!!error}
        aria-describedby={error ? `${inputId}-error` : undefined}
      />
      {error && (
        <p id={`${inputId}-error`} className="mt-1 text-sm text-red-500">
          {error}
        </p>
      )}
    </div>
  );
}
