import { Link } from "react-router-dom";
import { ThreePaneLayout } from "@/features/common/components/Layout/ThreePaneLayout";
import { VariableForm } from "@/features/common/components/Form/VariableForm";
import { DocumentContentPane } from "@/features/document/components/DocumentContentPane";
import { ExecutionRecordPane } from "@/features/document/components/ExecutionRecordPane";
import { useDocumentView } from "@/features/document/hooks/useDocumentView";
import { Page } from "@/shared/types/Page";

export const DocumentViewPage: Page<"docId"> = ({ path: { docId } }) => {
  const {
    document,
    isLoading,
    error,
    variableValues,
    processedContent,
    handleVariableChange,
    handleValidate,
  } = useDocumentView(docId);

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

  return (
    <div className="h-screen flex flex-col">
      <ThreePaneLayout
        leftPane={
          hasVariables ? (
            <div className="p-4">
              <VariableForm
                variables={currentVersion.variables}
                values={variableValues}
                onChange={handleVariableChange}
                onValidate={handleValidate}
              />
            </div>
          ) : (
            <div className="p-4 text-sm text-gray-500">
              No variables defined
            </div>
          )
        }
        centerPane={
          <DocumentContentPane
            document={document}
            processedContent={processedContent}
          />
        }
        rightPane={<ExecutionRecordPane />}
        leftWidth="w-80"
        rightWidth="w-80"
      />
    </div>
  );
};
