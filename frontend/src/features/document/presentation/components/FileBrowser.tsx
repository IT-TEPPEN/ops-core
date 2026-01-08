import { useOAuthQueryService } from "@/features/oauth";
import { Content } from "@/features/oauth/application/dto";
import { LoadingSpinner } from "@/ui";
import { useQuery } from "@tanstack/react-query";
import { Link, useSearchParams } from "react-router-dom";

interface FileBrowserProps {
  connectionId: string;
  repositoryId: string;
  repositoryFullName: string;
}

const formatFileSize = (bytes: number): string => {
  if (bytes < 1024) return `${bytes} B`;
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
};

/**
 * FileBrowser Component
 *
 * リポジトリのファイル/ディレクトリをブラウズするコンポーネント。
 * パンくずリストによるナビゲーションを提供し、Markdownファイルを選択できる。
 *
 * 責任:
 * - ディレクトリ構造の表示
 * - ディレクトリナビゲーション
 * - ファイルの選択（Markdownファイルのみ）
 */
export function FileBrowser({
  connectionId,
  repositoryId,
  repositoryFullName,
}: FileBrowserProps) {
  const [searchParams, setSearchParams] = useSearchParams();
  const prefix = searchParams.get("prefix") || "";
  const oauthQueryService = useOAuthQueryService();
  const query = useQuery({
    queryKey: ["repositoryContents", connectionId, repositoryFullName, prefix],
    queryFn: async () =>
      oauthQueryService.listRepositoryContents(
        connectionId,
        repositoryId,
        repositoryFullName,
        prefix
      ),
  });

  if (query.isLoading) {
    return <LoadingSpinner message="Loading..." />;
  }

  if (query.isError || !query.data) {
    return (
      <div className="text-sm text-red-500">
        Failed to load files. Please try again.
      </div>
    );
  }

  const fileNodes = query.data;

  const handleNodeClick = (node: Content) => {
    if (node.type === "directory") {
      setSearchParams((searchParams) => {
        searchParams.set("prefix", node.path);
        return searchParams;
      });
    } else if (node.name.endsWith(".md") || node.name.endsWith(".markdown")) {
      setSearchParams((searchParams) => {
        searchParams.set("selected_file", node.path);
        return searchParams;
      });
    }
  };

  const getFileIcon = (node: Content) => {
    if (node.type === "directory") {
      return (
        <svg
          className="w-5 h-5 text-blue-500"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z"
          />
        </svg>
      );
    }

    const isMarkdown =
      node.name.endsWith(".md") || node.name.endsWith(".markdown");
    return (
      <svg
        className={`w-5 h-5 ${isMarkdown ? "text-green-500" : "text-gray-400"}`}
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
    );
  };

  const isClickable = (node: Content) => {
    return (
      node.type === "directory" ||
      node.name.endsWith(".md") ||
      node.name.endsWith(".markdown")
    );
  };

  return (
    <div className="h-full flex flex-col">
      {/* Header with breadcrumbs */}
      <header className="flex-none bg-white dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700 px-6 py-4">
        <h2 className="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-2">
          Browse Files
        </h2>
        {/* <nav className="flex items-center space-x-2 text-sm">
          {breadcrumbs.map((crumb, index) => (
            <div key={crumb.path} className="flex items-center">
              {index > 0 && (
                <svg
                  className="w-4 h-4 text-gray-400 mx-2"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M9 5l7 7-7 7"
                  />
                </svg>
              )}
              <button
                onClick={() => setCurrentPath(crumb.path)}
                className={`${
                  index === breadcrumbs.length - 1
                    ? "text-gray-900 dark:text-gray-100 font-medium"
                    : "text-blue-600 hover:text-blue-700"
                }`}
              >
                {crumb.name}
              </button>
            </div>
          ))}
        </nav> */}
      </header>

      {/* File list */}
      <div className="flex-1 overflow-y-auto p-6">
        <div className="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 overflow-hidden">
          {fileNodes.length === 0 ? (
            <div className="text-center text-gray-500 dark:text-gray-400 py-8">
              <p>This directory is empty</p>
            </div>
          ) : (
            <div className="divide-y divide-gray-200 dark:divide-gray-700">
              {fileNodes.map((node) => (
                <button
                  key={node.path}
                  onClick={() => handleNodeClick(node)}
                  disabled={!isClickable(node)}
                  className={`w-full flex items-center gap-3 px-4 py-3 text-left transition-colors ${
                    isClickable(node)
                      ? "hover:bg-gray-50 dark:hover:bg-gray-700 cursor-pointer"
                      : "cursor-not-allowed opacity-50"
                  }`}
                >
                  <div className="shrink-0">{getFileIcon(node)}</div>
                  <div className="flex-1 min-w-0">
                    <div className="text-sm font-medium text-gray-900 dark:text-gray-100 truncate">
                      {node.name}
                    </div>
                    {node.type === "file" && node.size && (
                      <div className="text-xs text-gray-500 dark:text-gray-400">
                        {formatFileSize(node.size)}
                      </div>
                    )}

                    {isClickable(node) && node.type === "file" && (
                      <Link
                        to={`/documents/preview?connection_id=${connectionId}&repository_id=${repositoryId}&repository_full_name=${encodeURIComponent(
                          repositoryFullName
                        )}&path=${encodeURIComponent(node.path)}`}
                        className="text-blue-500 hover:text-blue-700 font-medium"
                      >
                        Preview
                      </Link>
                    )}
                  </div>

                  {node.type === "directory" && (
                    <svg
                      className="w-4 h-4 text-gray-400 shrink-0"
                      fill="none"
                      viewBox="0 0 24 24"
                      stroke="currentColor"
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth={2}
                        d="M9 5l7 7-7 7"
                      />
                    </svg>
                  )}
                </button>
              ))}
            </div>
          )}
        </div>

        <div className="mt-4 text-xs text-gray-500 dark:text-gray-400">
          <p>💡 Tip: Only Markdown files (.md, .markdown) can be published</p>
        </div>
      </div>
    </div>
  );
}
