import { useState, useEffect } from "react";
import { useParams, Link } from "react-router-dom";
import ReactMarkdown from "react-markdown";
import { Document, VariableDefinition } from "../types/domain";

function DocumentDetailPage() {
  const { docId } = useParams<{ docId: string }>();

  const [document, setDocument] = useState<Document | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [variableValues, setVariableValues] = useState<Record<string, string | number | boolean>>({});
  const [processedContent, setProcessedContent] = useState<string>("");

  // API base URL
  const apiHost = import.meta.env.VITE_API_HOST;
  const apiUrl = apiHost ? `http://${apiHost}/api/v1` : "/api/v1";

  // Fetch document on component mount
  useEffect(() => {
    if (docId) {
      fetchDocument();
    }
  }, [docId]);

  // Initialize variable values when document is loaded
  useEffect(() => {
    if (document?.current_version?.variables) {
      const initialValues: Record<string, string | number | boolean> = {};
      document.current_version.variables.forEach((v) => {
        initialValues[v.name] = v.default_value ?? "";
      });
      setVariableValues(initialValues);
    }
  }, [document]);

  // Process content with variable substitution
  useEffect(() => {
    if (document?.current_version?.content) {
      let content = document.current_version.content;
      Object.entries(variableValues).forEach(([name, value]) => {
        const escapedName = name.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
        const regex = new RegExp(`\\{\\{\\s*${escapedName}\\s*\\}\\}`, "g");
        content = content.replace(regex, String(value));
      });
      setProcessedContent(content);
    }
  }, [document, variableValues]);

  const fetchDocument = async () => {
    setIsLoading(true);
    setError(null);

    try {
      const response = await fetch(`${apiUrl}/documents/${docId}`);
      if (!response.ok) {
        if (response.status === 404) {
          setError("Document not found");
        } else {
          throw new Error(`HTTP error! status: ${response.status}`);
        }
        return;
      }
      const data = await response.json();
      setDocument(data);
    } catch (err) {
      setError("Failed to load document. Please try again later.");
      console.error("Error fetching document:", err);
    } finally {
      setIsLoading(false);
    }
  };

  const handleVariableChange = (name: string, value: string | number | boolean) => {
    setVariableValues((prev) => ({
      ...prev,
      [name]: value,
    }));
  };

  const renderVariableInput = (variable: VariableDefinition) => {
    const value = variableValues[variable.name] ?? "";

    switch (variable.type) {
      case "boolean":
        return (
          <input
            type="checkbox"
            checked={Boolean(value)}
            onChange={(e) => handleVariableChange(variable.name, e.target.checked)}
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
                handleVariableChange(variable.name, "");
              } else {
                const numValue = parseFloat(inputValue);
                handleVariableChange(variable.name, isNaN(numValue) ? 0 : numValue);
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
            onChange={(e) => handleVariableChange(variable.name, e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 dark:border-gray-700 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-white"
          />
        );
      default:
        return (
          <input
            type="text"
            value={String(value)}
            onChange={(e) => handleVariableChange(variable.name, e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 dark:border-gray-700 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-white"
          />
        );
    }
  };

  if (isLoading) {
    return (
      <div className="text-center py-12">
        <p className="text-gray-500">Loading document...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="space-y-4">
        <div className="p-4 bg-red-100 text-red-800 dark:bg-red-800 dark:text-red-100 rounded">
          {error}
        </div>
        <Link
          to="/documents"
          className="inline-block px-4 py-2 bg-gray-200 text-gray-800 rounded hover:bg-gray-300 dark:bg-gray-700 dark:text-gray-200"
        >
          Back to Documents
        </Link>
      </div>
    );
  }

  if (!document) {
    return null;
  }

  const currentVersion = document.current_version;

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold">{currentVersion?.title || "Untitled Document"}</h1>
          <p className="text-sm text-gray-500 dark:text-gray-400">
            Version {currentVersion?.version_number || 1} • by {document.owner}
          </p>
        </div>
        <div className="flex gap-2">
          <Link
            to={`/documents/${docId}/versions`}
            className="px-4 py-2 bg-gray-200 text-gray-800 rounded hover:bg-gray-300 dark:bg-gray-700 dark:text-gray-200"
          >
            Version History
          </Link>
          <Link
            to="/documents"
            className="px-4 py-2 bg-gray-200 text-gray-800 rounded hover:bg-gray-300 dark:bg-gray-700 dark:text-gray-200"
          >
            Back to List
          </Link>
        </div>
      </div>

      {/* Metadata */}
      <div className="bg-white dark:bg-gray-800 p-4 rounded-lg shadow">
        <div className="flex flex-wrap gap-4 text-sm">
          <div>
            <span className="font-medium">Type:</span>{" "}
            <span
              className={`px-2 py-0.5 rounded-full ${
                currentVersion?.doc_type === "procedure"
                  ? "bg-purple-100 text-purple-800 dark:bg-purple-900 dark:text-purple-200"
                  : "bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200"
              }`}
            >
              {currentVersion?.doc_type || "unknown"}
            </span>
          </div>
          <div>
            <span className="font-medium">Status:</span>{" "}
            <span
              className={`px-2 py-0.5 rounded-full ${
                document.is_published
                  ? "bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200"
                  : "bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-200"
              }`}
            >
              {document.is_published ? "Published" : "Draft"}
            </span>
          </div>
          <div>
            <span className="font-medium">Access:</span>{" "}
            <span className="capitalize">{document.access_scope}</span>
          </div>
          <div>
            <span className="font-medium">Auto Update:</span>{" "}
            {document.is_auto_update ? "Enabled" : "Disabled"}
          </div>
        </div>

        {/* Tags */}
        {currentVersion?.tags && currentVersion.tags.length > 0 && (
          <div className="mt-4">
            <span className="font-medium text-sm">Tags:</span>
            <div className="flex flex-wrap gap-1 mt-1">
              {currentVersion.tags.map((tag) => (
                <span
                  key={tag}
                  className="px-2 py-0.5 text-xs bg-gray-100 dark:bg-gray-700 rounded"
                >
                  {tag}
                </span>
              ))}
            </div>
          </div>
        )}
      </div>

      <div className="flex gap-6">
        {/* Variable Input Panel */}
        {currentVersion?.variables && currentVersion.variables.length > 0 && (
          <div className="w-80 flex-shrink-0">
            <div className="bg-white dark:bg-gray-800 p-4 rounded-lg shadow sticky top-20">
              <h2 className="text-lg font-semibold mb-4">Variables</h2>
              <div className="space-y-4">
                {currentVersion.variables.map((variable) => (
                  <div key={variable.name}>
                    <label className="block text-sm font-medium mb-1">
                      {variable.label}
                      {variable.required && (
                        <span className="text-red-500 ml-1">*</span>
                      )}
                    </label>
                    {renderVariableInput(variable)}
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
        )}

        {/* Document Content */}
        <div className="flex-1">
          <div className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow prose dark:prose-invert max-w-none">
            <ReactMarkdown>{processedContent}</ReactMarkdown>
          </div>
        </div>
      </div>
    </div>
  );
}

export default DocumentDetailPage;
