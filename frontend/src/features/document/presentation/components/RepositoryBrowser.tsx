import { useState, useMemo } from "react";

interface Repository {
  id: string;
  owner: string;
  name: string;
  fullName: string;
  description: string | null;
  private: boolean;
  defaultBranch: string;
}

interface RepositoryBrowserProps {
  connectionId: string;
  selectedRepository: {
    owner: string;
    name: string;
    fullName: string;
  } | null;
  onRepositorySelect: (repository: {
    owner: string;
    name: string;
    fullName: string;
  }) => void;
}

// Mock data - will be replaced with actual API call
const getMockRepositories = (connectionId: string): Repository[] => {
  if (connectionId === "conn-1") {
    // GitHub repositories
    return [
      {
        id: "repo-1",
        owner: "john-doe",
        name: "ops-procedures",
        fullName: "john-doe/ops-procedures",
        description: "Operational procedures and runbooks",
        private: false,
        defaultBranch: "main",
      },
      {
        id: "repo-2",
        owner: "john-doe",
        name: "incident-response",
        fullName: "john-doe/incident-response",
        description: "Incident response documentation",
        private: true,
        defaultBranch: "main",
      },
      {
        id: "repo-3",
        owner: "john-doe",
        name: "knowledge-base",
        fullName: "john-doe/knowledge-base",
        description: "Team knowledge base",
        private: false,
        defaultBranch: "main",
      },
    ];
  } else if (connectionId === "conn-2") {
    // GitLab repositories
    return [
      {
        id: "repo-4",
        owner: "jane-smith",
        name: "deploy-scripts",
        fullName: "jane-smith/deploy-scripts",
        description: "Deployment automation scripts",
        private: true,
        defaultBranch: "master",
      },
      {
        id: "repo-5",
        owner: "jane-smith",
        name: "monitoring-guides",
        fullName: "jane-smith/monitoring-guides",
        description: null,
        private: false,
        defaultBranch: "main",
      },
    ];
  }
  return [];
};

/**
 * RepositoryBrowser Component
 *
 * 選択されたOAuth接続に紐づくリポジトリ一覧を表示し、選択できるようにする。
 *
 * 責任:
 * - リポジトリ一覧の表示
 * - リポジトリの検索
 * - リポジトリの選択
 */
export function RepositoryBrowser({
  connectionId,
  selectedRepository,
  onRepositorySelect,
}: RepositoryBrowserProps) {
  const [searchQuery, setSearchQuery] = useState("");
  const repositories = useMemo(
    () => getMockRepositories(connectionId),
    [connectionId]
  );

  const filteredRepositories = useMemo(() => {
    if (!searchQuery) return repositories;
    const query = searchQuery.toLowerCase();
    return repositories.filter(
      (repo) =>
        repo.name.toLowerCase().includes(query) ||
        repo.fullName.toLowerCase().includes(query) ||
        repo.description?.toLowerCase().includes(query)
    );
  }, [repositories, searchQuery]);

  const isSelected = (repo: Repository) => {
    return selectedRepository?.fullName === repo.fullName;
  };

  return (
    <div className="space-y-3">
      <h3 className="text-sm font-medium text-gray-700 dark:text-gray-300">
        Repository
      </h3>

      {/* Search Box */}
      <div className="relative">
        <input
          type="text"
          placeholder="Search repositories..."
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          className="w-full px-3 py-2 pl-9 text-sm border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 placeholder-gray-500 dark:placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
        <svg
          className="absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-gray-400"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
          />
        </svg>
      </div>

      {/* Repository List */}
      <div className="space-y-2 max-h-96 overflow-y-auto">
        {filteredRepositories.length === 0 ? (
          <div className="text-xs text-gray-500 dark:text-gray-400 text-center py-4">
            {searchQuery
              ? "No repositories found matching your search"
              : "No repositories available"}
          </div>
        ) : (
          filteredRepositories.map((repo) => (
            <button
              key={repo.id}
              onClick={() =>
                onRepositorySelect({
                  owner: repo.owner,
                  name: repo.name,
                  fullName: repo.fullName,
                })
              }
              className={`w-full text-left px-3 py-2 rounded-lg border transition-colors ${
                isSelected(repo)
                  ? "border-blue-500 bg-blue-50 dark:bg-blue-900/30"
                  : "border-gray-200 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-700"
              }`}
            >
              <div className="flex items-start gap-2">
                <div className="flex-1 min-w-0">
                  <div className="flex items-center gap-2">
                    <span className="text-sm font-medium text-gray-900 dark:text-gray-100 truncate">
                      {repo.name}
                    </span>
                    {repo.private && (
                      <span className="flex-shrink-0 px-1.5 py-0.5 text-xs bg-gray-200 dark:bg-gray-600 text-gray-700 dark:text-gray-300 rounded">
                        Private
                      </span>
                    )}
                  </div>
                  {repo.description && (
                    <p className="text-xs text-gray-600 dark:text-gray-400 truncate mt-0.5">
                      {repo.description}
                    </p>
                  )}
                  <p className="text-xs text-gray-500 dark:text-gray-500 mt-0.5">
                    Branch: {repo.defaultBranch}
                  </p>
                </div>
                {isSelected(repo) && (
                  <svg
                    className="w-4 h-4 text-blue-600 flex-shrink-0 mt-0.5"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M5 13l4 4L19 7"
                    />
                  </svg>
                )}
              </div>
            </button>
          ))
        )}
      </div>

      <div className="text-xs text-gray-500 dark:text-gray-400">
        {filteredRepositories.length} of {repositories.length} repositories
      </div>
    </div>
  );
}
