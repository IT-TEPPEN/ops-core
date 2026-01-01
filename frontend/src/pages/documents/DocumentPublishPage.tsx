import { useState } from "react";
import { useNavigate, useSearchParams } from "react-router-dom";
import {
  ConnectionSelector,
  RepositoryBrowser,
  FileBrowser,
  PublishDialog,
} from "@/features/document";
import { ConnectionList } from "@/features/oauth";

/**
 * Document Publish Page
 *
 * ドキュメント公開ページ。OAuth接続からリポジトリを選択し、
 * ファイルをブラウズして公開する。
 *
 * フロー:
 * 1. Connection選択（または新規OAuth接続）
 * 2. リポジトリ選択
 * 3. ファイル選択
 * 4. 公開確認 → 公開実行
 */
export function DocumentPublishPage() {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const selectedConnectionId = searchParams.get("selected_connection_id");

  const [selectedRepository, setSelectedRepository] = useState<{
    owner: string;
    name: string;
    fullName: string;
  } | null>(null);
  const [selectedFile, setSelectedFile] = useState<{
    path: string;
    url: string;
  } | null>(null);
  const [isPublishDialogOpen, setIsPublishDialogOpen] = useState(false);

  const handleRepositorySelect = (repository: {
    owner: string;
    name: string;
    fullName: string;
  }) => {
    setSelectedRepository(repository);
    // Reset file selection when repository changes
    setSelectedFile(null);
  };

  const handleFileSelect = (file: { path: string; url: string }) => {
    setSelectedFile(file);
    setIsPublishDialogOpen(true);
  };

  const handlePublish = () => {
    navigate("/documents");
  };

  return (
    <div className="h-screen flex flex-col">
      <header className="bg-white dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700 px-6 py-4">
        <h1 className="text-2xl font-bold text-gray-900 dark:text-gray-100">
          Publish Document
        </h1>
        <p className="text-sm text-gray-600 dark:text-gray-400 mt-1">
          Select a file from your Git repository to publish as a document
        </p>
      </header>

      <div className="flex-1 flex overflow-hidden">
        {/* Left Pane: Connection & Repository Selection */}
        <aside className="w-80 bg-white dark:bg-gray-800 border-r border-gray-200 dark:border-gray-700 flex flex-col overflow-hidden">
          <div className="flex-none p-4 border-b border-gray-200 dark:border-gray-700">
            <h2 className="text-lg font-semibold text-gray-900 dark:text-gray-100">
              Step 1 & 2
            </h2>
          </div>

          {/* Connection Selector */}
          <div className="flex-none p-4 border-b border-gray-200 dark:border-gray-700">
            <div className="space-y-3">
              <ConnectionSelector />
              <ConnectionList />
            </div>
          </div>

          {/* Repository Browser */}
          <div className="flex-1 overflow-y-auto p-4">
            {selectedConnectionId ? (
              <RepositoryBrowser
                connectionId={selectedConnectionId}
                selectedRepository={selectedRepository}
                onRepositorySelect={handleRepositorySelect}
              />
            ) : (
              <div className="text-center text-gray-500 dark:text-gray-400 py-8">
                <p>Select a connection to view repositories</p>
              </div>
            )}
          </div>
        </aside>

        {/* Right Pane: File Browser */}
        <main className="flex-1 bg-gray-50 dark:bg-gray-900 overflow-hidden">
          {selectedRepository ? (
            <FileBrowser
              connectionId={selectedConnectionId!}
              repository={selectedRepository}
              onFileSelect={handleFileSelect}
            />
          ) : (
            <div className="flex items-center justify-center h-full">
              <div className="text-center text-gray-500 dark:text-gray-400">
                <svg
                  className="mx-auto h-12 w-12 text-gray-400"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
                  />
                </svg>
                <h3 className="mt-2 text-sm font-medium">
                  No repository selected
                </h3>
                <p className="mt-1 text-sm">
                  Select a repository from the left panel to browse files
                </p>
              </div>
            </div>
          )}
        </main>
      </div>

      {/* Publish Dialog */}
      <PublishDialog
        isOpen={isPublishDialogOpen}
        file={selectedFile}
        onClose={() => setIsPublishDialogOpen(false)}
        onConfirm={handlePublish}
      />
    </div>
  );
}
