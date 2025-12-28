import ReactMarkdown from "react-markdown";
import { DocumentVariable } from "@/shared/types/domain";

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
  return (
    <div className="bg-white p-6 rounded-lg shadow-md">
      <h2 className="text-xl font-bold mb-4">{documentTitle}</h2>

      {/* Variables Section */}
      {variables.length > 0 && (
        <div className="mb-6 p-4 bg-gray-50 rounded">
          <h3 className="font-semibold mb-3">Variables</h3>
          <div className="space-y-3">
            {variables.map((variable) => (
              <div key={variable.id}>
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
                  value={variableValues[variable.id] || ""}
                  onChange={(e) =>
                    onVariableChange(variable.id, e.target.value)
                  }
                  disabled={isExecutionStarted}
                  placeholder={variable.default_value || ""}
                  className="w-full px-3 py-2 border rounded focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:bg-gray-100"
                />
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Document Content */}
      <div className="prose max-w-none">
        <ReactMarkdown>{processedContent}</ReactMarkdown>
      </div>
    </div>
  );
}
