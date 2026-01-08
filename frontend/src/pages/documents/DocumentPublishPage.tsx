import { Suspense } from "react";
import { useNavigate, useSearchParams } from "react-router-dom";
import { ErrorBoundary } from "react-error-boundary";
import {
  ConnectionSelector,
  RepositoryBrowser,
  FileBrowser,
  PublishDialog,
  useDocumentCommandService,
} from "@/features/document";
import { ConnectionList } from "@/features/oauth";
import { LoadingSpinner, UI_Icon_Document } from "@/ui";

export function DocumentPublishPage() {
  const navigate = useNavigate();
  const documentService = useDocumentCommandService();
  const [searchParams, setSearchParams] = useSearchParams();
  const selectedConnectionId = searchParams.get("selected_connection_id");
  const selectedRepositoryId = searchParams.get("selected_repository_id");
  const selectedRepositoryFullName = searchParams.get(
    "selected_repository_full_name"
  );
  const selectedFile = searchParams.get("selected_file");

  const isSelectedFile = !!selectedFile;

  const onClose = () => {
    setSearchParams((searchParams) => {
      searchParams.delete("selected_file");
      return searchParams;
    });
  };

  const handlePublish = (options: { accessScope: "public" | "private" }) => {
    if (
      !selectedFile ||
      !selectedConnectionId ||
      !selectedRepositoryId ||
      !selectedRepositoryFullName
    )
      return;

    const parts = selectedRepositoryFullName.split("/");
    const repository = parts.pop()!;
    const owner = parts.join("/");

    documentService
      .publishFromOAuth({
        filePath: selectedFile,
        connectionId: selectedConnectionId,
        providerRepositoryId: selectedRepositoryId,
        repository,
        owner,
        accessScope: options.accessScope,
        ref: "main",
        isAutoUpdate: true,
      })
      .then(() => {
        navigate("/documents");
      });
  };

  return (
    <div className="h-full flex flex-col p-4">
      <header className="bg-white dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700 px-6 py-4">
        <h1 className="text-2xl font-bold text-gray-900 dark:text-gray-100">
          Publish Document
        </h1>
        <p className="text-sm text-gray-600 dark:text-gray-400 mt-1">
          Select a file from your Git repository to publish as a document
        </p>
      </header>

      <div className="flex-1 flex overflow-hidden">
        <aside className="w-80 bg-white dark:bg-gray-800 border-r border-gray-200 dark:border-gray-700 flex flex-col overflow-hidden">
          <div className="flex-none p-4 border-b border-gray-200 dark:border-gray-700">
            <h2 className="text-lg font-semibold text-gray-900 dark:text-gray-100">
              Select Connection & Repository
            </h2>
          </div>

          <div className="flex-none p-4 border-b border-gray-200 dark:border-gray-700">
            <div className="space-y-3">
              <ConnectionSelector />
              <ErrorBoundary
                fallback={
                  <p className="text-red-600">Failed to load connections.</p>
                }
              >
                <Suspense fallback={<LoadingSpinner message="Loading..." />}>
                  <ConnectionList />
                </Suspense>
              </ErrorBoundary>
            </div>
          </div>

          <div className="flex-1 overflow-y-auto p-4">
            {selectedConnectionId ? (
              <ErrorBoundary
                fallback={
                  <p className="text-red-600">
                    Failed to load repositories. Please try again.
                  </p>
                }
              >
                <Suspense fallback={<LoadingSpinner message="Loading..." />}>
                  <RepositoryBrowser connectionId={selectedConnectionId} />
                </Suspense>
              </ErrorBoundary>
            ) : (
              <div className="text-center text-gray-500 dark:text-gray-400 py-8">
                <p>Select a connection to view repositories</p>
              </div>
            )}
          </div>
        </aside>

        <main className="flex-1 bg-gray-50 dark:bg-gray-900 overflow-hidden">
          {selectedConnectionId &&
          selectedRepositoryId &&
          selectedRepositoryFullName ? (
            <FileBrowser
              connectionId={selectedConnectionId}
              repositoryId={selectedRepositoryId}
              repositoryFullName={selectedRepositoryFullName}
            />
          ) : (
            <div className="flex items-center justify-center h-full">
              <div className="flex flex-col items-center gap-2 text-center text-gray-500 dark:text-gray-400">
                <UI_Icon_Document size={40} />
                <div className="text-center">
                  <h3 className="mt-2 text-sm font-medium">
                    No repository selected
                  </h3>
                  <p className="mt-1 text-sm">
                    Select a repository from the left panel to browse files
                  </p>
                </div>
              </div>
            </div>
          )}
        </main>
      </div>

      {isSelectedFile && (
        <PublishDialog
          filePath={selectedFile}
          onClose={onClose}
          onConfirm={handlePublish}
        />
      )}
    </div>
  );
}
