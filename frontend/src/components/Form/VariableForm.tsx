import { useState } from "react";
import { VariableDefinition } from "../../types/domain";

export interface VariableFormProps {
  variables: VariableDefinition[];
  values: { [key: string]: any };
  onChange: (name: string, value: any) => void;
  onValidate?: () => Promise<boolean>;
}

export function VariableForm({ 
  variables, 
  values, 
  onChange, 
  onValidate 
}: VariableFormProps) {
  const [validationErrors, setValidationErrors] = useState<{ [key: string]: string }>({});
  const [isValidating, setIsValidating] = useState(false);

  const handleValidate = async () => {
    if (!onValidate) return true;
    
    setIsValidating(true);
    setValidationErrors({});
    
    try {
      const isValid = await onValidate();
      return isValid;
    } catch (error) {
      console.error("Validation error:", error);
      return false;
    } finally {
      setIsValidating(false);
    }
  };

  const renderInput = (variable: VariableDefinition) => {
    const value = values[variable.name] ?? variable.default_value ?? "";
    const hasError = !!validationErrors[variable.name];

    const inputClassName = `w-full px-3 py-2 border rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-white ${
      hasError 
        ? "border-red-500 dark:border-red-500" 
        : "border-gray-300 dark:border-gray-700"
    }`;

    switch (variable.type) {
      case "string":
        return (
          <input
            type="text"
            value={String(value)}
            onChange={(e) => onChange(variable.name, e.target.value)}
            className={inputClassName}
            aria-label={variable.label}
            aria-describedby={variable.description ? `${variable.name}-description` : undefined}
            aria-invalid={hasError}
          />
        );

      case "number":
        return (
          <input
            type="number"
            value={value === "" ? "" : Number(value)}
            onChange={(e) => {
              const inputValue = e.target.value;
              if (inputValue === "") {
                onChange(variable.name, "");
              } else {
                const numValue = parseFloat(inputValue);
                onChange(variable.name, isNaN(numValue) ? "" : numValue);
              }
            }}
            className={inputClassName}
            aria-label={variable.label}
            aria-describedby={variable.description ? `${variable.name}-description` : undefined}
            aria-invalid={hasError}
          />
        );

      case "boolean":
        return (
          <div className="flex items-center">
            <input
              type="checkbox"
              checked={Boolean(value)}
              onChange={(e) => onChange(variable.name, e.target.checked)}
              className="h-4 w-4 text-blue-500 focus:ring-blue-500 border-gray-300 rounded"
              aria-label={variable.label}
              aria-describedby={variable.description ? `${variable.name}-description` : undefined}
            />
          </div>
        );

      case "date":
        return (
          <input
            type="date"
            value={String(value)}
            onChange={(e) => onChange(variable.name, e.target.value)}
            className={inputClassName}
            aria-label={variable.label}
            aria-describedby={variable.description ? `${variable.name}-description` : undefined}
            aria-invalid={hasError}
          />
        );

      default:
        return null;
    }
  };

  if (variables.length === 0) {
    return null;
  }

  return (
    <div className="bg-white dark:bg-gray-800 p-4 rounded-lg shadow">
      <h2 className="text-lg font-semibold mb-4">Variables</h2>
      
      <div className="space-y-4">
        {variables.map((variable) => (
          <div key={variable.name}>
            <label className="block text-sm font-medium mb-1">
              {variable.label}
              {variable.required && (
                <span className="text-red-500 ml-1" aria-label="required">*</span>
              )}
            </label>
            
            {renderInput(variable)}
            
            {variable.description && (
              <p 
                id={`${variable.name}-description`}
                className="mt-1 text-xs text-gray-500 dark:text-gray-400"
              >
                {variable.description}
              </p>
            )}
            
            {validationErrors[variable.name] && (
              <p className="mt-1 text-xs text-red-500" role="alert">
                {validationErrors[variable.name]}
              </p>
            )}
          </div>
        ))}
      </div>

      {onValidate && (
        <button
          onClick={handleValidate}
          disabled={isValidating}
          className="mt-4 w-full px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 disabled:bg-gray-400 disabled:cursor-not-allowed"
        >
          {isValidating ? "Validating..." : "Validate"}
        </button>
      )}
    </div>
  );
}
