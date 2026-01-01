import { VariableDefinition } from "@/shared/types/domain";

interface VariableInputPanelProps {
  variables: VariableDefinition[];
  values: Record<string, string | number | boolean>;
  onChange: (name: string, value: string | number | boolean) => void;
}

export function VariableInputPanel({
  variables,
  values,
  onChange,
}: VariableInputPanelProps) {
  const renderInput = (variable: VariableDefinition) => {
    const value = values[variable.name] ?? "";

    switch (variable.type) {
      case "boolean":
        return (
          <input
            type="checkbox"
            checked={Boolean(value)}
            onChange={(e) => onChange(variable.name, e.target.checked)}
            className="h-4 w-4 text-blue-500 focus:ring-blue-500 border-gray-300 rounded"
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
            className="w-full px-3 py-2 border border-gray-300 dark:border-gray-700 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-white"
          />
        );
      case "date":
        return (
          <input
            type="date"
            value={String(value)}
            onChange={(e) => onChange(variable.name, e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 dark:border-gray-700 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-white"
          />
        );
      default:
        return (
          <input
            type="text"
            value={String(value)}
            onChange={(e) => onChange(variable.name, e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 dark:border-gray-700 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-white"
          />
        );
    }
  };

  if (variables.length === 0) {
    return null;
  }

  return (
    <div className="w-80 shrink-0">
      <div className="bg-white dark:bg-gray-800 p-4 rounded-lg shadow sticky top-20">
        <h2 className="text-lg font-semibold mb-4">Variables</h2>
        <div className="space-y-4">
          {variables.map((variable) => (
            <div key={variable.name}>
              <label className="block text-sm font-medium mb-1">
                {variable.label}
                {variable.required && (
                  <span className="text-red-500 ml-1">*</span>
                )}
              </label>
              {renderInput(variable)}
              {variable.description && (
                <p className="mt-1 text-xs text-gray-500 dark:text-gray-400">
                  {variable.description}
                </p>
              )}
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
