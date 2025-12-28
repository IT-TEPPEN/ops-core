import { Link } from "react-router-dom";
import ReactMarkdown from "react-markdown";
import { useDocumentExecution } from "@/features/document/hooks/useDocumentExecution";
import { DocumentMetadata } from "@/features/document/components/DocumentMetadata";
import { VariableInputPanel } from "@/features/document/components/VariableInputPanel";
import { Page } from "@/shared/types/Page";

export const DocumentDetailPage: Page<"docId"> = ({ path: { docId } }) => {
  const {
    document,
    isLoading,
    error,
    variableValues,
    processedContent,
    handleVariableChange,
  } = useDocumentExecution(docId);

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
          <h1 className="text-2xl font-bold">
            {currentVersion?.title || "Untitled Document"}
          </h1>
          <p className="text-sm text-gray-500 dark:text-gray-400">
            Version {currentVersion?.version_number || 1} • by {document.owner}
          </p>
        </div>
        <div className="flex gap-2">
          <Link
            to={`/documents/${docId}/view`}
            className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
          >
            Three-Pane View
          </Link>
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

      <DocumentMetadata document={document} />

      <div className="flex gap-6">
        {currentVersion?.variables && currentVersion.variables.length > 0 && (
          <VariableInputPanel
            variables={currentVersion.variables}
            values={variableValues}
            onChange={handleVariableChange}
          />
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
};
