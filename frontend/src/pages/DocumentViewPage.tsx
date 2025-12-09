import { useState, useEffect } from "react";
import { useParams, Link } from "react-router-dom";
import ReactMarkdown from "react-markdown";
import { ThreePaneLayout } from "../components/Layout/ThreePaneLayout";
import { VariableForm } from "../components/Form/VariableForm";
import { Document, VariableDefinition } from "../types/domain";
import { substituteVariables } from "../utils/variableSubstitution";

function DocumentViewPage() {
  const { docId } = useParams<{ docId: string }>();

  const [document, setDocument] = useState<Document | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [variableValues, setVariableValues] = useState<Record<string, any>>({});
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
      const initialValues: Record<string, any> = {};
      document.current_version.variables.forEach((v: VariableDefinition) => {
        initialValues[v.name] = v.default_value ?? "";
      });
      setVariableValues(initialValues);
    }
  }, [document]);

  // Process content with variable substitution when variables change
  useEffect(() => {
    if (document?.current_version?.content) {
      const substituted = substituteVariables(
        document.current_version.content,
        variableValues
      );
      setProcessedContent(substituted);
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

  const handleVariableChange = (name: string, value: any) => {
    setVariableValues((prev) => ({
      ...prev,
      [name]: value,
    }));
  };

  const handleValidate = async (): Promise<boolean> => {
    if (!docId) return false;

    try {
      const values = Object.entries(variableValues).map(([name, value]) => ({
        name,
        value,
      }));

      const response = await fetch(
        `${apiUrl}/documents/${docId}/validate-variables`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({ values }),
        }
      );

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const result = await response.json();
      return result.valid;
    } catch (err) {
      console.error("Error validating variables:", err);
      return false;
    }
  };

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-screen">
        <p className="text-gray-500">Loading document...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="p-8 space-y-4">
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
  const hasVariables =
    currentVersion?.variables && currentVersion.variables.length > 0;

  // Left pane: Variable input form
  const leftPane = hasVariables ? (
    <div className="p-4">
      <VariableForm
        variables={currentVersion.variables}
        values={variableValues}
        onChange={handleVariableChange}
        onValidate={handleValidate}
      />
    </div>
  ) : (
    <div className="p-4 text-sm text-gray-500">No variables defined</div>
  );

  // Center pane: Document content
  const centerPane = (
    <div className="p-6">
      {/* Header */}
      <div className="mb-6">
        <div className="flex items-center justify-between mb-2">
          <h1 className="text-3xl font-bold">
            {currentVersion?.title || "Untitled Document"}
          </h1>
          <Link
            to="/documents"
            className="px-4 py-2 bg-gray-200 text-gray-800 rounded hover:bg-gray-300 dark:bg-gray-700 dark:text-gray-200"
          >
            Back to List
          </Link>
        </div>
        <p className="text-sm text-gray-500 dark:text-gray-400">
          Version {currentVersion?.version_number || 1} • by {document.owner}
        </p>
      </div>

      {/* Metadata */}
      <div className="bg-white dark:bg-gray-800 p-4 rounded-lg shadow mb-6">
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
        </div>

        {/* Tags */}
        {currentVersion?.tags && currentVersion.tags.length > 0 && (
          <div className="mt-4">
            <span className="font-medium text-sm">Tags:</span>
            <div className="flex flex-wrap gap-1 mt-1">
              {currentVersion.tags.map((tag: string) => (
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

      {/* Document Content */}
      <div className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow prose dark:prose-invert max-w-none">
        <ReactMarkdown>{processedContent}</ReactMarkdown>
      </div>
    </div>
  );

  // Right pane: Execution record / Work evidence panel (placeholder)
  const rightPane = (
    <div className="p-4">
      <h2 className="text-lg font-semibold mb-4">Execution Record</h2>
      <div className="text-sm text-gray-500">
        <p>Work evidence tracking will be displayed here.</p>
        <p className="mt-2">Features:</p>
        <ul className="mt-1 ml-4 list-disc space-y-1">
          <li>Execution status</li>
          <li>Step-by-step notes</li>
          <li>Evidence attachments</li>
          <li>Completion timestamps</li>
        </ul>
      </div>
    </div>
  );

  return (
    <div className="h-screen flex flex-col">
      <ThreePaneLayout
        leftPane={leftPane}
        centerPane={centerPane}
        rightPane={rightPane}
        leftWidth="w-80"
        rightWidth="w-80"
      />
    </div>
  );
}

export default DocumentViewPage;
