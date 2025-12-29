import { useState } from "react";
import { GitRepository } from "@/shared/api/gitProviderApi";

interface RepositorySelectorProps {
  repositories: GitRepository[];
  isLoading: boolean;
  onSelect: (repository: GitRepository) => void;
  selectedRepository?: GitRepository | null;
}

/**
 * GitHubなどから取得したリポジトリ一覧から選択するコンポーネント
 */
export function RepositorySelector({
  repositories,
  isLoading,
  onSelect,
  selectedRepository,
}: RepositorySelectorProps) {
  const [searchQuery, setSearchQuery] = useState("");

  const filteredRepositories = repositories.filter(
    (repo) =>
      repo.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      repo.fullName.toLowerCase().includes(searchQuery.toLowerCase()) ||
      (repo.description &&
        repo.description.toLowerCase().includes(searchQuery.toLowerCase()))
  );

  if (isLoading) {
    return (
      <div className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow">
        <h3 className="text-lg font-semibold mb-4">Select Repository</h3>
        <div className="flex items-center justify-center py-8">
          <svg
            className="animate-spin h-8 w-8 text-blue-600"
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
          >
            <circle
              className="opacity-25"
              cx="12"
              cy="12"
              r="10"
              stroke="currentColor"
              strokeWidth="4"
            ></circle>
            <path
              className="opacity-75"
              fill="currentColor"
              d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
            ></path>
          </svg>
          <span className="ml-3 text-gray-600 dark:text-gray-400">
            Loading repositories...
          </span>
        </div>
      </div>
    );
  }

  return (
    <div className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow">
      <h3 className="text-lg font-semibold mb-4">Select Repository</h3>

      {/* 検索ボックス */}
      <div className="mb-4">
        <input
          type="text"
          placeholder="Search repositories..."
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          className="w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
      </div>

      {/* リポジトリ一覧 */}
      <div className="max-h-96 overflow-y-auto space-y-2">
        {filteredRepositories.length === 0 ? (
          <p className="text-center text-gray-500 dark:text-gray-400 py-4">
            {searchQuery
              ? "No repositories found matching your search"
              : "No repositories available"}
          </p>
        ) : (
          filteredRepositories.map((repo) => (
            <button
              key={repo.id}
              type="button"
              onClick={() => onSelect(repo)}
              className={`w-full text-left p-4 rounded-lg border transition-colors ${
                selectedRepository?.id === repo.id
                  ? "border-blue-500 bg-blue-50 dark:bg-blue-900/30"
                  : "border-gray-200 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-700"
              }`}
            >
              <div className="flex items-start gap-3">
                <img
                  src={repo.owner.avatarUrl}
                  alt={repo.owner.login}
                  className="w-10 h-10 rounded-full"
                />
                <div className="flex-1 min-w-0">
                  <div className="flex items-center gap-2">
                    <span className="font-medium text-gray-900 dark:text-gray-100 truncate">
                      {repo.fullName}
                    </span>
                    {repo.private && (
                      <span className="px-2 py-0.5 text-xs bg-gray-200 dark:bg-gray-600 text-gray-700 dark:text-gray-300 rounded">
                        Private
                      </span>
                    )}
                  </div>
                  {repo.description && (
                    <p className="text-sm text-gray-600 dark:text-gray-400 truncate mt-1">
                      {repo.description}
                    </p>
                  )}
                  <p className="text-xs text-gray-500 dark:text-gray-500 mt-1">
                    Default branch: {repo.defaultBranch}
                  </p>
                </div>
                {selectedRepository?.id === repo.id && (
                  <svg
                    className="w-5 h-5 text-blue-600 shrink-0"
                    fill="currentColor"
                    viewBox="0 0 20 20"
                  >
                    <path
                      fillRule="evenodd"
                      d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
                      clipRule="evenodd"
                    />
                  </svg>
                )}
              </div>
            </button>
          ))
        )}
      </div>

      <p className="text-sm text-gray-500 dark:text-gray-400 mt-4">
        {filteredRepositories.length} of {repositories.length} repositories
      </p>
    </div>
  );
}
