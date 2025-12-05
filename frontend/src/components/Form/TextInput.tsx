/**
 * Text input component
 */

import type { BaseInputProps } from "../../types/ui";

export interface TextInputProps extends BaseInputProps {
  /** Input value */
  value: string;
  /** Change handler */
  onChange: (value: string) => void;
  /** Input type */
  type?: "text" | "email" | "password" | "url" | "tel" | "search";
  /** Placeholder text */
  placeholder?: string;
  /** Blur handler */
  onBlur?: () => void;
  /** Maximum length */
  maxLength?: number;
  /** Auto-complete attribute */
  autoComplete?: string;
}

/**
 * A styled text input component
 */
export function TextInput({
  id,
  name,
  label,
  value,
  onChange,
  type = "text",
  placeholder,
  error,
  required = false,
  disabled = false,
  className = "",
  onBlur,
  maxLength,
  autoComplete,
}: TextInputProps) {
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
        type={type}
        value={value}
        onChange={(e) => onChange(e.target.value)}
        onBlur={onBlur}
        placeholder={placeholder}
        disabled={disabled}
        required={required}
        maxLength={maxLength}
        autoComplete={autoComplete}
        className={`
          w-full px-3 py-2 rounded-md shadow-sm
          border ${error ? "border-red-500" : "border-gray-300 dark:border-gray-600"}
          bg-white dark:bg-gray-700
          text-gray-900 dark:text-gray-100
          placeholder-gray-400 dark:placeholder-gray-500
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
