import { MarkdownProcessor } from "@/features/markdown";
import { DocumentVariable } from "@/features/repository";
import { useEffect, useState } from "react";

interface DocumentContentDisplayProps {
  documentTitle: string;
  processedContent: string;
  variables: DocumentVariable[];
  variableValues: Record<string, string>;
  onVariableChange: (varId: string, value: string) => void;
  isExecutionStarted: boolean;
}

export function DocumentContentDisplay({
  documentTitle,
  processedContent,
  variables,
  variableValues,
  onVariableChange,
  isExecutionStarted,
}: DocumentContentDisplayProps) {
  const [Component, setComponent] = useState<React.ReactElement | null>(null);

  useEffect(() => {
    MarkdownProcessor.process(processedContent).then(
      (file: { result: unknown }) => {
        setComponent(file.result as React.ReactElement);
      }
    );
  }, [processedContent]);

  return (
    <div className="bg-white p-6 rounded-lg shadow-md">
      <h2 className="text-xl font-bold mb-4">{documentTitle}</h2>

      {/* Variables Section */}
      {variables.length > 0 && (
        <div className="mb-6 p-4 bg-gray-50 rounded">
          <h3 className="font-semibold mb-3">Variables</h3>
          <div className="space-y-3">
            {variables.map((variable) => (
              <div key={variable.name}>
                <label className="block text-sm font-medium mb-1">
                  {variable.name}
                  {variable.description && (
                    <span className="text-gray-500 ml-2">
                      ({variable.description})
                    </span>
                  )}
                </label>
                <input
                  type="text"
                  value={variableValues[variable.name] || ""}
                  onChange={(e) =>
                    onVariableChange(variable.name, e.target.value)
                  }
                  disabled={isExecutionStarted}
                  placeholder={
                    typeof variable.defaultValue === "string"
                      ? variable.defaultValue
                      : typeof variable.defaultValue === "number"
                      ? variable.defaultValue.toString()
                      : typeof variable.defaultValue === "boolean"
                      ? variable.defaultValue
                        ? "true"
                        : "false"
                      : ""
                  }
                  className="w-full px-3 py-2 border rounded focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:bg-gray-100"
                />
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Document Content */}
      <div className="prose max-w-none">{Component}</div>
    </div>
  );
}
