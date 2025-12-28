import { Link } from "react-router-dom";
import { useDocumentList } from "@/features/document/hooks/useDocumentList";
import { DocumentFilterBar } from "@/features/document/components/DocumentFilterBar";
import { DocumentTable } from "@/features/document/components/DocumentTable";

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
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">Documents</h1>
        <Link
          to="/documents/new"
          className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
        >
          Create Document
        </Link>
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

      {isLoading && <p className="text-gray-500">Loading documents...</p>}

      {!isLoading && !error && filteredDocuments.length === 0 && (
        <div className="text-center py-12 bg-white dark:bg-gray-800 rounded-lg shadow">
          <h2 className="text-xl font-semibold mb-2">No Documents Found</h2>
          <p className="text-gray-600 dark:text-gray-400 mb-4">
            {documents.length === 0
              ? "Get started by creating your first document."
              : "No documents match your filters."}
          </p>
          {documents.length === 0 && (
            <Link
              to="/documents/new"
              className="inline-block px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
            >
              Create Document
            </Link>
          )}
        </div>
      )}

      {!isLoading && !error && filteredDocuments.length > 0 && (
        <DocumentTable documents={filteredDocuments} />
      )}
    </div>
  );
}
