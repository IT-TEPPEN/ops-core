import { useQuery } from "@tanstack/react-query";
import { useSearchParams } from "react-router-dom";
import { useOAuthQueryService } from "../contexts";

/**
 * コネクション一覧コンポーネント
 * 登録済みのOAuthコネクションを一覧表示し、選択できるようにする
 */
export function ConnectionList() {
  const [searchParams, setSearchParams] = useSearchParams();
  const oauthQueryService = useOAuthQueryService();
  const query = useQuery({
    queryKey: ["connections"],
    queryFn: async () => oauthQueryService.listConnections(),
  });

  const selected_connection_id = searchParams.get("selected_connection_id");

  if (query.isLoading) {
    return <div className="text-sm text-gray-500">Loading connections...</div>;
  }

  if (query.isError || !query.data) {
    return (
      <div className="text-sm text-red-500">
        Failed to load connections. Please try again.
      </div>
    );
  }

  const connections = query.data;

  if (connections.length === 0) {
    return (
      <div className="text-sm text-gray-500">
        No OAuth connections found. Click &quot;Add&quot; to connect to a Git
        provider.
      </div>
    );
  }

  return (
    <div className="space-y-2">
      {connections.map((connection, i) => {
        const isSelected =
          connection.id === selected_connection_id ||
          (selected_connection_id === "first" && i === 0);
        return (
          <button
            key={connection.id}
            onClick={(e) => {
              e.preventDefault();
              setSearchParams((searchParams) => {
                searchParams.set("selected_connection_id", connection.id);
                return searchParams;
              });
            }}
            className={`w-full text-left p-4 rounded-lg border transition-all ${
              isSelected
                ? "border-blue-500 bg-blue-50 dark:bg-gray-800 shadow-md"
                : "border-gray-200 hover:border-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700"
            }`}
          >
            <div className="flex items-start space-x-3">
              <span className="text-2xl" aria-label={connection.provider}>
                {getProviderIcon(connection.provider)}
              </span>
              <div className="flex-1 min-w-0">
                <div className="flex items-center justify-between">
                  <h3 className="font-semibold text-base truncate">
                    {connection.providerHost}
                  </h3>
                  {isSelected && (
                    <span className="text-xs text-blue-600 font-medium">
                      Selected
                    </span>
                  )}
                </div>
                <p className="text-sm text-gray-600 dark:text-gray-300 mt-1 truncate">
                  @{connection.providerUsername}
                </p>
                <p className="text-xs text-gray-600 dark:text-gray-300">
                  Connected:{" "}
                  {new Date(connection.connectedAt).toLocaleDateString()}
                </p>
              </div>
            </div>
          </button>
        );
      })}
    </div>
  );
}

const getProviderIcon = (provider: string) => {
  switch (provider) {
    case "github":
      return "🐙"; // GitHub Octocat placeholder
    case "gitlab":
    case "gitlab-self-hosted":
      return "🦊"; // GitLab Fox placeholder
    default:
      return "🔗";
  }
};
