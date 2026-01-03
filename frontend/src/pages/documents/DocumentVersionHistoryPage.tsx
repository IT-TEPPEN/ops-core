import { Link } from "react-router-dom";
import { useVersionHistory } from "@/features/document/hooks/useVersionHistory";
import { VersionTable } from "@/features/document/components/VersionTable";
import { Page } from "@/shared/types/Page";
import { LoadingSpinner } from "@/ui";

export const DocumentVersionHistoryPage: Page<"docId"> = ({
  path: { docId },
}) => {
  const {
    versions,
    isLoading,
    error,
    selectedVersions,
    actionMessage,
    isProcessing,
    handleRollback,
    handlePublish,
    toggleVersionSelection,
  } = useVersionHistory(docId);

  if (isLoading) {
    return <LoadingSpinner message="Loading..." />;
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

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">Version History</h1>
        <div className="flex gap-2">
          <Link
            to={`/documents/${docId}`}
            className="px-4 py-2 bg-gray-200 text-gray-800 rounded hover:bg-gray-300 dark:bg-gray-700 dark:text-gray-200"
          >
            Back to Document
          </Link>
        </div>
      </div>

      {/* Action message */}
      {actionMessage && (
        <div
          className={`p-3 rounded ${
            actionMessage.type === "success"
              ? "bg-green-100 text-green-800 dark:bg-green-800 dark:text-green-100"
              : "bg-red-100 text-red-800 dark:bg-red-800 dark:text-red-100"
          }`}
        >
          {actionMessage.text}
        </div>
      )}

      {/* Version comparison hint */}
      {selectedVersions[0] !== null && selectedVersions[1] === null && (
        <div className="p-3 bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200 rounded">
          Select another version to compare
        </div>
      )}

      <VersionTable
        versions={versions}
        selectedVersions={selectedVersions}
        isProcessing={isProcessing}
        onToggleVersion={toggleVersionSelection}
        onRollback={handleRollback}
        onPublish={handlePublish}
        docId={docId}
      />
    </div>
  );
};
