import { Link } from "react-router-dom";
import { useDocumentList } from "@/features/document/hooks/useDocumentList";
import { DocumentFilterBar } from "@/features/document/components/DocumentFilterBar";
import { DocumentTable } from "@/features/document/components/DocumentTable";
import { LoadingSpinner, UI_Icon_PlusIcon } from "@/ui";

export function DocumentListPage() {
  const {
    documents,
    filteredDocuments,
    isLoading,
    error,
    filterType,
    searchQuery,
    setFilterType,
    setSearchQuery,
  } = useDocumentList();

  return (
    <div className="space-y-6 p-4">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">Documents</h1>
        <PublishDocumentButton />
      </div>

      <DocumentFilterBar
        filterType={filterType}
        searchQuery={searchQuery}
        onFilterChange={setFilterType}
        onSearchChange={setSearchQuery}
      />

      {error && (
        <div className="p-3 bg-red-100 text-red-800 dark:bg-red-800 dark:text-red-100 rounded">
          {error}
        </div>
      )}

      {isLoading && <LoadingSpinner message="Loading..." />}

      {!isLoading && !error && filteredDocuments.length === 0 && (
        <div className="flex flex-col justify-center items-center gap-4 py-12 bg-white dark:bg-gray-800 rounded-lg shadow">
          <div className="text-center">
            <h2 className="text-xl font-semibold mb-2">No Documents Found</h2>
            <p className="text-gray-600 dark:text-gray-400 mb-4">
              {documents.length === 0
                ? "Get started by creating your first document."
                : "No documents match your filters."}
            </p>
          </div>
          {documents.length === 0 && <PublishDocumentButton />}
        </div>
      )}

      {!isLoading && !error && filteredDocuments.length > 0 && (
        <DocumentTable documents={filteredDocuments} />
      )}
    </div>
  );
}

function PublishDocumentButton() {
  return (
    <Link
      to="/documents/publish"
      className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors flex items-center gap-2 w-fit"
    >
      <UI_Icon_PlusIcon />
      Publish Document
    </Link>
  );
}
