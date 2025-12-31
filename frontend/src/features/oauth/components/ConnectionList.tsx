import { Connection } from "../application/dto";

interface ConnectionListProps {
  connections: Connection[];
  selectedConnectionId: string | null;
  onSelectConnection: (connection: Connection) => void;
  onAddConnection: () => void;
  isLoading: boolean;
}

/**
 * コネクション一覧コンポーネント
 * 登録済みのOAuthコネクションを一覧表示し、選択できるようにする
 */
export function ConnectionList({
  connections,
  selectedConnectionId,
  onSelectConnection,
  onAddConnection,
  isLoading,
}: ConnectionListProps) {
  if (isLoading) {
    return (
      <div className="space-y-4">
        <div className="flex items-center justify-between">
          <h2 className="text-lg font-semibold">Connections</h2>
          <button
            onClick={onAddConnection}
            className="px-3 py-1.5 text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 rounded-lg transition-colors"
          >
            + Add
          </button>
        </div>
        <div className="text-sm text-gray-500">Loading connections...</div>
      </div>
    );
  }

  if (connections.length === 0) {
    return (
      <div className="space-y-4">
        <div className="flex items-center justify-between">
          <h2 className="text-lg font-semibold">Connections</h2>
          <button
            onClick={onAddConnection}
            className="px-3 py-1.5 text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 rounded-lg transition-colors"
          >
            + Add
          </button>
        </div>
        <div className="text-sm text-gray-500">
          No OAuth connections found. Click &quot;Add&quot; to connect to a Git
          provider.
        </div>
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

  const getProviderLabel = (provider: string) => {
    switch (provider) {
      case "github":
        return "GitHub";
      case "gitlab":
        return "GitLab";
      case "gitlab-self-hosted":
        return "GitLab (Self-Hosted)";
      default:
        return provider;
    }
  };

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h2 className="text-lg font-semibold">Connections</h2>
        <button
          onClick={onAddConnection}
          className="px-3 py-1.5 text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 rounded-lg transition-colors"
        >
          + Add
        </button>
      </div>
      <div className="space-y-2">
        {connections.map((connection) => {
          const isSelected = connection.id === selectedConnectionId;
          return (
            <button
              key={connection.id}
              onClick={() => onSelectConnection(connection)}
              className={`w-full text-left p-4 rounded-lg border transition-all ${
                isSelected
                  ? "border-blue-500 bg-blue-50 dark:bg-gray-800 shadow-md"
                  : "border-gray-200 hover:border-gray-300 hover:bg-gray-50"
              }`}
            >
              <div className="flex items-start space-x-3">
                <span className="text-2xl" aria-label={connection.provider}>
                  {getProviderIcon(connection.provider)}
                </span>
                <div className="flex-1 min-w-0">
                  <div className="flex items-center justify-between">
                    <h3 className="font-semibold text-sm truncate">
                      {getProviderLabel(connection.provider)}
                    </h3>
                    {isSelected && (
                      <span className="text-xs text-blue-600 font-medium">
                        Selected
                      </span>
                    )}
                  </div>
                  <p className="text-sm text-gray-600 truncate">
                    @{connection.providerUsername}
                  </p>
                  <p className="text-xs text-gray-400 mt-1">
                    Connected:{" "}
                    {new Date(connection.connectedAt).toLocaleDateString()}
                  </p>
                </div>
              </div>
            </button>
          );
        })}
      </div>
    </div>
  );
}
