import { useQuery } from "@tanstack/react-query";
import { useRepositoryQueryService } from "../contexts/RepositoryQueryServiceContext";
import type { FileCommitViewData } from "../../application/dto";

interface FileCommitHistoryProps {
  repoId: string;
  filePath: string;
  currentCommit?: string;
  onCommitSelect: (commitHash: string) => void;
}

export function FileCommitHistory({
  repoId,
  filePath,
  currentCommit,
  onCommitSelect,
}: FileCommitHistoryProps) {
  const queryService = useRepositoryQueryService();

  const { data: commits, isLoading, error } = useQuery({
    queryKey: ["repositories", repoId, "files", "history", filePath],
    queryFn: () => queryService.getFileHistory(repoId, filePath),
  });

  if (isLoading) {
    return (
      <div className="p-4">
        <p className="text-sm text-gray-500 dark:text-gray-400">
          Loading commit history...
        </p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="p-4">
        <p className="text-sm text-red-600 dark:text-red-400">
          Failed to load commit history
        </p>
      </div>
    );
  }

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
        {commits.map((commit: FileCommitViewData) => (
          <button
            key={commit.commitHash}
            onClick={() => onCommitSelect(commit.commitHash)}
            className={`w-full text-left p-3 rounded border transition-colors ${
              currentCommit === commit.commitHash
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
              {currentCommit === commit.commitHash && (
                <span className="ml-2 px-2 py-1 text-xs font-medium bg-blue-100 dark:bg-blue-900/40 text-blue-800 dark:text-blue-300 rounded">
                  Current
                </span>
              )}
            </div>
            <p className="text-xs font-mono text-gray-500 dark:text-gray-500 mt-1 truncate">
              {commit.commitHash.substring(0, 7)}
            </p>
          </button>
        ))}
      </div>
    </div>
  );
}
