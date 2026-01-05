import { useSuspenseQuery } from "@tanstack/react-query";
import { useOAuthQueryService } from "@/features/oauth";
import { useSearchParams } from "react-router-dom";
import { useEffect } from "react";

interface FileCommitHistoryProps {
  connectionId: string;
  repositoryId: string;
  repositoryFullName: string;
  filePath: string;
}

export function FileCommitHistory({
  connectionId,
  repositoryId,
  repositoryFullName,
  filePath,
}: FileCommitHistoryProps) {
  const [searchParams, setSearchParams] = useSearchParams();
  const currentCommit = searchParams.get("commit") || undefined;
  const queryService = useOAuthQueryService();

  const { data: commits } = useSuspenseQuery({
    queryKey: ["repositories", repositoryId, "files", "history", filePath],
    queryFn: () =>
      queryService.getFileCommitHistory(
        connectionId,
        repositoryId,
        repositoryFullName,
        filePath
      ),
  });

  useEffect(() => {
    if (
      commits &&
      commits.length > 0 &&
      (!currentCommit || currentCommit === "")
    ) {
      setSearchParams((searchParams) => {
        searchParams.set("commit", commits[0].hash);
        return searchParams;
      });
    }
  }, [commits, currentCommit, setSearchParams]);

  if (!commits || commits.length === 0) {
    return (
      <div className="p-4">
        <p className="text-sm text-gray-500 dark:text-gray-400">
          No commit history available
        </p>
      </div>
    );
  }

  return (
    <div className="space-y-2">
      <h3 className="text-lg font-semibold mb-3">Commit History</h3>
      <div className="space-y-2 max-h-96 overflow-y-auto">
        {commits.map((commit) => (
          <button
            key={commit.hash}
            onClick={(e) => {
              e.preventDefault();
              setSearchParams((searchParams) => {
                searchParams.set("commit", commit.hash);
                return searchParams;
              });
            }}
            className={`w-full text-left p-3 rounded border transition-colors ${
              currentCommit === commit.hash
                ? "bg-blue-50 dark:bg-blue-900/20 border-blue-500 dark:border-blue-400"
                : "bg-white dark:bg-gray-800 border-gray-200 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-700"
            }`}
          >
            <div className="flex items-start justify-between">
              <div className="flex-1 min-w-0">
                <p className="text-sm font-medium text-gray-900 dark:text-gray-100 truncate">
                  {commit.message}
                </p>
                <p className="text-xs text-gray-600 dark:text-gray-400 mt-1">
                  {commit.author} • {commit.date.toLocaleDateString()}
                </p>
              </div>
              {currentCommit === commit.hash && (
                <span className="ml-2 px-2 py-1 text-xs font-medium bg-blue-100 dark:bg-blue-900/40 text-blue-800 dark:text-blue-300 rounded">
                  Current
                </span>
              )}
            </div>
            <p className="text-xs font-mono text-gray-500 dark:text-gray-500 mt-1 truncate">
              {commit.hash.substring(0, 7)}
            </p>
          </button>
        ))}
      </div>
    </div>
  );
}
