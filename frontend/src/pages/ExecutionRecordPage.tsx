import { useParams, Link } from "react-router-dom";
import { useEffect, useState } from "react";
import { ThreePaneLayout, VariableForm } from "@/features/common/components";
import { useExecutionRecord } from "@/features/execution/hooks/useExecutionRecord";
import { useDocumentExecution } from "@/features/document/hooks/useDocumentExecution";
import { MarkdownProcessor } from "@/features/markdown";
import { LoadingSpinner } from "@/ui";

function ExecutionRecordPage() {
  const { docId, recordId } = useParams<{
    docId: string;
    recordId?: string;
  }>();
  const [Component, setComponent] = useState<React.ReactElement | null>(null);

  // Use custom hooks for business logic - must be called before any conditional returns
  const {
    document,
    isLoading: isDocumentLoading,
    error: documentError,
    variableValues,
    processedContent,
    handleVariableChange,
  } = useDocumentExecution({ docId, recordId });

  useEffect(() => {
    MarkdownProcessor.process(processedContent).then(
      (file: { result: unknown }) => {
        setComponent(file.result as React.ReactElement);
      }
    );
  }, [processedContent]);

  const {
    executionRecord,
    error: executionError,
    isCreating,
    isSaving,
    handleStartExecution,
    handleUpdateTitle,
    handleUpdateNotes,
    handleComplete,
    handleFail,
  } = useExecutionRecord({
    docId,
    recordId,
    documentVersionId: document?.current_version?.id,
    variableValues,
    executionTitle: document?.current_version?.title
      ? `Execution of ${
          document.current_version.title
        } - ${new Date().toLocaleString()}`
      : `Execution - ${new Date().toLocaleString()}`,
  });

  // Local state for title and notes editing
  // Initialize with empty strings, sync with executionRecord when available
  const [localTitle, setLocalTitle] = useState("");
  const [localNotes, setLocalNotes] = useState("");
  const [hasInitialized, setHasInitialized] = useState(false);

  // Sync local state when executionRecord first loads or changes ID
  if (executionRecord && !hasInitialized) {
    setLocalTitle(executionRecord.title);
    setLocalNotes(executionRecord.notes);
    setHasInitialized(true);
  }

  if (!docId) {
    return (
      <div className="p-8">
        <p className="text-red-500">Document ID is required</p>
        <Link to="/documents" className="text-blue-500 hover:underline">
          Back to Documents
        </Link>
      </div>
    );
  }

  if (isDocumentLoading) {
    return <LoadingSpinner message="Loading..." />;
  }

  // Error state
  const error = documentError || executionError;
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
        {Component}
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
              value={localTitle}
              onChange={(e) => setLocalTitle(e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 dark:bg-gray-700 dark:text-white"
            />
          </div>
          <button
            onClick={() => handleStartExecution()}
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
                value={localTitle}
                onChange={(e) => setLocalTitle(e.target.value)}
                disabled={executionRecord.status !== "in_progress"}
                className="flex-1 px-3 py-2 text-sm border border-gray-300 dark:border-gray-600 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 dark:bg-gray-700 dark:text-white disabled:opacity-50"
              />
              {executionRecord.status === "in_progress" && (
                <button
                  onClick={() => handleUpdateTitle(localTitle)}
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
              value={localNotes}
              onChange={(e) => setLocalNotes(e.target.value)}
              disabled={executionRecord.status !== "in_progress"}
              className="w-full px-3 py-2 text-sm border border-gray-300 dark:border-gray-600 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 dark:bg-gray-700 dark:text-white disabled:opacity-50"
              rows={3}
            />
            {executionRecord.status === "in_progress" && (
              <button
                onClick={() => handleUpdateNotes(localNotes)}
                disabled={isSaving}
                className="mt-2 px-3 py-2 text-sm bg-blue-600 text-white rounded hover:bg-blue-700 disabled:opacity-50"
              >
                {isSaving ? "Saving..." : "Save Notes"}
              </button>
            )}
          </div>

          {/* Steps */}
          {/* {executionRecord.status === "in_progress" && (
            <ExecutionStepPanel
              steps={executionRecord.steps}
              onAddStep={handleAddStep}
              onUpdateStepNotes={handleUpdateStepNotes}
            />
          )} */}

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
            <p>
              Started: {new Date(executionRecord.createdAt).toLocaleString()}
            </p>
            {executionRecord.completedAt && (
              <p>
                Completed:{" "}
                {new Date(executionRecord.completedAt).toLocaleString()}
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
