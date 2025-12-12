import { useState, useEffect } from "react";
import { useParams, Link } from "react-router-dom";
import ReactMarkdown from "react-markdown";
import { ThreePaneLayout } from "../components/Layout/ThreePaneLayout";
import { VariableForm } from "../components/Form/VariableForm";
import { ExecutionStepPanel } from "../components/Form/ExecutionStepPanel";
import {
  Document,
  VariableDefinition,
  ExecutionRecord,
} from "../types/domain";
import { substituteVariables } from "../utils/variableSubstitution";
import {
  createExecutionRecord,
  getExecutionRecord,
  updateExecutionRecordTitle,
  updateExecutionRecordNotes,
  addExecutionStep,
  updateStepNotes,
  completeExecutionRecord,
  failExecutionRecord,
} from "../api/executionRecordApi";

function ExecutionRecordPage() {
  const { docId, recordId } = useParams<{
    docId: string;
    recordId?: string;
  }>();

  const [document, setDocument] = useState<Document | null>(null);
  const [executionRecord, setExecutionRecord] =
    useState<ExecutionRecord | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [variableValues, setVariableValues] = useState<Record<string, string | number | boolean>>({});
  const [processedContent, setProcessedContent] = useState<string>("");
  const [executionTitle, setExecutionTitle] = useState<string>("");
  const [executionNotes, setExecutionNotes] = useState<string>("");
  const [isCreating, setIsCreating] = useState(false);
  const [isSaving, setIsSaving] = useState(false);

  // API base URL
  const apiHost = import.meta.env.VITE_API_HOST;
  const apiUrl = apiHost ? `http://${apiHost}/api/v1` : "/api/v1";

  // Fetch document on component mount
  useEffect(() => {
    if (!docId) return;

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

        // Set default title
        if (!executionTitle) {
          setExecutionTitle(
            `Execution of ${data.current_version?.title || "Document"} - ${new Date().toLocaleString()}`
          );
        }
      } catch (err) {
        setError("Failed to load document. Please try again later.");
        console.error("Error fetching document:", err);
      } finally {
        setIsLoading(false);
      }
    };

    fetchDocument();
  }, [docId, apiUrl]);

  // Fetch execution record if recordId is provided
  useEffect(() => {
    if (!recordId) return;

    const fetchRecord = async () => {
      try {
        const record = await getExecutionRecord(recordId);
        setExecutionRecord(record);
        setExecutionTitle(record.title);
        setExecutionNotes(record.notes);

        // Populate variable values from record
        const values: Record<string, string | number | boolean> = {};
        record.variable_values.forEach((vv) => {
          values[vv.name] = vv.value;
        });
        setVariableValues(values);
      } catch (err) {
        console.error("Error fetching execution record:", err);
        setError("Failed to load execution record");
      }
    };

    fetchRecord();
  }, [recordId]);

  // Initialize variable values when document is loaded
  useEffect(() => {
    if (document?.current_version?.variables && !recordId) {
      const initialValues: Record<string, string | number | boolean> = {};
      document.current_version.variables.forEach((v: VariableDefinition) => {
        // Use type-specific defaults
        if (v.default_value !== undefined && v.default_value !== null) {
          initialValues[v.name] = v.default_value;
        } else {
          switch (v.type) {
            case "number":
              initialValues[v.name] = 0;
              break;
            case "boolean":
              initialValues[v.name] = false;
              break;
            default:
              initialValues[v.name] = "";
          }
        }
      });
      setVariableValues(initialValues);
    }
  }, [document, recordId]);

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

  const handleVariableChange = (name: string, value: string | number | boolean) => {
    setVariableValues((prev) => ({
      ...prev,
      [name]: value,
    }));
  };

  const handleStartExecution = async () => {
    if (!docId || !document?.current_version) return;

    setIsCreating(true);
    try {
      const variableValuesList = Object.entries(variableValues).map(
        ([name, value]) => ({ name, value })
      );

      const record = await createExecutionRecord({
        document_id: docId,
        document_version_id: document.current_version.id,
        title: executionTitle,
        variable_values: variableValuesList,
      });

      setExecutionRecord(record);
      setError(null);
    } catch (err) {
      console.error("Failed to create execution record:", err);
      setError("Failed to start execution");
    } finally {
      setIsCreating(false);
    }
  };

  const handleUpdateTitle = async () => {
    if (!executionRecord) return;

    setIsSaving(true);
    try {
      const updated = await updateExecutionRecordTitle(
        executionRecord.id,
        executionTitle
      );
      setExecutionRecord(updated);
    } catch (err) {
      console.error("Failed to update title:", err);
    } finally {
      setIsSaving(false);
    }
  };

  const handleUpdateNotes = async () => {
    if (!executionRecord) return;

    setIsSaving(true);
    try {
      const updated = await updateExecutionRecordNotes(
        executionRecord.id,
        executionNotes
      );
      setExecutionRecord(updated);
    } catch (err) {
      console.error("Failed to update notes:", err);
    } finally {
      setIsSaving(false);
    }
  };

  const handleAddStep = async (stepNumber: number, description: string) => {
    if (!executionRecord) return;

    const updated = await addExecutionStep(
      executionRecord.id,
      stepNumber,
      description
    );
    setExecutionRecord(updated);
  };

  const handleUpdateStepNotes = async (stepNumber: number, notes: string) => {
    if (!executionRecord) return;

    const updated = await updateStepNotes(
      executionRecord.id,
      stepNumber,
      notes
    );
    setExecutionRecord(updated);
  };

  const handleComplete = async () => {
    if (!executionRecord) return;

    setIsSaving(true);
    try {
      const updated = await completeExecutionRecord(executionRecord.id);
      setExecutionRecord(updated);
    } catch (err) {
      console.error("Failed to complete execution:", err);
    } finally {
      setIsSaving(false);
    }
  };

  const handleFail = async () => {
    if (!executionRecord) return;

    setIsSaving(true);
    try {
      const updated = await failExecutionRecord(executionRecord.id);
      setExecutionRecord(updated);
    } catch (err) {
      console.error("Failed to mark as failed:", err);
    } finally {
      setIsSaving(false);
    }
  };

  if (isLoading) {
    return (
      <div
        className="flex items-center justify-center h-screen"
        role="status"
        aria-live="polite"
      >
        <p className="text-gray-500">Loading document...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="p-8 space-y-4">
        <div
          className="p-4 bg-red-100 text-red-800 dark:bg-red-800 dark:text-red-100 rounded"
          role="alert"
          aria-live="assertive"
        >
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
      <h2 className="text-lg font-semibold mb-4">Variables</h2>
      <VariableForm
        variables={currentVersion.variables}
        values={variableValues}
        onChange={handleVariableChange}
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

      {/* Document Content */}
      <div className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow prose dark:prose-invert max-w-none">
        <ReactMarkdown>{processedContent}</ReactMarkdown>
      </div>
    </div>
  );

  // Right pane: Execution record panel
  const rightPane = (
    <div className="p-4 space-y-4">
      <h2 className="text-lg font-semibold">Execution Record</h2>

      {!executionRecord ? (
        <div className="space-y-4">
          <div>
            <label className="block text-sm font-medium mb-1">
              Execution Title
            </label>
            <input
              type="text"
              value={executionTitle}
              onChange={(e) => setExecutionTitle(e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 dark:bg-gray-700 dark:text-white"
            />
          </div>
          <button
            onClick={handleStartExecution}
            disabled={isCreating}
            className="w-full px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 disabled:opacity-50"
          >
            {isCreating ? "Starting..." : "Start Execution"}
          </button>
        </div>
      ) : (
        <div className="space-y-4">
          {/* Status Badge */}
          <div>
            <span
              className={`inline-block px-3 py-1 text-sm font-medium rounded-full ${
                executionRecord.status === "completed"
                  ? "bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200"
                  : executionRecord.status === "failed"
                    ? "bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200"
                    : "bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200"
              }`}
            >
              {executionRecord.status}
            </span>
          </div>

          {/* Title */}
          <div>
            <label className="block text-sm font-medium mb-1">Title</label>
            <div className="flex gap-2">
              <input
                type="text"
                value={executionTitle}
                onChange={(e) => setExecutionTitle(e.target.value)}
                disabled={executionRecord.status !== "in_progress"}
                className="flex-1 px-3 py-2 text-sm border border-gray-300 dark:border-gray-600 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 dark:bg-gray-700 dark:text-white disabled:opacity-50"
              />
              {executionRecord.status === "in_progress" && (
                <button
                  onClick={handleUpdateTitle}
                  disabled={isSaving}
                  className="px-3 py-2 text-sm bg-blue-600 text-white rounded hover:bg-blue-700 disabled:opacity-50"
                >
                  Save
                </button>
              )}
            </div>
          </div>

          {/* Notes */}
          <div>
            <label className="block text-sm font-medium mb-1">
              Overall Notes
            </label>
            <textarea
              value={executionNotes}
              onChange={(e) => setExecutionNotes(e.target.value)}
              disabled={executionRecord.status !== "in_progress"}
              className="w-full px-3 py-2 text-sm border border-gray-300 dark:border-gray-600 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 dark:bg-gray-700 dark:text-white disabled:opacity-50"
              rows={3}
            />
            {executionRecord.status === "in_progress" && (
              <button
                onClick={handleUpdateNotes}
                disabled={isSaving}
                className="mt-2 px-3 py-2 text-sm bg-blue-600 text-white rounded hover:bg-blue-700 disabled:opacity-50"
              >
                {isSaving ? "Saving..." : "Save Notes"}
              </button>
            )}
          </div>

          {/* Steps */}
          {executionRecord.status === "in_progress" && (
            <ExecutionStepPanel
              steps={executionRecord.steps}
              onAddStep={handleAddStep}
              onUpdateStepNotes={handleUpdateStepNotes}
            />
          )}

          {/* Action Buttons */}
          {executionRecord.status === "in_progress" && (
            <div className="flex gap-2 pt-4 border-t border-gray-200 dark:border-gray-700">
              <button
                onClick={handleComplete}
                disabled={isSaving}
                className="flex-1 px-4 py-2 bg-green-600 text-white rounded hover:bg-green-700 disabled:opacity-50"
              >
                Complete
              </button>
              <button
                onClick={handleFail}
                disabled={isSaving}
                className="flex-1 px-4 py-2 bg-red-600 text-white rounded hover:bg-red-700 disabled:opacity-50"
              >
                Mark as Failed
              </button>
            </div>
          )}

          {/* Timestamps */}
          <div className="text-xs text-gray-500 dark:text-gray-400 space-y-1 pt-4 border-t border-gray-200 dark:border-gray-700">
            <p>Started: {new Date(executionRecord.started_at).toLocaleString()}</p>
            {executionRecord.completed_at && (
              <p>
                Completed:{" "}
                {new Date(executionRecord.completed_at).toLocaleString()}
              </p>
            )}
          </div>
        </div>
      )}
    </div>
  );

  return (
    <div className="h-screen flex flex-col">
      <ThreePaneLayout
        leftPane={leftPane}
        centerPane={centerPane}
        rightPane={rightPane}
        leftWidth="w-80"
        rightWidth="w-96"
      />
    </div>
  );
}

export default ExecutionRecordPage;
