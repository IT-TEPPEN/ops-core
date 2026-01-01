import { useQuery } from "@tanstack/react-query";
import { useSearchParams } from "react-router-dom";
import { useOAuthQueryService } from "../contexts";

/**
 * コネクション一覧コンポーネント
 * 登録済みのOAuthコネクションを一覧表示し、選択できるようにする
 */
export function ConnectionList() {
  const [searchParams, setSearchParams] = useSearchParams();
  const selected_connection_id = searchParams.get("selected_connection_id");

  const oauthQueryService = useOAuthQueryService();
  const query = useQuery({
    queryKey: ["connections"],
    queryFn: async () => oauthQueryService.listConnections(),
  });

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
      return (
        <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
          <path d="M12 0C5.37 0 0 5.37 0 12c0 5.31 3.435 9.795 8.205 11.385.6.105.825-.255.825-.57 0-.285-.015-1.23-.015-2.235-3.015.555-3.795-.735-4.035-1.41-.135-.345-.72-1.41-1.23-1.695-.42-.225-1.02-.78-.015-.795.945-.015 1.62.87 1.845 1.23 1.08 1.815 2.805 1.305 3.495.99.105-.78.42-1.305.765-1.605-2.67-.3-5.46-1.335-5.46-5.925 0-1.305.465-2.385 1.23-3.225-.12-.3-.54-1.53.12-3.18 0 0 1.005-.315 3.3 1.23.96-.27 1.98-.405 3-.405s2.04.135 3 .405c2.295-1.56 3.3-1.23 3.3-1.23.66 1.65.24 2.88.12 3.18.765.84 1.23 1.905 1.23 3.225 0 4.605-2.805 5.625-5.475 5.925.435.375.81 1.095.81 2.22 0 1.605-.015 2.895-.015 3.3 0 .315.225.69.825.57A12.02 12.02 0 0024 12c0-6.63-5.37-12-12-12z" />
        </svg>
      );
    case "gitlab":
      return (
        <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
          <path d="M23.955 13.587l-1.342-4.135-2.664-8.189c-.135-.423-.73-.423-.867 0L16.418 9.45H7.582L4.919 1.263C4.783.84 4.185.84 4.05 1.26L1.386 9.449.044 13.587c-.121.375.014.789.331 1.023L12 23.054l11.625-8.443c.318-.235.453-.647.33-1.024" />
        </svg>
      );
    default:
      return null;
  }
};
