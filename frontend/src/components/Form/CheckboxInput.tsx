/**
 * Checkbox input component
 */

import type { BaseInputProps } from "../../types/ui";

export interface CheckboxInputProps extends Omit<BaseInputProps, "error"> {
  /** Whether the checkbox is checked */
  checked: boolean;
  /** Change handler */
  onChange: (checked: boolean) => void;
  /** Helper text displayed below the label */
  helperText?: string;
}

/**
 * A styled checkbox input component
 */
export function CheckboxInput({
  id,
  name,
  label,
  checked,
  onChange,
  helperText,
  required = false,
  disabled = false,
  className = "",
}: CheckboxInputProps) {
  const inputId = id || name;

  return (
    <div className={`flex items-start ${className}`}>
      <div className="flex items-center h-5">
        <input
          id={inputId}
          name={name}
          type="checkbox"
          checked={checked}
          onChange={(e) => onChange(e.target.checked)}
          disabled={disabled}
          required={required}
          className={`
            h-4 w-4 rounded
            border-gray-300 dark:border-gray-600
            text-blue-600
            focus:ring-2 focus:ring-blue-500 focus:ring-offset-0
            disabled:opacity-50 disabled:cursor-not-allowed
            transition-colors
          `}
        />
      </div>
      {(label || helperText) && (
        <div className="ml-3">
          {label && (
            <label
              htmlFor={inputId}
              className={`
                text-sm font-medium
                ${disabled ? "text-gray-400 dark:text-gray-500" : "text-gray-700 dark:text-gray-300"}
              `}
            >
              {label}
              {required && <span className="text-red-500 ml-1">*</span>}
            </label>
          )}
          {helperText && (
            <p className="text-sm text-gray-500 dark:text-gray-400">
              {helperText}
            </p>
          )}
        </div>
      )}
    </div>
  );
}
